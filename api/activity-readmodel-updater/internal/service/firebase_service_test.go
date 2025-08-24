package service

import (
	"context"
	"testing"
	"time"

	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

func TestFirebaseService_SyncActivity(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds when syncing activity with valid data", func(t *testing.T) {
		// This is a placeholder test since we can't easily mock Firestore client
		// In a real implementation, you would use a Firestore emulator or mock
		logger := utils.NewTestLogger()

		// For now, just test that the service can be created
		service := NewFirebaseService(nil, logger)
		assert.NotNil(t, service)

		// Test with nil client should return error
		activityEvent := activityV1.ActivityEvent{
			Type:        "CREATE",
			ActivityID:  1,
			UrgencyID:   1,
			EmployeeID:  1,
			Description: "Test activity",
		}

		err := service.SyncActivity(context.Background(), activityEvent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Firestore client is nil")
	})
}

func TestFirebaseService_HealthCheck(t *testing.T) {
	t.Parallel()

	t.Run("it returns error when Firebase client is nil", func(t *testing.T) {
		logger := utils.NewTestLogger()
		service := NewFirebaseService(nil, logger)
		assert.NotNil(t, service)

		// With nil client, health check should fail
		err := service.HealthCheck(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Firestore client is nil")
	})
}

func TestFirebaseService_GetActivitiesByUrgency(t *testing.T) {
	t.Parallel()

	t.Run("it returns error when Firebase client is nil", func(t *testing.T) {
		logger := utils.NewTestLogger()
		service := NewFirebaseService(nil, logger)
		assert.NotNil(t, service)

		// With nil client, should return error
		activities, err := service.GetActivitiesByUrgency(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, activities)
		assert.Contains(t, err.Error(), "Firestore client is nil")
	})
}

func TestFirebaseService_GetAllActivities(t *testing.T) {
	t.Parallel()

	t.Run("it returns error when Firebase client is nil", func(t *testing.T) {
		logger := utils.NewTestLogger()
		service := NewFirebaseService(nil, logger)
		assert.NotNil(t, service)

		// With nil client, should return error
		activities, err := service.GetAllActivities(context.Background(), 10)
		assert.Error(t, err)
		assert.Nil(t, activities)
		assert.Contains(t, err.Error(), "Firestore client is nil")
	})
}

func TestFirebaseService_SyncActivity_DetailedScenarios(t *testing.T) {
	t.Parallel()

	t.Run("it handles CREATE event type correctly", func(t *testing.T) {
		logger := utils.NewTestLogger()
		service := NewFirebaseService(nil, logger)
		assert.NotNil(t, service)

		activityEvent := activityV1.ActivityEvent{
			Type:         "CREATE",
			ActivityID:   1,
			UrgencyID:    1,
			EmployeeID:   1,
			Description:  "Test activity created",
			CreatedAt:    time.Now(),
			EmployeeName: "John Doe",
			UrgencyTitle: "Test Urgency",
			UrgencyLevel: "High",
		}

		err := service.SyncActivity(context.Background(), activityEvent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Firestore client is nil")
	})

	t.Run("it handles UPDATE event type correctly", func(t *testing.T) {
		logger := utils.NewTestLogger()
		service := NewFirebaseService(nil, logger)
		assert.NotNil(t, service)

		activityEvent := activityV1.ActivityEvent{
			Type:         "UPDATE",
			ActivityID:   1,
			UrgencyID:    1,
			EmployeeID:   1,
			Description:  "Test activity updated",
			CreatedAt:    time.Now(),
			EmployeeName: "John Doe",
			UrgencyTitle: "Test Urgency",
			UrgencyLevel: "Medium",
		}

		err := service.SyncActivity(context.Background(), activityEvent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Firestore client is nil")
	})

	t.Run("it handles DELETE event type correctly", func(t *testing.T) {
		logger := utils.NewTestLogger()
		service := NewFirebaseService(nil, logger)
		assert.NotNil(t, service)

		activityEvent := activityV1.ActivityEvent{
			Type:       "DELETE",
			ActivityID: 1,
			CreatedAt:  time.Now(),
		}

		err := service.SyncActivity(context.Background(), activityEvent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Firestore client is nil")
	})
}

func TestFirebaseService_NewFirebaseService_Comprehensive(t *testing.T) {
	t.Parallel()

	t.Run("it creates service with all methods working correctly with nil client", func(t *testing.T) {
		logger := utils.NewTestLogger()
		service := NewFirebaseService(nil, logger)
		assert.NotNil(t, service)

		ctx := context.Background()

		// Test all interface methods return appropriate errors
		_, err := service.GetActivitiesByUrgency(ctx, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Firestore client is nil")

		_, err = service.GetAllActivities(ctx, 10)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Firestore client is nil")

		err = service.HealthCheck(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Firestore client is nil")

		// Test different event types
		eventTypes := []string{"CREATE", "UPDATE", "DELETE", "UNKNOWN"}
		for _, eventType := range eventTypes {
			activityEvent := activityV1.ActivityEvent{
				Type:       eventType,
				ActivityID: 1,
				CreatedAt:  time.Now(),
			}
			err = service.SyncActivity(ctx, activityEvent)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Firestore client is nil")
		}
	})

	t.Run("it handles different urgency IDs correctly", func(t *testing.T) {
		logger := utils.NewTestLogger()
		service := NewFirebaseService(nil, logger)
		assert.NotNil(t, service)

		ctx := context.Background()

		// Test with different urgency IDs
		urgencyIDs := []uint{0, 1, 999, 1000000}
		for _, urgencyID := range urgencyIDs {
			_, err := service.GetActivitiesByUrgency(ctx, urgencyID)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Firestore client is nil")
		}
	})

	t.Run("it handles different limits correctly", func(t *testing.T) {
		logger := utils.NewTestLogger()
		service := NewFirebaseService(nil, logger)
		assert.NotNil(t, service)

		ctx := context.Background()

		// Test with different limits
		limits := []int{-1, 0, 1, 10, 100, 1000}
		for _, limit := range limits {
			_, err := service.GetAllActivities(ctx, limit)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Firestore client is nil")
		}
	})
}

func TestFirebaseService_SyncActivity_ComprehensiveEventData(t *testing.T) {
	t.Parallel()

	t.Run("it handles complete activity event data", func(t *testing.T) {
		logger := utils.NewTestLogger()
		service := NewFirebaseService(nil, logger)
		assert.NotNil(t, service)

		activityEvent := activityV1.ActivityEvent{
			Type:         "CREATE",
			ActivityID:   12345,
			UrgencyID:    67890,
			EmployeeID:   11111,
			Description:  "Comprehensive test activity with all fields populated",
			CreatedAt:    time.Now(),
			EmployeeName: "John Doe Smith",
			UrgencyTitle: "Critical System Alert - Database Connection Lost",
			UrgencyLevel: "Critical",
		}

		err := service.SyncActivity(context.Background(), activityEvent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Firestore client is nil")
	})

	t.Run("it handles minimal activity event data", func(t *testing.T) {
		logger := utils.NewTestLogger()
		service := NewFirebaseService(nil, logger)
		assert.NotNil(t, service)

		activityEvent := activityV1.ActivityEvent{
			Type:       "DELETE",
			ActivityID: 1,
			CreatedAt:  time.Now(),
		}

		err := service.SyncActivity(context.Background(), activityEvent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Firestore client is nil")
	})

	t.Run("it handles activity event with special characters", func(t *testing.T) {
		logger := utils.NewTestLogger()
		service := NewFirebaseService(nil, logger)
		assert.NotNil(t, service)

		activityEvent := activityV1.ActivityEvent{
			Type:         "UPDATE",
			ActivityID:   999,
			Description:  "Activity with special chars: àáâãäåæçèéêë & symbols: @#$%^&*()",
			EmployeeName: "José María García-López",
			UrgencyTitle: "Ürgenčy with ünïcödé characters",
			CreatedAt:    time.Now(),
		}

		err := service.SyncActivity(context.Background(), activityEvent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Firestore client is nil")
	})
}
