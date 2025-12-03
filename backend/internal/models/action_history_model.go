package models

import (
	"time"
)

// ActionHistory tracks user actions for audit and history purposes
type ActionHistory struct {
	Id           uint      `gorm:"primaryKey" json:"id"`
	UserId       uint      `gorm:"index;not null" json:"user_id"`
	User         User      `gorm:"foreignKey:UserId" json:"user,omitempty"`
	Action       string    `gorm:"size:50;not null" json:"action"`        // created, updated, deleted, completed, failed, executed, downloaded
	ResourceType string    `gorm:"size:50;not null" json:"resource_type"` // database, backup, schedule, restore
	ResourceId   uint      `gorm:"not null" json:"resource_id"`
	Description  string    `gorm:"type:text;not null" json:"description"`
	Metadata     string    `gorm:"type:jsonb" json:"metadata,omitempty"` // JSON field for additional data
	IpAddress    string    `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent    string    `gorm:"type:text" json:"user_agent,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName specifies the table name for ActionHistory
func (ActionHistory) TableName() string {
	return "action_histories"
}
