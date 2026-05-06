package classifier_test

import (
	"errors"
	"testing"

	"askoc-mvp/internal/classifier"
	"askoc-mvp/internal/domain"
)

func TestParseLLMClassificationOutputValidJSONParses(t *testing.T) {
	got, err := classifier.ParseLLMClassificationOutput(`{"intent":"transcript_status","intent_confidence":0.82,"sentiment":"urgent","urgency":"high","needs_handoff":false,"reason":"Learner asks for transcript status."}`)
	if err != nil {
		t.Fatalf("ParseLLMClassificationOutput returned error: %v", err)
	}

	if got.Intent != domain.IntentTranscriptStatus {
		t.Fatalf("intent = %q, want %q", got.Intent, domain.IntentTranscriptStatus)
	}
	if got.Confidence != 0.82 {
		t.Fatalf("confidence = %.2f, want 0.82", got.Confidence)
	}
	if got.Sentiment != domain.SentimentUrgent {
		t.Fatalf("sentiment = %q, want %q", got.Sentiment, domain.SentimentUrgent)
	}
}

func TestParseLLMClassificationOutputMalformedJSONFailsSafely(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("ParseLLMClassificationOutput panicked on malformed JSON: %v", recovered)
		}
	}()

	got, err := classifier.ParseLLMClassificationOutput(`{"intent":`)
	if err == nil {
		t.Fatal("ParseLLMClassificationOutput returned nil error for malformed JSON")
	}
	if !errors.Is(err, classifier.ErrInvalidClassificationOutput) {
		t.Fatalf("error = %v, want ErrInvalidClassificationOutput", err)
	}
	assertSafeClassification(t, got)
}

func TestParseLLMClassificationOutputRejectsNonStrictJSON(t *testing.T) {
	tests := []struct {
		name string
		raw  string
	}{
		{
			name: "unknown field",
			raw:  `{"intent":"transcript_request","intent_confidence":0.91,"sentiment":"neutral","urgency":"low","needs_handoff":false,"reason":"Learner asks for transcript steps.","tool":"banner"}`,
		},
		{
			name: "trailing payload",
			raw:  `{"intent":"transcript_request","intent_confidence":0.91,"sentiment":"neutral","urgency":"low","needs_handoff":false,"reason":"Learner asks for transcript steps."} {"intent":"fee_payment","intent_confidence":1,"sentiment":"neutral","urgency":"low","needs_handoff":false,"reason":"second payload"}`,
		},
		{
			name: "missing sentiment",
			raw:  `{"intent":"transcript_request","intent_confidence":0.91,"urgency":"low","needs_handoff":false,"reason":"Learner asks for transcript steps."}`,
		},
		{
			name: "legacy confidence key",
			raw:  `{"intent":"transcript_request","confidence":0.91,"sentiment":"neutral","urgency":"low","needs_handoff":false,"reason":"Learner asks for transcript steps."}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := classifier.ParseLLMClassificationOutput(tt.raw)
			if err == nil {
				t.Fatal("ParseLLMClassificationOutput returned nil error for non-strict JSON")
			}
			if !errors.Is(err, classifier.ErrInvalidClassificationOutput) {
				t.Fatalf("error = %v, want ErrInvalidClassificationOutput", err)
			}
			assertSafeClassification(t, got)
		})
	}
}

func TestParseLLMClassificationOutputUnknownIntentMapsToUnknown(t *testing.T) {
	got, err := classifier.ParseLLMClassificationOutput(`{"intent":"housing_question","intent_confidence":0.94,"sentiment":"neutral","urgency":"low","needs_handoff":false,"reason":"Learner asks about housing."}`)
	if err != nil {
		t.Fatalf("ParseLLMClassificationOutput returned error: %v", err)
	}

	if got.Intent != domain.IntentUnknown {
		t.Fatalf("intent = %q, want %q", got.Intent, domain.IntentUnknown)
	}
	if got.Confidence >= classifier.SensitiveToolConfidence {
		t.Fatalf("unknown intent confidence = %.2f, want below sensitive tool threshold", got.Confidence)
	}
	if got.CanTriggerSensitiveTools() {
		t.Fatalf("unknown intent should not trigger sensitive tools: %+v", got)
	}
}

func TestParseLLMClassificationOutputRejectsOutOfRangeConfidence(t *testing.T) {
	tests := []struct {
		name string
		raw  string
	}{
		{
			name: "below zero",
			raw:  `{"intent":"transcript_status","intent_confidence":-0.01,"sentiment":"neutral","urgency":"low","needs_handoff":false,"reason":"Learner asks for transcript status."}`,
		},
		{
			name: "above one",
			raw:  `{"intent":"transcript_status","intent_confidence":1.01,"sentiment":"neutral","urgency":"low","needs_handoff":false,"reason":"Learner asks for transcript status."}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := classifier.ParseLLMClassificationOutput(tt.raw)
			if err == nil {
				t.Fatal("ParseLLMClassificationOutput returned nil error for out-of-range confidence")
			}
			if !errors.Is(err, classifier.ErrInvalidClassificationOutput) {
				t.Fatalf("error = %v, want ErrInvalidClassificationOutput", err)
			}
			assertSafeClassification(t, got)
		})
	}
}

func TestParseLLMClassificationOutputDisablesToolsBelowThreshold(t *testing.T) {
	got, err := classifier.ParseLLMClassificationOutput(`{"intent":"transcript_status","intent_confidence":0.69,"sentiment":"neutral","urgency":"low","needs_handoff":false,"reason":"Learner asks for transcript status."}`)
	if err != nil {
		t.Fatalf("ParseLLMClassificationOutput returned error: %v", err)
	}

	if got.CanTriggerSensitiveTools() {
		t.Fatalf("low-confidence parser result should not trigger sensitive tools: %+v", got)
	}
}

func TestParseLLMClassificationOutputPreservesHandoffReasonAndUrgency(t *testing.T) {
	got, err := classifier.ParseLLMClassificationOutput(`{"intent":"transcript_status","intent_confidence":0.86,"sentiment":"negative","urgency":"high","needs_handoff":true,"reason":"Learner reports a blocked transcript and needs staff review."}`)
	if err != nil {
		t.Fatalf("ParseLLMClassificationOutput returned error: %v", err)
	}

	if !got.NeedsHandoff {
		t.Fatalf("NeedsHandoff = false, want true")
	}
	if got.Urgency != "high" {
		t.Fatalf("Urgency = %q, want high", got.Urgency)
	}
	if got.Reason != "Learner reports a blocked transcript and needs staff review." {
		t.Fatalf("Reason = %q", got.Reason)
	}
}

func assertSafeClassification(t *testing.T, got classifier.Result) {
	t.Helper()

	if got.Intent != domain.IntentUnknown {
		t.Fatalf("safe intent = %q, want %q", got.Intent, domain.IntentUnknown)
	}
	if got.Confidence >= classifier.SensitiveToolConfidence {
		t.Fatalf("safe confidence = %.2f, want below %.2f", got.Confidence, classifier.SensitiveToolConfidence)
	}
	if got.Sentiment != domain.SentimentNeutral {
		t.Fatalf("safe sentiment = %q, want %q", got.Sentiment, domain.SentimentNeutral)
	}
	if got.CanTriggerSensitiveTools() {
		t.Fatalf("safe classification should not trigger sensitive tools: %+v", got)
	}
}
