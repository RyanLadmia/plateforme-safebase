package routes

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/handlers"
	"github.com/RyanLadmia/plateforme-safebase/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func ProfileRoutes(router *gin.Engine, profileHandler *handlers.ProfileHandler, authMiddleware *middlewares.AuthMiddleware) {
	// Profile routes - require authentication
	profile := router.Group("/api/profile")
	profile.Use(authMiddleware.RequireAuth())
	{
		profile.GET("", profileHandler.GetProfile)
		profile.PUT("", profileHandler.UpdateProfile)
		profile.PUT("/password", profileHandler.ChangePassword)
	}
}

