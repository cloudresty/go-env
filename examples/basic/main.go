package main

import (
	"fmt"
	"log"

	"github.com/cloudresty/go-env"
)

func main() {
	fmt.Println("=== Basic Environment Operations ===")

	// 1. Setting environment variables
	fmt.Println("1. Setting environment variables:")
	err := env.Set("APP_NAME", "My Go Application")
	if err != nil {
		log.Fatalf("Failed to set APP_NAME: %v", err)
	}

	err = env.Set("PORT", "8080")
	if err != nil {
		log.Fatalf("Failed to set PORT: %v", err)
	}

	err = env.Set("DEBUG", "true")
	if err != nil {
		log.Fatalf("Failed to set DEBUG: %v", err)
	}

	fmt.Println("  ✓ Set APP_NAME, PORT, and DEBUG")

	// 2. Getting environment variables
	fmt.Println("\n2. Getting environment variables:")

	appName := env.Get("APP_NAME")
	fmt.Printf("  APP_NAME: %s\n", appName)

	port := env.Get("PORT")
	fmt.Printf("  PORT: %s\n", port)

	debug := env.Get("DEBUG")
	fmt.Printf("  DEBUG: %s\n", debug)

	// 3. Getting with default values
	fmt.Println("\n3. Getting with default values:")

	dbHost := env.Get("DB_HOST", "localhost")
	fmt.Printf("  DB_HOST (with default): %s\n", dbHost)

	maxConnections := env.Get("MAX_CONNECTIONS", "100")
	fmt.Printf("  MAX_CONNECTIONS (with default): %s\n", maxConnections)

	// 4. Checking if variables exist
	fmt.Println("\n4. Checking if variables exist:")

	if value, exists := env.Lookup("APP_NAME"); exists {
		fmt.Printf("  ✓ APP_NAME exists: %s\n", value)
	} else {
		fmt.Printf("  ✗ APP_NAME does not exist\n")
	}

	if value, exists := env.Lookup("NON_EXISTENT_VAR"); exists {
		fmt.Printf("  ✓ NON_EXISTENT_VAR exists: %s\n", value)
	} else {
		fmt.Printf("  ✗ NON_EXISTENT_VAR does not exist\n")
	}

	// 5. Unsetting variables
	fmt.Println("\n5. Unsetting variables:")

	fmt.Printf("  DEBUG before unset: %s\n", env.Get("DEBUG"))

	err = env.Unset("DEBUG")
	if err != nil {
		log.Printf("Failed to unset DEBUG: %v", err)
	} else {
		fmt.Println("  ✓ Unset DEBUG variable")
	}

	debugAfter := env.Get("DEBUG", "not_set")
	fmt.Printf("  DEBUG after unset (with default): %s\n", debugAfter)

	// Clean up
	_ = env.Unset("APP_NAME")
	_ = env.Unset("PORT")

	fmt.Println("\n=== Basic Operations Complete ===")
}
