package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"askoc-mvp/internal/middleware"
)

func TestWriteJSONWritesHeadersStatusAndBody(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/example", nil)

	WriteJSON(rec, req, http.StatusCreated, map[string]string{"status": "ok"})

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	if strings.TrimSpace(rec.Body.String()) != `{"status":"ok"}` {
		t.Fatalf("body = %q", rec.Body.String())
	}
}

func TestWriteErrorUsesStableSafeShape(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	req = req.WithContext(middleware.WithTraceID(req.Context(), "trace-123"))

	WriteError(rec, req, http.StatusBadRequest, "bad_request", "check the request")

	var got struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			TraceID string `json:"trace_id"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if got.Error.Code != "bad_request" || got.Error.Message != "check the request" || got.Error.TraceID != "trace-123" {
		t.Fatalf("error response = %+v", got.Error)
	}
}

func TestWriteJSONHandlesUnsupportedValuesSafely(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/example", nil)
	req = req.WithContext(middleware.WithTraceID(req.Context(), "trace-456"))

	WriteJSON(rec, req, http.StatusOK, map[string]any{"bad": func() {}})

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	if strings.Contains(rec.Body.String(), "func") || strings.Contains(rec.Body.String(), "unsupported type") {
		t.Fatalf("response leaked raw encoder details: %s", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"trace_id":"trace-456"`) {
		t.Fatalf("response did not include trace ID: %s", rec.Body.String())
	}
}
