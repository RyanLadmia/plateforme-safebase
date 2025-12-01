package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// EncryptionService handles AES-256-GCM encryption/decryption of files
type EncryptionService struct {
	key []byte
}

// NewEncryptionService creates a new encryption service with a derived key
func NewEncryptionService(masterKey string) *EncryptionService {
	// Derive a 32-byte key from the master key using SHA-256
	hash := sha256.Sum256([]byte(masterKey))
	return &EncryptionService{
		key: hash[:],
	}
}

// EncryptFile encrypts a file and returns the encrypted data
func (e *EncryptionService) EncryptFile(filePath string) ([]byte, error) {
	// Read the original file
	plainData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file for encryption: %v", err)
	}

	// Encrypt the data
	encryptedData, err := e.encryptData(plainData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}

	return encryptedData, nil
}

// DecryptData decrypts encrypted data and returns the original data
func (e *EncryptionService) DecryptData(encryptedData []byte) ([]byte, error) {
	return e.decryptData(encryptedData)
}

// encryptData performs AES-256-GCM encryption
func (e *EncryptionService) encryptData(plainData []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plainData, nil)
	return ciphertext, nil
}

// decryptData performs AES-256-GCM decryption
func (e *EncryptionService) decryptData(encryptedData []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("encrypted data too short")
	}

	nonce := encryptedData[:nonceSize]
	ciphertext := encryptedData[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %v", err)
	}

	return plaintext, nil
}

// GenerateUserKey generates a unique encryption key for a user
// This creates a deterministic key based on user ID and a master salt
func GenerateUserKey(userID uint, masterSalt string) string {
	data := fmt.Sprintf("%d:%s", userID, masterSalt)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
