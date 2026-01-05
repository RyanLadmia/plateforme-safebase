package functionals

import (
	"testing"
	"time"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ============================================================================
// FUNCTIONAL TEST FRAMEWORK - Professional Setup
// ============================================================================

// TestSuite represents a complete functional test environment
type TestSuite struct {
	DB                   *gorm.DB
	AuthService          *services.AuthService
	DatabaseService      *services.DatabaseService
	ScheduleService      *services.ScheduleService
	BackupService        *services.BackupService
	ActionHistoryService *services.ActionHistoryService
	UserService          *services.UserService
	
	// Repositories
	UserRepo         *repositories.UserRepository
	SessionRepo      *repositories.SessionRepository
	DatabaseRepo     *repositories.DatabaseRepository
	BackupRepo       *repositories.BackupRepository
	ScheduleRepo     *repositories.ScheduleRepository
	RestoreRepo      *repositories.RestoreRepository
	ActionHistoryRepo *repositories.ActionHistoryRepository
	
	// Test data to cleanup
	createdUsers     []uint
	createdDatabases []uint
	createdBackups   []uint
	createdSchedules []uint
	createdSessions  []string
}

// SetupTestSuite initializes a complete test environment with all services
func SetupTestSuite(t *testing.T) *TestSuite {
	// Create in-memory SQLite database for functional tests
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		// Disable foreign key constraints during cleanup
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	require.NoError(t, err, "Failed to connect to test database")

	// Auto migrate all models
	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Session{},
		&models.Database{},
		&models.Backup{},
		&models.Schedule{},
		&models.Restore{},
		&models.ActionHistory{},
		&models.Alert{},
	)
	require.NoError(t, err, "Failed to migrate test database")

	// Create default roles
	roles := []models.Role{
		{Id: 1, Name: "admin"},
		{Id: 2, Name: "user"},
	}
	for _, role := range roles {
		db.Create(&role)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	databaseRepo := repositories.NewDatabaseRepository(db)
	backupRepo := repositories.NewBackupRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	restoreRepo := repositories.NewRestoreRepository(db)
	actionHistoryRepo := repositories.NewActionHistoryRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, sessionRepo, "test-secret-key-functional", 24*time.Hour)
	userService := services.NewUserService(userRepo, nil, actionHistoryRepo)
	
	// Create database service
	databaseService := services.NewDatabaseService(databaseRepo, backupRepo, restoreRepo, scheduleRepo, nil)
	
	// Create backup service with temp directory
	backupService := services.NewBackupService(backupRepo, databaseService, userService, "/tmp/test-backups")
	
	// Create schedule service
	scheduleService := services.NewScheduleService(scheduleRepo, databaseRepo, backupService)
	
	// Create action history service
	actionHistoryService := services.NewActionHistoryService(actionHistoryRepo)
	
	// Connect action history to other services
	databaseService.SetActionHistoryService(actionHistoryService)
	scheduleService.SetActionHistoryService(actionHistoryService)
	backupService.SetActionHistoryService(actionHistoryService)

	return &TestSuite{
		DB:                   db,
		AuthService:          authService,
		DatabaseService:      databaseService,
		ScheduleService:      scheduleService,
		BackupService:        backupService,
		ActionHistoryService: actionHistoryService,
		UserService:          userService,
		UserRepo:             userRepo,
		SessionRepo:          sessionRepo,
		DatabaseRepo:         databaseRepo,
		BackupRepo:           backupRepo,
		ScheduleRepo:         scheduleRepo,
		RestoreRepo:          restoreRepo,
		ActionHistoryRepo:    actionHistoryRepo,
		createdUsers:         []uint{},
		createdDatabases:     []uint{},
		createdBackups:       []uint{},
		createdSchedules:     []uint{},
		createdSessions:      []string{},
	}
}

