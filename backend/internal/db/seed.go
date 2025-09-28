package db

import (
	"log"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"gorm.io/gorm"
)

// SeedRoles : assure que les rôles par défaut existent
func SeedRoles(db *gorm.DB) {
	var count int64

	// Crée le rôle "admin" si il n'existe pas
	db.Model(&models.Role{}).Where("name = ?", "admin").Count(&count)
	if count == 0 {
		if err := db.Create(&models.Role{Name: "admin"}).Error; err != nil {
			log.Fatalf("Failed to seed role 'admin': %v", err)
		}
	}

	// Crée le rôle "user" si il n'existe pas
	db.Model(&models.Role{}).Where("name = ?", "user").Count(&count)
	if count == 0 {
		if err := db.Create(&models.Role{Name: "user"}).Error; err != nil {
			log.Fatalf("Failed to seed role 'user': %v", err)
		}
	}

}
