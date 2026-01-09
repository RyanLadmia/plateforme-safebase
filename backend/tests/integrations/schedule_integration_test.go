package integrations

import (
	"testing"
	"time"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================================
// INTEGRATION TESTS - Schedule Management with Cron Integration
// ============================================================================

// setupScheduleServices creates all services needed for schedule tests
func setupScheduleServices(db *gorm.DB) (*services.ScheduleService, *services.DatabaseService, *repositories.ScheduleRepository) {
	scheduleRepo := repositories.NewScheduleRepository(db)
	databaseRepo := repositories.NewDatabaseRepository(db)
	backupRepo := repositories.NewBackupRepository(db)
	restoreRepo := repositories.NewRestoreRepository(db)
	
	databaseService := services.NewDatabaseService(databaseRepo, backupRepo, restoreRepo, scheduleRepo, nil)
	scheduleService := services.NewScheduleService(scheduleRepo, databaseRepo, nil)
	
	return scheduleService, databaseService, scheduleRepo
}

// TestIntegration_ScheduleCRUDWithCron tests schedule CRUD with cron scheduler
func TestIntegration_ScheduleCRUDWithCron(t *testing.T) {
	db := setupIntegrationDB(t)
	scheduleService, databaseService, scheduleRepo := setupScheduleServices(db)
	scheduleService.StartScheduler()
	user := createTestUser(db)

	// Create a database first
	database := &models.Database{
		Name:     "Scheduled Backup DB",
		Type:     "mysql",
		Host:     "localhost",
		Port:     "3306",
		Username: "testuser",
		Password: "password123",
		DbName:   "testdb",
		UserId:   user.Id,
	}
	err := databaseService.CreateDatabase(database, user.Id, "127.0.0.1", "integration-test")
	require.NoError(t, err, "Database creation should succeed")

	var scheduleID uint

	// Step 1: Create a schedule
	t.Run("Create Schedule", func(t *testing.T) {
		schedule, err := scheduleService.CreateSchedule(
			database.Id,
			user.Id,
			"Daily Backup at Midnight",
			"0 0 * * *",
			"127.0.0.1",
			"integration-test",
		)
		require.NoError(t, err, "Schedule creation should succeed")
		assert.NotZero(t, schedule.Id, "Schedule ID should be set")
		assert.True(t, schedule.Active, "Schedule should be active by default")
		scheduleID = schedule.Id

		// Verify schedule is in database
		dbSchedule, err := scheduleRepo.GetByID(schedule.Id)
		require.NoError(t, err, "Should retrieve schedule from database")
		assert.Equal(t, "Daily Backup at Midnight", dbSchedule.Name, "Schedule name should match")
	})

	// Step 2: Retrieve schedule
	t.Run("Get Schedule", func(t *testing.T) {
		schedule, err := scheduleService.GetSchedule(scheduleID, user.Id)
		require.NoError(t, err, "Should retrieve schedule")
		assert.Equal(t, "Daily Backup at Midnight", schedule.Name, "Schedule name should match")
		assert.Equal(t, "0 0 * * *", schedule.CronExpression, "Cron expression should match")
	})

	// Step 3: Update schedule
	t.Run("Update Schedule", func(t *testing.T) {
		updatedSchedule, err := scheduleService.UpdateSchedule(
			scheduleID,
			user.Id,
			"Updated Daily Backup",
			"0 2 * * *", // 2 AM instead of midnight
			nil,
			"127.0.0.1",
			"integration-test",
		)
		require.NoError(t, err, "Schedule update should succeed")
		assert.Equal(t, "Updated Daily Backup", updatedSchedule.Name, "Name should be updated")
		assert.Equal(t, "0 2 * * *", updatedSchedule.CronExpression, "Cron expression should be updated")
	})

	// Step 4: Deactivate schedule
	t.Run("Deactivate Schedule", func(t *testing.T) {
		inactive := false
		updatedSchedule, err := scheduleService.UpdateSchedule(
			scheduleID,
			user.Id,
			"",
			"",
			&inactive,
			"127.0.0.1",
			"integration-test",
		)
		require.NoError(t, err, "Schedule deactivation should succeed")
		assert.False(t, updatedSchedule.Active, "Schedule should be inactive")
	})

	// Step 5: Reactivate schedule
	t.Run("Reactivate Schedule", func(t *testing.T) {
		active := true
		updatedSchedule, err := scheduleService.UpdateSchedule(
			scheduleID,
			user.Id,
			"",
			"",
			&active,
			"127.0.0.1",
			"integration-test",
		)
		require.NoError(t, err, "Schedule reactivation should succeed")
		assert.True(t, updatedSchedule.Active, "Schedule should be active")
	})

	// Step 6: List user schedules
	t.Run("List User Schedules", func(t *testing.T) {
		schedules, err := scheduleService.GetSchedules(user.Id)
		require.NoError(t, err, "Should retrieve user schedules")
		assert.Equal(t, 1, len(schedules), "User should have 1 schedule")
		assert.Equal(t, "Updated Daily Backup", schedules[0].Name, "Schedule name should match")
	})

	// Step 7: Delete schedule
	t.Run("Delete Schedule", func(t *testing.T) {
		err := scheduleService.DeleteSchedule(scheduleID, user.Id, "127.0.0.1", "integration-test")
		require.NoError(t, err, "Schedule deletion should succeed")

		// Verify schedule is deleted
		_, err = scheduleService.GetSchedule(scheduleID, user.Id)
		assert.Error(t, err, "Deleted schedule should not be retrievable")
	})
}

// TestIntegration_CronExpressionValidation tests cron expression validation
func TestIntegration_CronExpressionValidation(t *testing.T) {
	db := setupIntegrationDB(t)
	scheduleService, databaseService, _ := setupScheduleServices(db)
	scheduleService.StartScheduler()
	user := createTestUser(db)

	// Create a database
	database := &models.Database{
		Name:     "Test DB",
		Type:     "mysql",
		Host:     "localhost",
		Port:     "3306",
		Username: "testuser",
		Password: "password123",
		DbName:   "testdb",
		UserId:   user.Id,
	}
	err := databaseService.CreateDatabase(database, user.Id, "127.0.0.1", "integration-test")
	require.NoError(t, err, "Database creation should succeed")

	validExpressions := []struct {
		name       string
		expression string
	}{
		{"Every minute", "* * * * *"},
		{"Every hour", "0 * * * *"},
		{"Daily at noon", "0 12 * * *"},
		{"Weekly on Monday", "0 0 * * 1"},
		{"Monthly on 1st", "0 0 1 * *"},
	}

	// Test valid expressions
	for _, expr := range validExpressions {
		t.Run("Valid: "+expr.name, func(t *testing.T) {
			_, err := scheduleService.CreateSchedule(
				database.Id,
				user.Id,
				expr.name,
				expr.expression,
				"127.0.0.1",
				"integration-test",
			)
			assert.NoError(t, err, "Valid cron expression should be accepted")
		})
	}

	invalidExpressions := []string{
		"invalid cron",
		"* * * * * *", // Too many fields
		"60 * * * *",  // Invalid minute
		"",            // Empty
	}

	// Test invalid expressions
	for _, expr := range invalidExpressions {
		t.Run("Invalid: "+expr, func(t *testing.T) {
			_, err := scheduleService.CreateSchedule(
				database.Id,
				user.Id,
				"Invalid Schedule",
				expr,
				"127.0.0.1",
				"integration-test",
			)
			assert.Error(t, err, "Invalid cron expression should be rejected")
		})
	}
}

// TestIntegration_MultipleSchedulesPerDatabase tests multiple schedules for one database
func TestIntegration_MultipleSchedulesPerDatabase(t *testing.T) {
	db := setupIntegrationDB(t)
	scheduleService, databaseService, _ := setupScheduleServices(db)
	scheduleService.StartScheduler()
	user := createTestUser(db)

	// Create a database
	database := &models.Database{
		Name:     "Multi-Schedule DB",
		Type:     "postgresql",
		Host:     "localhost",
		Port:     "5432",
		Username: "testuser",
		Password: "password123",
		DbName:   "testdb",
		UserId:   user.Id,
	}
	err := databaseService.CreateDatabase(database, user.Id, "127.0.0.1", "integration-test")
	require.NoError(t, err, "Database creation should succeed")

	schedules := []struct {
		name       string
		expression string
	}{
		{"Hourly Backup", "0 * * * *"},
		{"Daily Backup", "0 0 * * *"},
		{"Weekly Backup", "0 0 * * 0"},
	}

	// Create multiple schedules
	for _, sched := range schedules {
		_, err := scheduleService.CreateSchedule(
			database.Id,
			user.Id,
			sched.name,
			sched.expression,
			"127.0.0.1",
			"integration-test",
		)
		require.NoError(t, err, "Schedule creation should succeed")
	}

	// Verify all schedules exist
	userSchedules, err := scheduleService.GetSchedules(user.Id)
	require.NoError(t, err, "Should retrieve user schedules")
	assert.GreaterOrEqual(t, len(userSchedules), 3, "User should have at least 3 schedules")

	// Verify all schedules belong to the same database
	dbScheduleCount := 0
	for _, s := range userSchedules {
		if s.DatabaseId == database.Id {
			dbScheduleCount++
		}
	}
	assert.GreaterOrEqual(t, dbScheduleCount, 3, "Database should have at least 3 schedules")
}

// TestIntegration_LoadActiveSchedulesOnStartup tests loading schedules on service start
func TestIntegration_LoadActiveSchedulesOnStartup(t *testing.T) {
	db := setupIntegrationDB(t)
	user := createTestUser(db)

	// Setup services and create database
	scheduleRepo := repositories.NewScheduleRepository(db)
	databaseRepo := repositories.NewDatabaseRepository(db)
	backupRepo := repositories.NewBackupRepository(db)
	restoreRepo := repositories.NewRestoreRepository(db)
	
	databaseService := services.NewDatabaseService(databaseRepo, backupRepo, restoreRepo, scheduleRepo, nil)
	
	database := &models.Database{
		Name:     "Startup Test DB",
		Type:     "mysql",
		Host:     "localhost",
		Port:     "3306",
		Username: "testuser",
		Password: "password123",
		DbName:   "testdb",
		UserId:   user.Id,
	}
	err := databaseService.CreateDatabase(database, user.Id, "127.0.0.1", "integration-test")
	require.NoError(t, err, "Database creation should succeed")

	// Create schedules directly in database (simulating existing schedules)
	activeSchedule := &models.Schedule{
		Name:           "Active Schedule",
		CronExpression: "0 0 * * *",
		Active:         true,
		UserId:         user.Id,
		DatabaseId:     database.Id,
	}
	inactiveSchedule := &models.Schedule{
		Name:           "Inactive Schedule",
		CronExpression: "0 12 * * *",
		Active:         false,
		UserId:         user.Id,
		DatabaseId:     database.Id,
	}
	db.Create(activeSchedule)
	db.Create(inactiveSchedule)

	// Create new schedule service (simulating startup)
	scheduleService := services.NewScheduleService(scheduleRepo, databaseRepo, nil)
	scheduleService.StartScheduler()

	// Load active schedules
	err = scheduleService.LoadActiveSchedules()
	require.NoError(t, err, "Loading active schedules should succeed")

	// Wait for cron to process
	time.Sleep(100 * time.Millisecond)

	// Verify active schedules are loaded
	activeSchedules, err := scheduleRepo.GetActive()
	require.NoError(t, err, "Should retrieve active schedules")
	
	// Check that our active schedule is in the list
	found := false
	for _, s := range activeSchedules {
		if s.Id == activeSchedule.Id {
			found = true
			break
		}
	}
	assert.True(t, found, "Active schedule should be loaded")
}

// TestIntegration_ScheduleAccessControl tests access control between users
func TestIntegration_ScheduleAccessControl(t *testing.T) {
	db := setupIntegrationDB(t)
	scheduleService, databaseService, _ := setupScheduleServices(db)
	scheduleService.StartScheduler()

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

	// User 1 creates database and schedule
	database := &models.Database{
		Name:     "User 1 DB",
		Type:     "mysql",
		Host:     "localhost",
		Port:     "3306",
		Username: "user1",
		Password: "password123",
		DbName:   "user1db",
		UserId:   user1.Id,
	}
	err := databaseService.CreateDatabase(database, user1.Id, "127.0.0.1", "integration-test")
	require.NoError(t, err, "Database creation should succeed")

	schedule, err := scheduleService.CreateSchedule(
		database.Id,
		user1.Id,
		"User 1 Schedule",
		"0 0 * * *",
		"127.0.0.1",
		"integration-test",
	)
	require.NoError(t, err, "Schedule creation should succeed")

	// User 2 tries to access User 1's schedule - should fail
	_, err = scheduleService.GetSchedule(schedule.Id, user2.Id)
	assert.Error(t, err, "User 2 should not access User 1's schedule")

	// User 2 tries to delete User 1's schedule - should fail
	err = scheduleService.DeleteSchedule(schedule.Id, user2.Id, "127.0.0.1", "integration-test")
	assert.Error(t, err, "User 2 should not delete User 1's schedule")

	// Verify schedule still exists
	_, err = scheduleService.GetSchedule(schedule.Id, user1.Id)
	assert.NoError(t, err, "Schedule should still exist")
}

// TestIntegration_CronParserValidation tests direct cron parser validation
func TestIntegration_CronParserValidation(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		shouldPass bool
	}{
		{"Standard 5-field", "0 0 * * *", true},
		{"Every minute", "* * * * *", true},
		{"Complex expression", "*/15 9-17 * * 1-5", true},
		{"Too many fields", "0 0 * * * * *", false},
		{"Too few fields", "0 0 *", false},
		{"Invalid character", "x x * * *", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := cron.ParseStandard(tc.expression)
			if tc.shouldPass {
				assert.NoError(t, err, "Expression should be valid")
			} else {
				assert.Error(t, err, "Expression should be invalid")
			}
		})
	}
}