// Cleanup removes all test data created during the test
func (ts *TestSuite) Cleanup(t *testing.T) {
	t.Log("Starting cleanup of functional test data...")

	// Delete sessions
	for _, token := range ts.createdSessions {
		ts.SessionRepo.DeleteByToken(token)
	}
	t.Logf("Cleaned %d sessions", len(ts.createdSessions))

	// Delete schedules (must be done before databases due to foreign keys)
	for _, scheduleID := range ts.createdSchedules {
		ts.DB.Unscoped().Delete(&models.Schedule{}, scheduleID)
	}
	t.Logf("Cleaned %d schedules", len(ts.createdSchedules))

	// Delete backups
	for _, backupID := range ts.createdBackups {
		ts.DB.Unscoped().Delete(&models.Backup{}, backupID)
	}
	t.Logf("Cleaned %d backups", len(ts.createdBackups))

	// Delete databases
	for _, databaseID := range ts.createdDatabases {
		ts.DB.Unscoped().Delete(&models.Database{}, databaseID)
	}
	t.Logf("Cleaned %d databases", len(ts.createdDatabases))

	// Delete action history
	ts.DB.Unscoped().Where("user_id IN ?", ts.createdUsers).Delete(&models.ActionHistory{})

	// Delete users
	for _, userID := range ts.createdUsers {
		ts.DB.Unscoped().Delete(&models.User{}, userID)
	}
	t.Logf("Cleaned %d users", len(ts.createdUsers))

	t.Log("Cleanup completed successfully")
}

// RegisterUser registers a test user and tracks it for cleanup
func (ts *TestSuite) RegisterUser(t *testing.T, email, password, firstname, lastname string) (*models.User, error) {
	user := &models.User{
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Password:  password,
	}
	
	err := ts.AuthService.Register(user)
	if err != nil {
		return nil, err
	}
	
	// Track for cleanup
	ts.createdUsers = append(ts.createdUsers, user.Id)
	t.Logf("Registered test user: %s (ID: %d)", email, user.Id)
	
	return user, nil
}

// LoginUser logs in a user and tracks the session for cleanup
func (ts *TestSuite) LoginUser(t *testing.T, email, password string) (string, error) {
	token, err := ts.AuthService.Login(email, password)
	if err != nil {
		return "", err
	}
	
	// Track for cleanup
	ts.createdSessions = append(ts.createdSessions, token)
	t.Logf("Logged in user: %s", email)
	
	return token, nil
}

// CreateDatabase creates a test database and tracks it for cleanup
func (ts *TestSuite) CreateDatabase(t *testing.T, userID uint, name, dbType string) (*models.Database, error) {
	database := &models.Database{
		Name:     name,
		Type:     dbType,
		Host:     "localhost",
		Port:     "3306",
		Username: "test_user",
		Password: "test_password",
		DbName:   "test_db",
		UserId:   userID,
	}
	
	err := ts.DatabaseService.CreateDatabase(database, userID, "127.0.0.1", "functional-test")
	if err != nil {
		return nil, err
	}
	
	// Track for cleanup
	ts.createdDatabases = append(ts.createdDatabases, database.Id)
	t.Logf("Created test database: %s (ID: %d)", name, database.Id)
	
	return database, nil
}

// CreateSchedule creates a test schedule and tracks it for cleanup
func (ts *TestSuite) CreateSchedule(t *testing.T, databaseID, userID uint, name, cronExpr string) (*models.Schedule, error) {
	schedule, err := ts.ScheduleService.CreateSchedule(databaseID, userID, name, cronExpr, "127.0.0.1", "functional-test")
	if err != nil {
		return nil, err
	}
	
	// Track for cleanup
	ts.createdSchedules = append(ts.createdSchedules, schedule.Id)
	t.Logf("Created test schedule: %s (ID: %d)", name, schedule.Id)
	
	return schedule, nil
}

// ============================================================================
// FUNCTIONAL TESTS - Complete User Workflows
// ============================================================================

