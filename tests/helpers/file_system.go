// Package helpers provides testing utilities for file system operations.
package helpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

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
