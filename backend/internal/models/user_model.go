package models

import (
	"time"

	"gorm.io/gorm"
)

// Structure of the user model in the database
type User struct {
	Id        uint           `gorm:"primaryKey" json:"id"`
	Lastname  string         `gorm:"size:100;not null" json:"lastname"`
	Firstname string         `gorm:"size:100;not null" json:"firstname"`
	Email     string         `gorm:"uniqueIndex;size:255;not null;" json:"email"`
	Password  string         `gorm:"type:text;not null" json:"-"` // Not shown in JSON response
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	RoleID    *uint          `gorm:"index;default:2" json:"role_id"`
	Role      *Role          `gorm:"foreignKey:RoleID" json:"role,omitempty"` // Not shown in JSON response if RoleID is null
	Backups   []Backup       `gorm:"constraint:OnDelete:CASCADE;" json:"backups,omitempty"`
	Restores  []Restore      `gorm:"constraint:OnDelete:CASCADE;" json:"restores,omitempty"`
	Sessions  []Session      `gorm:"constraint:OnDelete:CASCADE;" json:"sessions,omitempty"`
	Alerts    []Alert        `gorm:"constraint:OnDelete:CASCADE;" json:"alerts,omitempty"`
}
