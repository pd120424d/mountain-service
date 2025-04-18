package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapUpdateRequestToEmployee(t *testing.T) {
	t.Run("it updates allowed fields from the request into the employee struct", func(t *testing.T) {
		req := &EmployeeUpdateRequest{
			FirstName:      "Bruce",
			LastName:       "Lee",
			Email:          "test-user@example.com",
			Gender:         "M",
			Phone:          "123456789",
			ProfilePicture: "https://example.com/profile.jpg",
			ProfileType:    "Medic",
		}
		existing := &Employee{
			FirstName:      "Alice",
			LastName:       "Smith",
			Email:          "asmith@example.com",
			Gender:         "F",
			Phone:          "987654321",
			ProfilePicture: "https://example.com/old-profile.jpg",
			ProfileType:    Technical,
		}

		MapUpdateRequestToEmployee(req, existing)

		assert.Equal(t, "Bruce", existing.FirstName)
		assert.Equal(t, "Lee", existing.LastName)
		assert.Equal(t, "test-user@example.com", existing.Email)
		assert.Equal(t, "M", existing.Gender)
		assert.Equal(t, "123456789", existing.Phone)
		assert.Equal(t, "https://example.com/profile.jpg", existing.ProfilePicture)
		assert.Equal(t, Medic, existing.ProfileType)
	})
}
