package handlers

import (
	"net/http"
	"strconv"

	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/gin-gonic/gin"
)

type BackupHandler struct {
	backupService *services.BackupService
}

// Constructor for BackupHandler
func NewBackupHandler(backupService *services.BackupService) *BackupHandler {
	return &BackupHandler{
		backupService: backupService,
	}
}

// CreateBackup creates a new backup for a database
func (h *BackupHandler) CreateBackup(c *gin.Context) {
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

	// Extract IP address and User-Agent for logging
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	backup, err := h.backupService.CreateBackupWithLogging(uint(databaseID), userID.(uint), ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la sauvegarde: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Sauvegarde créée avec succès. Le processus de sauvegarde a commencé.",
		"backup":  backup,
	})
}

// GetBackups returns all backups for the authenticated user
func (h *BackupHandler) GetBackups(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	backups, err := h.backupService.GetBackupsByUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des sauvegardes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"backups": backups,
	})
}

// GetBackupsByDatabase returns all backups for a specific database
func (h *BackupHandler) GetBackupsByDatabase(c *gin.Context) {
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

	backups, err := h.backupService.GetBackupsByDatabase(uint(databaseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des sauvegardes: " + err.Error()})
		return
	}

	// Filter backups to only show those belonging to the user
	var userBackups []interface{}
	for _, backup := range backups {
		if backup.UserId == userID.(uint) {
			userBackups = append(userBackups, backup)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"backups": userBackups,
	})
}

// GetBackup returns a specific backup by ID
func (h *BackupHandler) GetBackup(c *gin.Context) {
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

	backup, err := h.backupService.GetBackupByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sauvegarde introuvable"})
		return
	}

	// Verify user ownership
	if backup.UserId != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"backup": backup,
	})
}

// DeleteBackup deletes a backup
func (h *BackupHandler) DeleteBackup(c *gin.Context) {
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

	// Extract IP address and User-Agent for logging
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	if err := h.backupService.DeleteBackupWithLogging(uint(id), userID.(uint), ipAddress, userAgent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sauvegarde supprimée avec succès",
	})
}

// DownloadBackup allows downloading a backup file
func (h *BackupHandler) DownloadBackup(c *gin.Context) {
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

	backup, err := h.backupService.GetBackupByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sauvegarde introuvable"})
		return
	}

	// Verify user ownership
	if backup.UserId != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}

	// Check if backup is completed
	if backup.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La sauvegarde n'est pas encore terminée"})
		return
	}

	// Extract IP address and User-Agent for logging
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Download the file data from the service (MinIO or local)
	fileData, err := h.backupService.DownloadBackupWithLogging(uint(id), userID.(uint), ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du téléchargement: " + err.Error()})
		return
	}

	// Serve the file data directly
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+backup.Filename)
	c.Header("Content-Type", "application/zip")
	c.Data(http.StatusOK, "application/zip", fileData)
}
