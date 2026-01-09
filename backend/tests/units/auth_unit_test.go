package units

import (
	"testing"
	"time"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/RyanLadmia/plateforme-safebase/pkg/security"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// setupAuthTestDB initializes an in-memory SQLite database for testing
func setupAuthTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Role{}, &models.Session{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Create default roles
	roles := []models.Role{
		{Id: 1, Name: "admin"},
		{Id: 2, Name: "user"},
	}
	for _, role := range roles {
		db.Create(&role)
	}

	return db
}

// createTestUser creates a test user with hashed password
func createTestUser(db *gorm.DB, email, password string, roleID uint) *models.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &models.User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     email,
		Password:  string(hashedPassword),
		Active:    true,
		RoleID:    &roleID,
	}
	db.Create(user)
	return user
}

// ============================================================================
// TESTS - Authentication Service
// ============================================================================

// Test 1: Validate Password - Should validate strong passwords
func TestAuthService_ValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid strong password",
			password: "StrongP@ssw0rd",
			wantErr:  false,
		},
		{
			name:     "Password too short",
			password: "Short1!",
			wantErr:  true,
		},
		{
			name:     "Password without uppercase",
			password: "weakpassword123!",
			wantErr:  true,
		},
		{
			name:     "Password without special character",
			password: "WeakPassword123",
			wantErr:  true,
		},
		{
			name:     "Password without number",
			password: "WeakPassword!",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := services.ValidatePassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err, "Expected error for weak password")
			} else {
				assert.NoError(t, err, "Expected no error for strong password")
			}
		})
	}
}

// Test 2: Login - Should authenticate user and create session with JWT
func TestAuthService_Login(t *testing.T) {
	db := setupAuthTestDB(t)

	// Create repositories and service
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	authService := services.NewAuthService(userRepo, sessionRepo, "test-secret-key", 24*time.Hour)

	// Create test user
	testPassword := "TestP@ssw0rd123"
	testUser := createTestUser(db, "test@example.com", testPassword, 2)

	// Test successful login
	token, err := authService.Login(testUser.Email, testPassword)
	assert.NoError(t, err, "Login should succeed with correct credentials")
	assert.NotEmpty(t, token, "Token should be generated")

	// Verify session was created
	session, err := sessionRepo.GetSessionByToken(token)
	assert.NoError(t, err, "Session should exist")
	assert.Equal(t, testUser.Id, session.UserId, "Session should belong to the user")

	// Test failed login with wrong password
	_, err = authService.Login(testUser.Email, "WrongPassword123!")
	assert.Error(t, err, "Login should fail with incorrect password")

	// Test failed login with non-existent email
	_, err = authService.Login("nonexistent@example.com", testPassword)
	assert.Error(t, err, "Login should fail with non-existent email")
}

// Test 3: JWT Token Generation and Validation
func TestAuthService_JWTTokenValidation(t *testing.T) {
	secretKey := "test-secret-key-for-jwt"
	userID := uint(1)
	email := "test@example.com"
	role := "user"
	duration := 24 * time.Hour

	// Generate token
	token, err := security.GenerateJWT(secretKey, userID, email, role, duration)
	assert.NoError(t, err, "Token generation should succeed")
	assert.NotEmpty(t, token, "Token should not be empty")

	// Verify token
	claims, err := security.VerifyJWT(secretKey, token)
	assert.NoError(t, err, "Token verification should succeed")
	assert.Equal(t, userID, claims.UserID, "User ID should match")
	assert.Equal(t, email, claims.Email, "Email should match")
	assert.Equal(t, role, claims.Role, "Role should match")

	// Test invalid token
	_, err = security.VerifyJWT(secretKey, "invalid.token.here")
	assert.Error(t, err, "Invalid token should fail verification")

	// Test expired token
	expiredToken, _ := security.GenerateJWT(secretKey, userID, email, role, -1*time.Hour)
	_, err = security.VerifyJWT(secretKey, expiredToken)
	assert.Error(t, err, "Expired token should fail verification")
}

// Test 4: Logout - Should delete session
func TestAuthService_Logout(t *testing.T) {
	db := setupAuthTestDB(t)

	// Create repositories and service
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	authService := services.NewAuthService(userRepo, sessionRepo, "test-secret-key", 24*time.Hour)

	// Create test user and login
	testPassword := "TestP@ssw0rd123"
	testUser := createTestUser(db, "logout@example.com", testPassword, 2)
	token, _ := authService.Login(testUser.Email, testPassword)

	// Verify session exists
	session, err := sessionRepo.GetSessionByToken(token)
	assert.NoError(t, err, "Session should exist before logout")
	assert.NotNil(t, session, "Session should not be nil")

	// Logout
	err = authService.Logout(token)
	assert.NoError(t, err, "Logout should succeed")

	// Verify session is deleted
	_, err = sessionRepo.GetSessionByToken(token)
	assert.Error(t, err, "Session should not exist after logout")
}

// Test 5: Password Hashing - Should hash and verify passwords correctly
func TestAuthService_PasswordHashing(t *testing.T) {
	testPassword := "MySecureP@ssw0rd123"

	// Hash password
	hashedPassword, err := security.HashPassword(testPassword)
	assert.NoError(t, err, "Password hashing should succeed")
	assert.NotEmpty(t, hashedPassword, "Hashed password should not be empty")
	assert.NotEqual(t, testPassword, hashedPassword, "Hashed password should be different from plain password")

	// Verify correct password
	isValid := security.CheckPassword(hashedPassword, testPassword)
	assert.True(t, isValid, "Correct password should be valid")

	// Verify incorrect password
	isValid = security.CheckPassword(hashedPassword, "WrongPassword123!")
	assert.False(t, isValid, "Incorrect password should be invalid")

	// Test that same password generates different hashes (due to salt)
	hashedPassword2, _ := security.HashPassword(testPassword)
	assert.NotEqual(t, hashedPassword, hashedPassword2, "Same password should generate different hashes")
}
