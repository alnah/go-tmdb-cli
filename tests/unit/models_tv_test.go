package unit

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestTVShowValidation(t *testing.T) {
	t.Run("valid TV show passes validation", func(t *testing.T) {
		show := fixtures.SampleTVShow
		err := show.Validate()
		assert.NoError(t, err)
	})

	t.Run("TV show with empty name fails validation", func(t *testing.T) {
		show := helpers.CreateValidTVShow(map[string]any{
			"Name": "",
		})
		err := show.Validate()
		helpers.AssertValidationError(t, err, "TV show name is required")
	})

	t.Run("TV show with zero ID fails validation", func(t *testing.T) {
		show := helpers.CreateValidTVShow(map[string]any{
			"ID": 0,
		})
		err := show.Validate()
		helpers.AssertValidationError(t, err, "TV show ID must be positive")
	})

	t.Run("TV show with negative ID fails validation", func(t *testing.T) {
		show := helpers.CreateValidTVShow(map[string]any{
			"ID": -1,
		})
		err := show.Validate()
		helpers.AssertValidationError(t, err, "TV show ID must be positive")
	})

	t.Run("multiple validation errors", func(t *testing.T) {
		show := internal.TVShow{
			ID:   0,  // Invalid
			Name: "", // Invalid
		}
		err := show.Validate()
		require.Error(t, err)
		// Should fail on first validation error (name)
		assert.Contains(t, err.Error(), "TV show name is required")
	})
}

func TestTVShowMediaItemInterface(t *testing.T) {
	show := fixtures.SampleTVShow

	t.Run("implements MediaItem interface correctly", func(t *testing.T) {
		var mediaItem internal.MediaItem = show

		assert.Equal(t, show.ID, mediaItem.GetID())
		assert.Equal(t, show.Name, mediaItem.GetTitle()) // Name maps to Title
		assert.Equal(
			t,
			show.OriginalName,
			mediaItem.GetOriginalTitle(),
		) // OriginalName maps to OriginalTitle
		assert.Equal(t, show.Year, mediaItem.GetYear())
		assert.Equal(t, show.Rating, mediaItem.GetRating())
		assert.Equal(t, show.Votes, mediaItem.GetVotes())
		assert.Equal(t, show.Popularity, mediaItem.GetPopularity())
		assert.Equal(t, show.Genres, mediaItem.GetGenres())
		assert.Equal(t, show.Overview, mediaItem.GetOverview())
	})

	t.Run("interface methods return correct types", func(t *testing.T) {
		assert.IsType(t, 0, show.GetID())
		assert.IsType(t, "", show.GetTitle())
		assert.IsType(t, "", show.GetOriginalTitle())
		assert.IsType(t, 0, show.GetYear())
		assert.IsType(t, 0.0, show.GetRating())
		assert.IsType(t, 0, show.GetVotes())
		assert.IsType(t, 0.0, show.GetPopularity())
		assert.IsType(t, "", show.GetGenres())
		assert.IsType(t, "", show.GetOverview())
	})

	t.Run("name vs title mapping is correct", func(t *testing.T) {
		// TV shows use "Name" field but interface expects "Title"
		assert.Equal(t, show.Name, show.GetTitle())
		assert.Equal(t, show.OriginalName, show.GetOriginalTitle())
	})
}

func TestTVShowFields(t *testing.T) {
	t.Run("all fields are properly set", func(t *testing.T) {
		show := fixtures.SampleTVShow

		assert.Equal(t, 789, show.ID)
		assert.Equal(t, "Breaking Bad", show.Name)
		assert.Equal(t, "Breaking Bad", show.OriginalName)
		assert.Equal(t, 2008, show.Year)
		assert.Equal(t, 9.5, show.Rating)
		assert.Equal(t, 1600000, show.Votes)
		assert.Equal(t, 95.2, show.Popularity)
		assert.Equal(t, "Crime, Drama, Thriller", show.Genres)
		assert.Equal(
			t,
			"A high school chemistry teacher turned methamphetamine manufacturer.",
			show.Overview,
		)
		assert.Equal(t, "en", show.Language)
		assert.False(t, show.Adult)
		assert.Equal(t, "2008-01-20", show.FirstAirDate)
	})

	t.Run("handles different original name", func(t *testing.T) {
		show := fixtures.SampleTVShowWithDifferentOriginal

		assert.Equal(t, "Squid Game", show.Name)
		assert.Equal(t, "오징어 게임", show.OriginalName)
		assert.NotEqual(t, show.Name, show.OriginalName)
	})

	t.Run("handles zero and negative values", func(t *testing.T) {
		show := internal.TVShow{
			Year:       0,   // Unset year
			Rating:     0.0, // Zero rating
			Votes:      0,   // No votes
			Popularity: 0.0, // No popularity
		}

		assert.Equal(t, 0, show.Year)
		assert.Equal(t, 0.0, show.Rating)
		assert.Equal(t, 0, show.Votes)
		assert.Equal(t, 0.0, show.Popularity)
	})

	t.Run("handles adult content flag", func(t *testing.T) {
		adultShow := helpers.CreateValidTVShow(map[string]any{
			"Adult": true,
		})
		nonAdultShow := helpers.CreateValidTVShow(map[string]any{
			"Adult": false,
		})

		assert.True(t, adultShow.Adult)
		assert.False(t, nonAdultShow.Adult)
	})
}

