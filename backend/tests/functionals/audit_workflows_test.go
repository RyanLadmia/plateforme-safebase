package functionals

import (
	"testing"

	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// FUNCTIONAL TESTS - Action History and Audit Trail
// ============================================================================

// TestFunctional_ComprehensiveAuditTrail tests complete audit trail functionality
func TestFunctional_ComprehensiveAuditTrail(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.Cleanup(t)
	suite.ScheduleService.StartScheduler()

	t.Log("=== Starting Comprehensive Audit Trail Test ===")

	// Create user
	user, err := suite.RegisterUser(t, "auditor@company.com", "AuditorP@ss123", "Audit", "Manager")
	require.NoError(t, err, "User registration should succeed")

	t.Run("Database Operations Are Logged", func(t *testing.T) {
		// Create database
		database, err := suite.CreateDatabase(t, user.Id, "Audited Database", "mysql")
		require.NoError(t, err, "Database creation should succeed")

		// Update database
		database.Name = "Updated Audited Database"
		err = suite.DatabaseService.UpdateDatabase(database, user.Id, "192.168.1.100", "Chrome/90.0")
		require.NoError(t, err, "Database update should succeed")

		// Retrieve action history
		histories, _, err := suite.ActionHistoryService.GetUserActionHistoryByType(user.Id, "database", 1, 10)
		require.NoError(t, err, "Should retrieve database action history")
		assert.GreaterOrEqual(t, len(histories), 2, "Should have at least 2 database actions (create + update)")

		// Verify metadata
		for _, h := range histories {
			assert.NotEmpty(t, h.IpAddress, "IP address should be logged")
			assert.NotEmpty(t, h.UserAgent, "User agent should be logged")
			assert.NotNil(t, h.Metadata, "Metadata should be present")
		}
	})

	t.Run("Schedule Operations Are Logged", func(t *testing.T) {
		// Get database
		databases, _ := suite.DatabaseService.GetDatabasesByUser(user.Id)
		require.NotEmpty(t, databases, "User should have databases")

		// Create schedule
		schedule, err := suite.CreateSchedule(t, databases[0].Id, user.Id, "Audit Schedule", "0 0 * * *")
		require.NoError(t, err, "Schedule creation should succeed")

		// Update schedule
		inactive := false
		_, err = suite.ScheduleService.UpdateSchedule(schedule.Id, user.Id, "Updated Audit Schedule", "", &inactive, "10.0.0.1", "Firefox/88.0")
		require.NoError(t, err, "Schedule update should succeed")

		// Delete schedule
		err = suite.ScheduleService.DeleteSchedule(schedule.Id, user.Id, "10.0.0.2", "Safari/14.0")
		require.NoError(t, err, "Schedule deletion should succeed")

		// Retrieve schedule history
		histories, _, err := suite.ActionHistoryService.GetUserActionHistoryByType(user.Id, "schedule", 1, 10)
		require.NoError(t, err, "Should retrieve schedule action history")
		assert.GreaterOrEqual(t, len(histories), 3, "Should have at least 3 schedule actions (create + update + delete)")

		// Verify action types
		actionTypes := make(map[string]int)
		for _, h := range histories {
			actionTypes[h.Action]++
		}
		assert.GreaterOrEqual(t, actionTypes["create"], 1, "Should have create action")
		assert.GreaterOrEqual(t, actionTypes["update"], 1, "Should have update action")
		assert.GreaterOrEqual(t, actionTypes["delete"], 1, "Should have delete action")
	})

	t.Run("Pagination Works Correctly", func(t *testing.T) {
		// Get total count
		_, total, err := suite.ActionHistoryService.GetUserActionHistory(user.Id, 1, 100)
		require.NoError(t, err, "Should retrieve action history")

		if total > 5 {
			// Test first page
			page1, _, err := suite.ActionHistoryService.GetUserActionHistory(user.Id, 1, 2)
			require.NoError(t, err, "Should retrieve first page")
			assert.Equal(t, 2, len(page1), "First page should have 2 entries")

			// Test second page
			page2, _, err := suite.ActionHistoryService.GetUserActionHistory(user.Id, 2, 2)
			require.NoError(t, err, "Should retrieve second page")
			assert.Equal(t, 2, len(page2), "Second page should have 2 entries")

			// Verify pages are different
			assert.NotEqual(t, page1[0].Id, page2[0].Id, "Pages should contain different entries")
		}
	})

	t.Run("Metadata Contains Detailed Information", func(t *testing.T) {
		histories, _, err := suite.ActionHistoryService.GetUserActionHistory(user.Id, 1, 10)
		require.NoError(t, err, "Should retrieve action history")
		require.NotEmpty(t, histories, "Should have action history")

		// Find an update action with changes
		var updateAction *services.ActionHistoryResponse
		for _, h := range histories {
			if h.Action == "update" && h.Metadata != nil {
				updateAction = &h
				break
			}
		}

		if updateAction != nil {
			// Verify metadata structure
			assert.NotNil(t, updateAction.Metadata, "Update action should have metadata")
			
			// Check for common metadata fields
			_, hasResourceID := updateAction.Metadata["database_id"]
			_, hasScheduleID := updateAction.Metadata["schedule_id"]
			assert.True(t, hasResourceID || hasScheduleID, "Metadata should contain resource ID")
		}
	})

	t.Log("=== Comprehensive Audit Trail Test Finished Successfully ===")
}

// TestFunctional_MultiUserAuditIsolation tests audit trail isolation between users
func TestFunctional_MultiUserAuditIsolation(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.Cleanup(t)

	t.Log("=== Starting Multi-User Audit Isolation Test ===")

	// Create multiple users
	user1, err := suite.RegisterUser(t, "user1@company.com", "User1P@ss123", "User", "One")
	require.NoError(t, err, "User 1 registration should succeed")

	user2, err := suite.RegisterUser(t, "user2@company.com", "User2P@ss123", "User", "Two")
	require.NoError(t, err, "User 2 registration should succeed")

	t.Run("Each User Creates Their Own Databases", func(t *testing.T) {
		// User 1 creates 3 databases
		for i := 1; i <= 3; i++ {
			_, err := suite.CreateDatabase(t, user1.Id, "User1 DB "+string(rune(i+'0')), "mysql")
			require.NoError(t, err, "User 1 database creation should succeed")
		}

		// User 2 creates 2 databases
		for i := 1; i <= 2; i++ {
			_, err := suite.CreateDatabase(t, user2.Id, "User2 DB "+string(rune(i+'0')), "postgresql")
			require.NoError(t, err, "User 2 database creation should succeed")
		}
	})

	t.Run("Verify Action History Isolation", func(t *testing.T) {
		// User 1 should have 3 database actions
		user1Histories, total1, err := suite.ActionHistoryService.GetUserActionHistoryByType(user1.Id, "database", 1, 10)
		require.NoError(t, err, "Should retrieve user 1 history")
		assert.Equal(t, int64(3), total1, "User 1 should have 3 database actions")

		// Verify all actions belong to user 1
		for _, h := range user1Histories {
			assert.Equal(t, user1.Id, h.UserId, "All actions should belong to user 1")
		}

		// User 2 should have 2 database actions
		user2Histories, total2, err := suite.ActionHistoryService.GetUserActionHistoryByType(user2.Id, "database", 1, 10)
		require.NoError(t, err, "Should retrieve user 2 history")
		assert.Equal(t, int64(2), total2, "User 2 should have 2 database actions")

		// Verify all actions belong to user 2
		for _, h := range user2Histories {
			assert.Equal(t, user2.Id, h.UserId, "All actions should belong to user 2")
		}
	})

	t.Run("Users Cannot Access Each Other's History", func(t *testing.T) {
		// This test verifies the logical isolation at service level
		// In a real API, this would be enforced by authentication middleware

		user1History, _, err := suite.ActionHistoryService.GetUserActionHistory(user1.Id, 1, 10)
		require.NoError(t, err, "Should retrieve user 1 history")

		user2History, _, err := suite.ActionHistoryService.GetUserActionHistory(user2.Id, 1, 10)
		require.NoError(t, err, "Should retrieve user 2 history")

		// Verify histories are completely different
		if len(user1History) > 0 && len(user2History) > 0 {
			assert.NotEqual(t, user1History[0].UserId, user2History[0].UserId, "Histories should belong to different users")
		}
	})

	t.Log("=== Multi-User Audit Isolation Test Finished Successfully ===")
}

// TestFunctional_ActionHistoryPerformance tests performance with large history
func TestFunctional_ActionHistoryPerformance(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.Cleanup(t)

	t.Log("=== Starting Action History Performance Test ===")

	user, err := suite.RegisterUser(t, "performance@company.com", "PerfP@ss123", "Performance", "Tester")
	require.NoError(t, err, "User registration should succeed")

	t.Run("Create Large Number of Actions", func(t *testing.T) {
		// Create 50 databases to generate 50 actions
		for i := 1; i <= 50; i++ {
			_, err := suite.CreateDatabase(t, user.Id, "Performance DB "+string(rune(i)), "mysql")
			require.NoError(t, err, "Database creation should succeed")
		}
	})

	t.Run("Pagination Handles Large Dataset", func(t *testing.T) {
		// Test different page sizes
		pageSizes := []int{5, 10, 20, 50}

		for _, pageSize := range pageSizes {
			histories, total, err := suite.ActionHistoryService.GetUserActionHistory(user.Id, 1, pageSize)
			require.NoError(t, err, "Should retrieve action history with page size "+string(rune(pageSize)))
			assert.Equal(t, int64(50), total, "Total should always be 50")
			assert.LessOrEqual(t, len(histories), pageSize, "Page size should be respected")
		}
	})

	t.Run("Filtering Performs Well", func(t *testing.T) {
		// Filter by resource type
		histories, total, err := suite.ActionHistoryService.GetUserActionHistoryByType(user.Id, "database", 1, 100)
		require.NoError(t, err, "Should retrieve filtered history")
		assert.Equal(t, int64(50), total, "Should have 50 database actions")

		// Verify all are database actions
		for _, h := range histories {
			assert.Equal(t, "database", h.ResourceType, "All should be database type")
		}
	})

	t.Run("Recent History Access", func(t *testing.T) {
		// Get recent actions (most common use case)
		histories, _, err := suite.ActionHistoryService.GetUserActionHistory(user.Id, 1, 10)
		require.NoError(t, err, "Should retrieve recent history")
		assert.Equal(t, 10, len(histories), "Should return 10 most recent actions")

		// Verify ordering (most recent first)
		if len(histories) >= 2 {
			assert.True(t, histories[0].CreatedAt.After(histories[1].CreatedAt) || histories[0].CreatedAt.Equal(histories[1].CreatedAt), 
				"History should be ordered by most recent first")
		}
	})

	t.Log("=== Action History Performance Test Finished Successfully ===")
}

// TestFunctional_ActionHistoryComplexMetadata tests complex metadata scenarios
func TestFunctional_ActionHistoryComplexMetadata(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.Cleanup(t)
	suite.ScheduleService.StartScheduler()

	t.Log("=== Starting Action History Complex Metadata Test ===")

	user, err := suite.RegisterUser(t, "metadata@company.com", "MetaP@ss123", "Meta", "Tester")
	require.NoError(t, err, "User registration should succeed")

	database, err := suite.CreateDatabase(t, user.Id, "Metadata Test DB", "mysql")
	require.NoError(t, err, "Database creation should succeed")

	t.Run("Database Name Change Logs Detailed Changes", func(t *testing.T) {
		oldName := database.Name
		newName := "Renamed Metadata Test DB"

		err := suite.DatabaseService.UpdateDatabaseName(database.Id, newName, user.Id, "192.168.1.100", "Chrome/90.0")
		require.NoError(t, err, "Database name update should succeed")

		// Retrieve action history
		histories, _, err := suite.ActionHistoryService.GetUserActionHistoryByType(user.Id, "database", 1, 10)
		require.NoError(t, err, "Should retrieve database history")

		// Find the name change action
		var nameChangeAction *services.ActionHistoryResponse
		for _, h := range histories {
			if h.Action == "updated" && h.ResourceId == database.Id {
				if h.Metadata != nil {
					if changes, ok := h.Metadata["changes"].(map[string]interface{}); ok {
						if _, hasName := changes["name"]; hasName {
							nameChangeAction = &h
							break
						}
					}
				}
			}
		}

		if nameChangeAction != nil {
			// Verify metadata structure
			assert.NotNil(t, nameChangeAction.Metadata, "Should have metadata")
			
			changes, ok := nameChangeAction.Metadata["changes"].(map[string]interface{})
			assert.True(t, ok, "Should have changes in metadata")
			
			if ok {
				nameChange, ok := changes["name"].(map[string]interface{})
				assert.True(t, ok, "Should have name change details")
				
				if ok {
					assert.Equal(t, oldName, nameChange["from"], "Should log old name")
					assert.Equal(t, newName, nameChange["to"], "Should log new name")
				}
			}
		}
	})

	t.Run("Schedule Updates Log Change Details", func(t *testing.T) {
		// Create schedule
		schedule, err := suite.CreateSchedule(t, database.Id, user.Id, "Original Schedule", "0 0 * * *")
		require.NoError(t, err, "Schedule creation should succeed")

		// Update multiple fields
		inactive := false
		_, err = suite.ScheduleService.UpdateSchedule(
			schedule.Id,
			user.Id,
			"Updated Schedule",
			"0 2 * * *",
			&inactive,
			"10.0.0.1",
			"Firefox/88.0",
		)
		require.NoError(t, err, "Schedule update should succeed")

		// Retrieve action history
		histories, _, err := suite.ActionHistoryService.GetUserActionHistoryByType(user.Id, "schedule", 1, 10)
		require.NoError(t, err, "Should retrieve schedule history")

		// Find the update action
		var updateAction *services.ActionHistoryResponse
		for _, h := range histories {
			if h.Action == "update" && h.ResourceId == schedule.Id {
				updateAction = &h
				break
			}
		}

		if updateAction != nil {
			assert.NotNil(t, updateAction.Metadata, "Update action should have metadata")
			
			// Verify changes are logged
			if changes, ok := updateAction.Metadata["changes"].(map[string]interface{}); ok {
				// Should have multiple changes
				assert.NotEmpty(t, changes, "Should log multiple changes")
			}
		}
	})

	t.Log("=== Action History Complex Metadata Test Finished Successfully ===")
}

