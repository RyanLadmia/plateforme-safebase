package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

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

// Structure pour la requête d'inscription (inclut le mot de passe)
type RegisterRequest struct {
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

// Register endpoint: POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	// Bind JSON body to request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Erreur de binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Log des données reçues (sans le mot de passe pour la sécurité)
	log.Printf("Tentative d'inscription pour: %s %s (%s)", req.Firstname, req.Lastname, req.Email)
	log.Printf("Longueur du mot de passe reçu: %d caractères", len(req.Password))

	// Créer le modèle User à partir de la requête
	user := models.User{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Email:     req.Email,
		Password:  req.Password,
	}

	// Call service to register user
	if err := h.authService.Register(&user); err != nil {
		log.Printf("Erreur lors de l'inscription: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Utilisateur créé avec succès: ID %d", user.Id)

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

	// Définir le cookie HTTP-only sécurisé
	// En production, le cookie sera sécurisé avec HTTPS
	isProduction := os.Getenv("GO_ENV") == "production"
	c.SetCookie(
		"auth_token",                // nom du cookie
		token,                       // valeur (JWT token)
		int(24*time.Hour.Seconds()), // maxAge en secondes (24h)
		"/",                         // path
		"",                          // domain (vide = domaine actuel)
		isProduction,                // secure (true en production avec HTTPS)
		true,                        // httpOnly (empêche l'accès via JavaScript)
	)

	// Récupérer les infos utilisateur pour la réponse
	user, err := h.authService.GetUserFromToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}

	// Respond with user info (sans le token dans le JSON)
	c.JSON(http.StatusOK, gin.H{
		"message": "Connexion réussie",
		"user": gin.H{
			"id":        user.Id,
			"firstname": user.Firstname,
			"lastname":  user.Lastname,
			"email":     user.Email,
			"role_id":   user.RoleID,
		},
	})
}

// Logout endpoint: POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Récupérer le token depuis le cookie
	token, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no token provided"})
		return
	}

	// Appeler le service de déconnexion
	if err := h.authService.Logout("Bearer " + token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Supprimer le cookie en définissant une expiration passée
	c.SetCookie(
		"auth_token",
		"",
		-1, // maxAge négatif pour supprimer le cookie
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Déconnexion réussie"})
}

// GetCurrentUser endpoint: GET /auth/me
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Récupérer le token depuis le cookie
	token, err := c.Cookie("auth_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	// Récupérer les infos utilisateur
	user, err := h.authService.GetUserFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":        user.Id,
			"firstname": user.Firstname,
			"lastname":  user.Lastname,
			"email":     user.Email,
			"role_id":   user.RoleID,
		},
	})
}

// GetSessionsStats endpoint: GET /auth/sessions/stats (pour monitoring)
func (h *AuthHandler) GetSessionsStats(c *gin.Context) {
	count, err := h.authService.GetActiveSessionsCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sessions count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"active_sessions": count,
		"timestamp":       time.Now().Unix(),
	})
}
