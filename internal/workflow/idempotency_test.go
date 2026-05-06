package workflow_test

import (
	"strings"
	"testing"

	"askoc-mvp/internal/workflow"
)

func TestIdempotencyKeyHashIsDeterministicAndRedactionSafe(t *testing.T) {
	key := workflow.PaymentReminderKey("trace-p8", "S100002", "official_transcript")

	first := workflow.IdempotencyKeyHash(key)
	second := workflow.IdempotencyKeyHash(key)

	if first == "" {
		t.Fatal("IdempotencyKeyHash returned empty string, want deterministic audit hash")
	}
	if second != first {
		t.Fatalf("IdempotencyKeyHash is not deterministic: first %q second %q", first, second)
	}
	for _, unsafe := range []string{"trace-p8", "S100002", "official_transcript", key} {
		if strings.Contains(first, unsafe) {
			t.Fatalf("IdempotencyKeyHash = %q, want redaction-safe hash without %q", first, unsafe)
		}
	}
}

func TestIdempotencyKeyHashReturnsNonEmptyHashForEmptyKey(t *testing.T) {
	if got := workflow.IdempotencyKeyHash(""); got == "" {
		t.Fatal("IdempotencyKeyHash returned empty string for empty key")
	}
}
