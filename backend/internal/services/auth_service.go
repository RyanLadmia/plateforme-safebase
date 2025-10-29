package services

import (
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/RyanLadmia/plateforme-safebase/pkg/security"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService manages everything related to authentication
type AuthService struct {
	userRepo    *repositories.UserRepository
	sessionRepo *repositories.SessionRepository
	jwtSecret   string
	tokenTTL    time.Duration
}

// NewAuthService constructor
func NewAuthService(userRepo *repositories.UserRepository, sessionRepo *repositories.SessionRepository, jwtSecret string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
		tokenTTL:    tokenTTL,
	}
}

// ValidatePassword check the password according to the rules
func ValidatePassword(password string) error {
	if len(password) < 10 {
		return errors.New("password must be at least 10 characters")
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("password must include uppercase, lowercase, number, and special character")
	}
	return nil
}

// Register create a new user and assign the "user" role by default
func (s *AuthService) Register(user *models.User) error {
	// Check if the email already exists
	existingUser, _ := s.userRepo.GetUserByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email invalid")
	}

	// Check password
	if err := ValidatePassword(user.Password); err != nil {
		return err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Automatic assignment of the user role if not specified
	if user.RoleID == nil {
		// Get the ID of the "user" role from the database
		var userRole models.Role
		if err := s.userRepo.GetDB().Where("name = ?", "user").First(&userRole).Error; err != nil {
			return errors.New("default user role not found")
		}
		user.RoleID = &userRole.Id
	}

	// Creation in DB, here we use hte userRepository function CreateUser
	log.Printf("Tentative de création de l'utilisateur en base de données...")
	if err := s.userRepo.CreateUser(user); err != nil {
		log.Printf("Erreur lors de la création en DB: %v", err)
		return err
	}
	log.Printf("Utilisateur créé en DB avec l'ID: %d", user.Id)
	return nil
}

// Login check the credentials and create a session with a JWT token
func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := security.GenerateJWT(s.jwtSecret, user.Id, user.Email, user.Role.Name, s.tokenTTL)
	if err != nil {
		return "", err
	}

	// Create session in DB
	session := &models.Session{
		UserId:    user.Id,
		Token:     token,
		ExpiresAt: time.Now().Add(s.tokenTTL),
	}
	if err := s.sessionRepo.CreateSession(session); err != nil {
		return "", err
	}

	return token, nil
}

// GetUserFromToken get a user from a JWT token
func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Get the user ID from the claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}
	userID := uint(userIDFloat)

	// Get the user from the database
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// Logout delete the session associated with the token to disconnect the user
func (s *AuthService) Logout(token string) error {
	// Delete the session in DB (blacklist the token)
	return s.sessionRepo.DeleteByToken(token)
}
