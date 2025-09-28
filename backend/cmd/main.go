package main

import (
	"fmt"
	"log"

	"github.com/RyanLadmia/plateforme-safebase/internal/config"
	"github.com/RyanLadmia/plateforme-safebase/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to PostgreSQL
	database := db.ConnectPostgres(cfg)
	_ = database // To avoid linter error

	// Seed roles
	db.SeedRoles(database)

	// Gin server
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.SetTrustedProxies([]string{"127.0.0.1"})

	// Test route
	server.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Safebase is here!"})
	})

	// Run server
	port := cfg.PORT
	if port == "" {
		port = "8080"
	}
	fmt.Printf(config.Green+"Server running on port %s\n", port+config.Reset)
	if err := server.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
