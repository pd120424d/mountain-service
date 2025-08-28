package internal

import (
	"testing"

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
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)
		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
		mockEmployeeClient.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return([]employeeV1.EmployeeResponse{}, nil)

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
		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		mockRepo.EXPECT().Create(gomock.Any()).Return(assert.AnError)

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(&model.Urgency{})
		assert.Error(t, err)
	})

	t.Run("it returns an error when it fails to fetch on-call employees", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)
		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
		mockEmployeeClient.EXPECT().GetOnCallEmployees(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(&model.Urgency{})
		assert.Error(t, err)
	})

	t.Run("it logs an error and continue when it fails to create assignment and notification", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
		mockNotificationRepo := repositories.NewMockNotificationRepository(mockCtrl)
		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
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

		mockAssignmentRepo.EXPECT().Create(gomock.Any()).Return(assert.AnError)

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(&model.Urgency{})
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
		mockEmployeeClient := clients.NewMockEmployeeClient(mockCtrl)

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)
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
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
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

		mockAssignmentRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(assignment *model.EmergencyAssignment) error {
			assert.Equal(t, urgency.ID, assignment.UrgencyID)
			assert.Equal(t, employee.ID, assignment.EmployeeID)
			assert.Equal(t, model.AssignmentPending, assignment.Status)
			assignment.ID = 1
			return nil
		})

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

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)

		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

		err := svc.CreateUrgency(urgency)
		assert.NoError(t, err)
	})

	t.Run("it successfully creates assignment and email notification for employee without phone", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
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

		mockAssignmentRepo.EXPECT().Create(gomock.Any()).Return(nil)

		mockNotificationRepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(notification *model.Notification) error {
			assert.Equal(t, model.NotificationEmail, notification.NotificationType)
			assert.Equal(t, "marko@example.com", notification.Recipient)
			return nil
		})

		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(urgency)
		assert.NoError(t, err)
	})

	t.Run("it returns an error when it fails to create assignment", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
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

		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
		mockAssignmentRepo.EXPECT().Create(gomock.Any()).Return(assert.AnError)

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(urgency)
		assert.NoError(t, err) // CreateUrgency should not return error, it logs and continues
	})

	t.Run("it returns an error when it fails to create sms notification", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
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

		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
		mockAssignmentRepo.EXPECT().Create(gomock.Any()).Return(nil)
		mockNotificationRepo.EXPECT().Create(gomock.Any()).Return(assert.AnError).AnyTimes()

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(urgency)
		assert.NoError(t, err) // CreateUrgency should not return error, it logs and continues
	})

	t.Run("it returns an error when it fails to create email notification", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
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

		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
		mockAssignmentRepo.EXPECT().Create(gomock.Any()).Return(nil)
		mockNotificationRepo.EXPECT().Create(gomock.Any()).Return(assert.AnError)

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(urgency)
		assert.NoError(t, err) // CreateUrgency should not return error, it logs and continues
	})

	t.Run("it handles employee with no phone and no email", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
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

		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
		mockAssignmentRepo.EXPECT().Create(gomock.Any()).Return(nil)
		// No notification expectations since employee has no contact info

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(urgency)
		assert.NoError(t, err)
	})

	t.Run("it creates both SMS and email notifications for employee with both contacts", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
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

		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)
		mockAssignmentRepo.EXPECT().Create(gomock.Any()).Return(nil)

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

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(urgency)
		assert.NoError(t, err)
	})

	t.Run("it handles multiple employees with mixed success and failures", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
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

		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

		// First employee - assignment succeeds, notifications succeed
		mockAssignmentRepo.EXPECT().Create(gomock.Any()).Return(nil)
		mockNotificationRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(2) // SMS + Email

		// Second employee - assignment fails
		mockAssignmentRepo.EXPECT().Create(gomock.Any()).Return(assert.AnError)

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(urgency)
		assert.NoError(t, err) // Should not return error even if some assignments fail
	})
}

