package internal

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	commonv1 "github.com/pd120424d/mountain-service/api/contracts/common/v1"
	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
)

type UrgencyHandler interface {
	CreateUrgency(ctx *gin.Context)
	ListUrgencies(ctx *gin.Context)
	GetUrgency(ctx *gin.Context)
	UpdateUrgency(ctx *gin.Context)
	DeleteUrgency(ctx *gin.Context)
	ResetAllData(ctx *gin.Context)

	AssignUrgency(ctx *gin.Context)
	UnassignUrgency(ctx *gin.Context)
	CloseUrgency(ctx *gin.Context)
}

type urgencyHandler struct {
	log utils.Logger
	svc UrgencyService
}

func NewUrgencyHandler(log utils.Logger, svc UrgencyService) UrgencyHandler {
	return &urgencyHandler{log: log.WithName("urgencyHandler"), svc: svc}
}

// CreateUrgency Креирање нове ургентне ситуације
// @Summary Креирање нове ургентне ситуације
// @Description Креирање нове ургентне ситуације са свим потребним подацима
// @Tags urgency
// @Security OAuth2Password
// @Accept  json
// @Produce  json
// @Param urgency body urgencyV1.UrgencyCreateRequest true "Urgency data"
// @Success 201 {object} urgencyV1.UrgencyResponse
// @Router /urgencies [post]
func (h *urgencyHandler) CreateUrgency(ctx *gin.Context) {
	log := h.log.WithContext(requestContext(ctx))
	log.Info("Received Create Urgency request")

	var req urgencyV1.UrgencyCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.log.Errorf("failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		h.log.Errorf("validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	urgency := model.Urgency{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		ContactPhone: req.ContactPhone,
		Location:     req.Location,
		Description:  req.Description,
		Level:        req.Level,
		Status:       urgencyV1.Open,
	}

	if err := h.svc.CreateUrgency(requestContext(ctx), &urgency); err != nil {
		log.Errorf("failed to create urgency: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "URGENCY_ERRORS.CREATE_FAILED", "details": err.Error()})
		return
	}

	response := urgency.ToResponse()
	ctx.JSON(http.StatusCreated, response)

	log.Infof("Successfully created urgency with ID %d", urgency.ID)
}

// ListUrgencies Извлачење листе ургентних ситуација
// @Summary Извлачење листе ургентних ситуација
// @Description Извлачење свих ургентних ситуација
// @Tags urgency
// @Security OAuth2Password
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Success 200 {object} urgencyV1.UrgencyListResponse
// @Router /urgencies [get]
func (h *urgencyHandler) ListUrgencies(ctx *gin.Context) {
	reqLog := h.log.WithContext(requestContext(ctx))
	reqLog.Info("Received List Urgencies request")

	// Parse pagination params
	page := 1
	pageSize := 20
	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if sizeStr := ctx.Query("pageSize"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	var assignedEmployeeID *uint
	if ctx.Query("myUrgencies") == "true" {
		if v, exists := ctx.Get("employeeID"); exists {
			if id, ok := v.(uint); ok {
				assignedEmployeeID = &id
			}
		}
	}

	urgencies, total, err := h.svc.ListUrgencies(requestContext(ctx), page, pageSize, assignedEmployeeID)
	if err != nil {
		reqLog.Errorf("failed to retrieve urgencies: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "URGENCY_ERRORS.LIST_FAILED", "details": err.Error()})
		return
	}

	items := make([]urgencyV1.UrgencyResponse, 0, len(urgencies))
	for _, u := range urgencies {
		items = append(items, u.ToResponse())
	}
	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	resp := urgencyV1.UrgencyListResponse{Urgencies: items, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}
	reqLog.Infof("Successfully retrieved %d urgencies out of %d total", len(items), total)
	ctx.JSON(http.StatusOK, resp)
}

// GetUrgency Извлачење ургентне ситуације по ID
// @Summary Извлачење ургентне ситуације по ID
// @Description Извлачење ургентне ситуације по њеном ID
// @Tags urgency
// @Security OAuth2Password
// @Produce  json
// @Param id path int true "Urgency ID"
// @Success 200 {object} urgencyV1.UrgencyResponse
// @Router /urgencies/{id} [get]
func (h *urgencyHandler) GetUrgency(ctx *gin.Context) {
	reqLog := h.log.WithContext(requestContext(ctx))
	reqLog.Info("Received Get Urgency request")

	idParam := ctx.Param("id")
	urgencyID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		reqLog.Errorf("invalid urgency ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid urgency ID"})
		return
	}

	urgency, err := h.svc.GetUrgencyByID(requestContext(ctx), uint(urgencyID))
	if err != nil {
		reqLog.Errorf("failed to get urgency with ID %d: %v", urgencyID, err)
		if aerr, ok := err.(*commonv1.AppError); ok {
			if aerr.Code == "URGENCY_ERRORS.NOT_FOUND" {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "urgency not found"})
				return
			}
		}
		// Also handle gorm.ErrRecordNotFound from service mocks in tests
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "urgency not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "URGENCY_ERRORS.DB_ERROR"})
		return
	}

	response := urgency.ToResponse()
	reqLog.Infof("Successfully retrieved urgency with ID %d", urgencyID)
	ctx.JSON(http.StatusOK, response)
}

