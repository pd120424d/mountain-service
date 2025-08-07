package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActivityType_Valid(t *testing.T) {
	t.Parallel()

	validTypes := []ActivityType{
		ActivityEmployeeCreated, ActivityEmployeeUpdated, ActivityEmployeeDeleted, ActivityEmployeeLogin,
		ActivityShiftAssigned, ActivityShiftRemoved,
		ActivityUrgencyCreated, ActivityUrgencyUpdated, ActivityUrgencyDeleted,
		ActivityEmergencyAssigned, ActivityEmergencyAccepted, ActivityEmergencyDeclined,
		ActivityNotificationSent, ActivityNotificationFailed,
		ActivitySystemReset,
	}

	for _, validType := range validTypes {
		t.Run(string(validType), func(t *testing.T) {
			assert.True(t, validType.Valid())
		})
	}

	t.Run("invalid type", func(t *testing.T) {
		invalidType := ActivityType("invalid_type")
		assert.False(t, invalidType.Valid())
	})

	t.Run("empty type", func(t *testing.T) {
		emptyType := ActivityType("")
		assert.False(t, emptyType.Valid())
	})
}

func TestActivityLevel_Valid(t *testing.T) {
	t.Parallel()

	validLevels := []ActivityLevel{
		ActivityLevelInfo, ActivityLevelWarning, ActivityLevelError, ActivityLevelCritical,
	}

	for _, validLevel := range validLevels {
		t.Run(string(validLevel), func(t *testing.T) {
			assert.True(t, validLevel.Valid())
		})
	}

	t.Run("invalid level", func(t *testing.T) {
		invalidLevel := ActivityLevel("invalid_level")
		assert.False(t, invalidLevel.Valid())
	})

	t.Run("empty level", func(t *testing.T) {
		emptyLevel := ActivityLevel("")
		assert.False(t, emptyLevel.Valid())
	})
}

func TestActivityCreateRequest_Validate(t *testing.T) {
	t.Parallel()

	t.Run("it returns no error for a valid request", func(t *testing.T) {
		req := &ActivityCreateRequest{
			Type:        ActivityEmployeeCreated,
			Level:       ActivityLevelInfo,
			Title:       "Employee Created",
			Description: "A new employee was created",
			ActorName:   "admin",
			TargetType:  "employee",
			Metadata:    `{"employeeId": 123}`,
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("it returns no error for minimal valid request", func(t *testing.T) {
		req := &ActivityCreateRequest{
			Type:        ActivitySystemReset,
			Level:       ActivityLevelWarning,
			Title:       "System Reset",
			Description: "System was reset",
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("it returns an error for invalid activity type", func(t *testing.T) {
		req := &ActivityCreateRequest{
			Type:        ActivityType("invalid_type"),
			Level:       ActivityLevelInfo,
			Title:       "Test Title",
			Description: "Test Description",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid activity type")
	})

	t.Run("it returns an error for invalid activity level", func(t *testing.T) {
		req := &ActivityCreateRequest{
			Type:        ActivityEmployeeCreated,
			Level:       ActivityLevel("invalid_level"),
			Title:       "Test Title",
			Description: "Test Description",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid activity level")
	})

	t.Run("it returns an error for missing title", func(t *testing.T) {
		req := &ActivityCreateRequest{
			Type:        ActivityEmployeeCreated,
			Level:       ActivityLevelInfo,
			Description: "Test Description",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("it returns an error for empty title", func(t *testing.T) {
		req := &ActivityCreateRequest{
			Type:        ActivityEmployeeCreated,
			Level:       ActivityLevelInfo,
			Title:       "   ",
			Description: "Test Description",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("it returns an error for missing description", func(t *testing.T) {
		req := &ActivityCreateRequest{
			Type:  ActivityEmployeeCreated,
			Level: ActivityLevelInfo,
			Title: "Test Title",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "description is required")
	})

	t.Run("it returns an error for empty description", func(t *testing.T) {
		req := &ActivityCreateRequest{
			Type:        ActivityEmployeeCreated,
			Level:       ActivityLevelInfo,
			Title:       "Test Title",
			Description: "   ",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "description is required")
	})

	t.Run("it returns an error for invalid JSON metadata", func(t *testing.T) {
		req := &ActivityCreateRequest{
			Type:        ActivityEmployeeCreated,
			Level:       ActivityLevelInfo,
			Title:       "Test Title",
			Description: "Test Description",
			Metadata:    `{"invalid": json}`,
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "metadata must be valid JSON")
	})

	t.Run("it returns no error for valid JSON metadata", func(t *testing.T) {
		req := &ActivityCreateRequest{
			Type:        ActivityEmployeeCreated,
			Level:       ActivityLevelInfo,
			Title:       "Test Title",
			Description: "Test Description",
			Metadata:    `{"employeeId": 123, "action": "create"}`,
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("it returns no error for empty metadata", func(t *testing.T) {
		req := &ActivityCreateRequest{
			Type:        ActivityEmployeeCreated,
			Level:       ActivityLevelInfo,
			Title:       "Test Title",
			Description: "Test Description",
			Metadata:    "",
		}

		err := req.Validate()
		assert.NoError(t, err)
	})
}

func TestActivityListRequest_Validate(t *testing.T) {
	t.Parallel()

	t.Run("it returns no error for valid request", func(t *testing.T) {
		req := &ActivityListRequest{
			Type:      ActivityEmployeeCreated,
			Level:     ActivityLevelInfo,
			StartDate: "2023-01-01T00:00:00Z",
			EndDate:   "2023-12-31T23:59:59Z",
			Page:      1,
			PageSize:  50,
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("it returns no error for minimal request", func(t *testing.T) {
		req := &ActivityListRequest{}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("it returns an error for negative page", func(t *testing.T) {
		req := &ActivityListRequest{
			Page: -1,
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "page must be non-negative")
	})

	t.Run("it returns an error for negative page size", func(t *testing.T) {
		req := &ActivityListRequest{
			PageSize: -1,
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "pageSize must be non-negative")
	})

	t.Run("it returns an error for page size exceeding limit", func(t *testing.T) {
		req := &ActivityListRequest{
			PageSize: 1001,
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "pageSize cannot exceed 1000")
	})

	t.Run("it returns an error for invalid start date", func(t *testing.T) {
		req := &ActivityListRequest{
			StartDate: "invalid-date",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid startDate format")
	})

	t.Run("it returns an error for invalid end date", func(t *testing.T) {
		req := &ActivityListRequest{
			EndDate: "invalid-date",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid endDate format")
	})
}
