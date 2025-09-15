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
		createdAt := time.Now()
		updatedAt := time.Now()

		activity := &Activity{
			ID:          123,
			Description: "Test Description",
			EmployeeID:  1,
			UrgencyID:   2,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		response := activity.ToResponse()

		assert.Equal(t, uint(123), response.ID)
		assert.Equal(t, "Test Description", response.Description)
		assert.Equal(t, uint(1), response.EmployeeID)
		assert.Equal(t, uint(2), response.UrgencyID)
		assert.Equal(t, createdAt.Format(time.RFC3339), response.CreatedAt)
		assert.Equal(t, updatedAt.Format(time.RFC3339), response.UpdatedAt)
	})

	t.Run("converts activity with different level", func(t *testing.T) {
		activity := &Activity{
			ID:          456,
			Description: "System was reset",
			EmployeeID:  3,
			UrgencyID:   4,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		response := activity.ToResponse()

		assert.Equal(t, uint(456), response.ID)
		assert.Equal(t, "System was reset", response.Description)
		assert.Equal(t, uint(3), response.EmployeeID)
		assert.Equal(t, uint(4), response.UrgencyID)
	})
}

func TestFromCreateRequest(t *testing.T) {
	t.Parallel()

	t.Run("creates activity from request with all fields", func(t *testing.T) {
		req := &activityV1.ActivityCreateRequest{
			Description: "Test Description",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		activity := FromCreateRequest(req)

		assert.Equal(t, "Test Description", activity.Description)
		assert.Equal(t, uint(1), activity.EmployeeID)
		assert.Equal(t, uint(2), activity.UrgencyID)
	})

	t.Run("creates activity from request with different level", func(t *testing.T) {
		req := &activityV1.ActivityCreateRequest{
			Description: "System was reset",
			EmployeeID:  3,
			UrgencyID:   4,
		}

		activity := FromCreateRequest(req)

		assert.Equal(t, "System was reset", activity.Description)
		assert.Equal(t, uint(3), activity.EmployeeID)
		assert.Equal(t, uint(4), activity.UrgencyID)
	})
}

func TestNewActivityFilter(t *testing.T) {
	t.Parallel()

	filter := NewActivityFilter()

	assert.Equal(t, 1, filter.Page)
	assert.Equal(t, DefaultPageSize, filter.PageSize)
	assert.Nil(t, filter.EmployeeID)
	assert.Nil(t, filter.UrgencyID)
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

func TestNewActivity(t *testing.T) {
	t.Run("it creates a new activity with correct values", func(t *testing.T) {
		activity := NewActivity("Test Description", 1, 2)

		assert.Equal(t, "Test Description", activity.Description)
		assert.Equal(t, uint(1), activity.EmployeeID)
		assert.Equal(t, uint(2), activity.UrgencyID)
	})
}
