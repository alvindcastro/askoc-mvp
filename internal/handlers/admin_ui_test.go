package handlers

import (
	"net/http"
	"net/http/httptest"
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

func TestAdminPageRejectsUnsupportedMethod(t *testing.T) {
	handler := AdminPageHandler("../../web/templates/admin.html")
	req := httptest.NewRequest(http.MethodPost, "/admin", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
