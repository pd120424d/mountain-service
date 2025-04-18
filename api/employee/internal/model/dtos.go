package model

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
)

// EmployeeLogin DTO for employee login
// swagger:model
type EmployeeLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

// EmployeeUpdateRequest DTO for updating an existing employee
// swagger:model
type EmployeeUpdateRequest struct {
	FirstName      string `json:"firstName,omitempty"`
	LastName       string `json:"lastName,omitempty"`
	Email          string `json:"email,omitempty"`
	Gender         string `json:"gender,omitempty"`
	Phone          string `json:"phone,omitempty"`
	ProfilePicture string `json:"profilePicture,omitempty"`
	ProfileType    string `json:"profileType,omitempty"`
}

// ShiftResponse DTO for returning shift data for a certain employee
// swagger:model
type ShiftResponse struct {
	ID        uint      `gorm:"primaryKey"`
	ShiftDate time.Time `gorm:"not null"`
	ShiftType int       `gorm:"not null"` // 1: 6am-2pm, 2: 2pm-10pm, 3: 10pm-6am, < 1 or > 3: invalid
	CreatedAt time.Time
}

// AssignShiftRequest DTO for assigning a shift to an employee
// swagger:model
type AssignShiftRequest struct {
	ShiftDate   string `json:"shiftDate" binding:"required"`
	ShiftType   int    `json:"shiftType" binding:"required,min=1,max=3"`
	ProfileType string `json:"profileType" binding:"required,oneof=Medic Technical"`
}

// AssignShiftResponse DTO for returning the shift data after assigning it to an employee
// swagger:model
type AssignShiftResponse struct {
	ID          uint   `json:"id"`
	ShiftDate   string `json:"shiftDate"`
	ShiftType   int    `json:"shiftType"`
	ProfileType string `json:"profileType"`
}

// RemoveShiftRequest DTO for removing a shift from an employee
// swagger:model
type RemoveShiftRequest struct {
	ID uint `json:"id" cinding:"required"`
}

func (e *EmployeeCreateRequest) ToString() string {
	return fmt.Sprintf(
		"EmployeeCreateRequest { FirstName: %s, LastName: %s, Username: %s, Password: %s,"+
			" Email: %s, Gender: %s, Phone: %s, ProfilePicture: %s, ProfileType: %s }",
		e.FirstName,
		e.LastName,
		e.Username,
		sanitizePassword(e.Password),
		e.Email,
		e.Gender,
		e.Phone,
		e.ProfilePicture,
		e.ProfileType,
	)
}

func (e *EmployeeUpdateRequest) Validate() error {
	if e.Email == "" {
		return nil
	}
	_, err := mail.ParseAddress(e.Email)
	return err
}

// Function to sanitize the password by masking it with asterisks
func sanitizePassword(password string) string {
	if password == "" {
		return ""
	}
	return strings.Repeat("*", len(password)) // Replace each character with an asterisk
}
