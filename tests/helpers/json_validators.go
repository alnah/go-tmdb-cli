// Package helpers provides testing utilities for the TMDB CLI test suite.
package helpers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// ValidateJSONOutput validates that output is valid JSON with expected structure for movies/TV shows.
func ValidateJSONOutput(t *testing.T, output string, expectMovies bool) {
	t.Helper()

	require.NotEmpty(t, output, "JSON output should not be empty")

	var result map[string]any
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err, "Output should be valid JSON")

	if expectMovies {
		require.Contains(t, result, "movies", "JSON should contain movies key")
		require.Contains(t, result, "count", "JSON should contain count key")

		movies, ok := result["movies"].([]any)
		require.True(t, ok, "movies should be an array")

		count, ok := result["count"].(float64)
		require.True(t, ok, "count should be a number")
		require.Equal(t, float64(len(movies)), count, "count should match movies array length")
	} else {
		require.Contains(t, result, "tv_shows", "JSON should contain tv_shows key")
		require.Contains(t, result, "count", "JSON should contain count key")

		shows, ok := result["tv_shows"].([]any)
		require.True(t, ok, "tv_shows should be an array")

		count, ok := result["count"].(float64)
		require.True(t, ok, "count should be a number")
		require.Equal(t, float64(len(shows)), count, "count should match tv_shows array length")
	}
}

// ValidateJSONStructure validates basic JSON structure without content verification.
func ValidateJSONStructure(t *testing.T, output string) map[string]any {
	t.Helper()

	require.NotEmpty(t, output, "JSON output should not be empty")

	var result map[string]any
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err, "Output should be valid JSON")

	return result
}

// ValidateJSONFormatting validates that JSON is properly formatted with indentation.
func ValidateJSONFormatting(t *testing.T, output string) {
	t.Helper()

	require.Contains(t, output, "\n", "JSON should contain newlines for formatting")
	require.Contains(t, output, "  ", "JSON should contain indentation")

	// Validate that it's properly formatted JSON by re-marshaling
	var data any
	err := json.Unmarshal([]byte(output), &data)
	require.NoError(t, err, "JSON should be valid")

	formatted, err := json.MarshalIndent(data, "", "  ")
	require.NoError(t, err, "Should be able to re-format JSON")

	// The lengths should be similar (allowing for different newline styles)
	require.InDelta(t, len(formatted), len(output), 10,
		"JSON should be properly formatted")
}

// ValidateJSONSearchResult validates search result JSON structure.
func ValidateJSONSearchResult(t *testing.T, output string, expectMovies bool) {
	t.Helper()

	result := ValidateJSONStructure(t, output)

	if expectMovies {
		require.Contains(t, result, "movies", "Search result should contain movies")
		require.Contains(t, result, "page", "Search result should contain page")
		require.Contains(t, result, "total_pages", "Search result should contain total_pages")
		require.Contains(t, result, "total_results", "Search result should contain total_results")

		// Validate data types
		_, ok := result["movies"].([]any)
		require.True(t, ok, "movies should be an array")

		_, ok = result["page"].(float64)
		require.True(t, ok, "page should be a number")

		_, ok = result["total_pages"].(float64)
		require.True(t, ok, "total_pages should be a number")

		_, ok = result["total_results"].(float64)
		require.True(t, ok, "total_results should be a number")
	} else {
		require.Contains(t, result, "tv_shows", "Search result should contain tv_shows")
		require.Contains(t, result, "page", "Search result should contain page")
		require.Contains(t, result, "total_pages", "Search result should contain total_pages")
		require.Contains(t, result, "total_results", "Search result should contain total_results")

		// Validate data types
		_, ok := result["tv_shows"].([]any)
		require.True(t, ok, "tv_shows should be an array")

		_, ok = result["page"].(float64)
		require.True(t, ok, "page should be a number")

		_, ok = result["total_pages"].(float64)
		require.True(t, ok, "total_pages should be a number")

		_, ok = result["total_results"].(float64)
		require.True(t, ok, "total_results should be a number")
	}
}

// ValidateJSONMovieFields validates that a movie object contains required fields.
func ValidateJSONMovieFields(t *testing.T, movieObj map[string]any) {
	t.Helper()

	requiredFields := []string{"id", "title", "year", "rating", "votes", "popularity", "genres", "overview"}

	for _, field := range requiredFields {
		require.Contains(t, movieObj, field, "Movie object should contain %s field", field)
	}

	// Validate specific field types
	_, ok := movieObj["id"].(float64)
	require.True(t, ok, "Movie ID should be a number")

	_, ok = movieObj["title"].(string)
	require.True(t, ok, "Movie title should be a string")

	_, ok = movieObj["rating"].(float64)
	require.True(t, ok, "Movie rating should be a number")

	_, ok = movieObj["votes"].(float64)
	require.True(t, ok, "Movie votes should be a number")
}

