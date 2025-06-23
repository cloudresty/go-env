# Go Env Examples

This directory contains comprehensive examples demonstrating all features of the Go Env package. Each example is self-contained and focuses on specific functionality.

## Examples Overview

### 🔰 [Basic Operations](./basic/)

**Perfect for getting started**

Demonstrates the fundamental environment variable operations:

- Setting and getting environment variables
- Using default values
- Checking if variables exist
- Unsetting variables

```bash
cd examples/basic && go run main.go
```

### 📦 [Base64 Operations](./base64/)

**For binary data and encoding**

Shows how to work with base64-encoded data:

- Storing and retrieving base64-encoded values
- Handling binary data
- Working with JSON configurations
- Error handling for invalid base64

```bash
cd examples/base64 && go run main.go
```

### 🔐 [Sealed Secrets](./sealed-secrets/)

**The killer feature - encryption/decryption**

Demonstrates the revolutionary sealed secrets functionality:

- Encrypting sensitive values with AES-GCM
- Using different encryption keys
- Safe storage of secrets in version control
- Production-ready security features

```bash
cd examples/sealed-secrets && go run main.go
```

### 📁 [File Operations](./file-operations/)

**Working with .env files**

Shows comprehensive file handling:

- Loading from multiple .env files
- Environment-specific configurations
- Saving variables to files
- Handling special characters and quotes
- Error handling for missing files

```bash
cd examples/file-operations && go run main.go
```

### 🏢 [Real-World Application](./real-world-app/)

**Production-ready configuration management**

A complete example showing how to structure configuration in a real application:

- Multi-environment configuration (dev/staging/prod)
- Configuration structs and validation
- Sealed secrets workflow
- Best practices for production deployment

```bash
cd examples/real-world-app && go run main.go
```

## Running All Examples

You can run all examples at once to see the full feature set:

```bash
# From the repository root
for example in examples/*/; do
    echo "=== Running $(basename "$example") ==="
    (cd "$example" && go run main.go)
    echo ""
done
```

## Example Progression

We recommend exploring the examples in this order:

1. **Basic Operations** - Learn the fundamentals
2. **Base64 Operations** - Understand encoding features
3. **File Operations** - Master .env file handling
4. **Sealed Secrets** - Explore the encryption capabilities
5. **Real-World Application** - See it all come together

## Key Features Demonstrated

### 🔐 Sealed Secrets (Game Changer)

The sealed secrets feature allows you to:

- Safely commit encrypted secrets to version control
- Use the same .env files across all environments
- Automatically decrypt secrets at runtime
- Avoid complex secret management infrastructure

### 📦 Base64 Support

Perfect for:

- Binary data in environment variables
- Encoded certificates and keys
- JSON configurations
- Any data that needs encoding

### 🛠️ Production Ready

All examples demonstrate:

- Proper error handling
- Security best practices
- Performance considerations
- Real-world usage patterns

## Integration Examples

Each example can be easily integrated into your existing projects:

```go
// Quick integration example
import "github.com/cloudresty/go-env"

// Basic usage
dbHost := env.Get("DB_HOST", "localhost")

// With base64
cert, err := env.GetB64("TLS_CERT_B64")

// With sealed secrets
apiKey, err := env.GetSealed("API_KEY_SEALED", "ENV_SEAL_KEY")
```

## Next Steps

After exploring these examples:

1. Check out [ADVANCED_USAGE.md](../ADVANCED_USAGE.md) for production patterns
2. Read the main [README.md](../README.md) for complete API documentation
3. Look at the [test files](../env_test.go) for more usage examples

## Contributing Examples

Have a great example that showcases Go Env features? We'd love to include it! Please submit a pull request with:

- A focused, self-contained example
- Clear documentation
- Proper error handling
- Comments explaining the concepts

Each example should demonstrate specific features while being practical and educational.
