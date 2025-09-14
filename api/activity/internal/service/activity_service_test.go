package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/pd120424d/mountain-service/api/activity/internal/clients"

	"github.com/pd120424d/mountain-service/api/activity/internal/model"
	"github.com/pd120424d/mountain-service/api/activity/internal/repositories"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	commonv1 "github.com/pd120424d/mountain-service/api/contracts/common/v1"
	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
	"github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// test fakes for EmployeeClient
type fakeEmployeeClient struct {
	first, last string
	err         error
}

func (f *fakeEmployeeClient) GetEmployeeByID(ctx context.Context, employeeID uint) (*employeeV1.EmployeeResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &employeeV1.EmployeeResponse{ID: employeeID, FirstName: f.first, LastName: f.last}, nil
}

func TestActivityService_CreateActivity(t *testing.T) {
	t.Parallel()

	t.Run("it successfully creates activity with all fields", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		req := &activityV1.ActivityCreateRequest{
			Description: "User completed onboarding process",
			EmployeeID:  123,
			UrgencyID:   456,
		}

		mockRepo.EXPECT().CreateWithOutbox(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, activity *model.Activity, _ any) error {
			assert.Equal(t, "User completed onboarding process", activity.Description)
			assert.Equal(t, uint(123), activity.EmployeeID)
			assert.Equal(t, uint(456), activity.UrgencyID)
			// Set ID for response
			activity.ID = 1
			activity.CreatedAt = time.Now()
			activity.UpdatedAt = time.Now()
			return nil
		})

		response, err := service.CreateActivity(t.Context(), req)
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "User completed onboarding process", response.Description)
		assert.Equal(t, uint(123), response.EmployeeID)
		assert.Equal(t, uint(456), response.UrgencyID)
		assert.NotEmpty(t, response.CreatedAt)
		assert.NotEmpty(t, response.UpdatedAt)
	})

	t.Run("it returns error when description is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		req := &activityV1.ActivityCreateRequest{
			Description: "",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		response, err := service.CreateActivity(t.Context(), req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("it returns error when employee ID is zero", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		req := &activityV1.ActivityCreateRequest{
			Description: "Test activity",
			EmployeeID:  0,
			UrgencyID:   2,
		}

		response, err := service.CreateActivity(t.Context(), req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("it returns error when urgency ID is zero", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		req := &activityV1.ActivityCreateRequest{
			Description: "Test activity",
			EmployeeID:  1,
			UrgencyID:   0,
		}

		response, err := service.CreateActivity(t.Context(), req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("it returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		req := &activityV1.ActivityCreateRequest{
			Description: "Test activity",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		mockRepo.EXPECT().CreateWithOutbox(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("database connection failed"))

		response, err := service.CreateActivity(t.Context(), req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to create activity")
	})

	t.Run("it handles nil request gracefully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		response, err := service.CreateActivity(t.Context(), nil)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "request cannot be nil")
	})

	t.Run("it successfully creates activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		req := &activityV1.ActivityCreateRequest{
			Description: "New employee was created",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		mockRepo.EXPECT().CreateWithOutbox(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, activity *model.Activity, _ any) error {
			activity.ID = 1
			activity.CreatedAt = time.Now()
			activity.UpdatedAt = time.Now()
			return nil
		})

		response, err := service.CreateActivity(t.Context(), req)
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "New employee was created", response.Description)
		assert.Equal(t, uint(1), response.EmployeeID)
		assert.Equal(t, uint(2), response.UrgencyID)
	})

	t.Run("it returns error when validation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		req := &activityV1.ActivityCreateRequest{
			Description: "",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		response, err := service.CreateActivity(t.Context(), req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("it returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		req := &activityV1.ActivityCreateRequest{
			Description: "Test",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		mockRepo.EXPECT().CreateWithOutbox(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("database error"))

		response, err := service.CreateActivity(t.Context(), req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to create activity")
	})

	t.Run("it denies when urgency has no assignee and actor is not admin", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		repo := repositories.NewMockActivityRepository(ctrl)
		mockUrg := clients.NewMockUrgencyClient(ctrl)
		mockUrg.EXPECT().GetUrgencyByID(gomock.Any(), uint(7)).Return(&urgencyV1.UrgencyResponse{ID: 7, Status: urgencyV1.InProgress}, nil)
		svc := NewActivityService(log, repo, mockUrg)

		req := &activityV1.ActivityCreateRequest{Description: "x", EmployeeID: 999, UrgencyID: 7}
		ctx := context.WithValue(t.Context(), "employeeID", uint(1))
		ctx = context.WithValue(ctx, "role", "Medic")

		resp, err := svc.CreateActivity(ctx, req)
		assert.Nil(t, resp)
		if appErr, ok := err.(*commonv1.AppError); ok {
			assert.Equal(t, "VALIDATION.MISSING_ASSIGNEE", appErr.Code)
		} else {
			t.Fatalf("expected AppError, got %T", err)
		}
	})

	// Simulate context with non-admin actor id 1
	t.Run("denies when actor is not the assignee", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		repo := repositories.NewMockActivityRepository(ctrl)
		assigned := uint(42)
		mockUrg := clients.NewMockUrgencyClient(ctrl)
		mockUrg.EXPECT().GetUrgencyByID(gomock.Any(), uint(7)).Return(&urgencyV1.UrgencyResponse{ID: 7, Status: urgencyV1.InProgress, AssignedEmployeeId: &assigned}, nil)
		svc := NewActivityService(log, repo, mockUrg)

		req := &activityV1.ActivityCreateRequest{Description: "x", EmployeeID: 999, UrgencyID: 7}
		ctx := context.WithValue(t.Context(), "employeeID", uint(1))
		ctx = context.WithValue(ctx, "role", "Medic")

		resp, err := svc.CreateActivity(ctx, req)
		assert.Nil(t, resp)
		if appErr, ok := err.(*commonv1.AppError); ok {
			assert.Equal(t, "AUTH_ERRORS.FORBIDDEN", appErr.Code)
		} else {
			t.Fatalf("expected AppError, got %T", err)
		}
	})

	t.Run("allows when actor is the assignee (and overrides employeeID)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		repo := repositories.NewMockActivityRepository(ctrl)
		assigned := uint(7)
		mockUrg := clients.NewMockUrgencyClient(ctrl)
		mockUrg.EXPECT().GetUrgencyByID(gomock.Any(), uint(9)).Return(&urgencyV1.UrgencyResponse{ID: 9, Status: urgencyV1.InProgress, AssignedEmployeeId: &assigned}, nil)
		svc := NewActivityService(log, repo, mockUrg)

		repo.EXPECT().CreateWithOutbox(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		req := &activityV1.ActivityCreateRequest{Description: "x", EmployeeID: 999, UrgencyID: 9}
		ctx := context.WithValue(t.Context(), "employeeID", uint(7))
		ctx = context.WithValue(ctx, "role", "Medic")

		resp, err := svc.CreateActivity(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, uint(7), resp.EmployeeID)
		assert.Equal(t, uint(9), resp.UrgencyID)
	})

	t.Run("it allows admin regardless of assignee and does not override employeeID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		repo := repositories.NewMockActivityRepository(ctrl)
		// No assignee set
		mockUrg := clients.NewMockUrgencyClient(ctrl)
		mockUrg.EXPECT().GetUrgencyByID(gomock.Any(), uint(11)).Return(&urgencyV1.UrgencyResponse{ID: 11, Status: urgencyV1.InProgress}, nil)
		svc := NewActivityService(log, repo, mockUrg)

		repo.EXPECT().CreateWithOutbox(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		req := &activityV1.ActivityCreateRequest{Description: "x", EmployeeID: 123, UrgencyID: 11}
		ctx := context.WithValue(t.Context(), "employeeID", uint(999))
		ctx = context.WithValue(ctx, "role", "Administrator")

		resp, err := svc.CreateActivity(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		// For admin we keep the request's EmployeeID (actor may be a system user)
		assert.Equal(t, uint(123), resp.EmployeeID)
		assert.Equal(t, uint(11), resp.UrgencyID)
	})

	t.Run("it returns error when urgency client fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		repo := repositories.NewMockActivityRepository(ctrl)
		mockUrg := clients.NewMockUrgencyClient(ctrl)
		mockUrg.EXPECT().GetUrgencyByID(gomock.Any(), uint(77)).Return(nil, fmt.Errorf("boom"))
		svc := NewActivityService(log, repo, mockUrg)

		req := &activityV1.ActivityCreateRequest{Description: "x", EmployeeID: 1, UrgencyID: 77}
		resp, err := svc.CreateActivity(t.Context(), req)
		assert.Nil(t, resp)
		if appErr, ok := err.(*commonv1.AppError); ok {
			assert.Equal(t, "ACTIVITY_ERRORS.URGENCY_FETCH_FAILED", appErr.Code)
		} else {
			t.Fatalf("expected AppError, got %T", err)
		}
	})

	t.Run("it returns error when urgency is nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		repo := repositories.NewMockActivityRepository(ctrl)
		mockUrg := clients.NewMockUrgencyClient(ctrl)
		mockUrg.EXPECT().GetUrgencyByID(gomock.Any(), uint(88)).Return(nil, nil)
		svc := NewActivityService(log, repo, mockUrg)

		req := &activityV1.ActivityCreateRequest{Description: "x", EmployeeID: 1, UrgencyID: 88}
		resp, err := svc.CreateActivity(t.Context(), req)
		assert.Nil(t, resp)
		if appErr, ok := err.(*commonv1.AppError); ok {
			assert.Equal(t, "ACTIVITY_ERRORS.INVALID_URGENCY_STATE", appErr.Code)
		} else {
			t.Fatalf("expected AppError, got %T", err)
		}
	})

	t.Run("it returns error when urgency is not in_progress", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		repo := repositories.NewMockActivityRepository(ctrl)
		mockUrg := clients.NewMockUrgencyClient(ctrl)
		mockUrg.EXPECT().GetUrgencyByID(gomock.Any(), uint(99)).Return(&urgencyV1.UrgencyResponse{ID: 99, Status: urgencyV1.Open}, nil)
		svc := NewActivityService(log, repo, mockUrg)

		req := &activityV1.ActivityCreateRequest{Description: "x", EmployeeID: 1, UrgencyID: 99}
		resp, err := svc.CreateActivity(t.Context(), req)
		assert.Nil(t, resp)
		if appErr, ok := err.(*commonv1.AppError); ok {
			assert.Equal(t, "ACTIVITY_ERRORS.INVALID_URGENCY_STATE", appErr.Code)
		} else {
			t.Fatalf("expected AppError, got %T", err)
		}
	})

}

