package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Structure of the configuration
type Config struct {
	PORT        string
	JWT_SECRET  string
	DB_URL      string
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
}

// ANSI color codes for console output
const (
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Reset   = "\033[0m"
)

// Load config from .env file and return a Config struct
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(Red + " .env not found, using system variables" + Reset)
	}

	return &Config{
		PORT:        os.Getenv("PORT"),
		JWT_SECRET:  os.Getenv("JWT_SECRET"),
		DB_URL:      os.Getenv("DB_URL"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_NAME:     os.Getenv("DB_NAME"),
	}
}
