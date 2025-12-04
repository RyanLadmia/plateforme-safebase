package handlers

import (
	"net/http"
	"strconv"

	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/gin-gonic/gin"
)

type ActionHistoryHandler struct {
	actionHistoryService *services.ActionHistoryService
}

// NewActionHistoryHandler creates a new ActionHistoryHandler
func NewActionHistoryHandler(actionHistoryService *services.ActionHistoryService) *ActionHistoryHandler {
	return &ActionHistoryHandler{
		actionHistoryService: actionHistoryService,
	}
}

// GetUserActionHistory returns action history for the authenticated user
func (h *ActionHistoryHandler) GetUserActionHistory(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Parse pagination parameters
	page := 1
	limit := 20

	if pageParam := c.Query("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get action history
	histories, total, err := h.actionHistoryService.GetUserActionHistory(userID.(uint), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de l'historique: " + err.Error()})
		return
	}

	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, gin.H{
		"history":     histories,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}

// GetActionHistoryByType returns action history filtered by resource type
func (h *ActionHistoryHandler) GetActionHistoryByType(c *gin.Context) {
	resourceType := c.Param("type")

	// Validate resource type
	validTypes := map[string]bool{
		"database": true,
		"backup":   true,
		"schedule": true,
		"restore":  true,
	}

	if !validTypes[resourceType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type de ressource invalide"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Parse pagination parameters
	page := 1
	limit := 20

	if pageParam := c.Query("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get action history by type for the authenticated user
	histories, total, err := h.actionHistoryService.GetUserActionHistoryByType(userID.(uint), resourceType, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de l'historique: " + err.Error()})
		return
	}

	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, gin.H{
		"history":     histories,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
		"type":        resourceType,
	})
}

// GetResourceActionHistory returns action history for a specific resource
func (h *ActionHistoryHandler) GetResourceActionHistory(c *gin.Context) {
	resourceType := c.Param("type")
	resourceIDParam := c.Param("id")

	resourceID, err := strconv.ParseUint(resourceIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ressource invalide"})
		return
	}

	// Validate resource type
	validTypes := map[string]bool{
		"database": true,
		"backup":   true,
		"schedule": true,
		"restore":  true,
	}

	if !validTypes[resourceType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type de ressource invalide"})
		return
	}

	// Get action history for resource
	histories, err := h.actionHistoryService.GetResourceActionHistory(resourceType, uint(resourceID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de l'historique: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history":     histories,
		"type":        resourceType,
		"resource_id": resourceID,
	})
}

// GetRecentActionHistory returns recent action history (admin only)
func (h *ActionHistoryHandler) GetRecentActionHistory(c *gin.Context) {
	// Check if user is admin
	isAdmin, exists := c.Get("is_admin")
	if !exists || !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès réservé aux administrateurs"})
		return
	}

	// Parse pagination parameters
	page := 1
	limit := 20

	if pageParam := c.Query("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get recent action history
	histories, total, err := h.actionHistoryService.GetRecentActionHistory(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de l'historique récent: " + err.Error()})
		return
	}

	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, gin.H{
		"history":     histories,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	})
}
