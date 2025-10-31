package handlers

import (
	"net/http"
	"strconv"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/gin-gonic/gin"
)

type DatabaseHandler struct {
	databaseService *services.DatabaseService
}

// Constructor for DatabaseHandler
func NewDatabaseHandler(databaseService *services.DatabaseService) *DatabaseHandler {
	return &DatabaseHandler{
		databaseService: databaseService,
	}
}

// CreateDatabase creates a new database configuration
func (h *DatabaseHandler) CreateDatabase(c *gin.Context) {
	var request models.DatabaseCreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Convert request to Database model
	database := &models.Database{
		Name:     request.Name,
		Type:     request.Type,
		Host:     request.Host,
		Port:     request.Port,
		Username: request.Username,
		Password: request.Password,
		DbName:   request.DbName,
		UserId:   userID.(uint),
	}

	if err := h.databaseService.CreateDatabase(database); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la base de données: " + err.Error()})
		return
	}

	// Return database without password for security
	responseDb := *database
	responseDb.Password = "" // Don't expose password in response

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Base de données créée avec succès",
		"database": responseDb,
	})
}

// GetDatabases returns all databases for the authenticated user
func (h *DatabaseHandler) GetDatabases(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	databases, err := h.databaseService.GetDatabasesByUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des bases de données: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"databases": databases,
	})
}

// GetDatabase returns a specific database by ID
func (h *DatabaseHandler) GetDatabase(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	database, err := h.databaseService.GetDatabaseByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Base de données introuvable"})
		return
	}

	// Verify user ownership
	if database.UserId != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}

	// Return database without password for security
	responseDb := *database
	responseDb.Password = "" // Don't expose password in response

	c.JSON(http.StatusOK, gin.H{
		"database": responseDb,
	})
}

// UpdateDatabase updates a database configuration
func (h *DatabaseHandler) UpdateDatabase(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Get existing database
	existingDatabase, err := h.databaseService.GetDatabaseByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Base de données introuvable"})
		return
	}

	// Verify user ownership
	if existingDatabase.UserId != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}

	var request models.DatabaseUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	// Update fields
	existingDatabase.Name = request.Name
	existingDatabase.Type = request.Type
	existingDatabase.Host = request.Host
	existingDatabase.Port = request.Port
	existingDatabase.Username = request.Username
	existingDatabase.DbName = request.DbName
	if request.Password != "" {
		existingDatabase.Password = request.Password
	}

	if err := h.databaseService.UpdateDatabase(existingDatabase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour: " + err.Error()})
		return
	}

	// Return database without password for security
	responseDb := *existingDatabase
	responseDb.Password = "" // Don't expose password in response

	c.JSON(http.StatusOK, gin.H{
		"message":  "Base de données mise à jour avec succès",
		"database": responseDb,
	})
}

// DeleteDatabase deletes a database configuration
func (h *DatabaseHandler) DeleteDatabase(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	if err := h.databaseService.DeleteDatabase(uint(id), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Base de données supprimée avec succès",
	})
}
