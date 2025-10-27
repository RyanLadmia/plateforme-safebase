package repositories

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role
	if err := r.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepository) GetByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *UserRepository) UpdateUserRole(userID uint, newRoleID uint) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("role_id", newRoleID).
		Error
}

func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}
