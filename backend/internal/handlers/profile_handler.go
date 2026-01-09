package handlers

import (
	"net/http"

	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/gin-gonic/gin"
)

// ProfileHandler handles user profile operations
type ProfileHandler struct {
	userService *services.UserService
	authService *services.AuthService
}

// NewProfileHandler constructor
func NewProfileHandler(userService *services.UserService, authService *services.AuthService) *ProfileHandler {
	return &ProfileHandler{
		userService: userService,
		authService: authService,
	}
}

// UpdateProfileRequest represents the request body for updating profile
type UpdateProfileRequest struct {
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}

// UpdateProfile PUT /api/profile
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Extract IP address and User-Agent for logging
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	// Update profile (without role change - users can't change their own role)
	err := h.userService.UpdateUserProfile(userID.(uint), req.Firstname, req.Lastname, req.Email, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve updated user info
	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des données"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profil mis à jour avec succès",
		"user": gin.H{
			"id":        user.Id,
			"firstname": user.Firstname,
			"lastname":  user.Lastname,
			"email":     user.Email,
			"role":      user.Role,
		},
	})
}

// ChangePasswordRequest represents the request body for changing password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// ChangePassword PUT /api/profile/password
func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Extract IP address and User-Agent for logging
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	// Validate passwords match
	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Les mots de passe ne correspondent pas"})
		return
	}

	// Validate new password strength (optional - add your rules)
	if len(req.NewPassword) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le mot de passe doit contenir au moins 8 caractères"})
		return
	}

	// Change password
	err := h.userService.ChangeUserPassword(userID.(uint), req.CurrentPassword, req.NewPassword, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mot de passe changé avec succès",
	})
}

// GetProfile GET /api/profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Retrieve user info
	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur introuvable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":        user.Id,
			"firstname": user.Firstname,
			"lastname":  user.Lastname,
			"email":     user.Email,
			"role":      user.Role,
			"active":    user.Active,
		},
	})
}

