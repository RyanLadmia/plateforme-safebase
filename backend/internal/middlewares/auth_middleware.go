package middleware

import (
	"net/http"
	"strings"

	"github.com/RyanLadmia/plateforme-safebase/pkg/security"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware vérifie que le token JWT est présent et valide
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupère le header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
			c.Abort()
			return
		}

		// Format attendu : "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header format must be Bearer <token>"})
			c.Abort()
			return
		}

		token := parts[1]

		// Vérifie le token
		claims, err := security.VerifyJWT(jwtSecret, token)
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
