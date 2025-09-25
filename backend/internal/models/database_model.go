package models

import (
	"time"

	"gorm.io/gorm"
)

type Database struct {
	Id        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Type      string         `gorm:"size:50;not null" json:"type"` // mysql, postgres
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserId    uint           `gorm:"index;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserId" json:"-"`
	Backups   []Backup       `gorm:"constraint:OnDelete:CASCADE;" json:"backups,omitempty"`
	Restores  []Restore      `gorm:"constraint:OnDelete:CASCADE;" json:"restores,omitempty"`
}
