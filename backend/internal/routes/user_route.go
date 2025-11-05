package routes

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/handlers"
	"github.com/RyanLadmia/plateforme-safebase/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userHandler *handlers.UserHandler, authMiddleware *middlewares.AuthMiddleware) {
	// Admin routes - require admin role
	admin := router.Group("/api/admin/users")
	admin.Use(authMiddleware.RequireAuth())
	admin.Use(authMiddleware.RequireRole("admin"))
	{
		admin.GET("", userHandler.GetAllUsers)
		admin.GET("/active", userHandler.GetAllActiveUsers)
		admin.GET("/:id", userHandler.GetUser)
		admin.PUT("/:id", userHandler.UpdateUser)
		admin.PUT("/:id/role", userHandler.ChangeUserRole)
		admin.PUT("/:id/deactivate", userHandler.DeactivateUser)
		admin.PUT("/:id/activate", userHandler.ActivateUser)
	}
}
