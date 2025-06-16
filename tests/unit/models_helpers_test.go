package unit

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
)

func TestFormatRating(t *testing.T) {
	for _, tc := range fixtures.RatingTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := internal.FormatRating(tc.Rating, tc.Votes)
			assert.Equal(t, tc.Expected, result, tc.Description)
		})
	}
}

func TestFormatRatingEdgeCases(t *testing.T) {
	t.Run("extreme ratings", func(t *testing.T) {
		tests := []struct {
			rating   float64
			votes    int
			expected string
		}{
			{0.0, 1000, "0.0"},
			{10.0, 1000, "10.0"},
			{-1.0, 1000, "-1.0"}, // Edge case - negative rating
			{11.0, 1000, "11.0"}, // Edge case - rating over 10
		}

		for _, test := range tests {
			result := internal.FormatRating(test.rating, test.votes)
			assert.Equal(t, test.expected, result)
		}
	})

	t.Run("boundary vote counts", func(t *testing.T) {
		tests := []struct {
			rating   float64
			votes    int
			expected string
		}{
			{8.5, 9, "N/A"},   // Just below MinVotesForRating
			{8.5, 10, "8.5?"}, // Exactly MinVotesForRating
			{8.5, 99, "8.5?"}, // Just below MinVotesForUncertain
			{8.5, 100, "8.5"}, // Exactly MinVotesForUncertain
			{8.5, 101, "8.5"}, // Just above MinVotesForUncertain
		}

		for _, test := range tests {
			result := internal.FormatRating(test.rating, test.votes)
			assert.Equal(t, test.expected, result)
		}
	})
}

func TestFormatVotes(t *testing.T) {
	for _, tc := range fixtures.VoteTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := internal.FormatVotes(tc.Votes)
			assert.Equal(t, tc.Expected, result, tc.Description)
		})
	}
}

func TestFormatVotesEdgeCases(t *testing.T) {
	t.Run("boundary values", func(t *testing.T) {
		tests := []struct {
			votes    int
			expected string
		}{
			{999, "999"},        // Just below 1K
			{1000, "1.0K"},      // Exactly 1K
			{1001, "1.0K"},      // Just above 1K
			{999999, "1000.0K"}, // Just below 1M
			{1000000, "1.0M"},   // Exactly 1M
			{1000001, "1.0M"},   // Just above 1M
		}

		for _, test := range tests {
			result := internal.FormatVotes(test.votes)
			assert.Equal(t, test.expected, result)
		}
	})

	t.Run("large numbers", func(t *testing.T) {
		tests := []struct {
			votes    int
			expected string
		}{
			{5500000, "5.5M"},
			{12345678, "12.3M"},
			{999999999, "1000.0M"},
		}

		for _, test := range tests {
			result := internal.FormatVotes(test.votes)
			assert.Equal(t, test.expected, result)
		}
	})

	t.Run("negative votes", func(t *testing.T) {
		// Edge case - negative votes shouldn't happen but test behavior
		result := internal.FormatVotes(-100)
		assert.Equal(t, "-100", result)
	})
}

func TestFormatGenres(t *testing.T) {
	t.Run("short genres unchanged", func(t *testing.T) {
		tests := []struct {
			genres   string
			maxWidth int
			expected string
		}{
			{"Action", 25, "Action"},
			{"Action, Comedy", 25, "Action, Comedy"},
			{"", 25, "N/A"},
		}

		for _, test := range tests {
			result := internal.FormatGenres(test.genres, test.maxWidth)
			assert.Equal(t, test.expected, result)
		}
	})

	t.Run("long genres truncated", func(t *testing.T) {
		longGenres := "Action, Adventure, Comedy, Drama, Fantasy, Science Fiction"
		result := internal.FormatGenres(longGenres, 20)

		assert.Equal(t, "Action, Adventure...", result)
		assert.LessOrEqual(t, len(result), 20)
	})

	t.Run("empty genres returns N/A", func(t *testing.T) {
		result := internal.FormatGenres("", 25)
		assert.Equal(t, "N/A", result)
	})

	t.Run("exact width boundary", func(t *testing.T) {
		genres := "Action, Comedy" // 14 characters

		result1 := internal.FormatGenres(genres, 14)
		assert.Equal(t, genres, result1)

		result2 := internal.FormatGenres(genres, 13)
		assert.Equal(t, "Action, Co...", result2)
	})
}

