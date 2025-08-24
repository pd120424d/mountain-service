package v1

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/shared/validation"
)

// ActivityType represents the type of activity
type ActivityType string

// TODO: cleanup this since it is not used
const (
	// Employee activities
	ActivityEmployeeCreated ActivityType = "employee_created"
	ActivityEmployeeUpdated ActivityType = "employee_updated"
	ActivityEmployeeDeleted ActivityType = "employee_deleted"
	ActivityEmployeeLogin   ActivityType = "employee_login"

	// Shift activities
	ActivityShiftAssigned ActivityType = "shift_assigned"
	ActivityShiftRemoved  ActivityType = "shift_removed"

	// Emergency activities
	ActivityUrgencyCreated     ActivityType = "urgency_created"
	ActivityUrgencyUpdated     ActivityType = "urgency_updated"
	ActivityUrgencyDeleted     ActivityType = "urgency_deleted"
	ActivityEmergencyAssigned  ActivityType = "emergency_assigned"
	ActivityEmergencyAccepted  ActivityType = "emergency_accepted"
	ActivityEmergencyDeclined  ActivityType = "emergency_declined"
	ActivityNotificationSent   ActivityType = "notification_sent"
	ActivityNotificationFailed ActivityType = "notification_failed"

	// System activities
	ActivitySystemReset ActivityType = "system_reset"
)

// ActivityLevel represents the importance level of an activity
type ActivityLevel string

