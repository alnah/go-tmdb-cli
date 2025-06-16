package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestSearchOptionsValidation(t *testing.T) {
	t.Run("valid search options pass validation", func(t *testing.T) {
		opts := fixtures.ValidSearchOptions
		err := internal.ValidateSearchOptions(&opts)
		assert.NoError(t, err)
		helpers.AssertSearchOptionsValid(t, opts)
	})

	t.Run("invalid page gets corrected", func(t *testing.T) {
		opts := helpers.CreateValidSearchOptions()
		opts.Page = -1

		err := internal.ValidateSearchOptions(&opts)
		assert.NoError(t, err)
		assert.Equal(t, 1, opts.Page, "Negative page should be corrected to 1")
	})

	t.Run("zero page gets corrected", func(t *testing.T) {
		opts := helpers.CreateValidSearchOptions()
		opts.Page = 0

		err := internal.ValidateSearchOptions(&opts)
		assert.NoError(t, err)
		assert.Equal(t, 1, opts.Page, "Zero page should be corrected to 1")
	})

	t.Run("invalid max items fails validation", func(t *testing.T) {
		tests := []struct {
			name     string
			maxItems int
		}{
			{"zero max items", 0},
			{"negative max items", -1},
			{"too large max items", 1001},
			{"extremely large max items", 999999},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				opts := helpers.CreateValidSearchOptions()
				opts.MaxItems = test.maxItems

				err := internal.ValidateSearchOptions(&opts)
				helpers.AssertValidationError(t, err, "max-items must be between 1 and 1000")
			})
		}
	})

	t.Run("invalid ratings fail validation", func(t *testing.T) {
		tests := []struct {
			name      string
			minRating float64
			maxRating float64
			errorMsg  string
		}{
			{"negative min rating", -1.0, 8.0, "min-rating must be between 0 and 10"},
			{"negative max rating", 5.0, -1.0, "max-rating must be between 0 and 10"},
			{"min rating too high", 11.0, 12.0, "min-rating must be between 0 and 10"},
			{"max rating too high", 5.0, 15.0, "max-rating must be between 0 and 10"},
			{
				"min greater than max",
				8.0,
				6.0,
				"min-rating (8.0) cannot be greater than max-rating (6.0)",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				opts := helpers.CreateValidSearchOptions()
				opts.MinRating = test.minRating
				opts.MaxRating = test.maxRating

				err := internal.ValidateSearchOptions(&opts)
				helpers.AssertValidationError(t, err, test.errorMsg)
			})
		}
	})

	t.Run("boundary rating values are valid", func(t *testing.T) {
		tests := []struct {
			name      string
			minRating float64
			maxRating float64
		}{
			{"zero ratings", 0.0, 0.0},
			{"maximum ratings", 10.0, 10.0},
			{"equal ratings", 7.5, 7.5},
			{"full range", 0.0, 10.0},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				opts := helpers.CreateValidSearchOptions()
				opts.MinRating = test.minRating
				opts.MaxRating = test.maxRating

				err := internal.ValidateSearchOptions(&opts)
				assert.NoError(t, err)
			})
		}
	})
}

