package integrations

import (
	"encoding/json"
	"testing"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================================
// INTEGRATION TESTS - Action History Logging Across Services
// ============================================================================

// setupHistoryServices creates all services with action history integration
func setupHistoryServices(db *gorm.DB) (*services.ActionHistoryService, *services.DatabaseService, *services.ScheduleService) {
	historyRepo := repositories.NewActionHistoryRepository(db)
	databaseRepo := repositories.NewDatabaseRepository(db)
	backupRepo := repositories.NewBackupRepository(db)
	restoreRepo := repositories.NewRestoreRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	
	historyService := services.NewActionHistoryService(historyRepo)
	databaseService := services.NewDatabaseService(databaseRepo, backupRepo, restoreRepo, scheduleRepo, nil)
	scheduleService := services.NewScheduleService(scheduleRepo, databaseRepo, nil)
	
	// Connect action history to other services
	databaseService.SetActionHistoryService(historyService)
	scheduleService.SetActionHistoryService(historyService)
	
	return historyService, databaseService, scheduleService
}

// TestIntegration_DatabaseActionHistoryLogging tests action history for database operations
func TestIntegration_DatabaseActionHistoryLogging(t *testing.T) {
	db := setupIntegrationDB(t)
	historyService, databaseService, _ := setupHistoryServices(db)
	user := createTestUser(db)

	// Create database
	database := &models.Database{
		Name:     "History Test DB",
		Type:     "mysql",
		Host:     "localhost",
		Port:     "3306",
		Username: "testuser",
		Password: "password123",
		DbName:   "testdb",
		UserId:   user.Id,
	}

	t.Run("Create Database Logs Action", func(t *testing.T) {
		err := databaseService.CreateDatabase(database, user.Id, "192.168.1.100", "Mozilla/5.0")
		require.NoError(t, err, "Database creation should succeed")

		// Verify action was logged
		histories, total, err := historyService.GetUserActionHistory(user.Id, 1, 10)
		require.NoError(t, err, "Should retrieve action history")
		assert.Equal(t, int64(1), total, "Should have 1 action logged")
		assert.Equal(t, "created", histories[0].Action, "Action should be 'created'")
		assert.Equal(t, "database", histories[0].ResourceType, "Resource type should be 'database'")
		assert.Equal(t, database.Id, histories[0].ResourceId, "Resource ID should match database ID")
		assert.Equal(t, "192.168.1.100", histories[0].IpAddress, "IP address should be logged")
		assert.Equal(t, "Mozilla/5.0", histories[0].UserAgent, "User agent should be logged")
	})

	t.Run("Update Database Logs Action", func(t *testing.T) {
		database.Name = "Updated History Test DB"
		err := databaseService.UpdateDatabase(database, user.Id, "192.168.1.101", "Chrome/90.0")
		require.NoError(t, err, "Database update should succeed")

		// Verify update action was logged
		histories, total, err := historyService.GetUserActionHistory(user.Id, 1, 10)
		require.NoError(t, err, "Should retrieve action history")
		assert.Equal(t, int64(2), total, "Should have 2 actions logged")
		assert.Equal(t, "updated", histories[0].Action, "Latest action should be 'updated'")
	})

	t.Run("Delete Database Logs Action", func(t *testing.T) {
		err := databaseService.DeleteDatabase(database.Id, user.Id, "192.168.1.102", "Firefox/88.0")
		require.NoError(t, err, "Database deletion should succeed")

		// Verify delete action was logged
		histories, total, err := historyService.GetUserActionHistory(user.Id, 1, 10)
		require.NoError(t, err, "Should retrieve action history")
		assert.Equal(t, int64(3), total, "Should have 3 actions logged")
		assert.Equal(t, "deleted", histories[0].Action, "Latest action should be 'deleted'")
		assert.Equal(t, "192.168.1.102", histories[0].IpAddress, "IP address should be logged")
	})
}

// TestIntegration_ScheduleActionHistoryLogging tests action history for schedule operations
func TestIntegration_ScheduleActionHistoryLogging(t *testing.T) {
	db := setupIntegrationDB(t)
	historyService, databaseService, scheduleService := setupHistoryServices(db)
	scheduleService.StartScheduler()
	user := createTestUser(db)

	// Create database first
	database := &models.Database{
		Name:     "Schedule History DB",
		Type:     "postgresql",
		Host:     "localhost",
		Port:     "5432",
		Username: "testuser",
		Password: "password123",
		DbName:   "testdb",
		UserId:   user.Id,
	}
	err := databaseService.CreateDatabase(database, user.Id, "127.0.0.1", "test-agent")
	require.NoError(t, err, "Database creation should succeed")

	var scheduleID uint

	t.Run("Create Schedule Logs Action", func(t *testing.T) {
		schedule, err := scheduleService.CreateSchedule(
			database.Id,
			user.Id,
			"Test Schedule",
			"0 0 * * *",
			"10.0.0.1",
			"Safari/14.0",
		)
		require.NoError(t, err, "Schedule creation should succeed")
		scheduleID = schedule.Id

		// Verify action was logged
		histories, _, err := historyService.GetUserActionHistoryByType(user.Id, "schedule", 1, 10)
		require.NoError(t, err, "Should retrieve schedule action history")
		assert.NotEmpty(t, histories, "Should have schedule actions logged")
		
		var createAction *services.ActionHistoryResponse
		for _, h := range histories {
			if h.Action == "create" && h.ResourceId == schedule.Id {
				createAction = &h
				break
			}
		}
		require.NotNil(t, createAction, "Create action should be logged")
		assert.Equal(t, "10.0.0.1", createAction.IpAddress, "IP address should be logged")
	})

	t.Run("Update Schedule Logs Action With Changes", func(t *testing.T) {
		inactive := false
		_, err := scheduleService.UpdateSchedule(
			scheduleID,
			user.Id,
			"Updated Schedule",
			"0 2 * * *",
			&inactive,
			"10.0.0.2",
			"Edge/90.0",
		)
		require.NoError(t, err, "Schedule update should succeed")

		// Verify update action was logged with changes metadata
		histories, _, err := historyService.GetUserActionHistoryByType(user.Id, "schedule", 1, 10)
		require.NoError(t, err, "Should retrieve schedule action history")
		
		var updateAction *services.ActionHistoryResponse
		for _, h := range histories {
			if h.Action == "update" && h.ResourceId == scheduleID {
				updateAction = &h
				break
			}
		}
		require.NotNil(t, updateAction, "Update action should be logged")
		
		// Verify metadata contains changes
		if updateAction.Metadata != nil {
			changes, ok := updateAction.Metadata["changes"].(map[string]interface{})
			assert.True(t, ok, "Metadata should contain changes")
			if ok {
				assert.NotNil(t, changes, "Changes should be recorded")
			}
		}
	})

	t.Run("Delete Schedule Logs Action", func(t *testing.T) {
		err := scheduleService.DeleteSchedule(scheduleID, user.Id, "10.0.0.3", "Opera/75.0")
		require.NoError(t, err, "Schedule deletion should succeed")

		// Verify delete action was logged
		histories, _, err := historyService.GetUserActionHistoryByType(user.Id, "schedule", 1, 10)
		require.NoError(t, err, "Should retrieve schedule action history")
		
		var deleteAction *services.ActionHistoryResponse
		for _, h := range histories {
			if h.Action == "delete" && h.ResourceId == scheduleID {
				deleteAction = &h
				break
			}
		}
		require.NotNil(t, deleteAction, "Delete action should be logged")
	})
}

// TestIntegration_ActionHistoryPagination tests pagination across different resource types
func TestIntegration_ActionHistoryPagination(t *testing.T) {
	db := setupIntegrationDB(t)
	historyService, databaseService, scheduleService := setupHistoryServices(db)
	scheduleService.StartScheduler()
	user := createTestUser(db)

	// Create multiple databases to generate history
	for i := 1; i <= 15; i++ {
		database := &models.Database{
			Name:     "Test DB " + string(rune(i+'0')),
			Type:     "mysql",
			Host:     "localhost",
			Port:     "3306",
			Username: "user",
			Password: "pass",
			DbName:   "testdb",
			UserId:   user.Id,
		}
		err := databaseService.CreateDatabase(database, user.Id, "127.0.0.1", "test-agent")
		require.NoError(t, err, "Database creation should succeed")
	}

	t.Run("First Page", func(t *testing.T) {
		histories, total, err := historyService.GetUserActionHistory(user.Id, 1, 5)
		require.NoError(t, err, "Should retrieve first page")
		assert.Equal(t, 5, len(histories), "First page should have 5 entries")
		assert.Equal(t, int64(15), total, "Total should be 15")
	})

	t.Run("Second Page", func(t *testing.T) {
		histories, total, err := historyService.GetUserActionHistory(user.Id, 2, 5)
		require.NoError(t, err, "Should retrieve second page")
		assert.Equal(t, 5, len(histories), "Second page should have 5 entries")
		assert.Equal(t, int64(15), total, "Total should still be 15")
	})

	t.Run("Last Page", func(t *testing.T) {
		histories, total, err := historyService.GetUserActionHistory(user.Id, 3, 5)
		require.NoError(t, err, "Should retrieve last page")
		assert.Equal(t, 5, len(histories), "Last page should have 5 entries")
		assert.Equal(t, int64(15), total, "Total should still be 15")
	})

	t.Run("Beyond Last Page", func(t *testing.T) {
		histories, total, err := historyService.GetUserActionHistory(user.Id, 4, 5)
		require.NoError(t, err, "Should handle page beyond last")
		assert.Equal(t, 0, len(histories), "Should return empty array")
		assert.Equal(t, int64(15), total, "Total should still be 15")
	})
}

// TestIntegration_ActionHistoryMetadata tests metadata storage and retrieval
func TestIntegration_ActionHistoryMetadata(t *testing.T) {
	db := setupIntegrationDB(t)
	historyService, _, _ := setupHistoryServices(db)
	user := createTestUser(db)

	// Log action with complex metadata
	metadata := map[string]interface{}{
		"database_name": "Production DB",
		"database_type": "postgresql",
		"host":          "prod.example.com",
		"port":          5432,
		"changes": map[string]interface{}{
			"name": map[string]string{
				"from": "Old Name",
				"to":   "New Name",
			},
			"host": map[string]string{
				"from": "old.example.com",
				"to":   "new.example.com",
			},
		},
		"tags":   []string{"production", "critical", "backup"},
		"active": true,
	}

	err := historyService.LogAction(
		user.Id,
		"updated",
		"database",
		1,
		"Complex metadata test",
		metadata,
		"192.168.1.50",
		"TestAgent/1.0",
	)
	require.NoError(t, err, "Should log action with metadata")

	// Retrieve and verify metadata
	histories, _, err := historyService.GetUserActionHistory(user.Id, 1, 10)
	require.NoError(t, err, "Should retrieve action history")
	assert.NotEmpty(t, histories, "Should have actions logged")

	action := histories[0]
	assert.NotNil(t, action.Metadata, "Metadata should be present")
	
	// Verify metadata fields
	assert.Equal(t, "Production DB", action.Metadata["database_name"], "Database name should match")
	assert.Equal(t, "postgresql", action.Metadata["database_type"], "Database type should match")
	assert.Equal(t, float64(5432), action.Metadata["port"], "Port should match (JSON numbers are float64)")
	assert.Equal(t, true, action.Metadata["active"], "Active flag should match")
	
	// Verify nested changes
	changes, ok := action.Metadata["changes"].(map[string]interface{})
	assert.True(t, ok, "Changes should be a map")
	if ok {
		nameChange, ok := changes["name"].(map[string]interface{})
		assert.True(t, ok, "Name change should be a map")
		if ok {
			assert.Equal(t, "Old Name", nameChange["from"], "Old name should match")
			assert.Equal(t, "New Name", nameChange["to"], "New name should match")
		}
	}
	
	// Verify array
	tags, ok := action.Metadata["tags"].([]interface{})
	assert.True(t, ok, "Tags should be an array")
	if ok {
		assert.Equal(t, 3, len(tags), "Should have 3 tags")
	}
}

// TestIntegration_ActionHistoryFiltering tests filtering by resource type and action
func TestIntegration_ActionHistoryFiltering(t *testing.T) {
	db := setupIntegrationDB(t)
	historyService, databaseService, scheduleService := setupHistoryServices(db)
	scheduleService.StartScheduler()
	user := createTestUser(db)

	// Create databases
	for i := 1; i <= 3; i++ {
		database := &models.Database{
			Name:     "Test DB " + string(rune(i+'0')),
			Type:     "mysql",
			Host:     "localhost",
			Port:     "3306",
			Username: "user",
			Password: "pass",
			DbName:   "testdb",
			UserId:   user.Id,
		}
		err := databaseService.CreateDatabase(database, user.Id, "127.0.0.1", "test-agent")
		require.NoError(t, err, "Database creation should succeed")
		
		// Also create a schedule for each database
		_, err = scheduleService.CreateSchedule(
			database.Id,
			user.Id,
			"Schedule "+string(rune(i+'0')),
			"0 0 * * *",
			"127.0.0.1",
			"test-agent",
		)
		require.NoError(t, err, "Schedule creation should succeed")
	}

	t.Run("Filter By Database Type", func(t *testing.T) {
		histories, total, err := historyService.GetUserActionHistoryByType(user.Id, "database", 1, 10)
		require.NoError(t, err, "Should retrieve database actions")
		assert.Equal(t, int64(3), total, "Should have 3 database actions")
		
		// Verify all are database actions
		for _, h := range histories {
			assert.Equal(t, "database", h.ResourceType, "All should be database type")
		}
	})

	t.Run("Filter By Schedule Type", func(t *testing.T) {
		histories, total, err := historyService.GetUserActionHistoryByType(user.Id, "schedule", 1, 10)
		require.NoError(t, err, "Should retrieve schedule actions")
		assert.Equal(t, int64(3), total, "Should have 3 schedule actions")
		
		// Verify all are schedule actions
		for _, h := range histories {
			assert.Equal(t, "schedule", h.ResourceType, "All should be schedule type")
		}
	})

	t.Run("Get All Resource Types", func(t *testing.T) {
		_, total, err := historyService.GetUserActionHistory(user.Id, 1, 20)
		require.NoError(t, err, "Should retrieve all actions")
		assert.Equal(t, int64(6), total, "Should have 6 total actions (3 databases + 3 schedules)")
	})
}

// TestIntegration_ActionHistoryUserIsolation tests that users cannot see each other's history
func TestIntegration_ActionHistoryUserIsolation(t *testing.T) {
	db := setupIntegrationDB(t)
	historyService, databaseService, _ := setupHistoryServices(db)

	// Create two users
	user1 := createTestUser(db)
	
	roleID := uint(2)
	user2 := &models.User{
		Firstname: "User",
		Lastname:  "Two",
		Email:     "user2@example.com",
		Password:  "hashedpassword",
		Active:    true,
		RoleID:    &roleID,
	}
	db.Create(user2)

	// User 1 creates databases
	for i := 1; i <= 3; i++ {
		database := &models.Database{
			Name:     "User1 DB " + string(rune(i+'0')),
			Type:     "mysql",
			Host:     "localhost",
			Port:     "3306",
			Username: "user1",
			Password: "pass",
			DbName:   "user1db",
			UserId:   user1.Id,
		}
		err := databaseService.CreateDatabase(database, user1.Id, "127.0.0.1", "test-agent")
		require.NoError(t, err, "Database creation should succeed")
	}

	// User 2 creates databases
	for i := 1; i <= 2; i++ {
		database := &models.Database{
			Name:     "User2 DB " + string(rune(i+'0')),
			Type:     "postgresql",
			Host:     "localhost",
			Port:     "5432",
			Username: "user2",
			Password: "pass",
			DbName:   "user2db",
			UserId:   user2.Id,
		}
		err := databaseService.CreateDatabase(database, user2.Id, "127.0.0.1", "test-agent")
		require.NoError(t, err, "Database creation should succeed")
	}

	// Verify User 1 only sees their own history
	_, total, err := historyService.GetUserActionHistory(user1.Id, 1, 10)
	require.NoError(t, err, "Should retrieve user1 history")
	assert.Equal(t, int64(3), total, "User 1 should have 3 actions")

	// Verify User 2 only sees their own history
	user2Histories, total, err := historyService.GetUserActionHistory(user2.Id, 1, 10)
	require.NoError(t, err, "Should retrieve user2 history")
	assert.Equal(t, int64(2), total, "User 2 should have 2 actions")
	for _, h := range user2Histories {
		assert.Equal(t, user2.Id, h.UserId, "All actions should belong to User 2")
	}
}

// TestIntegration_ActionHistoryJSONSerialization tests JSON metadata serialization
func TestIntegration_ActionHistoryJSONSerialization(t *testing.T) {
	db := setupIntegrationDB(t)
	historyRepo := repositories.NewActionHistoryRepository(db)
	historyService := services.NewActionHistoryService(historyRepo)
	user := createTestUser(db)

	// Create metadata with various data types
	metadata := map[string]interface{}{
		"string_field": "test value",
		"int_field":    42,
		"float_field":  3.14159,
		"bool_field":   true,
		"null_field":   nil,
		"array_field":  []string{"a", "b", "c"},
		"object_field": map[string]interface{}{
			"nested": "value",
		},
	}

	// Log action
	err := historyService.LogAction(user.Id, "test", "resource", 1, "JSON test", metadata, "127.0.0.1", "test")
	require.NoError(t, err, "Should log action")

	// Retrieve from database directly to check JSON storage
	var history models.ActionHistory
	err = db.Where("user_id = ?", user.Id).First(&history).Error
	require.NoError(t, err, "Should retrieve history from DB")

	// Verify JSON can be unmarshaled
	var storedMetadata map[string]interface{}
	err = json.Unmarshal([]byte(history.Metadata), &storedMetadata)
	require.NoError(t, err, "Should unmarshal JSON metadata")
	
	// Verify values
	assert.Equal(t, "test value", storedMetadata["string_field"])
	assert.Equal(t, float64(42), storedMetadata["int_field"]) // JSON numbers are float64
	assert.InDelta(t, 3.14159, storedMetadata["float_field"].(float64), 0.00001)
	assert.Equal(t, true, storedMetadata["bool_field"])
}