// UpdateUrgency Ажурирање ургентне ситуације
// @Summary Ажурирање ургентне ситуације
// @Description Ажурирање постојеће ургентне ситуације
// @Tags urgency
// @Security OAuth2Password
// @Accept  json
// @Produce  json
// @Param id path int true "Urgency ID"
// @Param urgency body urgencyV1.UrgencyUpdateRequest true "Updated urgency data"
// @Success 200 {object} urgencyV1.UrgencyResponse
// @Router /urgencies/{id} [put]
func (h *urgencyHandler) UpdateUrgency(ctx *gin.Context) {
	reqLog := h.log.WithContext(requestContext(ctx))
	reqLog.Info("Received Update Urgency request")

	idParam := ctx.Param("id")
	urgencyID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		reqLog.Errorf("invalid urgency ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid urgency ID"})
		return
	}

	var req urgencyV1.UrgencyUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLog.Errorf("failed to bind JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		reqLog.Errorf("validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	urgency, err := h.svc.GetUrgencyByID(requestContext(ctx), uint(urgencyID))
	if err != nil {
		reqLog.Errorf("failed to get urgency with ID %d: %v", urgencyID, err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "urgency not found"})
		return
	}

	urgency.UpdateWithRequest(&req)

	if err := h.svc.UpdateUrgency(requestContext(ctx), urgency); err != nil {
		reqLog.Errorf("failed to update urgency: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "URGENCY_ERRORS.UPDATE_FAILED", "details": err.Error()})
		return
	}

	response := urgency.ToResponse()
	reqLog.Infof("Successfully updated urgency with ID %d", urgencyID)
	ctx.JSON(http.StatusOK, response)
}

// DeleteUrgency Брисање ургентне ситуације
// @Summary Брисање ургентне ситуације
// @Description Брисање ургентне ситуације по ID
// @Tags urgency
// @Security OAuth2Password
// @Param id path int true "Urgency ID"
// @Success 204
// @Router /urgencies/{id} [delete]
func (h *urgencyHandler) DeleteUrgency(ctx *gin.Context) {
	reqLog := h.log.WithContext(requestContext(ctx))
	reqLog.Info("Received Delete Urgency request")

	idParam := ctx.Param("id")
	urgencyID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		reqLog.Errorf("invalid urgency ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid urgency ID"})
		return
	}

	if err := h.svc.DeleteUrgency(requestContext(ctx), uint(urgencyID)); err != nil {
		reqLog.Errorf("failed to delete urgency with ID %d: %v", urgencyID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "URGENCY_ERRORS.DELETE_FAILED", "details": err.Error()})
		return
	}

	reqLog.Infof("Successfully deleted urgency with ID %d", urgencyID)
	ctx.JSON(http.StatusNoContent, nil)
}

