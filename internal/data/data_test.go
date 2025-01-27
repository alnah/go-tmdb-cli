package data_test

import (
	"testing"

	d "github.com/alnah/go-tmdb-cli/internal/data"
	"github.com/stretchr/testify/assert"
)

var (
	fakeMovie1 = d.Movie{
		Id:            1,
		Date:          "2023-01-01",
		OriginalTitle: "L'Aube de l'Aventure",
		Title:         "Epic Journey Begins",
		Average:       8.5,
		Votes:         100,
	}

	fakeMovie2 = d.Movie{
		Id:            2,
		Date:          "2023-02-01",
		OriginalTitle: "Rise of the Heroes",
		Title:         "Rise of the Heroes",
		Average:       7.0,
		Votes:         50,
	}

	fakeMovie3 = d.Movie{
		Id:            3,
		Date:          "2023-03-01",
		OriginalTitle: "O Confronto Final",
		Title:         "Clash of Titans",
		Average:       9.0,
		Votes:         200,
	}
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
