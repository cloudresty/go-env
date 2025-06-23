package env

import (
	"os"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	// Test with existing environment variable
	_ = os.Setenv("TEST_VAR", "test_value")
	defer func() { _ = os.Unsetenv("TEST_VAR") }()

	result := Get("TEST_VAR")
	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}

	// Test with default value
	result = Get("NON_EXISTENT", "default")
	if result != "default" {
		t.Errorf("Expected 'default', got '%s'", result)
	}

	// Test without default value
	result = Get("NON_EXISTENT")
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestSet(t *testing.T) {
	err := Set("TEST_SET_VAR", "set_value")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	result := os.Getenv("TEST_SET_VAR")
	if result != "set_value" {
		t.Errorf("Expected 'set_value', got '%s'", result)
	}

	// Clean up
	_ = os.Unsetenv("TEST_SET_VAR")
}

func TestUnset(t *testing.T) {
	_ = os.Setenv("TEST_UNSET_VAR", "value")

	err := Unset("TEST_UNSET_VAR")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	result := os.Getenv("TEST_UNSET_VAR")
	if result != "" {
		t.Errorf("Expected empty string after unset, got '%s'", result)
	}
}

func TestLookup(t *testing.T) {
	_ = os.Setenv("TEST_LOOKUP_VAR", "lookup_value")
	defer func() { _ = os.Unsetenv("TEST_LOOKUP_VAR") }()

	value, exists := Lookup("TEST_LOOKUP_VAR")
	if !exists {
		t.Error("Expected variable to exist")
	}
	if value != "lookup_value" {
		t.Errorf("Expected 'lookup_value', got '%s'", value)
	}

	_, exists = Lookup("NON_EXISTENT_VAR")
	if exists {
		t.Error("Expected variable to not exist")
	}
}

func TestGetB64(t *testing.T) {
	// Test valid base64
	original := "Hello, World! 🌍"
	err := SetB64("TEST_B64", original)
	if err != nil {
		t.Errorf("Unexpected error setting base64: %v", err)
	}
	defer func() { _ = os.Unsetenv("TEST_B64") }()

	decoded, err := GetB64("TEST_B64")
	if err != nil {
		t.Errorf("Unexpected error getting base64: %v", err)
	}
	if decoded != original {
		t.Errorf("Expected '%s', got '%s'", original, decoded)
	}

	// Test invalid base64
	_ = os.Setenv("INVALID_B64", "not valid base64!")
	defer func() { _ = os.Unsetenv("INVALID_B64") }()

	_, err = GetB64("INVALID_B64")
	if err == nil {
		t.Error("Expected error for invalid base64")
	}

	// Test non-existent variable with default
	result, err := GetB64("NON_EXISTENT", "default")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "default" {
		t.Errorf("Expected 'default', got '%s'", result)
	}
}

func TestSetB64(t *testing.T) {
	original := "Secret data with special chars: @#$%^&*()"

	err := SetB64("TEST_SET_B64", original)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	defer func() { _ = os.Unsetenv("TEST_SET_B64") }()

	// Verify it was encoded properly
	rawValue := os.Getenv("TEST_SET_B64")
	if rawValue == original {
		t.Error("Value should have been base64 encoded")
	}

	// Verify we can decode it back
	decoded, err := GetB64("TEST_SET_B64")
	if err != nil {
		t.Errorf("Unexpected error decoding: %v", err)
	}
	if decoded != original {
		t.Errorf("Expected '%s', got '%s'", original, decoded)
	}
}

func TestSealedSecrets(t *testing.T) {
	// Set up encryption key
	keyEnv := "TEST_SEAL_KEY"
	key := "test-encryption-key-123"
	_ = os.Setenv(keyEnv, key)
	defer func() { _ = os.Unsetenv(keyEnv) }()

	original := "super-secret-password-123!"

	// Test sealing
	err := SetSealed("TEST_SEALED", original, keyEnv)
	if err != nil {
		t.Errorf("Unexpected error sealing: %v", err)
	}
	defer func() { _ = os.Unsetenv("TEST_SEALED") }()

	// Verify it was encrypted (shouldn't match original)
	rawValue := os.Getenv("TEST_SEALED")
	if rawValue == original {
		t.Error("Value should have been encrypted")
	}

	// Test unsealing
	decrypted, err := GetSealed("TEST_SEALED", keyEnv)
	if err != nil {
		t.Errorf("Unexpected error unsealing: %v", err)
	}
	if decrypted != original {
		t.Errorf("Expected '%s', got '%s'", original, decrypted)
	}

	// Test with default key source
	_ = os.Setenv("ENV_SEAL_KEY", key)
	defer func() { _ = os.Unsetenv("ENV_SEAL_KEY") }()

	err = SetSealed("TEST_SEALED_DEFAULT", original, "")
	if err != nil {
		t.Errorf("Unexpected error sealing with default key: %v", err)
	}
	defer func() { _ = os.Unsetenv("TEST_SEALED_DEFAULT") }()

	decrypted, err = GetSealed("TEST_SEALED_DEFAULT", "")
	if err != nil {
		t.Errorf("Unexpected error unsealing with default key: %v", err)
	}
	if decrypted != original {
		t.Errorf("Expected '%s', got '%s'", original, decrypted)
	}
}

