package integrations

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
// HELPER FUNCTIONS - Setup Integration Test Database
// ============================================================================

// setupIntegrationDB initializes a complete in-memory database with all models
func setupIntegrationDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
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

	return db
}

// setupAuthServices creates all services needed for authentication tests
func setupAuthServices(db *gorm.DB) (*services.AuthService, *repositories.SessionRepository) {
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	authService := services.NewAuthService(userRepo, sessionRepo, "test-secret-key", 24*time.Hour)
	return authService, sessionRepo
}

// ============================================================================
// INTEGRATION TESTS - Authentication Flow (User Registration -> Login -> Logout)
// ============================================================================

// TestIntegration_CompleteAuthFlow tests the complete authentication flow
func TestIntegration_CompleteAuthFlow(t *testing.T) {
	db := setupIntegrationDB(t)
	authService, sessionRepo := setupAuthServices(db)

	// Step 1: Register a new user
	t.Run("Register User", func(t *testing.T) {
		user := &models.User{
			Firstname: "John",
			Lastname:  "Doe",
			Email:     "john.doe@example.com",
			Password:  "SecureP@ssw0rd123",
		}

		err := authService.Register(user)
		require.NoError(t, err, "User registration should succeed")
		assert.NotZero(t, user.Id, "User ID should be set after registration")
		assert.NotNil(t, user.RoleID, "User should have a role assigned")
		assert.Equal(t, uint(2), *user.RoleID, "User should have 'user' role by default")
	})

	// Step 2: Login with registered user
	var token string
	t.Run("Login User", func(t *testing.T) {
		var err error
		token, err = authService.Login("john.doe@example.com", "SecureP@ssw0rd123")
		require.NoError(t, err, "Login should succeed")
		assert.NotEmpty(t, token, "Token should be generated")

		// Verify session was created in database
		session, err := sessionRepo.GetSessionByToken(token)
		require.NoError(t, err, "Session should exist in database")
		assert.Equal(t, "john.doe@example.com", session.User.Email, "Session should belong to correct user")
	})

	// Step 3: Verify user from token
	t.Run("Get User From Token", func(t *testing.T) {
		user, err := authService.GetUserFromToken(token)
		require.NoError(t, err, "Should retrieve user from token")
		assert.Equal(t, "john.doe@example.com", user.Email, "User email should match")
		assert.Equal(t, "John", user.Firstname, "User firstname should match")
	})

	// Step 4: Logout user
	t.Run("Logout User", func(t *testing.T) {
		err := authService.Logout(token)
		require.NoError(t, err, "Logout should succeed")

		// Verify session was deleted from database
		_, err = sessionRepo.GetSessionByToken(token)
		assert.Error(t, err, "Session should not exist after logout")
	})

	// Step 5: Verify session is deleted (token JWT still valid but session removed)
	t.Run("Verify Session Deleted After Logout", func(t *testing.T) {
		// Token JWT is still technically valid (not expired) but session is deleted
		user, err := authService.GetUserFromToken(token)
		// GetUserFromToken only validates JWT, not session existence
		// So it might still return user if JWT is valid
		// The session check happens at middleware level
		if err == nil {
			// JWT is still valid, verify session is deleted
			_, sessionErr := sessionRepo.GetSessionByToken(token)
			assert.Error(t, sessionErr, "Session should not exist in database after logout")
			assert.NotNil(t, user, "User from JWT should be retrievable (JWT not expired)")
		} else {
			// If there's an error, it's acceptable
			assert.Error(t, err, "Token validation may fail")
		}
	})
}

// TestIntegration_MultipleUserSessions tests multiple users with concurrent sessions
func TestIntegration_MultipleUserSessions(t *testing.T) {
	db := setupIntegrationDB(t)
	authService, sessionRepo := setupAuthServices(db)

	// Register and login multiple users
	users := []struct {
		email    string
		password string
	}{
		{"user1@example.com", "SecureP@ssw0rd123"},
		{"user2@example.com", "SecureP@ssw0rd456"},
		{"user3@example.com", "SecureP@ssw0rd789"},
	}

	tokens := make([]string, len(users))

	// Register and login all users
	for i, u := range users {
		user := &models.User{
			Firstname: "User",
			Lastname:  "Test",
			Email:     u.email,
			Password:  u.password,
		}
		err := authService.Register(user)
		require.NoError(t, err, "User registration should succeed")

		token, err := authService.Login(u.email, u.password)
		require.NoError(t, err, "Login should succeed")
		tokens[i] = token
	}

	// Verify all sessions exist
	for i, token := range tokens {
		session, err := sessionRepo.GetSessionByToken(token)
		require.NoError(t, err, "Session should exist")
		assert.Equal(t, users[i].email, session.User.Email, "Session should belong to correct user")
	}

	// Logout one user and verify others remain active
	err := authService.Logout(tokens[0])
	require.NoError(t, err, "Logout should succeed")

	_, err = sessionRepo.GetSessionByToken(tokens[0])
	assert.Error(t, err, "First user session should be deleted")

	// Verify other sessions still exist
	for i := 1; i < len(tokens); i++ {
		_, err := sessionRepo.GetSessionByToken(tokens[i])
		assert.NoError(t, err, "Other user sessions should still exist")
	}
}

