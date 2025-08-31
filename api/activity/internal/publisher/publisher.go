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
		cfg.Interval = 10 * time.Second
	}
	p := &Publisher{log: log.WithName("publisher"), repo: repo, pubsub: pubsubClient, config: cfg}
	p.topicFactory = func(name string) topic { return realTopic{t: p.pubsub.Topic(name)} }
	return p
}

func (p *Publisher) Start(ctx context.Context) {
	p.log.Infof("Starting outbox publisher: topic=%s interval=%s", p.config.TopicName, p.config.Interval)
	ticker := time.NewTicker(p.config.Interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				p.log.Info("Stopping outbox publisher")
				return
			case <-ticker.C:
				if err := p.processOnce(ctx); err != nil {
					p.log.Errorf("publisher cycle error: %v", err)
				}
			}
		}
	}()
}

func (p *Publisher) processOnce(ctx context.Context) error {
	log := p.log.WithContext(ctx)
	log.Info("Processing outbox events")
	events, err := p.repo.GetUnpublishedEvents(ctx, 100)
	if err != nil {
		return fmt.Errorf("get unpublished: %w", err)
	}
	if len(events) == 0 {
		log.Info("No unpublished events")
		return nil
	}

	topic := p.topicFactory(p.config.TopicName)
	defer topic.Stop()

	sent := 0
	for _, e := range events {
		// Build the payload from the contract dto
		payload := activityV1.OutboxEvent{
			ID:          e.ID,
			EventType:   e.EventType,
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
			Attributes: map[string]string{"eventType": e.EventType, "aggregateId": e.AggregateID},
		})
		if _, err := res.Get(ctx); err != nil {
			log.Errorf("failed to publish event id=%d: %v", e.ID, err)
			continue
		}
		if err := p.repo.MarkAsPublished(ctx, e.ID); err != nil {
			log.Errorf("failed to mark published id=%d: %v", e.ID, err)
			continue
		}
		sent++
	}
	log.Infof("Published %d/%d events", sent, len(events))
	return nil
}
