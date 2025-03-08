# goenv

`goenv` is a lightweight and dependency-free Go package for loading environment variables from `.env` files. It provides a simple and efficient alternative to other `.env` loading libraries, focusing on core functionality and ease of use.

&nbsp;

## Installation

```bash
go get https://github.com/cloudresty/goenv
```

&nbsp;

## Usage

&nbsp;

**Create a `.env` file:**

In the root directory of your project, create a file named `.env` and add your environment variables in the format `KEY=VALUE`, one per line.

&nbsp;

**Example:**

```bash
DB_HOST=localhost
DB_PORT=3000
DB_USER=myuser
DB_PASS=mypassword
API_KEY=your_secret_api_key
# MY_DISABLED_VAR=value_will_be_ignored
```

**Note:** Lines starting with `#` are treated as comments and ignored.

&nbsp;

**Import and use the `goenv` package:**

```go
package main

import (
    "fmt"
    "https://github.com/cloudresty/goenv"
)

func main() {
    // Load environment variables from .env
    err := goenv.Load(".env")
    if err != nil {
        fmt.Println("Error loading .env:", err)
        return
    }

    // Access environment variables
    dbHost := goenv.Get("DB_HOST", "localhost")     // Get with default value
    apiKey, exists := goenv.Lookup("API_KEY")       // Check existence and get value

    fmt.Println("DB_HOST:", dbHost)
    if exists {
        fmt.Println("API_KEY:", apiKey)
    } else {
        fmt.Println("API_KEY not found.")
    }

    // Example using MustLoad (panics on error)
    // goenv.MustLoad(".env")
}
```

&nbsp;

## Functions

&nbsp;

`func Load(filename string) error`

Loads environment variables from the specified .env file.

* `filename`: The path to the `.env` file.
* Returns an `error` if the file cannot be opened or read.
* Ignores empty lines and lines starting with `#` (comments).
* Only sets environment variables that are not already set in the current environment.

&nbsp;

`Lookup(key string) (string, bool)`

Retrieves the value of an environment variable.

* `key`: The name of the environment variable.
* Returns the value of the variable and a boolean indicating whether the variable exists.

&nbsp;

`Get(key string, defaultValue string) string`

Retrieves the value of an environment variable, providing a default value if the variable is not set.

* `key`: The name of the environment variable.
* `defaultValue`: The default value to return if the variable is not set.
* Returns the value of the variable or the default value.

&nbsp;

`MustLoad(filename string)`

Loads environment variables from the specified `.env` file and panics if an error occurs.

* `filename`: The path to the `.env` file.
* Panics if the file cannot be opened or read.

&nbsp;

## Features

* **Simple and lightweight**: No external dependencies.
* **Comment support**: Ignores lines starting with `#`.
* **Empty line handling**: Ignores empty lines.
* **Existing environment variable preservation**: Only sets variables that are not already set.
* **Convenient functions**: `Lookup` and `Get` for easy access to environment variables.
* **Error handling**: Provides error return values for robust applications.
* **`MustLoad` function**: for cases where the env file is critical to application functionality.

&nbsp;

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

&nbsp;

---
Made with ❤️ by [Cloudresty](https://cloudresty.com)


