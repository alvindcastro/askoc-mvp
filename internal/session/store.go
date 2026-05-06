package session

import (
	"errors"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

var ErrInvalidConversationID = errors.New("conversation ID must not be empty")

type Message struct {
	Role      Role      `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	ConversationID string    `json:"conversation_id"`
	Messages       []Message `json:"messages"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

type Store struct {
	mu       sync.RWMutex
	ttl      time.Duration
	now      func() time.Time
	redactor func(string) string
	sessions map[string]Session
}

type Option func(*Store)

func WithClock(now func() time.Time) Option {
	return func(s *Store) {
		if now != nil {
			s.now = now
		}
	}
}

func WithRedactor(redactor func(string) string) Option {
	return func(s *Store) {
		if redactor != nil {
			s.redactor = redactor
		}
	}
}

func NewStore(ttl time.Duration, opts ...Option) *Store {
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}
	store := &Store{
		ttl:      ttl,
		now:      time.Now,
		redactor: DefaultRedactor,
		sessions: make(map[string]Session),
	}
	for _, opt := range opts {
		opt(store)
	}
	return store
}

func (s *Store) Create(conversationID string, messages ...Message) error {
	conversationID = strings.TrimSpace(conversationID)
	if conversationID == "" {
		return ErrInvalidConversationID
	}

	now := s.now()
	session := Session{
		ConversationID: conversationID,
		Messages:       make([]Message, 0, len(messages)),
		CreatedAt:      now,
		UpdatedAt:      now,
		ExpiresAt:      now.Add(s.ttl),
	}
	for _, msg := range messages {
		session.Messages = append(session.Messages, s.prepareMessage(msg, now))
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[conversationID] = session
	return nil
}

func (s *Store) Append(conversationID string, message Message) error {
	conversationID = strings.TrimSpace(conversationID)
	if conversationID == "" {
		return ErrInvalidConversationID
	}

	now := s.now()
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[conversationID]
	if !ok || !now.Before(session.ExpiresAt) {
		session = Session{
			ConversationID: conversationID,
			CreatedAt:      now,
			Messages:       []Message{},
		}
	}
	session.Messages = append(session.Messages, s.prepareMessage(message, now))
	session.UpdatedAt = now
	session.ExpiresAt = now.Add(s.ttl)
	s.sessions[conversationID] = session
	return nil
}

func (s *Store) Get(conversationID string) (Session, bool) {
	conversationID = strings.TrimSpace(conversationID)
	if conversationID == "" {
		return Session{}, false
	}

	now := s.now()
	s.mu.RLock()
	session, ok := s.sessions[conversationID]
	s.mu.RUnlock()
	if !ok {
		return Session{}, false
	}
	if !now.Before(session.ExpiresAt) {
		s.mu.Lock()
		delete(s.sessions, conversationID)
		s.mu.Unlock()
		return Session{}, false
	}

	session.Messages = append([]Message(nil), session.Messages...)
	return session, true
}

func (s *Store) prepareMessage(message Message, now time.Time) Message {
	if message.CreatedAt.IsZero() {
		message.CreatedAt = now
	}
	message.Content = s.redactor(message.Content)
	return message
}

var (
	emailPattern  = regexp.MustCompile(`[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}`)
	phonePattern  = regexp.MustCompile(`\b(?:\+?1[-.\s]?)?(?:\(?[0-9]{3}\)?[-.\s]?)?[0-9]{3}[-.\s]?[0-9]{4}\b`)
	secretPattern = regexp.MustCompile(`(?i)\bpassword\s+is\s+\S+`)
	realIDPattern = regexp.MustCompile(`\b[0-9]{7,}\b`)
	spacePattern  = regexp.MustCompile(`[ \t]{2,}`)
)

func DefaultRedactor(value string) string {
	value = emailPattern.ReplaceAllString(value, "[REDACTED_EMAIL]")
	value = phonePattern.ReplaceAllString(value, "[REDACTED_PHONE]")
	value = secretPattern.ReplaceAllString(value, "[REDACTED_SECRET]")
	value = realIDPattern.ReplaceAllString(value, "[REDACTED_ID]")
	return strings.TrimSpace(spacePattern.ReplaceAllString(value, " "))
}