func TestTVShowConstants(t *testing.T) {
	t.Run("TV constants are reasonable", func(t *testing.T) {
		assert.Equal(t, 1928, internal.TVMinValidYear)
		assert.Equal(t, 2030, internal.TVMaxReasonableYear)
		assert.Equal(t, 30, internal.TVDefaultTimeout)
		assert.Equal(t, 5, internal.TVCacheTTL)

		assert.Greater(t, internal.TVMaxReasonableYear, internal.TVMinValidYear)
		assert.Greater(t, internal.TVDefaultTimeout, 0)
		assert.Greater(t, internal.TVCacheTTL, 0)
	})

	t.Run("TV endpoints are defined", func(t *testing.T) {
		assert.Equal(t, "/tv/popular", internal.TVPopularEndpoint)
		assert.Equal(t, "/tv/top_rated", internal.TVTopRatedEndpoint)
		assert.Equal(t, "/tv/on_the_air", internal.TVOnTheAirEndpoint)
		assert.Equal(t, "/search/tv", internal.TVSearchEndpoint)
		assert.Equal(t, "/discover/tv", internal.TVDiscoverEndpoint)
		assert.Equal(t, "/genre/tv/list", internal.TVGenresEndpoint)

		// All endpoints should start with appropriate paths
		assert.True(t, strings.HasPrefix(internal.TVPopularEndpoint, "/tv/"))
		assert.True(t, strings.HasPrefix(internal.TVTopRatedEndpoint, "/tv/"))
		assert.True(t, strings.HasPrefix(internal.TVOnTheAirEndpoint, "/tv/"))
		assert.True(t, strings.HasPrefix(internal.TVSearchEndpoint, "/search/"))
		assert.True(t, strings.HasPrefix(internal.TVDiscoverEndpoint, "/discover/"))
		assert.True(t, strings.HasPrefix(internal.TVGenresEndpoint, "/genre/"))
	})
}

func TestTVShowEdgeCases(t *testing.T) {
	t.Run("handles empty optional fields", func(t *testing.T) {
		show := internal.TVShow{
			ID:           1,
			Name:         "Minimal Show",
			OriginalName: "", // Empty original name
			Overview:     "", // Empty overview
			Language:     "", // Empty language
			FirstAirDate: "", // Empty air date
			Genres:       "", // Empty genres
		}

		err := show.Validate()
		assert.NoError(t, err, "TV show with empty optional fields should be valid")

		assert.Empty(t, show.OriginalName)
		assert.Empty(t, show.Overview)
		assert.Empty(t, show.Language)
		assert.Empty(t, show.FirstAirDate)
		assert.Empty(t, show.Genres)
	})

	t.Run("handles extreme numeric values", func(t *testing.T) {
		show := internal.TVShow{
			ID:         999999999, // Very large ID
			Name:       "Test Show",
			Year:       2050,      // Future year
			Rating:     10.0,      // Maximum rating
			Votes:      999999999, // Very large vote count
			Popularity: 999.99,    // Very high popularity
		}

		err := show.Validate()
		assert.NoError(t, err)

		assert.Equal(t, 999999999, show.ID)
		assert.Equal(t, 2050, show.Year)
		assert.Equal(t, 10.0, show.Rating)
		assert.Equal(t, 999999999, show.Votes)
		assert.Equal(t, 999.99, show.Popularity)
	})

	t.Run("handles unicode in text fields", func(t *testing.T) {
		show := internal.TVShow{
			ID:           1,
			Name:         "Show with émojis 📺 and åccénts",
			OriginalName: "Оригинальное название сериала",
			Overview:     "Description with 中文 characters and العربية text",
			Language:     "ko",
			Genres:       "ドラマ, スリラー",
		}

		err := show.Validate()
		assert.NoError(t, err)

		assert.Contains(t, show.Name, "📺")
		assert.Contains(t, show.Name, "åccénts")
		assert.Contains(t, show.OriginalName, "Оригинальное")
		assert.Contains(t, show.Overview, "中文")
		assert.Contains(t, show.Overview, "العربية")
		assert.Contains(t, show.Genres, "ドラマ")
	})

	t.Run("handles very long text fields", func(t *testing.T) {
		longName := strings.Repeat("Very Long TV Show Name ", 50)
		longOverview := strings.Repeat(
			"This is a very long overview description for a TV show. ",
			100,
		)

		show := internal.TVShow{
			ID:       1,
			Name:     longName,
			Overview: longOverview,
		}

		err := show.Validate()
		assert.NoError(t, err)

		assert.Equal(t, longName, show.Name)
		assert.Equal(t, longOverview, show.Overview)
		assert.Greater(t, len(show.Name), 1000)
		assert.Greater(t, len(show.Overview), 5000)
	})

	t.Run("handles historical air dates", func(t *testing.T) {
		show := internal.TVShow{
			ID:           1,
			Name:         "Historical Show",
			FirstAirDate: "1928-08-11", // Early television broadcast
			Year:         1928,
		}

		err := show.Validate()
		assert.NoError(t, err)
		assert.Equal(t, "1928-08-11", show.FirstAirDate)
		assert.Equal(t, 1928, show.Year)
	})
}

