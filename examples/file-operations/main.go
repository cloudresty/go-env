package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudresty/go-env"
)

func main() {
	fmt.Println("=== File Operations Example ===")

	// 1. Creating sample .env files
	fmt.Println("1. Creating sample .env files:")

	// Create a basic .env file
	basicEnvContent := `# Basic application configuration
APP_NAME=File Operations Demo
APP_VERSION=1.0.0
DEBUG=true
PORT=8080

# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp

# API configuration
API_BASE_URL=https://api.example.com
API_TIMEOUT=30s
`

	err := writeFile("basic.env", basicEnvContent)
	if err != nil {
		log.Fatalf("Failed to create basic.env: %v", err)
	}
	fmt.Println("  ✓ Created basic.env")

	// Create an environment-specific file
	productionEnvContent := `# Production overrides
DEBUG=false
PORT=80

# Production database
DB_HOST=prod-db.example.com
DB_PORT=5432

# Production API
API_BASE_URL=https://api.production.com
API_TIMEOUT=60s

# Production-specific settings
LOG_LEVEL=info
CACHE_TTL=3600
`

	err = writeFile("production.env", productionEnvContent)
	if err != nil {
		log.Fatalf("Failed to create production.env: %v", err)
	}
	fmt.Println("  ✓ Created production.env")

	// 2. Loading from files
	fmt.Println("\n2. Loading from .env files:")

	// Load basic configuration
	err = env.Load("basic.env")
	if err != nil {
		log.Fatalf("Failed to load basic.env: %v", err)
	}
	fmt.Println("  ✓ Loaded basic.env")

	// Show loaded values
	fmt.Printf("    APP_NAME: %s\n", env.Get("APP_NAME"))
	fmt.Printf("    DEBUG: %s\n", env.Get("DEBUG"))
	fmt.Printf("    DB_HOST: %s\n", env.Get("DB_HOST"))

	// Load production overrides
	err = env.Load("production.env")
	if err != nil {
		log.Fatalf("Failed to load production.env: %v", err)
	}
	fmt.Println("  ✓ Loaded production.env (overrides applied)")

	// Show overridden values
	fmt.Printf("    DEBUG (after override): %s\n", env.Get("DEBUG"))
	fmt.Printf("    PORT (after override): %s\n", env.Get("PORT"))
	fmt.Printf("    DB_HOST (after override): %s\n", env.Get("DB_HOST"))

	// 3. Loading with MustLoad
	fmt.Println("\n3. Loading with MustLoad:")

	// This would panic if the file doesn't exist or has errors
	env.MustLoad("basic.env")
	fmt.Println("  ✓ Successfully used MustLoad (would panic on error)")

	// 4. Loading non-existent files gracefully
	fmt.Println("\n4. Handling non-existent files:")

	// Try to load a file that doesn't exist
	err = env.Load("non-existent.env")
	if err != nil {
		fmt.Printf("  Expected error loading non-existent file: %v\n", err)
	}

	// Load with default behavior (no error if .env doesn't exist)
	err = env.Load()
	if err != nil {
		fmt.Printf("  Error loading default .env: %v\n", err)
	} else {
		fmt.Println("  ✓ Default .env load completed (no error if file missing)")
	}

	// 5. Saving environment variables to files
	fmt.Println("\n5. Saving environment variables:")

	// Set some additional variables
	_ = env.Set("RUNTIME_CONFIG", "generated_at_runtime")
	_ = env.Set("USER_PREFERENCE", "dark_mode")
	_ = env.Set("SESSION_ID", "abc123xyz")

	// Save specific variables to a file
	keysToSave := []string{
		"APP_NAME",
		"APP_VERSION",
		"DEBUG",
		"PORT",
		"RUNTIME_CONFIG",
		"USER_PREFERENCE",
		"NON_EXISTENT_KEY", // This won't be saved since it doesn't exist
	}

	err = env.Save("exported.env", keysToSave)
	if err != nil {
		log.Fatalf("Failed to save variables: %v", err)
	}
	fmt.Println("  ✓ Saved selected variables to exported.env")

	// Read and display the saved file
	fmt.Println("  Content of exported.env:")
	content, err := os.ReadFile("exported.env")
	if err != nil {
		log.Printf("    Failed to read exported.env: %v", err)
	} else {
		fmt.Printf("    %s", string(content))
	}

	// 6. Handling special characters and quotes
	fmt.Println("\n6. Handling special characters:")

	// Create a file with various quote styles and special characters
	specialCharsContent := `# Testing special characters and quotes
SIMPLE_VALUE=hello
QUOTED_VALUE="hello world"
SINGLE_QUOTED='single quotes'
VALUE_WITH_SPACES=   value with spaces   
EMPTY_VALUE=
VALUE_WITH_EQUALS=key=value=another
VALUE_WITH_HASH=value#with#hash
MULTILINE_JSON={"key": "value", "number": 42}
`

	err = writeFile("special.env", specialCharsContent)
	if err != nil {
		log.Fatalf("Failed to create special.env: %v", err)
	}

	// Clear existing values first
	specialKeys := []string{"SIMPLE_VALUE", "QUOTED_VALUE", "SINGLE_QUOTED", "VALUE_WITH_SPACES", "EMPTY_VALUE", "VALUE_WITH_EQUALS", "VALUE_WITH_HASH", "MULTILINE_JSON"}
	for _, key := range specialKeys {
		_ = env.Unset(key)
	}

	err = env.Load("special.env")
	if err != nil {
		log.Fatalf("Failed to load special.env: %v", err)
	}
	fmt.Println("  ✓ Loaded file with special characters")

	// Show how values are parsed
	fmt.Printf("    SIMPLE_VALUE: '%s'\n", env.Get("SIMPLE_VALUE"))
	fmt.Printf("    QUOTED_VALUE: '%s'\n", env.Get("QUOTED_VALUE"))
	fmt.Printf("    SINGLE_QUOTED: '%s'\n", env.Get("SINGLE_QUOTED"))
	fmt.Printf("    VALUE_WITH_SPACES: '%s'\n", env.Get("VALUE_WITH_SPACES"))
	fmt.Printf("    EMPTY_VALUE: '%s'\n", env.Get("EMPTY_VALUE"))
	fmt.Printf("    VALUE_WITH_EQUALS: '%s'\n", env.Get("VALUE_WITH_EQUALS"))
	fmt.Printf("    MULTILINE_JSON: '%s'\n", env.Get("MULTILINE_JSON"))

	// 7. Error handling with Save
	fmt.Println("\n7. Error handling with Save:")

	// Try to save with no keys (should fail)
	err = env.Save("empty.env", []string{})
	if err != nil {
		fmt.Printf("  ✓ Properly failed when trying to save with no keys: %v\n", err)
	}

	// Clean up files
	filesToClean := []string{"basic.env", "production.env", "exported.env", "special.env"}
	fmt.Println("\n8. Cleaning up files:")
	for _, filename := range filesToClean {
		err := os.Remove(filename)
		if err != nil {
			fmt.Printf("  Warning: Failed to remove %s: %v\n", filename, err)
		} else {
			fmt.Printf("  ✓ Removed %s\n", filename)
		}
	}

	fmt.Println("\n=== File Operations Example Complete ===")
}

// Helper function to write files
func writeFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	_, err = file.WriteString(content)
	return err
}
