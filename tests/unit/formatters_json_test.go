package unit

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestFormatMoviesJSON(t *testing.T) {
	tests := []struct {
		name         string
		movies       []internal.Movie
		validateFunc func(t *testing.T, output string, movies []internal.Movie)
	}{
		{
			name:   "format single movie",
			movies: fixtures.SingleMovie,
			validateFunc: func(t *testing.T, output string, movies []internal.Movie) {
				t.Helper()
				var result struct {
					Movies []internal.Movie `json:"movies"`
					Count  int              `json:"count"`
				}

				err := json.Unmarshal([]byte(output), &result)
				require.NoError(t, err, "Output should be valid JSON")

				assert.Equal(t, 1, result.Count)
				assert.Len(t, result.Movies, 1)
				assert.Equal(t, movies[0].ID, result.Movies[0].ID)
				assert.Equal(t, movies[0].Title, result.Movies[0].Title)
				assert.Equal(t, movies[0].Year, result.Movies[0].Year)
			},
		},
		{
			name:   "format multiple movies",
			movies: fixtures.SampleMovies[:3],
			validateFunc: func(t *testing.T, output string, movies []internal.Movie) {
				t.Helper()
				var result struct {
					Movies []internal.Movie `json:"movies"`
					Count  int              `json:"count"`
				}

				err := json.Unmarshal([]byte(output), &result)
				require.NoError(t, err, "Output should be valid JSON")

				assert.Equal(t, 3, result.Count)
				assert.Len(t, result.Movies, 3)

				for i, movie := range movies {
					assert.Equal(t, movie.ID, result.Movies[i].ID)
					assert.Equal(t, movie.Title, result.Movies[i].Title)
				}
			},
		},
		{
			name:   "format empty movies collection",
			movies: fixtures.EmptyMovies,
			validateFunc: func(t *testing.T, output string, movies []internal.Movie) {
				t.Helper()
				var result struct {
					Movies []internal.Movie `json:"movies"`
					Count  int              `json:"count"`
				}

				err := json.Unmarshal([]byte(output), &result)
				require.NoError(t, err, "Output should be valid JSON")

				assert.Equal(t, 0, result.Count)
				assert.Len(t, result.Movies, 0)
				assert.NotNil(t, result.Movies) // Should be empty array, not nil
			},
		},
		{
			name:   "format movies with edge case data",
			movies: fixtures.EdgeCaseMovies,
			validateFunc: func(t *testing.T, output string, movies []internal.Movie) {
				t.Helper()
				var result struct {
					Movies []internal.Movie `json:"movies"`
					Count  int              `json:"count"`
				}

				err := json.Unmarshal([]byte(output), &result)
				require.NoError(t, err, "Output should be valid JSON")

				assert.Equal(t, len(movies), result.Count)
				assert.Len(t, result.Movies, len(movies))

				// Verify Unicode handling
				found := false
				for _, movie := range result.Movies {
					if movie.Title == "Unicode Title 🎬" {
						found = true
						assert.Equal(t, "Título Unicódé ñáéíóú", movie.OriginalTitle)
						break
					}
				}
				assert.True(t, found, "Unicode movie should be present")
			},
		},
		{
			name: "format movies with special characters",
			movies: []internal.Movie{
				{
					ID:            100,
					Title:         "Movie with \"quotes\" and 'apostrophes'",
					OriginalTitle: "Original with\nnewlines\tand\ttabs",
					Year:          2023,
					Rating:        7.5,
					Votes:         1000,
					Popularity:    50.0,
					Genres:        "Action, Comedy",
					Overview:      "Overview with <HTML> & special chars",
					Language:      "en",
					Adult:         false,
				},
			},
			validateFunc: func(t *testing.T, output string, movies []internal.Movie) {
				t.Helper()
				var result struct {
					Movies []internal.Movie `json:"movies"`
					Count  int              `json:"count"`
				}

				err := json.Unmarshal([]byte(output), &result)
				require.NoError(t, err, "Output should be valid JSON")

				assert.Equal(t, 1, result.Count)
				movie := result.Movies[0]
				assert.Contains(t, movie.Title, "quotes")
				assert.Contains(t, movie.OriginalTitle, "newlines")
				assert.Contains(t, movie.Overview, "<HTML>")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := internal.FormatMovies(
				&buf,
				tt.movies,
				internal.FormatOptions{Format: internal.FormatJSON},
			)
			require.NoError(t, err)

			output := buf.String()
			require.NotEmpty(t, output)

			// Validate JSON structure
			helpers.ValidateJSONOutput(t, output, true)

			// Run custom validation
			tt.validateFunc(t, output, tt.movies)

			// Validate JSON formatting (indentation)
			assert.Contains(t, output, "\n")
			assert.Contains(t, output, "  ")
		})
	}
}

