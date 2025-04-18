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
		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees", nil)

		funcToTest(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Unauthorized")
	})

	t.Run("it returns an error when token is invalid", func(t *testing.T) {
		log := utils.NewTestLogger()
		funcToTest := AuthMiddleware(log)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees", nil)
		ctx.Request.Header.Set("Authorization", "Bearer invalid-token")

		funcToTest(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid token")
	})

	t.Run("it sets claims in context when token is valid", func(t *testing.T) {
		log := utils.NewTestLogger()
		funcToTest := AuthMiddleware(log)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/employees", nil)

		token, _ := GenerateJWT(123, "Medic")
		ctx.Request.Header.Set("Authorization", "Bearer "+token)

		funcToTest(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, uint(123), ctx.GetUint("employeeID"))
		assert.Equal(t, "Medic", ctx.GetString("role"))
	})
}
