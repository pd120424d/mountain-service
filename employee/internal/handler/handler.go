package handler

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mountain-service/employee/internal/model"
	"mountain-service/employee/internal/repositories"
	"mountain-service/shared/utils"
)

type EmployeeHandler interface {
	RegisterEmployee(ctx *gin.Context)
	ListEmployees(c *gin.Context)
	DeleteEmployee(c *gin.Context)
}

type employeeHandler struct {
	log  utils.Logger
	repo repositories.EmployeeRepository
}

func NewEmployeeHandler(log utils.Logger, repo repositories.EmployeeRepository) EmployeeHandler {
	return &employeeHandler{log: log, repo: repo}
}

// RegisterEmployee Креирање новог запосленог
// @Summary Креирање новог запосленог
// @Description Креирање новог запосленог у систему
// @Tags запослени
// @Accept  json
// @Produce  json
// @Param employee body model.EmployeeCreateRequest true "Подаци о новом запосленом"
// @Success 201 {object} model.EmployeeResponse
// @Failure 400 {object} gin.H
// @Router /employees [post]
func (h *employeeHandler) RegisterEmployee(ctx *gin.Context) {
	var req model.EmployeeCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee := model.Employee{
		Username:       req.Username,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Password:       req.Password,
		Gender:         req.Gender,
		Phone:          req.Phone,
		Email:          req.Email,
		ProfilePicture: req.ProfilePicture,
		ProfileType:    req.ProfileType,
	}

	// Check for unique username
	usernameFilter := map[string]interface{}{
		"username": employee.Username,
	}
	existingEmployees, err := h.repo.ListEmployees(usernameFilter)
	if err != nil {
		h.log.Error("Failed to check for existing username", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for existing username"})
		return
	}
	if len(existingEmployees) > 0 {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Check for unique email
	emailFilter := map[string]interface{}{
		"email": employee.Email,
	}
	existingEmployees, err = h.repo.ListEmployees(emailFilter)
	if err != nil {
		h.log.Error("Failed to check for existing email", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for existing email"})
		return
	}
	if len(existingEmployees) > 0 {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Validate the password
	if err := utils.ValidatePassword(employee.Password); err != nil {
		h.log.Errorf("password validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(&employee); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := model.EmployeeResponse{
		ID:             employee.ID,
		Username:       employee.Username,
		FirstName:      employee.FirstName,
		LastName:       employee.LastName,
		Gender:         employee.Gender,
		Phone:          employee.Phone,
		Email:          employee.Email,
		ProfilePicture: employee.ProfilePicture,
		ProfileType:    employee.ProfileType,
	}

	ctx.JSON(http.StatusOK, response)
}

// ListEmployees Преузимање листе запослених
// @Summary Преузимање листе запослених
// @Description Преузимање свих запослених
// @Tags запослени
// @Produce  json
// @Success 200 {array} model.EmployeeResponse
// @Router /employees [get]
func (h *employeeHandler) ListEmployees(c *gin.Context) {
	employees, err := h.repo.GetAll()
	if err != nil {
		h.log.Errorf("failed to retrieve employees: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Infof("Successfully retrieved %d employees", len(employees))
	var response []model.EmployeeResponse
	for _, emp := range employees {
		response = append(response, model.EmployeeResponse{
			ID:             emp.ID,
			Username:       emp.Username,
			FirstName:      emp.FirstName,
			LastName:       emp.LastName,
			Gender:         emp.Gender,
			Phone:          emp.Phone,
			Email:          emp.Email,
			ProfilePicture: emp.ProfilePicture,
		})
	}
	c.JSON(http.StatusOK, employees)
}

func (h *employeeHandler) UpdateEmployee(db *gorm.DB, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var employee model.Employee
		if err := db.First(&employee, id).Error; err != nil {
			logger.Error("Employee not found", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}

		if err := c.ShouldBindJSON(&employee); err != nil {
			logger.Error("Failed to bind employee data", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		if err := db.Save(&employee).Error; err != nil {
			logger.Error("Failed to update employee", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
			return
		}

		c.JSON(http.StatusOK, employee)
	}
}

// DeleteEmployee Брисање запосленог
// @Summary Брисање запосленог
// @Description Брисање запосленог по ID-ју
// @Tags запослени
// @Param id path int true "ID запосленог"
// @Success 204
// @Failure 404 {object} gin.H
// @Router /employees/{id} [delete]
func (h *employeeHandler) DeleteEmployee(c *gin.Context) {
	idParam := c.Param("id")
	employeeID, err := strconv.Atoi(idParam)
	if err != nil {
		h.log.Errorf("failed to convert employee ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	if err := h.repo.Delete(uint(employeeID)); err != nil {
		h.log.Errorf("failed to delete employee: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
		return
	}

	h.log.Infof("Employee with ID %d was soft deleted", employeeID)
	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}
