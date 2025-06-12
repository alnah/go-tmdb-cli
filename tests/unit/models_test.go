// tests/unit/models_test.go
package unit_test

import (
	"testing"
	"time"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/stretchr/testify/assert"
)

func TestFormatRating(t *testing.T) {
	tests := []struct {
		name     string
		rating   float64
		votes    int
		expected string
	}{
		{
			name:     "too few votes shows N/A",
			rating:   8.5,
			votes:    5,
			expected: "N/A",
		},
		{
			name:     "uncertain rating with question mark",
			rating:   7.2,
			votes:    50,
			expected: "7.2?",
		},
		{
			name:     "reliable rating without question mark",
			rating:   8.9,
			votes:    1000,
			expected: "8.9",
		},
		{
			name:     "exactly 10 votes threshold",
			rating:   6.5,
			votes:    10,
			expected: "6.5?",
		},
		{
			name:     "exactly 100 votes threshold",
			rating:   7.8,
			votes:    100,
			expected: "7.8",
		},
		{
			name:     "zero votes",
			rating:   9.0,
			votes:    0,
			expected: "N/A",
		},
		{
			name:     "high rating with many votes",
			rating:   9.5,
			votes:    50000,
			expected: "9.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.FormatRating(tt.rating, tt.votes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatVotes(t *testing.T) {
	tests := []struct {
		name     string
		votes    int
		expected string
	}{
		{
			name:     "small number unchanged",
			votes:    123,
			expected: "123",
		},
		{
			name:     "exactly 1000 shows K",
			votes:    1000,
			expected: "1.0K",
		},
		{
			name:     "thousands with decimal",
			votes:    1500,
			expected: "1.5K",
		},
		{
			name:     "large thousands",
			votes:    999999,
			expected: "1000.0K",
		},
		{
			name:     "exactly 1 million shows M",
			votes:    1000000,
			expected: "1.0M",
		},
		{
			name:     "millions with decimal",
			votes:    2500000,
			expected: "2.5M",
		},
		{
			name:     "zero votes",
			votes:    0,
			expected: "0",
		},
		{
			name:     "large millions",
			votes:    15750000,
			expected: "15.8M",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.FormatVotes(tt.votes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatGenres(t *testing.T) {
	tests := []struct {
		name     string
		genres   string
		maxWidth int
		expected string
	}{
		{
			name:     "empty genres returns N/A",
			genres:   "",
			maxWidth: 20,
			expected: "N/A",
		},
		{
			name:     "short genres unchanged",
			genres:   "Action, Comedy",
			maxWidth: 20,
			expected: "Action, Comedy",
		},
		{
			name:     "genres exactly at max width",
			genres:   "Action, Comedy, Drama",
			maxWidth: 21,
			expected: "Action, Comedy, Drama",
		},
		{
			name:     "long genres truncated",
			genres:   "Action, Comedy, Drama, Thriller, Adventure",
			maxWidth: 20,
			expected: "Action, Comedy, D...",
		},
		{
			name:     "very small max width",
			genres:   "Action",
			maxWidth: 5,
			expected: "Ac...",
		},
		{
			name:     "single character genres",
			genres:   "A",
			maxWidth: 10,
			expected: "A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.FormatGenres(tt.genres, tt.maxWidth)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseYear(t *testing.T) {
	tests := []struct {
		name        string
		releaseDate string
		expected    int
	}{
		{
			name:        "valid full date",
			releaseDate: "1999-10-15",
			expected:    1999,
		},
		{
			name:        "valid date different year",
			releaseDate: "2024-06-01",
			expected:    2024,
		},
		{
			name:        "empty date returns zero",
			releaseDate: "",
			expected:    0,
		},
		{
			name:        "partial date",
			releaseDate: "2023",
			expected:    2023,
		},
		{
			name:        "short date",
			releaseDate: "99",
			expected:    0,
		},
		{
			name:        "invalid date format",
			releaseDate: "invalid-date",
			expected:    0,
		},
		{
			name:        "date with only year and month",
			releaseDate: "2022-12",
			expected:    2022,
		},
		{
			name:        "future date",
			releaseDate: "2030-01-01",
			expected:    2030,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.ParseYear(tt.releaseDate)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestShortenOverview(t *testing.T) {
	tests := []struct {
		name      string
		overview  string
		maxLength int
		expected  string
	}{
		{
			name:      "short overview unchanged",
			overview:  "A short movie description.",
			maxLength: 50,
			expected:  "A short movie description.",
		},
		{
			name:      "overview exactly at max length",
			overview:  "Exactly twenty chars",
			maxLength: 20,
			expected:  "Exactly twenty chars",
		},
		{
			name:      "long overview truncated at word boundary",
			overview:  "This is a very long movie description that needs to be truncated properly.",
			maxLength: 30,
			expected:  "This is a very long movie...",
		},
		{
			name:      "very small max length",
			overview:  "This is a test description",
			maxLength: 8,
			expected:  "This is ...",
		},
		{
			name:      "max length less than 10",
			overview:  "Test description",
			maxLength: 5,
			expected:  "Test ...",
		},
		{
			name:      "no spaces in overview",
			overview:  "Wordwithoutspaces",
			maxLength: 10,
			expected:  "Wordwit...",
		},
		{
			name:      "empty overview",
			overview:  "",
			maxLength: 20,
			expected:  "",
		},
		{
			name:      "overview with last space beyond half",
			overview:  "This description has a space",
			maxLength: 15,
			expected:  "This descrip...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.ShortenOverview(tt.overview, tt.maxLength)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsHighlyRated(t *testing.T) {
	tests := []struct {
		name     string
		rating   float64
		votes    int
		expected bool
	}{
		{
			name:     "highly rated with enough votes",
			rating:   8.5,
			votes:    1000,
			expected: true,
		},
		{
			name:     "highly rated but not enough votes",
			rating:   9.0,
			votes:    50,
			expected: false,
		},
		{
			name:     "not highly rated with enough votes",
			rating:   7.5,
			votes:    1000,
			expected: false,
		},
		{
			name:     "exactly 8.0 rating threshold",
			rating:   8.0,
			votes:    1000,
			expected: true,
		},
		{
			name:     "exactly 100 votes threshold",
			rating:   8.5,
			votes:    100,
			expected: true,
		},
		{
			name:     "just below rating threshold",
			rating:   7.9,
			votes:    1000,
			expected: false,
		},
		{
			name:     "just below votes threshold",
			rating:   8.5,
			votes:    99,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.IsHighlyRated(tt.rating, tt.votes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsPopular(t *testing.T) {
	tests := []struct {
		name       string
		popularity float64
		votes      int
		expected   bool
	}{
		{
			name:       "high popularity",
			popularity: 75.0,
			votes:      500,
			expected:   true,
		},
		{
			name:       "low popularity but many votes",
			popularity: 25.0,
			votes:      1500,
			expected:   true,
		},
		{
			name:       "low popularity and few votes",
			popularity: 25.0,
			votes:      500,
			expected:   false,
		},
		{
			name:       "exactly 50.0 popularity threshold",
			popularity: 50.0,
			votes:      500,
			expected:   true,
		},
		{
			name:       "exactly 1000 votes threshold",
			popularity: 25.0,
			votes:      1000,
			expected:   true,
		},
		{
			name:       "just below popularity threshold",
			popularity: 49.9,
			votes:      500,
			expected:   false,
		},
		{
			name:       "just below votes threshold",
			popularity: 25.0,
			votes:      999,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.IsPopular(tt.popularity, tt.votes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsRecent(t *testing.T) {
	currentYear := time.Now().Year()

	tests := []struct {
		name     string
		year     int
		expected bool
	}{
		{
			name:     "current year",
			year:     currentYear,
			expected: true,
		},
		{
			name:     "last year",
			year:     currentYear - 1,
			expected: true,
		},
		{
			name:     "two years ago",
			year:     currentYear - 2,
			expected: true,
		},
		{
			name:     "three years ago",
			year:     currentYear - 3,
			expected: false,
		},
		{
			name:     "future year",
			year:     currentYear + 1,
			expected: true,
		},
		{
			name:     "far future year",
			year:     currentYear + 5,
			expected: true,
		},
		{
			name:     "old movie",
			year:     1999,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.IsRecent(tt.year)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildDisplayTitle(t *testing.T) {
	tests := []struct {
		name          string
		title         string
		originalTitle string
		expected      string
	}{
		{
			name:          "same titles",
			title:         "The Matrix",
			originalTitle: "The Matrix",
			expected:      "The Matrix",
		},
		{
			name:          "different titles",
			title:         "The Matrix",
			originalTitle: "Matrix",
			expected:      "The Matrix (Matrix)",
		},
		{
			name:          "empty original title",
			title:         "The Matrix",
			originalTitle: "",
			expected:      "The Matrix",
		},
		{
			name:          "foreign film with english title",
			title:         "Seven Samurai",
			originalTitle: "Shichinin no Samurai",
			expected:      "Seven Samurai (Shichinin no Samurai)",
		},
		{
			name:          "both titles empty",
			title:         "",
			originalTitle: "",
			expected:      "",
		},
		{
			name:          "original title empty string",
			title:         "Fight Club",
			originalTitle: "",
			expected:      "Fight Club",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internal.BuildDisplayTitle(tt.title, tt.originalTitle)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateSearchOptions(t *testing.T) {
	tests := []struct {
		name        string
		opts        internal.SearchOptions
		shouldError bool
		errorMsg    string
	}{
		{
			name: "valid options",
			opts: internal.SearchOptions{
				MaxItems:  20,
				MinRating: 7.0,
				MaxRating: 9.0,
			},
			shouldError: false,
		},
		{
			name: "zero page gets corrected",
			opts: internal.SearchOptions{
				Page:     0,
				MaxItems: 20,
			},
			shouldError: false,
		},
		{
			name: "negative max items",
			opts: internal.SearchOptions{
				MaxItems: -5,
			},
			shouldError: true,
			errorMsg:    "max-items must be between 1 and 1000",
		},
		{
			name: "too many max items",
			opts: internal.SearchOptions{
				MaxItems: 1500,
			},
			shouldError: true,
			errorMsg:    "max-items must be between 1 and 1000",
		},
		{
			name: "invalid min rating negative",
			opts: internal.SearchOptions{
				MaxItems:  20,
				MinRating: -1.0,
			},
			shouldError: true,
			errorMsg:    "min-rating must be between 0 and 10",
		},
		{
			name: "invalid min rating too high",
			opts: internal.SearchOptions{
				MaxItems:  20,
				MinRating: 11.0,
			},
			shouldError: true,
			errorMsg:    "min-rating must be between 0 and 10",
		},
		{
			name: "invalid max rating negative",
			opts: internal.SearchOptions{
				MaxItems:  20,
				MaxRating: -1.0,
			},
			shouldError: true,
			errorMsg:    "max-rating must be between 0 and 10",
		},
		{
			name: "invalid max rating too high",
			opts: internal.SearchOptions{
				MaxItems:  20,
				MaxRating: 15.0,
			},
			shouldError: true,
			errorMsg:    "max-rating must be between 0 and 10",
		},
		{
			name: "min rating greater than max rating",
			opts: internal.SearchOptions{
				MaxItems:  20,
				MinRating: 8.0,
				MaxRating: 6.0,
			},
			shouldError: true,
			errorMsg:    "min-rating (8.0) cannot be greater than max-rating (6.0)",
		},
		{
			name: "boundary values",
			opts: internal.SearchOptions{
				MaxItems:  1,
				MinRating: 0.0,
				MaxRating: 10.0,
			},
			shouldError: false,
		},
		{
			name: "max boundary values",
			opts: internal.SearchOptions{
				MaxItems:  1000,
				MinRating: 10.0,
				MaxRating: 10.0,
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalOpts := tt.opts // Keep copy for comparison
			err := internal.ValidateSearchOptions(&tt.opts)

			if tt.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
				// Check that zero page was corrected to 1
				if originalOpts.Page == 0 {
					assert.Equal(t, 1, tt.opts.Page)
				}
			}
		})
	}
}

func TestParseGenres(t *testing.T) {
	tests := []struct {
		name        string
		genreNames  []string
		expectedIDs []int
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "valid single genre",
			genreNames:  []string{"action"},
			expectedIDs: []int{28},
			shouldError: false,
		},
		{
			name:        "valid multiple genres",
			genreNames:  []string{"action", "comedy", "drama"},
			expectedIDs: []int{28, 35, 18},
			shouldError: false,
		},
		{
			name:        "sci-fi alias",
			genreNames:  []string{"sci-fi"},
			expectedIDs: []int{878},
			shouldError: false,
		},
		{
			name:        "science-fiction full name",
			genreNames:  []string{"science-fiction"},
			expectedIDs: []int{878},
			shouldError: false,
		},
		{
			name:        "case insensitive",
			genreNames:  []string{"ACTION", "Comedy", "dRaMa"},
			expectedIDs: []int{28, 35, 18},
			shouldError: false,
		},
		{
			name:        "genres with whitespace",
			genreNames:  []string{" action ", "  comedy  "},
			expectedIDs: []int{28, 35},
			shouldError: false,
		},
		{
			name:        "invalid genre",
			genreNames:  []string{"invalid-genre"},
			shouldError: true,
			errorMsg:    "unknown genre: invalid-genre",
		},
		{
			name:        "mixed valid and invalid",
			genreNames:  []string{"action", "invalid-genre"},
			shouldError: true,
			errorMsg:    "unknown genre: invalid-genre",
		},
		{
			name:        "empty slice",
			genreNames:  []string{},
			expectedIDs: nil,
			shouldError: false,
		},
		{
			name: "all major genres",
			genreNames: []string{
				"action",
				"adventure",
				"animation",
				"comedy",
				"crime",
				"documentary",
				"drama",
				"family",
				"fantasy",
				"history",
				"horror",
				"music",
				"mystery",
				"romance",
				"science-fiction",
				"thriller",
				"war",
				"western",
			},
			expectedIDs: []int{
				28,
				12,
				16,
				35,
				80,
				99,
				18,
				10751,
				14,
				36,
				27,
				10402,
				9648,
				10749,
				878,
				53,
				10752,
				37,
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := internal.ParseGenres(tt.genreNames)

			if tt.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedIDs, result)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := internal.DefaultConfig()

	assert.Equal(t, "https://api.themoviedb.org/3", config.BaseURL)
	assert.Equal(t, 30*time.Second, config.Timeout)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, 5*time.Minute, config.CacheTTL)
	assert.Equal(t, "info", config.LogLevel)
	assert.Equal(t, "table", config.Format)
	assert.Empty(t, config.APIKey) // Should be empty in default config
}

// Integration test for the genre mapping
func TestGenreMapCompleteness(t *testing.T) {
	genreMap := internal.GenreMap

	// Test that all expected genres are present
	expectedGenres := map[string]int{
		"action":          28,
		"adventure":       12,
		"animation":       16,
		"comedy":          35,
		"crime":           80,
		"documentary":     99,
		"drama":           18,
		"family":          10751,
		"fantasy":         14,
		"history":         36,
		"horror":          27,
		"music":           10402,
		"mystery":         9648,
		"romance":         10749,
		"science-fiction": 878,
		"sci-fi":          878, // Alias
		"thriller":        53,
		"war":             10752,
		"western":         37,
	}

	for name, expectedID := range expectedGenres {
		actualID, exists := genreMap[name]
		assert.True(t, exists, "Genre %s should exist in map", name)
		assert.Equal(t, expectedID, actualID, "Genre %s should have correct ID", name)
	}

	// Ensure we don't have unexpected extra genres
	assert.Len(t, genreMap, len(expectedGenres), "GenreMap should only contain expected genres")
}

// Edge case tests for error conditions
func TestEdgeCases(t *testing.T) {
	t.Run("FormatRating with extreme values", func(t *testing.T) {
		// Test with very large numbers
		result := internal.FormatRating(10.0, 1000000)
		assert.Equal(t, "10.0", result)

		// Test with zero rating
		result = internal.FormatRating(0.0, 1000)
		assert.Equal(t, "0.0", result)

		// Test with negative votes (shouldn't happen in real data)
		result = internal.FormatRating(8.0, -1)
		assert.Equal(t, "N/A", result)
	})

	t.Run("FormatVotes with extreme values", func(t *testing.T) {
		// Test with very large numbers
		result := internal.FormatVotes(999999999)
		assert.Equal(t, "1000.0M", result)

		// Test with negative votes
		result = internal.FormatVotes(-100)
		assert.Equal(t, "-100", result) // Should handle gracefully
	})

	t.Run("ParseYear with malformed dates", func(t *testing.T) {
		malformedDates := []string{
			"abcd-12-25",
			"20xx-01-01",
			"99999999999999",
			"2024-13-45", // Invalid month/day
			"not-a-date",
		}

		for _, date := range malformedDates {
			result := internal.ParseYear(date)
			// Should either return 0 or the year portion if parseable
			assert.True(
				t,
				result == 0 || result > 1800,
				"ParseYear should handle malformed date: %s",
				date,
			)
		}
	})
}

// Benchmark tests to ensure performance
func BenchmarkFormatRating(b *testing.B) {
	for i := 0; i < b.N; i++ {
		internal.FormatRating(8.5, 1000)
	}
}

func BenchmarkFormatVotes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		internal.FormatVotes(1500000)
	}
}

func BenchmarkParseGenres(b *testing.B) {
	genres := []string{"action", "comedy", "drama", "thriller"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = internal.ParseGenres(genres)
	}
}

// Property-based testing for format functions
func TestFormatRatingProperties(t *testing.T) {
	// Property: FormatRating should never return empty string for valid inputs
	for rating := 0.0; rating <= 10.0; rating += 0.1 {
		for votes := 0; votes <= 1000; votes += 100 {
			result := internal.FormatRating(rating, votes)
			assert.NotEmpty(t, result, "FormatRating should never return empty string")
		}
	}
}

func TestFormatVotesProperties(t *testing.T) {
	// Property: FormatVotes should always return a string representation
	testVotes := []int{0, 1, 999, 1000, 1500, 999999, 1000000, 2500000}

	for _, votes := range testVotes {
		result := internal.FormatVotes(votes)
		assert.NotEmpty(t, result, "FormatVotes should never return empty string")
		assert.NotContains(t, result, " ", "FormatVotes should not contain spaces")
	}
}
