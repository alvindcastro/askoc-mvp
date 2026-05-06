package eval

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"askoc-mvp/internal/domain"
	"askoc-mvp/internal/privacy"
)

type ReviewQueue struct {
	mu    sync.RWMutex
	now   func() time.Time
	items map[string]ReviewItem
}

type ReviewItem struct {
	ID              string          `json:"id"`
	CaseID          string          `json:"case_id,omitempty"`
	Reason          string          `json:"reason"`
	Question        string          `json:"question"`
	Sources         []domain.Source `json:"sources,omitempty"`
	Actions         []domain.Action `json:"actions,omitempty"`
	Critical        bool            `json:"critical"`
	Failures        []string        `json:"failures,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	ResolvedAt      *time.Time      `json:"resolved_at,omitempty"`
	OccurrenceCount int             `json:"occurrence_count"`
}

type ReviewQueueOption func(*ReviewQueue)

func WithReviewQueueClock(now func() time.Time) ReviewQueueOption {
	return func(q *ReviewQueue) {
		if now != nil {
			q.now = now
		}
	}
}

func NewReviewQueue(opts ...ReviewQueueOption) *ReviewQueue {
	q := &ReviewQueue{
		now:   time.Now,
		items: map[string]ReviewItem{},
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

func (q *ReviewQueue) AddFromResult(ctx context.Context, result CaseResult) bool {
	if err := ctx.Err(); err != nil {
		return false
	}
	reason, critical := reviewReason(result)
	if reason == "" {
		return false
	}
	normalized := normalizeQuestion(result.Prompt)
	if normalized == "" {
		normalized = strings.ToLower(strings.TrimSpace(result.ID))
	}
	key := reason + ":" + normalized
	item := ReviewItem{
		ID:              reviewItemID(key),
		CaseID:          result.ID,
		Reason:          reason,
		Question:        privacy.Redact(result.Prompt),
		Sources:         sanitizeSources(result.Response.Sources),
		Actions:         sanitizeActions(result.Response.Actions),
		Critical:        critical,
		Failures:        append([]string{}, result.Score.CriticalFailures...),
		CreatedAt:       q.now().UTC(),
		OccurrenceCount: 1,
	}
	if len(item.Failures) == 0 {
		item.Failures = append(item.Failures, result.Score.MinorFailures...)
	}

	q.mu.Lock()
	defer q.mu.Unlock()
	if existing, ok := q.items[key]; ok {
		existing.OccurrenceCount++
		q.items[key] = existing
		return false
	}
	q.items[key] = item
	return true
}

func (q *ReviewQueue) Open(ctx context.Context) []ReviewItem {
	if err := ctx.Err(); err != nil {
		return nil
	}
	q.mu.RLock()
	defer q.mu.RUnlock()
	items := make([]ReviewItem, 0, len(q.items))
	for _, item := range q.items {
		if item.ResolvedAt == nil {
			items = append(items, cloneReviewItem(item))
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Critical == items[j].Critical {
			return items[i].CreatedAt.Before(items[j].CreatedAt)
		}
		return items[i].Critical && !items[j].Critical
	})
	return items
}

func (q *ReviewQueue) OpenReviewItems(ctx context.Context) any {
	return q.Open(ctx)
}

func (q *ReviewQueue) Resolve(ctx context.Context, id string) bool {
	if err := ctx.Err(); err != nil {
		return false
	}
	id = strings.TrimSpace(id)
	q.mu.Lock()
	defer q.mu.Unlock()
	for key, item := range q.items {
		if item.ID != id || item.ResolvedAt != nil {
			continue
		}
		resolvedAt := q.now().UTC()
		item.ResolvedAt = &resolvedAt
		q.items[key] = item
		return true
	}
	return false
}

func reviewReason(result CaseResult) (string, bool) {
	if len(result.Score.CriticalFailures) > 0 {
		return "critical_eval_failure", true
	}
	if result.Response.Intent.Confidence > 0 && result.Response.Intent.Confidence < 0.5 {
		return "low_confidence", false
	}
	if len(result.Score.MinorFailures) > 0 {
		return "eval_miss", false
	}
	return "", false
}

func sanitizeSources(sources []domain.Source) []domain.Source {
	out := make([]domain.Source, len(sources))
	for i, source := range sources {
		out[i] = source
		out[i].Title = privacy.Redact(out[i].Title)
		out[i].URL = privacy.Redact(out[i].URL)
		out[i].Caution = privacy.Redact(out[i].Caution)
	}
	return out
}

func sanitizeActions(actions []domain.Action) []domain.Action {
	out := make([]domain.Action, len(actions))
	for i, action := range actions {
		out[i] = action
		out[i].Message = privacy.Redact(out[i].Message)
		out[i].ReferenceID = privacy.Redact(out[i].ReferenceID)
		out[i].IdempotencyKey = privacy.Redact(out[i].IdempotencyKey)
	}
	return out
}

func cloneReviewItem(item ReviewItem) ReviewItem {
	item.Sources = append([]domain.Source(nil), item.Sources...)
	item.Actions = append([]domain.Action(nil), item.Actions...)
	item.Failures = append([]string(nil), item.Failures...)
	if item.ResolvedAt != nil {
		resolvedAt := *item.ResolvedAt
		item.ResolvedAt = &resolvedAt
	}
	return item
}

var spacePattern = regexp.MustCompile(`\s+`)

func normalizeQuestion(question string) string {
	question = strings.TrimSpace(strings.ToLower(question))
	return spacePattern.ReplaceAllString(question, " ")
}

func reviewItemID(key string) string {
	hash := sha1.Sum([]byte(key))
	return "REV-" + strings.ToUpper(hex.EncodeToString(hash[:])[:12])
}