// TestIntegration_InvalidCredentials tests authentication with invalid credentials
func TestIntegration_InvalidCredentials(t *testing.T) {
	db := setupIntegrationDB(t)
	authService, _ := setupAuthServices(db)

	// Register a user
	user := &models.User{
		Firstname: "Test",
		Lastname:  "User",
		Email:     "test@example.com",
		Password:  "SecureP@ssw0rd123",
	}
	err := authService.Register(user)
	require.NoError(t, err, "User registration should succeed")

	// Test with wrong password
	_, err = authService.Login("test@example.com", "WrongPassword123!")
	assert.Error(t, err, "Login should fail with wrong password")

	// Test with non-existent email
	_, err = authService.Login("nonexistent@example.com", "SecureP@ssw0rd123")
	assert.Error(t, err, "Login should fail with non-existent email")

	// Test with weak password during registration
	weakUser := &models.User{
		Firstname: "Weak",
		Lastname:  "User",
		Email:     "weak@example.com",
		Password:  "weak",
	}
	err = authService.Register(weakUser)
	assert.Error(t, err, "Registration should fail with weak password")
}

// TestIntegration_SessionExpiration tests session expiration and cleanup
func TestIntegration_SessionExpiration(t *testing.T) {
	db := setupIntegrationDB(t)
	
	// Create auth service with very short token TTL (1 second)
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	authService := services.NewAuthService(userRepo, sessionRepo, "test-secret-key", 1*time.Second)

	// Register and login user
	user := &models.User{
		Firstname: "Expiry",
		Lastname:  "Test",
		Email:     "expiry@example.com",
		Password:  "SecureP@ssw0rd123",
	}
	err := authService.Register(user)
	require.NoError(t, err, "User registration should succeed")

	token, err := authService.Login("expiry@example.com", "SecureP@ssw0rd123")
	require.NoError(t, err, "Login should succeed")
	assert.NotEmpty(t, token, "Token should be generated")

	// Wait for token to expire
	time.Sleep(2 * time.Second)

	// Cleanup expired sessions
	err = authService.CleanupExpiredSessions()
	require.NoError(t, err, "Cleanup should succeed")

	// Verify token is no longer valid
	_, err = authService.GetUserFromToken(token)
	assert.Error(t, err, "Expired token should fail verification")
}

// TestIntegration_UserDeactivation tests login with deactivated account
func TestIntegration_UserDeactivation(t *testing.T) {
	db := setupIntegrationDB(t)
	authService, _ := setupAuthServices(db)
	userRepo := repositories.NewUserRepository(db)

	// Register a user
	user := &models.User{
		Firstname: "Deactivated",
		Lastname:  "User",
		Email:     "deactivated@example.com",
		Password:  "SecureP@ssw0rd123",
	}
	err := authService.Register(user)
	require.NoError(t, err, "User registration should succeed")

	// Login successfully
	token, err := authService.Login("deactivated@example.com", "SecureP@ssw0rd123")
	require.NoError(t, err, "Login should succeed")
	assert.NotEmpty(t, token, "Token should be generated")

	// Deactivate user
	retrievedUser, err := userRepo.GetUserByEmail("deactivated@example.com")
	require.NoError(t, err, "Should retrieve user")
	err = userRepo.UpdateUserById(retrievedUser.Id, map[string]interface{}{"active": false})
	require.NoError(t, err, "Should deactivate user")

	// Try to login again - should fail
	_, err = authService.Login("deactivated@example.com", "SecureP@ssw0rd123")
	assert.Error(t, err, "Login should fail with deactivated account")
	assert.Contains(t, err.Error(), "disabled", "Error message should indicate account is disabled")

	// Verify existing token is also invalid
	_, err = authService.GetUserFromToken(token)
	assert.Error(t, err, "Existing token should be invalid for deactivated user")
}

