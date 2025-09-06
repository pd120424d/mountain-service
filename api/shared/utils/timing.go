package utils

import (
	"time"

	"go.uber.org/zap"
)

// TimeOperation is a small helper for timing blocks of code.
// Usage:
//
//	defer utils.TimeOperation(log, "UrgencyRepository.ListPaginated")()
//
// The passed logger should already have context attached via log.WithContext(ctx).
func TimeOperation(log Logger, operation string, extraFields ...zap.Field) func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)

		fields := []zap.Field{
			zap.String("op", operation),
			zap.Duration("duration", elapsed),
			zap.String("duration_human", elapsed.String()),
		}
		if len(extraFields) > 0 {
			fields = append(fields, extraFields...)
		}
		log.Info("Timing", fields...)
	}
}
