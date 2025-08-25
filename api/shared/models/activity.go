package models

import (
	"time"

	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
)

type Activity struct {
	ID          uint      `json:"id"`
	Description string    `json:"description"`
	EmployeeID  uint      `json:"employee_id"`
	UrgencyID   uint      `json:"urgency_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (a *Activity) ToResponse() *activityV1.ActivityResponse {
	return &activityV1.ActivityResponse{
		ID:          a.ID,
		Description: a.Description,
		EmployeeID:  a.EmployeeID,
		UrgencyID:   a.UrgencyID,
		CreatedAt:   a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   a.UpdatedAt.Format(time.RFC3339),
	}
}
