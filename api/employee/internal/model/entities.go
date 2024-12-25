package model

import (
	"gorm.io/gorm"
	"time"
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
	ProfileType    string `gorm:"not null"`
}

type Shift struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ShiftDate    time.Time `json:"shiftDate"`    // Represents the date of the shift
	ShiftType    int       `json:"shiftType"`    // 1: 6am-2pm, 2: 2pm-10pm, 3: 10pm-6am, < 1 or > 3: invalid
	EmployeeID   uint      `json:"employeeId"`   // Employee assigned to the shift
	EmployeeRole string    `json:"employeeRole"` // Role of the employee (e.g., Medic, Technical)
}
