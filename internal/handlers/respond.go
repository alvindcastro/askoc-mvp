package handlers

import (
	"encoding/json"
	"net/http"

	"askoc-mvp/internal/middleware"
)

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"trace_id,omitempty"`
}

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, "internal_error", "unable to encode response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	traceID := ""
	if r != nil {
		traceID = middleware.TraceIDFromContext(r.Context())
	}

	body, err := json.Marshal(errorEnvelope{
		Error: errorBody{
			Code:    code,
			Message: message,
			TraceID: traceID,
		},
	})
	if err != nil {
		body = []byte(`{"error":{"code":"internal_error","message":"internal server error"}}`)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
