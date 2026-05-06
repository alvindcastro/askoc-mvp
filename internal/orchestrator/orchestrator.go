package orchestrator

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"askoc-mvp/internal/audit"
	"askoc-mvp/internal/classifier"
	"askoc-mvp/internal/domain"
	"askoc-mvp/internal/middleware"
	"askoc-mvp/internal/privacy"
	"askoc-mvp/internal/tools"
	"askoc-mvp/internal/workflow"
)

var syntheticIDPattern = regexp.MustCompile(`\bS[0-9]{6}\b`)

type IntentClassifier interface {
	Classify(context.Context, string) (classifier.Result, error)
}

type Retriever interface {
	RetrieveSources(context.Context, string) ([]domain.Source, error)
}

type LLM interface {
	GenerateAnswer(context.Context, string) (string, error)
}

type BannerTool interface {
	GetTranscriptStatus(context.Context, string) (tools.BannerTranscriptStatus, error)
}

type PaymentTool interface {
	GetPaymentStatus(context.Context, string) (tools.PaymentStatus, error)
}

type WorkflowTool interface {
	SendPaymentReminder(context.Context, workflow.PaymentReminderRequest) (workflow.PaymentReminderResponse, error)
}

type CRMTool interface {
	CreateCase(context.Context, tools.CRMCaseRequest) (tools.CRMCaseResponse, error)
}

type AuditRecorder interface {
	Record(context.Context, audit.Event) error
}

type Dependencies struct {
	Classifier IntentClassifier
	Retriever  Retriever
	LLM        LLM
	Banner     BannerTool
	Payment    PaymentTool
	Workflow   WorkflowTool
	CRM        CRMTool
	Audit      AuditRecorder
	Redact     func(string) string
}

type Orchestrator struct {
	classifier IntentClassifier
	retriever  Retriever
	llm        LLM
	banner     BannerTool
	payment    PaymentTool
	workflow   WorkflowTool
	crm        CRMTool
	audit      AuditRecorder
	redact     func(string) string
}

func New(deps Dependencies) (*Orchestrator, error) {
	if deps.Classifier == nil {
		return nil, errors.New("orchestrator classifier dependency is required")
	}
	if deps.Retriever == nil {
		return nil, errors.New("orchestrator retriever dependency is required")
	}
	if deps.LLM == nil {
		return nil, errors.New("orchestrator llm dependency is required")
	}
	if deps.Banner == nil {
		return nil, errors.New("orchestrator banner dependency is required")
	}
	if deps.Payment == nil {
		return nil, errors.New("orchestrator payment dependency is required")
	}
	if deps.Workflow == nil {
		return nil, errors.New("orchestrator workflow dependency is required")
	}
	if deps.CRM == nil {
		return nil, errors.New("orchestrator crm dependency is required")
	}
	if deps.Audit == nil {
		return nil, errors.New("orchestrator audit dependency is required")
	}
	if deps.Redact == nil {
		deps.Redact = privacy.Redact
	}

	return &Orchestrator{
		classifier: deps.Classifier,
		retriever:  deps.Retriever,
		llm:        deps.LLM,
		banner:     deps.Banner,
		payment:    deps.Payment,
		workflow:   deps.Workflow,
		crm:        deps.CRM,
		audit:      deps.Audit,
		redact:     deps.Redact,
	}, nil
}

func (o *Orchestrator) HandleChat(ctx context.Context, req domain.ChatRequest) (resp domain.ChatResponse, err error) {
	traceID := traceID(ctx)
	conversationID := conversationIDFor(req)

	result, err := o.classifier.Classify(ctx, req.Message)
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf("classify chat: %w", err)
	}
	defer func() {
		if err == nil {
			o.recordResponseAudit(ctx, req, result, resp)
		}
	}()

	resp = domain.ChatResponse{
		ConversationID: conversationID,
		TraceID:        traceID,
		Intent: domain.IntentResult{
			Name:       result.Intent,
			Confidence: result.Confidence,
		},
		Sentiment: result.Sentiment,
		Actions: []domain.Action{
			o.action(ctx, "intent_classified", domain.ActionStatusCompleted, "Message classified by validated classifier logic.", ""),
		},
	}

	if !result.CanTriggerSensitiveTools() {
		resp.Answer = "I could not validate that classification strongly enough for synthetic tool checks, so I created a normal staff handoff instead of checking transcript or payment records."
		resp.Actions = append(resp.Actions, o.action(ctx, "classification_guardrail", domain.ActionStatusPending, "Low-confidence classification blocked sensitive synthetic tool calls.", ""))
		return o.createHandoff(ctx, req, result, resp, handoffRequest{
			queue:    "learner_support",
			priority: "normal",
			reason:   "classification below sensitive tool threshold",
		})
	}

	switch {
	case result.Intent == domain.IntentTranscriptStatus || result.Intent == domain.IntentFeePayment:
		resp = o.attachGroundingIfAvailable(ctx, req, resp)
		return o.handleTranscriptStatus(ctx, req, result, resp)
	case result.Intent == domain.IntentTranscriptRequest:
		return o.handleGroundedAnswer(ctx, req, resp), nil
	case needsEscalation(result):
		return o.createHandoff(ctx, req, result, resp, handoffRequest{
			queue:    queueForResult(result),
			priority: priorityForResult(result),
			reason:   "conversation needs staff review",
		})
	default:
		resp.Answer = "I could not confidently match that to the transcript/payment demo flow. I queued a normal synthetic handoff for staff review."
		return o.createHandoff(ctx, req, result, resp, handoffRequest{
			queue:    "learner_support",
			priority: "normal",
			reason:   "low confidence fallback",
		})
	}
}

