package units

import (
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestUnit_MockDB : test just the connection logic with a fake DB (SQLite in memory)
func TestUnit_MockDB(t *testing.T) {
	// Use SQLite in memory to simulate a DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Échec de connexion à la DB mock: %v\n", err)
		fmt.Printf("Échec de connexion à la DB mock: %v\n", err)
	}

	// Check that the DB is not nil
	if db == nil {
		t.Fatal("La connexion DB mock est nil\n")
		fmt.Printf("La connexion DB mock est nil\n")
	}
	fmt.Printf("Connexion à la DB mock réussie\n")
}
