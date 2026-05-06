package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const TraceHeader = "X-Trace-ID"

type traceKey struct{}

func TraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := strings.TrimSpace(r.Header.Get(TraceHeader))
		if traceID == "" {
			traceID = newTraceID()
		}

		w.Header().Set(TraceHeader, traceID)
		next.ServeHTTP(w, r.WithContext(WithTraceID(r.Context(), traceID)))
	})
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceKey{}, traceID)
}

func TraceIDFromContext(ctx context.Context) string {
	traceID, _ := ctx.Value(traceKey{}).(string)
	return traceID
}

func newTraceID() string {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err == nil {
		return hex.EncodeToString(buf[:])
	}
	return fmt.Sprintf("trace-%d", time.Now().UnixNano())
}
