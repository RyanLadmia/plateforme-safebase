package config

import (
	"os"

	"github.com/RyanLadmia/plateforme-safebase/internal/services"
)

// GetMegaConfig returns Mega configuration from environment variables
func GetMegaConfig() *services.MegaConfig {
	return &services.MegaConfig{
		Email:    os.Getenv("MEGA_EMAIL"),
		Password: os.Getenv("MEGA_PASSWORD"),
	}
}
