// Package helpers provides testing utilities for output capture.
package helpers

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// CaptureStderr captures stderr output during function execution.
// This is useful for testing logging output that goes to stderr.
func CaptureStderr(t *testing.T, fn func()) string {
	t.Helper()

	// Save original stderr
	originalStderr := os.Stderr

	// Create pipe to capture output
	r, w, err := os.Pipe()
	require.NoError(t, err, "Failed to create pipe for stderr capture")

	// Replace stderr with write end of pipe
	os.Stderr = w

	// Execute function
	fn()

	// Close write end and restore stderr
	err = w.Close()
	require.NoError(t, err, "Failed to close write end of pipe")
	os.Stderr = originalStderr

	// Read captured output
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err, "Failed to read from pipe")

	return buf.String()
}

// CaptureStdout captures stdout output during function execution.
// This is useful for testing output that goes to stdout.
func CaptureStdout(t *testing.T, fn func()) string {
	t.Helper()

	// Save original stdout
	originalStdout := os.Stdout

	// Create pipe to capture output
	r, w, err := os.Pipe()
	require.NoError(t, err, "Failed to create pipe for stdout capture")

	// Replace stdout with write end of pipe
	os.Stdout = w

	// Execute function
	fn()

	// Close write end and restore stdout
	err = w.Close()
	require.NoError(t, err, "Failed to close write end of pipe")
	os.Stdout = originalStdout

	// Read captured output
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err, "Failed to read from pipe")

	return buf.String()
}
