package handlers

import (
	"context"
	"net/http"

	"askoc-mvp/internal/middleware"
)

type DependencyCheck struct {
	Name  string
	Check func(context.Context) error
}

type healthResponse struct {
	Status  string `json:"status"`
	TraceID string `json:"trace_id,omitempty"`
}

type readyResponse struct {
	Status       string            `json:"status"`
	TraceID      string            `json:"trace_id,omitempty"`
	Dependencies map[string]string `json:"dependencies"`
}

func HealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}

		WriteJSON(w, r, http.StatusOK, healthResponse{
			Status:  "ok",
			TraceID: middleware.TraceIDFromContext(r.Context()),
		})
	})
}

func ReadyHandler(checks ...DependencyCheck) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}

		status := "ready"
		code := http.StatusOK
		dependencies := make(map[string]string, len(checks))
		for _, check := range checks {
			if check.Name == "" {
				continue
			}
			dependencies[check.Name] = "ok"
			if check.Check != nil {
				if err := check.Check(r.Context()); err != nil {
					dependencies[check.Name] = "unavailable"
					status = "not_ready"
					code = http.StatusServiceUnavailable
				}
			}
		}

		WriteJSON(w, r, code, readyResponse{
			Status:       status,
			TraceID:      middleware.TraceIDFromContext(r.Context()),
			Dependencies: dependencies,
		})
	})
}
