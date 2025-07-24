package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

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

type ActivityHandler interface {
	CreateActivity(ctx *gin.Context)
	GetActivity(ctx *gin.Context)
	ListActivities(ctx *gin.Context)
	GetActivityStats(ctx *gin.Context)
	DeleteActivity(ctx *gin.Context)
	ResetAllData(ctx *gin.Context)
}

type activityHandler struct {
	log utils.Logger
	svc ActivityService
}

func NewActivityHandler(log utils.Logger, svc ActivityService) ActivityHandler {
	return &activityHandler{log: log.WithName("activityHandler"), svc: svc}
}

// CreateActivity creates a new activity
// @Summary Create a new activity
// @Description Create a new activity in the system
// @Tags activities
// @Accept json
// @Produce json
// @Param activity body activityV1.ActivityCreateRequest true "Activity data"
// @Success 201 {object} activityV1.ActivityResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /activities [post]
func (h *activityHandler) CreateActivity(ctx *gin.Context) {
	var req activityV1.ActivityCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.log.Errorf("Failed to bind request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	response, err := h.svc.CreateActivity(&req)
	if err != nil {
		h.log.Errorf("Failed to create activity: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity", "details": err.Error()})
		return
	}

	h.log.Infof("Activity created successfully with ID: %d", response.ID)
	ctx.JSON(http.StatusCreated, response)
}

// GetActivity retrieves an activity by ID
// @Summary Get activity by ID
// @Description Get a specific activity by its ID
// @Tags activities
// @Produce json
// @Param id path int true "Activity ID"
// @Success 200 {object} activityV1.ActivityResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /activities/{id} [get]
func (h *activityHandler) GetActivity(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.log.Errorf("Invalid activity ID: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	response, err := h.svc.GetActivityByID(uint(id))
	if err != nil {
		h.log.Errorf("Failed to get activity: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ListActivities retrieves a list of activities with filtering and pagination
// @Summary List activities
// @Description Get a paginated list of activities with optional filtering
// @Tags activities
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Param type query string false "Activity type filter"
// @Param level query string false "Activity level filter"
// @Success 200 {object} activityV1.ActivityListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /activities [get]
func (h *activityHandler) ListActivities(ctx *gin.Context) {
	var req activityV1.ActivityListRequest

	// Parse query parameters
	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			req.Page = page
		}
	}
	if pageSizeStr := ctx.Query("pageSize"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			req.PageSize = pageSize
		}
	}

	if typeStr := ctx.Query("type"); typeStr != "" {
		req.Type = activityV1.ActivityType(typeStr)
	}
	if levelStr := ctx.Query("level"); levelStr != "" {
		req.Level = activityV1.ActivityLevel(levelStr)
	}
	req.TargetType = ctx.Query("targetType")
	req.StartDate = ctx.Query("startDate")
	req.EndDate = ctx.Query("endDate")

	response, err := h.svc.ListActivities(&req)
	if err != nil {
		h.log.Errorf("Failed to list activities: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list activities", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetActivityStats retrieves activity statistics
// @Summary Get activity statistics
// @Description Get comprehensive activity statistics
// @Tags activities
// @Produce json
// @Success 200 {object} activityV1.ActivityStatsResponse
// @Failure 500 {object} map[string]interface{}
// @Router /activities/stats [get]
func (h *activityHandler) GetActivityStats(ctx *gin.Context) {
	response, err := h.svc.GetActivityStats()
	if err != nil {
		h.log.Errorf("Failed to get activity stats: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get activity stats", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteActivity deletes an activity by ID
// @Summary Delete activity
// @Description Delete a specific activity by its ID
// @Tags activities
// @Param id path int true "Activity ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /activities/{id} [delete]
func (h *activityHandler) DeleteActivity(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.log.Errorf("Invalid activity ID: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	err = h.svc.DeleteActivity(uint(id))
	if err != nil {
		h.log.Errorf("Failed to delete activity: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	h.log.Infof("Activity deleted successfully with ID: %d", id)
	ctx.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})
}

// ResetAllData resets all activity data
// @Summary Reset all activity data
// @Description Delete all activities from the system
// @Tags activities
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /activities/reset [delete]
func (h *activityHandler) ResetAllData(ctx *gin.Context) {
	err := h.svc.ResetAllData()
	if err != nil {
		h.log.Errorf("Failed to reset activity data: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset activity data", "details": err.Error()})
		return
	}

	h.log.Info("All activity data reset successfully")
	ctx.JSON(http.StatusOK, gin.H{"message": "All activity data reset successfully"})
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

	response := &activityV1.ActivityStatsResponse{
		TotalActivities:      stats.TotalActivities,
		ActivitiesLast24h:    stats.ActivitiesLast24h,
		ActivitiesLast7Days:  stats.ActivitiesLast7Days,
		ActivitiesLast30Days: stats.ActivitiesLast30Days,
	}

	return response, nil
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
