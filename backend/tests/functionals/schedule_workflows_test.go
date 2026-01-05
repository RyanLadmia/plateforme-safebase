package functionals

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// FUNCTIONAL TESTS - Schedule Management Workflows
// ============================================================================

// TestFunctional_ScheduleManagementWorkflow tests complete schedule management
func TestFunctional_ScheduleManagementWorkflow(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.Cleanup(t)
	suite.ScheduleService.StartScheduler()

	t.Log("=== Starting Schedule Management Workflow Test ===")

	// Setup user and database
	user, err := suite.RegisterUser(t, "scheduler@company.com", "ScheduleP@ss123", "Schedule", "Manager")
	require.NoError(t, err, "User registration should succeed")

	database, err := suite.CreateDatabase(t, user.Id, "Scheduled Database", "postgresql")
	require.NoError(t, err, "Database creation should succeed")

	t.Run("Create Multiple Schedules", func(t *testing.T) {
		schedules := []struct {
			name string
			cron string
		}{
			{"Hourly Backup", "0 * * * *"},
			{"Daily Backup at 2 AM", "0 2 * * *"},
			{"Weekly Backup on Sunday", "0 0 * * 0"},
			{"Monthly Backup on 1st", "0 0 1 * *"},
		}

		for _, s := range schedules {
			schedule, err := suite.CreateSchedule(t, database.Id, user.Id, s.name, s.cron)
			require.NoError(t, err, "Schedule creation should succeed: "+s.name)
			assert.True(t, schedule.Active, "Schedule should be active by default")
		}

		// Verify all schedules were created
		userSchedules, err := suite.ScheduleService.GetSchedules(user.Id)
		require.NoError(t, err, "Should retrieve user schedules")
		assert.GreaterOrEqual(t, len(userSchedules), 4, "User should have at least 4 schedules")
	})

	t.Run("Update Schedule Frequency", func(t *testing.T) {
		schedules, err := suite.ScheduleService.GetSchedules(user.Id)
		require.NoError(t, err, "Should retrieve schedules")
		require.NotEmpty(t, schedules, "User should have schedules")

		// Update first schedule
		schedule := schedules[0]
		updated, err := suite.ScheduleService.UpdateSchedule(
			schedule.Id,
			user.Id,
			"Modified Hourly Backup",
			"*/30 * * * *", // Every 30 minutes
			nil,
			"127.0.0.1",
			"functional-test",
		)
		require.NoError(t, err, "Schedule update should succeed")
		assert.Equal(t, "Modified Hourly Backup", updated.Name, "Name should be updated")
		assert.Equal(t, "*/30 * * * *", updated.CronExpression, "Cron should be updated")
	})

	t.Run("Activate and Deactivate Schedules", func(t *testing.T) {
		schedules, err := suite.ScheduleService.GetSchedules(user.Id)
		require.NoError(t, err, "Should retrieve schedules")
		require.NotEmpty(t, schedules, "User should have schedules")

		// Deactivate first schedule
		inactive := false
		deactivated, err := suite.ScheduleService.UpdateSchedule(
			schedules[0].Id,
			user.Id,
			"",
			"",
			&inactive,
			"127.0.0.1",
			"functional-test",
		)
		require.NoError(t, err, "Deactivation should succeed")
		assert.False(t, deactivated.Active, "Schedule should be inactive")

		// Reactivate schedule
		active := true
		reactivated, err := suite.ScheduleService.UpdateSchedule(
			schedules[0].Id,
			user.Id,
			"",
			"",
			&active,
			"127.0.0.1",
			"functional-test",
		)
		require.NoError(t, err, "Reactivation should succeed")
		assert.True(t, reactivated.Active, "Schedule should be active")
	})

	t.Run("Schedule Validation", func(t *testing.T) {
		// Test invalid cron expressions
		invalidExpressions := []string{
			"invalid cron",
			"60 * * * *",  // Invalid minute
			"* * * * * *", // Too many fields
			"",            // Empty
		}

		for _, expr := range invalidExpressions {
			_, err := suite.CreateSchedule(t, database.Id, user.Id, "Invalid Schedule", expr)
			assert.Error(t, err, "Invalid cron expression should be rejected: "+expr)
		}
	})

	t.Run("Delete All Schedules", func(t *testing.T) {
		schedules, err := suite.ScheduleService.GetSchedules(user.Id)
		require.NoError(t, err, "Should retrieve schedules")

		initialCount := len(schedules)
		assert.Greater(t, initialCount, 0, "Should have schedules to delete")

		// Delete all schedules
		for _, schedule := range schedules {
			err := suite.ScheduleService.DeleteSchedule(schedule.Id, user.Id, "127.0.0.1", "functional-test")
			require.NoError(t, err, "Schedule deletion should succeed")
		}

		// Verify all deleted
		remainingSchedules, err := suite.ScheduleService.GetSchedules(user.Id)
		require.NoError(t, err, "Should retrieve schedules")
		assert.Equal(t, 0, len(remainingSchedules), "All schedules should be deleted")
	})

	t.Log("=== Schedule Management Workflow Test Finished Successfully ===")
}

