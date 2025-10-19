package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/pd120424d/mountain-service/api/activity/internal/clients"
	"github.com/pd120424d/mountain-service/api/activity/internal/service"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/config"
	sharedModels "github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

const defaultPageSize = 50
const maximumPageSize = 100
const maxBatchSize = 100

type ActivityHandler interface {
	CreateActivity(ctx *gin.Context)
	GetActivity(ctx *gin.Context)
	ListActivities(ctx *gin.Context)
	GetActivityCounts(ctx *gin.Context)
	DeleteActivity(ctx *gin.Context)
	ResetAllData(ctx *gin.Context)

	// Admin-only endpoints
	AddActivitiesBatch(ctx *gin.Context)
	GetActivitySourceFlag(ctx *gin.Context)
	SetActivitySourceFlag(ctx *gin.Context)

	SetFeatureFlagService(svc service.FeatureFlagService)
}

type ActivityHandlerConfig struct {
	DefaultSource  string
	AdminCanToggle bool
}

type activityHandler struct {
	log       utils.Logger
	svc       service.ActivityService
	readModel service.FirestoreService

	urgencyClient clients.UrgencyClient
	config        ActivityHandlerConfig
	flags         service.FeatureFlagService
}

func NewActivityHandler(log utils.Logger, svc service.ActivityService, readModel service.FirestoreService, urgencyClient clients.UrgencyClient, defaultSource string, adminCanToggle bool) ActivityHandler {
	return &activityHandler{
		log:           log.WithName("activityHandler"),
		svc:           svc,
		readModel:     readModel,
		urgencyClient: urgencyClient,
		config: ActivityHandlerConfig{
			DefaultSource:  defaultSource,
			AdminCanToggle: adminCanToggle,
		},
	}
}

// ActivitySourceFlagRequest represents the admin toggle payload
// swagger:model ActivitySourceFlagRequest
type ActivitySourceFlagRequest struct {
	UsePostgres bool `json:"usePostgres" example:"true"`
}

// ActivitySourceFlagResponse represents the flag state returned to clients
// swagger:model ActivitySourceFlagResponse
type ActivitySourceFlagResponse struct {
	UsePostgres bool `json:"usePostgres"`
}

// SetFeatureFlagService wires a shared feature flag provider (optional).
func (h *activityHandler) SetFeatureFlagService(svc service.FeatureFlagService) {
	h.flags = svc
}

// GetActivitySourceFlag Админ: враћа тренутну вредност feature флага (да ли се користи Postgres за листање)
// @Summary Админ: одакле се читају активности
// @Description Враћа глобалну вредност флага: ако је true, користи се Postgres за читање активности; иначе Firestore
// @Tags admin
// @Produce json
// @Success 200 {object} ActivitySourceFlagResponse
// @Failure 503 {object} map[string]string
// @Router /admin/feature-flags/activity-source [get]
func (h *activityHandler) GetActivitySourceFlag(ctx *gin.Context) {
	if h.flags == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Feature flag service not available"})
		return
	}
	usePg, _ := h.flags.GetUsePostgresForActivities(ctx.Request.Context())
	ctx.JSON(http.StatusOK, ActivitySourceFlagResponse{UsePostgres: usePg})
}

// SetActivitySourceFlag Админ: поставља глобални feature флаг
// @Summary Админ: постави извор за листање активности
// @Description Када се укључи (true), листање активности користи Postgres. Подразумевано је false (Firestore).
// @Tags admin
// @Accept json
// @Produce json
// @Param flag body ActivitySourceFlagRequest true "Toggle Postgres as activity source"
// @Success 200 {object} ActivitySourceFlagResponse
// @Failure 400 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /admin/feature-flags/activity-source [put]
func (h *activityHandler) SetActivitySourceFlag(ctx *gin.Context) {
	if h.flags == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Feature flag service not available"})
		return
	}
	var req ActivitySourceFlagRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}
	actor := "admin"
	if v, ok := ctx.Get("email"); ok {
		if s, ok2 := v.(string); ok2 && s != "" {
			actor = s
		}
	}
	_ = h.flags.SetUsePostgresForActivities(ctx.Request.Context(), req.UsePostgres, actor)
	ctx.JSON(http.StatusOK, ActivitySourceFlagResponse{UsePostgres: req.UsePostgres})
}