// ValidateJSONTVShowFields validates that a TV show object contains required fields.
func ValidateJSONTVShowFields(t *testing.T, showObj map[string]any) {
	t.Helper()

	requiredFields := []string{"id", "name", "year", "rating", "votes", "popularity", "genres", "overview"}

	for _, field := range requiredFields {
		require.Contains(t, showObj, field, "TV show object should contain %s field", field)
	}

	// Validate specific field types
	_, ok := showObj["id"].(float64)
	require.True(t, ok, "TV show ID should be a number")

	_, ok = showObj["name"].(string)
	require.True(t, ok, "TV show name should be a string")

	_, ok = showObj["rating"].(float64)
	require.True(t, ok, "TV show rating should be a number")

	_, ok = showObj["votes"].(float64)
	require.True(t, ok, "TV show votes should be a number")
}

// ValidateJSONArrayContents validates the contents of a JSON array.
func ValidateJSONArrayContents(t *testing.T, output string, expectedCount int, expectMovies bool) {
	t.Helper()

	result := ValidateJSONStructure(t, output)

	var items []any
	var ok bool

	if expectMovies {
		items, ok = result["movies"].([]any)
		require.True(t, ok, "movies should be an array")
	} else {
		items, ok = result["tv_shows"].([]any)
		require.True(t, ok, "tv_shows should be an array")
	}

	require.Equal(t, expectedCount, len(items),
		"JSON array should contain expected number of items")

	// Validate each item structure
	for i, item := range items {
		itemObj, ok := item.(map[string]any)
		require.True(t, ok, "Item %d should be an object", i)

		if expectMovies {
			ValidateJSONMovieFields(t, itemObj)
		} else {
			ValidateJSONTVShowFields(t, itemObj)
		}
	}
}

// ValidateJSONEmpty validates JSON output for empty collections.
func ValidateJSONEmpty(t *testing.T, output string, expectMovies bool) {
	t.Helper()

	result := ValidateJSONStructure(t, output)

	if expectMovies {
		movies, ok := result["movies"].([]any)
		require.True(t, ok, "movies should be an array")
		require.Empty(t, movies, "movies array should be empty")

		count, ok := result["count"].(float64)
		require.True(t, ok, "count should be a number")
		require.Equal(t, float64(0), count, "count should be 0 for empty collection")
	} else {
		shows, ok := result["tv_shows"].([]any)
		require.True(t, ok, "tv_shows should be an array")
		require.Empty(t, shows, "tv_shows array should be empty")

		count, ok := result["count"].(float64)
		require.True(t, ok, "count should be a number")
		require.Equal(t, float64(0), count, "count should be 0 for empty collection")
	}
}

// ValidateJSONContainsField validates that JSON output contains a specific field with expected value.
func ValidateJSONContainsField(t *testing.T, output string, fieldPath []string, expectedValue any) {
	t.Helper()

	result := ValidateJSONStructure(t, output)

	current := result
	for i, field := range fieldPath {
		if i == len(fieldPath)-1 {
			// Last field - check value
			require.Contains(t, current, field, "JSON should contain field %s", field)
			require.Equal(t, expectedValue, current[field],
				"Field %s should have expected value", field)
		} else {
			// Intermediate field - navigate deeper
			require.Contains(t, current, field, "JSON should contain field %s", field)
			next, ok := current[field].(map[string]any)
			require.True(t, ok, "Field %s should be an object", field)
			current = next
		}
	}
}

// ValidateJSONErrorResponse validates JSON error response structure.
func ValidateJSONErrorResponse(t *testing.T, output string) {
	t.Helper()

	result := ValidateJSONStructure(t, output)

	require.Contains(t, result, "error", "Error response should contain error field")
	require.Contains(t, result, "message", "Error response should contain message field")

	_, ok := result["error"].(string)
	require.True(t, ok, "error should be a string")

	_, ok = result["message"].(string)
	require.True(t, ok, "message should be a string")
}

// ValidateJSONConfigResponse validates JSON config response structure.
func ValidateJSONConfigResponse(t *testing.T, output string) {
	t.Helper()

	result := ValidateJSONStructure(t, output)

	configFields := []string{"api_key", "base_url", "timeout", "max_retries", "cache_ttl", "log_level", "format"}

	for _, field := range configFields {
		require.Contains(t, result, field, "Config response should contain %s field", field)
	}
}

// ExtractJSONField extracts a specific field value from JSON output.
func ExtractJSONField(t *testing.T, output string, fieldPath []string) any {
	t.Helper()

	result := ValidateJSONStructure(t, output)

	current := result
	for i, field := range fieldPath {
		require.Contains(t, current, field, "JSON should contain field %s", field)

		if i == len(fieldPath)-1 {
			return current[field]
		}

		next, ok := current[field].(map[string]any)
		require.True(t, ok, "Field %s should be an object", field)
		current = next
	}

	return nil
}
