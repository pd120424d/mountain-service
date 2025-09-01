package models

import (
	"time"
)

// OutboxEvent represents an event in the outbox pattern for eventual consistency
// Note: EventType field removed as deprecated/unused.
type OutboxEvent struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	AggregateID string     `gorm:"not null" json:"aggregateId"`
	EventData   string     `gorm:"type:text" json:"eventData"`
	Published   bool       `gorm:"default:false" json:"published"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
}

func (OutboxEvent) TableName() string {
	return "outbox_events"
}
