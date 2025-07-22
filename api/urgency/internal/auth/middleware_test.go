package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	t.Run("it returns an error when Authorization header is missing", func(t *testing.T) {
		log := utils.NewTestLogger()
		funcToTest := AuthMiddleware(log)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/urgencies", nil)

		funcToTest(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")
	})

	t.Run("it returns an error when Authorization header is invalid", func(t *testing.T) {
		log := utils.NewTestLogger()
		funcToTest := AuthMiddleware(log)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/urgencies", nil)
		ctx.Request.Header.Set("Authorization", "Invalid token")

		funcToTest(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")
	})
}

func TestAdminMiddleware(t *testing.T) {
	t.Run("it returns an error when Authorization header is missing", func(t *testing.T) {
		log := utils.NewTestLogger()
		funcToTest := AdminMiddleware(log)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/admin/urgencies/reset", nil)

		funcToTest(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")
	})

	t.Run("it returns an error when Authorization header is invalid", func(t *testing.T) {
		log := utils.NewTestLogger()
		funcToTest := AdminMiddleware(log)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/admin/urgencies/reset", nil)
		ctx.Request.Header.Set("Authorization", "Invalid token")

		funcToTest(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")
	})
}
