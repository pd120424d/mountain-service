package v1

import (
	"fmt"

	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/shared/validation"
)

// UrgencyLevel represents the urgency level
type UrgencyLevel string

const (
	Low      UrgencyLevel = "low"
	Medium   UrgencyLevel = "medium"
	High     UrgencyLevel = "high"
	Critical UrgencyLevel = "critical"
)

// UrgencyStatus represents the urgency status
type UrgencyStatus string

const (
	Open       UrgencyStatus = "open"
	InProgress UrgencyStatus = "in_progress"
	Resolved   UrgencyStatus = "resolved"
	Closed     UrgencyStatus = "closed"
)

// UrgencyCreateRequest DTO for creating a new urgency
// swagger:model
type UrgencyCreateRequest struct {
	FirstName    string       `json:"firstName" binding:"required"`
	LastName     string       `json:"lastName" binding:"required"`
	Email        string       `json:"email"`
	ContactPhone string       `json:"contactPhone" binding:"required"`
	Location     string       `json:"location" binding:"required"`
	Description  string       `json:"description" binding:"required"`
	Level        UrgencyLevel `json:"level"`
}

// UrgencyUpdateRequest DTO for updating an urgency
// swagger:model
type UrgencyUpdateRequest struct {
	FirstName    string        `json:"firstName"`
	LastName     string        `json:"lastName"`
	Email        string        `json:"email"`
	ContactPhone string        `json:"contactPhone"`
	Location     string        `json:"location"`
	Description  string        `json:"description"`
	Level        UrgencyLevel  `json:"level"`
	Status       UrgencyStatus `json:"status"`
}

// UrgencyResponse DTO for returning an urgency
// swagger:model
type UrgencyResponse struct {
	ID                 uint          `json:"id"`
	FirstName          string        `json:"firstName"`
	LastName           string        `json:"lastName"`
	Email              string        `json:"email"`
	ContactPhone       string        `json:"contactPhone"`
	Location           string        `json:"location"`
	Description        string        `json:"description"`
	Level              UrgencyLevel  `json:"level"`
	Status             UrgencyStatus `json:"status"`
	AssignedEmployeeId *uint         `json:"assignedEmployeeId,omitempty"`
	AssignedAt         string        `json:"assignedAt,omitempty"`
	CreatedAt          string        `json:"createdAt"`
	UpdatedAt          string        `json:"updatedAt"`
}

// UrgencyList DTO for returning a list of urgencies
// swagger:model
type UrgencyList struct {
	Urgencies []UrgencyResponse `json:"urgencies"`
}

// UrgencyListRequest DTO for listing urgencies with pagination
// swagger:model
type UrgencyListRequest struct {
	Page     int `json:"page,omitempty" form:"page"`
	PageSize int `json:"pageSize,omitempty" form:"pageSize"`
}

// UrgencyListResponse DTO for returning paginated urgencies
// swagger:model
type UrgencyListResponse struct {
	Urgencies  []UrgencyResponse `json:"urgencies"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"pageSize"`
	TotalPages int               `json:"totalPages"`
}

// AssignmentCreateRequest DTO for direct assignment to an urgency
// swagger:model
type AssignmentCreateRequest struct {
	EmployeeID uint `json:"employeeId" binding:"required"`
}

// NotificationResponse DTO for returning notification data
// swagger:model
type NotificationResponse struct {
	ID               uint   `json:"id"`
	UrgencyID        uint   `json:"urgencyId"`
	EmployeeID       uint   `json:"employeeId"`
	NotificationType string `json:"notificationType"`
	Recipient        string `json:"recipient"`
	Message          string `json:"message"`
	Status           string `json:"status"`
	Attempts         int    `json:"attempts"`
	LastAttemptAt    string `json:"lastAttemptAt,omitempty"`
	SentAt           string `json:"sentAt,omitempty"`
	ErrorMessage     string `json:"errorMessage,omitempty"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

// AssignmentResponse DTO for returning minimal assignment info
// swagger:model
type AssignmentResponse struct {
	UrgencyID        uint   `json:"urgencyId"`
	AssignedEmployee uint   `json:"assignedEmployeeId"`
	AssignedAt       string `json:"assignedAt"`
}

// Helper methods

func (l UrgencyLevel) Valid() bool {
	for _, v := range []UrgencyLevel{Low, Medium, High, Critical} {
		if l == v {
			return true
		}
	}
	return false
}

func (s UrgencyStatus) Valid() bool {
	for _, v := range []UrgencyStatus{Open, InProgress, Resolved, Closed} {
		if s == v {
			return true
		}
	}
	return false
}

func (r *UrgencyCreateRequest) Validate() error {
	// Validate enum types first
	if r.Level != "" && !r.Level.Valid() {
		return fmt.Errorf("invalid urgency level")
	}

	// Validate other fields
	var errors validation.ValidationErrors

	if err := utils.ValidateRequiredField(r.FirstName, "first name"); err != nil {
		errors.AddError("firstName", err)
	}
	if err := utils.ValidateRequiredField(r.LastName, "last name"); err != nil {
		errors.AddError("lastName", err)
	}
	if err := utils.ValidateRequiredField(r.Description, "description"); err != nil {
		errors.AddError("description", err)
	}

	if err := utils.ValidateOptionalEmail(r.Email); err != nil {
		errors.Add("email", "invalid email format")
	}

	if err := utils.ValidatePhone(r.ContactPhone); err != nil {
		errors.AddError("contactPhone", err)
	}

	if err := utils.ValidateCoordinates(r.Location); err != nil {
		errors.AddError("location", err)
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

func (r *UrgencyUpdateRequest) Validate() error {
	// Validate enum types first
	if r.Level != "" && !r.Level.Valid() {
		return fmt.Errorf("invalid urgency level")
	}
	if r.Status != "" && !r.Status.Valid() {
		return fmt.Errorf("invalid status")
	}

	// Validate other fields
	var errors validation.ValidationErrors

	if err := utils.ValidateOptionalEmail(r.Email); err != nil {
		errors.Add("email", "invalid email format")
	}

	if r.ContactPhone != "" {
		if err := utils.ValidatePhone(r.ContactPhone); err != nil {
			errors.AddError("contactPhone", err)
		}
	}

	if r.Location != "" {
		if err := utils.ValidateCoordinates(r.Location); err != nil {
			errors.AddError("location", err)
		}
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}
