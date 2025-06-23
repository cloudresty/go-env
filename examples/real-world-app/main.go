package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudresty/go-env"
)

// Config represents our application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	API      APIConfig
	Security SecurityConfig
}

type AppConfig struct {
	Name        string
	Version     string
	Environment string
	Debug       bool
	Port        string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
	SSLMode  string
}

type APIConfig struct {
	BaseURL    string
	Timeout    string
	APIKey     string
	RateLimit  int
	RetryCount int
}

type SecurityConfig struct {
	JWTSecret     string
	EncryptionKey string
	SessionSecret string
}

func main() {
	fmt.Println("=== Real-World Application Configuration Example ===")

	// 1. Set up the environment
	fmt.Println("1. Setting up environment:")
	setupEnvironment()

	// 2. Load configuration
	fmt.Println("\n2. Loading application configuration:")
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 3. Display the loaded configuration
	fmt.Println("\n3. Loaded configuration:")
	displayConfig(config)

	// 4. Demonstrate environment-specific loading
	fmt.Println("\n4. Environment-specific configuration:")
	demonstrateEnvironmentSpecificConfig()

	// 5. Show sealed secrets in action
	fmt.Println("\n5. Sealed secrets in production simulation:")
	demonstrateProductionSecrets()

	// 6. Configuration validation
	fmt.Println("\n6. Configuration validation:")
	err = validateConfiguration(config)
	if err != nil {
		fmt.Printf("  ✗ Configuration validation failed: %v\n", err)
	} else {
		fmt.Println("  ✓ Configuration validation passed")
	}

	fmt.Println("\n=== Real-World Example Complete ===")
}

func setupEnvironment() {
	// Create different environment files
	createEnvironmentFiles()

	// Set up encryption key for sealed secrets
	err := env.Set("ENV_SEAL_KEY", "production-grade-encryption-key-2024-secure")
	if err != nil {
		log.Fatalf("Failed to set encryption key: %v", err)
	}
	fmt.Println("  ✓ Set up encryption key for sealed secrets")

	// Set current environment
	err = env.Set("ENVIRONMENT", "development")
	if err != nil {
		log.Fatalf("Failed to set environment: %v", err)
	}
	fmt.Println("  ✓ Set environment to development")
}

func createEnvironmentFiles() {
	// Base configuration (common to all environments)
	baseConfig := `# Base application configuration
APP_NAME=Real World Go App
APP_VERSION=2.1.0

# Database defaults
DB_NAME=myapp
DB_USERNAME=postgres
DB_SSL_MODE=prefer

# API defaults
API_TIMEOUT=30s
API_RATE_LIMIT=1000
API_RETRY_COUNT=3
`

	// Development configuration
	devConfig := `# Development environment overrides
ENVIRONMENT=development
DEBUG=true
PORT=3000

# Development database
DB_HOST=localhost
DB_PORT=5432
DB_PASSWORD=dev_password_123

# Development API
API_BASE_URL=https://api-dev.example.com
API_KEY=dev_api_key_xyz

# Development secrets (in real app, these would be sealed)
JWT_SECRET=dev_jwt_secret_key
ENCRYPTION_KEY=dev_encryption_key
SESSION_SECRET=dev_session_secret
`

	// Production configuration (with sealed secrets)
	prodConfig := `# Production environment overrides
ENVIRONMENT=production
DEBUG=false
PORT=80

# Production database
DB_HOST=prod-db.cluster.amazonaws.com
DB_PORT=5432

# Production API
API_BASE_URL=https://api.example.com
API_RATE_LIMIT=5000

# Production secrets (sealed for security)
# These would be the actual encrypted values in a real deployment
DB_PASSWORD_SEALED=encrypted_db_password_here
API_KEY_SEALED=encrypted_api_key_here
JWT_SECRET_SEALED=encrypted_jwt_secret_here
ENCRYPTION_KEY_SEALED=encrypted_encryption_key_here
SESSION_SECRET_SEALED=encrypted_session_secret_here
`

	files := map[string]string{
		"app.env":         baseConfig,
		"development.env": devConfig,
		"production.env":  prodConfig,
	}

	for filename, content := range files {
		err := writeFile(filename, content)
		if err != nil {
			log.Fatalf("Failed to create %s: %v", filename, err)
		}
	}
	fmt.Println("  ✓ Created environment configuration files")
}

