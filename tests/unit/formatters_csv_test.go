package unit

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestFormatMoviesCSV(t *testing.T) {
	tests := []struct {
		name           string
		movies         []internal.Movie
		options        internal.FormatOptions
		expectedHeader []string
		validateFunc   func(t *testing.T, output string)
	}{
		{
			name:   "format single movie with header",
			movies: fixtures.SingleMovie,
			options: internal.FormatOptions{
				Format:   internal.FormatCSV,
				NoHeader: false,
			},
			expectedHeader: []string{
				"ID",
				"Title",
				"Original Title",
				"Year",
				"Rating",
				"Votes",
				"Popularity",
				"Genres",
				"Overview",
				"Language",
				"Adult",
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				helpers.ValidateCSVRecordCount(t, output, 1, true)
				helpers.ValidateCSVColumnCount(t, output)
			},
		},
		{
			name:   "format multiple movies with header",
			movies: fixtures.SampleMovies[:3],
			options: internal.FormatOptions{
				Format:   internal.FormatCSV,
				NoHeader: false,
			},
			expectedHeader: []string{
				"ID",
				"Title",
				"Original Title",
				"Year",
				"Rating",
				"Votes",
				"Popularity",
				"Genres",
				"Overview",
				"Language",
				"Adult",
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				helpers.ValidateCSVRecordCount(t, output, 3, true)
				helpers.ValidateCSVColumnCount(t, output)

				// Validate specific movie data
				helpers.ValidateCSVFieldContent(t, output, 1, 1, "The Matrix")
				helpers.ValidateCSVFieldContent(t, output, 2, 1, "Spirited Away")
			},
		},
		{
			name:   "format movies without header",
			movies: fixtures.SampleMovies[:2],
			options: internal.FormatOptions{
				Format:   internal.FormatCSV,
				NoHeader: true,
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				helpers.ValidateCSVRecordCount(t, output, 2, false)

				// First row should be data, not header
				records := helpers.ValidateCSVStructure(t, output)
				require.Equal(t, "2", records[0][0], "First field should be movie ID, not header")
			},
		},
		{
			name:   "format empty movies collection",
			movies: fixtures.EmptyMovies,
			options: internal.FormatOptions{
				Format:   internal.FormatCSV,
				NoHeader: false,
			},
			expectedHeader: []string{
				"ID",
				"Title",
				"Original Title",
				"Year",
				"Rating",
				"Votes",
				"Popularity",
				"Genres",
				"Overview",
				"Language",
				"Adult",
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				helpers.ValidateCSVRecordCount(t, output, 0, true)
			},
		},
		{
			name:   "format movies with special characters",
			movies: fixtures.EdgeCaseMovies,
			options: internal.FormatOptions{
				Format:   internal.FormatCSV,
				NoHeader: false,
			},
			expectedHeader: []string{
				"ID",
				"Title",
				"Original Title",
				"Year",
				"Rating",
				"Votes",
				"Popularity",
				"Genres",
				"Overview",
				"Language",
				"Adult",
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				helpers.ValidateCSVRecordCount(t, output, len(fixtures.EdgeCaseMovies), true)

				// Verify Unicode handling
				require.Contains(t, output, "Unicode Title 🎬")
				require.Contains(t, output, "Título Unicódé ñáéíóú")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := internal.FormatMovies(&buf, tt.movies, tt.options)
			require.NoError(t, err)

			output := buf.String()
			require.NotEmpty(t, output)

			// Validate CSV structure
			if tt.expectedHeader != nil {
				helpers.ValidateCSVOutput(t, output, tt.expectedHeader, !tt.options.NoHeader)
			}

			// Run custom validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, output)
			}
		})
	}
}

