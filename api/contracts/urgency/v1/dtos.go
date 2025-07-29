package v1

import (
	"fmt"

	"github.com/pd120424d/mountain-service/api/shared/utils"
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
	ID           uint          `json:"id"`
	FirstName    string        `json:"firstName"`
	LastName     string        `json:"lastName"`
	Email        string        `json:"email"`
	ContactPhone string        `json:"contactPhone"`
	Location     string        `json:"location"`
	Description  string        `json:"description"`
	Level        UrgencyLevel  `json:"level"`
	Status       UrgencyStatus `json:"status"`
	CreatedAt    string        `json:"createdAt"`
	UpdatedAt    string        `json:"updatedAt"`
}

// UrgencyList DTO for returning a list of urgencies
// swagger:model
type UrgencyList struct {
	Urgencies []UrgencyResponse `json:"urgencies"`
}

// EmergencyAssignmentResponse DTO for returning assignment data
// swagger:model
type EmergencyAssignmentResponse struct {
	ID         uint   `json:"id"`
	UrgencyID  uint   `json:"urgencyId"`
	EmployeeID uint   `json:"employeeId"`
	Status     string `json:"status"`
	AssignedAt string `json:"assignedAt"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

// AssignmentAcceptRequest DTO for accepting an assignment
// swagger:model
type AssignmentAcceptRequest struct {
	AssignmentID uint `json:"assignmentId" binding:"required"`
}

// AssignmentDeclineRequest DTO for declining an assignment
// swagger:model
type AssignmentDeclineRequest struct {
	AssignmentID uint   `json:"assignmentId" binding:"required"`
	Reason       string `json:"reason"`
}

// EmployeeAssignmentsResponse DTO for returning employee's assignments
// swagger:model
type EmployeeAssignmentsResponse struct {
	Assignments []EmergencyAssignmentResponse `json:"assignments"`
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
	if err := utils.ValidateRequiredField(r.FirstName, "first name"); err != nil {
		return err
	}
	if err := utils.ValidateRequiredField(r.LastName, "last name"); err != nil {
		return err
	}
	if err := utils.ValidateOptionalEmail(r.Email); err != nil {
		return fmt.Errorf("invalid email format")
	}
	if err := utils.ValidateRequiredField(r.ContactPhone, "contact phone"); err != nil {
		return err
	}
	if err := utils.ValidateRequiredField(r.Location, "location"); err != nil {
		return err
	}
	if err := utils.ValidateRequiredField(r.Description, "description"); err != nil {
		return err
	}
	if r.Level != "" && !r.Level.Valid() {
		return fmt.Errorf("invalid urgency level")
	}
	return nil
}

func (r *UrgencyUpdateRequest) Validate() error {
	if err := utils.ValidateOptionalEmail(r.Email); err != nil {
		return fmt.Errorf("invalid email format")
	}
	if r.Level != "" && !r.Level.Valid() {
		return fmt.Errorf("invalid urgency level")
	}
	if r.Status != "" && !r.Status.Valid() {
		return fmt.Errorf("invalid status")
	}
	return nil
}
