package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/employee/internal/service"
	sharedAuth "github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/config"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

// Local type aliases for Swagger documentation to avoid long namespace names
type EmployeeCreateRequest = employeeV1.EmployeeCreateRequest
type EmployeeResponse = employeeV1.EmployeeResponse
type EmployeeLogin = employeeV1.EmployeeLogin
type TokenResponse = employeeV1.TokenResponse
type EmployeeUpdateRequest = employeeV1.EmployeeUpdateRequest
type AssignShiftRequest = employeeV1.AssignShiftRequest
type AssignShiftResponse = employeeV1.AssignShiftResponse
type ShiftResponse = employeeV1.ShiftResponse
type ShiftAvailabilityResponse = employeeV1.ShiftAvailabilityResponse
type RemoveShiftRequest = employeeV1.RemoveShiftRequest
type OnCallEmployeesResponse = employeeV1.OnCallEmployeesResponse
type ActiveEmergenciesResponse = employeeV1.ActiveEmergenciesResponse
type ErrorResponse = employeeV1.ErrorResponse
type MessageResponse = employeeV1.MessageResponse

type EmployeeHandler interface {
	// Crud operations, Register is create
	RegisterEmployee(ctx *gin.Context)
	LoginEmployee(ctx *gin.Context)
	LogoutEmployee(ctx *gin.Context)
	OAuth2Token(ctx *gin.Context)
	ListEmployees(ctx *gin.Context)
	GetEmployee(ctx *gin.Context)
	UpdateEmployee(ctx *gin.Context)
	DeleteEmployee(ctx *gin.Context)

	// Shift operations
	AssignShift(ctx *gin.Context)
	GetShifts(ctx *gin.Context)
	GetShiftsAvailability(ctx *gin.Context)
	RemoveShift(ctx *gin.Context)
	GetShiftWarnings(ctx *gin.Context)

	// Emergency operations
	GetOnCallEmployees(ctx *gin.Context)
	CheckActiveEmergencies(ctx *gin.Context)

	// Admin operations
	ResetAllData(ctx *gin.Context)
	GetAdminShiftsAvailability(ctx *gin.Context)

	// Catalog and metadata
	GetErrorCatalog(ctx *gin.Context)
}

type employeeHandler struct {
	log          utils.Logger
	emplService  service.EmployeeService
	shiftService service.ShiftService
}

func NewEmployeeHandler(log utils.Logger, employeeService service.EmployeeService, shiftService service.ShiftService) EmployeeHandler {
	return &employeeHandler{
		log:          log.WithName("employeeHandler"),
		emplService:  employeeService,
		shiftService: shiftService,
	}
}

// RegisterEmployee Креирање новог запосленог
// @Summary Креирање новог запосленог
// @Description Креирање новог запосленог у систему
// @Tags запослени
// @Accept  json
// @Produce  json
// @Param employee body EmployeeCreateRequest true "Подаци о новом запосленом (JSON)"

