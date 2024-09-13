# Safezone: 🛡️ Bulletproof Error Handling for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/crazywolf132/safezone)](https://goreportcard.com/report/github.com/crazywolf132/safezone)
[![GoDoc](https://godoc.org/github.com/crazywolf132/safezone?status.svg)](https://godoc.org/github.com/crazywolf132/safezone)

**Safezone** brings the power and elegance of modern error handling to Go, taking it to the next level! Whether you're building microservices, CLI tools, or complex applications, `safezone` makes error management in Go more intuitive, expressive, and robust. With features like rich error context, a Result type, fluent error handling, and even retry mechanisms, you'll wonder how you ever handled errors without it! 🚀

## Why Safezone?

Go's built-in error handling, while simple, can lead to verbose code and doesn't provide rich context out of the box. This can make debugging and error management challenging, especially in large codebases.

**Safezone** aims to solve these problems by:

- **Simplifying error handling**: Use a fluent API for cleaner, more readable code.
- **Enhancing debugging**: Rich error context and stack traces make identifying issues a breeze.
- **Adding advanced features**: Result types, retry mechanisms, and panic recovery give you more control.
- **Improving concurrency**: Built-in support for handling errors from multiple goroutines.

## 🌟 Key Features

- **🛡️ Rich Error Context**: Easily add and retrieve contextual information for errors.
- **🎁 Result Type**: A generic type for handling operations that might fail, inspired by Rust.
- **🔗 Fluent Error Handling**: Chain error handlers for different error types with a clean syntax.
- **🔄 Retry Mechanism**: Built-in support for retrying operations with exponential backoff.
- **🦺 Panic Recovery**: Safely convert panics to errors.
- **🧵 Concurrent Error Handling**: Manage errors from multiple goroutines with ease.
- **🔍 Stack Traces**: Automatically capture stack traces for enhanced debugging.
- **🔧 Customizable**: Extend with your own error types and handling logic.

## 📦 Installation

To add `safezone` to your project, simply run:

```bash
go get github.com/crazywolf132/safezone
```

## 🚀 Quick Start Guide

Start using `safezone` with just a few lines of code:

```go
package main

import (
    "fmt"
    "github.com/crazywolf132/safezone"
)

func main() {
    // Basic error creation with context
    err := safezone.New("something went wrong").With("details", "more info")
    fmt.Println(err) // Prints error with context and stack trace

    // Using the Result type
    result := divide(10, 2)
    if value, err := result.Unwrap(); err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Result:", value)
    }

    // Fluent error handling
    safezone.Do(func() error {
        return someRiskyOperation()
    }).On(ErrNotFound, func(err error) {
        fmt.Println("Resource not found:", err)
    }).Else(func(err error) {
        fmt.Println("An unexpected error occurred:", err)
    })

    // Retry mechanism
    err = safezone.Retry(context.Background(), func() error {
        return unreliableOperation()
    }, 3)
    if err != nil {
        fmt.Println("Operation failed after 3 retries:", err)
    }
}

func divide(a, b int) safezone.Result[int] {
    if b == 0 {
        return safezone.Err[int](fmt.Errorf("division by zero"))
    }
    return safezone.Ok(a / b)
}
```

## 📚 Comprehensive Feature Guide

### 1. Rich Error Context

Add context to your errors for better debugging:

```go
err := safezone.New("database connection failed").
    With("host", "example.com").
    With("port", 5432)
```

### 2. Result Type

Handle operations that might fail with the `Result` type:

```go
func fetchUser(id int) safezone.Result[User] {
    // ... implementation ...
}

result := fetchUser(123)
user := result.UnwrapOr(defaultUser)
```

### 3. Fluent Error Handling

Chain error handlers for clean, readable code:

```go
safezone.Do(func() error {
    return someOperation()
}).On(ErrNotFound, func(err error) {
    // Handle not found error
}).On(ErrPermission, func(err error) {
    // Handle permission error
}).Else(func(err error) {
    // Handle any other error
})
```

### 4. Retry Mechanism

Easily retry operations with exponential backoff:

```go
err := safezone.Retry(ctx, func() error {
    return unreliableNetworkCall()
}, 5) // Retry up to 5 times
```

### 5. Panic Recovery

Safely convert panics to errors:

```go
func riskyFunction() (err error) {
    defer safezone.Recover(&err)
    // ... implementation that might panic ...
}
```

### 6. Concurrent Error Handling

Manage errors from multiple goroutines:

```go
var g safezone.Group
g.Go(func() error {
    return operation1()
})
g.Go(func() error {
    return operation2()
})
if err := g.Wait(); err != nil {
    fmt.Println("One or more operations failed:", err)
}
```

## 🤝 Contributing to Safezone

We welcome contributions! Whether it's a bug report, new feature, or a pull request, we're happy to have your input. Check out our [contribution guidelines](CONTRIBUTING.md) to get started.

## 📄 License

Safezone is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

✨ **Safezone**: Empowering Go developers with robust, expressive, and painless error handling. Say goodbye to cryptic errors and hello to clear, contextual problem solving!