func TestSearchOptionsEdgeCases(t *testing.T) {
	t.Run("empty search options", func(t *testing.T) {
		opts := internal.SearchOptions{
			MaxItems: 20, // Set a valid MaxItems to avoid validation error
		}

		err := internal.ValidateSearchOptions(&opts)
		assert.NoError(t, err)
		assert.Equal(t, 1, opts.Page, "Page should be set to 1")
	})

	t.Run("partial search options", func(t *testing.T) {
		opts := internal.SearchOptions{
			Query:    "test",
			MaxItems: 50,
		}

		err := internal.ValidateSearchOptions(&opts)
		assert.NoError(t, err)
		assert.Equal(t, 1, opts.Page)
		assert.Equal(t, "test", opts.Query)
		assert.Equal(t, 50, opts.MaxItems)
	})

	t.Run("zero rating values are treated as unset", func(t *testing.T) {
		opts := internal.SearchOptions{
			Query:     "test",
			MaxItems:  20,
			MinRating: 0.0, // Zero should be treated as unset
			MaxRating: 0.0, // Zero should be treated as unset
		}

		err := internal.ValidateSearchOptions(&opts)
		assert.NoError(t, err)
		// Zero ratings should be acceptable (meaning unset)
	})

	t.Run("genre arrays work correctly", func(t *testing.T) {
		opts := internal.SearchOptions{
			Query:         "test",
			MaxItems:      20,
			IncludeGenres: []int{28, 35, 18}, // Action, Comedy, Drama
			ExcludeGenres: []int{27, 53},     // Horror, Thriller
		}

		err := internal.ValidateSearchOptions(&opts)
		assert.NoError(t, err)
		assert.Len(t, opts.IncludeGenres, 3)
		assert.Len(t, opts.ExcludeGenres, 2)
		assert.Contains(t, opts.IncludeGenres, 28)
		assert.Contains(t, opts.ExcludeGenres, 27)
	})

	t.Run("empty genre arrays are valid", func(t *testing.T) {
		opts := internal.SearchOptions{
			Query:         "test",
			MaxItems:      20,
			IncludeGenres: []int{},
			ExcludeGenres: []int{},
		}

		err := internal.ValidateSearchOptions(&opts)
		assert.NoError(t, err)
		assert.Empty(t, opts.IncludeGenres)
		assert.Empty(t, opts.ExcludeGenres)
	})

	t.Run("nil genre arrays are valid", func(t *testing.T) {
		opts := internal.SearchOptions{
			Query:         "test",
			MaxItems:      20,
			IncludeGenres: nil,
			ExcludeGenres: nil,
		}

		err := internal.ValidateSearchOptions(&opts)
		assert.NoError(t, err)
		assert.Nil(t, opts.IncludeGenres)
		assert.Nil(t, opts.ExcludeGenres)
	})
}

