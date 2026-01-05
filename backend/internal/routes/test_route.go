package routes

import (
	"github.com/RyanLadmia/plateforme-safebase/internal/handlers"
	"github.com/gin-gonic/gin"
)

// TestRoutes sets up the routes for test utilities (cleanup, etc.)
// These routes should only be available in non-production environments
func TestRoutes(router *gin.Engine, testHandler *handlers.TestHandler) {
	test := router.Group("/api/test")
	{
		// POST instead of DELETE for easier Cypress usage
		test.POST("/cleanup-users", testHandler.CleanupTestUsers)
	}
}

