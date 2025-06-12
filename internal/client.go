// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Client constants.
const (
	IdleConnTimeoutSec = 90
	RateLimitPerSec    = 3.5
	RateLimitBurst     = 10
	LoadGenresTimeout  = 30
	BadRequestStatus   = 400
	UnauthorizedStatus = 401
	NotFoundStatus     = 404
	RateLimitStatus    = 429
)

// Client handles all TMDB API interactions with caching and rate limiting.
type Client struct {
	httpClient  *http.Client
	config      Config
	rateLimiter *rate.Limiter

	// Simple in-memory cache
	cache   map[string]cachedItem
	cacheMu sync.RWMutex

	// Genre cache for movies
	genres   map[int]string
	genresMu sync.RWMutex

	// Genre cache for TV shows
	tvGenres   map[int]string
	tvGenresMu sync.RWMutex
}

type cachedItem struct {
	data   []byte
	expiry time.Time
}

// NewClient creates a new TMDB client.
func NewClient(config Config) *Client {
	// Configure HTTP client with reasonable timeouts
	httpClient := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    IdleConnTimeoutSec * time.Second,
			DisableCompression: false,
		},
	}

	// Rate limiter: TMDB allows 40 requests per 10 seconds, be conservative
	rateLimiter := rate.NewLimiter(RateLimitPerSec, RateLimitBurst) // 3.5 req/sec with burst of 10

	client := &Client{
		httpClient:  httpClient,
		config:      config,
		rateLimiter: rateLimiter,
		cache:       make(map[string]cachedItem),
		genres:      make(map[int]string),
		tvGenres:    make(map[int]string),
	}

	// Load genres in background
	go client.loadGenres()
	go client.loadTVGenres()

	return client
}

// Movie API methods

// GetPopularMovies fetches popular movies.
func (c *Client) GetPopularMovies(ctx context.Context, maxItems int) ([]Movie, error) {
	return c.fetchMoviePages(ctx, "/movie/popular", nil, maxItems)
}

// GetTopRatedMovies fetches top-rated movies.
func (c *Client) GetTopRatedMovies(ctx context.Context, maxItems int) ([]Movie, error) {
	return c.fetchMoviePages(ctx, "/movie/top_rated", nil, maxItems)
}

// GetNowPlayingMovies fetches movies currently in theaters.
func (c *Client) GetNowPlayingMovies(ctx context.Context, maxItems int) ([]Movie, error) {
	return c.fetchMoviePages(ctx, "/movie/now_playing", nil, maxItems)
}

// GetUpcomingMovies fetches upcoming movie releases.
func (c *Client) GetUpcomingMovies(ctx context.Context, maxItems int) ([]Movie, error) {
	return c.fetchMoviePages(ctx, "/movie/upcoming", nil, maxItems)
}

// SearchMovies searches for movies by query.
func (c *Client) SearchMovies(ctx context.Context, query string, maxItems int) ([]Movie, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	params := url.Values{}
	params.Set("query", query)

	return c.fetchMoviePages(ctx, "/search/movie", params, maxItems)
}

// DiscoverMovies discovers movies with filters.
func (c *Client) DiscoverMovies(ctx context.Context, opts SearchOptions) ([]Movie, error) {
	params := c.buildDiscoverParams(opts)
	return c.fetchMoviePages(ctx, "/discover/movie", params, opts.MaxItems)
}

// TV Show API methods

// GetPopularTVShows fetches popular TV shows.
func (c *Client) GetPopularTVShows(ctx context.Context, maxItems int) ([]TVShow, error) {
	return c.fetchTVShowPages(ctx, TVPopularEndpoint, nil, maxItems)
}

// GetTopRatedTVShows fetches top-rated TV shows.
func (c *Client) GetTopRatedTVShows(ctx context.Context, maxItems int) ([]TVShow, error) {
	return c.fetchTVShowPages(ctx, TVTopRatedEndpoint, nil, maxItems)
}

// GetOnTheAirTVShows fetches TV shows currently on the air.
func (c *Client) GetOnTheAirTVShows(ctx context.Context, maxItems int) ([]TVShow, error) {
	return c.fetchTVShowPages(ctx, TVOnTheAirEndpoint, nil, maxItems)
}

// SearchTVShows searches for TV shows by query.
func (c *Client) SearchTVShows(ctx context.Context, query string, maxItems int) ([]TVShow, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	params := url.Values{}
	params.Set("query", query)

	return c.fetchTVShowPages(ctx, TVSearchEndpoint, params, maxItems)
}

