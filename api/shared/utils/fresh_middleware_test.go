package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pd120424d/mountain-service/api/shared/config"
	"github.com/stretchr/testify/assert"
)

func TestFreshReadWindowMiddleware_AllowsChainWithoutHeader(t *testing.T) {
	t.Parallel()

	t.Run("it allows the request to continue when the header is missing", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.Use(FreshReadWindowMiddleware())
		r.GET("/ok", func(c *gin.Context) { c.String(200, "ok") })

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "ok", w.Body.String())
	})

	t.Run("it allows the request to continue when the header is invalid", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.Use(FreshReadWindowMiddleware())
		r.GET("/ok", func(c *gin.Context) { c.String(200, "ok") })

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		req.Header.Set(config.FreshWindowHeader, "invalid")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "ok", w.Body.String())
	})

	t.Run("it allows the request to continue when the header is in the past", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.Use(FreshReadWindowMiddleware())
		r.GET("/ok", func(c *gin.Context) { c.String(200, "ok") })

		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		req.Header.Set(config.FreshWindowHeader, time.Now().Add(-1*time.Second).UTC().Format(time.RFC3339))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "ok", w.Body.String())
	})

	t.Run("it sets the context when the header is in the future", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.Use(FreshReadWindowMiddleware())
		var hadFresh bool
		r.GET("/check", func(c *gin.Context) {
			hadFresh = IsFreshRequired(c.Request.Context())
			c.String(200, "ok")
		})

		req := httptest.NewRequest(http.MethodGet, "/check", nil)
		req.Header.Set(config.FreshWindowHeader, time.Now().Add(2*time.Second).UTC().Format(time.RFC3339))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "ok", w.Body.String())
		assert.True(t, hadFresh)
	})
}
