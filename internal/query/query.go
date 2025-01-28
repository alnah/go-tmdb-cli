package query

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

var genreMap = map[string]int{
	"action":          28,
	"adventure":       12,
	"animation":       16,
	"comedy":          35,
	"crime":           80,
	"documentary":     99,
	"drama":           18,
	"family":          10751,
	"fantasy":         14,
	"history":         36,
	"horror":          27,
	"music":           10402,
	"mystery":         9648,
	"romance":         10749,
	"science-fiction": 878,
	"tv-movie":        10770,
	"thriller":        53,
	"war":             10752,
	"western":         37,
}

type (
	Query interface {
		BuildQuery() (string, error)
	}

	Genres []string

	Date struct {
		StartDate   string
		StartOption string
		EndDate     string
		EndOption   string
	}

	Average struct {
		StartAverage float64
		StartOption  string
		EndAverage   float64
		EndOption    string
	}

	Vote struct {
		StartVote   int
		StartOption string
		EndVote     int
		EndOption   string
	}

	QueryParams struct {
		MovieListPath string
		Language      string
		MaxItems      int
		Year          int
		Date          *Date
		Average       *Average
		Vote          *Vote
		Genres        *Genres
	}
)

func (qp QueryParams) BuildQuery() (string, error) {
	var query strings.Builder
	query.WriteString(BaseURL)

	if qp.MovieListPath != "" {
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
		{qp.SetLanguage},
		{qp.SetYearFilter},
		{qp.SetDateFilter},
		{qp.SetAverageFilter},
		{qp.SetVoteFilter},
		{qp.SetGenresFilter},
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
	return fmt.Sprintf("%s %s.", e.Filter, e.Message)
}

func (qp QueryParams) SetMoviesList() (string, error) {
	if qp.MovieListPath == "" {
		return "", nil
	}

	wantPaths := []string{"now_playing", "popular", "top_rated", "upcoming"}
	for _, p := range wantPaths {
		if qp.MovieListPath == p {
			return fmt.Sprintf(MoviesListURL, qp.MovieListPath), nil
		}
	}

	allowedPaths := allowedValues(wantPaths)
	allowedPaths = strings.ReplaceAll(allowedPaths, "_", " ")
	return "", &FilterError{
		Filter:  "Movies List",
		Message: fmt.Sprintf("must be %s", allowedPaths),
	}
}

func (qp *QueryParams) SetLanguage() (string, error) {
	if qp.Language == "" {
		return "", nil
	}

	if len(qp.Language) != 2 {
		return "", &FilterError{
			Filter:  "Language",
			Message: "must be a valid ISO 639-1 language code",
		}
	}

	return fmt.Sprintf("with_original_language=%s", qp.Language), nil
}

func (qp QueryParams) SetYearFilter() (string, error) {
	if qp.Year == 0 {
		return "", nil
	}

	nowYear := time.Now().Year()
	if err := validateRange(qp.Year, firstYear, nowYear, "Year"); err != nil {
		return "", err
	}

	return fmt.Sprintf("primary_release_year=%d", qp.Year), nil
}

func (qp QueryParams) SetDateFilter() (string, error) {
	if qp.Date == nil {
		return "", nil
	}

	date := *qp.Date
	for _, d := range []string{date.StartDate, date.EndDate} {
		if d != "" {
			parsedDate, err := time.Parse(time.DateOnly, d)
			if err != nil {
				return "", &FilterError{
					Filter:  "Date",
					Message: `value must be a valid "YYYY-MM-DD" format`,
				}
			}

			year := parsedDate.Year()
			nowYear := time.Now().Year()
			err = validateRange(year, firstYear, nowYear, "Year Date")
			if err != nil {
				return "", err
			}
		}
	}

	return handleGteOrLte(date, "primary_release_date", "Date")
}

func (qp QueryParams) SetAverageFilter() (string, error) {
	if qp.Average == nil {
		return "", nil
	}

	average := *qp.Average
	for _, a := range []float64{average.StartAverage, average.EndAverage} {
		if a != 0 {
			err := validateRange(a, minAverage, maxAverage, "Average")
			if err != nil {
				return "", err
			}
		}
	}

	return handleGteOrLte(average, "vote_average", "Average")
}

func (qp QueryParams) SetVoteFilter() (string, error) {
	if qp.Vote == nil {
		return "", nil
	}

	vote := *qp.Vote
	for _, v := range []int{vote.StartVote, vote.EndVote} {
		if v != 0 {
			if err := validateRange(v, minInt, maxInt, "Vote"); err != nil {
				return "", err
			}
		}
	}

	return handleGteOrLte(vote, "vote_count", "Vote")
}

func (qp QueryParams) SetGenresFilter() (string, error) {
	if qp.Genres == nil {
		return "", nil
	}

	allGenres := getKeys(genreMap)
	allowedGenres := allowedValues(allGenres)

	var formattedGenres Genres
	for _, g := range *qp.Genres {
		id, found := genreMap[g]

		if !found {
			return "", &FilterError{
				Filter:  "Genres",
				Message: fmt.Sprintf("must be %s", allowedGenres),
			}
		}

		formattedGenres = append(formattedGenres, strconv.Itoa(id))
	}

	return fmt.Sprintf("with_genres=%s", strings.Join(formattedGenres, ",")), nil
}

func validateRange[T constraints.Ordered](val, min, max T, filter string) error {
	if val < min || val > max {
		message := "must be between %v and %v"
		return &FilterError{
			Filter:  filter,
			Message: fmt.Sprintf(message, min, max),
		}
	}
	return nil
}

func getKeys(hashmap map[string]int) []string {
	keys := make([]string, len(hashmap))
	i := 0
	for k := range hashmap {
		keys[i] = k
		i++
	}

	return keys
}

func allowedValues(values []string) string {
	var formatted string

	formatted = strings.Join(values, ", ")
	lastIndex := strings.LastIndex(formatted, ", ")
	if lastIndex != -1 {
		formatted = formatted[:lastIndex] + ", or " + formatted[lastIndex+2:]
	}

	return formatted
}

func handleGteOrLte[T any](structure T, param, filter string) (string, error) {
	var startValue, endValue, startPart string
	var startOption, endOption, endPart string

	switch v := any(structure).(type) {
	case Date:
		startValue, endValue = v.StartDate, v.EndDate
		startOption, endOption = v.StartOption, v.EndOption
	case Average:
		startValue = fmt.Sprintf("%.1f", v.StartAverage)
		endValue = fmt.Sprintf("%.1f", v.EndAverage)
		startOption, endOption = v.StartOption, v.EndOption
	case Vote:
		startValue = fmt.Sprintf("%d", v.StartVote)
		endValue = fmt.Sprintf("%d", v.EndVote)
		startOption, endOption = v.StartOption, v.EndOption
	}

	if startOption == endOption {
		return "", &FilterError{
			Filter:  filter,
			Message: `The same option can't be used twice`,
		}
	}

	if startValue == endValue {
		return "", &FilterError{
			Filter:  filter,
			Message: `The same value can't be used twice`,
		}
	}

	buildQuery := func(option, param, filter string, value any) (string, error) {
		var part string

		switch option {
		case "gte":
			part = fmt.Sprintf("%s.gte=%s", param, value)
		case "lte":
			part = fmt.Sprintf("%s.lte=%s", param, value)
		default:
			return "", &FilterError{
				Filter:  filter,
				Message: `Option must be "gte" or "lte"`,
			}
		}
		return part, nil
	}

	startPart, err := buildQuery(startOption, param, filter, startValue)
	if err != nil {
		return "", err
	}

	if endValue != "" && endOption != "" {
		endPart, err = buildQuery(endOption, param, filter, endValue)
		if err != nil {
			return "", err
		}
	}

	if endPart != "" {
		return startPart + "&" + endPart, nil
	}

	return startPart, nil
}
