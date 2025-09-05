package internal

import (
	"context"
	"testing"
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/clients"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
	"github.com/pd120424d/mountain-service/api/urgency/internal/repositories"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUrgencyService_CreateUrgency(t *testing.T) {
	t.Parallel()

	t.Run("it successfully creates an urgency", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)
		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mockEmployeeClient.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return([]employeeV1.EmployeeResponse{}, nil)

		svc := NewUrgencyService(log, mockRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(context.Background(), &model.Urgency{})
		assert.NoError(t, err)
	})

	t.Run("it returns an error when repository call fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)
		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(assert.AnError)

		svc := NewUrgencyService(log, mockRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(context.Background(), &model.Urgency{})
		assert.Error(t, err)
	})

	t.Run("it returns an error when it fails to fetch on-call employees", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)
		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mockEmployeeClient.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)

		svc := NewUrgencyService(log, mockRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(context.Background(), &model.Urgency{})
		assert.Error(t, err)
	})

	t.Run("it logs an error and continue when it fails to create assignment and notification", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)
		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mockEmployeeClient.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return([]employeeV1.EmployeeResponse{
			{
				ID:        1,
				FirstName: "Marko",
				LastName:  "Markovic",
				Phone:     "+1987654321",
				Email:     "marko@example.com",
				Username:  "Marko",
			},
		}, nil)

		// Since we no longer create assignments on urgency creation, only notifications are attempted.
		// Simulate failures in notification creation and ensure service logs and continues without error.
		mockNotificationRepo.EXPECT().Create(gomock.Any()).AnyTimes().Return(assert.AnError)

		svc := NewUrgencyService(log, mockRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(context.Background(), &model.Urgency{})
		assert.NoError(t, err)
	})
}

func TestUrgencyService_GetAllUrgencies(t *testing.T) {
	t.Parallel()

	t.Run("it successfully retrieves all urgencies", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockRepo.EXPECT().GetAll(gomock.Any()).Return([]model.Urgency{}, nil)

		svc := &urgencyService{log: log, repo: mockRepo}

		urgencies, err := svc.GetAllUrgencies(context.Background())
		assert.NoError(t, err)
		assert.Len(t, urgencies, 0)
	})

	t.Run("it returns an error when repository call fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockRepo.EXPECT().GetAll(gomock.Any()).Return(nil, assert.AnError)

		svc := &urgencyService{log: log, repo: mockRepo}

		_, err := svc.GetAllUrgencies(context.Background())
		assert.Error(t, err)
	})
}

func TestUrgencyService_ListUrgencies(t *testing.T) {
	log := utils.NewTestLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositories.NewMockUrgencyRepository(ctrl)
	repo.EXPECT().ListPaginated(gomock.Any(), 2, 10, gomock.Nil()).Return([]model.Urgency{{ID: 1}}, int64(11), nil)

	svc := &urgencyService{log: log, repo: repo}
	items, total, err := svc.ListUrgencies(context.Background(), 2, 10, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(11), total)
	assert.Len(t, items, 1)
}

func TestUrgencyService_GetUrgencyByID(t *testing.T) {
	t.Parallel()

	t.Run("it successfully retrieves an urgency by ID", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)

		mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, urgency *model.Urgency) error {
			*urgency = model.Urgency{ID: id}
			return nil
		})

		svc := &urgencyService{log: log, repo: mockRepo}

		urgency, err := svc.GetUrgencyByID(context.Background(), 1)
		assert.NoError(t, err)

		assert.Equal(t, uint(1), urgency.ID)
	})

	t.Run("it returns an error when repository call fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(assert.AnError)

		svc := &urgencyService{log: log, repo: mockRepo}

		_, err := svc.GetUrgencyByID(context.Background(), 1)
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
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

		svc := &urgencyService{log: log, repo: mockRepo}

		err := svc.UpdateUrgency(context.Background(), &model.Urgency{})
		assert.NoError(t, err)
	})

	t.Run("it returns an error when repository call fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(assert.AnError)

		svc := &urgencyService{log: log, repo: mockRepo}

		err := svc.UpdateUrgency(context.Background(), &model.Urgency{})
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
		mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)

		svc := &urgencyService{log: log, repo: mockRepo}

		err := svc.DeleteUrgency(context.Background(), 1)
		assert.NoError(t, err)
	})

	t.Run("it returns an error when repository call fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(assert.AnError)

		svc := &urgencyService{log: log, repo: mockRepo}

		err := svc.DeleteUrgency(context.Background(), 1)
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
		mockRepo.EXPECT().ResetAllData(gomock.Any()).Return(nil)

		svc := &urgencyService{log: log, repo: mockRepo}

		err := svc.ResetAllData(context.Background())
		assert.NoError(t, err)
	})

	t.Run("it returns an error when repository call fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockRepo.EXPECT().ResetAllData(gomock.Any()).Return(assert.AnError)

		svc := &urgencyService{log: log, repo: mockRepo}

		err := svc.ResetAllData(context.Background())
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
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)
		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		svc := NewUrgencyService(log, mockRepo, mockNotificationRepo, mockEmployeeClient)
		assert.NotNil(t, svc)
		assert.IsType(t, &urgencyService{}, svc)
	})
}