const (
	ActivityLevelInfo     ActivityLevel = "info"
	ActivityLevelWarning  ActivityLevel = "warning"
	ActivityLevelError    ActivityLevel = "error"
	ActivityLevelCritical ActivityLevel = "critical"
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
	ID          uint          `json:"id"`
	Type        ActivityType  `json:"type"`
	Level       ActivityLevel `json:"level"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	ActorID     *uint         `json:"actorId,omitempty"`    // ID of the user who performed the action
	ActorName   string        `json:"actorName,omitempty"`  // Name of the user who performed the action
	TargetID    *uint         `json:"targetId,omitempty"`   // ID of the target entity
	TargetType  string        `json:"targetType,omitempty"` // Type of the target entity (employee, urgency, etc.)
	Metadata    string        `json:"metadata,omitempty"`   // JSON string with additional data
	CreatedAt   string        `json:"createdAt"`
	UpdatedAt   string        `json:"updatedAt"`
}

// ActivityCreateRequest DTO for creating a new activity
// swagger:model
type ActivityCreateRequest struct {
	Type        ActivityType  `json:"type" binding:"required"`
	Level       ActivityLevel `json:"level" binding:"required"`
	Title       string        `json:"title" binding:"required"`
	Description string        `json:"description" binding:"required"`
	ActorID     *uint         `json:"actorId,omitempty"`
	ActorName   string        `json:"actorName,omitempty"`
	TargetID    *uint         `json:"targetId,omitempty"`
	TargetType  string        `json:"targetType,omitempty"`
	Metadata    string        `json:"metadata,omitempty"`
}

// ActivityListRequest DTO for listing activities with filters
// swagger:model
type ActivityListRequest struct {
	Type       ActivityType  `json:"type,omitempty" form:"type"`
	Level      ActivityLevel `json:"level,omitempty" form:"level"`
	ActorID    *uint         `json:"actorId,omitempty" form:"actorId"`
	TargetID   *uint         `json:"targetId,omitempty" form:"targetId"`
	TargetType string        `json:"targetType,omitempty" form:"targetType"`
	StartDate  string        `json:"startDate,omitempty" form:"startDate"` // RFC3339 format
	EndDate    string        `json:"endDate,omitempty" form:"endDate"`     // RFC3339 format
	Page       int           `json:"page,omitempty" form:"page"`
	PageSize   int           `json:"pageSize,omitempty" form:"pageSize"`
}

// ActivityListResponse DTO for returning paginated activities
// swagger:model
type ActivityListResponse struct {
	Activities []ActivityResponse `json:"activities"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"pageSize"`
	TotalPages int                `json:"totalPages"`
}

// ActivityStatsResponse DTO for returning activity statistics
// swagger:model
type ActivityStatsResponse struct {
	TotalActivities      int64                   `json:"totalActivities"`
	ActivitiesByType     map[ActivityType]int64  `json:"activitiesByType"`
	ActivitiesByLevel    map[ActivityLevel]int64 `json:"activitiesByLevel"`
	RecentActivities     []ActivityResponse      `json:"recentActivities"`
	ActivitiesLast24h    int64                   `json:"activitiesLast24h"`
	ActivitiesLast7Days  int64                   `json:"activitiesLast7Days"`
	ActivitiesLast30Days int64                   `json:"activitiesLast30Days"`
}

// Helper methods

func (r *ActivityCreateRequest) Validate() error {
	if !r.Type.Valid() {
		return fmt.Errorf("invalid activity type: %s", r.Type)
	}
	if !r.Level.Valid() {
		return fmt.Errorf("invalid activity level: %s", r.Level)
	}

	var errors validation.ValidationErrors

	if err := utils.ValidateRequiredField(r.Title, "title"); err != nil {
		errors.AddError("title", err)
	}
	if err := utils.ValidateRequiredField(r.Description, "description"); err != nil {
		errors.AddError("description", err)
	}

	if r.Metadata != "" {
		var temp interface{}
		if err := json.Unmarshal([]byte(r.Metadata), &temp); err != nil {
			errors.Add("metadata", "metadata must be valid JSON")
		}
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
		"ActivityCreateRequest { Type: %s, Level: %s, Title: %s, Description: %s, ActorID: %v, ActorName: %s, TargetID: %v, TargetType: %s }",
		r.Type,
		r.Level,
		r.Title,
		description,
		r.ActorID,
		r.ActorName,
		r.TargetID,
		r.TargetType,
	)
}

func ActivityTypeFromString(s string) ActivityType {
	return ActivityType(s)
}

func ActivityLevelFromString(s string) ActivityLevel {
	return ActivityLevel(s)
}

func (t ActivityType) Valid() bool {
	validTypes := []ActivityType{
		ActivityEmployeeCreated, ActivityEmployeeUpdated, ActivityEmployeeDeleted, ActivityEmployeeLogin,
		ActivityShiftAssigned, ActivityShiftRemoved,
		ActivityUrgencyCreated, ActivityUrgencyUpdated, ActivityUrgencyDeleted,
		ActivityEmergencyAssigned, ActivityEmergencyAccepted, ActivityEmergencyDeclined,
		ActivityNotificationSent, ActivityNotificationFailed,
		ActivitySystemReset,
	}
	for _, validType := range validTypes {
		if t == validType {
			return true
		}
	}
	return false
}

func (l ActivityLevel) Valid() bool {
	validLevels := []ActivityLevel{
		ActivityLevelInfo, ActivityLevelWarning, ActivityLevelError, ActivityLevelCritical,
	}
	for _, validLevel := range validLevels {
		if l == validLevel {
			return true
		}
	}
	return false
}

func (at ActivityType) String() string {
	return string(at)
}

func (al ActivityLevel) String() string {
	return string(al)
}

// CQRS-related contracts

// OutboxEvent represents an event to be published to Pub/Sub for CQRS
type OutboxEvent struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	EventType   string     `json:"event_type" gorm:"not null"`
	AggregateID string     `json:"aggregate_id" gorm:"not null"`
	EventData   string     `json:"event_data" gorm:"type:text"`
	Published   bool       `json:"published" gorm:"default:false"`
	CreatedAt   time.Time  `json:"created_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

// ActivityEvent represents the event data for activity operations (CQRS)
type ActivityEvent struct {
	Type        string    `json:"type"` // CREATE, UPDATE, DELETE
	ActivityID  uint      `json:"activity_id"`
	UrgencyID   uint      `json:"urgency_id"`
	EmployeeID  uint      `json:"employee_id"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	// Denormalized data for read model
	EmployeeName string `json:"employee_name,omitempty"`
	UrgencyTitle string `json:"urgency_title,omitempty"`
	UrgencyLevel string `json:"urgency_level,omitempty"`
}

// ActivityEventType represents the type of CQRS event
type ActivityEventType string

const (
	ActivityEventCreated ActivityEventType = "activity.created"
	ActivityEventUpdated ActivityEventType = "activity.updated"
	ActivityEventDeleted ActivityEventType = "activity.deleted"
)

func (aet ActivityEventType) String() string {
	return string(aet)
}

func CreateOutboxEvent(eventType ActivityEventType, activityID uint, eventData ActivityEvent) *OutboxEvent {
	data, _ := json.Marshal(eventData)

	return &OutboxEvent{
		EventType:   eventType.String(),
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

// Simplified Activity model for CQRS
type SimpleActivity struct {
	ID          uint      `json:"id"`
	UrgencyID   uint      `json:"urgency_id"`
	EmployeeID  uint      `json:"employee_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type SimpleActivityCreateRequest struct {
	UrgencyID   uint   `json:"urgency_id" validate:"required"`
	EmployeeID  uint   `json:"employee_id" validate:"required"`
	Description string `json:"description" validate:"required,min=10"`
}

func (r *SimpleActivityCreateRequest) ToSimpleActivity() *SimpleActivity {
	return &SimpleActivity{
		UrgencyID:   r.UrgencyID,
		EmployeeID:  r.EmployeeID,
		Description: r.Description,
		CreatedAt:   time.Now(),
	}
}
