package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/pd120424d/mountain-service/api/activity/internal/service"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	sharedModels "github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type readModelFake struct {
	items []sharedModels.Activity
	err   error
}

func (f *readModelFake) ListByUrgency(_ context.Context, urgencyID uint, limit int) ([]sharedModels.Activity, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.items, nil
}

func (f *readModelFake) ListAll(_ context.Context, limit int) ([]sharedModels.Activity, error) {
	return nil, fmt.Errorf("not used")
}

func TestActivityHandler_CreateActivity(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it returns an error when request payload is invalid json", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		invalidPayload := `{
			"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/activities", strings.NewReader(invalidPayload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewActivityHandler(log, nil, nil)
		handler.CreateActivity(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Invalid request payload: invalid character '\\\\n' in string literal\"}")
	})

	t.Run("it returns an error when validation fails", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		invalidPayload := `{
			"description": "   ",
			"employeeId": 1,
			"urgencyId": 2
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/activities", strings.NewReader(invalidPayload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		handler := NewActivityHandler(log, nil, nil)
		handler.CreateActivity(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "description is required")
	})

	t.Run("it returns an error when service fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		validPayload := `{
			"description": "Test",
			"employeeId": 1,
			"urgencyId": 2
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/activities", strings.NewReader(validPayload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().CreateActivity(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("database error"))

		handler := NewActivityHandler(log, svcMock, nil)
		handler.CreateActivity(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "{\"details\":\"database error\",\"error\":\"Failed to create activity\"}")
	})

	t.Run("it successfully creates activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		validPayload := `{
			"description": "Test",
			"employeeId": 1,
			"urgencyId": 2
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/activities", strings.NewReader(validPayload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().CreateActivity(gomock.Any(), gomock.Any()).Return(&activityV1.ActivityResponse{ID: 1}, nil)

		handler := NewActivityHandler(log, svcMock, nil)
		handler.CreateActivity(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "{\"id\":1")
	})
}

func TestActivityHandler_GetActivity(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it returns an error when ID is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/activities/invalid", nil)

		ctx.Params = []gin.Param{{Key: "id", Value: "invalid"}}
		handler := NewActivityHandler(log, nil, nil)
		handler.GetActivity(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Invalid activity ID\"}")
	})

	t.Run("it returns an error when activity is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/activities/1", nil)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().GetActivityByID(gomock.Any(), uint(1)).Return(nil, fmt.Errorf("not found"))

		handler := NewActivityHandler(log, svcMock, nil)
		handler.GetActivity(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Activity not found\"}")
	})

	t.Run("it successfully retrieves activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/activities/1", nil)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().GetActivityByID(gomock.Any(), uint(1)).Return(&activityV1.ActivityResponse{ID: 1}, nil)

		handler := NewActivityHandler(log, svcMock, nil)
		handler.GetActivity(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "{\"id\":1")
	})
}

func TestActivityHandler_ListActivities(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it returns an error when service fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/activities", nil)

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().ListActivities(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("database error"))

		handler := NewActivityHandler(log, svcMock, nil)
		handler.ListActivities(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "{\"details\":\"database error\",\"error\":\"Failed to list activities\"}")
	})

	t.Run("it successfully retrieves activities", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/activities", nil)

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().ListActivities(gomock.Any(), gomock.Any()).Return(&activityV1.ActivityListResponse{
			Activities: []activityV1.ActivityResponse{{ID: 1}},
			Total:      1,
			Page:       1,
			PageSize:   10,
		}, nil)

		handler := NewActivityHandler(log, svcMock, nil)
		handler.ListActivities(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "{\"activities\":[{\"id\":1")
	})

	t.Run("it handles service error correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/activities", nil)

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().ListActivities(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("service error"))

		handler := NewActivityHandler(log, svcMock, nil)
		handler.ListActivities(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")
	})

	t.Run("it handles query parameters correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		// Set query parameters
		ctx.Request = httptest.NewRequest("GET", "/activities?page=2&pageSize=25&type=employee_created&level=info", nil)

		svcMock := service.NewMockActivityService(ctrl)

		svcMock.EXPECT().ListActivities(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, req *activityV1.ActivityListRequest) (*activityV1.ActivityListResponse, error) {
			assert.Equal(t, 2, req.Page)
			assert.Equal(t, 25, req.PageSize)
			return &activityV1.ActivityListResponse{
				Activities: []activityV1.ActivityResponse{},
				Total:      0,
				Page:       2,
				PageSize:   25,
				TotalPages: 0,
			}, nil
		})

		handler := NewActivityHandler(log, svcMock, nil)
		handler.ListActivities(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
	})

}

