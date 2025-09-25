package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	Id               uint           `gorm:"primaryKey" json:"id"`
	Token            string         `gorm:"type:text;not null" json:"-"`
	ExpiresAt        time.Time      `gorm:"not null" json:"expires_at"`
	RefreshToken     string         `gorm:"type:text;not null" json:"-"`
	RefreshExpiresAt time.Time      `gorm:"not null" json:"refresh_expires_at"`
	ResetToken       string         `gorm:"type:text;not null" json:"-"`
	ResetExpiresAt   time.Time      `gorm:"not null" json:"reset_expires_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	UserId           uint           `gorm:"index;not null" json:"user_id"`
	User             User           `gorm:"foreignKey:UserId" json:"-"`
}