// TestFunctional_CompleteUserJourney tests a complete user workflow from registration to operations
func TestFunctional_CompleteUserJourney(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.Cleanup(t)

	t.Log("=== Starting Complete User Journey Test ===")

	// Step 1: User Registration
	t.Run("User Registration", func(t *testing.T) {
		user, err := suite.RegisterUser(t, "alice@company.com", "SecureP@ssw0rd123", "Alice", "Johnson")
		require.NoError(t, err, "User registration should succeed")
		assert.NotZero(t, user.Id, "User should have an ID")
		assert.NotNil(t, user.RoleID, "User should have a role")
		assert.Equal(t, uint(2), *user.RoleID, "User should have 'user' role by default")
	})

	// Step 2: User Login
	var token string
	t.Run("User Login", func(t *testing.T) {
		var err error
		token, err = suite.LoginUser(t, "alice@company.com", "SecureP@ssw0rd123")
		require.NoError(t, err, "Login should succeed")
		assert.NotEmpty(t, token, "Token should be generated")
	})

	// Step 3: Verify User Session
	var userID uint
	t.Run("Verify Active Session", func(t *testing.T) {
		user, err := suite.AuthService.GetUserFromToken(token)
		require.NoError(t, err, "Should retrieve user from token")
		assert.Equal(t, "alice@company.com", user.Email, "User email should match")
		userID = user.Id
	})

	// Step 4: Create Database Configuration
	var databaseID uint
	t.Run("Create Database Configuration", func(t *testing.T) {
		database, err := suite.CreateDatabase(t, userID, "Production MySQL", "mysql")
		require.NoError(t, err, "Database creation should succeed")
		assert.NotZero(t, database.Id, "Database should have an ID")
		databaseID = database.Id
	})

	// Step 5: Create Backup Schedule
	var scheduleID uint
	t.Run("Create Backup Schedule", func(t *testing.T) {
		schedule, err := suite.CreateSchedule(t, databaseID, userID, "Daily Backup", "0 2 * * *")
		require.NoError(t, err, "Schedule creation should succeed")
		assert.NotZero(t, schedule.Id, "Schedule should have an ID")
		assert.True(t, schedule.Active, "Schedule should be active")
		scheduleID = schedule.Id
	})

	// Step 6: Verify Action History
	t.Run("Verify Action History Logged", func(t *testing.T) {
		histories, total, err := suite.ActionHistoryService.GetUserActionHistory(userID, 1, 10)
		require.NoError(t, err, "Should retrieve action history")
		assert.GreaterOrEqual(t, int(total), 2, "Should have at least 2 actions logged (database + schedule)")
		
		// Verify action types
		actionTypes := make(map[string]int)
		for _, h := range histories {
			actionTypes[h.ResourceType]++
		}
		assert.GreaterOrEqual(t, actionTypes["database"], 1, "Should have database actions")
		assert.GreaterOrEqual(t, actionTypes["schedule"], 1, "Should have schedule actions")
	})

	// Step 7: Update Schedule
	t.Run("Update Schedule", func(t *testing.T) {
		inactive := false
		updated, err := suite.ScheduleService.UpdateSchedule(scheduleID, userID, "Updated Daily Backup", "0 3 * * *", &inactive, "127.0.0.1", "functional-test")
		require.NoError(t, err, "Schedule update should succeed")
		assert.Equal(t, "Updated Daily Backup", updated.Name, "Name should be updated")
		assert.False(t, updated.Active, "Schedule should be inactive")
	})

	// Step 8: List User's Databases
	t.Run("List User Databases", func(t *testing.T) {
		databases, err := suite.DatabaseService.GetDatabasesByUser(userID)
		require.NoError(t, err, "Should retrieve user databases")
		assert.GreaterOrEqual(t, len(databases), 1, "User should have at least 1 database")
	})

	// Step 9: Delete Schedule
	t.Run("Delete Schedule", func(t *testing.T) {
		err := suite.ScheduleService.DeleteSchedule(scheduleID, userID, "127.0.0.1", "functional-test")
		require.NoError(t, err, "Schedule deletion should succeed")
	})

	// Step 10: User Logout
	t.Run("User Logout", func(t *testing.T) {
		err := suite.AuthService.Logout(token)
		require.NoError(t, err, "Logout should succeed")
		
		// Verify session is deleted
		_, err = suite.SessionRepo.GetSessionByToken(token)
		assert.Error(t, err, "Session should not exist after logout")
	})

	t.Log("=== Complete User Journey Test Finished Successfully ===")
}

