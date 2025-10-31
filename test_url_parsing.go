package main

import (
	"fmt"
	"github.com/RyanLadmia/plateforme-safebase/internal/models"
)

func main() {
	// Test URL parsing for different database types
	testCases := []struct {
		dbType string
		url    string
	}{
		{"mysql", "mysql://testuser:testpass@remotehost:3306/testdb"},
		{"postgresql", "postgresql://testuser:testpass@remotehost:5432/testdb"},
		{"mysql", "mysql://user:pass@tcp(localhost:8889)/mydb"},
	}

	for _, tc := range testCases {
		fmt.Printf("\nTesting %s URL: %s\n", tc.dbType, tc.url)

		host, port, username, password, dbName, err := models.ParseDatabaseURL(tc.dbType, tc.url)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}

		fmt.Printf("✅ Parsed successfully:\n")
		fmt.Printf("   Host: %s\n", host)
		fmt.Printf("   Port: %s\n", port)
		fmt.Printf("   Username: %s\n", username)
		fmt.Printf("   Password: %s\n", password)
		fmt.Printf("   Database: %s\n", dbName)
	}

	// Test validation
	fmt.Printf("\n--- Testing validation ---\n")

	req := &models.DatabaseCreateRequest{
		Name:     "Test DB",
		Type:     "mysql",
		URL:      "mysql://user:pass@host:3306/db",
	}

	err := models.ValidateAndNormalizeDatabaseData(req)
	if err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Printf("✅ Validation passed. Normalized data:\n")
		fmt.Printf("   Host: %s\n", req.Host)
		fmt.Printf("   Port: %s\n", req.Port)
		fmt.Printf("   Username: %s\n", req.Username)
		fmt.Printf("   Password length: %d\n", len(req.Password))
		fmt.Printf("   Database: %s\n", req.DbName)
	}
}
