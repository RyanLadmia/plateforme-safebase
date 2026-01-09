package handlers

import (
	"net/http"
	"os"
	"strings"

	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/gin-gonic/gin"
)

// TestHandler handles test-related operations (cleanup, etc.)
type TestHandler struct {
	userRepo *repositories.UserRepository
}

// NewTestHandler creates a new test handler
func NewTestHandler(userRepo *repositories.UserRepository) *TestHandler {
	return &TestHandler{
		userRepo: userRepo,
	}
}

// CleanupTestUsers DELETE /api/test/cleanup-users
// Deletes all users with email ending in @e2e.com
// This endpoint should only be available in non-production environments
func (h *TestHandler) CleanupTestUsers(c *gin.Context) {
	// Security check: only allow in non-production environments
	env := os.Getenv("GO_ENV")
	if env == "production" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Test cleanup not allowed in production"})
		return
	}

	// Get all users
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	// Filter and delete test users
	deletedCount := 0
	for _, user := range users {
		if strings.HasSuffix(user.Email, "@e2e.com") {
			// Delete user by ID
			err := h.userRepo.DeleteUser(user.Id)
			if err == nil {
				deletedCount++
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Test users cleaned up successfully",
		"deleted_count": deletedCount,
	})
}

