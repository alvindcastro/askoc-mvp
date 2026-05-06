package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoverConvertsPanicToSafeInternalError(t *testing.T) {
	handler := TraceID(Recover(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("database password leaked")
	})))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(TraceHeader, "trace-recovery")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	body := rec.Body.String()
	if strings.Contains(body, "database password leaked") || strings.Contains(body, "panic") {
		t.Fatalf("response leaked panic details: %s", body)
	}
	if !strings.Contains(body, `"trace_id":"trace-recovery"`) {
		t.Fatalf("response missing trace ID: %s", body)
	}
}
