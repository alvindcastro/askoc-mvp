package crm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerCreatesCaseWithRedactedSummary(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crm/cases", strings.NewReader(`{
		"student_id": "S100003",
		"conversation_id": "conv_123",
		"intent": "transcript_status",
		"queue": "registrar_student_accounts",
		"priority": "normal",
		"summary": "Learner S100003 emailed demo.learner@example.test about a transcript hold.",
		"source_trace_id": "trace-crm"
	}`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	var got CaseResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode case response: %v", err)
	}
	if !strings.HasPrefix(got.CaseID, "MOCK-CRM-") || got.Queue != "registrar_student_accounts" || got.Priority != "normal" {
		t.Fatalf("case response = %+v", got)
	}
	if strings.Contains(got.Summary, "S100003") || strings.Contains(got.Summary, "demo.learner@example.test") {
		t.Fatalf("summary was not redacted: %q", got.Summary)
	}
}

func TestHandlerCreatesPriorityCase(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crm/cases", strings.NewReader(`{
		"student_id": "S100002",
		"conversation_id": "conv_urgent",
		"intent": "human_handoff",
		"queue": "learner_support",
		"priority": "priority",
		"summary": "Learner needs staff follow-up for an urgent transcript question.",
		"source_trace_id": "trace-priority"
	}`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	var got CaseResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode case response: %v", err)
	}
	if got.Priority != "priority" || got.Status != "created" {
		t.Fatalf("case response = %+v", got)
	}
}

func TestHandlerRejectsEmptySummary(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/crm/cases", strings.NewReader(`{
		"student_id": "S100003",
		"conversation_id": "conv_empty",
		"intent": "transcript_status",
		"queue": "registrar_student_accounts",
		"priority": "normal",
		"summary": "   ",
		"source_trace_id": "trace-empty"
	}`))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	assertErrorCode(t, rec.Body.Bytes(), "invalid_case_summary")
}

func assertErrorCode(t *testing.T, body []byte, want string) {
	t.Helper()
	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("decode error response: %v; body=%s", err, body)
	}
	if got.Error.Code != want {
		t.Fatalf("error code = %q, want %q; body=%s", got.Error.Code, want, body)
	}
}
