package privacy_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"askoc-mvp/internal/middleware"
	"askoc-mvp/internal/privacy"
)

func TestRequestLoggerRedactsPathAndDoesNotLogChatBody(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{}))
	handler := middleware.RequestLogger(logger, privacy.Redact)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/chat?email=learner@example.test&token=abc123", strings.NewReader(`{"message":"call 250-555-0199"}`))
	req = req.WithContext(middleware.WithTraceID(req.Context(), "trace-privacy-log"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	logs := buf.String()
	for _, leaked := range []string{"learner@example.test", "abc123", "250-555-0199"} {
		if strings.Contains(logs, leaked) {
			t.Fatalf("request log leaked %q: %s", leaked, logs)
		}
	}
	for _, want := range []string{"trace-privacy-log", "202", "[REDACTED_EMAIL]", "[REDACTED_SECRET]"} {
		if !strings.Contains(logs, want) {
			t.Fatalf("request log missing %q: %s", want, logs)
		}
	}
}