func TestFormatTVShowsJSON(t *testing.T) {
	tests := []struct {
		name         string
		tvShows      []internal.TVShow
		validateFunc func(t *testing.T, output string, tvShows []internal.TVShow)
	}{
		{
			name:    "format single TV show",
			tvShows: fixtures.SingleTVShow,
			validateFunc: func(t *testing.T, output string, tvShows []internal.TVShow) {
				t.Helper()
				var result struct {
					TVShows []internal.TVShow `json:"tv_shows"`
					Count   int               `json:"count"`
				}

				err := json.Unmarshal([]byte(output), &result)
				require.NoError(t, err, "Output should be valid JSON")

				assert.Equal(t, 1, result.Count)
				assert.Len(t, result.TVShows, 1)
				assert.Equal(t, tvShows[0].ID, result.TVShows[0].ID)
				assert.Equal(t, tvShows[0].Name, result.TVShows[0].Name)
				assert.Equal(t, tvShows[0].Year, result.TVShows[0].Year)
			},
		},
		{
			name:    "format multiple TV shows",
			tvShows: fixtures.SampleTVShows[:2],
			validateFunc: func(t *testing.T, output string, tvShows []internal.TVShow) {
				t.Helper()
				var result struct {
					TVShows []internal.TVShow `json:"tv_shows"`
					Count   int               `json:"count"`
				}

				err := json.Unmarshal([]byte(output), &result)
				require.NoError(t, err, "Output should be valid JSON")

				assert.Equal(t, 2, result.Count)
				assert.Len(t, result.TVShows, 2)

				for i, show := range tvShows {
					assert.Equal(t, show.ID, result.TVShows[i].ID)
					assert.Equal(t, show.Name, result.TVShows[i].Name)
				}
			},
		},
		{
			name:    "format empty TV shows collection",
			tvShows: fixtures.EmptyTVShows,
			validateFunc: func(t *testing.T, output string, tvShows []internal.TVShow) {
				t.Helper()
				var result struct {
					TVShows []internal.TVShow `json:"tv_shows"`
					Count   int               `json:"count"`
				}

				err := json.Unmarshal([]byte(output), &result)
				require.NoError(t, err, "Output should be valid JSON")

				assert.Equal(t, 0, result.Count)
				assert.Len(t, result.TVShows, 0)
				assert.NotNil(t, result.TVShows) // Should be empty array, not nil
			},
		},
		{
			name:    "format TV shows with edge case data",
			tvShows: fixtures.EdgeCaseTVShows,
			validateFunc: func(t *testing.T, output string, tvShows []internal.TVShow) {
				t.Helper()
				var result struct {
					TVShows []internal.TVShow `json:"tv_shows"`
					Count   int               `json:"count"`
				}

				err := json.Unmarshal([]byte(output), &result)
				require.NoError(t, err, "Output should be valid JSON")

				assert.Equal(t, len(tvShows), result.Count)
				assert.Len(t, result.TVShows, len(tvShows))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := internal.FormatTVShows(
				&buf,
				tt.tvShows,
				internal.FormatOptions{Format: internal.FormatJSON},
			)
			require.NoError(t, err)

			output := buf.String()
			require.NotEmpty(t, output)

			// Validate JSON structure
			helpers.ValidateJSONOutput(t, output, false)

			// Run custom validation
			tt.validateFunc(t, output, tt.tvShows)

			// Validate JSON formatting (indentation)
			assert.Contains(t, output, "\n")
			assert.Contains(t, output, "  ")
		})
	}
}

func TestFormatSearchResultJSON(t *testing.T) {
	tests := []struct {
		name         string
		searchResult internal.SearchResult
		validateFunc func(t *testing.T, output string)
	}{
		{
			name: "format complete search result",
			searchResult: internal.SearchResult{
				Movies:       fixtures.SampleMovies[:2],
				Page:         1,
				TotalPages:   5,
				TotalResults: 100,
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				var result internal.SearchResult
				err := json.Unmarshal([]byte(output), &result)
				require.NoError(t, err)

				assert.Equal(t, 1, result.Page)
				assert.Equal(t, 5, result.TotalPages)
				assert.Equal(t, 100, result.TotalResults)
				assert.Len(t, result.Movies, 2)
			},
		},
		{
			name: "format empty search result",
			searchResult: internal.SearchResult{
				Movies:       []internal.Movie{},
				Page:         1,
				TotalPages:   0,
				TotalResults: 0,
			},
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				var result internal.SearchResult
				err := json.Unmarshal([]byte(output), &result)
				require.NoError(t, err)

				assert.Equal(t, 1, result.Page)
				assert.Equal(t, 0, result.TotalPages)
				assert.Equal(t, 0, result.TotalResults)
				assert.Len(t, result.Movies, 0)
				assert.NotNil(t, result.Movies)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := internal.FormatSearchResult(
				&buf,
				tt.searchResult,
				internal.FormatOptions{Format: internal.FormatJSON},
			)
			require.NoError(t, err)

			output := buf.String()
			require.NotEmpty(t, output)

			tt.validateFunc(t, output)

			// Validate JSON formatting
			assert.Contains(t, output, "\n")
			assert.Contains(t, output, "  ")
		})
	}
}

func TestJSONOutputFormatting(t *testing.T) {
	tests := []struct {
		name         string
		validateFunc func(t *testing.T, output string)
	}{
		{
			name: "proper indentation",
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				assert.Contains(t, output, "{\n")
				assert.Contains(t, output, "  \"movies\":")
				assert.Contains(t, output, "  \"count\":")
			},
		},
		{
			name: "valid JSON structure",
			validateFunc: func(t *testing.T, output string) {
				t.Helper()
				var result map[string]any
				err := json.Unmarshal([]byte(output), &result)
				require.NoError(t, err, "Output should be valid JSON")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := internal.FormatMovies(
				&buf,
				fixtures.SampleMovies[:1],
				internal.FormatOptions{Format: internal.FormatJSON},
			)
			require.NoError(t, err)

			output := buf.String()
			tt.validateFunc(t, output)
		})
	}
}
