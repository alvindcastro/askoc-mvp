package workflow

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
)

type PaymentReminderRequest struct {
	StudentID      string  `json:"student_id"`
	ConversationID string  `json:"conversation_id,omitempty"`
	TraceID        string  `json:"trace_id,omitempty"`
	Item           string  `json:"item"`
	AmountDue      float64 `json:"amount_due,omitempty"`
	Currency       string  `json:"currency,omitempty"`
	Reason         string  `json:"reason"`
	IdempotencyKey string  `json:"idempotency_key"`
}

type PaymentReminderResponse struct {
	WorkflowID     string `json:"workflow_id"`
	Status         string `json:"status"`
	Message        string `json:"message,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	Synthetic      bool   `json:"synthetic"`
}

type PaymentReminderSender interface {
	SendPaymentReminder(context.Context, PaymentReminderRequest) (PaymentReminderResponse, error)
}

func PaymentReminderKey(traceID, studentID, item string) string {
	return "payment-reminder:" + strings.TrimSpace(traceID) + ":" + strings.TrimSpace(studentID) + ":" + strings.TrimSpace(item)
}

type InMemoryClient struct {
	mu        sync.Mutex
	responses map[string]PaymentReminderResponse
}

func NewInMemoryClient() *InMemoryClient {
	return &InMemoryClient{responses: make(map[string]PaymentReminderResponse)}
}

func (c *InMemoryClient) SendPaymentReminder(ctx context.Context, req PaymentReminderRequest) (PaymentReminderResponse, error) {
	if err := ctx.Err(); err != nil {
		return PaymentReminderResponse{}, err
	}
	if strings.TrimSpace(req.IdempotencyKey) == "" {
		return PaymentReminderResponse{}, errors.New("payment reminder idempotency key is required")
	}
	if strings.TrimSpace(req.StudentID) == "" {
		return PaymentReminderResponse{}, errors.New("payment reminder student ID is required")
	}
	if strings.TrimSpace(req.Item) == "" {
		return PaymentReminderResponse{}, errors.New("payment reminder item is required")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if existing, ok := c.responses[req.IdempotencyKey]; ok {
		return existing, nil
	}

	resp := PaymentReminderResponse{
		WorkflowID:     localWorkflowID(req.IdempotencyKey),
		Status:         "accepted",
		Message:        "Payment reminder workflow accepted by local deterministic P4 client.",
		IdempotencyKey: req.IdempotencyKey,
		Synthetic:      true,
	}
	c.responses[req.IdempotencyKey] = resp
	return resp, nil
}

func localWorkflowID(key string) string {
	hash := sha1.Sum([]byte(key))
	return fmt.Sprintf("LOCAL-WF-%s", strings.ToUpper(hex.EncodeToString(hash[:])[:12]))
}
