package repositories

import (
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

// Update backup size
func (r *BackupRepository) UpdateSize(id uint, size int64) error {
	return r.db.Model(&models.Backup{}).Where("id = ?", id).Update("size", size).Error
}

// Delete backup
func (r *BackupRepository) Delete(id uint) error {
	return r.db.Delete(&models.Backup{}, id).Error
}
