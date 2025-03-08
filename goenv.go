package goenv

import (
	"bufio"
	"os"
	"strings"
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

// Lookup retrieves the value of an environment variable.
// It returns the value and a boolean indicating whether the variable exists.
func Lookup(key string) (string, bool) {

	value, exists := os.LookupEnv(key)
	return value, exists

}

// Load loads environment variables from a .env file.
// If filename is empty, it loads from the default ".env" file.
// If no filename is provided, it attempts to load ".env" automatically.
func Load(filename ...string) error {

	var fileToLoad string

	if len(filename) > 0 && filename[0] != "" {
		fileToLoad = filename[0]
	} else {
		fileToLoad = ".env"
	}

	file, err := os.Open(fileToLoad)
	if err != nil {
		if os.IsNotExist(err) && fileToLoad == ".env" && (len(filename) == 0 || filename[0] == "") {
			// If .env doesn't exist and it was the default, don't return an error.
			return nil
		}
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

// MustLoad loads environment variables from a .env file and panics on error.
// If filename is empty, it loads from the default ".env" file.
func MustLoad(filename ...string) {

	if err := Load(filename...); err != nil {
		panic(err)
	}

}
