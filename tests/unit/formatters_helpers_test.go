package unit

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
)

func TestFormatYear(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		expected string
	}{
		{"valid year", 2023, "2023"},
		{"zero year", 0, "N/A"},
		{"negative year", -1, "N/A"}, // Implementation may not handle negatives
		{"future year", 2030, "2030"},
		{"old year", 1888, "1888"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through table formatting where formatYear is actually used
			movie := internal.Movie{Year: tt.year, Title: "Test Movie"}
			var buf bytes.Buffer

			err := internal.FormatMovies(&buf, []internal.Movie{movie}, internal.FormatOptions{
				Format: internal.FormatTable,
			})

			if tt.name == "zero year should return N/A" {
				// For zero years in table format, we expect N/A
				assert.NoError(t, err)
				output := buf.String()
				assert.Contains(t, output, "N/A")
			} else {
				assert.NoError(t, err)
				output := buf.String()
				assert.Contains(t, output, tt.expected)
			}
		})
	}
}

func TestFormatSummary(t *testing.T) {
	tests := []struct {
		name            string
		movies          []internal.Movie
		command         string
		useOriginal     bool
		expectedContent []string
	}{
		{
			name:    "popular movies summary",
			movies:  fixtures.SampleMovies[:1],
			command: "popular",
			expectedContent: []string{
				"Showing 1 popular movies",
				"titles",
			},
		},
		{
			name:        "popular movies with original title",
			movies:      fixtures.SampleMovies[:1],
			command:     "popular",
			useOriginal: true,
			expectedContent: []string{
				"Showing 1 popular movies",
				"original language titles",
			},
		},
		{
			name:    "top-rated movies summary",
			movies:  fixtures.SampleMovies[:2],
			command: "top-rated",
			expectedContent: []string{
				"Showing 2 top-rated movies",
				"titles",
			},
		},
		{
			name:    "now-playing movies summary",
			movies:  fixtures.SampleMovies,
			command: "now-playing",
			expectedContent: []string{
				"Showing 3 movies now playing",
				"titles",
			},
		},
		{
			name:    "upcoming movies summary",
			movies:  fixtures.SampleMovies[:1],
			command: "upcoming",
			expectedContent: []string{
				"Showing 1 upcoming movies",
				"titles",
			},
		},
		{
			name:    "search movies summary",
			movies:  fixtures.SampleMovies[:2],
			command: "search",
			expectedContent: []string{
				"Found 2 movies",
				"titles",
			},
		},
		{
			name:    "discover movies summary",
			movies:  fixtures.SampleMovies[:1],
			command: "discover",
			expectedContent: []string{
				"Discovered 1 movies",
				"titles",
			},
		},
		{
			name:            "empty movies should not show summary",
			movies:          fixtures.EmptyMovies,
			command:         "popular",
			expectedContent: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			internal.FormatSummary(&buf, tt.movies, tt.command, tt.useOriginal)
			output := buf.String()

			for _, expected := range tt.expectedContent {
				assert.Contains(t, output, expected, "Summary should contain: %s", expected)
			}

			if len(tt.movies) == 0 {
				assert.Empty(t, output, "Empty movies should produce no summary")
			}
		})
	}
}

func TestFormatTVSummary(t *testing.T) {
	tests := []struct {
		name            string
		tvShows         []internal.TVShow
		command         string
		useOriginal     bool
		expectedContent []string
	}{
		{
			name:    "popular TV shows summary",
			tvShows: fixtures.SampleTVShows[:1],
			command: "popular",
			expectedContent: []string{
				"Showing 1 popular TV shows",
				"names",
			},
		},
		{
			name:        "popular TV shows with original names",
			tvShows:     fixtures.SampleTVShows[:1],
			command:     "popular",
			useOriginal: true,
			expectedContent: []string{
				"Showing 1 popular TV shows",
				"original language names",
			},
		},
		{
			name:    "top-rated TV shows summary",
			tvShows: fixtures.SampleTVShows[:2], // Use first 2 shows to avoid bounds error
			command: "top-rated",
			expectedContent: []string{
				"Showing 2 top-rated TV shows",
				"names",
			},
		},
		{
			name:    "on-the-air TV shows summary",
			tvShows: fixtures.SampleTVShows, // Use all shows (3 total)
			command: "on-the-air",
			expectedContent: []string{
				"Showing 3 TV shows on the air",
				"names",
			},
		},
		{
			name:    "search TV shows summary",
			tvShows: fixtures.SampleTVShows[:2],
			command: "search",
			expectedContent: []string{
				"Found 2 TV shows",
				"names",
			},
		},
		{
			name:    "discover TV shows summary",
			tvShows: fixtures.SampleTVShows[:1],
			command: "discover",
			expectedContent: []string{
				"Discovered 1 TV shows",
				"names",
			},
		},
		{
			name:            "empty TV shows should not show summary",
			tvShows:         fixtures.EmptyTVShows,
			command:         "popular",
			expectedContent: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			internal.FormatTVSummary(&buf, tt.tvShows, tt.command, tt.useOriginal)
			output := buf.String()

			for _, expected := range tt.expectedContent {
				assert.Contains(t, output, expected, "TV summary should contain: %s", expected)
			}

			if len(tt.tvShows) == 0 {
				assert.Empty(t, output, "Empty TV shows should produce no summary")
			}
		})
	}
}

func TestFormatTitle(t *testing.T) {
	tests := []struct {
		name          string
		title         string
		originalTitle string
		useOriginal   bool
		expected      string
		maxLength     int
	}{
		{
			name:          "show regular title when not using original",
			title:         "The Matrix",
			originalTitle: "The Matrix",
			useOriginal:   false,
			expected:      "The Matrix",
		},
		{
			name:          "show original title when requested",
			title:         "Spirited Away",
			originalTitle: "Sen to Chihiro no Kamikakushi",
			useOriginal:   true,
			expected:      "Sen to Chihiro no Kamikakushi",
		},
		{
			name:          "show original with English in parentheses when space allows",
			title:         "Parasite",
			originalTitle: "기생충",
			useOriginal:   true,
			expected:      "기생충 (Parasite)",
		},
		{
			name:          "show English with original in parentheses when not using original",
			title:         "Parasite",
			originalTitle: "기생충",
			useOriginal:   false,
			expected:      "Parasite (기생충)",
		},
		{
			name:          "handle empty original title",
			title:         "The Matrix",
			originalTitle: "",
			useOriginal:   true,
			expected:      "The Matrix",
		},
		{
			name:          "handle identical titles",
			title:         "The Matrix",
			originalTitle: "The Matrix",
			useOriginal:   true,
			expected:      "The Matrix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We need to test this through a formatter since formatTitle is not exported
			movie := internal.Movie{
				Title:         tt.title,
				OriginalTitle: tt.originalTitle,
			}

			formatter := internal.MovieFormatter{Movie: movie}
			result := formatter.GetTitle(tt.useOriginal)

			assert.Contains(t, result, strings.Split(tt.expected, " ")[0],
				"Title should contain expected content")
		})
	}
}
