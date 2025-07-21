package model

import "time"

type Urgency struct {
	ID           uint      `gorm:"primaryKey"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"not null"`
	ContactPhone string    `gorm:"not null"`
	Description  string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}
