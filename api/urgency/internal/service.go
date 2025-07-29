package internal

//go:generate mockgen -source=service.go -destination=service_gomock.go -package=internal mountain_service/urgency/internal -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"fmt"
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/clients"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
	"github.com/pd120424d/mountain-service/api/urgency/internal/repositories"
)

type UrgencyService interface {
	CreateUrgency(urgency *model.Urgency) error
	GetAllUrgencies() ([]model.Urgency, error)
	GetUrgencyByID(id uint) (*model.Urgency, error)
	UpdateUrgency(urgency *model.Urgency) error
	DeleteUrgency(id uint) error
	ResetAllData() error
}

type urgencyService struct {
	log              utils.Logger
	repo             repositories.UrgencyRepository
	assignmentRepo   repositories.AssignmentRepository
	notificationRepo repositories.NotificationRepository
	employeeClient   clients.EmployeeClient
}

func NewUrgencyService(
	log utils.Logger,
	repo repositories.UrgencyRepository,
	assignmentRepo repositories.AssignmentRepository,
	notificationRepo repositories.NotificationRepository,
	employeeClient clients.EmployeeClient,
) UrgencyService {
	return &urgencyService{
		log:              log.WithName("urgencyService"),
		repo:             repo,
		assignmentRepo:   assignmentRepo,
		notificationRepo: notificationRepo,
		employeeClient:   employeeClient,
	}
}

func (s *urgencyService) CreateUrgency(urgency *model.Urgency) error {
	s.log.Infof("Creating urgency: %s %s", urgency.FirstName, urgency.LastName)
	err := s.repo.Create(urgency)

	if err != nil {
		s.log.Errorf("Failed to create urgency: %v", err)
		return err
	}

	ctx := context.Background()
	shiftBuffer := 1 * time.Hour // Include employees from next shift if current shift ends within 1 hour
	onCallEmployees, err := s.employeeClient.GetOnCallEmployees(ctx, shiftBuffer)
	if err != nil {
		s.log.Errorf("Failed to fetch on-call employees: %v", err)
		return fmt.Errorf("failed to fetch on-call employees: %w", err)
	}

	s.log.Infof("Found %d on-call employees for urgency %d", len(onCallEmployees), urgency.ID)

	for _, employee := range onCallEmployees {
		if err := s.createAssignmentAndNotification(urgency, employee); err != nil {
			s.log.Errorf("Failed to create assignment/notification for employee %d: %v", employee.ID, err)
			// Continue with other employees even if one fails
		}
	}

	return nil
}

func (s *urgencyService) GetAllUrgencies() ([]model.Urgency, error) {
	return s.repo.GetAll()
}

func (s *urgencyService) GetUrgencyByID(id uint) (*model.Urgency, error) {
	var urgency model.Urgency
	if err := s.repo.GetByID(id, &urgency); err != nil {
		return nil, err
	}
	return &urgency, nil
}

func (s *urgencyService) UpdateUrgency(urgency *model.Urgency) error {
	return s.repo.Update(urgency)
}

func (s *urgencyService) DeleteUrgency(id uint) error {
	return s.repo.Delete(id)
}

func (s *urgencyService) ResetAllData() error {
	return s.repo.ResetAllData()
}

func (s *urgencyService) createAssignmentAndNotification(urgency *model.Urgency, employee employeeV1.EmployeeResponse) error {
	assignment := &model.EmergencyAssignment{
		UrgencyID:  urgency.ID,
		EmployeeID: employee.ID,
		Status:     model.AssignmentPending,
		AssignedAt: time.Now(),
	}

	if err := s.assignmentRepo.Create(assignment); err != nil {
		return fmt.Errorf("failed to create assignment: %w", err)
	}

	s.log.Infof("Created assignment %d for employee %d and urgency %d", assignment.ID, employee.ID, urgency.ID)

	if employee.Phone != "" {
		smsNotification := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       employee.ID,
			NotificationType: model.NotificationSMS,
			Recipient:        employee.Phone,
			Message:          s.buildNotificationMessage(urgency, employee, "SMS"),
			Status:           model.NotificationPending,
		}

		if err := s.notificationRepo.Create(smsNotification); err != nil {
			s.log.Errorf("Failed to create SMS notification: %v", err)
		} else {
			s.log.Infof("Created SMS notification %d for employee %d", smsNotification.ID, employee.ID)
		}
	}

	if employee.Email != "" {
		emailNotification := &model.Notification{
			UrgencyID:        urgency.ID,
			EmployeeID:       employee.ID,
			NotificationType: model.NotificationEmail,
			Recipient:        employee.Email,
			Message:          s.buildNotificationMessage(urgency, employee, "Email"),
			Status:           model.NotificationPending,
		}

		if err := s.notificationRepo.Create(emailNotification); err != nil {
			s.log.Errorf("Failed to create email notification: %v", err)
		} else {
			s.log.Infof("Created email notification %d for employee %d", emailNotification.ID, employee.ID)
		}
	}

	return nil
}

func (s *urgencyService) buildNotificationMessage(urgency *model.Urgency, employee employeeV1.EmployeeResponse, notificationType string) string {
	// TODO: this is temporary and will have to be moved to a template which will be rendered by a frontend and translated accordingly
	baseMessage := fmt.Sprintf(
		"🚨 EMERGENCY ALERT 🚨\n\n"+
			"Hello %s %s,\n\n"+
			"You have been assigned to an emergency situation:\n\n"+
			"📍 Location: %s\n"+
			"📞 Contact: %s (%s)\n"+
			"📝 Description: %s\n"+
			"⚠️ Priority: %s\n\n"+
			"Please respond immediately by accepting or declining this assignment.",
		employee.FirstName,
		employee.LastName,
		urgency.Location,
		urgency.FirstName+" "+urgency.LastName,
		urgency.ContactPhone,
		urgency.Description,
		urgency.Level,
	)

	if notificationType == "SMS" {
		// TODO: this is temporary and will have to be moved to a template which will be rendered by a frontend and translated accordingly
		return fmt.Sprintf(
			"🚨 EMERGENCY: %s at %s. Contact: %s (%s). Priority: %s. Please respond ASAP.",
			urgency.Description,
			urgency.Location,
			urgency.FirstName+" "+urgency.LastName,
			urgency.ContactPhone,
			urgency.Level,
		)
	}

	return baseMessage
}
