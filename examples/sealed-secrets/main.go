package main

import (
	"fmt"
	"log"

	"github.com/cloudresty/go-env"
)

func main() {
	fmt.Println("=== Sealed Secrets (Encryption) Example ===")

	// 1. Setting up encryption
	fmt.Println("1. Setting up encryption:")

	encryptionKey := "my-super-secret-encryption-key-2024-secure"
	err := env.Set("ENV_SEAL_KEY", encryptionKey)
	if err != nil {
		log.Fatalf("Failed to set encryption key: %v", err)
	}

	fmt.Println("  ✓ Set encryption key (ENV_SEAL_KEY)")

	// 2. Encrypting and storing secrets
	fmt.Println("\n2. Encrypting and storing secrets:")

	secrets := map[string]string{
		"DATABASE_PASSWORD": "super-secret-db-password-123!",
		"API_TOKEN":         "api-token-xyz789-very-secret",
		"JWT_SECRET":        "jwt-signing-key-that-must-be-secure",
	}

	for key, value := range secrets {
		err := env.SetSealed(key+"_SEALED", value, "ENV_SEAL_KEY")
		if err != nil {
			log.Fatalf("Failed to seal %s: %v", key, err)
		}
		fmt.Printf("  ✓ Sealed %s\n", key)

		// Show the encrypted value (truncated for display)
		encryptedValue := env.Get(key + "_SEALED")
		if len(encryptedValue) > 50 {
			fmt.Printf("    Encrypted value: %s...\n", encryptedValue[:50])
		} else {
			fmt.Printf("    Encrypted value: %s\n", encryptedValue)
		}
	}

	// 3. Decrypting and using secrets
	fmt.Println("\n3. Decrypting and using secrets:")

	for key := range secrets {
		decrypted, err := env.GetSealed(key+"_SEALED", "ENV_SEAL_KEY")
		if err != nil {
			log.Fatalf("Failed to unseal %s: %v", key, err)
		}

		fmt.Printf("  ✓ Unsealed %s: %s\n", key, decrypted)

		// Verify it matches the original
		if decrypted == secrets[key] {
			fmt.Printf("    ✓ Matches original value\n")
		} else {
			fmt.Printf("    ✗ Does not match original value\n")
		}
	}

	// 4. Using custom key sources
	fmt.Println("\n4. Using custom key sources:")

	// Set up a different encryption key for API secrets
	apiKey := "api-specific-encryption-key-2024"
	err = env.Set("API_SEAL_KEY", apiKey)
	if err != nil {
		log.Fatalf("Failed to set API key: %v", err)
	}

	// Encrypt with the custom key
	apiSecret := "third-party-api-secret-token"
	err = env.SetSealed("THIRD_PARTY_API_SECRET", apiSecret, "API_SEAL_KEY")
	if err != nil {
		log.Fatalf("Failed to seal API secret: %v", err)
	}

	fmt.Println("  ✓ Sealed API secret with custom key")

	// Decrypt with the custom key
	decryptedApiSecret, err := env.GetSealed("THIRD_PARTY_API_SECRET", "API_SEAL_KEY")
	if err != nil {
		log.Fatalf("Failed to unseal API secret: %v", err)
	}

	fmt.Printf("  ✓ Unsealed API secret: %s\n", decryptedApiSecret)

	// 5. Default key source behavior
	fmt.Println("\n5. Default key source behavior:")

	// When no key source is specified, it defaults to "ENV_SEAL_KEY"
	err = env.SetSealed("DEFAULT_SECRET", "secret-with-default-key", "")
	if err != nil {
		log.Fatalf("Failed to seal with default key: %v", err)
	}

	fmt.Println("  ✓ Sealed secret using default key source (ENV_SEAL_KEY)")

	decryptedDefault, err := env.GetSealed("DEFAULT_SECRET", "")
	if err != nil {
		log.Fatalf("Failed to unseal with default key: %v", err)
	}

	fmt.Printf("  ✓ Unsealed with default key: %s\n", decryptedDefault)

	// 6. Error handling
	fmt.Println("\n6. Error handling:")

	// Try to decrypt with wrong key
	_, err = env.GetSealed("DATABASE_PASSWORD_SEALED", "WRONG_KEY")
	if err != nil {
		fmt.Printf("  ✓ Properly failed with wrong key: %v\n", err)
	}

	// Try to decrypt non-existent secret
	_, err = env.GetSealed("NON_EXISTENT_SECRET", "ENV_SEAL_KEY")
	if err != nil {
		fmt.Printf("  ✓ Properly failed with non-existent secret: %v\n", err)
	}

	// Try with missing encryption key
	_, err = env.GetSealed("DATABASE_PASSWORD_SEALED", "MISSING_KEY")
	if err != nil {
		fmt.Printf("  ✓ Properly failed with missing encryption key: %v\n", err)
	}

	// 7. Practical workflow demonstration
	fmt.Println("\n7. Practical workflow demonstration:")
	fmt.Println("   This is how you'd use sealed secrets in a real application:")
	fmt.Println()
	fmt.Println("   Step 1: Developer encrypts secrets locally")
	fmt.Println("   Step 2: Encrypted secrets are committed to repository")
	fmt.Println("   Step 3: Production environment has the same encryption key")
	fmt.Println("   Step 4: Application automatically decrypts secrets at runtime")
	fmt.Println()
	fmt.Println("   Benefits:")
	fmt.Println("   - Secrets are safe in version control")
	fmt.Println("   - Same .env files work in all environments")
	fmt.Println("   - No complex secret management infrastructure needed")
	fmt.Println("   - Enterprise-grade encryption (AES-GCM)")

	// Clean up
	for key := range secrets {
		_ = env.Unset(key + "_SEALED")
	}
	_ = env.Unset("THIRD_PARTY_API_SECRET")
	_ = env.Unset("DEFAULT_SECRET")
	_ = env.Unset("ENV_SEAL_KEY")
	_ = env.Unset("API_SEAL_KEY")

	fmt.Println("\n=== Sealed Secrets Example Complete ===")
}
