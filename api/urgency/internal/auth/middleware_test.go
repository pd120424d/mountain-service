package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	// Remove t.Parallel() to avoid race conditions with environment variables

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

	t.Run("it returns an error when JWT token is malformed", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret-key")
		defer os.Unsetenv("JWT_SECRET")

		log := utils.NewTestLogger()
		funcToTest := AuthMiddleware(log)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/urgencies", nil)
		ctx.Request.Header.Set("Authorization", "Bearer invalid.jwt.token")

		funcToTest(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid token")
	})

	t.Run("it returns an error when JWT token is expired", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret-key")
		defer os.Unsetenv("JWT_SECRET")

		expiredToken, err := generateExpiredTestJWT(1, "Medic")
		require.NoError(t, err)

		log := utils.NewTestLogger()
		funcToTest := AuthMiddleware(log)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/urgencies", nil)
		ctx.Request.Header.Set("Authorization", "Bearer "+expiredToken)

		funcToTest(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid token")
	})

	t.Run("it succeeds with valid JWT token", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret-key")
		defer os.Unsetenv("JWT_SECRET")

		validToken, err := generateTestJWT(1, "Medic")
		require.NoError(t, err)

		log := utils.NewTestLogger()

		router := gin.New()
		router.Use(AuthMiddleware(log))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})
}

func TestAdminMiddleware(t *testing.T) {
	// Remove t.Parallel() to avoid race conditions with environment variables

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

	t.Run("it returns an error when JWT token is malformed", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret-key")
		defer os.Unsetenv("JWT_SECRET")

		log := utils.NewTestLogger()
		funcToTest := AdminMiddleware(log)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/admin/urgencies/reset", nil)
		ctx.Request.Header.Set("Authorization", "Bearer invalid.jwt.token")

		funcToTest(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid token")
	})

	t.Run("it returns forbidden when user is not administrator", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret-key")
		defer os.Unsetenv("JWT_SECRET")

		nonAdminToken, err := generateTestJWT(1, "Medic")
		require.NoError(t, err)

		log := utils.NewTestLogger()
		funcToTest := AdminMiddleware(log)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/admin/urgencies/reset", nil)
		ctx.Request.Header.Set("Authorization", "Bearer "+nonAdminToken)

		funcToTest(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "Administrator access required")
	})

	t.Run("it succeeds with valid admin JWT token", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret-key")
		defer os.Unsetenv("JWT_SECRET")

		adminToken, err := generateTestJWT(1, "Administrator")
		require.NoError(t, err)

		log := utils.NewTestLogger()

		router := gin.New()
		router.Use(AdminMiddleware(log))
		router.DELETE("/admin/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/admin/test", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})
}

func generateTestJWT(employeeID uint, role string) (string, error) {
	claims := EmployeeClaims{
		ID:   employeeID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "test-secret-key"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func generateExpiredTestJWT(employeeID uint, role string) (string, error) {
	claims := EmployeeClaims{
		ID:   employeeID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired 1 hour ago
		},
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "test-secret-key"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
