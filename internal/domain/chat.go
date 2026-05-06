package domain

type Intent string

const (
	IntentTranscriptRequest Intent = "transcript_request"
	IntentTranscriptStatus  Intent = "transcript_status"
	IntentFeePayment        Intent = "fee_payment"
	IntentHumanHandoff      Intent = "human_handoff"
	IntentEscalationRequest Intent = "escalation_request"
	IntentUnknown           Intent = "unknown"
)

type Sentiment string

const (
	SentimentNeutral        Sentiment = "neutral"
	SentimentNegative       Sentiment = "negative"
	SentimentUrgent         Sentiment = "urgent"
	SentimentUrgentNegative Sentiment = "urgent_negative"
)

type ActionStatus string

const (
	ActionStatusCompleted ActionStatus = "completed"
	ActionStatusPending   ActionStatus = "pending"
	ActionStatusSkipped   ActionStatus = "skipped"
)

type HandoffStatus string

const (
	HandoffNone    HandoffStatus = "none"
	HandoffPending HandoffStatus = "pending"
)

type ChatRequest struct {
	ConversationID string `json:"conversation_id,omitempty"`
	Channel        string `json:"channel"`
	Message        string `json:"message"`
	StudentID      string `json:"student_id,omitempty"`
}

type ChatResponse struct {
	ConversationID string       `json:"conversation_id"`
	TraceID        string       `json:"trace_id"`
	Answer         string       `json:"answer"`
	Intent         IntentResult `json:"intent"`
	Sentiment      Sentiment    `json:"sentiment"`
	Sources        []Source     `json:"sources"`
	Actions        []Action     `json:"actions"`
	Escalation     *Escalation  `json:"escalation,omitempty"`
}

type IntentResult struct {
	Name       Intent  `json:"name"`
	Confidence float64 `json:"confidence"`
}

type Source struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	ChunkID string `json:"chunk_id"`
}

type Action struct {
	Type           string       `json:"type"`
	Status         ActionStatus `json:"status"`
	Message        string       `json:"message,omitempty"`
	ReferenceID    string       `json:"reference_id,omitempty"`
	TraceID        string       `json:"trace_id,omitempty"`
	IdempotencyKey string       `json:"idempotency_key,omitempty"`
}

type Escalation struct {
	Required bool          `json:"required"`
	Status   HandoffStatus `json:"status"`
	Queue    string        `json:"queue,omitempty"`
	Priority string        `json:"priority,omitempty"`
	CaseID   string        `json:"case_id,omitempty"`
	Reason   string        `json:"reason,omitempty"`
}
