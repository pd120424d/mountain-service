package model

import (
	"testing"
	"time"

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

func TestMapShiftsAvailabilityToResponse(t *testing.T) {
	t.Run("it maps shifts availability to the response format", func(t *testing.T) {
		availability := &ShiftsAvailabilityRange{
			Days: map[time.Time][]map[ProfileType]int{
				time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC): {
					{Medic: 2, Technical: 4},
					{Medic: 2, Technical: 4},
					{Medic: 2, Technical: 4},
				},
			},
		}

		response := MapShiftsAvailabilityToResponse(availability)

		assert.Equal(t, 1, len(response.Days))
		assert.Equal(t, 2, response.Days[time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)].FirstShift.Medic)
		assert.Equal(t, 4, response.Days[time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)].FirstShift.Technical)
		assert.Equal(t, 2, response.Days[time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)].SecondShift.Medic)
		assert.Equal(t, 4, response.Days[time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)].SecondShift.Technical)
		assert.Equal(t, 2, response.Days[time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)].ThirdShift.Medic)
		assert.Equal(t, 4, response.Days[time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)].ThirdShift.Technical)
	})
}
