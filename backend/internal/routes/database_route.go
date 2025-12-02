package routes

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/handlers"
	"github.com/RyanLadmia/plateforme-safebase/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// SetupDatabaseRoutes configures all database-related routes
func SetupDatabaseRoutes(router *gin.Engine, databaseHandler *handlers.DatabaseHandler, authMiddleware *middlewares.AuthMiddleware) {
	// Database routes group
	databaseRoutes := router.Group("/api/databases")
	databaseRoutes.Use(authMiddleware.RequireAuth())
	{
		databaseRoutes.POST("", databaseHandler.CreateDatabase)
		databaseRoutes.GET("", databaseHandler.GetDatabases)
		databaseRoutes.GET("/:id", databaseHandler.GetDatabase)
		databaseRoutes.PUT("/:id", databaseHandler.UpdateDatabase)
		databaseRoutes.PUT("/:id/partial", databaseHandler.UpdateDatabasePartial) // Secure partial update
		// databaseRoutes.DELETE("/:id", databaseHandler.DeleteDatabase) // TODO: Implement DeleteDatabase method
	}
}
