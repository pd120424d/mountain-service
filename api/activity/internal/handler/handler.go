package handler

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/pd120424d/mountain-service/api/activity/internal/service"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	sharedModels "github.com/pd120424d/mountain-service/api/shared/models"
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
	log       utils.Logger
	svc       service.ActivityService
	readModel service.FirestoreService
}

func NewActivityHandler(log utils.Logger, svc service.ActivityService, readModel service.FirestoreService) ActivityHandler {
	return &activityHandler{log: log.WithName("activityHandler"), svc: svc, readModel: readModel}
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
	log := h.log.WithContext(ctx.Request.Context())
	defer utils.TimeOperation(log, "ActivityHandler.CreateActivity")()
	log.Info("Received Create Activity request")

	var req activityV1.ActivityCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorf("Failed to bind request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		log.Errorf("validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.svc.CreateActivity(ctx.Request.Context(), &req)
	if err != nil {
		log.Errorf("Failed to create activity: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity", "details": err.Error()})
		return
	}

	log.Infof("Successfully created activity with ID: %d", response.ID)
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
	log := h.log.WithContext(ctx.Request.Context())
	defer utils.TimeOperation(log, "ActivityHandler.GetActivity")()
	log.Info("Received Get Activity request")

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Errorf("Invalid activity ID: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	response, err := h.svc.GetActivityByID(ctx.Request.Context(), uint(id))
	if err != nil {
		log.Errorf("Failed to get activity: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	log.Infof("Successfully retrieved activity with ID %d", id)
	ctx.JSON(http.StatusOK, response)
}

func buildActivityListRequest(ctx *gin.Context) activityV1.ActivityListRequest {
	var req activityV1.ActivityListRequest
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
	req.PageToken = ctx.Query("pageToken")
	return req
}

// Local cursor token encoder for exposing nextPageToken in page-based responses
// Matches service-side format: base64(JSON{"createdAt": RFC3339 UTC})
type cursorToken struct {
	CreatedAt string `json:"createdAt"`
	ID        uint   `json:"id,omitempty"`
}

func encodeCursorToken(t time.Time, id uint) string {
	if t.IsZero() {
		return ""
	}
	b, _ := json.Marshal(cursorToken{CreatedAt: t.UTC().Format(time.RFC3339), ID: id})
	return base64.RawURLEncoding.EncodeToString(b)
}

// ListActivities Преузимање листе активности са филтрирањем и пагинацијом
// @Summary Листа активности
// @Description Извлачење листе активности са опционим филтрирањем и страничењем
// @Tags activities
// @Produce json
// @Param pageToken query string false "Курсор за наставак (вредност nextPageToken из претходног одговора)"
// @Param page query int false "Број стране (за класично страничење)" default(1)
// @Param pageSize query int false "Број ставки по страни" default(10)
// @Param urgencyId query int false "Филтер по ургенцији"
// @Param employeeId query int false "Филтер по запосленом"
// @Param startDate query string false "Почетни датум (RFC3339)"
// @Param endDate query string false "Крајњи датум (RFC3339)"
// @Success 200 {object} activityV1.ActivityListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /activities [get]
func (h *activityHandler) ListActivities(ctx *gin.Context) {
	var req activityV1.ActivityListRequest
	log := h.log.WithContext(ctx.Request.Context())
	defer utils.TimeOperation(log, "ActivityHandler.ListActivities")()
	log.Info("Received List Activities request")

	req = buildActivityListRequest(ctx)

	// Cursor-based pagination (preferred when pageToken provided)
	if h.readModel != nil && req.PageToken != "" {
		size := req.PageSize
		if size <= 0 {
			size = 10
		}
		if size > 100 {
			size = 100
		}

		var (
			activities []sharedModels.Activity
			nextToken  string
			err        error
		)
		if req.UrgencyID != nil {
			activities, nextToken, err = h.readModel.ListByUrgencyCursor(ctx.Request.Context(), *req.UrgencyID, size, req.PageToken)
		} else {
			activities, nextToken, err = h.readModel.ListAllCursor(ctx.Request.Context(), size, req.PageToken)
		}
		if err != nil {
			log.Warnf("Cursor read-model fetch failed, falling back: %v", err)
		} else {
			resp := &activityV1.ActivityListResponse{Activities: make([]activityV1.ActivityResponse, 0, len(activities))}
			for _, a := range activities {
				ar := a.ToResponse()
				resp.Activities = append(resp.Activities, *ar)
			}
			resp.Total = int64(len(activities))
			resp.PageSize = size
			resp.NextPageToken = nextToken
			log.Infof("Listed %d activities using Firestore cursor. nextToken set? %v", len(resp.Activities), nextToken != "")
			ctx.JSON(http.StatusOK, resp)
			return
		}
	}

	if h.readModel != nil && req.UrgencyID != nil {
		page := req.Page
		size := req.PageSize
		if page <= 0 {
			page = 1
		}
		if size <= 0 {
			size = 10
		}

		limit := page * size
		activities, err := h.readModel.ListByUrgency(ctx.Request.Context(), *req.UrgencyID, limit)
		if err != nil {
			// If Firestore has no data we don't want to fail
			// but just return empty result
			if strings.Contains(strings.ToLower(err.Error()), "no more items") {
				resp := &activityV1.ActivityListResponse{Activities: []activityV1.ActivityResponse{}, Total: 0, Page: page, PageSize: size, TotalPages: 1}
				log.Infof("No Firestore items for urgency %d; returning empty list from read model", *req.UrgencyID)
				ctx.JSON(http.StatusOK, resp)
				return
			}
			log.Warnf("Read-model fetch failed, falling back to DB: %v", err)
		} else {
			start := (page - 1) * size
			end := start + size
			if start > len(activities) {
				start = len(activities)
			}
			if end > len(activities) {
				end = len(activities)
			}

			resp := &activityV1.ActivityListResponse{Activities: make([]activityV1.ActivityResponse, 0, end-start)}
			for _, a := range activities[start:end] {
				ar := a.ToResponse()
				resp.Activities = append(resp.Activities, *ar)
			}
			// Approximate total (how many we fetched up to this page); UI can treat as lower bound
			resp.Total = int64(len(activities))
			resp.Page = page
			resp.PageSize = size
			if resp.PageSize > 0 {
				resp.TotalPages = int((resp.Total + int64(resp.PageSize) - 1) / int64(resp.PageSize))
			} else {
				resp.TotalPages = 1
			}

			// Provide a nextPageToken for infinite scroll to switch to cursor mode
			if end-start > 0 {
				last := activities[end-1]
				resp.NextPageToken = encodeCursorToken(last.CreatedAt, last.ID)
			}

			log.Infof("Listed %d activities out of approx %d total. Used Firestore read model.", len(resp.Activities), resp.Total)
			ctx.JSON(http.StatusOK, resp)
			return
		}

	} else if h.readModel != nil && req.UrgencyID == nil {
		// No urgencyId: still prefer Firestore read model for the main feed (supports pagination)
		page := req.Page
		size := req.PageSize
		if page <= 0 {
			page = 1
		}
		if size <= 0 {
			size = 10
		}
		limit := page * size
		activities, err := h.readModel.ListAll(ctx.Request.Context(), limit)
		if err != nil {
			log.Warnf("Read-model fetch (all) failed, falling back to DB: %v", err)
		} else {

			start := (page - 1) * size
			end := start + size
			if start > len(activities) {
				start = len(activities)
			}
			if end > len(activities) {
				end = len(activities)
			}
			resp := &activityV1.ActivityListResponse{Activities: make([]activityV1.ActivityResponse, 0, end-start)}
			for _, a := range activities[start:end] {
				ar := a.ToResponse()
				resp.Activities = append(resp.Activities, *ar)
			}
			resp.Total = int64(len(activities))
			resp.Page = page
			resp.PageSize = size
			if resp.PageSize > 0 {
				resp.TotalPages = int((resp.Total + int64(resp.PageSize) - 1) / int64(resp.PageSize))
			} else {
				resp.TotalPages = 1
			}

			if end-start > 0 {
				last := activities[end-1]
				resp.NextPageToken = encodeCursorToken(last.CreatedAt, last.ID)
			}

			log.Infof("Listed %d activities out of approx %d total. Used Firestore read model (all).", len(resp.Activities), resp.Total)
			ctx.JSON(http.StatusOK, resp)
			return
		}
	}

	response, err := h.svc.ListActivities(ctx.Request.Context(), &req)
	if err != nil {
		log.Errorf("Failed to list activities: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list activities", "details": err.Error()})
		return
	}

	log.Infof("Listed %d activities out of %d total. Used PostgreSQL write model.", len(response.Activities), response.Total)

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
	log := h.log.WithContext(ctx.Request.Context())
	defer utils.TimeOperation(log, "ActivityHandler.GetActivityStats")()
	log.Info("Received Get Activity Stats request")

	if ctx.Request != nil {
		log = h.log.WithContext(ctx.Request.Context())
	}
	log.Info("Received Get Activity Stats request")

	response, err := h.svc.GetActivityStats(ctx.Request.Context())
	if err != nil {
		log.Errorf("Failed to get activity stats: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get activity stats", "details": err.Error()})
		return
	}

	log.Info("Successfully retrieved activity stats")

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
	log := h.log.WithContext(ctx.Request.Context())
	defer utils.TimeOperation(log, "ActivityHandler.DeleteActivity")()
	log.Info("Received Delete Activity request")

	if ctx.Request != nil {
		log = h.log.WithContext(ctx.Request.Context())
	}
	log.Info("Received Delete Activity request")

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Errorf("Invalid activity ID: %s", idStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	err = h.svc.DeleteActivity(ctx.Request.Context(), uint(id))
	if err != nil {
		log.Errorf("Failed to delete activity: %v", err)
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity", "details": err.Error()})
		return
	}

	log.Infof("Successfully deleted activity with ID: %d", id)
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
	log := h.log.WithContext(ctx.Request.Context())
	defer utils.TimeOperation(log, "ActivityHandler.ResetAllData")()
	log.Info("Received Reset All Data request")

	if ctx.Request != nil {
		log = h.log.WithContext(ctx.Request.Context())
	}
	log.Info("Received Reset All Data request")

	err := h.svc.ResetAllData(ctx.Request.Context())
	if err != nil {
		log.Errorf("Failed to reset activity data: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset activity data", "details": err.Error()})
		return
	}

	log.Info("Successfully reset all activity data")
	ctx.JSON(http.StatusOK, gin.H{"message": "All activity data reset successfully"})
}