func TestActivityService_CreateActivity_Enrichment(t *testing.T) {
	t.Parallel()

	t.Run("enriches employeeName and urgency fields in outbox event", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		repo := repositories.NewMockActivityRepository(ctrl)
		assigned := uint(7)
		mockUrg := clients.NewMockUrgencyClient(ctrl)
		mockUrg.EXPECT().GetUrgencyByID(gomock.Any(), uint(9)).Return(&urgencyV1.UrgencyResponse{ID: 9, Status: urgencyV1.InProgress, AssignedEmployeeId: &assigned, FirstName: "Petar", LastName: "Petrovic", Level: urgencyV1.High}, nil)
		// employee client fake that returns a fixed name
		empClient := &fakeEmployeeClient{first: "Mika", last: "Mikic"}
		svc := NewActivityServiceWithDeps(log, repo, mockUrg, empClient)

		repo.EXPECT().CreateWithOutbox(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, activity *model.Activity, ev any) error {
			assert.Equal(t, uint(7), activity.EmployeeID) // overridden by actor
			ob, ok := ev.(*models.OutboxEvent)
			require.True(t, ok)
			var data activityV1.ActivityEvent
			require.NoError(t, json.Unmarshal([]byte(ob.EventData), &data))
			assert.Equal(t, "CREATE", data.Type)
			assert.Equal(t, uint(9), data.UrgencyID)
			assert.Equal(t, uint(7), data.EmployeeID)
			assert.Equal(t, "Mika Mikic", data.EmployeeName)
			assert.Equal(t, "Petar Petrovic", data.UrgencyTitle)
			assert.Equal(t, string(urgencyV1.High), data.UrgencyLevel)
			return nil
		})

		req := &activityV1.ActivityCreateRequest{Description: "note", EmployeeID: 999, UrgencyID: 9}
		ctx := context.WithValue(t.Context(), "employeeID", uint(7))
		ctx = context.WithValue(ctx, "role", "Medic")
		resp, err := svc.CreateActivity(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("continues when employee client fails (no name)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		repo := repositories.NewMockActivityRepository(ctrl)
		assigned := uint(5)
		mockUrg := clients.NewMockUrgencyClient(ctrl)
		mockUrg.EXPECT().GetUrgencyByID(gomock.Any(), uint(3)).Return(&urgencyV1.UrgencyResponse{ID: 3, Status: urgencyV1.InProgress, AssignedEmployeeId: &assigned, FirstName: "Ana", LastName: "Anic", Level: urgencyV1.Medium}, nil)
		// failing employee client
		empClient := &fakeEmployeeClient{err: fmt.Errorf("boom")}
		svc := NewActivityServiceWithDeps(log, repo, mockUrg, empClient)

		repo.EXPECT().CreateWithOutbox(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, activity *model.Activity, ev any) error {
			ob := ev.(*models.OutboxEvent)
			var data activityV1.ActivityEvent
			require.NoError(t, json.Unmarshal([]byte(ob.EventData), &data))
			assert.Equal(t, "", data.EmployeeName) // still proceeds, empty name
			assert.Equal(t, "Ana Anic", data.UrgencyTitle)
			assert.Equal(t, string(urgencyV1.Medium), data.UrgencyLevel)
			return nil
		})

		req := &activityV1.ActivityCreateRequest{Description: "x", EmployeeID: 5, UrgencyID: 3}
		ctx := context.WithValue(t.Context(), "employeeID", uint(5))
		ctx = context.WithValue(ctx, "role", "Medic")
		_, err := svc.CreateActivity(ctx, req)
		assert.NoError(t, err)
	})
}

