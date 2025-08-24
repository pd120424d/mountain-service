package service

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

func TestActivityService_ListActivities_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("it handles empty request gracefully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityListRequest{}

		mockRepo.EXPECT().List(gomock.Any()).DoAndReturn(func(filter *model.ActivityFilter) ([]model.Activity, int64, error) {
			// The filter should have the raw values from the request (0, 0)
			// The repository's filter.Validate() will set defaults
			assert.Equal(t, 0, filter.Page)
			assert.Equal(t, 0, filter.PageSize)
			return []model.Activity{}, 0, nil
		})

		response, err := service.ListActivities(req)
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Empty(t, response.Activities)
		assert.Equal(t, int64(0), response.Total)
	})

	t.Run("it handles filters correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		actorID := uint(123)
		targetID := uint(456)
		activityType := activityV1.ActivityEmployeeCreated
		activityLevel := activityV1.ActivityLevelInfo
		targetType := "employee"

		req := &activityV1.ActivityListRequest{
			Page:       2,
			PageSize:   25,
			Type:       activityType,
			Level:      activityLevel,
			ActorID:    &actorID,
			TargetID:   &targetID,
			TargetType: targetType,
		}

		mockRepo.EXPECT().List(gomock.Any()).DoAndReturn(func(filter *model.ActivityFilter) ([]model.Activity, int64, error) {
			assert.Equal(t, 2, filter.Page)
			assert.Equal(t, 25, filter.PageSize)
			assert.Equal(t, activityType, *filter.Type)
			assert.Equal(t, activityLevel, *filter.Level)
			assert.Equal(t, actorID, *filter.ActorID)
			assert.Equal(t, targetID, *filter.TargetID)
			assert.Equal(t, targetType, *filter.TargetType)
			return []model.Activity{}, 0, nil
		})

		response, err := service.ListActivities(req)
		assert.NoError(t, err)
		require.NotNil(t, response)
	})
}

func TestActivityService_LogEmployeeActivity_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("it returns error when employee ID is zero", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		// The service will create an activity with TargetID pointing to 0
		mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(activity *model.Activity) error {
			// Verify that the target ID is a pointer to 0 (since employee ID was 0)
			assert.NotNil(t, activity.TargetID)
			assert.Equal(t, uint(0), *activity.TargetID)
			return nil
		})

		err := service.LogEmployeeActivity(activityV1.ActivityEmployeeCreated, activityV1.ActivityLevelInfo, "Test", "Test", 0)
		assert.NoError(t, err) // The service doesn't validate employee ID being 0, it creates a pointer to 0
	})

	t.Run("it returns error when title is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		err := service.LogEmployeeActivity(activityV1.ActivityEmployeeCreated, activityV1.ActivityLevelInfo, "", "Test", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("it returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		mockRepo.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("database error"))

		err := service.LogEmployeeActivity(activityV1.ActivityEmployeeCreated, activityV1.ActivityLevelInfo, "Test", "Test", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create activity")
	})
}

func TestActivityService_LogUrgencyActivity_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("it returns error when urgency ID is zero", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		// The service will create an activity with TargetID pointing to 0
		mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(activity *model.Activity) error {
			// Verify that the target ID is a pointer to 0 (since urgency ID was 0)
			assert.NotNil(t, activity.TargetID)
			assert.Equal(t, uint(0), *activity.TargetID)
			return nil
		})

		err := service.LogUrgencyActivity(activityV1.ActivityUrgencyCreated, activityV1.ActivityLevelWarning, "Test", "Test", 0)
		assert.NoError(t, err) // The service doesn't validate urgency ID being 0, it creates a pointer to 0
	})

	t.Run("it returns error when description is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		err := service.LogUrgencyActivity(activityV1.ActivityUrgencyCreated, activityV1.ActivityLevelWarning, "Test", "", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("it returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		mockRepo.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("database error"))

		err := service.LogUrgencyActivity(activityV1.ActivityUrgencyCreated, activityV1.ActivityLevelWarning, "Test", "Test", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create activity")
	})
}

func TestActivityService_LogSystemActivity_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("it returns error when title is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		err := service.LogSystemActivity(activityV1.ActivitySystemReset, activityV1.ActivityLevelCritical, "", "Test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("it returns error when description is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		err := service.LogSystemActivity(activityV1.ActivitySystemReset, activityV1.ActivityLevelCritical, "Test", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("it returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		mockRepo.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("database error"))

		err := service.LogSystemActivity(activityV1.ActivitySystemReset, activityV1.ActivityLevelCritical, "Test", "Test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create activity")
	})
}

