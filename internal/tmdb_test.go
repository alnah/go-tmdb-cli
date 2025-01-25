package tmdb_test

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	i "github.com/alnah/go-tmdb-cli/internal"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUnitGetTMDBToken(t *testing.T) {
	t.Run("returns TMDB API Token from env file", func(t *testing.T) {
		file := createTempFile(t)

		_, err := file.WriteString(`TOKEN="test"`)
		assert.NoError(t, err)

		got, err := i.GetTMDBToken(file.Name())
		assert.NoError(t, err)
		assert.Equal(t, "test", got)
	})

	t.Run("returns path error when env file doesn't exist", func(t *testing.T) {
		var want *fs.PathError
		_, err := i.GetTMDBToken("test.env")
		assert.ErrorAs(t, err, &want)
	})

	t.Run("returns config parse error when invalid data", func(t *testing.T) {
		file := createTempFile(t)

		_, err := file.WriteString(`invalid_data`)
		assert.NoError(t, err)

		var want viper.ConfigParseError
		_, err = i.GetTMDBToken(file.Name())
		assert.ErrorAs(t, err, &want)
	})
}

func TestUnitBuildURL(t *testing.T) {
	t.Run("return query", func(t *testing.T) {
		testCases := []struct {
			desc  string
			query *i.QueryParams
			want  string
		}{
			{
				desc: "page year date average",
				query: &i.QueryParams{
					Page:    toPtr(2),
					Year:    toPtr(2000),
					Date:    &i.Date{Value: "2000-06-01", Option: "gte"},
					Average: &i.Average{Value: 8.0, Option: "gte"},
					Vote:    &i.Vote{Value: 1000, Option: "gte"},
				},
				want: i.BaseURL + i.DiscoverURL +
					"page=2" +
					"&primary_release_year=2000" +
					"&primary_release_date.gte=2000-06-01" +
					"&vote_average.gte=8.0&vote_count.gte=1000",
			},
			{
				desc: "popular movies list",
				query: &i.QueryParams{
					MovieListPath: toPtr("popular"),
				},
				want: fmt.Sprintf(i.BaseURL+i.MoviesListURL, "popular"),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.desc, func(t *testing.T) {
				got, err := tc.query.BuildURL()
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("return filter error when error occured", func(t *testing.T) {
		testCases := []struct {
			desc  string
			query *i.QueryParams
		}{
			{
				desc: "on discover",
				query: &i.QueryParams{
					Date: &i.Date{
						Value:  "1800-01-01", // movies don't exist at that time
						Option: "gte",
					},
				},
			},
			{
				desc: "on movies list",
				query: &i.QueryParams{
					MovieListPath: toPtr("invalid"),
				},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.desc, func(t *testing.T) {
				var want *i.FilterError
				got, err := tc.query.BuildURL()
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			})
		}
	})
}

func TestUnitFilterError(t *testing.T) {
	t.Run("return filter error", func(t *testing.T) {
		err := &i.FilterError{Filter: "Test", Message: "Description"}
		assert.Equal(t, "Test: Description.", err.Error())
	})
}

func TestUnitSetMoviesList(t *testing.T) {
	testCases := []struct {
		desc    string
		path    string
		wantErr bool
	}{
		{
			desc: "return now playing",
			path: "now_playing",
		},
		{
			desc: "return popular",
			path: "popular",
		},
		{
			desc: "return top rated",
			path: "top_rated",
		},
		{
			desc: "return upcoming",
			path: "upcoming",
		},
		{
			desc:    "return filter error when doesn't match wanted paths",
			path:    "invalid",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			switch {
			case tc.wantErr:
				query := i.QueryParams{MovieListPath: &tc.path}
				var want *i.FilterError
				got, err := query.SetMoviesList()
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				query := i.QueryParams{MovieListPath: &tc.path}
				want := fmt.Sprintf(i.MoviesListURL, tc.path)
				got, err := query.SetMoviesList()
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}

	t.Run("return empty string when no movies list", func(t *testing.T) {
		query := &i.QueryParams{MovieListPath: nil}
		got, err := query.SetMoviesList()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func TestUnitSetPageFilter(t *testing.T) {
	testCases := []struct {
		desc    string
		value   *int
		want    string
		wantErr bool
	}{
		{
			desc:  "return page",
			value: toPtr(20),
			want:  "page=20",
		},
		{
			desc:    "return empty string if no page",
			value:   nil,
			wantErr: false,
		},
		{
			desc:    "return filter error if out of range",
			value:   toPtr(2147483648), // out of range
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &i.QueryParams{Page: tc.value}
			got, err := query.SetPageFilter()
			switch {
			case tc.wantErr:
				var want *i.FilterError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestUnitSetYearFilter(t *testing.T) {
	testCases := []struct {
		desc    string
		value   *int
		want    string
		wantErr bool
	}{
		{
			desc:  "return year",
			value: toPtr(2001),
			want:  "primary_release_year=2001",
		},
		{
			desc:    "return empty string when no year",
			value:   nil,
			wantErr: false,
		},
		{
			desc:    "return filter error when year out of range",
			value:   toPtr(2147483648), // out of range
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &i.QueryParams{Year: tc.value}
			got, err := query.SetYearFilter()
			switch {
			case tc.wantErr:
				var want *i.FilterError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestUnitSetDateFilter(t *testing.T) {
	testCases := []struct {
		desc    string
		date    i.Date
		want    string
		wantErr bool
	}{
		{
			desc: "return date gte",
			date: i.Date{
				Value:  "2001-01-01",
				Option: "gte",
			},
			want: "primary_release_date.gte=2001-01-01",
		},
		{
			desc: "return date lte",
			date: i.Date{
				Value:  "2001-01-01",
				Option: "lte",
			},
			want: "primary_release_date.lte=2001-01-01",
		},
		{
			desc: "return filter error when invalid date value format",
			date: i.Date{
				Value:  "01-01-2001", // invalid date format
				Option: "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid date value range",
			date: i.Date{
				Value:  "1800-01-01", // movies don't exist at that time
				Option: "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid date option",
			date: i.Date{
				Value:  "2001-01-01",
				Option: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &i.QueryParams{Date: &tc.date}
			got, err := query.SetDateFilter()
			switch {
			case tc.wantErr:
				var want *i.FilterError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

	t.Run("return empty string when no date", func(t *testing.T) {
		query := &i.QueryParams{Date: nil}
		got, err := query.SetDateFilter()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func TestUnitSetAverageFilter(t *testing.T) {
	testCases := []struct {
		desc    string
		average i.Average
		want    string
		wantErr bool
	}{
		{
			desc: "return average gte",
			average: i.Average{
				Value:  8.0,
				Option: "gte",
			},
			want: "vote_average.gte=8.0",
		},
		{
			desc: "return average lte",
			average: i.Average{
				Value:  8.0,
				Option: "lte",
			},
			want: "vote_average.lte=8.0",
		},
		{
			desc: "return filter error when invalid option",
			average: i.Average{
				Value:  8.0,
				Option: "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when value out of range",
			average: i.Average{
				Value:  100.0,
				Option: "gte",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &i.QueryParams{Average: &tc.average}
			got, err := query.SetAverageFilter()
			switch {
			case tc.wantErr:
				var want *i.FilterError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

	t.Run("return empty string when no average", func(t *testing.T) {
		query := &i.QueryParams{Average: nil}
		got, err := query.SetAverageFilter()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func TestUnitSetVoteFilter(t *testing.T) {
	testCases := []struct {
		desc    string
		vote    i.Vote
		want    string
		wantErr bool
	}{
		{
			desc: "return vote gte",
			vote: i.Vote{
				Value:  1000,
				Option: "gte",
			},
			want: "vote_count.gte=1000",
		},
		{
			desc: "return vote lte",
			vote: i.Vote{
				Value:  1000,
				Option: "lte",
			},
			want: "vote_count.lte=1000",
		},
		{
			desc: "return filter error when invalid option",
			vote: i.Vote{
				Value:  1000,
				Option: "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when value out of range",
			vote: i.Vote{
				Value:  2147483648, // out of range,
				Option: "gte",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &i.QueryParams{Vote: &tc.vote}
			got, err := query.SetVoteFilter()
			switch {
			case tc.wantErr:
				var want *i.FilterError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

	t.Run("return empty string when no date", func(t *testing.T) {
		query := &i.QueryParams{Average: nil}
		got, err := query.SetVoteFilter()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func createTempFile(t *testing.T) *os.File {
	t.Helper()

	file, err := os.CreateTemp("", "test.env")
	assert.NoError(t, err)
	t.Cleanup(func() {
		file.Close()
		os.Remove(file.Name())
	})

	return file
}

func toPtr[T any](value T) *T {
	return &value
}
