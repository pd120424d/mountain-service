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

		// Debug: Always log that timing is being called
		l.Infof("TIMING_DEBUG: About to log timing for operation: %s", operation)

		fields := []zap.Field{
			zap.String("op", operation),
			zap.Duration("duration", elapsed),
		}
		if len(extraFields) > 0 {
			fields = append(fields, extraFields...)
		}
		l.Info("Timing", fields...)

		// Debug: Confirm timing was logged
		l.Infof("TIMING_DEBUG: Finished logging timing for operation: %s (took %v)", operation, elapsed)
	}
}
