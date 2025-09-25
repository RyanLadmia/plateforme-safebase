package models

import (
	"time"

	"gorm.io/gorm"
)

type Alert struct {
	Id        uint           `gorm:"primaryKey" json:"id"`
	Message   string         `gorm:"type:text;not null" json:"message"`
	Read      bool           `gorm:"default:false" json:"read"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserId    uint           `gorm:"index;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserId" json:"-"`
}
