package units

import (
	"os"
	"testing"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/RyanLadmia/plateforme-safebase/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ============================================================================
// HELPER FUNCTIONS & MOCKS
// ============================================================================

// setupBackupTestDB initializes an in-memory SQLite database for testing
func setupBackupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.Backup{}, &models.Database{}, &models.User{}, &models.Role{})
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

	// Create test database with encrypted password
	encryptedPassword, _ := security.EncryptDatabasePassword("testpassword")
	database := &models.Database{
		Id:       1,
		Name:     "Test Database",
		Type:     "mysql",
		Host:     "localhost",
		Port:     "3306",
		Username: "testuser",
		Password: encryptedPassword,
		DbName:   "testdb",
		UserId:   1,
	}
	db.Create(database)

	return db
}

// MockCloudStorage is a mock implementation of CloudStorageService
type MockCloudStorage struct {
	uploadedFiles   map[string][]byte
	shouldFailOnOp  string // "upload", "download", "delete"
}

func NewMockCloudStorage() *MockCloudStorage {
	return &MockCloudStorage{
		uploadedFiles: make(map[string][]byte),
	}
}

func (m *MockCloudStorage) UploadFile(localPath, remotePath string) error {
	if m.shouldFailOnOp == "upload" {
		return assert.AnError
	}
	data, _ := os.ReadFile(localPath)
	m.uploadedFiles[remotePath] = data
	return nil
}

func (m *MockCloudStorage) DownloadFile(remotePath string) ([]byte, error) {
	if m.shouldFailOnOp == "download" {
		return nil, assert.AnError
	}
	if data, exists := m.uploadedFiles[remotePath]; exists {
		return data, nil
	}
	return nil, assert.AnError
}

func (m *MockCloudStorage) DeleteFile(remotePath string) error {
	if m.shouldFailOnOp == "delete" {
		return assert.AnError
	}
	delete(m.uploadedFiles, remotePath)
	return nil
}

func (m *MockCloudStorage) GenerateRemotePath(userName, dbType, filename string) string {
	return userName + "/" + dbType + "/" + filename
}

func (m *MockCloudStorage) FileExists(remotePath string) (bool, error) {
	if m.shouldFailOnOp == "exists" {
		return false, assert.AnError
	}
	_, exists := m.uploadedFiles[remotePath]
	return exists, nil
}

// createTestBackup creates a test backup record
func createTestBackup(db *gorm.DB, databaseID, userID uint, status string) *models.Backup {
	backup := &models.Backup{
		Filename:   "test_backup_2024-01-01_12-00-00.zip",
		Filepath:   "/test/path/test_backup_2024-01-01_12-00-00.zip",
		Size:       1024,
		Status:     status,
		UserId:     userID,
		DatabaseId: databaseID,
	}
	db.Create(backup)
	return backup
}

// ============================================================================
// TESTS - Backup Service
// ============================================================================

// Test 1: Create Backup - Should create backup record with pending status
func TestBackupService_CreateBackup(t *testing.T) {
	db := setupBackupTestDB(t)

	// Create repositories and services
	backupRepo := repositories.NewBackupRepository(db)
	databaseRepo := repositories.NewDatabaseRepository(db)
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo, nil, nil)
	databaseService := services.NewDatabaseService(databaseRepo, backupRepo, nil, nil, nil)
	backupService := services.NewBackupService(backupRepo, databaseService, userService, "/tmp/backups")

	// Test backup creation
	backup, err := backupService.CreateBackup(1, 1, "127.0.0.1", "test-agent")
	assert.NoError(t, err, "Backup creation should succeed")
	assert.NotNil(t, backup, "Backup should not be nil")
	assert.Equal(t, "pending", backup.Status, "Backup status should be pending")
	assert.Equal(t, uint(1), backup.UserId, "Backup should belong to user 1")
	assert.Equal(t, uint(1), backup.DatabaseId, "Backup should be for database 1")
	assert.NotEmpty(t, backup.Filename, "Backup filename should be generated")

	// Test with non-existent database
	_, err = backupService.CreateBackup(999, 1, "127.0.0.1", "test-agent")
	assert.Error(t, err, "Should fail with non-existent database")
}

// Test 2: Get Backup by ID - Should retrieve backup record
func TestBackupService_GetBackupByID(t *testing.T) {
	db := setupBackupTestDB(t)

	// Create repositories and service
	backupRepo := repositories.NewBackupRepository(db)
	backupService := services.NewBackupService(backupRepo, nil, nil, "/tmp/backups")

	// Create test backup
	testBackup := createTestBackup(db, 1, 1, "completed")

	// Test successful retrieval
	backup, err := backupService.GetBackupByID(testBackup.Id)
	assert.NoError(t, err, "Should retrieve backup successfully")
	assert.Equal(t, testBackup.Id, backup.Id, "Backup ID should match")
	assert.Equal(t, "completed", backup.Status, "Backup status should match")

	// Test non-existent backup
	_, err = backupService.GetBackupByID(9999)
	assert.Error(t, err, "Should return error for non-existent backup")
}

