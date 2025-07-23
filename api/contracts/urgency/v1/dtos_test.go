package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrgencyCreateRequest_Validate(t *testing.T) {
	t.Parallel()

	t.Run("it returns no error for a valid request", func(t *testing.T) {
		req := &UrgencyCreateRequest{
			Name:         "Test Urgency",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Location:     "N 43.401123 E 22.662756",
			Description:  "Test description",
			Level:        High,
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("it returns an error for a missing name", func(t *testing.T) {
		req := &UrgencyCreateRequest{
			Name:         "",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Location:     "N 43.401123 E 22.662756",
			Description:  "Test description",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("it returns an error for a missing email", func(t *testing.T) {
		req := &UrgencyCreateRequest{
			Name:         "Test Urgency",
			Email:        "",
			ContactPhone: "123456789",
			Location:     "N 43.401123 E 22.662756",
			Description:  "Test description",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email is required")
	})

	t.Run("it returns an error for an invalid email", func(t *testing.T) {
		req := &UrgencyCreateRequest{
			Name:         "Test Urgency",
			Email:        "invalid-email",
			ContactPhone: "123456789",
			Location:     "N 43.401123 E 22.662756",
			Description:  "Test description",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email format")
	})

	t.Run("it returns an error for a missing contact phone", func(t *testing.T) {
		req := &UrgencyCreateRequest{
			Name:         "Test Urgency",
			Email:        "test@example.com",
			ContactPhone: "",
			Location:     "N 43.401123 E 22.662756",
			Description:  "Test description",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "contact phone is required")
	})

	t.Run("it returns an error for a missing location", func(t *testing.T) {
		req := &UrgencyCreateRequest{
			Name:         "Test Urgency",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Location:     "",
			Description:  "Test description",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "location is required")
	})

	t.Run("it returns an error for a missing description", func(t *testing.T) {
		req := &UrgencyCreateRequest{
			Name:         "Test Urgency",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Location:     "N 43.401123 E 22.662756",
			Description:  "",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "description is required")
	})

	t.Run("it returns an error for an invalid urgency level", func(t *testing.T) {
		req := &UrgencyCreateRequest{
			Name:         "Test Urgency",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Location:     "N 43.401123 E 22.662756",
			Description:  "Test description",
			Level:        "Invalid Level",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid urgency level")
	})
}

func TestUrgencyUpdateRequest_Validate(t *testing.T) {
	t.Parallel()

	t.Run("it returns no error for a valid request", func(t *testing.T) {
		req := &UrgencyUpdateRequest{
			Name:         "Updated Urgency",
			Email:        "updated@example.com",
			ContactPhone: "987654321",
			Location:     "N 44.401123 E 23.662756",
			Description:  "Updated description",
			Level:        Critical,
			Status:       InProgress,
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("it returns an error for an invalid email", func(t *testing.T) {
		req := &UrgencyUpdateRequest{
			Email: "invalid-email",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email format")
	})

	t.Run("it returns an error for an invalid urgency level", func(t *testing.T) {
		req := &UrgencyUpdateRequest{
			Level: "Invalid Level",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid urgency level")
	})

	t.Run("it returns an error for an invalid status", func(t *testing.T) {
		req := &UrgencyUpdateRequest{
			Status: "Invalid",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status")
	})
}

func TestUrgencyLevel_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level    UrgencyLevel
		expected bool
	}{
		{Low, true},
		{Medium, true},
		{High, true},
		{Critical, true},
		{UrgencyLevel("Invalid"), false},
		{UrgencyLevel(""), false},
	}

	for _, test := range tests {
		t.Run(string(test.level), func(t *testing.T) {
			result := test.level.Valid()
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestStatus_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status   Status
		expected bool
	}{
		{Open, true},
		{InProgress, true},
		{Resolved, true},
		{Closed, true},
		{Status("Invalid"), false},
		{Status(""), false},
	}

	for _, test := range tests {
		t.Run(string(test.status), func(t *testing.T) {
			result := test.status.Valid()
			assert.Equal(t, test.expected, result)
		})
	}
}
