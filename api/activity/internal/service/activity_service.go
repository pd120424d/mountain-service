package service

//go:generate mockgen -source=activity_service.go -destination=activity_service_gomock.go -package=service mountain_service/activity/internal/service -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pd120424d/mountain-service/api/activity/internal/model"
	"github.com/pd120424d/mountain-service/api/activity/internal/repositories"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	commonv1 "github.com/pd120424d/mountain-service/api/contracts/common/v1"
	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
	"github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type ActivityService interface {
	CreateActivity(ctx context.Context, req *activityV1.ActivityCreateRequest) (*activityV1.ActivityResponse, error)
	GetActivityByID(ctx context.Context, id uint) (*activityV1.ActivityResponse, error)
	ListActivities(ctx context.Context, req *activityV1.ActivityListRequest) (*activityV1.ActivityListResponse, error)
	DeleteActivity(ctx context.Context, id uint) error
	ResetAllData(ctx context.Context) error

	LogActivity(ctx context.Context, description string, employeeID, urgencyID uint) error
}

type activityService struct {
	log  utils.Logger
	repo repositories.ActivityRepository
	// clients are optional; when nil, enrichment is skipped
	urgencyClient interface {
		GetUrgencyByID(ctx context.Context, id uint) (*urgencyV1.UrgencyResponse, error)
	}
	employeeClient interface {
		GetEmployeeByID(ctx context.Context, employeeID uint) (*employeeV1.EmployeeResponse, error)
	}
}

func NewActivityService(log utils.Logger, repo repositories.ActivityRepository, urgencyClient interface {
	GetUrgencyByID(ctx context.Context, id uint) (*urgencyV1.UrgencyResponse, error)
}) ActivityService {
	return &activityService{log: log.WithName("activityService"), repo: repo, urgencyClient: urgencyClient}
}

// NewActivityServiceWithDeps allows injecting both urgency and employee clients
func NewActivityServiceWithDeps(
	log utils.Logger,
	repo repositories.ActivityRepository,
	urgencyClient interface {
		GetUrgencyByID(ctx context.Context, id uint) (*urgencyV1.UrgencyResponse, error)
	},
	employeeClient interface {
		GetEmployeeByID(ctx context.Context, employeeID uint) (*employeeV1.EmployeeResponse, error)
	},
) ActivityService {
	return &activityService{log: log.WithName("activityService"), repo: repo, urgencyClient: urgencyClient, employeeClient: employeeClient}
}

