package model

import (
	"testing"
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/stretchr/testify/assert"
)

func TestMapUpdateRequestToEmployee(t *testing.T) {
	t.Run("it updates allowed fields from the request into the employee struct", func(t *testing.T) {
		req := &employeeV1.EmployeeUpdateRequest{
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
		testDate := time.Date(2025, 2, 3, 0, 0, 0, 0, time.UTC)

		// Check that all shifts have full availability (2 medics, 4 technical)
		assert.Equal(t, 2, response.Days[testDate].Shift1.MedicSlotsAvailable)
		assert.Equal(t, 4, response.Days[testDate].Shift1.TechnicalSlotsAvailable)
		assert.Equal(t, 2, response.Days[testDate].Shift2.MedicSlotsAvailable)
		assert.Equal(t, 4, response.Days[testDate].Shift2.TechnicalSlotsAvailable)
		assert.Equal(t, 2, response.Days[testDate].Shift3.MedicSlotsAvailable)
		assert.Equal(t, 4, response.Days[testDate].Shift3.TechnicalSlotsAvailable)
	})
}
