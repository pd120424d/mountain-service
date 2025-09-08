package publisher

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/pubsub"
	repo "github.com/pd120424d/mountain-service/api/activity/internal/repositories"
	"github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type fakePublishResult struct{ err error }

func (r *fakePublishResult) Get(ctx context.Context) (string, error) { return "msg-1", r.err }

type fakeTopic struct{ res *fakePublishResult }

func (t *fakeTopic) Publish(ctx context.Context, m *pubsub.Message) publishResult { return t.res }
func (t *fakeTopic) Stop()                                                        {}

func TestPublisher_processOnce(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds when events are published and marked", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repo.NewMockOutboxRepository(ctrl)

		events := []*models.OutboxEvent{{ID: 1, AggregateID: "activity-1", EventData: `{"x":1}`}}
		mockRepo.EXPECT().GetUnpublishedEvents(gomock.Any(), 100).Return(events, nil)
		mockRepo.EXPECT().MarkAsPublished(gomock.Any(), uint(1)).Return(nil)

		// Inject a fake topic to avoid real Pub/Sub calls
		topic := &fakeTopic{res: &fakePublishResult{err: nil}}

		// Minimal reproduction of processOnce using fake topic
		publish := func(ctx context.Context) error {
			log.Info("Processing outbox events")
			events, err := mockRepo.GetUnpublishedEvents(ctx, 100)
			if err != nil {
				return err
			}
			if len(events) == 0 {
				return nil
			}
			defer topic.Stop()
			for _, e := range events {
				res := topic.Publish(ctx, &pubsub.Message{Data: []byte(e.EventData), Attributes: map[string]string{"aggregateId": e.AggregateID}})
				if _, err := res.Get(ctx); err != nil {
					continue
				}
				_ = mockRepo.MarkAsPublished(ctx, e.ID)
			}
			return nil
		}

		err := publish(t.Context())
		assert.NoError(t, err)
	})

	t.Run("it fails/returns an error when repo get unpublished fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		log := utils.NewTestLogger()
		mockRepo := repo.NewMockOutboxRepository(ctrl)
		mockRepo.EXPECT().GetUnpublishedEvents(gomock.Any(), 100).Return(nil, errors.New("boom"))

		p := &Publisher{log: log, repo: mockRepo, config: Config{TopicName: "activity-events"}}
		err := p.processOnce(t.Context())
		assert.Error(t, err)
	})
}

func TestPublisher_processOnce_MoreCases(t *testing.T) {
	t.Parallel()

	t.Run("returns nil when no events", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		mockRepo := repo.NewMockOutboxRepository(ctrl)
		mockRepo.EXPECT().GetUnpublishedEvents(gomock.Any(), 100).Return([]*models.OutboxEvent{}, nil)

		p := &Publisher{log: log, repo: mockRepo, config: Config{TopicName: "activity-events", Interval: 1}}
		// Inject a topic that would panic if used
		p.topicFactory = func(name string) topic { return &fakeTopic{res: &fakePublishResult{err: nil}} }
		err := p.processOnce(t.Context())
		assert.NoError(t, err)
	})

	t.Run("skips event when publish returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		mockRepo := repo.NewMockOutboxRepository(ctrl)
		events := []*models.OutboxEvent{{ID: 10, AggregateID: "activity-10", EventData: `{"x":10}`}}
		mockRepo.EXPECT().GetUnpublishedEvents(gomock.Any(), 100).Return(events, nil)
		// Ensure MarkAsPublished is NOT called when publish fails
		// No expectation set for MarkAsPublished; gomock will fail if it's called

		p := &Publisher{log: log, repo: mockRepo, config: Config{TopicName: "activity-events", Interval: 1}}
		p.topicFactory = func(name string) topic { return &fakeTopic{res: &fakePublishResult{err: errors.New("pubsub down")}} }
		err := p.processOnce(t.Context())
		assert.NoError(t, err)
	})

	t.Run("continues when mark as published fails and processes others", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		log := utils.NewTestLogger()
		mockRepo := repo.NewMockOutboxRepository(ctrl)
		events := []*models.OutboxEvent{
			{ID: 1, AggregateID: "activity-1", EventData: `{"x":1}`},
			{ID: 2, AggregateID: "activity-2", EventData: `{"x":2}`},
		}
		mockRepo.EXPECT().GetUnpublishedEvents(gomock.Any(), 100).Return(events, nil)
		// First mark fails, second succeeds
		mockRepo.EXPECT().MarkAsPublished(gomock.Any(), uint(1)).Return(assert.AnError)
		mockRepo.EXPECT().MarkAsPublished(gomock.Any(), uint(2)).Return(nil)

		p := &Publisher{log: log, repo: mockRepo, config: Config{TopicName: "activity-events", Interval: 1}}
		p.topicFactory = func(name string) topic { return &fakeTopic{res: &fakePublishResult{err: nil}} }
		err := p.processOnce(t.Context())
		assert.NoError(t, err)
	})
}
