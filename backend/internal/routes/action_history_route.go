package routes

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/handlers"
	"github.com/RyanLadmia/plateforme-safebase/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// SetupActionHistoryRoutes configures all action history-related routes
func SetupActionHistoryRoutes(router *gin.Engine, actionHistoryHandler *handlers.ActionHistoryHandler, authMiddleware *middlewares.AuthMiddleware) {
	// Action history routes group
	historyRoutes := router.Group("/api/history")
	historyRoutes.Use(authMiddleware.RequireAuth())
	{
		// User action history
		historyRoutes.GET("", actionHistoryHandler.GetUserActionHistory)

		// Action history by resource type
		historyRoutes.GET("/type/:type", actionHistoryHandler.GetActionHistoryByType)

		// Action history for specific resource
		historyRoutes.GET("/resource/:type/:id", actionHistoryHandler.GetResourceActionHistory)

		// Recent action history (admin only)
		historyRoutes.GET("/recent", actionHistoryHandler.GetRecentActionHistory)
	}
}