# Crocodile

**Crocodile** is an in-memory cache written in Go, designed to offer a lightweight and efficient caching solution. Its key feature is the ability to set a memory consumption limit per instance, ensuring optimal resource usage and preventing overconsumption.

## Features

- In-memory caching for fast data retrieval
- Configurable memory limit to control cache size
- Simple and intuitive API
- Thread-safe operations

## Installation

To install Crocodile, use `go get`:

```bash
go get github.com/Toolnado/crocodile
```

## Usage

Here's a basic example of how to use Crocodile:

```go
package main

import (
    "fmt"
    "github.com/Toolnado/crocodile"
)

func main() {
    // Create a new cache instance with a memory limit of 64MB
    cache := crocodile.New(64 * 1024 * 1024)

    // Set a value in the cache
    cache.Set("key",[]byte("value"))

    // Get a value from the cache
    value, found := cache.Get("key")
    if found {
        fmt.Println("Found value:", value)
    } else {
        fmt.Println("Value not found")
    }
}
```

## Configuration

You can configure the memory limit when creating a new cache instance:

```go
// Create a new cache instance with a memory limit of 128MB
cache := crocodile.New(128 * 1024 * 1024)
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes. Be sure to follow the project's coding guidelines and include tests for any new features or bug fixes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
