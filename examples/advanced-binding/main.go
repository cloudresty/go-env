package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudresty/go-env"
)

// Config demonstrates advanced struct binding with various data types
type Config struct {
	// Basic types with defaults
	AppName string `env:"APP_NAME,default=MyApp"`
	Port    int    `env:"PORT,default=8080"`
	Debug   bool   `env:"DEBUG,default=false"`
	Host    string `env:"HOST,default=localhost"`

	// Numeric types
	MaxConnections int64   `env:"MAX_CONNECTIONS,default=100"`
	Timeout        float64 `env:"TIMEOUT_SECONDS,default=30.5"`
	BufferSize     uint    `env:"BUFFER_SIZE,default=1024"`

	// Duration support
	RequestTimeout time.Duration `env:"REQUEST_TIMEOUT,default=30s"`
	SessionTimeout time.Duration `env:"SESSION_TIMEOUT,default=24h"`

	// Slices
	AllowedHosts []string `env:"ALLOWED_HOSTS,default=localhost,127.0.0.1"`
	TrustedPorts []int    `env:"TRUSTED_PORTS,default=80,443,8080"`

	// Required fields (will fail if not provided)
	DatabaseURL string `env:"DATABASE_URL,required"`

	// Nested struct
	Redis RedisConfig `env:"REDIS_CONFIG"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST,default=localhost"`
	Port     int    `env:"REDIS_PORT,default=6379"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB,default=0"`
}

func main() {
	fmt.Println("=== Advanced Struct Binding Example ===")

	// Set some environment variables for demonstration
	_ = os.Setenv("APP_NAME", "AdvancedEnvDemo")
	_ = os.Setenv("PORT", "3000")
	_ = os.Setenv("DEBUG", "true")
	_ = os.Setenv("MAX_CONNECTIONS", "500")
	_ = os.Setenv("REQUEST_TIMEOUT", "45s")
	_ = os.Setenv("SESSION_TIMEOUT", "2h30m")
	_ = os.Setenv("ALLOWED_HOSTS", "example.com,api.example.com,localhost")
	_ = os.Setenv("TRUSTED_PORTS", "80,443,3000,8080")
	_ = os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/mydb")
	_ = os.Setenv("REDIS_HOST", "redis.example.com")
	_ = os.Setenv("REDIS_PORT", "6380")
	_ = os.Setenv("REDIS_PASSWORD", "secret123")

	// Example 1: Basic binding with default options
	fmt.Println("1. Basic struct binding:")
	var config Config
	if err := env.Bind(&config, env.DefaultBindingOptions()); err != nil {
		log.Fatal(err)
	}
	printConfig(config)

	// Example 2: Binding with custom options
	fmt.Println("\n2. Binding with prefix:")
	var prefixedConfig Config
	options := env.BindingOptions{
		Tag:      "env",
		Prefix:   "MYAPP_",
		Required: false,
	}

	// Set some prefixed environment variables
	_ = os.Setenv("MYAPP_APP_NAME", "PrefixedApp")
	_ = os.Setenv("MYAPP_PORT", "4000")
	_ = os.Setenv("MYAPP_DATABASE_URL", "postgres://localhost/prefixed_db")

	if err := env.Bind(&prefixedConfig, options); err != nil {
		log.Fatal(err)
	}
	printPrefixedConfig(prefixedConfig)

	// Example 3: Loading from JSON config with env overrides
	fmt.Println("\n3. JSON config with environment overrides:")
	demoJSONConfig()

	// Example 4: Error handling for required fields
	fmt.Println("\n4. Error handling for required fields:")
	demoRequiredFieldsError()

	fmt.Println("\n=== Example completed successfully! ===")
}

func printConfig(config Config) {
	fmt.Printf("  App Name: %s\n", config.AppName)
	fmt.Printf("  Port: %d\n", config.Port)
	fmt.Printf("  Debug: %t\n", config.Debug)
	fmt.Printf("  Host: %s\n", config.Host)
	fmt.Printf("  Max Connections: %d\n", config.MaxConnections)
	fmt.Printf("  Timeout: %.1f seconds\n", config.Timeout)
	fmt.Printf("  Buffer Size: %d\n", config.BufferSize)
	fmt.Printf("  Request Timeout: %v\n", config.RequestTimeout)
	fmt.Printf("  Session Timeout: %v\n", config.SessionTimeout)
	fmt.Printf("  Allowed Hosts: %v\n", config.AllowedHosts)
	fmt.Printf("  Trusted Ports: %v\n", config.TrustedPorts)
	fmt.Printf("  Database URL: %s\n", config.DatabaseURL)
	fmt.Printf("  Redis Host: %s\n", config.Redis.Host)
	fmt.Printf("  Redis Port: %d\n", config.Redis.Port)
	fmt.Printf("  Redis Password: %s\n", config.Redis.Password)
	fmt.Printf("  Redis DB: %d\n", config.Redis.DB)
}

func printPrefixedConfig(config Config) {
	fmt.Printf("  App Name: %s\n", config.AppName)
	fmt.Printf("  Port: %d\n", config.Port)
	fmt.Printf("  Database URL: %s\n", config.DatabaseURL)
	fmt.Printf("  (Other fields use defaults since MYAPP_ prefixed vars not set)\n")
}

func demoJSONConfig() {
	// Create a sample JSON config file
	jsonContent := `{
		"AppName": "JSONApp",
		"Port": 5000,
		"Debug": false,
		"DatabaseURL": "postgres://localhost/json_db"
	}`

	// Create a temporary JSON file
	if err := os.WriteFile("demo_config.json", []byte(jsonContent), 0644); err != nil {
		log.Printf("  Failed to create demo JSON file: %v\n", err)
		return
	}
	defer func() { _ = os.Remove("demo_config.json") }()

	// Override some values with environment variables
	_ = os.Setenv("PORT", "6000")  // This will override the JSON value
	_ = os.Setenv("DEBUG", "true") // This will override the JSON value

	var config Config
	if err := env.LoadFromFile("demo_config.json", &config); err != nil {
		log.Printf("  Failed to load JSON config: %v\n", err)
		return
	}

	fmt.Printf("  JSON App Name: %s (from JSON)\n", config.AppName)
	fmt.Printf("  JSON Port: %d (overridden by ENV)\n", config.Port)
	fmt.Printf("  JSON Debug: %t (overridden by ENV)\n", config.Debug)
	fmt.Printf("  JSON Database URL: %s (from JSON)\n", config.DatabaseURL)
}

func demoRequiredFieldsError() {
	// Temporarily unset the required DATABASE_URL
	originalURL := os.Getenv("DATABASE_URL")
	_ = os.Unsetenv("DATABASE_URL")
	defer func() { _ = os.Setenv("DATABASE_URL", originalURL) }()

	var config Config
	options := env.BindingOptions{
		Tag:      "env",
		Required: true, // This will make all fields required
	}

	if err := env.Bind(&config, options); err != nil {
		fmt.Printf("  Expected error occurred: %v\n", err)
	} else {
		fmt.Printf("  ERROR: Should have failed due to missing required field!\n")
	}
}