func (s *activityService) CreateActivity(ctx context.Context, req *activityV1.ActivityCreateRequest) (*activityV1.ActivityResponse, error) {
	log := s.log.WithContext(ctx)
	defer utils.TimeOperation(log, "ActivityService.CreateActivity")()

	if req == nil {
		log.Error("Activity create request is nil")
		return nil, commonv1.NewAppError("VALIDATION.INVALID_REQUEST", "request cannot be nil", nil)
	}

	log.Infof("Creating activity: %s", req.ToString())

	if err := req.Validate(); err != nil {
		log.Errorf("Activity validation failed: %v", err)
		return nil, commonv1.NewAppError("VALIDATION.INVALID_REQUEST", fmt.Sprintf("validation failed: %v", err), nil)
	}

	var urgencyTitle string
	var urgencyLevel string
	if s.urgencyClient != nil {
		urg, err := s.urgencyClient.GetUrgencyByID(ctx, req.UrgencyID)
		if err != nil {
			log.Errorf("Failed to fetch urgency %d: %v", req.UrgencyID, err)
			return nil, commonv1.NewAppError("ACTIVITY_ERRORS.URGENCY_FETCH_FAILED", "failed to validate urgency", map[string]interface{}{"cause": err.Error()})
		}
		if urg == nil || urg.Status != urgencyV1.InProgress {
			log.Warnf("CreateActivity denied: urgency invalid state. urgencyId=%d status=%v", req.UrgencyID, func() string {
				if urg != nil {
					return string(urg.Status)
				}
				return ""
			}())
			return nil, commonv1.NewAppError("ACTIVITY_ERRORS.INVALID_URGENCY_STATE", "activities can be added only to in_progress urgencies", map[string]interface{}{"urgencyId": req.UrgencyID, "status": func() string {
				if urg != nil {
					return string(urg.Status)
				}
				return ""
			}()})
		}

		urgencyTitle = fmt.Sprintf("%s %s", strings.TrimSpace(urg.FirstName), strings.TrimSpace(urg.LastName))
		urgencyTitle = strings.TrimSpace(urgencyTitle)
		urgencyLevel = string(urg.Level)

		actorID := ctx.Value("employeeID")
		role := ctx.Value("role")
		roleStr, _ := role.(string)

		// Validation: urgency must have assignee for non-admins
		if urg.AssignedEmployeeId == nil && roleStr != "Administrator" {
			log.Warnf("CreateActivity denied: missing assignee. urgencyId=%d actorId=%v role=%s", req.UrgencyID, actorID, roleStr)
			return nil, commonv1.NewAppError("VALIDATION.MISSING_ASSIGNEE", "urgency must have an assigned employee before adding activities", map[string]interface{}{"urgencyId": req.UrgencyID})
		}

		// Enforce assignee-only unless admin
		if actorID != nil && roleStr != "Administrator" && urg.AssignedEmployeeId != nil {
			if actID, ok := actorID.(uint); ok {
				if *urg.AssignedEmployeeId != actID {
					log.Warnf("CreateActivity denied: actor is not assignee. urgencyId=%d assignedEmployeeId=%d actorId=%d role=%s", req.UrgencyID, *urg.AssignedEmployeeId, actID, roleStr)
					return nil, commonv1.NewAppError("AUTH_ERRORS.FORBIDDEN", "only assignee or admin can add activities", map[string]interface{}{"urgencyId": req.UrgencyID, "actorId": actID, "assignedEmployeeId": *urg.AssignedEmployeeId})
				}
			}
		}

		// Override payload employeeId with actorId when present, to prevent spoofing (non-admin only)
		if roleStr != "Administrator" {
			if actID, ok := actorID.(uint); ok {
				req.EmployeeID = actID
			}
		}

		log.Infof("CreateActivity allowed: urgencyId=%d status=%s assignedEmployeeId=%v actorId=%v role=%s", req.UrgencyID, urg.Status, urg.AssignedEmployeeId, actorID, roleStr)
	}

	activity := model.FromCreateRequest(req)

	// Denormalize employee name for read-model if client is available
	var employeeName string
	if s.employeeClient != nil {
		if emp, err := s.employeeClient.GetEmployeeByID(ctx, req.EmployeeID); err != nil {
			log.Warnf("Failed to fetch employee %d for denormalization: %v", req.EmployeeID, err)
		} else if emp != nil {
			fullName := strings.TrimSpace(strings.TrimSpace(emp.FirstName) + " " + strings.TrimSpace(emp.LastName))
			employeeName = fullName
		}
	}

	// Build outbox event payload for CQRS (enriched)
	event := activityV1.CreateOutboxEvent(
		activity.ID,
		activityV1.ActivityEvent{
			Type:         "CREATE", // used by Firestore updater to determine action
			ActivityID:   activity.ID,
			UrgencyID:    activity.UrgencyID,
			EmployeeID:   activity.EmployeeID,
			Description:  activity.Description,
			CreatedAt:    activity.CreatedAt,
			EmployeeName: employeeName,
			UrgencyTitle: urgencyTitle,
			UrgencyLevel: urgencyLevel,
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
	defer utils.TimeOperation(log, "ActivityService.GetActivityByID")()
	log.Infof("Getting activity by ID: %d", id)

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

	log.Infof("Activity retrieved successfully with ID: %d", id)
	response := activity.ToResponse()
	return &response, nil
}

func (s *activityService) DeleteActivity(ctx context.Context, id uint) error {
	log := s.log.WithContext(ctx)
	defer utils.TimeOperation(log, "ActivityService.DeleteActivity")()
	log.Infof("Deleting activity with ID: %d", id)

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
	defer utils.TimeOperation(log, "ActivityService.ListActivities")()
	log.Infof("Listing activities with filters: %+v", req)

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

func (s *activityService) ResetAllData(ctx context.Context) error {
	log := s.log.WithContext(ctx)
	defer utils.TimeOperation(log, "ActivityService.ResetAllData")()
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
