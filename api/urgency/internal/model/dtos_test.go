package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrgencyCreateRequest_Validate(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		req := &UrgencyCreateRequest{
			Name:         "Test Urgency",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Description:  "Test description",
			Level:        "High",
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("missing name", func(t *testing.T) {
		req := &UrgencyCreateRequest{
			Name:         "",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Description:  "Test description",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("missing email", func(t *testing.T) {
		req := &UrgencyCreateRequest{
			Name:         "Test Urgency",
			Email:        "",
			ContactPhone: "123456789",
			Description:  "Test description",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email is required")
	})

	t.Run("invalid email", func(t *testing.T) {
		req := &UrgencyCreateRequest{
			Name:         "Test Urgency",
			Email:        "invalid-email",
			ContactPhone: "123456789",
			Description:  "Test description",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email format")
	})

	t.Run("missing contact phone", func(t *testing.T) {
		req := &UrgencyCreateRequest{
			Name:         "Test Urgency",
			Email:        "test@example.com",
			ContactPhone: "",
			Description:  "Test description",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "contact phone is required")
	})

	t.Run("missing description", func(t *testing.T) {
		req := &UrgencyCreateRequest{
			Name:         "Test Urgency",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Description:  "",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "description is required")
	})

	t.Run("invalid urgency level", func(t *testing.T) {
		req := &UrgencyCreateRequest{
			Name:         "Test Urgency",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Description:  "Test description",
			Level:        "Invalid",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid urgency level")
	})
}

func TestUrgencyUpdateRequest_Validate(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		req := &UrgencyUpdateRequest{
			Name:         "Updated Urgency",
			Email:        "updated@example.com",
			ContactPhone: "987654321",
			Description:  "Updated description",
			Level:        "Critical",
			Status:       "In Progress",
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid email", func(t *testing.T) {
		req := &UrgencyUpdateRequest{
			Email: "invalid-email",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email format")
	})

	t.Run("invalid urgency level", func(t *testing.T) {
		req := &UrgencyUpdateRequest{
			Level: "Invalid",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid urgency level")
	})

	t.Run("invalid status", func(t *testing.T) {
		req := &UrgencyUpdateRequest{
			Status: "Invalid Status",
		}

		err := req.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid status")
	})
}

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		status   string
		expected bool
	}{
		{"Open", true},
		{"In Progress", true},
		{"Resolved", true},
		{"Closed", true},
		{"Invalid", false},
		{"", false},
	}

	for _, test := range tests {
		t.Run(test.status, func(t *testing.T) {
			result := isValidStatus(test.status)
			assert.Equal(t, test.expected, result)
		})
	}
}
