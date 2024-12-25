package model

import (
	"gorm.io/gorm"
	"time"
)

type ProfileType string

const (
	Medic         ProfileType = "Medic"
	Technical     ProfileType = "Technical"
	Administrator ProfileType = "Administrator"
)

type Employee struct {
	gorm.Model
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
	ID           uint      `gorm:"primaryKey" json:"id"`
	ShiftDate    time.Time `json:"shiftDate"`    // Represents the date of the shift
	ShiftType    int       `json:"shiftType"`    // 1: 6am-2pm, 2: 2pm-10pm, 3: 10pm-6am, < 1 or > 3: invalid
	EmployeeID   uint      `json:"employeeId"`   // Employee assigned to the shift
	EmployeeRole string    `json:"employeeRole"` // Role of the employee (e.g., Medic, Technical)
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