func TestActivityService_GetActivityByID(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves activity by ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		expectedActivity := &model.Activity{
			ID:          42,
			Description: "Emergency response completed",
			EmployeeID:  100,
			UrgencyID:   200,
			CreatedAt:   time.Now().Add(-1 * time.Hour),
			UpdatedAt:   time.Now(),
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), uint(42)).Return(expectedActivity, nil)

		response, err := service.GetActivityByID(t.Context(), 42)
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
		service := NewActivityService(log, mockRepo, nil)

		mockRepo.EXPECT().GetByID(gomock.Any(), uint(999)).Return(nil, fmt.Errorf("activity not found"))

		response, err := service.GetActivityByID(t.Context(), 999)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to get activity")
	})

	t.Run("returns error for zero ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		response, err := service.GetActivityByID(t.Context(), 0)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "invalid activity ID")
	})

	t.Run("successfully retrieves activity by ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		activity := &model.Activity{
			ID:          1,
			Description: "New employee was created",
			EmployeeID:  1,
			UrgencyID:   2,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockRepo.EXPECT().GetByID(gomock.Any(), uint(1)).Return(activity, nil)

		response, err := service.GetActivityByID(t.Context(), 1)
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
		service := NewActivityService(log, mockRepo, nil)

		mockRepo.EXPECT().GetByID(gomock.Any(), uint(999)).Return(nil, fmt.Errorf("activity not found"))

		response, err := service.GetActivityByID(t.Context(), 999)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to get activity")
	})
}