// DiscoverTVShows discovers TV shows with filters.
func (c *Client) DiscoverTVShows(ctx context.Context, opts SearchOptions) ([]TVShow, error) {
	params := c.buildTVDiscoverParams(opts)
	return c.fetchTVShowPages(ctx, TVDiscoverEndpoint, params, opts.MaxItems)
}

// fetchMoviePages handles pagination and returns consolidated movie results.
//
//nolint:dupl // Similar structure to fetchTVShowPages but different types
func (c *Client) fetchMoviePages(
	ctx context.Context,
	endpoint string,
	params url.Values,
	maxItems int,
) ([]Movie, error) {
	if maxItems <= 0 {
		maxItems = 20
	}

	var allMovies []Movie
	page := 1

	for len(allMovies) < maxItems {
		pageData, err := c.fetchPageData(ctx, endpoint, params, page, "")
		if err != nil {
			return nil, err
		}
		if pageData == nil {
			break
		}

		var response TMDBResponse
		if err := json.Unmarshal(pageData, &response); err != nil {
			return nil, fmt.Errorf("parse response: %w", err)
		}

		movies := c.convertMovies(response.Results)
		allMovies = append(allMovies, movies...)

		if page >= response.TotalPages || len(allMovies) >= maxItems {
			break
		}
		page++
	}

	if len(allMovies) > maxItems {
		allMovies = allMovies[:maxItems]
	}

	return allMovies, nil
}

// fetchTVShowPages handles pagination and returns consolidated TV show results.
//
//nolint:dupl // Similar structure to fetchMoviePages but different types
func (c *Client) fetchTVShowPages(
	ctx context.Context,
	endpoint string,
	params url.Values,
	maxItems int,
) ([]TVShow, error) {
	if maxItems <= 0 {
		maxItems = 20
	}

	var allTVShows []TVShow
	page := 1

	for len(allTVShows) < maxItems {
		pageData, err := c.fetchPageData(ctx, endpoint, params, page, "tv_")
		if err != nil {
			return nil, err
		}
		if pageData == nil {
			break
		}

		var response TMDBTVResponse
		if err := json.Unmarshal(pageData, &response); err != nil {
			return nil, fmt.Errorf("parse response: %w", err)
		}

		tvShows := c.convertTVShows(response.Results)
		allTVShows = append(allTVShows, tvShows...)

		if page >= response.TotalPages || len(allTVShows) >= maxItems {
			break
		}
		page++
	}

	if len(allTVShows) > maxItems {
		allTVShows = allTVShows[:maxItems]
	}

	return allTVShows, nil
}

// fetchPageData fetches a single page of data, either from cache or API.
func (c *Client) fetchPageData(
	ctx context.Context,
	endpoint string,
	params url.Values,
	page int,
	cachePrefix string,
) ([]byte, error) {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	cacheKey := c.buildCacheKey(cachePrefix, endpoint, params, page)

	// Try cache first
	if data := c.getFromCache(cacheKey); data != nil {
		return data, nil
	}

	// Make API call with rate limiting
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	// Make request
	body, err := c.makeGenericRequest(ctx, endpoint, params, page)
	if err != nil {
		return nil, err
	}

	// Cache the response
	c.putInCache(cacheKey, body)

	return body, nil
}

// buildCacheKey creates a cache key for the given parameters.
func (c *Client) buildCacheKey(prefix, endpoint string, params url.Values, page int) string {
	pageParams := make(url.Values)
	maps.Copy(pageParams, params)
	pageParams.Set("page", fmt.Sprintf("%d", page))

	return fmt.Sprintf("%s%s?%s", prefix, endpoint, pageParams.Encode())
}

// makeGenericRequest is a generic function to make HTTP requests.
func (c *Client) makeGenericRequest(
	ctx context.Context,
	endpoint string,
	params url.Values,
	page int,
) ([]byte, error) {
	pageParams := make(url.Values)
	maps.Copy(pageParams, params)
	pageParams.Set("page", fmt.Sprintf("%d", page))

	return c.makeHTTPRequest(ctx, endpoint, pageParams)
}

