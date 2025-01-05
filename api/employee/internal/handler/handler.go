package handler

import (
	"api/employee/internal/model"
	"api/employee/internal/repositories"
	"api/shared/utils"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

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
	return &employeeHandler{log: log.WithName("employeeHandler"), emplRepo: emplRepo, shiftsRepo: shiftsRepo}
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

	ctx.JSON(http.StatusCreated, response)
}

// ListEmployees Преузимање листе запослених
// @Summary Преузимање листе запослених
// @Description Преузимање свих запослених
// @Tags запослени
// @Produce  json
// @Success 200 {array} []model.EmployeeResponse
// @Router /employees [get]
func (h *employeeHandler) ListEmployees(ctx *gin.Context) {
	h.log.Info("Received List Employees request")

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

	h.log.Info("Successfully mapped employees to response format")
	h.log.Infof("Returning %d employees", len(response))
	ctx.JSON(http.StatusOK, response)
}

// UpdateEmployee Ажурирање запосленог
// @Summary Ажурирање запосленог
// @Description Ажурирање запосленог по ID-ју
// @Tags запослени
// @Param id path int true "ID запосленог"
// @Param employee body model.EmployeeUpdateRequest true "Подаци за ажурирање запосленог"
// @Success 200 {object} model.EmployeeResponse
// @Failure 400 {object} gin.H
// @Router /employees/{id} [put]
func (h *employeeHandler) UpdateEmployee(db *gorm.DB, logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.log.Info("Received Update Employee request")

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

		h.log.Infof("Successfully updated employee with ID %v", employee.ID)
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
	h.log.Info("Received Delete Employee request")

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

	h.log.Infof("Employee with ID %d was deleted", employeeID)
	ctx.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}

// AssignShift Додељује смену запосленом
// @Summary Додељује смену запосленом
// @Description Додељује смену запосленом по ID-ју
// @Tags запослени
// @Param id path int true "ID запосленог"
// @Param shift body model.AssignShiftRequest true "Подаци о смени"
// @Success 201 {object} model.AssignShiftResponse
// @Failure 400 {object} gin.H
// @Router /employees/{id}/shifts [post]
func (h *employeeHandler) AssignShift(ctx *gin.Context) {
	employeeIDParam := ctx.Param("id")
	h.log.Infof("Received Assign Shift request for employee ID %s", employeeIDParam)

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
	shiftDate, err := time.Parse(time.DateOnly, req.ShiftDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid shiftDate format, expected YYYY-MM-DD"})
		return
	}

	// Call the repository method
	assignmentId, err := h.shiftsRepo.AssignEmployee(shiftDate, req.ShiftType, uint(employeeID), req.ProfileType)
	if errors.Is(err, model.ErrCapacityReached) {
		ctx.JSON(http.StatusConflict, gin.H{"error": "maximum capacity for this role reached in the selected shift"})
	} else if errors.Is(err, model.ErrAlreadyAssigned) {
		ctx.JSON(http.StatusConflict, gin.H{"error": "employee is already assigned to this shift"})
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	// Return success response
	h.log.Infof("Successfully assigned shift for employee ID %d", employeeID)

	response := model.AssignShiftResponse{
		ID: assignmentId,
	}

	ctx.JSON(http.StatusCreated, response)
}

// GetShifts Дохватање смена за запосленог
// @Summary Дохватање смена за запосленог
// @Description Дохватање смена за запосленог по ID-ју
// @Tags запослени
// @Param id path int true "ID запосленог"
// @Success 200 {object} []model.ShiftResponse
// @Router /employees/{id}/shifts [get]
func (h *employeeHandler) GetShifts(ctx *gin.Context) {
	employeeIDParam := ctx.Param("id")

	employeeID, err := strconv.Atoi(employeeIDParam)
	if err != nil || employeeID <= 0 {
		h.log.Error("failed to extract url param: invalid employee ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	// Fetch shifts from the repository
	var shifts []model.Shift
	err = h.shiftsRepo.GetShiftsByEmployeeID(uint(employeeID), &shifts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	h.log.Infof("Successfully retrieved shifts for employee ID %d", employeeID)

	var response []model.ShiftResponse
	for _, shift := range shifts {
		response = append(response, model.ShiftResponse{
			ID:        shift.ID,
			ShiftDate: shift.ShiftDate,
			ShiftType: shift.ShiftType,
		})
	}

	h.log.Info("Successfully mapped shifts to response format")
	h.log.Infof("Returning %d shifts", len(response))
	ctx.JSON(http.StatusOK, response)
}

// GetShiftsAvailability Дохватање доступности смена
// @Summary Дохватање доступности смена
// @Description Дохватање доступности смена за одређени дан
// @Tags запослени
// @Param date query string false "Дан за који се проверава доступност смена"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /shifts/availability [get]
func (h *employeeHandler) GetShiftsAvailability(ctx *gin.Context) {
	h.log.Infof("Received Get Shifts Availability request for date %s", ctx.Query("date"))

	dateParam := ctx.Query("date")
	var date time.Time
	var err error
	if dateParam != "" {
		date, err = time.Parse(time.DateOnly, dateParam)
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

	h.log.Infof("Successfully retrieved shifts availability for date %s", date.Format(time.DateOnly))
	ctx.JSON(http.StatusOK, gin.H{"availability": availability})
}

// RemoveShift Уклањање смене за запосленог
// @Summary Уклањање смене за запосленог
// @Description Уклањање смене за запосленог по ID-ју и подацима о смени
// @Tags запослени
// @Param id path int true "ID запосленог"
// @Param shift body model.RemoveShiftRequest true "Подаци о смени"
// @Success 204
// @Failure 400 {object} gin.H
// @Router /employees/{id}/shifts [delete]
func (h *employeeHandler) RemoveShift(ctx *gin.Context) {
	h.log.Infof("Received Remove Shift request for employee ID %s", ctx.Param("id"))

	var req model.RemoveShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call repository method to remove the shift
	err := h.shiftsRepo.RemoveEmployeeFromShift(req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.log.Infof("Successfully removed shift for employee ID %d", req.ID)
	ctx.JSON(http.StatusNoContent, nil)
}