func LoadConfig() (*Config, error) {
	// Load base configuration first
	err := env.Load("app.env")
	if err != nil {
		return nil, fmt.Errorf("failed to load base config: %w", err)
	}

	// Load environment-specific configuration
	environment := env.Get("ENVIRONMENT", "development")
	envFile := environment + ".env"

	err = env.Load(envFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s: %w", envFile, err)
	}

	// Load local overrides (if they exist, ignore errors)
	_ = env.Load(".env.local")

	// Build configuration struct
	config := &Config{
		App:      loadAppConfig(),
		Database: loadDatabaseConfig(),
		API:      loadAPIConfig(),
		Security: loadSecurityConfig(),
	}

	return config, nil
}

func loadAppConfig() AppConfig {
	debug := strings.ToLower(env.Get("DEBUG", "false")) == "true"

	return AppConfig{
		Name:        env.Get("APP_NAME", "Unknown App"),
		Version:     env.Get("APP_VERSION", "0.0.0"),
		Environment: env.Get("ENVIRONMENT", "development"),
		Debug:       debug,
		Port:        env.Get("PORT", "8080"),
	}
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     env.Get("DB_HOST", "localhost"),
		Port:     env.Get("DB_PORT", "5432"),
		Name:     env.Get("DB_NAME", "myapp"),
		Username: env.Get("DB_USERNAME", "postgres"),
		Password: getSecretValue("DB_PASSWORD"),
		SSLMode:  env.Get("DB_SSL_MODE", "prefer"),
	}
}

func loadAPIConfig() APIConfig {
	return APIConfig{
		BaseURL:    env.Get("API_BASE_URL", "https://api.example.com"),
		Timeout:    env.Get("API_TIMEOUT", "30s"),
		APIKey:     getSecretValue("API_KEY"),
		RateLimit:  parseInt(env.Get("API_RATE_LIMIT", "1000")),
		RetryCount: parseInt(env.Get("API_RETRY_COUNT", "3")),
	}
}

func loadSecurityConfig() SecurityConfig {
	return SecurityConfig{
		JWTSecret:     getSecretValue("JWT_SECRET"),
		EncryptionKey: getSecretValue("ENCRYPTION_KEY"),
		SessionSecret: getSecretValue("SESSION_SECRET"),
	}
}

// getSecretValue tries to get a sealed secret first, then falls back to plain text
func getSecretValue(key string) string {
	// Try sealed version first (for production)
	if sealed, err := env.GetSealed(key+"_SEALED", "ENV_SEAL_KEY"); err == nil && sealed != "" {
		return sealed
	}

	// Fall back to plain text (for development)
	value := env.Get(key)
	if value == "" {
		log.Printf("Warning: %s is not set (neither sealed nor plain text)", key)
	}

	return value
}

func parseInt(s string) int {
	// Simple int parsing for demo (in real app, use strconv.Atoi with error handling)
	if s == "" {
		return 0
	}
	// This is simplified - in real code, properly handle conversion errors
	var result int
	_, _ = fmt.Sscanf(s, "%d", &result)
	return result
}