func (h *activityHandler) determineSource(ctx *gin.Context) string {
	if override, exists := ctx.Get("activity_source_override"); exists {
		if source, ok := override.(string); ok {
			return source
		}
	}
	if h.flags != nil {
		if usePg, _ := h.flags.GetUsePostgresForActivities(ctx.Request.Context()); usePg {
			return "postgres"
		}
	}
	return h.config.DefaultSource
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

	if h.urgencyClient != nil {
		cctx, cancel := context.WithTimeout(ctx.Request.Context(), config.DefaultListTimeout)
		defer cancel()
		if _, uerr := h.urgencyClient.GetUrgencyByID(cctx, req.UrgencyID); uerr != nil {
			if strings.Contains(strings.ToLower(uerr.Error()), "not found") {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "urgency has been deleted or does not exist"})
				return
			}
			log.Warnf("Urgency validation failed; continuing without blocking: %v", uerr)
		}
	}

	response, err := h.svc.CreateActivity(ctx.Request.Context(), &req)
	if err != nil {
		log.Errorf("Failed to create activity: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity", "details": err.Error()})
		return
	}

	log.Infof("Successfully created activity with ID: %d", response.ID)

	utils.WriteFreshWindow(ctx, config.DefaultFreshWindow)

	ctx.JSON(http.StatusCreated, response)
}

// AddActivitiesBatch Админ: креирање више активности у једном захтеву (само за тестирање)
// @Summary Админ: серијско додавање активности
// @Description Креира више активности у једном захтеву. Само за администраторе.
// @Tags admin
// @Accept json
// @Produce json
// @Param request body activityV1.BatchAddActivitiesRequest true "Пакет захтева за креирање активности"
// @Success 200 {object} activityV1.BatchAddActivitiesResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/activities/batch [post]
func (h *activityHandler) AddActivitiesBatch(ctx *gin.Context) {
	log := h.log.WithContext(ctx.Request.Context())
	defer utils.TimeOperation(log, "ActivityHandler.AddActivitiesBatch")()

	var req activityV1.BatchAddActivitiesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorf("Failed to bind batch request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if len(req.Items) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "items must not be empty"})
		return
	}
	if len(req.Items) > maxBatchSize {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("items cannot exceed %d", maxBatchSize)})
		return
	}

	results, err := h.svc.CreateActivitiesBatch(ctx.Request.Context(), req.Items)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Batch operation failed", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, activityV1.BatchAddActivitiesResponse{Results: results})
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

	if h.urgencyClient != nil {
		cctx, cancel := context.WithTimeout(ctx.Request.Context(), config.DefaultListTimeout)
		defer cancel()
		if _, uerr := h.urgencyClient.GetUrgencyByID(cctx, response.UrgencyID); uerr != nil {
			if strings.Contains(strings.ToLower(uerr.Error()), "not found") {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "urgency has been deleted or does not exist"})
				return
			}
			log.Warnf("Urgency validation failed; continuing without blocking: %v", uerr)
		}
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
// Matches service-side format: base64(JSON{"createdAt": RFC3339Nano UTC})
type cursorToken struct {
	CreatedAt string `json:"createdAt"`
	ID        uint   `json:"id,omitempty"`
}

func encodeCursorToken(t time.Time, id uint) string {
	if t.IsZero() {
		return ""
	}
	b, _ := json.Marshal(cursorToken{CreatedAt: t.UTC().Format(time.RFC3339Nano), ID: id})
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
	// Parse query first, then choose timeout (Postgres and cursor paths get a longer timeout)
	req := buildActivityListRequest(ctx)
	baseCtx := ctx.Request.Context()

	source := h.determineSource(ctx)

	to := config.DefaultListTimeout
	if req.PageToken != "" {
		to = config.CursorListTimeout
	} else if source == "postgres" {
		to = config.PostgresListTimeout
	}
	cctx, cancel := context.WithTimeout(baseCtx, to)
	defer cancel()
	log := h.log.WithContext(cctx)
	defer utils.TimeOperation(log, "ActivityHandler.ListActivities",
		zap.String("source", source),
		zap.Int("page_size", req.PageSize),
		zap.Bool("has_page_token", req.PageToken != ""),
	)()

	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		log.Infof("ActivityHandler.ListActivities completed: source=%s pageSize=%d hasPageToken=%v duration=%dms",
			source, req.PageSize, req.PageToken != "", duration.Milliseconds())
	}()

	log.Infof("Received List Activities request: source=%s urgencyId=%v pageToken=%v page=%d pageSize=%d",
		source, req.UrgencyID, req.PageToken != "", req.Page, req.PageSize)

	if req.UrgencyID != nil && h.urgencyClient != nil {
		if _, err := h.urgencyClient.GetUrgencyByID(cctx, *req.UrgencyID); err != nil {
			// Treat not found as deleted or missing
			if strings.Contains(strings.ToLower(err.Error()), "not found") {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "urgency has been deleted or does not exist"})
				return
			}
			// On transient client errors, proceed to avoid breaking UX; logs will capture the failure
			log.Warnf("Urgency validation failed; continuing without blocking: %v", err)
		}
	}

	// If source is explicitly set to postgres, use PostgreSQL directly
	if source == "postgres" {
		h.listFromPostgres(ctx, cctx, log, &req)
		return
	}

	// Cursor-based pagination (preferred when pageToken provided)
	if h.readModel != nil && req.PageToken != "" {
		size := req.PageSize
		if size <= 0 {
			size = defaultPageSize
		}
		if size > maximumPageSize {
			size = maximumPageSize
		}

		var (
			activities []sharedModels.Activity
			nextToken  string
			err        error
		)
		if req.UrgencyID != nil {
			activities, nextToken, err = h.readModel.ListByUrgencyCursor(cctx, *req.UrgencyID, size, req.PageToken)
		} else {
			activities, nextToken, err = h.readModel.ListAllCursor(cctx, size, req.PageToken)
		}
		if err != nil {
			log.Warnf("Cursor read-model fetch failed; returning empty page to stop repetition: %v", err)
			resp := &activityV1.ActivityListResponse{
				Activities:    []activityV1.ActivityResponse{},
				Total:         0,
				PageSize:      size,
				NextPageToken: "",
			}
			ctx.JSON(http.StatusOK, resp)
			return
		} else {
			marshalStart := time.Now()
			resp := &activityV1.ActivityListResponse{Activities: make([]activityV1.ActivityResponse, 0, len(activities))}
			for _, a := range activities {
				ar := a.ToResponse()
				resp.Activities = append(resp.Activities, *ar)
			}
			resp.Total = int64(len(activities))
			resp.PageSize = size
			resp.NextPageToken = nextToken
			marshalMs := time.Since(marshalStart).Milliseconds()
			log.Infof("Firestore cursor marshal: marshalMs=%d count=%d nextToken=%v", marshalMs, len(resp.Activities), nextToken != "")
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
			size = defaultPageSize
		}

		limit := page * size
		activities, err := h.readModel.ListByUrgency(cctx, *req.UrgencyID, limit)
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
			size = defaultPageSize
		}
		limit := page * size
		activities, err := h.readModel.ListAll(cctx, limit)
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

	response, err := h.svc.ListActivities(cctx, &req)
	if err != nil {
		log.Errorf("Failed to list activities: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list activities", "details": err.Error()})
		return
	}

	log.Infof("Listed %d activities out of %d total. Used PostgreSQL write model.", len(response.Activities), response.Total)

	ctx.JSON(http.StatusOK, response)
}

