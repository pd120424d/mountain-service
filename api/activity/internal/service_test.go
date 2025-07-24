package internal

import (
	"fmt"
	"testing"
	"time"

	"github.com/pd120424d/mountain-service/api/activity/internal/model"
	"github.com/pd120424d/mountain-service/api/activity/internal/repositories"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestActivityService_CreateActivity(t *testing.T) {
	t.Parallel()

	t.Run("successfully creates activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		actorID := uint(1)
		targetID := uint(2)
		req := &activityV1.ActivityCreateRequest{
			Type:        activityV1.ActivityEmployeeCreated,
			Level:       activityV1.ActivityLevelInfo,
			Title:       "Employee Created",
			Description: "New employee was created",
			ActorID:     &actorID,
			ActorName:   "Admin",
			TargetID:    &targetID,
			TargetType:  "employee",
			Metadata:    "{}",
		}

		mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(activity *model.Activity) error {
			activity.ID = 1
			activity.CreatedAt = time.Now()
			activity.UpdatedAt = time.Now()
			return nil
		})

		response, err := service.CreateActivity(req)
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, activityV1.ActivityEmployeeCreated, response.Type)
		assert.Equal(t, activityV1.ActivityLevelInfo, response.Level)
		assert.Equal(t, "Employee Created", response.Title)
		assert.Equal(t, "New employee was created", response.Description)
		assert.Equal(t, &actorID, response.ActorID)
		assert.Equal(t, "Admin", response.ActorName)
		assert.Equal(t, &targetID, response.TargetID)
		assert.Equal(t, "employee", response.TargetType)
	})

	t.Run("returns error when validation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityCreateRequest{
			Type:        "", // Invalid - empty type
			Level:       activityV1.ActivityLevelInfo,
			Title:       "Test",
			Description: "Test",
		}

		response, err := service.CreateActivity(req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityCreateRequest{
			Type:        activityV1.ActivityEmployeeCreated,
			Level:       activityV1.ActivityLevelInfo,
			Title:       "Test",
			Description: "Test",
		}

		mockRepo.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("database error"))

		response, err := service.CreateActivity(req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to create activity")
	})
}

func TestActivityService_GetActivityByID(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves activity by ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		actorID := uint(1)
		targetID := uint(2)
		activity := &model.Activity{
			ID:          1,
			Type:        activityV1.ActivityEmployeeCreated,
			Level:       activityV1.ActivityLevelInfo,
			Title:       "Employee Created",
			Description: "New employee was created",
			ActorID:     &actorID,
			ActorName:   "Admin",
			TargetID:    &targetID,
			TargetType:  "employee",
			Metadata:    "{}",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockRepo.EXPECT().GetByID(uint(1)).Return(activity, nil)

		response, err := service.GetActivityByID(1)
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, activityV1.ActivityEmployeeCreated, response.Type)
		assert.Equal(t, activityV1.ActivityLevelInfo, response.Level)
	})

	t.Run("returns error when activity not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		mockRepo.EXPECT().GetByID(uint(999)).Return(nil, fmt.Errorf("activity not found"))

		response, err := service.GetActivityByID(999)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to get activity")
	})
}

func TestActivityService_DeleteActivity(t *testing.T) {
	t.Parallel()

	t.Run("successfully deletes activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		mockRepo.EXPECT().Delete(uint(1)).Return(nil)

		err := service.DeleteActivity(1)
		assert.NoError(t, err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		mockRepo.EXPECT().Delete(uint(999)).Return(fmt.Errorf("activity not found"))

		err := service.DeleteActivity(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to delete activity")
	})
}

func TestActivityService_ResetAllData(t *testing.T) {
	t.Parallel()

	t.Run("successfully resets all data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		mockRepo.EXPECT().ResetAllData().Return(nil)

		err := service.ResetAllData()
		assert.NoError(t, err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		mockRepo.EXPECT().ResetAllData().Return(fmt.Errorf("database error"))

		err := service.ResetAllData()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to reset activity data")
	})
}

func TestActivityService_LogEmployeeActivity(t *testing.T) {
	t.Parallel()

	t.Run("successfully logs employee activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(activity *model.Activity) error {
			assert.Equal(t, activityV1.ActivityEmployeeCreated, activity.Type)
			assert.Equal(t, activityV1.ActivityLevelInfo, activity.Level)
			assert.Equal(t, "Employee Created", activity.Title)
			assert.Equal(t, "New employee was created", activity.Description)
			assert.Equal(t, uint(1), *activity.TargetID)
			assert.Equal(t, "employee", activity.TargetType)
			return nil
		})

		err := service.LogEmployeeActivity(
			activityV1.ActivityEmployeeCreated,
			activityV1.ActivityLevelInfo,
			"Employee Created",
			"New employee was created",
			1,
		)
		assert.NoError(t, err)
	})
}

