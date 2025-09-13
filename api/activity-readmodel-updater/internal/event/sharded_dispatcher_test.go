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
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestShardedDispatcher_Process(t *testing.T) {
	t.Parallel()
	logger := utils.NewTestLogger()

	t.Run("orders per activity across shards", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFB := service.NewMockFirebaseService(ctrl)
		dispatcher := NewShardedDispatcher(mockFB, logger, 4, 64)

		var mu sync.Mutex
		calls := map[uint][]int{}

		mockFB.EXPECT().SyncActivity(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
			func(ctx context.Context, ev activityV1.ActivityEvent) error {
				seq, _ := strconv.Atoi(ev.Description)
				mu.Lock()
				calls[ev.ActivityID] = append(calls[ev.ActivityID], seq)
				mu.Unlock()
				return nil
			},
		)

		mkMsg := func(id int, seq int) *pubsub.Message {
			b, _ := json.Marshal(activityV1.ActivityEvent{Type: "UPDATE", ActivityID: uint(id), Description: strconv.Itoa(seq)})
			return &pubsub.Message{Data: b, Attributes: map[string]string{}}
		}

		const N = 10
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			for i := 1; i <= N; i++ {
				_ = dispatcher.Process(t.Context(), mkMsg(1, i))
			}
		}()
		go func() {
			defer wg.Done()
			for i := 1; i <= N; i++ {
				_ = dispatcher.Process(t.Context(), mkMsg(2, i))
			}
		}()
		wg.Wait()

		mu.Lock()
		defer mu.Unlock()
		if assert.Len(t, calls[1], N) && assert.Len(t, calls[2], N) {
			for i := 1; i <= N; i++ {
				assert.Equal(t, i, calls[1][i-1])
				assert.Equal(t, i, calls[2][i-1])
			}
		}
	})

	t.Run("returns error on unparseable payload", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFB := service.NewMockFirebaseService(ctrl)
		d := NewShardedDispatcher(mockFB, logger, 2, 8)
		msg := &pubsub.Message{Data: []byte("notjson"), Attributes: map[string]string{}}
		if err := d.Process(t.Context(), msg); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("returns error when shard saturated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFB := service.NewMockFirebaseService(ctrl)
		mockFB.EXPECT().SyncActivity(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
			func(ctx context.Context, ev activityV1.ActivityEvent) error {
				time.Sleep(20 * time.Millisecond)
				return nil
			},
		)
		di := NewShardedDispatcher(mockFB, logger, 1, 1)
		sd := di.(*shardedDispatcher)
		sd.enqueueTimeout = 5 * time.Millisecond

		mkMsg := func(seq int) *pubsub.Message {
			b, _ := json.Marshal(activityV1.ActivityEvent{Type: "UPDATE", ActivityID: 1, Description: strconv.Itoa(seq)})
			return &pubsub.Message{Data: b}
		}

		go func() { _ = di.Process(t.Context(), mkMsg(1)) }()
		go func() { _ = di.Process(t.Context(), mkMsg(2)) }()

		deadline := time.Now().Add(120 * time.Millisecond)
		for len(sd.chans[0]) < 1 && time.Now().Before(deadline) {
			time.Sleep(1 * time.Millisecond)
		}
		if len(sd.chans[0]) < 1 {
			t.Fatalf("failed to fill shard buffer for saturation test")
		}

		err := di.Process(t.Context(), mkMsg(3))
		if err == nil {
			t.Fatalf("expected enqueue timeout error")
		}
	})

	t.Run("surfaces context deadline exceeded when work timeout elapses", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFB := service.NewMockFirebaseService(ctrl)
		mockFB.EXPECT().SyncActivity(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, ev activityV1.ActivityEvent) error {
				<-ctx.Done()
				return ctx.Err()
			},
		)
		di := NewShardedDispatcher(mockFB, logger, 1, 8)
		sd := di.(*shardedDispatcher)
		sd.workTimeout = 10 * time.Millisecond

		b, _ := json.Marshal(activityV1.ActivityEvent{Type: "UPDATE", ActivityID: 42, Description: "1"})
		err := di.Process(t.Context(), &pubsub.Message{Data: b})
		if err == nil {
			t.Fatalf("expected context deadline exceeded error")
		}
	})

	t.Run("returns error when context canceled before enqueue", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFB := service.NewMockFirebaseService(ctrl)
		di := NewShardedDispatcher(mockFB, logger, 1, 1)
		sd := di.(*shardedDispatcher)
		sd.enqueueTimeout = 50 * time.Millisecond

		b, _ := json.Marshal(activityV1.ActivityEvent{Type: "UPDATE", ActivityID: 7, Description: "1"})
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		err := di.Process(ctx, &pubsub.Message{Data: b})
		if err == nil {
			t.Fatalf("expected context canceled error")
		}
	})

	t.Run("it recovers from panic in SyncActivity and returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockFB := service.NewMockFirebaseService(ctrl)
		mockFB.EXPECT().SyncActivity(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, ev activityV1.ActivityEvent) error {
				panic("boom")
			},
		)
		di := NewShardedDispatcher(mockFB, logger, 1, 8)
		b, _ := json.Marshal(activityV1.ActivityEvent{Type: "UPDATE", ActivityID: 9, Description: "1"})
		err := di.Process(t.Context(), &pubsub.Message{Data: b})
		if err == nil {
			t.Fatalf("expected error from panic recovery")
		}
	})
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
