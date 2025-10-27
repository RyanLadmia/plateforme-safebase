package security

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword : to generate a hashed password
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost = 10 → good balance between security and performance
	// It's not the number of caracter in the password, it's the cost of the hash
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Erreur de hashage du mot de passe: %v", err)
		return "", err
	}
	return string(hashed), nil
}

// CheckPassword : to check if a password is correct
func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil // true = password is correct
}
