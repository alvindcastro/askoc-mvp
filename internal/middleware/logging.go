package middleware

import (
	"log/slog"
	"net/http"
	"regexp"
	"strings"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func RequestLogger(logger *slog.Logger, redact func(string) string) func(http.Handler) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}
	if redact == nil {
		redact = func(value string) string { return value }
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)

			logger.Info(
				"http request",
				"method", r.Method,
				"path", redact(r.URL.RequestURI()),
				"status", rec.status,
				"trace_id", TraceIDFromContext(r.Context()),
			)
		})
	}
}

func BasicRedactor(value string) string {
	email := regexp.MustCompile(`[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}`)
	studentID := regexp.MustCompile(`S[0-9]{6}`)

	value = email.ReplaceAllString(value, "[redacted-email]")
	value = studentID.ReplaceAllString(value, "[redacted-student-id]")
	value = strings.ReplaceAll(value, "Authorization", "[redacted-header]")
	return value
}