func TestParseYear(t *testing.T) {
	for _, tc := range fixtures.YearTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := internal.ParseYear(tc.Input)
			assert.Equal(t, tc.Expected, result, tc.Description)
		})
	}
}

func TestParseYearEdgeCases(t *testing.T) {
	t.Run("various date formats", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int
		}{
			{"2023", 2023},
			{"2023-12", 2023},
			{"2023-12-25", 2023},
			{"2023-12-25T10:30:00Z", 2023},
			{"23", 0},       // Too short
			{"123", 0},      // Too short
			{"12345", 1234}, // Takes first 4 digits
			{"abcd", 0},     // Non-numeric
			{"202a", 0},     // Mixed alphanumeric
		}

		for _, test := range tests {
			result := internal.ParseYear(test.input)
			assert.Equal(t, test.expected, result, "Input: %s", test.input)
		}
	})
}

func TestParseTVYear(t *testing.T) {
	t.Run("TV year parsing delegates to ParseYear", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int
		}{
			{"2008-01-20", 2008},
			{"2021-09-17", 2021},
			{"", 0},
			{"invalid", 0},
		}

		for _, test := range tests {
			result := internal.ParseTVYear(test.input)
			expected := internal.ParseYear(test.input)
			assert.Equal(t, expected, result)
			assert.Equal(t, test.expected, result)
		}
	})
}

func TestShortenOverview(t *testing.T) {
	t.Run("short overview unchanged", func(t *testing.T) {
		overview := "Short overview"
		result := internal.ShortenOverview(overview, 50)
		assert.Equal(t, overview, result)
	})

	t.Run("long overview truncated", func(t *testing.T) {
		overview := "This is a very long overview that should be truncated because it exceeds the maximum length"
		result := internal.ShortenOverview(overview, 30)

		assert.True(t, len(result) <= 30)
		assert.True(t, strings.HasSuffix(result, "..."))
		assert.True(t, strings.HasPrefix(result, "This is a very"))
	})

	t.Run("truncation respects word boundaries", func(t *testing.T) {
		overview := "This is a test overview with multiple words"
		result := internal.ShortenOverview(overview, 20)

		assert.True(t, strings.HasSuffix(result, "..."))
		assert.False(t, strings.Contains(result, "w..."), "Should not break in middle of word")
	})

	t.Run("very small max length", func(t *testing.T) {
		overview := "This is a test"
		result := internal.ShortenOverview(overview, 5)

		// Based on the actual implementation, when maxLength < 10,
		// it still tries to truncate at maxLength-3 + "..."
		assert.True(t, strings.HasSuffix(result, "..."))
		assert.True(t, len(result) <= 5)
	})

	t.Run("max length less than 10", func(t *testing.T) {
		overview := "This is a longer overview"
		result := internal.ShortenOverview(overview, 8)

		// When maxLength < 10, it still uses maxLength-3 + "..."
		assert.True(t, strings.HasSuffix(result, "..."))
		assert.True(t, len(result) <= 8)
	})

	t.Run("exact length boundary", func(t *testing.T) {
		overview := "Exactly twenty chars" // 20 characters

		result1 := internal.ShortenOverview(overview, 20)
		assert.Equal(t, overview, result1)

		result2 := internal.ShortenOverview(overview, 19)
		assert.True(t, strings.HasSuffix(result2, "..."))
		assert.True(t, len(result2) <= 19)
	})
}

func TestIsHighlyRated(t *testing.T) {
	for _, tc := range fixtures.HelperTestCases.HighlyRated {
		result := internal.IsHighlyRated(tc.Rating, tc.Votes)
		assert.Equal(t, tc.Expected, result,
			"Rating: %.1f, Votes: %d", tc.Rating, tc.Votes)
	}
}

func TestIsHighlyRatedEdgeCases(t *testing.T) {
	t.Run("boundary conditions", func(t *testing.T) {
		tests := []struct {
			rating   float64
			votes    int
			expected bool
		}{
			{7.9, 100, false},    // Just below rating threshold
			{8.0, 100, true},     // Exactly rating threshold
			{8.0, 99, false},     // Just below votes threshold
			{8.0, 100, true},     // Exactly votes threshold
			{10.0, 100, true},    // Maximum rating
			{8.0, 1000000, true}, // High votes
		}

		for _, test := range tests {
			result := internal.IsHighlyRated(test.rating, test.votes)
			assert.Equal(t, test.expected, result)
		}
	})
}

