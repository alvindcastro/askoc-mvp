package audit

import (
	"context"
	"strings"
	"sync"
	"time"

	"askoc-mvp/internal/privacy"
)

type MemoryStore struct {
	mu     sync.RWMutex
	now    func() time.Time
	events []Event
}

type StoreOption func(*MemoryStore)

func WithClock(now func() time.Time) StoreOption {
	return func(s *MemoryStore) {
		if now != nil {
			s.now = now
		}
	}
}

func NewMemoryStore(opts ...StoreOption) *MemoryStore {
	store := &MemoryStore{
		now: time.Now,
	}
	for _, opt := range opts {
		opt(store)
	}
	return store
}

func (s *MemoryStore) Record(_ context.Context, event Event) error {
	prepared := sanitizeEvent(event, s.now())

	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, prepared)
	return nil
}

func (s *MemoryStore) List(context.Context) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneEvents(s.events, false)
}

func (s *MemoryStore) ListByTraceID(_ context.Context, traceID string) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []Event
	for _, event := range s.events {
		if event.TraceID == traceID {
			filtered = append(filtered, event)
		}
	}
	return cloneEvents(filtered, false)
}

func (s *MemoryStore) Export(context.Context) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneEvents(s.events, true)
}

func (s *MemoryStore) Reset(context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = nil
}

func (s *MemoryStore) PruneBefore(_ context.Context, cutoff time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	kept := s.events[:0]
	pruned := 0
	for _, event := range s.events {
		if event.RecordedAt.Before(cutoff) {
			pruned++
			continue
		}
		kept = append(kept, event)
	}
	s.events = kept
	return pruned
}

func sanitizeEvent(event Event, now time.Time) Event {
	if event.RecordedAt.IsZero() {
		event.RecordedAt = now
	}
	event.TraceID = strings.TrimSpace(event.TraceID)
	event.ConversationID = strings.TrimSpace(event.ConversationID)
	event.StudentID = redact(event.StudentID)
	event.Type = strings.TrimSpace(event.Type)
	event.Action = strings.TrimSpace(event.Action)
	event.Status = strings.TrimSpace(event.Status)
	event.ReferenceID = redact(event.ReferenceID)
	event.Message = redact(event.Message)
	event.Metadata = cloneMetadata(event.Metadata)
	for key, value := range event.Metadata {
		event.Metadata[key] = redact(value)
	}
	return event
}

func cloneEvents(events []Event, omitMessage bool) []Event {
	cloned := make([]Event, len(events))
	for i, event := range events {
		cloned[i] = event
		cloned[i].Metadata = cloneMetadata(event.Metadata)
		if omitMessage {
			cloned[i].Message = ""
		}
	}
	return cloned
}

func cloneMetadata(metadata map[string]string) map[string]string {
	if len(metadata) == 0 {
		return nil
	}
	cloned := make(map[string]string, len(metadata))
	for key, value := range metadata {
		cloned[key] = value
	}
	return cloned
}

func redact(value string) string {
	return privacy.Redact(value)
}
