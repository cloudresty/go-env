# Advanced Usage Guide

This guide demonstrates advanced usage patterns and best practices for the Go Env package.

**💡 Tip:** Before diving into advanced patterns, explore our [focused examples](./examples/) to understand the basics.

**🚀 Performance Note:** Go Env includes significant optimizations (40% faster, 42% less memory) that work automatically. See [Performance Considerations](#performance-considerations) for details.

&nbsp;

## Table of Contents

1. [Struct Tag Format](#struct-tag-format)
2. [Sealed Secrets Workflow](#sealed-secrets-workflow)
3. [Environment-Specific Configurations](#environment-specific-configurations)
4. [Security Best Practices](#security-best-practices)
5. [Error Handling](#error-handling)
6. [Performance Considerations](#performance-considerations)

&nbsp;

## Struct Tag Format

Go Env uses a clean, single tag format for defining environment variable bindings in your structs:

&nbsp;

### Combined Tag Format

All configuration is specified in the `env` tag using a comma-separated format:

```go
type Config struct {
    Host     string        `env:"HOST,default=localhost"`
    Port     int           `env:"PORT,default=8080"`
    Debug    bool          `env:"DEBUG,default=false"`
    Timeout  time.Duration `env:"TIMEOUT,default=30s"`
    APIKey   string        `env:"API_KEY,required"`
}
```

&nbsp;

### Tag Format Features

**Supported options:**

- **Environment variable name**: First parameter (required)
- **Default values**: `default=value` - used when environment variable is not set
- **Required fields**: `required` - causes binding to fail if environment variable is missing
- **All Go types**: strings, numbers, booleans, time.Duration, slices
- **Nested structs and embedded fields**

&nbsp;

### Complex Default Values

The tag format handles complex default values gracefully:

```go
type AppConfig struct {
    // Default values with commas
    Tags         string `env:"TAGS,default=web,api,backend"`

    // Default values with equals signs
    Config       string `env:"CONFIG,default=key=value,debug=true"`

    // Multiple options with default
    DatabaseURL  string `env:"DATABASE_URL,required,default=postgres://localhost/myapp"`

    // Empty defaults
    OptionalPath string `env:"OPTIONAL_PATH,default="`
}
```

&nbsp;

## Sealed Secrets Workflow

&nbsp;

### Development Workflow

```go
package main

import (
    "log"
    "github.com/cloudresty/go-env"
)

func main() {
    // 1. Load base configuration
    env.MustLoad(".env.base")

    // 2. Set up encryption key (in real app, this would be from secure source)
    if env.Get("ENV_SEAL_KEY") == "" {
        log.Fatal("ENV_SEAL_KEY must be set for sealed secrets")
    }

    // 3. Use sealed secrets in your application
    dbPassword, err := env.GetSealed("DB_PASSWORD_SEALED", "ENV_SEAL_KEY")
    if err != nil {
        log.Fatalf("Failed to unseal database password: %v", err)
    }

    // 4. Use the unsealed value
    connectToDatabase(dbPassword)
}
```

&nbsp;

### Creating Sealed Secrets

Create a utility script to seal your secrets:

```go
// tools/seal-secrets/main.go
package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/cloudresty/go-env"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Usage: go run main.go <env-file>")
    }

    filename := os.Args[1]

    // Load existing environment file
    err := env.Load(filename)
    if err != nil {
        log.Fatalf("Failed to load %s: %v", filename, err)
    }

    // Get encryption key
    encKey := env.Get("ENV_SEAL_KEY")
    if encKey == "" {
        log.Fatal("ENV_SEAL_KEY must be set")
    }

    // Read file and process secrets
    file, err := os.Open(filename)
    if err != nil {
        log.Fatalf("Failed to open file: %v", err)
    }
    defer func() { _ = file.Close() }()

    var lines []string
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()

        // Check if this line should be sealed (ends with _SECRET=)
        if strings.Contains(line, "_SECRET=") || strings.Contains(line, "_PASSWORD=") {
            parts := strings.SplitN(line, "=", 2)
            if len(parts) == 2 {
                key := parts[0]
                value := parts[1]

                // Seal the value
                err := env.SetSealed(key+"_SEALED", value, "ENV_SEAL_KEY")
                if err != nil {
                    log.Printf("Failed to seal %s: %v", key, err)
                    lines = append(lines, line)
                } else {
                    sealedValue := env.Get(key + "_SEALED")
                    lines = append(lines, fmt.Sprintf("# Original: %s", line))
                    lines = append(lines, fmt.Sprintf("%s_SEALED=%s", key, sealedValue))
                }
                continue
            }
        }

        lines = append(lines, line)
    }

    // Write back to file
    outFile, err := os.Create(filename + ".sealed")
    if err != nil {
        log.Fatalf("Failed to create output file: %v", err)
    }
    defer func() { _ = outFile.Close() }()

    for _, line := range lines {
        fmt.Fprintln(outFile, line)
    }

    fmt.Printf("Sealed secrets written to %s.sealed\n", filename)
}
```

&nbsp;

## Environment-Specific Configurations

&nbsp;

### Loading Multiple Environment Files

```go
package config

import (
    "log"
    "os"
    "github.com/cloudresty/go-env"
)

type Config struct {
    Database DatabaseConfig
    API      APIConfig
    Security SecurityConfig
}

type DatabaseConfig struct {
    Host     string
    Port     string
    Name     string
    Username string
    Password string
}

type APIConfig struct {
    BaseURL string
    Timeout string
    APIKey  string
}

type SecurityConfig struct {
    JWTSecret     string
    EncryptionKey string
}

func LoadConfig() *Config {
    // Load base configuration first
    env.MustLoad(".env")

    // Load environment-specific overrides
    environment := env.Get("ENVIRONMENT", "development")

    switch environment {
    case "development":
        env.Load(".env.development")
    case "staging":
        env.Load(".env.staging")
    case "production":
        env.Load(".env.production")
    }

    // Load local overrides (should not be committed)
    env.Load(".env.local")

    return &Config{
        Database: loadDatabaseConfig(),
        API:      loadAPIConfig(),
        Security: loadSecurityConfig(),
    }
}

func loadDatabaseConfig() DatabaseConfig {
    return DatabaseConfig{
        Host:     env.Get("DB_HOST", "localhost"),
        Port:     env.Get("DB_PORT", "5432"),
        Name:     env.Get("DB_NAME", "myapp"),
        Username: env.Get("DB_USERNAME", "postgres"),
        Password: getSecretValue("DB_PASSWORD"),
    }
}

func loadAPIConfig() APIConfig {
    return APIConfig{
        BaseURL: env.Get("API_BASE_URL", "https://api.example.com"),
        Timeout: env.Get("API_TIMEOUT", "30s"),
        APIKey:  getSecretValue("API_KEY"),
    }
}

func loadSecurityConfig() SecurityConfig {
    return SecurityConfig{
        JWTSecret:     getSecretValue("JWT_SECRET"),
        EncryptionKey: getSecretValue("ENCRYPTION_KEY"),
    }
}

// getSecretValue tries to get a sealed secret first, then falls back to plain text
func getSecretValue(key string) string {
    // Try sealed version first
    if sealed, err := env.GetSealed(key+"_SEALED", "ENV_SEAL_KEY"); err == nil && sealed != "" {
        return sealed
    }

    // Fall back to plain text
    value := env.Get(key)
    if value == "" {
        log.Printf("Warning: %s is not set", key)
    }

    return value
}
```

&nbsp;

## Security Best Practices

&nbsp;

### 1. Key Management

```go
// Good: Load encryption key from secure source
func getEncryptionKey() string {
    // In production, load from AWS SSM, HashiCorp Vault, etc.
    if key := env.Get("AWS_SSM_ENCRYPTION_KEY_PATH"); key != "" {
        return loadFromSSM(key)
    }

    // For development
    return env.Get("ENV_SEAL_KEY")
}

// Good: Validate key strength
func validateEncryptionKey(key string) error {
    if len(key) < 32 {
        return fmt.Errorf("encryption key must be at least 32 characters")
    }
    return nil
}
```

&nbsp;

### 2. Secret Rotation

```go
// Utility for rotating sealed secrets
func rotateSeals(oldKeyEnv, newKeyEnv string) error {
    // Get both keys
    oldKey := env.Get(oldKeyEnv)
    newKey := env.Get(newKeyEnv)

    if oldKey == "" || newKey == "" {
        return fmt.Errorf("both old and new keys must be provided")
    }

    // Find all sealed variables
    sealedVars := findSealedVariables()

    for _, varName := range sealedVars {
        // Unseal with old key
        value, err := env.GetSealed(varName, oldKeyEnv)
        if err != nil {
            log.Printf("Failed to unseal %s with old key: %v", varName, err)
            continue
        }

        // Re-seal with new key
        err = env.SetSealed(varName, value, newKeyEnv)
        if err != nil {
            log.Printf("Failed to re-seal %s with new key: %v", varName, err)
            continue
        }

        log.Printf("Successfully rotated %s", varName)
    }

    return nil
}
```

&nbsp;

## Error Handling

&nbsp;

### Graceful Degradation

```go
func getConfigWithFallback() *Config {
    cfg := &Config{}

    // Critical configuration - must succeed
    cfg.Database.Host = env.Get("DB_HOST")
    if cfg.Database.Host == "" {
        log.Fatal("DB_HOST is required")
    }

    // Secret configuration - try sealed first, then plain
    if password, err := env.GetSealed("DB_PASSWORD_SEALED", "ENV_SEAL_KEY"); err != nil {
        log.Printf("Warning: Failed to unseal DB password: %v", err)
        cfg.Database.Password = env.Get("DB_PASSWORD")
        if cfg.Database.Password == "" {
            log.Fatal("DB_PASSWORD or DB_PASSWORD_SEALED is required")
        }
    } else {
        cfg.Database.Password = password
    }

    // Optional configuration with defaults
    cfg.API.Timeout = env.Get("API_TIMEOUT", "30s")

    // Base64 configuration with error handling
    if cert, err := env.GetB64("TLS_CERT_B64"); err != nil {
        log.Printf("Warning: Failed to decode TLS cert: %v", err)
        cfg.TLS.CertPath = env.Get("TLS_CERT_PATH", "./cert.pem")
    } else if cert != "" {
        cfg.TLS.CertData = cert
    }

    return cfg
}
```

&nbsp;

### Validation and Health Checks

```go
func validateConfiguration(cfg *Config) error {
    var errors []string

    // Validate database config
    if cfg.Database.Host == "" {
        errors = append(errors, "database host is required")
    }

    if cfg.Database.Password == "" {
        errors = append(errors, "database password is required")
    }

    // Validate API config
    if cfg.API.APIKey == "" {
        errors = append(errors, "API key is required")
    }

    // Test sealed secrets
    if testSecret, err := env.GetSealed("HEALTH_CHECK_SECRET", "ENV_SEAL_KEY"); err != nil {
        errors = append(errors, fmt.Sprintf("sealed secrets not working: %v", err))
    } else if testSecret == "" {
        errors = append(errors, "health check secret is empty")
    }

    if len(errors) > 0 {
        return fmt.Errorf("configuration validation failed: %s", strings.Join(errors, ", "))
    }

    return nil
}
```

&nbsp;

## Performance Considerations

&nbsp;

### Built-in Optimizations

Go Env includes several performance optimizations that are automatically enabled:

&nbsp;

#### **Smart Caching System**

```go
// The library automatically caches expensive operations
func demonstrateCaching() {
    // First call: parses struct reflection (slower)
    var config Config
    env.Bind(&config, env.DefaultBindingOptions())

    // Subsequent calls: uses cached reflection (80%+ faster)
    var config2 Config
    env.Bind(&config2, env.DefaultBindingOptions()) // Much faster

    // View cache performance
    stats := env.GetCacheStats()
    fmt.Printf("Type cache hit rate: %.1f%%\n", stats["type_cache_hit_rate"])
    fmt.Printf("Tag cache hit rate: %.1f%%\n", stats["tag_cache_hit_rate"])
}
```

&nbsp;

#### **Memory Pool Optimization**

```go
// The library uses memory pools for efficient slice parsing
type Config struct {
    SmallList []string `env:"SMALL_FEATURES"`  // Uses small pool (efficient)
    LargeList []string `env:"LARGE_FEATURES"`  // Uses large pool (efficient)
}

// Memory pools are managed automatically - no configuration needed
```

&nbsp;

#### **Performance Monitoring**

```go
// Monitor performance in production
func monitorPerformance() {
    // Get detailed cache statistics
    stats := env.GetCacheStats()

    fmt.Printf("Performance Metrics:\n")
    fmt.Printf("- Type cache hits: %d\n", stats["type_cache_hits"])
    fmt.Printf("- Type cache misses: %d\n", stats["type_cache_misses"])
    fmt.Printf("- Tag cache hits: %d\n", stats["tag_cache_hits"])
    fmt.Printf("- Tag cache misses: %d\n", stats["tag_cache_misses"])

    // Reset counters for next measurement period
    env.ResetCacheCounters()
}
```

&nbsp;

### Performance Best Practices

&nbsp;

#### **1. Reuse Configuration Structs**

```go
// Good: Reuse the same struct type for maximum cache benefit
type AppConfig struct {
    Database DatabaseConfig
    API      APIConfig
}

func loadConfig() *AppConfig {
    var config AppConfig
    env.Bind(&config, env.DefaultBindingOptions()) // Cached after first use
    return &config
}

// Avoid: Creating many different struct types reduces cache efficiency
```

&nbsp;

#### **2. Batch Operations for Sealed Secrets**

```go
// Efficient: Load multiple secrets in one operation
func loadAllSecrets(keySource string, secretKeys []string) (map[string]string, error) {
    secrets := make(map[string]string)
    var errors []string

    for _, key := range secretKeys {
        if value, err := env.GetSealed(key, keySource); err != nil {
            errors = append(errors, fmt.Sprintf("%s: %v", key, err))
        } else {
            secrets[key] = value
        }
    }

    if len(errors) > 0 {
        return secrets, fmt.Errorf("failed to load some secrets: %s", strings.Join(errors, ", "))
    }

    return secrets, nil
}
```

&nbsp;

#### **3. Cache Management**

```go
// Clear caches if memory usage becomes a concern (rare)
func manageMemory() {
    // Clear type information cache
    env.ClearTypeCache()

    // Clear all caches
    env.ClearAllCaches()

    // Note: Only do this if you have memory constraints
    // The caches are typically small and very beneficial
}
```

&nbsp;

### Caching Unsealed Values

```go
type SecretCache struct {
    cache map[string]string
    mutex sync.RWMutex
}

func NewSecretCache() *SecretCache {
    return &SecretCache{
        cache: make(map[string]string),
    }
}

func (sc *SecretCache) GetSealed(key, keySource string) (string, error) {
    sc.mutex.RLock()
    if value, exists := sc.cache[key]; exists {
        sc.mutex.RUnlock()
        return value, nil
    }
    sc.mutex.RUnlock()

    // Unseal the value
    value, err := env.GetSealed(key, keySource)
    if err != nil {
        return "", err
    }

    // Cache the result
    sc.mutex.Lock()
    sc.cache[key] = value
    sc.mutex.Unlock()

    return value, nil
}

func (sc *SecretCache) Clear() {
    sc.mutex.Lock()
    defer sc.mutex.Unlock()

    // Clear sensitive data from memory
    for k := range sc.cache {
        delete(sc.cache, k)
    }
}
```

&nbsp;

### Batch Operations

```go
// Load multiple sealed secrets efficiently
func loadAllSecrets(keySource string, secretKeys []string) (map[string]string, error) {
    secrets := make(map[string]string)
    var errors []string

    for _, key := range secretKeys {
        if value, err := env.GetSealed(key, keySource); err != nil {
            errors = append(errors, fmt.Sprintf("%s: %v", key, err))
        } else {
            secrets[key] = value
        }
    }

    if len(errors) > 0 {
        return secrets, fmt.Errorf("failed to load some secrets: %s", strings.Join(errors, ", "))
    }

    return secrets, nil
}
```

This guide covers advanced patterns for using the Go Env package in production applications with a focus on security, reliability, and performance.

**🔧 Latest Update:** All performance optimizations have been integrated into the main API. You get 40% better performance automatically without any code changes.
