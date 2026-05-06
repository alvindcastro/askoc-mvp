package classifier

import (
	"context"
	"strings"

	"askoc-mvp/internal/domain"
)

const SensitiveToolConfidence = 0.70

type Result struct {
	Intent       domain.Intent
	Confidence   float64
	Sentiment    domain.Sentiment
	Urgency      string
	NeedsHandoff bool
	Reason       string
}

func (r Result) CanTriggerSensitiveTools() bool {
	return r.Intent != domain.IntentUnknown && r.Confidence >= SensitiveToolConfidence
}

type Fallback struct{}

func (Fallback) Classify(ctx context.Context, message string) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}

	normalized := strings.ToLower(strings.TrimSpace(message))
	result := Result{
		Intent:     domain.IntentUnknown,
		Confidence: 0.35,
		Sentiment:  classifySentiment(normalized),
	}

	switch {
	case containsAny(normalized, "person", "human", "staff", "advisor", "representative", "talk to someone", "connect me", "learner services"):
		result.Intent = domain.IntentHumanHandoff
		result.Confidence = 0.78
	case containsAny(normalized, "payment", "paid", "fee", "balance", "owing", "charge"):
		result.Intent = domain.IntentFeePayment
		result.Confidence = 0.76
	case result.Sentiment == domain.SentimentUrgentNegative && containsAny(normalized, "transcript"):
		result.Intent = domain.IntentEscalationRequest
		result.Confidence = 0.74
	case containsAny(normalized, "processed", "status", "arrived", "check", "moving", "blocked", "not received", "where is", "does not have") && containsAny(normalized, "transcript"):
		result.Intent = domain.IntentTranscriptStatus
		result.Confidence = 0.86
	case containsAny(normalized, "order", "request", "how do i get", "send", "instructions", "copy") && containsAny(normalized, "transcript"):
		result.Intent = domain.IntentTranscriptRequest
		result.Confidence = 0.80
	}

	return result, nil
}

func classifySentiment(message string) domain.Sentiment {
	negative := containsAny(message, "frustrating", "frustrated", "upset", "angry", "extremely", "not acceptable", "unacceptable")
	urgent := containsAny(message, "urgent", "today", "asap", "deadline", "job application", "right away", "immediately")
	switch {
	case negative && urgent:
		return domain.SentimentUrgentNegative
	case negative:
		return domain.SentimentNegative
	case urgent:
		return domain.SentimentUrgent
	default:
		return domain.SentimentNeutral
	}
}

func containsAny(value string, terms ...string) bool {
	for _, term := range terms {
		if strings.Contains(value, term) {
			return true
		}
	}
	return false
}
