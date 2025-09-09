package events

//go:generate mockgen -source=sharded_dispatcher.go -destination=sharded_dispatcher_gomock.go -package=events github.com/pd120424d/mountain-service/api/activity-readmodel-updater/internal/event -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"hash/fnv"

	"cloud.google.com/go/pubsub"
	"github.com/pd120424d/mountain-service/api/activity-readmodel-updater/internal/service"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
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
	fb     service.FirebaseService
	logger utils.Logger
	chans  []chan workItem
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
	d := &shardedDispatcher{fb: fb, logger: logger.WithName("shardedDispatcher"), chans: make([]chan workItem, shards)}
	for i := 0; i < shards; i++ {
		ch := make(chan workItem, queueSize)
		d.chans[i] = ch
		go d.worker(ch)
	}
	return d
}

func (d *shardedDispatcher) worker(ch <-chan workItem) {
	for wi := range ch {
		err := d.fb.SyncActivity(wi.ctx, wi.ev)
		wi.done <- err
		close(wi.done)
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
	select {
	case d.chans[idx] <- wi:
		// wait for completion to preserve Pub/Sub semantics
		return <-wi.done
	case <-ctx.Done():
		return ctx.Err()
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
