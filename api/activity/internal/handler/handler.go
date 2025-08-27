package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/pd120424d/mountain-service/api/activity/internal/service"
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
	svc service.ActivityService
}

func NewActivityHandler(log utils.Logger, svc service.ActivityService) ActivityHandler {
	return &activityHandler{log: log.WithName("activityHandler"), svc: svc}
}

// CreateActivity Креирање нове активности
// @Summary Креирање нове активности
// @Description Креирање нове активности у систему
// @Tags activities
// @Accept json
// @Produce json
// @Param activity body activityV1.ActivityCreateRequest true "Activity data"
// @Success 201 {object} activityV1.ActivityResponse
// @Failure 400 {object} activityV1.ErrorResponse
// @Failure 500 {object} activityV1.ErrorResponse
// @Router /activities [post]
func (h *activityHandler) CreateActivity(ctx *gin.Context) {
	var req activityV1.ActivityCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.log.Errorf("Failed to bind request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		h.log.Errorf("validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

// GetActivity Преузимање активности по ID
// @Summary Преузимање активности по ID
// @Description Преузимање одређене активности по њеном ID
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

// ListActivities Преузимање листе активности са филтрирањем и пагинацијом
// @Summary Листа активности
// @Description Извлачење листе активности са опционим филтрирањем и страничењем
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

	if employeeIDStr := ctx.Query("employeeId"); employeeIDStr != "" {
		if employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32); err == nil {
			employeeIDUint := uint(employeeID)
			req.EmployeeID = &employeeIDUint
		}
	}
	if urgencyIDStr := ctx.Query("urgencyId"); urgencyIDStr != "" {
		if urgencyID, err := strconv.ParseUint(urgencyIDStr, 10, 32); err == nil {
			urgencyIDUint := uint(urgencyID)
			req.UrgencyID = &urgencyIDUint
		}
	}
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

// GetActivityStats Преузимање статистика активности
// @Summary Статистике активности
// @Description Преузимање свеобухватних статистика активности
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

// DeleteActivity Брисање активности по ID
// @Summary Брисање активности
// @Description Брисање одређене активности по њеном ID
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
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity", "details": err.Error()})
		return
	}

	h.log.Infof("Activity deleted successfully with ID: %d", id)
	ctx.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})
}

// ResetAllData Ресетовање свих података о активностима
// @Summary Ресетовање свих података о активностима
// @Description Брисање свих активности из система
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
