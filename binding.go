package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Optimization: Cache reflection types to avoid repeated TypeOf calls
var (
	durationType   = reflect.TypeOf(time.Duration(0))
	typeCache      = make(map[reflect.Type]*typeInfo)
	typeCacheMutex sync.RWMutex

	tagCache      = make(map[string]tagInfo)
	tagCacheMutex sync.RWMutex

	// Buffer pools optimized for environment variables (typically small values)
	// Small pool for most env vars (1-4 elements)
	smallStringSlicePool = sync.Pool{
		New: func() interface{} {
			return make([]string, 0, 4)
		},
	}
	// Medium pool for larger env vars (5-16 elements)
	mediumStringSlicePool = sync.Pool{
		New: func() interface{} {
			return make([]string, 0, 16)
		},
	}

	// Lock-free caching using atomic operations for frequently accessed data
	// Cache hit counters for monitoring
	typeCacheHits   int64
	typeCacheMisses int64
	tagCacheHits    int64
	tagCacheMisses  int64
)

// tagInfo caches parsed tag data
type tagInfo struct {
	envKey     string
	isRequired bool
}

// typeInfo caches reflection information for a struct type
type typeInfo struct {
	fields []fieldInfo
}

// fieldInfo caches field-specific reflection data
type fieldInfo struct {
	index       int
	name        string
	tag         string
	kind        reflect.Kind
	fieldType   reflect.Type
	isDuration  bool
	isAnonymous bool
}

// BindingOptions provides configuration for struct binding
type BindingOptions struct {
	Tag          string // Tag name to look for (default: "env")
	Required     bool   // Whether to fail if required fields are missing
	Prefix       string // Prefix to add to all environment variable names
	DefaultValue string // Default tag name for default values (default: "default")
}

// DefaultBindingOptions returns default binding options
func DefaultBindingOptions() BindingOptions {
	return BindingOptions{
		Tag:          "env",
		Required:     false,
		Prefix:       "",
		DefaultValue: "default",
	}
}

// Bind binds environment variables to struct fields with configurable options.
// Optimized version with reflection caching and reduced allocations.
// Example usage:
//
//	type Config struct {
//		Port     int    `env:"PORT" default:"8080"`
//		Host     string `env:"HOST" default:"localhost"`
//		Debug    bool   `env:"DEBUG" default:"false"`
//		Features []string `env:"FEATURES"`
//		Timeout  time.Duration `env:"TIMEOUT" default:"30s"`
//	}
//	var config Config
//	err := env.Bind(&config, env.DefaultBindingOptions())
func Bind(out interface{}, options BindingOptions) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("output parameter must be a pointer to a struct")
	}

	return bindStruct(v.Elem(), v.Elem().Type(), options)
}

// bindStruct uses cached reflection data for better performance
func bindStruct(structValue reflect.Value, structType reflect.Type, options BindingOptions) error {
	// Get or cache type information
	info := getTypeInfo(structType)

	for _, fieldInfo := range info.fields {
		field := structValue.Field(fieldInfo.index)

		if !field.CanSet() {
			continue // Skip unexported fields
		}

		// Handle nested structs
		if fieldInfo.kind == reflect.Struct {
			if fieldInfo.isAnonymous {
				// Embedded struct
				if err := bindStruct(field, fieldInfo.fieldType, options); err != nil {
					return err
				}
			} else {
				// Named nested struct
				if err := bindStruct(field, fieldInfo.fieldType, options); err != nil {
					return err
				}
			}
			continue
		}

		if fieldInfo.tag == "" {
			continue // Skip if no tag value
		}

		// Parse tag options efficiently (avoid repeated allocations)
		envKey, isRequired := parseTag(fieldInfo.tag, options)

		// Get the environment variable value
		envValue, exists := os.LookupEnv(envKey)

		// Handle missing values
		if !exists {
			// Check for default value using cached struct field
			structField := structType.Field(fieldInfo.index)
			defaultValue := structField.Tag.Get(options.DefaultValue)
			if defaultValue != "" {
				envValue = defaultValue
			} else if isRequired || options.Required {
				return fmt.Errorf("required environment variable %s is not set for field %s", envKey, fieldInfo.name)
			} else {
				continue // Skip if not required and no default
			}
		}

		// Convert the value to the appropriate type with optimizations
		if err := setField(field, envValue, &fieldInfo); err != nil {
			return fmt.Errorf("failed to set field %s from env var %s: %w", fieldInfo.name, envKey, err)
		}
	}

	return nil
}

