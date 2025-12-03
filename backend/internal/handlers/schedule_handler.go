package handlers

import (
	"net/http"
	"strconv"

	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	scheduleService *services.ScheduleService
}

// Constructor for ScheduleHandler
func NewScheduleHandler(scheduleService *services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
	}
}

// CreateSchedule creates a new schedule
func (h *ScheduleHandler) CreateSchedule(c *gin.Context) {
	var request struct {
		DatabaseID     uint   `json:"database_id" binding:"required"`
		CronExpression string `json:"cron_expression" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Extract IP address and User-Agent for logging
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	schedule, err := h.scheduleService.CreateSchedule(request.DatabaseID, userID.(uint), request.CronExpression, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du schedule: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Schedule créé avec succès",
		"schedule": schedule,
	})
}

// GetSchedules returns all schedules for the authenticated user
func (h *ScheduleHandler) GetSchedules(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	schedules, err := h.scheduleService.GetSchedules(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des schedules: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"schedules": schedules,
	})
}

// GetSchedule returns a specific schedule by ID
func (h *ScheduleHandler) GetSchedule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	schedule, err := h.scheduleService.GetSchedule(uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du schedule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"schedule": schedule,
	})
}

// UpdateSchedule updates a schedule
func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var request struct {
		CronExpression *string `json:"cron_expression"`
		Active         *bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides: " + err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Extract IP address and User-Agent for logging
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	cronExpr := ""
	if request.CronExpression != nil {
		cronExpr = *request.CronExpression
	}
	active := true
	if request.Active != nil {
		active = *request.Active
	}

	schedule, err := h.scheduleService.UpdateSchedule(uint(id), userID.(uint), cronExpr, active, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du schedule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Schedule mis à jour avec succès",
		"schedule": schedule,
	})
}

// DeleteSchedule deletes a schedule
func (h *ScheduleHandler) DeleteSchedule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	// Extract IP address and User-Agent for logging
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	if err := h.scheduleService.DeleteSchedule(uint(id), userID.(uint), ipAddress, userAgent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression du schedule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Schedule supprimé avec succès",
	})
}
