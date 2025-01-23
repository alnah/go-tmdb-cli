package main

import (
	"io/fs"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUnitGetTMDBToken(t *testing.T) {
	t.Run("returns TMDB API Token from env file", func(t *testing.T) {
		file := createTempFile(t)

		_, err := file.WriteString(`TOKEN="test"`)
		assert.NoError(t, err)

		got, err := getTMDBToken(file.Name())
		assert.NoError(t, err)
		assert.Equal(t, "test", got)
	})

	t.Run("returns path error when env file doesn't exist", func(t *testing.T) {
		var want *fs.PathError
		_, err := getTMDBToken("test.env")
		assert.ErrorAs(t, err, &want)
	})

	t.Run("returns config parse error when invalid data", func(t *testing.T) {
		file := createTempFile(t)

		_, err := file.WriteString(`invalid_data`)
		assert.NoError(t, err)

		var want viper.ConfigParseError
		_, err = getTMDBToken(file.Name())
		assert.ErrorAs(t, err, &want)
	})
}

func TestUnitBuildQuery(t *testing.T) {
	t.Run("return query", func(t *testing.T) {
		query := &queryParams{
			page:    intPtr(2),
			year:    intPtr(2000),
			date:    &date{value: "2000-06-01", option: "gte"},
			average: &average{value: 8.0, option: "gte"},
			vote:    &vote{value: 1000, option: "gte"},
		}

		want := "https://api.themoviedb.org/3/discover/movie?language=en-US" +
			"&page=2" +
			"&primary_release_year=2000" +
			"&primary_release_date.gte=2000-06-01" +
			"&vote_average.gte=8.0&vote_count.gte=1000"
		got, err := query.buildQuery()
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("return filter error when error occured", func(t *testing.T) {
		query := &queryParams{
			date: &date{
				value:  "1800-01-01", // movies don't exist at that time
				option: "gte",
			},
		}

		var want *filterError
		got, err := query.buildQuery()
		assert.Empty(t, got)
		assert.ErrorAs(t, err, &want)
	})
}

func TestUnitFilterError(t *testing.T) {
	t.Run("return filter error", func(t *testing.T) {
		err := &filterError{filter: "Test", message: "Description"}
		assert.Equal(t, "Test: Description.", err.Error())
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
			value: intPtr(20),
			want:  "page=20",
		},
		{
			desc:    "return empty string if no page",
			value:   nil,
			wantErr: false,
		},
		{
			desc:    "return filter error if out of range",
			value:   intPtr(2147483648), // out of range
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &queryParams{page: tc.value}
			got, err := query.setPageFilter()
			switch {
			case tc.wantErr:
				var want *filterError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestHandleGteOrLTe(t *testing.T) {
	t.Run("return filter error for unsupported type", func(t *testing.T) {
		type fakeType struct {
			value  string
			option string
		}
		invalid := &fakeType{value: "", option: ""}

		var want *filterError
		got, err := handleGteOrLte(invalid, "Random", "invalid")
		assert.Empty(t, got)
		assert.ErrorAs(t, err, &want)
	})
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
			value: intPtr(2001),
			want:  "primary_release_year=2001",
		},
		{
			desc:    "return empty string when no year",
			value:   nil,
			wantErr: false,
		},
		{
			desc:    "return filter error when year out of range",
			value:   intPtr(2147483648), // out of range
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &queryParams{year: tc.value}
			got, err := query.setYearFilter()
			switch {
			case tc.wantErr:
				var want *filterError
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
		date    date
		want    string
		wantErr bool
	}{
		{
			desc: "return date gte",
			date: date{
				value:  "2001-01-01",
				option: "gte",
			},
			want: "primary_release_date.gte=2001-01-01",
		},
		{
			desc: "return date lte",
			date: date{
				value:  "2001-01-01",
				option: "lte",
			},
			want: "primary_release_date.lte=2001-01-01",
		},
		{
			desc: "return filter error when invalid date value format",
			date: date{
				value:  "01-01-2001", // invalid date format
				option: "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid date value range",
			date: date{
				value:  "1800-01-01", // movies don't exist at that time
				option: "gte",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when invalid date option",
			date: date{
				value:  "2001-01-01",
				option: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &queryParams{date: &tc.date}
			got, err := query.setDateFilter()
			switch {
			case tc.wantErr:
				var want *filterError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

	t.Run("return empty string when no date", func(t *testing.T) {
		query := &queryParams{date: nil}
		got, err := query.setDateFilter()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func TestUnitSetAverageFilter(t *testing.T) {
	testCases := []struct {
		desc    string
		average average
		want    string
		wantErr bool
	}{
		{
			desc: "return average gte",
			average: average{
				value:  8.0,
				option: "gte",
			},
			want: "vote_average.gte=8.0",
		},
		{
			desc: "return average lte",
			average: average{
				value:  8.0,
				option: "lte",
			},
			want: "vote_average.lte=8.0",
		},
		{
			desc: "return filter error when invalid option",
			average: average{
				value:  8.0,
				option: "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when value out of range",
			average: average{
				value:  100.0,
				option: "gte",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &queryParams{average: &tc.average}
			got, err := query.setAverageFilter()
			switch {
			case tc.wantErr:
				var want *filterError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

	t.Run("return empty string when no average", func(t *testing.T) {
		query := &queryParams{average: nil}
		got, err := query.setAverageFilter()
		assert.Empty(t, got)
		assert.NoError(t, err)
	})
}

func TestUnitSetVoteFilter(t *testing.T) {
	testCases := []struct {
		desc    string
		vote    vote
		want    string
		wantErr bool
	}{
		{
			desc: "return vote gte",
			vote: vote{
				value:  1000,
				option: "gte",
			},
			want: "vote_count.gte=1000",
		},
		{
			desc: "return vote lte",
			vote: vote{
				value:  1000,
				option: "lte",
			},
			want: "vote_count.lte=1000",
		},
		{
			desc: "return filter error when invalid option",
			vote: vote{
				value:  1000,
				option: "invalid",
			},
			wantErr: true,
		},
		{
			desc: "return filter error when value out of range",
			vote: vote{
				value:  2147483648, // out of range,
				option: "gte",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			query := &queryParams{vote: &tc.vote}
			got, err := query.setVoteFilter()
			switch {
			case tc.wantErr:
				var want *filterError
				assert.Empty(t, got)
				assert.ErrorAs(t, err, &want)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}

	t.Run("return empty string when no date", func(t *testing.T) {
		query := &queryParams{average: nil}
		got, err := query.setVoteFilter()
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

func intPtr(value int) *int {
	return &value
}
