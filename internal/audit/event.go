package audit

import (
	"context"
	"time"
)

const (
	EventTypeIntent     = "intent"
	EventTypeTool       = "tool"
	EventTypeWorkflow   = "workflow"
	EventTypeEscalation = "escalation"
	EventTypeGuardrail  = "guardrail"

	ActionClassify        = "intent_classified"
	ActionBannerLookup    = "banner_status_checked"
	ActionPaymentLookup   = "payment_status_checked"
	ActionPaymentReminder = "workflow_payment_reminder"
	ActionCreateCRMCase   = "crm_case_created"
	ActionSourceCheck     = "source_confirmation_required"

	StatusAttempted = "attempted"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusBlocked   = "blocked"
)

type Event struct {
	TraceID        string            `json:"trace_id,omitempty"`
	ConversationID string            `json:"conversation_id,omitempty"`
	StudentID      string            `json:"student_id,omitempty"`
	Type           string            `json:"type,omitempty"`
	Action         string            `json:"action,omitempty"`
	Status         string            `json:"status,omitempty"`
	ReferenceID    string            `json:"reference_id,omitempty"`
	Message        string            `json:"message,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	RecordedAt     time.Time         `json:"recorded_at,omitempty"`
}

type Recorder interface {
	Record(context.Context, Event) error
}

type NopRecorder struct{}

func (NopRecorder) Record(context.Context, Event) error {
	return nil
}