func TestActivityHandler_ListActivities_ReadModelPreference(t *testing.T) {
	t.Parallel()
	log := utils.NewTestLogger()
	ctrl := gomock.NewController(t)

	// Service mock for fallback
	svcMock := service.NewMockActivityService(ctrl)

	// fake read model implementing ActivityReadModel
	readModel := &readModelFake{
		items: []sharedModels.Activity{{ID: 42, Description: "rm", UrgencyID: 7, EmployeeID: 3}},
	}

	t.Run("it successfully retrieves activities via read-model", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/activities?urgencyId=7", nil)

		h := NewActivityHandler(log, svcMock, readModel)
		h.ListActivities(ctx)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "\"id\":42")
	})

	t.Run("it successfully retrieves empty list via read-model", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/activities?urgencyId=7", nil)
		readModel.items = []sharedModels.Activity{}

		h := NewActivityHandler(log, svcMock, readModel)
		h.ListActivities(ctx)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "\"activities\":[]")
	})

	t.Run("it falls back to service on read-model error", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/activities?urgencyId=8", nil)

		readModel.err = fmt.Errorf("rm error")
		svcMock.EXPECT().ListActivities(gomock.Any(), gomock.Any()).Return(&activityV1.ActivityListResponse{Activities: []activityV1.ActivityResponse{{ID: 77}}, Total: 1, Page: 1, PageSize: 10}, nil)

		h := NewActivityHandler(log, svcMock, readModel)
		h.ListActivities(ctx)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "\"id\":77")
	})
}
func TestActivityHandler_GetActivity_EdgeCases(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it returns error for invalid ID parameter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/activities/invalid", nil)
		ctx.Params = gin.Params{{Key: "id", Value: "invalid"}}

		svcMock := service.NewMockActivityService(ctrl)
		handler := NewActivityHandler(log, svcMock, nil)
		handler.GetActivity(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid activity ID")
	})

	t.Run("it returns error for zero ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/activities/0", nil)
		ctx.Params = gin.Params{{Key: "id", Value: "0"}}

		svcMock := service.NewMockActivityService(ctrl)
		// The handler will call the service with ID 0, so we need to expect it
		svcMock.EXPECT().GetActivityByID(gomock.Any(), uint(0)).Return(nil, fmt.Errorf("activity not found"))

		handler := NewActivityHandler(log, svcMock, nil)
		handler.GetActivity(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Activity not found")
	})

	t.Run("it handles service error correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/activities/999", nil)
		ctx.Params = gin.Params{{Key: "id", Value: "999"}}

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().GetActivityByID(gomock.Any(), uint(999)).Return(nil, fmt.Errorf("activity not found"))

		handler := NewActivityHandler(log, svcMock, nil)
		handler.GetActivity(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Activity not found")
	})
}

func TestActivityHandler_DeleteActivity_EdgeCases(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it returns error for invalid ID parameter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/activities/invalid", nil)
		ctx.Params = gin.Params{{Key: "id", Value: "invalid"}}

		svcMock := service.NewMockActivityService(ctrl)
		handler := NewActivityHandler(log, svcMock, nil)
		handler.DeleteActivity(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid activity ID")
	})

	t.Run("it returns error for zero ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/activities/0", nil)
		ctx.Params = gin.Params{{Key: "id", Value: "0"}}

		svcMock := service.NewMockActivityService(ctrl)
		// The handler will call the service with ID 0, so we need to expect it
		svcMock.EXPECT().DeleteActivity(gomock.Any(), uint(0)).Return(fmt.Errorf("activity not found"))

		handler := NewActivityHandler(log, svcMock, nil)
		handler.DeleteActivity(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Activity not found")
	})

	t.Run("it handles service error correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/activities/999", nil)
		ctx.Params = gin.Params{{Key: "id", Value: "999"}}

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().DeleteActivity(gomock.Any(), uint(999)).Return(fmt.Errorf("activity not found"))

		handler := NewActivityHandler(log, svcMock, nil)
		handler.DeleteActivity(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Activity not found")
	})
}

