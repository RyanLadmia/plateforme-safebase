package services

import (
	"errors"
	"fmt"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
)

// UserService manages user-related operations for admins
type UserService struct {
	userRepo *repositories.UserRepository
	roleRepo *repositories.RoleRepository
}

// NewUserService constructor
func NewUserService(userRepo *repositories.UserRepository, roleRepo *repositories.RoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// GetAllUsers retrieves all users with their roles
func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAllUsers()
}

// GetAllActiveUsers retrieves only active users
func (s *UserService) GetAllActiveUsers() ([]models.User, error) {
	return s.userRepo.GetAllActiveUsers()
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetUserById(id)
}

// UpdateUser updates user information (firstname, lastname, email, role)
func (s *UserService) UpdateUser(id uint, firstname, lastname, email string, roleID *uint) error {
	// Validate input
	if firstname == "" || lastname == "" || email == "" {
		return errors.New("firstname, lastname, and email are required")
	}

	// Check if email is already taken by another user
	existingUser, err := s.userRepo.GetUserByEmail(email)
	if err == nil && existingUser.Id != id {
		return errors.New("email already exists")
	}

	// If roleID is provided, validate it exists
	if roleID != nil {
		_, err := s.roleRepo.GetByID(*roleID)
		if err != nil {
			return fmt.Errorf("invalid role ID: %w", err)
		}
	}

	// Prepare updates
	updates := map[string]interface{}{
		"firstname": firstname,
		"lastname":  lastname,
		"email":     email,
	}

	if roleID != nil {
		updates["role_id"] = *roleID
	}

	return s.userRepo.UpdateUserById(id, updates)
}

// DeactivateUser deactivates a user account
func (s *UserService) DeactivateUser(id uint) error {
	return s.userRepo.DeactivateUser(id)
}

// ActivateUser activates a user account
func (s *UserService) ActivateUser(id uint) error {
	return s.userRepo.ActivateUser(id)
}

// ChangeUserRole changes the role of a user
func (s *UserService) ChangeUserRole(userID, roleID uint) error {
	// Validate role exists
	_, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return fmt.Errorf("invalid role ID: %w", err)
	}

	updates := map[string]interface{}{
		"role_id": roleID,
	}

	return s.userRepo.UpdateUserById(userID, updates)
}
