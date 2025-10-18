package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

func TestAdminToggleMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("allows request without X-Activity-Source header", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/activities", nil)

		middleware := AdminToggleMiddleware(AdminToggleConfig{
			Logger:         utils.NewTestLogger(),
			AdminCanToggle: true,
		})

		middleware(c)

		assert.Equal(t, http.StatusOK, w.Code)
		_, exists := c.Get("activity_source_override")
		assert.False(t, exists)
	})

	t.Run("rejects toggle when feature is disabled", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/activities", nil)
		c.Request.Header.Set("X-Activity-Source", "postgres")

		middleware := AdminToggleMiddleware(AdminToggleConfig{
			Logger:         utils.NewTestLogger(),
			AdminCanToggle: false,
		})

		middleware(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.True(t, c.IsAborted())
	})

	t.Run("rejects non-admin user", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/activities", nil)
		c.Request.Header.Set("X-Activity-Source", "postgres")
		c.Set("role", "user")

		middleware := AdminToggleMiddleware(AdminToggleConfig{
			Logger:         utils.NewTestLogger(),
			AdminCanToggle: true,
		})

		middleware(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.True(t, c.IsAborted())
	})

	t.Run("rejects invalid source value", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/activities", nil)
		c.Request.Header.Set("X-Activity-Source", "invalid")
		c.Set("role", "admin")

		middleware := AdminToggleMiddleware(AdminToggleConfig{
			Logger:         utils.NewTestLogger(),
			AdminCanToggle: true,
		})

		middleware(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.True(t, c.IsAborted())
	})

	t.Run("allows admin to toggle to postgres", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/activities", nil)
		c.Request.Header.Set("X-Activity-Source", "postgres")
		c.Set("role", "admin")

		middleware := AdminToggleMiddleware(AdminToggleConfig{
			Logger:         utils.NewTestLogger(),
			AdminCanToggle: true,
		})

		middleware(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.False(t, c.IsAborted())
		override, exists := c.Get("activity_source_override")
		assert.True(t, exists)
		assert.Equal(t, "postgres", override)
	})

	t.Run("allows admin to toggle to firestore", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/activities", nil)
		c.Request.Header.Set("X-Activity-Source", "firestore")
		c.Set("role", "admin")

		middleware := AdminToggleMiddleware(AdminToggleConfig{
			Logger:         utils.NewTestLogger(),
			AdminCanToggle: true,
		})

		middleware(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.False(t, c.IsAborted())
		override, exists := c.Get("activity_source_override")
		assert.True(t, exists)
		assert.Equal(t, "firestore", override)
	})

	t.Run("rejects when role is missing", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/activities", nil)
		c.Request.Header.Set("X-Activity-Source", "postgres")

		middleware := AdminToggleMiddleware(AdminToggleConfig{
			Logger:         utils.NewTestLogger(),
			AdminCanToggle: true,
		})

		middleware(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.True(t, c.IsAborted())
	})
}

