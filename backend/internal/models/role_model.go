package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	Id        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;uniqueIndex;not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // if deleted, User.RoleID will be set to null
	Users     []User         `gorm:"constraint:OnDelete:SET NULL;" json:"users,omitempty"`
}
