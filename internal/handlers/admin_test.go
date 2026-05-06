package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"askoc-mvp/internal/audit"
	"askoc-mvp/internal/domain"
)

func TestAdminMetricsRequiresBearerToken(t *testing.T) {
	store := audit.NewMemoryStore()
	handler := AdminMetricsHandler(store, "demo-admin-token")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/metrics", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAdminMetricsReturnsZerosForEmptyStore(t *testing.T) {
	store := audit.NewMemoryStore()
	handler := AdminMetricsHandler(store, "demo-admin-token")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/metrics", nil)
	req.Header.Set("Authorization", "Bearer demo-admin-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var got audit.Summary
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.TotalConversations != 0 || got.ContainmentRate != 0 || got.EscalationRate != 0 {
		t.Fatalf("empty summary = %+v, want zero values", got)
	}
	if len(got.TopIntents) != 0 {
		t.Fatalf("TopIntents = %+v, want empty", got.TopIntents)
	}
}

func TestAdminMetricsSummarizesSeededAuditEventsAndRedactsReviewQueue(t *testing.T) {
	now := time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC)
	store := audit.NewMemoryStore(audit.WithClock(func() time.Time { return now }))
	seedEvents(t, store, []audit.Event{
		{
			TraceID:        "trace-contained",
			ConversationID: "conv-contained",
			Type:           audit.EventTypeIntent,
			Action:         audit.ActionClassify,
			Status:         audit.StatusCompleted,
			Message:        "How do I order a transcript?",
			Metadata:       map[string]string{"intent": "transcript_request", "confidence": "0.92"},
		},
		{
			TraceID:        "trace-low",
			ConversationID: "conv-low",
			Type:           audit.EventTypeIntent,
			Action:         audit.ActionClassify,
			Status:         audit.StatusCompleted,
			Message:        "Email learner@example.test and check 12345678",
			Metadata:       map[string]string{"intent": "unknown", "confidence": "0.31", "low_confidence": "true"},
		},
		{
			TraceID:        "trace-low",
			ConversationID: "conv-low",
			Type:           audit.EventTypeEscalation,
			Action:         audit.ActionCreateCRMCase,
			Status:         audit.StatusCompleted,
		},
		{
			TraceID:        "trace-workflow",
			ConversationID: "conv-workflow",
			Type:           audit.EventTypeWorkflow,
			Action:         audit.ActionPaymentReminder,
			Status:         audit.StatusCompleted,
			ReferenceID:    "WF-1",
		},
		{
			TraceID:        "trace-workflow-fail",
			ConversationID: "conv-workflow-fail",
			Type:           audit.EventTypeWorkflow,
			Action:         audit.ActionPaymentReminder,
			Status:         audit.StatusFailed,
		},
		{
			TraceID:        "trace-stale",
			ConversationID: "conv-stale",
			Type:           audit.EventTypeGuardrail,
			Action:         audit.ActionSourceCheck,
			Status:         audit.StatusBlocked,
			Message:        "Can I rely on this stale source?",
			Metadata:       map[string]string{"stale_source": "true"},
		},
	})

	handler := AdminMetricsHandler(store, "demo-admin-token")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/metrics", nil)
	req.Header.Set("Authorization", "Bearer demo-admin-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var got audit.Summary
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.TotalConversations != 5 {
		t.Fatalf("TotalConversations = %d, want 5", got.TotalConversations)
	}
	if got.Escalations != 1 || got.Automation.WorkflowFailures != 1 || got.Automation.PaymentRemindersSent != 1 {
		t.Fatalf("summary = %+v, want escalation and workflow counts", got)
	}
	if len(got.TopIntents) == 0 || got.TopIntents[0].Intent != "transcript_request" {
		t.Fatalf("TopIntents = %+v, want transcript_request first", got.TopIntents)
	}
	if got.ReviewQueue.LowConfidenceAnswers != 1 || got.ReviewQueue.StaleSourceQuestions != 1 {
		t.Fatalf("ReviewQueue = %+v, want low-confidence and stale-source counts", got.ReviewQueue)
	}
	rendered := rec.Body.String()
	for _, leaked := range []string{"learner@example.test", "12345678"} {
		if strings.Contains(rendered, leaked) {
			t.Fatalf("metrics response leaked %q: %s", leaked, rendered)
		}
	}
	if !strings.Contains(rendered, "[REDACTED_EMAIL]") || !strings.Contains(rendered, "[REDACTED_ID]") {
		t.Fatalf("metrics response missing redaction markers: %s", rendered)
	}
}

func TestAdminAuditExportAndResetRequireAuthAndOmitMessages(t *testing.T) {
	store := audit.NewMemoryStore()
	seedEvents(t, store, []audit.Event{{
		TraceID:        "trace-export",
		ConversationID: "conv-export",
		Type:           audit.EventTypeIntent,
		Action:         audit.ActionClassify,
		Status:         audit.StatusCompleted,
		Message:        "Email learner@example.test",
	}})

	exportHandler := AdminAuditExportHandler(store, "demo-admin-token")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/audit/export", nil)
	req.Header.Set("Authorization", "Bearer demo-admin-token")
	rec := httptest.NewRecorder()

	exportHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("export status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "learner@example.test") || strings.Contains(rec.Body.String(), "Email") {
		t.Fatalf("export leaked message content: %s", rec.Body.String())
	}

	resetHandler := AdminAuditResetHandler(store, "demo-admin-token")
	resetReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/audit/reset", nil)
	resetReq.Header.Set("Authorization", "Bearer demo-admin-token")
	resetRec := httptest.NewRecorder()

	resetHandler.ServeHTTP(resetRec, resetReq)

	if resetRec.Code != http.StatusOK {
		t.Fatalf("reset status = %d, want %d; body=%s", resetRec.Code, http.StatusOK, resetRec.Body.String())
	}
	if got := store.List(context.Background()); len(got) != 0 {
		t.Fatalf("store count after reset = %d, want 0", len(got))
	}
}

func TestAdminAuditPurgeAppliesDefaultRetention(t *testing.T) {
	store := audit.NewMemoryStore()
	now := time.Now().UTC()
	seedEvents(t, store, []audit.Event{
		{TraceID: "trace-old", Type: audit.EventTypeIntent, RecordedAt: now.Add(-8 * 24 * time.Hour)},
		{TraceID: "trace-new", Type: audit.EventTypeIntent, RecordedAt: now},
	})
	handler := AdminAuditPurgeHandler(store, "demo-admin-token")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/audit/purge", nil)
	req.Header.Set("Authorization", "Bearer demo-admin-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"pruned":1`) {
		t.Fatalf("purge response = %s, want pruned count", rec.Body.String())
	}
	got := store.List(context.Background())
	if len(got) != 1 || got[0].TraceID != "trace-new" {
		t.Fatalf("remaining events = %+v, want only trace-new", got)
	}
}

func TestAdminReviewQueueReturnsOpenRedactedItems(t *testing.T) {
	handler := AdminReviewQueueHandler(fakeReviewQueue{}, "demo-admin-token")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/review-items", nil)
	req.Header.Set("Authorization", "Bearer demo-admin-token")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, "classification_guardrail") || !strings.Contains(body, "oc-registrar-office") {
		t.Fatalf("review queue body missing sources/actions: %s", body)
	}
	if strings.Contains(body, "learner@example.test") || strings.Contains(body, "12345678") {
		t.Fatalf("review queue leaked PII: %s", body)
	}
}

type fakeReviewQueue struct{}

func (fakeReviewQueue) OpenReviewItems(context.Context) any {
	return []struct {
		Reason   string          `json:"reason"`
		Question string          `json:"question"`
		Sources  []domain.Source `json:"sources"`
		Actions  []domain.Action `json:"actions"`
	}{
		{
			Reason:   "low_confidence",
			Question: "Email [REDACTED_EMAIL] and check [REDACTED_ID]",
			Sources:  []domain.Source{{ID: "oc-registrar-office", Title: "Office of the Registrar"}},
			Actions:  []domain.Action{{Type: "classification_guardrail", Status: domain.ActionStatusPending}},
		},
	}
}

func seedEvents(t *testing.T, store *audit.MemoryStore, events []audit.Event) {
	t.Helper()
	for _, event := range events {
		if err := store.Record(context.Background(), event); err != nil {
			t.Fatalf("Record error = %v", err)
		}
	}
}
