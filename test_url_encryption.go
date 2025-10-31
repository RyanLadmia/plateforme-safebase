package main

import (
	"fmt"
	"github.com/RyanLadmia/plateforme-safebase/pkg/security"
)

func main() {
	// Test URL encryption/decryption
	originalURL := "mysql://root:secret@localhost:8889/mydatabase"

	fmt.Printf("Testing URL encryption/decryption...\n")
	fmt.Printf("Original URL: %s\n", originalURL)

	// Encrypt
	encryptedURL, err := security.EncryptDatabasePassword(originalURL)
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}
	fmt.Printf("Encrypted URL: %s\n", encryptedURL)
	fmt.Printf("Encrypted length: %d bytes\n", len(encryptedURL))

	// Decrypt
	decryptedURL, err := security.DecryptDatabasePassword(encryptedURL)
	if err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		return
	}
	fmt.Printf("Decrypted URL: %s\n", decryptedURL)

	// Verify
	if originalURL == decryptedURL {
		fmt.Printf("✅ URL encryption/decryption test PASSED\n")
		fmt.Printf("✅ URLs are securely stored and can be decrypted for backups\n")
	} else {
		fmt.Printf("❌ URL encryption/decryption test FAILED\n")
	}
}
