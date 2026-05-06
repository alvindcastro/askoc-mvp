package banner

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"askoc-mvp/internal/fixtures"
)

func TestHandlerReturnsKnownSyntheticStudent(t *testing.T) {
	handler := testHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/S100002", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var got StudentProfile
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode profile: %v", err)
	}
	if got.StudentID != "S100002" || got.PreferredName != "Demo Learner Two" || !got.Synthetic {
		t.Fatalf("profile = %+v", got)
	}
}

func TestHandlerReturnsNotFoundForUnknownStudent(t *testing.T) {
	handler := testHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/S999999", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusNotFound, rec.Body.String())
	}
	assertErrorCode(t, rec.Body.Bytes(), "student_not_found")
}

func TestHandlerReturnsTranscriptStatusAndFinancialHold(t *testing.T) {
	handler := testHandler(t)

	tests := []struct {
		id       string
		status   string
		hold     string
		eligible bool
	}{
		{id: "S100001", status: "ready_for_processing", hold: "none", eligible: true},
		{id: "S100003", status: "needs_staff_review", hold: "financial", eligible: false},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/students/"+tt.id+"/transcript-status", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
			}
			var got TranscriptStatus
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("decode transcript status: %v", err)
			}
			if got.StudentID != tt.id || got.TranscriptRequestStatus != tt.status || got.Hold != tt.hold || got.EligibleForProcessing != tt.eligible {
				t.Fatalf("transcript status = %+v, want status=%q hold=%q eligible=%v", got, tt.status, tt.hold, tt.eligible)
			}
		})
	}
}

func TestTranscriptStatusContractIncludesStableFields(t *testing.T) {
	handler := testHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/S100002/transcript-status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode contract response: %v", err)
	}
	for _, field := range []string{"student_id", "transcript_request_id", "transcript_request_status", "eligible_for_processing", "hold", "holds", "synthetic"} {
		if _, ok := got[field]; !ok {
			t.Fatalf("field %q missing from response %v", field, got)
		}
	}
}

func testHandler(t *testing.T) http.Handler {
	t.Helper()
	records, err := fixtures.Load(context.Background(), "../../../data/synthetic-students.json")
	if err != nil {
		t.Fatalf("load fixtures: %v", err)
	}
	return NewHandler(records)
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
