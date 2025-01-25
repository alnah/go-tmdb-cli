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
					Page:    2,
					Year:    2000,
					Date:    &i.Date{StartDate: "2000-06-01", StartOption: "gte"},
					Average: &i.Average{StartAverage: 8.0, StartOption: "gte"},
					Vote:    &i.Vote{StartVotes: 1000, StartOption: "gte"},
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
					MovieListPath: "popular",
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
						StartDate:   "1800-01-01", // movies don't exist at that time
						StartOption: "gte",
					},
				},
			},
			{
				desc: "on movies list",
				query: &i.QueryParams{
					MovieListPath: "invalid",
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
				query := i.QueryParams{MovieListPath: tc.path}
				var want *i.FilterError
				got, err := query.SetMoviesList()
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				query := i.QueryParams{MovieListPath: tc.path}
				want := fmt.Sprintf(i.MoviesListURL, tc.path)
				got, err := query.SetMoviesList()
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}

	t.Run("return empty string when no movies list", func(t *testing.T) {
		query := &i.QueryParams{MovieListPath: ""}
		got, err := query.SetMoviesList()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func TestUnitSetPageFilter(t *testing.T) {
	testCases := []struct {
		desc    string
		value   int
		want    string
		wantErr bool
	}{
		{
			desc:  "return page",
			value: 20,
			want:  "page=20",
		},
		{
			desc:    "return empty string if no page",
			wantErr: false,
		},
		{
			desc:    "return filter error if out of range",
			value:   2147483648, // out of range
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
		value   int
		want    string
		wantErr bool
	}{
		{
			desc:  "return year",
			value: 2001,
			want:  "primary_release_year=2001",
		},
		{
			desc:    "return empty string when no year",
			wantErr: false,
		},
		{
			desc:    "return filter error when year out of range",
			value:   2147483648, // out of range
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
			desc: "return start date gte",
			date: i.Date{
				StartDate:   "2001-01-01",
				StartOption: "gte",
			},
			want: "primary_release_date.gte=2001-01-01",
		},
		{
			desc: "return start date lte",
			date: i.Date{
				StartDate:   "2001-01-01",
				StartOption: "lte",
			},
			want: "primary_release_date.lte=2001-01-01",
		},
		{
			desc: "return start date gte and end date lte",
			date: i.Date{
				StartDate:   "2001-01-01",
				StartOption: "gte",
				EndDate:     "2002-02-02",
				EndOption:   "lte",
			},
			want: "primary_release_date.gte=2001-01-01&" +
				"primary_release_date.lte=2002-02-02",
		},
		{
			desc: "return start date lte and end date gte",
			date: i.Date{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "2002-02-02",
				EndOption:   "gte",
			},
			want: "primary_release_date.lte=2001-01-01&" +
				"primary_release_date.gte=2002-02-02",
		},
		{
			desc: "return filter error when invalid start date format",
			date: i.Date{
				StartDate:   "01-01-2001", // invalid date format
				StartOption: "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid start date range",
			date: i.Date{
				StartDate:   "1800-01-01", // movies don't exist at that time
				StartOption: "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid start date option",
			date: i.Date{
				StartDate:   "2001-01-01",
				StartOption: "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid end date format",
			date: i.Date{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "01-01-2001", // invalid date format
				EndOption:   "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid end date range",
			date: i.Date{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "1800-01-01", // movies don't exist at that time
				EndOption:   "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid end date option",
			date: i.Date{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "2002-02-02",
				EndOption:   "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when same options",
			date: i.Date{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "2002-02-02",
				EndOption:   "lte",
			},
			wantErr: true,
		},

		{
			desc: "return filter error when same dates",
			date: i.Date{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "2001-01-01",
				EndOption:   "gte",
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
			desc: "return start average gte",
			average: i.Average{
				StartAverage: 8.0,
				StartOption:  "gte",
			},
			want: "vote_average.gte=8.0",
		},
		{
			desc: "return start average lte",
			average: i.Average{
				StartAverage: 8.0,
				StartOption:  "lte",
			},
			want: "vote_average.lte=8.0",
		},
		{
			desc: "return start average gte and end average lte",
			average: i.Average{
				StartAverage: 7.0,
				StartOption:  "gte",
				EndAverage:   8.0,
				EndOption:    "lte",
			},
			want: "vote_average.gte=7.0&" +
				"vote_average.lte=8.0",
		},
		{
			desc: "return start average lte and end average gte",
			average: i.Average{
				StartAverage: 7.0,
				StartOption:  "lte",
				EndAverage:   8.0,
				EndOption:    "gte",
			},
			want: "vote_average.lte=7.0&" +
				"vote_average.gte=8.0",
		},
		{
			desc: "return filter error when invalid start option",
			average: i.Average{
				StartAverage: 8.0,
				StartOption:  "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when start average out of range",
			average: i.Average{
				StartAverage: 100.0,
				StartOption:  "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when end average out of range",
			average: i.Average{
				StartAverage: 7.0,
				StartOption:  "lte",
				EndAverage:   100.0, // invalid average
				EndOption:    "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid end option",
			average: i.Average{
				StartAverage: 7.0,
				StartOption:  "lte",
				EndAverage:   8.0,
				EndOption:    "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when same options",
			average: i.Average{
				StartAverage: 7.0,
				StartOption:  "lte",
				EndAverage:   8.0,
				EndOption:    "lte",
			},
			wantErr: true,
		},

		{
			desc: "return filter error when same averages",
			average: i.Average{
				StartAverage: 8.0,
				StartOption:  "lte",
				EndAverage:   8.0,
				EndOption:    "gte",
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
		votes   i.Vote
		want    string
		wantErr bool
	}{
		{
			desc: "return start votes gte",
			votes: i.Vote{
				StartVotes:  1000,
				StartOption: "gte",
			},
			want: "vote_count.gte=1000",
		},
		{
			desc: "return start votes lte",
			votes: i.Vote{
				StartVotes:  1000,
				StartOption: "lte",
			},
			want: "vote_count.lte=1000",
		},
		{
			desc: "return start votes gte and end votes lte",
			votes: i.Vote{
				StartVotes:  500,
				StartOption: "gte",
				EndVotes:    1000,
				EndOption:   "lte",
			},
			want: "vote_count.gte=500&" +
				"vote_count.lte=1000",
		},
		{
			desc: "return start votes lte and end votes gte",
			votes: i.Vote{
				StartVotes:  500,
				StartOption: "lte",
				EndVotes:    1000,
				EndOption:   "gte",
			},
			want: "vote_count.lte=500&" +
				"vote_count.gte=1000",
		},
		{
			desc: "return filter error when invalid start option",
			votes: i.Vote{
				StartVotes:  1000,
				StartOption: "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when start votes out of range",
			votes: i.Vote{
				StartVotes:  2147483648, // invalid votes
				StartOption: "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when end votes out of range",
			votes: i.Vote{
				StartVotes:  500,
				StartOption: "lte",
				EndVotes:    2147483648, // invalid votes
				EndOption:   "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid end option",
			votes: i.Vote{
				StartVotes:  500,
				StartOption: "lte",
				EndVotes:    1000,
				EndOption:   "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when same options",
			votes: i.Vote{
				StartVotes:  500,
				StartOption: "lte",
				EndVotes:    1000,
				EndOption:   "lte",
			},
			wantErr: true,
		},

		{
			desc: "return filter error when same votes",
			votes: i.Vote{
				StartVotes:  1000,
				StartOption: "lte",
				EndVotes:    1000,
				EndOption:   "gte",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &i.QueryParams{Vote: &tc.votes}
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
