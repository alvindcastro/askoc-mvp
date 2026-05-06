package audit

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMemoryStoreRecordsRedactedEventsAndListsDeterministically(t *testing.T) {
	now := time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC)
	store := NewMemoryStore(WithClock(func() time.Time { return now }))

	events := []Event{
		{
			TraceID:        "trace-002",
			ConversationID: "conv-1",
			StudentID:      "12345678",
			Type:           EventTypeTool,
			Action:         ActionPaymentLookup,
			Status:         StatusCompleted,
			Message:        "Email learner@example.com or 250-555-1212, password is secret.",
			Metadata: map[string]string{
				"raw_email":  "learner@example.com",
				"student_id": "12345678",
				"safe_id":    "S100002",
			},
		},
		{
			TraceID: "trace-001",
			Type:    EventTypeIntent,
			Action:  ActionClassify,
			Status:  StatusCompleted,
		},
	}
	for _, event := range events {
		if err := store.Record(context.Background(), event); err != nil {
			t.Fatalf("Record error = %v", err)
		}
		now = now.Add(time.Second)
	}

	got := store.List(context.Background())
	if len(got) != 2 {
		t.Fatalf("List count = %d, want 2", len(got))
	}
	if got[0].TraceID != "trace-002" || got[1].TraceID != "trace-001" {
		t.Fatalf("List order = %+v, want insertion order", got)
	}
	if !got[0].RecordedAt.Equal(time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC)) {
		t.Fatalf("RecordedAt = %s, want injected clock", got[0].RecordedAt)
	}
	for _, leaked := range []string{"learner@example.com", "250-555-1212", "password is secret", "12345678"} {
		if strings.Contains(fmt.Sprintf("%+v", got[0]), leaked) {
			t.Fatalf("stored event leaked %q: %+v", leaked, got[0])
		}
	}
	if got[0].Metadata["safe_id"] != "S100002" {
		t.Fatalf("synthetic ID metadata = %q, want preserved", got[0].Metadata["safe_id"])
	}

	got[0].Metadata["safe_id"] = "mutated"
	again := store.List(context.Background())
	if again[0].Metadata["safe_id"] != "S100002" {
		t.Fatalf("List exposed mutable metadata: %+v", again[0].Metadata)
	}
}

func TestMemoryStoreSupportsEmptyResetTraceQueriesExportAndPrune(t *testing.T) {
	now := time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC)
	store := NewMemoryStore(WithClock(func() time.Time { return now }))

	if got := store.List(context.Background()); len(got) != 0 {
		t.Fatalf("empty List count = %d, want 0", len(got))
	}
	if got := store.Export(context.Background()); len(got) != 0 {
		t.Fatalf("empty Export count = %d, want 0", len(got))
	}

	seed := []Event{
		{TraceID: "trace-old", Type: EventTypeWorkflow, Action: ActionPaymentReminder, Status: StatusCompleted, Message: "old"},
		{TraceID: "trace-keep", Type: EventTypeWorkflow, Action: ActionPaymentReminder, Status: StatusFailed, Message: "new"},
		{TraceID: "trace-keep", Type: EventTypeEscalation, Action: ActionCreateCRMCase, Status: StatusCompleted, Message: "Email learner@example.com"},
	}
	for _, event := range seed {
		if err := store.Record(context.Background(), event); err != nil {
			t.Fatalf("Record error = %v", err)
		}
		now = now.Add(time.Hour)
	}

	if got := store.ListByTraceID(context.Background(), "trace-keep"); len(got) != 2 {
		t.Fatalf("ListByTraceID count = %d, want 2", len(got))
	}
	exported := store.Export(context.Background())
	if len(exported) != 3 {
		t.Fatalf("Export count = %d, want 3", len(exported))
	}
	if exported[2].Message != "" {
		t.Fatalf("Export message = %q, want raw/sanitized message omitted", exported[2].Message)
	}

	pruned := store.PruneBefore(context.Background(), time.Date(2026, 5, 6, 13, 0, 0, 0, time.UTC))
	if pruned != 1 {
		t.Fatalf("PruneBefore removed = %d, want 1", pruned)
	}
	if got := store.List(context.Background()); len(got) != 2 {
		t.Fatalf("List after prune count = %d, want 2", len(got))
	}

	store.Reset(context.Background())
	if got := store.List(context.Background()); len(got) != 0 {
		t.Fatalf("List after reset count = %d, want 0", len(got))
	}
}

func TestMemoryStoreConcurrentRecordsAreSafe(t *testing.T) {
	store := NewMemoryStore()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := store.Record(context.Background(), Event{
				TraceID: fmt.Sprintf("trace-%02d", i),
				Type:    EventTypeTool,
				Action:  ActionBannerLookup,
				Status:  StatusCompleted,
			}); err != nil {
				t.Errorf("Record error = %v", err)
			}
		}(i)
	}
	wg.Wait()

	if got := store.List(context.Background()); len(got) != 50 {
		t.Fatalf("List count = %d, want 50", len(got))
	}
}

func TestDashboardMetricsFromSeededEvents(t *testing.T) {
	events := []Event{
		{Type: EventTypeIntent, Action: ActionClassify, Status: StatusCompleted},
		{Type: EventTypeIntent, Action: ActionClassify, Status: StatusCompleted, Metadata: map[string]string{"confidence": "0.41"}},
		{Type: EventTypeGuardrail, Action: ActionSourceCheck, Status: StatusBlocked},
		{Type: EventTypeWorkflow, Action: ActionPaymentReminder, Status: StatusAttempted},
		{Type: EventTypeWorkflow, Action: ActionPaymentReminder, Status: StatusCompleted},
		{Type: EventTypeWorkflow, Action: ActionPaymentReminder, Status: StatusFailed},
		{Type: EventTypeEscalation, Action: ActionCreateCRMCase, Status: StatusCompleted},
	}

	got := DashboardMetricsFromEvents(events)
	wantByType := map[string]int{
		EventTypeIntent:     2,
		EventTypeGuardrail:  1,
		EventTypeWorkflow:   3,
		EventTypeEscalation: 1,
	}
	if !reflect.DeepEqual(got.ByType, wantByType) {
		t.Fatalf("ByType = %+v, want %+v", got.ByType, wantByType)
	}
	if got.ByAction[ActionPaymentReminder] != 3 {
		t.Fatalf("payment reminder count = %d, want 3", got.ByAction[ActionPaymentReminder])
	}
	if got.ByStatus[StatusFailed] != 1 || got.ByStatus[StatusBlocked] != 1 {
		t.Fatalf("ByStatus = %+v, want failed and blocked counts", got.ByStatus)
	}
	if got.Escalations != 1 {
		t.Fatalf("Escalations = %d, want 1", got.Escalations)
	}
	if got.LowConfidence != 1 {
		t.Fatalf("LowConfidence = %d, want 1", got.LowConfidence)
	}
	if got.GuardrailEvents != 1 {
		t.Fatalf("GuardrailEvents = %d, want 1", got.GuardrailEvents)
	}
	if got.WorkflowEvents != 3 {
		t.Fatalf("WorkflowEvents = %d, want 3", got.WorkflowEvents)
	}

	empty := DashboardMetricsFromEvents(nil)
	if len(empty.ByType) != 0 || empty.Escalations != 0 || empty.LowConfidence != 0 {
		t.Fatalf("empty metrics = %+v, want zero values", empty)
	}
}
