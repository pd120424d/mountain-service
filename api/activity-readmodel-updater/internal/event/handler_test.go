package events

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type firebaseMock struct {
	err  error
	last activityV1.ActivityEvent
}

func (m *firebaseMock) GetActivitiesByUrgency(ctx context.Context, urgencyID uint) ([]*models.Activity, error) {
	return nil, nil
}
func (m *firebaseMock) GetAllActivities(ctx context.Context, limit int) ([]*models.Activity, error) {
	return nil, nil
}
func (m *firebaseMock) SyncActivity(ctx context.Context, ev activityV1.ActivityEvent) error {
	m.last = ev
	return m.err
}
func (m *firebaseMock) HealthCheck(ctx context.Context) error { return nil }

func Test_Handler_it_succeeds_when_valid_event_is_processed(t *testing.T) {
	logger := utils.NewTestLogger()
	fb := &firebaseMock{}
	h := NewHandler(fb, logger)

	msg := &pubsub.Message{Data: []byte(`{"eventType":"activity.created","aggregateId":"activity-10","eventData":"{\"type\":\"activity.created\",\"activityId\":10,\"urgencyId\":2,\"employeeId\":3,\"description\":\"x\",\"createdAt\":\"2025-01-01T00:00:00Z\"}"}`)}
	if err := h.Handle(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fb.last.ActivityID != 10 || fb.last.Type != "CREATE" {
		t.Fatalf("unexpected event: %+v", fb.last)
	}
}

func Test_Handler_it_fails_when_activity_id_is_zero(t *testing.T) {
	logger := utils.NewTestLogger()
	fb := &firebaseMock{err: errors.New("invalid activity id: 0")}
	h := NewHandler(fb, logger)

	msg := &pubsub.Message{Data: []byte(`{"eventType":"activity.created","aggregateId":"activity-0","eventData":"{\"type\":\"activity.created\",\"activityId\":0,\"urgencyId\":2,\"employeeId\":3,\"description\":\"x\",\"createdAt\":\"2025-01-01T00:00:00Z\"}"}`)}
	if err := h.Handle(context.Background(), msg); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func Test_Handler_it_succeeds_when_legacy_event_is_processed(t *testing.T) {
	logger := utils.NewTestLogger()
	fb := &firebaseMock{}
	h := NewHandler(fb, logger)
	ev := activityV1.ActivityEvent{Type: "CREATE", ActivityID: 42, UrgencyID: 2, EmployeeID: 7, Description: "y", CreatedAt: time.Now()}
	b, _ := json.Marshal(ev)
	msg := &pubsub.Message{Data: b}
	if err := h.Handle(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fb.last.ActivityID != 42 {
		t.Fatalf("unexpected id: %d", fb.last.ActivityID)
	}
}
