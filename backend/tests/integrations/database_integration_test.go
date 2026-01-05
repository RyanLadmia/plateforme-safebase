package integrations

import (
	"testing"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/RyanLadmia/plateforme-safebase/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ============================================================================
// INTEGRATION TESTS - Database Management Flow
// ============================================================================

// setupDatabaseServices creates all services needed for database management tests
func setupDatabaseServices(db *gorm.DB) (*services.DatabaseService, *services.ActionHistoryService, *repositories.DatabaseRepository) {
	databaseRepo := repositories.NewDatabaseRepository(db)
	backupRepo := repositories.NewBackupRepository(db)
	restoreRepo := repositories.NewRestoreRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	historyRepo := repositories.NewActionHistoryRepository(db)
	
	databaseService := services.NewDatabaseService(databaseRepo, backupRepo, restoreRepo, scheduleRepo, nil)
	historyService := services.NewActionHistoryService(historyRepo)
	databaseService.SetActionHistoryService(historyService)
	
	return databaseService, historyService, databaseRepo
}

// createTestUser creates a test user for integration tests
func createTestUser(db *gorm.DB) *models.User {
	roleID := uint(2)
	user := &models.User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     "testuser@example.com",
		Password:  "hashedpassword",
		Active:    true,
		RoleID:    &roleID,
	}
	db.Create(user)
	return user
}

// TestIntegration_DatabaseCRUDFlow tests the complete CRUD flow for databases
func TestIntegration_DatabaseCRUDFlow(t *testing.T) {
	db := setupIntegrationDB(t)
	databaseService, historyService, databaseRepo := setupDatabaseServices(db)
	user := createTestUser(db)

	var databaseID uint

	// Step 1: Create a database
	t.Run("Create Database", func(t *testing.T) {
		database := &models.Database{
			Name:     "Production MySQL DB",
			Type:     "mysql",
			Host:     "prod-mysql.example.com",
			Port:     "3306",
			Username: "prod_user",
			Password: "SecurePassword123!",
			DbName:   "production_db",
			UserId:   user.Id,
		}

		err := databaseService.CreateDatabase(database, user.Id, "127.0.0.1", "integration-test")
		require.NoError(t, err, "Database creation should succeed")
		assert.NotZero(t, database.Id, "Database ID should be set")
		databaseID = database.Id

		// Verify password was encrypted
		dbFromRepo, err := databaseRepo.GetByID(database.Id)
		require.NoError(t, err, "Should retrieve database from repository")
		assert.NotEqual(t, "SecurePassword123!", dbFromRepo.Password, "Password should be encrypted")

		// Verify action was logged
		histories, _, err := historyService.GetUserActionHistory(user.Id, 1, 10)
		require.NoError(t, err, "Should retrieve action history")
		assert.NotEmpty(t, histories, "Action history should be recorded")
		assert.Equal(t, "created", histories[0].Action, "Action should be 'created'")
		assert.Equal(t, "database", histories[0].ResourceType, "Resource type should be 'database'")
	})

	// Step 2: Retrieve database by ID (with password decryption)
	t.Run("Get Database By ID", func(t *testing.T) {
		database, err := databaseService.GetDatabaseByID(databaseID)
		require.NoError(t, err, "Should retrieve database")
		assert.Equal(t, "Production MySQL DB", database.Name, "Database name should match")
		assert.Equal(t, "SecurePassword123!", database.Password, "Password should be decrypted")
		assert.Equal(t, user.Id, database.UserId, "Database should belong to user")
	})

	// Step 3: Update database
	t.Run("Update Database", func(t *testing.T) {
		database, err := databaseService.GetDatabaseByID(databaseID)
		require.NoError(t, err, "Should retrieve database")

		database.Name = "Updated Production MySQL DB"
		database.Host = "new-prod-mysql.example.com"
		database.Port = "3307"

		err = databaseService.UpdateDatabase(database, user.Id, "127.0.0.1", "integration-test")
		require.NoError(t, err, "Database update should succeed")

		// Verify updates
		updatedDB, err := databaseService.GetDatabaseByID(databaseID)
		require.NoError(t, err, "Should retrieve updated database")
		assert.Equal(t, "Updated Production MySQL DB", updatedDB.Name, "Name should be updated")
		assert.Equal(t, "new-prod-mysql.example.com", updatedDB.Host, "Host should be updated")
		assert.Equal(t, "3307", updatedDB.Port, "Port should be updated")

		// Verify action was logged
		histories, _, err := historyService.GetUserActionHistory(user.Id, 1, 10)
		require.NoError(t, err, "Should retrieve action history")
		assert.GreaterOrEqual(t, len(histories), 2, "Should have at least 2 actions logged")
	})

	// Step 4: Update database name specifically
	t.Run("Update Database Name", func(t *testing.T) {
		err := databaseService.UpdateDatabaseName(databaseID, "Final Database Name", user.Id, "127.0.0.1", "integration-test")
		require.NoError(t, err, "Database name update should succeed")

		// Verify name change
		database, err := databaseService.GetDatabaseByID(databaseID)
		require.NoError(t, err, "Should retrieve database")
		assert.Equal(t, "Final Database Name", database.Name, "Name should be updated")
	})

	// Step 5: List user's databases
	t.Run("List User Databases", func(t *testing.T) {
		databases, err := databaseService.GetDatabasesByUser(user.Id)
		require.NoError(t, err, "Should retrieve user databases")
		assert.Equal(t, 1, len(databases), "User should have 1 database")
		assert.Equal(t, "Final Database Name", databases[0].Name, "Database name should match")
	})

	// Step 6: Delete database (soft delete)
	t.Run("Delete Database", func(t *testing.T) {
		err := databaseService.DeleteDatabase(databaseID, user.Id, "127.0.0.1", "integration-test")
		require.NoError(t, err, "Database deletion should succeed")

		// Verify database is soft deleted
		_, err = databaseService.GetDatabaseByID(databaseID)
		assert.Error(t, err, "Deleted database should not be retrievable")

		// Verify action was logged
		histories, _, err := historyService.GetUserActionHistory(user.Id, 1, 10)
		require.NoError(t, err, "Should retrieve action history")
		
		// Find delete action
		var deleteFound bool
		for _, h := range histories {
			if h.Action == "deleted" && h.ResourceType == "database" {
				deleteFound = true
				break
			}
		}
		assert.True(t, deleteFound, "Delete action should be logged")
	})
}

