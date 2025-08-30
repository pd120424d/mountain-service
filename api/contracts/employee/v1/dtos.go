package v1

import (
	"fmt"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/shared/validation"
)

// Common response types
// ErrorResponse represents an error response
// swagger:model
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse represents a simple message response
// swagger:model
type MessageResponse struct {
	Message string `json:"message"`
}

// EmployeeLogin DTO for employee login
// swagger:model
type EmployeeLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// EmployeeResponse DTO for returning employee data
// swagger:model
type EmployeeResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Gender    string `json:"gender"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	// this may be represented as a byte array if we read the picture from somewhere for an example
	ProfilePicture string `json:"profilePicture"`
	ProfileType    string `json:"profileType"`
}

// EmployeeCreateRequest DTO for creating a new employee
// swagger:model
type EmployeeCreateRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Gender    string `json:"gender" binding:"required"`
	Phone     string `json:"phone" binding:"required"`

	ProfileType string `json:"profileType" binding:"required"`
}

// Validate validates the EmployeeCreateRequest
func (r *EmployeeCreateRequest) Validate() error {
	var errors validation.ValidationErrors

	if err := utils.ValidateRequiredField(r.FirstName, "first name"); err != nil {
		errors.AddError("firstName", err)
	}
	if err := utils.ValidateRequiredField(r.LastName, "last name"); err != nil {
		errors.AddError("lastName", err)
	}
	if err := utils.ValidateRequiredField(r.Username, "username"); err != nil {
		errors.AddError("username", err)
	}
	if err := utils.ValidateRequiredField(r.Password, "password"); err != nil {
		errors.AddError("password", err)
	}
	if err := utils.ValidateRequiredField(r.ProfileType, "profile type"); err != nil {
		errors.AddError("profileType", err)
	}

	if err := utils.ValidateEmail(r.Email); err != nil {
		errors.AddError("email", err)
	}

	if err := utils.ValidateGender(r.Gender); err != nil {
		errors.AddError("gender", err)
	}

	if err := utils.ValidatePhone(r.Phone); err != nil {
		errors.AddError("phone", err)
	}

	if err := r.validateProfileType(); err != nil {
		errors.AddError("profileType", err)
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

func (r *EmployeeCreateRequest) validateProfileType() error {
	validProfileTypes := []string{"Medic", "Technical", "Administrator"}
	for _, validType := range validProfileTypes {
		if r.ProfileType == validType {
			return nil
		}
	}
	return fmt.Errorf("profile type must be one of: Medic, Technical, Administrator")
}

// EmployeeUpdateRequest DTO for updating an employee
// swagger:model
type EmployeeUpdateRequest struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Gender         string `json:"gender"`
	Phone          string `json:"phone"`
	ProfilePicture string `json:"profilePicture"`
	ProfileType    string `json:"profileType"`
}

// ShiftResponse DTO for returning shift data for a certain employee
// swagger:model
type ShiftResponse struct {
	ID         uint      `json:"id"`
	ShiftDate  time.Time `json:"shiftDate"`
	ShiftType  int       `json:"shiftType"` // 1: 6am-2pm, 2: 2pm-10pm, 3: 10pm-6am, < 1 or > 3: invalid
	CreatedAt  time.Time `json:"createdAt"`
	AssignedAt time.Time `json:"assignedAt"`
}

// AssignShiftRequest DTO for assigning a shift to an employee
// swagger:model
type AssignShiftRequest struct {
	ShiftDate string `json:"shiftDate" binding:"required"`
	ShiftType int    `json:"shiftType" binding:"required,min=1,max=3"`
}

// AssignShiftResponse DTO for returning the shift data after assigning it to an employee
// swagger:model
type AssignShiftResponse struct {
	ID        uint   `json:"id"  binding:"required"`
	ShiftDate string `json:"shiftDate"  binding:"required"`
	ShiftType int    `json:"shiftType"  binding:"required"`
}

// RemoveShiftRequest DTO for removing a shift from an employee
// swagger:model
type RemoveShiftRequest struct {
	ShiftType int    `json:"shiftType" binding:"required,min=1,max=3"`
	ShiftDate string `json:"shiftDate" binding:"required"`
}

// ShiftAvailabilityResponse DTO for returning the shift availability for a certain date
// swagger:model
type ShiftAvailabilityResponse struct {
	Days map[time.Time]ShiftAvailabilityPerDay `json:"days"`
}

// ShiftAvailabilityPerDay DTO for returning the shift availability for a certain day
// swagger:model
type ShiftAvailabilityPerDay struct {
	FirstShift  ShiftAvailability `json:"firstShift"`
	SecondShift ShiftAvailability `json:"secondShift"`
	ThirdShift  ShiftAvailability `json:"thirdShift"`
}

// ShiftAvailability DTO for returning the shift availability for a certain shift
// swagger:model
type ShiftAvailability struct {
	MedicSlotsAvailable     int  `json:"medicSlotsAvailable"`     // Available slots for medics (0-2)
	TechnicalSlotsAvailable int  `json:"technicalSlotsAvailable"` // Available slots for technical staff (0-4)
	IsAssignedToEmployee    bool `json:"isAssignedToEmployee"`    // Whether the requesting employee is assigned to this shift
	IsFullyBooked           bool `json:"isFullyBooked"`           // Whether the shift is at full capacity (2 medics + 4 technicians)
}

// OnCallEmployeesResponse DTO for returning on-call employees
// swagger:model
type OnCallEmployeesResponse struct {
	Employees []EmployeeResponse `json:"employees"`
}

// AllEmployeesResponse DTO for returning all employees
// swagger:model
type AllEmployeesResponse struct {
	Employees []EmployeeResponse `json:"employees"`
}

// TokenResponse DTO for returning a JWT token
// swagger:model
type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ActiveEmergenciesResponse DTO for returning active emergencies status
// swagger:model
type ActiveEmergenciesResponse struct {
	HasActiveEmergencies bool `json:"hasActiveEmergencies"`
}

// Helper methods

func (r *RemoveShiftRequest) String() string {
	return fmt.Sprintf("RemoveShiftRequest { ShiftType: %d, ShiftDate: %s }", r.ShiftType, r.ShiftDate)
}

// Validate validates the EmployeeUpdateRequest
func (r *EmployeeUpdateRequest) Validate() error {
	var errors validation.ValidationErrors

	if r.Email != "" {
		if err := utils.ValidateEmail(r.Email); err != nil {
			errors.AddError("email", err)
		}
	}

	if r.Gender != "" {
		if err := utils.ValidateGender(r.Gender); err != nil {
			errors.AddError("gender", err)
		}
	}

	if r.Phone != "" {
		if err := utils.ValidatePhone(r.Phone); err != nil {
			errors.AddError("phone", err)
		}
	}

	if r.ProfileType != "" {
		if err := r.validateProfileType(); err != nil {
			errors.AddError("profileType", err)
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

func (r *EmployeeUpdateRequest) validateProfileType() error {
	validProfileTypes := []string{"Medic", "Technical", "Administrator"}
	for _, validType := range validProfileTypes {
		if r.ProfileType == validType {
			return nil
		}
	}
	return fmt.Errorf("profile type must be one of: Medic, Technical, Administrator")
}

func (e *EmployeeCreateRequest) ToString() string {
	return fmt.Sprintf(
		"EmployeeCreateRequest { FirstName: %s, LastName: %s, Username: %s, Password: %s,"+
			" Email: %s, Gender: %s, Phone: %s, ProfileType: %s }",
		e.FirstName,
		e.LastName,
		e.Username,
		utils.SanitizePassword(e.Password),
		e.Email,
		e.Gender,
		e.Phone,
		e.ProfileType,
	)
}

// Validate validates the EmployeeLogin request
func (r *EmployeeLogin) Validate() error {
	var errors validation.ValidationErrors

	if err := utils.ValidateRequiredField(r.Username, "username"); err != nil {
		errors.AddError("username", err)
	}

	if err := utils.ValidateRequiredField(r.Password, "password"); err != nil {
		errors.AddError("password", err)
	}
	// Note: Password format validation is intentionally skipped for login
	// Admin passwords may not follow employee password rules
	// Password validation is handled during authentication, not at DTO level

	if errors.HasErrors() {
		return errors
	}
	return nil
}
