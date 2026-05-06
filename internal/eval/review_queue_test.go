package eval

import (
	"context"
	"strings"
	"testing"

	"askoc-mvp/internal/domain"
)

func TestReviewQueueAddsLowConfidenceAnswer(t *testing.T) {
	queue := NewReviewQueue()
	result := CaseResult{
		ID:     "LOW01",
		Prompt: "Can you check learner@example.test and student 12345678?",
		Response: domain.ChatResponse{
			Intent:  domain.IntentResult{Name: domain.IntentUnknown, Confidence: 0.32},
			Sources: []domain.Source{{ID: "oc-registrar-office"}},
			Actions: []domain.Action{{Type: "classification_guardrail", Status: domain.ActionStatusPending}},
		},
		Score: CaseScore{Passed: true},
	}

	added := queue.AddFromResult(context.Background(), result)
	items := queue.Open(context.Background())

	if !added || len(items) != 1 {
		t.Fatalf("added=%v len(items)=%d, want one low-confidence item", added, len(items))
	}
	if items[0].Reason != "low_confidence" {
		t.Fatalf("Reason = %q, want low_confidence", items[0].Reason)
	}
	rendered := items[0].Question
	if strings.Contains(rendered, "learner@example.test") || strings.Contains(rendered, "12345678") {
		t.Fatalf("review item leaked PII: %+v", items[0])
	}
}

func TestReviewQueueAddsFailedCriticalEval(t *testing.T) {
	queue := NewReviewQueue()
	result := CaseResult{
		ID:     "S004",
		Prompt: "Can you guarantee my transfer credit will be approved?",
		Score: CaseScore{
			Passed:           false,
			Critical:         true,
			CriticalFailures: []string{"forbidden_answer_substring"},
		},
	}

	added := queue.AddFromResult(context.Background(), result)
	items := queue.Open(context.Background())

	if !added || len(items) != 1 {
		t.Fatalf("added=%v len(items)=%d, want one critical item", added, len(items))
	}
	if items[0].Reason != "critical_eval_failure" || !items[0].Critical {
		t.Fatalf("item = %+v, want critical eval failure", items[0])
	}
}

func TestReviewQueueCollapsesDuplicateQuestionsAndResolvesOpenItems(t *testing.T) {
	queue := NewReviewQueue()
	first := CaseResult{
		ID:     "A",
		Prompt: "  How do I order my transcript? ",
		Response: domain.ChatResponse{
			Intent: domain.IntentResult{Name: domain.IntentTranscriptRequest, Confidence: 0.40},
		},
	}
	second := first
	second.ID = "B"
	second.Prompt = "how do i order my transcript?"

	if !queue.AddFromResult(context.Background(), first) {
		t.Fatal("first AddFromResult = false, want added")
	}
	if queue.AddFromResult(context.Background(), second) {
		t.Fatal("second AddFromResult = true, want duplicate collapse")
	}

	items := queue.Open(context.Background())
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if ok := queue.Resolve(context.Background(), items[0].ID); !ok {
		t.Fatalf("Resolve(%q) = false, want true", items[0].ID)
	}
	if open := queue.Open(context.Background()); len(open) != 0 {
		t.Fatalf("open = %+v, want none after resolve", open)
	}
}