func TestActivityService_ListActivities(t *testing.T) {
	t.Parallel()

	t.Run("successfully lists activities with filters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

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

		mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, filter *model.ActivityFilter) ([]model.Activity, int64, error) {
			assert.Equal(t, &employeeID, filter.EmployeeID)
			assert.Equal(t, &urgencyID, filter.UrgencyID)
			assert.Equal(t, 1, filter.Page)
			assert.Equal(t, 10, filter.PageSize)
			assert.NotNil(t, filter.StartDate)
			assert.NotNil(t, filter.EndDate)
			return expectedActivities, 2, nil
		})

		response, err := service.ListActivities(t.Context(), req)
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
		service := NewActivityService(log, mockRepo, nil)

		req := &activityV1.ActivityListRequest{
			Page:     1,
			PageSize: 10,
		}

		mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]model.Activity{}, int64(0), nil)

		response, err := service.ListActivities(t.Context(), req)
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
		service := NewActivityService(log, mockRepo, nil)

		req := &activityV1.ActivityListRequest{
			Page:     1,
			PageSize: 10,
		}

		mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, int64(0), fmt.Errorf("database timeout"))

		response, err := service.ListActivities(t.Context(), req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to list activities")
	})

	t.Run("returns error for invalid date formats", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		req := &activityV1.ActivityListRequest{
			Page:      1,
			PageSize:  10,
			StartDate: "invalid-date",
			EndDate:   "also-invalid",
		}

		// Service should return validation error for invalid dates
		response, err := service.ListActivities(t.Context(), req)
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("validates request parameters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		req := &activityV1.ActivityListRequest{
			Page:     1,
			PageSize: 10,
		}

		mockRepo.EXPECT().List(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, filter *model.ActivityFilter) ([]model.Activity, int64, error) {
			// Validate that request validation was called
			assert.Equal(t, 1, filter.Page)
			assert.Equal(t, 10, filter.PageSize)
			return []model.Activity{}, 0, nil
		})

		response, err := service.ListActivities(t.Context(), req)
		assert.NoError(t, err)
		require.NotNil(t, response)
	})
}

