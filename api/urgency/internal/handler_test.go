package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
)

func setupTestHandler(t *testing.T) (*urgencyHandler, *MockUrgencyService, *gin.Engine) {
	ctrl := gomock.NewController(t)
	mockSvc := NewMockUrgencyService(ctrl)
	log := utils.NewTestLogger()

	handler := &urgencyHandler{
		log: log,
		svc: mockSvc,
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()

	return handler, mockSvc, router
}

func TestUrgencyHandler_CreateUrgency(t *testing.T) {
	t.Parallel()

	handler, mockSvc, router := setupTestHandler(t)
	router.POST("/urgencies", handler.CreateUrgency)

	t.Run("it creates a new urgency successfully", func(t *testing.T) {
		req := model.UrgencyCreateRequest{
			Name:         "Test Urgency",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Description:  "Test description",
			Level:        model.High,
		}

		mockSvc.EXPECT().CreateUrgency(gomock.Any()).DoAndReturn(func(urgency *model.Urgency) error {
			urgency.ID = 1
			urgency.CreatedAt = time.Now()
			urgency.UpdatedAt = time.Now()
			return nil
		})

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/urgencies", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response model.UrgencyResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "Test Urgency", response.Name)
		assert.Equal(t, "Open", string(response.Status))
	})

	t.Run("it returns an error when JSON is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/urgencies", bytes.NewBuffer([]byte("invalid json")))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("it returns an error when validation fails", func(t *testing.T) {
		req := model.UrgencyCreateRequest{
			Name:         "", // Missing required field
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Description:  "Test description",
			Level:        model.High,
		}

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/urgencies", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "required")
	})

	t.Run("it returns an error when repository fails", func(t *testing.T) {
		req := model.UrgencyCreateRequest{
			Name:         "Test Urgency",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Description:  "Test description",
			Level:        model.High,
		}

		mockSvc.EXPECT().CreateUrgency(gomock.Any()).Return(errors.New("database error"))

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/urgencies", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})
}

func TestUrgencyHandler_ListUrgencies(t *testing.T) {
	t.Parallel()

	handler, mockSvc, router := setupTestHandler(t)
	router.GET("/urgencies", handler.ListUrgencies)

	t.Run("it lists all urgencies successfully", func(t *testing.T) {
		urgencies := []model.Urgency{
			{
				ID:           1,
				Name:         "Urgency 1",
				Email:        "test1@example.com",
				ContactPhone: "123456789",
				Description:  "Description 1",
				Level:        model.High,
				Status:       model.Open,
			},
			{
				ID:           2,
				Name:         "Urgency 2",
				Email:        "test2@example.com",
				ContactPhone: "987654321",
				Description:  "Description 2",
				Level:        model.Medium,
				Status:       model.InProgress,
			},
		}
		urgencies[0].CreatedAt = time.Now()
		urgencies[0].UpdatedAt = time.Now()
		urgencies[1].CreatedAt = time.Now()
		urgencies[1].UpdatedAt = time.Now()

		mockSvc.EXPECT().GetAllUrgencies().Return(urgencies, nil)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/urgencies", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []model.UrgencyResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, "Urgency 1", response[0].Name)
		assert.Equal(t, "Urgency 2", response[1].Name)
	})

	t.Run("it lists an empty list when no urgencies exist", func(t *testing.T) {
		mockSvc.EXPECT().GetAllUrgencies().Return([]model.Urgency{}, nil)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/urgencies", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []model.UrgencyResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response, 0)
	})

	t.Run("it returns an error when repository fails", func(t *testing.T) {
		mockSvc.EXPECT().GetAllUrgencies().Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/urgencies", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})
}

func TestUrgencyHandler_GetUrgency(t *testing.T) {
	t.Parallel()

	handler, mockSvc, router := setupTestHandler(t)
	router.GET("/urgencies/:id", handler.GetUrgency)

	t.Run("it gets an urgency successfully", func(t *testing.T) {
		urgency := model.Urgency{
			ID:           1,
			Name:         "Test Urgency",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Description:  "Test description",
			Level:        model.High,
			Status:       model.Open,
		}
		urgency.CreatedAt = time.Now()
		urgency.UpdatedAt = time.Now()

		mockSvc.EXPECT().GetUrgencyByID(uint(1)).DoAndReturn(func(id uint) (*model.Urgency, error) {
			urgency := &model.Urgency{
				ID:           1,
				Name:         "Test Urgency",
				Email:        "test@example.com",
				ContactPhone: "123456789",
				Description:  "Test description",
				Level:        model.High,
				Status:       model.Open,
			}
			urgency.CreatedAt = time.Now()
			urgency.UpdatedAt = time.Now()
			return urgency, nil
		})

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/urgencies/1", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.UrgencyResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "Test Urgency", response.Name)
	})

	t.Run("it returns an error when ID is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/urgencies/invalid", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid urgency ID")
	})

	t.Run("it returns an error when urgency is not found", func(t *testing.T) {
		mockSvc.EXPECT().GetUrgencyByID(uint(999)).Return(nil, gorm.ErrRecordNotFound)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/urgencies/999", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "urgency not found")
	})
}

func TestUrgencyHandler_UpdateUrgency(t *testing.T) {
	t.Parallel()

	handler, mockSvc, router := setupTestHandler(t)
	router.PUT("/urgencies/:id", handler.UpdateUrgency)

	t.Run("it updates an urgency successfully", func(t *testing.T) {
		req := model.UrgencyUpdateRequest{
			Name:   "Updated Urgency",
			Email:  "updated@example.com", // Include valid email to avoid Gin validation error
			Status: model.InProgress,
		}

		existingUrgency := model.Urgency{
			ID:           1,
			Name:         "Original Urgency",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Description:  "Test description",
			Level:        model.High,
			Status:       model.Open,
		}
		existingUrgency.CreatedAt = time.Now()
		existingUrgency.UpdatedAt = time.Now()

		mockSvc.EXPECT().GetUrgencyByID(uint(1)).DoAndReturn(func(id uint) (*model.Urgency, error) {
			return &existingUrgency, nil
		})

		mockSvc.EXPECT().UpdateUrgency(gomock.Any()).DoAndReturn(func(urgency *model.Urgency) error {
			urgency.UpdatedAt = time.Now()
			return nil
		})

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/urgencies/1", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.UrgencyResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Updated Urgency", response.Name)
		assert.Equal(t, "In Progress", string(response.Status))
	})

	t.Run("it returns an error when ID is invalid", func(t *testing.T) {
		req := model.UrgencyUpdateRequest{Name: "Updated"}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/urgencies/invalid", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid urgency ID")
	})

	t.Run("it returns an error when JSON is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/urgencies/1", bytes.NewBuffer([]byte("invalid json")))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("it returns an error when validation fails", func(t *testing.T) {
		req := model.UrgencyUpdateRequest{
			Email: "invalid-email", // Invalid email format
		}

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/urgencies/1", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "email")
	})

	t.Run("it returns an error when urgency is not found", func(t *testing.T) {
		req := model.UrgencyUpdateRequest{
			Name:  "Updated",
			Email: "valid@example.com", // Include valid email
		}

		mockSvc.EXPECT().GetUrgencyByID(uint(999)).Return(nil, gorm.ErrRecordNotFound)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/urgencies/999", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "urgency not found")
	})

	t.Run("it returns an error when repository fails", func(t *testing.T) {
		req := model.UrgencyUpdateRequest{
			Name:  "Updated",
			Email: "valid@example.com", // Include valid email
		}

		existingUrgency := model.Urgency{ID: 1, Name: "Original"}
		mockSvc.EXPECT().GetUrgencyByID(uint(1)).DoAndReturn(func(id uint) (*model.Urgency, error) {
			return &existingUrgency, nil
		})
		mockSvc.EXPECT().UpdateUrgency(gomock.Any()).Return(errors.New("database error"))

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("PUT", "/urgencies/1", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})
}

func TestUrgencyHandler_DeleteUrgency(t *testing.T) {
	t.Parallel()

	handler, mockSvc, router := setupTestHandler(t)
	router.DELETE("/urgencies/:id", handler.DeleteUrgency)

	t.Run("it deletes an urgency successfully", func(t *testing.T) {
		mockSvc.EXPECT().DeleteUrgency(uint(1)).Return(nil)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("DELETE", "/urgencies/1", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())
	})

	t.Run("it returns an error when ID is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("DELETE", "/urgencies/invalid", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid urgency ID")
	})

	t.Run("it returns an error when urgency is not found", func(t *testing.T) {
		mockSvc.EXPECT().DeleteUrgency(uint(999)).Return(gorm.ErrRecordNotFound)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("DELETE", "/urgencies/999", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "record not found")
	})

	t.Run("it returns an error when repository fails", func(t *testing.T) {
		mockSvc.EXPECT().DeleteUrgency(uint(1)).Return(errors.New("database error"))

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("DELETE", "/urgencies/1", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})
}

func TestUrgencyHandler_ResetAllData(t *testing.T) {
	t.Parallel()

	handler, mockSvc, router := setupTestHandler(t)
	router.DELETE("/admin/urgencies/reset", handler.ResetAllData)

	t.Run("it resets all urgencies successfully", func(t *testing.T) {
		mockSvc.EXPECT().ResetAllData().Return(nil)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("DELETE", "/admin/urgencies/reset", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())
	})

	t.Run("it returns an error when repository fails", func(t *testing.T) {
		mockSvc.EXPECT().ResetAllData().Return(errors.New("database error"))

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("DELETE", "/admin/urgencies/reset", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})
}

func TestNewUrgencyHandler(t *testing.T) {
	log := utils.NewTestLogger()
	ctrl := gomock.NewController(t)
	mockSvc := NewMockUrgencyService(ctrl)

	handler := NewUrgencyHandler(log, mockSvc)

	assert.NotNil(t, handler)
	assert.IsType(t, &urgencyHandler{}, handler)
}
