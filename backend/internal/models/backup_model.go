package models

import (
	"time"

	"gorm.io/gorm"
)

type Backup struct {
	Id         uint           `gorm:"primaryKey" json:"id"`
	Filename   string         `gorm:"size:255;not null" json:"filename"`
	Filepath   string         `gorm:"type:text;not null" json:"filepath"`
	Size       int64          `gorm:"not null;default:0" json:"size"`                   // Taille en bytes
	Status     string         `gorm:"size:50;not null;default:'pending'" json:"status"` // pending, completed, failed
	ErrorMsg   string         `gorm:"type:text" json:"error_msg,omitempty"`
	UserAgent  string         `gorm:"size:255" json:"user_agent,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserId     uint           `gorm:"index;not null" json:"user_id"`
	User       User           `gorm:"foreignKey:UserId" json:"-"`
	DatabaseId uint           `gorm:"index;not null" json:"database_id"`
	Database   Database       `gorm:"foreignKey:DatabaseId" json:"-"`
	Restores   []Restore      `gorm:"foreignKey:BackupId;constraint:OnDelete:CASCADE;" json:"restores,omitempty"`
}