func TestUrgencyService_createAssignmentAndNotification(t *testing.T) {
	t.Parallel()

	t.Run("it successfully creates assignment and SMS notification for employee with phone", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)

		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)
		employees := []employeeV1.EmployeeResponse{
			{
				ID:        1,
				FirstName: "Marko",
				LastName:  "Markovic",
				Phone:     "+1987654321",
				Email:     "marko@example.com",
			},
		}
		mockEmployeeClient.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(employees, nil)

		urgency := &model.Urgency{
			ID:           1,
			FirstName:    "Marko",
			LastName:     "Markovic",
			Location:     "Mountain Peak",
			ContactPhone: "+1234567890",
			Description:  "Lost hiker",
			Level:        "High",
		}

		employee := employeeV1.EmployeeResponse{
			ID:        1,
			FirstName: "Marko",
			LastName:  "Markovic",
			Phone:     "+1987654321",
			Email:     "marko@example.com",
		}

		mockNotificationRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(notification *model.Notification) error {
			assert.Equal(t, urgency.ID, notification.UrgencyID)
			assert.Equal(t, employee.ID, notification.EmployeeID)
			assert.Equal(t, model.NotificationSMS, notification.NotificationType)
			assert.Equal(t, employee.Phone, notification.Recipient)
			assert.Equal(t, model.NotificationPending, notification.Status)
			assert.Contains(t, notification.Message, "EMERGENCY")
			assert.Contains(t, notification.Message, urgency.Description)
			notification.ID = 1
			return nil
		})

		mockNotificationRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(notification *model.Notification) error {
			assert.Equal(t, model.NotificationEmail, notification.NotificationType)
			assert.Equal(t, "marko@example.com", notification.Recipient)
			return nil
		})

		svc := NewUrgencyService(log, mockRepo, mockNotificationRepo, mockEmployeeClient)

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		err := svc.CreateUrgency(context.Background(), urgency)
		assert.NoError(t, err)
	})

	t.Run("it successfully creates assignment and email notification for employee without phone", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)

		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)
		employees := []employeeV1.EmployeeResponse{
			{
				ID:        1,
				FirstName: "Marko",
				LastName:  "Markovic",
				Phone:     "",
				Email:     "marko@example.com",
			},
		}
		mockEmployeeClient.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(employees, nil)

		urgency := &model.Urgency{
			ID:           1,
			FirstName:    "Marko",
			LastName:     "Markovic",
			Location:     "Mountain Peak",
			ContactPhone: "+1234567890",
			Description:  "Lost hiker",
			Level:        "High",
		}

		mockNotificationRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(notification *model.Notification) error {
			assert.Equal(t, model.NotificationEmail, notification.NotificationType)
			assert.Equal(t, "marko@example.com", notification.Recipient)
			return nil
		})

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		svc := NewUrgencyService(log, mockRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(context.Background(), urgency)
		assert.NoError(t, err)
	})

	t.Run("it returns an error when it fails to create assignment", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)

		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		employees := []employeeV1.EmployeeResponse{
			{
				ID:        1,
				FirstName: "Marko",
				LastName:  "Markovic",
				Phone:     "+1987654321",
				Email:     "marko@example.com",
			},
		}
		mockEmployeeClient.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(employees, nil)

		urgency := &model.Urgency{
			ID:           1,
			FirstName:    "Marko",
			LastName:     "Markovic",
			Location:     "Mountain Peak",
			ContactPhone: "+1234567890",
			Description:  "Lost hiker",
			Level:        "High",
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		// Notifications are created on create (no assignments). Accept any number of creates.
		mockNotificationRepo.EXPECT().Create(gomock.Any()).AnyTimes().Return(nil)

		svc := NewUrgencyService(log, mockRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(context.Background(), urgency)
		assert.NoError(t, err) // CreateUrgency should not return error, it logs and continues
	})

	t.Run("it returns an error when it fails to create sms notification", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)

		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		employees := []employeeV1.EmployeeResponse{
			{
				ID:        1,
				FirstName: "Marko",
				LastName:  "Markovic",
				Phone:     "+1987654321",
				Email:     "marko@example.com",
			},
		}
		mockEmployeeClient.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(employees, nil)

		urgency := &model.Urgency{
			ID:           1,
			FirstName:    "Marko",
			LastName:     "Markovic",
			Location:     "Mountain Peak",
			ContactPhone: "+1234567890",
			Description:  "Lost hiker",
			Level:        "High",
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mockNotificationRepo.EXPECT().Create(gomock.Any()).Return(assert.AnError).AnyTimes()

		svc := NewUrgencyService(log, mockRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(context.Background(), urgency)
		assert.NoError(t, err) // CreateUrgency should not return error, it logs and continues
	})

	t.Run("it returns an error when it fails to create email notification", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)

		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		employees := []employeeV1.EmployeeResponse{
			{
				ID:        1,
				FirstName: "Marko",
				LastName:  "Markovic",
				Phone:     "", // No phone, only email
				Email:     "marko@example.com",
			},
		}
		mockEmployeeClient.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(employees, nil)

		urgency := &model.Urgency{
			ID:           1,
			FirstName:    "Marko",
			LastName:     "Markovic",
			Location:     "Mountain Peak",
			ContactPhone: "+1234567890",
			Description:  "Lost hiker",
			Level:        "High",
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		mockNotificationRepo.EXPECT().Create(gomock.Any()).Return(assert.AnError)

		svc := NewUrgencyService(log, mockRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(context.Background(), urgency)
		assert.NoError(t, err) // CreateUrgency should not return error, it logs and continues
	})

	t.Run("it handles employee with no phone and no email", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)

		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		employees := []employeeV1.EmployeeResponse{
			{
				ID:        1,
				FirstName: "Marko",
				LastName:  "Markovic",
				Phone:     "", // No phone
				Email:     "", // No email
			},
		}
		mockEmployeeClient.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(employees, nil)

		urgency := &model.Urgency{
			ID:           1,
			FirstName:    "Marko",
			LastName:     "Markovic",
			Location:     "Mountain Peak",
			ContactPhone: "+1234567890",
			Description:  "Lost hiker",
			Level:        "High",
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		// No notification expectations since employee has no contact info

		svc := NewUrgencyService(log, mockRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(context.Background(), urgency)
		assert.NoError(t, err)
	})

	t.Run("it creates both SMS and email notifications for employee with both contacts", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)

		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		employees := []employeeV1.EmployeeResponse{
			{
				ID:        1,
				FirstName: "Marko",
				LastName:  "Markovic",
				Phone:     "+1987654321",
				Email:     "marko@example.com",
			},
		}
		mockEmployeeClient.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(employees, nil)

		urgency := &model.Urgency{
			ID:           1,
			FirstName:    "Marko",
			LastName:     "Markovic",
			Location:     "Mountain Peak",
			ContactPhone: "+1234567890",
			Description:  "Lost hiker",
			Level:        "High",
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		// Expect two notification creations: SMS and Email
		mockNotificationRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(notification *model.Notification) error {
			assert.Equal(t, model.NotificationSMS, notification.NotificationType)
			assert.Equal(t, "+1987654321", notification.Recipient)
			return nil
		})
		mockNotificationRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(notification *model.Notification) error {
			assert.Equal(t, model.NotificationEmail, notification.NotificationType)
			assert.Equal(t, "marko@example.com", notification.Recipient)
			return nil
		})

		svc := NewUrgencyService(log, mockRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(context.Background(), urgency)
		assert.NoError(t, err)
	})

	t.Run("it handles multiple employees with mixed success and failures", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)

		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		employees := []employeeV1.EmployeeResponse{
			{
				ID:        1,
				FirstName: "Marko",
				LastName:  "Markovic",
				Phone:     "+1987654321",
				Email:     "marko@example.com",
			},
			{
				ID:        2,
				FirstName: "Marko",
				LastName:  "Markovic",
				Phone:     "+1123456789",
				Email:     "john@example.com",
			},
		}
		mockEmployeeClient.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(employees, nil)

		urgency := &model.Urgency{
			ID:           1,
			FirstName:    "Emergency",
			LastName:     "Situation",
			Location:     "Mountain Peak",
			ContactPhone: "+1234567890",
			Description:  "Lost hiker",
			Level:        "High",
		}

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		// First employee - notifications succeed
		mockNotificationRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(2) // SMS + Email

		// Second employee - notifications fail (both SMS and Email), service should log and continue
		mockNotificationRepo.EXPECT().Create(gomock.Any()).Return(assert.AnError).Times(2)

		svc := NewUrgencyService(log, mockRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(context.Background(), urgency)
		assert.NoError(t, err) // Should not return error even if some notifications fail
	})
}

