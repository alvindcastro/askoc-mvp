package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
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

func TestChatPageRendersRevampShellAndCopy(t *testing.T) {
	handler := ChatPageHandler("../../web/templates/chat.html")
	req := httptest.NewRequest(http.MethodGet, "/chat", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	for _, stale := range []string{"P2 placeholder", "deterministic placeholder responses"} {
		if strings.Contains(body, stale) {
			t.Fatalf("chat page still contains stale copy %q: %s", stale, body)
		}
	}
	for _, want := range []string{
		`<nav class="app-nav" aria-label="AskOC routes">`,
		`href="/chat" aria-current="page"`,
		`href="/admin"`,
		"Transcript and payment support",
		"Source-grounded answer",
		"Action trace",
		"Safe synthetic data only",
		"Trace ID",
		"High risk",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("chat page missing revamp shell/copy %q: %s", want, body)
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

func TestChatStaticAssetsExposeThemeAndEvidenceContracts(t *testing.T) {
	css := readTestFile(t, "../../web/static/app.css")
	js := readTestFile(t, "../../web/static/app.js")

	for _, want := range []string{
		"--color-ink: #0A0A0A",
		"--color-paper: #FAFAFA",
		"--color-accent: #EF4444",
		"--radius-none: 0px",
		"--border-strong: 2px solid #0A0A0A",
		"--focus-ring: 0 0 0 2px #FAFAFA, 0 0 0 4px #0A0A0A",
		"font-family: var(--font-body)",
		".app-header h1",
		"color: var(--color-accent)",
		"border-radius: var(--radius-none)",
		"@media (max-width: 820px)",
	} {
		if !strings.Contains(css, want) {
			t.Fatalf("app.css missing theme contract %q", want)
		}
	}
	for _, forbidden := range []string{"linear-gradient", "radial-gradient", "border-radius: 8px", "border-radius: 999px"} {
		if strings.Contains(css, forbidden) {
			t.Fatalf("app.css contains forbidden visual pattern %q", forbidden)
		}
	}
	for _, want := range []string{
		"renderSourceEvidence",
		"renderActionEvidence",
		"confidence",
		"risk_level",
		"freshness_status",
		"trace_id",
		"workflow_id",
		"crm_case_id",
		"priority",
		"idempotency_key",
		"Low-confidence fallback",
		"No approved source",
		"Synthetic integration",
	} {
		if !strings.Contains(js, want) {
			t.Fatalf("app.js missing evidence rendering contract %q", want)
		}
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(body)
}
