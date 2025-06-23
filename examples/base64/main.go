package main

import (
	"fmt"
	"log"

	"github.com/cloudresty/go-env"
)

func main() {
	fmt.Println("=== Base64 Encoding/Decoding Example ===")

	// 1. Basic base64 operations
	fmt.Println("1. Basic base64 operations:")

	originalData := "Hello, World! This is some data that will be base64 encoded."

	// Store data as base64
	err := env.SetB64("ENCODED_MESSAGE", originalData)
	if err != nil {
		log.Fatalf("Failed to set base64 data: %v", err)
	}

	fmt.Printf("  Original data: %s\n", originalData)

	// Show the raw encoded value
	rawEncoded := env.Get("ENCODED_MESSAGE")
	fmt.Printf("  Raw base64 value: %s\n", rawEncoded)

	// Decode the data
	decoded, err := env.GetB64("ENCODED_MESSAGE")
	if err != nil {
		log.Fatalf("Failed to decode base64 data: %v", err)
	}

	fmt.Printf("  Decoded data: %s\n", decoded)
	fmt.Printf("  ✓ Data matches: %t\n", originalData == decoded)

	// 2. Handling binary-like data
	fmt.Println("\n2. Handling binary-like data:")

	binaryData := "This contains special chars: 🔒 🌟 ♠ ♥ ♦ ♣ \x00\x01\x02"

	err = env.SetB64("BINARY_DATA", binaryData)
	if err != nil {
		log.Fatalf("Failed to set binary data: %v", err)
	}

	fmt.Printf("  Original binary data: %q\n", binaryData)

	decodedBinary, err := env.GetB64("BINARY_DATA")
	if err != nil {
		log.Fatalf("Failed to decode binary data: %v", err)
	}

	fmt.Printf("  Decoded binary data: %q\n", decodedBinary)
	fmt.Printf("  ✓ Binary data matches: %t\n", binaryData == decodedBinary)

	// 3. Using with default values
	fmt.Println("\n3. Using with default values:")

	// Try to get a non-existent base64 variable with a default
	defaultMessage := "This is the default message"
	result, err := env.GetB64("NON_EXISTENT_B64", defaultMessage)
	if err != nil {
		log.Printf("Error getting with default: %v", err)
	} else {
		fmt.Printf("  Non-existent variable with default: %s\n", result)
	}

	// 4. Error handling - invalid base64
	fmt.Println("\n4. Error handling:")

	// Manually set an invalid base64 value
	_ = env.Set("INVALID_B64", "This is not valid base64!")

	_, err = env.GetB64("INVALID_B64")
	if err != nil {
		fmt.Printf("  ✓ Properly caught invalid base64 error: %v\n", err)
	} else {
		fmt.Printf("  ✗ Should have failed for invalid base64\n")
	}

	// 5. Practical example - storing configuration
	fmt.Println("\n5. Practical example - storing JSON configuration:")

	jsonConfig := `{
		"database": {
			"host": "localhost",
			"port": 5432,
			"ssl": true
		},
		"features": ["auth", "logging", "metrics"]
	}`

	err = env.SetB64("APP_CONFIG", jsonConfig)
	if err != nil {
		log.Fatalf("Failed to store JSON config: %v", err)
	}

	fmt.Println("  ✓ Stored JSON configuration as base64")

	retrievedConfig, err := env.GetB64("APP_CONFIG")
	if err != nil {
		log.Fatalf("Failed to retrieve JSON config: %v", err)
	}

	fmt.Println("  ✓ Retrieved JSON configuration:")
	fmt.Printf("  %s\n", retrievedConfig)

	// Clean up
	_ = env.Unset("ENCODED_MESSAGE")
	_ = env.Unset("BINARY_DATA")
	_ = env.Unset("INVALID_B64")
	_ = env.Unset("APP_CONFIG")

	fmt.Println("\n=== Base64 Example Complete ===")
}
