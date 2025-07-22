package model

import (
	"time"

	"gorm.io/gorm"
)

type (
	Status       string
	UrgencyLevel string
)

const (
	Open       Status = "Open"
	InProgress Status = "In Progress"
	Resolved   Status = "Resolved"
	Closed     Status = "Closed"

	Low      UrgencyLevel = "Low"
	Medium   UrgencyLevel = "Medium"
	High     UrgencyLevel = "High"
	Critical UrgencyLevel = "Critical"
)

type Urgency struct {
	gorm.Model
	ID           uint         `gorm:"primaryKey"`
	Name         string       `gorm:"not null"`
	Email        string       `gorm:"not null"`
	ContactPhone string       `gorm:"not null"`
	Location     string       `gorm:"not null"`
	Description  string       `gorm:"not null"`
	Level        UrgencyLevel `gorm:"type:text;not null;default:'Medium'"`
	Status       Status       `gorm:"type:text;not null;default:'Open'"`
}

func (u UrgencyLevel) String() string {
	switch u {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	case Critical:
		return "Critical"
	default:
		return string(u)
	}
}

func (u UrgencyLevel) Valid() bool {
	for _, v := range []UrgencyLevel{Low, Medium, High, Critical} {
		if v == u {
			return true
		}
	}
	return false
}

func UrgencyLevelFromString(s string) UrgencyLevel {
	switch s {
	case "Low":
		return Low
	case "Medium":
		return Medium
	case "High":
		return High
	case "Critical":
		return Critical
	default:
		return Medium
	}
}

func (u *Urgency) ToResponse() UrgencyResponse {
	return UrgencyResponse{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		ContactPhone: u.ContactPhone,
		Location:     u.Location,
		Description:  u.Description,
		Level:        u.Level,
		Status:       u.Status,
		CreatedAt:    u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    u.UpdatedAt.Format(time.RFC3339),
	}
}

func (s Status) Valid() bool {
	for _, v := range []Status{Open, InProgress, Resolved, Closed} {
		if v == s {
			return true
		}
	}
	return false
}
