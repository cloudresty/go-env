package env

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sync"
)

// Optimization: Pool for reusable hash objects
var (
	hashPool = sync.Pool{
		New: func() interface{} {
			return sha256.New()
		},
	}
)

// encrypt provides optimized encryption with object pooling
func encrypt(plaintext, key string) (string, error) {
	// Get hasher from pool
	hasher := hashPool.Get().(interface {
		io.Writer
		Sum([]byte) []byte
		Reset()
	})
	defer func() {
		hasher.Reset()
		hashPool.Put(hasher)
	}()

	// Hash the key
	_, _ = hasher.Write([]byte(key))
	keyBytes := hasher.Sum(nil)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Pre-allocate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Pre-allocate ciphertext with exact capacity
	plaintextBytes := []byte(plaintext)
	ciphertext := make([]byte, 0, len(nonce)+len(plaintextBytes)+gcm.Overhead())
	ciphertext = append(ciphertext, nonce...)
	ciphertext = gcm.Seal(ciphertext, nonce, plaintextBytes, nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt provides optimized decryption with object pooling
func decrypt(ciphertext, key string) (string, error) {
	// Get hasher from pool
	hasher := hashPool.Get().(interface {
		io.Writer
		Sum([]byte) []byte
		Reset()
	})
	defer func() {
		hasher.Reset()
		hashPool.Put(hasher)
	}()

	// Hash the key
	_, _ = hasher.Write([]byte(key))
	keyBytes := hasher.Sum(nil)

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GetSealed retrieves and decrypts a sealed (encrypted) environment variable.
// The encryption key is derived from the provided keySource environment variable.
// If keySource is empty, it defaults to "ENV_SEAL_KEY".
func GetSealed(key string, keySource string, defaultValue ...string) (string, error) {
	if keySource == "" {
		keySource = "ENV_SEAL_KEY"
	}

	// Get the encryption key
	encKey := Get(keySource)
	if encKey == "" {
		return "", fmt.Errorf("encryption key not found in environment variable: %s", keySource)
	}

	// Get the encrypted value
	encryptedValue := Get(key, defaultValue...)
	if encryptedValue == "" {
		return "", nil
	}

	// Decrypt the value
	decrypted, err := decrypt(encryptedValue, encKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt value for key %s: %w", key, err)
	}

	return decrypted, nil
}

// SetSealed encrypts a value and sets it as an environment variable.
// The encryption key is derived from the provided keySource environment variable.
// If keySource is empty, it defaults to "ENV_SEAL_KEY".
func SetSealed(key, value, keySource string) error {
	if keySource == "" {
		keySource = "ENV_SEAL_KEY"
	}

	// Get the encryption key
	encKey := Get(keySource)
	if encKey == "" {
		return fmt.Errorf("encryption key not found in environment variable: %s", keySource)
	}

	// Encrypt the value
	encrypted, err := encrypt(value, encKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt value for key %s: %w", key, err)
	}

	return Set(key, encrypted)
}
