package repositories

import (
	"log"
	"time"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"gorm.io/gorm"
)

// SessionRepository gère les opérations de base de données pour les sessions
type SessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository constructeur pour créer un nouveau SessionRepository
func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// CreateSession crée une nouvelle session en base de données
func (r *SessionRepository) CreateSession(session *models.Session) error {
	return r.db.Create(session).Error
}

// GetSessionByToken récupère une session par son token JWT (vérifie aussi l'expiration)
func (r *SessionRepository) GetSessionByToken(token string) (*models.Session, error) {
	var session models.Session
	if err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).
		Preload("User").
		First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// DeleteByToken supprime physiquement une session par son token (utilisé pour la déconnexion)
func (r *SessionRepository) DeleteByToken(token string) error {
	// Extrait le token du format "Bearer <token>" si nécessaire
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	// Suppression physique directe (plus de soft delete)
	result := r.db.Unscoped().Where("token = ?", token).Delete(&models.Session{})
	if result.Error != nil {
		return result.Error
	}
	log.Printf("Sessions supprimées physiquement: %d", result.RowsAffected)
	return nil
}

// DeleteByUserId supprime physiquement toutes les sessions d'un utilisateur
func (r *SessionRepository) DeleteByUserId(userId uint) error {
	// Suppression physique directe (plus de soft delete)
	result := r.db.Unscoped().Where("user_id = ?", userId).Delete(&models.Session{})
	if result.Error != nil {
		return result.Error
	}
	log.Printf("Sessions utilisateur %d supprimées physiquement: %d", userId, result.RowsAffected)
	return nil
}

// DeleteExpiredSessions supprime physiquement toutes les sessions expirées (à appeler périodiquement)
func (r *SessionRepository) DeleteExpiredSessions() error {
	// Suppression physique directe (plus de soft delete)
	result := r.db.Unscoped().Where("expires_at < ?", time.Now()).Delete(&models.Session{})
	if result.Error != nil {
		return result.Error
	}
	log.Printf("Sessions expirées supprimées physiquement: %d", result.RowsAffected)
	return nil
}

// GetActiveSessionsForUser récupère toutes les sessions actives d'un utilisateur
func (r *SessionRepository) GetActiveSessionsForUser(userId uint) ([]models.Session, error) {
	var sessions []models.Session
	if err := r.db.Where("user_id = ? AND expires_at > ?", userId, time.Now()).
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// FindByToken alias pour la compatibilité (utilise GetSessionByToken)
func (r *SessionRepository) FindByToken(token string) (*models.Session, error) {
	return r.GetSessionByToken(token)
}

// GetActiveSessionsCount retourne le nombre de sessions actives
func (r *SessionRepository) GetActiveSessionsCount() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Session{}).
		Where("expires_at > ?", time.Now()).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
