package repositories

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"gorm.io/gorm"
)

// Keep an instance of the database connection
type UserRepository struct {
	db *gorm.DB
}

// Constructor of the UserRepository, can be used to create a new UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetDB return the database connection (for complex queries)
func (r *UserRepository) GetDB() *gorm.DB {
	return r.db
}

// Create a new user
func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// Get all users
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := r.db.Preload("Role").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Get all active users
func (r *UserRepository) GetAllActiveUsers() ([]models.User, error) {
	var users []models.User
	if err := r.db.Preload("Role").Where("active = ?", true).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Deactivate user (soft delete alternative)
func (r *UserRepository) DeactivateUser(id uint) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("active", false).
		Error
}

// Activate user
func (r *UserRepository) ActivateUser(id uint) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("active", true).
		Error
}

// Hard delete user (use with caution)
func (r *UserRepository) DeleteUser(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// Get user by id
func (r *UserRepository) GetUserById(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Get user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).Preload("Role").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateU user infos
func (r *UserRepository) UpdateUserById(id uint, updates map[string]interface{}) error {
	// We use Model(&User{}) to specify the table, Where to select the user
	// Updates only affect fields present in the map
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Updates(updates).
		Error
}
