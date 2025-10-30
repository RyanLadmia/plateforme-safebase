package services

import (
	"fmt"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
)

type DatabaseService struct {
	databaseRepo *repositories.DatabaseRepository
}

// Constructor for DatabaseService
func NewDatabaseService(databaseRepo *repositories.DatabaseRepository) *DatabaseService {
	return &DatabaseService{
		databaseRepo: databaseRepo,
	}
}

// CreateDatabase creates a new database record
func (s *DatabaseService) CreateDatabase(database *models.Database) error {
	// Validate database type
	if database.Type != "mysql" && database.Type != "postgres" && database.Type != "postgresql" {
		return fmt.Errorf("type de base de données non supporté: %s", database.Type)
	}

	// Normalize postgres type
	if database.Type == "postgres" {
		database.Type = "postgresql"
	}

	return s.databaseRepo.Create(database)
}

// GetDatabasesByUser returns all databases for a user
func (s *DatabaseService) GetDatabasesByUser(userID uint) ([]models.Database, error) {
	return s.databaseRepo.GetByUserID(userID)
}

// GetDatabaseByID returns a database by ID
func (s *DatabaseService) GetDatabaseByID(id uint) (*models.Database, error) {
	return s.databaseRepo.GetByID(id)
}

// UpdateDatabase updates a database record
func (s *DatabaseService) UpdateDatabase(database *models.Database) error {
	// Validate database type
	if database.Type != "mysql" && database.Type != "postgres" && database.Type != "postgresql" {
		return fmt.Errorf("type de base de données non supporté: %s", database.Type)
	}

	// Normalize postgres type
	if database.Type == "postgres" {
		database.Type = "postgresql"
	}

	return s.databaseRepo.Update(database)
}

// DeleteDatabase deletes a database record
func (s *DatabaseService) DeleteDatabase(id uint, userID uint) error {
	database, err := s.databaseRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("base de données introuvable: %v", err)
	}

	// Verify user ownership
	if database.UserId != userID {
		return fmt.Errorf("accès non autorisé à cette base de données")
	}

	return s.databaseRepo.Delete(id)
}
