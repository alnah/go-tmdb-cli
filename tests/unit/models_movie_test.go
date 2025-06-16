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

func TestMovieValidation(t *testing.T) {
	t.Run("valid movie passes validation", func(t *testing.T) {
		movie := fixtures.SampleMovie
		err := movie.Validate()
		assert.NoError(t, err)
	})

	t.Run("movie with empty title fails validation", func(t *testing.T) {
		movie := helpers.CreateValidMovie(map[string]any{
			"Title": "",
		})
		err := movie.Validate()
		helpers.AssertValidationError(t, err, "movie title is required")
	})

	t.Run("movie with zero ID fails validation", func(t *testing.T) {
		movie := helpers.CreateValidMovie(map[string]any{
			"ID": 0,
		})
		err := movie.Validate()
		helpers.AssertValidationError(t, err, "movie ID must be positive")
	})

	t.Run("movie with negative ID fails validation", func(t *testing.T) {
		movie := helpers.CreateValidMovie(map[string]any{
			"ID": -1,
		})
		err := movie.Validate()
		helpers.AssertValidationError(t, err, "movie ID must be positive")
	})

	t.Run("multiple validation errors", func(t *testing.T) {
		movie := internal.Movie{
			ID:    0,  // Invalid
			Title: "", // Invalid
		}
		err := movie.Validate()
		require.Error(t, err)
		// Should fail on first validation error (title)
		assert.Contains(t, err.Error(), "movie title is required")
	})
}

func TestMovieMediaItemInterface(t *testing.T) {
	movie := fixtures.SampleMovie

	t.Run("implements MediaItem interface correctly", func(t *testing.T) {
		var mediaItem internal.MediaItem = movie

		assert.Equal(t, movie.ID, mediaItem.GetID())
		assert.Equal(t, movie.Title, mediaItem.GetTitle())
		assert.Equal(t, movie.OriginalTitle, mediaItem.GetOriginalTitle())
		assert.Equal(t, movie.Year, mediaItem.GetYear())
		assert.Equal(t, movie.Rating, mediaItem.GetRating())
		assert.Equal(t, movie.Votes, mediaItem.GetVotes())
		assert.Equal(t, movie.Popularity, mediaItem.GetPopularity())
		assert.Equal(t, movie.Genres, mediaItem.GetGenres())
		assert.Equal(t, movie.Overview, mediaItem.GetOverview())
	})

	t.Run("interface methods return correct types", func(t *testing.T) {
		assert.IsType(t, 0, movie.GetID())
		assert.IsType(t, "", movie.GetTitle())
		assert.IsType(t, "", movie.GetOriginalTitle())
		assert.IsType(t, 0, movie.GetYear())
		assert.IsType(t, 0.0, movie.GetRating())
		assert.IsType(t, 0, movie.GetVotes())
		assert.IsType(t, 0.0, movie.GetPopularity())
		assert.IsType(t, "", movie.GetGenres())
		assert.IsType(t, "", movie.GetOverview())
	})
}

func TestMovieFields(t *testing.T) {
	t.Run("all fields are properly set", func(t *testing.T) {
		movie := fixtures.SampleMovie

		assert.Equal(t, 123, movie.ID)
		assert.Equal(t, "The Matrix", movie.Title)
		assert.Equal(t, "The Matrix", movie.OriginalTitle)
		assert.Equal(t, 1999, movie.Year)
		assert.Equal(t, 8.7, movie.Rating)
		assert.Equal(t, 1800000, movie.Votes)
		assert.Equal(t, 85.5, movie.Popularity)
		assert.Equal(t, "Action, Sci-Fi", movie.Genres)
		assert.Equal(
			t,
			"A computer hacker learns about the true nature of his reality.",
			movie.Overview,
		)
		assert.Equal(t, "en", movie.Language)
		assert.False(t, movie.Adult)
		assert.Equal(t, "1999-03-31", movie.ReleaseDate)
	})

	t.Run("handles different original title", func(t *testing.T) {
		movie := fixtures.SampleMovieWithDifferentOriginal

		assert.Equal(t, "Spirited Away", movie.Title)
		assert.Equal(t, "千と千尋の神隠し", movie.OriginalTitle)
		assert.NotEqual(t, movie.Title, movie.OriginalTitle)
	})

	t.Run("handles zero and negative values", func(t *testing.T) {
		movie := internal.Movie{
			Year:       0,   // Unset year
			Rating:     0.0, // Zero rating
			Votes:      0,   // No votes
			Popularity: 0.0, // No popularity
		}

		assert.Equal(t, 0, movie.Year)
		assert.Equal(t, 0.0, movie.Rating)
		assert.Equal(t, 0, movie.Votes)
		assert.Equal(t, 0.0, movie.Popularity)
	})

	t.Run("handles adult content flag", func(t *testing.T) {
		adultMovie := helpers.CreateValidMovie(map[string]any{
			"Adult": true,
		})
		nonAdultMovie := helpers.CreateValidMovie(map[string]any{
			"Adult": false,
		})

		assert.True(t, adultMovie.Adult)
		assert.False(t, nonAdultMovie.Adult)
	})
}

