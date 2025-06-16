package unit

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
)

func TestMovieFormatterBasicFields(t *testing.T) {
	movie := fixtures.SampleMovies[0] // The Matrix
	formatter := internal.MovieFormatter{Movie: movie}

	t.Run("GetTitle", func(t *testing.T) {
		title := formatter.GetTitle(false)
		assert.Equal(t, "The Matrix", title)

		originalTitle := formatter.GetTitle(true)
		assert.Equal(t, "The Matrix", originalTitle) // Same for this movie
	})

	t.Run("GetOriginalTitle", func(t *testing.T) {
		originalTitle := formatter.GetOriginalTitle()
		assert.Equal(t, "The Matrix", originalTitle)
	})

	t.Run("GetFields", func(t *testing.T) {
		fields := formatter.GetFields()
		expectedFields := []string{
			"2",                       // ID
			"The Matrix",              // Title
			"The Matrix",              // Original Title
			"1999",                    // Year
			"8.7",                     // Rating
			"1847304",                 // Votes
			"78.45",                   // Popularity
			"Action, Science Fiction", // Genres
			"A computer hacker learns from mysterious rebels about the true nature of his reality.", // Overview
			"en",    // Language
			"false", // Adult
		}

		require.Len(t, fields, len(expectedFields), "Should have correct number of fields")
		for i, expected := range expectedFields {
			assert.Equal(t, expected, fields[i], "Field %d should match", i)
		}
	})

	t.Run("GetEmptyMessage", func(t *testing.T) {
		message := formatter.GetEmptyMessage()
		assert.Equal(t, "No movies found.", message)
	})

	t.Run("GetHeaderName", func(t *testing.T) {
		headerName := formatter.GetHeaderName()
		assert.Equal(t, "Title", headerName)
	})
}

func TestMovieFormatterWithDifferentTitles(t *testing.T) {
	movie := fixtures.SampleMovies[1] // Spirited Away
	formatter := internal.MovieFormatter{Movie: movie}

	t.Run("GetTitle shows original when requested", func(t *testing.T) {
		title := formatter.GetTitle(true)
		// Should show original title with English in parentheses if space allows
		assert.Contains(t, title, "Sen to Chihiro no Kamikakushi")
	})

	t.Run("GetTitle shows English when not requested", func(t *testing.T) {
		title := formatter.GetTitle(false)
		// Should show English title with original in parentheses if space allows
		assert.Contains(t, title, "Spirited Away")
	})
}

func TestMovieFormatterEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		movie    internal.Movie
		testFunc func(t *testing.T, formatter internal.MovieFormatter)
	}{
		{
			name: "empty titles",
			movie: internal.Movie{
				ID:            1,
				Title:         "",
				OriginalTitle: "",
			},
			testFunc: func(t *testing.T, formatter internal.MovieFormatter) {
				t.Helper()
				title := formatter.GetTitle(false)
				assert.Equal(t, "", title)

				originalTitle := formatter.GetOriginalTitle()
				assert.Equal(t, "", originalTitle)
			},
		},
		{
			name: "only original title",
			movie: internal.Movie{
				ID:            2,
				Title:         "",
				OriginalTitle: "Original Only",
			},
			testFunc: func(t *testing.T, formatter internal.MovieFormatter) {
				t.Helper()
				title := formatter.GetTitle(true)
				assert.Equal(t, "Original Only", title)

				title = formatter.GetTitle(false)
				assert.Equal(t, "", title) // Empty title
			},
		},
		{
			name: "only regular title",
			movie: internal.Movie{
				ID:            3,
				Title:         "Regular Only",
				OriginalTitle: "",
			},
			testFunc: func(t *testing.T, formatter internal.MovieFormatter) {
				t.Helper()
				title := formatter.GetTitle(false)
				assert.Equal(t, "Regular Only", title)

				title = formatter.GetTitle(true)
				assert.Equal(t, "Regular Only", title) // Falls back to regular title
			},
		},
		{
			name: "very long titles",
			movie: internal.Movie{
				ID:            4,
				Title:         "This Is An Extremely Long Movie Title That Should Be Truncated When Displayed",
				OriginalTitle: "This Is An Even Longer Original Movie Title That Should Also Be Truncated",
			},
			testFunc: func(t *testing.T, formatter internal.MovieFormatter) {
				t.Helper()
				title := formatter.GetTitle(false)
				// Should be truncated
				assert.Contains(t, title, "...")

				originalTitle := formatter.GetTitle(true)
				// Should also be truncated
				assert.Contains(t, originalTitle, "...")
			},
		},
		{
			name: "unicode characters",
			movie: internal.Movie{
				ID:            5,
				Title:         "Movie 🎬 with émojis",
				OriginalTitle: "Película con émojis 🎬",
			},
			testFunc: func(t *testing.T, formatter internal.MovieFormatter) {
				t.Helper()
				title := formatter.GetTitle(false)
				assert.Contains(t, title, "🎬")
				assert.Contains(t, title, "émojis")

				originalTitle := formatter.GetTitle(true)
				assert.Contains(t, originalTitle, "🎬")
				assert.Contains(t, originalTitle, "Película")
			},
		},
		{
			name: "zero and negative values",
			movie: internal.Movie{
				ID:         6,
				Title:      "Edge Case Movie",
				Year:       0,
				Rating:     0.0,
				Votes:      0,
				Popularity: 0.0,
			},
			testFunc: func(t *testing.T, formatter internal.MovieFormatter) {
				t.Helper()
				fields := formatter.GetFields()
				assert.Equal(t, "0", fields[3])    // Year
				assert.Equal(t, "0.0", fields[4])  // Rating
				assert.Equal(t, "0", fields[5])    // Votes
				assert.Equal(t, "0.00", fields[6]) // Popularity
			},
		},
		{
			name: "adult content flag",
			movie: internal.Movie{
				ID:    7,
				Title: "Adult Movie",
				Adult: true,
			},
			testFunc: func(t *testing.T, formatter internal.MovieFormatter) {
				t.Helper()
				fields := formatter.GetFields()
				assert.Equal(t, "true", fields[10]) // Adult field
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := internal.MovieFormatter{Movie: tt.movie}
			tt.testFunc(t, formatter)
		})
	}
}

