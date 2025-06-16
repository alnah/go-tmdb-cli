// Package helpers provides testing utilities for panic handling.
package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// AssertPanic asserts that a function panics with the expected message.
func AssertPanic(t *testing.T, expectedMessage string, fn func()) {
	t.Helper()

	defer func() {
		r := recover()
		require.NotNil(t, r, "Expected function to panic but it did not")

		if expectedMessage != "" {
			require.Contains(t, r, expectedMessage, "Panic message does not contain expected text")
		}
	}()

	fn()
}

// AssertNoPanic asserts that a function does not panic.
func AssertNoPanic(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		r := recover()
		require.Nil(t, r, "Expected function not to panic but it panicked with: %v", r)
	}()

	fn()
}
