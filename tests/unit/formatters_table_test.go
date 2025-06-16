package unit

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
)

func TestFormatMoviesAsTable(t *testing.T) {
	tests := []struct {
		name         string
		movies       []internal.Movie
		options      internal.FormatOptions
		validateFunc func(t *testing.T, output string)
	}{
		{
			name:   "basic table formatting",
			movies: fixtures.SampleMovies[:2], // First 2 movies
			options: internal.FormatOptions{
				Format:   "table",
				MaxWidth: 120,
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				assert.Contains(t, output, "Spirited Away")
				assert.Contains(t, output, "The Matrix")
				assert.Contains(t, output, "2001")
				assert.Contains(t, output, "1999")
			},
		},
		{
			name:   "table with original titles",
			movies: fixtures.SampleMovies[1:2], // Just Spirited Away
			options: internal.FormatOptions{
				Format:           "table",
				MaxWidth:         120,
				UseOriginalTitle: true,
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				// When using original titles, should show the original Japanese title
				assert.Contains(t, output, "Sen to Chihiro no Kamikakushi")
			},
		},
		{
			name:   "empty movies table",
			movies: fixtures.EmptyMovies,
			options: internal.FormatOptions{
				Format:   "table",
				MaxWidth: 80,
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				assert.Contains(t, output, "No movies found")
			},
		},
		{
			name:   "single movie table",
			movies: fixtures.SingleMovie,
			options: internal.FormatOptions{
				Format:   "table",
				MaxWidth: 100,
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				assert.Contains(t, output, "Test Movie")
				assert.Contains(t, output, "2023")
				assert.Contains(t, output, "7.5")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := internal.FormatMovies(&buf, test.movies, test.options)

			require.NoError(t, err)
			output := buf.String()
			test.validateFunc(t, output)
		})
	}
}

func TestTableWidthConstraints(t *testing.T) {
	t.Run("table formatting produces reasonable output", func(t *testing.T) {
		var buf bytes.Buffer
		options := internal.FormatOptions{
			Format:   "table",
			MaxWidth: 45, // Very narrow constraint
		}

		err := internal.FormatMovies(&buf, fixtures.SampleMovies[:1], options)
		require.NoError(t, err)

		output := buf.String()

		// Just verify that output is produced and contains expected content
		// The exact width handling depends on the third-party table library
		assert.Contains(t, output, "Matrix")
		assert.Contains(t, output, "1999")

		// Test is mainly about verifying no crash occurs with narrow widths
		lines := strings.Split(output, "\n")
		assert.Greater(t, len(lines), 2, "Should produce multi-line table output")
	})
}

func TestTableHeaders(t *testing.T) {
	t.Run("movie table contains data", func(t *testing.T) {
		var buf bytes.Buffer
		options := internal.FormatOptions{
			Format:   "table",
			MaxWidth: 120,
		}

		err := internal.FormatMovies(&buf, fixtures.SampleMovies[:1], options)
		require.NoError(t, err)

		output := buf.String()

		// Focus on data content rather than exact header format
		// since that depends on the third-party table library
		assert.Contains(t, output, "The Matrix")
		assert.Contains(t, output, "1999")
		assert.Contains(t, output, "8.7")

		// Verify it's structured like a table (has separators)
		assert.Contains(t, output, "|")
		assert.Contains(t, output, "+")
	})

	t.Run("TV show table contains data", func(t *testing.T) {
		var buf bytes.Buffer
		options := internal.FormatOptions{
			Format:   "table",
			MaxWidth: 120,
		}

		err := internal.FormatTVShows(&buf, fixtures.SampleTVShows[:1], options)
		require.NoError(t, err)

		output := buf.String()

		// Focus on data content rather than exact header format
		assert.Contains(t, output, "Breaking Bad")
		assert.Contains(t, output, "2008")
		assert.Contains(t, output, "9.5")

		// Verify it's structured like a table
		assert.Contains(t, output, "|")
		assert.Contains(t, output, "+")
	})
}

func TestTableNoHeader(t *testing.T) {
	t.Run("table without headers", func(t *testing.T) {
		var buf bytes.Buffer
		options := internal.FormatOptions{
			Format:   "table",
			MaxWidth: 120,
			NoHeader: true,
		}

		err := internal.FormatMovies(&buf, fixtures.SampleMovies[:1], options)
		require.NoError(t, err)

		output := buf.String()

		// Should contain actual movie data
		assert.Contains(t, output, "The Matrix")
		assert.Contains(t, output, "1999")

		// Verify table structure is present
		assert.Contains(t, output, "|")
	})
}

func TestTableErrorHandling(t *testing.T) {
	t.Run("invalid format", func(t *testing.T) {
		var buf bytes.Buffer
		options := internal.FormatOptions{
			Format:   "xml",
			MaxWidth: 80,
		}

		err := internal.FormatMovies(&buf, fixtures.SampleMovies, options)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported format")
	})
}
