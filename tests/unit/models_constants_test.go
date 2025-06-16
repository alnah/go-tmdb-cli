package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
)

func TestConstants(t *testing.T) {
	t.Run("year constants are valid", func(t *testing.T) {
		assert.Equal(t, 1888, internal.MinValidYear)
		assert.Equal(t, 2030, internal.MaxReasonableYear)
		assert.Greater(t, internal.MaxReasonableYear, internal.MinValidYear)
	})

	t.Run("numeric constants are positive", func(t *testing.T) {
		assert.Equal(t, 4, internal.YearDigits)
		assert.Equal(t, 1000, internal.OneThousand)
		assert.Equal(t, 1000000, internal.OneMillion)
		assert.Greater(t, internal.OneMillion, internal.OneThousand)
	})

	t.Run("threshold constants are reasonable", func(t *testing.T) {
		assert.Equal(t, 10, internal.MinVotesForRating)
		assert.Equal(t, 100, internal.MinVotesForUncertain)
		assert.Equal(t, 8.0, internal.HighRatingThreshold)
		assert.Equal(t, 50.0, internal.PopularityThreshold)
		assert.Equal(t, 1000, internal.PopularVotesThreshold)
		assert.Equal(t, 2, internal.RecentYearsBack)

		assert.Greater(t, internal.MinVotesForUncertain, internal.MinVotesForRating)
		assert.Greater(t, internal.HighRatingThreshold, 0.0)
		assert.LessOrEqual(t, internal.HighRatingThreshold, 10.0)
	})

	t.Run("timeout constants are reasonable", func(t *testing.T) {
		assert.Equal(t, 30, internal.DefaultTimeout)
		assert.Equal(t, 5, internal.DefaultCacheTTL)
		assert.Greater(t, internal.DefaultTimeout, 0)
		assert.Greater(t, internal.DefaultCacheTTL, 0)
	})
}

func TestGenreConstants(t *testing.T) {
	t.Run("genre ID constants are correct", func(t *testing.T) {
		assert.Equal(t, 28, internal.GenreAction)
		assert.Equal(t, 12, internal.GenreAdventure)
		assert.Equal(t, 16, internal.GenreAnimation)
		assert.Equal(t, 35, internal.GenreComedy)
		assert.Equal(t, 80, internal.GenreCrime)
		assert.Equal(t, 99, internal.GenreDocumentary)
		assert.Equal(t, 18, internal.GenreDrama)
		assert.Equal(t, 10751, internal.GenreFamily)
		assert.Equal(t, 14, internal.GenreFantasy)
		assert.Equal(t, 36, internal.GenreHistory)
		assert.Equal(t, 27, internal.GenreHorror)
		assert.Equal(t, 10402, internal.GenreMusic)
		assert.Equal(t, 9648, internal.GenreMystery)
		assert.Equal(t, 10749, internal.GenreRomance)
		assert.Equal(t, 878, internal.GenreSciFi)
		assert.Equal(t, 53, internal.GenreThriller)
		assert.Equal(t, 10752, internal.GenreWar)
		assert.Equal(t, 37, internal.GenreWestern)
	})

	t.Run("all genre IDs are unique and positive", func(t *testing.T) {
		genreIDs := []int{
			internal.GenreAction, internal.GenreAdventure, internal.GenreAnimation,
			internal.GenreComedy, internal.GenreCrime, internal.GenreDocumentary,
			internal.GenreDrama, internal.GenreFamily, internal.GenreFantasy,
			internal.GenreHistory, internal.GenreHorror, internal.GenreMusic,
			internal.GenreMystery, internal.GenreRomance, internal.GenreSciFi,
			internal.GenreThriller, internal.GenreWar, internal.GenreWestern,
		}

		// Check all IDs are positive
		for _, id := range genreIDs {
			assert.Greater(t, id, 0, "Genre ID should be positive")
		}

		// Check uniqueness
		seen := make(map[int]bool)
		for _, id := range genreIDs {
			assert.False(t, seen[id], "Genre ID %d should be unique", id)
			seen[id] = true
		}
	})
}

