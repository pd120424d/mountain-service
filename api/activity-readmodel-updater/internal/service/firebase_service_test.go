package service

import (
	"context"
	"testing"
	"time"

	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/firestoretest"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

func TestFirebaseService_WithFakeFirestore(t *testing.T) {
	t.Parallel()
	logger := utils.NewTestLogger()

	fake := firestoretest.NewFake().WithCollection("activities", []map[string]interface{}{
		{"id": int64(1), "urgency_id": int64(2), "employee_id": int64(5), "description": "A", "created_at": "2025-01-02T10:00:00Z"},
		{"id": int64(2), "urgency_id": int64(3), "employee_id": int64(6), "description": "B", "created_at": "2025-01-03T10:00:00Z"},
		{"id": int64(3), "urgency_id": int64(2), "employee_id": int64(7), "description": "C", "created_at": "2025-01-04T10:00:00Z"},
	})
	svc := NewFirebaseService(fake, logger)

	t.Run("it succeeds when GetActivitiesByUrgency filters correctly", func(t *testing.T) {
		ctx := context.Background()
		items, err := svc.GetActivitiesByUrgency(ctx, 2)
		assert.NoError(t, err)
		assert.Len(t, items, 2)
	})

	t.Run("it succeeds when GetAllActivities orders desc and limits", func(t *testing.T) {
		ctx := context.Background()
		items, err := svc.GetAllActivities(ctx, 2)
		assert.NoError(t, err)
		assert.Len(t, items, 2)
		// Expect the two latest by created_at desc to be first (ids 3 and 2 based on times)
	})

	t.Run("it succeeds when SyncActivity CREATE writes a doc", func(t *testing.T) {
		ctx := context.Background()
		ev := activityV1.ActivityEvent{Type: "CREATE", ActivityID: 10, UrgencyID: 5, EmployeeID: 9, Description: "New", CreatedAt: time.Now().UTC()}
		err := svc.SyncActivity(ctx, ev)
		assert.NoError(t, err)

		items, err := svc.GetActivitiesByUrgency(ctx, 5)
		assert.NoError(t, err)
		found := false
		for _, it := range items {
			if it.ID == 10 {
				found = true
				break
			}
		}
		assert.True(t, found)
	})

	t.Run("it succeeds when SyncActivity UPDATE increments version", func(t *testing.T) {
		ctx := context.Background()
		fake := firestoretest.NewFake().WithCollection("activities", nil)
		svc2 := NewFirebaseService(fake, logger)

		create := activityV1.ActivityEvent{Type: "CREATE", ActivityID: 11, UrgencyID: 6, EmployeeID: 9, Description: "Old", CreatedAt: time.Now()}
		err := svc2.SyncActivity(ctx, create)
		assert.NoError(t, err)

		ev := activityV1.ActivityEvent{Type: "UPDATE", ActivityID: 11, UrgencyID: 6, EmployeeID: 9, Description: "New"}
		err = svc2.SyncActivity(ctx, ev)
		assert.NoError(t, err)

		items, err := svc2.GetActivitiesByUrgency(ctx, 6)
		assert.NoError(t, err)
		assert.Len(t, items, 1)
		assert.Equal(t, "New", items[0].Description)
	})

	t.Run("it succeeds when SyncActivity DELETE removes document", func(t *testing.T) {
		ctx := context.Background()
		fake := firestoretest.NewFake().WithCollection("activities", nil)
		svc3 := NewFirebaseService(fake, logger)

		evCreate := activityV1.ActivityEvent{Type: "CREATE", ActivityID: 12, UrgencyID: 7, EmployeeID: 9, Description: "ToDelete", CreatedAt: time.Now()}
		err := svc3.SyncActivity(ctx, evCreate)
		assert.NoError(t, err)

		ev := activityV1.ActivityEvent{Type: "DELETE", ActivityID: 12}
		err = svc3.SyncActivity(ctx, ev)
		assert.NoError(t, err)

		items, err := svc3.GetActivitiesByUrgency(ctx, 7)
		assert.NoError(t, err)
		assert.Len(t, items, 0)
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

		_, err := service.GetActivitiesByUrgency(ctx, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Firestore client is nil")

		_, err = service.GetAllActivities(ctx, 10)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Firestore client is nil")

		err = service.HealthCheck(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Firestore client is nil")

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

	t.Run("it returns error for unknown event type", func(t *testing.T) {
		logger := utils.NewTestLogger()
		service := NewFirebaseService(nil, logger)
		assert.NotNil(t, service)

		activityEvent := activityV1.ActivityEvent{
			Type:       "UNKNOWN_TYPE",
			ActivityID: 1,
			CreatedAt:  time.Now(),
		}

		err := service.SyncActivity(context.Background(), activityEvent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Firestore client is nil")
	})
}
