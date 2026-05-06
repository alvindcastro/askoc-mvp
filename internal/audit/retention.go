package audit

import (
	"context"
	"time"
)

const DefaultDemoRetention = 7 * 24 * time.Hour

type RetentionPolicy struct {
	MaxAge time.Duration
}

type Pruner interface {
	PruneBefore(context.Context, time.Time) int
}

func DefaultRetentionPolicy() RetentionPolicy {
	return RetentionPolicy{MaxAge: DefaultDemoRetention}
}

func (p RetentionPolicy) PurgeExpired(ctx context.Context, store Pruner, now time.Time) int {
	if store == nil {
		return 0
	}
	maxAge := p.MaxAge
	if maxAge <= 0 {
		maxAge = DefaultDemoRetention
	}
	return store.PruneBefore(ctx, now.Add(-maxAge))
}
