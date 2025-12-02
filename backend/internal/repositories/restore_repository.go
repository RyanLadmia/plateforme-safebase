package repositories

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"gorm.io/gorm"
)

// Keep an instance of the database connection
type RestoreRepository struct {
	db *gorm.DB
}

// Constructor of the RestoreRepository, can be used to create a new RestoreRepository
func NewRestoreRepository(db *gorm.DB) *RestoreRepository {
	return &RestoreRepository{db: db}
}

// Create a new restore record
func (r *RestoreRepository) Create(restore *models.Restore) error {
	return r.db.Create(restore).Error
}

// Get restore by ID
func (r *RestoreRepository) GetByID(id uint) (*models.Restore, error) {
	var restore models.Restore
	err := r.db.Preload("Database").Preload("User").Preload("Backup").First(&restore, id).Error
	if err != nil {
		return nil, err
	}
	return &restore, nil
}

// Get all restores for a user
func (r *RestoreRepository) GetByUserID(userID uint) ([]models.Restore, error) {
	var restores []models.Restore
	err := r.db.Preload("Database").Preload("User").Preload("Backup").Where("user_id = ?", userID).Find(&restores).Error
	return restores, err
}

// Get all restores for a database
func (r *RestoreRepository) GetByDatabaseID(databaseID uint) ([]models.Restore, error) {
	var restores []models.Restore
	err := r.db.Preload("Database").Preload("User").Preload("Backup").Where("database_id = ?", databaseID).Find(&restores).Error
	return restores, err
}

// Get all restores for a backup
func (r *RestoreRepository) GetByBackupID(backupID uint) ([]models.Restore, error) {
	var restores []models.Restore
	err := r.db.Preload("Database").Preload("User").Preload("Backup").Where("backup_id = ?", backupID).Find(&restores).Error
	return restores, err
}

// Update restore status
func (r *RestoreRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Restore{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": status,
	}).Error
}

// GetDB returns the database connection (for advanced queries)
func (r *RestoreRepository) GetDB() *gorm.DB {
	return r.db
}

// Delete restore
func (r *RestoreRepository) Delete(id uint) error {
	return r.db.Delete(&models.Restore{}, id).Error
}
