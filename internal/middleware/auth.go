package middleware

import (
	"net/http"
	"strings"
)

func MockAuth(enabled bool, token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !enabled {
				next.ServeHTTP(w, r)
				return
			}

			if token == "" || strings.TrimSpace(r.Header.Get("Authorization")) != "Bearer "+token {
				writeMiddlewareError(w, r, http.StatusUnauthorized, "unauthorized", "missing or invalid bearer token")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
