package safezone

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

var ErrTest = errors.New("test error")

func TestError(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		err := New("test error")
		if err.Error() == "" {
			t.Error("Expected non-empty error message")
		}
		if !strings.Contains(err.Error(), "test error") {
			t.Error("Error message does not contain original message")
		}
		if !strings.Contains(err.Error(), "Stack Trace") {
			t.Error("Error does not contain stack trace")
		}
	})

	t.Run("Wrap", func(t *testing.T) {
		originalErr := errors.New("original error")
		wrappedErr := Wrap(originalErr, "wrapped message")
		if !strings.Contains(wrappedErr.Error(), "original error") {
			t.Error("Wrapped error does not contain original error")
		}
		if !strings.Contains(wrappedErr.Error(), "wrapped message") {
			t.Error("Wrapped error does not contain wrapping message")
		}
	})

	t.Run("With", func(t *testing.T) {
		err := New("test error").With("key", "value")
		if !strings.Contains(err.Error(), "key") || !strings.Contains(err.Error(), "value") {
			t.Error("Error does not contain added context")
		}
	})
}

func TestResult(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		result := Ok(42)
		if result.Check() != nil {
			t.Error("Ok result should not have an error")
		}
		if result.Unwrap() != 42 {
			t.Error("Ok result should contain the correct value")
		}
	})

	t.Run("Err", func(t *testing.T) {
		result := Err[int](errors.New("test error"))
		if result.Check() == nil {
			t.Error("Err result should have an error")
		}
	})

	t.Run("UnwrapOr", func(t *testing.T) {
		okResult := Ok(42)
		if okResult.UnwrapOr(0) != 42 {
			t.Error("UnwrapOr should return the value for Ok results")
		}

		errResult := Err[int](errors.New("test error"))
		if errResult.UnwrapOr(0) != 0 {
			t.Error("UnwrapOr should return the default value for Err results")
		}
	})

	t.Run("Map", func(t *testing.T) {
		result := Ok(21).Map(func(i int) int { return i * 2 })
		if result.Unwrap() != 42 {
			t.Error("Map should apply the function to the value")
		}
	})

	t.Run("FlatMap", func(t *testing.T) {
		result := Ok(21).FlatMap(func(i int) Result[int] { return Ok(i * 2) })
		if result.Unwrap() != 42 {
			t.Error("FlatMap should apply the function to the value")
		}
	})
}

func TestHandle(t *testing.T) {
	t.Run("On", func(t *testing.T) {
		var handled bool
		Do(func() error {
			return ErrTest
		}).On(ErrTest, func(err error) {
			handled = true
		})
		if !handled {
			t.Error("On should handle matching errors")
		}
	})

	t.Run("Else", func(t *testing.T) {
		var handled bool
		Do(func() error {
			return errors.New("unexpected error")
		}).On(ErrTest, func(err error) {
			t.Error("This handler should not be called")
		}).Else(func(err error) {
			handled = true
		})
		if !handled {
			t.Error("Else should handle unmatched errors")
		}
	})

	t.Run("MultipleOn", func(t *testing.T) {
		var handled1, handled2 bool
		Do(func() error {
			return ErrTest
		}).On(errors.New("other error"), func(err error) {
			t.Error("This handler should not be called")
		}).On(ErrTest, func(err error) {
			handled1 = true
		}).On(ErrTest, func(err error) {
			handled2 = true
		})
		if !handled1 {
			t.Error("First matching On should handle the error")
		}
		if handled2 {
			t.Error("Subsequent matching On should not handle the error")
		}
	})
}

func TestGroup(t *testing.T) {
	t.Run("NoErrors", func(t *testing.T) {
		var g Group
		g.Go(func() error { return nil })
		g.Go(func() error { return nil })
		if err := g.Wait(); err != nil {
			t.Error("Wait should return nil when no errors occur")
		}
	})

	t.Run("WithErrors", func(t *testing.T) {
		var g Group
		g.Go(func() error { return errors.New("error 1") })
		g.Go(func() error { return errors.New("error 2") })
		err := g.Wait()
		if err == nil {
			t.Error("Wait should return an error when errors occur")
		}
		if !strings.Contains(err.Error(), "error 1") || !strings.Contains(err.Error(), "error 2") {
			t.Error("Combined error should contain all error messages")
		}
	})
}

func TestTry(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		result := Try(func() (int, error) { return 42, nil })
		if result.Check() != nil {
			t.Error("Try should return Ok result for successful operations")
		}
		if result.Unwrap() != 42 {
			t.Error("Try should return correct value for successful operations")
		}
	})

	t.Run("Failure", func(t *testing.T) {
		result := Try(func() (int, error) { return 0, errors.New("test error") })
		if result.Check() == nil {
			t.Error("Try should return Err result for failed operations")
		}
	})
}

func TestMust(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		value := Must(42, nil)
		if value != 42 {
			t.Error("Must should return the value when no error occurs")
		}
	})

	t.Run("Failure", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Must should panic when an error occurs")
			}
		}()
		Must(0, errors.New("test error"))
	})
}

func TestRetry(t *testing.T) {
	t.Run("EventualSuccess", func(t *testing.T) {
		attempts := 0
		err := Retry(context.Background(), func() error {
			attempts++
			if attempts < 3 {
				return errors.New("temporary error")
			}
			return nil
		}, 5)
		if err != nil {
			t.Error("Retry should eventually succeed")
		}
		if attempts != 3 {
			t.Errorf("Expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("Failure", func(t *testing.T) {
		err := Retry(context.Background(), func() error {
			return errors.New("persistent error")
		}, 3)
		if err == nil {
			t.Error("Retry should fail after max attempts")
		}
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()
		err := Retry(ctx, func() error {
			time.Sleep(100 * time.Millisecond)
			return errors.New("slow operation")
		}, 5)
		if !errors.Is(err, context.Canceled) {
			t.Error("Retry should respect context cancellation")
		}
	})
}

func TestRecover(t *testing.T) {
	t.Run("RecoverPanic", func(t *testing.T) {
		var err error
		func() {
			defer Recover(&err)
			panic("test panic")
		}()
		if err == nil {
			t.Error("Recover should catch panics and convert them to errors")
		}
		if !strings.Contains(err.Error(), "test panic") {
			t.Error("Recovered error should contain panic message")
		}
	})

	t.Run("NoPanic", func(t *testing.T) {
		var err error
		func() {
			defer Recover(&err)
			// No panic occurs
		}()
		if err != nil {
			t.Error("Recover should not set error when no panic occurs")
		}
	})
}
