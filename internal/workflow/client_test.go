package workflow_test

import (
	"context"
	"strings"
	"testing"

	"askoc-mvp/internal/workflow"
)

func TestPaymentReminderKeyIsStableAndScoped(t *testing.T) {
	got := workflow.PaymentReminderKey("trace-p4", "S100002", "official_transcript")
	want := "payment-reminder:trace-p4:S100002:official_transcript"
	if got != want {
		t.Fatalf("PaymentReminderKey = %q, want %q", got, want)
	}
}

func TestInMemoryClientReturnsSameWorkflowForDuplicateIdempotencyKey(t *testing.T) {
	client := workflow.NewInMemoryClient()
	req := workflow.PaymentReminderRequest{
		StudentID:      "S100002",
		ConversationID: "conv-p4",
		TraceID:        "trace-p4",
		Item:           "official_transcript",
		AmountDue:      15,
		Currency:       "CAD",
		Reason:         "Transcript request cannot be processed until payment is complete.",
		IdempotencyKey: workflow.PaymentReminderKey("trace-p4", "S100002", "official_transcript"),
	}

	first, err := client.SendPaymentReminder(context.Background(), req)
	if err != nil {
		t.Fatalf("first SendPaymentReminder returned error: %v", err)
	}
	second, err := client.SendPaymentReminder(context.Background(), req)
	if err != nil {
		t.Fatalf("second SendPaymentReminder returned error: %v", err)
	}

	if first.WorkflowID == "" || !strings.HasPrefix(first.WorkflowID, "LOCAL-WF-") {
		t.Fatalf("workflow ID = %q, want local synthetic workflow ID", first.WorkflowID)
	}
	if second.WorkflowID != first.WorkflowID {
		t.Fatalf("duplicate workflow ID = %q, want %q", second.WorkflowID, first.WorkflowID)
	}
	if second.IdempotencyKey != req.IdempotencyKey {
		t.Fatalf("idempotency key = %q, want %q", second.IdempotencyKey, req.IdempotencyKey)
	}
}

func TestInMemoryClientRejectsMissingIdempotencyKey(t *testing.T) {
	client := workflow.NewInMemoryClient()
	_, err := client.SendPaymentReminder(context.Background(), workflow.PaymentReminderRequest{
		StudentID: "S100002",
		Item:      "official_transcript",
	})
	if err == nil {
		t.Fatal("SendPaymentReminder returned nil error, want validation error")
	}
	if !strings.Contains(err.Error(), "idempotency") {
		t.Fatalf("error = %q, want idempotency guidance", err.Error())
	}
}
