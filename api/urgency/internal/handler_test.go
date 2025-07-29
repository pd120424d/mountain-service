package internal

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
)

func TestUrgencyHandler_CreateUrgency(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it returns an error when request payload is invalid json", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		invalidPayload := `{
			"name": "Test Urgency",
			"email": "test@example.com"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/urgencies", strings.NewReader(invalidPayload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewUrgencyHandler(log, nil)
		handler.CreateUrgency(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("it returns an error when validation fails - missing name", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"name": "",
			"email": "test@example.com",
			"contactPhone": "123456789",
			"location": "N 43.401123 E 22.662756",
			"description": "Test description",
			"level": "high"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/urgencies", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewUrgencyHandler(log, nil)
		handler.CreateUrgency(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Name")
	})

	t.Run("it returns an error when validation fails - invalid email", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"firstName": "Marko",
			"lastName": "Markovic",
			"email": "invalid-email",
			"contactPhone": "123456789",
			"location": "N 43.401123 E 22.662756",
			"description": "Test description",
			"level": "high"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/urgencies", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewUrgencyHandler(log, nil)
		handler.CreateUrgency(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "email")
	})

	t.Run("it returns an error when validation fails - missing contact phone", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"name": "Test Urgency",
			"email": "test@example.com",
			"contactPhone": "",
			"location": "N 43.401123 E 22.662756",
			"description": "Test description",
			"level": "high"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/urgencies", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewUrgencyHandler(log, nil)
		handler.CreateUrgency(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "ContactPhone")
	})

	t.Run("it returns an error when validation fails - missing location", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"name": "Test Urgency",
			"email": "test@example.com",
			"contactPhone": "123456789",
			"location": "",
			"description": "Test description",
			"level": "high"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/urgencies", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewUrgencyHandler(log, nil)
		handler.CreateUrgency(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Location")
	})

	t.Run("it returns an error when validation fails - missing description", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"name": "Test Urgency",
			"email": "test@example.com",
			"contactPhone": "123456789",
			"location": "N 43.401123 E 22.662756",
			"description": "",
			"level": "high"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/urgencies", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewUrgencyHandler(log, nil)
		handler.CreateUrgency(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Description")
	})

	t.Run("it returns an error when validation fails - invalid urgency level", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"firstName": "John",
			"lastName": "Doe",
			"email": "test@example.com",
			"contactPhone": "123456789",
			"location": "N 43.401123 E 22.662756",
			"description": "Test description",
			"level": "invalid"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/urgencies", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewUrgencyHandler(log, nil)
		handler.CreateUrgency(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid urgency level")
	})

	t.Run("it returns an error when service fails to create urgency", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"firstName": "John",
			"lastName": "Doe",
			"email": "test@example.com",
			"contactPhone": "123456789",
			"location": "N 43.401123 E 22.662756",
			"description": "Test description",
			"level": "high"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/urgencies", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		mockService := NewMockUrgencyService(ctrl)
		mockService.EXPECT().CreateUrgency(gomock.Any()).Return(errors.New("database error")).Times(1)

		handler := NewUrgencyHandler(log, mockService)
		handler.CreateUrgency(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})

	t.Run("it successfully creates urgency when data is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		payload := `{
			"firstName": "John",
			"lastName": "Doe",
			"email": "test@example.com",
			"contactPhone": "123456789",
			"location": "N 43.401123 E 22.662756",
			"description": "Test description",
			"level": "high"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/urgencies", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		mockService := NewMockUrgencyService(ctrl)
		mockService.EXPECT().CreateUrgency(gomock.Any()).DoAndReturn(func(urgency *model.Urgency) error {
			urgency.ID = 1
			urgency.CreatedAt = time.Now()
			urgency.UpdatedAt = time.Now()
			return nil
		}).Times(1)

		handler := NewUrgencyHandler(log, mockService)
		handler.CreateUrgency(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "John")
		assert.Contains(t, w.Body.String(), "Doe")
		assert.Contains(t, w.Body.String(), "test@example.com")
		assert.Contains(t, w.Body.String(), "high")
	})
}

