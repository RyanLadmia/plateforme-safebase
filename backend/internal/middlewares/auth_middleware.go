package middlewares

import (
	"net/http"
	"strings"

	"github.com/RyanLadmia/plateforme-safebase/pkg/security"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware structure pour encapsuler le middleware d'authentification
type AuthMiddleware struct {
	jwtSecret string
}

// NewAuthMiddleware crée une nouvelle instance d'AuthMiddleware
func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}

// RequireAuth vérifie que le token JWT est présent et valide
func (am *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// Essayer d'abord de récupérer le token depuis le header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// Format attendu : "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// Si pas de token dans le header, essayer le cookie HTTP-only (plus sécurisé)
		if token == "" {
			cookieToken, err := c.Cookie("auth_token")
			if err == nil && cookieToken != "" {
				token = cookieToken
			}
		}

		// Si aucun token trouvé
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		// Vérifie le token
		claims, err := security.VerifyJWT(am.jwtSecret, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Stocke les claims dans le contexte pour qu'ils soient accessibles dans le handler
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// RequireRole vérifie que l'utilisateur a le rôle requis
func (am *AuthMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupère le rôle depuis le contexte (défini par RequireAuth)
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "role information missing"})
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid role format"})
			c.Abort()
			return
		}

		if roleStr != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
