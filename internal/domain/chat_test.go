package domain_test

import (
	"encoding/json"
	"strings"
	"testing"

	"askoc-mvp/internal/domain"
	"askoc-mvp/internal/validation"
)

func TestChatRequestJSONRoundTrip(t *testing.T) {
	in := domain.ChatRequest{
		ConversationID: "conv_existing",
		Channel:        "web",
		Message:        "I ordered my transcript but it has not been processed.",
		StudentID:      "S100002",
	}

	body, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal chat request: %v", err)
	}

	var got domain.ChatRequest
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal chat request: %v", err)
	}

	if got != in {
		t.Fatalf("chat request = %+v, want %+v", got, in)
	}
	if !strings.Contains(string(body), `"student_id":"S100002"`) {
		t.Fatalf("student ID JSON field missing: %s", body)
	}
}

func TestChatResponseJSONIncludesSourcesActionsAndEscalation(t *testing.T) {
	in := domain.ChatResponse{
		ConversationID: "conv_123",
		TraceID:        "trace-123",
		Answer:         "This is a deterministic demo response.",
		Intent: domain.IntentResult{
			Name:       domain.IntentTranscriptStatus,
			Confidence: 0.77,
		},
		Sentiment: domain.SentimentNeutral,
		Sources: []domain.Source{
			{
				Title:   "Transcript Request - 2005 Onwards",
				URL:     "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
				ChunkID: "transcript-001",
			},
		},
		Actions: []domain.Action{
			{
				Type:   "placeholder_response",
				Status: domain.ActionStatusCompleted,
			},
		},
		Escalation: &domain.Escalation{
			Required: true,
			Status:   domain.HandoffPending,
			Queue:    "demo_review",
			Reason:   "learner requested staff support",
		},
	}

	body, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal chat response: %v", err)
	}

	var got domain.ChatResponse
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal chat response: %v", err)
	}

	if got.Intent.Name != domain.IntentTranscriptStatus || got.Intent.Confidence != 0.77 {
		t.Fatalf("intent = %+v", got.Intent)
	}
	if len(got.Sources) != 1 || got.Sources[0].ChunkID != "transcript-001" {
		t.Fatalf("sources = %+v", got.Sources)
	}
	if len(got.Actions) != 1 || got.Actions[0].Type != "placeholder_response" {
		t.Fatalf("actions = %+v", got.Actions)
	}
	if got.Escalation == nil || got.Escalation.Status != domain.HandoffPending {
		t.Fatalf("escalation = %+v", got.Escalation)
	}
}

func TestValidateChatRequestRejectsMissingMessageAndInvalidStudentID(t *testing.T) {
	tests := []struct {
		name string
		req  domain.ChatRequest
		code string
	}{
		{
			name: "missing message",
			req: domain.ChatRequest{
				Channel: "web",
			},
			code: validation.CodeMissingMessage,
		},
		{
			name: "invalid student ID shape",
			req: domain.ChatRequest{
				Channel:   "web",
				Message:   "Please check my transcript.",
				StudentID: "123456",
			},
			code: validation.CodeInvalidStudentID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateChatRequest(tt.req)
			if err == nil {
				t.Fatal("ValidateChatRequest error = nil, want validation error")
			}
			if code := validation.Code(err); code != tt.code {
				t.Fatalf("validation code = %q, want %q", code, tt.code)
			}
		})
	}
}
