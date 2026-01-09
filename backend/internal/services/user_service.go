package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

// UserService manages user-related operations for admins
type UserService struct {
	userRepo      *repositories.UserRepository
	roleRepo      *repositories.RoleRepository
	historyRepo   *repositories.ActionHistoryRepository
}

// NewUserService constructor
func NewUserService(userRepo *repositories.UserRepository, roleRepo *repositories.RoleRepository, historyRepo *repositories.ActionHistoryRepository) *UserService {
	return &UserService{
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		historyRepo: historyRepo,
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

// UpdateUserProfile updates user profile information (for regular users updating their own profile)
func (s *UserService) UpdateUserProfile(userID uint, firstname, lastname, email, ipAddress, userAgent string) error {
	// Validate input
	if firstname == "" || lastname == "" || email == "" {
		return errors.New("le prénom, le nom et l'email sont requis")
	}

	// Check if email is already taken by another user
	existingUser, err := s.userRepo.GetUserByEmail(email)
	if err == nil && existingUser.Id != userID {
		return errors.New("cet email est déjà utilisé")
	}

	// Get current user for history
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return fmt.Errorf("utilisateur introuvable: %w", err)
	}

	// Prepare updates
	updates := map[string]interface{}{
		"firstname": firstname,
		"lastname":  lastname,
		"email":     email,
	}

	// Update user
	if err := s.userRepo.UpdateUserById(userID, updates); err != nil {
		return err
	}

	// Log action in history
	metadata := map[string]interface{}{
		"previous_firstname": user.Firstname,
		"previous_lastname":  user.Lastname,
		"previous_email":     user.Email,
		"new_firstname":      firstname,
		"new_lastname":       lastname,
		"new_email":          email,
	}

	metadataJSON, _ := json.Marshal(metadata)

	history := &models.ActionHistory{
		UserId:       userID,
		Action:       "updated",
		ResourceType: "profile",
		ResourceId:   userID,
		Description:  fmt.Sprintf("Profil mis à jour: %s %s", firstname, lastname),
		IpAddress:    ipAddress,
		UserAgent:    userAgent,
		Metadata:     string(metadataJSON),
	}

	// Don't fail the whole operation if history logging fails
	_ = s.historyRepo.Create(history)

	return nil
}

// ChangeUserPassword changes a user's password
func (s *UserService) ChangeUserPassword(userID uint, currentPassword, newPassword, ipAddress, userAgent string) error {
	// Get user
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return fmt.Errorf("utilisateur introuvable: %w", err)
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("mot de passe actuel incorrect")
	}

	// Validate new password
	if len(newPassword) < 8 {
		return errors.New("le nouveau mot de passe doit contenir au moins 8 caractères")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("erreur lors du hachage du mot de passe: %w", err)
	}

	// Update password
	updates := map[string]interface{}{
		"password": string(hashedPassword),
	}

	if err := s.userRepo.UpdateUserById(userID, updates); err != nil {
		return err
	}

	// Log action in history
	metadata := map[string]interface{}{
		"changed_at": "profile",
	}

	metadataJSON, _ := json.Marshal(metadata)

	history := &models.ActionHistory{
		UserId:       userID,
		Action:       "updated",
		ResourceType: "password",
		ResourceId:   userID,
		Description:  "Mot de passe modifié",
		IpAddress:    ipAddress,
		UserAgent:    userAgent,
		Metadata:     string(metadataJSON),
	}

	// Don't fail the whole operation if history logging fails
	_ = s.historyRepo.Create(history)

	return nil
}
