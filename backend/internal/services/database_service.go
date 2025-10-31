package services

import (
	"fmt"
	"strings"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/RyanLadmia/plateforme-safebase/pkg/security"
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

	// Encrypt the database password before storing
	if database.Password != "" {
		encryptedPassword, err := security.EncryptDatabasePassword(database.Password)
		if err != nil {
			return fmt.Errorf("erreur lors du chiffrement du mot de passe: %v", err)
		}
		database.Password = encryptedPassword
	}

	return s.databaseRepo.Create(database)
}

// GetDatabasesByUser returns all databases for a user (without decrypted passwords for security)
func (s *DatabaseService) GetDatabasesByUser(userID uint) ([]models.Database, error) {
	return s.databaseRepo.GetByUserID(userID)
}

// GetDatabaseByID returns a database by ID
func (s *DatabaseService) GetDatabaseByID(id uint) (*models.Database, error) {
	database, err := s.databaseRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Decrypt the database password
	if database.Password != "" {
		decryptedPassword, err := security.DecryptDatabasePassword(database.Password)
		if err != nil {
			return nil, fmt.Errorf("erreur lors du déchiffrement du mot de passe: %v", err)
		}
		database.Password = decryptedPassword
	}

	return database, nil
}

// GetDatabaseByIDForBackup returns a database by ID with decrypted password for backup operations
func (s *DatabaseService) GetDatabaseByIDForBackup(id uint) (*models.Database, error) {
	return s.GetDatabaseByID(id) // Already decrypts password
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

	// Encrypt the database password before storing
	if database.Password != "" {
		// Check if password is already encrypted (contains base64 characters)
		if !isEncryptedPassword(database.Password) {
			encryptedPassword, err := security.EncryptDatabasePassword(database.Password)
			if err != nil {
				return fmt.Errorf("erreur lors du chiffrement du mot de passe: %v", err)
			}
			database.Password = encryptedPassword
		}
	}

	return s.databaseRepo.Update(database)
}

// isEncryptedPassword checks if a password string appears to be encrypted (base64)
func isEncryptedPassword(password string) bool {
	// Simple check: encrypted passwords are base64 encoded and contain special chars
	return len(password) > 20 && (password[len(password)-1] == '=' || strings.ContainsAny(password, "+/"))
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