func TestUrgencyService_buildNotificationMessage(t *testing.T) {
	t.Parallel()

	t.Run("it builds SMS notification message correctly", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)

		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)
		employees := []employeeV1.EmployeeResponse{
			{
				ID:        1,
				FirstName: "Marko",
				LastName:  "Markovic",
				Phone:     "+1987654321",
				Email:     "marko@example.com",
			},
		}
		mockEmployeeClient.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(employees, nil)

		urgency := &model.Urgency{
			ID:           1,
			FirstName:    "Marko",
			LastName:     "Markovic",
			Location:     "Mountain Peak",
			ContactPhone: "+1234567890",
			Description:  "Lost hiker",
			Level:        urgencyV1.High,
		}

		mockNotificationRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(notification *model.Notification) error {
			assert.Equal(t, model.NotificationSMS, notification.NotificationType)

			expectedContent := []string{
				"🚨 EMERGENCY:",
				urgency.Description,
				urgency.Location,
				urgency.FirstName + " " + urgency.LastName,
				urgency.ContactPhone,
				string(urgency.Level),
			}

			for _, content := range expectedContent {
				assert.Contains(t, notification.Message, content)
			}

			return nil
		})

		mockNotificationRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(notification *model.Notification) error {
			assert.Equal(t, model.NotificationEmail, notification.NotificationType)
			assert.Equal(t, "marko@example.com", notification.Recipient)
			return nil
		})

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		svc := NewUrgencyService(log, mockRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(context.Background(), urgency)
		assert.NoError(t, err)
	})

	t.Run("it builds email notification message correctly", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)

		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)
		employees := []employeeV1.EmployeeResponse{
			{
				ID:        1,
				FirstName: "Marko",
				LastName:  "Markovic",
				Phone:     "",
				Email:     "marko@example.com",
			},
		}
		mockEmployeeClient.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(employees, nil)

		urgency := &model.Urgency{
			ID:           1,
			FirstName:    "Marko",
			LastName:     "Markovic",
			Location:     "Mountain Peak",
			ContactPhone: "+1234567890",
			Description:  "Lost hiker",
			Level:        urgencyV1.High,
		}

		mockNotificationRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(notification *model.Notification) error {
			assert.Equal(t, model.NotificationEmail, notification.NotificationType)

			expectedContent := []string{
				"🚨 EMERGENCY ALERT 🚨",
				"Hello Marko Markovic",
				"📍 Location: " + urgency.Location,
				"📞 Contact: " + urgency.FirstName + " " + urgency.LastName,
				urgency.ContactPhone,
				"📝 Description: " + urgency.Description,
				"⚠️ Priority: " + string(urgency.Level),
				"Please respond immediately",
			}

			for _, content := range expectedContent {
				assert.Contains(t, notification.Message, content)
			}

			return nil
		})

		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		svc := NewUrgencyService(log, mockRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(context.Background(), urgency)
		assert.NoError(t, err)
	})
}

