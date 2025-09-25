package units

import (
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestUnitaire_MockDB : teste juste la logique de connexion avec une DB factice (SQLite en mémoire)
func TestUnit_MockDB(t *testing.T) {
	// On utilise SQLite en mémoire pour simuler une DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Échec de connexion à la DB mock: %v\n", err)
		fmt.Printf("Échec de connexion à la DB mock: %v\n", err)
	}

	// Vérifie que la DB n'est pas nil
	if db == nil {
		t.Fatal("La connexion DB mock est nil\n")
		fmt.Printf("La connexion DB mock est nil\n")
	}
	fmt.Printf("Connexion à la DB mock réussie\n")
}
