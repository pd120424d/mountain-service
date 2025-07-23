package model

import (
	"net/mail"
)

// Swagger documentation types - these mirror the contract types exactly
// This allows Swagger to generate proper documentation while the actual
// implementation uses the contract types from the contracts package

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
	Username string `json:"username"`
	Password string `json:"password"`
}

// TokenResponse DTO for returning authentication token
// swagger:model
type TokenResponse struct {
	Token string `json:"token"`
}

// EmployeeResponse DTO for returning employee data
// swagger:model
type EmployeeResponse struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Gender         string `json:"gender"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
	ProfileType    string `json:"profileType"`
}

// EmployeeCreateRequest DTO for creating a new employee
// swagger:model
type EmployeeCreateRequest struct {
	FirstName      string `json:"firstName" binding:"required"`
	LastName       string `json:"lastName" binding:"required"`
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Gender         string `json:"gender" binding:"required"`
	Phone          string `json:"phone" binding:"required"`
	ProfilePicture string `json:"profilePicture"`
	ProfileType    string `json:"profileType"`
}

// EmployeeUpdateRequest DTO for updating an employee
// swagger:model
type EmployeeUpdateRequest struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email" binding:"email"`
	Gender         string `json:"gender"`
	Phone          string `json:"phone"`
	ProfilePicture string `json:"profilePicture"`
	ProfileType    string `json:"profileType"`
}

// AllEmployeesResponse DTO for returning all employees
// swagger:model
type AllEmployeesResponse struct {
	Employees []EmployeeResponse `json:"employees"`
}

// AssignShiftRequest DTO for assigning a shift to an employee
// swagger:model
type AssignShiftRequest struct {
	ShiftDate string `json:"shiftDate" binding:"required"`
	ShiftType string `json:"shiftType" binding:"required"`
}

// AssignShiftResponse DTO for returning shift assignment result
// swagger:model
type AssignShiftResponse struct {
	ID        uint   `json:"id"`
	ShiftDate string `json:"shiftDate"`
	ShiftType string `json:"shiftType"`
}

// RemoveShiftRequest DTO for removing a shift from an employee
// swagger:model
type RemoveShiftRequest struct {
	ShiftDate string `json:"shiftDate" binding:"required"`
	ShiftType string `json:"shiftType" binding:"required"`
}

// ShiftResponse DTO for returning shift data
// swagger:model
type ShiftResponse struct {
	ID        uint   `json:"id"`
	ShiftDate string `json:"shiftDate"`
	ShiftType string `json:"shiftType"`
	CreatedAt string `json:"createdAt"`
}

// ShiftAvailabilityPerDay represents availability for a single day
// swagger:model
type ShiftAvailabilityPerDay struct {
	Available bool     `json:"available"`
	Employees []string `json:"employees"`
}

// ShiftAvailabilityResponse DTO for returning shift availability
// swagger:model
type ShiftAvailabilityResponse struct {
	Days map[string]struct {
		Shift1 ShiftAvailabilityPerDay `json:"shift1"`
		Shift2 ShiftAvailabilityPerDay `json:"shift2"`
		Shift3 ShiftAvailabilityPerDay `json:"shift3"`
	} `json:"days"`
}

// OnCallEmployeesResponse DTO for returning on-call employees
// swagger:model
type OnCallEmployeesResponse struct {
	Employees []EmployeeResponse `json:"employees"`
}

// ActiveEmergenciesResponse DTO for returning active emergencies status
// swagger:model
type ActiveEmergenciesResponse struct {
	HasActiveEmergencies bool `json:"hasActiveEmergencies"`
}

// Validate validates the EmployeeUpdateRequest
func (r *EmployeeUpdateRequest) Validate() error {
	if r.Email != "" {
		_, err := mail.ParseAddress(r.Email)
		return err
	}
	return nil
}

// sanitizePassword masks the password with asterisks
func sanitizePassword(password string) string {
	if password == "" {
		return ""
	}
	return "********"
}
