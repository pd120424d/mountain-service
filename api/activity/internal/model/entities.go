package model

import (
	"time"

	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"gorm.io/gorm"
)

// Activity represents an activity/event in the system
type Activity struct {
	gorm.Model
	ID          uint                      `gorm:"primaryKey"`
	Type        activityV1.ActivityType   `gorm:"type:text;not null;index"`
	Level       activityV1.ActivityLevel  `gorm:"type:text;not null;index"`
	Title       string                    `gorm:"not null"`
	Description string                    `gorm:"type:text;not null"`
	ActorID     *uint                     `gorm:"index"`                    // ID of the user who performed the action
	ActorName   string                    `gorm:"index"`                    // Name of the user who performed the action
	TargetID    *uint                     `gorm:"index"`                    // ID of the target entity
	TargetType  string                    `gorm:"index"`                    // Type of the target entity (employee, urgency, etc.)
	Metadata    string                    `gorm:"type:text"`                // JSON string with additional data
	CreatedAt   time.Time                 `gorm:"index"`
	UpdatedAt   time.Time
}

// TableName returns the table name for the Activity model
func (Activity) TableName() string {
	return "activities"
}

// ToResponse converts the Activity model to ActivityResponse DTO
func (a *Activity) ToResponse() activityV1.ActivityResponse {
	response := activityV1.ActivityResponse{
		ID:          a.ID,
		Type:        a.Type,
		Level:       a.Level,
		Title:       a.Title,
		Description: a.Description,
		ActorName:   a.ActorName,
		TargetType:  a.TargetType,
		Metadata:    a.Metadata,
		CreatedAt:   a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   a.UpdatedAt.Format(time.RFC3339),
	}

	// Only include ActorID if it's not nil
	if a.ActorID != nil {
		response.ActorID = a.ActorID
	}

	// Only include TargetID if it's not nil
	if a.TargetID != nil {
		response.TargetID = a.TargetID
	}

	return response
}

// FromCreateRequest creates an Activity model from ActivityCreateRequest DTO
func FromCreateRequest(req *activityV1.ActivityCreateRequest) *Activity {
	activity := &Activity{
		Type:        req.Type,
		Level:       req.Level,
		Title:       req.Title,
		Description: req.Description,
		ActorName:   req.ActorName,
		TargetType:  req.TargetType,
		Metadata:    req.Metadata,
	}

	// Only set ActorID if provided
	if req.ActorID != nil {
		activity.ActorID = req.ActorID
	}

	// Only set TargetID if provided
	if req.TargetID != nil {
		activity.TargetID = req.TargetID
	}

	return activity
}

// ActivityFilter represents filters for querying activities
type ActivityFilter struct {
	Type       *activityV1.ActivityType
	Level      *activityV1.ActivityLevel
	ActorID    *uint
	TargetID   *uint
	TargetType *string
	StartDate  *time.Time
	EndDate    *time.Time
	Page       int
	PageSize   int
}

// DefaultPageSize is the default number of activities per page
const DefaultPageSize = 50

// MaxPageSize is the maximum number of activities per page
const MaxPageSize = 1000

// NewActivityFilter creates a new ActivityFilter with default values
func NewActivityFilter() *ActivityFilter {
	return &ActivityFilter{
		Page:     1,
		PageSize: DefaultPageSize,
	}
}

// Validate validates the activity filter
func (f *ActivityFilter) Validate() error {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = DefaultPageSize
	}
	if f.PageSize > MaxPageSize {
		f.PageSize = MaxPageSize
	}
	return nil
}

// GetOffset calculates the offset for pagination
func (f *ActivityFilter) GetOffset() int {
	return (f.Page - 1) * f.PageSize
}

// GetLimit returns the page size
func (f *ActivityFilter) GetLimit() int {
	return f.PageSize
}

// ActivityStats represents statistics about activities
type ActivityStats struct {
	TotalActivities      int64
	ActivitiesByType     map[activityV1.ActivityType]int64
	ActivitiesByLevel    map[activityV1.ActivityLevel]int64
	RecentActivities     []Activity
	ActivitiesLast24h    int64
	ActivitiesLast7Days  int64
	ActivitiesLast30Days int64
}

// ToResponse converts ActivityStats to ActivityStatsResponse DTO
func (s *ActivityStats) ToResponse() activityV1.ActivityStatsResponse {
	recentActivities := make([]activityV1.ActivityResponse, len(s.RecentActivities))
	for i, activity := range s.RecentActivities {
		recentActivities[i] = activity.ToResponse()
	}

	return activityV1.ActivityStatsResponse{
		TotalActivities:      s.TotalActivities,
		ActivitiesByType:     s.ActivitiesByType,
		ActivitiesByLevel:    s.ActivitiesByLevel,
		RecentActivities:     recentActivities,
		ActivitiesLast24h:    s.ActivitiesLast24h,
		ActivitiesLast7Days:  s.ActivitiesLast7Days,
		ActivitiesLast30Days: s.ActivitiesLast30Days,
	}
}

// Helper functions for creating common activities

// NewEmployeeActivity creates a new employee-related activity
func NewEmployeeActivity(activityType activityV1.ActivityType, employeeID uint, employeeName, title, description string) *Activity {
	return &Activity{
		Type:        activityType,
		Level:       activityV1.ActivityLevelInfo,
		Title:       title,
		Description: description,
		TargetID:    &employeeID,
		TargetType:  "employee",
		ActorName:   employeeName,
	}
}

// NewUrgencyActivity creates a new urgency-related activity
func NewUrgencyActivity(activityType activityV1.ActivityType, urgencyID uint, actorName, title, description string) *Activity {
	return &Activity{
		Type:        activityType,
		Level:       activityV1.ActivityLevelWarning,
		Title:       title,
		Description: description,
		TargetID:    &urgencyID,
		TargetType:  "urgency",
		ActorName:   actorName,
	}
}

// NewSystemActivity creates a new system-related activity
func NewSystemActivity(activityType activityV1.ActivityType, level activityV1.ActivityLevel, title, description string) *Activity {
	return &Activity{
		Type:        activityType,
		Level:       level,
		Title:       title,
		Description: description,
		TargetType:  "system",
		ActorName:   "system",
	}
}

// NewNotificationActivity creates a new notification-related activity
func NewNotificationActivity(activityType activityV1.ActivityType, level activityV1.ActivityLevel, employeeID, urgencyID uint, title, description string) *Activity {
	return &Activity{
		Type:        activityType,
		Level:       level,
		Title:       title,
		Description: description,
		TargetID:    &urgencyID,
		TargetType:  "notification",
		ActorID:     &employeeID,
		ActorName:   "system",
	}
}
