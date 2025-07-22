package internal

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"

	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/auth"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
	"github.com/pd120424d/mountain-service/api/urgency/internal/repositories"
)

func TestIntegration_UrgencyLifecycle(t *testing.T) {
	router, _, cleanup := setupIntegrationTest(t)
	defer cleanup()

	token, err := generateTestJWT(1, "Medic")
	require.NoError(t, err)

	authHeader := "Bearer " + token

	t.Run("it successfully completes the urgency lifecycle (POST, GET, PUT and DELETE)", func(t *testing.T) {
		createReq := model.UrgencyCreateRequest{
			Name:         "Mountain Rescue Emergency",
			Email:        "rescue@example.com",
			ContactPhone: "123456789",
			Location:     "N 43.401123 E 22.662756",
			Description:  "Hiker injured on mountain trail",
			Level:        model.Critical,
		}

		body, _ := json.Marshal(createReq)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/v1/urgencies", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var createResponse model.UrgencyResponse
		err := json.Unmarshal(w.Body.Bytes(), &createResponse)
		require.NoError(t, err)
		assert.Equal(t, "Mountain Rescue Emergency", createResponse.Name)
		assert.Equal(t, "N 43.401123 E 22.662756", createResponse.Location)
		assert.Equal(t, "Open", string(createResponse.Status))
		urgencyID := createResponse.ID

		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "/api/v1/urgencies", nil)
		req.Header.Set("Authorization", authHeader)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var listResponse []model.UrgencyResponse
		err = json.Unmarshal(w.Body.Bytes(), &listResponse)
		require.NoError(t, err)
		assert.Len(t, listResponse, 1)
		assert.Equal(t, urgencyID, listResponse[0].ID)

		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/urgencies/%d", urgencyID), nil)
		req.Header.Set("Authorization", authHeader)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var getResponse model.UrgencyResponse
		err = json.Unmarshal(w.Body.Bytes(), &getResponse)
		require.NoError(t, err)
		assert.Equal(t, urgencyID, getResponse.ID)
		assert.Equal(t, "Mountain Rescue Emergency", getResponse.Name)

		updateReq := model.UrgencyUpdateRequest{
			Status: model.InProgress,
			Email:  "updated@example.com",
		}

		body, _ = json.Marshal(updateReq)
		w = httptest.NewRecorder()
		req = httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/urgencies/%d", urgencyID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", authHeader)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var updateResponse model.UrgencyResponse
		err = json.Unmarshal(w.Body.Bytes(), &updateResponse)
		require.NoError(t, err)
		assert.Equal(t, "In Progress", string(updateResponse.Status))
		assert.Equal(t, "updated@example.com", updateResponse.Email)

		w = httptest.NewRecorder()
		req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/urgencies/%d", urgencyID), nil)
		req.Header.Set("Authorization", authHeader)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/urgencies/%d", urgencyID), nil)
		req.Header.Set("Authorization", authHeader)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "/api/v1/urgencies", nil)
		req.Header.Set("Authorization", authHeader)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &listResponse)
		require.NoError(t, err)
		assert.Len(t, listResponse, 0)
	})
}

func TestIntegration_AdminOperations(t *testing.T) {
	router, db, cleanup := setupIntegrationTest(t)
	defer cleanup()

	adminToken, err := generateTestJWT(1, "Administrator")
	require.NoError(t, err)

	adminAuthHeader := "Bearer " + adminToken

	t.Run("it successfully resets all urgencies as admininistrator", func(t *testing.T) {
		urgency1 := model.Urgency{
			Name:         "Test 1",
			Email:        "test1@example.com",
			ContactPhone: "123456789",
			Description:  "Test description 1",
			Level:        model.High,
			Status:       model.Open,
		}
		urgency2 := model.Urgency{
			Name:         "Test 2",
			Email:        "test2@example.com",
			ContactPhone: "987654321",
			Description:  "Test description 2",
			Level:        model.Medium,
			Status:       model.InProgress,
		}

		db.Create(&urgency1)
		db.Create(&urgency2)

		var count int64
		db.Model(&model.Urgency{}).Count(&count)
		assert.Equal(t, int64(2), count)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/api/v1/admin/urgencies/reset", nil)
		req.Header.Set("Authorization", adminAuthHeader)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		w = httptest.NewRecorder()
		req = httptest.NewRequest("DELETE", "/api/v1/admin/urgencies/reset", nil)
		req.Header.Set("Authorization", adminAuthHeader)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		db.Model(&model.Urgency{}).Count(&count)
		assert.Equal(t, int64(0), count)
	})
}

func TestIntegration_AuthenticationErrors(t *testing.T) {
	router, _, cleanup := setupIntegrationTest(t)
	defer cleanup()

	t.Run("it returns an error when Authorization header is missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/urgencies", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("it returns an error when Authorization header is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/v1/urgencies", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestIntegration_HealthCheck(t *testing.T) {
	router, _, cleanup := setupIntegrationTest(t)
	defer cleanup()

	t.Run("it returns pong for health check endpoint", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/ping", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "pong")
	})
}

// setupIntegrationTest sets up a complete test environment with real database and handlers
func setupIntegrationTest(t *testing.T) (*gin.Engine, *gorm.DB, func()) {
	os.Setenv("JWT_SECRET", "test-secret-key")

	sqlDB, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	db, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Urgency{})
	require.NoError(t, err)

	log := utils.NewTestLogger()

	svc := NewUrgencyService(repositories.NewUrgencyRepository(log, db))
	urgencyHandler := NewUrgencyHandler(log, svc)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	authorized := router.Group("/api/v1").Use(auth.AuthMiddleware(log))
	{
		authorized.POST("/urgencies", urgencyHandler.CreateUrgency)
		authorized.GET("/urgencies", urgencyHandler.ListUrgencies)
		authorized.GET("/urgencies/:id", urgencyHandler.GetUrgency)
		authorized.PUT("/urgencies/:id", urgencyHandler.UpdateUrgency)
		authorized.DELETE("/urgencies/:id", urgencyHandler.DeleteUrgency)
	}

	admin := router.Group("/api/v1/admin").Use(auth.AdminMiddleware(log))
	{
		admin.DELETE("/urgencies/reset", urgencyHandler.ResetAllData)
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	cleanup := func() {
		os.Unsetenv("JWT_SECRET")
	}

	return router, db, cleanup
}

func generateTestJWT(employeeID uint, role string) (string, error) {
	claims := auth.EmployeeClaims{
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
