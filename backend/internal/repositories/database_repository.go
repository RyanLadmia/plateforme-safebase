package repositories

import (
	"fmt"

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
	fmt.Printf("[DEBUG] DatabaseRepository.GetByID: Querying database ID %d\n", id)
	var database models.Database
	err := r.db.Preload("User").First(&database, id).Error
	if err != nil {
		fmt.Printf("[DEBUG] DatabaseRepository.GetByID: Database ID %d not found in DB - error: %v\n", id, err)
		return nil, err
	}
	fmt.Printf("[DEBUG] DatabaseRepository.GetByID: Found database ID %d, name '%s'\n", database.Id, database.Name)
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

// UpdateDatabaseName updates only the name field of a database
func (r *DatabaseRepository) UpdateDatabaseName(id uint, name string) error {
	fmt.Printf("[DEBUG] DatabaseRepository.UpdateDatabaseName: Updating database ID %d with name '%s'\n", id, name)
	err := r.db.Model(&models.Database{}).Where("id = ?", id).Update("name", name).Error
	if err != nil {
		fmt.Printf("[DEBUG] DatabaseRepository.UpdateDatabaseName: Update failed - error: %v\n", err)
		return err
	}
	fmt.Printf("[DEBUG] DatabaseRepository.UpdateDatabaseName: Update succeeded for database ID %d\n", id)
	return nil
}

// Delete database
func (r *DatabaseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Database{}, id).Error
}

// Soft delete database
func (r *DatabaseRepository) SoftDelete(id uint) error {
	return r.db.Delete(&models.Database{}, id).Error // GORM automatically does soft delete with DeletedAt field
}
