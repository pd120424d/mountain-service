package model

import (
	"testing"
	"time"

	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/stretchr/testify/assert"
)

func TestActivity_TableName(t *testing.T) {
	t.Parallel()

	activity := Activity{}
	assert.Equal(t, "activities", activity.TableName())
}

func TestActivity_ToResponse(t *testing.T) {
	t.Parallel()

	t.Run("converts activity with all fields", func(t *testing.T) {
		actorID := uint(1)
		targetID := uint(2)
		createdAt := time.Now()
		updatedAt := time.Now()

		activity := &Activity{
			ID:          123,
			Type:        activityV1.ActivityEmployeeCreated,
			Level:       activityV1.ActivityLevelInfo,
			Title:       "Test Title",
			Description: "Test Description",
			ActorID:     &actorID,
			ActorName:   "Test Actor",
			TargetID:    &targetID,
			TargetType:  "employee",
			Metadata:    "{}",
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		response := activity.ToResponse()

		assert.Equal(t, uint(123), response.ID)
		assert.Equal(t, activityV1.ActivityEmployeeCreated, response.Type)
		assert.Equal(t, activityV1.ActivityLevelInfo, response.Level)
		assert.Equal(t, "Test Title", response.Title)
		assert.Equal(t, "Test Description", response.Description)
		assert.Equal(t, &actorID, response.ActorID)
		assert.Equal(t, "Test Actor", response.ActorName)
		assert.Equal(t, &targetID, response.TargetID)
		assert.Equal(t, "employee", response.TargetType)
		assert.Equal(t, "{}", response.Metadata)
		assert.Equal(t, createdAt.Format(time.RFC3339), response.CreatedAt)
		assert.Equal(t, updatedAt.Format(time.RFC3339), response.UpdatedAt)
	})

	t.Run("converts activity with nil ActorID and TargetID", func(t *testing.T) {
		activity := &Activity{
			ID:          123,
			Type:        activityV1.ActivitySystemReset,
			Level:       activityV1.ActivityLevelCritical,
			Title:       "System Reset",
			Description: "System was reset",
			ActorID:     nil,
			ActorName:   "system",
			TargetID:    nil,
			TargetType:  "system",
			Metadata:    "{}",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		response := activity.ToResponse()

		assert.Equal(t, uint(123), response.ID)
		assert.Nil(t, response.ActorID)
		assert.Nil(t, response.TargetID)
		assert.Equal(t, "system", response.ActorName)
		assert.Equal(t, "system", response.TargetType)
	})
}

func TestFromCreateRequest(t *testing.T) {
	t.Parallel()

	t.Run("creates activity from request with all fields", func(t *testing.T) {
		actorID := uint(1)
		targetID := uint(2)

		req := &activityV1.ActivityCreateRequest{
			Type:        activityV1.ActivityEmployeeCreated,
			Level:       activityV1.ActivityLevelInfo,
			Title:       "Test Title",
			Description: "Test Description",
			ActorID:     &actorID,
			ActorName:   "Test Actor",
			TargetID:    &targetID,
			TargetType:  "employee",
			Metadata:    "{}",
		}

		activity := FromCreateRequest(req)

		assert.Equal(t, activityV1.ActivityEmployeeCreated, activity.Type)
		assert.Equal(t, activityV1.ActivityLevelInfo, activity.Level)
		assert.Equal(t, "Test Title", activity.Title)
		assert.Equal(t, "Test Description", activity.Description)
		assert.Equal(t, &actorID, activity.ActorID)
		assert.Equal(t, "Test Actor", activity.ActorName)
		assert.Equal(t, &targetID, activity.TargetID)
		assert.Equal(t, "employee", activity.TargetType)
		assert.Equal(t, "{}", activity.Metadata)
	})

	t.Run("creates activity from request with nil ActorID and TargetID", func(t *testing.T) {
		req := &activityV1.ActivityCreateRequest{
			Type:        activityV1.ActivitySystemReset,
			Level:       activityV1.ActivityLevelCritical,
			Title:       "System Reset",
			Description: "System was reset",
			ActorID:     nil,
			ActorName:   "system",
			TargetID:    nil,
			TargetType:  "system",
			Metadata:    "{}",
		}

		activity := FromCreateRequest(req)

		assert.Equal(t, activityV1.ActivitySystemReset, activity.Type)
		assert.Nil(t, activity.ActorID)
		assert.Nil(t, activity.TargetID)
		assert.Equal(t, "system", activity.ActorName)
		assert.Equal(t, "system", activity.TargetType)
	})
}

func TestNewActivityFilter(t *testing.T) {
	t.Parallel()

	filter := NewActivityFilter()

	assert.Equal(t, 1, filter.Page)
	assert.Equal(t, DefaultPageSize, filter.PageSize)
	assert.Nil(t, filter.Type)
	assert.Nil(t, filter.Level)
	assert.Nil(t, filter.ActorID)
	assert.Nil(t, filter.TargetID)
	assert.Nil(t, filter.TargetType)
	assert.Nil(t, filter.StartDate)
	assert.Nil(t, filter.EndDate)
}

func TestActivityFilter_Validate(t *testing.T) {
	t.Parallel()

	t.Run("corrects invalid page and pageSize", func(t *testing.T) {
		filter := &ActivityFilter{
			Page:     -1,
			PageSize: -1,
		}

		err := filter.Validate()
		assert.NoError(t, err)
		assert.Equal(t, 1, filter.Page)
		assert.Equal(t, DefaultPageSize, filter.PageSize)
	})

	t.Run("corrects pageSize exceeding maximum", func(t *testing.T) {
		filter := &ActivityFilter{
			Page:     1,
			PageSize: MaxPageSize + 100,
		}

		err := filter.Validate()
		assert.NoError(t, err)
		assert.Equal(t, 1, filter.Page)
		assert.Equal(t, MaxPageSize, filter.PageSize)
	})

	t.Run("keeps valid values unchanged", func(t *testing.T) {
		filter := &ActivityFilter{
			Page:     2,
			PageSize: 25,
		}

		err := filter.Validate()
		assert.NoError(t, err)
		assert.Equal(t, 2, filter.Page)
		assert.Equal(t, 25, filter.PageSize)
	})
}

func TestActivityFilter_GetOffset(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{"first page", 1, 10, 0},
		{"second page", 2, 10, 10},
		{"third page", 3, 25, 50},
		{"large page", 10, 100, 900},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := &ActivityFilter{
				Page:     tt.page,
				PageSize: tt.pageSize,
			}

			offset := filter.GetOffset()
			assert.Equal(t, tt.expected, offset)
		})
	}
}

