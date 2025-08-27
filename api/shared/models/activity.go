package models

import (
	"time"

	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
)

type Activity struct {
	ID          uint      `json:"id"`
	Description string    `json:"description"`
	EmployeeID  uint      `json:"employeeId"`
	UrgencyID   uint      `json:"urgencyId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
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