func TestUrgencyHandler_ListUrgencies(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it returns an error when service fails to retrieve urgencies", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		mockService := NewMockUrgencyService(ctrl)
		mockService.EXPECT().GetAllUrgencies().Return(nil, errors.New("database error")).Times(1)

		handler := NewUrgencyHandler(log, mockService)
		handler.ListUrgencies(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})

	t.Run("it returns an empty list when no urgencies exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		mockService := NewMockUrgencyService(ctrl)
		mockService.EXPECT().GetAllUrgencies().Return([]model.Urgency{}, nil).Times(1)

		handler := NewUrgencyHandler(log, mockService)
		handler.ListUrgencies(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "[]", w.Body.String())
	})

	t.Run("it returns a list of urgencies when urgencies exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		urgencies := []model.Urgency{
			{
				ID:           1,
				FirstName:    "Marko",
				LastName:     "Markovic",
				Email:        "test1@example.com",
				ContactPhone: "123456789",
				Location:     "N 43.401123 E 22.662756",
				Description:  "Test description 1",
				Level:        urgencyV1.High,
				Status:       urgencyV1.Open,
			},
			{
				ID:           2,
				FirstName:    "Marko",
				LastName:     "Markovic",
				Email:        "test2@example.com",
				ContactPhone: "987654321",
				Location:     "N 44.401123 E 23.662756",
				Description:  "Test description 2",
				Level:        urgencyV1.Critical,
				Status:       urgencyV1.InProgress,
			},
		}

		mockService := NewMockUrgencyService(ctrl)
		mockService.EXPECT().GetAllUrgencies().Return(urgencies, nil).Times(1)

		handler := NewUrgencyHandler(log, mockService)
		handler.ListUrgencies(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Marko")
		assert.Contains(t, w.Body.String(), "Markovic")
		assert.Contains(t, w.Body.String(), "Marko")
		assert.Contains(t, w.Body.String(), "Markovic")
		assert.Contains(t, w.Body.String(), "test1@example.com")
		assert.Contains(t, w.Body.String(), "test2@example.com")
	})
}

func TestUrgencyHandler_GetUrgency(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it returns an error when urgency ID is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		handler := NewUrgencyHandler(log, nil)
		handler.GetUrgency(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid urgency ID")
	})

	t.Run("it returns an error when urgency does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := NewMockUrgencyService(ctrl)
		mockService.EXPECT().GetUrgencyByID(uint(1)).Return(nil, gorm.ErrRecordNotFound).Times(1)

		handler := NewUrgencyHandler(log, mockService)
		handler.GetUrgency(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "urgency not found")
	})

	t.Run("it returns an error when service fails to get urgency", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := NewMockUrgencyService(ctrl)
		mockService.EXPECT().GetUrgencyByID(uint(1)).Return(nil, errors.New("database error")).Times(1)

		handler := NewUrgencyHandler(log, mockService)
		handler.GetUrgency(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "urgency not found")
	})

	t.Run("it successfully returns urgency when it exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		urgency := &model.Urgency{
			ID:           1,
			FirstName:    "Marko",
			LastName:     "Markovic",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Location:     "N 43.401123 E 22.662756",
			Description:  "Test description",
			Level:        urgencyV1.High,
			Status:       urgencyV1.Open,
		}

		mockService := NewMockUrgencyService(ctrl)
		mockService.EXPECT().GetUrgencyByID(uint(1)).Return(urgency, nil).Times(1)

		handler := NewUrgencyHandler(log, mockService)
		handler.GetUrgency(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Marko")
		assert.Contains(t, w.Body.String(), "Markovic")
		assert.Contains(t, w.Body.String(), "test@example.com")
		assert.Contains(t, w.Body.String(), "high")
	})
}

