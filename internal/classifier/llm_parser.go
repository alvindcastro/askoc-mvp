package classifier

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"askoc-mvp/internal/domain"
)

const safeUnknownConfidence = 0.35

var ErrInvalidClassificationOutput = errors.New("invalid classification output")

type llmClassificationPayload struct {
	Intent           *string  `json:"intent"`
	IntentConfidence *float64 `json:"intent_confidence"`
	Sentiment        *string  `json:"sentiment"`
	Urgency          *string  `json:"urgency"`
	NeedsHandoff     *bool    `json:"needs_handoff"`
	Reason           *string  `json:"reason"`
}

func ParseLLMClassificationOutput(raw string) (Result, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return safeClassification(), invalidClassificationError("empty output")
	}

	var payload llmClassificationPayload
	decoder := json.NewDecoder(strings.NewReader(trimmed))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		return safeClassification(), fmt.Errorf("%w: %v", ErrInvalidClassificationOutput, err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return safeClassification(), invalidClassificationError("trailing JSON content")
	}

	if payload.Intent == nil || payload.IntentConfidence == nil || payload.Sentiment == nil || payload.Urgency == nil || payload.NeedsHandoff == nil || payload.Reason == nil {
		return safeClassification(), invalidClassificationError("missing required field")
	}
	if *payload.IntentConfidence < 0 || *payload.IntentConfidence > 1 {
		return safeClassification(), invalidClassificationError("confidence outside 0..1")
	}
	urgency := strings.TrimSpace(*payload.Urgency)
	if !validUrgency(urgency) {
		return safeClassification(), invalidClassificationError("unsupported urgency")
	}
	if strings.TrimSpace(*payload.Reason) == "" {
		return safeClassification(), invalidClassificationError("empty reason")
	}

	sentiment, ok := parseLLMSentiment(*payload.Sentiment)
	if !ok {
		return safeClassification(), invalidClassificationError("unsupported sentiment")
	}

	intent := parseLLMIntent(*payload.Intent)
	confidence := *payload.IntentConfidence
	if intent == domain.IntentUnknown {
		confidence = safeUnknownConfidence
	}
	if *payload.NeedsHandoff && intent == domain.IntentUnknown && confidence >= SensitiveToolConfidence {
		confidence = safeUnknownConfidence
	}

	return Result{
		Intent:       intent,
		Confidence:   confidence,
		Sentiment:    sentiment,
		Urgency:      urgency,
		NeedsHandoff: *payload.NeedsHandoff,
		Reason:       strings.TrimSpace(*payload.Reason),
	}, nil
}

func invalidClassificationError(reason string) error {
	return fmt.Errorf("%w: %s", ErrInvalidClassificationOutput, reason)
}

func safeClassification() Result {
	return Result{
		Intent:     domain.IntentUnknown,
		Confidence: safeUnknownConfidence,
		Sentiment:  domain.SentimentNeutral,
	}
}

func parseLLMIntent(value string) domain.Intent {
	switch domain.Intent(strings.TrimSpace(value)) {
	case domain.IntentTranscriptRequest:
		return domain.IntentTranscriptRequest
	case domain.IntentTranscriptStatus:
		return domain.IntentTranscriptStatus
	case domain.IntentFeePayment:
		return domain.IntentFeePayment
	case domain.IntentHumanHandoff:
		return domain.IntentHumanHandoff
	case domain.IntentEscalationRequest:
		return domain.IntentEscalationRequest
	case domain.IntentUnknown:
		return domain.IntentUnknown
	default:
		return domain.IntentUnknown
	}
}

func parseLLMSentiment(value string) (domain.Sentiment, bool) {
	switch domain.Sentiment(strings.TrimSpace(value)) {
	case domain.SentimentNeutral:
		return domain.SentimentNeutral, true
	case domain.SentimentNegative:
		return domain.SentimentNegative, true
	case domain.SentimentUrgent:
		return domain.SentimentUrgent, true
	case domain.SentimentUrgentNegative:
		return domain.SentimentUrgentNegative, true
	default:
		return domain.SentimentNeutral, false
	}
}

func validUrgency(value string) bool {
	switch value {
	case "low", "medium", "high":
		return true
	default:
		return false
	}
}
