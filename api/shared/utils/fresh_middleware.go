package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pd120424d/mountain-service/api/shared/config"
)

// FreshReadWindowMiddleware reads X-Fresh-Until header and applies a fresh-read window
// to the request context so downstream repositories can route reads to primary.
// If the header timestamp is in the future, it sets a window for the remaining duration.
func FreshReadWindowMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		expires := c.GetHeader(config.FreshWindowHeader)

		if expires == "" {
			c.Next()
			return
		}

		t, err := time.Parse(time.RFC3339, expires)

		if err != nil {
			c.Next()
			return
		}

		now := time.Now().UTC()
		if t.After(now) {
			d := time.Until(t)
			ctx := WithFreshWindow(c.Request.Context(), d)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}
