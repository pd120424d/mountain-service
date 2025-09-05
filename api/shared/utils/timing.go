package utils

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// TimeOperation is a small helper for timing blocks of code.
// Usage:
//
//	defer utils.TimeOperation(ctx, log, "UrgencyRepository.ListPaginated")()
//
// It attaches request_id automatically via log.WithContext(ctx).
func TimeOperation(ctx context.Context, log Logger, operation string, extraFields ...zap.Field) func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		l := log.WithContext(ctx)
		fields := []zap.Field{
			zap.String("op", operation),
			zap.Duration("duration", elapsed),
		}
		if len(extraFields) > 0 {
			fields = append(fields, extraFields...)
		}
		l.Info("timing", fields...)
	}
}
