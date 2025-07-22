package model

import (
	"time"

	"gorm.io/gorm"
)

type UrgencyLevel string

const (
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
	Description  string       `gorm:"not null"`
	Level        UrgencyLevel `gorm:"type:text;not null;default:'Medium'"`
	Status       string       `gorm:"type:text;not null;default:'Open'"`
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
		return "Medium"
	}
}

func (u UrgencyLevel) Valid() bool {
	return u == Low || u == Medium || u == High || u == Critical
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
		Description:  u.Description,
		Level:        u.Level.String(),
		Status:       u.Status,
		CreatedAt:    u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    u.UpdatedAt.Format(time.RFC3339),
	}
}
