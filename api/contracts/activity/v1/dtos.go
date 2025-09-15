package v1

import (
	"encoding/json"
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

// ActivityResponse DTO for returning activity data
// swagger:model
type ActivityResponse struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	EmployeeID  uint   `json:"employeeId"` // ID of the employee who created the activity
	UrgencyID   uint   `json:"urgencyId"`  // ID of the urgency this activity relates to
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// ActivityCreateRequest DTO for creating a new activity
// swagger:model
type ActivityCreateRequest struct {
	Description string `json:"description" binding:"required"`
	EmployeeID  uint   `json:"employeeId" binding:"required"`
	UrgencyID   uint   `json:"urgencyId" binding:"required"`
}

// ActivityListRequest DTO for listing activities with filters
// swagger:model
type ActivityListRequest struct {
	EmployeeID *uint  `json:"employeeId,omitempty" form:"employeeId"`
	UrgencyID  *uint  `json:"urgencyId,omitempty" form:"urgencyId"`
	StartDate  string `json:"startDate,omitempty" form:"startDate"` // RFC3339 format
	EndDate    string `json:"endDate,omitempty" form:"endDate"`     // RFC3339 format
	Page       int    `json:"page,omitempty" form:"page"`
	PageSize   int    `json:"pageSize,omitempty" form:"pageSize"`
	PageToken  string `json:"pageToken,omitempty" form:"pageToken"`
}

// ActivityListResponse DTO for returning paginated activities
// swagger:model
type ActivityListResponse struct {
	Activities    []ActivityResponse `json:"activities"`
	Total         int64              `json:"total"`
	Page          int                `json:"page"`
	PageSize      int                `json:"pageSize"`
	TotalPages    int                `json:"totalPages"`
	NextPageToken string             `json:"nextPageToken,omitempty"`
}

// ActivityCountsRequest DTO for requesting counts by urgency IDs
// swagger:model
type ActivityCountsRequest struct {
	UrgencyIDs []uint `json:"urgencyIds" binding:"required"`
}

func (r *ActivityCountsRequest) Validate() error {
	if len(r.UrgencyIDs) == 0 {
		return fmt.Errorf("urgencyIds must not be empty")
	}
	if len(r.UrgencyIDs) > 100 {
		return fmt.Errorf("urgencyIds cannot exceed 100 per request")
	}
	for i, id := range r.UrgencyIDs {
		if id == 0 {
			return fmt.Errorf("urgencyIds[%d] must be > 0", i)
		}
	}
	return nil
}

// ActivityCountsResponse DTO for returning counts keyed by urgencyId (as string keys for JSON)
// swagger:model
type ActivityCountsResponse struct {
	Counts map[string]int64 `json:"counts"`
}

// Helper methods

func (r *ActivityCreateRequest) Validate() error {
	var errors validation.ValidationErrors

	if err := utils.ValidateRequiredField(r.Description, "description"); err != nil {
		errors.AddError("description", err)
	}

	if r.EmployeeID == 0 {
		errors.Add("employeeId", "employeeId is required and must be greater than 0")
	}

	if r.UrgencyID == 0 {
		errors.Add("urgencyId", "urgencyId is required and must be greater than 0")
	}

	if errors.HasErrors() {
		return errors
	}
	return nil
}

func (r *ActivityListRequest) Validate() error {
	if r.Page < 0 {
		return fmt.Errorf("page must be non-negative")
	}
	if r.PageSize < 0 {
		return fmt.Errorf("pageSize must be non-negative")
	}
	if r.PageSize > 1000 {
		return fmt.Errorf("pageSize cannot exceed 1000")
	}

	// Validate date formats if provided
	if r.StartDate != "" {
		if _, err := time.Parse(time.RFC3339, r.StartDate); err != nil {
			return fmt.Errorf("invalid startDate format, expected RFC3339: %v", err)
		}
	}
	if r.EndDate != "" {
		if _, err := time.Parse(time.RFC3339, r.EndDate); err != nil {
			return fmt.Errorf("invalid endDate format, expected RFC3339: %v", err)
		}
	}

	return nil
}

func (r *ActivityCreateRequest) ToString() string {
	description := r.Description
	if len(description) > 50 {
		description = description[:50] + "..."
	}
	return fmt.Sprintf(
		"ActivityCreateRequest { Description: %s, EmployeeID: %d, UrgencyID: %d }",
		description,
		r.EmployeeID,
		r.UrgencyID,
	)
}

// CQRS-related contracts

// OutboxEvent represents an event to be published to Pub/Sub for CQRS
type OutboxEvent struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	AggregateID string     `json:"aggregateId" gorm:"not null"`
	EventData   string     `json:"eventData" gorm:"type:text"`
	Published   bool       `json:"published" gorm:"default:false"`
	CreatedAt   time.Time  `json:"createdAt"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
}

// ActivityEvent represents the event data for activity operations (CQRS)
// Note: Type is used by the read-model (Firestore) updater to decide the action (CREATE/UPDATE/DELETE).
type ActivityEvent struct {
	Type        string    `json:"type"` // CREATE, UPDATE, DELETE
	ActivityID  uint      `json:"activityId"`
	UrgencyID   uint      `json:"urgencyId"`
	EmployeeID  uint      `json:"employeeId"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	// Denormalized data for read model
	EmployeeName string `json:"employeeName,omitempty"`
	UrgencyTitle string `json:"urgencyTitle,omitempty"`
	UrgencyLevel string `json:"urgencyLevel,omitempty"`
}

func CreateOutboxEvent(activityID uint, eventData ActivityEvent) *OutboxEvent {
	data, _ := json.Marshal(eventData)

	return &OutboxEvent{
		AggregateID: fmt.Sprintf("activity-%d", activityID),
		EventData:   string(data),
		Published:   false,
		CreatedAt:   time.Now(),
	}
}

// GetEventData unmarshals the event data from an outbox event
func (e *OutboxEvent) GetEventData() (*ActivityEvent, error) {
	var eventData ActivityEvent
	err := json.Unmarshal([]byte(e.EventData), &eventData)
	return &eventData, err
}
