package llm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Provider interface {
	GenerateAnswer(context.Context, AnswerRequest) (AnswerResponse, error)
}

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type GroundingSource struct {
	ID              string  `json:"id,omitempty"`
	Title           string  `json:"title"`
	URL             string  `json:"url"`
	ChunkID         string  `json:"chunk_id"`
	Confidence      float64 `json:"confidence,omitempty"`
	RiskLevel       string  `json:"risk_level,omitempty"`
	FreshnessStatus string  `json:"freshness_status,omitempty"`
	Caution         string  `json:"caution,omitempty"`
}

type Citation struct {
	Title      string  `json:"title"`
	URL        string  `json:"url"`
	ChunkID    string  `json:"chunk_id"`
	Confidence float64 `json:"confidence,omitempty"`
}

type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

type AnswerRequest struct {
	ConversationID string            `json:"conversation_id,omitempty"`
	TraceID        string            `json:"trace_id,omitempty"`
	Messages       []Message         `json:"messages"`
	Sources        []GroundingSource `json:"sources,omitempty"`
	MaxTokens      int               `json:"max_tokens,omitempty"`
	Temperature    float64           `json:"temperature,omitempty"`
	Timeout        time.Duration     `json:"-"`
}

func (r AnswerRequest) Validate() error {
	if r.Timeout <= 0 {
		return &ProviderError{
			Kind:     KindValidation,
			Provider: "llm",
			Message:  "timeout must be greater than zero",
		}
	}
	if len(r.Messages) == 0 {
		return &ProviderError{
			Kind:     KindValidation,
			Provider: "llm",
			Message:  "at least one message is required",
		}
	}
	for _, msg := range r.Messages {
		if strings.TrimSpace(string(msg.Role)) == "" || strings.TrimSpace(msg.Content) == "" {
			return &ProviderError{
				Kind:     KindValidation,
				Provider: "llm",
				Message:  "messages require role and content",
			}
		}
	}
	return nil
}

type AnswerResponse struct {
	Answer       string     `json:"answer"`
	Sources      []Citation `json:"sources,omitempty"`
	Model        string     `json:"model,omitempty"`
	FinishReason string     `json:"finish_reason,omitempty"`
	Usage        TokenUsage `json:"usage,omitempty"`
}

type ErrorKind string

const (
	KindValidation  ErrorKind = "validation"
	KindRateLimited ErrorKind = "rate_limited"
	KindRetryable   ErrorKind = "retryable"
	KindExternal    ErrorKind = "external_service"
	KindParse       ErrorKind = "parse"
	KindTimeout     ErrorKind = "timeout"
)

type ProviderError struct {
	Kind       ErrorKind
	Provider   string
	StatusCode int
	Message    string
	Err        error
}

func (e *ProviderError) Error() string {
	if e == nil {
		return ""
	}
	provider := e.Provider
	if provider == "" {
		provider = "llm"
	}
	if e.Message != "" {
		return fmt.Sprintf("%s %s: %s", provider, e.Kind, e.Message)
	}
	return fmt.Sprintf("%s %s", provider, e.Kind)
}

func (e *ProviderError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func IsProviderErrorKind(err error, kind ErrorKind) bool {
	var providerErr *ProviderError
	return errors.As(err, &providerErr) && providerErr.Kind == kind
}
