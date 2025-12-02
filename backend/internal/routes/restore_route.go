package routes

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/handlers"
	"github.com/RyanLadmia/plateforme-safebase/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// SetupRestoreRoutes configures all restore-related routes
func SetupRestoreRoutes(router *gin.Engine, restoreHandler *handlers.RestoreHandler, authMiddleware *middlewares.AuthMiddleware) {
	// Restore routes group
	restoreRoutes := router.Group("/api/restores")
	restoreRoutes.Use(authMiddleware.RequireAuth())
	{
		// General restore routes
		restoreRoutes.GET("", restoreHandler.GetRestores)
		restoreRoutes.GET("/:id", restoreHandler.GetRestore)

		// Database-specific restore routes
		restoreRoutes.GET("/database/:database_id", restoreHandler.GetRestoresByDatabase)

		// Backup-specific restore routes
		restoreRoutes.GET("/backup/:backup_id", restoreHandler.GetRestoresByBackup)
		restoreRoutes.POST("/backup/:backup_id/database/:database_id", restoreHandler.CreateRestore)
	}
}
