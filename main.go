package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/exp/constraints"
)

func main() {
	token, err := getTMDBToken(".env")
	if err != nil {
		fmt.Println(err.Error())
	}

	intPtr := func(value int) *int {
		return &value
	}

	query := &queryParams{
		page:    intPtr(1),
		year:    intPtr(2000),
		date:    &date{value: "2000-06-01", option: "gte"},
		average: &average{value: 8.0, option: "gte"},
		vote:    &vote{value: 1000, option: "gte"},
	}

	url, err := query.buildURL()
	if err != nil {
		fmt.Println(err.Error())
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
}

func getTMDBToken(filepath string) (string, error) {
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

const (
	discoverURL = "https://api.themoviedb.org/3/discover/movie"
	language    = "?language=en-US&"
	firstYear   = 1888 // first movie of all time "Roundhay Garden Scene" (1888)
	minInt      = -2147483648
	maxInt      = 2147483647
	minAverage  = 0.0
	maxAverage  = 10.0
)

type date struct {
	value  string
	option string
}

type average struct {
	value  float64
	option string
}

type vote struct {
	value  int
	option string
}

type queryParams struct {
	page    *int
	year    *int
	date    *date
	average *average
	vote    *vote
}

func (qp queryParams) buildURL() (string, error) {
	var query strings.Builder
	query.WriteString(discoverURL + language)

	filters := []struct {
		setFilterFunc func() (string, error)
	}{
		{qp.setPageFilter},
		{qp.setYearFilter},
		{qp.setDateFilter},
		{qp.setAverageFilter},
		{qp.setVoteFilter},
	}

	for _, filter := range filters {
		if str, err := filter.setFilterFunc(); err != nil {
			return "", err
		} else {
			query.WriteString(str + "&")
		}
	}

	return strings.TrimSuffix(query.String(), "&"), nil
}

type filterError struct {
	filter  string
	message string
}

func (e *filterError) Error() string {
	return fmt.Sprintf("%s: %s.", e.filter, e.message)
}

func (qp queryParams) setPageFilter() (string, error) {
	if qp.page == nil {
		return "", nil
	}

	page := *qp.page
	if page < minInt || page > maxInt {
		return "", &filterError{
			filter:  "Page",
			message: fmt.Sprintf("Must be between %d and %d", minInt, maxInt),
		}
	}

	return fmt.Sprintf("page=%d", page), nil
}

func (qp queryParams) setYearFilter() (string, error) {
	if qp.year == nil {
		return "", nil
	}

	year := *qp.year
	nowYear := time.Now().Year()
	if err := validateRange(year, firstYear, nowYear, "Year"); err != nil {
		return "", err
	}

	return fmt.Sprintf("primary_release_year=%d", year), nil
}

func (qp queryParams) setDateFilter() (string, error) {
	if qp.date == nil {
		return "", nil
	}

	date := *qp.date
	parsedDate, err := time.Parse(time.DateOnly, date.value)
	if err != nil {
		return "", &filterError{
			filter:  "Date",
			message: `Date value must be a valid "YYYY-MM-DD" format`,
		}
	}

	year := parsedDate.Year()
	nowYear := time.Now().Year()
	if err = validateRange(year, firstYear, nowYear, "Year Date"); err != nil {
		return "", err
	}

	return handleGteOrLte(date, "primary_release_date", "Date")
}

func (qp queryParams) setAverageFilter() (string, error) {
	if qp.average == nil {
		return "", nil
	}

	average := *qp.average
	err := validateRange(average.value, minAverage, maxAverage, "Average")
	if err != nil {
		return "", err
	}

	return handleGteOrLte(average, "vote_average", "Average")
}

func (qp queryParams) setVoteFilter() (string, error) {
	if qp.vote == nil {
		return "", nil
	}

	vote := *qp.vote
	if err := validateRange(vote.value, minInt, maxInt, "Vote"); err != nil {
		return "", err
	}

	return handleGteOrLte(vote, "vote_count", "Vote")
}

func validateRange[T constraints.Ordered](val, min, max T, filter string) error {
	if val < min || val > max {
		message := "Must be between %v and %v"
		return &filterError{
			filter:  filter,
			message: fmt.Sprintf(message, min, max),
		}
	}
	return nil
}

func handleGteOrLte[T any](structure T, param, filter string) (string, error) {
	var value string
	var option string

	switch v := any(structure).(type) {
	case date:
		value = v.value
		option = v.option
	case average:
		value = fmt.Sprintf("%.1f", v.value)
		option = v.option
	case vote:
		value = fmt.Sprintf("%d", v.value)
		option = v.option
	default:
		return "", &filterError{
			filter:  filter,
			message: "Unsupported type",
		}
	}

	switch option {
	case "gte":
		return fmt.Sprintf("%s.gte=%s", param, value), nil
	case "lte":
		return fmt.Sprintf("%s.lte=%s", param, value), nil
	default:
		return "", &filterError{
			filter:  filter,
			message: `Option must be "gte" or "lte"`,
		}
	}
}
