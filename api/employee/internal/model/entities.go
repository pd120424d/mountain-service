package model

import (
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
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
	ProfileType    ProfileType `gorm:"type:text;not null"`
	Shifts         []Shift     `gorm:"many2many:employee_shifts;"`
}

type Shift struct {
	ID        uint      `gorm:"primaryKey"`
	ShiftDate time.Time `gorm:"not null"`
	ShiftType int       `gorm:"not null"` // 1: 6am-2pm, 2: 2pm-10pm, 3: 10pm-6am, < 1 or > 3: invalid
	CreatedAt time.Time
	Employees []Employee `gorm:"many2many:employee_shifts;"`
}

type EmployeeShift struct {
	ID         uint `gorm:"primaryKey"`
	EmployeeID uint `gorm:"not null"`
	ShiftID    uint `gorm:"not null"`
}

type ShiftsAvailabilityRange struct {
	Days map[time.Time][]map[ProfileType]int `json:"days"`
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

func (p ProfileType) Valid() bool {
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

// Methods for BasicAuth interface compatibility
func (e *Employee) GetID() uint {
	return e.ID
}

func (e *Employee) GetPassword() string {
	return e.Password
}

func (e *Employee) GetRole() string {
	return e.Role()
}

func (e *Employee) UpdateResponseFromEmployee() employeeV1.EmployeeResponse {
	return employeeV1.EmployeeResponse{
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
