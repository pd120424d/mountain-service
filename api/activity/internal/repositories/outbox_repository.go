package repositories

//go:generate mockgen -source=outbox_repository.go -destination=outbox_repository_gomock.go -package=repositories mountain_service/activity/internal/repositories -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"fmt"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"gorm.io/gorm"
)

type OutboxRepository interface {
	GetUnpublishedEvents(ctx context.Context, limit int) ([]*models.OutboxEvent, error)
	MarkAsPublished(ctx context.Context, eventID uint) error
	MarkOutboxEventAsPublished(ctx context.Context, event *models.OutboxEvent) error
}

type outboxRepository struct {
	log utils.Logger
	db  *gorm.DB
}

func NewOutboxRepository(log utils.Logger, db *gorm.DB) OutboxRepository {
	return &outboxRepository{log: log.WithName("outboxRepository"), db: db}
}

func (r *outboxRepository) GetUnpublishedEvents(ctx context.Context, limit int) ([]*models.OutboxEvent, error) {
	log := r.log.WithContext(ctx)
	log.Infof("Getting unpublished events with limit: %d", limit)

	var events []*models.OutboxEvent
	err := r.db.WithContext(ctx).Where("published = ?", false).
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error

	if err != nil {
		log.Errorf("Failed to get unpublished events: %v", err)
		return nil, fmt.Errorf("failed to get unpublished events: %w", err)
	}

	log.Infof("Found %d unpublished events", len(events))
	return events, nil
}

func (r *outboxRepository) MarkAsPublished(ctx context.Context, eventID uint) error {
	log := r.log.WithContext(ctx)
	log.Infof("Marking event as published: %d", eventID)

	err := r.db.WithContext(ctx).Model(&models.OutboxEvent{}).
		Where("id = ?", eventID).
		Updates(map[string]interface{}{
			"published":    true,
			"published_at": time.Now(),
		}).Error

	if err != nil {
		log.Errorf("Failed to mark event as published: event_id=%d, error=%v", eventID, err)
		return fmt.Errorf("failed to mark event as published: %w", err)
	}

	log.Infof("Event marked as published successfully: %d", eventID)
	return nil
}

func (r *outboxRepository) MarkOutboxEventAsPublished(ctx context.Context, event *models.OutboxEvent) error {
	log := r.log.WithContext(ctx)
	log.Infof("Marking outbox event as published: event_id=%d", event.ID)

	now := time.Now()
	event.Published = true
	event.PublishedAt = &now

	if err := r.db.WithContext(ctx).Save(event).Error; err != nil {
		log.Errorf("Failed to mark outbox event as published: event_id=%d, error=%v", event.ID, err)
		return fmt.Errorf("failed to mark outbox event as published: %w", err)
	}

	log.Infof("Outbox event marked as published successfully: event_id=%d", event.ID)
	return nil
}
