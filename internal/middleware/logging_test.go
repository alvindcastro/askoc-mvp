package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestLoggerUsesRedactionHook(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{}))
	redact := func(value string) string {
		return strings.ReplaceAll(value, "learner@example.test", "[redacted-email]")
	}
	handler := RequestLogger(logger, redact)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	req := httptest.NewRequest(http.MethodGet, "/healthz?email=learner@example.test", nil)
	req = req.WithContext(WithTraceID(req.Context(), "trace-log"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	logs := buf.String()
	if strings.Contains(logs, "learner@example.test") {
		t.Fatalf("log was not redacted: %s", logs)
	}
	if !strings.Contains(logs, "[redacted-email]") {
		t.Fatalf("redaction hook was not used: %s", logs)
	}
	if !strings.Contains(logs, "trace-log") || !strings.Contains(logs, "202") {
		t.Fatalf("log missing trace or status: %s", logs)
	}
}
