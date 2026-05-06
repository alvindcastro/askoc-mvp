package workflow

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"askoc-mvp/internal/audit"
	"askoc-mvp/internal/handlers"
)

const PaymentReminderPath = "/api/v1/automation/payment-reminder"

type SimulatorHandler struct {
	sender PaymentReminderSender
	audit  audit.Recorder
}

func NewSimulatorHandler(sender PaymentReminderSender, recorder audit.Recorder) http.Handler {
	if sender == nil {
		sender = NewInMemoryClient()
	}
	if recorder == nil {
		recorder = audit.NopRecorder{}
	}
	return &SimulatorHandler{
		sender: sender,
		audit:  recorder,
	}
}

func (h *SimulatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != PaymentReminderPath {
		handlers.WriteError(w, r, http.StatusNotFound, "not_found", "route not found")
		return
	}
	if r.Method != http.MethodPost {
		handlers.WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	var req PaymentReminderRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		handlers.WriteError(w, r, http.StatusBadRequest, "invalid_workflow_request", "payment reminder request must be valid JSON")
		return
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		handlers.WriteError(w, r, http.StatusBadRequest, "invalid_workflow_request", "payment reminder request must contain one JSON object")
		return
	}

	resp, err := h.sender.SendPaymentReminder(r.Context(), req)
	if err != nil {
		handlers.WriteError(w, r, http.StatusBadRequest, "invalid_workflow_request", err.Error())
		return
	}

	if err := h.audit.Record(r.Context(), audit.Event{
		TraceID:        req.TraceID,
		ConversationID: req.ConversationID,
		StudentID:      req.StudentID,
		Type:           audit.EventTypeWorkflow,
		Action:         audit.ActionPaymentReminder,
		Status:         audit.StatusCompleted,
		ReferenceID:    resp.WorkflowID,
		Message:        "payment reminder workflow accepted",
		Metadata: map[string]string{
			"idempotency_key_hash": IdempotencyKeyHash(req.IdempotencyKey),
			"item":                 req.Item,
		},
	}); err != nil {
		handlers.WriteError(w, r, http.StatusInternalServerError, "workflow_audit_failed", "unable to record workflow event")
		return
	}

	handlers.WriteJSON(w, r, http.StatusAccepted, resp)
}