func displayConfig(config *Config) {
	fmt.Printf("  App Configuration:\n")
	fmt.Printf("    Name: %s\n", config.App.Name)
	fmt.Printf("    Version: %s\n", config.App.Version)
	fmt.Printf("    Environment: %s\n", config.App.Environment)
	fmt.Printf("    Debug: %t\n", config.App.Debug)
	fmt.Printf("    Port: %s\n", config.App.Port)

	fmt.Printf("\n  Database Configuration:\n")
	fmt.Printf("    Host: %s\n", config.Database.Host)
	fmt.Printf("    Port: %s\n", config.Database.Port)
	fmt.Printf("    Name: %s\n", config.Database.Name)
	fmt.Printf("    Username: %s\n", config.Database.Username)
	fmt.Printf("    Password: %s\n", maskSecret(config.Database.Password))
	fmt.Printf("    SSL Mode: %s\n", config.Database.SSLMode)

	fmt.Printf("\n  API Configuration:\n")
	fmt.Printf("    Base URL: %s\n", config.API.BaseURL)
	fmt.Printf("    Timeout: %s\n", config.API.Timeout)
	fmt.Printf("    API Key: %s\n", maskSecret(config.API.APIKey))
	fmt.Printf("    Rate Limit: %d\n", config.API.RateLimit)
	fmt.Printf("    Retry Count: %d\n", config.API.RetryCount)

	fmt.Printf("\n  Security Configuration:\n")
	fmt.Printf("    JWT Secret: %s\n", maskSecret(config.Security.JWTSecret))
	fmt.Printf("    Encryption Key: %s\n", maskSecret(config.Security.EncryptionKey))
	fmt.Printf("    Session Secret: %s\n", maskSecret(config.Security.SessionSecret))
}

func maskSecret(secret string) string {
	if secret == "" {
		return "(not set)"
	}
	if len(secret) <= 8 {
		return "***masked***"
	}
	return secret[:4] + "***" + secret[len(secret)-4:]
}

func demonstrateEnvironmentSpecificConfig() {
	// Simulate switching to production
	fmt.Println("  Switching to production environment...")

	_ = env.Set("ENVIRONMENT", "production")

	prodConfig, err := LoadConfig()
	if err != nil {
		log.Printf("  Failed to load production config: %v", err)
		return
	}

	fmt.Printf("  Production vs Development differences:\n")
	fmt.Printf("    Debug: %t (vs true in dev)\n", prodConfig.App.Debug)
	fmt.Printf("    Port: %s (vs 3000 in dev)\n", prodConfig.App.Port)
	fmt.Printf("    DB Host: %s (vs localhost in dev)\n", prodConfig.Database.Host)
	fmt.Printf("    API Rate Limit: %d (vs 1000 in dev)\n", prodConfig.API.RateLimit)
}

func demonstrateProductionSecrets() {
	// Simulate sealing secrets for production
	fmt.Println("  Sealing production secrets...")

	prodSecrets := map[string]string{
		"DB_PASSWORD":    "prod-super-secret-db-password!",
		"API_KEY":        "prod-api-key-abc123xyz789",
		"JWT_SECRET":     "prod-jwt-secret-key-very-secure",
		"ENCRYPTION_KEY": "prod-encryption-key-aes-256",
		"SESSION_SECRET": "prod-session-secret-for-cookies",
	}

	// Seal all production secrets
	for key, value := range prodSecrets {
		err := env.SetSealed(key+"_SEALED", value, "ENV_SEAL_KEY")
		if err != nil {
			log.Printf("  Failed to seal %s: %v", key, err)
			continue
		}
		fmt.Printf("  ✓ Sealed %s\n", key)
	}

	// Now demonstrate retrieving them
	fmt.Println("\n  Retrieving sealed secrets:")
	for key := range prodSecrets {
		retrieved, err := env.GetSealed(key+"_SEALED", "ENV_SEAL_KEY")
		if err != nil {
			log.Printf("  Failed to retrieve %s: %v", key, err)
			continue
		}
		fmt.Printf("  ✓ Retrieved %s: %s\n", key, maskSecret(retrieved))
	}

	// Clean up sealed secrets
	for key := range prodSecrets {
		_ = env.Unset(key + "_SEALED")
	}
}

func validateConfiguration(config *Config) error {
	var errors []string

	// Validate app config
	if config.App.Name == "" {
		errors = append(errors, "app name is required")
	}

	// Validate database config
	if config.Database.Host == "" {
		errors = append(errors, "database host is required")
	}
	if config.Database.Password == "" {
		errors = append(errors, "database password is required")
	}

	// Validate API config
	if config.API.BaseURL == "" {
		errors = append(errors, "API base URL is required")
	}
	if config.API.APIKey == "" {
		errors = append(errors, "API key is required")
	}

	// Validate security config
	if config.Security.JWTSecret == "" {
		errors = append(errors, "JWT secret is required")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errors, ", "))
	}

	return nil
}

func writeFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	_, err = file.WriteString(content)
	return err
}
