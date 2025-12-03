package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
)

// ActionHistoryService manages action history operations
type ActionHistoryService struct {
	actionHistoryRepo *repositories.ActionHistoryRepository
}

// NewActionHistoryService creates a new ActionHistoryService
func NewActionHistoryService(actionHistoryRepo *repositories.ActionHistoryRepository) *ActionHistoryService {
	return &ActionHistoryService{
		actionHistoryRepo: actionHistoryRepo,
	}
}

// LogAction logs a user action with metadata
func (s *ActionHistoryService) LogAction(userID uint, action, resourceType string, resourceID uint, description string, metadata map[string]interface{}, ipAddress, userAgent string) error {
	// Convert metadata to JSON string
	var metadataJSON string
	if metadata != nil {
		jsonBytes, err := json.Marshal(metadata)
		if err != nil {
			// Log the error but continue without metadata
			fmt.Printf("Error marshaling metadata: %v\n", err)
		} else {
			metadataJSON = string(jsonBytes)
		}
	}

	actionHistory := &models.ActionHistory{
		UserId:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceId:   resourceID,
		Description:  description,
		Metadata:     metadataJSON,
		IpAddress:    ipAddress,
		UserAgent:    userAgent,
		CreatedAt:    time.Now(),
	}

	return s.actionHistoryRepo.Create(actionHistory)
}

// GetUserActionHistory gets action history for a user with pagination
func (s *ActionHistoryService) GetUserActionHistory(userID uint, page, limit int) ([]models.ActionHistory, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20 // Default limit
	}

	offset := (page - 1) * limit

	histories, err := s.actionHistoryRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.actionHistoryRepo.GetByUserIDCount(userID)
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// GetResourceActionHistory gets action history for a specific resource
func (s *ActionHistoryService) GetResourceActionHistory(resourceType string, resourceID uint) ([]models.ActionHistory, error) {
	return s.actionHistoryRepo.GetByResource(resourceType, resourceID)
}

// GetActionHistoryByType gets action history for a resource type with pagination
func (s *ActionHistoryService) GetActionHistoryByType(resourceType string, page, limit int) ([]models.ActionHistory, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	histories, err := s.actionHistoryRepo.GetByResourceType(resourceType, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.actionHistoryRepo.GetByResourceTypeCount(resourceType)
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// GetRecentActionHistory gets recent action history with pagination
func (s *ActionHistoryService) GetRecentActionHistory(page, limit int) ([]models.ActionHistory, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	histories, err := s.actionHistoryRepo.GetRecent(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.actionHistoryRepo.GetRecentCount()
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

// Helper methods for common actions

// LogDatabaseAction logs database-related actions
func (s *ActionHistoryService) LogDatabaseAction(userID uint, action string, databaseID uint, databaseName string, ipAddress, userAgent string) error {
	description := fmt.Sprintf("Base de données '%s' %s", databaseName, s.getActionDescription(action))
	metadata := map[string]interface{}{
		"database_name": databaseName,
	}

	return s.LogAction(userID, action, "database", databaseID, description, metadata, ipAddress, userAgent)
}

// LogBackupAction logs backup-related actions
func (s *ActionHistoryService) LogBackupAction(userID uint, action string, backupID uint, databaseName string, size int64, ipAddress, userAgent string) error {
	description := fmt.Sprintf("Sauvegarde de '%s' %s", databaseName, s.getActionDescription(action))
	metadata := map[string]interface{}{
		"database_name": databaseName,
		"size":          size,
	}

	return s.LogAction(userID, action, "backup", backupID, description, metadata, ipAddress, userAgent)
}

// LogScheduleAction logs schedule-related actions
func (s *ActionHistoryService) LogScheduleAction(userID uint, action string, scheduleID uint, databaseName string, frequency string, ipAddress, userAgent string) error {
	description := fmt.Sprintf("Planification pour '%s' %s", databaseName, s.getActionDescription(action))
	metadata := map[string]interface{}{
		"database_name": databaseName,
		"frequency":     frequency,
	}

	return s.LogAction(userID, action, "schedule", scheduleID, description, metadata, ipAddress, userAgent)
}

// LogRestoreAction logs restore-related actions
func (s *ActionHistoryService) LogRestoreAction(userID uint, action string, restoreID uint, databaseName string, backupFile string, ipAddress, userAgent string) error {
	description := fmt.Sprintf("Restauration de '%s' %s", databaseName, s.getActionDescription(action))
	metadata := map[string]interface{}{
		"database_name": databaseName,
		"backup_file":   backupFile,
	}

	return s.LogAction(userID, action, "restore", restoreID, description, metadata, ipAddress, userAgent)
}

// getActionDescription returns a human-readable description for an action
func (s *ActionHistoryService) getActionDescription(action string) string {
	descriptions := map[string]string{
		"created":   "créée",
		"updated":   "modifiée",
		"deleted":   "supprimée",
		"completed": "terminée",
		"failed":    "échouée",
		"executed":  "exécutée",
	}

	if desc, exists := descriptions[action]; exists {
		return desc
	}
	return action
}