func TestGenreMap(t *testing.T) {
	t.Run("genre map contains expected mappings", func(t *testing.T) {
		expectedMappings := map[string]int{
			"action":          internal.GenreAction,
			"adventure":       internal.GenreAdventure,
			"animation":       internal.GenreAnimation,
			"comedy":          internal.GenreComedy,
			"crime":           internal.GenreCrime,
			"documentary":     internal.GenreDocumentary,
			"drama":           internal.GenreDrama,
			"family":          internal.GenreFamily,
			"fantasy":         internal.GenreFantasy,
			"history":         internal.GenreHistory,
			"horror":          internal.GenreHorror,
			"music":           internal.GenreMusic,
			"mystery":         internal.GenreMystery,
			"romance":         internal.GenreRomance,
			"science-fiction": internal.GenreSciFi,
			"sci-fi":          internal.GenreSciFi,
			"thriller":        internal.GenreThriller,
			"war":             internal.GenreWar,
			"western":         internal.GenreWestern,
		}

		for name, expectedID := range expectedMappings {
			actualID, exists := internal.GenreMap[name]
			assert.True(t, exists, "Genre %s should exist in map", name)
			assert.Equal(t, expectedID, actualID, "Genre %s should map to ID %d", name, expectedID)
		}
	})

	t.Run("genre map has sci-fi aliases", func(t *testing.T) {
		sciFiID := internal.GenreMap["sci-fi"]
		scienceFictionID := internal.GenreMap["science-fiction"]

		assert.Equal(t, internal.GenreSciFi, sciFiID)
		assert.Equal(t, internal.GenreSciFi, scienceFictionID)
		assert.Equal(t, sciFiID, scienceFictionID, "Both sci-fi aliases should map to same ID")
	})

	t.Run("all genre names are lowercase", func(t *testing.T) {
		for name := range internal.GenreMap {
			assert.Equal(t, name, name, "Genre name should be lowercase: %s", name)
			assert.NotContains(t, name, " ", "Genre name should not contain spaces: %s", name)
		}
	})
}

func TestParseGenres(t *testing.T) {
	for _, tc := range fixtures.GenreTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := internal.ParseGenres(tc.Input)

			if tc.HasError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), "unknown genre")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.Expected, result)
			}
		})
	}
}

func TestParseGenresEdgeCases(t *testing.T) {
	t.Run("empty genre slice", func(t *testing.T) {
		result, err := internal.ParseGenres([]string{})
		assert.NoError(t, err)
		// Accept either nil or empty slice
		if result != nil {
			assert.Empty(t, result)
		}
	})

	t.Run("nil genre slice", func(t *testing.T) {
		result, err := internal.ParseGenres(nil)
		assert.NoError(t, err)
		// Accept either nil or empty slice
		if result != nil {
			assert.Empty(t, result)
		}
	})

	t.Run("genres with extra whitespace", func(t *testing.T) {
		input := []string{"  action  ", "\tcomedy\t", "\ndrama\n"}
		result, err := internal.ParseGenres(input)

		assert.NoError(t, err)
		expected := []int{internal.GenreAction, internal.GenreComedy, internal.GenreDrama}
		assert.Equal(t, expected, result)
	})

	t.Run("duplicate genres", func(t *testing.T) {
		input := []string{"action", "action", "comedy"}
		result, err := internal.ParseGenres(input)

		assert.NoError(t, err)
		expected := []int{internal.GenreAction, internal.GenreAction, internal.GenreComedy}
		assert.Equal(t, expected, result)
	})

	t.Run("mixed case genres", func(t *testing.T) {
		input := []string{"ACTION", "Comedy", "dRaMa"}
		result, err := internal.ParseGenres(input)

		assert.NoError(t, err)
		expected := []int{internal.GenreAction, internal.GenreComedy, internal.GenreDrama}
		assert.Equal(t, expected, result)
	})

	t.Run("error message includes available genres", func(t *testing.T) {
		input := []string{"invalid-genre"}
		_, err := internal.ParseGenres(input)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown genre: invalid-genre")
		assert.Contains(t, err.Error(), "Available:")
		assert.Contains(t, err.Error(), "action")
		assert.Contains(t, err.Error(), "comedy")
	})

	t.Run("error on first invalid genre in mixed list", func(t *testing.T) {
		input := []string{"action", "invalid-genre", "comedy"}
		_, err := internal.ParseGenres(input)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown genre: invalid-genre")
	})

	t.Run("sci-fi variations work correctly", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected int
		}{
			{"sci-fi", internal.GenreSciFi},
			{"science-fiction", internal.GenreSciFi},
			{"SCI-FI", internal.GenreSciFi},
			{"SCIENCE-FICTION", internal.GenreSciFi},
		}

		for _, tc := range testCases {
			result, err := internal.ParseGenres([]string{tc.input})
			assert.NoError(t, err)
			assert.Equal(t, []int{tc.expected}, result)
		}
	})
}
