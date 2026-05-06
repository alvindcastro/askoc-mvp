package tools

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"askoc-mvp/internal/middleware"
)

func doJSON(ctx context.Context, client *http.Client, service, method, endpoint string, requestPayload any, responsePayload any) error {
	if err := ctx.Err(); err != nil {
		return &ToolError{Kind: KindTimeout, Service: service, Message: "request context ended before tool call", Err: err}
	}

	var body io.Reader
	if requestPayload != nil {
		encoded, err := json.Marshal(requestPayload)
		if err != nil {
			return &ToolError{Kind: KindParse, Service: service, Message: "unable to encode request", Err: err}
		}
		body = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return &ToolError{Kind: KindExternal, Service: service, Message: "unable to create request", Err: err}
	}
	req.Header.Set(middleware.TraceHeader, traceHeader(ctx))
	if requestPayload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() != nil || isTimeout(err) {
			return &ToolError{Kind: KindTimeout, Service: service, Message: "tool request timed out or was canceled", Err: err}
		}
		return &ToolError{Kind: KindExternal, Service: service, Message: "tool request failed", Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, resp.Body)
		return statusError(service, resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(responsePayload); err != nil {
		return &ToolError{Kind: KindParse, Service: service, StatusCode: resp.StatusCode, Message: "tool response was not valid JSON", Err: err}
	}
	return nil
}

func statusError(service string, status int) error {
	kind := KindExternal
	switch {
	case status == http.StatusNotFound:
		kind = KindNotFound
	case status >= 500:
		kind = KindRetryable
	}
	return &ToolError{
		Kind:       kind,
		Service:    service,
		StatusCode: status,
		Message:    fmt.Sprintf("tool returned HTTP %d", status),
	}
}

func isTimeout(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func traceHeader(ctx context.Context) string {
	if traceID := strings.TrimSpace(middleware.TraceIDFromContext(ctx)); traceID != "" {
		return traceID
	}
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err == nil {
		return "tool-" + hex.EncodeToString(buf[:])
	}
	return fmt.Sprintf("tool-%d", time.Now().UnixNano())
}
