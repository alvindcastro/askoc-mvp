package llm

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestAnswerRequestJSONSchemaMarshalsProviderNeutralPayload(t *testing.T) {
	req := AnswerRequest{
		ConversationID: "conv-demo",
		TraceID:        "trace-demo",
		Messages: []Message{
			{Role: RoleSystem, Content: "Use approved public sources only."},
			{Role: RoleUser, Content: "How do I request an official transcript?"},
		},
		Sources: []GroundingSource{
			{
				ID:              "oc-transcript-request",
				Title:           "Transcript Request Guidance",
				URL:             "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
				ChunkID:         "transcript-chunk-1",
				Confidence:      0.91,
				RiskLevel:       "high",
				FreshnessStatus: "fresh",
			},
		},
		MaxTokens:   400,
		Temperature: 0.2,
		Timeout:     3 * time.Second,
	}

	got, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal(AnswerRequest) error = %v", err)
	}

	want := `{"conversation_id":"conv-demo","trace_id":"trace-demo","messages":[{"role":"system","content":"Use approved public sources only."},{"role":"user","content":"How do I request an official transcript?"}],"sources":[{"id":"oc-transcript-request","title":"Transcript Request Guidance","url":"https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards","chunk_id":"transcript-chunk-1","confidence":0.91,"risk_level":"high","freshness_status":"fresh"}],"max_tokens":400,"temperature":0.2}`
	if string(got) != want {
		t.Fatalf("JSON = %s, want %s", got, want)
	}
	if strings.Contains(string(got), "timeout") || strings.Contains(string(got), "Timeout") {
		t.Fatalf("provider-neutral JSON leaked runtime timeout: %s", got)
	}
}

func TestAnswerResponseJSONSchemaMarshalsProviderNeutralPayload(t *testing.T) {
	resp := AnswerResponse{
		Answer: "Use the approved public transcript request page for official steps.",
		Sources: []Citation{
			{
				Title:      "Transcript Request Guidance",
				URL:        "https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards",
				ChunkID:    "transcript-chunk-1",
				Confidence: 0.91,
			},
		},
		Model:        "demo-model",
		FinishReason: "stop",
		Usage: TokenUsage{
			PromptTokens:     12,
			CompletionTokens: 9,
			TotalTokens:      21,
		},
	}

	got, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Marshal(AnswerResponse) error = %v", err)
	}

	want := `{"answer":"Use the approved public transcript request page for official steps.","sources":[{"title":"Transcript Request Guidance","url":"https://www.okanagancollege.ca/ask-oc/transcript-request-2005-onwards","chunk_id":"transcript-chunk-1","confidence":0.91}],"model":"demo-model","finish_reason":"stop","usage":{"prompt_tokens":12,"completion_tokens":9,"total_tokens":21}}`
	if string(got) != want {
		t.Fatalf("JSON = %s, want %s", got, want)
	}
}

func TestProviderErrorSupportsKindChecksAndUnwrap(t *testing.T) {
	cause := errors.New("upstream throttled")
	err := &ProviderError{
		Kind:       KindRateLimited,
		Provider:   "openai-compatible",
		StatusCode: 429,
		Message:    "provider request was rate limited",
		Err:        cause,
	}

	if !IsProviderErrorKind(err, KindRateLimited) {
		t.Fatalf("IsProviderErrorKind() = false, want true")
	}
	if !errors.Is(err, cause) {
		t.Fatalf("errors.Is() did not unwrap cause")
	}

	var got *ProviderError
	if !errors.As(err, &got) {
		t.Fatalf("errors.As() did not find ProviderError")
	}
	if got.Provider != "openai-compatible" || got.StatusCode != 429 {
		t.Fatalf("ProviderError = %+v, want provider and status metadata", got)
	}
	if strings.Contains(err.Error(), "upstream throttled") {
		t.Fatalf("Error() leaked wrapped provider detail: %q", err.Error())
	}
}

func TestAnswerRequestValidateRequiresPositiveTimeout(t *testing.T) {
	valid := validAnswerRequest()
	if err := valid.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	for _, timeout := range []time.Duration{0, -1 * time.Second} {
		t.Run(timeout.String(), func(t *testing.T) {
			req := validAnswerRequest()
			req.Timeout = timeout

			err := req.Validate()
			if !IsProviderErrorKind(err, KindValidation) {
				t.Fatalf("Validate() error = %v, want kind %s", err, KindValidation)
			}
		})
	}
}

func validAnswerRequest() AnswerRequest {
	return AnswerRequest{
		Messages: []Message{
			{Role: RoleSystem, Content: "Use approved public sources only."},
			{Role: RoleUser, Content: "How do I request an official transcript?"},
		},
		MaxTokens:   400,
		Temperature: 0.2,
		Timeout:     3 * time.Second,
	}
}
