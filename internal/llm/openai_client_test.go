package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestOpenAIClientSendsPayloadAndParsesSuccessResponse(t *testing.T) {
	var gotRequest openAIChatRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/openai/deployments/demo/chat/completions" {
			t.Fatalf("path = %q, want Azure-compatible chat completions path", r.URL.Path)
		}
		if r.URL.Query().Get("api-version") != "2024-02-15-preview" {
			t.Fatalf("api-version query = %q", r.URL.RawQuery)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer config-key" {
			t.Fatalf("Authorization header = %q, want Bearer config-key", got)
		}
		if got := r.Header.Get("api-key"); got != "config-key" {
			t.Fatalf("api-key header = %q, want config-key", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want application/json", got)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotRequest); err != nil {
			t.Fatalf("decode request: %v", err)
		}

		writeOpenAIResponse(t, w, openAIChatResponse{
			ID:    "chatcmpl-demo",
			Model: "gpt-demo",
			Choices: []openAIChoice{
				{
					Index:        0,
					Message:      Message{Role: RoleAssistant, Content: "Use the approved public transcript request page."},
					FinishReason: "stop",
				},
			},
			Usage: TokenUsage{PromptTokens: 10, CompletionTokens: 7, TotalTokens: 17},
		})
	}))
	defer server.Close()

	t.Setenv("ASKOC_PROVIDER_API_KEY", "env-secret-must-not-be-used")
	client, err := NewOpenAIClient(OpenAIClientConfig{
		Endpoint:   server.URL + "/openai/deployments/demo/chat/completions?api-version=2024-02-15-preview",
		APIKey:     "config-key",
		Model:      "gpt-demo",
		HTTPClient: server.Client(),
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatalf("NewOpenAIClient() error = %v", err)
	}

	resp, err := client.GenerateAnswer(context.Background(), AnswerRequest{
		Messages: []Message{
			{Role: RoleSystem, Content: "Use approved public sources only."},
			{Role: RoleUser, Content: "How do I request an official transcript?"},
		},
		MaxTokens:   256,
		Temperature: 0.2,
		Timeout:     time.Second,
	})
	if err != nil {
		t.Fatalf("GenerateAnswer() error = %v", err)
	}

	if gotRequest.Model != "gpt-demo" {
		t.Fatalf("model = %q, want gpt-demo", gotRequest.Model)
	}
	if len(gotRequest.Messages) != 2 || gotRequest.Messages[0].Role != RoleSystem || gotRequest.Messages[1].Role != RoleUser {
		t.Fatalf("messages = %+v, want provider-neutral roles preserved", gotRequest.Messages)
	}
	if gotRequest.MaxTokens != 256 || gotRequest.Temperature != 0.2 {
		t.Fatalf("request tuning = max_tokens:%d temperature:%v", gotRequest.MaxTokens, gotRequest.Temperature)
	}
	if resp.Answer != "Use the approved public transcript request page." {
		t.Fatalf("answer = %q", resp.Answer)
	}
	if resp.Model != "gpt-demo" || resp.FinishReason != "stop" {
		t.Fatalf("response metadata = model:%q finish:%q", resp.Model, resp.FinishReason)
	}
	if resp.Usage.TotalTokens != 17 {
		t.Fatalf("usage = %+v, want total token count parsed", resp.Usage)
	}
}

func TestOpenAIClientMaps429And500WithoutLeakingPromptText(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		wantKind ErrorKind
	}{
		{name: "rate limited", status: http.StatusTooManyRequests, wantKind: KindRateLimited},
		{name: "server error", status: http.StatusInternalServerError, wantKind: KindRetryable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(`{"error":{"message":"model saw synthetic student S100002 asking about payment"}}`))
			}))
			defer server.Close()

			client := mustOpenAIClient(t, server.URL, server.Client())
			_, err := client.GenerateAnswer(context.Background(), validAnswerRequest())
			if !IsProviderErrorKind(err, tt.wantKind) {
				t.Fatalf("GenerateAnswer() error = %v, want kind %s", err, tt.wantKind)
			}
			if strings.Contains(err.Error(), "S100002") || strings.Contains(err.Error(), "payment") {
				t.Fatalf("error leaked prompt/provider body detail: %q", err.Error())
			}
		})
	}
}

func TestOpenAIClientMapsRequestTimeout(t *testing.T) {
	started := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := mustOpenAIClient(t, server.URL, server.Client())
	req := validAnswerRequest()
	req.Timeout = 20 * time.Millisecond

	_, err := client.GenerateAnswer(context.Background(), req)
	if !IsProviderErrorKind(err, KindTimeout) {
		t.Fatalf("GenerateAnswer() error = %v, want kind %s", err, KindTimeout)
	}
	select {
	case <-started:
	default:
		t.Fatalf("server did not receive timeout test request")
	}
}

func TestOpenAIClientRespectsAlreadyCanceledContext(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := mustOpenAIClient(t, server.URL, server.Client())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GenerateAnswer(ctx, validAnswerRequest())
	if !IsProviderErrorKind(err, KindTimeout) {
		t.Fatalf("GenerateAnswer() error = %v, want kind %s", err, KindTimeout)
	}
	if got := atomic.LoadInt32(&requests); got != 0 {
		t.Fatalf("server received %d requests after context cancellation, want 0", got)
	}
}

func mustOpenAIClient(t *testing.T, endpoint string, httpClient *http.Client) *OpenAIClient {
	t.Helper()
	client, err := NewOpenAIClient(OpenAIClientConfig{
		Endpoint:   endpoint,
		APIKey:     "config-key",
		Model:      "gpt-demo",
		HTTPClient: httpClient,
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatalf("NewOpenAIClient() error = %v", err)
	}
	return client
}

func writeOpenAIResponse(t *testing.T, w http.ResponseWriter, payload openAIChatResponse) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Fatalf("encode response: %v", err)
	}
}