func TestUrgencyService_AssignUrgency(t *testing.T) {
	t.Parallel()

	t.Run("it returns error on invalid parameters", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockUrgencyRepository(ctrl)
		nrepo := repositories.NewMockNotificationRepository(ctrl)
		ecli := clients.NewMockEmployeeClient(ctrl)

		svc := NewUrgencyService(log, repo, nrepo, ecli)
		err := svc.AssignUrgency(context.Background(), 0, 10)
		assert.Error(t, err)
		err = svc.AssignUrgency(context.Background(), 10, 0)
		assert.Error(t, err)
	})

	t.Run("it returns error when urgency not found", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockUrgencyRepository(ctrl)
		nrepo := repositories.NewMockNotificationRepository(ctrl)
		ecli := clients.NewMockEmployeeClient(ctrl)

		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(1), gomock.Any()).Return(assert.AnError)
		svc := NewUrgencyService(log, repo, nrepo, ecli)
		err := svc.AssignUrgency(context.Background(), 1, 2)
		assert.Error(t, err)
	})

	t.Run("it returns error when urgency already has assignee", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockUrgencyRepository(ctrl)
		nrepo := repositories.NewMockNotificationRepository(ctrl)
		ecli := clients.NewMockEmployeeClient(ctrl)

		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(1), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error {
			emp := uint(5)
			u.AssignedEmployeeID = &emp
			u.ID = id
			return nil
		})
		svc := NewUrgencyService(log, repo, nrepo, ecli)
		err := svc.AssignUrgency(context.Background(), 1, 2)
		assert.Error(t, err)
	})

	t.Run("it succeeds when urgency is unassigned and updates with assigned fields", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockUrgencyRepository(ctrl)
		nrepo := repositories.NewMockNotificationRepository(ctrl)
		ecli := clients.NewMockEmployeeClient(ctrl)

		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(1), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error { *u = model.Urgency{ID: id}; return nil })
		repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		ecli.EXPECT().GetEmployeeByID(gomock.Any(), uint(2)).Return(&employeeV1.EmployeeResponse{ID: 2}, nil)

		svc := NewUrgencyService(log, repo, nrepo, ecli)
		err := svc.AssignUrgency(context.Background(), 1, 2)
		assert.NoError(t, err)

		t.Run("it returns error when employee does not exist", func(t *testing.T) {
			log := utils.NewTestLogger()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := repositories.NewMockUrgencyRepository(ctrl)
			nrepo := repositories.NewMockNotificationRepository(ctrl)
			ecli := clients.NewMockEmployeeClient(ctrl)

			repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(1), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error { *u = model.Urgency{ID: id}; return nil })
			ecli.EXPECT().GetEmployeeByID(gomock.Any(), uint(2)).Return(nil, assert.AnError)

			svc := NewUrgencyService(log, repo, nrepo, ecli)
			err := svc.AssignUrgency(context.Background(), 1, 2)
			assert.Error(t, err)
		})

	})
}

