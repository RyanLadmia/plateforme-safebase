package config

import (
	"os"
	"strconv"

	"github.com/RyanLadmia/plateforme-safebase/internal/services"
)

// GetMinIOConfig returns MinIO configuration from environment variables
func GetMinIOConfig() *services.MinIOConfig {
	return &services.MinIOConfig{
		Endpoint:        getEnv("MINIO_ENDPOINT", "localhost:9000"),
		AccessKeyID:     getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		SecretAccessKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
		BucketName:      getEnv("MINIO_BUCKET", "safebase"),
		UseSSL:          getEnvAsBool("MINIO_USE_SSL", false),
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsBool gets an environment variable as boolean with a fallback value
func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}
