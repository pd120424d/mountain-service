package v1

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActivityCreateRequest_Validate(t *testing.T) {
	t.Parallel()

	t.Run("it returns no error for a valid request", func(t *testing.T) {
		req := &ActivityCreateRequest{
			Description: "Employee was assigned to urgency",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("it returns an error for missing description", func(t *testing.T) {
		req := &ActivityCreateRequest{}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "description is required")
	})

	t.Run("it returns an error for empty description", func(t *testing.T) {
		req := &ActivityCreateRequest{
			Description: "   ",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "description is required")
	})
}

func TestActivityListRequest_Validate(t *testing.T) {
	t.Parallel()

	employeeID := uint(1)
	urgencyID := uint(2)

	t.Run("it returns no error for valid request", func(t *testing.T) {
		req := &ActivityListRequest{
			EmployeeID: &employeeID,
			UrgencyID:  &urgencyID,
			StartDate:  "2023-01-01T00:00:00Z",
			EndDate:    "2023-12-31T23:59:59Z",
			Page:       1,
			PageSize:   50,
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

func TestActivityCreateRequest_ToString(t *testing.T) {
	t.Parallel()

	t.Run("it returns a string representation of the request", func(t *testing.T) {
		req := &ActivityCreateRequest{
			Description: "Employee was assigned to urgency",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		expected := "ActivityCreateRequest { Description: Employee was assigned to urgency, EmployeeID: 1, UrgencyID: 2 }"
		assert.Equal(t, expected, req.ToString())
	})

	t.Run("it truncates the description if it exceeds 50 characters", func(t *testing.T) {
		req := &ActivityCreateRequest{
			Description: "This is a very long description that exceeds the maximum length of 50 characters",
			EmployeeID:  1,
			UrgencyID:   2,
		}

		expected := "ActivityCreateRequest { Description: This is a very long description that exceeds the m..., EmployeeID: 1, UrgencyID: 2 }"
		assert.Equal(t, expected, req.ToString())
	})
}

func TestCreateOutboxEvent(t *testing.T) {
	t.Parallel()

	t.Run("it creates an outbox event with correct data", func(t *testing.T) {
		activityID := uint(1)
		activityEvent := ActivityEvent{
			Type:        "test_type",
			ActivityID:  activityID,
			UrgencyID:   2,
			EmployeeID:  3,
			Description: "test description",
			CreatedAt:   time.Now(),
		}

		outboxEvent := CreateOutboxEvent(ActivityEventCreated, activityID, activityEvent)

		assert.Equal(t, string(ActivityEventCreated), outboxEvent.EventType)
		assert.Equal(t, fmt.Sprintf("activity-%d", activityID), outboxEvent.AggregateID)
		assert.Equal(t, false, outboxEvent.Published)
		assert.NotZero(t, outboxEvent.CreatedAt)
		assert.Empty(t, outboxEvent.PublishedAt)
	})
}

func TestOutboxEvent_GetEventData(t *testing.T) {
	t.Parallel()

	t.Run("it unmarshals the event data correctly", func(t *testing.T) {
		activityID := uint(1)
		activityEvent := ActivityEvent{
			Type:        "test_type",
			ActivityID:  activityID,
			UrgencyID:   2,
			EmployeeID:  3,
			Description: "test description",
			CreatedAt:   time.Now(),
		}

		outboxEvent := CreateOutboxEvent(ActivityEventCreated, activityID, activityEvent)

		eventData, err := outboxEvent.GetEventData()
		assert.NoError(t, err)
		assert.Equal(t, &activityEvent, eventData)
	})
}
