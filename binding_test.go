package env

import (
	"os"
	"testing"
	"time"
)

type TestConfig struct {
	StringField   string        `env:"STRING_FIELD,default=default_string"`
	IntField      int           `env:"INT_FIELD,default=42"`
	BoolField     bool          `env:"BOOL_FIELD,default=true"`
	FloatField    float64       `env:"FLOAT_FIELD,default=3.14"`
	DurationField time.Duration `env:"DURATION_FIELD,default=5m"`
	SliceField    []string      `env:"SLICE_FIELD,default=a,b,c"`
	RequiredField string        `env:"REQUIRED_FIELD,required"`
}

type NestedConfig struct {
	Parent ParentConfig `env:"PARENT"`
	Child  ChildConfig  `env:"CHILD"`
}

type ParentConfig struct {
	Name string `env:"PARENT_NAME,default=parent"`
	Age  int    `env:"PARENT_AGE,default=30"`
}

type ChildConfig struct {
	Name string `env:"CHILD_NAME,default=child"`
	Age  int    `env:"CHILD_AGE,default=5"`
}

func TestBind(t *testing.T) {
	// Setup environment variables
	_ = os.Setenv("STRING_FIELD", "test_string")
	_ = os.Setenv("INT_FIELD", "123")
	_ = os.Setenv("BOOL_FIELD", "false")
	_ = os.Setenv("FLOAT_FIELD", "2.71")
	_ = os.Setenv("DURATION_FIELD", "10m30s")
	_ = os.Setenv("SLICE_FIELD", "x,y,z")
	_ = os.Setenv("REQUIRED_FIELD", "required_value")

	defer func() {
		_ = os.Unsetenv("STRING_FIELD")
		_ = os.Unsetenv("INT_FIELD")
		_ = os.Unsetenv("BOOL_FIELD")
		_ = os.Unsetenv("FLOAT_FIELD")
		_ = os.Unsetenv("DURATION_FIELD")
		_ = os.Unsetenv("SLICE_FIELD")
		_ = os.Unsetenv("REQUIRED_FIELD")
	}()

	var config TestConfig
	err := Bind(&config, DefaultBindingOptions())
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	// Test string field
	if config.StringField != "test_string" {
		t.Errorf("Expected StringField to be 'test_string', got '%s'", config.StringField)
	}

	// Test int field
	if config.IntField != 123 {
		t.Errorf("Expected IntField to be 123, got %d", config.IntField)
	}

	// Test bool field
	if config.BoolField != false {
		t.Errorf("Expected BoolField to be false, got %t", config.BoolField)
	}

	// Test float field
	if config.FloatField != 2.71 {
		t.Errorf("Expected FloatField to be 2.71, got %f", config.FloatField)
	}

	// Test duration field
	expected := 10*time.Minute + 30*time.Second
	if config.DurationField != expected {
		t.Errorf("Expected DurationField to be %v, got %v", expected, config.DurationField)
	}

	// Test slice field
	expectedSlice := []string{"x", "y", "z"}
	if len(config.SliceField) != len(expectedSlice) {
		t.Errorf("Expected SliceField length to be %d, got %d", len(expectedSlice), len(config.SliceField))
	}
	for i, v := range expectedSlice {
		if config.SliceField[i] != v {
			t.Errorf("Expected SliceField[%d] to be '%s', got '%s'", i, v, config.SliceField[i])
		}
	}

	// Test required field
	if config.RequiredField != "required_value" {
		t.Errorf("Expected RequiredField to be 'required_value', got '%s'", config.RequiredField)
	}
}

func TestBindWithDefaults(t *testing.T) {
	// Clear environment variables to test defaults
	_ = os.Unsetenv("STRING_FIELD")
	_ = os.Unsetenv("INT_FIELD")
	_ = os.Unsetenv("BOOL_FIELD")
	_ = os.Unsetenv("FLOAT_FIELD")
	_ = os.Unsetenv("DURATION_FIELD")
	_ = os.Unsetenv("SLICE_FIELD")
	_ = os.Setenv("REQUIRED_FIELD", "required_value") // Still need required field

	defer func() { _ = os.Unsetenv("REQUIRED_FIELD") }()

	var config TestConfig
	err := Bind(&config, DefaultBindingOptions())
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	// Test default values
	if config.StringField != "default_string" {
		t.Errorf("Expected StringField default to be 'default_string', got '%s'", config.StringField)
	}

	if config.IntField != 42 {
		t.Errorf("Expected IntField default to be 42, got %d", config.IntField)
	}

	if config.BoolField != true {
		t.Errorf("Expected BoolField default to be true, got %t", config.BoolField)
	}

	if config.FloatField != 3.14 {
		t.Errorf("Expected FloatField default to be 3.14, got %f", config.FloatField)
	}

	expected := 5 * time.Minute
	if config.DurationField != expected {
		t.Errorf("Expected DurationField default to be %v, got %v", expected, config.DurationField)
	}

	expectedSlice := []string{"a", "b", "c"}
	if len(config.SliceField) != len(expectedSlice) {
		t.Errorf("Expected SliceField default length to be %d, got %d", len(expectedSlice), len(config.SliceField))
	}
}