// TestIntegration_DatabasePasswordEncryption tests password encryption/decryption flow
func TestIntegration_DatabasePasswordEncryption(t *testing.T) {
	db := setupIntegrationDB(t)
	databaseService, _, databaseRepo := setupDatabaseServices(db)
	user := createTestUser(db)

	plainPassword := "SuperSecretP@ssw0rd123"

	// Create database with plain password
	database := &models.Database{
		Name:     "Test Encryption DB",
		Type:     "postgresql",
		Host:     "localhost",
		Port:     "5432",
		Username: "testuser",
		Password: plainPassword,
		DbName:   "testdb",
		UserId:   user.Id,
	}

	err := databaseService.CreateDatabase(database, user.Id, "127.0.0.1", "integration-test")
	require.NoError(t, err, "Database creation should succeed")

	// Verify password is encrypted in database
	dbFromRepo, err := databaseRepo.GetByID(database.Id)
	require.NoError(t, err, "Should retrieve database")
	assert.NotEqual(t, plainPassword, dbFromRepo.Password, "Password should be encrypted in storage")

	// Verify password can be decrypted
	decryptedPassword, err := security.DecryptDatabasePassword(dbFromRepo.Password)
	require.NoError(t, err, "Should decrypt password")
	assert.Equal(t, plainPassword, decryptedPassword, "Decrypted password should match original")

	// Verify GetDatabaseByID returns decrypted password
	retrievedDB, err := databaseService.GetDatabaseByID(database.Id)
	require.NoError(t, err, "Should retrieve database")
	assert.Equal(t, plainPassword, retrievedDB.Password, "Service should return decrypted password")
}

// TestIntegration_MultipleDatabasesPerUser tests multiple databases for a single user
func TestIntegration_MultipleDatabasesPerUser(t *testing.T) {
	db := setupIntegrationDB(t)
	databaseService, _, _ := setupDatabaseServices(db)
	user := createTestUser(db)

	databases := []struct {
		name   string
		dbType string
		host   string
		port   string
	}{
		{"MySQL Dev", "mysql", "localhost", "3306"},
		{"MySQL Staging", "mysql", "staging.example.com", "3306"},
		{"PostgreSQL Production", "postgresql", "prod.example.com", "5432"},
		{"PostgreSQL Backup", "postgresql", "backup.example.com", "5432"},
	}

	// Create multiple databases
	for _, dbConfig := range databases {
		database := &models.Database{
			Name:     dbConfig.name,
			Type:     dbConfig.dbType,
			Host:     dbConfig.host,
			Port:     dbConfig.port,
			Username: "user",
			Password: "password123",
			DbName:   "testdb",
			UserId:   user.Id,
		}
		err := databaseService.CreateDatabase(database, user.Id, "127.0.0.1", "integration-test")
		require.NoError(t, err, "Database creation should succeed")
	}

	// Retrieve all user databases
	userDatabases, err := databaseService.GetDatabasesByUser(user.Id)
	require.NoError(t, err, "Should retrieve user databases")
	assert.Equal(t, 4, len(userDatabases), "User should have 4 databases")

	// Verify all databases belong to user
	for _, db := range userDatabases {
		assert.Equal(t, user.Id, db.UserId, "All databases should belong to user")
	}

	// Verify database types
	mysqlCount := 0
	postgresCount := 0
	for _, db := range userDatabases {
		if db.Type == "mysql" {
			mysqlCount++
		} else if db.Type == "postgresql" {
			postgresCount++
		}
	}
	assert.Equal(t, 2, mysqlCount, "Should have 2 MySQL databases")
	assert.Equal(t, 2, postgresCount, "Should have 2 PostgreSQL databases")
}

// TestIntegration_DatabaseAccessControl tests access control between users
func TestIntegration_DatabaseAccessControl(t *testing.T) {
	db := setupIntegrationDB(t)
	databaseService, _, _ := setupDatabaseServices(db)

	// Create two users
	user1 := createTestUser(db)
	
	roleID := uint(2)
	user2 := &models.User{
		Firstname: "Another",
		Lastname:  "User",
		Email:     "anotheruser@example.com",
		Password:  "hashedpassword",
		Active:    true,
		RoleID:    &roleID,
	}
	db.Create(user2)

	// User 1 creates a database
	database := &models.Database{
		Name:     "User 1 Database",
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

	// User 2 tries to delete User 1's database - should fail
	err = databaseService.DeleteDatabase(database.Id, user2.Id, "127.0.0.1", "integration-test")
	assert.Error(t, err, "User 2 should not be able to delete User 1's database")
	assert.Contains(t, err.Error(), "non autorisé", "Error should indicate unauthorized access")

	// Verify database still exists
	_, err = databaseService.GetDatabaseByID(database.Id)
	assert.NoError(t, err, "Database should still exist")

	// Verify User 2 has no databases
	user2Databases, err := databaseService.GetDatabasesByUser(user2.Id)
	require.NoError(t, err, "Should retrieve user2 databases")
	assert.Equal(t, 0, len(user2Databases), "User 2 should have no databases")
}

