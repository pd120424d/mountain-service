package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/pd120424d/mountain-service/api/activity/internal"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

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
	svc internal.ActivityService
}

func NewActivityHandler(log utils.Logger, svc internal.ActivityService) ActivityHandler {
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request payload: %v", err)})
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
