package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
)

func TestUrgencyLevelFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected urgencyV1.UrgencyLevel
	}{
		{"Low", urgencyV1.Low},
		{"Medium", urgencyV1.Medium},
		{"High", urgencyV1.High},
		{"Critical", urgencyV1.Critical},
		{"invalid", urgencyV1.Medium},
		{"", urgencyV1.Medium},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := UrgencyLevelFromString(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestUrgency_ToResponse(t *testing.T) {
	t.Parallel()

	t.Run("it converts urgency to response correctly", func(t *testing.T) {
		createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)

		urgency := &Urgency{
			Model: gorm.Model{
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			ID:           1,
			FirstName:    "Marko",
			LastName:     "Markovic",
			Email:        "test@example.com",
			ContactPhone: "+1234567890",
			Location:     "Test Location",
			Description:  "Test Description",
			Level:        urgencyV1.UrgencyLevel(urgencyV1.High),
			Status:       urgencyV1.UrgencyStatus(urgencyV1.Open),
		}

		response := urgency.ToResponse()

		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "Marko", response.FirstName)
		assert.Equal(t, "Markovic", response.LastName)
		assert.Equal(t, "test@example.com", response.Email)
		assert.Equal(t, "+1234567890", response.ContactPhone)
		assert.Equal(t, "Test Location", response.Location)
		assert.Equal(t, "Test Description", response.Description)
		assert.Equal(t, urgencyV1.UrgencyLevel(urgencyV1.High), response.Level)
		assert.Equal(t, urgencyV1.UrgencyStatus(urgencyV1.Open), response.Status)
		assert.Equal(t, createdAt.Format(time.RFC3339), response.CreatedAt)
		assert.Equal(t, updatedAt.Format(time.RFC3339), response.UpdatedAt)
	})
}

func TestEmergencyAssignment_ToResponse(t *testing.T) {
	t.Parallel()

	t.Run("it converts emergency assignment to response correctly", func(t *testing.T) {
		createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)
		assignedAt := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)

		assignment := &EmergencyAssignment{
			Model: gorm.Model{
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			ID:         1,
			UrgencyID:  123,
			EmployeeID: 456,
			Status:     AssignmentAccepted,
			AssignedAt: assignedAt,
		}

		response := assignment.ToResponse()

		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, uint(123), response.UrgencyID)
		assert.Equal(t, uint(456), response.EmployeeID)
		assert.Equal(t, string(AssignmentAccepted), response.Status)
		assert.Equal(t, assignedAt.Format(time.RFC3339), response.AssignedAt)
		assert.Equal(t, createdAt.Format(time.RFC3339), response.CreatedAt)
		assert.Equal(t, updatedAt.Format(time.RFC3339), response.UpdatedAt)
	})
}

func TestNotification_ToResponse(t *testing.T) {
	t.Parallel()

	t.Run("it converts notification to response correctly with all fields", func(t *testing.T) {
		createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)
		lastAttemptAt := time.Date(2023, 1, 1, 11, 0, 0, 0, time.UTC)
		sentAt := time.Date(2023, 1, 1, 11, 30, 0, 0, time.UTC)

		notification := &Notification{
			Model: gorm.Model{
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			ID:               1,
			UrgencyID:        123,
			EmployeeID:       456,
			NotificationType: NotificationSMS,
			Recipient:        "+1234567890",
			Message:          "Emergency notification",
			Status:           NotificationSent,
			Attempts:         2,
			LastAttemptAt:    &lastAttemptAt,
			SentAt:           &sentAt,
			ErrorMessage:     "",
		}

		response := notification.ToResponse()

		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, uint(123), response.UrgencyID)
		assert.Equal(t, uint(456), response.EmployeeID)
		assert.Equal(t, string(NotificationSMS), response.NotificationType)
		assert.Equal(t, "+1234567890", response.Recipient)
		assert.Equal(t, "Emergency notification", response.Message)
		assert.Equal(t, string(NotificationSent), response.Status)
		assert.Equal(t, 2, response.Attempts)
		assert.Equal(t, "", response.ErrorMessage)
		assert.Equal(t, lastAttemptAt.Format(time.RFC3339), response.LastAttemptAt)
		assert.Equal(t, sentAt.Format(time.RFC3339), response.SentAt)
		assert.Equal(t, createdAt.Format(time.RFC3339), response.CreatedAt)
		assert.Equal(t, updatedAt.Format(time.RFC3339), response.UpdatedAt)
	})

	t.Run("it converts notification to response correctly with nil timestamps", func(t *testing.T) {
		createdAt := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC)

		notification := &Notification{
			Model: gorm.Model{
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			ID:               2,
			UrgencyID:        123,
			EmployeeID:       456,
			NotificationType: NotificationEmail,
			Recipient:        "test@example.com",
			Message:          "Emergency notification",
			Status:           NotificationPending,
			Attempts:         0,
			LastAttemptAt:    nil,
			SentAt:           nil,
			ErrorMessage:     "Connection failed",
		}

		response := notification.ToResponse()

		assert.Equal(t, uint(2), response.ID)
		assert.Equal(t, uint(123), response.UrgencyID)
		assert.Equal(t, uint(456), response.EmployeeID)
		assert.Equal(t, string(NotificationEmail), response.NotificationType)
		assert.Equal(t, "test@example.com", response.Recipient)
		assert.Equal(t, "Emergency notification", response.Message)
		assert.Equal(t, string(NotificationPending), response.Status)
		assert.Equal(t, 0, response.Attempts)
		assert.Equal(t, "Connection failed", response.ErrorMessage)
		assert.Equal(t, "", response.LastAttemptAt)
		assert.Equal(t, "", response.SentAt)
		assert.Equal(t, createdAt.Format(time.RFC3339), response.CreatedAt)
		assert.Equal(t, updatedAt.Format(time.RFC3339), response.UpdatedAt)
	})
}

func TestUrgency_UpdateWithRequest(t *testing.T) {
	t.Parallel()

	t.Run("it updates urgency with request correctly", func(t *testing.T) {
		urgency := &Urgency{
			FirstName:    "Marko",
			LastName:     "Markovic",
			Email:        "test@example.com",
			ContactPhone: "123456789",
			Location:     "Test Location",
			Description:  "Test Description",
			Level:        urgencyV1.UrgencyLevel(urgencyV1.High),
			Status:       urgencyV1.UrgencyStatus(urgencyV1.Open),
		}

		req := &urgencyV1.UrgencyUpdateRequest{
			FirstName:    "Updated",
			LastName:     "Name",
			Email:        "updated@example.com",
			ContactPhone: "987654321",
			Location:     "Updated Location",
			Description:  "Updated Description",
			Level:        urgencyV1.Critical,
			Status:       urgencyV1.InProgress,
		}

		urgency.UpdateWithRequest(req)

		assert.Equal(t, "Updated", urgency.FirstName)
		assert.Equal(t, "Name", urgency.LastName)
		assert.Equal(t, "updated@example.com", urgency.Email)
		assert.Equal(t, "987654321", urgency.ContactPhone)
		assert.Equal(t, "Updated Location", urgency.Location)
		assert.Equal(t, "Updated Description", urgency.Description)
		assert.Equal(t, urgencyV1.UrgencyLevel(urgencyV1.Critical), urgency.Level)
		assert.Equal(t, urgencyV1.UrgencyStatus(urgencyV1.InProgress), urgency.Status)
	})
}
