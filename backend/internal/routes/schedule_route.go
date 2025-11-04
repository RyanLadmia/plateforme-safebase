package routes

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/handlers"
	"github.com/RyanLadmia/plateforme-safebase/internal/middlewares"
	"github.com/gin-gonic/gin"
)

// SetupScheduleRoutes configures all schedule-related routes
func SetupScheduleRoutes(router *gin.Engine, scheduleHandler *handlers.ScheduleHandler, authMiddleware *middlewares.AuthMiddleware) {
	// Schedule routes group
	scheduleRoutes := router.Group("/api/schedules")
	scheduleRoutes.Use(authMiddleware.RequireAuth())
	{
		scheduleRoutes.POST("", scheduleHandler.CreateSchedule)
		scheduleRoutes.GET("", scheduleHandler.GetSchedules)
		scheduleRoutes.GET("/:id", scheduleHandler.GetSchedule)
		scheduleRoutes.PUT("/:id", scheduleHandler.UpdateSchedule)
		scheduleRoutes.DELETE("/:id", scheduleHandler.DeleteSchedule)
	}
}