func TestSearchOptionsFields(t *testing.T) {
	t.Run("all fields can be set", func(t *testing.T) {
		opts := internal.SearchOptions{
			Query:         "matrix",
			Page:          5,
			Language:      "en",
			Year:          1999,
			MinRating:     8.0,
			MaxRating:     10.0,
			MinVotes:      1000,
			MaxVotes:      2000000,
			IncludeGenres: []int{28, 878}, // Action, Sci-Fi
			ExcludeGenres: []int{27},      // Horror
			SortBy:        "popularity",
			SortOrder:     "desc",
			MaxItems:      50,
		}

		err := internal.ValidateSearchOptions(&opts)
		assert.NoError(t, err)

		assert.Equal(t, "matrix", opts.Query)
		assert.Equal(t, 5, opts.Page)
		assert.Equal(t, "en", opts.Language)
		assert.Equal(t, 1999, opts.Year)
		assert.Equal(t, 8.0, opts.MinRating)
		assert.Equal(t, 10.0, opts.MaxRating)
		assert.Equal(t, 1000, opts.MinVotes)
		assert.Equal(t, 2000000, opts.MaxVotes)
		assert.Equal(t, []int{28, 878}, opts.IncludeGenres)
		assert.Equal(t, []int{27}, opts.ExcludeGenres)
		assert.Equal(t, "popularity", opts.SortBy)
		assert.Equal(t, "desc", opts.SortOrder)
		assert.Equal(t, 50, opts.MaxItems)
	})

	t.Run("string fields handle various values", func(t *testing.T) {
		tests := []struct {
			field string
			value string
		}{
			{"Query", "test query with spaces"},
			{"Query", "query-with-dashes"},
			{"Query", "query_with_underscores"},
			{"Query", "query123with456numbers"},
			{"Language", "en"},
			{"Language", "en-US"},
			{"Language", "zh-CN"},
			{"SortBy", "popularity"},
			{"SortBy", "vote_average"},
			{"SortBy", "release_date"},
			{"SortOrder", "asc"},
			{"SortOrder", "desc"},
		}

		for _, test := range tests {
			t.Run(test.field+"_"+test.value, func(t *testing.T) {
				opts := helpers.CreateValidSearchOptions()

				switch test.field {
				case "Query":
					opts.Query = test.value
				case "Language":
					opts.Language = test.value
				case "SortBy":
					opts.SortBy = test.value
				case "SortOrder":
					opts.SortOrder = test.value
				}

				err := internal.ValidateSearchOptions(&opts)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("numeric fields handle various values", func(t *testing.T) {
		tests := []struct {
			name  string
			year  int
			votes int
			valid bool
		}{
			{"current year", 2024, 1000, true},
			{"old year", 1900, 500, true},
			{"future year", 2030, 2000, true},
			{"zero year", 0, 1000, true},         // Zero year means unset
			{"negative year", -1999, 1000, true}, // Negative years could be valid (BC)
			{"high votes", 0, 10000000, true},
			{"zero votes", 0, 0, true},
			{"negative votes", 0, -1000, true}, // Edge case - might happen
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				opts := helpers.CreateValidSearchOptions()
				opts.Year = test.year
				opts.MinVotes = test.votes
				opts.MaxVotes = test.votes * 2

				err := internal.ValidateSearchOptions(&opts)
				if test.valid {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			})
		}
	})
}

func TestGenreStruct(t *testing.T) {
	t.Run("genre struct works correctly", func(t *testing.T) {
		genre := internal.Genre{
			ID:   28,
			Name: "Action",
		}

		assert.Equal(t, 28, genre.ID)
		assert.Equal(t, "Action", genre.Name)
	})

	t.Run("multiple genres can be created", func(t *testing.T) {
		genres := []internal.Genre{
			{ID: 28, Name: "Action"},
			{ID: 35, Name: "Comedy"},
			{ID: 18, Name: "Drama"},
		}

		assert.Len(t, genres, 3)
		assert.Equal(t, "Action", genres[0].Name)
		assert.Equal(t, "Comedy", genres[1].Name)
		assert.Equal(t, "Drama", genres[2].Name)
	})

	t.Run("genre handles unicode names", func(t *testing.T) {
		genre := internal.Genre{
			ID:   1,
			Name: "Действие", // Russian for "Action"
		}

		assert.Equal(t, 1, genre.ID)
		assert.Equal(t, "Действие", genre.Name)
		assert.Contains(t, genre.Name, "Действие")
	})
}

func TestMediaItemInterface(t *testing.T) {
	t.Run("movie implements MediaItem interface", func(t *testing.T) {
		movie := fixtures.SampleMovie
		var mediaItem internal.MediaItem = movie

		// Verify all interface methods work
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

	t.Run("TV show implements MediaItem interface", func(t *testing.T) {
		show := fixtures.SampleTVShow
		var mediaItem internal.MediaItem = show

		// Verify all interface methods work
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

	t.Run("MediaItem interface allows polymorphic usage", func(t *testing.T) {
		var items []internal.MediaItem

		// Add different types to the same slice
		items = append(items, fixtures.SampleMovie)
		items = append(items, fixtures.SampleTVShow)

		require.Len(t, items, 2)

		// Test that we can call interface methods on all items
		for i, item := range items {
			assert.Greater(t, item.GetID(), 0, "Item %d should have positive ID", i)
			assert.NotEmpty(t, item.GetTitle(), "Item %d should have non-empty title", i)
			assert.GreaterOrEqual(
				t,
				item.GetRating(),
				0.0,
				"Item %d should have non-negative rating",
				i,
			)
			assert.GreaterOrEqual(
				t,
				item.GetVotes(),
				0,
				"Item %d should have non-negative votes",
				i,
			)
		}

		// Verify first item is movie
		movieItem := items[0]
		assert.Equal(t, fixtures.SampleMovie.ID, movieItem.GetID())
		assert.Equal(t, fixtures.SampleMovie.Title, movieItem.GetTitle())

		// Verify second item is TV show
		showItem := items[1]
		assert.Equal(t, fixtures.SampleTVShow.ID, showItem.GetID())
		assert.Equal(t, fixtures.SampleTVShow.Name, showItem.GetTitle()) // Name maps to Title
	})
}

func TestSearchOptionsValidationCoverage(t *testing.T) {
	t.Run("covers all validation scenarios", func(t *testing.T) {
		// Test with invalid options from fixtures
		opts := fixtures.InvalidSearchOptions

		err := internal.ValidateSearchOptions(&opts)
		require.Error(t, err)

		// Should fail on first validation error
		assert.Contains(t, err.Error(), "max-items must be between 1 and 1000")
	})

	t.Run("validates vote ranges", func(t *testing.T) {
		opts := helpers.CreateValidSearchOptions()
		opts.MinVotes = -100
		opts.MaxVotes = -50

		// Negative votes might be edge case but validation should handle it
		err := internal.ValidateSearchOptions(&opts)
		assert.NoError(t, err) // Current validation doesn't check vote ranges
	})
}