// TestFunctional_ComplexSchedulingScenario tests realistic scheduling scenarios
func TestFunctional_ComplexSchedulingScenario(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.Cleanup(t)
	suite.ScheduleService.StartScheduler()

	t.Log("=== Starting Complex Scheduling Scenario Test ===")

	// Create company with multiple databases
	user, err := suite.RegisterUser(t, "devops@startup.com", "DevOpsP@ss123", "DevOps", "Team")
	require.NoError(t, err, "User registration should succeed")

	// Create databases for different environments
	environments := []struct {
		name   string
		dbType string
	}{
		{"Production MySQL", "mysql"},
		{"Staging PostgreSQL", "postgresql"},
		{"Development MySQL", "mysql"},
		{"Analytics PostgreSQL", "postgresql"},
	}

	databaseIDs := make([]uint, len(environments))
	for i, env := range environments {
		db, err := suite.CreateDatabase(t, user.Id, env.name, env.dbType)
		require.NoError(t, err, "Database creation should succeed")
		databaseIDs[i] = db.Id
	}

	t.Run("Create Environment-Specific Backup Schedules", func(t *testing.T) {
		// Production: Hourly backups
		_, err := suite.CreateSchedule(t, databaseIDs[0], user.Id, "Production Hourly", "0 * * * *")
		require.NoError(t, err, "Production schedule should succeed")

		// Staging: Every 6 hours
		_, err = suite.CreateSchedule(t, databaseIDs[1], user.Id, "Staging 6-Hourly", "0 */6 * * *")
		require.NoError(t, err, "Staging schedule should succeed")

		// Development: Daily at midnight
		_, err = suite.CreateSchedule(t, databaseIDs[2], user.Id, "Development Daily", "0 0 * * *")
		require.NoError(t, err, "Development schedule should succeed")

		// Analytics: Weekly on Sunday
		_, err = suite.CreateSchedule(t, databaseIDs[3], user.Id, "Analytics Weekly", "0 0 * * 0")
		require.NoError(t, err, "Analytics schedule should succeed")
	})

	t.Run("Verify Schedules Per Database", func(t *testing.T) {
		allSchedules, err := suite.ScheduleService.GetSchedules(user.Id)
		require.NoError(t, err, "Should retrieve schedules")
		assert.Equal(t, 4, len(allSchedules), "Should have 4 schedules (one per database)")

		// Verify each database has exactly one schedule
		schedulesByDB := make(map[uint]int)
		for _, s := range allSchedules {
			schedulesByDB[s.DatabaseId]++
		}

		for _, dbID := range databaseIDs {
			assert.Equal(t, 1, schedulesByDB[dbID], "Each database should have exactly 1 schedule")
		}
	})

	t.Run("Simulate Production Issue - Disable Backups", func(t *testing.T) {
		// Find production schedule
		allSchedules, err := suite.ScheduleService.GetSchedules(user.Id)
		require.NoError(t, err, "Should retrieve schedules")

		var prodScheduleID uint
		for _, s := range allSchedules {
			if s.DatabaseId == databaseIDs[0] { // Production database
				prodScheduleID = s.Id
				break
			}
		}

		require.NotZero(t, prodScheduleID, "Should find production schedule")

		// Disable production backups
		inactive := false
		updated, err := suite.ScheduleService.UpdateSchedule(
			prodScheduleID,
			user.Id,
			"",
			"",
			&inactive,
			"127.0.0.1",
			"functional-test",
		)
		require.NoError(t, err, "Should deactivate production schedule")
		assert.False(t, updated.Active, "Production schedule should be inactive")

		// Verify other schedules remain active
		activeSchedules, err := suite.ScheduleRepo.GetActive()
		require.NoError(t, err, "Should retrieve active schedules")
		
		for _, s := range activeSchedules {
			assert.NotEqual(t, prodScheduleID, s.Id, "Production schedule should not be in active list")
		}
	})

	t.Run("Verify Action History for All Operations", func(t *testing.T) {
		histories, total, err := suite.ActionHistoryService.GetUserActionHistoryByType(user.Id, "schedule", 1, 20)
		require.NoError(t, err, "Should retrieve schedule history")
		assert.Greater(t, int(total), 4, "Should have multiple schedule actions logged")

		// Count action types
		actionCounts := make(map[string]int)
		for _, h := range histories {
			actionCounts[h.Action]++
		}

		assert.Greater(t, actionCounts["create"], 0, "Should have create actions")
		assert.Greater(t, actionCounts["update"], 0, "Should have update actions")
	})

	t.Log("=== Complex Scheduling Scenario Test Finished Successfully ===")
}

