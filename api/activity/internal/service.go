package internal

//go:generate mockgen -source=service.go -destination=service_gomock.go -package=internal mountain_service/activity/internal -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"fmt"
	"time"

	"github.com/pd120424d/mountain-service/api/activity/internal/model"
	"github.com/pd120424d/mountain-service/api/activity/internal/repositories"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type ActivityService interface {
	CreateActivity(req *activityV1.ActivityCreateRequest) (*activityV1.ActivityResponse, error)
	GetActivityByID(id uint) (*activityV1.ActivityResponse, error)
	ListActivities(req *activityV1.ActivityListRequest) (*activityV1.ActivityListResponse, error)
	GetActivityStats() (*activityV1.ActivityStatsResponse, error)
	DeleteActivity(id uint) error
	ResetAllData() error

	// Helper methods for logging activities
	LogEmployeeActivity(activityType activityV1.ActivityType, level activityV1.ActivityLevel, title, description string, employeeID uint) error
	LogUrgencyActivity(activityType activityV1.ActivityType, level activityV1.ActivityLevel, title, description string, urgencyID uint) error
	LogSystemActivity(activityType activityV1.ActivityType, level activityV1.ActivityLevel, title, description string) error
}

// Service implementation
type activityService struct {
	log  utils.Logger
	repo repositories.ActivityRepository
}

func NewActivityService(log utils.Logger, repo repositories.ActivityRepository) ActivityService {
	return &activityService{log: log.WithName("activityService"), repo: repo}
}

