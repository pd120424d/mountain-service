package models

import "gorm.io/gorm"

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
