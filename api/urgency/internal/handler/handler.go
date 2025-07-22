package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
	"github.com/pd120424d/mountain-service/api/urgency/internal/repositories"
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
	log  utils.Logger
	repo repositories.UrgencyRepository
}

func NewUrgencyHandler(log utils.Logger, repo repositories.UrgencyRepository) UrgencyHandler {
	return &urgencyHandler{log: log.WithName("urgencyHandler"), repo: repo}
}

// CreateUrgency Креирање нове ургентне ситуације
// @Summary Креирање нове ургентне ситуације
// @Description Креирање нове ургентне ситуације са свим потребним подацима
// @Tags urgency
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param urgency body model.UrgencyCreateRequest true "Urgency data"
// @Success 201 {object} model.UrgencyResponse
// @Router /urgencies [post]
func (h *urgencyHandler) CreateUrgency(ctx *gin.Context) {
	h.log.Info("Received Create Urgency request")

	var req model.UrgencyCreateRequest
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
		Name:         req.Name,
		Email:        req.Email,
		ContactPhone: req.ContactPhone,
		Description:  req.Description,
		Level:        req.Level,
		Status:       "Open",
	}

	if err := h.repo.Create(&urgency); err != nil {
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
// @Security BearerAuth
// @Produce  json
// @Success 200 {array} []model.UrgencyResponse
// @Router /urgencies [get]
func (h *urgencyHandler) ListUrgencies(ctx *gin.Context) {
	h.log.Info("Received List Urgencies request")

	urgencies, err := h.repo.GetAll()
	if err != nil {
		h.log.Errorf("failed to retrieve urgencies: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]model.UrgencyResponse, 0)
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
// @Security BearerAuth
// @Produce  json
// @Param id path int true "Urgency ID"
// @Success 200 {object} model.UrgencyResponse
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

	var urgency model.Urgency
	if err := h.repo.GetByID(uint(urgencyID), &urgency); err != nil {
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
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Urgency ID"
// @Param urgency body model.UrgencyUpdateRequest true "Updated urgency data"
// @Success 200 {object} model.UrgencyResponse
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

	var req model.UrgencyUpdateRequest
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

	var urgency model.Urgency
	if err := h.repo.GetByID(uint(urgencyID), &urgency); err != nil {
		h.log.Errorf("failed to get urgency with ID %d: %v", urgencyID, err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "urgency not found"})
		return
	}

	// Update fields if provided
	if req.Name != "" {
		urgency.Name = req.Name
	}
	if req.Email != "" {
		urgency.Email = req.Email
	}
	if req.ContactPhone != "" {
		urgency.ContactPhone = req.ContactPhone
	}
	if req.Description != "" {
		urgency.Description = req.Description
	}
	if req.Level != "" {
		urgency.Level = req.Level
	}
	if req.Status != "" {
		urgency.Status = req.Status
	}

	if err := h.repo.Update(&urgency); err != nil {
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
// @Security BearerAuth
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

	// Check if urgency exists
	var urgency model.Urgency
	if err := h.repo.GetByID(uint(urgencyID), &urgency); err != nil {
		h.log.Errorf("failed to get urgency with ID %d: %v", urgencyID, err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "urgency not found"})
		return
	}

	if err := h.repo.Delete(uint(urgencyID)); err != nil {
		h.log.Errorf("failed to delete urgency with ID %d: %v", urgencyID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Infof("Successfully deleted urgency with ID %d", urgencyID)
	ctx.Status(http.StatusNoContent)
}

// ResetAllData Ресетовање свих података
// @Summary Ресетовање свих података
// @Description Брисање свих ургентних ситуација (само за администраторе)
// @Tags urgency
// @Security BearerAuth
// @Success 204
// @Router /admin/urgencies/reset [delete]
func (h *urgencyHandler) ResetAllData(ctx *gin.Context) {
	h.log.Info("Received Reset All Data request")

	if err := h.repo.ResetAllData(); err != nil {
		h.log.Errorf("failed to reset all data: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Successfully reset all urgency data")
	ctx.Status(http.StatusNoContent)
}