func TestUrgencyService_UnassignUrgency(t *testing.T) {
	t.Parallel()

	t.Run("it returns error on invalid urgencyID", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := &urgencyService{log: log}
		err := svc.UnassignUrgency(context.Background(), 0, 10, false)
		assert.Error(t, err)
	})

	t.Run("it returns error when urgency repo GetByID fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockUrgencyRepository(ctrl)
		svc := &urgencyService{log: log, repo: repo}
		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(1), gomock.Any()).Return(assert.AnError)
		err := svc.UnassignUrgency(context.Background(), 1, 99, false)
		assert.Error(t, err)
	})

	t.Run("it returns error when urgency is not assigned", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockUrgencyRepository(ctrl)
		svc := &urgencyService{log: log, repo: repo}
		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(1), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error { *u = model.Urgency{ID: id}; return nil })
		err := svc.UnassignUrgency(context.Background(), 1, 99, false)
		assert.Error(t, err)
	})

	t.Run("it returns forbidden when actor is not assignee and not admin", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockUrgencyRepository(ctrl)
		svc := &urgencyService{log: log, repo: repo}
		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(1), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error {
			emp := uint(55)
			*u = model.Urgency{ID: id, AssignedEmployeeID: &emp}
			return nil
		})
		err := svc.UnassignUrgency(context.Background(), 1, 99, false)
		assert.Error(t, err)
	})

	t.Run("it returns error when repo update fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockUrgencyRepository(ctrl)
		svc := &urgencyService{log: log, repo: repo}
		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(1), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error {
			emp := uint(55)
			*u = model.Urgency{ID: id, AssignedEmployeeID: &emp}
			return nil
		})
		repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(assert.AnError)
		err := svc.UnassignUrgency(context.Background(), 1, 55, false)
		assert.Error(t, err)
	})

	t.Run("it succeeds when assignee unassigns", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockUrgencyRepository(ctrl)
		svc := &urgencyService{log: log, repo: repo}
		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(1), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error {
			emp := uint(55)
			*u = model.Urgency{ID: id, AssignedEmployeeID: &emp}
			return nil
		})
		repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		err := svc.UnassignUrgency(context.Background(), 1, 55, false)
		assert.NoError(t, err)
	})

	t.Run("it succeeds when admin unassigns", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockUrgencyRepository(ctrl)
		svc := &urgencyService{log: log, repo: repo}
		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(1), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error {
			emp := uint(55)
			*u = model.Urgency{ID: id, AssignedEmployeeID: &emp}
			return nil
		})
		repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		err := svc.UnassignUrgency(context.Background(), 1, 99, true)
		assert.NoError(t, err)
	})
}

