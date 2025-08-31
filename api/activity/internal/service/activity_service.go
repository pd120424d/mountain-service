package service

//go:generate mockgen -source=activity_service.go -destination=activity_service_gomock.go -package=service mountain_service/activity/internal/service -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
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
	CreateActivity(ctx context.Context, req *activityV1.ActivityCreateRequest) (*activityV1.ActivityResponse, error)
	GetActivityByID(ctx context.Context, id uint) (*activityV1.ActivityResponse, error)
	ListActivities(ctx context.Context, req *activityV1.ActivityListRequest) (*activityV1.ActivityListResponse, error)
	GetActivityStats(ctx context.Context) (*activityV1.ActivityStatsResponse, error)
	DeleteActivity(ctx context.Context, id uint) error
	ResetAllData(ctx context.Context) error

	LogActivity(ctx context.Context, description string, employeeID, urgencyID uint) error
}

type activityService struct {
	log  utils.Logger
	repo repositories.ActivityRepository
}

func NewActivityService(log utils.Logger, repo repositories.ActivityRepository) ActivityService {
	return &activityService{log: log.WithName("activityService"), repo: repo}
}

func (s *activityService) CreateActivity(ctx context.Context, req *activityV1.ActivityCreateRequest) (*activityV1.ActivityResponse, error) {
	log := s.log.WithContext(ctx)
	if req == nil {
		log.Error("Activity create request is nil")
		return nil, commonv1.NewAppError("VALIDATION.INVALID_REQUEST", "request cannot be nil", nil)
	}

	log.Infof("Creating activity: %s", req.ToString())

	if err := req.Validate(); err != nil {
		log.Errorf("Activity validation failed: %v", err)
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

	if err := s.repo.CreateWithOutbox(ctx, activity, (*models.OutboxEvent)(event)); err != nil {
		log.Errorf("Failed to create activity with outbox: %v", err)
		return nil, commonv1.NewAppError("ACTIVITY_ERRORS.CREATE_FAILED", "failed to create activity", map[string]interface{}{"cause": err.Error()})
	}

	response := activity.ToResponse()

	log.Infof("Activity created successfully with ID: %d", activity.ID)
	return &response, nil
}

func (s *activityService) GetActivityByID(ctx context.Context, id uint) (*activityV1.ActivityResponse, error) {
	log := s.log.WithContext(ctx)
	if id == 0 {
		log.Error("Invalid activity ID: 0")
		return nil, commonv1.NewAppError("VALIDATION.INVALID_ID", "invalid activity ID: cannot be zero", nil)
	}

	log.Infof("Getting activity by ID: %d", id)

	activity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Errorf("Failed to get activity: %v", err)
		return nil, commonv1.NewAppError("ACTIVITY_ERRORS.NOT_FOUND", "failed to get activity", map[string]interface{}{"cause": err.Error()})
	}

	response := activity.ToResponse()
	return &response, nil
}

func (s *activityService) DeleteActivity(ctx context.Context, id uint) error {
	log := s.log.WithContext(ctx)
	if id == 0 {
		log.Error("Invalid activity ID: 0")
		return commonv1.NewAppError("VALIDATION.INVALID_ID", "invalid activity ID: cannot be zero", nil)
	}

	log.Infof("Deleting activity with ID: %d", id)

	if err := s.repo.Delete(ctx, id); err != nil {
		log.Errorf("Failed to delete activity: %v", err)
		return commonv1.NewAppError("ACTIVITY_ERRORS.DELETE_FAILED", "failed to delete activity", map[string]interface{}{"cause": err.Error()})
	}

	log.Infof("Activity deleted successfully with ID: %d", id)
	return nil
}

func (s *activityService) ListActivities(ctx context.Context, req *activityV1.ActivityListRequest) (*activityV1.ActivityListResponse, error) {
	log := s.log.WithContext(ctx)
	log.Infof("Listing activities with filters: %+v", req)

	// Validate request
	if err := req.Validate(); err != nil {
		log.Errorf("Activity list validation failed: %v", err)
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

	activities, total, err := s.repo.List(ctx, filter)
	if err != nil {
		log.Errorf("Failed to list activities: %v", err)
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

	log.Infof("Listed %d activities out of %d total", len(activities), total)
	return response, nil
}

func (s *activityService) GetActivityStats(ctx context.Context) (*activityV1.ActivityStatsResponse, error) {
	log := s.log.WithContext(ctx)
	log.Info("Getting activity statistics")

	stats, err := s.repo.GetStats(ctx)
	if err != nil {
		log.Errorf("Failed to get activity stats: %v", err)
		return nil, commonv1.NewAppError("ACTIVITY_ERRORS.STATS_FAILED", "failed to get activity stats", map[string]interface{}{"cause": err.Error()})
	}

	response := stats.ToResponse()
	return &response, nil
}

func (s *activityService) ResetAllData(ctx context.Context) error {
	log := s.log.WithContext(ctx)
	log.Info("Resetting all activity data")

	if err := s.repo.ResetAllData(ctx); err != nil {
		log.Errorf("Failed to reset activity data: %v", err)
		return commonv1.NewAppError("ACTIVITY_ERRORS.RESET_FAILED", "failed to reset activity data", map[string]interface{}{"cause": err.Error()})
	}

	log.Info("All activity data reset successfully")
	return nil
}

func (s *activityService) LogActivity(ctx context.Context, description string, employeeID, urgencyID uint) error {
	req := &activityV1.ActivityCreateRequest{
		Description: description,
		EmployeeID:  employeeID,
		UrgencyID:   urgencyID,
	}

	_, err := s.CreateActivity(ctx, req)
	return err
}
