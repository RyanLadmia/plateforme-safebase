package db

import (
	"fmt"
	"log"

	"github.com/RyanLadmia/plateforme-safebase/internal/config"
	"gorm.io/driver/postgres"

	// "github.com/go-sql-driver/mysql" // Pour MySQL
	"gorm.io/gorm"
)

// Connection to PostgreSQL, variables order doesn't matter for PostgreSQL
func ConnectPostgres(cfg *config.Config) *gorm.DB {
	dns := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME,
	)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf(config.Red+"Error connecting to PostgreSQL : %v"+config.Reset, err)
	}

	log.Println(config.Green + "Connection to Database successfully established" + config.Reset)
	return db
}
