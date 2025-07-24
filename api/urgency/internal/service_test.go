package internal

import (
	"context"
	"testing"
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
	"github.com/pd120424d/mountain-service/api/urgency/internal/repositories"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// mockEmployeeClientForTest is a simple mock for unit tests
type mockEmployeeClientForTest struct{}

func (m *mockEmployeeClientForTest) GetOnCallEmployees(ctx context.Context, shiftBuffer time.Duration) ([]employeeV1.EmployeeResponse, error) {
	return []employeeV1.EmployeeResponse{}, nil
}

func (m *mockEmployeeClientForTest) GetAllEmployees(ctx context.Context) ([]employeeV1.EmployeeResponse, error) {
	return []employeeV1.EmployeeResponse{}, nil
}

func (m *mockEmployeeClientForTest) GetEmployeeByID(ctx context.Context, employeeID uint) (*employeeV1.EmployeeResponse, error) {
	return nil, nil
}

func (m *mockEmployeeClientForTest) CheckActiveEmergencies(ctx context.Context, employeeID uint) (bool, error) {
	return false, nil
}

func TestUrgencyService_CreateUrgency(t *testing.T) {
	t.Parallel()

	t.Run("it successfully creates an urgency", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)
		mockEmployeeClient := &mockEmployeeClientForTest{}

		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
		// Expect calls to assignment and notification repos for each employee
		// Since our mock returns empty list, no calls should be made

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(&model.Urgency{})
		assert.NoError(t, err)
	})

	t.Run("it returns an error when repository call fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)
		mockEmployeeClient := &mockEmployeeClientForTest{}

		mockRepo.EXPECT().Create(gomock.Any()).Return(assert.AnError)

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(&model.Urgency{})
		assert.Error(t, err)
	})
}

func TestUrgencyService_GetAllUrgencies(t *testing.T) {
	t.Parallel()

	t.Run("it successfully retrieves all urgencies", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockRepo.EXPECT().GetAll().Return([]model.Urgency{}, nil)

		svc := &urgencyService{log: log, repo: mockRepo}

		urgencies, err := svc.GetAllUrgencies()
		assert.NoError(t, err)
		assert.Len(t, urgencies, 0)
	})

	t.Run("it returns an error when repository call fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockRepo.EXPECT().GetAll().Return(nil, assert.AnError)

		svc := &urgencyService{log: log, repo: mockRepo}

		_, err := svc.GetAllUrgencies()
		assert.Error(t, err)
	})
}

func TestUrgencyService_GetUrgencyByID(t *testing.T) {
	t.Parallel()

	t.Run("it successfully retrieves an urgency by ID", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)

		mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).DoAndReturn(func(id uint, urgency *model.Urgency) error {
			*urgency = model.Urgency{ID: id}
			return nil
		})

		svc := &urgencyService{log: log, repo: mockRepo}

		urgency, err := svc.GetUrgencyByID(1)
		assert.NoError(t, err)

		assert.Equal(t, uint(1), urgency.ID)
	})

	t.Run("it returns an error when repository call fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(assert.AnError)

		svc := &urgencyService{log: log, repo: mockRepo}

		_, err := svc.GetUrgencyByID(1)
		assert.Error(t, err)
	})
}

func TestUrgencyService_UpdateUrgency(t *testing.T) {
	t.Parallel()

	t.Run("it successfully updates an urgency", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockRepo.EXPECT().Update(gomock.Any()).Return(nil)

		svc := &urgencyService{log: log, repo: mockRepo}

		err := svc.UpdateUrgency(&model.Urgency{})
		assert.NoError(t, err)
	})

	t.Run("it returns an error when repository call fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockRepo.EXPECT().Update(gomock.Any()).Return(assert.AnError)

		svc := &urgencyService{log: log, repo: mockRepo}

		err := svc.UpdateUrgency(&model.Urgency{})
		assert.Error(t, err)
	})
}

func TestUrgencyService_DeleteUrgency(t *testing.T) {
	t.Parallel()

	t.Run("it successfully deletes an urgency", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockRepo.EXPECT().Delete(gomock.Any()).Return(nil)

		svc := &urgencyService{log: log, repo: mockRepo}

		err := svc.DeleteUrgency(1)
		assert.NoError(t, err)
	})

	t.Run("it returns an error when repository call fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockRepo.EXPECT().Delete(gomock.Any()).Return(assert.AnError)

		svc := &urgencyService{log: log, repo: mockRepo}

		err := svc.DeleteUrgency(1)
		assert.Error(t, err)
	})
}

func TestUrgencyService_ResetAllData(t *testing.T) {
	t.Parallel()

	t.Run("it successfully resets all data", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockRepo.EXPECT().ResetAllData().Return(nil)

		svc := &urgencyService{log: log, repo: mockRepo}

		err := svc.ResetAllData()
		assert.NoError(t, err)
	})

	t.Run("it returns an error when repository call fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockRepo.EXPECT().ResetAllData().Return(assert.AnError)

		svc := &urgencyService{log: log, repo: mockRepo}

		err := svc.ResetAllData()
		assert.Error(t, err)
	})
}

func TestNewUrgencyService(t *testing.T) {
	t.Parallel()

	t.Run("it creates a new urgency service", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)
		mockEmployeeClient := &mockEmployeeClientForTest{}

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)
		assert.NotNil(t, svc)
		assert.IsType(t, &urgencyService{}, svc)
	})
}
