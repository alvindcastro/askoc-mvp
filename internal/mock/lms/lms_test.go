package lms

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"askoc-mvp/internal/fixtures"
)

func TestHandlerReturnsKnownSyntheticStudentLMSAccess(t *testing.T) {
	handler := testHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/S100001/lms-access?course_id=DEMO-LMS-101", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var got AccessStatus
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode LMS access: %v", err)
	}
	if got.StudentID != "S100001" || got.CourseID != "DEMO-LMS-101" || got.AccessStatus != "available" || !got.Synthetic {
		t.Fatalf("LMS access = %+v", got)
	}
}

func TestHandlerReturnsSafeFallbackForUnknownCourse(t *testing.T) {
	handler := testHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/S100001/lms-access?course_id=DEMO-LMS-999", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var got AccessStatus
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode LMS fallback: %v", err)
	}
	if got.AccessStatus != "unknown_demo_course" || got.Message == "" {
		t.Fatalf("LMS fallback = %+v", got)
	}
}

func TestHandlerReturnsNotFoundForUnknownStudent(t *testing.T) {
	handler := testHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/S999999/lms-access?course_id=DEMO-LMS-101", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusNotFound, rec.Body.String())
	}
	assertErrorCode(t, rec.Body.Bytes(), "lms_student_not_found")
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
