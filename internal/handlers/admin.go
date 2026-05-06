package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"askoc-mvp/internal/audit"
)

type auditEventStore interface {
	List(context.Context) []audit.Event
	Export(context.Context) []audit.Event
	Reset(context.Context)
}

type auditListStore interface {
	List(context.Context) []audit.Event
}

type auditExportStore interface {
	Export(context.Context) []audit.Event
}

type auditResetStore interface {
	Reset(context.Context)
}

type auditExportResponse struct {
	Events []audit.Event `json:"events"`
}

type auditResetResponse struct {
	Status string `json:"status"`
}

type auditPurgeResponse struct {
	Pruned int `json:"pruned"`
}

func AdminMetricsHandler(store auditListStore, adminToken string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}
		if !authorizedAdmin(r, adminToken) {
			WriteError(w, r, http.StatusUnauthorized, "unauthorized", "missing or invalid admin bearer token")
			return
		}
		if store == nil {
			WriteJSON(w, r, http.StatusOK, audit.SummaryFromEvents(nil))
			return
		}
		WriteJSON(w, r, http.StatusOK, audit.SummaryFromEvents(store.List(r.Context())))
	})
}

func AdminAuditExportHandler(store auditExportStore, adminToken string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}
		if !authorizedAdmin(r, adminToken) {
			WriteError(w, r, http.StatusUnauthorized, "unauthorized", "missing or invalid admin bearer token")
			return
		}
		var events []audit.Event
		if store != nil {
			events = store.Export(r.Context())
		}
		WriteJSON(w, r, http.StatusOK, auditExportResponse{Events: events})
	})
}

func AdminAuditResetHandler(store auditResetStore, adminToken string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}
		if !authorizedAdmin(r, adminToken) {
			WriteError(w, r, http.StatusUnauthorized, "unauthorized", "missing or invalid admin bearer token")
			return
		}
		if store != nil {
			store.Reset(r.Context())
		}
		WriteJSON(w, r, http.StatusOK, auditResetResponse{Status: "reset"})
	})
}

func AdminAuditPurgeHandler(store audit.Pruner, adminToken string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			WriteError(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}
		if !authorizedAdmin(r, adminToken) {
			WriteError(w, r, http.StatusUnauthorized, "unauthorized", "missing or invalid admin bearer token")
			return
		}
		pruned := audit.DefaultRetentionPolicy().PurgeExpired(r.Context(), store, time.Now().UTC())
		WriteJSON(w, r, http.StatusOK, auditPurgeResponse{Pruned: pruned})
	})
}

func authorizedAdmin(r *http.Request, token string) bool {
	token = strings.TrimSpace(token)
	if token == "" {
		return false
	}
	return strings.TrimSpace(r.Header.Get("Authorization")) == "Bearer "+token
}
