package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v5"
)

const (
	earliestMovie  = 1888
	minVoteAverage = 0
	maxVoteAverage = 10
	minVoteCount   = 0
	yearFormat     = "2006"
	helpISO6391    = "https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes"
	firstPage      = 1
	resultsPerPage = 20
	maxAPICalls    = 20
	APIMaxItems    = resultsPerPage * maxAPICalls
)

var (
	yearNow   = time.Now().Year()
	genresMap = map[string]int{
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
)

type (
	// movies represents a collection of TMDB film entries for processing.
	movies []movie
	// movie contains essential metadata for a single TMDB film record.
	movie struct {
		ID            int     `json:"id"`
		OriginalTitle string  `json:"original_title"`
		ReleaseDate   string  `json:"release_date"`
		Title         string  `json:"title"`
		VoteAverage   float64 `json:"vote_average"`
		VoteCount     int     `json:"vote_count"`
	}
)

// deduplicate removes repeated movie entries while preserving order.
func (m movies) deduplicate() movies {
	seen := make(map[int]bool)
	result := make(movies, 0, len(m))
	for _, movie := range m {
		if !seen[movie.ID] {
			seen[movie.ID] = true
			result = append(result, movie)
		}
	}
	return result
}

// sortByField organizes movies by specified criteria and direction.
func (m movies) sortByField(param string) (movies, error) {
	param = cleanString(param)
	parts := strings.Split(param, ",")
	if len(parts) != 2 {
		return m, fmt.Errorf(`sort format: expected "field, order", e.g. "average,desc" or "date,asc"`)
	}
	compareFunc, err := m.getCompareFunc(parts[0])
	if err != nil {
		return m, err
	}
	if err := m.sortHelper(parts[1], compareFunc); err != nil {
		return m, err
	}
	return m, nil
}

func (m movies) compareReleaseDate(i, j int) bool {
	iDate, _ := time.Parse(time.DateOnly, m[i].ReleaseDate)
	jDate, _ := time.Parse(time.DateOnly, m[j].ReleaseDate)
	return iDate.Before(jDate)
}

func (m movies) compareOriginalTitle(i, j int) bool { return m[i].OriginalTitle < m[j].OriginalTitle }
func (m movies) compareTitle(i, j int) bool         { return m[i].Title < m[j].Title }
func (m movies) compareVoteAverage(i, j int) bool   { return m[i].VoteAverage < m[j].VoteAverage }
func (m movies) compareVoteCount(i, j int) bool     { return m[i].VoteCount < m[j].VoteCount }

func (m movies) getCompareFunc(field string) (func(i, j int) bool, error) {
	mapCompareFunc := map[string]func(i, j int) bool{
		"date":    m.compareReleaseDate,
		"otitle":  m.compareOriginalTitle,
		"title":   m.compareTitle,
		"average": m.compareVoteAverage,
		"votes":   m.compareVoteCount,
	}
	compareFunc, ok := mapCompareFunc[field]
	if !ok {
		return nil, fmt.Errorf("validation error: movie list parameter must be one of: %v",
			[]string{"date", "otitle", "title", "average", "votes"})
	}
	return compareFunc, nil
}

func (m movies) sortHelper(order string, compare func(i, j int) bool) error {
	if err := validateOrder(order); err != nil {
		return err
	}
	sort.Slice(m, func(i, j int) bool {
		if order == "asc" {
			return compare(i, j)
		}
		return !compare(i, j)
	})
	return nil
}

func validateOrder(order string) error {
	if order != "asc" && order != "desc" {
		return fmt.Errorf("validation error: order parameter must be one of: %v",
			[]string{"asc", "desc"})
	}
	return nil
}

type (
	// httpClient manages authenticated requests and error handling for GitHub API.
	httpClient struct {
		url    string
		APIKey string
		Method string
		Client *http.Client
	}
	// tmdbResponse represents paginated results from TMDB's API endpoints.
	tmdbResponse struct {
		Page         int    `json:"page"`
		Results      movies `json:"results"`
		TotalPages   int    `json:"total_pages"`
		TotalResults int    `json:"total_results"`
	}
)

// newHTTPClient configures secure defaults for TMDB API communication.
func newHTTPClient(apiKey string) *httpClient {
	return &httpClient{
		APIKey: apiKey,
		Method: "GET",
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// asyncFetchMovies efficiently retrieves multiple pages of movie results.
func asyncFetchMovies(hc *httpClient, url string, maxItems int) (movies, error) {
	if maxItems > APIMaxItems {
		return movies{}, fmt.Errorf("validation error: movies can't be more than %d", APIMaxItems)
	}
	var (
		allResults movies
		mu         sync.Mutex
		wg         sync.WaitGroup
	)
	firstPageURL := fmt.Sprintf("%s&page=%d", url, firstPage)
	firstRes, err := fetchTMDBResponse(hc, firstPageURL)
	if err != nil {
		return movies{}, err
	}
	if maxItems < len(firstRes.Results) {
		firstRes.Results = firstRes.Results[:maxItems]
		return firstRes.Results, nil
	}
	totalPages := (maxItems + resultsPerPage - firstPage) / resultsPerPage
	errChan := make(chan error, totalPages-firstPage)
	for page := 2; page <= totalPages; page++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			fetchUrl := fmt.Sprintf("%s&page=%d", url, p)
			pageRes, err := fetchTMDBResponse(hc, fetchUrl)
			if err != nil {
				errChan <- err
				return
			}
			mu.Lock()
			allResults = append(allResults, pageRes.Results...)
			mu.Unlock()
		}(page)
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return movies{}, err
		}
	}
	allResults = append(firstRes.Results, allResults...)
	if len(allResults) > maxItems {
		allResults = allResults[:maxItems]
	}
	return allResults.deduplicate(), nil
}

