package crm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"askoc-mvp/internal/handlers"
	"askoc-mvp/internal/middleware"
)

type Handler struct {
	mu     sync.Mutex
	nextID int
	cases  []CaseResponse
}

type CaseRequest struct {
	StudentID      string `json:"student_id"`
	ConversationID string `json:"conversation_id"`
	Intent         string `json:"intent"`
	Priority       string `json:"priority"`
	Queue          string `json:"queue"`
	Summary        string `json:"summary"`
	SourceTraceID  string `json:"source_trace_id"`
}

type CaseResponse struct {
	CaseID         string `json:"case_id"`
	Status         string `json:"status"`
	Queue          string `json:"queue"`
	Priority       string `json:"priority"`
	Summary        string `json:"summary"`
	ConversationID string `json:"conversation_id,omitempty"`
	SourceTraceID  string `json:"source_trace_id,omitempty"`
	Synthetic      bool   `json:"synthetic"`
}

func NewHandler() http.Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/v1/crm/cases" {
		handlers.WriteError(w, r, http.StatusNotFound, "not_found", "route not found")
		return
	}
	if r.Method != http.MethodPost {
		handlers.WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	var req CaseRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		handlers.WriteError(w, r, http.StatusBadRequest, "invalid_case_request", "case request must be valid JSON")
		return
	}
	if strings.TrimSpace(req.Summary) == "" {
		handlers.WriteError(w, r, http.StatusBadRequest, "invalid_case_summary", "case summary is required")
		return
	}
	if strings.TrimSpace(req.Queue) == "" {
		req.Queue = "learner_support"
	}
	if strings.TrimSpace(req.Priority) == "" {
		req.Priority = "normal"
	}

	resp := h.createCase(req)
	handlers.WriteJSON(w, r, http.StatusCreated, resp)
}

func (h *Handler) createCase(req CaseRequest) CaseResponse {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.nextID++
	resp := CaseResponse{
		CaseID:         caseID(h.nextID),
		Status:         "created",
		Queue:          strings.TrimSpace(req.Queue),
		Priority:       strings.TrimSpace(req.Priority),
		Summary:        middleware.BasicRedactor(strings.TrimSpace(req.Summary)),
		ConversationID: strings.TrimSpace(req.ConversationID),
		SourceTraceID:  strings.TrimSpace(req.SourceTraceID),
		Synthetic:      true,
	}
	h.cases = append(h.cases, resp)
	return resp
}

func caseID(sequence int) string {
	return fmt.Sprintf("MOCK-CRM-%06d", sequence)
}
