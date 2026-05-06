package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const openAIProviderName = "openai-compatible"

type OpenAIClientConfig struct {
	Endpoint   string
	APIKey     string
	Model      string
	HTTPClient *http.Client
	Timeout    time.Duration
}

type OpenAIClient struct {
	endpoint   string
	apiKey     string
	model      string
	httpClient *http.Client
	timeout    time.Duration
}

func NewOpenAIClient(cfg OpenAIClientConfig) (*OpenAIClient, error) {
	endpoint := strings.TrimSpace(cfg.Endpoint)
	if endpoint == "" {
		return nil, validationError("endpoint is required")
	}
	if _, err := url.ParseRequestURI(endpoint); err != nil {
		return nil, validationError("endpoint must be an absolute URL")
	}
	apiKey := strings.TrimSpace(cfg.APIKey)
	if apiKey == "" {
		return nil, validationError("api key is required")
	}
	model := strings.TrimSpace(cfg.Model)
	if model == "" {
		return nil, validationError("model is required")
	}
	if cfg.Timeout <= 0 {
		return nil, validationError("timeout must be greater than zero")
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	clientCopy := *httpClient
	if clientCopy.Timeout == 0 || clientCopy.Timeout > cfg.Timeout {
		clientCopy.Timeout = cfg.Timeout
	}

	return &OpenAIClient{
		endpoint:   endpoint,
		apiKey:     apiKey,
		model:      model,
		httpClient: &clientCopy,
		timeout:    cfg.Timeout,
	}, nil
}

func (c *OpenAIClient) GenerateAnswer(ctx context.Context, req AnswerRequest) (AnswerResponse, error) {
	if err := ctx.Err(); err != nil {
		return AnswerResponse{}, timeoutError("request context ended before provider call", err)
	}
	if err := req.Validate(); err != nil {
		return AnswerResponse{}, err
	}

	timeout := req.Timeout
	if timeout <= 0 {
		timeout = c.timeout
	}
	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	payload := openAIChatRequest{
		Model:       c.model,
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return AnswerResponse{}, &ProviderError{
			Kind:     KindParse,
			Provider: openAIProviderName,
			Message:  "unable to encode provider request",
			Err:      err,
		}
	}

	httpReq, err := http.NewRequestWithContext(callCtx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return AnswerResponse{}, &ProviderError{
			Kind:     KindExternal,
			Provider: openAIProviderName,
			Message:  "unable to create provider request",
			Err:      err,
		}
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("api-key", c.apiKey)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		if callCtx.Err() != nil || isTimeout(err) {
			return AnswerResponse{}, timeoutError("provider request timed out or was canceled", err)
		}
		return AnswerResponse{}, &ProviderError{
			Kind:     KindExternal,
			Provider: openAIProviderName,
			Message:  "provider request failed",
			Err:      err,
		}
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, httpResp.Body)
		return AnswerResponse{}, statusError(httpResp.StatusCode)
	}

	var decoded openAIChatResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&decoded); err != nil {
		return AnswerResponse{}, &ProviderError{
			Kind:       KindParse,
			Provider:   openAIProviderName,
			StatusCode: httpResp.StatusCode,
			Message:    "provider response was not valid JSON",
			Err:        err,
		}
	}
	if len(decoded.Choices) == 0 {
		return AnswerResponse{}, &ProviderError{
			Kind:       KindParse,
			Provider:   openAIProviderName,
			StatusCode: httpResp.StatusCode,
			Message:    "provider response did not include a choice",
		}
	}

	choice := decoded.Choices[0]
	return AnswerResponse{
		Answer:       choice.Message.Content,
		Model:        decoded.Model,
		FinishReason: choice.FinishReason,
		Usage:        decoded.Usage,
	}, nil
}

type openAIChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type openAIChatResponse struct {
	ID      string         `json:"id,omitempty"`
	Model   string         `json:"model,omitempty"`
	Choices []openAIChoice `json:"choices"`
	Usage   TokenUsage     `json:"usage,omitempty"`
}

type openAIChoice struct {
	Index        int     `json:"index,omitempty"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason,omitempty"`
}

func validationError(message string) error {
	return &ProviderError{
		Kind:     KindValidation,
		Provider: openAIProviderName,
		Message:  message,
	}
}

func statusError(status int) error {
	kind := KindExternal
	switch {
	case status == http.StatusTooManyRequests:
		kind = KindRateLimited
	case status >= 500:
		kind = KindRetryable
	}
	return &ProviderError{
		Kind:       kind,
		Provider:   openAIProviderName,
		StatusCode: status,
		Message:    fmt.Sprintf("provider returned HTTP %d", status),
	}
}

func timeoutError(message string, err error) error {
	return &ProviderError{
		Kind:     KindTimeout,
		Provider: openAIProviderName,
		Message:  message,
		Err:      err,
	}
}

func isTimeout(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