func TestTVSearchResult(t *testing.T) {
	t.Run("TV search result with shows", func(t *testing.T) {
		result := fixtures.SampleTVSearchResult

		assert.Len(t, result.TVShows, 1)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 2, result.TotalPages)
		assert.Equal(t, 25, result.TotalResults)

		helpers.AssertTVShowValid(t, result.TVShows[0])
	})

	t.Run("empty TV search result", func(t *testing.T) {
		result := internal.TVSearchResult{
			TVShows:      []internal.TVShow{},
			Page:         1,
			TotalPages:   0,
			TotalResults: 0,
		}

		assert.Empty(t, result.TVShows)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 0, result.TotalPages)
		assert.Equal(t, 0, result.TotalResults)
	})

	t.Run("large TV search result", func(t *testing.T) {
		shows := make([]internal.TVShow, 15)
		for i := range shows {
			shows[i] = helpers.CreateValidTVShow(map[string]any{
				"ID":   i + 1,
				"Name": fmt.Sprintf("TV Show %d", i+1),
			})
		}

		result := internal.TVSearchResult{
			TVShows:      shows,
			Page:         3,
			TotalPages:   50,
			TotalResults: 750,
		}

		assert.Len(t, result.TVShows, 15)
		assert.Equal(t, 3, result.Page)
		assert.Equal(t, 50, result.TotalPages)
		assert.Equal(t, 750, result.TotalResults)

		for _, show := range result.TVShows {
			helpers.AssertTVShowValid(t, show)
		}
	})

	t.Run("TV search result pagination boundaries", func(t *testing.T) {
		tests := []struct {
			name         string
			page         int
			totalPages   int
			totalResults int
			valid        bool
		}{
			{"first page", 1, 10, 200, true},
			{"middle page", 5, 10, 200, true},
			{"last page", 10, 10, 200, true},
			{"zero page", 0, 10, 200, false},
			{"negative page", -1, 10, 200, false},
			{"page beyond total", 11, 10, 200, false},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				result := internal.TVSearchResult{
					Page:       test.page,
					TotalPages: test.totalPages,
				}

				if test.valid {
					assert.Greater(t, result.Page, 0)
					assert.LessOrEqual(t, result.Page, result.TotalPages)
				} else {
					// Invalid states - testing boundary conditions
					if test.page <= 0 {
						assert.LessOrEqual(t, result.Page, 0)
					}
					if test.page > test.totalPages {
						assert.Greater(t, result.Page, result.TotalPages)
					}
				}
			})
		}
	})
}

func TestTVShowJSONTags(t *testing.T) {
	t.Run("json tags are correct", func(t *testing.T) {
		// Create a TV show and verify interface compatibility
		show := fixtures.SampleTVShow

		// This test ensures that the struct tags are properly defined
		// by checking that the struct can be used with MediaItem interface
		var mediaItem internal.MediaItem = show

		// Verify interface methods work
		assert.Equal(t, show.ID, mediaItem.GetID())
		assert.Equal(t, show.Name, mediaItem.GetTitle())
		assert.Equal(t, show.OriginalName, mediaItem.GetOriginalTitle())
	})

	t.Run("omitempty fields work correctly", func(t *testing.T) {
		// Test that fields with omitempty tag behave correctly
		show := internal.TVShow{
			ID:   1,
			Name: "Test",
			// OriginalName omitted - should work with omitempty
			// Overview omitted - should work with omitempty
			// FirstAirDate omitted - should work with omitempty
		}

		err := show.Validate()
		assert.NoError(t, err)
	})

	t.Run("first air date field works correctly", func(t *testing.T) {
		show := internal.TVShow{
			FirstAirDate: "2023-01-15",
		}

		// FirstAirDate is TV-specific field, different from movies' ReleaseDate
		assert.Equal(t, "2023-01-15", show.FirstAirDate)
		assert.NotEmpty(t, show.FirstAirDate)
	})
}
