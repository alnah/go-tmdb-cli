// tests/helpers/assertions.go
package helpers

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/stretchr/testify/assert"
)

// AssertMovieValid validates that a movie has all required fields and sensible values
func AssertMovieValid(t *testing.T, movie internal.Movie) {
	t.Helper()

	assert.Greater(t, movie.ID, 0, "Movie ID should be positive")
	assert.NotEmpty(t, movie.Title, "Movie should have a title")
	assert.GreaterOrEqual(t, movie.Rating, 0.0, "Rating should be non-negative")
	assert.LessOrEqual(t, movie.Rating, 10.0, "Rating should not exceed 10")
	assert.GreaterOrEqual(t, movie.Votes, 0, "Votes should be non-negative")
	assert.GreaterOrEqual(t, movie.Popularity, 0.0, "Popularity should be non-negative")

	// Year validation (0 is acceptable for movies without release date)
	if movie.Year != 0 {
		assert.GreaterOrEqual(t, movie.Year, 1888, "Movie year should be after invention of cinema")
		assert.LessOrEqual(t, movie.Year, 2030, "Movie year should be reasonable")
	}
}

// AssertValidJSON validates that data is properly formatted JSON
func AssertValidJSON(t *testing.T, data []byte) {
	t.Helper()

	var result map[string]any
	err := json.Unmarshal(data, &result)
	assert.NoError(t, err, "Output should be valid JSON")
}

// AssertValidMovieJSON validates JSON contains expected movie structure
func AssertValidMovieJSON(t *testing.T, data []byte) {
	t.Helper()

	var result struct {
		Movies []internal.Movie `json:"movies"`
		Count  int              `json:"count"`
	}

	err := json.Unmarshal(data, &result)
	assert.NoError(t, err, "Should be valid movie JSON structure")
	assert.Equal(t, len(result.Movies), result.Count, "Count should match movies length")

	for i, movie := range result.Movies {
		AssertMovieValid(t, movie)
		assert.NotEmptyf(t, movie.Title, "Movie %d should have a title", i)
	}
}

// AssertSearchResultValid validates search result structure
func AssertSearchResultValid(t *testing.T, result internal.SearchResult) {
	t.Helper()

	assert.GreaterOrEqual(t, result.Page, 1, "Page should be at least 1")
	assert.GreaterOrEqual(t, result.TotalPages, 0, "Total pages should be non-negative")
	assert.GreaterOrEqual(t, result.TotalResults, 0, "Total results should be non-negative")
	assert.LessOrEqual(
		t,
		len(result.Movies),
		result.TotalResults,
		"Movies count should not exceed total results",
	)

	for _, movie := range result.Movies {
		AssertMovieValid(t, movie)
	}
}

// AssertConfigValid validates configuration structure
func AssertConfigValid(t *testing.T, config internal.Config) {
	t.Helper()

	assert.NotEmpty(t, config.APIKey, "API key should not be empty")
	assert.NotEmpty(t, config.BaseURL, "Base URL should not be empty")
	assert.Greater(t, config.Timeout, 0, "Timeout should be positive")
	assert.GreaterOrEqual(t, config.MaxRetries, 0, "Max retries should be non-negative")
	assert.Greater(t, config.CacheTTL, 0, "Cache TTL should be positive")
	assert.NotEmpty(t, config.LogLevel, "Log level should not be empty")
	assert.Contains(t, []string{"table", "json", "csv"}, config.Format, "Format should be valid")
}

// AssertGenreMapValid validates genre mapping correctness
func AssertGenreMapValid(t *testing.T, genreMap map[string]int) {
	t.Helper()

	// Test known genre mappings
	expectedGenres := map[string]int{
		"action":   28,
		"comedy":   35,
		"drama":    18,
		"horror":   27,
		"sci-fi":   878,
		"thriller": 53,
	}

	for name, expectedID := range expectedGenres {
		actualID, exists := genreMap[name]
		assert.True(t, exists, "Genre %s should exist in map", name)
		assert.Equal(t, expectedID, actualID, "Genre %s should have correct ID", name)
	}
}

// AssertFormatOptions validates format options
func AssertFormatOptionsValid(t *testing.T, options internal.FormatOptions) {
	t.Helper()

	validFormats := []string{"table", "json", "csv"}
	assert.Contains(t, validFormats, options.Format, "Format should be valid")
	assert.Greater(t, options.MaxWidth, 0, "MaxWidth should be positive")
}

// AssertErrorContains checks that error contains expected message
func AssertErrorContains(t *testing.T, err error, expectedMsg string) {
	t.Helper()

	assert.Error(t, err, "Expected an error")
	assert.Contains(t, err.Error(), expectedMsg, "Error should contain expected message")
}

// AssertNoErrorAndNotNil checks no error and result is not nil
func AssertNoErrorAndNotNil(t *testing.T, result any, err error) {
	t.Helper()

	assert.NoError(t, err, "Should not have error")
	assert.NotNil(t, result, "Result should not be nil")
}

// AssertSliceNotEmpty checks slice is not empty
func AssertSliceNotEmpty(t *testing.T, slice any, msgAndArgs ...any) {
	t.Helper()

	assert.NotEmpty(t, slice, msgAndArgs...)
}

// AssertStringNotEmptyOrDefault checks string is not empty unless it's an expected default
func AssertStringNotEmptyOrDefault(t *testing.T, value, defaultValue, fieldName string) {
	t.Helper()

	if value != defaultValue {
		assert.NotEmpty(t, value, "%s should not be empty when not default", fieldName)
	}
}

