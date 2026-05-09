package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAdminPageRendersDashboardShell(t *testing.T) {
	handler := AdminPageHandler("../../web/templates/admin.html")
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{"AskOC Admin Dashboard", "Top intents", "Escalations", "Workflows", "Low-confidence review", "Stale-source warnings", "/api/v1/admin/metrics"} {
		if !strings.Contains(body, want) {
			t.Fatalf("admin page missing %q: %s", want, body)
		}
	}
}

func TestAdminPageRendersRevampDashboardShell(t *testing.T) {
	handler := AdminPageHandler("../../web/templates/admin.html")
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	for _, want := range []string{
		`<nav class="app-nav" aria-label="AskOC routes">`,
		`href="/chat"`,
		`href="/admin" aria-current="page"`,
		"Synthetic operations mode",
		"Redacted aggregate metrics",
		"Evaluation gate",
		"Review queue filter",
		"Admin token",
		"Export audit",
		"Purge expired",
		"Reset demo data",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("admin page missing revamp dashboard shell %q: %s", want, body)
		}
	}
}

func TestAdminPageRejectsUnsupportedMethod(t *testing.T) {
	handler := AdminPageHandler("../../web/templates/admin.html")
	req := httptest.NewRequest(http.MethodPost, "/admin", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestAdminStaticAssetsExposeThemeAndReviewContracts(t *testing.T) {
	css := readAdminTestFile(t, "../../web/static/admin.css")
	js := readAdminTestFile(t, "../../web/static/admin.js")

	for _, want := range []string{
		"--color-ink: #0A0A0A",
		"--color-paper: #FAFAFA",
		"--color-accent: #EF4444",
		"--radius-none: 0px",
		"--border-strong: 2px solid #0A0A0A",
		"--focus-ring: 0 0 0 2px #FAFAFA, 0 0 0 4px #0A0A0A",
		"border-radius: var(--radius-none)",
		"@media (max-width: 820px)",
	} {
		if !strings.Contains(css, want) {
			t.Fatalf("admin.css missing theme contract %q", want)
		}
	}
	for _, forbidden := range []string{"linear-gradient", "radial-gradient", "border-radius: 8px", "border-radius: 999px"} {
		if strings.Contains(css, forbidden) {
			t.Fatalf("admin.css contains forbidden visual pattern %q", forbidden)
		}
	}
	for _, want := range []string{
		"renderReviewItems",
		"trace_id",
		"queue",
		"priority",
		"status",
		"redacted",
		"reviewItems.length ? reviewItems",
		"Evaluation gate",
		"No review items",
	} {
		if !strings.Contains(js, want) {
			t.Fatalf("admin.js missing review rendering contract %q", want)
		}
	}
}

func readAdminTestFile(t *testing.T, path string) string {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(body)
}
