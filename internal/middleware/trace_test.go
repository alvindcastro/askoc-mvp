package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTraceIDPreservesInboundTraceID(t *testing.T) {
	var contextTrace string
	handler := TraceID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextTrace = TraceIDFromContext(r.Context())
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(TraceHeader, "trace-inbound")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if contextTrace != "trace-inbound" {
		t.Fatalf("context trace = %q, want trace-inbound", contextTrace)
	}
	if rec.Header().Get(TraceHeader) != "trace-inbound" {
		t.Fatalf("response trace header = %q", rec.Header().Get(TraceHeader))
	}
}

func TestTraceIDGeneratesTraceIDWhenMissing(t *testing.T) {
	var contextTrace string
	handler := TraceID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextTrace = TraceIDFromContext(r.Context())
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if contextTrace == "" {
		t.Fatal("context trace ID was empty")
	}
	if rec.Header().Get(TraceHeader) != contextTrace {
		t.Fatalf("response trace header = %q, want %q", rec.Header().Get(TraceHeader), contextTrace)
	}
}
