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

	// Encrypt the database URL before storing (contains sensitive connection info)
	if database.URL != "" {
		encryptedURL, err := security.EncryptDatabasePassword(database.URL)
		if err != nil {
			return fmt.Errorf("erreur lors du chiffrement de l'URL: %v", err)
		}
		database.URL = encryptedURL
	}

	return s.databaseRepo.Create(database)
}

// GetDatabasesByUser returns all databases for a user (without decrypted passwords for security)
func (s *DatabaseService) GetDatabasesByUser(userID uint) ([]models.Database, error) {
	return s.databaseRepo.GetByUserID(userID)
}

// GetDatabaseByID returns a database by ID
func (s *DatabaseService) GetDatabaseByID(id uint) (*models.Database, error) {
	fmt.Printf("[DEBUG] GetDatabaseByID: Looking for database ID %d\n", id)
	database, err := s.databaseRepo.GetByID(id)
	if err != nil {
		fmt.Printf("[DEBUG] GetDatabaseByID: Database ID %d not found - error: %v\n", id, err)
		return nil, err
	}
	fmt.Printf("[DEBUG] GetDatabaseByID: Found database ID %d, name '%s', user_id %d\n", database.Id, database.Name, database.UserId)

	// Decrypt the database password
	if database.Password != "" {
		decryptedPassword, err := security.DecryptDatabasePassword(database.Password)
		if err != nil {
			return nil, fmt.Errorf("erreur lors du déchiffrement du mot de passe: %v", err)
		}
		database.Password = decryptedPassword
	}

	// Decrypt the database URL
	if database.URL != "" {
		decryptedURL, err := security.DecryptDatabasePassword(database.URL)
		if err != nil {
			// Log the error but don't fail - URL decryption failure shouldn't prevent database access
			fmt.Printf("[DEBUG] GetDatabaseByID: Failed to decrypt URL for database ID %d - error: %v\n", database.Id, err)
			// Keep the encrypted URL as-is or set to empty to avoid issues
			database.URL = ""
		} else {
			database.URL = decryptedURL
		}
	}
	fmt.Print("ddddd", database, "ddddd")
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
		if !isMyEncryptedPassword(database.Password) {
			encryptedPassword, err := security.EncryptDatabasePassword(database.Password)
			if err != nil {
				return fmt.Errorf("erreur lors du chiffrement du mot de passe: %v", err)
			}
			database.Password = encryptedPassword
		}
	}

	// Encrypt the database URL before storing
	if database.URL != "" {
		// Check if URL is already encrypted
		if !isMyEncryptedPassword(database.URL) {
			encryptedURL, err := security.EncryptDatabasePassword(database.URL)
			if err != nil {
				return fmt.Errorf("erreur lors du chiffrement de l'URL: %v", err)
			}
			database.URL = encryptedURL
		}
	}

	return s.databaseRepo.Update(database)
}

// UpdateDatabaseName updates only the name of a database (secure partial update)
func (s *DatabaseService) UpdateDatabaseName(id uint, name string) error {
	return s.databaseRepo.UpdateDatabaseName(id, strings.TrimSpace(name))
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

// OLD FUNCTION TO CHECK IF PASSWORD IS ENCRYPTED
// isEncryptedPassword checks if a password string appears to be encrypted (base64)
//func isEncryptedPassword(password string) bool {
// Simple check: encrypted passwords are base64 encoded and contain special chars
//return len(password) > 20 && (password[len(password)-1] == '=' || strings.ContainsAny(password, "+/"))
//}

// NEW FUNCTION TO CHECK IF PASSWORD IS ENCRYPTED
func isMyEncryptedPassword(password string) bool {
	if password == "" {
		return false
	}
	_, err := security.DecryptDatabasePassword(password)
	return err == nil
}
