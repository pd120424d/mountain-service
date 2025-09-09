package events

import (
	"context"
	"encoding/json"
	"hash/fnv"
	"strconv"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pd120424d/mountain-service/api/activity-readmodel-updater/internal/service"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

type fakeFB struct {
	mu    sync.Mutex
	calls map[int][]int // activityID -> sequence of "description as int"
	delay time.Duration
}

func (f *fakeFB) GetActivitiesByUrgency(ctx context.Context, urgencyID uint) ([]*models.Activity, error) {
	return nil, nil
}
func (f *fakeFB) GetAllActivities(ctx context.Context, limit int) ([]*models.Activity, error) {
	return nil, nil
}

func (f *fakeFB) SyncActivity(ctx context.Context, ev activityV1.ActivityEvent) error {
	if f.delay > 0 {
		time.Sleep(f.delay)
	}
	n, _ := strconv.Atoi(ev.Description)
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.calls == nil {
		f.calls = map[int][]int{}
	}
	f.calls[int(ev.ActivityID)] = append(f.calls[int(ev.ActivityID)], n)
	return nil
}
func (f *fakeFB) HealthCheck(ctx context.Context) error { return nil }

func TestShardedDispatcher_Process(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds when message is processed", func(t *testing.T) {

		logger := utils.NewTestLogger()
		fb := &fakeFB{delay: 5 * time.Millisecond}
		dispatcher := NewShardedDispatcher(fb, logger, 4, 64)

		mkMsg := func(id int, seq int) *pubsub.Message {
			b, _ := json.Marshal(activityV1.ActivityEvent{Type: "UPDATE", ActivityID: uint(id), Description: strconv.Itoa(seq)})
			return &pubsub.Message{Data: b, Attributes: map[string]string{}}
		}

		// Two activities processed concurrently; ensure per-activity order preserved
		const N = 50
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			for i := 1; i <= N; i++ {
				_ = dispatcher.Process(context.Background(), mkMsg(1, i))
			}
		}()
		go func() {
			defer wg.Done()
			for i := 1; i <= N; i++ {
				_ = dispatcher.Process(context.Background(), mkMsg(2, i))
			}
		}()
		wg.Wait()

		fb.mu.Lock()
		defer fb.mu.Unlock()
		seq1 := fb.calls[1]
		seq2 := fb.calls[2]
		if assert.Len(t, seq1, N) && assert.Len(t, seq2, N) {
			for i := 1; i <= N; i++ {
				assert.Equal(t, i, seq1[i-1])
				assert.Equal(t, i, seq2[i-1])
			}
		}

		var _ service.FirebaseService = (*fakeFB)(nil)
	})

}

func TestShardedDispatcher_Process_ErrorOnUnparseable(t *testing.T) {
	logger := utils.NewTestLogger()
	fb := &fakeFB{}
	d := NewShardedDispatcher(fb, logger, 2, 8)
	msg := &pubsub.Message{Data: []byte("notjson"), Attributes: map[string]string{}}
	if err := d.Process(context.Background(), msg); err == nil {
		t.Fatalf("expected error")
	}
}

func Test_shardKey_FallbackToMessageID(t *testing.T) {
	// ActivityID = 0 -> use hashed message ID
	msg := &pubsub.Message{ID: "abc-123"}
	var ev activityV1.ActivityEvent // zero ActivityID
	k := shardKey(ev, msg)
	// compute expected FNV-1a 64 of msg.ID
	h := fnv.New64a()
	_, _ = h.Write([]byte(msg.ID))
	expected := h.Sum64()
	assert.Equal(t, expected, k)
}
