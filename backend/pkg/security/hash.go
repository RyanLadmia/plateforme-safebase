package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
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
// plainPassword is the password that the user entered (visible in the input field not in database)
func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil // true = password is correct
}

// EncryptDatabasePassword encrypts a database password using AES-GCM (AEAD mode)
func EncryptDatabasePassword(password string) (string, error) {
	// In production, this key should come from environment variables
	key := []byte("your-32-byte-secret-key-here!!!!") // 32 bytes for AES-256

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext := []byte(password)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptDatabasePassword decrypts a database password using AES-GCM (AEAD mode)
func DecryptDatabasePassword(encryptedPassword string) (string, error) {
	// In production, this key should come from environment variables
	key := []byte("your-32-byte-secret-key-here!!!!") // 32 bytes for AES-256

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