// makeHTTPRequest makes the actual HTTP request to TMDB API.
func (c *Client) makeHTTPRequest(
	ctx context.Context,
	endpoint string,
	params url.Values,
) ([]byte, error) {
	// Build URL
	apiURL, err := url.Parse(c.config.BaseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Add API key and params
	q := apiURL.Query()
	q.Set("api_key", c.config.APIKey)
	for key, values := range params {
		for _, value := range values {
			q.Add(key, value)
		}
	}
	apiURL.RawQuery = q.Encode()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "tmdb-cli/2.0")

	// Make request with retries
	var resp *http.Response
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			if attempt == c.config.MaxRetries {
				return nil, fmt.Errorf(
					"request failed after %d attempts: %w",
					c.config.MaxRetries+1,
					err,
				)
			}
			time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}
		break
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Check for API errors
	if resp.StatusCode >= BadRequestStatus {
		return nil, c.handleAPIError(resp.StatusCode, body)
	}

	return body, nil
}

// convertMovies converts TMDB API response to our simplified format.
func (c *Client) convertMovies(tmdbMovies []TMDBMovie) []Movie {
	movies := make([]Movie, len(tmdbMovies))

	for i, tm := range tmdbMovies {
		// Map genres
		genreNames := c.mapGenres(tm.GenreIDs)

		movies[i] = Movie{
			ID:            tm.ID,
			Title:         tm.Title,
			OriginalTitle: tm.OriginalTitle,
			Year:          ParseYear(tm.ReleaseDate),
			Rating:        tm.VoteAverage,
			Votes:         tm.VoteCount,
			Popularity:    tm.Popularity,
			Genres:        strings.Join(genreNames, ", "),
			Overview:      tm.Overview,
			Language:      "en", // Simplified for now
			Adult:         tm.Adult,
			ReleaseDate:   tm.ReleaseDate,
		}
	}

	return movies
}

// convertTVShows converts TMDB TV API response to our simplified format.
func (c *Client) convertTVShows(tmdbTVShows []TMDBTVShow) []TVShow {
	shows := make([]TVShow, len(tmdbTVShows))

	for i, tm := range tmdbTVShows {
		shows[i] = c.convertTVShowSafely(tm)
	}

	return shows
}

// convertTVShowSafely converts a single TMDB TV show with error handling.
func (c *Client) convertTVShowSafely(tmdbShow TMDBTVShow) TVShow {
	// Map genres
	genreNames := c.mapTVGenres(tmdbShow.GenreIDs)

	show := TVShow{
		ID:           tmdbShow.ID,
		Name:         tmdbShow.Name,
		OriginalName: tmdbShow.OriginalName,
		Year:         ParseTVYear(tmdbShow.FirstAirDate),
		Rating:       tmdbShow.VoteAverage,
		Votes:        tmdbShow.VoteCount,
		Popularity:   tmdbShow.Popularity,
		Genres:       strings.Join(genreNames, ", "),
		Overview:     tmdbShow.Overview,
		Language:     tmdbShow.OriginalLanguage,
		Adult:        tmdbShow.Adult,
		FirstAirDate: tmdbShow.FirstAirDate,
	}

	// Fallback for empty name
	if show.Name == "" {
		show.Name = "Unknown TV Show"
	}

	return show
}

// mapGenres converts genre IDs to names for movies.
func (c *Client) mapGenres(genreIDs []int) []string {
	c.genresMu.RLock()
	defer c.genresMu.RUnlock()

	names := make([]string, 0, len(genreIDs))
	for _, id := range genreIDs {
		if name, ok := c.genres[id]; ok {
			names = append(names, name)
		}
	}
	return names
}

// mapTVGenres converts genre IDs to names for TV shows.
func (c *Client) mapTVGenres(genreIDs []int) []string {
	c.tvGenresMu.RLock()
	defer c.tvGenresMu.RUnlock()

	names := make([]string, 0, len(genreIDs))
	for _, id := range genreIDs {
		if name, ok := c.tvGenres[id]; ok {
			names = append(names, name)
		}
	}
	return names
}

// loadGenres loads movie genre mappings from API.
func (c *Client) loadGenres() {
	c.loadGenreMapping("/genre/movie/list", func(genres []Genre) {
		c.genresMu.Lock()
		defer c.genresMu.Unlock()
		for _, genre := range genres {
			c.genres[genre.ID] = genre.Name
		}
	})
}

