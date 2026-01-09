package units

import (
	"testing"
	"time"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// setupScheduleTestDB initializes an in-memory SQLite database for testing
func setupScheduleTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.Schedule{}, &models.Database{}, &models.User{}, &models.Role{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Create test user
	roleID := uint(2)
	user := &models.User{
		Id:        1,
		Firstname: "Test",
		Lastname:  "User",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Active:    true,
		RoleID:    &roleID,
	}
	db.Create(user)

	// Create test database
	database := &models.Database{
		Id:       1,
		Name:     "Test Database",
		Type:     "mysql",
		Host:     "localhost",
		Port:     "3306",
		Username: "testuser",
		Password: "testpassword",
		DbName:   "testdb",
		UserId:   1,
	}
	db.Create(database)

	return db
}

// createTestSchedule creates a test schedule
func createTestSchedule(db *gorm.DB, name, cronExpression string, databaseID, userID uint, active bool) *models.Schedule {
	schedule := &models.Schedule{
		Name:           name,
		CronExpression: cronExpression,
		Active:         active,
		UserId:         userID,
		DatabaseId:     databaseID,
	}
	db.Create(schedule)
	return schedule
}

// ============================================================================
// TESTS - Schedule Service
// ============================================================================

// Test 1: Create Schedule - Should create a new schedule with valid cron expression
func TestScheduleService_CreateSchedule(t *testing.T) {
	db := setupScheduleTestDB(t)

	// Create repositories and service
	scheduleRepo := repositories.NewScheduleRepository(db)
	databaseRepo := repositories.NewDatabaseRepository(db)
	scheduleService := services.NewScheduleService(scheduleRepo, databaseRepo, nil)
	scheduleService.StartScheduler()

	tests := []struct {
		name           string
		scheduleName   string
		cronExpression string
		wantErr        bool
	}{
		{
			name:           "Daily at midnight",
			scheduleName:   "Daily Backup",
			cronExpression: "0 0 * * *",
			wantErr:        false,
		},
		{
			name:           "Every 6 hours",
			scheduleName:   "6 Hour Backup",
			cronExpression: "0 */6 * * *",
			wantErr:        false,
		},
		{
			name:           "Weekly on Sunday",
			scheduleName:   "Weekly Backup",
			cronExpression: "0 0 * * 0",
			wantErr:        false,
		},
		{
			name:           "Invalid cron expression",
			scheduleName:   "Invalid Schedule",
			cronExpression: "invalid cron",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule, err := scheduleService.CreateSchedule(1, 1, tt.scheduleName, tt.cronExpression, "127.0.0.1", "test-agent")
			if tt.wantErr {
				assert.Error(t, err, "Expected error for invalid cron expression")
			} else {
				assert.NoError(t, err, "Schedule creation should succeed")
				assert.NotNil(t, schedule, "Schedule should not be nil")
				assert.Equal(t, tt.scheduleName, schedule.Name, "Schedule name should match")
				assert.Equal(t, tt.cronExpression, schedule.CronExpression, "Cron expression should match")
				assert.True(t, schedule.Active, "Schedule should be active by default")
			}
		})
	}
}

// Test 2: Validate Cron Expression - Should validate cron syntax
func TestScheduleService_ValidateCronExpression(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantErr    bool
	}{
		{
			name:       "Valid: Every minute",
			expression: "* * * * *",
			wantErr:    false,
		},
		{
			name:       "Valid: Daily at 2am",
			expression: "0 2 * * *",
			wantErr:    false,
		},
		{
			name:       "Valid: Monthly on 1st",
			expression: "0 0 1 * *",
			wantErr:    false,
		},
		{
			name:       "Invalid: Too many fields",
			expression: "0 0 * * * * *",
			wantErr:    true,
		},
		{
			name:       "Invalid: Wrong syntax",
			expression: "every day",
			wantErr:    true,
		},
		{
			name:       "Invalid: Empty",
			expression: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cron.ParseStandard(tt.expression)
			if tt.wantErr {
				assert.Error(t, err, "Expected error for invalid cron expression")
			} else {
				assert.NoError(t, err, "Expected no error for valid cron expression")
			}
		})
	}
}

// Test 3: Update Schedule - Should update schedule properties
func TestScheduleService_UpdateSchedule(t *testing.T) {
	db := setupScheduleTestDB(t)

	// Create repositories and service
	scheduleRepo := repositories.NewScheduleRepository(db)
	databaseRepo := repositories.NewDatabaseRepository(db)
	scheduleService := services.NewScheduleService(scheduleRepo, databaseRepo, nil)
	scheduleService.StartScheduler()

	// Create test schedule
	testSchedule := createTestSchedule(db, "Original Schedule", "0 0 * * *", 1, 1, true)

	// Test update name and cron expression
	newActive := false
	updated, err := scheduleService.UpdateSchedule(testSchedule.Id, 1, "Updated Schedule", "0 12 * * *", &newActive, "127.0.0.1", "test-agent")
	assert.NoError(t, err, "Schedule update should succeed")
	assert.Equal(t, "Updated Schedule", updated.Name, "Name should be updated")
	assert.Equal(t, "0 12 * * *", updated.CronExpression, "Cron expression should be updated")
	assert.False(t, updated.Active, "Active status should be updated")

	// Test update with invalid cron expression
	_, err = scheduleService.UpdateSchedule(testSchedule.Id, 1, "Invalid Update", "invalid cron", nil, "127.0.0.1", "test-agent")
	assert.Error(t, err, "Should fail with invalid cron expression")
}

