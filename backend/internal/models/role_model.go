package models

import "time"

type Role struct {
	Id          uint `gorm:"primaryKey"`
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Users       []User `gorm:"foreignKey:RoleID"`
}
