package workflow_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"askoc-mvp/internal/workflow"
)

func TestPowerAutomateClientSendsPaymentReminderJSONAndHeaders(t *testing.T) {
	req := workflow.PaymentReminderRequest{
		StudentID:      "S100002",
		ConversationID: "conv-p8",
		TraceID:        "trace-p8",
		Item:           "official_transcript",
		AmountDue:      15,
		Currency:       "CAD",
		Reason:         "Transcript request cannot be processed until payment is complete.",
		IdempotencyKey: workflow.PaymentReminderKey("trace-p8", "S100002", "official_transcript"),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want application/json", got)
		}
		if got := r.Header.Get("X-Trace-ID"); got != req.TraceID {
			t.Fatalf("X-Trace-ID = %q, want %q", got, req.TraceID)
		}
		if got := r.Header.Get("Idempotency-Key"); got != req.IdempotencyKey {
			t.Fatalf("Idempotency-Key = %q, want %q", got, req.IdempotencyKey)
		}
		if got := r.Header.Get("X-AskOC-Workflow-Signature"); got != "test-signature" {
			t.Fatalf("signature header = %q, want configured value", got)
		}

		var got workflow.PaymentReminderRequest
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if got != req {
			t.Fatalf("request body = %+v, want %+v", got, req)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(workflow.PaymentReminderResponse{
			WorkflowID:     "WF-2026-000789",
			Status:         "accepted",
			Message:        "Payment reminder workflow accepted.",
			IdempotencyKey: req.IdempotencyKey,
		}); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
	defer server.Close()

	client, err := workflow.NewPowerAutomateClient(workflow.PowerAutomateClientConfig{
		WebhookURL:      server.URL,
		HTTPClient:      server.Client(),
		Signature:       "test-signature",
		SignatureHeader: "X-AskOC-Workflow-Signature",
		MaxRetries:      1,
	})
	if err != nil {
		t.Fatalf("NewPowerAutomateClient returned error: %v", err)
	}

	resp, err := client.SendPaymentReminder(context.Background(), req)
	if err != nil {
		t.Fatalf("SendPaymentReminder returned error: %v", err)
	}
	if resp.WorkflowID != "WF-2026-000789" || resp.Status != "accepted" || resp.IdempotencyKey != req.IdempotencyKey {
		t.Fatalf("response = %+v, want accepted workflow response", resp)
	}
}

func TestPowerAutomateClientRetriesTransient5xxWithinLimit(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := atomic.AddInt32(&attempts, 1)
		if attempt == 1 {
			http.Error(w, "temporary outage", http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(workflow.PaymentReminderResponse{
			WorkflowID: "WF-RETRY-OK",
			Status:     "accepted",
		}); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
	defer server.Close()

	client, err := workflow.NewPowerAutomateClient(workflow.PowerAutomateClientConfig{
		WebhookURL: server.URL,
		HTTPClient: server.Client(),
		MaxRetries: 1,
	})
	if err != nil {
		t.Fatalf("NewPowerAutomateClient returned error: %v", err)
	}

	resp, err := client.SendPaymentReminder(context.Background(), validPaymentReminderRequest())
	if err != nil {
		t.Fatalf("SendPaymentReminder returned error: %v", err)
	}
	if resp.WorkflowID != "WF-RETRY-OK" {
		t.Fatalf("WorkflowID = %q, want retry success response", resp.WorkflowID)
	}
	if resp.AttemptCount != 2 {
		t.Fatalf("AttemptCount = %d, want initial request plus one retry", resp.AttemptCount)
	}
	if got := atomic.LoadInt32(&attempts); got != 2 {
		t.Fatalf("attempts = %d, want initial request plus one retry", got)
	}
}

func TestPowerAutomateClientDoesNotRetryPermanent400OrLeakSecrets(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		http.Error(w, "token=body-secret", http.StatusBadRequest)
	}))
	defer server.Close()

	client, err := workflow.NewPowerAutomateClient(workflow.PowerAutomateClientConfig{
		WebhookURL: server.URL + "?sig=query-secret",
		HTTPClient: server.Client(),
		Signature:  "header-secret",
		MaxRetries: 2,
	})
	if err != nil {
		t.Fatalf("NewPowerAutomateClient returned error: %v", err)
	}

	_, err = client.SendPaymentReminder(context.Background(), validPaymentReminderRequest())
	if err == nil {
		t.Fatal("SendPaymentReminder returned nil error, want permanent status error")
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Fatalf("attempts = %d, want no retry for HTTP 400", got)
	}
	for _, secret := range []string{"query-secret", "body-secret", "header-secret"} {
		if strings.Contains(err.Error(), secret) {
			t.Fatalf("error %q leaked secret %q", err.Error(), secret)
		}
	}
}

func TestPowerAutomateClientStopsOnContextCancellation(t *testing.T) {
	var attempts int32
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		cancel()
		http.Error(w, "temporary outage", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client, err := workflow.NewPowerAutomateClient(workflow.PowerAutomateClientConfig{
		WebhookURL: server.URL,
		HTTPClient: server.Client(),
		MaxRetries: 3,
	})
	if err != nil {
		t.Fatalf("NewPowerAutomateClient returned error: %v", err)
	}

	_, err = client.SendPaymentReminder(ctx, validPaymentReminderRequest())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Fatalf("attempts = %d, want cancellation to stop retries", got)
	}
}

func TestNewPowerAutomateClientRejectsMissingWebhookURLSafely(t *testing.T) {
	_, err := workflow.NewPowerAutomateClient(workflow.PowerAutomateClientConfig{
		Signature: "secret-signature",
	})
	if err == nil {
		t.Fatal("NewPowerAutomateClient returned nil error, want configuration error")
	}
	if !strings.Contains(err.Error(), "workflow webhook URL") {
		t.Fatalf("error = %q, want safe webhook URL configuration guidance", err.Error())
	}
	if strings.Contains(err.Error(), "secret-signature") {
		t.Fatalf("error %q leaked signature", err.Error())
	}
}

func validPaymentReminderRequest() workflow.PaymentReminderRequest {
	return workflow.PaymentReminderRequest{
		StudentID:      "S100002",
		ConversationID: "conv-p8",
		TraceID:        "trace-p8",
		Item:           "official_transcript",
		AmountDue:      15,
		Currency:       "CAD",
		Reason:         "Transcript request cannot be processed until payment is complete.",
		IdempotencyKey: workflow.PaymentReminderKey("trace-p8", "S100002", "official_transcript"),
	}
}
