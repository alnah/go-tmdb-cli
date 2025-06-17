// Package helpers provides testing utilities for the TMDB CLI test suite.
package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
)

// AssertValidConfig validates that a config struct has all required fields set correctly.
func AssertValidConfig(t *testing.T, config internal.Config) {
	t.Helper()

	assert.NotEmpty(t, config.APIKey, "API key should not be empty")
	assert.NotEmpty(t, config.BaseURL, "Base URL should not be empty")
	assert.True(t, config.Timeout > 0, "Timeout should be positive")
	assert.True(t, config.MaxRetries >= 0, "Max retries should be non-negative")
	assert.True(t, config.CacheTTL > 0, "Cache TTL should be positive")
	assert.NotEmpty(t, config.LogLevel, "Log level should not be empty")
	assert.NotEmpty(t, config.Format, "Format should not be empty")
}

// AssertValidMovie validates that a movie struct has all required fields set correctly.
// Alias for AssertMovieValid for backward compatibility.
func AssertValidMovie(t *testing.T, movie internal.Movie) {
	t.Helper()
	AssertMovieValid(t, movie)
}

// AssertValidTVShow validates that a TV show struct has all required fields set correctly.
// Alias for AssertTVShowValid for backward compatibility.
func AssertValidTVShow(t *testing.T, show internal.TVShow) {
	t.Helper()
	AssertTVShowValid(t, show)
}

// AssertValidSearchOptions validates that search options are properly configured.
func AssertValidSearchOptions(t *testing.T, opts internal.SearchOptions) {
	t.Helper()

	assert.True(
		t,
		opts.MaxItems >= 1 && opts.MaxItems <= 1000,
		"MaxItems should be between 1 and 1000",
	)
	assert.True(
		t,
		opts.MinRating >= 0.0 && opts.MinRating <= 10.0,
		"MinRating should be between 0 and 10",
	)
	assert.True(
		t,
		opts.MaxRating >= 0.0 && opts.MaxRating <= 10.0,
		"MaxRating should be between 0 and 10",
	)

	if opts.MinRating > 0 && opts.MaxRating > 0 {
		assert.LessOrEqual(
			t,
			opts.MinRating,
			opts.MaxRating,
			"MinRating should be less than or equal to MaxRating",
		)
	}

	// Note: SearchOptions doesn't have MinYear/MaxYear fields, only Year
	if opts.Year > 0 {
		assert.Greater(t, opts.Year, 1800, "Year should be reasonable")
	}
}

// AssertValidationError checks if a validation error contains expected message.
func AssertValidationError(t *testing.T, err error, expectedMessage string) {
	t.Helper()

	require.Error(t, err, "Expected validation error")
	assert.Contains(t, err.Error(), expectedMessage, "Error message should contain expected text")
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

// AssertMediaItemInterface validates that a struct properly implements MediaItem interface.
func AssertMediaItemInterface(t *testing.T, item internal.MediaItem) {
	t.Helper()

	assert.Greater(t, item.GetID(), 0, "MediaItem ID should be positive")
	assert.NotEmpty(t, item.GetTitle(), "MediaItem title should not be empty")
	assert.GreaterOrEqual(t, item.GetRating(), 0.0, "MediaItem rating should be non-negative")
	assert.GreaterOrEqual(t, item.GetVotes(), 0, "MediaItem votes should be non-negative")
	assert.GreaterOrEqual(
		t,
		item.GetPopularity(),
		0.0,
		"MediaItem popularity should be non-negative",
	)
	assert.NotNil(t, item.GetGenres(), "MediaItem genres should not be nil")
	assert.NotEmpty(t, item.GetOverview(), "MediaItem overview should not be empty")
}

// AssertEmptyCollection validates that a collection of items is empty and properly formatted.
func AssertEmptyCollection(t *testing.T, output string, format string) {
	t.Helper()

	assert.NotEmpty(t, output, "Output should not be empty even for empty collections")

	switch format {
	case FormatJSON:
		ValidateJSONStructure(t, output)
	case FormatCSV:
		// For CSV, empty collections should still have headers
		assert.Contains(t, output, "ID", "CSV should contain header even when empty")
	case FormatTable:
		assert.Contains(t, output, "No", "Table should indicate no items found")
	}
}

// AssertNonEmptyCollection validates that a collection contains the expected number of items.
func AssertNonEmptyCollection(t *testing.T, output string, format string, expectedCount int) {
	t.Helper()

	assert.NotEmpty(t, output, "Output should not be empty")

	switch format {
	case FormatJSON:
		result := ValidateJSONStructure(t, output)
		if movies, ok := result["movies"]; ok {
			if movieList, ok := movies.([]any); ok {
				assert.Len(
					t,
					movieList,
					expectedCount,
					"JSON should contain expected number of movies",
				)
			}
		} else if shows, ok := result["tv_shows"]; ok {
			if showList, ok := shows.([]any); ok {
				assert.Len(t, showList, expectedCount, "JSON should contain expected number of TV shows")
			}
		}
	case FormatCSV:
		records := ValidateCSVStructure(t, output)
		dataRecords := len(records) - 1 // Subtract header row
		assert.Equal(
			t,
			expectedCount,
			dataRecords,
			"CSV should contain expected number of data records",
		)
	}
}