// Test 4: Delete Schedule - Should delete schedule and remove from cron
func TestScheduleService_DeleteSchedule(t *testing.T) {
	db := setupScheduleTestDB(t)

	// Create repositories and service
	scheduleRepo := repositories.NewScheduleRepository(db)
	databaseRepo := repositories.NewDatabaseRepository(db)
	scheduleService := services.NewScheduleService(scheduleRepo, databaseRepo, nil)
	scheduleService.StartScheduler()

	// Create test schedule
	testSchedule := createTestSchedule(db, "Schedule to Delete", "0 0 * * *", 1, 1, true)

	// Delete schedule
	err := scheduleService.DeleteSchedule(testSchedule.Id, 1, "127.0.0.1", "test-agent")
	assert.NoError(t, err, "Schedule deletion should succeed")

	// Verify schedule is deleted
	_, err = scheduleRepo.GetByID(testSchedule.Id)
	assert.Error(t, err, "Deleted schedule should not be retrievable")

	// Test unauthorized deletion
	testSchedule2 := createTestSchedule(db, "Another Schedule", "0 0 * * *", 1, 1, true)
	err = scheduleService.DeleteSchedule(testSchedule2.Id, 999, "127.0.0.1", "test-agent")
	assert.Error(t, err, "Should fail with unauthorized user")
}

// Test 5: List Schedules by User - Should return all schedules for a user
func TestScheduleService_GetSchedules(t *testing.T) {
	db := setupScheduleTestDB(t)

	// Create repositories and service
	scheduleRepo := repositories.NewScheduleRepository(db)
	databaseRepo := repositories.NewDatabaseRepository(db)
	scheduleService := services.NewScheduleService(scheduleRepo, databaseRepo, nil)

	// Create multiple test schedules
	createTestSchedule(db, "Daily Backup", "0 0 * * *", 1, 1, true)
	createTestSchedule(db, "Weekly Backup", "0 0 * * 0", 1, 1, true)
	createTestSchedule(db, "Monthly Backup", "0 0 1 * *", 1, 1, false)

	// Get all schedules for user
	schedules, err := scheduleService.GetSchedules(1)
	assert.NoError(t, err, "Should retrieve schedules successfully")
	assert.Equal(t, 3, len(schedules), "Should return 3 schedules")

	// Verify each schedule belongs to the user
	for _, schedule := range schedules {
		assert.Equal(t, uint(1), schedule.UserId, "All schedules should belong to user 1")
	}

	// Test with non-existent user
	schedules, err = scheduleService.GetSchedules(999)
	assert.NoError(t, err, "Should not error for non-existent user")
	assert.Equal(t, 0, len(schedules), "Should return empty list for non-existent user")
}

// Test 6: Activate/Deactivate Schedule - Should toggle schedule active status
func TestScheduleService_ToggleScheduleStatus(t *testing.T) {
	db := setupScheduleTestDB(t)

	// Create repositories and service
	scheduleRepo := repositories.NewScheduleRepository(db)
	databaseRepo := repositories.NewDatabaseRepository(db)
	scheduleService := services.NewScheduleService(scheduleRepo, databaseRepo, nil)
	scheduleService.StartScheduler()

	// Create test schedule (active)
	testSchedule := createTestSchedule(db, "Toggle Schedule", "0 0 * * *", 1, 1, true)
	assert.True(t, testSchedule.Active, "Schedule should start as active")

	// Deactivate schedule
	newActive := false
	updated, err := scheduleService.UpdateSchedule(testSchedule.Id, 1, "", "", &newActive, "127.0.0.1", "test-agent")
	assert.NoError(t, err, "Should deactivate schedule successfully")
	assert.False(t, updated.Active, "Schedule should be deactivated")

	// Reactivate schedule
	newActive = true
	updated, err = scheduleService.UpdateSchedule(testSchedule.Id, 1, "", "", &newActive, "127.0.0.1", "test-agent")
	assert.NoError(t, err, "Should reactivate schedule successfully")
	assert.True(t, updated.Active, "Schedule should be reactivated")
}

// Test 7: Load Active Schedules - Should load all active schedules into cron
func TestScheduleService_LoadActiveSchedules(t *testing.T) {
	db := setupScheduleTestDB(t)

	// Create repositories and service
	scheduleRepo := repositories.NewScheduleRepository(db)
	databaseRepo := repositories.NewDatabaseRepository(db)
	scheduleService := services.NewScheduleService(scheduleRepo, databaseRepo, nil)
	scheduleService.StartScheduler()

	// Create multiple schedules with different states
	sched1 := createTestSchedule(db, "Active Schedule 1", "0 0 * * *", 1, 1, true)
	sched2 := createTestSchedule(db, "Active Schedule 2", "0 12 * * *", 1, 1, true)
	createTestSchedule(db, "Inactive Schedule", "0 6 * * *", 1, 1, false)

	// Load active schedules
	err := scheduleService.LoadActiveSchedules()
	assert.NoError(t, err, "Loading active schedules should succeed")

	// Wait a moment for cron to process
	time.Sleep(100 * time.Millisecond)

	// Verify that active schedules were loaded
	// We can verify by checking the schedules by ID
	activeSchedules, _ := scheduleRepo.GetActive()
	assert.NotEmpty(t, activeSchedules, "Should have active schedules")
	
	// Verify our specific active schedules exist
	foundSched1 := false
	foundSched2 := false
	for _, s := range activeSchedules {
		if s.Id == sched1.Id {
			foundSched1 = true
		}
		if s.Id == sched2.Id {
			foundSched2 = true
		}
	}
	assert.True(t, foundSched1, "Active Schedule 1 should be in active schedules")
	assert.True(t, foundSched2, "Active Schedule 2 should be in active schedules")
}

