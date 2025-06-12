// internal/client.go
package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Client handles all TMDB API interactions with caching and rate limiting
type Client struct {
	httpClient  *http.Client
	config      Config
	rateLimiter *rate.Limiter

	// Simple in-memory cache
	cache   map[string]cachedItem
	cacheMu sync.RWMutex

	// Genre cache
	genres   map[int]string
	genresMu sync.RWMutex
}

type cachedItem struct {
	data   []byte
	expiry time.Time
}

// TMDBResponse represents raw API response structure
type TMDBResponse struct {
	Page         int         `json:"page"`
	Results      []TMDBMovie `json:"results"`
	TotalPages   int         `json:"total_pages"`
	TotalResults int         `json:"total_results"`
}

type TMDBMovie struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title"`
	Overview      string  `json:"overview"`
	ReleaseDate   string  `json:"release_date"`
	VoteAverage   float64 `json:"vote_average"`
	VoteCount     int     `json:"vote_count"`
	GenreIDs      []int   `json:"genre_ids"`
	Popularity    float64 `json:"popularity"`
	Adult         bool    `json:"adult"`
	Video         bool    `json:"video"`
}

type TMDBGenresResponse struct {
	Genres []Genre `json:"genres"`
}

// NewClient creates a new TMDB client
func NewClient(config Config) *Client {
	// Configure HTTP client with reasonable timeouts
	httpClient := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: false,
		},
	}

	// Rate limiter: TMDB allows 40 requests per 10 seconds, be conservative
	rateLimiter := rate.NewLimiter(3.5, 10) // 3.5 req/sec with burst of 10

	client := &Client{
		httpClient:  httpClient,
		config:      config,
		rateLimiter: rateLimiter,
		cache:       make(map[string]cachedItem),
		genres:      make(map[int]string),
	}

	// Load genres in background
	go client.loadGenres()

	return client
}

// GetPopularMovies fetches popular movies
func (c *Client) GetPopularMovies(ctx context.Context, maxItems int) ([]Movie, error) {
	return c.fetchMoviePages(ctx, "/movie/popular", nil, maxItems)
}

// GetTopRatedMovies fetches top-rated movies
func (c *Client) GetTopRatedMovies(ctx context.Context, maxItems int) ([]Movie, error) {
	return c.fetchMoviePages(ctx, "/movie/top_rated", nil, maxItems)
}

// GetNowPlayingMovies fetches movies currently in theaters
func (c *Client) GetNowPlayingMovies(ctx context.Context, maxItems int) ([]Movie, error) {
	return c.fetchMoviePages(ctx, "/movie/now_playing", nil, maxItems)
}

// GetUpcomingMovies fetches upcoming movie releases
func (c *Client) GetUpcomingMovies(ctx context.Context, maxItems int) ([]Movie, error) {
	return c.fetchMoviePages(ctx, "/movie/upcoming", nil, maxItems)
}

// SearchMovies searches for movies by query
func (c *Client) SearchMovies(ctx context.Context, query string, maxItems int) ([]Movie, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	params := url.Values{}
	params.Set("query", query)

	return c.fetchMoviePages(ctx, "/search/movie", params, maxItems)
}

// DiscoverMovies discovers movies with filters
func (c *Client) DiscoverMovies(ctx context.Context, opts SearchOptions) ([]Movie, error) {
	params := c.buildDiscoverParams(opts)
	return c.fetchMoviePages(ctx, "/discover/movie", params, opts.MaxItems)
}

// fetchMoviePages handles pagination and returns consolidated results
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
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Prepare params for this page
		pageParams := make(url.Values)
		if params != nil {
			pageParams = params
		}
		pageParams.Set("page", fmt.Sprintf("%d", page))

		// Build cache key
		cacheKey := fmt.Sprintf("%s?%s", endpoint, pageParams.Encode())

		// Try cache first
		if data := c.getFromCache(cacheKey); data != nil {
			var response TMDBResponse
			if err := json.Unmarshal(data, &response); err == nil {
				movies := c.convertMovies(response.Results)
				allMovies = append(allMovies, movies...)

				// Check if we have more pages and need more items
				if page >= response.TotalPages || len(allMovies) >= maxItems {
					break
				}
				page++
				continue
			}
		}

		// Make API call with rate limiting
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limit: %w", err)
		}

		response, err := c.makeRequest(ctx, endpoint, pageParams)
		if err != nil {
			return nil, err
		}

		// Cache the response
		if data, err := json.Marshal(response); err == nil {
			c.putInCache(cacheKey, data)
		}

		// Convert and append movies
		movies := c.convertMovies(response.Results)
		allMovies = append(allMovies, movies...)

		// Check if we have more pages and need more items
		if page >= response.TotalPages || len(allMovies) >= maxItems {
			break
		}
		page++
	}

	// Trim to requested size
	if len(allMovies) > maxItems {
		allMovies = allMovies[:maxItems]
	}

	return allMovies, nil
}

// makeRequest makes HTTP request to TMDB API
func (c *Client) makeRequest(
	ctx context.Context,
	endpoint string,
	params url.Values,
) (*TMDBResponse, error) {
	// Build URL
	u, err := url.Parse(c.config.BaseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Add API key and params
	q := u.Query()
	q.Set("api_key", c.config.APIKey)
	for key, values := range params {
		for _, value := range values {
			q.Add(key, value)
		}
	}
	u.RawQuery = q.Encode()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
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
	defer func() { _ = resp.Body.Close() }()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Check for API errors
	if resp.StatusCode >= 400 {
		return nil, c.handleAPIError(resp.StatusCode, body)
	}

	// Parse response
	var response TMDBResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &response, nil
}

// convertMovies converts TMDB API response to our simplified format
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

// mapGenres converts genre IDs to names
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

// loadGenres loads genre mappings from API
func (c *Client) loadGenres() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.config.BaseURL+"/genre/movie/list?api_key="+c.config.APIKey,
		nil,
	)
	if err != nil {
		return
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var genresResp TMDBGenresResponse
	if err := json.NewDecoder(resp.Body).Decode(&genresResp); err != nil {
		return
	}

	c.genresMu.Lock()
	defer c.genresMu.Unlock()

	for _, genre := range genresResp.Genres {
		c.genres[genre.ID] = genre.Name
	}
}

// buildDiscoverParams converts search options to API parameters
func (c *Client) buildDiscoverParams(opts SearchOptions) url.Values {
	params := make(url.Values)

	if opts.Language != "" {
		params.Set("with_original_language", opts.Language)
	}
	if opts.Year > 0 {
		params.Set("primary_release_year", fmt.Sprintf("%d", opts.Year))
	}
	if opts.MinRating > 0 {
		params.Set("vote_average.gte", fmt.Sprintf("%.1f", opts.MinRating))
	}
	if opts.MaxRating > 0 {
		params.Set("vote_average.lte", fmt.Sprintf("%.1f", opts.MaxRating))
	}
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

// Cache management
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
	if len(c.cache) > 100 {
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

// handleAPIError creates user-friendly error messages
func (c *Client) handleAPIError(statusCode int, body []byte) error {
	switch statusCode {
	case 401:
		return fmt.Errorf(
			"invalid TMDB API key - get one from https://www.themoviedb.org/settings/api",
		)
	case 404:
		return fmt.Errorf("resource not found")
	case 429:
		return fmt.Errorf("rate limit exceeded - please wait and try again")
	default:
		return fmt.Errorf("API error (status %d): %s", statusCode, string(body))
	}
}
