package safezone

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Error represents an error with additional context and stack trace
type Error struct {
	err        error
	context    map[string]interface{}
	stackTrace string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v\nContext: %v\nStack Trace:\n%s", e.err, e.context, e.stackTrace)
}

func (e *Error) Unwrap() error { return e.err }

// New creates a new Error with stack trace
func New(message string) *Error {
	return &Error{
		err:        errors.New(message),
		context:    make(map[string]interface{}),
		stackTrace: getStackTrace(),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, message string) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		err:        fmt.Errorf("%s: %w", message, err),
		context:    make(map[string]interface{}),
		stackTrace: getStackTrace(),
	}
}

// With adds context to the error
func (e *Error) With(key string, value interface{}) *Error {
	e.context[key] = value
	return e
}

func getStackTrace() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return string(buf[:n])
		}
		buf = make([]byte, 2*len(buf))
	}
}

// Result represents the outcome of an operation that might fail
type Result[T any] struct {
	value T
	err   error
}

// Ok creates a successful Result
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

// Err creates a failed Result
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// Unwrap returns the value if there's no error, otherwise panics
func (r Result[T]) Unwrap() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}

// UnwrapOr returns the value if there's no error, otherwise returns the default value
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.err != nil {
		return defaultValue
	}
	return r.value
}

// UnwrapOrElse returns the value if there's no error, otherwise calls the provided function
func (r Result[T]) UnwrapOrElse(f func(error) T) T {
	if r.err != nil {
		return f(r.err)
	}
	return r.value
}

// Map applies a function to the value if there's no error
func (r Result[T]) Map(f func(T) T) Result[T] {
	if r.err != nil {
		return r
	}
	return Ok(f(r.value))
}

// FlatMap applies a function that returns a Result
func (r Result[T]) FlatMap(f func(T) Result[T]) Result[T] {
	if r.err != nil {
		return r
	}
	return f(r.value)
}

// Check returns the error if there is one, otherwise returns nil
func (r Result[T]) Check() error {
	return r.err
}

// Try attempts to execute a function and returns a Result
func Try[T any](f func() (T, error)) Result[T] {
	value, err := f()
	if err != nil {
		return Err[T](Wrap(err, "operation failed"))
	}
	return Ok(value)
}

// Must panics if err is not nil, otherwise returns the value
func Must[T any](value T, err error) T {
	if err != nil {
		panic(Wrap(err, "assertion failed"))
	}
	return value
}

// Handle provides a fluent interface for error handling
type Handle struct {
	err error
}

// On registers an error handler for a specific error type
func (h Handle) On(target error, handler func(error)) Handle {
	if h.err != nil && errors.Is(h.err, target) {
		handler(h.err)
		h.err = nil
	}
	return h
}

// Else handles any remaining error
func (h Handle) Else(handler func(error)) {
	if h.err != nil {
		handler(h.err)
	}
}

// Do executes a function and returns a Handle for error handling
func Do(f func() error) Handle {
	return Handle{err: f()}
}

// Retry retries a function with exponential backoff
func Retry(ctx context.Context, f func() error, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		if err = f(); err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(1<<uint(i)) * time.Second):
		}
	}
	return Wrap(err, fmt.Sprintf("operation failed after %d retries", maxRetries))
}

// Group runs functions concurrently and collects their errors
type Group struct {
	wg     sync.WaitGroup
	errMux sync.Mutex
	errs   []error
}

// Go runs the given function in a goroutine
func (g *Group) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(); err != nil {
			g.errMux.Lock()
			g.errs = append(g.errs, err)
			g.errMux.Unlock()
		}
	}()
}

// Wait waits for all goroutines to complete and returns a combined error
func (g *Group) Wait() error {
	g.wg.Wait()
	if len(g.errs) == 0 {
		return nil
	}
	return Wrap(errors.Join(g.errs...), "multiple errors occurred")
}

// Recover is a function that can be used in a defer statement to recover from panics
func Recover(errPtr *error) {
	if r := recover(); r != nil {
		*errPtr = Wrap(fmt.Errorf("%v", r), "panic recovered")
	}
}
