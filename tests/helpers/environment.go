// Package helpers provides testing utilities for environment management.
package helpers

import (
	"os"
	"testing"
)

// MockEnvironment temporarily sets environment variables for testing.
type MockEnvironment struct {
	originalVars map[string]string
	setVars      map[string]string
}

// NewMockEnvironment creates a new mock environment manager.
func NewMockEnvironment() *MockEnvironment {
	return &MockEnvironment{
		originalVars: make(map[string]string),
		setVars:      make(map[string]string),
	}
}

// Set sets an environment variable and remembers the original value.
func (me *MockEnvironment) Set(key, value string) {
	if original, exists := os.LookupEnv(key); exists {
		me.originalVars[key] = original
	}
	me.setVars[key] = value
	_ = os.Setenv(key, value)
}

// Restore restores all environment variables to their original state.
func (me *MockEnvironment) Restore() {
	// Remove variables that were set during testing
	for key := range me.setVars {
		if original, hadOriginal := me.originalVars[key]; hadOriginal {
			_ = os.Setenv(key, original)
		} else {
			_ = os.Unsetenv(key)
		}
	}
}

// WithMockEnvironment temporarily sets environment variables for a test function.
func WithMockEnvironment(t *testing.T, vars map[string]string, fn func()) {
	t.Helper()

	env := NewMockEnvironment()
	defer env.Restore()

	for key, value := range vars {
		env.Set(key, value)
	}

	fn()
}
