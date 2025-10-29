package routes

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/handlers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, jwtSecret string) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/me", authHandler.GetCurrentUser)
		auth.GET("/sessions/stats", authHandler.GetSessionsStats) // Monitoring
	}
}
