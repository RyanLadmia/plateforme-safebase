package units

import (
	"encoding/json"
	"testing"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// setupHistoryTestDB initializes an in-memory SQLite database for testing
func setupHistoryTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.ActionHistory{}, &models.User{}, &models.Role{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Create test users
	roleID := uint(2)
	users := []models.User{
		{
			Id:        1,
			Firstname: "Test",
			Lastname:  "User",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			Active:    true,
			RoleID:    &roleID,
		},
		{
			Id:        2,
			Firstname: "Another",
			Lastname:  "User",
			Email:     "another@example.com",
			Password:  "hashedpassword",
			Active:    true,
			RoleID:    &roleID,
		},
	}
	for _, user := range users {
		db.Create(&user)
	}

	return db
}

// createTestActionHistory creates a test action history entry
func createTestActionHistory(db *gorm.DB, userID uint, action, resourceType string, resourceID uint, description string) *models.ActionHistory {
	metadata := map[string]interface{}{
		"test_key": "test_value",
	}
	metadataJSON, _ := json.Marshal(metadata)

	history := &models.ActionHistory{
		UserId:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceId:   resourceID,
		Description:  description,
		Metadata:     string(metadataJSON),
		IpAddress:    "127.0.0.1",
		UserAgent:    "test-agent",
	}
	db.Create(history)
	return history
}

// ============================================================================
// TESTS - Action History Service
// ============================================================================

// Test 1: Log Action - Should create action history entry
func TestActionHistoryService_LogAction(t *testing.T) {
	db := setupHistoryTestDB(t)

	// Create repository and service
	historyRepo := repositories.NewActionHistoryRepository(db)
	historyService := services.NewActionHistoryService(historyRepo)

	metadata := map[string]interface{}{
		"database_name": "Test Database",
		"database_type": "mysql",
	}

	// Test logging an action
	err := historyService.LogAction(1, "created", "database", 1, "Database 'Test Database' created", metadata, "127.0.0.1", "test-agent")
	assert.NoError(t, err, "Action logging should succeed")

	// Verify the action was logged
	var history models.ActionHistory
	db.First(&history)
	assert.Equal(t, uint(1), history.UserId, "User ID should match")
	assert.Equal(t, "created", history.Action, "Action should match")
	assert.Equal(t, "database", history.ResourceType, "Resource type should match")
	assert.Equal(t, uint(1), history.ResourceId, "Resource ID should match")
	assert.Contains(t, history.Description, "Test Database", "Description should contain database name")

	// Verify metadata was stored as JSON
	var storedMetadata map[string]interface{}
	err = json.Unmarshal([]byte(history.Metadata), &storedMetadata)
	assert.NoError(t, err, "Should unmarshal metadata successfully")
	assert.Equal(t, "Test Database", storedMetadata["database_name"], "Metadata should be stored correctly")
}

// Test 2: Get User Action History - Should return paginated history for a user
func TestActionHistoryService_GetUserActionHistory(t *testing.T) {
	db := setupHistoryTestDB(t)

	// Create repository and service
	historyRepo := repositories.NewActionHistoryRepository(db)
	historyService := services.NewActionHistoryService(historyRepo)

	// Create multiple history entries for user 1
	createTestActionHistory(db, 1, "created", "database", 1, "Database created")
	createTestActionHistory(db, 1, "updated", "database", 1, "Database updated")
	createTestActionHistory(db, 1, "deleted", "backup", 2, "Backup deleted")
	createTestActionHistory(db, 2, "created", "schedule", 3, "Schedule created") // Different user

	// Get user 1's history
	histories, total, err := historyService.GetUserActionHistory(1, 1, 10)
	assert.NoError(t, err, "Should retrieve history successfully")
	assert.Equal(t, int64(3), total, "Total should be 3 for user 1")
	assert.Equal(t, 3, len(histories), "Should return 3 history entries")

	// Verify all entries belong to user 1
	for _, history := range histories {
		assert.Equal(t, uint(1), history.UserId, "All entries should belong to user 1")
	}

	// Test pagination
	histories, total, err = historyService.GetUserActionHistory(1, 1, 2)
	assert.NoError(t, err, "Should retrieve paginated history successfully")
	assert.Equal(t, int64(3), total, "Total should still be 3")
	assert.Equal(t, 2, len(histories), "Should return only 2 entries (page size)")
}

// Test 3: Get Action History by Resource Type - Should filter by resource type
func TestActionHistoryService_GetActionHistoryByType(t *testing.T) {
	db := setupHistoryTestDB(t)

	// Create repository and service
	historyRepo := repositories.NewActionHistoryRepository(db)
	historyService := services.NewActionHistoryService(historyRepo)

	// Create history entries for different resource types
	createTestActionHistory(db, 1, "created", "database", 1, "Database created")
	createTestActionHistory(db, 1, "updated", "database", 2, "Database updated")
	createTestActionHistory(db, 1, "created", "backup", 1, "Backup created")
	createTestActionHistory(db, 1, "created", "schedule", 1, "Schedule created")

	// Get history for database resource type
	histories, total, err := historyService.GetActionHistoryByType("database", 1, 10)
	assert.NoError(t, err, "Should retrieve history by type successfully")
	assert.Equal(t, int64(2), total, "Total should be 2 for database type")
	assert.Equal(t, 2, len(histories), "Should return 2 database history entries")

	// Verify all entries are for database resource type
	for _, history := range histories {
		assert.Equal(t, "database", history.ResourceType, "All entries should be database type")
	}

	// Get history for backup resource type
	_, total, err = historyService.GetActionHistoryByType("backup", 1, 10)
	assert.NoError(t, err, "Should retrieve backup history successfully")
	assert.Equal(t, int64(1), total, "Total should be 1 for backup type")
}

// Test 4: Get Action History by Resource - Should return history for specific resource
func TestActionHistoryService_GetResourceActionHistory(t *testing.T) {
	db := setupHistoryTestDB(t)

	// Create repository and service
	historyRepo := repositories.NewActionHistoryRepository(db)
	historyService := services.NewActionHistoryService(historyRepo)

	// Create history entries for the same resource
	createTestActionHistory(db, 1, "created", "database", 5, "Database created")
	createTestActionHistory(db, 1, "updated", "database", 5, "Database name changed")
	createTestActionHistory(db, 1, "updated", "database", 5, "Database credentials updated")
	createTestActionHistory(db, 1, "deleted", "database", 10, "Different database deleted")

	// Get history for database ID 5
	histories, err := historyService.GetResourceActionHistory("database", 5)
	assert.NoError(t, err, "Should retrieve resource history successfully")
	assert.Equal(t, 3, len(histories), "Should return 3 history entries for database 5")

	// Verify all entries are for resource ID 5
	for _, history := range histories {
		assert.Equal(t, uint(5), history.ResourceId, "All entries should be for resource ID 5")
		assert.Equal(t, "database", history.ResourceType, "All entries should be database type")
	}
}

// Test 5: Filter Action History by Action Type - Should return specific actions
func TestActionHistoryService_GetUserActionHistoryByType(t *testing.T) {
	db := setupHistoryTestDB(t)

	// Create repository and service
	historyRepo := repositories.NewActionHistoryRepository(db)
	historyService := services.NewActionHistoryService(historyRepo)

	// Create history entries with different actions for user 1
	createTestActionHistory(db, 1, "created", "database", 1, "Database created")
	createTestActionHistory(db, 1, "created", "database", 2, "Database created")
	createTestActionHistory(db, 1, "updated", "database", 1, "Database updated")
	createTestActionHistory(db, 1, "created", "backup", 1, "Backup created")
	createTestActionHistory(db, 1, "created", "backup", 2, "Backup created")

	// Get database history for user 1
	histories, total, err := historyService.GetUserActionHistoryByType(1, "database", 1, 10)
	assert.NoError(t, err, "Should retrieve filtered history successfully")
	assert.Equal(t, int64(3), total, "Total should be 3 database entries")
	assert.Equal(t, 3, len(histories), "Should return 3 database history entries")

	// Verify filtering
	for _, history := range histories {
		assert.Equal(t, uint(1), history.UserId, "All entries should belong to user 1")
		assert.Equal(t, "database", history.ResourceType, "All entries should be database type")
	}

	// Get backup history for user 1
	_, total, err = historyService.GetUserActionHistoryByType(1, "backup", 1, 10)
	assert.NoError(t, err, "Should retrieve backup history successfully")
	assert.Equal(t, int64(2), total, "Total should be 2 backup entries")
}

// Test 6: Helper Methods - Should log specific resource actions correctly
func TestActionHistoryService_HelperMethods(t *testing.T) {
	db := setupHistoryTestDB(t)

	// Create repository and service
	historyRepo := repositories.NewActionHistoryRepository(db)
	historyService := services.NewActionHistoryService(historyRepo)

	// Test LogDatabaseAction
	err := historyService.LogDatabaseAction(1, "created", 1, "Test Database", "mysql", "127.0.0.1", "test-agent")
	assert.NoError(t, err, "LogDatabaseAction should succeed")

	// Test LogBackupAction
	err = historyService.LogBackupAction(1, "created", 1, "Test Database", 1024, "127.0.0.1", "test-agent")
	assert.NoError(t, err, "LogBackupAction should succeed")

	// Test LogScheduleAction
	err = historyService.LogScheduleAction(1, "created", 1, "Test Database", "0 0 * * *", "127.0.0.1", "test-agent")
	assert.NoError(t, err, "LogScheduleAction should succeed")

	// Verify all actions were logged
	var count int64
	db.Model(&models.ActionHistory{}).Count(&count)
	assert.Equal(t, int64(3), count, "Should have 3 action history entries")
}

// Test 7: Metadata Parsing - Should correctly parse JSON metadata
func TestActionHistoryService_MetadataParsing(t *testing.T) {
	db := setupHistoryTestDB(t)

	// Create repository and service
	historyRepo := repositories.NewActionHistoryRepository(db)
	historyService := services.NewActionHistoryService(historyRepo)

	// Log action with complex metadata
	metadata := map[string]interface{}{
		"database_name": "Production Database",
		"database_type": "postgresql",
		"changes": map[string]interface{}{
			"name": map[string]string{
				"from": "Old Name",
				"to":   "New Name",
			},
		},
		"size": 2048,
	}

	err := historyService.LogAction(1, "updated", "database", 1, "Database updated", metadata, "127.0.0.1", "test-agent")
	assert.NoError(t, err, "Action logging should succeed")

	// Retrieve and verify metadata parsing
	histories, _, err := historyService.GetUserActionHistory(1, 1, 10)
	assert.NoError(t, err, "Should retrieve history successfully")
	assert.Equal(t, 1, len(histories), "Should return 1 history entry")

	// Verify metadata was parsed correctly
	retrievedMetadata := histories[0].Metadata
	assert.NotNil(t, retrievedMetadata, "Metadata should not be nil")
	assert.Equal(t, "Production Database", retrievedMetadata["database_name"], "Database name should match")
	assert.Equal(t, "postgresql", retrievedMetadata["database_type"], "Database type should match")
	assert.Equal(t, float64(2048), retrievedMetadata["size"], "Size should match")
}