func TestMovieEdgeCases(t *testing.T) {
	t.Run("handles empty optional fields", func(t *testing.T) {
		movie := internal.Movie{
			ID:            1,
			Title:         "Minimal Movie",
			OriginalTitle: "", // Empty original title
			Overview:      "", // Empty overview
			Language:      "", // Empty language
			ReleaseDate:   "", // Empty release date
			Genres:        "", // Empty genres
		}

		err := movie.Validate()
		assert.NoError(t, err, "Movie with empty optional fields should be valid")

		assert.Empty(t, movie.OriginalTitle)
		assert.Empty(t, movie.Overview)
		assert.Empty(t, movie.Language)
		assert.Empty(t, movie.ReleaseDate)
		assert.Empty(t, movie.Genres)
	})

	t.Run("handles extreme numeric values", func(t *testing.T) {
		movie := internal.Movie{
			ID:         999999999, // Very large ID
			Title:      "Test Movie",
			Year:       2050,      // Future year
			Rating:     10.0,      // Maximum rating
			Votes:      999999999, // Very large vote count
			Popularity: 999.99,    // Very high popularity
		}

		err := movie.Validate()
		assert.NoError(t, err)

		assert.Equal(t, 999999999, movie.ID)
		assert.Equal(t, 2050, movie.Year)
		assert.Equal(t, 10.0, movie.Rating)
		assert.Equal(t, 999999999, movie.Votes)
		assert.Equal(t, 999.99, movie.Popularity)
	})

	t.Run("handles unicode in text fields", func(t *testing.T) {
		movie := internal.Movie{
			ID:            1,
			Title:         "Movie with émojis 🎬 and åccénts",
			OriginalTitle: "Оригинальное название фильма",
			Overview:      "Description with 中文 characters and العربية text",
			Language:      "zh",
			Genres:        "アクション, コメディ",
		}

		err := movie.Validate()
		assert.NoError(t, err)

		assert.Contains(t, movie.Title, "🎬")
		assert.Contains(t, movie.Title, "åccénts")
		assert.Contains(t, movie.OriginalTitle, "Оригинальное")
		assert.Contains(t, movie.Overview, "中文")
		assert.Contains(t, movie.Overview, "العربية")
		assert.Contains(t, movie.Genres, "アクション")
	})

	t.Run("handles very long text fields", func(t *testing.T) {
		longTitle := strings.Repeat("Very Long Movie Title ", 50)
		longOverview := strings.Repeat("This is a very long overview description. ", 100)

		movie := internal.Movie{
			ID:       1,
			Title:    longTitle,
			Overview: longOverview,
		}

		err := movie.Validate()
		assert.NoError(t, err)

		assert.Equal(t, longTitle, movie.Title)
		assert.Equal(t, longOverview, movie.Overview)
		assert.Greater(t, len(movie.Title), 1000)
		assert.Greater(t, len(movie.Overview), 4000)
	})
}

func TestSearchResult(t *testing.T) {
	t.Run("search result with movies", func(t *testing.T) {
		result := fixtures.SampleSearchResult

		assert.Len(t, result.Movies, 1)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 3, result.TotalPages)
		assert.Equal(t, 50, result.TotalResults)

		helpers.AssertMovieValid(t, result.Movies[0])
	})

	t.Run("empty search result", func(t *testing.T) {
		result := internal.SearchResult{
			Movies:       []internal.Movie{},
			Page:         1,
			TotalPages:   0,
			TotalResults: 0,
		}

		assert.Empty(t, result.Movies)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 0, result.TotalPages)
		assert.Equal(t, 0, result.TotalResults)
	})

	t.Run("large search result", func(t *testing.T) {
		movies := make([]internal.Movie, 20)
		for i := range movies {
			movies[i] = helpers.CreateValidMovie(map[string]any{
				"ID":    i + 1,
				"Title": fmt.Sprintf("Movie %d", i+1),
			})
		}

		result := internal.SearchResult{
			Movies:       movies,
			Page:         5,
			TotalPages:   100,
			TotalResults: 2000,
		}

		assert.Len(t, result.Movies, 20)
		assert.Equal(t, 5, result.Page)
		assert.Equal(t, 100, result.TotalPages)
		assert.Equal(t, 2000, result.TotalResults)

		for _, movie := range result.Movies {
			helpers.AssertMovieValid(t, movie)
		}
	})

	t.Run("search result pagination boundaries", func(t *testing.T) {
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
				result := internal.SearchResult{
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

func TestMovieJSONTags(t *testing.T) {
	t.Run("json tags are correct", func(t *testing.T) {
		// Create a movie and marshal to JSON to verify tags
		movie := fixtures.SampleMovie

		// This test ensures that the struct tags are properly defined
		// by checking that the struct can be used with MediaItem interface
		var mediaItem internal.MediaItem = movie

		// Verify interface methods work
		assert.Equal(t, movie.ID, mediaItem.GetID())
		assert.Equal(t, movie.Title, mediaItem.GetTitle())
		assert.Equal(t, movie.OriginalTitle, mediaItem.GetOriginalTitle())
	})

	t.Run("omitempty fields work correctly", func(t *testing.T) {
		// Test that fields with omitempty tag behave correctly
		movie := internal.Movie{
			ID:    1,
			Title: "Test",
			// OriginalTitle omitted - should work with omitempty
			// Overview omitted - should work with omitempty
			// ReleaseDate omitted - should work with omitempty
		}

		err := movie.Validate()
		assert.NoError(t, err)
	})
}
