package models

import "time"

// Structure of the user model in the database
type User struct {
	Id        uint `gorm:"primaryKey"`
	Lastname  string
	Firstname string
	Email     string `gorm:"unique"`
	Password  string
	RoleID    uint
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}
