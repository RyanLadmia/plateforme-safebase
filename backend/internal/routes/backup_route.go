package routes

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/handlers"
	"github.com/RyanLadmia/plateforme-safebase/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// SetupBackupRoutes configures all backup-related routes
func SetupBackupRoutes(router *gin.Engine, backupHandler *handlers.BackupHandler, authMiddleware *middlewares.AuthMiddleware) {
	// Backup routes group
	backupRoutes := router.Group("/api/backups")
	backupRoutes.Use(authMiddleware.RequireAuth())
	{
		// General backup routes
		backupRoutes.GET("", backupHandler.GetBackups)
		backupRoutes.GET("/:id", backupHandler.GetBackup)
		backupRoutes.DELETE("/:id", backupHandler.DeleteBackup)
		backupRoutes.GET("/:id/download", backupHandler.DownloadBackup)

		// Database-specific backup routes
		backupRoutes.POST("/database/:database_id", backupHandler.CreateBackup)
		backupRoutes.GET("/database/:database_id", backupHandler.GetBackupsByDatabase)
	}
}