// GetActivityCounts Број активности по ургенцији (Firestore)
// @Summary Број активности по ургенцији
// @Description Враћа број активности за наведене ургенције, помоћу Firestore агрегатних упита.
// @Tags activities
// @Produce json
// @Param urgencyId query []int true "Идентификатори ургенција" collectionFormat(multi)
// @Success 200 {object} activityV1.ActivityCountsResponse
// @Failure 400 {object} activityV1.ErrorResponse
// @Failure 503 {object} activityV1.ErrorResponse
// @Router /activities/counts [get]
func (h *activityHandler) GetActivityCounts(ctx *gin.Context) {
	log := h.log.WithContext(ctx.Request.Context())
	defer utils.TimeOperation(log, "ActivityHandler.GetActivityCounts")()

	if h.readModel == nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Read model (Firestore) is not available"})
		return
	}

	raw := ctx.QueryArray("urgencyId")
	if len(raw) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "at least one urgencyId is required"})
		return
	}
	if len(raw) > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "urgencyId cannot exceed 100 per request"})
		return
	}
	ids := make([]uint, 0, len(raw))
	for i, s := range raw {
		v, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil || v <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("urgencyId[%d] must be a positive integer", i)})
			return
		}
		ids = append(ids, uint(v))
	}

	cctx, cancel := context.WithTimeout(ctx.Request.Context(), config.CountTimeout)
	defer cancel()

	counts, err := h.readModel.CountByUrgencyIDs(cctx, ids)
	if err != nil {
		log.Errorf("Failed to get activity counts: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compute counts", "details": err.Error()})
		return
	}

	out := make(map[string]int64, len(counts))
	for id, c := range counts {
		out[strconv.Itoa(int(id))] = c
	}
	ctx.JSON(http.StatusOK, activityV1.ActivityCountsResponse{Counts: out})
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

	utils.WriteFreshWindow(ctx, config.DefaultFreshWindow)

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

	utils.WriteFreshWindow(ctx, config.DefaultFreshWindow)

	ctx.JSON(http.StatusOK, gin.H{"message": "All activity data reset successfully"})
}

func (h *activityHandler) listFromPostgres(ctx *gin.Context, cctx context.Context, log utils.Logger, req *activityV1.ActivityListRequest) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		log.Infof("PostgreSQL query completed: duration=%dms urgencyId=%v page=%d pageSize=%d",
			duration.Milliseconds(), req.UrgencyID, req.Page, req.PageSize)
	}()

	response, err := h.svc.ListActivities(cctx, req)
	if err != nil {
		log.Errorf("Failed to list activities from PostgreSQL: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list activities", "details": err.Error()})
		return
	}

	log.Infof("Listed %d activities out of %d total. Used PostgreSQL write model.", len(response.Activities), response.Total)
	ctx.JSON(http.StatusOK, response)
}
