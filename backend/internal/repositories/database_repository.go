package repositories

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"gorm.io/gorm"
)

// Keep an instance of the database connection
type DatabaseRepository struct {
	db *gorm.DB
}

// Constructor of the DatabaseRepository
func NewDatabaseRepository(db *gorm.DB) *DatabaseRepository {
	return &DatabaseRepository{db: db}
}

// Create a new database record
func (r *DatabaseRepository) Create(database *models.Database) error {
	return r.db.Create(database).Error
}

// Get database by ID
func (r *DatabaseRepository) GetByID(id uint) (*models.Database, error) {
	var database models.Database
	err := r.db.Preload("User").First(&database, id).Error
	if err != nil {
		return nil, err
	}
	return &database, nil
}

// Get all databases for a user
func (r *DatabaseRepository) GetByUserID(userID uint) ([]models.Database, error) {
	var databases []models.Database
	err := r.db.Where("user_id = ?", userID).Find(&databases).Error
	return databases, err
}

// Update database
func (r *DatabaseRepository) Update(database *models.Database) error {
	return r.db.Save(database).Error
}

// Delete database
func (r *DatabaseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Database{}, id).Error
}
