package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ctxKey string

const (
	RequestIDKey    ctxKey = "request_id"
	HeaderRequestID string = "X-Request-ID"
)

// RequestIDFromContext extracts the request ID from the context, if present.
func RequestIDFromContext(ctx context.Context) string {
	if v := ctx.Value(RequestIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

func ContextWithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}

// RequestIDMiddleware ensures every request has an X-Request-ID and places it into the context.
// - If the incoming request already has X-Request-ID, it is preserved.
// - Otherwise, a new cryptographically-strong random ID is generated.
// The value is available in:
//   - gin.Context via c.Request.Context()
//   - HTTP response header X-Request-ID
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(HeaderRequestID)
		if id == "" {
			id = generateRequestID()
		}

		ctx := ContextWithRequestID(c.Request.Context(), id)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(HeaderRequestID, id)

		c.Next()
	}
}

// generateRequestID returns a 32-char hex string based on 16 random bytes.
func generateRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback to a timestamp-based ID if crypto fails, it should never happen
		return hex.EncodeToString([]byte("fallback-request-id"))
	}
	return hex.EncodeToString(b)
}

func SetRequestIDHeader(req *http.Request, id string) {
	if req == nil || id == "" {
		return
	}
	req.Header.Set(HeaderRequestID, id)
}
