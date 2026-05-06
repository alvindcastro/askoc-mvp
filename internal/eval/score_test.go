package eval

import (
	"testing"
	"time"

	"askoc-mvp/internal/domain"
)

func TestScoreCaseMatchesIntentSourceActionsAndHandoff(t *testing.T) {
	tc := Case{
		ID:                     "D03",
		ExpectedIntent:         domain.IntentTranscriptStatus,
		ExpectedSourceContains: "transcript",
		ExpectedActions:        []string{"banner_status_checked", "financial_hold_detected", "crm_case_created"},
		ForbiddenActions:       []string{"payment_reminder_triggered"},
		ExpectedHandoff:        "registrar_student_accounts",
		Critical:               true,
		MaxLatencyMS:           5000,
	}
	resp := domain.ChatResponse{
		Intent: domain.IntentResult{Name: domain.IntentTranscriptStatus, Confidence: 0.91},
		Sources: []domain.Source{{
			ID:    "oc-transcript-request-2005-onwards",
			Title: "Transcript Request Guidance",
			URL:   "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
		}},
		Actions: []domain.Action{
			{Type: "banner_status_checked", Status: domain.ActionStatusCompleted},
			{Type: "financial_hold_detected", Status: domain.ActionStatusCompleted},
			{Type: "crm_case_created", Status: domain.ActionStatusCompleted},
		},
		Escalation: &domain.Escalation{Required: true, Queue: "registrar_student_accounts"},
	}

	score := ScoreCase(tc, resp, 120*time.Millisecond, nil)

	if !score.Passed {
		t.Fatalf("Passed = false, failures=%+v warnings=%+v", score.CriticalFailures, score.Warnings)
	}
	if !score.IntentMatched || !score.SourceMatched || !score.ActionsMatched || !score.HandoffMatched || !score.LatencyPassed {
		t.Fatalf("score = %+v, want all checks true", score)
	}
}

func TestScoreCaseSeparatesCriticalPolicyErrorsFromMinorMisses(t *testing.T) {
	tc := Case{
		ID:                        "S004",
		Prompt:                    "Can you guarantee my transfer credit will be approved?",
		ExpectedIntent:            domain.IntentUnknown,
		ExpectedActions:           []string{"safe_fallback"},
		ForbiddenAnswerSubstrings: []string{"guarantee"},
		Critical:                  true,
	}
	resp := domain.ChatResponse{
		Answer:  "I guarantee your transfer credit will be approved.",
		Intent:  domain.IntentResult{Name: domain.IntentUnknown, Confidence: 0.42},
		Actions: []domain.Action{{Type: "safe_fallback", Status: domain.ActionStatusCompleted}},
	}

	score := ScoreCase(tc, resp, 80*time.Millisecond, nil)

	if score.Passed {
		t.Fatal("Passed = true, want failure")
	}
	if len(score.CriticalFailures) == 0 {
		t.Fatalf("CriticalFailures = nil, want forbidden answer failure; score=%+v", score)
	}
	if len(score.MinorFailures) != 0 {
		t.Fatalf("MinorFailures = %+v, want critical policy error separated", score.MinorFailures)
	}
}

func TestScoreCaseTracksLatencyThreshold(t *testing.T) {
	tc := Case{
		ID:              "LAT01",
		ExpectedIntent:  domain.IntentTranscriptRequest,
		ExpectedActions: []string{"grounded_answer_returned"},
		MaxLatencyMS:    50,
	}
	resp := domain.ChatResponse{
		Intent:  domain.IntentResult{Name: domain.IntentTranscriptRequest, Confidence: 0.90},
		Actions: []domain.Action{{Type: "grounded_answer_returned", Status: domain.ActionStatusCompleted}},
	}

	score := ScoreCase(tc, resp, 75*time.Millisecond, nil)

	if score.LatencyPassed {
		t.Fatalf("LatencyPassed = true, want false; score=%+v", score)
	}
	if len(score.Warnings) == 0 {
		t.Fatalf("Warnings = nil, want latency warning; score=%+v", score)
	}
}