func TestActivityFilter_GetLimit(t *testing.T) {
	t.Parallel()

	filter := &ActivityFilter{
		Page:     2,
		PageSize: 25,
	}

	limit := filter.GetLimit()
	assert.Equal(t, 25, limit)
}

func TestActivityStats_ToResponse(t *testing.T) {
	t.Parallel()

	stats := &ActivityStats{
		TotalActivities: 100,
		ActivitiesByType: map[activityV1.ActivityType]int64{
			activityV1.ActivityEmployeeCreated: 50,
			activityV1.ActivityUrgencyCreated:  30,
		},
		ActivitiesByLevel: map[activityV1.ActivityLevel]int64{
			activityV1.ActivityLevelInfo:    70,
			activityV1.ActivityLevelWarning: 30,
		},
		RecentActivities: []Activity{
			{
				ID:          1,
				Type:        activityV1.ActivityEmployeeCreated,
				Level:       activityV1.ActivityLevelInfo,
				Title:       "Recent Activity",
				Description: "Recent activity description",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
		ActivitiesLast24h:    10,
		ActivitiesLast7Days:  50,
		ActivitiesLast30Days: 100,
	}

	response := stats.ToResponse()

	assert.Equal(t, int64(100), response.TotalActivities)
	assert.Equal(t, int64(50), response.ActivitiesByType[activityV1.ActivityEmployeeCreated])
	assert.Equal(t, int64(30), response.ActivitiesByType[activityV1.ActivityUrgencyCreated])
	assert.Equal(t, int64(70), response.ActivitiesByLevel[activityV1.ActivityLevelInfo])
	assert.Equal(t, int64(30), response.ActivitiesByLevel[activityV1.ActivityLevelWarning])
	assert.Len(t, response.RecentActivities, 1)
	assert.Equal(t, uint(1), response.RecentActivities[0].ID)
	assert.Equal(t, int64(10), response.ActivitiesLast24h)
	assert.Equal(t, int64(50), response.ActivitiesLast7Days)
	assert.Equal(t, int64(100), response.ActivitiesLast30Days)
}

func TestNewEmployeeActivity(t *testing.T) {
	t.Parallel()

	activity := NewEmployeeActivity(
		activityV1.ActivityEmployeeCreated,
		123,
		"John Doe",
		"Employee Created",
		"New employee was created",
	)

	assert.Equal(t, activityV1.ActivityEmployeeCreated, activity.Type)
	assert.Equal(t, activityV1.ActivityLevelInfo, activity.Level)
	assert.Equal(t, "Employee Created", activity.Title)
	assert.Equal(t, "New employee was created", activity.Description)
	assert.Equal(t, uint(123), *activity.TargetID)
	assert.Equal(t, "employee", activity.TargetType)
	assert.Equal(t, "John Doe", activity.ActorName)
}

func TestNewUrgencyActivity(t *testing.T) {
	t.Parallel()

	activity := NewUrgencyActivity(
		activityV1.ActivityUrgencyCreated,
		456,
		"Admin User",
		"Urgency Created",
		"New urgency was created",
	)

	assert.Equal(t, activityV1.ActivityUrgencyCreated, activity.Type)
	assert.Equal(t, activityV1.ActivityLevelWarning, activity.Level)
	assert.Equal(t, "Urgency Created", activity.Title)
	assert.Equal(t, "New urgency was created", activity.Description)
	assert.Equal(t, uint(456), *activity.TargetID)
	assert.Equal(t, "urgency", activity.TargetType)
	assert.Equal(t, "Admin User", activity.ActorName)
}

func TestNewSystemActivity(t *testing.T) {
	t.Parallel()

	activity := NewSystemActivity(
		activityV1.ActivitySystemReset,
		activityV1.ActivityLevelCritical,
		"System Reset",
		"System was reset",
	)

	assert.Equal(t, activityV1.ActivitySystemReset, activity.Type)
	assert.Equal(t, activityV1.ActivityLevelCritical, activity.Level)
	assert.Equal(t, "System Reset", activity.Title)
	assert.Equal(t, "System was reset", activity.Description)
	assert.Nil(t, activity.TargetID)
	assert.Equal(t, "system", activity.TargetType)
	assert.Equal(t, "system", activity.ActorName)
}

func TestNewNotificationActivity(t *testing.T) {
	t.Parallel()

	activity := NewNotificationActivity(
		activityV1.ActivityNotificationSent,
		activityV1.ActivityLevelInfo,
		123,
		456,
		"Notification Sent",
		"Notification was sent to employee",
	)

	assert.Equal(t, activityV1.ActivityNotificationSent, activity.Type)
	assert.Equal(t, activityV1.ActivityLevelInfo, activity.Level)
	assert.Equal(t, "Notification Sent", activity.Title)
	assert.Equal(t, "Notification was sent to employee", activity.Description)
	assert.Equal(t, uint(456), *activity.TargetID)
	assert.Equal(t, "notification", activity.TargetType)
	assert.Equal(t, uint(123), *activity.ActorID)
	assert.Equal(t, "system", activity.ActorName)
}

func TestActivityFilter_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("it handles negative values correctly", func(t *testing.T) {
		filter := &ActivityFilter{
			Page:     -10,
			PageSize: -5,
		}

		err := filter.Validate()
		assert.NoError(t, err)
		assert.Equal(t, 1, filter.Page)
		assert.Equal(t, 50, filter.PageSize)
	})

	t.Run("it handles zero values correctly", func(t *testing.T) {
		filter := &ActivityFilter{
			Page:     0,
			PageSize: 0,
		}

		err := filter.Validate()
		assert.NoError(t, err)
		assert.Equal(t, 1, filter.Page)
		assert.Equal(t, 50, filter.PageSize)
	})

	t.Run("it handles extremely large page size", func(t *testing.T) {
		filter := &ActivityFilter{
			Page:     1,
			PageSize: 10000,
		}

		err := filter.Validate()
		assert.NoError(t, err)
		assert.Equal(t, 1, filter.Page)
		assert.Equal(t, 1000, filter.PageSize) // Should be capped at maximum (MaxPageSize = 1000)
	})

	t.Run("it calculates offset correctly for large pages", func(t *testing.T) {
		filter := &ActivityFilter{
			Page:     1000,
			PageSize: 50,
		}

		filter.Validate()
		offset := filter.GetOffset()
		assert.Equal(t, 49950, offset) // (1000-1) * 50
	})

	t.Run("it calculates limit correctly", func(t *testing.T) {
		filter := &ActivityFilter{
			Page:     1,
			PageSize: 25,
		}

		filter.Validate()
		limit := filter.GetLimit()
		assert.Equal(t, 25, limit)
	})
}

func TestActivityStats_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("it converts stats with zero values", func(t *testing.T) {
		stats := &ActivityStats{
			TotalActivities:      0,
			ActivitiesByType:     make(map[activityV1.ActivityType]int64),
			ActivitiesByLevel:    make(map[activityV1.ActivityLevel]int64),
			RecentActivities:     []Activity{},
			ActivitiesLast24h:    0,
			ActivitiesLast7Days:  0,
			ActivitiesLast30Days: 0,
		}

		response := stats.ToResponse()
		assert.Equal(t, int64(0), response.TotalActivities)
		assert.Equal(t, int64(0), response.ActivitiesLast24h)
		assert.Equal(t, int64(0), response.ActivitiesLast7Days)
		assert.Equal(t, int64(0), response.ActivitiesLast30Days)
		assert.Empty(t, response.RecentActivities)
	})

	t.Run("it converts stats with large values", func(t *testing.T) {
		stats := &ActivityStats{
			TotalActivities: 999999,
			ActivitiesByType: map[activityV1.ActivityType]int64{
				activityV1.ActivityEmployeeCreated: 100000,
				activityV1.ActivityUrgencyCreated:  200000,
			},
			ActivitiesByLevel: map[activityV1.ActivityLevel]int64{
				activityV1.ActivityLevelInfo:    300000,
				activityV1.ActivityLevelWarning: 400000,
			},
			RecentActivities:     []Activity{},
			ActivitiesLast24h:    50000,
			ActivitiesLast7Days:  200000,
			ActivitiesLast30Days: 500000,
		}

		response := stats.ToResponse()
		assert.Equal(t, int64(999999), response.TotalActivities)
		assert.Equal(t, int64(50000), response.ActivitiesLast24h)
		assert.Equal(t, int64(200000), response.ActivitiesLast7Days)
		assert.Equal(t, int64(500000), response.ActivitiesLast30Days)
		assert.NotNil(t, response.ActivitiesByType)
		assert.NotNil(t, response.ActivitiesByLevel)
	})

	t.Run("it handles nil maps correctly", func(t *testing.T) {
		stats := &ActivityStats{
			TotalActivities:      100,
			ActivitiesByType:     nil,
			ActivitiesByLevel:    nil,
			RecentActivities:     nil,
			ActivitiesLast24h:    10,
			ActivitiesLast7Days:  50,
			ActivitiesLast30Days: 100,
		}

		response := stats.ToResponse()
		assert.Equal(t, int64(100), response.TotalActivities)
		// The ToResponse method directly assigns nil maps, so they will be nil in response
		assert.Nil(t, response.ActivitiesByType)
		assert.Nil(t, response.ActivitiesByLevel)
		assert.NotNil(t, response.RecentActivities) // This is created as empty slice
	})
}
