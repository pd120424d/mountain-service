package handler

import (
	"api/employee/internal/model"
	"api/employee/internal/repositories"
	"api/shared/utils"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler interface {
	RegisterEmployee(ctx *gin.Context)
	ListEmployees(ctx *gin.Context)
	DeleteEmployee(ctx *gin.Context)
	AssignShift(ctx *gin.Context)
	GetShifts(ctx *gin.Context)
	GetShiftsAvailability(ctx *gin.Context)
}

type employeeHandler struct {
	log        utils.Logger
	emplRepo   repositories.EmployeeRepository
	shiftsRepo repositories.ShiftRepository
}

func NewEmployeeHandler(log utils.Logger, emplRepo repositories.EmployeeRepository, shiftsRepo repositories.ShiftRepository) EmployeeHandler {
	return &employeeHandler{log: log, emplRepo: emplRepo, shiftsRepo: shiftsRepo}
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
	h.log.Info("Received Register Employee request")
	req := &model.EmployeeCreateRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Infof("Creating new employee with data: %s", req.ToString())

	profileType := model.ProfileTypeFromString(req.ProfileType)
	if !profileType.IsValid() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile type"})
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
		ProfileType:    profileType,
	}

	// Check for unique username
	usernameFilter := map[string]interface{}{
		"username": employee.Username,
	}
	existingEmployees, err := h.emplRepo.ListEmployees(usernameFilter)
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
	existingEmployees, err = h.emplRepo.ListEmployees(emailFilter)
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

	if err := h.emplRepo.Create(&employee); err != nil {
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
		ProfileType:    employee.ProfileType.String(),
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
func (h *employeeHandler) ListEmployees(ctx *gin.Context) {
	employees, err := h.emplRepo.GetAll()
	if err != nil {
		h.log.Errorf("failed to retrieve employees: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
			ProfileType:    emp.ProfileType.String(),
		})
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *employeeHandler) UpdateEmployee(db *gorm.DB, logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		var employee model.Employee
		if err := db.First(&employee, id).Error; err != nil {
			logger.Error("Employee not found", zap.Error(err))
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}

		if err := ctx.ShouldBindJSON(&employee); err != nil {
			logger.Error("Failed to bind employee data", zap.Error(err))
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		if err := db.Save(&employee).Error; err != nil {
			logger.Error("Failed to update employee", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
			return
		}

		ctx.JSON(http.StatusOK, employee)
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
func (h *employeeHandler) DeleteEmployee(ctx *gin.Context) {
	idParam := ctx.Param("id")
	employeeID, err := strconv.Atoi(idParam)
	if err != nil {
		h.log.Errorf("failed to convert employee ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	if err := h.emplRepo.Delete(uint(employeeID)); err != nil {
		h.log.Errorf("failed to delete employee: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
		return
	}

	h.log.Infof("Employee with ID %d was soft deleted", employeeID)
	ctx.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}

// AssignShift Додељује смену запосленом
// @Summary Додељује смену запосленом
// @Description Додељује смену запосленом по ID-ју
// @Tags запослени
// @Param id path int true "ID запосленог"
// @Param shift body model.ShiftRequest true "Подаци о смени"
// @Success 201 {object} model.ShiftResponse
// @Failure 400 {object} gin.H
// @Router /employees/{id}/shifts [post]
func (h *employeeHandler) AssignShift(ctx *gin.Context) {
	// Extract employee ID from the URL
	employeeIDParam := ctx.Param("id")
	employeeID, err := strconv.Atoi(employeeIDParam)
	if err != nil || employeeID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	// Parse and validate request body
	var req model.AssignShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse and validate profile type
	profileType := model.ProfileTypeFromString(req.ProfileType)
	if !profileType.IsValid() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile type"})
		return
	}

	// Parse the shift date
	shiftDate, err := time.Parse(model.DateFormat, req.ShiftDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid shiftDate format, expected YYYY-MM-DD"})
		return
	}

	// Call the repository method
	err = h.shiftsRepo.AssignEmployee(shiftDate, req.ShiftType, uint(employeeID), req.ProfileType)
	if errors.Is(err, model.ErrCapacityReached) {
		ctx.JSON(http.StatusConflict, gin.H{"error": "maximum capacity for this role reached in the selected shift"})
	} else if errors.Is(err, model.ErrAlreadyAssigned) {
		ctx.JSON(http.StatusConflict, gin.H{"error": "employee is already assigned to this shift"})
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	// Return success response
	ctx.JSON(http.StatusCreated, gin.H{"message": "shift assigned successfully"})
}

func (h *employeeHandler) GetShifts(ctx *gin.Context) {
	// Extract employee ID from the URL
	employeeIDParam := ctx.Param("id")
	employeeID, err := strconv.Atoi(employeeIDParam)
	if err != nil || employeeID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	// Fetch shifts from the repository
	shifts, err := h.shiftsRepo.GetShiftsByEmployeeID(uint(employeeID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Return the shifts
	ctx.JSON(http.StatusOK, gin.H{"shifts": shifts})
}

func (h *employeeHandler) GetShiftsAvailability(ctx *gin.Context) {
	// Extract and validate date query parameter
	dateParam := ctx.Query("date")
	var date time.Time
	var err error
	if dateParam != "" {
		date, err = time.Parse("2006-01-02", dateParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, expected YYYY-MM-DD"})
			return
		}
	} else {
		date = time.Now() // Default to today
	}

	// Fetch availability from the repository
	availability, err := h.shiftsRepo.GetShiftAvailability(date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Return availability
	ctx.JSON(http.StatusOK, gin.H{"availability": availability})
}

func (h *employeeHandler) RemoveShift(ctx *gin.Context) {
	// Parse and validate request body
	var req model.RemoveShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract employee ID from the URL
	employeeIDParam := ctx.Param("id")
	employeeID, err := strconv.Atoi(employeeIDParam)
	if err != nil || employeeID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	// Parse the shift date
	shiftDate, err := time.Parse(model.DateFormat, req.ShiftDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid shiftDate format, expected YYYY-MM-DD"})
		return
	}

	// Call repository method to remove the shift
	err = h.shiftsRepo.RemoveEmployeeFromShift(shiftDate, req.ShiftType, uint(employeeID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Return success response
	ctx.JSON(http.StatusNoContent, nil)
}
