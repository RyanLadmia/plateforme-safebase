package services

import (
	"fmt"
	"strings"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/RyanLadmia/plateforme-safebase/pkg/security"
)

type DatabaseService struct {
	databaseRepo         *repositories.DatabaseRepository
	backupRepo           *repositories.BackupRepository
	restoreRepo          *repositories.RestoreRepository
	scheduleRepo         *repositories.ScheduleRepository
	backupService        *BackupService
	actionHistoryService *ActionHistoryService
}

// Constructor for DatabaseService
func NewDatabaseService(databaseRepo *repositories.DatabaseRepository, backupRepo *repositories.BackupRepository, restoreRepo *repositories.RestoreRepository, scheduleRepo *repositories.ScheduleRepository, backupService *BackupService) *DatabaseService {
	return &DatabaseService{
		databaseRepo:  databaseRepo,
		backupRepo:    backupRepo,
		restoreRepo:   restoreRepo,
		scheduleRepo:  scheduleRepo,
		backupService: backupService,
	}
}

// SetBackupService sets the backup service reference for cascade deletion
func (s *DatabaseService) SetBackupService(backupService *BackupService) {
	s.backupService = backupService
}

// SetActionHistoryService sets the action history service reference for logging
func (s *DatabaseService) SetActionHistoryService(actionHistoryService *ActionHistoryService) {
	s.actionHistoryService = actionHistoryService
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

// LogDatabaseAction logs a database action if action history service is available
func (s *DatabaseService) logDatabaseAction(userID uint, action string, databaseID uint, databaseName string, ipAddress, userAgent string) {
	if s.actionHistoryService != nil {
		err := s.actionHistoryService.LogDatabaseAction(userID, action, databaseID, databaseName, ipAddress, userAgent)
		if err != nil {
			fmt.Printf("[HISTORY] Failed to log database action %s for database %d: %v\n", action, databaseID, err)
		}
	}
}

// GetDatabasesByUser returns all databases for a user (without decrypted passwords for security)
func (s *DatabaseService) GetDatabasesByUser(userID uint) ([]models.Database, error) {
	return s.databaseRepo.GetByUserID(userID)
}

// GetBackupsByDatabase returns all backups for a database
func (s *DatabaseService) GetBackupsByDatabase(databaseID uint) ([]models.Backup, error) {
	return s.backupRepo.GetByDatabaseID(databaseID)
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

// CreateDatabaseWithLogging creates a new database record with action logging
func (s *DatabaseService) CreateDatabaseWithLogging(database *models.Database, userID uint, ipAddress, userAgent string) error {
	err := s.CreateDatabase(database)
	if err != nil {
		return err
	}

	// Log the action
	s.logDatabaseAction(userID, "created", database.Id, database.Name, ipAddress, userAgent)
	return nil
}

// UpdateDatabaseWithLogging updates a database record with action logging
func (s *DatabaseService) UpdateDatabaseWithLogging(database *models.Database, userID uint, ipAddress, userAgent string) error {
	err := s.UpdateDatabase(database)
	if err != nil {
		return err
	}

	// Log the action
	s.logDatabaseAction(userID, "updated", database.Id, database.Name, ipAddress, userAgent)
	return nil
}

// UpdateDatabaseNameWithLogging updates only the name of a database with action logging
func (s *DatabaseService) UpdateDatabaseNameWithLogging(id uint, name string, userID uint, ipAddress, userAgent string) error {
	err := s.UpdateDatabaseName(id, name)
	if err != nil {
		return err
	}

	// Log the action
	s.logDatabaseAction(userID, "updated", id, name, ipAddress, userAgent)
	return nil
}

// DeleteDatabaseWithLogging soft deletes a database record with action logging
func (s *DatabaseService) DeleteDatabaseWithLogging(id uint, userID uint, ipAddress, userAgent string) error {
	// Get the database name for logging before deletion
	database, err := s.databaseRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("base de données introuvable: %v", err)
	}

	err = s.DeleteDatabase(id, userID)
	if err != nil {
		return err
	}

	// Log the action
	s.logDatabaseAction(userID, "deleted", id, database.Name, ipAddress, userAgent)
	return nil
}

// DeleteDatabase soft deletes a database record and all associated records (backups, schedules, restores)
func (s *DatabaseService) DeleteDatabase(id uint, userID uint) error {
	database, err := s.databaseRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("base de données introuvable: %v", err)
	}

	// Verify user ownership
	if database.UserId != userID {
		return fmt.Errorf("accès non autorisé à cette base de données")
	}

	fmt.Printf("[DELETE] Soft deleting database %d and all associated records\n", id)

	// Soft delete all associated backups (keep files in Mega for storage optimization)
	if err := s.backupRepo.SoftDeleteByDatabaseID(id); err != nil {
		fmt.Printf("[DELETE] Warning: failed to soft delete backups for database %d: %v\n", id, err)
		// Continue even if soft delete fails
	} else {
		fmt.Printf("[DELETE] Soft deleted backups for database %d\n", id)
	}

	// Soft delete all associated schedules
	if err := s.scheduleRepo.SoftDeleteByDatabaseID(id); err != nil {
		fmt.Printf("[DELETE] Warning: failed to soft delete schedules for database %d: %v\n", id, err)
		// Continue even if soft delete fails
	} else {
		fmt.Printf("[DELETE] Soft deleted schedules for database %d\n", id)
	}

	// Soft delete all associated restores
	if err := s.restoreRepo.SoftDeleteByDatabaseID(id); err != nil {
		fmt.Printf("[DELETE] Warning: failed to soft delete restores for database %d: %v\n", id, err)
		// Continue even if soft delete fails
	} else {
		fmt.Printf("[DELETE] Soft deleted restores for database %d\n", id)
	}

	// Finally, soft delete the database record
	fmt.Printf("[DELETE] Soft deleting database record ID %d\n", id)
	return s.databaseRepo.SoftDelete(id)
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
