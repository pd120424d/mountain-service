package service

//go:generate mockgen -source=activity_service.go -destination=activity_service_gomock.go -package=service mountain_service/activity/internal/service -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"fmt"
	"time"

	"github.com/pd120424d/mountain-service/api/activity/internal/model"
	"github.com/pd120424d/mountain-service/api/activity/internal/repositories"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	commonv1 "github.com/pd120424d/mountain-service/api/contracts/common/v1"
	"github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type ActivityService interface {
	CreateActivity(req *activityV1.ActivityCreateRequest) (*activityV1.ActivityResponse, error)
	GetActivityByID(id uint) (*activityV1.ActivityResponse, error)
	ListActivities(req *activityV1.ActivityListRequest) (*activityV1.ActivityListResponse, error)
	GetActivityStats() (*activityV1.ActivityStatsResponse, error)
	DeleteActivity(id uint) error
	ResetAllData() error

	LogActivity(description string, employeeID, urgencyID uint) error
}

type activityService struct {
	log  utils.Logger
	repo repositories.ActivityRepository
}

func NewActivityService(log utils.Logger, repo repositories.ActivityRepository) ActivityService {
	return &activityService{log: log.WithName("activityService"), repo: repo}
}

func (s *activityService) CreateActivity(req *activityV1.ActivityCreateRequest) (*activityV1.ActivityResponse, error) {
	if req == nil {
		s.log.Error("Activity create request is nil")
		return nil, commonv1.NewAppError("VALIDATION.INVALID_REQUEST", "request cannot be nil", nil)
	}

	s.log.Infof("Creating activity: %s", req.ToString())

	if err := req.Validate(); err != nil {
		s.log.Errorf("Activity validation failed: %v", err)
		return nil, commonv1.NewAppError("VALIDATION.INVALID_REQUEST", fmt.Sprintf("validation failed: %v", err), nil)
	}

	activity := model.FromCreateRequest(req)

	// Build outbox event payload for CQRS
	event := activityV1.CreateOutboxEvent(
		activityV1.ActivityEventCreated,
		activity.ID,
		activityV1.ActivityEvent{
			Type:        string(activityV1.ActivityEventCreated),
			ActivityID:  activity.ID,
			UrgencyID:   activity.UrgencyID,
			EmployeeID:  activity.EmployeeID,
			Description: activity.Description,
			CreatedAt:   activity.CreatedAt,
		},
	)

	if err := s.repo.CreateWithOutbox(activity, (*models.OutboxEvent)(event)); err != nil {
		s.log.Errorf("Failed to create activity with outbox: %v", err)
		return nil, commonv1.NewAppError("ACTIVITY_ERRORS.CREATE_FAILED", "failed to create activity", map[string]interface{}{"cause": err.Error()})
	}

	response := activity.ToResponse()

	s.log.Infof("Activity created successfully with ID: %d", activity.ID)
	return &response, nil
}

func (s *activityService) GetActivityByID(id uint) (*activityV1.ActivityResponse, error) {
	if id == 0 {
		s.log.Error("Invalid activity ID: 0")
		return nil, commonv1.NewAppError("VALIDATION.INVALID_ID", "invalid activity ID: cannot be zero", nil)
	}

	s.log.Infof("Getting activity by ID: %d", id)

	activity, err := s.repo.GetByID(id)
	if err != nil {
		s.log.Errorf("Failed to get activity: %v", err)
		return nil, commonv1.NewAppError("ACTIVITY_ERRORS.NOT_FOUND", "failed to get activity", map[string]interface{}{"cause": err.Error()})
	}

	response := activity.ToResponse()
	return &response, nil
}

func (s *activityService) DeleteActivity(id uint) error {
	if id == 0 {
		s.log.Error("Invalid activity ID: 0")
		return commonv1.NewAppError("VALIDATION.INVALID_ID", "invalid activity ID: cannot be zero", nil)
	}

	s.log.Infof("Deleting activity with ID: %d", id)

	if err := s.repo.Delete(id); err != nil {
		s.log.Errorf("Failed to delete activity: %v", err)
		return commonv1.NewAppError("ACTIVITY_ERRORS.DELETE_FAILED", "failed to delete activity", map[string]interface{}{"cause": err.Error()})
	}

	s.log.Infof("Activity deleted successfully with ID: %d", id)
	return nil
}

func (s *activityService) ListActivities(req *activityV1.ActivityListRequest) (*activityV1.ActivityListResponse, error) {
	s.log.Infof("Listing activities with filters: %+v", req)

	// Validate request
	if err := req.Validate(); err != nil {
		s.log.Errorf("Activity list validation failed: %v", err)
		return nil, commonv1.NewAppError("VALIDATION.INVALID_REQUEST", fmt.Sprintf("validation failed: %v", err), nil)
	}

	// Convert DTO to filter
	filter := &model.ActivityFilter{
		EmployeeID: req.EmployeeID,
		UrgencyID:  req.UrgencyID,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	// No additional filters needed for simplified model

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
		return nil, commonv1.NewAppError("ACTIVITY_ERRORS.LIST_FAILED", "failed to list activities", map[string]interface{}{"cause": err.Error()})
	}

	// Convert entities to response DTOs
	activityResponses := make([]activityV1.ActivityResponse, len(activities))
	for i, activity := range activities {
		activityResponses[i] = activity.ToResponse()
	}

	// Calculate total pages
	totalPages := 0
	if req.PageSize > 0 {
		totalPages = int((total + int64(req.PageSize) - 1) / int64(req.PageSize)) // Ceiling division
	}

	response := &activityV1.ActivityListResponse{
		Activities: activityResponses,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}

	s.log.Infof("Listed %d activities out of %d total", len(activities), total)
	return response, nil
}

func (s *activityService) GetActivityStats() (*activityV1.ActivityStatsResponse, error) {
	s.log.Info("Getting activity statistics")

	stats, err := s.repo.GetStats()
	if err != nil {
		s.log.Errorf("Failed to get activity stats: %v", err)
		return nil, commonv1.NewAppError("ACTIVITY_ERRORS.STATS_FAILED", "failed to get activity stats", map[string]interface{}{"cause": err.Error()})
	}

	response := stats.ToResponse()
	return &response, nil
}

func (s *activityService) ResetAllData() error {
	s.log.Info("Resetting all activity data")

	if err := s.repo.ResetAllData(); err != nil {
		s.log.Errorf("Failed to reset activity data: %v", err)
		return commonv1.NewAppError("ACTIVITY_ERRORS.RESET_FAILED", "failed to reset activity data", map[string]interface{}{"cause": err.Error()})
	}

	s.log.Info("All activity data reset successfully")
	return nil
}

func (s *activityService) LogActivity(description string, employeeID, urgencyID uint) error {
	req := &activityV1.ActivityCreateRequest{
		Description: description,
		EmployeeID:  employeeID,
		UrgencyID:   urgencyID,
	}

	_, err := s.CreateActivity(req)
	return err
}
