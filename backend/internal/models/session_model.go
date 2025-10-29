package models

import (
	"time"

	"gorm.io/gorm"
)

// Session représente une session utilisateur avec token JWT et fonctionnalités futures
type Session struct {
	Id               uint           `gorm:"primaryKey" json:"id"`
	Token            string         `gorm:"type:text;not null;uniqueIndex" json:"-"` // Token JWT principal
	ExpiresAt        time.Time      `gorm:"not null" json:"expires_at"`
	RefreshToken     string         `gorm:"type:text" json:"-"` // Pour refresh automatique (futur)
	RefreshExpiresAt time.Time      `json:"refresh_expires_at"` // Expiration refresh token
	ResetToken       string         `gorm:"type:text" json:"-"` // Pour reset password (futur)
	ResetExpiresAt   time.Time      `json:"reset_expires_at"`   // Expiration reset token
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	UserId           uint           `gorm:"index;not null" json:"user_id"`
	User             User           `gorm:"foreignKey:UserId" json:"-"`
}