func TestActivityHandler_GetActivityStats_EdgeCases(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it handles service error correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/activities/stats", nil)

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().GetActivityStats(gomock.Any()).Return(nil, fmt.Errorf("database error"))

		handler := NewActivityHandler(log, svcMock, nil)
		handler.GetActivityStats(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})
}

func TestActivityHandler_ResetAllData_EdgeCases(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it handles service error correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/activities/reset", nil)

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().ResetAllData(gomock.Any()).Return(fmt.Errorf("database error"))

		handler := NewActivityHandler(log, svcMock, nil)
		handler.ResetAllData(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "database error")
	})
}

func TestActivityHandler_GetActivityStats(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it returns an error when service fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/activities/stats", nil)

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().GetActivityStats(gomock.Any()).Return(nil, fmt.Errorf("database error"))

		handler := NewActivityHandler(log, svcMock, nil)
		handler.GetActivityStats(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "{\"details\":\"database error\",\"error\":\"Failed to get activity stats\"}")
	})

	t.Run("it successfully retrieves activity stats", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/activities/stats", nil)

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().GetActivityStats(gomock.Any()).Return(&activityV1.ActivityStatsResponse{
			TotalActivities:      1,
			RecentActivities:     []activityV1.ActivityResponse{{ID: 5}},
			ActivitiesLast24h:    1,
			ActivitiesLast7Days:  2,
			ActivitiesLast30Days: 3,
		}, nil)

		handler := NewActivityHandler(log, svcMock, nil)
		handler.GetActivityStats(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		// Check that the response contains the expected data (without strict JSON format matching)
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "\"totalActivities\":1")
		assert.Contains(t, responseBody, "\"id\":5")
		assert.Contains(t, responseBody, "\"activitiesLast24h\":1")
		assert.Contains(t, responseBody, "\"activitiesLast7Days\":2")
		assert.Contains(t, responseBody, "\"activitiesLast30Days\":3")
	})
}

func TestActivityHandler_DeleteActivity(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it returns an error when ID is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/activities/invalid", nil)

		ctx.Params = []gin.Param{{Key: "id", Value: "invalid"}}
		handler := NewActivityHandler(log, nil, nil)
		handler.DeleteActivity(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Invalid activity ID\"}")
	})

	t.Run("it returns an error when activity is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/activities/1", nil)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().DeleteActivity(gomock.Any(), uint(1)).Return(fmt.Errorf("not found"))

		handler := NewActivityHandler(log, svcMock, nil)
		handler.DeleteActivity(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Activity not found\"}")
	})

	t.Run("it successfully deletes activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/activities/1", nil)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().DeleteActivity(gomock.Any(), uint(1)).Return(nil)

		handler := NewActivityHandler(log, svcMock, nil)
		handler.DeleteActivity(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "{\"message\":\"Activity deleted successfully\"}")
	})
}

func TestActivityHandler_ResetAllData(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it returns an error when service fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/activities/reset", nil)

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().ResetAllData(gomock.Any()).Return(fmt.Errorf("database error"))

		handler := NewActivityHandler(log, svcMock, nil)
		handler.ResetAllData(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "{\"details\":\"database error\",\"error\":\"Failed to reset activity data\"}")
	})

	t.Run("it successfully resets all activity data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest(http.MethodDelete, "/activities/reset", nil)

		svcMock := service.NewMockActivityService(ctrl)
		svcMock.EXPECT().ResetAllData(gomock.Any()).Return(nil)

		handler := NewActivityHandler(log, svcMock, nil)
		handler.ResetAllData(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "{\"message\":\"All activity data reset successfully\"}")
	})
}
