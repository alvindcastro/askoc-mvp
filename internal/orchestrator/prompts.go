package orchestrator

import (
	"fmt"
	"strings"

	"askoc-mvp/internal/domain"
)

const PromptVersion = "p6.1"

func ClassificationPrompt(message string) string {
	return strings.Join([]string{
		"AskOC classifier prompt " + PromptVersion,
		"Return strict JSON only with keys: intent, intent_confidence, sentiment, urgency, needs_handoff, reason.",
		"Allowed intents: transcript_request, transcript_status, fee_payment, human_handoff, escalation_request, unknown.",
		"Allowed sentiments: neutral, negative, urgent, urgent_negative.",
		"Use only synthetic demo data. Do not request, reveal, or infer real learner records.",
		"Do not call tools or claim that tools were called. The Go orchestrator decides tool actions after validation.",
		"Classify the learner message below.",
		"Message: " + strings.TrimSpace(message),
	}, "\n")
}

func GroundedAnswerPrompt(question string, sources []domain.Source) string {
	lines := []string{
		"AskOC grounded answer prompt " + PromptVersion,
		"Use synthetic demo data only and never infer real learner records.",
		"You must answer only from the provided sources and include no unsupported claims.",
		"If account-specific details, fees, deadlines, eligibility, or private records are needed, ask staff for account-specific details.",
		"Question: " + strings.TrimSpace(question),
		"Sources:",
	}
	for i, source := range sources {
		lines = append(lines, fmt.Sprintf(
			"%d. title=%q url=%q chunk_id=%q confidence=%.2f risk=%q freshness=%q",
			i+1,
			source.Title,
			source.URL,
			source.ChunkID,
			source.Confidence,
			source.RiskLevel,
			source.FreshnessStatus,
		))
	}
	return strings.Join(lines, "\n")
}
