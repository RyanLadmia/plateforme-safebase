package units

import (
	"testing"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/RyanLadmia/plateforme-safebase/pkg/security"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// setupDatabaseTestDB initializes an in-memory SQLite database for testing
func setupDatabaseTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.Database{}, &models.User{}, &models.Role{}, &models.Backup{}, &models.Schedule{}, &models.Restore{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Create a test user
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

	return db
}

// createTestDatabase creates a test database configuration
func createTestDatabase(db *gorm.DB, name, dbType string, userID uint) *models.Database {
	database := &models.Database{
		Name:     name,
		Type:     dbType,
		Host:     "localhost",
		Port:     "3306",
		Username: "testuser",
		Password: "testpassword",
		DbName:   "testdb",
		UserId:   userID,
	}
	db.Create(database)
	return database
}

// ============================================================================
// TESTS - Database Service
// ============================================================================

// Test 1: Create Database - Should create a new database configuration
func TestDatabaseService_CreateDatabase(t *testing.T) {
	db := setupDatabaseTestDB(t)

	// Create repositories and service
	databaseRepo := repositories.NewDatabaseRepository(db)
	backupRepo := repositories.NewBackupRepository(db)
	restoreRepo := repositories.NewRestoreRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	databaseService := services.NewDatabaseService(databaseRepo, backupRepo, restoreRepo, scheduleRepo, nil)

	tests := []struct {
		name     string
		database *models.Database
		wantErr  bool
	}{
		{
			name: "Create MySQL database",
			database: &models.Database{
				Name:     "Test MySQL DB",
				Type:     "mysql",
				Host:     "localhost",
				Port:     "3306",
				Username: "root",
				Password: "password123",
				DbName:   "testdb",
				UserId:   1,
			},
			wantErr: false,
		},
		{
			name: "Create PostgreSQL database",
			database: &models.Database{
				Name:     "Test PostgreSQL DB",
				Type:     "postgresql",
				Host:     "localhost",
				Port:     "5432",
				Username: "postgres",
				Password: "password123",
				DbName:   "testdb",
				UserId:   1,
			},
			wantErr: false,
		},
		{
			name: "Invalid database type",
			database: &models.Database{
				Name:     "Invalid DB",
				Type:     "mongodb",
				Host:     "localhost",
				Port:     "27017",
				Username: "mongo",
				Password: "password123",
				DbName:   "testdb",
				UserId:   1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := databaseService.CreateDatabase(tt.database, 1, "127.0.0.1", "test-agent")
			if tt.wantErr {
				assert.Error(t, err, "Expected error for invalid database type")
			} else {
				assert.NoError(t, err, "Database creation should succeed")
				assert.NotZero(t, tt.database.Id, "Database ID should be set")

				// Verify password was encrypted
				assert.NotEqual(t, "password123", tt.database.Password, "Password should be encrypted")
			}
		})
	}
}

// Test 2: Get Database by ID - Should retrieve database configuration
func TestDatabaseService_GetDatabaseByID(t *testing.T) {
	db := setupDatabaseTestDB(t)

	// Create repositories and service
	databaseRepo := repositories.NewDatabaseRepository(db)
	backupRepo := repositories.NewBackupRepository(db)
	restoreRepo := repositories.NewRestoreRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	databaseService := services.NewDatabaseService(databaseRepo, backupRepo, restoreRepo, scheduleRepo, nil)

	// Create test database with encrypted password
	testDB := &models.Database{
		Name:     "Test Database",
		Type:     "mysql",
		Host:     "localhost",
		Port:     "3306",
		Username: "testuser",
		Password: "password123",
		DbName:   "testdb",
		UserId:   1,
	}
	encryptedPassword, _ := security.EncryptDatabasePassword(testDB.Password)
	testDB.Password = encryptedPassword
	db.Create(testDB)

	// Test successful retrieval
	retrieved, err := databaseService.GetDatabaseByID(testDB.Id)
	assert.NoError(t, err, "Should retrieve database successfully")
	assert.Equal(t, testDB.Id, retrieved.Id, "Database ID should match")
	assert.Equal(t, testDB.Name, retrieved.Name, "Database name should match")
	assert.Equal(t, "password123", retrieved.Password, "Password should be decrypted")

	// Test non-existent database
	_, err = databaseService.GetDatabaseByID(9999)
	assert.Error(t, err, "Should return error for non-existent database")
}

// Test 3: Update Database - Should update database configuration
func TestDatabaseService_UpdateDatabase(t *testing.T) {
	db := setupDatabaseTestDB(t)

	// Create repositories and service
	databaseRepo := repositories.NewDatabaseRepository(db)
	backupRepo := repositories.NewBackupRepository(db)
	restoreRepo := repositories.NewRestoreRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	databaseService := services.NewDatabaseService(databaseRepo, backupRepo, restoreRepo, scheduleRepo, nil)

	// Create test database
	testDB := createTestDatabase(db, "Original Name", "mysql", 1)

	// Update database
	testDB.Name = "Updated Name"
	testDB.Host = "newhost.example.com"
	testDB.Port = "3307"

	err := databaseService.UpdateDatabase(testDB, 1, "127.0.0.1", "test-agent")
	assert.NoError(t, err, "Database update should succeed")

	// Verify updates
	updated, err := databaseRepo.GetByID(testDB.Id)
	assert.NoError(t, err, "Should retrieve updated database")
	assert.Equal(t, "Updated Name", updated.Name, "Name should be updated")
	assert.Equal(t, "newhost.example.com", updated.Host, "Host should be updated")
	assert.Equal(t, "3307", updated.Port, "Port should be updated")
}

// Test 4: Delete Database - Should soft delete database configuration
func TestDatabaseService_DeleteDatabase(t *testing.T) {
	db := setupDatabaseTestDB(t)

	// Create repositories and service
	databaseRepo := repositories.NewDatabaseRepository(db)
	backupRepo := repositories.NewBackupRepository(db)
	restoreRepo := repositories.NewRestoreRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	databaseService := services.NewDatabaseService(databaseRepo, backupRepo, restoreRepo, scheduleRepo, nil)

	// Create test database
	testDB := createTestDatabase(db, "Database to Delete", "postgresql", 1)

	// Delete database
	err := databaseService.DeleteDatabase(testDB.Id, 1, "127.0.0.1", "test-agent")
	assert.NoError(t, err, "Database deletion should succeed")

	// Verify database is soft deleted
	_, err = databaseRepo.GetByID(testDB.Id)
	assert.Error(t, err, "Deleted database should not be retrievable")

	// Verify soft delete (check DeletedAt is set)
	var deletedDB models.Database
	db.Unscoped().First(&deletedDB, testDB.Id)
	assert.NotNil(t, deletedDB.DeletedAt, "DeletedAt should be set for soft delete")
}

// Test 5: List Databases by User - Should return all databases for a user
func TestDatabaseService_GetDatabasesByUser(t *testing.T) {
	db := setupDatabaseTestDB(t)

	// Create repositories and service
	databaseRepo := repositories.NewDatabaseRepository(db)
	backupRepo := repositories.NewBackupRepository(db)
	restoreRepo := repositories.NewRestoreRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	databaseService := services.NewDatabaseService(databaseRepo, backupRepo, restoreRepo, scheduleRepo, nil)

	// Create multiple test databases
	createTestDatabase(db, "Database 1", "mysql", 1)
	createTestDatabase(db, "Database 2", "postgresql", 1)
	createTestDatabase(db, "Database 3", "mysql", 1)

	// Get all databases for user
	databases, err := databaseService.GetDatabasesByUser(1)
	assert.NoError(t, err, "Should retrieve databases successfully")
	assert.Equal(t, 3, len(databases), "Should return 3 databases")

	// Verify each database belongs to the user
	for _, database := range databases {
		assert.Equal(t, uint(1), database.UserId, "All databases should belong to user 1")
	}

	// Test with non-existent user
	databases, err = databaseService.GetDatabasesByUser(999)
	assert.NoError(t, err, "Should not error for non-existent user")
	assert.Equal(t, 0, len(databases), "Should return empty list for non-existent user")
}

