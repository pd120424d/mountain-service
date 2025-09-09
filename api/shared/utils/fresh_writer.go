package utils

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pd120424d/mountain-service/api/shared/config"
)

// WriteFreshWindow applies a fresh-read window to the request context and
// writes the X-Fresh-Until response header with the given duration.
// Safe for tests where ctx.Request may be nil.
func WriteFreshWindow(c *gin.Context, d time.Duration) {
	if c == nil || d <= 0 {
		return
	}

	base := context.Background()
	if c.Request != nil {
		base = c.Request.Context()
	}

	fresh := WithFreshWindow(base, d)
	if c.Request != nil {
		c.Request = c.Request.WithContext(fresh)
	}

	expires := time.Now().Add(d).UTC().Format(time.RFC3339)
	c.Writer.Header().Set(config.FreshWindowHeader, expires)
}

