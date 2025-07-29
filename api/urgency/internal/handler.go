package internal

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

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
	h.log.Info("Received Create Urgency request")

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
		Level:        urgencyV1.High,
		Status:       urgencyV1.Open,
	}

	if err := h.svc.CreateUrgency(&urgency); err != nil {
		h.log.Errorf("failed to create urgency: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := urgency.ToResponse()
	h.log.Infof("Successfully created urgency with ID %d", urgency.ID)
	ctx.JSON(http.StatusCreated, response)
}

// ListUrgencies Извлачење листе ургентних ситуација
// @Summary Извлачење листе ургентних ситуација
// @Description Извлачење свих ургентних ситуација
// @Tags urgency
// @Security OAuth2Password
// @Produce  json
// @Success 200 {array} []urgencyV1.UrgencyResponse
// @Router /urgencies [get]
func (h *urgencyHandler) ListUrgencies(ctx *gin.Context) {
	h.log.Info("Received List Urgencies request")

	urgencies, err := h.svc.GetAllUrgencies()
	if err != nil {
		h.log.Errorf("failed to retrieve urgencies: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]urgencyV1.UrgencyResponse, 0)
	for _, urgency := range urgencies {
		response = append(response, urgency.ToResponse())
	}

	h.log.Infof("Successfully retrieved %d urgencies", len(response))
	ctx.JSON(http.StatusOK, response)
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
	h.log.Info("Received Get Urgency request")

	idParam := ctx.Param("id")
	urgencyID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.log.Errorf("invalid urgency ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid urgency ID"})
		return
	}

	urgency, err := h.svc.GetUrgencyByID(uint(urgencyID))
	if err != nil {
		h.log.Errorf("failed to get urgency with ID %d: %v", urgencyID, err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "urgency not found"})
		return
	}

	response := urgency.ToResponse()
	h.log.Infof("Successfully retrieved urgency with ID %d", urgencyID)
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
	h.log.Info("Received Update Urgency request")

	idParam := ctx.Param("id")
	urgencyID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.log.Errorf("invalid urgency ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid urgency ID"})
		return
	}

	var req urgencyV1.UrgencyUpdateRequest
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

	urgency, err := h.svc.GetUrgencyByID(uint(urgencyID))
	if err != nil {
		h.log.Errorf("failed to get urgency with ID %d: %v", urgencyID, err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "urgency not found"})
		return
	}

	// Update fields if provided
	if req.FirstName != "" {
		urgency.FirstName = req.FirstName
	}
	if req.LastName != "" {
		urgency.LastName = req.LastName
	}
	if req.Email != "" {
		urgency.Email = req.Email
	}
	if req.ContactPhone != "" {
		urgency.ContactPhone = req.ContactPhone
	}
	if req.Location != "" {
		urgency.Location = req.Location
	}
	if req.Description != "" {
		urgency.Description = req.Description
	}
	if req.Level != "" {
		urgency.Level = urgencyV1.UrgencyLevel(req.Level)
	}
	if req.Status != "" {
		urgency.Status = urgencyV1.UrgencyStatus(req.Status)
	}

	if err := h.svc.UpdateUrgency(urgency); err != nil {
		h.log.Errorf("failed to update urgency: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := urgency.ToResponse()
	h.log.Infof("Successfully updated urgency with ID %d", urgencyID)
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
	h.log.Info("Received Delete Urgency request")

	idParam := ctx.Param("id")
	urgencyID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.log.Errorf("invalid urgency ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid urgency ID"})
		return
	}

	if err := h.svc.DeleteUrgency(uint(urgencyID)); err != nil {
		h.log.Errorf("failed to delete urgency with ID %d: %v", urgencyID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Infof("Successfully deleted urgency with ID %d", urgencyID)
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
	h.log.Info("Received Reset All Data request")

	if err := h.svc.ResetAllData(); err != nil {
		h.log.Errorf("failed to reset all data: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Successfully reset all urgency data")
	ctx.JSON(http.StatusNoContent, nil)
}
