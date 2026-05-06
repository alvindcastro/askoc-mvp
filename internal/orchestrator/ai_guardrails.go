package orchestrator

import (
	"context"
	"errors"
	"strings"

	"askoc-mvp/internal/audit"
	"askoc-mvp/internal/classifier"
	"askoc-mvp/internal/domain"
)

var (
	ErrUngroundedAnswer = errors.New("grounded answer requires approved sources")
	ErrEmptyLLMAnswer   = errors.New("llm answer is empty")
)

type ClassificationParser func(string) (classifier.Result, error)

type LLMBackedClassifier struct {
	LLM      LLM
	Fallback IntentClassifier
	Parse    ClassificationParser
	Audit    AuditRecorder
}

func (c LLMBackedClassifier) Classify(ctx context.Context, message string) (classifier.Result, error) {
	if c.Fallback == nil {
		return classifier.Result{}, errors.New("fallback classifier is required")
	}
	if c.LLM == nil {
		return c.fallback(ctx, message, "model_unavailable")
	}

	output, err := c.LLM.GenerateAnswer(ctx, ClassificationPrompt(message))
	if err != nil {
		return c.fallback(ctx, message, "model_error")
	}
	if strings.TrimSpace(output) == "" {
		return c.fallback(ctx, message, "empty_model_output")
	}
	if c.Parse == nil {
		return c.fallback(ctx, message, "parser_unavailable")
	}

	result, err := c.Parse(output)
	if err != nil {
		return c.fallback(ctx, message, "parser_error")
	}
	c.record(ctx, "llm_classification_validated", "completed", "classification accepted", map[string]string{
		"prompt_version": PromptVersion,
		"intent":         string(result.Intent),
	})
	return result, nil
}

func (c LLMBackedClassifier) fallback(ctx context.Context, message, reason string) (classifier.Result, error) {
	c.record(ctx, "llm_classification_fallback", "completed", "deterministic classifier used after LLM guardrail", map[string]string{
		"prompt_version": PromptVersion,
		"reason":         reason,
	})
	return c.Fallback.Classify(ctx, message)
}

func (c LLMBackedClassifier) record(ctx context.Context, action, status, message string, metadata map[string]string) {
	if c.Audit == nil {
		return
	}
	_ = c.Audit.Record(ctx, audit.Event{
		TraceID:  traceID(ctx),
		Type:     "guardrail",
		Action:   action,
		Status:   status,
		Message:  message,
		Metadata: metadata,
	})
}

func ValidateGroundedAnswer(answer string, sources []domain.Source) (string, error) {
	trimmed := strings.TrimSpace(answer)
	if trimmed == "" {
		return "", ErrEmptyLLMAnswer
	}
	if len(sources) == 0 {
		return "", ErrUngroundedAnswer
	}
	return trimmed, nil
}

func (o *Orchestrator) generateGroundedAnswer(ctx context.Context, question string, sources []domain.Source) (string, bool) {
	if o.llm == nil {
		return "", false
	}
	output, err := o.llm.GenerateAnswer(ctx, GroundedAnswerPrompt(question, sources))
	if err != nil {
		o.recordGuardrail(ctx, "llm_answer_fallback", "model_error")
		return "", false
	}
	answer, err := ValidateGroundedAnswer(output, sources)
	if err != nil {
		o.recordGuardrail(ctx, "llm_answer_rejected", "ungrounded_or_empty")
		return "", false
	}
	return answer, true
}

func (o *Orchestrator) recordGuardrail(ctx context.Context, action, reason string) {
	o.recordAudit(ctx, audit.Event{
		TraceID: traceID(ctx),
		Type:    "guardrail",
		Action:  action,
		Status:  "completed",
		Message: "LLM guardrail event",
		Metadata: map[string]string{
			"prompt_version": PromptVersion,
			"reason":         reason,
		},
	})
}
