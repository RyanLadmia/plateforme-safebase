package handlers

import (
	"net/http"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/gin-gonic/gin"
)

// AuthHandler contient le service d'authentification
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler constructeur
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register endpoint: POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var user models.User
	// Bind JSON body to user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Call service to register user
	if err := h.authService.Register(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Respond with created user info (without password)
	c.JSON(http.StatusCreated, gin.H{
		"id":        user.Id,
		"firstname": user.Firstname,
		"lastname":  user.Lastname,
		"email":     user.Email,
		"role_id":   user.RoleID,
	})
}

// Login endpoint: POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Call service to login
	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Respond with JWT token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Logout endpoint: POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Token peut être passé dans header Authorization: Bearer <token>
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no token provided"})
		return
	}

	if err := h.authService.Logout(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
