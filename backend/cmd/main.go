package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/RyanLadmia/plateforme-safebase/internal/config"
	"github.com/RyanLadmia/plateforme-safebase/internal/db"
	"github.com/RyanLadmia/plateforme-safebase/internal/handlers"
	"github.com/RyanLadmia/plateforme-safebase/internal/middlewares"
	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"github.com/RyanLadmia/plateforme-safebase/internal/repositories"
	"github.com/RyanLadmia/plateforme-safebase/internal/routes"
	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/RyanLadmia/plateforme-safebase/utils"
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

	if err := database.AutoMigrate(
		&models.Role{},     // Role table
		&models.User{},     // User table
		&models.Session{},  // Session table
		&models.Database{}, // Database table
		&models.Backup{},   // Backup table
		&models.Schedule{}, // Schedule table
		&models.Restore{},  // Restore table
		// &models.Alert{},   // Alert table (for later)
	); err != nil {
		log.Fatalf(config.Red+"Failed to migrate database: %v"+config.Reset, err)
	}
	log.Println(config.Green + "Database migrations completed successfully" + config.Reset)

	// Initialize default roles (admin, user)
	db.SeedRoles(database)

	// Initialize repositories (data access layer)
	userRepo := repositories.NewUserRepository(database)
	sessionRepo := repositories.NewSessionRepository(database)
	databaseRepo := repositories.NewDatabaseRepository(database)
	backupRepo := repositories.NewBackupRepository(database)
	scheduleRepo := repositories.NewScheduleRepository(database)
	roleRepo := repositories.NewRoleRepository(database)
	restoreRepo := repositories.NewRestoreRepository(database)

	// Initialize services (business logic)
	authService := services.NewAuthService(
		userRepo,
		sessionRepo,
		cfg.JWT_SECRET, // Secret key to sign JWT tokens
		24*time.Hour,   // Token validity duration (24h)
	)

	// Initialize backup service with backup directory
	backupDir := filepath.Join(".", "db", "backups")
	databaseService := services.NewDatabaseService(databaseRepo)
	userService := services.NewUserService(userRepo, roleRepo)
	backupService := services.NewBackupService(backupRepo, databaseService, userService, backupDir)
	scheduleService := services.NewScheduleService(scheduleRepo, databaseRepo, backupService)
	restoreService := services.NewRestoreService(restoreRepo, backupService, databaseService, userService)

	// Initialize Mega service for cloud storage
	megaConfig := config.GetMegaConfig()
	if megaConfig.Email != "" && megaConfig.Password != "" {
		megaService, err := services.NewMegaService(*megaConfig)
		if err != nil {
			log.Printf(config.Yellow+"Avertissement: Impossible d'initialiser Mega: %v"+config.Reset, err)
			log.Println(config.Yellow + "Les sauvegardes seront stockées localement uniquement" + config.Reset)
		} else {
			log.Println(config.Green + "Service Mega initialisé avec succès" + config.Reset)
			backupService.SetCloudStorage(megaService)

			// Initialize encryption service
			encryptionService := services.NewEncryptionService("SafeBaseMasterKey2025!")
			backupService.SetEncryptionService(encryptionService)
			log.Println(config.Green + "Service de chiffrement AES-256 initialisé" + config.Reset)
		}
	} else {
		log.Println(config.Yellow + "Configuration Mega manquante - stockage local uniquement" + config.Reset)
	}

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
	databaseHandler := handlers.NewDatabaseHandler(databaseService)
	backupHandler := handlers.NewBackupHandler(backupService)
	scheduleHandler := handlers.NewScheduleHandler(scheduleService)
	userHandler := handlers.NewUserHandler(userService)
	restoreHandler := handlers.NewRestoreHandler(restoreService)

	// Initialize middleware
	authMiddleware := middlewares.NewAuthMiddleware(cfg.JWT_SECRET)

	// Configure the Gin server
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.SetTrustedProxies([]string{"127.0.0.1"})

	// Secure CORS middleware to allow cookies
	server.Use(utils.CORSMiddleware())

	// Test route to check if the server is running
	server.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Safebase API is running!"})
	})

	// Integration of authentication routes (/auth/register, /auth/login, /auth/logout)
	routes.AuthRoutes(server, authHandler, cfg.JWT_SECRET)

	// Integration of database and backup routes
	routes.SetupDatabaseRoutes(server, databaseHandler, authMiddleware)
	routes.SetupBackupRoutes(server, backupHandler, authMiddleware)
	routes.SetupScheduleRoutes(server, scheduleHandler, authMiddleware)
	routes.SetupRestoreRoutes(server, restoreHandler, authMiddleware)
	routes.UserRoutes(server, userHandler, authMiddleware)

	// Initialize worker pool for background tasks
	workerPool := utils.NewWorkerPool(5) // 5 workers
	workerPool.Start()

	// Start background workers
	go utils.StartSessionCleanupWorker(sessionRepo)
	go utils.StartBackupCleanupWorker(backupRepo, workerPool)

	// Pass worker pool to backup service
	backupService.SetWorkerPool(workerPool)
	restoreService.SetWorkerPool(workerPool)

	// Start the cron scheduler and load active schedules
	scheduleService.StartScheduler()
	if err := scheduleService.LoadActiveSchedules(); err != nil {
		log.Printf(config.Yellow+"Avertissement lors du chargement des schedules: %v"+config.Reset, err)
	}

	// Start the server
	port := cfg.PORT
	if port == "" {
		port = "3000"
	}
	fmt.Printf(config.Green+"Server running on port %s\n", port+config.Reset)
	utils.DisplayEndpoints(port)

	if err := server.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