type DisabledRetriever struct{}

func (DisabledRetriever) RetrieveSources(context.Context, string) ([]domain.Source, error) {
	return nil, nil
}

type DisabledLLM struct{}

func (DisabledLLM) GenerateAnswer(context.Context, string) (string, error) {
	return "", nil
}

func conversationIDFor(req domain.ChatRequest) string {
	if got := strings.TrimSpace(req.ConversationID); got != "" {
		return got
	}
	key := strings.TrimSpace(req.StudentID)
	if key == "" {
		key = strings.TrimSpace(req.Channel) + ":" + strings.TrimSpace(req.Message)
	}
	hash := sha1.Sum([]byte(key))
	return "conv_" + hex.EncodeToString(hash[:])[:12]
}

func traceID(ctx context.Context) string {
	if got := strings.TrimSpace(middleware.TraceIDFromContext(ctx)); got != "" {
		return got
	}
	return "trace-unset"
}

func (o *Orchestrator) action(ctx context.Context, actionType string, status domain.ActionStatus, message, referenceID string) domain.Action {
	return domain.Action{
		Type:        actionType,
		Status:      status,
		Message:     message,
		ReferenceID: referenceID,
		TraceID:     traceID(ctx),
	}
}

func withIdempotency(action domain.Action, key string) domain.Action {
	action.IdempotencyKey = key
	return action
}

func (o *Orchestrator) recordResponseAudit(ctx context.Context, req domain.ChatRequest, result classifier.Result, resp domain.ChatResponse) {
	for _, action := range resp.Actions {
		metadata := map[string]string{
			"intent":     string(resp.Intent.Name),
			"confidence": strconv.FormatFloat(resp.Intent.Confidence, 'f', 2, 64),
			"sentiment":  string(resp.Sentiment),
		}
		if strings.TrimSpace(req.Message) != "" {
			metadata["question"] = o.redact(req.Message)
		}
		if resp.Intent.Confidence > 0 && resp.Intent.Confidence < classifier.SensitiveToolConfidence {
			metadata["low_confidence"] = "true"
		}
		if action.Type == "source_confirmation_required" {
			metadata["stale_source"] = "true"
		}
		if strings.TrimSpace(action.IdempotencyKey) != "" {
			for key, value := range workflowAuditMetadata(action.IdempotencyKey, 0) {
				metadata[key] = value
			}
		}
		if resp.Escalation != nil {
			metadata["queue"] = resp.Escalation.Queue
			metadata["priority"] = resp.Escalation.Priority
			metadata["case_id"] = resp.Escalation.CaseID
		}
		if result.NeedsHandoff {
			metadata["needs_handoff"] = "true"
		}
		o.recordAudit(ctx, audit.Event{
			TraceID:        traceID(ctx),
			ConversationID: resp.ConversationID,
			StudentID:      req.StudentID,
			Type:           auditTypeForAction(action.Type),
			Action:         action.Type,
			Status:         string(action.Status),
			ReferenceID:    action.ReferenceID,
			Message:        o.redact(action.Message),
			Metadata:       metadata,
		})
	}
}

func auditTypeForAction(actionType string) string {
	switch actionType {
	case "intent_classified":
		return audit.EventTypeIntent
	case "classification_guardrail", "source_confirmation_required":
		return audit.EventTypeGuardrail
	case "payment_reminder_triggered":
		return audit.EventTypeWorkflow
	case "crm_case_created":
		return audit.EventTypeEscalation
	default:
		return audit.EventTypeTool
	}
}