func TestFormatTVShowsCSV(t *testing.T) {
	tests := []struct {
		name           string
		tvShows        []internal.TVShow
		options        internal.FormatOptions
		expectedHeader []string
		validateFunc   func(t *testing.T, output string)
	}{
		{
			name:    "format single TV show with header",
			tvShows: fixtures.SingleTVShow,
			options: internal.FormatOptions{
				Format:   internal.FormatCSV,
				NoHeader: false,
			},
			expectedHeader: []string{
				"ID",
				"Name",
				"Original Name",
				"Year",
				"Rating",
				"Votes",
				"Popularity",
				"Genres",
				"Overview",
				"Language",
				"Adult",
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				helpers.ValidateCSVRecordCount(t, output, 1, true)
				helpers.ValidateCSVColumnCount(t, output)
			},
		},
		{
			name:    "format multiple TV shows with header",
			tvShows: fixtures.SampleTVShows[:2],
			options: internal.FormatOptions{
				Format:   internal.FormatCSV,
				NoHeader: false,
			},
			expectedHeader: []string{
				"ID",
				"Name",
				"Original Name",
				"Year",
				"Rating",
				"Votes",
				"Popularity",
				"Genres",
				"Overview",
				"Language",
				"Adult",
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				helpers.ValidateCSVRecordCount(t, output, 2, true)
				helpers.ValidateCSVColumnCount(t, output)

				// Validate specific TV show data
				helpers.ValidateCSVFieldContent(t, output, 1, 1, "Breaking Bad")
				helpers.ValidateCSVFieldContent(t, output, 2, 1, "Squid Game")
			},
		},
		{
			name:    "format TV shows without header",
			tvShows: fixtures.SampleTVShows[:1],
			options: internal.FormatOptions{
				Format:   internal.FormatCSV,
				NoHeader: true,
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				helpers.ValidateCSVRecordCount(t, output, 1, false)

				// First row should be data, not header
				records := helpers.ValidateCSVStructure(t, output)
				require.Equal(t, "1", records[0][0], "First field should be TV show ID, not header")
			},
		},
		{
			name:    "format empty TV shows collection",
			tvShows: fixtures.EmptyTVShows,
			options: internal.FormatOptions{
				Format:   internal.FormatCSV,
				NoHeader: false,
			},
			expectedHeader: []string{
				"ID",
				"Name",
				"Original Name",
				"Year",
				"Rating",
				"Votes",
				"Popularity",
				"Genres",
				"Overview",
				"Language",
				"Adult",
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				helpers.ValidateCSVRecordCount(t, output, 0, true)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := internal.FormatTVShows(&buf, tt.tvShows, tt.options)
			require.NoError(t, err)

			output := buf.String()
			require.NotEmpty(t, output)

			// Validate CSV structure
			if tt.expectedHeader != nil {
				helpers.ValidateCSVOutput(t, output, tt.expectedHeader, !tt.options.NoHeader)
			}

			// Run custom validation
			if tt.validateFunc != nil {
				tt.validateFunc(t, output)
			}
		})
	}
}

func TestCSVSpecialCharacterHandling(t *testing.T) {
	tests := []struct {
		name    string
		movie   internal.Movie
		checkIn string
	}{
		{
			name: "quotes in CSV fields",
			movie: internal.Movie{
				ID:       1,
				Title:    `Movie with "quotes" in title`,
				Overview: `Overview with "double quotes" and 'single quotes'`,
			},
			checkIn: `"Movie with ""quotes"" in title"`,
		},
		{
			name: "commas in CSV fields",
			movie: internal.Movie{
				ID:     2,
				Title:  "Movie, with, commas",
				Genres: "Action, Comedy, Drama",
			},
			checkIn: `"Movie, with, commas"`,
		},
		{
			name: "newlines in CSV fields",
			movie: internal.Movie{
				ID:       3,
				Title:    "Movie with\nnewlines",
				Overview: "Overview with\nmultiple\nlines",
			},
			checkIn: "Movie with\nnewlines",
		},
		{
			name: "unicode characters",
			movie: internal.Movie{
				ID:    4,
				Title: "Movie with émojis 🎬 and ñáéíóú",
			},
			checkIn: "🎬",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := internal.FormatMovies(&buf, []internal.Movie{tt.movie}, internal.FormatOptions{
				Format:   internal.FormatCSV,
				NoHeader: true,
			})
			require.NoError(t, err)

			output := buf.String()
			require.NotEmpty(t, output)

			// Validate that the special character is properly handled
			require.Contains(t, output, tt.checkIn, "CSV should properly handle special characters")
		})
	}
}