// getTypeInfo returns cached type information or creates and caches it
func getTypeInfo(t reflect.Type) *typeInfo {
	typeCacheMutex.RLock()
	if info, exists := typeCache[t]; exists {
		typeCacheMutex.RUnlock()
		atomic.AddInt64(&typeCacheHits, 1)
		return info
	}
	typeCacheMutex.RUnlock()
	atomic.AddInt64(&typeCacheMisses, 1)

	// Build type info
	typeCacheMutex.Lock()
	defer typeCacheMutex.Unlock()

	// Double-check after acquiring write lock
	if info, exists := typeCache[t]; exists {
		return info
	}

	info := &typeInfo{
		fields: make([]fieldInfo, 0, t.NumField()),
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("env") // Default tag, will be overridden by options

		fieldInfo := fieldInfo{
			index:       i,
			name:        field.Name,
			tag:         tag,
			kind:        field.Type.Kind(),
			fieldType:   field.Type,
			isDuration:  field.Type == durationType,
			isAnonymous: field.Anonymous,
		}

		info.fields = append(info.fields, fieldInfo)
	}

	typeCache[t] = info
	return info
}

// parseTag parses tag options efficiently with caching and minimal allocations
func parseTag(tag string, options BindingOptions) (envKey string, isRequired bool) {
	if tag == "" {
		return "", false
	}

	// Create cache key including prefix for proper caching
	cacheKey := tag
	if options.Prefix != "" {
		cacheKey = options.Prefix + ":" + tag
	}

	// Check cache first
	tagCacheMutex.RLock()
	if info, exists := tagCache[cacheKey]; exists {
		tagCacheMutex.RUnlock()
		atomic.AddInt64(&tagCacheHits, 1)
		return info.envKey, info.isRequired
	}
	tagCacheMutex.RUnlock()
	atomic.AddInt64(&tagCacheMisses, 1)

	// Parse and cache result
	tagCacheMutex.Lock()
	defer tagCacheMutex.Unlock()

	// Double-check after acquiring write lock
	if info, exists := tagCache[cacheKey]; exists {
		return info.envKey, info.isRequired
	}

	// Fast path for simple tags (no comma) - optimized with IndexByte
	commaIndex := strings.IndexByte(tag, ',')
	if commaIndex == -1 {
		envKey = tag
		if options.Prefix != "" {
			envKey = options.Prefix + envKey
		}
		isRequired = false
	} else {
		// Parse complex tags
		envKey = tag[:commaIndex]
		if options.Prefix != "" {
			envKey = options.Prefix + envKey
		}

		// Optimized required check using byte-level search
		remainingTag := tag[commaIndex+1:]
		isRequired = strings.Contains(remainingTag, "required")
	}

	// Cache the result for future use
	tagCache[cacheKey] = tagInfo{
		envKey:     envKey,
		isRequired: isRequired,
	}

	return envKey, isRequired
}

// setField sets field values with optimizations for common types
func setField(field reflect.Value, value string, info *fieldInfo) error {
	switch info.kind {
	case reflect.String:
		field.SetString(value)
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if info.isDuration {
			duration, err := time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("invalid duration format: %w", err)
			}
			field.SetInt(int64(duration))
			return nil
		}

		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintValue)
		return nil

	case reflect.Bool:
		// Optimized boolean parsing for common values
		switch value {
		case "true", "1", "TRUE", "True":
			field.SetBool(true)
			return nil
		case "false", "0", "FALSE", "False", "":
			field.SetBool(false)
			return nil
		default:
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			field.SetBool(boolValue)
			return nil
		}

	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatValue)
		return nil

	case reflect.Slice:
		return setSlice(field, value)

	case reflect.Map:
		mapValue := reflect.MakeMap(field.Type())
		if err := json.Unmarshal([]byte(value), mapValue.Addr().Interface()); err != nil {
			return fmt.Errorf("failed to parse map from JSON: %w", err)
		}
		field.Set(mapValue)
		return nil

	case reflect.Struct:
		if err := json.Unmarshal([]byte(value), field.Addr().Interface()); err != nil {
			return fmt.Errorf("failed to parse struct from JSON: %w", err)
		}
		return nil

	case reflect.Ptr:
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		// Create a temporary fieldInfo for the pointer element
		elemInfo := &fieldInfo{
			kind:       field.Elem().Kind(),
			fieldType:  field.Type().Elem(),
			isDuration: field.Type().Elem() == durationType,
		}
		return setField(field.Elem(), value, elemInfo)

	default:
		return fmt.Errorf("unsupported field type: %s", field.Type())
	}
}

