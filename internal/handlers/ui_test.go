package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChatPageRendersDemoUI(t *testing.T) {
	handler := ChatPageHandler("../../web/templates/chat.html")
	req := httptest.NewRequest(http.MethodGet, "/chat", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); !strings.Contains(got, "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", got)
	}
	body := rec.Body.String()
	for _, want := range []string{"AskOC AI Concierge", "/api/v1/chat", "Synthetic demo mode"} {
		if !strings.Contains(body, want) {
			t.Fatalf("chat page missing %q: %s", want, body)
		}
	}
}

func TestChatPageRejectsUnsupportedMethod(t *testing.T) {
	handler := ChatPageHandler("../../web/templates/chat.html")
	req := httptest.NewRequest(http.MethodPost, "/chat", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestStaticFileRouteReturnsExpectedContentType(t *testing.T) {
	handler := http.StripPrefix("/static/", StaticFileHandler("../../web/static"))
	req := httptest.NewRequest(http.MethodGet, "/static/app.js", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); !strings.Contains(got, "javascript") {
		t.Fatalf("Content-Type = %q, want javascript", got)
	}
}
