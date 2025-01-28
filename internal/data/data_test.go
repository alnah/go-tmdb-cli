package data_test

import (
	"testing"

	d "github.com/alnah/go-tmdb-cli/internal/data"
	tf "github.com/alnah/go-tmdb-cli/test/fake"

	"github.com/stretchr/testify/assert"
)

var (
	fakeMovie1 = tf.FakeMovies[0]
	fakeMovie2 = tf.FakeMovies[1]
	fakeMovie3 = tf.FakeMovies[2]
)

func TestUnitSortError(t *testing.T) {
	t.Run("return error message", func(t *testing.T) {
		err := &d.OrderError{}
		assert.Equal(t, `Sorting order must be "asc" or "desc".`, err.Error())
	})
}

func TestUnitMoviesSort(t *testing.T) {
	t.Run("return sorted movies", func(t *testing.T) {
		fakeMovies := d.Movies{fakeMovie1, fakeMovie2, fakeMovie3}
		got, err := fakeMovies.Sort()
		assert.NoError(t, err)
		assert.Equal(t, fakeMovies, got)
	})
}

func TestUnitMoviesSortByDate(t *testing.T) {
	fakeMovies := d.Movies{fakeMovie1, fakeMovie2, fakeMovie3}

	testCases := []struct {
		description string
		order       string
		want        d.Movies
		wantErr     bool
	}{
		{
			description: "ascending order",
			order:       "asc",
			want:        d.Movies{fakeMovie1, fakeMovie2, fakeMovie3},
		},
		{
			description: "descending order",
			order:       "desc",
			want:        d.Movies{fakeMovie3, fakeMovie2, fakeMovie1},
		},
		{
			description: "sort error when invalid order",
			order:       "invalid",
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := fakeMovies.SortByDate(tc.order)
			switch {
			case tc.wantErr:
				var want *d.OrderError
				assert.Equal(t, fakeMovies, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestUnitMoviesSortByOriginalTitle(t *testing.T) {
	fakeMovies := d.Movies{fakeMovie1, fakeMovie2, fakeMovie3}

	testCases := []struct {
		description string
		order       string
		want        d.Movies
		wantErr     bool
	}{
		{
			description: "ascending order",
			order:       "asc",
			want:        d.Movies{fakeMovie1, fakeMovie3, fakeMovie2},
		},
		{
			description: "descending order",
			order:       "desc",
			want:        d.Movies{fakeMovie2, fakeMovie3, fakeMovie1},
		},
		{
			description: "sort error when invalid order",
			order:       "invalid",
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := fakeMovies.SortByOriginalTitle(tc.order)
			switch {
			case tc.wantErr:
				var want *d.OrderError
				assert.Equal(t, fakeMovies, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestUnitMoviesSortByTitle(t *testing.T) {
	fakeMovies := d.Movies{fakeMovie1, fakeMovie2, fakeMovie3}

	testCases := []struct {
		description string
		order       string
		want        d.Movies
		wantErr     bool
	}{
		{
			description: "ascending order",
			order:       "asc",
			want:        d.Movies{fakeMovie3, fakeMovie1, fakeMovie2},
		},
		{
			description: "descending order",
			order:       "desc",
			want:        d.Movies{fakeMovie2, fakeMovie1, fakeMovie3},
		},
		{
			description: "sort error when invalid order",
			order:       "invalid",
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := fakeMovies.SortByTitle(tc.order)
			switch {
			case tc.wantErr:
				var want *d.OrderError
				assert.Equal(t, fakeMovies, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestUnitMoviesSortByAverage(t *testing.T) {
	fakeMovies := d.Movies{fakeMovie1, fakeMovie2, fakeMovie3}

	testCases := []struct {
		description string
		order       string
		want        d.Movies
		wantErr     bool
	}{
		{
			description: "ascending order",
			order:       "asc",
			want:        d.Movies{fakeMovie2, fakeMovie1, fakeMovie3},
		},
		{
			description: "descending order",
			order:       "desc",
			want:        d.Movies{fakeMovie3, fakeMovie1, fakeMovie2},
		},
		{
			description: "sort error when invalid order",
			order:       "invalid",
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := fakeMovies.SortByAverage(tc.order)
			switch {
			case tc.wantErr:
				var want *d.OrderError
				assert.Equal(t, fakeMovies, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestUnitMoviesSortByVotes(t *testing.T) {
	fakeMovies := d.Movies{fakeMovie1, fakeMovie2, fakeMovie3}

	testCases := []struct {
		description string
		order       string
		want        d.Movies
		wantErr     bool
	}{
		{
			description: "ascending order",
			order:       "asc",
			want:        d.Movies{fakeMovie2, fakeMovie1, fakeMovie3},
		},
		{
			description: "descending order",
			order:       "desc",
			want:        d.Movies{fakeMovie3, fakeMovie1, fakeMovie2},
		},
		{
			description: "sort error when invalid order",
			order:       "invalid",
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			got, err := fakeMovies.SortByVotes(tc.order)
			switch {
			case tc.wantErr:
				var want *d.OrderError
				assert.Equal(t, fakeMovies, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