func TestUrgencyHandler_UpdateUrgency(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it returns an error when urgency ID is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		handler := NewUrgencyHandler(log, nil)
		handler.UpdateUrgency(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid urgency ID")
	})

	t.Run("it returns an error when request payload is invalid json", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		invalidPayload := `{
			"name": "Updated Urgency",
			"email": "updated@example.com"
		`
		ctx.Request = httptest.NewRequest(http.MethodPut, "/urgencies/1", strings.NewReader(invalidPayload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewUrgencyHandler(log, nil)
		handler.UpdateUrgency(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})

	t.Run("it returns an error when validation fails - invalid email", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		payload := `{
			"email": "invalid-email"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPut, "/urgencies/1", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewUrgencyHandler(log, nil)
		handler.UpdateUrgency(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid email format")
	})

	t.Run("it returns an error when validation fails - invalid urgency level", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		payload := `{
			"level": "invalid"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPut, "/urgencies/1", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewUrgencyHandler(log, nil)
		handler.UpdateUrgency(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid urgency level")
	})

	t.Run("it returns an error when validation fails - invalid status", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		payload := `{
			"status": "invalid"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPut, "/urgencies/1", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewUrgencyHandler(log, nil)
		handler.UpdateUrgency(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid status")
	})

	t.Run("it returns an error when urgency does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		payload := `{
			"name": "Updated Urgency"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPut, "/urgencies/1", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		mockService := NewMockUrgencyService(ctrl)
		mockService.EXPECT().GetUrgencyByID(uint(1)).Return(nil, gorm.ErrRecordNotFound).Times(1)

		handler := NewUrgencyHandler(log, mockService)
		handler.UpdateUrgency(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "urgency not found")
	})

	t.Run("it returns an error when service fails to update urgency", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		payload := `{
			"name": "Updated Urgency"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPut, "/urgencies/1", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		existingUrgency := &model.Urgency{
			ID:           1,
			FirstName:    "Marko",
			LastName:     "Markovic",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Location:     "N 43.401123 E 22.662756",
			Description:  "Test description",
			Level:        urgencyV1.High,
			Status:       urgencyV1.Open,
		}

		mockService := NewMockUrgencyService(ctrl)
		mockService.EXPECT().GetUrgencyByID(uint(1)).Return(existingUrgency, nil).Times(1)
		mockService.EXPECT().UpdateUrgency(gomock.Any()).Return(errors.New("database error")).Times(1)

		handler := NewUrgencyHandler(log, mockService)
		handler.UpdateUrgency(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})

	t.Run("it successfully updates urgency when data is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		payload := `{
			"firstName": "Updated",
			"lastName": "Urgency",
			"email": "updated@example.com",
			"status": "in_progress"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPut, "/urgencies/1", strings.NewReader(payload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		existingUrgency := &model.Urgency{
			ID:           1,
			FirstName:    "Marko",
			LastName:     "Markovic",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Location:     "N 43.401123 E 22.662756",
			Description:  "Test description",
			Level:        urgencyV1.High,
			Status:       urgencyV1.Open,
		}

		mockService := NewMockUrgencyService(ctrl)
		mockService.EXPECT().GetUrgencyByID(uint(1)).Return(existingUrgency, nil).Times(1)
		mockService.EXPECT().UpdateUrgency(gomock.Any()).Return(nil).Times(1)

		handler := NewUrgencyHandler(log, mockService)
		handler.UpdateUrgency(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Updated")
		assert.Contains(t, w.Body.String(), "Urgency")
		assert.Contains(t, w.Body.String(), "updated@example.com")
		assert.Contains(t, w.Body.String(), "in_progress")
	})
}

func TestUrgencyHandler_DeleteUrgency(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it returns an error when urgency ID is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "invalid"}}

		handler := NewUrgencyHandler(log, nil)
		handler.DeleteUrgency(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid urgency ID")
	})

	t.Run("it returns an error when service fails to delete urgency", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := NewMockUrgencyService(ctrl)
		mockService.EXPECT().DeleteUrgency(uint(1)).Return(errors.New("database error")).Times(1)

		handler := NewUrgencyHandler(log, mockService)
		handler.DeleteUrgency(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})

	t.Run("it successfully deletes urgency when it exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		mockService := NewMockUrgencyService(ctrl)
		mockService.EXPECT().DeleteUrgency(uint(1)).Return(nil).Times(1)

		handler := NewUrgencyHandler(log, mockService)
		handler.DeleteUrgency(ctx)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func TestUrgencyHandler_ResetAllData(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it returns an error when service fails to reset data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		mockService := NewMockUrgencyService(ctrl)
		mockService.EXPECT().ResetAllData().Return(errors.New("database error")).Times(1)

		handler := NewUrgencyHandler(log, mockService)
		handler.ResetAllData(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})

	t.Run("it successfully resets all data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		mockService := NewMockUrgencyService(ctrl)
		mockService.EXPECT().ResetAllData().Return(nil).Times(1)

		handler := NewUrgencyHandler(log, mockService)
		handler.ResetAllData(ctx)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
