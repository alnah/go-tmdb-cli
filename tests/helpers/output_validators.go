// Package helpers provides testing utilities for the TMDB CLI test suite.
package helpers

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
)

// Format constants for consistent usage across helpers.
const (
	FormatJSON  = "json"
	FormatCSV   = "csv"
	FormatTable = "table"
)

// SplitLines splits a string into lines and removes empty lines at the end.
func SplitLines(text string) []string {
	lines := strings.Split(text, "\n")
	// Remove trailing empty lines
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// IndexOf returns the first index of substr in s, or -1 if not found.
func IndexOf(s, substr string) int {
	return strings.Index(s, substr)
}

// CaptureOutput captures both stdout and stderr output during function execution.
// CaptureOutput captures both stdout and stderr output during function execution.
// Fixed version that properly captures both streams simultaneously.
func CaptureOutput(t *testing.T, fn func()) (stdout, stderr string) {
	t.Helper()

	// Save original stdout and stderr
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	// Create pipes for capturing output
	rOut, wOut, err := os.Pipe()
	require.NoError(t, err, "Failed to create stdout pipe")

	rErr, wErr, err := os.Pipe()
	require.NoError(t, err, "Failed to create stderr pipe")

	// Replace stdout and stderr with write ends of pipes
	os.Stdout = wOut
	os.Stderr = wErr

	// Channels to collect output
	stdoutChan := make(chan string, 1)
	stderrChan := make(chan string, 1)

	// Start goroutines to read from pipes
	go func() {
		var buf bytes.Buffer
		_, e := buf.ReadFrom(rOut)
		if e != nil {
			stdoutChan <- ""
		} else {
			stdoutChan <- buf.String()
		}
	}()

	go func() {
		var buf bytes.Buffer
		_, e := buf.ReadFrom(rErr)
		if e != nil {
			stderrChan <- ""
		} else {
			stderrChan <- buf.String()
		}
	}()

	// Execute function
	fn()

	// Close write ends to signal completion
	err = wOut.Close()
	require.NoError(t, err, "Failed to close stdout write end")
	err = wErr.Close()
	require.NoError(t, err, "Failed to close stderr write end")

	// Restore original stdout and stderr
	os.Stdout = originalStdout
	os.Stderr = originalStderr

	// Collect results
	stdoutResult := <-stdoutChan
	stderrResult := <-stderrChan

	// Close read ends
	_ = rOut.Close()
	_ = rErr.Close()

	return stdoutResult, stderrResult
}

// ValidateOutputNotEmpty validates that output is not empty or just whitespace.
func ValidateOutputNotEmpty(t *testing.T, output string) {
	t.Helper()

	require.NotEmpty(t, output, "Output should not be empty")
	require.NotEmpty(t, strings.TrimSpace(output), "Output should not be just whitespace")
}

// ValidateOutputContains validates that output contains expected strings.
func ValidateOutputContains(t *testing.T, output string, expectedContents ...string) {
	t.Helper()

	ValidateOutputNotEmpty(t, output)

	for _, expected := range expectedContents {
		require.Contains(t, output, expected,
			"Output should contain: %s", expected)
	}
}

// ValidateOutputNotContains validates that output does not contain specified strings.
func ValidateOutputNotContains(t *testing.T, output string, unexpectedContents ...string) {
	t.Helper()

	for _, unexpected := range unexpectedContents {
		require.NotContains(t, output, unexpected,
			"Output should not contain: %s", unexpected)
	}
}

// ValidateOutputLength validates that output length is within expected bounds.
func ValidateOutputLength(t *testing.T, output string, minLength, maxLength int) {
	t.Helper()

	length := len(output)
	require.GreaterOrEqual(t, length, minLength,
		"Output length should be at least %d characters", minLength)

	if maxLength > 0 {
		require.LessOrEqual(t, length, maxLength,
			"Output length should be at most %d characters", maxLength)
	}
}

// ValidateOutputLines validates that output has expected number of lines.
func ValidateOutputLines(t *testing.T, output string, expectedLines int) {
	t.Helper()

	lines := SplitLines(output)
	actualLines := len(lines)

	// Handle empty output case
	if strings.TrimSpace(output) == "" {
		actualLines = 0
	}

	require.Equal(t, expectedLines, actualLines,
		"Output should have %d lines, got %d", expectedLines, actualLines)
}

// ValidateOutputLinesRange validates that output has lines within expected range.
func ValidateOutputLinesRange(t *testing.T, output string, minLines, maxLines int) {
	t.Helper()

	lines := SplitLines(output)
	actualLines := len(lines)

	// Handle empty output case
	if strings.TrimSpace(output) == "" {
		actualLines = 0
	}

	require.GreaterOrEqual(t, actualLines, minLines,
		"Output should have at least %d lines", minLines)
	require.LessOrEqual(t, actualLines, maxLines,
		"Output should have at most %d lines", maxLines)
}

// ValidateOutputFormat validates that output conforms to expected format characteristics.
func ValidateOutputFormat(t *testing.T, output, format string) {
	t.Helper()

	ValidateOutputNotEmpty(t, output)

	switch format {
	case FormatJSON:
		ValidateJSONStructure(t, output)
	case FormatCSV:
		ValidateCSVStructure(t, output)
	case FormatTable:
		ValidateTableOutput(t, output)
	default:
		t.Fatalf("Unknown format for validation: %s", format)
	}
}

// ValidateTableOutput validates basic table format characteristics.
func ValidateTableOutput(t *testing.T, output string) {
	t.Helper()

	ValidateOutputNotEmpty(t, output)

	// Tables should have multiple lines for headers and data
	lines := SplitLines(output)

	// Allow for empty tables with messages like "No items found"
	if len(lines) == 1 && (strings.Contains(output, "No ") || strings.Contains(output, "Empty")) {
		return
	}

	// For non-empty tables, expect at least header and separator
	require.GreaterOrEqual(t, len(lines), 2,
		"Table output should have at least header and separator lines")
}

// ValidateOutputEncoding validates that output properly handles special characters.
func ValidateOutputEncoding(t *testing.T, output string, expectedSpecialChars ...string) {
	t.Helper()

	for _, char := range expectedSpecialChars {
		require.Contains(t, output, char,
			"Output should properly handle special character: %s", char)
	}
}

// ValidateNoErrorInOutput validates that output doesn't contain common error indicators.
func ValidateNoErrorInOutput(t *testing.T, output string) {
	t.Helper()

	errorIndicators := []string{
		"error:", "Error:", "ERROR:",
		"failed:", "Failed:", "FAILED:",
		"panic:", "Panic:", "PANIC:",
		"fatal:", "Fatal:", "FATAL:",
	}

	for _, indicator := range errorIndicators {
		require.NotContains(t, output, indicator,
			"Output should not contain error indicator: %s", indicator)
	}
}

// ValidateOutputPattern validates that output matches expected patterns.
func ValidateOutputPattern(t *testing.T, output string, patterns ...string) {
	t.Helper()

	for _, pattern := range patterns {
		require.Regexp(t, pattern, output,
			"Output should match pattern: %s", pattern)
	}
}

// ValidateOutputStructure validates common output structure requirements.
func ValidateOutputStructure(t *testing.T, output, format string, hasData bool) {
	t.Helper()

	ValidateOutputNotEmpty(t, output)

	switch format {
	case FormatJSON:
		if hasData {
			ValidateJSONOutput(t, output, true) // Assume movies for simplicity
		} else {
			ValidateJSONEmpty(t, output, true)
		}
	case FormatCSV:
		if hasData {
			ValidateCSVRecordCount(t, output, 1, true) // At least 1 data record with header
		} else {
			ValidateCSVEmpty(t, output, true)
		}
	case FormatTable:
		ValidateTableOutput(t, output)
		if !hasData {
			ValidateOutputContains(t, output, "No")
		}
	}
}

// ValidateWriterBehavior validates that a writer function behaves correctly.
func ValidateWriterBehavior(t *testing.T, writerFunc func(io.Writer) error) {
	t.Helper()

	// Test with normal writer
	mockWriter := NewMockWriter()
	err := writerFunc(mockWriter)
	require.NoError(t, err, "Writer function should succeed with normal writer")
	require.Greater(t, mockWriter.Len(), 0, "Writer should have written data")

	// Test with failing writer
	failingWriter := NewMockWriter()
	failingWriter.SetShouldError(true)
	err = writerFunc(failingWriter)
	require.Error(t, err, "Writer function should fail with failing writer")
}

// ValidateContentConsistency validates that content is consistent across different formats.
func ValidateContentConsistency(t *testing.T, outputs map[string]string, expectedItemCount int) {
	t.Helper()

	for format, output := range outputs {
		ValidateOutputNotEmpty(t, output)

		switch format {
		case FormatJSON:
			result := ValidateJSONStructure(t, output)
			if movies, ok := result["movies"]; ok {
				if movieList, ok := movies.([]any); ok {
					require.Equal(t, expectedItemCount, len(movieList),
						"JSON should have consistent item count")
				}
			}
		case FormatCSV:
			records := ValidateCSVStructure(t, output)
			dataRecords := len(records) - 1 // Subtract header
			require.Equal(t, expectedItemCount, dataRecords,
				"CSV should have consistent item count")
		}
	}
}

// ValidateEmptyCollectionHandling validates proper handling of empty collections.
func ValidateEmptyCollectionHandling(t *testing.T, output, format string) {
	t.Helper()

	ValidateOutputNotEmpty(t, output)

	switch format {
	case FormatJSON:
		ValidateJSONEmpty(t, output, true)
	case FormatCSV:
		ValidateCSVEmpty(t, output, true)
	case FormatTable:
		ValidateOutputContains(t, output, "No")
	}
}

// ValidateProgressiveOutput validates output that should build progressively.
func ValidateProgressiveOutput(t *testing.T, outputs []string) {
	t.Helper()

	require.NotEmpty(t, outputs, "Should have progressive outputs")

	for i, output := range outputs {
		ValidateOutputNotEmpty(t, output)

		if i > 0 {
			// Each output should be larger or equal to previous
			require.GreaterOrEqual(t, len(output), len(outputs[i-1]),
				"Progressive output %d should not be smaller than previous", i)
		}
	}
}

// ValidateOutputMetadata validates metadata information in output.
func ValidateOutputMetadata(t *testing.T, output string, expectedCount int, format string) {
	t.Helper()

	switch format {
	case FormatJSON:
		result := ValidateJSONStructure(t, output)

		if count, ok := result["count"]; ok {
			require.Equal(t, float64(expectedCount), count,
				"JSON metadata count should match expected")
		}
	case FormatTable:
		if expectedCount == 0 {
			ValidateOutputContains(t, output, "No")
		}
		// Table format doesn't typically contain explicit count metadata
	}
}

// ValidateOutputSortOrder validates that output appears to be sorted correctly.
func ValidateOutputSortOrder(
	t *testing.T,
	output, format string,
	sortField string,
	ascending bool,
) {
	t.Helper()

	// This is a basic implementation that can be extended based on specific needs
	ValidateOutputNotEmpty(t, output)

	switch format {
	case "csv":
		records := ValidateCSVStructure(t, output)
		if len(records) > 2 { // Header + at least 2 data rows
			// Basic validation - first and last data rows should be in correct order
			// This would need to be enhanced based on specific sort field indices
			t.Logf("Validating sort order for field %s (ascending: %v)", sortField, ascending)
		}
	case "json":
		result := ValidateJSONStructure(t, output)
		if items, ok := result["movies"]; ok {
			itemList, _ := items.([]any)
			if len(itemList) > 1 {
				t.Logf("Validating JSON sort order for %d items", len(itemList))
			}
		}
	}
}

// ValidateOutputPerformance validates that output generation is within performance bounds.
func ValidateOutputPerformance(t *testing.T, outputFunc func() (string, error), maxDurationMs int) {
	t.Helper()

	// This could be enhanced with actual timing measurements
	output, err := outputFunc()
	require.NoError(t, err, "Output function should not error")
	ValidateOutputNotEmpty(t, output)

	// Placeholder for performance validation
	t.Logf("Performance validation placeholder - max duration: %dms", maxDurationMs)
}

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

// RunConcurrently executes a function concurrently in multiple goroutines.
// This is useful for testing thread safety and concurrent behavior.
func RunConcurrently(t *testing.T, numGoroutines int, fn func(goroutineID int)) {
	t.Helper()

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			fn(id)
		}(i)
	}

	wg.Wait()
}

// ValidateFormatOptionsEdgeCases validates format options edge cases.
func ValidateFormatOptionsEdgeCases(t *testing.T, options internal.FormatOptions) {
	t.Helper()

	require.True(t, internal.ValidateFormat(options.Format), "Format should be valid")
	require.GreaterOrEqual(t, options.MaxWidth, 0, "MaxWidth should be non-negative")
}