// AssertValidRating checks rating is within expected bounds
func AssertValidRating(t *testing.T, rating float64) {
	t.Helper()

	assert.GreaterOrEqual(t, rating, 0.0, "Rating should be >= 0")
	assert.LessOrEqual(t, rating, 10.0, "Rating should be <= 10")
}

// AssertValidYear checks year is reasonable
func AssertValidYear(t *testing.T, year int) {
	t.Helper()

	if year != 0 { // 0 is acceptable for unknown years
		assert.GreaterOrEqual(t, year, 1888, "Year should be after invention of cinema")
		assert.LessOrEqual(t, year, 2030, "Year should be reasonable")
	}
}

// AssertValidPopularity checks popularity is non-negative
func AssertValidPopularity(t *testing.T, popularity float64) {
	t.Helper()

	assert.GreaterOrEqual(t, popularity, 0.0, "Popularity should be non-negative")
}

// AssertValidVoteCount checks vote count is non-negative
func AssertValidVoteCount(t *testing.T, votes int) {
	t.Helper()

	assert.GreaterOrEqual(t, votes, 0, "Vote count should be non-negative")
}

// AssertTableOutput validates table output contains expected elements
func AssertTableOutput(t *testing.T, output string) {
	t.Helper()

	// Should contain table headers
	assert.Contains(t, output, "#", "Table should contain row numbers")
	assert.Contains(t, output, "Title", "Table should contain Title header")
	assert.Contains(t, output, "Year", "Table should contain Year header")
	assert.Contains(t, output, "Rating", "Table should contain Rating header")

	// Should not be empty
	assert.NotEmpty(t, output, "Table output should not be empty")

	// Should have multiple lines for actual table
	lines := len(strings.Split(output, "\n"))
	assert.Greater(t, lines, 2, "Table should have header and content lines")
}

// AssertCSVOutput validates CSV output format
func AssertCSVOutput(t *testing.T, output string, expectHeader bool) {
	t.Helper()

	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.NotEmpty(t, lines, "CSV should have content")

	if expectHeader {
		headerLine := lines[0]
		assert.Contains(t, headerLine, "ID", "CSV header should contain ID")
		assert.Contains(t, headerLine, "Title", "CSV header should contain Title")
		assert.Contains(t, headerLine, "Year", "CSV header should contain Year")
	}

	// Check that lines have consistent comma count
	if len(lines) > 1 {
		firstLineCommas := strings.Count(lines[0], ",")
		for i, line := range lines[1:] {
			commas := strings.Count(line, ",")
			assert.Equal(
				t,
				firstLineCommas,
				commas,
				"Line %d should have same comma count as header",
				i+1,
			)
		}
	}
}

// AssertMovieListSorted checks if movies are sorted by specified criteria
func AssertMovieListSorted(t *testing.T, movies []internal.Movie, sortBy string, ascending bool) {
	t.Helper()

	if len(movies) < 2 {
		return // Cannot validate sorting with less than 2 items
	}

	for i := 1; i < len(movies); i++ {
		prev := movies[i-1]
		curr := movies[i]

		switch sortBy {
		case "rating":
			if ascending {
				assert.LessOrEqual(
					t,
					prev.Rating,
					curr.Rating,
					"Movies should be sorted by rating ascending",
				)
			} else {
				assert.GreaterOrEqual(t, prev.Rating, curr.Rating, "Movies should be sorted by rating descending")
			}
		case "year":
			if ascending {
				assert.LessOrEqual(
					t,
					prev.Year,
					curr.Year,
					"Movies should be sorted by year ascending",
				)
			} else {
				assert.GreaterOrEqual(t, prev.Year, curr.Year, "Movies should be sorted by year descending")
			}
		case "popularity":
			if ascending {
				assert.LessOrEqual(
					t,
					prev.Popularity,
					curr.Popularity,
					"Movies should be sorted by popularity ascending",
				)
			} else {
				assert.GreaterOrEqual(t, prev.Popularity, curr.Popularity, "Movies should be sorted by popularity descending")
			}
		case "title":
			if ascending {
				assert.LessOrEqual(
					t,
					prev.Title,
					curr.Title,
					"Movies should be sorted by title ascending",
				)
			} else {
				assert.GreaterOrEqual(t, prev.Title, curr.Title, "Movies should be sorted by title descending")
			}
		}
	}
}

// Helper functions for common patterns
func CreateMovieSlice(count int, baseName string) []internal.Movie {
	movies := make([]internal.Movie, count)
	for i := range count {
		movies[i] = internal.Movie{
			ID:         i + 1,
			Title:      fmt.Sprintf("%s %d", baseName, i+1),
			Year:       2020 + (i % 5),
			Rating:     5.0 + float64(i%6),
			Votes:      1000 + i*100,
			Popularity: float64(50 + i*10),
			Genres:     "Action, Drama",
			Language:   "en",
		}
	}
	return movies
}

func CreateTestMovie(id int, title string) internal.Movie {
	return internal.Movie{
		ID:          id,
		Title:       title,
		Year:        2023,
		Rating:      7.5,
		Votes:       1500,
		Popularity:  85.5,
		Genres:      "Action, Thriller",
		Overview:    "A test movie for testing purposes",
		Language:    "en",
		Adult:       false,
		ReleaseDate: "2023-06-15",
	}
}