func (s *activityService) CreateActivity(req *activityV1.ActivityCreateRequest) (*activityV1.ActivityResponse, error) {
	s.log.Infof("Creating activity: %s", req.ToString())

	// Validate request
	if err := req.Validate(); err != nil {
		s.log.Errorf("Activity validation failed: %v", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert DTO to entity
	activity := &model.Activity{
		Type:        req.Type,
		Level:       req.Level,
		Title:       req.Title,
		Description: req.Description,
		ActorID:     req.ActorID,
		ActorName:   req.ActorName,
		TargetID:    req.TargetID,
		TargetType:  req.TargetType,
		Metadata:    req.Metadata,
		CreatedAt:   time.Now(),
	}

	// Save to repository
	if err := s.repo.Create(activity); err != nil {
		s.log.Errorf("Failed to create activity: %v", err)
		return nil, fmt.Errorf("failed to create activity: %w", err)
	}

	// Convert entity to response DTO
	response := &activityV1.ActivityResponse{
		ID:          activity.ID,
		Type:        activity.Type,
		Level:       activity.Level,
		Title:       activity.Title,
		Description: activity.Description,
		ActorID:     activity.ActorID,
		ActorName:   activity.ActorName,
		TargetID:    activity.TargetID,
		TargetType:  activity.TargetType,
		Metadata:    activity.Metadata,
		CreatedAt:   activity.CreatedAt.Format(time.RFC3339),
	}

	s.log.Infof("Activity created successfully with ID: %d", activity.ID)
	return response, nil
}

func (s *activityService) GetActivityByID(id uint) (*activityV1.ActivityResponse, error) {
	s.log.Infof("Getting activity by ID: %d", id)

	activity, err := s.repo.GetByID(id)
	if err != nil {
		s.log.Errorf("Failed to get activity: %v", err)
		return nil, fmt.Errorf("failed to get activity: %w", err)
	}

	response := &activityV1.ActivityResponse{
		ID:          activity.ID,
		Type:        activity.Type,
		Level:       activity.Level,
		Title:       activity.Title,
		Description: activity.Description,
		ActorID:     activity.ActorID,
		ActorName:   activity.ActorName,
		TargetID:    activity.TargetID,
		TargetType:  activity.TargetType,
		Metadata:    activity.Metadata,
		CreatedAt:   activity.CreatedAt.Format(time.RFC3339),
	}

	return response, nil
}

func (s *activityService) DeleteActivity(id uint) error {
	s.log.Infof("Deleting activity with ID: %d", id)

	if err := s.repo.Delete(id); err != nil {
		s.log.Errorf("Failed to delete activity: %v", err)
		return fmt.Errorf("failed to delete activity: %w", err)
	}

	s.log.Infof("Activity deleted successfully with ID: %d", id)
	return nil
}

func (s *activityService) ResetAllData() error {
	s.log.Info("Resetting all activity data")

	if err := s.repo.ResetAllData(); err != nil {
		s.log.Errorf("Failed to reset activity data: %v", err)
		return fmt.Errorf("failed to reset activity data: %w", err)
	}

	s.log.Info("All activity data reset successfully")
	return nil
}

func (s *activityService) ListActivities(req *activityV1.ActivityListRequest) (*activityV1.ActivityListResponse, error) {
	s.log.Infof("Listing activities with filters: %+v", req)

	// Validate request
	if err := req.Validate(); err != nil {
		s.log.Errorf("Activity list validation failed: %v", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert DTO to filter
	filter := &model.ActivityFilter{
		ActorID:  req.ActorID,
		TargetID: req.TargetID,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// Handle optional string fields
	if req.Type != "" {
		filter.Type = &req.Type
	}
	if req.Level != "" {
		filter.Level = &req.Level
	}
	if req.TargetType != "" {
		filter.TargetType = &req.TargetType
	}

	// Handle date parsing
	if req.StartDate != "" {
		if startDate, err := time.Parse(time.RFC3339, req.StartDate); err == nil {
			filter.StartDate = &startDate
		}
	}
	if req.EndDate != "" {
		if endDate, err := time.Parse(time.RFC3339, req.EndDate); err == nil {
			filter.EndDate = &endDate
		}
	}

	activities, total, err := s.repo.List(filter)
	if err != nil {
		s.log.Errorf("Failed to list activities: %v", err)
		return nil, fmt.Errorf("failed to list activities: %w", err)
	}

	// Convert entities to response DTOs
	activityResponses := make([]activityV1.ActivityResponse, len(activities))
	for i, activity := range activities {
		activityResponses[i] = activityV1.ActivityResponse{
			ID:          activity.ID,
			Type:        activity.Type,
			Level:       activity.Level,
			Title:       activity.Title,
			Description: activity.Description,
			ActorID:     activity.ActorID,
			ActorName:   activity.ActorName,
			TargetID:    activity.TargetID,
			TargetType:  activity.TargetType,
			Metadata:    activity.Metadata,
			CreatedAt:   activity.CreatedAt.Format(time.RFC3339),
		}
	}

	response := &activityV1.ActivityListResponse{
		Activities: activityResponses,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	s.log.Infof("Listed %d activities out of %d total", len(activities), total)
	return response, nil
}

func (s *activityService) GetActivityStats() (*activityV1.ActivityStatsResponse, error) {
	s.log.Info("Getting activity statistics")

	stats, err := s.repo.GetStats()
	if err != nil {
		s.log.Errorf("Failed to get activity stats: %v", err)
		return nil, fmt.Errorf("failed to get activity stats: %w", err)
	}

	// Use the model's ToResponse method for proper conversion
	response := stats.ToResponse()
	return &response, nil
}

// Helper methods for logging activities
func (s *activityService) LogEmployeeActivity(activityType activityV1.ActivityType, level activityV1.ActivityLevel, title, description string, employeeID uint) error {
	req := &activityV1.ActivityCreateRequest{
		Type:        activityType,
		Level:       level,
		Title:       title,
		Description: description,
		ActorID:     &employeeID,
		TargetID:    &employeeID,
		TargetType:  "employee",
	}

	_, err := s.CreateActivity(req)
	return err
}

func (s *activityService) LogUrgencyActivity(activityType activityV1.ActivityType, level activityV1.ActivityLevel, title, description string, urgencyID uint) error {
	req := &activityV1.ActivityCreateRequest{
		Type:        activityType,
		Level:       level,
		Title:       title,
		Description: description,
		TargetID:    &urgencyID,
		TargetType:  "urgency",
	}

	_, err := s.CreateActivity(req)
	return err
}

func (s *activityService) LogSystemActivity(activityType activityV1.ActivityType, level activityV1.ActivityLevel, title, description string) error {
	req := &activityV1.ActivityCreateRequest{
		Type:        activityType,
		Level:       level,
		Title:       title,
		Description: description,
		ActorName:   "system",
		TargetType:  "system",
	}

	_, err := s.CreateActivity(req)
	return err
}
