package tmdb

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/exp/constraints"
)

const (
	BaseURL       string  = "https://api.themoviedb.org/3/"
	MoviesListURL string  = "movie/%s?language=en-US&page=1"
	DiscoverURL   string  = "discover/movie?language=en-US&"
	firstYear     int     = 1888 // "Roundhay Garden Scene" (1888), first movie
	minInt        int     = -2147483648
	maxInt        int     = 2147483647
	minAverage    float64 = 0.0
	maxAverage    float64 = 10.0
)

func GetTMDBToken(filepath string) (string, error) {
	byt, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	viper.SetConfigType("ENV")
	if err = viper.ReadConfig(bytes.NewBuffer(byt)); err != nil {
		return "", err
	}

	return viper.GetString("TOKEN"), nil
}

type Date struct {
	Value  string
	Option string
}

type Average struct {
	Value  float64
	Option string
}

type Vote struct {
	Value  int
	Option string
}

type QueryParams struct {
	MovieListPath *string
	Page          *int
	Year          *int
	Date          *Date
	Average       *Average
	Vote          *Vote
}

func (qp QueryParams) BuildURL() (string, error) {
	var query strings.Builder
	query.WriteString(BaseURL)

	if qp.MovieListPath != nil {
		filter, err := qp.SetMoviesList()
		if err != nil {
			return "", err
		}
		query.WriteString(filter)
		return query.String(), nil
	}

	query.WriteString(DiscoverURL)
	filters := []struct {
		setFilterFunc func() (string, error)
	}{
		{qp.SetMoviesList},
		{qp.SetPageFilter},
		{qp.SetYearFilter},
		{qp.SetDateFilter},
		{qp.SetAverageFilter},
		{qp.SetVoteFilter},
	}

	for _, filter := range filters {
		if str, err := filter.setFilterFunc(); err != nil {
			return "", err
		} else if str != "" {
			query.WriteString(str + "&")
		}
	}

	return strings.TrimSuffix(query.String(), "&"), nil
}

type FilterError struct {
	Filter  string
	Message string
}

func (e FilterError) Error() string {
	return fmt.Sprintf("%s: %s.", e.Filter, e.Message)
}

func (qp QueryParams) SetMoviesList() (string, error) {
	if qp.MovieListPath == nil {
		return "", nil
	}

	wantPaths := []string{"now_playing", "popular", "top_rated", "upcoming"}
	path := *qp.MovieListPath
	for _, p := range wantPaths {
		if path == p {
			return fmt.Sprintf(MoviesListURL, path), nil
		}
	}

	var paths string
	for _, p := range wantPaths {
		paths += `"` + p + `"` + ", or "
	}
	paths = strings.ReplaceAll(paths, "_", " ")
	paths = strings.TrimSuffix(paths, ", or ")
	return "", &FilterError{
		Filter:  "MoviesList",
		Message: fmt.Sprintf("Path must be %v", paths),
	}
}

func (qp QueryParams) SetPageFilter() (string, error) {
	if qp.Page == nil {
		return "", nil
	}

	page := *qp.Page
	if page < minInt || page > maxInt {
		return "", &FilterError{
			Filter:  "Page",
			Message: fmt.Sprintf("Must be between %d and %d", minInt, maxInt),
		}
	}

	return fmt.Sprintf("page=%d", page), nil
}

func (qp QueryParams) SetYearFilter() (string, error) {
	if qp.Year == nil {
		return "", nil
	}

	year := *qp.Year
	nowYear := time.Now().Year()
	if err := validateRange(year, firstYear, nowYear, "Year"); err != nil {
		return "", err
	}

	return fmt.Sprintf("primary_release_year=%d", year), nil
}

func (qp QueryParams) SetDateFilter() (string, error) {
	if qp.Date == nil {
		return "", nil
	}

	date := *qp.Date
	parsedDate, err := time.Parse(time.DateOnly, date.Value)
	if err != nil {
		return "", &FilterError{
			Filter:  "Date",
			Message: `Date value must be a valid "YYYY-MM-DD" format`,
		}
	}

	year := parsedDate.Year()
	nowYear := time.Now().Year()
	if err = validateRange(year, firstYear, nowYear, "Year Date"); err != nil {
		return "", err
	}

	return handleGteOrLte(date, "primary_release_date", "Date")
}

func (qp QueryParams) SetAverageFilter() (string, error) {
	if qp.Average == nil {
		return "", nil
	}

	average := *qp.Average
	err := validateRange(average.Value, minAverage, maxAverage, "Average")
	if err != nil {
		return "", err
	}

	return handleGteOrLte(average, "vote_average", "Average")
}

func (qp QueryParams) SetVoteFilter() (string, error) {
	if qp.Vote == nil {
		return "", nil
	}

	vote := *qp.Vote
	if err := validateRange(vote.Value, minInt, maxInt, "Vote"); err != nil {
		return "", err
	}

	return handleGteOrLte(vote, "vote_count", "Vote")
}

func validateRange[T constraints.Ordered](val, min, max T, filter string) error {
	if val < min || val > max {
		message := "Must be between %v and %v"
		return &FilterError{
			Filter:  filter,
			Message: fmt.Sprintf(message, min, max),
		}
	}
	return nil
}

func handleGteOrLte[T any](structure T, param, filter string) (string, error) {
	var value string
	var option string

	switch v := any(structure).(type) {
	case Date:
		value = v.Value
		option = v.Option
	case Average:
		value = fmt.Sprintf("%.1f", v.Value)
		option = v.Option
	case Vote:
		value = fmt.Sprintf("%d", v.Value)
		option = v.Option
	}

	switch option {
	case "gte":
		return fmt.Sprintf("%s.gte=%s", param, value), nil
	case "lte":
		return fmt.Sprintf("%s.lte=%s", param, value), nil
	default:
		return "", &FilterError{
			Filter:  filter,
			Message: `Option must be "gte" or "lte"`,
		}
	}
}
