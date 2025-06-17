package unit

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected bool
	}{
		{"valid table format", "table", true},
		{"valid json format", "json", true},
		{"valid csv format", "csv", true},
		{"invalid format", "xml", false},
		{"empty format", "", false},
		{"case sensitive check", "TABLE", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := internal.ValidateFormat(test.format)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestFormatMovies(t *testing.T) {
	tests := []struct {
		name         string
		movies       []internal.Movie
		options      internal.FormatOptions
		expectedErr  string
		validateFunc func(t *testing.T, output string)
	}{
		{
			name:   "format empty movie list",
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
			name:   "format single movie as table",
			movies: fixtures.SingleMovie,
			options: internal.FormatOptions{
				Format:   "table",
				MaxWidth: 120,
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				assert.Contains(t, output, "Test Movie")
				assert.Contains(t, output, "2023")
				assert.Contains(t, output, "7.5")
			},
		},
		{
			name:   "format movies as JSON",
			movies: fixtures.SampleMovies,
			options: internal.FormatOptions{
				Format: "json",
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				helpers.ValidateJSONOutput(t, output, true)
				assert.Contains(t, output, "Spirited Away")
				assert.Contains(t, output, "The Matrix")
			},
		},
		{
			name:   "format movies as CSV",
			movies: fixtures.SampleMovies,
			options: internal.FormatOptions{
				Format: "csv",
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				lines := strings.Split(strings.TrimSpace(output), "\n")
				assert.GreaterOrEqual(t, len(lines), 2) // Header + at least one data row
				assert.Contains(t, lines[0], "ID,Title")
				assert.Contains(t, output, "Spirited Away")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := internal.FormatMovies(&buf, test.movies, test.options)

			if test.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.expectedErr)
				return
			}

			require.NoError(t, err)
			output := buf.String()
			test.validateFunc(t, output)
		})
	}
}

func TestFormatTVShows(t *testing.T) {
	tests := []struct {
		name         string
		tvShows      []internal.TVShow
		options      internal.FormatOptions
		expectedErr  string
		validateFunc func(t *testing.T, output string)
	}{
		{
			name:    "format empty TV shows list",
			tvShows: fixtures.EmptyTVShows,
			options: internal.FormatOptions{
				Format:   "table",
				MaxWidth: 80,
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				assert.Contains(t, output, "No TV shows found")
			},
		},
		{
			name:    "format single TV show as table",
			tvShows: fixtures.SingleTVShow,
			options: internal.FormatOptions{
				Format:   "table",
				MaxWidth: 120,
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				assert.Contains(t, output, "Test Show")
				assert.Contains(t, output, "2023")
				assert.Contains(t, output, "8.0")
			},
		},
		{
			name:    "format TV shows as JSON",
			tvShows: fixtures.SampleTVShows,
			options: internal.FormatOptions{
				Format: "json",
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				helpers.ValidateJSONOutput(t, output, false)
				assert.Contains(t, output, "Breaking Bad")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := internal.FormatTVShows(&buf, test.tvShows, test.options)

			if test.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.expectedErr)
				return
			}

			require.NoError(t, err)
			output := buf.String()
			test.validateFunc(t, output)
		})
	}
}

func TestFormatError(t *testing.T) {
	tests := []struct {
		name         string
		options      internal.FormatOptions
		expectedErr  string
		validateFunc func(t *testing.T, output string)
	}{
		{
			name: "invalid format",
			options: internal.FormatOptions{
				Format: "xml",
			},
			expectedErr: "unsupported format",
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				assert.Empty(t, output)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := internal.FormatMovies(&buf, fixtures.SampleMovies, test.options)

			if test.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.expectedErr)
			}

			output := buf.String()
			test.validateFunc(t, output)
		})
	}
}

func TestFormatOptionsValidation(t *testing.T) {
	tests := []struct {
		name         string
		options      internal.FormatOptions
		validateFunc func(t *testing.T, output string)
	}{
		{
			name: "default options validation",
			options: internal.FormatOptions{
				Format:   "table",
				MaxWidth: 80,
			},
			validateFunc: func(t *testing.T, _ string) {
				t.Helper()
				helpers.ValidateFormatOptionsEdgeCases(t, internal.FormatOptions{
					Format:   "table",
					MaxWidth: 80,
				})
			},
		},
		{
			name: "zero width validation",
			options: internal.FormatOptions{
				Format:   "json",
				MaxWidth: 0,
			},
			validateFunc: func(t *testing.T, _ string) {
				t.Helper()
				helpers.ValidateFormatOptionsEdgeCases(t, internal.FormatOptions{
					Format:   "json",
					MaxWidth: 0,
				})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := internal.FormatMovies(&buf, fixtures.SampleMovies, test.options)
			require.NoError(t, err)

			output := buf.String()
			test.validateFunc(t, output)
		})
	}
}
