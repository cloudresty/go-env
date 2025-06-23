package env

import (
	"encoding/base64"
	"fmt"
	"os"
)

// GetB64 retrieves a base64-encoded environment variable and decodes it.
// If the variable does not exist, it returns the provided default value.
// Returns an error if the value is not valid base64.
func GetB64(key string, defaultValue ...string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return "", nil
	}

	if value == "" {
		return "", nil
	}

	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 value for key %s: %w", key, err)
	}

	return string(decoded), nil
}

// SetB64 base64-encodes a value and sets it as an environment variable.
// It returns an error if the operation fails.
func SetB64(key, value string) error {
	encoded := base64.StdEncoding.EncodeToString([]byte(value))
	return Set(key, encoded)
}
