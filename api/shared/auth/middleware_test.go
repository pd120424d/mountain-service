package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewServiceAuthMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("it returns an error when Authorization header is missing", func(t *testing.T) {
		serviceAuth := NewServiceAuth(ServiceAuthConfig{
			Secret:      "test-secret",
			ServiceName: "test-service",
			TokenTTL:    1 * time.Hour,
		})
		funcToTest := NewServiceAuthMiddleware(serviceAuth)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		funcToTest(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header required")
	})

	t.Run("it returns an error when Authorization header is invalid", func(t *testing.T) {
		serviceAuth := NewServiceAuth(ServiceAuthConfig{
			Secret:      "test-secret",
			ServiceName: "test-service",
			TokenTTL:    1 * time.Hour,
		})
		funcToTest := NewServiceAuthMiddleware(serviceAuth)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx.Request.Header.Set("Authorization", "Invalid token")

		funcToTest(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid authorization header format")
	})

	t.Run("it returns an error when JWT token is invalid", func(t *testing.T) {
		serviceAuth := NewServiceAuth(ServiceAuthConfig{
			Secret:      "test-secret",
			ServiceName: "test-service",
			TokenTTL:    1 * time.Hour,
		})
		funcToTest := NewServiceAuthMiddleware(serviceAuth)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx.Request.Header.Set("Authorization", "Bearer invalid-token")

		funcToTest(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid token")
	})

	t.Run("it succeeds with valid JWT token", func(t *testing.T) {
		serviceAuth := NewServiceAuth(ServiceAuthConfig{
			Secret:      "test-secret",
			ServiceName: "test-service",
			TokenTTL:    1 * time.Hour,
		})
		token, err := serviceAuth.GenerateToken()
		assert.NoError(t, err)

		funcToTest := NewServiceAuthMiddleware(serviceAuth)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx.Request.Header.Set("Authorization", "Bearer "+token)

		funcToTest(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "test-service", ctx.GetString("service_name"))
	})
}

func TestOptionalServiceAuthMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("it does not return an error when Authorization header is missing", func(t *testing.T) {
		serviceAuth := NewServiceAuth(ServiceAuthConfig{
			Secret:      "test-secret",
			ServiceName: "test-service",
			TokenTTL:    1 * time.Hour,
		})
		funcToTest := OptionalServiceAuthMiddleware(serviceAuth)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

		funcToTest(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, false, ctx.GetBool("is_service_request"))
	})

	t.Run("it does not return an error when Authorization header is invalid", func(t *testing.T) {
		serviceAuth := NewServiceAuth(ServiceAuthConfig{
			Secret:      "test-secret",
			ServiceName: "test-service",
			TokenTTL:    1 * time.Hour,
		})
		funcToTest := OptionalServiceAuthMiddleware(serviceAuth)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx.Request.Header.Set("Authorization", "Invalid token")

		funcToTest(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, false, ctx.GetBool("is_service_request"))
	})

	t.Run("it does not return an error when JWT token is invalid", func(t *testing.T) {
		serviceAuth := NewServiceAuth(ServiceAuthConfig{
			Secret:      "test-secret",
			ServiceName: "test-service",
			TokenTTL:    1 * time.Hour,
		})
		funcToTest := OptionalServiceAuthMiddleware(serviceAuth)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx.Request.Header.Set("Authorization", "Bearer invalid-token")

		funcToTest(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, false, ctx.GetBool("is_service_request"))
	})

	t.Run("it succeeds with valid JWT token", func(t *testing.T) {
		serviceAuth := NewServiceAuth(ServiceAuthConfig{
			Secret:      "test-secret",
			ServiceName: "test-service",
			TokenTTL:    1 * time.Hour,
		})
		token, err := serviceAuth.GenerateToken()
		assert.NoError(t, err)

		funcToTest := OptionalServiceAuthMiddleware(serviceAuth)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx.Request.Header.Set("Authorization", "Bearer "+token)

		funcToTest(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, true, ctx.GetBool("is_service_request"))
		assert.Equal(t, "test-service", ctx.GetString("service_name"))
	})
}