func TestActivityService_LogUrgencyActivity(t *testing.T) {
	t.Parallel()

	t.Run("successfully logs urgency activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(activity *model.Activity) error {
			assert.Equal(t, activityV1.ActivityUrgencyCreated, activity.Type)
			assert.Equal(t, activityV1.ActivityLevelWarning, activity.Level)
			assert.Equal(t, "Urgency Created", activity.Title)
			assert.Equal(t, "New urgency was created", activity.Description)
			assert.Equal(t, uint(1), *activity.TargetID)
			assert.Equal(t, "urgency", activity.TargetType)
			return nil
		})

		err := service.LogUrgencyActivity(
			activityV1.ActivityUrgencyCreated,
			activityV1.ActivityLevelWarning,
			"Urgency Created",
			"New urgency was created",
			1,
		)
		assert.NoError(t, err)
	})
}

func TestActivityService_LogSystemActivity(t *testing.T) {
	t.Parallel()

	t.Run("successfully logs system activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(activity *model.Activity) error {
			assert.Equal(t, activityV1.ActivitySystemReset, activity.Type)
			assert.Equal(t, activityV1.ActivityLevelCritical, activity.Level)
			assert.Equal(t, "System Reset", activity.Title)
			assert.Equal(t, "System was reset", activity.Description)
			assert.Nil(t, activity.TargetID)
			assert.Equal(t, "system", activity.TargetType)
			assert.Equal(t, "system", activity.ActorName)
			return nil
		})

		err := service.LogSystemActivity(
			activityV1.ActivitySystemReset,
			activityV1.ActivityLevelCritical,
			"System Reset",
			"System was reset",
		)
		assert.NoError(t, err)
	})
}

func TestActivityService_ListActivities(t *testing.T) {
	t.Parallel()

	t.Run("successfully lists activities", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityListRequest{
			Page:     1,
			PageSize: 10,
		}

		actorID := uint(1)
		targetID := uint(2)
		activities := []model.Activity{
			{
				ID:          1,
				Type:        activityV1.ActivityEmployeeCreated,
				Level:       activityV1.ActivityLevelInfo,
				Title:       "Employee Created",
				Description: "New employee was created",
				ActorID:     &actorID,
				ActorName:   "Admin",
				TargetID:    &targetID,
				TargetType:  "employee",
				Metadata:    "{}",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		mockRepo.EXPECT().List(gomock.Any()).DoAndReturn(func(filter *model.ActivityFilter) ([]model.Activity, int64, error) {
			assert.Equal(t, 1, filter.Page)
			assert.Equal(t, 10, filter.PageSize)
			return activities, 1, nil
		})

		response, err := service.ListActivities(req)
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Len(t, response.Activities, 1)
		assert.Equal(t, int64(1), response.Total)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.PageSize)
		assert.Equal(t, uint(1), response.Activities[0].ID)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityListRequest{
			Page:     1,
			PageSize: 10,
		}

		mockRepo.EXPECT().List(gomock.Any()).Return(nil, int64(0), fmt.Errorf("database error"))

		response, err := service.ListActivities(req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to list activities")
	})
}

func TestActivityService_GetActivityStats(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves activity stats", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		stats := &model.ActivityStats{
			TotalActivities: 10,
			ActivitiesByType: map[activityV1.ActivityType]int64{
				activityV1.ActivityEmployeeCreated: 5,
				activityV1.ActivityUrgencyCreated:  3,
			},
			ActivitiesByLevel: map[activityV1.ActivityLevel]int64{
				activityV1.ActivityLevelInfo:    6,
				activityV1.ActivityLevelWarning: 4,
			},
			RecentActivities: []model.Activity{
				{
					ID:          1,
					Type:        activityV1.ActivityEmployeeCreated,
					Level:       activityV1.ActivityLevelInfo,
					Title:       "Recent Activity",
					Description: "Recent activity description",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			ActivitiesLast24h:    5,
			ActivitiesLast7Days:  8,
			ActivitiesLast30Days: 10,
		}

		mockRepo.EXPECT().GetStats().Return(stats, nil)

		response, err := service.GetActivityStats()
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, int64(10), response.TotalActivities)
		assert.Equal(t, int64(5), response.ActivitiesByType[activityV1.ActivityEmployeeCreated])
		assert.Equal(t, int64(3), response.ActivitiesByType[activityV1.ActivityUrgencyCreated])
		assert.Equal(t, int64(6), response.ActivitiesByLevel[activityV1.ActivityLevelInfo])
		assert.Equal(t, int64(4), response.ActivitiesByLevel[activityV1.ActivityLevelWarning])
		assert.Len(t, response.RecentActivities, 1)
		assert.Equal(t, uint(1), response.RecentActivities[0].ID)
		assert.Equal(t, int64(5), response.ActivitiesLast24h)
		assert.Equal(t, int64(8), response.ActivitiesLast7Days)
		assert.Equal(t, int64(10), response.ActivitiesLast30Days)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		mockRepo.EXPECT().GetStats().Return(nil, fmt.Errorf("database error"))

		response, err := service.GetActivityStats()
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to get activity stats")
	})
}
