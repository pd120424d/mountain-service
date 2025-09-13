package events

//go:generate mockgen -source=sharded_dispatcher.go -destination=sharded_dispatcher_gomock.go -package=events github.com/pd120424d/mountain-service/api/activity-readmodel-updater/internal/event -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pd120424d/mountain-service/api/activity-readmodel-updater/internal/service"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"go.uber.org/zap"
)

const (
	DefaultMaxParallelWorkers = 16
	DefaultEnqueueTimeout     = 2 * time.Second
	DefaultWorkTimeout        = 60 * time.Second
)

type ShardedDispatcher interface {
	Process(ctx context.Context, msg *pubsub.Message) error
}

// shardedDispatcher serializes processing per activity ID while allowing
// parallelism across different activities.
//
// Flow: Receive callback parses -> chooses shard -> enqueues -> waits for result.
// This preserves Pub/Sub acking semantics (ack/nack after processing) and keeps
// ordering per key within a shard.
type shardedDispatcher struct {
	fb                 service.FirebaseService
	logger             utils.Logger
	enqueueTimeout     time.Duration
	workTimeout        time.Duration
	limiter            chan struct{}
	maxParallelWorkers int
	locks              sync.Map
}

func NewShardedDispatcher(fb service.FirebaseService, logger utils.Logger, shards int, queueSize int) ShardedDispatcher {
	if shards <= 0 {
		shards = DefaultMaxParallelWorkers
	}
	d := &shardedDispatcher{
		fb:                 fb,
		logger:             logger.WithName("shardedDispatcher"),
		enqueueTimeout:     DefaultEnqueueTimeout,
		workTimeout:        DefaultWorkTimeout,
		limiter:            make(chan struct{}, shards),
		maxParallelWorkers: shards,
	}
	return d
}

// Process parses the message and handles it with keyed per-activity ordering and a global concurrency limit.
func (d *shardedDispatcher) Process(ctx context.Context, msg *pubsub.Message) error {
	ev, strat, err := Parse(msg.Data, msg.Attributes)
	if err != nil {
		d.logger.WithContext(ctx).Errorf("Unrecognized event payload format, cannot parse message_id=%s", msg.ID)
		return err
	}
	_ = strat
	normalizeType(&ev)

	if err := ctx.Err(); err != nil {
		return err
	}

	select {
	case d.limiter <- struct{}{}:
		defer func() { <-d.limiter }()
	case <-ctx.Done():
		return ctx.Err()
	}

	key := shardKey(ev, msg)
	lk := d.getLock(uint(key))
	lk.Lock()
	defer lk.Unlock()

	wctx, cancel := context.WithTimeout(ctx, d.workTimeout)
	log := d.logger.WithContext(wctx)
	stop := utils.TimeOperation(log, "KeyedDispatcher.SyncActivity", zap.Int("activity_id", int(ev.ActivityID)), zap.String("type", ev.Type))
	defer stop()
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errCh <- fmt.Errorf("panic: %v", r)
			}
		}()
		errCh <- d.fb.SyncActivity(wctx, ev)
	}()

	select {
	case err := <-errCh:
		return err
	case <-wctx.Done():
		log.Warnf("work timeout: activity_id=%d type=%s timeout=%s", int(ev.ActivityID), ev.Type, d.workTimeout)
		return wctx.Err()
	}
}

func shardKey(ev activityV1.ActivityEvent, msg *pubsub.Message) uint64 {
	if ev.ActivityID != 0 {
		return uint64(ev.ActivityID)
	}
	// fallback to message ID hashing in case producer omitted ActivityID
	h := fnv.New64a()
	_, _ = h.Write([]byte(msg.ID))
	return h.Sum64()
}

func (d *shardedDispatcher) getLock(k uint) *sync.Mutex {
	if v, ok := d.locks.Load(k); ok {
		return v.(*sync.Mutex)
	}
	m := &sync.Mutex{}
	actual, _ := d.locks.LoadOrStore(k, m)
	return actual.(*sync.Mutex)
}
