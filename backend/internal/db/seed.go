package db

import (
	"log"

	"github.com/RyanLadmia/plateforme-safebase/internal/models"
	"gorm.io/gorm"
)

// SeedRoles : ensure that the default roles exist
func SeedRoles(db *gorm.DB) {
	var count int64

	// Create the admin role if it doesn't exist
	db.Model(&models.Role{}).Where("name = ?", "admin").Count(&count)
	if count == 0 {
		if err := db.Create(&models.Role{Name: "admin"}).Error; err != nil {
			log.Fatalf("Failed to seed role 'admin': %v", err)
		}
	}

	// Create the user role if it doesn't exist
	db.Model(&models.Role{}).Where("name = ?", "user").Count(&count)
	if count == 0 {
		if err := db.Create(&models.Role{Name: "user"}).Error; err != nil {
			log.Fatalf("Failed to seed role 'user': %v", err)
		}
	}

}
