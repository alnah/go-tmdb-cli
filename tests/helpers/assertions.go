// Package helpers provides testing utilities for assertions.
package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
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

// AssertMovieValid checks if a movie has all required fields and valid values.
func AssertMovieValid(t *testing.T, movie internal.Movie) {
	t.Helper()

	assert.Greater(t, movie.ID, 0, "Movie ID should be positive")
	assert.NotEmpty(t, movie.Title, "Movie title should not be empty")
	assert.GreaterOrEqual(t, movie.Rating, 0.0, "Rating should be non-negative")
	assert.LessOrEqual(t, movie.Rating, 10.0, "Rating should not exceed 10")
	assert.GreaterOrEqual(t, movie.Votes, 0, "Votes should be non-negative")
	assert.GreaterOrEqual(t, movie.Popularity, 0.0, "Popularity should be non-negative")
}

// AssertTVShowValid checks if a TV show has all required fields and valid values.
func AssertTVShowValid(t *testing.T, show internal.TVShow) {
	t.Helper()

	assert.Greater(t, show.ID, 0, "TV show ID should be positive")
	assert.NotEmpty(t, show.Name, "TV show name should not be empty")
	assert.GreaterOrEqual(t, show.Rating, 0.0, "Rating should be non-negative")
	assert.LessOrEqual(t, show.Rating, 10.0, "Rating should not exceed 10")
	assert.GreaterOrEqual(t, show.Votes, 0, "Votes should be non-negative")
	assert.GreaterOrEqual(t, show.Popularity, 0.0, "Popularity should be non-negative")
}

// AssertConfigValid checks if configuration has all required fields and valid values.
func AssertConfigValid(t *testing.T, config internal.Config) {
	t.Helper()

	assert.NotEmpty(t, config.APIKey, "API key should not be empty")
	assert.NotEmpty(t, config.BaseURL, "Base URL should not be empty")
	assert.Greater(t, config.Timeout, 0, "Timeout should be positive")
	assert.GreaterOrEqual(t, config.MaxRetries, 0, "Max retries should be non-negative")
	assert.Greater(t, config.CacheTTL, 0, "Cache TTL should be positive")
	assert.NotEmpty(t, config.LogLevel, "Log level should not be empty")
	assert.NotEmpty(t, config.Format, "Format should not be empty")
}

// AssertSearchOptionsValid checks if search options have valid values.
func AssertSearchOptionsValid(t *testing.T, opts internal.SearchOptions) {
	t.Helper()

	if opts.Page > 0 {
		assert.Greater(t, opts.Page, 0, "Page should be positive")
	}
	if opts.MinRating > 0 {
		assert.GreaterOrEqual(t, opts.MinRating, 0.0, "Min rating should be non-negative")
		assert.LessOrEqual(t, opts.MinRating, 10.0, "Min rating should not exceed 10")
	}
	if opts.MaxRating > 0 {
		assert.GreaterOrEqual(t, opts.MaxRating, 0.0, "Max rating should be non-negative")
		assert.LessOrEqual(t, opts.MaxRating, 10.0, "Max rating should not exceed 10")
	}
	if opts.MinRating > 0 && opts.MaxRating > 0 {
		assert.LessOrEqual(
			t,
			opts.MinRating,
			opts.MaxRating,
			"Min rating should not exceed max rating",
		)
	}
	if opts.MaxItems > 0 {
		assert.Greater(t, opts.MaxItems, 0, "Max items should be positive")
		assert.LessOrEqual(t, opts.MaxItems, 1000, "Max items should not exceed 1000")
	}
}

// AssertValidationError checks if a validation error contains expected message.
func AssertValidationError(t *testing.T, err error, expectedMessage string) {
	t.Helper()

	require.Error(t, err, "Expected validation error")
	assert.Contains(t, err.Error(), expectedMessage, "Error message should contain expected text")
}
