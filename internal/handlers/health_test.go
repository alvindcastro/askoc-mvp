package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"askoc-mvp/internal/middleware"
)

func TestHealthReturnsOKJSONWithTraceID(t *testing.T) {
	handler := middleware.TraceID(HealthHandler())
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set(middleware.TraceHeader, "trace-health")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var got struct {
		Status  string `json:"status"`
		TraceID string `json:"trace_id"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode health response: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got.Status != "ok" || got.TraceID != "trace-health" {
		t.Fatalf("health response = %+v", got)
	}
}

func TestReadyReturnsDependencyStatus(t *testing.T) {
	handler := ReadyHandler(DependencyCheck{
		Name: "database",
		Check: func(context.Context) error {
			return nil
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	req = req.WithContext(middleware.WithTraceID(req.Context(), "trace-ready"))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	var got struct {
		Status       string            `json:"status"`
		TraceID      string            `json:"trace_id"`
		Dependencies map[string]string `json:"dependencies"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode ready response: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got.Status != "ready" || got.Dependencies["database"] != "ok" || got.TraceID != "trace-ready" {
		t.Fatalf("ready response = %+v", got)
	}
}

func TestReadyReturnsUnavailableWhenDependencyFails(t *testing.T) {
	handler := ReadyHandler(DependencyCheck{
		Name: "workflow",
		Check: func(context.Context) error {
			return errors.New("connection refused with private host detail")
		},
	})
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
	if strings.Contains(rec.Body.String(), "connection refused") || strings.Contains(rec.Body.String(), "private host detail") {
		t.Fatalf("ready response leaked dependency details: %s", rec.Body.String())
	}
	var got struct {
		Status       string            `json:"status"`
		Dependencies map[string]string `json:"dependencies"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode ready response: %v", err)
	}
	if got.Status != "not_ready" || got.Dependencies["workflow"] != "unavailable" {
		t.Fatalf("ready response = %+v", got)
	}
}

func TestHealthAndReadyRejectInvalidMethods(t *testing.T) {
	tests := []struct {
		name    string
		handler http.Handler
		path    string
	}{
		{name: "health", handler: HealthHandler(), path: "/healthz"},
		{name: "ready", handler: ReadyHandler(), path: "/readyz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.path, nil)
			rec := httptest.NewRecorder()

			tt.handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusMethodNotAllowed {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
			}
		})
	}
}