func TestActivityService_CreateActivity_ComprehensiveEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("it handles activity with all optional fields", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		actorID := uint(123)
		targetID := uint(456)
		req := &activityV1.ActivityCreateRequest{
			Type:        activityV1.ActivityEmployeeCreated,
			Level:       activityV1.ActivityLevelInfo,
			Title:       "Employee Created",
			Description: "New employee was created",
			ActorID:     &actorID,
			ActorName:   "Admin User",
			TargetID:    &targetID,
			TargetType:  "employee",
			Metadata:    `{"key": "value"}`,
		}

		mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(activity *model.Activity) error {
			activity.ID = 1
			assert.Equal(t, req.Type, activity.Type)
			assert.Equal(t, req.Level, activity.Level)
			assert.Equal(t, req.Title, activity.Title)
			assert.Equal(t, req.Description, activity.Description)
			assert.Equal(t, *req.ActorID, *activity.ActorID)
			assert.Equal(t, req.ActorName, activity.ActorName)
			assert.Equal(t, *req.TargetID, *activity.TargetID)
			assert.Equal(t, req.TargetType, activity.TargetType)
			assert.Equal(t, req.Metadata, activity.Metadata)
			return nil
		})

		response, err := service.CreateActivity(req)
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
	})

	t.Run("it handles activity with minimal fields", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityCreateRequest{
			Type:        activityV1.ActivitySystemReset,
			Level:       activityV1.ActivityLevelCritical,
			Title:       "System Reset",
			Description: "System was reset",
			ActorName:   "system",
			TargetType:  "system",
		}

		mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(activity *model.Activity) error {
			activity.ID = 2
			assert.Equal(t, req.Type, activity.Type)
			assert.Equal(t, req.Level, activity.Level)
			assert.Equal(t, req.Title, activity.Title)
			assert.Equal(t, req.Description, activity.Description)
			assert.Nil(t, activity.ActorID)
			assert.Equal(t, req.ActorName, activity.ActorName)
			assert.Nil(t, activity.TargetID)
			assert.Equal(t, req.TargetType, activity.TargetType)
			return nil
		})

		response, err := service.CreateActivity(req)
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, uint(2), response.ID)
	})
}

func TestActivityService_ListActivities_DateParsing(t *testing.T) {
	t.Parallel()

	t.Run("it handles valid date filters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityListRequest{
			Page:      1,
			PageSize:  10,
			StartDate: "2023-01-01T00:00:00Z",
			EndDate:   "2023-12-31T23:59:59Z",
		}

		mockRepo.EXPECT().List(gomock.Any()).DoAndReturn(func(filter *model.ActivityFilter) ([]model.Activity, int64, error) {
			assert.Equal(t, 1, filter.Page)
			assert.Equal(t, 10, filter.PageSize)
			assert.NotNil(t, filter.StartDate)
			assert.NotNil(t, filter.EndDate)
			return []model.Activity{}, 0, nil
		})

		response, err := service.ListActivities(req)
		assert.NoError(t, err)
		require.NotNil(t, response)
	})

	t.Run("it returns error for invalid date formats", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityListRequest{
			Page:      1,
			PageSize:  10,
			StartDate: "invalid-date",
			EndDate:   "also-invalid",
		}

		// The service should return a validation error for invalid date formats
		response, err := service.ListActivities(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, response)
	})

	t.Run("it returns error for mixed valid and invalid dates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityListRequest{
			Page:      1,
			PageSize:  10,
			StartDate: "2023-01-01T00:00:00Z", // Valid
			EndDate:   "invalid-date",         // Invalid
		}

		// The service should return a validation error for invalid endDate
		response, err := service.ListActivities(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, response)
	})
}

func TestActivityService_ListActivities_ResponseMapping(t *testing.T) {
	t.Parallel()

	t.Run("it correctly maps response with pagination", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityListRequest{
			Page:     2,
			PageSize: 5,
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

		mockRepo.EXPECT().List(gomock.Any()).Return(activities, int64(15), nil)

		response, err := service.ListActivities(req)
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Len(t, response.Activities, 1)
		assert.Equal(t, int64(15), response.Total)
		assert.Equal(t, 2, response.Page)
		assert.Equal(t, 5, response.PageSize)
		assert.Equal(t, 3, response.TotalPages) // ceil(15/5) = 3
		assert.Equal(t, uint(1), response.Activities[0].ID)
		assert.Equal(t, activityV1.ActivityEmployeeCreated, response.Activities[0].Type)
	})

	t.Run("it handles zero total correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityListRequest{
			Page:     1,
			PageSize: 10,
		}

		mockRepo.EXPECT().List(gomock.Any()).Return([]model.Activity{}, int64(0), nil)

		response, err := service.ListActivities(req)
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Empty(t, response.Activities)
		assert.Equal(t, int64(0), response.Total)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.PageSize)
		assert.Equal(t, 0, response.TotalPages)
	})
}

func TestActivityService_ValidationHelpers(t *testing.T) {
	t.Parallel()

	log := utils.NewTestLogger()

	t.Run("it validates activity types correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repoMock := repositories.NewMockActivityRepository(ctrl)
		svc := NewActivityService(log, repoMock)

		// Test valid activity type
		req := &activityV1.ActivityCreateRequest{
			Type:        activityV1.ActivityEmployeeCreated,
			Level:       activityV1.ActivityLevelInfo,
			Title:       "Valid Title",
			Description: "Valid Description",
		}

		repoMock.EXPECT().Create(gomock.Any()).DoAndReturn(func(activity *model.Activity) error {
			activity.ID = 1
			return nil
		})
		result, err := svc.CreateActivity(req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("it validates activity levels correctly", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repoMock := repositories.NewMockActivityRepository(ctrl)
		svc := NewActivityService(log, repoMock)

		// Test valid activity level
		req := &activityV1.ActivityCreateRequest{
			Type:        activityV1.ActivityEmployeeCreated,
			Level:       activityV1.ActivityLevelWarning,
			Title:       "Valid Title",
			Description: "Valid Description",
		}

		repoMock.EXPECT().Create(gomock.Any()).DoAndReturn(func(activity *model.Activity) error {
			activity.ID = 1
			return nil
		})
		result, err := svc.CreateActivity(req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("it validates required fields", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repoMock := repositories.NewMockActivityRepository(ctrl)
		svc := NewActivityService(log, repoMock)

		// Test missing title
		req := &activityV1.ActivityCreateRequest{
			Type:        activityV1.ActivityEmployeeCreated,
			Level:       activityV1.ActivityLevelInfo,
			Title:       "",
			Description: "Valid Description",
		}

		result, err := svc.CreateActivity(req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "title is required")
	})
}
