package utils

import "github.com/gin-gonic/gin"

// CORSMiddleware creates a CORS middleware to allow cookies and cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Allow development and production origins
		allowedOrigins := []string{
			"http://localhost:5173", // Vite dev server (frontend)
			"http://127.0.0.1:5173", // Alternative localhost (frontend)

			"http://localhost:3000", // Go dev server port (backend)
			"http://127.0.0.1:3000", // Alternative localhost (backend)

			// Replace with your production domain (backend and frontend if separated)
		}

		// Check if the origin is allowed
		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Cookie")
		c.Header("Access-Control-Expose-Headers", "Set-Cookie")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
