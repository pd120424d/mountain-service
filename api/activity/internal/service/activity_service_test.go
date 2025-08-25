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

func TestActivityService_CreateActivity_Comprehensive(t *testing.T) {
	t.Parallel()

	t.Run("successfully creates activity with all fields", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityCreateRequest{
			Description: "User completed onboarding process",
			EmployeeID:  123,
			UrgencyID:   456,
		}

		mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(activity *model.Activity) error {
			assert.Equal(t, "User completed onboarding process", activity.Description)
			assert.Equal(t, uint(123), activity.EmployeeID)
			assert.Equal(t, uint(456), activity.UrgencyID)
			// Set ID for response
			activity.ID = 1
			activity.CreatedAt = time.Now()
			activity.UpdatedAt = time.Now()
			return nil
		})

		response, err := service.CreateActivity(req)
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "User completed onboarding process", response.Description)
		assert.Equal(t, uint(123), response.EmployeeID)
		assert.Equal(t, uint(456), response.UrgencyID)
		assert.NotEmpty(t, response.CreatedAt)
		assert.NotEmpty(t, response.UpdatedAt)
	})

	t.Run("returns error when description is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityCreateRequest{
			Description: "", // Empty description
			EmployeeID:  1,
			UrgencyID:   2,
		}

		response, err := service.CreateActivity(req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("returns error when employee ID is zero", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityCreateRequest{
			Description: "Test activity",
			EmployeeID:  0, // Zero employee ID
			UrgencyID:   2,
		}

		response, err := service.CreateActivity(req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("returns error when urgency ID is zero", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityCreateRequest{
			Description: "Test activity",
			EmployeeID:  1,
			UrgencyID:   0, // Zero urgency ID
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
			Description: "Test activity",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		mockRepo.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("database connection failed"))

		response, err := service.CreateActivity(req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to create activity")
	})

	t.Run("handles nil request gracefully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		response, err := service.CreateActivity(nil)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "request cannot be nil")
	})
}

func TestActivityService_GetActivityByID_Comprehensive(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves activity by ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		expectedActivity := &model.Activity{
			ID:          42,
			Description: "Emergency response completed",
			EmployeeID:  100,
			UrgencyID:   200,
			CreatedAt:   time.Now().Add(-1 * time.Hour),
			UpdatedAt:   time.Now(),
		}

		mockRepo.EXPECT().GetByID(uint(42)).Return(expectedActivity, nil)

		response, err := service.GetActivityByID(42)
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, uint(42), response.ID)
		assert.Equal(t, "Emergency response completed", response.Description)
		assert.Equal(t, uint(100), response.EmployeeID)
		assert.Equal(t, uint(200), response.UrgencyID)
		assert.NotEmpty(t, response.CreatedAt)
		assert.NotEmpty(t, response.UpdatedAt)
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

	t.Run("returns error for zero ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		response, err := service.GetActivityByID(0)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "invalid activity ID")
	})
}

