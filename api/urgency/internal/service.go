package internal

//go:generate mockgen -source=service.go -destination=service_gomock.go -package=internal mountain_service/urgency/internal -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"github.com/pd120424d/mountain-service/api/shared/utils"
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
	log  utils.Logger
	repo repositories.UrgencyRepository
}

func NewUrgencyService(log utils.Logger, repo repositories.UrgencyRepository) UrgencyService {
	return &urgencyService{log: log.WithName("urgencyService"), repo: repo}
}

func (s *urgencyService) CreateUrgency(urgency *model.Urgency) error {
	s.log.Infof("Creating urgency: %s", urgency.Name)
	err := s.repo.Create(urgency)

	if err != nil {
		s.log.Errorf("Failed to create urgency: %v", err)
		return err
	}

	// TODO: Fetch employees which are on call from the employee service
	// TODO: Create assignments for the urgency and assign them to the employee
	// TODO: Try to send notifications for the assignments

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
