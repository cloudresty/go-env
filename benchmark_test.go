package env

import (
	"os"
	"reflect"
	"testing"
	"time"
)

// Benchmark structs
type SimpleConfig struct {
	Host    string `env:"HOST" default:"localhost"`
	Port    int    `env:"PORT" default:"8080"`
	Debug   bool   `env:"DEBUG" default:"false"`
	Timeout string `env:"TIMEOUT" default:"30s"`
}

type ComplexConfig struct {
	Host     string        `env:"HOST" default:"localhost"`
	Port     int           `env:"PORT" default:"8080"`
	Debug    bool          `env:"DEBUG" default:"false"`
	Timeout  time.Duration `env:"TIMEOUT" default:"30s"`
	Features []string      `env:"FEATURES"`
	MaxConns int64         `env:"MAX_CONNECTIONS" default:"100"`
	Rate     float64       `env:"RATE" default:"1.5"`
	Buffer   uint          `env:"BUFFER_SIZE" default:"1024"`
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" default:"localhost"`
	Port     int    `env:"DB_PORT" default:"5432"`
	Username string `env:"DB_USER" default:"postgres"`
	Password string `env:"DB_PASS"`
}

func setupBenchmarkEnv() {
	_ = os.Setenv("HOST", "example.com")
	_ = os.Setenv("PORT", "3000")
	_ = os.Setenv("DEBUG", "true")
	_ = os.Setenv("TIMEOUT", "45s")
	_ = os.Setenv("FEATURES", "auth,logging,metrics")
	_ = os.Setenv("MAX_CONNECTIONS", "500")
	_ = os.Setenv("RATE", "2.5")
	_ = os.Setenv("BUFFER_SIZE", "2048")
	_ = os.Setenv("DB_HOST", "db.example.com")
	_ = os.Setenv("DB_PORT", "5432")
	_ = os.Setenv("DB_USER", "admin")
	_ = os.Setenv("DB_PASS", "secret")
}

func cleanupBenchmarkEnv() {
	vars := []string{"HOST", "PORT", "DEBUG", "TIMEOUT", "FEATURES", "MAX_CONNECTIONS",
		"RATE", "BUFFER_SIZE", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASS"}
	for _, v := range vars {
		_ = os.Unsetenv(v)
	}
}

// Basic operations benchmarks
func BenchmarkGet(b *testing.B) {
	_ = os.Setenv("BENCH_KEY", "bench_value")
	defer func() { _ = os.Unsetenv("BENCH_KEY") }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Get("BENCH_KEY", "default")
	}
}

func BenchmarkSet(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Set("BENCH_KEY", "bench_value")
	}
	_ = os.Unsetenv("BENCH_KEY")
}

// Struct binding benchmarks
func BenchmarkBindSimpleStruct(b *testing.B) {
	setupBenchmarkEnv()
	defer cleanupBenchmarkEnv()

	options := DefaultBindingOptions()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var config SimpleConfig
		if err := Bind(&config, options); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBindComplexStruct(b *testing.B) {
	setupBenchmarkEnv()
	defer cleanupBenchmarkEnv()

	options := DefaultBindingOptions()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var config ComplexConfig
		if err := Bind(&config, options); err != nil {
			b.Fatal(err)
		}
	}
}

// Crypto benchmarks
func BenchmarkEncryption(b *testing.B) {
	_ = os.Setenv("ENV_SEAL_KEY", "test-encryption-key-12345")
	defer func() { _ = os.Unsetenv("ENV_SEAL_KEY") }()

	plaintext := "this is a secret message that needs to be encrypted"

	b.Run("Original", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := encrypt(plaintext, "test-encryption-key-12345")
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Optimized", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := encrypt(plaintext, "test-encryption-key-12345")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkDecryption(b *testing.B) {
	_ = os.Setenv("ENV_SEAL_KEY", "test-encryption-key-12345")
	defer func() { _ = os.Unsetenv("ENV_SEAL_KEY") }()

	plaintext := "this is a secret message that needs to be encrypted"
	key := "test-encryption-key-12345"

	// Pre-encrypt for benchmarking
	ciphertext, err := encrypt(plaintext, key)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("Original", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := decrypt(ciphertext, key)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Optimized", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := decrypt(ciphertext, key)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// File operations benchmarks
func BenchmarkFileOperations(b *testing.B) {
	// Create test .env file
	testContent := `# Test .env file
HOST=example.com
PORT=3000
DEBUG=true
DATABASE_URL=postgres://user:pass@localhost:5432/test
API_KEY=secret-key-123
TIMEOUT=30s
FEATURES=auth,logging,metrics
MAX_CONNECTIONS=100
`

	err := os.WriteFile("test_bench.env", []byte(testContent), 0644)
	if err != nil {
		b.Fatal(err)
	}
	defer func() { _ = os.Remove("test_bench.env") }()

	b.Run("LoadOriginal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cleanupBenchmarkEnv()
			if err := Load("test_bench.env"); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Load", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cleanupBenchmarkEnv()
			if err := Load("test_bench.env"); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Memory allocation benchmarks
func BenchmarkMemoryAllocations(b *testing.B) {
	setupBenchmarkEnv()
	defer cleanupBenchmarkEnv()

	options := DefaultBindingOptions()

	b.Run("StructBinding", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var config ComplexConfig
			if err := Bind(&config, options); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("BasicOperations", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = Get("HOST", "localhost")
			_ = Get("PORT", "8080")
			_ = Get("DEBUG", "false")
			_ = Get("TIMEOUT", "30s")
		}
	})
}

// Reflection overhead benchmark
func BenchmarkReflectionOverhead(b *testing.B) {
	setupBenchmarkEnv()
	defer cleanupBenchmarkEnv()

	var config ComplexConfig
	configType := reflect.TypeOf(config)
	configValue := reflect.ValueOf(&config).Elem()

	b.Run("ReflectionCalls", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < configType.NumField(); j++ {
				field := configType.Field(j)
				fieldValue := configValue.Field(j)
				_ = field.Tag.Get("env")
				_ = fieldValue.Kind()
			}
		}
	})
}

// Concurrent access benchmark
func BenchmarkConcurrentAccess(b *testing.B) {
	setupBenchmarkEnv()
	defer cleanupBenchmarkEnv()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Get("HOST", "localhost")
			_ = Get("PORT", "8080")
			_ = Get("DEBUG", "false")
		}
	})
}

// Base64 operations benchmark
func BenchmarkBase64Operations(b *testing.B) {
	testData := "This is some test data that will be base64 encoded"

	b.Run("SetB64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = SetB64("TEST_B64", testData)
		}
	})

	b.Run("GetB64", func(b *testing.B) {
		_ = SetB64("TEST_B64", testData)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := GetB64("TEST_B64")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
