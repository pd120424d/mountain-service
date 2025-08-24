package repositories

//go:generate mockgen -source=outbox_repository.go -destination=outbox_repository_gomock.go -package=repositories mountain_service/activity/internal/repositories -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"fmt"
	"time"

	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"gorm.io/gorm"
)

type OutboxRepository interface {
	GetUnpublishedEvents(limit int) ([]*models.OutboxEvent, error)
	MarkAsPublished(eventID uint) error
	MarkOutboxEventAsPublished(event *models.OutboxEvent) error
}

type outboxRepository struct {
	log utils.Logger
	db  *gorm.DB
}

func NewOutboxRepository(log utils.Logger, db *gorm.DB) OutboxRepository {
	return &outboxRepository{log: log.WithName("outboxRepository"), db: db}
}

func (r *outboxRepository) GetUnpublishedEvents(limit int) ([]*models.OutboxEvent, error) {
	r.log.Infof("Getting unpublished events with limit: %d", limit)

	var events []*models.OutboxEvent
	err := r.db.Where("published = ?", false).
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error

	if err != nil {
		r.log.Errorf("Failed to get unpublished events: %v", err)
		return nil, fmt.Errorf("failed to get unpublished events: %w", err)
	}

	r.log.Infof("Found %d unpublished events", len(events))
	return events, nil
}

func (r *outboxRepository) MarkAsPublished(eventID uint) error {
	r.log.Infof("Marking event as published: %d", eventID)

	err := r.db.Model(&activityV1.OutboxEvent{}).
		Where("id = ?", eventID).
		Updates(map[string]interface{}{
			"published":    true,
			"published_at": time.Now(),
		}).Error

	if err != nil {
		r.log.Errorf("Failed to mark event as published: event_id=%d, error=%v", eventID, err)
		return fmt.Errorf("failed to mark event as published: %w", err)
	}

	r.log.Infof("Event marked as published successfully: %d", eventID)
	return nil
}

// MarkOutboxEventAsPublished marks an outbox event as published and updates the timestamp
func (r *outboxRepository) MarkOutboxEventAsPublished(event *models.OutboxEvent) error {
	r.log.Infof("Marking outbox event as published: event_id=%d", event.ID)

	now := time.Now()
	event.Published = true
	event.PublishedAt = &now

	if err := r.db.Save(event).Error; err != nil {
		r.log.Errorf("Failed to mark outbox event as published: event_id=%d, error=%v", event.ID, err)
		return fmt.Errorf("failed to mark outbox event as published: %w", err)
	}

	r.log.Infof("Outbox event marked as published successfully: event_id=%d", event.ID)
	return nil
}
