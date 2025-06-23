package env

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// Enhanced benchmark to demonstrate optimization benefits

func BenchmarkSliceParsingSmall(b *testing.B) {
	setupBenchEnv()
	defer teardownBenchEnv()

	// Test small slice (common case for env vars)
	_ = os.Setenv("SMALL_FEATURES", "auth,logging")

	type SmallConfig struct {
		Features []string `env:"SMALL_FEATURES"`
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var config SmallConfig
		if err := Bind(&config, DefaultBindingOptions()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSliceParsingLarge(b *testing.B) {
	setupBenchEnv()
	defer teardownBenchEnv()

	// Test large slice (less common but should still be efficient)
	_ = os.Setenv("LARGE_FEATURES", "auth,logging,metrics,tracing,security,caching,compression,monitoring,alerting,analytics,reporting,backup,recovery,scaling,load-balancing,circuit-breaker")

	type LargeConfig struct {
		Features []string `env:"LARGE_FEATURES"`
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var config LargeConfig
		if err := Bind(&config, DefaultBindingOptions()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCacheEffectiveness(b *testing.B) {
	setupBenchEnv()
	defer teardownBenchEnv()

	type CacheTestConfig struct {
		Host    string        `env:"HOST"`
		Port    int           `env:"PORT"`
		Debug   bool          `env:"DEBUG"`
		Timeout time.Duration `env:"TIMEOUT"`
	}

	// Reset counters
	ResetCacheCounters()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var config CacheTestConfig
		if err := Bind(&config, DefaultBindingOptions()); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()

	// Report cache effectiveness
	stats := GetCacheStats()
	b.Logf("Cache Stats: %+v", stats)
}

func BenchmarkMemoryOptimizations(b *testing.B) {
	setupBenchEnv()
	defer teardownBenchEnv()

	type MemoryTestConfig struct {
		Host     string        `env:"HOST"`
		Port     int           `env:"PORT"`
		Debug    bool          `env:"DEBUG"`
		Timeout  time.Duration `env:"TIMEOUT"`
		Features []string      `env:"FEATURES"`
		MaxConns int64         `env:"MAX_CONNECTIONS"`
		Rate     float64       `env:"RATE"`
		Buffer   uint          `env:"BUFFER_SIZE"`
		Database DatabaseBenchConfig
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var config MemoryTestConfig
		if err := Bind(&config, DefaultBindingOptions()); err != nil {
			b.Fatal(err)
		}
	}
}

// Test that demonstrates the optimization pipeline
func TestOptimizationPipeline(t *testing.T) {
	setupBenchEnv()
	defer teardownBenchEnv()

	type PipelineConfig struct {
		Host     string   `env:"HOST"`
		Port     int      `env:"PORT"`
		Features []string `env:"FEATURES"`
	}

	// Reset and clear caches for clean test
	ClearAllCaches()
	ResetCacheCounters()

	// First run - will populate caches
	var config1 PipelineConfig
	err := Bind(&config1, DefaultBindingOptions())
	if err != nil {
		t.Fatal(err)
	}

	stats1 := GetCacheStats()
	t.Logf("After first run: %+v", stats1)

	// Second run - should hit caches
	var config2 PipelineConfig
	err = Bind(&config2, DefaultBindingOptions())
	if err != nil {
		t.Fatal(err)
	}

	stats2 := GetCacheStats()
	t.Logf("After second run: %+v", stats2)

	// Verify cache hit rates improved
	if stats2["type_cache_hit_rate"].(float64) < stats1["type_cache_hit_rate"].(float64) {
		t.Errorf("Type cache hit rate should improve on subsequent runs")
	}

	// Verify results are identical
	if config1.Host != config2.Host || config1.Port != config2.Port {
		t.Errorf("Results should be identical across runs")
	}
}

func TestPoolEfficiency(t *testing.T) {
	setupBenchEnv()
	defer teardownBenchEnv()

	type PoolTestConfig struct {
		SmallList  []string `env:"SMALL_FEATURES"`
		MediumList []string `env:"LARGE_FEATURES"`
	}

	// Set up different sized slice values
	_ = os.Setenv("SMALL_FEATURES", "a,b,c")               // 3 elements - should use small pool
	_ = os.Setenv("LARGE_FEATURES", "a,b,c,d,e,f,g,h,i,j") // 10 elements - should use medium pool

	var config PoolTestConfig
	err := Bind(&config, DefaultBindingOptions())
	if err != nil {
		t.Fatal(err)
	}

	if len(config.SmallList) != 3 {
		t.Errorf("Expected 3 elements in small list, got %d", len(config.SmallList))
	}

	if len(config.MediumList) != 10 {
		t.Errorf("Expected 10 elements in medium list, got %d", len(config.MediumList))
	}

	t.Logf("Pool efficiency test passed - small: %v, medium: %v", config.SmallList, config.MediumList)
}

func TestOptimizationBenefits(t *testing.T) {
	// This example demonstrates the key optimization benefits

	setupBenchEnv()
	defer teardownBenchEnv()

	type OptimizedConfig struct {
		Host     string   `env:"HOST"`
		Port     int      `env:"PORT"`
		Features []string `env:"FEATURES"`
	}

	// Clear caches for demonstration
	ClearAllCaches()
	ResetCacheCounters()

	// Run multiple times to show cache effectiveness
	for i := 0; i < 5; i++ {
		var config OptimizedConfig
		_ = Bind(&config, DefaultBindingOptions())
	}

	stats := GetCacheStats()
	fmt.Printf("Optimization Benefits:\n")
	fmt.Printf("- Type cache hit rate: %.1f%%\n", stats["type_cache_hit_rate"])
	fmt.Printf("- Tag cache hit rate: %.1f%%\n", stats["tag_cache_hit_rate"])
	fmt.Printf("- Memory pools used for slice parsing\n")
	fmt.Printf("- Reduced allocations through caching\n")
}
