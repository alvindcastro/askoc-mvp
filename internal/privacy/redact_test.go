package privacy

import (
	"strings"
	"testing"
)

func TestRedactRemovesConfiguredPIIPatterns(t *testing.T) {
	input := "Email learner@example.test or call 250-555-0199. My password is OceanBlue42 and student number 12345678 needs help."

	got := Redact(input)

	for _, leaked := range []string{"learner@example.test", "250-555-0199", "OceanBlue42", "12345678"} {
		if strings.Contains(got, leaked) {
			t.Fatalf("Redact leaked %q in %q", leaked, got)
		}
	}
	for _, marker := range []string{"[REDACTED_EMAIL]", "[REDACTED_PHONE]", "[REDACTED_SECRET]", "[REDACTED_ID]"} {
		if !strings.Contains(got, marker) {
			t.Fatalf("Redact(%q) missing marker %q: %q", input, marker, got)
		}
	}
}

func TestRedactPreservesSyntheticStudentIDs(t *testing.T) {
	input := "Synthetic student S100001 and request SYNTH-TRN-100001 are safe demo identifiers."

	got := Redact(input)

	if !strings.Contains(got, "S100001") {
		t.Fatalf("Redact removed synthetic student ID: %q", got)
	}
	if strings.Contains(got, "[REDACTED_ID]") {
		t.Fatalf("Redact treated demo identifiers as real IDs: %q", got)
	}
}

func TestRedactAvoidsCommonFalsePositives(t *testing.T) {
	input := "Course MATH 120, room A120, and term 2026S are not learner PII."

	got := Redact(input)

	if got != input {
		t.Fatalf("Redact changed non-PII text: got %q want %q", got, input)
	}
}

func TestRedactHandlesSecretLikeAssignments(t *testing.T) {
	input := "token=abc123 api_key: sk-demo-123 and passcode is 778899"

	got := Redact(input)

	for _, leaked := range []string{"abc123", "sk-demo-123", "778899"} {
		if strings.Contains(got, leaked) {
			t.Fatalf("Redact leaked secret-like value %q in %q", leaked, got)
		}
	}
}
