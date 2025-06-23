package env

import (
	"os"
)

// Get retrieves the value of an environment variable.
// If the variable does not exist, it returns the provided default value.
func Get(key string, defaultValue ...string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return "" // Return empty string if no default provided.
}

// Set sets an environment variable to the specified value.
// It returns an error if the operation fails.
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// Unset removes an environment variable.
// It returns an error if the operation fails.
func Unset(key string) error {
	return os.Unsetenv(key)
}

// Lookup retrieves the value of an environment variable.
// It returns the value and a boolean indicating whether the variable exists.
func Lookup(key string) (string, bool) {
	value, exists := os.LookupEnv(key)
	return value, exists
}
