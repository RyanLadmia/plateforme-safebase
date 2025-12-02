package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RyanLadmia/plateforme-safebase/internal/services"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Test Mega service
	megaConfig := services.MegaConfig{
		Email:    os.Getenv("MEGA_EMAIL"),
		Password: os.Getenv("MEGA_PASSWORD"),
	}

	if megaConfig.Email == "" || megaConfig.Password == "" {
		log.Fatal("MEGA_EMAIL and MEGA_PASSWORD environment variables must be set")
	}

	fmt.Println("Testing Mega service initialization...")
	megaService, err := services.NewMegaService(megaConfig)
	if err != nil {
		log.Fatalf("Failed to initialize Mega service: %v", err)
	}

	fmt.Println("Testing file upload...")
	err = megaService.UploadFile("test_backup.txt", "TestUser/mysql/test_backup.txt")
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}

	fmt.Println("Testing file existence check...")
	exists, err := megaService.FileExists("TestUser/mysql/test_backup.txt")
	if err != nil {
		log.Fatalf("Failed to check file existence: %v", err)
	}

	if exists {
		fmt.Println("✅ File exists in Mega!")
	} else {
		fmt.Println("❌ File does not exist in Mega")
	}

	fmt.Println("Testing file download...")
	data, err := megaService.DownloadFile("TestUser/mysql/test_backup.txt")
	if err != nil {
		log.Fatalf("Failed to download file: %v", err)
	}

	fmt.Printf("Downloaded content: %s\n", string(data))

	fmt.Println("Testing file deletion...")
	err = megaService.DeleteFile("TestUser/mysql/test_backup.txt")
	if err != nil {
		log.Fatalf("Failed to delete file: %v", err)
	}

	fmt.Println("✅ All Mega operations successful!")
}
