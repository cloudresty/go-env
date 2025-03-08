package goenv

import (
	"bufio"
	"os"
	"strings"
)

// Load loads environment variables from a .env file.
func Load(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip invalid lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if os.Getenv(key) == "" { // Only set if not already set
			os.Setenv(key, value)
		}

	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil

}

// Lookup retrieves the value of an environment variable.
// It returns the value and a boolean indicating whether the variable exists.
func Lookup(key string) (string, bool) {

	value, exists := os.LookupEnv(key)
	return value, exists

}

// Get retrieves the value of an environment variable.
// If the variable does not exist, it returns the provided default value.
func Get(key, defaultValue string) string {

	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue

}

// MustLoad loads environment variables from a .env file and panics on error.
func MustLoad(filename string) {

	if err := Load(filename); err != nil {
		panic(err)
	}

}
