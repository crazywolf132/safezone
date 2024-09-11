// Package safezone provides an elegant and powerful error handling solution for Go.
package safezone

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Zone represents a safe execution zone where errors are handled automatically.
type Zone struct {
	err                error
	wrapErrors         bool
	debug              bool
	recoveryStrategies []RecoveryStrategy
	plugins            []Plugin
	retryCount         int
}

// Option is a function that configures a Zone.
type Option func(*Zone)

// RecoveryStrategy is a function that attempts to recover from an error.
type RecoveryStrategy func(error) (bool, error)

// Plugin is an interface for extending Zone functionality.
type Plugin interface {
	Name() string
	OnExec(func() error) func() error
	OnTry(interface{}) interface{}
}

// Run executes the given function within a new Zone and returns any error that occurred.
func Run(f func(*Zone), opts ...Option) error {
	z := &Zone{wrapErrors: true}
	for _, opt := range opts {
		opt(z)
	}
	f(z)
	return z.err
}

// WithErrorWrapping sets whether errors should be wrapped with file and line information.
func WithErrorWrapping(wrap bool) Option {
	return func(z *Zone) {
		z.wrapErrors = wrap
	}
}

// WithDebug enables debug mode, which prints function execution information.
func WithDebug(debug bool) Option {
	return func(z *Zone) {
		z.debug = debug
	}
}

// WithRecovery adds recovery strategies to the Zone.
func WithRecovery(strategies ...RecoveryStrategy) Option {
	return func(z *Zone) {
		z.recoveryStrategies = append(z.recoveryStrategies, strategies...)
	}
}

// WithPlugins adds plugins to the Zone.
func WithPlugins(plugins ...Plugin) Option {
	return func(z *Zone) {
		z.plugins = append(z.plugins, plugins...)
	}
}

// Exec executes a function that may return an error.
func (z *Zone) Exec(f func() error) {
	if z.err != nil {
		return
	}

	var err error
	for {
		z.retryCount++ // Increment retry count

		wrappedF := f
		for _, p := range z.plugins {
			wrappedF = p.OnExec(wrappedF)
		}

		if z.debug {
			fmt.Printf("Executing function: %s (attempt %d)\n", functionName(wrappedF), z.retryCount)
		}

		err = wrappedF()
		if err == nil {
			z.retryCount = 0 // Reset retry count on success
			return
		}

		if !z.tryRecover(err) {
			break
		}
	}

	if z.wrapErrors {
		_, file, line, _ := runtime.Caller(1)
		z.err = fmt.Errorf("%s:%d: %w", file, line, err)
	} else {
		z.err = err
	}
}

// Try executes a function that returns a value and an error.
func (z *Zone) Try(f func() (interface{}, error)) interface{} {
	var zero interface{}
	if z.err != nil {
		return zero
	}

	wrapped := func() (interface{}, error) {
		return f()
	}
	for _, p := range z.plugins {
		wrapped = p.OnTry(wrapped).(func() (interface{}, error))
	}

	if z.debug {
		fmt.Printf("Trying function: %s\n", functionName(f))
	}

	result, err := wrapped()
	if err != nil {
		if z.tryRecover(err) {
			return zero
		}
		if z.wrapErrors {
			_, file, line, _ := runtime.Caller(1)
			z.err = fmt.Errorf("%s:%d: %w", file, line, err)
		} else {
			z.err = err
		}
		return zero
	}
	return result
}

// TryNamed is similar to Try but works with named return values.
func (z *Zone) TryNamed(f func() (interface{}, error)) interface{} {
	if z.err != nil {
		var zero interface{}
		return zero
	}

	result, err := f()
	if err != nil {
		if z.wrapErrors {
			_, file, line, _ := runtime.Caller(1)
			z.err = fmt.Errorf("%s:%d: %w", file, line, err)
		} else {
			z.err = err
		}
		var zero interface{}
		return zero
	}

	return result
}

// Recover clears the current error, allowing execution to continue.
func (z *Zone) Recover() {
	z.err = nil
}

// Error returns the current error, if any.
func (z *Zone) Error() error {
	return z.err
}

func (z *Zone) tryRecover(err error) bool {
	for _, strategy := range z.recoveryStrategies {
		if recovered, newErr := strategy(err); recovered {
			if newErr != nil {
				z.err = newErr
			} else {
				z.err = nil // Clear the error if recovery was successful
			}
			return true
		}
	}
	return false
}

// Utility functions

func functionName(i interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	parts := strings.Split(name, ".")
	return parts[len(parts)-1]
}

// Predefined recovery strategies

// RetryN returns a RecoveryStrategy that retries the operation n times.
func RetryN(n int) RecoveryStrategy {
	return func(err error) (bool, error) {
		if n > 0 {
			n--
			return true, nil // Return true to indicate the error is recoverable
		}
		return false, err // Return false after n retries are exhausted
	}
}

// RecoverFrom returns a RecoveryStrategy that recovers from specific error types.
func RecoverFrom(errorTypes ...error) RecoveryStrategy {
	return func(err error) (bool, error) {
		for _, et := range errorTypes {
			if reflect.TypeOf(err) == reflect.TypeOf(et) {
				return true, nil
			}
		}
		return false, err
	}
}
