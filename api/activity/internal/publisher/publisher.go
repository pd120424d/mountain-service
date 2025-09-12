package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/pd120424d/mountain-service/api/activity/internal/repositories"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type Config struct {
	TopicName string
	Interval  time.Duration
	BatchSize int
}

// publishResult abstracts Pub/Sub publish result for testability
type publishResult interface {
	Get(ctx context.Context) (serverID string, err error)
}

// topic abstracts Pub/Sub topic for testability
type topic interface {
	Publish(ctx context.Context, m *pubsub.Message) publishResult
	Stop()
}

// realTopic adapts *pubsub.Topic to the topic interface
type realTopic struct{ t *pubsub.Topic }

func (rt realTopic) Publish(ctx context.Context, m *pubsub.Message) publishResult {
	return rt.t.Publish(ctx, m)
}
func (rt realTopic) Stop() { rt.t.Stop() }

type Publisher struct {
	log          utils.Logger
	repo         repositories.OutboxRepository
	pubsub       *pubsub.Client
	config       Config
	topicFactory func(name string) topic
}

func New(log utils.Logger, repo repositories.OutboxRepository, pubsubClient *pubsub.Client, cfg Config) *Publisher {
	if cfg.Interval <= 0 {
		cfg.Interval = 1 * time.Second // lower default to reduce e2e latency
	}
	if cfg.BatchSize <= 0 {
		cfg.BatchSize = 100
	}
	p := &Publisher{log: log.WithName("publisher"), repo: repo, pubsub: pubsubClient, config: cfg}
	p.topicFactory = func(name string) topic {
		t := p.pubsub.Topic(name)
		// Moderate batching/concurrency to improve throughput while staying safe by default
		t.PublishSettings.NumGoroutines = 4
		t.PublishSettings.DelayThreshold = 25 * time.Millisecond // flush faster under low volume
		t.PublishSettings.CountThreshold = 100                   // smaller threshold improves latency
		return realTopic{t: t}
	}
	return p
}

func (p *Publisher) Start(ctx context.Context) {
	ctx, _ = utils.EnsureRequestID(ctx)
	// Adaptive polling: fast when backlog exists, exponential backoff when idle.
	minInterval := 1 * time.Second
	maxInterval := p.config.Interval
	if maxInterval < minInterval {
		maxInterval = minInterval
	}
	interval := minInterval

	minBatch := 50
	maxBatch := 2000
	batch := p.config.BatchSize
	if batch < minBatch {
		batch = minBatch
	}

	p.log.Infof("Starting outbox publisher: topic=%s minInterval=%s maxInterval=%s initialBatch=%d", p.config.TopicName, minInterval, maxInterval, batch)

	go func() {
		for {
			select {
			case <-ctx.Done():
				p.log.WithContext(ctx).Info("Stopping outbox publisher")
				return
			default:
			}

			// Peek to decide adaptively; tolerate errors by backing off
			events, err := p.repo.GetUnpublishedEvents(ctx, batch)
			if err != nil {
				p.log.WithContext(ctx).Errorf("get unpublished (peek) failed: %v", err)
				time.Sleep(interval)
				// back off a bit on error
				if interval < maxInterval {
					interval *= 2
					if interval > maxInterval {
						interval = maxInterval
					}
				}
				continue
			}

			if len(events) == 0 {
				// No backlog: back off interval up to max and slowly shrink batch down
				if interval < maxInterval {
					interval *= 2
					if interval > maxInterval {
						interval = maxInterval
					}
				}
				if batch > minBatch {
					batch = batch / 2
					if batch < minBatch {
						batch = minBatch
					}
				}
				time.Sleep(interval)
				continue
			}

			// Backlog present: run a publish cycle
			if err := p.processOnce(ctx); err != nil {
				p.log.WithContext(ctx).Errorf("publisher cycle error: %v", err)
			}

			// If we exactly filled the batch on peek, likely more backlog -> speed up and grow batch
			if len(events) == batch {
				if batch < maxBatch {
					batch *= 2
					if batch > maxBatch {
						batch = maxBatch
					}
				}
				interval = minInterval
				// Immediately continue to drain without sleeping
				continue
			}

			// Some backlog but less than batch: keep interval low for low latency and small sleep
			interval = minInterval
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (p *Publisher) processOnce(ctx context.Context) error {
	ctx, _ = utils.EnsureRequestID(ctx)
	log := p.log.WithContext(ctx)
	log.Info("Processing outbox events")
	batchSize := p.config.BatchSize
	if batchSize <= 0 {
		batchSize = 100
	}
	events, err := p.repo.GetUnpublishedEvents(ctx, batchSize)
	if err != nil {
		return fmt.Errorf("get unpublished: %w", err)
	}
	if len(events) == 0 {
		log.Info("No unpublished events")
		return nil
	}

	topic := p.topicFactory(p.config.TopicName)
	defer topic.Stop()

	type pending struct {
		id  uint
		res publishResult
	}
	pendings := make([]pending, 0, len(events))

	// Queue all publishes first to enable client-side batching
	for _, e := range events {
		payload := activityV1.OutboxEvent{
			ID:          e.ID,
			AggregateID: e.AggregateID,
			EventData:   e.EventData,
			Published:   e.Published,
			CreatedAt:   e.CreatedAt,
			PublishedAt: e.PublishedAt,
		}
		data, mErr := json.Marshal(payload)
		if mErr != nil {
			p.log.Errorf("failed to marshal outbox envelope id=%d: %v", e.ID, mErr)
			continue
		}
		res := topic.Publish(ctx, &pubsub.Message{
			Data:       data,
			Attributes: map[string]string{"aggregateId": e.AggregateID},
		})
		pendings = append(pendings, pending{id: e.ID, res: res})
	}

	sent := 0
	for _, pnd := range pendings {
		if _, err := pnd.res.Get(ctx); err != nil {
			log.Errorf("failed to publish event id=%d: %v", pnd.id, err)
			continue
		}
		if err := p.repo.MarkAsPublished(ctx, pnd.id); err != nil {
			log.Errorf("failed to mark published id=%d: %v", pnd.id, err)
			continue
		}
		sent++
	}
	log.Infof("Published %d/%d events", sent, len(events))
	return nil
}