func TestUrgencyService_GetAssignment(t *testing.T) {
	log := utils.NewTestLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repositories.NewMockUrgencyRepository(ctrl)
	nrepo := repositories.NewMockNotificationRepository(ctrl)
	ecli := clients.NewMockEmployeeClient(ctrl)
	svc := NewUrgencyService(log, repo, nrepo, ecli)

	t.Run("it returns nil when unassigned", func(t *testing.T) {
		repo.EXPECT().GetByID(gomock.Any(), uint(1), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error {
			*u = model.Urgency{ID: id}
			return nil
		})
		resp, err := svc.GetAssignment(context.Background(), 1)
		assert.NoError(t, err)
		assert.Nil(t, resp)
	})

	t.Run("it returns DTO when assigned", func(t *testing.T) {
		emp := uint(42)
		now := time.Now().UTC().Truncate(time.Second)
		repo.EXPECT().GetByID(gomock.Any(), uint(2), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error {
			*u = model.Urgency{ID: id, AssignedEmployeeID: &emp, AssignedAt: &now}
			return nil
		})
		resp, err := svc.GetAssignment(context.Background(), 2)
		assert.NoError(t, err)
		if assert.NotNil(t, resp) {
			assert.Equal(t, uint(2), resp.UrgencyID)
			assert.Equal(t, uint(42), resp.AssignedEmployee)
			assert.Equal(t, now.Format(time.RFC3339), resp.AssignedAt)
		}
	})
}

