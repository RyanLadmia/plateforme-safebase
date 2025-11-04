package models

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	Id             uint           `gorm:"primaryKey" json:"id"`
	CronExpression string         `gorm:"size:100;not null" json:"cron_expression"`
	Active         bool           `gorm:"default:true" json:"active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	UserId         uint           `gorm:"index;not null" json:"user_id"`
	User           User           `gorm:"foreignKey:UserId" json:"-"`
	DatabaseId     uint           `gorm:"index;not null" json:"database_id"`
	Database       Database       `gorm:"foreignKey:DatabaseId" json:"database"`
}