func TestMovieFormatterFieldsConsistency(t *testing.T) {
	// Test that GetFields always returns the same number of fields
	testMovies := []internal.Movie{
		fixtures.SampleMovies[0],
		fixtures.EdgeCaseMovies[0],
		{ID: 999}, // Minimal movie
	}

	expectedFieldCount := 11 // As defined in the CSV headers

	for i, movie := range testMovies {
		t.Run(fmt.Sprintf("movie_%d", i), func(t *testing.T) {
			formatter := internal.MovieFormatter{Movie: movie}
			fields := formatter.GetFields()
			assert.Len(t, fields, expectedFieldCount, "All movies should have same field count")
		})
	}
}

func TestMovieFormatterFloatFormatting(t *testing.T) {
	movie := internal.Movie{
		ID:         1,
		Title:      "Test Movie",
		Rating:     8.123456,  // Many decimal places
		Popularity: 75.987654, // Many decimal places
	}

	formatter := internal.MovieFormatter{Movie: movie}
	fields := formatter.GetFields()

	// Rating should be formatted to 1 decimal place
	assert.Equal(t, "8.1", fields[4])

	// Popularity should be formatted to 2 decimal places
	assert.Equal(t, "75.99", fields[6])
}

func TestMovieFormatterSpecialCharacters(t *testing.T) {
	movie := internal.Movie{
		ID:       1,
		Title:    "Movie with \"quotes\" and 'apostrophes'",
		Genres:   "Action, Comedy",
		Overview: "Overview with\nnewlines\tand\ttabs & special chars <html>",
		Language: "en-US",
	}

	formatter := internal.MovieFormatter{Movie: movie}
	fields := formatter.GetFields()

	// Fields should preserve special characters
	assert.Contains(t, fields[1], "quotes")   // Title
	assert.Contains(t, fields[7], "Action")   // Genres
	assert.Contains(t, fields[8], "newlines") // Overview
	assert.Contains(t, fields[8], "<html>")   // HTML in overview
	assert.Equal(t, "en-US", fields[9])       // Language
}

func TestGetMovieHeaders(t *testing.T) {
	// This tests the unexported function through the CSV formatter
	movie := fixtures.SampleMovies[0]

	var buf bytes.Buffer
	err := internal.FormatMovies(&buf, []internal.Movie{movie}, internal.FormatOptions{
		Format:   internal.FormatCSV,
		NoHeader: false,
	})
	require.NoError(t, err)

	output := buf.String()
	lines := strings.Split(output, "\n")
	require.NotEmpty(t, lines)

	// First line should contain headers
	headerLine := lines[0]
	expectedHeaders := []string{
		"ID", "Title", "Original Title", "Year", "Rating",
		"Votes", "Popularity", "Genres", "Overview", "Language", "Adult",
	}

	for _, header := range expectedHeaders {
		assert.Contains(t, headerLine, header, "Header line should contain: %s", header)
	}
}

func TestMovieFormatterInterface(t *testing.T) {
	movie := fixtures.SampleMovies[0]
	formatter := internal.MovieFormatter{Movie: movie}

	// Test that MovieFormatter implements MediaFormatter interface
	var mediaFormatter internal.MediaFormatter = formatter

	// Test interface methods
	title := mediaFormatter.GetTitle(false)
	assert.NotEmpty(t, title)

	originalTitle := mediaFormatter.GetOriginalTitle()
	assert.NotEmpty(t, originalTitle)

	fields := mediaFormatter.GetFields()
	assert.NotEmpty(t, fields)

	emptyMessage := mediaFormatter.GetEmptyMessage()
	assert.Equal(t, "No movies found.", emptyMessage)

	headerName := mediaFormatter.GetHeaderName()
	assert.Equal(t, "Title", headerName)
}
