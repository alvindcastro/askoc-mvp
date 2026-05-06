package handlers

import (
	"html/template"
	"net/http"
)

type chatPageData struct {
	ChatEndpoint string
}

func ChatPageHandler(templatePath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}

		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			WriteError(w, r, http.StatusInternalServerError, "template_unavailable", "chat page is unavailable")
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = tmpl.Execute(w, chatPageData{
			ChatEndpoint: "/api/v1/chat",
		})
	})
}

func StaticFileHandler(root string) http.Handler {
	return http.FileServer(http.Dir(root))
}
