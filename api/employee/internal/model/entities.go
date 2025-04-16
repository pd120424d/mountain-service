package model

import (
	"time"

	"gorm.io/gorm"
)

type ProfileType string

const (
	Medic         ProfileType = "Medic"
	Technical     ProfileType = "Technical"
	Administrator ProfileType = "Administrator"
)

type Employee struct {
	gorm.Model
	ID             uint   `gorm:"primaryKey"`
	Username       string `gorm:"unique;not null"`
	Password       string `gorm:"not null"`
	FirstName      string `gorm:"not null"`
	LastName       string `gorm:"not null"`
	Gender         string `gorm:"type:char(1);not null"`
	Phone          string
	Email          string `gorm:"unique;not null"`
	ProfilePicture string
	ProfileType    ProfileType `gorm:"not null"`
}

type Shift struct {
	ID        uint      `gorm:"primaryKey"`
	ShiftDate time.Time `gorm:"not null"`
	ShiftType int       `gorm:"not null"` // 1: 6am-2pm, 2: 2pm-10pm, 3: 10pm-6am, < 1 or > 3: invalid
	CreatedAt time.Time
}

type EmployeeShift struct {
	ID          uint   `gorm:"primaryKey"`
	EmployeeID  uint   `gorm:"not null"`
	ShiftID     uint   `gorm:"not null"`
	ProfileType string `gorm:"not null"` // e.g., "Medic", "Technical"
}

func (p ProfileType) String() string {
	switch {
	case p == Medic:
		return "Medic"
	case p == Technical:
		return "Technical"
	case p == Administrator:
		return "Administrator"
	default:
		return "Unknown"
	}
}

func (p ProfileType) IsValid() bool {
	return p == Medic || p == Technical || p == Administrator
}

func ProfileTypeFromString(s string) ProfileType {
	switch s {
	case "Medic":
		return Medic
	case "Technical":
		return Technical
	case "Administrator":
		return Administrator
	default:
		return ""
	}
}

func (e *Employee) Role() string {
	return string(e.ProfileType)
}

func (e *Employee) UpdateResponseFromEmployee() EmployeeResponse {
	return EmployeeResponse{
		ID:             e.ID,
		Username:       e.Username,
		FirstName:      e.FirstName,
		LastName:       e.LastName,
		Gender:         e.Gender,
		Phone:          e.Phone,
		Email:          e.Email,
		ProfilePicture: e.ProfilePicture,
		ProfileType:    e.ProfileType.String(),
	}
}
