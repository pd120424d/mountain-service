package models

import (
	"time"

	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
)

type Activity struct {
	ID          uint                     `json:"id"`
	Type        activityV1.ActivityType  `json:"type"`
	Level       activityV1.ActivityLevel `json:"level"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	ActorID     *uint                    `json:"actor_id,omitempty"`
	ActorName   string                   `json:"actor_name"`
	TargetID    *uint                    `json:"target_id,omitempty"`
	TargetType  string                   `json:"target_type"`
	Metadata    string                   `json:"metadata,omitempty"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

func (a *Activity) ToResponse() *activityV1.ActivityResponse {
	return &activityV1.ActivityResponse{
		ID:          a.ID,
		Type:        a.Type,
		Level:       a.Level,
		Title:       a.Title,
		Description: a.Description,
		ActorID:     a.ActorID,
		ActorName:   a.ActorName,
		TargetID:    a.TargetID,
		TargetType:  a.TargetType,
		Metadata:    a.Metadata,
		CreatedAt:   a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   a.UpdatedAt.Format(time.RFC3339),
	}
}
