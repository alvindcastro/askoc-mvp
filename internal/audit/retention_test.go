package audit

import (
	"context"
	"testing"
	"time"
)

func TestDefaultRetentionPolicyIsShortForDemo(t *testing.T) {
	policy := DefaultRetentionPolicy()

	if policy.MaxAge <= 0 || policy.MaxAge > 30*24*time.Hour {
		t.Fatalf("DefaultRetentionPolicy MaxAge = %s, want short demo retention", policy.MaxAge)
	}
}

func TestRetentionPolicyPurgesExpiredEvents(t *testing.T) {
	now := time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC)
	store := NewMemoryStore(WithClock(func() time.Time { return now }))
	if err := store.Record(context.Background(), Event{TraceID: "trace-old", Type: EventTypeIntent}); err != nil {
		t.Fatalf("Record old event: %v", err)
	}
	now = now.Add(8 * 24 * time.Hour)
	if err := store.Record(context.Background(), Event{TraceID: "trace-new", Type: EventTypeIntent}); err != nil {
		t.Fatalf("Record new event: %v", err)
	}

	pruned := DefaultRetentionPolicy().PurgeExpired(context.Background(), store, now)

	if pruned != 1 {
		t.Fatalf("PurgeExpired removed = %d, want 1", pruned)
	}
	got := store.List(context.Background())
	if len(got) != 1 || got[0].TraceID != "trace-new" {
		t.Fatalf("remaining events = %+v, want only trace-new", got)
	}
}