func TestActivityService_ListActivities_Comprehensive(t *testing.T) {
	t.Parallel()

	t.Run("successfully lists activities with filters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		employeeID := uint(10)
		urgencyID := uint(20)
		req := &activityV1.ActivityListRequest{
			EmployeeID: &employeeID,
			UrgencyID:  &urgencyID,
			Page:       1,
			PageSize:   10,
			StartDate:  "2023-01-01T00:00:00Z",
			EndDate:    "2023-12-31T23:59:59Z",
		}

		expectedActivities := []model.Activity{
			{
				ID:          1,
				Description: "First activity",
				EmployeeID:  10,
				UrgencyID:   20,
				CreatedAt:   time.Now().Add(-2 * time.Hour),
				UpdatedAt:   time.Now().Add(-2 * time.Hour),
			},
			{
				ID:          2,
				Description: "Second activity",
				EmployeeID:  10,
				UrgencyID:   20,
				CreatedAt:   time.Now().Add(-1 * time.Hour),
				UpdatedAt:   time.Now().Add(-1 * time.Hour),
			},
		}

		mockRepo.EXPECT().List(gomock.Any()).DoAndReturn(func(filter *model.ActivityFilter) ([]model.Activity, int64, error) {
			assert.Equal(t, &employeeID, filter.EmployeeID)
			assert.Equal(t, &urgencyID, filter.UrgencyID)
			assert.Equal(t, 1, filter.Page)
			assert.Equal(t, 10, filter.PageSize)
			assert.NotNil(t, filter.StartDate)
			assert.NotNil(t, filter.EndDate)
			return expectedActivities, 2, nil
		})

		response, err := service.ListActivities(req)
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Len(t, response.Activities, 2)
		assert.Equal(t, int64(2), response.Total)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.PageSize)
		assert.Equal(t, 1, response.TotalPages)
		assert.Equal(t, uint(1), response.Activities[0].ID)
		assert.Equal(t, "First activity", response.Activities[0].Description)
	})

	t.Run("handles empty results", func(t *testing.T) {
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
		assert.Equal(t, 0, response.TotalPages)
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

		mockRepo.EXPECT().List(gomock.Any()).Return(nil, int64(0), fmt.Errorf("database timeout"))

		response, err := service.ListActivities(req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to list activities")
	})

	t.Run("returns error for invalid date formats", func(t *testing.T) {
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

		// Service should return validation error for invalid dates
		response, err := service.ListActivities(req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("validates request parameters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityListRequest{
			Page:     1,
			PageSize: 10,
		}

		mockRepo.EXPECT().List(gomock.Any()).DoAndReturn(func(filter *model.ActivityFilter) ([]model.Activity, int64, error) {
			// Validate that request validation was called
			assert.Equal(t, 1, filter.Page)
			assert.Equal(t, 10, filter.PageSize)
			return []model.Activity{}, 0, nil
		})

		response, err := service.ListActivities(req)
		assert.NoError(t, err)
		require.NotNil(t, response)
	})
}

func TestActivityService_LogActivity_Comprehensive(t *testing.T) {
	t.Parallel()

	t.Run("successfully logs activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		mockRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(activity *model.Activity) error {
			assert.Equal(t, "User logged in successfully", activity.Description)
			assert.Equal(t, uint(50), activity.EmployeeID)
			assert.Equal(t, uint(75), activity.UrgencyID)
			return nil
		})

		err := service.LogActivity("User logged in successfully", 50, 75)
		assert.NoError(t, err)
	})

	t.Run("returns error for empty description", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		err := service.LogActivity("", 1, 2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("handles repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		mockRepo.EXPECT().Create(gomock.Any()).Return(fmt.Errorf("disk full"))

		err := service.LogActivity("Test activity", 1, 2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create activity")
	})
}

func TestActivityService_CreateActivity(t *testing.T) {
	t.Parallel()

	t.Run("successfully creates activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityCreateRequest{
			Description: "New employee was created",
			EmployeeID:  1,
			UrgencyID:   2,
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
		assert.Equal(t, "New employee was created", response.Description)
		assert.Equal(t, uint(1), response.EmployeeID)
		assert.Equal(t, uint(2), response.UrgencyID)
	})

	t.Run("returns error when validation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo)

		req := &activityV1.ActivityCreateRequest{
			Description: "", // Invalid - empty description
			EmployeeID:  1,
			UrgencyID:   2,
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
			Description: "Test",
			EmployeeID:  1,
			UrgencyID:   2,
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

		activity := &model.Activity{
			ID:          1,
			Description: "New employee was created",
			EmployeeID:  1,
			UrgencyID:   2,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockRepo.EXPECT().GetByID(uint(1)).Return(activity, nil)

		response, err := service.GetActivityByID(1)
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "New employee was created", response.Description)
		assert.Equal(t, uint(1), response.EmployeeID)
		assert.Equal(t, uint(2), response.UrgencyID)
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
