package repositories

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"gorm.io/gorm"
)

// ActionHistoryRepository manages action history records
type ActionHistoryRepository struct {
	db *gorm.DB
}

// NewActionHistoryRepository creates a new ActionHistoryRepository
func NewActionHistoryRepository(db *gorm.DB) *ActionHistoryRepository {
	return &ActionHistoryRepository{db: db}
}

// Create creates a new action history record
func (r *ActionHistoryRepository) Create(actionHistory *models.ActionHistory) error {
	return r.db.Create(actionHistory).Error
}

// GetByID gets an action history record by ID
func (r *ActionHistoryRepository) GetByID(id uint) (*models.ActionHistory, error) {
	var actionHistory models.ActionHistory
	err := r.db.Preload("User").First(&actionHistory, id).Error
	if err != nil {
		return nil, err
	}
	return &actionHistory, nil
}

// GetByUserID gets action history for a user with pagination
func (r *ActionHistoryRepository) GetByUserID(userID uint, limit int, offset int) ([]models.ActionHistory, error) {
	var actionHistories []models.ActionHistory
	err := r.db.Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&actionHistories).Error
	return actionHistories, err
}

// GetByUserIDCount gets the total count of action history records for a user
func (r *ActionHistoryRepository) GetByUserIDCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.ActionHistory{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// GetByResource gets action history for a specific resource
func (r *ActionHistoryRepository) GetByResource(resourceType string, resourceID uint) ([]models.ActionHistory, error) {
	var actionHistories []models.ActionHistory
	err := r.db.Preload("User").
		Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).
		Order("created_at DESC").
		Find(&actionHistories).Error
	return actionHistories, err
}

// GetByResourceType gets action history for a resource type with pagination
func (r *ActionHistoryRepository) GetByResourceType(resourceType string, limit int, offset int) ([]models.ActionHistory, error) {
	var actionHistories []models.ActionHistory
	err := r.db.Preload("User").
		Where("resource_type = ?", resourceType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&actionHistories).Error
	return actionHistories, err
}

// GetByResourceTypeCount gets the total count of action history records for a resource type
func (r *ActionHistoryRepository) GetByResourceTypeCount(resourceType string) (int64, error) {
	var count int64
	err := r.db.Model(&models.ActionHistory{}).
		Where("resource_type = ?", resourceType).
		Count(&count).Error
	return count, err
}

// GetByUserIDAndResourceType gets action history for a user and resource type with pagination
func (r *ActionHistoryRepository) GetByUserIDAndResourceType(userID uint, resourceType string, limit int, offset int) ([]models.ActionHistory, error) {
	var actionHistories []models.ActionHistory
	err := r.db.Preload("User").
		Where("user_id = ? AND resource_type = ?", userID, resourceType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&actionHistories).Error
	return actionHistories, err
}

// GetByUserIDAndResourceTypeCount gets the total count of action history records for a user and resource type
func (r *ActionHistoryRepository) GetByUserIDAndResourceTypeCount(userID uint, resourceType string) (int64, error) {
	var count int64
	err := r.db.Model(&models.ActionHistory{}).
		Where("user_id = ? AND resource_type = ?", userID, resourceType).
		Count(&count).Error
	return count, err
}

// GetRecent gets recent action history records with pagination
func (r *ActionHistoryRepository) GetRecent(limit int, offset int) ([]models.ActionHistory, error) {
	var actionHistories []models.ActionHistory
	err := r.db.Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&actionHistories).Error
	return actionHistories, err
}

// GetRecentCount gets the total count of all action history records
func (r *ActionHistoryRepository) GetRecentCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.ActionHistory{}).Count(&count).Error
	return count, err
}

// GetDB returns the database connection
func (r *ActionHistoryRepository) GetDB() *gorm.DB {
	return r.db
}
