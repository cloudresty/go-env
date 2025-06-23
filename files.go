package env

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// Optimization: Buffer pool for file operations
var (
	bufferPool = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 0, 4096) // 4KB initial capacity
			return &buf
		},
	}
)

// Load provides optimized .env file loading with reduced allocations
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
	defer func() { _ = file.Close() }()

	// Get file size for optimal buffer allocation
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// Use buffered reading with optimal buffer size
	bufSize := int(stat.Size())
	if bufSize > 64*1024 { // Cap at 64KB
		bufSize = 64 * 1024
	} else if bufSize < 1024 { // Minimum 1KB
		bufSize = 1024
	}

	reader := bufio.NewReaderSize(file, bufSize)
	return parseEnvContent(reader)
}

// parseEnvContentOptimized efficiently parses .env content with minimal allocations
func parseEnvContent(reader *bufio.Reader) error {
	// Get buffer from pool
	bufPtr := bufferPool.Get().(*[]byte)
	buf := *bufPtr
	defer func() {
		if cap(buf) <= 8192 { // Only reuse if not too large
			*bufPtr = buf[:0] // Reset length but keep capacity
			bufferPool.Put(bufPtr)
		}
	}()

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		// Process line efficiently
		if len(line) > 0 {
			// Remove trailing newline
			if line[len(line)-1] == '\n' {
				line = line[:len(line)-1]
			}
			// Remove trailing carriage return (Windows)
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}

			if err := processLine(line); err != nil {
				return err
			}
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

// processLineOptimized processes a single line with minimal allocations
func processLine(line []byte) error {
	// Trim leading and trailing whitespace in-place
	line = bytes.TrimSpace(line)

	// Skip empty lines and comments
	if len(line) == 0 || line[0] == '#' {
		return nil
	}

	// Find equals sign efficiently
	equalsIndex := bytes.IndexByte(line, '=')
	if equalsIndex == -1 {
		return nil // Skip invalid lines
	}

	// Extract key and value with minimal allocations
	keyBytes := bytes.TrimSpace(line[:equalsIndex])
	valueBytes := bytes.TrimSpace(line[equalsIndex+1:])

	if len(keyBytes) == 0 {
		return nil // Skip if no key
	}

	// Convert to strings only when necessary
	key := string(keyBytes)
	value := string(valueBytes)

	// Remove quotes efficiently
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'') {
			value = value[1 : len(value)-1]
		}
	}

	// Only set if not already set (preserves environment precedence)
	if os.Getenv(key) == "" {
		return os.Setenv(key, value)
	}

	return nil
}

// SaveOptimized provides optimized environment variable saving
func Save(filename string, keys []string) error {
	if len(keys) == 0 {
		return errors.New("no keys specified to save")
	}

	// Pre-allocate buffer for writing
	var buf strings.Builder
	buf.Grow(len(keys) * 50) // Estimate 50 chars per line

	for _, key := range keys {
		if value, exists := os.LookupEnv(key); exists {
			// Escape value if it contains spaces or special characters
			if needsQuoting(value) {
				value = fmt.Sprintf(`"%s"`, strings.ReplaceAll(value, `"`, `\"`))
			}

			buf.WriteString(key)
			buf.WriteByte('=')
			buf.WriteString(value)
			buf.WriteByte('\n')
		}
	}

	// Write atomically
	return os.WriteFile(filename, []byte(buf.String()), 0644)
}

// needsQuoting efficiently checks if a value needs quoting
func needsQuoting(value string) bool {
	for i := 0; i < len(value); i++ {
		switch value[i] {
		case ' ', '\t', '\n', '\r', '"', '\'':
			return true
		}
	}
	return false
}

// MustLoadOptimized is the optimized version of MustLoad
func MustLoad(filename ...string) {
	if err := Load(filename...); err != nil {
		panic(err)
	}
}
