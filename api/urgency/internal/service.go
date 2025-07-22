package internal

//go:generate mockgen -source=service.go -destination=service_gomock.go -package=internal mountain_service/urgency/internal -imports=gomock=go.uber.org/mock/gomock -typed

import (
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
	repo repositories.UrgencyRepository
}

func NewUrgencyService(repo repositories.UrgencyRepository) UrgencyService {
	return &urgencyService{repo: repo}
}

func (s *urgencyService) CreateUrgency(urgency *model.Urgency) error {
	return s.repo.Create(urgency)
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
