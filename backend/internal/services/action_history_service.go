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

// ActionHistoryResponse represents the response format for action history
type ActionHistoryResponse struct {
	Id           uint                   `json:"id"`
	UserId       uint                   `json:"user_id"`
	User         *models.User           `json:"user,omitempty"`
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resource_type"`
	ResourceId   uint                   `json:"resource_id"`
	Description  string                 `json:"description"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	IpAddress    string                 `json:"ip_address,omitempty"`
	UserAgent    string                 `json:"user_agent,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
}

// convertToResponse converts ActionHistory model to response format
func (s *ActionHistoryService) convertToResponse(history models.ActionHistory) ActionHistoryResponse {
	response := ActionHistoryResponse{
		Id:           history.Id,
		UserId:       history.UserId,
		User:         &history.User,
		Action:       history.Action,
		ResourceType: history.ResourceType,
		ResourceId:   history.ResourceId,
		Description:  history.Description,
		IpAddress:    history.IpAddress,
		UserAgent:    history.UserAgent,
		CreatedAt:    history.CreatedAt,
	}

	// Parse metadata JSON string to map
	if history.Metadata != "" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(history.Metadata), &metadata); err != nil {
			fmt.Printf("Error parsing metadata JSON for history ID %d: %v\n", history.Id, err)
		} else {
			response.Metadata = metadata
		}
	}

	return response
}

// GetUserActionHistory gets action history for a user with pagination
func (s *ActionHistoryService) GetUserActionHistory(userID uint, page, limit int) ([]ActionHistoryResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
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

	// Convert to response format
	responses := make([]ActionHistoryResponse, len(histories))
	for i, history := range histories {
		responses[i] = s.convertToResponse(history)
	}

	return responses, total, nil
}

// GetResourceActionHistory gets action history for a specific resource
func (s *ActionHistoryService) GetResourceActionHistory(resourceType string, resourceID uint) ([]ActionHistoryResponse, error) {
	histories, err := s.actionHistoryRepo.GetByResource(resourceType, resourceID)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	responses := make([]ActionHistoryResponse, len(histories))
	for i, history := range histories {
		responses[i] = s.convertToResponse(history)
	}

	return responses, nil
}

// GetActionHistoryByType gets action history for a resource type with pagination
func (s *ActionHistoryService) GetActionHistoryByType(resourceType string, page, limit int) ([]ActionHistoryResponse, int64, error) {
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

	// Convert to response format
	responses := make([]ActionHistoryResponse, len(histories))
	for i, history := range histories {
		responses[i] = s.convertToResponse(history)
	}

	return responses, total, nil
}

// GetRecentActionHistory gets recent action history with pagination
func (s *ActionHistoryService) GetRecentActionHistory(page, limit int) ([]ActionHistoryResponse, int64, error) {
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

	// Convert to response format
	responses := make([]ActionHistoryResponse, len(histories))
	for i, history := range histories {
		responses[i] = s.convertToResponse(history)
	}

	return responses, total, nil
}

// Helper methods for common actions

// LogDatabaseAction logs database-related actions
func (s *ActionHistoryService) LogDatabaseAction(userID uint, action string, databaseID uint, databaseName string, databaseType string, ipAddress, userAgent string) error {
	description := fmt.Sprintf("Base de données '%s' %s", databaseName, s.getActionDescription(action))
	metadata := map[string]interface{}{
		"database_name": databaseName,
		"database_type": databaseType,
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
		"created":   "ajoutée",
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
