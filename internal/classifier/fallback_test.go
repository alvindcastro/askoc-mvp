package classifier_test

import (
	"context"
	"testing"

	"askoc-mvp/internal/classifier"
	"askoc-mvp/internal/domain"
)

func TestFallbackClassifierMapsMessagesToIntent(t *testing.T) {
	tests := []struct {
		name          string
		message       string
		wantIntent    domain.Intent
		minConfidence float64
	}{
		{
			name:          "transcript request",
			message:       "How do I order my official transcript?",
			wantIntent:    domain.IntentTranscriptRequest,
			minConfidence: 0.70,
		},
		{
			name:          "transcript status",
			message:       "I ordered my transcript but it has not been processed.",
			wantIntent:    domain.IntentTranscriptStatus,
			minConfidence: 0.75,
		},
		{
			name:          "fee payment",
			message:       "I paid my fee but my balance still shows owing.",
			wantIntent:    domain.IntentFeePayment,
			minConfidence: 0.70,
		},
		{
			name:          "human handoff",
			message:       "Can I talk to a person on staff?",
			wantIntent:    domain.IntentHumanHandoff,
			minConfidence: 0.70,
		},
		{
			name:          "unknown",
			message:       "What food is available near campus tonight?",
			wantIntent:    domain.IntentUnknown,
			minConfidence: 0.00,
		},
	}

	fallback := classifier.Fallback{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fallback.Classify(context.Background(), tt.message)
			if err != nil {
				t.Fatalf("Classify returned error: %v", err)
			}
			if got.Intent != tt.wantIntent {
				t.Fatalf("intent = %q, want %q", got.Intent, tt.wantIntent)
			}
			if got.Confidence < tt.minConfidence {
				t.Fatalf("confidence = %.2f, want >= %.2f", got.Confidence, tt.minConfidence)
			}
			if tt.wantIntent == domain.IntentUnknown && got.Confidence >= classifier.SensitiveToolConfidence {
				t.Fatalf("unknown confidence = %.2f, want below sensitive tool threshold", got.Confidence)
			}
		})
	}
}

func TestFallbackClassifierMapsSentimentAndUrgency(t *testing.T) {
	tests := []struct {
		name          string
		message       string
		wantSentiment domain.Sentiment
	}{
		{
			name:          "neutral",
			message:       "Please check my transcript status.",
			wantSentiment: domain.SentimentNeutral,
		},
		{
			name:          "negative",
			message:       "This is frustrating and I am upset about the delay.",
			wantSentiment: domain.SentimentNegative,
		},
		{
			name:          "urgent negative",
			message:       "This is extremely frustrating. I need this transcript today for a job application.",
			wantSentiment: domain.SentimentUrgentNegative,
		},
	}

	fallback := classifier.Fallback{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fallback.Classify(context.Background(), tt.message)
			if err != nil {
				t.Fatalf("Classify returned error: %v", err)
			}
			if got.Sentiment != tt.wantSentiment {
				t.Fatalf("sentiment = %q, want %q", got.Sentiment, tt.wantSentiment)
			}
		})
	}
}

func TestLowConfidenceCannotTriggerSensitiveTools(t *testing.T) {
	fallback := classifier.Fallback{}
	unknown, err := fallback.Classify(context.Background(), "Can you recommend a restaurant?")
	if err != nil {
		t.Fatalf("Classify returned error: %v", err)
	}
	if unknown.CanTriggerSensitiveTools() {
		t.Fatalf("unknown low-confidence result should not trigger sensitive tools: %+v", unknown)
	}

	transcript, err := fallback.Classify(context.Background(), "I ordered my transcript and it has not arrived.")
	if err != nil {
		t.Fatalf("Classify returned error: %v", err)
	}
	if !transcript.CanTriggerSensitiveTools() {
		t.Fatalf("high-confidence transcript status should allow typed tool checks: %+v", transcript)
	}
}
