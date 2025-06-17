// Package helpers provides testing utilities for the TMDB CLI test suite.
package helpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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

// Unset removes an environment variable and remembers it was unset.
func (me *MockEnvironment) Unset(key string) {
	if original, exists := os.LookupEnv(key); exists {
		me.originalVars[key] = original
	}
	me.setVars[key] = "" // Mark as intentionally unset
	_ = os.Unsetenv(key)
}

// Get retrieves the current value of an environment variable.
func (me *MockEnvironment) Get(key string) string {
	return os.Getenv(key)
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

// Clear removes all tracked environment variables and resets state.
func (me *MockEnvironment) Clear() {
	me.Restore()
	me.originalVars = make(map[string]string)
	me.setVars = make(map[string]string)
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

// WithCleanEnvironment runs a test function with specified environment variables unset.
func WithCleanEnvironment(t *testing.T, varsToUnset []string, fn func()) {
	t.Helper()

	env := NewMockEnvironment()
	defer env.Restore()

	for _, key := range varsToUnset {
		env.Unset(key)
	}

	fn()
}

// SetupTestEnvironment sets up a complete test environment with common TMDB variables.
func SetupTestEnvironment(t *testing.T) *MockEnvironment {
	t.Helper()

	env := NewMockEnvironment()

	// Set common test environment variables
	env.Set("TMDB_API_KEY", "test-api-key")
	env.Set("TMDB_BASE_URL", "https://api.themoviedb.org/3")
	env.Set("TMDB_TIMEOUT", "30s")
	env.Set("TMDB_MAX_RETRIES", "3")
	env.Set("TMDB_CACHE_TTL", "5m")
	env.Set("TMDB_LOG_LEVEL", "info")
	env.Set("TMDB_FORMAT", "json")

	return env
}

// TempFileContent creates a temporary file with specified content for testing.
func TempFileContent(t *testing.T, content string) (string, func()) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "tmdb-cli-test-*")
	require.NoError(t, err, "Failed to create temporary file")

	_, err = tmpFile.WriteString(content)
	require.NoError(t, err, "Failed to write content to temporary file")

	err = tmpFile.Close()
	require.NoError(t, err, "Failed to close temporary file")

	cleanup := func() {
		_ = os.Remove(tmpFile.Name())
	}

	return tmpFile.Name(), cleanup
}

// TempDir creates a temporary directory for testing.
func TempDir(t *testing.T, pattern string) (string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", pattern)
	require.NoError(t, err, "Failed to create temporary directory")

	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// TempConfigFile creates a temporary config file with specified YAML content.
func TempConfigFile(t *testing.T, yamlContent string) (string, func()) {
	t.Helper()

	return TempFileContent(t, yamlContent)
}

// CreateTempFileWithPermissions creates a temporary file with specific permissions.
func CreateTempFileWithPermissions(t *testing.T, content string, perm os.FileMode) (string, func()) {
	t.Helper()

	filePath, cleanup := TempFileContent(t, content)

	err := os.Chmod(filePath, perm)
	require.NoError(t, err, "Failed to set file permissions")

	return filePath, cleanup
}

// WithTempFile runs a test function with a temporary file.
func WithTempFile(t *testing.T, content string, fn func(string)) {
	t.Helper()

	filePath, cleanup := TempFileContent(t, content)
	defer cleanup()

	fn(filePath)
}

// WithTempDir runs a test function with a temporary directory.
func WithTempDir(t *testing.T, pattern string, fn func(string)) {
	t.Helper()

	dirPath, cleanup := TempDir(t, pattern)
	defer cleanup()

	fn(dirPath)
}

// AssertFileExists validates that a file exists at the specified path.
func AssertFileExists(t *testing.T, filePath string) {
	t.Helper()

	_, err := os.Stat(filePath)
	require.NoError(t, err, "File should exist at path: %s", filePath)
}

// AssertFileNotExists validates that a file does not exist at the specified path.
func AssertFileNotExists(t *testing.T, filePath string) {
	t.Helper()

	_, err := os.Stat(filePath)
	require.True(t, os.IsNotExist(err), "File should not exist at path: %s", filePath)
}

// ReadTempFile reads the contents of a temporary file for validation.
func ReadTempFile(t *testing.T, filePath string) string {
	t.Helper()

	content, err := os.ReadFile(filePath)
	require.NoError(t, err, "Failed to read temporary file")

	return string(content)
}

// AssertFileContains validates that a file contains the expected content.
func AssertFileContains(t *testing.T, filePath, expectedContent string) {
	t.Helper()

	content := ReadTempFile(t, filePath)
	require.Contains(t, content, expectedContent,
		"File should contain expected content")
}

// AssertFileEquals validates that a file has exactly the expected content.
func AssertFileEquals(t *testing.T, filePath, expectedContent string) {
	t.Helper()

	content := ReadTempFile(t, filePath)
	require.Equal(t, expectedContent, content,
		"File should have exactly the expected content")
}

// CleanupTempFiles removes temporary files created during testing.
// This is a fallback cleanup function for cases where individual cleanups fail.
func CleanupTempFiles(t *testing.T, filePaths ...string) {
	t.Helper()

	for _, filePath := range filePaths {
		if filePath != "" {
			_ = os.Remove(filePath)
		}
	}
}

// CommonEnvironmentVariables returns a map of common TMDB environment variables for testing.
func CommonEnvironmentVariables() map[string]string {
	return map[string]string{
		"TMDB_API_KEY":     "test-api-key-12345",
		"TMDB_BASE_URL":    "https://api.themoviedb.org/3",
		"TMDB_TIMEOUT":     "30s",
		"TMDB_MAX_RETRIES": "3",
		"TMDB_CACHE_TTL":   "5m",
		"TMDB_LOG_LEVEL":   "info",
		"TMDB_FORMAT":      "table",
	}
}

// TestEnvironmentVariables returns environment variables suitable for testing.
func TestEnvironmentVariables() map[string]string {
	return map[string]string{
		"TMDB_API_KEY":     "test-key",
		"TMDB_BASE_URL":    "https://test.api.com",
		"TMDB_TIMEOUT":     "10s",
		"TMDB_MAX_RETRIES": "1",
		"TMDB_CACHE_TTL":   "1m",
		"TMDB_LOG_LEVEL":   "debug",
		"TMDB_FORMAT":      "json",
	}
}

// InvalidEnvironmentVariables returns invalid environment variables for testing validation.
func InvalidEnvironmentVariables() map[string]string {
	return map[string]string{
		"TMDB_API_KEY":     "", // Empty API key
		"TMDB_TIMEOUT":     "invalid-duration",
		"TMDB_MAX_RETRIES": "not-a-number",
		"TMDB_CACHE_TTL":   "invalid-duration",
		"TMDB_LOG_LEVEL":   "invalid-level",
		"TMDB_FORMAT":      "invalid-format",
	}
}
