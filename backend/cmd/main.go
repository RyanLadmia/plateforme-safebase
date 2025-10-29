package main

import (
	"fmt"
	"log"
	"time"

	"github.com/RyanLadmia/plateforme-safebase/internal/config"
	"github.com/RyanLadmia/plateforme-safebase/internal/db"
	"github.com/RyanLadmia/plateforme-safebase/internal/handlers"
	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/RyanLadmia/plateforme-safebase/internal/routes"
	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config from environment variables
	cfg := config.LoadConfig()

	// Connection to PostgreSQL database
	database := db.ConnectPostgres(cfg)

	// Automatic migration of tables (creation/update of structures)
	// For now, we only migrate the models necessary for authentication
	log.Println(config.Yellow + "Running database migrations..." + config.Reset)

	// Migration spéciale pour supprimer la colonne deleted_at de la table sessions
	log.Println(config.Yellow + "Removing deleted_at column from sessions table..." + config.Reset)
	if database.Migrator().HasColumn(&models.Session{}, "deleted_at") {
		if err := database.Migrator().DropColumn(&models.Session{}, "deleted_at"); err != nil {
			log.Printf(config.Yellow+"Warning: Could not drop deleted_at column: %v"+config.Reset, err)
		} else {
			log.Println(config.Green + "deleted_at column removed from sessions table" + config.Reset)
		}
	}

	if err := database.AutoMigrate(
		&models.Role{},    // Role table
		&models.User{},    // User table
		&models.Session{}, // Session table
	); err != nil {
		log.Fatalf(config.Red+"Failed to migrate database: %v"+config.Reset, err)
	}

	// Add other models later when the relationships are finalized
	// &models.Alert{},   // Alert table (for later)
	// &models.Database{}, // Database table (for later)
	// &models.Backup{},  // Backup table (for later)
	// &models.Restore{}, // Restore table (for later)
	log.Println(config.Green + "Database migrations completed successfully" + config.Reset)

	// Initialize default roles (admin, user)
	db.SeedRoles(database)

	// Initialize repositories (data access layer)
	userRepo := repositories.NewUserRepository(database)
	sessionRepo := repositories.NewSessionRepository(database)

	// Initialize services (business logic)
	authService := services.NewAuthService(
		userRepo,
		sessionRepo,
		cfg.JWT_SECRET, // Secret key to sign JWT tokens
		24*time.Hour,   // Token validity duration (24h)
	)

	// Nettoyage initial des sessions expirées au démarrage
	log.Println(config.Yellow + "Nettoyage des sessions expirées..." + config.Reset)
	if err := authService.CleanupExpiredSessions(); err != nil {
		log.Printf(config.Yellow+"Avertissement: %v"+config.Reset, err)
	}

	// Afficher le nombre de sessions actives
	if count, err := authService.GetActiveSessionsCount(); err == nil {
		log.Printf(config.Green+"Sessions actives: %d"+config.Reset, count)
	}

	// Initialize handlers (HTTP controllers)
	authHandler := handlers.NewAuthHandler(authService)

	// Configure the Gin server
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.SetTrustedProxies([]string{"127.0.0.1"})

	// Secure CORS middleware to allow cookies
	server.Use(func(c *gin.Context) {
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
	})

	// Test route to check if the server is running
	server.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Safebase API is running!"})
	})

	// Integration of authentication routes (/auth/register, /auth/login, /auth/logout)
	routes.AuthRoutes(server, authHandler, cfg.JWT_SECRET)

	// Start the server
	port := cfg.PORT
	if port == "" {
		port = "3000"
	}
	fmt.Printf(config.Green+"🚀 Server running on port %s\n", port+config.Reset)
	fmt.Printf(config.Cyan + "📋 Available endpoints:\n")
	fmt.Printf("   GET  /test            - Test endpoint\n")
	fmt.Printf("   POST /auth/register   - User registration\n")
	fmt.Printf("   POST /auth/login      - User login\n")
	fmt.Printf("   POST /auth/logout     - User logout\n")
	fmt.Printf("   GET  /auth/me         - Get current user\n" + config.Reset)

	if err := server.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
