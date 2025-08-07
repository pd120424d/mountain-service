package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/pd120424d/mountain-service/api/activity/internal"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

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

		handler := NewActivityHandler(log, nil)
		handler.CreateActivity(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Invalid request payload: invalid character '\\\\n' in string literal\"}")
	})

	t.Run("it returns an error when service fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		validPayload := `{
			"type": "employee_created",
			"level": "info",
			"title": "Test",
			"description": "Test"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/activities", strings.NewReader(validPayload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		svcMock := internal.NewMockActivityService(ctrl)
		svcMock.EXPECT().CreateActivity(gomock.Any()).Return(nil, fmt.Errorf("database error"))

		handler := NewActivityHandler(log, svcMock)
		handler.CreateActivity(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "{\"details\":\"database error\",\"error\":\"Failed to create activity\"}")
	})

	t.Run("it successfully creates activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		validPayload := `{
			"type": "employee_created",
			"level": "info",
			"title": "Test",
			"description": "Test"
		}`
		ctx.Request = httptest.NewRequest(http.MethodPost, "/activities", strings.NewReader(validPayload))
		ctx.Request.Header.Set("Content-Type", "application/json")

		svcMock := internal.NewMockActivityService(ctrl)
		svcMock.EXPECT().CreateActivity(gomock.Any()).Return(&activityV1.ActivityResponse{ID: 1}, nil)

		handler := NewActivityHandler(log, svcMock)
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

		ctx.Params = []gin.Param{{Key: "id", Value: "invalid"}}
		handler := NewActivityHandler(log, nil)
		handler.GetActivity(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Invalid activity ID\"}")
	})

	t.Run("it returns an error when activity is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		svcMock := internal.NewMockActivityService(ctrl)
		svcMock.EXPECT().GetActivityByID(uint(1)).Return(nil, fmt.Errorf("not found"))

		handler := NewActivityHandler(log, svcMock)
		handler.GetActivity(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Activity not found\"}")
	})

	t.Run("it successfully retrieves activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		svcMock := internal.NewMockActivityService(ctrl)
		svcMock.EXPECT().GetActivityByID(uint(1)).Return(&activityV1.ActivityResponse{ID: 1}, nil)

		handler := NewActivityHandler(log, svcMock)
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

		svcMock := internal.NewMockActivityService(ctrl)
		svcMock.EXPECT().ListActivities(gomock.Any()).Return(nil, fmt.Errorf("database error"))

		handler := NewActivityHandler(log, svcMock)
		handler.ListActivities(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "{\"details\":\"database error\",\"error\":\"Failed to list activities\"}")
	})

	t.Run("it successfully retrieves activities", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		svcMock := internal.NewMockActivityService(ctrl)
		svcMock.EXPECT().ListActivities(gomock.Any()).Return(&activityV1.ActivityListResponse{
			Activities: []activityV1.ActivityResponse{{ID: 1}},
			Total:      1,
			Page:       1,
			PageSize:   10,
		}, nil)

		handler := NewActivityHandler(log, svcMock)
		handler.ListActivities(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "{\"activities\":[{\"id\":1")
	})
}

func TestActivityHandler_GetActivityStats(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it returns an error when service fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		svcMock := internal.NewMockActivityService(ctrl)
		svcMock.EXPECT().GetActivityStats().Return(nil, fmt.Errorf("database error"))

		handler := NewActivityHandler(log, svcMock)
		handler.GetActivityStats(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "{\"details\":\"database error\",\"error\":\"Failed to get activity stats\"}")
	})

	t.Run("it successfully retrieves activity stats", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		svcMock := internal.NewMockActivityService(ctrl)
		svcMock.EXPECT().GetActivityStats().Return(&activityV1.ActivityStatsResponse{
			TotalActivities: 1,
			ActivitiesByType: map[activityV1.ActivityType]int64{
				activityV1.ActivityEmployeeCreated: 6,
			},
			ActivitiesByLevel: map[activityV1.ActivityLevel]int64{
				activityV1.ActivityLevelInfo: 4,
			},
			RecentActivities:     []activityV1.ActivityResponse{{ID: 5}},
			ActivitiesLast24h:    1,
			ActivitiesLast7Days:  2,
			ActivitiesLast30Days: 3,
		}, nil)

		handler := NewActivityHandler(log, svcMock)
		handler.GetActivityStats(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		// Check that the response contains the expected data (without strict JSON format matching)
		responseBody := w.Body.String()
		assert.Contains(t, responseBody, "\"totalActivities\":1")
		assert.Contains(t, responseBody, "\"employee_created\":6")
		assert.Contains(t, responseBody, "\"info\":4")
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

		ctx.Params = []gin.Param{{Key: "id", Value: "invalid"}}
		handler := NewActivityHandler(log, nil)
		handler.DeleteActivity(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Invalid activity ID\"}")
	})

	t.Run("it returns an error when activity is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		svcMock := internal.NewMockActivityService(ctrl)
		svcMock.EXPECT().DeleteActivity(uint(1)).Return(fmt.Errorf("not found"))

		handler := NewActivityHandler(log, svcMock)
		handler.DeleteActivity(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "{\"error\":\"Activity not found\"}")
	})

	t.Run("it successfully deletes activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Params = []gin.Param{{Key: "id", Value: "1"}}

		svcMock := internal.NewMockActivityService(ctrl)
		svcMock.EXPECT().DeleteActivity(uint(1)).Return(nil)

		handler := NewActivityHandler(log, svcMock)
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

		svcMock := internal.NewMockActivityService(ctrl)
		svcMock.EXPECT().ResetAllData().Return(fmt.Errorf("database error"))

		handler := NewActivityHandler(log, svcMock)
		handler.ResetAllData(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "{\"details\":\"database error\",\"error\":\"Failed to reset activity data\"}")
	})

	t.Run("it successfully resets all activity data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		svcMock := internal.NewMockActivityService(ctrl)
		svcMock.EXPECT().ResetAllData().Return(nil)

		handler := NewActivityHandler(log, svcMock)
		handler.ResetAllData(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "{\"message\":\"All activity data reset successfully\"}")
	})
}
