package middleware

import (
	"encoding/json"
	"net/http"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				writeMiddlewareError(w, r, http.StatusInternalServerError, "internal_error", "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func writeMiddlewareError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	traceID := ""
	if r != nil {
		traceID = TraceIDFromContext(r.Context())
	}

	body, err := json.Marshal(struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			TraceID string `json:"trace_id,omitempty"`
		} `json:"error"`
	}{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			TraceID string `json:"trace_id,omitempty"`
		}{
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