// Test 3: Delete Backup - Should delete backup from cloud and database
func TestBackupService_DeleteBackup(t *testing.T) {
	db := setupBackupTestDB(t)

	// Create repositories and service
	backupRepo := repositories.NewBackupRepository(db)
	backupService := services.NewBackupService(backupRepo, nil, nil, "/tmp/backups")

	// Create mock cloud storage
	mockCloud := NewMockCloudStorage()
	backupService.SetCloudStorage(mockCloud)

	// Create test backup
	testBackup := createTestBackup(db, 1, 1, "completed")
	mockCloud.uploadedFiles[testBackup.Filepath] = []byte("test data")

	// Test successful deletion
	err := backupService.DeleteBackup(testBackup.Id, 1)
	assert.NoError(t, err, "Backup deletion should succeed")

	// Verify backup is deleted from database
	_, err = backupRepo.GetByID(testBackup.Id)
	assert.Error(t, err, "Deleted backup should not be retrievable")

	// Verify file is deleted from cloud storage
	_, exists := mockCloud.uploadedFiles[testBackup.Filepath]
	assert.False(t, exists, "File should be deleted from cloud storage")
}

// Test 4: Download Backup - Should download and decrypt backup file
func TestBackupService_DownloadBackup(t *testing.T) {
	db := setupBackupTestDB(t)

	// Create repositories and service
	backupRepo := repositories.NewBackupRepository(db)
	backupService := services.NewBackupService(backupRepo, nil, nil, "/tmp/backups")

	// Create mock cloud storage and encryption
	mockCloud := NewMockCloudStorage()
	backupService.SetCloudStorage(mockCloud)

	// Generate user key and encryption service
	userKey := services.GenerateUserKey(1, "SafeBaseBackupSalt2025!")
	encryptionService := services.NewEncryptionService(userKey)
	backupService.SetEncryptionService(encryptionService)

	// Create test backup
	testBackup := createTestBackup(db, 1, 1, "completed")

	// Encrypt and store test data
	testData := []byte("test backup data")
	
	// Create a temporary file to encrypt
	tmpFile := "/tmp/test_backup_data.txt"
	err := os.WriteFile(tmpFile, testData, 0644)
	require.NoError(t, err, "Should write test file successfully")
	defer os.Remove(tmpFile)
	
	encryptedData, _ := encryptionService.EncryptFile(tmpFile)
	mockCloud.uploadedFiles[testBackup.Filepath] = encryptedData

	// Test successful download
	data, err := backupService.DownloadBackup(testBackup.Id, 1)
	assert.NoError(t, err, "Backup download should succeed")
	assert.Equal(t, testData, data, "Downloaded data should match original")

	// Test unauthorized access
	_, err = backupService.DownloadBackup(testBackup.Id, 999)
	assert.Error(t, err, "Should fail with unauthorized user")
}

// Test 5: List Backups by Database - Should return all backups for a database
func TestBackupService_GetBackupsByDatabase(t *testing.T) {
	db := setupBackupTestDB(t)

	// Create repositories and service
	backupRepo := repositories.NewBackupRepository(db)
	backupService := services.NewBackupService(backupRepo, nil, nil, "/tmp/backups")

	// Create multiple test backups
	createTestBackup(db, 1, 1, "completed")
	createTestBackup(db, 1, 1, "completed")
	createTestBackup(db, 1, 1, "failed")

	// Get all backups for database
	backups, err := backupService.GetBackupsByDatabase(1)
	assert.NoError(t, err, "Should retrieve backups successfully")
	assert.Equal(t, 3, len(backups), "Should return 3 backups")

	// Verify each backup belongs to the database
	for _, backup := range backups {
		assert.Equal(t, uint(1), backup.DatabaseId, "All backups should belong to database 1")
	}

	// Test with non-existent database
	backups, err = backupService.GetBackupsByDatabase(999)
	assert.NoError(t, err, "Should not error for non-existent database")
	assert.Equal(t, 0, len(backups), "Should return empty list for non-existent database")
}

// Test 6: Encryption and Decryption - Should encrypt and decrypt data correctly
func TestBackupService_EncryptionDecryption(t *testing.T) {
	// Generate user key
	userKey := services.GenerateUserKey(1, "SafeBaseBackupSalt2025!")
	encryptionService := services.NewEncryptionService(userKey)

	testData := []byte("This is sensitive backup data that needs encryption")

	// Create a temporary file to encrypt
	tmpFile := "/tmp/test_encryption_data.txt"
	err := os.WriteFile(tmpFile, testData, 0644)
	assert.NoError(t, err, "Should create temp file")
	defer os.Remove(tmpFile)

	// Test encryption
	encryptedData, err := encryptionService.EncryptFile(tmpFile)
	assert.NoError(t, err, "Encryption should succeed")
	assert.NotEqual(t, testData, encryptedData, "Encrypted data should be different from original")

	// Test decryption
	decryptedData, err := encryptionService.DecryptData(encryptedData)
	assert.NoError(t, err, "Decryption should succeed")
	assert.Equal(t, testData, decryptedData, "Decrypted data should match original")

	// Test decryption with wrong key
	wrongKey := services.GenerateUserKey(2, "SafeBaseBackupSalt2025!")
	wrongEncryptionService := services.NewEncryptionService(wrongKey)
	_, err = wrongEncryptionService.DecryptData(encryptedData)
	assert.Error(t, err, "Decryption with wrong key should fail")
}