func TestIsPopular(t *testing.T) {
	for _, tc := range fixtures.HelperTestCases.Popular {
		result := internal.IsPopular(tc.Popularity, tc.Votes)
		assert.Equal(t, tc.Expected, result,
			"Popularity: %.1f, Votes: %d", tc.Popularity, tc.Votes)
	}
}

func TestIsPopularEdgeCases(t *testing.T) {
	t.Run("boundary conditions", func(t *testing.T) {
		tests := []struct {
			popularity float64
			votes      int
			expected   bool
		}{
			{49.9, 999, false}, // Below both thresholds
			{50.0, 999, true},  // Popularity threshold met
			{49.9, 1000, true}, // Votes threshold met
			{50.0, 1000, true}, // Both thresholds met
			{0.0, 2000, true},  // High votes, low popularity
			{100.0, 0, true},   // High popularity, no votes
		}

		for _, test := range tests {
			result := internal.IsPopular(test.popularity, test.votes)
			assert.Equal(t, test.expected, result)
		}
	})
}

func TestIsRecent(t *testing.T) {
	currentYear := time.Now().Year()

	// Test with actual current year calculations
	tests := []struct {
		name     string
		year     int
		expected bool
	}{
		{"current year", currentYear, true},
		{"last year", currentYear - 1, true},
		{"two years ago", currentYear - 2, true},
		{"three years ago", currentYear - 3, false},
		{"ten years ago", currentYear - 10, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := internal.IsRecent(test.year)
			assert.Equal(t, test.expected, result,
				"Year: %d (current: %d)", test.year, currentYear)
		})
	}
}

func TestIsRecentWithCurrentYear(t *testing.T) {
	currentYear := time.Now().Year()

	t.Run("current and recent years", func(t *testing.T) {
		tests := []struct {
			year     int
			expected bool
		}{
			{currentYear, true},      // Current year
			{currentYear - 1, true},  // Last year
			{currentYear - 2, true},  // Two years ago
			{currentYear - 3, false}, // Three years ago
			{currentYear - 5, false}, // Five years ago
			{currentYear + 1, false}, // Future year should be false, not true
		}

		for _, test := range tests {
			result := internal.IsRecent(test.year)
			assert.Equal(t, test.expected, result,
				"Year: %d (current: %d)", test.year, currentYear)
		}
	})
}

func TestBuildDisplayTitle(t *testing.T) {
	t.Run("same title and original", func(t *testing.T) {
		result := internal.BuildDisplayTitle("The Matrix", "The Matrix")
		assert.Equal(t, "The Matrix", result)
	})

	t.Run("different title and original", func(t *testing.T) {
		result := internal.BuildDisplayTitle("Spirited Away", "千と千尋の神隠し")
		assert.Equal(t, "Spirited Away (千と千尋の神隠し)", result)
	})

	t.Run("empty original title", func(t *testing.T) {
		result := internal.BuildDisplayTitle("The Matrix", "")
		assert.Equal(t, "The Matrix", result)
	})

	t.Run("empty main title", func(t *testing.T) {
		result := internal.BuildDisplayTitle("", "Original")
		assert.Equal(t, " (Original)", result)
	})

	t.Run("both empty titles", func(t *testing.T) {
		result := internal.BuildDisplayTitle("", "")
		assert.Equal(t, "", result)
	})
}

func TestCleanQuery(t *testing.T) {
	t.Run("removes surrounding quotes", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{`"matrix"`, "matrix"},
			{`'matrix'`, "matrix"},
			{`"the matrix"`, "the matrix"},
			{`'the matrix'`, "the matrix"},
			{`"matrix`, `"matrix`}, // Only leading quote - should remain
			{`matrix"`, `matrix"`}, // Only trailing quote - should remain
			{`matrix`, "matrix"},   // No quotes
			{`""`, ""},             // Empty quotes
			{`''`, ""},             // Empty quotes
		}

		for _, test := range tests {
			result := internal.CleanQuery(test.input)
			assert.Equal(t, test.expected, result, "Input: %s", test.input)
		}
	})

	t.Run("preserves internal quotes", func(t *testing.T) {
		result := internal.CleanQuery(`"matrix "revolution""`)
		expected := `matrix "revolution"`
		assert.Equal(t, expected, result)
	})

	t.Run("handles mixed quotes", func(t *testing.T) {
		result := internal.CleanQuery(`"matrix'`)
		assert.Equal(t, "matrix", result)
	})
}
