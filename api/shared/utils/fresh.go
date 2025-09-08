package utils

import (
	"context"
	"time"
)

type freshKey struct{}

// WithFreshWindow marks the context to prefer primary DB reads until the given duration elapses.
func WithFreshWindow(ctx context.Context, d time.Duration) context.Context {
	if d <= 0 {
		return ctx
	}
	return context.WithValue(ctx, freshKey{}, time.Now().Add(d))
}

func RequireFresh(ctx context.Context) context.Context {
	return WithFreshWindow(ctx, 2*time.Second)
}

func IsFreshRequired(ctx context.Context) bool {
	v, ok := ctx.Value(freshKey{}).(time.Time)
	return ok && time.Now().Before(v)
}

func FreshUntil(ctx context.Context) (time.Time, bool) {
	v, ok := ctx.Value(freshKey{}).(time.Time)
	return v, ok
}
