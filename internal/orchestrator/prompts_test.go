package orchestrator

import (
	"strings"
	"testing"

	"askoc-mvp/internal/domain"
)

func TestClassificationPromptContainsStrictJSONPrivacyAndAllowedLabels(t *testing.T) {
	prompt := ClassificationPrompt("Can you check transcript S100002?")

	checks := []string{
		"strict JSON",
		"intent_confidence",
		string(domain.IntentTranscriptStatus),
		string(domain.IntentFeePayment),
		string(domain.IntentHumanHandoff),
		string(domain.SentimentUrgentNegative),
		"synthetic demo data",
		"Do not call tools",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Fatalf("classification prompt missing %q:\n%s", want, prompt)
		}
	}
}

func TestGroundedAnswerPromptRequiresSourceOnlyAnswersAndPrivacyBoundary(t *testing.T) {
	prompt := GroundedAnswerPrompt("How do I order my transcript?", []domain.Source{
		{
			Title:      "Transcript Request Guidance",
			URL:        "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
			ChunkID:    "chunk-1",
			Confidence: 0.92,
		},
	})

	checks := []string{
		"answer only from the provided sources",
		"ask staff for account-specific details",
		"synthetic demo data",
		"https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
		"chunk-1",
	}
	for _, want := range checks {
		if !strings.Contains(prompt, want) {
			t.Fatalf("grounded answer prompt missing %q:\n%s", want, prompt)
		}
	}
}

func TestPromptVersionAndGoldenPromptCatchDrift(t *testing.T) {
	if PromptVersion != "p6.1" {
		t.Fatalf("PromptVersion = %q, want p6.1", PromptVersion)
	}

	got := ClassificationPrompt("How do I order my transcript?")
	want := strings.Join([]string{
		"AskOC classifier prompt p6.1",
		"Return strict JSON only with keys: intent, intent_confidence, sentiment, urgency, needs_handoff, reason.",
		"Allowed intents: transcript_request, transcript_status, fee_payment, human_handoff, escalation_request, unknown.",
		"Allowed sentiments: neutral, negative, urgent, urgent_negative.",
		"Use only synthetic demo data. Do not request, reveal, or infer real learner records.",
		"Do not call tools or claim that tools were called. The Go orchestrator decides tool actions after validation.",
		"Classify the learner message below.",
		"Message: How do I order my transcript?",
	}, "\n")
	if got != want {
		t.Fatalf("classification prompt drifted\nwant:\n%s\n\ngot:\n%s", want, got)
	}
}
