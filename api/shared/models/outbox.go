package models

import (
	"time"
)

// OutboxEvent represents an event in the outbox pattern for eventual consistency
type OutboxEvent struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	EventType   string     `gorm:"not null" json:"event_type"`
	AggregateID string     `gorm:"not null" json:"aggregate_id"`
	EventData   string     `gorm:"type:text" json:"event_data"`
	Published   bool       `gorm:"default:false" json:"published"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

func (OutboxEvent) TableName() string {
	return "outbox_events"
}