// loadTVGenres loads TV show genre mappings from API.
func (c *Client) loadTVGenres() {
	c.loadGenreMapping(TVGenresEndpoint, func(genres []Genre) {
		c.tvGenresMu.Lock()
		defer c.tvGenresMu.Unlock()
		for _, genre := range genres {
			c.tvGenres[genre.ID] = genre.Name
		}
	})
}

// loadGenreMapping is a generic function to load genre mappings.
func (c *Client) loadGenreMapping(endpoint string, handler func([]Genre)) {
	ctx, cancel := context.WithTimeout(context.Background(), LoadGenresTimeout*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.config.BaseURL+endpoint+"?api_key="+c.config.APIKey,
		nil,
	)
	if err != nil {
		return
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var genresResp TMDBGenresResponse
	if err := json.NewDecoder(resp.Body).Decode(&genresResp); err != nil {
		return
	}

	handler(genresResp.Genres)
}

// buildDiscoverParams converts search options to API parameters for movies.
func (c *Client) buildDiscoverParams(opts SearchOptions) url.Values {
	params := c.buildBaseDiscoverParams(opts)

	// Movie-specific parameters
	if opts.Year > 0 {
		params.Set("primary_release_year", fmt.Sprintf("%d", opts.Year))
	}

	return params
}

// buildTVDiscoverParams converts search options to API parameters for TV shows.
func (c *Client) buildTVDiscoverParams(opts SearchOptions) url.Values {
	params := c.buildBaseDiscoverParams(opts)

	// TV-specific parameters
	if opts.Year > 0 {
		params.Set("first_air_date_year", fmt.Sprintf("%d", opts.Year))
	}

	return params
}

// buildBaseDiscoverParams builds common discover parameters for both movies and TV shows.
func (c *Client) buildBaseDiscoverParams(opts SearchOptions) url.Values {
	params := make(url.Values)

	if opts.Language != "" {
		params.Set("with_original_language", opts.Language)
	}
	if opts.MinRating > 0 {
		params.Set("vote_average.gte", fmt.Sprintf("%.1f", opts.MinRating))
	}
	if opts.MaxRating > 0 {
		params.Set("vote_average.lte", fmt.Sprintf("%.1f", opts.MaxRating))
	}

	// Genres
	if len(opts.IncludeGenres) > 0 {
		genreStr := make([]string, len(opts.IncludeGenres))
		for i, id := range opts.IncludeGenres {
			genreStr[i] = fmt.Sprintf("%d", id)
		}
		params.Set("with_genres", strings.Join(genreStr, ","))
	}

	if len(opts.ExcludeGenres) > 0 {
		genreStr := make([]string, len(opts.ExcludeGenres))
		for i, id := range opts.ExcludeGenres {
			genreStr[i] = fmt.Sprintf("%d", id)
		}
		params.Set("without_genres", strings.Join(genreStr, ","))
	}

	// Default sorting
	sortBy := "popularity"
	if opts.SortBy != "" {
		sortBy = opts.SortBy
	}
	sortOrder := "desc"
	if opts.SortOrder != "" {
		sortOrder = opts.SortOrder
	}
	params.Set("sort_by", sortBy+"."+sortOrder)

	return params
}

// Cache management.
func (c *Client) getFromCache(key string) []byte {
	c.cacheMu.RLock()
	defer c.cacheMu.RUnlock()

	item, ok := c.cache[key]
	if !ok || time.Now().After(item.expiry) {
		return nil
	}
	return item.data
}

func (c *Client) putInCache(key string, data []byte) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	c.cache[key] = cachedItem{
		data:   data,
		expiry: time.Now().Add(c.config.CacheTTL),
	}

	// Simple cleanup: remove expired items periodically
	const maxCacheSize = 100
	if len(c.cache) > maxCacheSize {
		go c.cleanupCache()
	}
}

func (c *Client) cleanupCache() {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	now := time.Now()
	for key, item := range c.cache {
		if now.After(item.expiry) {
			delete(c.cache, key)
		}
	}
}

// handleAPIError creates user-friendly error messages.
func (c *Client) handleAPIError(statusCode int, body []byte) error {
	switch statusCode {
	case UnauthorizedStatus:
		return fmt.Errorf(
			"invalid TMDB API key - get one from https://www.themoviedb.org/settings/api",
		)
	case NotFoundStatus:
		return fmt.Errorf("resource not found")
	case RateLimitStatus:
		return fmt.Errorf("rate limit exceeded - please wait and try again")
	default:
		return fmt.Errorf("API error (status %d): %s", statusCode, string(body))
	}
}
