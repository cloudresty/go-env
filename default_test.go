package env

import (
	"os"
	"testing"
	"time"
)

// Test struct for default value functionality
type DefaultTestConfig struct {
	StringField   string        `env:"STRING_FIELD,default=default_string"`
	IntField      int           `env:"INT_FIELD,default=42"`
	BoolField     bool          `env:"BOOL_FIELD,default=true"`
	FloatField    float64       `env:"FLOAT_FIELD,default=3.14"`
	DurationField time.Duration `env:"DURATION_FIELD,default=5m"`
	SliceField    []string      `env:"SLICE_FIELD,default=a,b,c"`
	EmptyDefault  string        `env:"EMPTY_DEFAULT,default="`
	NoDefault     string        `env:"NO_DEFAULT"`
}

func TestDefaultValues(t *testing.T) {
	// Clean up any existing environment variables
	vars := []string{"STRING_FIELD", "INT_FIELD", "BOOL_FIELD", "FLOAT_FIELD", "DURATION_FIELD", "SLICE_FIELD", "EMPTY_DEFAULT", "NO_DEFAULT"}
	for _, v := range vars {
		_ = os.Unsetenv(v)
	}

	var config DefaultTestConfig
	err := Bind(&config, DefaultBindingOptions())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Test default values are applied
	if config.StringField != "default_string" {
		t.Errorf("Expected StringField to be 'default_string', got '%s'", config.StringField)
	}

	if config.IntField != 42 {
		t.Errorf("Expected IntField to be 42, got %d", config.IntField)
	}

	if config.BoolField != true {
		t.Errorf("Expected BoolField to be true, got %v", config.BoolField)
	}

	if config.FloatField != 3.14 {
		t.Errorf("Expected FloatField to be 3.14, got %f", config.FloatField)
	}

	if config.DurationField != 5*time.Minute {
		t.Errorf("Expected DurationField to be 5m, got %v", config.DurationField)
	}

	if len(config.SliceField) != 3 || config.SliceField[0] != "a" || config.SliceField[1] != "b" || config.SliceField[2] != "c" {
		t.Errorf("Expected SliceField to be [a b c], got %v", config.SliceField)
	}

	if config.EmptyDefault != "" {
		t.Errorf("Expected EmptyDefault to be empty string, got '%s'", config.EmptyDefault)
	}

	if config.NoDefault != "" {
		t.Errorf("Expected NoDefault to be empty (zero value), got '%s'", config.NoDefault)
	}
}

func TestDefaultValuesWithEnvironmentOverrides(t *testing.T) {
	// Set some environment variables to override defaults
	_ = os.Setenv("STRING_FIELD", "env_override")
	_ = os.Setenv("INT_FIELD", "100")
	defer func() {
		_ = os.Unsetenv("STRING_FIELD")
		_ = os.Unsetenv("INT_FIELD")
		_ = os.Unsetenv("BOOL_FIELD")
	}()

	var config DefaultTestConfig
	err := Bind(&config, DefaultBindingOptions())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Environment variables should override defaults
	if config.StringField != "env_override" {
		t.Errorf("Expected StringField to be 'env_override', got '%s'", config.StringField)
	}

	if config.IntField != 100 {
		t.Errorf("Expected IntField to be 100, got %d", config.IntField)
	}

	// Non-overridden fields should use defaults
	if config.BoolField != true {
		t.Errorf("Expected BoolField to be true (default), got %v", config.BoolField)
	}
}
