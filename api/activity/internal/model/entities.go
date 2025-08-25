package model

import (
	"time"

	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"gorm.io/gorm"
)

// Activity represents an activity/event in the system
type Activity struct {
	gorm.Model
	ID          uint      `gorm:"primaryKey"`
	Description string    `gorm:"type:text;not null"`
	EmployeeID  uint      `gorm:"not null;index"`
	UrgencyID   uint      `gorm:"not null;index"`
	CreatedAt   time.Time `gorm:"index"`
	UpdatedAt   time.Time
}

func (Activity) TableName() string {
	return "activities"
}

// ToResponse converts the Activity model to ActivityResponse DTO
func (a *Activity) ToResponse() activityV1.ActivityResponse {
	return activityV1.ActivityResponse{
		ID:          a.ID,
		Description: a.Description,
		EmployeeID:  a.EmployeeID,
		UrgencyID:   a.UrgencyID,
		CreatedAt:   a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   a.UpdatedAt.Format(time.RFC3339),
	}
}

// FromCreateRequest creates an Activity model from ActivityCreateRequest DTO
func FromCreateRequest(req *activityV1.ActivityCreateRequest) *Activity {
	return &Activity{
		Description: req.Description,
		EmployeeID:  req.EmployeeID,
		UrgencyID:   req.UrgencyID,
	}
}

// ActivityFilter represents filters for querying activities
type ActivityFilter struct {
	EmployeeID *uint
	UrgencyID  *uint
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
		RecentActivities:     recentActivities,
		ActivitiesLast24h:    s.ActivitiesLast24h,
		ActivitiesLast7Days:  s.ActivitiesLast7Days,
		ActivitiesLast30Days: s.ActivitiesLast30Days,
	}
}

func NewActivity(description string, employeeID, urgencyID uint) *Activity {
	return &Activity{
		Description: description,
		EmployeeID:  employeeID,
		UrgencyID:   urgencyID,
	}
}