func TestSealedSecretsErrors(t *testing.T) {
	// Test missing encryption key
	_, err := GetSealed("SOME_VAR", "NON_EXISTENT_KEY")
	if err == nil {
		t.Error("Expected error for missing encryption key")
	}
	if !strings.Contains(err.Error(), "encryption key not found") {
		t.Errorf("Expected 'encryption key not found' error, got: %v", err)
	}

	// Test sealing with missing key
	err = SetSealed("SOME_VAR", "value", "NON_EXISTENT_KEY")
	if err == nil {
		t.Error("Expected error for missing encryption key")
	}

	// Test unsealing invalid data
	_ = os.Setenv("TEST_KEY", "key")
	_ = os.Setenv("INVALID_SEALED", "not-encrypted-data")
	defer func() { _ = os.Unsetenv("TEST_KEY") }()
	defer func() { _ = os.Unsetenv("INVALID_SEALED") }()

	_, err = GetSealed("INVALID_SEALED", "TEST_KEY")
	if err == nil {
		t.Error("Expected error for invalid encrypted data")
	}
}

func TestLoad(t *testing.T) {
	// Create a temporary .env file
	content := `TEST_LOAD_VAR1=value1
TEST_LOAD_VAR2=value2
# This is a comment
TEST_LOAD_VAR3="quoted value"
TEST_LOAD_VAR4='single quoted'
EMPTY_VAR=
SPACES_VAR = value with spaces
`

	filename := "test.env"
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(filename) }()

	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	_ = file.Close()

	// Clean up any existing values
	vars := []string{"TEST_LOAD_VAR1", "TEST_LOAD_VAR2", "TEST_LOAD_VAR3", "TEST_LOAD_VAR4", "EMPTY_VAR", "SPACES_VAR"}
	for _, v := range vars {
		_ = os.Unsetenv(v)
	}
	defer func() {
		for _, v := range vars {
			_ = os.Unsetenv(v)
		}
	}()

	// Load the file
	err = Load(filename)
	if err != nil {
		t.Errorf("Unexpected error loading file: %v", err)
	}

	// Check loaded values
	tests := []struct {
		key      string
		expected string
	}{
		{"TEST_LOAD_VAR1", "value1"},
		{"TEST_LOAD_VAR2", "value2"},
		{"TEST_LOAD_VAR3", "quoted value"},
		{"TEST_LOAD_VAR4", "single quoted"},
		{"EMPTY_VAR", ""},
		{"SPACES_VAR", "value with spaces"},
	}

	for _, test := range tests {
		actual := os.Getenv(test.key)
		if actual != test.expected {
			t.Errorf("For key %s, expected '%s', got '%s'", test.key, test.expected, actual)
		}
	}
}

func TestSave(t *testing.T) {
	// Set up some test variables
	vars := map[string]string{
		"SAVE_VAR1": "value1",
		"SAVE_VAR2": "value with spaces",
		"SAVE_VAR3": "value \"with quotes\"",
	}

	for k, v := range vars {
		_ = os.Setenv(k, v)
	}
	defer func() {
		for k := range vars {
			_ = os.Unsetenv(k)
		}
	}()

	// Save to file
	filename := "test_save.env"
	keys := []string{"SAVE_VAR1", "SAVE_VAR2", "SAVE_VAR3", "NON_EXISTENT"}

	err := Save(filename, keys)
	if err != nil {
		t.Errorf("Unexpected error saving: %v", err)
	}
	defer func() { _ = os.Remove(filename) }()

	// Read and verify file content
	file, err := os.Open(filename)
	if err != nil {
		t.Fatalf("Failed to open saved file: %v", err)
	}
	defer func() { _ = file.Close() }()

	content := make([]byte, 1024)
	n, _ := file.Read(content)
	fileContent := string(content[:n])

	// Check that saved variables are present
	for key, expectedValue := range vars {
		if !strings.Contains(fileContent, key+"=") {
			t.Errorf("Saved file should contain %s", key)
		}

		// For simple values without special characters
		if key == "SAVE_VAR1" && !strings.Contains(fileContent, expectedValue) {
			t.Errorf("Saved file should contain value %s", expectedValue)
		}
	}

	// Check that non-existent variable is not in file
	if strings.Contains(fileContent, "NON_EXISTENT") {
		t.Error("Saved file should not contain non-existent variable")
	}
}

func TestMustLoad(t *testing.T) {
	// Create a valid file
	filename := "test_must_load.env"
	content := "MUST_LOAD_VAR=success"

	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(filename) }()

	_, _ = file.WriteString(content)
	_ = file.Close()

	// Clean up
	_ = os.Unsetenv("MUST_LOAD_VAR")
	defer func() { _ = os.Unsetenv("MUST_LOAD_VAR") }()

	// This should not panic
	MustLoad(filename)

	// Verify it loaded
	value := os.Getenv("MUST_LOAD_VAR")
	if value != "success" {
		t.Errorf("Expected 'success', got '%s'", value)
	}
}

func TestSaveWithNoKeys(t *testing.T) {
	err := Save("test.env", []string{})
	if err == nil {
		t.Error("Expected error when saving with no keys")
	}
	if !strings.Contains(err.Error(), "no keys specified") {
		t.Errorf("Expected 'no keys specified' error, got: %v", err)
	}
}