// @Success 201 {object} EmployeeResponse
// @Failure 400 {object} ErrorResponse
// @Router /employees [post]
func (h *employeeHandler) RegisterEmployee(ctx *gin.Context) {
	reqLog := h.log.WithContext(requestContext(ctx))
	defer utils.TimeOperation(reqLog, "EmployeeHandler.RegisterEmployee")()
	reqLog.Info("Received Register Employee request")

	req := &employeeV1.EmployeeCreateRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request payload: %v", err)})
		return
	}

	if err := req.Validate(); err != nil {
		reqLog.Errorf("validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.emplService.RegisterEmployee(requestContext(ctx), *req)
	if err != nil {
		reqLog.Errorf("failed to register employee: %v", err)

		// Handle specific error types
		switch err.Error() {
		case "username already exists":
			ctx.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		case "email already exists":
			ctx.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		case "invalid profile type":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile type"})
		case "failed to check for existing username":
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for existing username"})
		case "failed to check for existing email":
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for existing email"})
		default:
			// Check if it's a password validation error
			if strings.Contains(err.Error(), "password must") {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else if strings.Contains(err.Error(), "invalid db") {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register employee"})
			}
		}
		return
	}

	reqLog.Infof("Successfully registered employee with username %s", response.Username)

	utils.WriteFreshWindow(ctx, config.DefaultFreshWindow)

	ctx.JSON(http.StatusCreated, response)
}

// LoginEmployee Пријавање запосленог
// @Summary Пријавање запосленог
// @Description Пријавање запосленог са корисничким именом и лозинком
// @Tags запослени
// @Accept  json
// @Produce  json
// @Param employee body EmployeeLogin true "Корисничко име и лозинка"
// @Success 200 {object} TokenResponse
// @Failure 401 {object} ErrorResponse
// @Router /login [post]
func (h *employeeHandler) LoginEmployee(ctx *gin.Context) {
	var req employeeV1.EmployeeLogin
	log := h.log.WithContext(requestContext(ctx))
	defer utils.TimeOperation(log, "EmployeeHandler.LoginEmployee")()
	log.Info("Received Login Employee request")

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorf("Failed to bind login request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request payload: %v", err)})
		return
	}

	if err := req.Validate(); err != nil {
		log.Errorf("validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if sharedAuth.IsAdminLogin(req.Username) {
		log.Info("Admin login attempt detected")

		if !sharedAuth.ValidateAdminPassword(req.Password) {
			log.Error("Invalid admin password")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := sharedAuth.GenerateAdminJWT()
		if err != nil {
			log.Errorf("failed to generate admin token: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		log.Info("Successfully authenticated admin user")
		ctx.JSON(http.StatusOK, gin.H{"token": token})
		return
	}

	token, err := h.emplService.LoginEmployee(requestContext(ctx), req)
	if err != nil {
		log.Errorf("failed to login employee: %v", err)

		if err.Error() == "invalid credentials" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
		}
		return
	}

	log.Info("Successfully validated employee and generated JWT token")
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// LogoutEmployee Одјављивање запосленог
// @Summary Одјављивање запосленог
// @Description Одјављивање запосленог и поништавање токена
// @Tags запослени
// @Security OAuth2Password
// @Produce json
// @Success 200 {object} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /logout [post]
func (h *employeeHandler) LogoutEmployee(ctx *gin.Context) {
	log := h.log.WithContext(requestContext(ctx))
	defer utils.TimeOperation(log, "EmployeeHandler.LogoutEmployee")()
	log.Info("Received Logout Employee request")

	tokenID, exists := ctx.Get("tokenID")
	if !exists {
		log.Error("Token ID not found in context")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	tokenIDStr, ok := tokenID.(string)
	if !ok {
		log.Error("Token ID is not a string")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Prefer expiration from context set by middleware to avoid re-parsing
	expiresAny, ok := ctx.Get("expiresAt")
	if !ok {
		// Fallback: parse token to get expiration (blacklist not used here)
		authHeader := ctx.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := sharedAuth.ValidateJWT(tokenString, nil)
		if err != nil {
			log.Errorf("failed to parse token for logout: %v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		expiresAny = claims.ExpiresAt.Time
	}
	expiresAt, _ := expiresAny.(time.Time)
	if err := h.emplService.LogoutEmployee(requestContext(ctx), tokenIDStr, expiresAt); err != nil {
		log.Errorf("failed to logout employee: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	log.Info("Successfully logged out employee")
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// OAuth2Token OAuth2 token endpoint for Swagger UI
// @Summary OAuth2 token endpoint
// @Description OAuth2 password flow token endpoint for Swagger UI authentication
// @Tags authentication
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {object} map[string]interface{} "OAuth2 token response"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /oauth/token [post]
func (h *employeeHandler) OAuth2Token(ctx *gin.Context) {
	log := h.log.WithContext(requestContext(ctx))
	defer utils.TimeOperation(log, "EmployeeHandler.OAuth2Token")()
	log.Info("Received OAuth2 Token request")

	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	req := employeeV1.EmployeeLogin{
		Username: username,
		Password: password,
	}

	if err := req.Validate(); err != nil {
		log.Errorf("validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if it's admin login
	if sharedAuth.IsAdminLogin(username) {
		log.Info("Admin OAuth2 login attempt detected")

		if !sharedAuth.ValidateAdminPassword(password) {
			log.Error("Invalid admin password")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := sharedAuth.GenerateAdminJWT()
		if err != nil {
			log.Errorf("failed to generate admin token: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		log.Info("Successfully authenticated admin user via OAuth2")
		ctx.JSON(http.StatusOK, gin.H{
			"access_token": token,
			"token_type":   "Bearer",
			"expires_in":   86400, // 24 hours in seconds
		})
		return
	}

	token, err := h.emplService.LoginEmployee(requestContext(ctx), req)
	if err != nil {
		log.Errorf("failed to login employee: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	log.Info("Successfully validated employee and generated JWT token via OAuth2")
	ctx.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   86400, // 24 hours in seconds
	})
}

// ListEmployees Преузимање листе запослених
// @Summary Преузимање листе запослених
// @Description Преузимање свих запослених
// @Tags запослени
// @Security OAuth2Password
// @Produce  json
// @Success 200 {array} []EmployeeResponse
// @Router /employees [get]
func (h *employeeHandler) ListEmployees(ctx *gin.Context) {
	base := requestContext(ctx)
	cctx, cancel := context.WithTimeout(base, config.DefaultListTimeout)
	defer cancel()
	log := h.log.WithContext(cctx)
	defer utils.TimeOperation(log, "EmployeeHandler.ListEmployees")()
	log.Info("Received List Employees request")

	employees, err := h.emplService.ListEmployees(cctx)
	if err != nil {
		log.Errorf("failed to retrieve employees: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve employees"})
		return
	}

	log.Infof("Successfully retrieved %d employees", len(employees))
	ctx.JSON(http.StatusOK, employees)
}

// GetEmployee Преузимање запосленог по ID-ју
// @Summary Преузимање запосленог по ID-ју
// @Description Преузимање запосленог по ID-ју
// @Tags запослени
// @Security OAuth2Password
// @Param id path int true "ID запосленог"
// @Produce  json
// @Success 200 {object} EmployeeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /employees/{id} [get]
func (h *employeeHandler) GetEmployee(ctx *gin.Context) {
	base := requestContext(ctx)
	cctx, cancel := context.WithTimeout(base, config.DefaultListTimeout)
	defer cancel()
	log := h.log.WithContext(cctx)
	defer utils.TimeOperation(log, "EmployeeHandler.GetEmployee")()
	log.Info("Received Get Employee request")

	idParam := ctx.Param("id")
	employeeID, err := strconv.Atoi(idParam)
	if err != nil {
		log.Errorf("failed to convert employee ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	employee, err := h.emplService.GetEmployeeByID(cctx, uint(employeeID))
	if err != nil {
		log.Errorf("failed to retrieve employee: %v", err)
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve employee"})
		}
		return
	}

	response := &employeeV1.EmployeeResponse{
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

	log.Infof("Successfully retrieved employee with ID %d", employeeID)
	ctx.JSON(http.StatusOK, response)
}

// UpdateEmployee Ажурирање запосленог
// @Summary Ажурирање запосленог
// @Description Ажурирање запосленог по ID-ју
// @Tags запослени
// @Security OAuth2Password
// @Param id path int true "ID запосленог"
// @Param employee body EmployeeUpdateRequest true "Подаци за ажурирање запосленог"
// @Success 200 {object} EmployeeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /employees/{id} [put]
func (h *employeeHandler) UpdateEmployee(ctx *gin.Context) {
	log := h.log.WithContext(requestContext(ctx))
	defer utils.TimeOperation(log, "EmployeeHandler.UpdateEmployee")()
	log.Info("Received Update Employee request")

	employeeID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Errorf("failed to convert employee ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var req employeeV1.EmployeeUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorf("failed to update employee, invalid employee update payload: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := req.Validate(); err != nil {
		log.Errorf("validation failed: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.emplService.UpdateEmployee(requestContext(ctx), uint(employeeID), req)
	if err != nil {
		log.Errorf("failed to update employee: %v", err)

		switch err.Error() {
		case "employee not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		default:
			// Check if it's a validation error
			if strings.Contains(err.Error(), "mail:") || strings.Contains(err.Error(), "@") {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {

				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
			}
		}
		return
	}

	utils.WriteFreshWindow(ctx, config.DefaultFreshWindow)

	ctx.JSON(http.StatusOK, response)
}

// DeleteEmployee Брисање запосленог
// @Summary Брисање запосленог
// @Description Брисање запосленог по ID-ју
// @Tags запослени
// @Security OAuth2Password
// @Param id path int true "ID запосленог"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Router /employees/{id} [delete]
func (h *employeeHandler) DeleteEmployee(ctx *gin.Context) {
	log := h.log.WithContext(requestContext(ctx))
	defer utils.TimeOperation(log, "EmployeeHandler.DeleteEmployee")()
	log.Info("Received Delete Employee request")

	idParam := ctx.Param("id")
	employeeID, err := strconv.Atoi(idParam)
	if err != nil {
		log.Errorf("failed to convert employee ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	err = h.emplService.DeleteEmployee(requestContext(ctx), uint(employeeID))
	if err != nil {
		log.Errorf("failed to delete employee: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
		return
	}

	log.Infof("Employee with ID %d was deleted", employeeID)

	utils.WriteFreshWindow(ctx, config.DefaultFreshWindow)

	ctx.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}

// AssignShift Додељује смену запосленом
// @Summary Додељује смену запосленом
// @Description Додељује смену запосленом по ID-ју
// @Tags запослени
// @Security OAuth2Password
// @Param id path int true "ID запосленог"
// @Param shift body AssignShiftRequest true "Подаци о смени"
// @Success 201 {object} AssignShiftResponse
// @Failure 400 {object} ErrorResponse
// @Router /employees/{id}/shifts [post]
func (h *employeeHandler) AssignShift(ctx *gin.Context) {
	employeeIDParam := ctx.Param("id")
	log := h.log.WithContext(requestContext(ctx))
	defer utils.TimeOperation(log, "EmployeeHandler.AssignShift")()
	log.Infof("Received Assign Shift request for employee ID %s", employeeIDParam)

	employeeID, err := strconv.Atoi(employeeIDParam)
	if err != nil || employeeID <= 0 {
		log.Errorf("failed to extract url param, invalid employee ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee ID"})
		return
	}

	var req employeeV1.AssignShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorf("failed to assign shift, invalid shift payload: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.shiftService.AssignShift(requestContext(ctx), uint(employeeID), req)
	if err != nil {
		log.Errorf("failed to assign shift: %v", err)

		if strings.HasPrefix(err.Error(), model.ErrorConsecutiveShiftsLimit) {
			parts := strings.Split(err.Error(), "|")
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "SHIFT_ERRORS.CONSECUTIVE_SHIFTS_LIMIT",
				"limit": parts[len(parts)-1],
			})

			return
		}
		if strings.HasPrefix(err.Error(), "shift capacity is full for ") || strings.HasPrefix(err.Error(), "maximum capacity for this role reached") {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}

		switch err.Error() {
		case "employee not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "invalid shift date format", "shift date must be in the future", "shift date cannot be more than 3 months in the future", "cannot assign shift in the past", "cannot assign shift more than 3 months in advance":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "employee is already assigned to this shift":
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	log.Infof("Successfully assigned shift for employee ID %d", employeeID)

	utils.WriteFreshWindow(ctx, config.DefaultFreshWindow)

	ctx.JSON(http.StatusCreated, response)
}

// GetShifts Дохватање смена за запосленог
// @Summary Дохватање смена за запосленог
// @Description Дохватање смена за запосленог по ID-ју
// @Tags запослени
// @Security OAuth2Password
// @Param id path int true "ID запосленог"
// @Success 200 {object} []ShiftResponse
// @Router /employees/{id}/shifts [get]
func (h *employeeHandler) GetShifts(ctx *gin.Context) {
	employeeIDParam := ctx.Param("id")
	base := requestContext(ctx)
	cctx, cancel := context.WithTimeout(base, config.DefaultListTimeout)
	defer cancel()
	log := h.log.WithContext(cctx)
	defer utils.TimeOperation(log, "EmployeeHandler.GetShifts")()
	log.Infof("Received Get Shifts request for employee ID %s", employeeIDParam)

	employeeID, err := strconv.Atoi(employeeIDParam)
	if err != nil || employeeID <= 0 {
		log.Errorf("failed to extract url param, invalid employee ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	response, err := h.shiftService.GetShifts(cctx, uint(employeeID))
	if err != nil {
		log.Errorf("failed to get shifts for employee ID %d: %v", employeeID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	log.Infof("Successfully retrieved %d shifts for employee ID %d", len(response), employeeID)
	ctx.JSON(http.StatusOK, response)
}

// GetShiftsAvailability Дохватање доступности смена
// @Summary Дохватање доступности смена
// @Description Дохватање доступности смена за одређени дан
// @Tags запослени
// @Security OAuth2Password
// @Param date query string false "Дан за који се проверава доступност смена"
// @Success 200 {object} ShiftAvailabilityResponse
// @Failure 400 {object} ErrorResponse
// @Router /shifts/availability [get]
func (h *employeeHandler) GetShiftsAvailability(ctx *gin.Context) {
	base := requestContext(ctx)
	cctx, cancel := context.WithTimeout(base, config.DefaultListTimeout)
	defer cancel()
	log := h.log.WithContext(cctx)
	defer utils.TimeOperation(log, "EmployeeHandler.GetShiftsAvailability")()
	log.Infof("Received Get Shifts Availability request for the next %s days", ctx.Query("days"))

	// Extract employee ID from authentication context
	employeeIDValue, exists := ctx.Get("employeeID")
	if !exists {
		log.Errorf("employee ID not found in context")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	employeeID, ok := employeeIDValue.(uint)
	if !ok || employeeID <= 0 {
		log.Errorf("invalid employee ID in context: %v", employeeIDValue)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid employee ID"})
		return
	}

	daysStr := ctx.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		log.Errorf("failed to extract url param, invalid days: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid days parameter"})
		return
	}

	response, err := h.shiftService.GetShiftsAvailability(cctx, employeeID, days)
	if err != nil {
		log.Errorf("failed to get shifts availability: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	log.Infof("Successfully retrieved shifts availability for employee %d for the next %v days", employeeID, days)
	ctx.JSON(http.StatusOK, response)
}

// RemoveShift Уклањање смене за запосленог
// @Summary Уклањање смене за запосленог
// @Description Уклањање смене за запосленог по ID-ју и подацима о смени
// @Tags запослени
// @Security OAuth2Password
// @Param id path int true "ID запосленог"
// @Param shift body RemoveShiftRequest true "Подаци о смени"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Router /employees/{id}/shifts [delete]
func (h *employeeHandler) RemoveShift(ctx *gin.Context) {
	employeeIDParam := ctx.Param("id")
	log := h.log.WithContext(requestContext(ctx))
	defer utils.TimeOperation(log, "EmployeeHandler.RemoveShift")()
	log.Infof("Received Remove Shift request for employee ID %s", employeeIDParam)

	employeeID, err := strconv.Atoi(employeeIDParam)
	if err != nil || employeeID <= 0 {
		log.Errorf("failed to extract url param, invalid employee ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var req employeeV1.RemoveShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorf("failed to remove shift, invalid shift payload: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.shiftService.RemoveShift(requestContext(ctx), uint(employeeID), req)
	if err != nil {
		log.Errorf("failed to remove shift: %v", err)

		if err.Error() == "invalid shift date format" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	log.Infof("Successfully removed shift for employee ID %d", employeeID)

	utils.WriteFreshWindow(ctx, config.DefaultFreshWindow)

	ctx.JSON(http.StatusNoContent, nil)
}

// ResetAllData Ресетовање свих података (само за админе)
// @Summary Ресетовање свих података

// @Description Брише све запослене, смене и повезане податке из система (само за админе)
// @Tags админ
// @Security OAuth2Password
// @Produce json
// @Success 200 {object} MessageResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/reset [delete]
func (h *employeeHandler) ResetAllData(ctx *gin.Context) {
	log := h.log.WithContext(requestContext(ctx))
	defer utils.TimeOperation(log, "EmployeeHandler.ResetAllData")()
	log.Warn("Admin data reset request received")

	err := h.emplService.ResetAllData(requestContext(ctx))
	if err != nil {
		log.Errorf("Failed to reset all data: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset data"})
		return
	}

	log.Info("Successfully reset all system data")
	utils.WriteFreshWindow(ctx, config.DefaultFreshWindow)

	ctx.JSON(http.StatusOK, gin.H{"message": "All data has been successfully reset"})
}

// GetAdminShiftsAvailability Дохватање доступности смена за админе
// @Summary Дохватање доступности смена за админе
// @Description Дохватање доступности смена за све запослене (само за админе)
// @Tags админ
// @Security OAuth2Password
// @Param days query int false "Број дана за које се проверава доступност (подразумевано 7)"
// @Success 200 {object} ShiftAvailabilityResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/shifts/availability [get]
func (h *employeeHandler) GetAdminShiftsAvailability(ctx *gin.Context) {
	log := h.log.WithContext(requestContext(ctx))
	defer utils.TimeOperation(log, "EmployeeHandler.GetAdminShiftsAvailability")()
	log.Info("Admin shifts availability request received")

	// Parse days parameter (default to 7)
	days := 7
	if daysStr := ctx.Query("days"); daysStr != "" {
		var err error
		days, err = strconv.Atoi(daysStr)
		if err != nil || days <= 0 || days > 90 {
			log.Errorf("Invalid days parameter: %s", daysStr)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Days must be a number between 1 and 90"})
			return
		}
	}

	// For admin, we can get availability for all employees or system-wide availability
	// For now, let's return system-wide availability (we can use employee ID 1 as a reference)
	response, err := h.shiftService.GetShiftsAvailability(requestContext(ctx), 1, days)
	if err != nil {
		log.Errorf("failed to get admin shifts availability: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve shifts availability"})
		return
	}

	log.Infof("Successfully retrieved admin shifts availability for %d days", days)
	ctx.JSON(http.StatusOK, response)
}

// GetOnCallEmployees Претрага запослених који су тренутно на дужности
// @Summary Претрага запослених који су тренутно на дужности
// @Description Враћа листу запослених који су тренутно на дужности, са опционим бафером у случају да се близу крај тренутне смене
// @Tags запослени
// @Security OAuth2Password
// @Accept  json
// @Produce  json
// @Param shift_buffer query string false "Бафер време пре краја смене (нпр. '1h')"
// @Success 200 {object} OnCallEmployeesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /employees/on-call [get]
func (h *employeeHandler) GetOnCallEmployees(ctx *gin.Context) {
	base := requestContext(ctx)
	cctx, cancel := context.WithTimeout(base, config.DefaultListTimeout)
	defer cancel()
	log := h.log.WithContext(cctx)
	defer utils.TimeOperation(log, "EmployeeHandler.GetOnCallEmployees")()
	log.Info("Getting on-call employees")

	var shiftBuffer time.Duration
	if bufferStr := ctx.Query("shift_buffer"); bufferStr != "" {
		var err error
		shiftBuffer, err = time.ParseDuration(bufferStr)
		if err != nil {
			log.Errorf("Invalid shift_buffer parameter: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shift_buffer format. Use format like '1h', '30m'"})
			return
		}
	}

	employeeResponses, err := h.shiftService.GetOnCallEmployees(cctx, time.Now().UTC(), shiftBuffer)
	if err != nil {
		log.Errorf("Failed to get on-call employees: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve on-call employees"})
		return
	}

	response := employeeV1.OnCallEmployeesResponse{
		Employees: employeeResponses,
	}

	log.Infof("Successfully retrieved %d on-call employees", len(employeeResponses))
	ctx.JSON(http.StatusOK, response)
}

// CheckActiveEmergencies Провера активних хитних случајева за запосленог
// @Summary Провера активних хитних случајева за запосленог
// @Description Проверава да ли запослени има активне хитне случајеве
// @Tags запослени
// @Security OAuth2Password
// @Accept  json
// @Produce  json
// @Param id path int true "ID запосленог"
// @Success 200 {object} ActiveEmergenciesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /employees/{id}/active-emergencies [get]
func (h *employeeHandler) CheckActiveEmergencies(ctx *gin.Context) {
	log := h.log.WithContext(requestContext(ctx))
	employeeIDStr := ctx.Param("id")
	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
	if err != nil {
		log.Errorf("Invalid employee ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	log.Infof("Checking active emergencies for employee %d", employeeID)

	// Check if employee exists
	_, err = h.emplService.GetEmployeeByID(requestContext(ctx), uint(employeeID))
	if err != nil {
		log.Errorf("Employee not found: %v", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// TODO: This will be implemented in the scope of integration with urgency service
	// For now, return false as placeholder
	response := employeeV1.ActiveEmergenciesResponse{
		HasActiveEmergencies: false,
	}

	log.Infof("Employee %d has active emergencies: %v", employeeID, response.HasActiveEmergencies)
	ctx.JSON(http.StatusOK, response)
}

// GetShiftWarnings Дохватање упозорења о сменама за запосленог
// @Summary Дохватање упозорења о сменама за запосленог
// @Description Враћа листу упозорења о сменама за запосленог (нпр. недостају смене, није испуњена норма)
// @Tags запослени
// @Security OAuth2Password
// @Param id path int true "ID запосленог"
// @Success 200 {object} map[string][]string
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /employees/{id}/shift-warnings [get]
func (h *employeeHandler) GetShiftWarnings(ctx *gin.Context) {
	log := h.log.WithContext(requestContext(ctx))
	employeeIDParam := ctx.Param("id")
	log.Infof("Received Get Shift Warnings request for employee ID %s", employeeIDParam)

	employeeID, err := strconv.Atoi(employeeIDParam)
	if err != nil || employeeID <= 0 {
		log.Errorf("failed to extract url param, invalid employee ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	warnings, err := h.shiftService.GetShiftWarnings(requestContext(ctx), uint(employeeID))
	if err != nil {
		log.Errorf("failed to get shift warnings for employee ID %d: %v", employeeID, err)

		if err.Error() == "employee not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	log.Infof("Successfully retrieved %d warnings for employee ID %d", len(warnings), employeeID)
	ctx.JSON(http.StatusOK, gin.H{"warnings": warnings})
}

// GetErrorCatalog returns a catalog of error codes/messages used by employee-service
func (h *employeeHandler) GetErrorCatalog(ctx *gin.Context) {
	type catalogEntry struct {
		Code          string            `json:"code"`
		Service       string            `json:"service"`
		HttpStatus    int               `json:"httpStatus"`
		DefaultMsg    string            `json:"defaultMessage"`
		DetailsSchema map[string]string `json:"detailsSchema,omitempty"`
	}
	errors := []catalogEntry{
		{Code: model.ErrorConsecutiveShiftsLimit, Service: "employee-service", HttpStatus: http.StatusConflict, DefaultMsg: "Exceeded consecutive days limit", DetailsSchema: map[string]string{"limit": "number"}},
		{Code: "SHIFT_ERRORS.ALREADY_ASSIGNED", Service: "employee-service", HttpStatus: http.StatusConflict, DefaultMsg: "Employee is already assigned to this shift"},
		{Code: "SHIFT_ERRORS.CAPACITY_FULL", Service: "employee-service", HttpStatus: http.StatusConflict, DefaultMsg: "Shift capacity is full for role"},
		{Code: "VALIDATION.INVALID_SHIFT_DATE", Service: "employee-service", HttpStatus: http.StatusBadRequest, DefaultMsg: "Invalid shift date format"},
		{Code: "VALIDATION.SHIFT_IN_PAST", Service: "employee-service", HttpStatus: http.StatusBadRequest, DefaultMsg: "Shift date must be in the future"},
		{Code: "VALIDATION.SHIFT_TOO_FAR", Service: "employee-service", HttpStatus: http.StatusBadRequest, DefaultMsg: "Shift date cannot be more than 3 months in the future"},
		{Code: "EMPLOYEE_ERRORS.NOT_FOUND", Service: "employee-service", HttpStatus: http.StatusNotFound, DefaultMsg: "Employee not found"},
	}
	ctx.JSON(http.StatusOK, gin.H{
		"service":  "employee-service",
		"errors":   errors,
		"warnings": []catalogEntry{{Code: model.WarningInsufficientShifts, Service: "employee-service", HttpStatus: http.StatusOK, DefaultMsg: "Insufficient shifts in the next period", DetailsSchema: map[string]string{"count": "number", "periodDays": "number", "perWeek": "number"}}},
	})
}

// requestContext safely extracts a context from gin.Context; falls back to Background.
func requestContext(ctx *gin.Context) context.Context {
	if ctx != nil && ctx.Request != nil {
		return ctx.Request.Context()
	}
	return context.Background()
}