// setSlice optimizes slice parsing with buffer pooling and pre-allocation
func setSlice(field reflect.Value, value string) error {
	if value == "" {
		return nil
	}

	// Count commas to decide which pool to use
	commaCount := strings.Count(value, ",")
	expectedLen := commaCount + 1

	// Choose appropriate pool based on expected size
	var elementsInterface interface{}
	if expectedLen <= 4 {
		elementsInterface = smallStringSlicePool.Get()
	} else {
		elementsInterface = mediumStringSlicePool.Get()
	}

	elements := elementsInterface.([]string)
	elements = elements[:0] // Reset length but keep capacity

	defer func() {
		if expectedLen <= 4 {
			smallStringSlicePool.Put(elementsInterface)
		} else {
			mediumStringSlicePool.Put(elementsInterface)
		}
	}()

	// Ensure capacity - if pool slice is too small, allocate new one
	if cap(elements) < expectedLen {
		elements = make([]string, 0, expectedLen)
	}

	// Parse elements using byte-level operations for better performance
	start := 0
	for i := 0; i <= len(value); i++ {
		if i == len(value) || value[i] == ',' {
			if i > start {
				elem := strings.TrimSpace(value[start:i])
				if elem != "" {
					elements = append(elements, elem)
				}
			}
			start = i + 1
		}
	}

	// Create slice with exact capacity needed
	slice := reflect.MakeSlice(field.Type(), 0, len(elements))

	for _, elem := range elements {
		elemValue := reflect.New(field.Type().Elem()).Elem()
		elemInfo := &fieldInfo{
			kind:       elemValue.Kind(),
			fieldType:  field.Type().Elem(),
			isDuration: field.Type().Elem() == durationType,
		}

		if err := setField(elemValue, elem, elemInfo); err != nil {
			return fmt.Errorf("failed to parse slice element '%s': %w", elem, err)
		}
		slice = reflect.Append(slice, elemValue)
	}

	field.Set(slice)
	return nil
}

// BindJSON binds environment variables to struct fields using the "env" tag.
// This is a simplified version of Bind for backward compatibility.
func BindJSON(tag string, out interface{}) error {
	options := BindingOptions{
		Tag:          tag,
		Required:     false,
		Prefix:       "",
		DefaultValue: "default",
	}
	return Bind(out, options)
}

// ClearTypeCache clears the reflection cache (useful for testing or memory management)
func ClearTypeCache() {
	typeCacheMutex.Lock()
	defer typeCacheMutex.Unlock()
	typeCache = make(map[reflect.Type]*typeInfo)
}

// ClearAllCaches clears all optimization caches
func ClearAllCaches() {
	typeCacheMutex.Lock()
	typeCache = make(map[reflect.Type]*typeInfo)
	typeCacheMutex.Unlock()

	tagCacheMutex.Lock()
	tagCache = make(map[string]tagInfo)
	tagCacheMutex.Unlock()
}

// GetCacheStats returns statistics about cache usage for monitoring
func GetCacheStats() map[string]interface{} {
	typeCacheMutex.RLock()
	typeCount := len(typeCache)
	typeCacheMutex.RUnlock()

	tagCacheMutex.RLock()
	tagCount := len(tagCache)
	tagCacheMutex.RUnlock()

	typeHits := atomic.LoadInt64(&typeCacheHits)
	typeMisses := atomic.LoadInt64(&typeCacheMisses)
	tagHits := atomic.LoadInt64(&tagCacheHits)
	tagMisses := atomic.LoadInt64(&tagCacheMisses)

	var typeHitRate, tagHitRate float64
	if typeHits+typeMisses > 0 {
		typeHitRate = float64(typeHits) / float64(typeHits+typeMisses) * 100
	}
	if tagHits+tagMisses > 0 {
		tagHitRate = float64(tagHits) / float64(tagHits+tagMisses) * 100
	}

	return map[string]interface{}{
		"type_cache_size":     typeCount,
		"tag_cache_size":      tagCount,
		"type_cache_hits":     typeHits,
		"type_cache_misses":   typeMisses,
		"type_cache_hit_rate": typeHitRate,
		"tag_cache_hits":      tagHits,
		"tag_cache_misses":    tagMisses,
		"tag_cache_hit_rate":  tagHitRate,
	}
}

// ResetCacheCounters resets all performance counters
func ResetCacheCounters() {
	atomic.StoreInt64(&typeCacheHits, 0)
	atomic.StoreInt64(&typeCacheMisses, 0)
	atomic.StoreInt64(&tagCacheHits, 0)
	atomic.StoreInt64(&tagCacheMisses, 0)
}
