package models

import "time"

// Structure of the user model in the database
type User struct {
	Id        uint `gorm:"primaryKey"`
	Lastname  string
	Fristname string
	Email     string `gorm:"unique"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
