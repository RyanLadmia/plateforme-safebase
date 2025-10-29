package models

import (
	"time"

	"gorm.io/gorm"
)

type Restore struct {
	Id         uint           `gorm:"primaryKey" json:"id"`
	Status     string         `gorm:"size:50;not null" json:"status"` // pending, success, failed
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	UserId     uint           `gorm:"index;not null" json:"user_id"`
	User       User           `gorm:"foreignKey:UserId" json:"-"`
	BackupId   uint           `gorm:"index;not null" json:"backup_id"`
	Backup     Backup         `gorm:"foreignKey:BackupId" json:"-"`
	DatabaseId uint           `gorm:"index;not null" json:"database_id"`
	Database   Database       `gorm:"foreignKey:DatabaseId" json:"-"`
}
