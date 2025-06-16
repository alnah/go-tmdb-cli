// Package helpers provides testing utilities for assertions.
package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// LoggerTestCase represents a test case structure for logger testing.
type LoggerTestCase struct {
	Name           string
	LogLevel       string
	Message        string
	Args           []any
	ExpectOutput   bool
	ExpectedPrefix string
	ExpectedInLog  []string
	NotInLog       []string
}

// AssertLogOutput validates that log output contains expected content.
func AssertLogOutput(t *testing.T, output string, testCase LoggerTestCase) {
	t.Helper()

	if testCase.ExpectOutput {
		require.NotEmpty(t, output, "Expected log output but got none for test: %s", testCase.Name)

		if testCase.ExpectedPrefix != "" {
			require.Contains(
				t,
				output,
				testCase.ExpectedPrefix,
				"Expected prefix '%s' not found in output for test: %s",
				testCase.ExpectedPrefix,
				testCase.Name,
			)
		}

		for _, expected := range testCase.ExpectedInLog {
			require.Contains(t, output, expected,
				"Expected content '%s' not found in output for test: %s", expected, testCase.Name)
		}

		for _, notExpected := range testCase.NotInLog {
			require.NotContains(t, output, notExpected,
				"Unexpected content '%s' found in output for test: %s", notExpected, testCase.Name)
		}
	} else {
		require.Empty(t, output, "Expected no log output but got: %s for test: %s", output, testCase.Name)
	}
}
