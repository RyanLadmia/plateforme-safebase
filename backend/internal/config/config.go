package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

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

func loadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(" .env non trouvé, on utilise les variables système")
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
