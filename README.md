# SafeZone: Tame Your Errors 🛡️

[![Go Report Card](https://goreportcard.com/badge/github.com/crazywolf132/safezone)](https://goreportcard.com/report/github.com/crazywolf132/safezone)
[![GoDoc](https://godoc.org/github.com/crazywolf132/safezone?status.svg)](https://godoc.org/github.com/crazywolf132/safezone)

SafeZone is a revolutionary error handling library for Go that makes dealing with errors a breeze. Say goodbye to repetitive `if err != nil` checks and hello to clean, expressive, and powerful error handling!

## Table of Contents

- [Features](#-features)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Core Concepts](#-core-concepts)
- [API Reference](#-api-reference)
- [Advanced Usage](#-advanced-usage)
- [Best Practices](#-best-practices)
- [Use Cases](#-use-cases)
- [Contributing](#-contributing)

## 🚀 Features

- **Simple Core Concept**: Wrap your code in a `safezone.Run` and forget about constant error checking.
- **Type-Safe Error Handling**: Use `Try` and `TryNamed` for functions that return values along with errors.
- **Smart Recovery System**: Define custom strategies to recover from errors automatically.
- **Debug Mode**: Trace execution flow with built-in debugging.
- **Flexible Configuration**: Use functional options to configure SafeZone to your needs.
- **Plugin System**: Extend SafeZone's functionality with custom plugins.
- **Performance Focused**: Designed to have minimal overhead compared to traditional error handling.

## 📦 Installation

To install SafeZone, use `go get`:

```
go get github.com/crazywolf132/safezone
```

## 🚦 Quick Start

Here's a simple example to get you started with SafeZone:

```go
package main

import (
    "fmt"
    "github.com/crazywolf132/safezone"
)

func main() {
    err := safezone.Run(func(z *safezone.Zone) {
        z.Exec(func() error {
            fmt.Println("Hello, SafeZone!")
            return nil
        })
    })

    if err != nil {
        fmt.Printf("An error occurred: %v\n", err)
    }
}
```

## 🧠 Core Concepts

SafeZone is built around two main concepts:

1. **Zones**: A `Zone` is a protected area where errors are automatically handled.
2. **Error Handling Methods**: Within a zone, you use methods like `Exec`, `Try`, and `TryNamed` to execute functions that might return errors.

## 📚 API Reference

### safezone.Run

```go
func Run(f func(*Zone), opts ...Option) error
```

`Run` creates a new `Zone` and executes the provided function within it. It returns any error that occurred during execution.

**When to use**: Use `Run` as the entry point for SafeZone in your application. It's ideal for wrapping main functions or significant portions of your code.

**Example**:
```go
err := safezone.Run(func(z *safezone.Zone) {
    // Your code here
})
```

### Zone.Exec

```go
func (z *Zone) Exec(f func() error)
```

`Exec` executes a function that returns only an error.

**When to use**: Use `Exec` when you have a function that returns only an error and you don't need its result.

**Example**:
```go
z.Exec(func() error {
    return os.Remove("temp.txt")
})
```

### Zone.Try

```go
func (z *Zone) Try[T any](f func() (T, error)) T
```

`Try` executes a function that returns a value and an error. It returns the value if no error occurred.

**When to use**: Use `Try` when you have a function that returns both a value and an error, and you need to use the returned value.

**Example**:
```go
content := z.Try(ioutil.ReadFile, "config.json")
```

### Zone.TryNamed

```go
func (z *Zone) TryNamed[T any](f func() (result T, err error)) (result T)
```

`TryNamed` is similar to `Try`, but works with functions that use named return values.

**When to use**: Use `TryNamed` when working with functions that have named return values, or when you want to make your code more explicit.

**Example**:
```go
sum := z.TryNamed(func() (result int, err error) {
    result = 1 + 2
    return
})
```

### Zone.Recover

```go
func (z *Zone) Recover()
```

`Recover` clears the current error in the zone, allowing execution to continue.

**When to use**: Use `Recover` when you want to clear an error and continue execution within the same zone.

**Example**:
```go
z.Exec(func() error {
    return fmt.Errorf("an error occurred")
})
z.Recover() // Clear the error
z.Exec(func() error {
    fmt.Println("This will execute")
    return nil
})
```

### Zone.Error

```go
func (z *Zone) Error() error
```

`Error` returns the current error in the zone, if any.

**When to use**: Use `Error` when you need to check the current error within a zone without ending the zone's execution.

**Example**:
```go
if z.Error() != nil {
    fmt.Println("An error has occurred, but we're continuing")
}
```

## 🔧 Advanced Usage

### Recovery Strategies

SafeZone allows you to define custom recovery strategies:

```go
safezone.Run(func(z *safezone.Zone) {
    z.Try(riskyOperation)
}, safezone.WithRecovery(
    safezone.RetryN(3),
    safezone.RecoverFrom(os.ErrNotExist),
))
```

### Plugins

Extend SafeZone's functionality with plugins:

```go
type LoggingPlugin struct{}

func (p LoggingPlugin) Name() string { return "LoggingPlugin" }
func (p LoggingPlugin) OnExec(f func() error) func() error {
    return func() error {
        fmt.Println("Before execution")
        err := f()
        fmt.Println("After execution")
        return err
    }
}
func (p LoggingPlugin) OnTry(f interface{}) interface{} { return f }

safezone.Run(func(z *safezone.Zone) {
    // Your code here
}, safezone.WithPlugins(LoggingPlugin{}))
```

### Debug Mode

Enable debug mode to trace execution flow:

```go
safezone.Run(func(z *safezone.Zone) {
    // Your code here
}, safezone.WithDebug(true))
```

## 👑 Best Practices

1. **Use `Run` at the highest level possible**: Ideally, wrap your `main` function or major components of your application with `safezone.Run`.

2. **Prefer `Try` over `Exec` when possible**: `Try` provides type-safety and allows you to use the returned value directly.

3. **Use recovery strategies for known error cases**: This can help your application gracefully handle expected error scenarios.

4. **Leverage plugins for cross-cutting concerns**: Use plugins for logging, metrics, or any functionality that needs to be applied across multiple operations.

5. **Use `Recover` judiciously**: While `Recover` can be useful, overuse can lead to ignoring important errors. Use it when you have a specific reason to continue execution after an error.

6. **Combine SafeZone with context for timeouts and cancellation**: SafeZone works well with Go's context package for managing timeouts and cancellation.

## 🎯 Use Cases

1. **Web Servers**: Streamline error handling in HTTP handlers.
   ```go
   http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
       safezone.Run(func(z *safezone.Zone) {
           // Handle request
       })
   })
   ```

2. **Data Processing Pipelines**: Manage errors across multiple processing stages effortlessly.
   ```go
   safezone.Run(func(z *safezone.Zone) {
       data := z.Try(readData)
       processedData := z.Try(processData, data)
       z.Exec(func() error { return writeData(processedData) })
   })
   ```

3. **CLI Applications**: Provide user-friendly error messages without cluttering your code.
   ```go
   func main() {
       err := safezone.Run(func(z *safezone.Zone) {
           // CLI logic here
       })
       if err != nil {
           fmt.Fprintf(os.Stderr, "Error: %v\n", err)
           os.Exit(1)
       }
   }
   ```

4. **Database Operations**: Simplify transaction management and error recovery.
   ```go
   safezone.Run(func(z *safezone.Zone) {
       tx := z.Try(db.Begin)
       defer z.Exec(tx.Rollback)
       
       z.Exec(func() error { return tx.Exec("INSERT INTO ...") })
       z.Exec(func() error { return tx.Exec("UPDATE ...") })
       
       z.Exec(tx.Commit)
   })
   ```

5. **Concurrent Operations**: Handle errors in goroutines with ease.
   ```go
   var wg sync.WaitGroup
   results := make(chan error, numWorkers)
   
   for i := 0; i < numWorkers; i++ {
       wg.Add(1)
       go func() {
           defer wg.Done()
           results <- safezone.Run(func(z *safezone.Zone) {
               // Worker logic here
           })
       }()
   }
   
   wg.Wait()
   close(results)
   
   for err := range results {
       if err != nil {
           fmt.Printf("Worker error: %v\n", err)
       }
   }
   ```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for more details.

## 🌟 Show Your Support

If SafeZone has made your Go programming life easier, consider giving it a star on GitHub! It helps others discover the project and motivates us to keep improving.

Remember, in the SafeZone, errors fear to tread! Happy coding! 🚀