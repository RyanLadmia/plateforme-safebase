package repositories

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"gorm.io/gorm"
)

// SessionRepository keeps a DB instance for session operations
type SessionRepository struct {
	db *gorm.DB
}

// Constructor
func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create a new session
func (r *SessionRepository) CreateSession(session *models.Session) error {
	return r.db.Create(session).Error
}

// Delete a session by token (logout)
func (r *SessionRepository) DeleteByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.Session{}).Error
}

// Find session by token
func (r *SessionRepository) FindByToken(token string) (*models.Session, error) {
	var session models.Session
	if err := r.db.Where("token = ?", token).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}