// ResetAllData Ресетовање свих података
// @Summary Ресетовање свих података
// @Description Брисање свих ургентних ситуација (само за администраторе)
// @Tags urgency
// @Security OAuth2Password
// @Success 204
// @Router /admin/urgencies/reset [delete]
func (h *urgencyHandler) ResetAllData(ctx *gin.Context) {
	log := h.log.WithContext(requestContext(ctx))
	log.Info("Received Reset All Data request")

	if err := h.svc.ResetAllData(requestContext(ctx)); err != nil {
		log.Errorf("failed to reset all data: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "URGENCY_ERRORS.RESET_FAILED", "details": err.Error()})
		return
	}

	log.Info("Successfully reset all urgency data")
	ctx.JSON(http.StatusNoContent, nil)
}

// AssignUrgency Додела ургентне ситуације запосленом
// @Summary Додела ургентне ситуације запосленом
// @Description Додела ургентне ситуације запосленом
// @Tags urgency
// @Security OAuth2Password
// @Param id path int true "Urgency ID"
// @Param payload body urgencyV1.AssignmentCreateRequest true "Assignment payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /urgencies/{id}/assign [post]
func (h *urgencyHandler) AssignUrgency(ctx *gin.Context) {
	log := h.log.WithContext(requestContext(ctx))
	log.Info("Received Assign Urgency request")

	idParam := ctx.Param("id")
	urgencyID64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || urgencyID64 == 0 {
		log.Errorf("invalid urgency ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid urgency ID"})
		return
	}
	var req urgencyV1.AssignmentCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorf("failed to bind json: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if appErr := h.svc.AssignUrgency(requestContext(ctx), uint(urgencyID64), req.EmployeeID); appErr != nil {
		log.Errorf("assign failed: %v", appErr)
		if aerr, ok := appErr.(*commonv1.AppError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": aerr.Code, "details": aerr.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": appErr.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "assigned"})

	log.Info("Successfully assigned urgency")
}

// UnassignUrgency Уклањање доделе ургентне ситуације (админ)
// @Summary Уклањање доделе ургентне ситуације
// @Description Уклањање доделе ургентне ситуације (админ)
// @Tags urgency
// @Security OAuth2Password
// @Param id path int true "Urgency ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Router /urgencies/{id}/assign [delete]
func (h *urgencyHandler) UnassignUrgency(ctx *gin.Context) {
	log := h.log.WithContext(requestContext(ctx))
	log.Info("Received Unassign Urgency request")

	idParam := ctx.Param("id")
	urgencyID64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || urgencyID64 == 0 {
		log.Errorf("invalid urgency ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid urgency ID"})
		return
	}
	actorIDVal, _ := ctx.Get("employeeID")
	roleVal, _ := ctx.Get("role")
	actorID, _ := actorIDVal.(uint)
	isAdmin := roleVal == "Administrator"
	if err := h.svc.UnassignUrgency(requestContext(ctx), uint(urgencyID64), actorID, isAdmin); err != nil {
		log.Errorf("unassign failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)

	log.Info("Successfully unassigned urgency")
}

// CloseUrgency Затварање ургентне ситуације (админ)
// @Summary Затварање ургентне ситуације (админ)
// @Description Затварање ургентне ситуације (админ)
// @Tags urgency
// @Security OAuth2Password
// @Param id path int true "Urgency ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Router /urgencies/{id}/close [put]
func (h *urgencyHandler) CloseUrgency(ctx *gin.Context) {
	log := h.log.WithContext(requestContext(ctx))
	log.Info("Received Close Urgency request")

	idParam := ctx.Param("id")
	urgencyID64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || urgencyID64 == 0 {
		log.Errorf("invalid urgency ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid urgency ID"})
		return
	}
	actorIDVal, _ := ctx.Get("employeeID")
	roleVal, _ := ctx.Get("role")
	actorID, _ := actorIDVal.(uint)
	isAdmin := roleVal == "Administrator"
	if err := h.svc.CloseUrgency(requestContext(ctx), uint(urgencyID64), actorID, isAdmin); err != nil {
		log.Errorf("close failed: %v", err)
		if aerr, ok := err.(*commonv1.AppError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": aerr.Code, "details": aerr.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)

	log.Info("Successfully closed urgency")
}

func requestContext(ctx *gin.Context) context.Context {
	if ctx != nil && ctx.Request != nil {
		return ctx.Request.Context()
	}
	return context.Background()
}
