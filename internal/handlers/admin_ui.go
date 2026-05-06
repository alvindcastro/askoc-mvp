package handlers

import (
	"html/template"
	"net/http"
)

type adminPageData struct {
	MetricsEndpoint string
	ExportEndpoint  string
	ResetEndpoint   string
	PurgeEndpoint   string
	ReviewEndpoint  string
}

func AdminPageHandler(templatePath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}

		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			WriteError(w, r, http.StatusInternalServerError, "template_unavailable", "admin page is unavailable")
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = tmpl.Execute(w, adminPageData{
			MetricsEndpoint: "/api/v1/admin/metrics",
			ExportEndpoint:  "/api/v1/admin/audit/export",
			ResetEndpoint:   "/api/v1/admin/audit/reset",
			PurgeEndpoint:   "/api/v1/admin/audit/purge",
			ReviewEndpoint:  "/api/v1/admin/review-items",
		})
	})
}
