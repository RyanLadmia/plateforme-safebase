package main

import (
	"fmt"
	"github.com/RyanLadmia/plateforme-safebase/pkg/security"
)

func main() {
	// Test password encryption/decryption
	originalPassword := "root"

	fmt.Printf("Testing password encryption/decryption...\n")
	fmt.Printf("Original password: %s\n", originalPassword)

	// Encrypt
	encrypted, err := security.EncryptDatabasePassword(originalPassword)
	if err != nil {
		fmt.Printf("Encryption failed: %v\n", err)
		return
	}
	fmt.Printf("Encrypted password: %s\n", encrypted)

	// Decrypt
	decrypted, err := security.DecryptDatabasePassword(encrypted)
	if err != nil {
		fmt.Printf("Decryption failed: %v\n", err)
		return
	}
	fmt.Printf("Decrypted password: %s\n", decrypted)

	// Verify
	if originalPassword == decrypted {
		fmt.Printf("✅ Encryption/decryption test PASSED\n")
	} else {
		fmt.Printf("❌ Encryption/decryption test FAILED\n")
	}
}
