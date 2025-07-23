package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/employee/internal/auth"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/employee/internal/repositories"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type EmployeeHandler interface {
	// Crud operations, Register is create
	RegisterEmployee(ctx *gin.Context)
	LoginEmployee(ctx *gin.Context)
	ListEmployees(ctx *gin.Context)
	UpdateEmployee(ctx *gin.Context)
	DeleteEmployee(ctx *gin.Context)

	// Shift operations
	AssignShift(ctx *gin.Context)
	GetShifts(ctx *gin.Context)
	GetShiftsAvailability(ctx *gin.Context)
	RemoveShift(ctx *gin.Context)

	// Emergency operations
	GetOnCallEmployees(ctx *gin.Context)
	CheckActiveEmergencies(ctx *gin.Context)

	// Admin operations
	ResetAllData(ctx *gin.Context)
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
// @Failure 400 {object} model.ErrorResponse
// @Router /employees [post]
func (h *employeeHandler) RegisterEmployee(ctx *gin.Context) {
	h.log.Info("Received Register Employee request")
	req := &employeeV1.EmployeeCreateRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request payload: %v", err)})
		return
	}

	h.log.Infof("Creating new employee with data: %s", req.ToString())

	profileType := model.ProfileTypeFromString(req.ProfileType)
	if !profileType.Valid() {
		h.log.Errorf("invalid profile type: %s", req.ProfileType)
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

	usernameFilter := map[string]any{
		"username": employee.Username,
	}
	existingEmployees, err := h.emplRepo.ListEmployees(usernameFilter)
	if err != nil {
		h.log.Error("failed to check for existing username", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for existing username"})
		return
	}
	if len(existingEmployees) > 0 {
		h.log.Errorf("username %s already exists", employee.Username)
		ctx.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	emailFilter := map[string]any{
		"email": employee.Email,
	}
	existingEmployees, err = h.emplRepo.ListEmployees(emailFilter)
	if err != nil {
		h.log.Error("failed to register employee, checking for employee with email failed: %v", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for existing email"})
		return
	}
	if len(existingEmployees) > 0 {
		h.log.Errorf("failed to register employee: email %s already exists", employee.Email)
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	if err := utils.ValidatePassword(employee.Password); err != nil {
		h.log.Errorf("failed to validate password: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.emplRepo.Create(&employee); err != nil {
		h.log.Errorf("failed to create employee: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := employeeV1.EmployeeResponse{
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

// LoginEmployee Пријавање запосленог
// @Summary Пријавање запосленог
// @Description Пријавање запосленог са корисничким именом и лозинком
// @Tags запослени
// @Accept  json
// @Produce  json
// @Param employee body model.EmployeeLogin true "Корисничко име и лозинка"
// @Success 200 {object} model.TokenResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /login [post]
func (h *employeeHandler) LoginEmployee(ctx *gin.Context) {
	var req employeeV1.EmployeeLogin
	h.log.Info("Received Login Employee request")

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request payload: %v", err)})
		return
	}

	if auth.IsAdminLogin(req.Username) {
		h.log.Info("Admin login attempt detected")

		if !auth.ValidateAdminPassword(req.Password) {
			h.log.Error("Invalid admin password")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := auth.GenerateAdminJWT()
		if err != nil {
			h.log.Errorf("failed to generate admin token: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		h.log.Info("Successfully authenticated admin user")
		ctx.JSON(http.StatusOK, gin.H{"token": token})
		return
	}

	employee, err := h.emplRepo.GetEmployeeByUsername(req.Username)
	if err != nil {
		h.log.Errorf("failed to retrieve employee: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !auth.CheckPassword(employee.Password, req.Password) {
		h.log.Error("failed to verify password")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateJWT(employee.ID, employee.Role())
	if err != nil {
		h.log.Errorf("failed to generate token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	h.log.Info("Successfully validate employee and generated JWT token")
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// ListEmployees Преузимање листе запослених
// @Summary Преузимање листе запослених
// @Description Преузимање свих запослених
// @Tags запослени
// @Security BearerAuth
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
	response := make([]employeeV1.EmployeeResponse, 0)
	for _, emp := range employees {
		response = append(response, employeeV1.EmployeeResponse{
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
// @Security BearerAuth
// @Param id path int true "ID запосленог"
// @Param employee body model.EmployeeUpdateRequest true "Подаци за ажурирање запосленог"
// @Success 200 {object} model.EmployeeResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /employees/{id} [put]
func (h *employeeHandler) UpdateEmployee(ctx *gin.Context) {
	h.log.Info("Received Update Employee request")

	employeeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.log.Errorf("failed to convert employee ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var req employeeV1.EmployeeUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.log.Errorf("failed to update employee, invalid employee update payload: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var employee model.Employee
	if err := h.emplRepo.GetEmployeeByID(uint(employeeID), &employee); err != nil {
		h.log.Error("failed to get employee: %v", zap.Error(err))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// Validate here or in middleware
	if validationErr := req.Validate(); validationErr != nil {
		h.log.Errorf("failed to update employee, validation failed: %v", validationErr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	model.MapUpdateRequestToEmployee(&req, &employee)

	if err := h.emplRepo.UpdateEmployee(&employee); err != nil {
		h.log.Errorf("failed to update employee: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
		return
	}

	resp := employee.UpdateResponseFromEmployee()
	ctx.JSON(http.StatusOK, resp)
}

// DeleteEmployee Брисање запосленог
// @Summary Брисање запосленог
// @Description Брисање запосленог по ID-ју
// @Tags запослени
// @Security BearerAuth
// @Param id path int true "ID запосленог"
// @Success 204
// @Failure 404 {object} model.ErrorResponse
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
// @Security BearerAuth
// @Param id path int true "ID запосленог"
// @Param shift body model.AssignShiftRequest true "Подаци о смени"
// @Success 201 {object} model.AssignShiftResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /employees/{id}/shifts [post]
func (h *employeeHandler) AssignShift(ctx *gin.Context) {
	employeeIDParam := ctx.Param("id")
	h.log.Infof("Received Assign Shift request for employee ID %s", employeeIDParam)

	employeeID, err := strconv.Atoi(employeeIDParam)
	if err != nil || employeeID <= 0 {
		h.log.Errorf("failed to extract url param, invalid employee ID: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	var req employeeV1.AssignShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.log.Errorf("failed to assign shift, invalid shift payload: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee := &model.Employee{}
	err = h.emplRepo.GetEmployeeByID(uint(employeeID), employee)
	if err != nil {
		h.log.Errorf("failed to assign shift, failed to get employee: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	profileType := employee.ProfileType

	shiftDate, err := time.Parse(time.DateOnly, req.ShiftDate)
	if err != nil {
		h.log.Errorf("failed to parse shift date: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid shiftDate format, expected YYYY-MM-DD"})
		return
	}

	// Step 1: Get or create shift
	shift, err := h.shiftsRepo.GetOrCreateShift(shiftDate, req.ShiftType)
	if err != nil {
		h.log.Errorf("failed to get/create shift: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not get or create shift"})
		return
	}

	// Step 2: Check existing assignment
	assigned, err := h.shiftsRepo.AssignedToShift(uint(employeeID), shift.ID)
	if err != nil {
		h.log.Errorf("failed to check assignment: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if assigned {
		h.log.Errorf("employee with ID %d is already assigned to the requested shift ID %d", employeeID, shift.ID)
		ctx.JSON(http.StatusConflict, gin.H{"error": "employee is already assigned to this shift"})
		return
	}

	// Step 3: Check capacity
	count, err := h.shiftsRepo.CountAssignmentsByProfile(shift.ID, profileType)
	if err != nil {
		h.log.Errorf("failed to count assignments by profile: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	if (profileType.String() == "Medic" && count >= 2) || (profileType.String() == "Technical" && count >= 4) {
		h.log.Errorf("failed to assign shift: maximum capacity for role %s reached in the selected shift", profileType.String())
		ctx.JSON(http.StatusConflict, gin.H{"error": "maximum capacity for this role reached in the selected shift"})
		return
	}

	// Step 4: Assign
	assignmentID, err := h.shiftsRepo.CreateAssignment(uint(employeeID), shift.ID)
	if err != nil {
		h.log.Errorf("failed to assign shift: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign employee"})
		return
	}

	h.log.Infof("Successfully assigned shift for employee ID %d", employeeID)
	resp := employeeV1.AssignShiftResponse{
		ID:        assignmentID,
		ShiftDate: shift.ShiftDate.Format(time.DateOnly),
		ShiftType: shift.ShiftType,
	}
	ctx.JSON(http.StatusCreated, resp)
}

// GetShifts Дохватање смена за запосленог
// @Summary Дохватање смена за запосленог
// @Description Дохватање смена за запосленог по ID-ју
// @Tags запослени
// @Security BearerAuth
// @Param id path int true "ID запосленог"
// @Success 200 {object} []model.ShiftResponse
// @Router /employees/{id}/shifts [get]
func (h *employeeHandler) GetShifts(ctx *gin.Context) {
	employeeIDParam := ctx.Param("id")
	h.log.Infof("Received Get Shifts request for employee ID %s", employeeIDParam)

	employeeID, err := strconv.Atoi(employeeIDParam)
	if err != nil || employeeID <= 0 {
		h.log.Errorf("failed to extract url param, invalid employee ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var shifts []model.Shift
	err = h.shiftsRepo.GetShiftsByEmployeeID(uint(employeeID), &shifts)
	if err != nil {
		h.log.Errorf("failed to get shifts for employee ID %d: %v", employeeID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	h.log.Infof("Successfully retrieved shifts for employee ID %d", employeeID)

	response := make([]employeeV1.ShiftResponse, 0)
	for _, shift := range shifts {
		response = append(response, employeeV1.ShiftResponse{
			ID:        shift.ID,
			ShiftDate: shift.ShiftDate,
			ShiftType: shift.ShiftType,
			CreatedAt: shift.CreatedAt,
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
// @Security BearerAuth
// @Param date query string false "Дан за који се проверава доступност смена"
// @Success 200 {object} model.ShiftAvailabilityResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /shifts/availability [get]
func (h *employeeHandler) GetShiftsAvailability(ctx *gin.Context) {
	h.log.Infof("Received Get Shifts Availability request for the next %s days", ctx.Query("days"))

	daysStr := ctx.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		h.log.Errorf("failed to extract url param, invalid days: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid days parameter"})
		return
	}

	start := time.Now().Truncate(24 * time.Hour)
	end := start.AddDate(0, 0, days)

	availability, err := h.shiftsRepo.GetShiftAvailability(start, end)
	if err != nil {
		h.log.Errorf("failed to get shifts availability: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.log.Infof("Successfully retrieved shifts availability for the next %v days", days)
	availabilityResponse := model.MapShiftsAvailabilityToResponse(availability)
	ctx.JSON(http.StatusOK, availabilityResponse)
}

// RemoveShift Уклањање смене за запосленог
// @Summary Уклањање смене за запосленог
// @Description Уклањање смене за запосленог по ID-ју и подацима о смени
// @Tags запослени
// @Security BearerAuth
// @Param id path int true "ID запосленог"
// @Param shift body model.RemoveShiftRequest true "Подаци о смени"
// @Success 204
// @Failure 400 {object} model.ErrorResponse
// @Router /employees/{id}/shifts [delete]
func (h *employeeHandler) RemoveShift(ctx *gin.Context) {
	employeeIDParam := ctx.Param("id")
	h.log.Infof("Received Remove Shift request for employee ID %s", employeeIDParam)

	employeeID, err := strconv.Atoi(employeeIDParam)
	if err != nil || employeeID <= 0 {
		h.log.Errorf("failed to extract url param, invalid employee ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var req employeeV1.RemoveShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.log.Errorf("failed to remove shift, invalid shift payload: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shiftDate, err := time.Parse(time.DateOnly, req.ShiftDate)
	if err != nil {
		h.log.Errorf("failed to parse shift date: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid shiftDate format, expected YYYY-MM-DD"})
		return
	}

	err = h.shiftsRepo.RemoveEmployeeFromShiftByDetails(uint(employeeID), shiftDate, req.ShiftType)
	if err != nil {
		h.log.Errorf("failed to remove shift: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.log.Infof("Successfully removed shift for employee ID %d", employeeID)
	ctx.JSON(http.StatusNoContent, nil)
}

// ResetAllData Ресетовање свих података (само за админе)
// @Summary Ресетовање свих података
// @Description Брише све запослене, смене и повезане податке из система (само за админе)
// @Tags админ
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.MessageResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /admin/reset [delete]
func (h *employeeHandler) ResetAllData(ctx *gin.Context) {
	h.log.Warn("Admin data reset request received")

	err := h.emplRepo.ResetAllData()
	if err != nil {
		h.log.Errorf("Failed to reset all data: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset data"})
		return
	}

	h.log.Info("Successfully reset all system data")
	ctx.JSON(http.StatusOK, gin.H{"message": "All data has been successfully reset"})
}

// GetOnCallEmployees Претрага запослених који су тренутно на дужности
// @Summary Претрага запослених који су тренутно на дужности
// @Description Враћа листу запослених који су тренутно на дужности, са опционим бафером у случају да се близу крај тренутне смене
// @Tags запослени
// @Accept  json
// @Produce  json
// @Param shift_buffer query string false "Бафер време пре краја смене (нпр. '1h')"
// @Success 200 {object} model.OnCallEmployeesResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /employees/on-call [get]
func (h *employeeHandler) GetOnCallEmployees(ctx *gin.Context) {
	h.log.Info("Getting on-call employees")

	var shiftBuffer time.Duration
	if bufferStr := ctx.Query("shift_buffer"); bufferStr != "" {
		var err error
		shiftBuffer, err = time.ParseDuration(bufferStr)
		if err != nil {
			h.log.Errorf("Invalid shift_buffer parameter: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shift_buffer format. Use format like '1h', '30m'"})
			return
		}
	}

	employees, err := h.shiftsRepo.GetOnCallEmployees(time.Now(), shiftBuffer)
	if err != nil {
		h.log.Errorf("Failed to get on-call employees: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve on-call employees"})
		return
	}

	var employeeResponses []employeeV1.EmployeeResponse
	for _, emp := range employees {
		employeeResponses = append(employeeResponses, emp.UpdateResponseFromEmployee())
	}

	response := employeeV1.OnCallEmployeesResponse{
		Employees: employeeResponses,
	}

	h.log.Infof("Successfully retrieved %d on-call employees", len(employeeResponses))
	ctx.JSON(http.StatusOK, response)
}

// CheckActiveEmergencies Провера активних хитних случајева за запосленог
// @Summary Провера активних хитних случајева за запосленог
// @Description Проверава да ли запослени има активне хитне случајеве
// @Tags запослени
// @Accept  json
// @Produce  json
// @Param id path int true "ID запосленог"
// @Success 200 {object} model.ActiveEmergenciesResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /employees/{id}/active-emergencies [get]
func (h *employeeHandler) CheckActiveEmergencies(ctx *gin.Context) {
	employeeIDStr := ctx.Param("id")
	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
	if err != nil {
		h.log.Errorf("Invalid employee ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	h.log.Infof("Checking active emergencies for employee %d", employeeID)

	// Check if employee exists
	var employee model.Employee
	if err := h.emplRepo.GetEmployeeByID(uint(employeeID), &employee); err != nil {
		h.log.Errorf("Employee not found: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// TODO: This will be implemented in the scope of integration with urgency service
	// For now, return false as placeholder
	response := employeeV1.ActiveEmergenciesResponse{
		HasActiveEmergencies: false,
	}

	h.log.Infof("Employee %d has active emergencies: %v", employeeID, response.HasActiveEmergencies)
	ctx.JSON(http.StatusOK, response)
}