// TestFunctional_MultiUserCollaboration tests multiple users working independently
func TestFunctional_MultiUserCollaboration(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.Cleanup(t)

	t.Log("=== Starting Multi-User Collaboration Test ===")

	// Create multiple users
	users := []struct {
		email     string
		password  string
		firstname string
		lastname  string
	}{
		{"bob@company.com", "BobP@ssw0rd123", "Bob", "Smith"},
		{"carol@company.com", "CarolP@ssw0rd123", "Carol", "Williams"},
		{"david@company.com", "DavidP@ssw0rd123", "David", "Brown"},
	}

	userIDs := make([]uint, len(users))
	tokens := make([]string, len(users))

	// Register and login all users
	t.Run("Register and Login Multiple Users", func(t *testing.T) {
		for i, u := range users {
			user, err := suite.RegisterUser(t, u.email, u.password, u.firstname, u.lastname)
			require.NoError(t, err, "User registration should succeed")
			userIDs[i] = user.Id

			token, err := suite.LoginUser(t, u.email, u.password)
			require.NoError(t, err, "Login should succeed")
			tokens[i] = token
		}
	})

	// Each user creates their own databases
	t.Run("Users Create Independent Databases", func(t *testing.T) {
		for i, userID := range userIDs {
			database, err := suite.CreateDatabase(t, userID, "User "+users[i].firstname+" DB", "mysql")
			require.NoError(t, err, "Database creation should succeed")
			assert.Equal(t, userID, database.UserId, "Database should belong to correct user")
		}
	})

	// Verify data isolation between users
	t.Run("Verify Data Isolation", func(t *testing.T) {
		for i, userID := range userIDs {
			databases, err := suite.DatabaseService.GetDatabasesByUser(userID)
			require.NoError(t, err, "Should retrieve user databases")
			assert.Equal(t, 1, len(databases), "Each user should have exactly 1 database")
			assert.Contains(t, databases[0].Name, users[i].firstname, "Database name should contain user's name")
		}
	})

	// Verify action history isolation
	t.Run("Verify Action History Isolation", func(t *testing.T) {
		for _, userID := range userIDs {
			histories, _, err := suite.ActionHistoryService.GetUserActionHistory(userID, 1, 10)
			require.NoError(t, err, "Should retrieve action history")
			
			// Verify all actions belong to this user
			for _, h := range histories {
				assert.Equal(t, userID, h.UserId, "All actions should belong to the user")
			}
		}
	})

	// Test unauthorized access
	t.Run("Test Access Control", func(t *testing.T) {
		// User 0 tries to access User 1's database
		user1Databases, _ := suite.DatabaseService.GetDatabasesByUser(userIDs[1])
		if len(user1Databases) > 0 {
			err := suite.DatabaseService.DeleteDatabase(user1Databases[0].Id, userIDs[0], "127.0.0.1", "functional-test")
			assert.Error(t, err, "User should not be able to delete another user's database")
			assert.Contains(t, err.Error(), "non autorisé", "Error should indicate unauthorized access")
		}
	})

	t.Log("=== Multi-User Collaboration Test Finished Successfully ===")
}

