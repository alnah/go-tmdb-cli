package query_test

import (
	"fmt"
	"testing"

	q "github.com/alnah/go-tmdb-cli/internal/query"

	"github.com/stretchr/testify/assert"
)

func TestUnitBuildQuery(t *testing.T) {
	t.Run("return query", func(t *testing.T) {
		testCases := []struct {
			desc  string
			query *q.QueryParams
			want  string
		}{
			{
				desc: "page year date average",
				query: &q.QueryParams{
					Year:    2000,
					Date:    &q.Date{StartDate: "2000-06-01", StartOption: "gte"},
					Average: &q.Average{StartAverage: 8.0, StartOption: "gte"},
					Vote:    &q.Vote{StartVote: 1000, StartOption: "gte"},
				},
				want: q.BaseURL + q.DiscoverURL +
					"primary_release_year=2000" +
					"&primary_release_date.gte=2000-06-01" +
					"&vote_average.gte=8.0&vote_count.gte=1000",
			},
			{
				desc: "popular movies list",
				query: &q.QueryParams{
					MovieListPath: "popular",
				},
				want: fmt.Sprintf(q.BaseURL+q.MoviesListURL, "popular"),
			},
		}

		for _, tc := range testCases {
			t.Run(tc.desc, func(t *testing.T) {
				got, err := tc.query.BuildQuery()
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("return filter error when error occured", func(t *testing.T) {
		testCases := []struct {
			desc  string
			query *q.QueryParams
		}{
			{
				desc: "on discover",
				query: &q.QueryParams{
					Date: &q.Date{
						StartDate:   "1800-01-01", // movies don't exist at that time
						StartOption: "gte",
					},
				},
			},
			{
				desc: "on movies list",
				query: &q.QueryParams{
					MovieListPath: "invalid",
				},
			},
		}
		for _, tc := range testCases {
			t.Run(tc.desc, func(t *testing.T) {
				var want *q.FilterError
				got, err := tc.query.BuildQuery()
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			})
		}
	})
}

func TestUnitFilterError(t *testing.T) {
	t.Run("return filter error", func(t *testing.T) {
		err := &q.FilterError{Filter: "Test", Message: "description"}
		assert.Equal(t, "Test description.", err.Error())
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
				query := q.QueryParams{MovieListPath: tc.path}
				var want *q.FilterError
				got, err := query.SetMoviesList()
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				query := q.QueryParams{MovieListPath: tc.path}
				want := fmt.Sprintf(q.MoviesListURL, tc.path)
				got, err := query.SetMoviesList()
				assert.NoError(t, err)
				assert.Equal(t, want, got)
			}
		})
	}

	t.Run("return empty string when no movies list", func(t *testing.T) {
		query := &q.QueryParams{MovieListPath: ""}
		got, err := query.SetMoviesList()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func TestUnitSetLanguage(t *testing.T) {
	testCases := []struct {
		desc    string
		iso     string
		want    string
		wantErr bool
	}{
		{
			desc: "return english",
			iso:  "en",
			want: "with_original_language=en",
		},
		{
			desc: "return french",
			iso:  "fr",
			want: "with_original_language=fr",
		},
		{
			desc: "return portuguese",
			iso:  "pt",
			want: "with_original_language=pt",
		},
		{
			desc: "return japanese",
			iso:  "jp",
			want: "with_original_language=jp",
		},
		{
			desc: "return mandarin",
			iso:  "zh",
			want: "with_original_language=zh",
		},
		{
			desc: "return guarani",
			iso:  "gn",
			want: "with_original_language=gn",
		},
		{
			desc:    "return filter error when not ISO 639-1 language code",
			iso:     "invalid",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &q.QueryParams{Language: tc.iso}
			got, err := query.SetLanguage()
			switch {
			case tc.wantErr:
				var want *q.FilterError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

	t.Run("return empty string when no language", func(t *testing.T) {
		query := &q.QueryParams{Language: ""}
		got, err := query.SetLanguage()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
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
			query := &q.QueryParams{Year: tc.value}
			got, err := query.SetYearFilter()
			switch {
			case tc.wantErr:
				var want *q.FilterError
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
		date    q.Date
		want    string
		wantErr bool
	}{
		{
			desc: "return start date gte",
			date: q.Date{
				StartDate:   "2001-01-01",
				StartOption: "gte",
			},
			want: "primary_release_date.gte=2001-01-01",
		},
		{
			desc: "return start date lte",
			date: q.Date{
				StartDate:   "2001-01-01",
				StartOption: "lte",
			},
			want: "primary_release_date.lte=2001-01-01",
		},
		{
			desc: "return start date gte and end date lte",
			date: q.Date{
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
			date: q.Date{
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
			date: q.Date{
				StartDate:   "01-01-2001", // invalid date format
				StartOption: "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid start date range",
			date: q.Date{
				StartDate:   "1800-01-01", // movies don't exist at that time
				StartOption: "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid start date option",
			date: q.Date{
				StartDate:   "2001-01-01",
				StartOption: "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid end date format",
			date: q.Date{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "01-01-2001", // invalid date format
				EndOption:   "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid end date range",
			date: q.Date{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "1800-01-01", // movies don't exist at that time
				EndOption:   "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid end date option",
			date: q.Date{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "2002-02-02",
				EndOption:   "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when same options",
			date: q.Date{
				StartDate:   "2001-01-01",
				StartOption: "lte",
				EndDate:     "2002-02-02",
				EndOption:   "lte",
			},
			wantErr: true,
		},

		{
			desc: "return filter error when same dates",
			date: q.Date{
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
			query := &q.QueryParams{Date: &tc.date}
			got, err := query.SetDateFilter()
			switch {
			case tc.wantErr:
				var want *q.FilterError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

	t.Run("return empty string when no date", func(t *testing.T) {
		query := &q.QueryParams{Date: nil}
		got, err := query.SetDateFilter()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func TestUnitSetAverageFilter(t *testing.T) {
	testCases := []struct {
		desc    string
		average q.Average
		want    string
		wantErr bool
	}{
		{
			desc: "return start average gte",
			average: q.Average{
				StartAverage: 8.0,
				StartOption:  "gte",
			},
			want: "vote_average.gte=8.0",
		},
		{
			desc: "return start average lte",
			average: q.Average{
				StartAverage: 8.0,
				StartOption:  "lte",
			},
			want: "vote_average.lte=8.0",
		},
		{
			desc: "return start average gte and end average lte",
			average: q.Average{
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
			average: q.Average{
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
			average: q.Average{
				StartAverage: 8.0,
				StartOption:  "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when start average out of range",
			average: q.Average{
				StartAverage: 100.0,
				StartOption:  "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when end average out of range",
			average: q.Average{
				StartAverage: 7.0,
				StartOption:  "lte",
				EndAverage:   100.0, // invalid average
				EndOption:    "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid end option",
			average: q.Average{
				StartAverage: 7.0,
				StartOption:  "lte",
				EndAverage:   8.0,
				EndOption:    "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when same options",
			average: q.Average{
				StartAverage: 7.0,
				StartOption:  "lte",
				EndAverage:   8.0,
				EndOption:    "lte",
			},
			wantErr: true,
		},

		{
			desc: "return filter error when same averages",
			average: q.Average{
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
			query := &q.QueryParams{Average: &tc.average}
			got, err := query.SetAverageFilter()
			switch {
			case tc.wantErr:
				var want *q.FilterError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

	t.Run("return empty string when no average", func(t *testing.T) {
		query := &q.QueryParams{Average: nil}
		got, err := query.SetAverageFilter()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func TestUnitSetVoteFilter(t *testing.T) {
	testCases := []struct {
		desc    string
		vote    q.Vote
		want    string
		wantErr bool
	}{
		{
			desc: "return start vote gte",
			vote: q.Vote{
				StartVote:   1000,
				StartOption: "gte",
			},
			want: "vote_count.gte=1000",
		},
		{
			desc: "return start vote lte",
			vote: q.Vote{
				StartVote:   1000,
				StartOption: "lte",
			},
			want: "vote_count.lte=1000",
		},
		{
			desc: "return start vote gte and end vote lte",
			vote: q.Vote{
				StartVote:   500,
				StartOption: "gte",
				EndVote:     1000,
				EndOption:   "lte",
			},
			want: "vote_count.gte=500&" +
				"vote_count.lte=1000",
		},
		{
			desc: "return start vote lte and end vote gte",
			vote: q.Vote{
				StartVote:   500,
				StartOption: "lte",
				EndVote:     1000,
				EndOption:   "gte",
			},
			want: "vote_count.lte=500&" +
				"vote_count.gte=1000",
		},
		{
			desc: "return filter error when invalid start option",
			vote: q.Vote{
				StartVote:   1000,
				StartOption: "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when start vote out of range",
			vote: q.Vote{
				StartVote:   2147483648, // invalid vote
				StartOption: "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when end vote out of range",
			vote: q.Vote{
				StartVote:   500,
				StartOption: "lte",
				EndVote:     2147483648, // invalid vote
				EndOption:   "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid end option",
			vote: q.Vote{
				StartVote:   500,
				StartOption: "lte",
				EndVote:     1000,
				EndOption:   "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when same options",
			vote: q.Vote{
				StartVote:   500,
				StartOption: "lte",
				EndVote:     1000,
				EndOption:   "lte",
			},
			wantErr: true,
		},

		{
			desc: "return filter error when same vote",
			vote: q.Vote{
				StartVote:   1000,
				StartOption: "lte",
				EndVote:     1000,
				EndOption:   "gte",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &q.QueryParams{Vote: &tc.vote}
			got, err := query.SetVoteFilter()
			switch {
			case tc.wantErr:
				var want *q.FilterError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

	t.Run("return empty string when no date", func(t *testing.T) {
		query := &q.QueryParams{Average: nil}
		got, err := query.SetVoteFilter()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func TestUnitSetGenresFilter(t *testing.T) {
	testCases := []struct {
		desc    string
		genres  q.Genres
		want    string
		wantErr bool
	}{
		{
			desc:   "return one genre ID",
			genres: q.Genres{"horror"},
			want:   "with_genres=27",
		},
		{
			desc:   "return many genre IDs",
			genres: q.Genres{"horror", "comedy"},
			want:   "with_genres=27,35",
		},
		{
			desc:    "return filter error when not allowed genre",
			genres:  q.Genres{"invalid"},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := q.QueryParams{Genres: &tc.genres}
			got, err := query.SetGenresFilter()
			switch {
			case tc.wantErr:
				var want *q.FilterError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
