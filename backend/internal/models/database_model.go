package models

import (
	"time"

	"gorm.io/gorm"
)

// DatabaseCreateRequest is used for JSON binding during creation
type DatabaseCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=mysql postgresql"`
	Host     string `json:"host" binding:"required"`
	Port     string `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	DbName   string `json:"db_name" binding:"required"`
}

// DatabaseUpdateRequest is used for JSON binding during updates
type DatabaseUpdateRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=mysql postgresql"`
	Host     string `json:"host" binding:"required"`
	Port     string `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password"` // Optional for updates
	DbName   string `json:"db_name" binding:"required"`
}

type Database struct {
	Id        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Type      string         `gorm:"size:50;not null" json:"type"` // mysql, postgres
	Host      string         `gorm:"size:255;not null" json:"host"`
	Port      string         `gorm:"size:10;not null" json:"port"`
	Username  string         `gorm:"size:100;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"` // Ne pas exposer le mot de passe dans le JSON
	DbName    string         `gorm:"size:100;not null" json:"db_name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserId    uint           `gorm:"index;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserId" json:"-"`
	Backups   []Backup       `gorm:"foreignKey:DatabaseId;constraint:OnDelete:CASCADE;" json:"backups,omitempty"`
	Restores  []Restore      `gorm:"foreignKey:DatabaseId;constraint:OnDelete:CASCADE;" json:"restores,omitempty"`
}