func TestBindWithPrefix(t *testing.T) {
	// Setup environment variables with prefix
	_ = os.Setenv("TEST_STRING_FIELD", "prefixed_string")
	_ = os.Setenv("TEST_INT_FIELD", "999")
	_ = os.Setenv("TEST_REQUIRED_FIELD", "prefixed_required")

	defer func() {
		_ = os.Unsetenv("TEST_STRING_FIELD")
		_ = os.Unsetenv("TEST_INT_FIELD")
		_ = os.Unsetenv("TEST_REQUIRED_FIELD")
	}()

	var config TestConfig
	options := BindingOptions{
		Tag:      "env",
		Prefix:   "TEST_",
		Required: false,
	}

	err := Bind(&config, options)
	if err != nil {
		t.Fatalf("Bind with prefix failed: %v", err)
	}

	if config.StringField != "prefixed_string" {
		t.Errorf("Expected StringField to be 'prefixed_string', got '%s'", config.StringField)
	}

	if config.IntField != 999 {
		t.Errorf("Expected IntField to be 999, got %d", config.IntField)
	}

	if config.RequiredField != "prefixed_required" {
		t.Errorf("Expected RequiredField to be 'prefixed_required', got '%s'", config.RequiredField)
	}
}

func TestBindRequiredFieldError(t *testing.T) {
	// Don't set required field
	_ = os.Unsetenv("REQUIRED_FIELD")

	var config TestConfig
	err := Bind(&config, DefaultBindingOptions())
	if err == nil {
		t.Error("Expected error for missing required field, got nil")
	}

	if !contains(err.Error(), "required environment variable") {
		t.Errorf("Expected error message to contain 'required environment variable', got: %v", err)
	}
}

func TestBindNestedStruct(t *testing.T) {
	// Setup environment variables for nested struct
	_ = os.Setenv("PARENT_NAME", "test_parent")
	_ = os.Setenv("PARENT_AGE", "40")
	_ = os.Setenv("CHILD_NAME", "test_child")
	_ = os.Setenv("CHILD_AGE", "10")

	defer func() {
		_ = os.Unsetenv("PARENT_NAME")
		_ = os.Unsetenv("PARENT_AGE")
		_ = os.Unsetenv("CHILD_NAME")
		_ = os.Unsetenv("CHILD_AGE")
	}()

	var config NestedConfig
	err := Bind(&config, DefaultBindingOptions())
	if err != nil {
		t.Fatalf("Bind nested struct failed: %v", err)
	}

	if config.Parent.Name != "test_parent" {
		t.Errorf("Expected Parent.Name to be 'test_parent', got '%s'", config.Parent.Name)
	}

	if config.Parent.Age != 40 {
		t.Errorf("Expected Parent.Age to be 40, got %d", config.Parent.Age)
	}

	if config.Child.Name != "test_child" {
		t.Errorf("Expected Child.Name to be 'test_child', got '%s'", config.Child.Name)
	}

	if config.Child.Age != 10 {
		t.Errorf("Expected Child.Age to be 10, got %d", config.Child.Age)
	}
}

func TestBindJSONCompatibility(t *testing.T) {
	// Test that BindJSON still works for backward compatibility
	_ = os.Setenv("STRING_FIELD", "json_compat_test")
	_ = os.Setenv("REQUIRED_FIELD", "json_required")

	defer func() {
		_ = os.Unsetenv("STRING_FIELD")
		_ = os.Unsetenv("REQUIRED_FIELD")
	}()

	var config TestConfig
	err := BindJSON("env", &config)
	if err != nil {
		t.Fatalf("BindJSON failed: %v", err)
	}

	if config.StringField != "json_compat_test" {
		t.Errorf("Expected StringField to be 'json_compat_test', got '%s'", config.StringField)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			indexOfSubstring(s, substr) >= 0)))
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
