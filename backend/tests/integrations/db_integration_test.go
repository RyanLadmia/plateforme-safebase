package integrations

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestIntegration_ConnectPostgres : connects to the real PostgreSQL database using .env or environment variables
func TestIntegration_ConnectPostgres(t *testing.T) {
	// Try to load .env file (for local development), but don't fail if it doesn't exist
	_ = godotenv.Load("../../.env")

	// Read environment variables (from .env or CI/CD environment)
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// If variables are not set, try to use PostgreSQL service variables from CI/CD
	if host == "" {
		host = os.Getenv("POSTGRES_HOST")
		if host == "" {
			host = "localhost"
		}
	}
	if port == "" {
		port = os.Getenv("POSTGRES_PORT")
		if port == "" {
			port = "5432"
		}
	}
	if user == "" {
		user = os.Getenv("POSTGRES_USER")
		if user == "" {
			user = "testuser"
		}
	}
	if password == "" {
		password = os.Getenv("POSTGRES_PASSWORD")
		if password == "" {
			password = "testpassword"
		}
	}
	if dbname == "" {
		dbname = os.Getenv("POSTGRES_DB")
		if dbname == "" {
			dbname = "testdb"
		}
	}

	// Build DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port,
	)

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v\n", err)
		fmt.Printf("Failed to connect to PostgreSQL: %v\n", err)
	}

	// Check connection
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get SQL DB: %v\n", err)
		fmt.Printf("Failed to get SQL DB: %v\n", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("Failed to ping PostgreSQL DB: %v\n", err)
		fmt.Printf("Failed to ping PostgreSQL DB: %v\n", err)
	}
	fmt.Printf("Connected to PostgreSQL DB\n")
}
