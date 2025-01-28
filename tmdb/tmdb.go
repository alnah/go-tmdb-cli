package tmdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/exp/constraints"
)

const (
	APIBaseURL        string  = "https://api.themoviedb.org/3/"
	MoviesListPathURL string  = "movie/%s?language=en-US&page=1"
	DiscoverPathURL   string  = "discover/movie?language=en-US&"
	FirstFilmYear     int     = 1888 // "Roundhay Garden Scene", first known film
	MinInt32          int     = -2147483648
	MaxInt32          int     = 2147483647
	MinAverage        float64 = 0.0
	MaxAverage        float64 = 10.0
)

var GenreIdMap = map[string]int{
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

type QueryBuilder interface {
	BuildQuery() (string, error)
}

type ValidationError struct {
	Filter  string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s %s.", e.Filter, e.Message)
}

type SortOrderError struct{}

func (e SortOrderError) Error() string {
	return `Sorting order must be "asc" or "desc".`
}

type (
	Genres []string

	Dates struct {
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

	Votes struct {
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
		Dates         *Dates
		Average       *Average
		Votes         *Votes
		Genres        *Genres
	}
)

func (qp QueryParams) BuildQuery() (string, error) {
	var query strings.Builder
	query.WriteString(APIBaseURL)

	if qp.MovieListPath != "" {
		filter, err := qp.WithMoviesList()
		if err != nil {
			return "", err
		}
		query.WriteString(filter)
		return query.String(), nil
	}

	query.WriteString(DiscoverPathURL)
	filters := []struct {
		setFilterFunc func() (string, error)
	}{
		{qp.WithMoviesList},
		{qp.WithLanguage},
		{qp.WithYear},
		{qp.WithDates},
		{qp.WithAverage},
		{qp.WithVotes},
		{qp.WithGenres},
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

func (qp QueryParams) WithMoviesList() (string, error) {
	if qp.MovieListPath == "" {
		return "", nil
	}

	wantPaths := []string{"now_playing", "popular", "top_rated", "upcoming"}
	for _, p := range wantPaths {
		if qp.MovieListPath == p {
			return fmt.Sprintf(MoviesListPathURL, qp.MovieListPath), nil
		}
	}

	allowedPaths := allowedValues(wantPaths)
	allowedPaths = strings.ReplaceAll(allowedPaths, "_", " ")
	return "", &ValidationError{
		Filter:  "Movies List",
		Message: fmt.Sprintf("must be %s", allowedPaths),
	}
}

func (qp *QueryParams) WithLanguage() (string, error) {
	if qp.Language == "" {
		return "", nil
	}

	if len(qp.Language) != 2 {
		return "", &ValidationError{
			Filter:  "Language",
			Message: "must be a valid ISO 639-1 language code",
		}
	}

	return fmt.Sprintf("with_original_language=%s", qp.Language), nil
}

func (qp QueryParams) WithYear() (string, error) {
	if qp.Year == 0 {
		return "", nil
	}

	nowYear := time.Now().Year()
	if err := validateRange(qp.Year, FirstFilmYear, nowYear, "Year"); err != nil {
		return "", err
	}

	return fmt.Sprintf("primary_release_year=%d", qp.Year), nil
}

func (qp QueryParams) WithDates() (string, error) {
	if qp.Dates == nil {
		return "", nil
	}

	dates := *qp.Dates
	for _, d := range []string{dates.StartDate, dates.EndDate} {
		if d != "" {
			parsedDate, err := time.Parse(time.DateOnly, d)
			if err != nil {
				return "", &ValidationError{
					Filter:  "Date",
					Message: `value must be a valid "YYYY-MM-DD" format`,
				}
			}

			year := parsedDate.Year()
			nowYear := time.Now().Year()
			err = validateRange(year, FirstFilmYear, nowYear, "Year Date")
			if err != nil {
				return "", err
			}
		}
	}

	return handleGteOrLte(dates, "primary_release_date", "Date")
}

func (qp QueryParams) WithAverage() (string, error) {
	if qp.Average == nil {
		return "", nil
	}

	average := *qp.Average
	for _, a := range []float64{average.StartAverage, average.EndAverage} {
		if a != 0 {
			err := validateRange(a, MinAverage, MaxAverage, "Average")
			if err != nil {
				return "", err
			}
		}
	}

	return handleGteOrLte(average, "vote_average", "Average")
}

func (qp QueryParams) WithVotes() (string, error) {
	if qp.Votes == nil {
		return "", nil
	}

	votes := *qp.Votes
	for _, v := range []int{votes.StartVote, votes.EndVote} {
		if v != 0 {
			if err := validateRange(v, MinInt32, MaxInt32, "Vote"); err != nil {
				return "", err
			}
		}
	}

	return handleGteOrLte(votes, "vote_count", "Vote")
}

func (qp QueryParams) WithGenres() (string, error) {
	if qp.Genres == nil {
		return "", nil
	}

	allGenres := getKeys(GenreIdMap)
	allowedGenres := allowedValues(allGenres)

	var formattedGenres Genres
	for _, g := range *qp.Genres {
		id, found := GenreIdMap[g]

		if !found {
			return "", &ValidationError{
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
		return &ValidationError{
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
	case Dates:
		startValue, endValue = v.StartDate, v.EndDate
		startOption, endOption = v.StartOption, v.EndOption
	case Average:
		startValue = fmt.Sprintf("%.1f", v.StartAverage)
		endValue = fmt.Sprintf("%.1f", v.EndAverage)
		startOption, endOption = v.StartOption, v.EndOption
	case Votes:
		startValue = fmt.Sprintf("%d", v.StartVote)
		endValue = fmt.Sprintf("%d", v.EndVote)
		startOption, endOption = v.StartOption, v.EndOption
	}

	if startOption == endOption {
		return "", &ValidationError{
			Filter:  filter,
			Message: `The same option can't be used twice`,
		}
	}

	if startValue == endValue {
		return "", &ValidationError{
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
			return "", &ValidationError{
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

type (
	Movie struct {
		Id            int     `json:"id"`
		Date          string  `json:"release_date"`
		OriginalTitle string  `json:"original_title"`
		Title         string  `json:"title"`
		Average       float64 `json:"vote_average"`
		Votes         int     `json:"vote_count"`
	}

	MovieCollection []Movie

	TMDBResponse struct {
		Results     MovieCollection `json:"results"`
		Page        int             `json:"page"`
		TotalPages  int             `json:"total_pages"`
		TotalMovies int             `json:"total_results"`
	}
)

func (m MovieCollection) DefaultSort() (MovieCollection, error) {
	sortOrder := []struct {
		sortFunc func(order string) (MovieCollection, error)
		order    string
	}{
		{m.SortByOriginalTitle, "asc"},
		{m.SortByTitle, "asc"},
		{m.SortByDate, "asc"},
		{m.SortByVotes, "desc"},
		{m.SortByAverage, "desc"},
	}

	var err error
	for _, s := range sortOrder {
		m, err = s.sortFunc(s.order)
		if err != nil {
			return m, err
		}
	}

	if err != nil {
		return m, err
	}

	return m, nil
}

func (m MovieCollection) SortByDate(order string) (MovieCollection, error) {
	return m.sort(order, func(i int, j int) bool {
		iDate, _ := time.Parse(time.DateOnly, m[i].Date)
		jDate, _ := time.Parse(time.DateOnly, m[j].Date)

		return iDate.Before(jDate)
	})
}

func (m MovieCollection) SortByOriginalTitle(order string) (MovieCollection, error) {
	return m.sort(order, func(i int, j int) bool {
		return m[i].OriginalTitle < m[j].OriginalTitle
	})
}

func (m MovieCollection) SortByTitle(order string) (MovieCollection, error) {
	return m.sort(order, func(i int, j int) bool {
		return m[i].Title < m[j].Title
	})
}

func (m MovieCollection) SortByAverage(order string) (MovieCollection, error) {
	return m.sort(order, func(i int, j int) bool {
		return m[i].Average < m[j].Average
	})
}

func (m MovieCollection) SortByVotes(order string) (MovieCollection, error) {
	return m.sort(order, func(i int, j int) bool {
		return m[i].Votes < m[j].Votes
	})
}

func (m MovieCollection) sort(
	order string,
	compare func(i, j int) bool,
) (MovieCollection, error) {
	if err := checkOrder(order); err != nil {
		return m, err
	}

	sort.Slice(m, func(i, j int) bool {
		if order == "asc" {
			return compare(i, j)
		}

		return !compare(i, j)
	})

	return m, nil
}

func checkOrder(order string) error {
	if order != "asc" && order != "desc" {
		return &SortOrderError{}
	}

	return nil
}

func FetchMovieCollection(
	query QueryBuilder,
	token string,
	maxItems int,
) (MovieCollection, error) {
	var url string
	url, err := query.BuildQuery()
	if err != nil {
		return MovieCollection{}, err
	}

	req, err := newRequest("GET", url, token)
	if err != nil {
		return MovieCollection{}, err
	}

	var data TMDBResponse
	var results MovieCollection
	for {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return MovieCollection{}, err
		}
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		if err = json.Unmarshal(body, &data); err != nil {
			return MovieCollection{}, err
		}

		results = append(results, data.Results...)
		if len(results) > maxItems || data.Page >= data.TotalPages {
			break
		}

		data.Page = data.Page + 1
		url, err = query.BuildQuery()
		if err != nil {
			return MovieCollection{}, err
		}

		req, err = newRequest("GET", url, token)
		if err != nil {
			return MovieCollection{}, err
		}
	}

	if len(results) > maxItems {
		results = results[:maxItems]
	}

	return results, nil
}

func newRequest(method, url, token string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	return req, nil
}

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
