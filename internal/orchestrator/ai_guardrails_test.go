package orchestrator

import (
	"context"
	"errors"
	"strings"
	"testing"

	"askoc-mvp/internal/classifier"
	"askoc-mvp/internal/domain"
)

func TestLLMClassificationTimeoutUsesDeterministicFallbackAndAuditsPromptVersion(t *testing.T) {
	fallback := &fakeClassifier{result: classifier.Result{
		Intent:     domain.IntentTranscriptStatus,
		Confidence: 0.86,
		Sentiment:  domain.SentimentNeutral,
	}}
	audit := &fakeAudit{}
	guarded := LLMBackedClassifier{
		LLM:      fakeLLM{err: context.DeadlineExceeded},
		Fallback: fallback,
		Audit:    audit,
	}

	got, err := guarded.Classify(traceContext("trace-llm-timeout"), "I ordered my transcript and it has not arrived.")
	if err != nil {
		t.Fatalf("Classify returned error: %v", err)
	}

	if got.Intent != domain.IntentTranscriptStatus || got.Confidence < classifier.SensitiveToolConfidence {
		t.Fatalf("classification = %+v, want deterministic fallback transcript status", got)
	}
	if len(audit.events) != 1 {
		t.Fatalf("audit events = %+v, want one guardrail event", audit.events)
	}
	event := audit.events[0]
	if event.Action != "llm_classification_fallback" || event.Status != "completed" {
		t.Fatalf("audit event = %+v, want completed fallback event", event)
	}
	if event.Metadata["prompt_version"] != PromptVersion || event.Metadata["reason"] != "model_error" {
		t.Fatalf("audit metadata = %+v, want prompt version and model_error reason", event.Metadata)
	}
}

func TestUnsafeGroundedAnswerWithoutSourcesIsRejected(t *testing.T) {
	_, err := ValidateGroundedAnswer("You can request a transcript from the portal.", nil)
	if err == nil {
		t.Fatal("ValidateGroundedAnswer returned nil error for answer without sources")
	}
	if !errors.Is(err, ErrUngroundedAnswer) {
		t.Fatalf("error = %v, want ErrUngroundedAnswer", err)
	}
}

func TestLowConfidenceValidatedIntentDoesNotTriggerSensitiveTools(t *testing.T) {
	deps := completeDeps()
	classifierFake(t, deps).result = classifier.Result{
		Intent:     domain.IntentTranscriptStatus,
		Confidence: classifier.SensitiveToolConfidence - 0.01,
		Sentiment:  domain.SentimentNeutral,
	}
	orch := mustNew(t, deps)

	resp, err := orch.HandleChat(traceContext("trace-low-validated"), domain.ChatRequest{
		ConversationID: "conv-low-validated",
		Channel:        "web",
		Message:        "I ordered my transcript and it might be delayed.",
		StudentID:      "S100002",
	})
	if err != nil {
		t.Fatalf("HandleChat returned error: %v", err)
	}

	if len(bannerFake(t, deps).calls) != 0 || len(paymentFake(t, deps).calls) != 0 || len(workflowFake(t, deps).calls) != 0 {
		t.Fatalf("sensitive tools should not run: banner=%v payment=%v workflow=%v", bannerFake(t, deps).calls, paymentFake(t, deps).calls, workflowFake(t, deps).calls)
	}
	if len(crmFake(t, deps).requests) != 1 {
		t.Fatalf("crm requests = %+v, want staff handoff", crmFake(t, deps).requests)
	}
	assertAction(t, resp, "classification_guardrail", domain.ActionStatusPending)
}

func TestGroundedAnswerUsesLLMOnlyWhenSourceGuardrailPasses(t *testing.T) {
	deps := completeDeps()
	classifierFake(t, deps).result = classifier.Result{
		Intent:     domain.IntentTranscriptRequest,
		Confidence: 0.92,
		Sentiment:  domain.SentimentNeutral,
	}
	deps.Retriever = &fakeRetriever{sources: []domain.Source{
		{
			Title:      "Transcript Request Guidance",
			URL:        "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
			ChunkID:    "transcript-source",
			Confidence: 0.91,
			RiskLevel:  "medium",
		},
	}}
	deps.LLM = &fakeLLM{answer: "Use the approved transcript request source and ask staff for account-specific details."}
	orch := mustNew(t, deps)

	resp, err := orch.HandleChat(traceContext("trace-llm-grounded"), domain.ChatRequest{
		ConversationID: "conv-llm-grounded",
		Channel:        "web",
		Message:        "How do I order my official transcript?",
	})
	if err != nil {
		t.Fatalf("HandleChat returned error: %v", err)
	}

	if !strings.Contains(resp.Answer, "approved transcript request source") {
		t.Fatalf("answer = %q, want guarded LLM answer", resp.Answer)
	}
	assertAction(t, resp, "llm_answer_generated", domain.ActionStatusCompleted)
}

type fakeLLM struct {
	answer string
	err    error
}

func (f fakeLLM) GenerateAnswer(context.Context, string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	return f.answer, nil
}