// TestFunctional_ScheduleAndDatabaseDeletion tests cascade deletion scenarios
func TestFunctional_ScheduleAndDatabaseDeletion(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.Cleanup(t)
	suite.ScheduleService.StartScheduler()

	t.Log("=== Starting Schedule and Database Deletion Test ===")

	user, err := suite.RegisterUser(t, "cleanup@company.com", "CleanupP@ss123", "Cleanup", "Tester")
	require.NoError(t, err, "User registration should succeed")

	// Create database with multiple schedules
	database, err := suite.CreateDatabase(t, user.Id, "Database with Schedules", "mysql")
	require.NoError(t, err, "Database creation should succeed")

	scheduleNames := []string{"Schedule 1", "Schedule 2", "Schedule 3"}
	scheduleIDs := make([]uint, len(scheduleNames))

	t.Run("Create Multiple Schedules for Same Database", func(t *testing.T) {
		for i, name := range scheduleNames {
			schedule, err := suite.CreateSchedule(t, database.Id, user.Id, name, "0 * * * *")
			require.NoError(t, err, "Schedule creation should succeed")
			scheduleIDs[i] = schedule.Id
		}
	})

	t.Run("Delete Database Should Soft Delete Associated Schedules", func(t *testing.T) {
		// Delete database
		err := suite.DatabaseService.DeleteDatabase(database.Id, user.Id, "127.0.0.1", "functional-test")
		require.NoError(t, err, "Database deletion should succeed")

		// Verify database is soft deleted
		_, err = suite.DatabaseService.GetDatabaseByID(database.Id)
		assert.Error(t, err, "Deleted database should not be retrievable")

		// Verify schedules are also soft deleted
		for _, scheduleID := range scheduleIDs {
			_, err := suite.ScheduleService.GetSchedule(scheduleID, user.Id)
			assert.Error(t, err, "Schedule should not be retrievable after database deletion")
		}

		// Verify no active schedules for this user
		activeSchedules, err := suite.ScheduleService.GetSchedules(user.Id)
		require.NoError(t, err, "Should retrieve schedules")
		assert.Equal(t, 0, len(activeSchedules), "Should have no active schedules")
	})

	t.Log("=== Schedule and Database Deletion Test Finished Successfully ===")
}

// TestFunctional_CronExpressionEdgeCases tests various cron expression scenarios
func TestFunctional_CronExpressionEdgeCases(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.Cleanup(t)
	suite.ScheduleService.StartScheduler()

	t.Log("=== Starting Cron Expression Edge Cases Test ===")

	user, err := suite.RegisterUser(t, "cron@company.com", "CronP@ss123", "Cron", "Expert")
	require.NoError(t, err, "User registration should succeed")

	database, err := suite.CreateDatabase(t, user.Id, "Cron Test DB", "mysql")
	require.NoError(t, err, "Database creation should succeed")

	t.Run("Valid Complex Cron Expressions", func(t *testing.T) {
		validExpressions := []struct {
			name string
			cron string
		}{
			{"Every 15 minutes", "*/15 * * * *"},
			{"Weekdays at 9 AM", "0 9 * * 1-5"},
			{"First day of month at midnight", "0 0 1 * *"},
			{"Every 2 hours during business hours", "0 9-17/2 * * *"},
			{"Sunday and Saturday at noon", "0 12 * * 0,6"},
		}

		for _, expr := range validExpressions {
			schedule, err := suite.CreateSchedule(t, database.Id, user.Id, expr.name, expr.cron)
			require.NoError(t, err, "Valid cron should be accepted: "+expr.cron)
			assert.NotZero(t, schedule.Id, "Schedule should be created")
		}
	})

	t.Run("Invalid Cron Expressions Should Fail", func(t *testing.T) {
		invalidExpressions := []struct {
			name   string
			cron   string
			reason string
		}{
			{"Empty", "", "Empty expression"},
			{"Too many fields", "* * * * * *", "6 fields instead of 5"},
			{"Invalid minute", "60 * * * *", "Minute out of range"},
			{"Invalid hour", "0 25 * * *", "Hour out of range"},
			{"Invalid day", "0 0 32 * *", "Day out of range"},
			{"Invalid month", "0 0 1 13 *", "Month out of range"},
			{"Invalid weekday", "0 0 * * 8", "Weekday out of range"},
			{"Random text", "every day", "Not a cron expression"},
		}

		for _, expr := range invalidExpressions {
			_, err := suite.CreateSchedule(t, database.Id, user.Id, expr.name, expr.cron)
			assert.Error(t, err, "Invalid cron should be rejected: "+expr.reason)
		}
	})

	t.Run("Update to Invalid Cron Should Fail", func(t *testing.T) {
		// Create valid schedule
		schedule, err := suite.CreateSchedule(t, database.Id, user.Id, "Valid Schedule", "0 0 * * *")
		require.NoError(t, err, "Valid schedule creation should succeed")

		// Try to update to invalid cron
		_, err = suite.ScheduleService.UpdateSchedule(
			schedule.Id,
			user.Id,
			"",
			"invalid cron",
			nil,
			"127.0.0.1",
			"functional-test",
		)
		assert.Error(t, err, "Update to invalid cron should fail")

		// Verify schedule remains unchanged
		unchanged, err := suite.ScheduleService.GetSchedule(schedule.Id, user.Id)
		require.NoError(t, err, "Should retrieve schedule")
		assert.Equal(t, "0 0 * * *", unchanged.CronExpression, "Cron should remain unchanged")
	})

	t.Log("=== Cron Expression Edge Cases Test Finished Successfully ===")
}

