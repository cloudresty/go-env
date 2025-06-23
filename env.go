// Package env provides a powerful, modern environment management library for Go.
//
// Features:
//
//   - Basic Operations (basic.go): Get, Set, Unset, Lookup environment variables
//   - Base64 Support (base64.go): Automatic encoding/decoding of base64 values
//   - Sealed Secrets (sealed.go): AES-GCM encryption/decryption for sensitive data
//   - File Operations (files.go): Load from and save to .env files
//   - Struct Binding (binding.go): Bind environment variables to struct fields with type conversion
//
// Example usage:
//
//	// Basic operations
//	value := env.Get("DATABASE_URL", "default_value")
//	env.Set("API_KEY", "secret123")
//
//	// Base64 operations
//	encoded, _ := env.GetB64("CONFIG_DATA")
//	env.SetB64("BINARY_DATA", "binary content")
//
//	// Sealed secrets
//	secret, _ := env.GetSealed("ENCRYPTED_PASSWORD", "MASTER_KEY")
//	env.SetSealed("CREDIT_CARD", "4111-1111-1111-1111", "MASTER_KEY")
//
//	// File operations
//	env.Load(".env")
//	env.Save("backup.env", []string{"API_KEY", "DATABASE_URL"})
//
//	// Struct binding
//	type Config struct {
//		Port     int    `env:"PORT"`
//		Host     string `env:"HOST"`
//		Debug    bool   `env:"DEBUG"`
//		Features []string `env:"FEATURES"`
//	}
//	var config Config
//	env.BindJSON("env", &config)
//
// Migration from other packages:
//
//	// From os package
//	value := os.Getenv("KEY")           // becomes: env.Get("KEY")
//	os.Setenv("KEY", "value")           // becomes: env.Set("KEY", "value")
//
//	// From godotenv
//	godotenv.Load()                     // becomes: env.Load()
//
//	// From caarlos0/env
//	env.Parse(&config)                  // becomes: env.BindJSON("env", &config)
//
//	// From viper
//	viper.GetString("key")              // becomes: env.Get("KEY")
//	viper.ReadInConfig()                // becomes: env.Load("config.env")
package env
