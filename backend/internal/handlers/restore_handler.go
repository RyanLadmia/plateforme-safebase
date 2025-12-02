package handlers

import (
	"net/http"
	"strconv"

	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/gin-gonic/gin"
)

type RestoreHandler struct {
	restoreService *services.RestoreService
}

// Constructor for RestoreHandler
func NewRestoreHandler(restoreService *services.RestoreService) *RestoreHandler {
	return &RestoreHandler{
		restoreService: restoreService,
	}
}

// CreateRestore creates a new restore operation
func (h *RestoreHandler) CreateRestore(c *gin.Context) {
	backupIDParam := c.Param("backup_id")
	backupID, err := strconv.ParseUint(backupIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de sauvegarde invalide"})
		return
	}

	databaseIDParam := c.Param("database_id")
	databaseID, err := strconv.ParseUint(databaseIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de base de données invalide"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	restore, err := h.restoreService.CreateRestore(uint(backupID), uint(databaseID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la restauration: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Restauration créée avec succès. Le processus de restauration a commencé.",
		"restore": restore,
	})
}

// GetRestores returns all restores for the authenticated user
func (h *RestoreHandler) GetRestores(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	restores, err := h.restoreService.GetRestoresByUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des restaurations: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"restores": restores,
	})
}

// GetRestoresByDatabase returns all restores for a specific database
func (h *RestoreHandler) GetRestoresByDatabase(c *gin.Context) {
	databaseIDParam := c.Param("database_id")
	databaseID, err := strconv.ParseUint(databaseIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de base de données invalide"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	restores, err := h.restoreService.GetRestoresByDatabase(uint(databaseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des restaurations: " + err.Error()})
		return
	}

	// Filter restores to only show those belonging to the user
	var userRestores []interface{}
	for _, restore := range restores {
		if restore.UserId == userID.(uint) {
			userRestores = append(userRestores, restore)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"restores": userRestores,
	})
}

// GetRestoresByBackup returns all restores for a specific backup
func (h *RestoreHandler) GetRestoresByBackup(c *gin.Context) {
	backupIDParam := c.Param("backup_id")
	backupID, err := strconv.ParseUint(backupIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de sauvegarde invalide"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	restores, err := h.restoreService.GetRestoresByBackup(uint(backupID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des restaurations: " + err.Error()})
		return
	}

	// Filter restores to only show those belonging to the user
	var userRestores []interface{}
	for _, restore := range restores {
		if restore.UserId == userID.(uint) {
			userRestores = append(userRestores, restore)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"restores": userRestores,
	})
}

// GetRestore returns a specific restore by ID
func (h *RestoreHandler) GetRestore(c *gin.Context) {
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

	restore, err := h.restoreService.GetRestoreByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restauration introuvable"})
		return
	}

	// Verify user ownership
	if restore.UserId != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"restore": restore,
	})
}
