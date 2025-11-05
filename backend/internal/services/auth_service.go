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

	// Check if user is active
	if !user.Active {
		return "", errors.New("account is disabled")
	}

	// OPTIMISATION 1: Nettoyer les sessions expirées avant de créer une nouvelle
	if err := s.sessionRepo.DeleteExpiredSessions(); err != nil {
		log.Printf("Avertissement: Impossible de nettoyer les sessions expirées: %v", err)
	}

	// OPTIMISATION 2: Supprimer les anciennes sessions de cet utilisateur (une seule session active par utilisateur)
	if err := s.sessionRepo.DeleteByUserId(user.Id); err != nil {
		log.Printf("Avertissement: Impossible de supprimer les anciennes sessions: %v", err)
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

	log.Printf("Session créée pour l'utilisateur %d, token: %s...", user.Id, token[:10])
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

	// Check if user is active
	if !user.Active {
		return nil, errors.New("account is disabled")
	}

	return user, nil
}

// Logout delete the session associated with the token to disconnect the user
func (s *AuthService) Logout(token string) error {
	// Extraire le token du format "Bearer <token>" si nécessaire
	cleanToken := token
	if len(token) > 7 && token[:7] == "Bearer " {
		cleanToken = token[7:]
	}

	// Vérifier que la session existe avant de la supprimer
	session, err := s.sessionRepo.GetSessionByToken(cleanToken)
	if err != nil {
		log.Printf("Session non trouvée pour la déconnexion: %v", err)
		return errors.New("session not found")
	}

	// Supprimer la session spécifique
	if err := s.sessionRepo.DeleteByToken(token); err != nil {
		log.Printf("Erreur lors de la suppression de la session: %v", err)
		return err
	}

	log.Printf("Session supprimée pour l'utilisateur %d", session.UserId)
	return nil
}

// CleanupExpiredSessions nettoie périodiquement les sessions expirées
func (s *AuthService) CleanupExpiredSessions() error {
	if err := s.sessionRepo.DeleteExpiredSessions(); err != nil {
		log.Printf("Erreur lors du nettoyage des sessions expirées: %v", err)
		return err
	}
	log.Printf("Nettoyage des sessions expirées effectué")
	return nil
}

// GetActiveSessionsCount retourne le nombre de sessions actives (pour monitoring)
func (s *AuthService) GetActiveSessionsCount() (int64, error) {
	return s.sessionRepo.GetActiveSessionsCount()
}