func (hc *httpClient) setURL(url string) {
	hc.url = url
}

// fetchTMDBResponse gets a single page of results from TMDB API.
func fetchTMDBResponse(hc *httpClient, url string) (tmdbResponse, error) {
	hc.setURL(url)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tmdbRes, err := hc.do(ctx)
	if err != nil {
		return tmdbResponse{}, err
	}
	return tmdbRes, nil
}

// do retrieves movie data from TMDB with a retry mechanism based on exponential backoff.
func (hc *httpClient) do(ctx context.Context) (tmdbResponse, error) {
	op := func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, hc.Method, hc.url, nil)
		if err != nil {
			return nil, backoff.Permanent(fmt.Errorf("request error: %w", err))
		}
		req.Header.Add("Authorization", "Bearer "+hc.APIKey)
		req.Header.Add("Content-Type", "application/json")
		cli := newHTTPClient(hc.APIKey)
		res, err := cli.Client.Do(req)
		if err != nil {
			return nil, backoff.Permanent(fmt.Errorf("request error: %w", err))
		}
		switch {
		case res.StatusCode >= 500:
			return nil, backoff.Permanent(fmt.Errorf("TMDB API server error: %q", res.Status))
		case res.StatusCode == 429:
			sec, err := strconv.ParseInt(res.Header.Get("Retry-After"), 10, 64)
			if err == nil {
				return nil, backoff.RetryAfter(int(sec))
			}
		case res.StatusCode >= 400:
			return nil, backoff.Permanent(fmt.Errorf("TMDB API client error: %q", res.Status))
		}
		return res, nil
	}
	res, err := backoff.Retry(ctx, op, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		return tmdbResponse{}, fmt.Errorf("fetch TMDB response: %w", err)
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()
	var results tmdbResponse
	if err = json.NewDecoder(res.Body).Decode(&results); err != nil {
		return tmdbResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return results, nil
}

type (
	// urlBuilder constructs valid TMDB API URLs with proper parameter encoding.
	urlBuilder struct {
		BaseURL      string
		ListPath     string
		DiscoverPath string
	}
	// queryParams encapsulates filter criteria for discover movie searches.
	queryParams struct {
		MaxItems      int
		Language      string
		Year          string
		VoteAverage   string
		VoteCount     string
		WithGenres    string
		WithoutGenres string
	}
)

// newURLBuilder initializes URL patterns for TMDB API endpoints.
func newURLBuilder() *urlBuilder {
	return &urlBuilder{
		BaseURL:      "https://api.themoviedb.org/3",
		ListPath:     "/movie/%s?",
		DiscoverPath: "/discover/movie?",
	}
}

// list generates URLs for TMDB's predefined movie list endpoints.
func (u *urlBuilder) list(param string) (string, error) {
	if param != "now_playing" && param != "popular" && param != "top_rated" && param != "upcoming" {
		return "", fmt.Errorf("validation error: movie list parameter must be one of: %v",
			[]string{"now_playing", "popular", "top_rated", "upcoming"})
	}
	return fmt.Sprintf(u.BaseURL+u.ListPath, param), nil
}

// discover builds complex query URLs for filtered movie searches.
func (ub *urlBuilder) discover(q queryParams) (string, error) {
	var query string
	var err error
	url := ub.BaseURL + ub.DiscoverPath
	for _, handler := range []struct {
		condition bool
		handle    func() (string, error)
	}{
		{q.Language != "", q.handleLanguage},
		{q.Year != "", q.handleYear},
		{q.VoteAverage != "", q.handleVoteAverage},
		{q.VoteCount != "", q.handleVoteCount},
		{q.WithGenres != "", q.handleWithGenres},
		{q.WithoutGenres != "", q.handleWithoutGenres},
	} {
		if handler.condition {
			if query, err = handler.handle(); err != nil {
				return "", err
			}
			url += query
		}
	}
	return strings.TrimSuffix(url, "&"), nil
}

func (qp *queryParams) handleLanguage() (string, error) {
	iso639_1Length := 2
	qp.Language = cleanString(qp.Language)
	if len(qp.Language) != iso639_1Length {
		return "", fmt.Errorf("validation error: language must be a 2-letter ISO 639-1 code (see %s)", helpISO6391)
	}
	return fmt.Sprintf("with_original_language=%s&", qp.Language), nil
}

func (qp *queryParams) handleYear() (string, error) {
	qp.Year = cleanString(qp.Year)
	parts := strings.Split(qp.Year, ",")
	if len(parts) > 2 {
		return "", fmt.Errorf(`year format: use "2000", "2000,gte", "2000,lte", or "2000,2010"`)
	}
	year, err := validateYear(parts[0])
	if err != nil {
		return "", err
	}
	if len(parts) == 1 {
		return fmt.Sprintf("primary_release_year=%s&", year), nil
	}
	if isValidComparison(parts[1]) {
		return fmt.Sprintf("primary_release_date.%s=%s-01-01&", parts[1], year), nil
	}
	year2, err := validateYear(parts[1])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("primary_release_date.gte=%s-01-01&primary_release_date.lte=%s-12-31&", year, year2), nil
}

func (qp *queryParams) handleVoteAverage() (string, error) {
	qp.VoteAverage = cleanString(qp.VoteAverage)
	parts := strings.Split(qp.VoteAverage, ",")
	if len(parts) != 2 {
		return "", fmt.Errorf(`vote average format: use "7.0,8.0", "7.5,gte" or "7.5,lte"`)
	}
	val, err := validateVote(parts[0])
	if err != nil {
		return "", err
	}
	if isValidComparison(parts[1]) {
		return fmt.Sprintf("vote_average.%s=%s&", parts[1], val), nil
	}
	val2, err := validateVote(parts[1])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("vote_average.gte=%s&vote_average.lte=%s&", val, val2), nil
}

func (qp *queryParams) handleVoteCount() (string, error) {
	qp.VoteCount = cleanString(qp.VoteCount)
	parts := strings.Split(qp.VoteCount, ",")
	if len(parts) > 2 {
		return "", fmt.Errorf(`vote count format: use "500,1000", "500,gte", or "500,lte"`)
	}
	val, err := validateVoteCount(parts[0])
	if err != nil {
		return "", err
	}
	if len(parts) == 1 {
		return "", fmt.Errorf(`vote count format: use "500,1000", "500,gte", or "500,lte"`)
	}
	if isValidComparison(parts[1]) {
		return fmt.Sprintf("vote_count.%s=%s&", parts[1], val), nil
	}
	val2, err := validateVoteCount(parts[1])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("vote_count.gte=%s&vote_count.lte=%s&", val, val2), nil
}

func (qp *queryParams) handleWithGenres() (string, error) {
	query, err := handleGenres(qp.WithGenres, "with")
	if err != nil {
		return "", err
	}
	return query, nil
}

func (qp *queryParams) handleWithoutGenres() (string, error) {
	query, err := handleGenres(qp.WithoutGenres, "without")
	if err != nil {
		return "", err
	}
	return query, nil
}

func handleGenres(genres, suffix string) (string, error) {
	if suffix != "with" && suffix != "without" {
		return "", fmt.Errorf(`validation error: suffix must be "with" or "without"`)
	}
	var strIDs strings.Builder
	genres = cleanString(genres)
	genresList := strings.Split(genres, ",")
	for _, g := range genresList {
		strId, err := validateGenre(g)
		if err != nil {
			return "", err
		}
		strIDs.WriteString(fmt.Sprintf("%s,", strId))
	}
	genreParam := strIDs.String()
	genreParam = strings.TrimSuffix(genreParam, ",")
	return fmt.Sprintf("%s_genres=%s&", suffix, genreParam), nil
}

func validateYear(v string) (string, error) {
	_, err := time.Parse(yearFormat, v)
	if err != nil {
		return "", fmt.Errorf(`year format: use "2000", "2000,2010", "2000,gte", or "2000,lte"`)
	}
	part1, _ := strconv.Atoi(v)
	if part1 < earliestMovie || part1 > yearNow {
		return "", fmt.Errorf("year must be between %d and %d", earliestMovie, yearNow)
	}
	return v, nil
}

func validateVote(v string) (string, error) {
	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return "", fmt.Errorf(`validation error: vote average must be a float, e.g. "7.5"`)
	}
	if value < minVoteAverage || value > maxVoteAverage {
		return "", fmt.Errorf(`vote average format: use "7.0,8.0", "7.5,gte", or "7.5,lte"`)
	}
	return v, nil
}

func validateVoteCount(v string) (string, error) {
	voteCount, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return "", fmt.Errorf(`validation error: vote count must be an integer, e.g. "1000"`)
	}
	if voteCount < minVoteCount {
		return "", fmt.Errorf("validation error: vote count must be ≥ %d", minVoteCount)
	}
	return v, nil
}

func validateGenre(v string) (string, error) {
	id, exists := genresMap[v]
	if !exists {
		var strGenres strings.Builder
		var sortedGenres []string
		for k := range genresMap {
			sortedGenres = append(sortedGenres, k)
		}
		sort.Strings(sortedGenres)
		for _, k := range sortedGenres {
			strGenres.WriteString(fmt.Sprintf("\t- %s\n", k))
		}
		return "", fmt.Errorf("validation error: genre must be one of these genres:\n%s", strGenres.String())
	}
	return strconv.Itoa(id), nil
}

func isValidComparison(v string) bool {
	return v == "gte" || v == "lte"
}

func cleanString(v string) string {
	cleanStr := strings.Trim(v, `",`)
	cleanStr = strings.TrimSpace(cleanStr)
	return cleanStr
}
