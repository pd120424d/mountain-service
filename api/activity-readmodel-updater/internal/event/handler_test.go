package events

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pd120424d/mountain-service/api/activity-readmodel-updater/internal/service"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"go.uber.org/mock/gomock"
)

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	logger := utils.NewTestLogger()
	ctx := t.Context()

	t.Run("it succeeds when valid envelope event is processed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFB := service.NewMockFirebaseService(ctrl)
		h := NewHandler(mockFB, logger)
		var captured activityV1.ActivityEvent
		mockFB.EXPECT().SyncActivity(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, ev activityV1.ActivityEvent) error {
				captured = ev
				return nil
			},
		)
		msg := &pubsub.Message{Data: []byte(`{"eventType":"activity.created","aggregateId":"activity-10","eventData":"{\"type\":\"activity.created\",\"activityId\":10,\"urgencyId\":2,\"employeeId\":3,\"description\":\"x\",\"createdAt\":\"2025-01-01T00:00:00Z\"}"}`)}
		if err := h.Handle(ctx, msg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if captured.ActivityID != 10 || captured.Type != "CREATE" {
			t.Fatalf("unexpected event: %+v", captured)
		}
	})

	t.Run("it fails when activity id is zero (SyncActivity returns error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFB := service.NewMockFirebaseService(ctrl)
		h := NewHandler(mockFB, logger)
		mockFB.EXPECT().SyncActivity(gomock.Any(), gomock.Any()).Return(errors.New("invalid activity id: 0"))
		msg := &pubsub.Message{Data: []byte(`{"eventType":"activity.created","aggregateId":"activity-0","eventData":"{\"type\":\"activity.created\",\"activityId\":0,\"urgencyId\":2,\"employeeId\":3,\"description\":\"x\",\"createdAt\":\"2025-01-01T00:00:00Z\"}"}`)}
		if err := h.Handle(ctx, msg); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("it succeeds when legacy event is processed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFB := service.NewMockFirebaseService(ctrl)
		h := NewHandler(mockFB, logger)
		var captured activityV1.ActivityEvent
		mockFB.EXPECT().SyncActivity(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, ev activityV1.ActivityEvent) error { captured = ev; return nil },
		)
		ev := activityV1.ActivityEvent{Type: "CREATE", ActivityID: 42, UrgencyID: 2, EmployeeID: 7, Description: "y", CreatedAt: time.Now()}
		b, _ := json.Marshal(ev)
		msg := &pubsub.Message{Data: b}
		if err := h.Handle(ctx, msg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if captured.ActivityID != 42 {
			t.Fatalf("unexpected id: %d", captured.ActivityID)
		}
	})

	t.Run("it returns error when SyncActivity fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFB := service.NewMockFirebaseService(ctrl)
		h := NewHandler(mockFB, logger)
		mockFB.EXPECT().SyncActivity(gomock.Any(), gomock.Any()).Return(errors.New("boom"))
		ev := activityV1.ActivityEvent{Type: "CREATE", ActivityID: 99, UrgencyID: 1, EmployeeID: 1, Description: "x", CreatedAt: time.Now()}
		b, _ := json.Marshal(ev)
		msg := &pubsub.Message{Data: b}
		if err := h.Handle(ctx, msg); err == nil {
			t.Fatalf("expected error when SyncActivity fails")
		}
	})

	t.Run("it returns error on invalid payload", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFB := service.NewMockFirebaseService(ctrl)
		h := NewHandler(mockFB, logger)
		msg := &pubsub.Message{Data: []byte("notjson")}
		if err := h.Handle(ctx, msg); err == nil {
			t.Fatalf("expected error for invalid payload")
		}
	})

	t.Run("it normalizes activity.updated to UPDATE", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFB := service.NewMockFirebaseService(ctrl)
		h := NewHandler(mockFB, logger)
		var captured activityV1.ActivityEvent
		mockFB.EXPECT().SyncActivity(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, ev activityV1.ActivityEvent) error { captured = ev; return nil },
		)
		msg := &pubsub.Message{Data: []byte(`{"eventType":"activity.updated","aggregateId":"activity-5","eventData":"{\"type\":\"activity.updated\",\"activityId\":5}"}`)}
		if err := h.Handle(ctx, msg); err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if captured.Type != "UPDATE" {
			t.Fatalf("expected UPDATE, got %s", captured.Type)
		}
	})

	t.Run("it normalizes activity.deleted to DELETE", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFB := service.NewMockFirebaseService(ctrl)
		h := NewHandler(mockFB, logger)
		var captured activityV1.ActivityEvent
		mockFB.EXPECT().SyncActivity(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, ev activityV1.ActivityEvent) error { captured = ev; return nil },
		)
		msg := &pubsub.Message{Data: []byte(`{"eventType":"activity.deleted","aggregateId":"activity-6","eventData":"{\"type\":\"activity.deleted\",\"activityId\":6}"}`)}
		if err := h.Handle(ctx, msg); err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if captured.Type != "DELETE" {
			t.Fatalf("expected DELETE, got %s", captured.Type)
		}
	})
}

func Test_normalizeType(t *testing.T) {
	t.Parallel()

	t.Run("it normalizes activity.created to CREATE", func(t *testing.T) {
		ev := activityV1.ActivityEvent{Type: "activity.created"}
		normalizeType(&ev)
		if ev.Type != "CREATE" {
			t.Fatalf("expected CREATE, got %s", ev.Type)
		}
	})

	t.Run("it normalizes activity.updated to UPDATE", func(t *testing.T) {
		ev := activityV1.ActivityEvent{Type: "activity.updated"}
		normalizeType(&ev)
		if ev.Type != "UPDATE" {
			t.Fatalf("expected UPDATE, got %s", ev.Type)
		}
	})

	t.Run("it normalizes activity.deleted to DELETE", func(t *testing.T) {
		ev := activityV1.ActivityEvent{Type: "activity.deleted"}
		normalizeType(&ev)
		if ev.Type != "DELETE" {
			t.Fatalf("expected DELETE, got %s", ev.Type)
		}
	})
}
