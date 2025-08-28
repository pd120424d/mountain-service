package model

import (
	"strings"
	"time"

	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
	"gorm.io/gorm"
)

// Use types from swagger_types.go for consistency

type Urgency struct {
	gorm.Model
	ID           uint                    `gorm:"primaryKey"`
	FirstName    string                  `gorm:"not null"`
	LastName     string                  `gorm:"not null"`
	Email        string                  `gorm:""`
	ContactPhone string                  `gorm:"not null"`
	Location     string                  `gorm:"not null"`
	Description  string                  `gorm:"not null"`
	Level        urgencyV1.UrgencyLevel  `gorm:"type:text;not null;default:'Medium'"`
	Status       urgencyV1.UrgencyStatus `gorm:"type:text;not null;default:'Open'"`

	AssignedEmployeeID *uint      `gorm:"index"`
	AssignedAt         *time.Time `gorm:"index"`
}

type (
	NotificationStatus string
	NotificationType   string
)

const (
	NotificationPending NotificationStatus = "pending"
	NotificationSent    NotificationStatus = "sent"
	NotificationFailed  NotificationStatus = "failed"

	NotificationSMS   NotificationType = "sms"
	NotificationEmail NotificationType = "email"
)

// Notification represents a notification to be sent to an employee
type Notification struct {
	gorm.Model
	ID               uint               `gorm:"primaryKey"`
	UrgencyID        uint               `gorm:"not null;index"`
	EmployeeID       uint               `gorm:"not null;index"`
	NotificationType NotificationType   `gorm:"type:text;not null"`
	Recipient        string             `gorm:"not null"` // phone or email
	Message          string             `gorm:"type:text;not null"`
	Status           NotificationStatus `gorm:"type:text;not null;default:'pending'"`
	Attempts         int                `gorm:"default:0"`
	LastAttemptAt    *time.Time
	SentAt           *time.Time
	ErrorMessage     string `gorm:"type:text"`

	Urgency *Urgency `gorm:"foreignKey:UrgencyID"`
}

func UrgencyLevelFromString(s string) urgencyV1.UrgencyLevel {
	switch strings.ToLower(s) {
	case "low":
		return urgencyV1.UrgencyLevel(urgencyV1.Low)
	case "medium":
		return urgencyV1.UrgencyLevel(urgencyV1.Medium)
	case "high":
		return urgencyV1.UrgencyLevel(urgencyV1.High)
	case "critical":
		return urgencyV1.UrgencyLevel(urgencyV1.Critical)
	default:
		return urgencyV1.UrgencyLevel(urgencyV1.Medium)
	}
}

func (u *Urgency) ToResponse() urgencyV1.UrgencyResponse {
	resp := urgencyV1.UrgencyResponse{
		ID:           u.ID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email,
		ContactPhone: u.ContactPhone,
		Location:     u.Location,
		Description:  u.Description,
		Level:        urgencyV1.UrgencyLevel(u.Level),
		Status:       urgencyV1.UrgencyStatus(u.Status),
		CreatedAt:    u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    u.UpdatedAt.Format(time.RFC3339),
	}
	if u.AssignedEmployeeID != nil {
		resp.AssignedEmployeeId = u.AssignedEmployeeID
	}
	if u.AssignedAt != nil {
		resp.AssignedAt = u.AssignedAt.Format(time.RFC3339)
	}
	return resp
}

func (n *Notification) ToResponse() urgencyV1.NotificationResponse {
	response := urgencyV1.NotificationResponse{
		ID:               n.ID,
		UrgencyID:        n.UrgencyID,
		EmployeeID:       n.EmployeeID,
		NotificationType: string(n.NotificationType),
		Recipient:        n.Recipient,
		Message:          n.Message,
		Status:           string(n.Status),
		Attempts:         n.Attempts,
		ErrorMessage:     n.ErrorMessage,
		CreatedAt:        n.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        n.UpdatedAt.Format(time.RFC3339),
	}

	if n.LastAttemptAt != nil {
		response.LastAttemptAt = n.LastAttemptAt.Format(time.RFC3339)
	}

	if n.SentAt != nil {
		response.SentAt = n.SentAt.Format(time.RFC3339)
	}

	return response
}

func (u *Urgency) UpdateWithRequest(req *urgencyV1.UrgencyUpdateRequest) {
	if req.FirstName != "" {
		u.FirstName = req.FirstName
	}
	if req.LastName != "" {
		u.LastName = req.LastName
	}
	if req.Email != "" {
		u.Email = req.Email
	}
	if req.ContactPhone != "" {
		u.ContactPhone = req.ContactPhone
	}
	if req.Location != "" {
		u.Location = req.Location
	}
	if req.Description != "" {
		u.Description = req.Description
	}
	if req.Level != "" {
		u.Level = urgencyV1.UrgencyLevel(req.Level)
	}
	if req.Status != "" {
		u.Status = urgencyV1.UrgencyStatus(req.Status)
	}
}