func TestActivityService_GetActivityStats(t *testing.T) {
	t.Parallel()

	t.Run("it returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		mockRepo.EXPECT().GetStats(gomock.Any()).Return(nil, fmt.Errorf("database timeout"))

		response, err := service.GetActivityStats(t.Context())
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "failed to get activity stats")
	})

	t.Run("it successfully retrieves activity statistics", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		expectedStats := &model.ActivityStats{
			TotalActivities: 100,
			RecentActivities: []model.Activity{
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
			},
			ActivitiesLast24h:    10,
			ActivitiesLast7Days:  50,
			ActivitiesLast30Days: 100,
		}

		mockRepo.EXPECT().GetStats(gomock.Any()).Return(expectedStats, nil)

		response, err := service.GetActivityStats(t.Context())
		assert.NoError(t, err)
		require.NotNil(t, response)
		assert.Equal(t, int64(100), response.TotalActivities)
		assert.Len(t, response.RecentActivities, 2)
		assert.Equal(t, int64(10), response.ActivitiesLast24h)
		assert.Equal(t, int64(50), response.ActivitiesLast7Days)
		assert.Equal(t, int64(100), response.ActivitiesLast30Days)
	})

}

func TestActivityService_DeleteActivity(t *testing.T) {
	t.Parallel()

	t.Run("it returns error for invalid ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		err := service.DeleteActivity(t.Context(), 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid activity ID")
	})

	t.Run("successfully deletes activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		mockRepo.EXPECT().Delete(gomock.Any(), uint(1)).Return(nil)

		err := service.DeleteActivity(t.Context(), 1)
		assert.NoError(t, err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		mockRepo.EXPECT().Delete(gomock.Any(), uint(999)).Return(fmt.Errorf("activity not found"))

		err := service.DeleteActivity(t.Context(), 999)
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
		service := NewActivityService(log, mockRepo, nil)

		mockRepo.EXPECT().ResetAllData(gomock.Any()).Return(nil)

		err := service.ResetAllData(t.Context())
		assert.NoError(t, err)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		mockRepo.EXPECT().ResetAllData(gomock.Any()).Return(fmt.Errorf("database error"))

		err := service.ResetAllData(t.Context())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to reset activity data")
	})
}

func TestActivityService_LogActivity(t *testing.T) {
	t.Parallel()

	t.Run("successfully logs activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		mockRepo.EXPECT().CreateWithOutbox(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, activity *model.Activity, _ any) error {
			assert.Equal(t, "User logged in successfully", activity.Description)
			assert.Equal(t, uint(50), activity.EmployeeID)
			assert.Equal(t, uint(75), activity.UrgencyID)
			return nil
		})

		err := service.LogActivity(t.Context(), "User logged in successfully", 50, 75)
		assert.NoError(t, err)
	})

	t.Run("returns error for empty description", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		err := service.LogActivity(t.Context(), "", 1, 2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
	})

	t.Run("handles repository error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		service := NewActivityService(log, mockRepo, nil)

		mockRepo.EXPECT().CreateWithOutbox(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("disk full"))

		err := service.LogActivity(t.Context(), "Test activity", 1, 2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create activity")
	})

	t.Run("it succeeds when CreateWithOutbox succeeds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		svc := NewActivityService(log, mockRepo, nil)

		mockRepo.EXPECT().CreateWithOutbox(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, a *model.Activity, _ any) error {
			// Ensure inputs are mapped correctly
			assert.Equal(t, "Quick note", a.Description)
			assert.Equal(t, uint(11), a.EmployeeID)
			assert.Equal(t, uint(22), a.UrgencyID)
			return nil
		})

		err := svc.LogActivity(t.Context(), "Quick note", 11, 22)
		assert.NoError(t, err)
	})

	t.Run("it fails/returns an error when CreateWithOutbox returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repositories.NewMockActivityRepository(ctrl)
		svc := NewActivityService(log, mockRepo, nil)

		mockRepo.EXPECT().CreateWithOutbox(gomock.Any(), gomock.Any(), gomock.Any()).Return(assert.AnError)

		err := svc.LogActivity(t.Context(), "Something happened", 1, 2)
		assert.Error(t, err)
	})
}
