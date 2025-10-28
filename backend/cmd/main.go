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
	// Chargement de la configuration depuis les variables d'environnement
	cfg := config.LoadConfig()

	// Connexion à PostgreSQL
	database := db.ConnectPostgres(cfg)

	// Migration automatique des tables (création/mise à jour des structures)
	// Pour l'instant, nous migrons seulement les modèles nécessaires à l'authentification
	log.Println(config.Yellow + "Running database migrations..." + config.Reset)
	if err := database.AutoMigrate(
		&models.Role{},    // Table des rôles
		&models.User{},    // Table des utilisateurs
		&models.Session{}, // Table des sessions
	); err != nil {
		log.Fatalf(config.Red+"Failed to migrate database: %v"+config.Reset, err)
	}

	// TODO: Ajouter plus tard les autres modèles quand les relations seront finalisées
	// &models.Alert{},   // Table des alertes (pour plus tard)
	// &models.Database{}, // Table des bases de données (pour plus tard)
	// &models.Backup{},  // Table des sauvegardes (pour plus tard)
	// &models.Restore{}, // Table des restaurations (pour plus tard)
	log.Println(config.Green + "Database migrations completed successfully" + config.Reset)

	// Initialisation des rôles par défaut (admin, user)
	db.SeedRoles(database)

	// Initialisation des repositories (couche d'accès aux données)
	userRepo := repositories.NewUserRepository(database)
	sessionRepo := repositories.NewSessionRepository(database)

	// Initialisation des services (logique métier)
	authService := services.NewAuthService(
		userRepo,
		sessionRepo,
		cfg.JWT_SECRET, // Clé secrète pour signer les JWT
		24*time.Hour,   // Durée de validité des tokens (24h)
	)

	// Initialisation des handlers (contrôleurs HTTP)
	authHandler := handlers.NewAuthHandler(authService)

	// Configuration du serveur Gin
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.SetTrustedProxies([]string{"127.0.0.1"})

	// Middleware CORS sécurisé pour permettre les cookies
	server.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Autoriser les origines de développement et production
		allowedOrigins := []string{
			"http://localhost:5173",      // Vite dev server
			"http://localhost:3000",      // Alternative dev port
			"http://127.0.0.1:5173",      // Alternative localhost
			"https://yourdomain.com",     // TODO: Remplacer par votre domaine de production
			"https://www.yourdomain.com", // TODO: Remplacer par votre domaine de production
		}

		// Vérifier si l'origine est autorisée
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

	// Route de test pour vérifier que le serveur fonctionne
	server.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Safebase API is running!"})
	})

	// Intégration des routes d'authentification (/auth/register, /auth/login, /auth/logout)
	routes.AuthRoutes(server, authHandler, cfg.JWT_SECRET)

	// Démarrage du serveur
	port := cfg.PORT
	if port == "" {
		port = "8080"
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