// TestFunctional_DatabaseLifecycle tests complete database management lifecycle
func TestFunctional_DatabaseLifecycle(t *testing.T) {
	suite := SetupTestSuite(t)
	defer suite.Cleanup(t)

	t.Log("=== Starting Database Lifecycle Test ===")

	// Setup user
	user, err := suite.RegisterUser(t, "admin@company.com", "AdminP@ssw0rd123", "Admin", "User")
	require.NoError(t, err, "User registration should succeed")

	token, err := suite.LoginUser(t, "admin@company.com", "AdminP@ssw0rd123")
	require.NoError(t, err, "Login should succeed")
	_ = token // Keep token for future use

	t.Run("Create Multiple Database Types", func(t *testing.T) {
		// Create MySQL database
		mysqlDB, err := suite.CreateDatabase(t, user.Id, "Production MySQL", "mysql")
		require.NoError(t, err, "MySQL database creation should succeed")
		assert.Equal(t, "mysql", mysqlDB.Type, "Database type should be mysql")

		// Create PostgreSQL database
		pgDB, err := suite.CreateDatabase(t, user.Id, "Production PostgreSQL", "postgresql")
		require.NoError(t, err, "PostgreSQL database creation should succeed")
		assert.Equal(t, "postgresql", pgDB.Type, "Database type should be postgresql")
	})

	t.Run("Update Database Configuration", func(t *testing.T) {
		databases, err := suite.DatabaseService.GetDatabasesByUser(user.Id)
		require.NoError(t, err, "Should retrieve databases")
		require.NotEmpty(t, databases, "User should have databases")

		// Update first database
		db := databases[0]
		originalName := db.Name
		db.Name = "Updated " + originalName
		db.Host = "new-host.example.com"

		err = suite.DatabaseService.UpdateDatabase(&db, user.Id, "127.0.0.1", "functional-test")
		require.NoError(t, err, "Database update should succeed")

		// Verify update
		updated, err := suite.DatabaseService.GetDatabaseByID(db.Id)
		require.NoError(t, err, "Should retrieve updated database")
		assert.Equal(t, "Updated "+originalName, updated.Name, "Name should be updated")
		assert.Equal(t, "new-host.example.com", updated.Host, "Host should be updated")
	})

	t.Run("Password Encryption Verification", func(t *testing.T) {
		databases, err := suite.DatabaseService.GetDatabasesByUser(user.Id)
		require.NoError(t, err, "Should retrieve databases")
		require.NotEmpty(t, databases, "User should have databases")

		// Get database from repository (encrypted)
		dbFromRepo, err := suite.DatabaseRepo.GetByID(databases[0].Id)
		require.NoError(t, err, "Should retrieve from repository")

		// Get database from service (decrypted)
		dbFromService, err := suite.DatabaseService.GetDatabaseByID(databases[0].Id)
		require.NoError(t, err, "Should retrieve from service")

		// Verify encryption
		assert.NotEqual(t, dbFromRepo.Password, dbFromService.Password, "Repository password should be encrypted")
		assert.Equal(t, "test_password", dbFromService.Password, "Service should return decrypted password")
	})

	t.Run("Soft Delete Verification", func(t *testing.T) {
		databases, err := suite.DatabaseService.GetDatabasesByUser(user.Id)
		require.NoError(t, err, "Should retrieve databases")
		require.NotEmpty(t, databases, "User should have databases")

		dbToDelete := databases[0]
		err = suite.DatabaseService.DeleteDatabase(dbToDelete.Id, user.Id, "127.0.0.1", "functional-test")
		require.NoError(t, err, "Database deletion should succeed")

		// Verify soft delete
		_, err = suite.DatabaseService.GetDatabaseByID(dbToDelete.Id)
		assert.Error(t, err, "Deleted database should not be retrievable")

		// Verify record still exists with DeletedAt set
		var deletedDB models.Database
		suite.DB.Unscoped().First(&deletedDB, dbToDelete.Id)
		assert.NotNil(t, deletedDB.DeletedAt, "DeletedAt should be set")
	})

	t.Log("=== Database Lifecycle Test Finished Successfully ===")
}

