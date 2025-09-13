package events

//go:generate mockgen -source=sharded_dispatcher.go -destination=sharded_dispatcher_gomock.go -package=events github.com/pd120424d/mountain-service/api/activity-readmodel-updater/internal/event -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"fmt"
	"hash/fnv"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pd120424d/mountain-service/api/activity-readmodel-updater/internal/service"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"go.uber.org/zap"
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
	fb             service.FirebaseService
	logger         utils.Logger
	chans          []chan workItem
	enqueueTimeout time.Duration
	workTimeout    time.Duration
}

type workItem struct {
	ctx  context.Context
	ev   activityV1.ActivityEvent
	done chan error
}

func NewShardedDispatcher(fb service.FirebaseService, logger utils.Logger, shards int, queueSize int) ShardedDispatcher {
	if shards <= 0 {
		shards = 8
	}
	if queueSize <= 0 {
		queueSize = 1024
	}
	d := &shardedDispatcher{
		fb:             fb,
		logger:         logger.WithName("shardedDispatcher"),
		chans:          make([]chan workItem, shards),
		enqueueTimeout: 2 * time.Second,
		workTimeout:    60 * time.Second,
	}
	for i := 0; i < shards; i++ {
		ch := make(chan workItem, queueSize)
		d.chans[i] = ch
		go d.worker(ch)
	}
	return d
}

func (d *shardedDispatcher) worker(ch <-chan workItem) {
	for wi := range ch {
		func(w workItem) {
			defer func() {
				if r := recover(); r != nil {
					d.logger.WithContext(w.ctx).Errorf("Panic in shard worker: %v", r)
					// Try to propagate error back to caller to avoid stalling Receive
					select {
					case w.done <- fmt.Errorf("panic: %v", r):
					default:
					}
					close(w.done)
				}
			}()

			ctx, cancel := context.WithTimeout(w.ctx, d.workTimeout)
			log := d.logger.WithContext(ctx) //
			stop := utils.TimeOperation(log, "ShardedDispatcher.worker.SyncActivity", zap.Int("activity_id", int(w.ev.ActivityID)), zap.String("type", w.ev.Type))

			errCh := make(chan error, 1)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						errCh <- fmt.Errorf("panic: %v", r)
					}
				}()
				errCh <- d.fb.SyncActivity(ctx, w.ev)
			}()
			var err error
			select {
			case err = <-errCh:
			case <-ctx.Done():
				err = ctx.Err()
				idx := int(hashKey(shardKey(w.ev, &pubsub.Message{ID: ""})) % uint64(len(d.chans)))
				qLen, qCap := len(d.chans[idx]), cap(d.chans[idx])
				log.Warnf("work timeout: activity_id=%d type=%s shard=%d shard_queue=%d/%d timeout=%s", int(w.ev.ActivityID), w.ev.Type, idx, qLen, qCap, d.workTimeout)
			}
			stop()
			cancel()
			w.done <- err
			close(w.done)
		}(wi)
	}
}

// Process parses the message, selects a shard by key and waits for completion.
// It returns the processing result so caller can ack/nack accordingly.
func (d *shardedDispatcher) Process(ctx context.Context, msg *pubsub.Message) error {
	// Parse once here to select shard by ActivityID
	ev, strat, err := Parse(msg.Data, msg.Attributes)
	if err != nil {
		// keep parity with Handler
		d.logger.WithContext(ctx).Errorf("Unrecognized event payload format, cannot parse message_id=%s", msg.ID)
		return err
	}
	_ = strat
	normalizeType(&ev)

	key := shardKey(ev, msg)
	idx := int(hashKey(key)) % len(d.chans)
	wi := workItem{ctx: ctx, ev: ev, done: make(chan error, 1)}

	// If the caller's context is already canceled returning that immediately
	if err := ctx.Err(); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case d.chans[idx] <- wi:
		// wait for completion to preserve Pub/Sub semantics
		return <-wi.done
	case <-time.After(d.enqueueTimeout):
		// Avoid stalling Receive loop when a shard queue is saturated
		qLen, qCap := len(d.chans[idx]), cap(d.chans[idx])
		d.logger.WithContext(ctx).Warnf("Shard queue saturated: shard=%d len=%d cap=%d; nacking for retry", idx, qLen, qCap)
		return fmt.Errorf("shard %d queue full (%d/%d)", idx, qLen, qCap)
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

func hashKey(k uint64) uint64 { return k }
