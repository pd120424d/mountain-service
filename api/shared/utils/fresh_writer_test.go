package utils

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pd120424d/mountain-service/api/shared/config"
	"github.com/stretchr/testify/assert"
)

func TestWriteFreshWindow(t *testing.T) {
	t.Parallel()

	t.Run("it does not set the header or context when duration is not positive", func(t *testing.T) {
		testSetGinMode()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/test", nil)

		WriteFreshWindow(c, 0)

		expiresStr := w.Header().Get(config.FreshWindowHeader)
		assert.Empty(t, expiresStr, "header should not be set")

		assert.False(t, IsFreshRequired(c.Request.Context()))
	})

	t.Run("it does not set the header or context when context is nil", func(t *testing.T) {
		testSetGinMode()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/test", nil)

		WriteFreshWindow(nil, 200*time.Millisecond)

		expiresStr := w.Header().Get(config.FreshWindowHeader)
		assert.Empty(t, expiresStr, "header should not be set")
	})

	t.Run("it sets the header and context when duration is positive", func(t *testing.T) {

		testSetGinMode()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/test", nil)

		d := 2 * time.Second
		before := time.Now().UTC()
		WriteFreshWindow(c, d)

		expiresStr := w.Header().Get(config.FreshWindowHeader)
		if assert.NotEmpty(t, expiresStr, "header should be set") {
			expires, err := time.Parse(time.RFC3339, expiresStr)
			assert.NoError(t, err)
			assert.True(t, expires.After(before), "expiry should be in the future")
		}

		assert.True(t, IsFreshRequired(c.Request.Context()))
	})
}
