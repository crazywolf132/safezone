package safezone

import (
	"errors"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("No error", func(t *testing.T) {
		err := Run(func(z *Zone) {
			z.Exec(func() error {
				return nil
			})
		})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("With error", func(t *testing.T) {
		err := Run(func(z *Zone) {
			z.Exec(func() error {
				return errors.New("test error")
			})
		})
		if err == nil {
			t.Error("Expected an error, got nil")
		}
	})
}

func TestExec(t *testing.T) {
	t.Run("No error", func(t *testing.T) {
		err := Run(func(z *Zone) {
			z.Exec(func() error {
				return nil
			})
		})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("With error", func(t *testing.T) {
		err := Run(func(z *Zone) {
			z.Exec(func() error {
				return errors.New("test error")
			})
		})
		if err == nil || !strings.Contains(err.Error(), "test error") {
			t.Errorf("Expected error containing 'test error', got %v", err)
		}
	})
}

func TestTry(t *testing.T) {
	t.Run("No error", func(t *testing.T) {
		var result int
		err := Run(func(z *Zone) {
			res := z.Try(func() (interface{}, error) {
				return 42, nil
			})
			result = res.(int)
		})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != 42 {
			t.Errorf("Expected 42, got %d", result)
		}
	})

	t.Run("With error", func(t *testing.T) {
		var result int
		err := Run(func(z *Zone) {
			res := z.Try(func() (interface{}, error) {
				return 0, errors.New("test error")
			})
			if res != nil {
				result = res.(int)
			}
		})
		if err == nil || !strings.Contains(err.Error(), "test error") {
			t.Errorf("Expected error containing 'test error', got %v", err)
		}
		if result != 0 {
			t.Errorf("Expected 0, got %d", result)
		}
	})
}

func TestTryNamed(t *testing.T) {
	t.Run("No error", func(t *testing.T) {
		var result int
		err := Run(func(z *Zone) {
			res := z.TryNamed(func() (interface{}, error) {
				return 42, nil
			})
			result = res.(int)
		})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != 42 {
			t.Errorf("Expected 42, got %d", result)
		}
	})

	t.Run("With error", func(t *testing.T) {
		var result int
		err := Run(func(z *Zone) {
			res := z.TryNamed(func() (interface{}, error) {
				return 0, errors.New("test error")
			})
			if res != nil {
				result = res.(int)
			}
		})
		if err == nil || !strings.Contains(err.Error(), "test error") {
			t.Errorf("Expected error containing 'test error', got %v", err)
		}
		if result != 0 {
			t.Errorf("Expected 0, got %d", result)
		}
	})
}

func TestRecover(t *testing.T) {
	err := Run(func(z *Zone) {
		z.Exec(func() error {
			return errors.New("first error")
		})
		z.Recover()
		z.Exec(func() error {
			return nil
		})
	})
	if err != nil {
		t.Errorf("Expected no error after recovery, got %v", err)
	}
}

func TestError(t *testing.T) {
	Run(func(z *Zone) {
		z.Exec(func() error {
			return errors.New("test error")
		})
		if z.Error() == nil || !strings.Contains(z.Error().Error(), "test error") {
			t.Errorf("Expected error containing 'test error', got %v", z.Error())
		}
	})
}

func TestWithErrorWrapping(t *testing.T) {
	err := Run(func(z *Zone) {
		z.Exec(func() error {
			return errors.New("test error")
		})
	}, WithErrorWrapping(true))
	if err == nil || !strings.Contains(err.Error(), "test error") {
		t.Errorf("Expected wrapped error containing 'test error', got %v", err)
	}
}

func TestWithDebug(t *testing.T) {
	// This is a basic test to ensure WithDebug doesn't cause errors
	// A more comprehensive test would involve capturing stdout and checking the output
	err := Run(func(z *Zone) {
		z.Exec(func() error {
			return nil
		})
	}, WithDebug(true))
	if err != nil {
		t.Errorf("Expected no error with debug mode, got %v", err)
	}
}

func TestWithRecovery(t *testing.T) {
	attempts := 0
	err := Run(func(z *Zone) {
		z.Exec(func() error {
			attempts++
			if attempts <= 3 {
				return errors.New("retry error")
			}
			return nil
		})
	}, WithRecovery(RetryN(3)), WithDebug(true))
	if err != nil {
		t.Errorf("Expected no error after retries, got %v", err)
	}
	if attempts != 4 {
		t.Errorf("Expected 4 attempts (including initial), got %d", attempts)
	}
}

func TestWithPlugins(t *testing.T) {
	callCount := 0
	testPlugin := PluginFunc(func(f func() error) func() error {
		return func() error {
			callCount++
			return f()
		}
	})

	err := Run(func(z *Zone) {
		z.Exec(func() error {
			return nil
		})
	}, WithPlugins(testPlugin))

	if err != nil {
		t.Errorf("Expected no error with plugin, got %v", err)
	}
	if callCount != 1 {
		t.Errorf("Expected plugin to be called once, got %d", callCount)
	}
}

// Helper function to check if an error contains a specific substring
func ErrorContains(err error, substr string) bool {
	return err != nil && strings.Contains(err.Error(), substr)
}

// Helper types and functions

type PluginFunc func(func() error) func() error

func (p PluginFunc) Name() string                       { return "TestPlugin" }
func (p PluginFunc) OnExec(f func() error) func() error { return p(f) }
func (p PluginFunc) OnTry(f interface{}) interface{}    { return f }
