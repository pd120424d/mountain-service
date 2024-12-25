package model

import (
	"fmt"
	"strings"
)

// EmployeeResponse DTO for returning employee data
type EmployeeResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Gender    string `json:"gender"`
	Phone     string `json:"phoneNumber"`
	Email     string `json:"email"`
	// this may be represented as a byte array if we read the picture from somewhere for an example
	ProfilePicture string `json:"profilePicture"`
	ProfileType    string `json:"profileType"`
}

// EmployeeCreateRequest DTO for creating a new employee
type EmployeeCreateRequest struct {
	FirstName      string `json:"firstName" binding:"required"`
	LastName       string `json:"lastName" binding:"required"`
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Gender         string `json:"gender" binding:"required"`
	Phone          string `json:"phoneNumber" binding:"required"`
	ProfilePicture string `json:"profilePicture"`
	ProfileType    string `json:"profileType"`
}

// EmployeeUpdateRequest DTO for updating an existing employee
type EmployeeUpdateRequest struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Age       int    `json:"age,omitempty"`
	Email     string `json:"email,omitempty"`
}

type AssignShiftRequest struct {
	ShiftDate   string `json:"shiftDate" binding:"required"`
	ShiftType   int    `json:"shiftType" binding:"required,min=1,max=3"`
	ProfileType string `json:"profileType" binding:"required,oneof=Medic Technical"`
}

type RemoveShiftRequest struct {
	ShiftDate string `json:"shiftDate" binding:"required"`
	ShiftType int    `json:"shiftType" binding:"required,min=1,max=3"`
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

// Function to sanitize the password by masking it with asterisks
func sanitizePassword(password string) string {
	if password == "" {
		return ""
	}
	return strings.Repeat("*", len(password)) // Replace each character with an asterisk
}
