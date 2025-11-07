package repositories

import (
	"time"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"gorm.io/gorm"
)

// Keep an instance of the database connection
type BackupRepository struct {
	db *gorm.DB
}

// Constructor of the BackupRepository, can be used to create a new BackupRepository
func NewBackupRepository(db *gorm.DB) *BackupRepository {
	return &BackupRepository{db: db}
}

// Create a new backup record
func (r *BackupRepository) Create(backup *models.Backup) error {
	return r.db.Create(backup).Error
}

// Get backup by ID
func (r *BackupRepository) GetByID(id uint) (*models.Backup, error) {
	var backup models.Backup
	err := r.db.Preload("Database").Preload("User").First(&backup, id).Error
	if err != nil {
		return nil, err
	}
	return &backup, nil
}

// Get all backups for a user
func (r *BackupRepository) GetByUserID(userID uint) ([]models.Backup, error) {
	var backups []models.Backup
	err := r.db.Preload("Database").Where("user_id = ?", userID).Find(&backups).Error
	return backups, err
}

// Get all backups for a database
func (r *BackupRepository) GetByDatabaseID(databaseID uint) ([]models.Backup, error) {
	var backups []models.Backup
	err := r.db.Preload("Database").Where("database_id = ?", databaseID).Find(&backups).Error
	return backups, err
}

// Update backup status
func (r *BackupRepository) UpdateStatus(id uint, status string, errorMsg string) error {
	return r.db.Model(&models.Backup{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    status,
		"error_msg": errorMsg,
	}).Error
}

// Update backup filepath and size
func (r *BackupRepository) UpdateFileInfo(id uint, filepath string, size int64) error {
	return r.db.Model(&models.Backup{}).Where("id = ?", id).Updates(map[string]interface{}{
		"filepath": filepath,
		"size":     size,
	}).Error
}

// Get old backups for cleanup (older than specified days)
func (r *BackupRepository) GetOldBackups(days int) ([]models.Backup, error) {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	var backups []models.Backup
	err := r.db.Where("created_at < ? AND status = ?", cutoffDate, "completed").Find(&backups).Error
	return backups, err
}

// GetDB returns the database connection (for advanced queries)
func (r *BackupRepository) GetDB() *gorm.DB {
	return r.db
}

// Delete backup
func (r *BackupRepository) Delete(id uint) error {
	return r.db.Delete(&models.Backup{}, id).Error
}
