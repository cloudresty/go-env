package env

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ConfigFormat represents supported configuration file formats
type ConfigFormat int

const (
	FormatJSON ConfigFormat = iota
	FormatYAML
	FormatTOML
	FormatINI
	FormatEnv
)

// String returns the string representation of the format
func (f ConfigFormat) String() string {
	switch f {
	case FormatJSON:
		return "json"
	case FormatYAML:
		return "yaml"
	case FormatTOML:
		return "toml"
	case FormatINI:
		return "ini"
	case FormatEnv:
		return "env"
	default:
		return "unknown"
	}
}

// LoadFromFile loads configuration from a file in various formats
// and optionally binds it to a struct using environment variable overrides.
func LoadFromFile(filename string, out interface{}) error {
	format, err := detectFormat(filename)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", filename, err)
	}

	switch format {
	case FormatJSON:
		return loadJSONConfig(data, out)
	case FormatEnv:
		// For .env files, just load them and then bind
		if err := Load(filename); err != nil {
			return err
		}
		return Bind(out, DefaultBindingOptions())
	default:
		return fmt.Errorf("format %s not yet implemented (coming soon!)", format)
	}
}

// detectFormat detects the configuration format based on file extension
func detectFormat(filename string) (ConfigFormat, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		return FormatJSON, nil
	case ".yaml", ".yml":
		return FormatYAML, nil
	case ".toml":
		return FormatTOML, nil
	case ".ini":
		return FormatINI, nil
	case ".env":
		return FormatEnv, nil
	default:
		return FormatEnv, nil // Default to env format
	}
}

// loadJSONConfig loads JSON configuration and allows environment variable overrides
func loadJSONConfig(data []byte, out interface{}) error {
	// First, load the JSON data
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("failed to parse JSON config: %w", err)
	}

	// Then, allow environment variables to override the loaded values
	// We use a simpler approach here - just bind env vars without requiring all fields
	options := DefaultBindingOptions()
	options.Required = false
	return Bind(out, options)
}

// SaveToFile saves configuration to a file in the specified format
func SaveToFile(filename string, data interface{}) error {
	format, err := detectFormat(filename)
	if err != nil {
		return err
	}

	switch format {
	case FormatJSON:
		return saveJSONConfig(filename, data)
	default:
		return fmt.Errorf("saving format %s not yet implemented", format)
	}
}

// saveJSONConfig saves data as JSON to a file
func saveJSONConfig(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON config to %s: %w", filename, err)
	}

	return nil
}