func TestUrgencyService_CloseUrgency(t *testing.T) {
	t.Parallel()

	t.Run("it returns error on invalid urgencyID", func(t *testing.T) {
		log := utils.NewTestLogger()
		svc := &urgencyService{log: log}
		err := svc.CloseUrgency(context.Background(), 0, 10, false)
		assert.Error(t, err)
	})

	t.Run("it returns error when GetByID fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := repositories.NewMockUrgencyRepository(ctrl)
		svc := &urgencyService{log: log, repo: repo}
		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(1), gomock.Any()).Return(assert.AnError)
		err := svc.CloseUrgency(context.Background(), 1, 99, false)
		assert.Error(t, err)
	})

	t.Run("it returns error when urgency has no assignee", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := repositories.NewMockUrgencyRepository(ctrl)
		svc := &urgencyService{log: log, repo: repo}
		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(1), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error { *u = model.Urgency{ID: id}; return nil })
		err := svc.CloseUrgency(context.Background(), 1, 99, false)
		assert.Error(t, err)
	})

	t.Run("it returns error when status is not in_progress", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := repositories.NewMockUrgencyRepository(ctrl)
		svc := &urgencyService{log: log, repo: repo}
		emp := uint(22)
		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(2), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error {
			*u = model.Urgency{ID: id, AssignedEmployeeID: &emp, Status: urgencyV1.Open}
			return nil
		})
		err := svc.CloseUrgency(context.Background(), 2, 22, false)
		assert.Error(t, err)
	})

	t.Run("it returns forbidden when actor is not assignee and not admin", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := repositories.NewMockUrgencyRepository(ctrl)
		svc := &urgencyService{log: log, repo: repo}
		emp := uint(55)
		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(3), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error {
			*u = model.Urgency{ID: id, AssignedEmployeeID: &emp, Status: urgencyV1.InProgress}
			return nil
		})
		err := svc.CloseUrgency(context.Background(), 3, 99, false)
		assert.Error(t, err)
	})

	t.Run("it returns error when assigned employee is invalid", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := repositories.NewMockUrgencyRepository(ctrl)
		ecli := clients.NewMockEmployeeClient(ctrl)
		svc := &urgencyService{log: log, repo: repo, employeeClient: ecli}
		emp := uint(77)
		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(4), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error {
			*u = model.Urgency{ID: id, AssignedEmployeeID: &emp, Status: urgencyV1.InProgress}
			return nil
		})
		ecli.EXPECT().GetEmployeeByID(gomock.Any(), emp).Return(nil, assert.AnError)
		err := svc.CloseUrgency(context.Background(), 4, emp, true)
		assert.Error(t, err)
	})

	t.Run("it returns error when repo update fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := repositories.NewMockUrgencyRepository(ctrl)
		ecli := clients.NewMockEmployeeClient(ctrl)
		svc := &urgencyService{log: log, repo: repo, employeeClient: ecli}
		emp := uint(5)
		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(5), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error {
			*u = model.Urgency{ID: id, AssignedEmployeeID: &emp, Status: urgencyV1.InProgress}
			return nil
		})
		ecli.EXPECT().GetEmployeeByID(gomock.Any(), emp).Return(&employeeV1.EmployeeResponse{ID: emp}, nil)
		repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(assert.AnError)
		err := svc.CloseUrgency(context.Background(), 5, emp, true)
		assert.Error(t, err)
	})

	t.Run("it succeeds when assignee closes", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := repositories.NewMockUrgencyRepository(ctrl)
		ecli := clients.NewMockEmployeeClient(ctrl)
		svc := &urgencyService{log: log, repo: repo, employeeClient: ecli}
		emp := uint(9)
		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(6), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error {
			*u = model.Urgency{ID: id, AssignedEmployeeID: &emp, Status: urgencyV1.InProgress}
			return nil
		})
		ecli.EXPECT().GetEmployeeByID(gomock.Any(), emp).Return(&employeeV1.EmployeeResponse{ID: emp}, nil)
		repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		err := svc.CloseUrgency(context.Background(), 6, emp, false)
		assert.NoError(t, err)
	})

	t.Run("it succeeds when admin closes", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repo := repositories.NewMockUrgencyRepository(ctrl)
		ecli := clients.NewMockEmployeeClient(ctrl)
		svc := &urgencyService{log: log, repo: repo, employeeClient: ecli}
		emp := uint(12)
		repo.EXPECT().GetByIDPrimary(gomock.Any(), uint(7), gomock.Any()).DoAndReturn(func(_ context.Context, id uint, u *model.Urgency) error {
			*u = model.Urgency{ID: id, AssignedEmployeeID: &emp, Status: urgencyV1.InProgress}
			return nil
		})
		ecli.EXPECT().GetEmployeeByID(gomock.Any(), emp).Return(&employeeV1.EmployeeResponse{ID: emp}, nil)
		repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		err := svc.CloseUrgency(context.Background(), 7, 999, true)
		assert.NoError(t, err)
	})
}
