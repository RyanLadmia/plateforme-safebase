package services

import (
	"errors"
	"regexp"
	"time"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/RyanLadmia/plateforme-safebase/pkg/security"
	"golang.org/x/crypto/bcrypt"
)

// AuthService gère tout ce qui concerne l'authentification
type AuthService struct {
	userRepo    *repositories.UserRepository
	sessionRepo *repositories.SessionRepository
	jwtSecret   string
	tokenTTL    time.Duration
}

// NewAuthService constructeur
func NewAuthService(userRepo *repositories.UserRepository, sessionRepo *repositories.SessionRepository, jwtSecret string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
		tokenTTL:    tokenTTL,
	}
}

// ValidatePassword vérifie le mot de passe selon les règles
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

// Register crée un nouvel utilisateur et lui assigne le rôle "user" par défaut
func (s *AuthService) Register(user *models.User) error {
	// Vérifier que l'email n'existe pas
	existingUser, _ := s.userRepo.GetUserByEmail(user.Email)
	if existingUser != nil {
		return errors.New("email invalid")
	}

	// Vérification mot de passe
	if err := ValidatePassword(user.Password); err != nil {
		return err
	}

	// Hash du mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Assignation automatique du rôle user si non précisé
	if user.RoleID == nil {
		defaultRoleID := uint(2) // correspond au rôle user seedé
		user.RoleID = &defaultRoleID
	}

	// Création en DB
	return s.userRepo.CreateUser(user)
}

// Login vérifie les identifiants et crée une session avec token JWT
func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Vérification du mot de passe
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Génération du token JWT
	token, err := security.GenerateJWT(s.jwtSecret, user.Id, user.Email, user.Role.Name, s.tokenTTL)
	if err != nil {
		return "", err
	}

	// Création de la session en DB
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

// Logout supprime la session associée au token pour déconnecter l'utilisateur
func (s *AuthService) Logout(token string) error {
	// Supprimer la session en DB (blacklist du token)
	return s.sessionRepo.DeleteByToken(token)
}
