package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestBuildBaseDiscoverParams_BasicFilters(t *testing.T) {
	t.Parallel()

	config := helpers.MockValidConfig()
	client := internal.NewClient(config)

	opts := internal.SearchOptions{
		Language:  "fr",
		MinRating: 7.5,
		MaxRating: 9.5,
		MinVotes:  1000,
		MaxVotes:  100000,
	}

	params := client.BuildBaseDiscoverParams(opts)

	assert.Equal(t, "fr", params.Get("with_original_language"))
	assert.Equal(t, "7.5", params.Get("vote_average.gte"))
	assert.Equal(t, "9.5", params.Get("vote_average.lte"))
	assert.Equal(t, "1000", params.Get("vote_count.gte"))
	assert.Equal(t, "100000", params.Get("vote_count.lte"))
}

func TestBuildBaseDiscoverParams_Genres(t *testing.T) {
	t.Parallel()

	config := helpers.MockValidConfig()
	client := internal.NewClient(config)

	opts := internal.SearchOptions{
		IncludeGenres: []int{28, 12, 16},
		ExcludeGenres: []int{27, 53},
	}

	params := client.BuildBaseDiscoverParams(opts)

	assert.Equal(t, "28,12,16", params.Get("with_genres"))
	assert.Equal(t, "27,53", params.Get("without_genres"))
}

func TestBuildBaseDiscoverParams_Sorting(t *testing.T) {
	t.Parallel()

	config := helpers.MockValidConfig()
	client := internal.NewClient(config)

	tests := []struct {
		sortBy    string
		sortOrder string
		expected  string
	}{
		{"popularity", "desc", "popularity.desc"},
		{"votes", "asc", "vote_count.asc"}, // votes -> vote_count
		{"", "", "popularity.desc"},        // defaults
	}

	for _, tt := range tests {
		opts := internal.SearchOptions{
			SortBy:    tt.sortBy,
			SortOrder: tt.sortOrder,
		}

		params := client.BuildBaseDiscoverParams(opts)
		assert.Equal(t, tt.expected, params.Get("sort_by"))
	}
}

func TestBuildBaseDiscoverParams_ZeroValues(t *testing.T) {
	t.Parallel()

	config := helpers.MockValidConfig()
	client := internal.NewClient(config)

	opts := internal.SearchOptions{
		MinRating: 0,
		MaxRating: 0,
		MinVotes:  0,
		MaxVotes:  0,
	}

	params := client.BuildBaseDiscoverParams(opts)

	assert.Empty(t, params.Get("vote_average.gte"))
	assert.Empty(t, params.Get("vote_average.lte"))
	assert.Empty(t, params.Get("vote_count.gte"))
	assert.Empty(t, params.Get("vote_count.lte"))
}

func TestSearchOptionsBuilder(t *testing.T) {
	t.Parallel()

	opts := helpers.NewSearchOptionsBuilder().
		WithMaxItems(50).
		WithLanguage("fr").
		WithMinRating(7.0).
		WithMaxRating(9.0).
		WithMinVotes(1000).
		WithMaxVotes(50000).
		WithYear(2023).
		WithIncludeGenres([]int{28, 12}).
		WithExcludeGenres([]int{27}).
		WithSortBy("vote_average").
		WithSortOrder("desc").
		Build()

	assert.Equal(t, 50, opts.MaxItems)
	assert.Equal(t, "fr", opts.Language)
	assert.Equal(t, 7.0, opts.MinRating)
	assert.Equal(t, 9.0, opts.MaxRating)
	assert.Equal(t, 1000, opts.MinVotes)
	assert.Equal(t, 50000, opts.MaxVotes)
	assert.Equal(t, 2023, opts.Year)
	assert.Equal(t, []int{28, 12}, opts.IncludeGenres)
	assert.Equal(t, []int{27}, opts.ExcludeGenres)
	assert.Equal(t, "vote_average", opts.SortBy)
	assert.Equal(t, "desc", opts.SortOrder)
}

func TestSearchOptionsBuilder_Minimal(t *testing.T) {
	t.Parallel()

	opts := helpers.NewSearchOptionsBuilder().
		WithMaxItems(20).
		Build()

	assert.Equal(t, 20, opts.MaxItems)
	assert.Empty(t, opts.Language)
	assert.Equal(t, 0.0, opts.MinRating)
	assert.Equal(t, 10.0, opts.MaxRating) // Default MaxRating is 10.0
	assert.Equal(t, 0, opts.MinVotes)
	assert.Equal(t, 0, opts.MaxVotes)
	assert.Equal(t, 0, opts.Year)
	assert.Nil(t, opts.IncludeGenres)
	assert.Nil(t, opts.ExcludeGenres)
	assert.Empty(t, opts.SortBy)
	assert.Empty(t, opts.SortOrder)
}
