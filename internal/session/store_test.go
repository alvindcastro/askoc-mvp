package session

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestStoreCreatesAppendsAndReadsConversation(t *testing.T) {
	now := time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC)
	store := NewStore(30*time.Minute, WithClock(func() time.Time { return now }))

	if err := store.Create("conv_123", Message{Role: RoleUser, Content: "How do I order a transcript?"}); err != nil {
		t.Fatalf("Create error = %v", err)
	}
	if err := store.Append("conv_123", Message{Role: RoleAssistant, Content: "Use the transcript request form."}); err != nil {
		t.Fatalf("Append error = %v", err)
	}

	got, ok := store.Get("conv_123")
	if !ok {
		t.Fatal("Get ok = false, want true")
	}
	if len(got.Messages) != 2 {
		t.Fatalf("messages = %+v, want 2 messages", got.Messages)
	}
	if got.Messages[0].Role != RoleUser || got.Messages[1].Role != RoleAssistant {
		t.Fatalf("messages = %+v", got.Messages)
	}
	if got.ExpiresAt != now.Add(30*time.Minute) {
		t.Fatalf("expires_at = %s, want %s", got.ExpiresAt, now.Add(30*time.Minute))
	}
}

func TestStoreExpiresConversation(t *testing.T) {
	now := time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC)
	store := NewStore(10*time.Second, WithClock(func() time.Time { return now }))

	if err := store.Create("conv_expire", Message{Role: RoleUser, Content: "hello"}); err != nil {
		t.Fatalf("Create error = %v", err)
	}
	now = now.Add(11 * time.Second)

	if _, ok := store.Get("conv_expire"); ok {
		t.Fatal("Get ok = true after TTL, want false")
	}
}

func TestStoreRedactsPIIBeforePersistence(t *testing.T) {
	store := NewStore(30 * time.Minute)

	if err := store.Create("conv_redact", Message{
		Role:    RoleUser,
		Content: "Email me at learner@example.com, phone 250-555-1212, password is secret, student 12345678, synthetic S100002.",
	}); err != nil {
		t.Fatalf("Create error = %v", err)
	}

	got, ok := store.Get("conv_redact")
	if !ok {
		t.Fatal("Get ok = false, want true")
	}
	content := got.Messages[0].Content
	for _, leaked := range []string{"learner@example.com", "250-555-1212", "password is secret", "12345678"} {
		if strings.Contains(content, leaked) {
			t.Fatalf("stored message leaked %q: %s", leaked, content)
		}
	}
	if !strings.Contains(content, "S100002") {
		t.Fatalf("synthetic student ID was not preserved: %s", content)
	}
}

func TestStoreConcurrentAccess(t *testing.T) {
	store := NewStore(30 * time.Minute)
	if err := store.Create("conv_race"); err != nil {
		t.Fatalf("Create error = %v", err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 25; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := store.Append("conv_race", Message{Role: RoleUser, Content: fmt.Sprintf("message %d", i)}); err != nil {
				t.Errorf("Append error = %v", err)
			}
			if _, ok := store.Get("conv_race"); !ok {
				t.Error("Get ok = false, want true")
			}
		}(i)
	}
	wg.Wait()

	got, ok := store.Get("conv_race")
	if !ok {
		t.Fatal("Get ok = false, want true")
	}
	if len(got.Messages) != 25 {
		t.Fatalf("message count = %d, want 25", len(got.Messages))
	}
}