func TestUrgencyService_buildNotificationMessage(t *testing.T) {
	t.Parallel()

	t.Run("it builds SMS notification message correctly", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
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

		mockAssignmentRepo.EXPECT().Create(gomock.Any()).Return(nil)

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

		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(urgency)
		assert.NoError(t, err)
	})

	t.Run("it builds email notification message correctly", func(t *testing.T) {
		log := utils.NewTestLogger()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockRepo := repositories.NewMockUrgencyRepository(mockCtrl)
		mockAssignmentRepo := repositories.NewMockAssignmentRepository(mockCtrl)
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

		mockAssignmentRepo.EXPECT().Create(gomock.Any()).Return(nil)

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

		mockRepo.EXPECT().Create(gomock.Any()).Return(nil)

		svc := NewUrgencyService(log, mockRepo, mockAssignmentRepo, mockNotificationRepo, mockEmployeeClient)

		err := svc.CreateUrgency(urgency)
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
		arepo := repositories.NewMockAssignmentRepository(ctrl)
		nrepo := repositories.NewMockNotificationRepository(ctrl)
		ecli := clients.NewMockEmployeeClient(ctrl)

		svc := NewUrgencyService(log, repo, arepo, nrepo, ecli)
		_, err := svc.AssignUrgency(0, 10)
		assert.Error(t, err)
		_, err = svc.AssignUrgency(10, 0)
		assert.Error(t, err)
	})

	t.Run("it returns error when urgency not found", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockUrgencyRepository(ctrl)
		arepo := repositories.NewMockAssignmentRepository(ctrl)
		nrepo := repositories.NewMockNotificationRepository(ctrl)
		ecli := clients.NewMockEmployeeClient(ctrl)

		repo.EXPECT().GetByID(uint(1), gomock.Any()).Return(assert.AnError)
		svc := NewUrgencyService(log, repo, arepo, nrepo, ecli)
		_, err := svc.AssignUrgency(1, 2)
		assert.Error(t, err)
	})

	t.Run("it returns error when urgency already has accepted assignment", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockUrgencyRepository(ctrl)
		arepo := repositories.NewMockAssignmentRepository(ctrl)
		nrepo := repositories.NewMockNotificationRepository(ctrl)
		ecli := clients.NewMockEmployeeClient(ctrl)

		repo.EXPECT().GetByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, u *model.Urgency) error { *u = model.Urgency{ID: id}; return nil })
		arepo.EXPECT().GetByUrgencyID(uint(1)).Return([]model.EmergencyAssignment{{Status: model.AssignmentAccepted}}, nil)

		svc := NewUrgencyService(log, repo, arepo, nrepo, ecli)
		_, err := svc.AssignUrgency(1, 2)
		assert.Error(t, err)
	})

	t.Run("it returns error when listing assignments fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockUrgencyRepository(ctrl)
		arepo := repositories.NewMockAssignmentRepository(ctrl)
		nrepo := repositories.NewMockNotificationRepository(ctrl)
		ecli := clients.NewMockEmployeeClient(ctrl)

		repo.EXPECT().GetByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, u *model.Urgency) error { *u = model.Urgency{ID: id}; return nil })
		arepo.EXPECT().GetByUrgencyID(uint(1)).Return(nil, assert.AnError)

		svc := NewUrgencyService(log, repo, arepo, nrepo, ecli)
		_, err := svc.AssignUrgency(1, 2)
		assert.Error(t, err)
	})

	t.Run("it creates accepted assignment and returns response", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := repositories.NewMockUrgencyRepository(ctrl)
		arepo := repositories.NewMockAssignmentRepository(ctrl)
		nrepo := repositories.NewMockNotificationRepository(ctrl)
		ecli := clients.NewMockEmployeeClient(ctrl)

		repo.EXPECT().GetByID(uint(1), gomock.Any()).DoAndReturn(func(id uint, u *model.Urgency) error { *u = model.Urgency{ID: id}; return nil })
		arepo.EXPECT().GetByUrgencyID(uint(1)).Return([]model.EmergencyAssignment{}, nil)
		arepo.EXPECT().Create(gomock.Any()).DoAndReturn(func(a *model.EmergencyAssignment) error {
			assert.Equal(t, uint(1), a.UrgencyID)
			assert.Equal(t, uint(2), a.EmployeeID)
			assert.Equal(t, model.AssignmentAccepted, a.Status)
			a.ID = 123
			return nil
		})

		svc := NewUrgencyService(log, repo, arepo, nrepo, ecli)
		resp, err := svc.AssignUrgency(1, 2)
		assert.NoError(t, err)
		assert.Equal(t, uint(123), resp.ID)
		assert.Equal(t, uint(1), resp.UrgencyID)
		assert.Equal(t, uint(2), resp.EmployeeID)
		assert.Equal(t, "accepted", resp.Status)
	})
}

