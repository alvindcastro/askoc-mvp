package workflow

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"askoc-mvp/internal/audit"
)

func TestSimulatorPaymentReminderPayloadAccepted(t *testing.T) {
	store := audit.NewMemoryStore()
	handler := NewSimulatorHandler(NewInMemoryClient(), store)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, newPaymentReminderRequest(t, validPaymentReminderPayload()))

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusAccepted, rec.Body.String())
	}
	var got PaymentReminderResponse
	decodeResponse(t, rec, &got)
	if got.WorkflowID == "" || got.Status != "accepted" || !got.Synthetic {
		t.Fatalf("response = %+v, want accepted synthetic workflow", got)
	}
	if got.IdempotencyKey != "payment-reminder:trace-sim:S100002:official_transcript" {
		t.Fatalf("idempotency key = %q", got.IdempotencyKey)
	}

	events := store.List(context.Background())
	if len(events) != 1 {
		t.Fatalf("audit event count = %d, want 1", len(events))
	}
	event := events[0]
	if event.Type != audit.EventTypeWorkflow || event.Action != audit.ActionPaymentReminder || event.Status != audit.StatusCompleted {
		t.Fatalf("audit event = %+v, want completed workflow payment reminder", event)
	}
	if event.ReferenceID != got.WorkflowID {
		t.Fatalf("audit reference = %q, want workflow ID %q", event.ReferenceID, got.WorkflowID)
	}
	if event.Metadata["idempotency_key_hash"] == "" {
		t.Fatalf("audit metadata = %+v, want idempotency key hash", event.Metadata)
	}
	if strings.Contains(event.Metadata["idempotency_key_hash"], got.IdempotencyKey) || event.Metadata["idempotency_key"] != "" {
		t.Fatalf("audit metadata = %+v, want no raw idempotency key", event.Metadata)
	}
}

func TestSimulatorPaymentReminderMissingIdempotencyKeyRejected(t *testing.T) {
	handler := NewSimulatorHandler(NewInMemoryClient(), audit.NopRecorder{})
	payload := validPaymentReminderPayload()
	payload.IdempotencyKey = ""

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, newPaymentReminderRequest(t, payload))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "idempotency") {
		t.Fatalf("body = %s, want idempotency guidance", rec.Body.String())
	}
}

func TestSimulatorPaymentReminderDuplicateIdempotencyKeyReturnsSameWorkflowID(t *testing.T) {
	handler := NewSimulatorHandler(NewInMemoryClient(), audit.NopRecorder{})
	payload := validPaymentReminderPayload()

	first := httptest.NewRecorder()
	handler.ServeHTTP(first, newPaymentReminderRequest(t, payload))
	if first.Code != http.StatusAccepted {
		t.Fatalf("first status = %d, want %d; body=%s", first.Code, http.StatusAccepted, first.Body.String())
	}
	var firstResp PaymentReminderResponse
	decodeResponse(t, first, &firstResp)

	second := httptest.NewRecorder()
	handler.ServeHTTP(second, newPaymentReminderRequest(t, payload))
	if second.Code != http.StatusAccepted {
		t.Fatalf("second status = %d, want %d; body=%s", second.Code, http.StatusAccepted, second.Body.String())
	}
	var secondResp PaymentReminderResponse
	decodeResponse(t, second, &secondResp)

	if secondResp.WorkflowID != firstResp.WorkflowID {
		t.Fatalf("duplicate workflow ID = %q, want %q", secondResp.WorkflowID, firstResp.WorkflowID)
	}
}

func TestSimulatorPaymentReminderInvalidPayloadReturnsBadRequest(t *testing.T) {
	handler := NewSimulatorHandler(NewInMemoryClient(), audit.NopRecorder{})
	req := httptest.NewRequest(http.MethodPost, PaymentReminderPath, strings.NewReader(`{"student_id":`))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func validPaymentReminderPayload() PaymentReminderRequest {
	return PaymentReminderRequest{
		StudentID:      "S100002",
		ConversationID: "conv-sim",
		TraceID:        "trace-sim",
		Item:           "official_transcript",
		AmountDue:      15,
		Currency:       "CAD",
		Reason:         "Transcript request cannot be processed until payment is complete.",
		IdempotencyKey: "payment-reminder:trace-sim:S100002:official_transcript",
	}
}

func newPaymentReminderRequest(t *testing.T, payload PaymentReminderRequest) *http.Request {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, PaymentReminderPath, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func decodeResponse(t *testing.T, rec *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.NewDecoder(rec.Body).Decode(target); err != nil {
		t.Fatalf("decode response: %v; body=%s", err, rec.Body.String())
	}
}