func TestUrgencyService_UnassignUrgency(t *testing.T) {
	t.Parallel()

	t.Run("it returns error on invalid urgencyID", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc := &urgencyService{log: log}
		err := svc.UnassignUrgency(0, 10, false)
		assert.Error(t, err)
	})

	t.Run("it returns error when listing assignments fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		arepo := repositories.NewMockAssignmentRepository(ctrl)
		svc := &urgencyService{log: log, assignmentRepo: arepo}
		arepo.EXPECT().GetByUrgencyID(uint(1)).Return(nil, assert.AnError)
		err := svc.UnassignUrgency(1, 99, false)
		assert.Error(t, err)
	})

	t.Run("it returns error when no accepted assignment exists", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		arepo := repositories.NewMockAssignmentRepository(ctrl)
		svc := &urgencyService{log: log, assignmentRepo: arepo}
		arepo.EXPECT().GetByUrgencyID(uint(1)).Return([]model.EmergencyAssignment{}, nil)
		err := svc.UnassignUrgency(1, 99, false)
		assert.Error(t, err)
	})

	t.Run("it returns forbidden when actor is not assignee and not admin", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		arepo := repositories.NewMockAssignmentRepository(ctrl)
		svc := &urgencyService{log: log, assignmentRepo: arepo}
		arepo.EXPECT().GetByUrgencyID(uint(1)).Return([]model.EmergencyAssignment{{ID: 7, EmployeeID: 55, Status: model.AssignmentAccepted}}, nil)
		err := svc.UnassignUrgency(1, 99, false)
		assert.Error(t, err)
	})

	t.Run("it returns error when delete fails", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		arepo := repositories.NewMockAssignmentRepository(ctrl)
		svc := &urgencyService{log: log, assignmentRepo: arepo}
		arepo.EXPECT().GetByUrgencyID(uint(1)).Return([]model.EmergencyAssignment{{ID: 7, EmployeeID: 55, Status: model.AssignmentAccepted}}, nil)
		arepo.EXPECT().Delete(uint(7)).Return(assert.AnError)
		err := svc.UnassignUrgency(1, 55, false)
		assert.Error(t, err)
	})

	t.Run("it succeeds when assignee unassigns", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		arepo := repositories.NewMockAssignmentRepository(ctrl)
		svc := &urgencyService{log: log, assignmentRepo: arepo}
		arepo.EXPECT().GetByUrgencyID(uint(1)).Return([]model.EmergencyAssignment{{ID: 7, EmployeeID: 55, Status: model.AssignmentAccepted}}, nil)
		arepo.EXPECT().Delete(uint(7)).Return(nil)
		err := svc.UnassignUrgency(1, 55, false)
		assert.NoError(t, err)
	})

	t.Run("it succeeds when admin unassigns", func(t *testing.T) {
		log := utils.NewTestLogger()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		arepo := repositories.NewMockAssignmentRepository(ctrl)
		svc := &urgencyService{log: log, assignmentRepo: arepo}
		arepo.EXPECT().GetByUrgencyID(uint(1)).Return([]model.EmergencyAssignment{{ID: 7, EmployeeID: 55, Status: model.AssignmentAccepted}}, nil)
		arepo.EXPECT().Delete(uint(7)).Return(nil)
		err := svc.UnassignUrgency(1, 99, true)
		assert.NoError(t, err)
	})
}
