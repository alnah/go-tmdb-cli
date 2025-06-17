package internal

import (
	"context"
	"net/http"
	"net/url"

	"golang.org/x/time/rate"
)

// ClientTestHelper provides access to internal client methods for testing.
type ClientTestHelper struct {
	client *Client
}

// NewClientTestHelper creates a test helper for the client.
func NewClientTestHelper(client *Client) *ClientTestHelper {
	return &ClientTestHelper{client: client}
}

// The following methods would ideally be added to the actual Client struct
// or made available through interfaces for testing:

// GetHTTPClient would return the underlying HTTP client (for testing transport configuration).
// GetConfig would return the client configuration.
// SetMovieGenre would allow setting movie genres for testing genre mapping.
// SetTVGenre would allow setting TV genres for testing genre mapping.
// BuildBaseDiscoverParams would expose the parameter building logic.
// ConvertMovies would expose the movie conversion logic.
// ConvertTVShows would expose the TV show conversion logic.
// FetchMoviePages would expose the movie pagination logic.
// FetchTVShowPages would expose the TV show pagination logic.
// MakeHTTPRequest would expose the HTTP request logic.
// FetchPageData would expose the page data fetching logic.
// WaitForRateLimit would expose the rate limiting logic.
// PutInCache/GetFromCache would expose cache operations.
// LoadGenres/LoadTVGenres would expose genre loading.
// MapGenres/MapTVGenres would expose genre mapping.

// For the actual implementation, these methods should be added to the Client struct:

// GetHTTPClient returns the underlying HTTP client for testing.
func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}

// GetConfig returns the client configuration for testing.
func (c *Client) GetConfig() Config {
	return c.config
}

// SetMovieGenre sets a movie genre for testing.
func (c *Client) SetMovieGenre(id int, name string) {
	c.genresMu.Lock()
	defer c.genresMu.Unlock()
	c.genres[id] = name
}

// SetTVGenre sets a TV genre for testing.
func (c *Client) SetTVGenre(id int, name string) {
	c.tvGenresMu.Lock()
	defer c.tvGenresMu.Unlock()
	c.tvGenres[id] = name
}

// Public methods to expose internal functionality for testing.

// BuildBaseDiscoverParams is for testing purpose.
func (c *Client) BuildBaseDiscoverParams(opts SearchOptions) url.Values {
	return c.buildBaseDiscoverParams(opts)
}

// ConvertMovies is for testing purpose.
func (c *Client) ConvertMovies(tmdbMovies []TMDBMovie) []Movie {
	return c.convertMovies(tmdbMovies)
}

// ConvertTVShows is for testing purpose.
func (c *Client) ConvertTVShows(tmdbTVShows []TMDBTVShow) []TVShow {
	return c.convertTVShows(tmdbTVShows)
}

// FetchMoviePages is for testing purpose.
func (c *Client) FetchMoviePages(
	ctx context.Context,
	endpoint string,
	params url.Values,
	maxItems int,
) ([]Movie, error) {
	return c.fetchMoviePages(ctx, endpoint, params, maxItems)
}

// FetchTVShowPages is for testing purpose.
func (c *Client) FetchTVShowPages(
	ctx context.Context,
	endpoint string,
	params url.Values,
	maxItems int,
) ([]TVShow, error) {
	return c.fetchTVShowPages(ctx, endpoint, params, maxItems)
}

// MakeHTTPRequest is for testing purpose.
func (c *Client) MakeHTTPRequest(
	ctx context.Context,
	endpoint string,
	params url.Values,
) ([]byte, error) {
	return c.makeHTTPRequest(ctx, endpoint, params)
}

// FetchPageData is for testing purpose.
func (c *Client) FetchPageData(
	ctx context.Context,
	endpoint string,
	params url.Values,
	page int,
	cachePrefix string,
) ([]byte, error) {
	return c.fetchPageData(ctx, endpoint, params, page, cachePrefix)
}

// WaitForRateLimit is for testing purpose.
func (c *Client) WaitForRateLimit(ctx context.Context) error {
	return c.rateLimiter.Wait(ctx)
}

// PutInCache is for testing purpose.
func (c *Client) PutInCache(key string, data []byte) {
	c.putInCache(key, data)
}

// GetFromCache is for testing purpose.
func (c *Client) GetFromCache(key string) []byte {
	return c.getFromCache(key)
}

// LoadGenres is for testing purpose.
func (c *Client) LoadGenres() {
	c.loadGenres()
}

// LoadTVGenres is for testing purpose.
func (c *Client) LoadTVGenres() {
	c.loadTVGenres()
}

// MapGenres is for testing purpose.
func (c *Client) MapGenres(genreIDs []int) []string {
	return c.mapGenres(genreIDs)
}

// MapTVGenres is for testing purpose.
func (c *Client) MapTVGenres(genreIDs []int) []string {
	return c.mapTVGenres(genreIDs)
}

// NewClientWithHTTPClient creates a client with a custom HTTP client for testing.
func NewClientWithHTTPClient(config Config, httpClient *http.Client) *Client {
	// Similar to NewClient but allows injecting a custom HTTP client
	rateLimiter := rate.NewLimiter(RateLimitPerSec, RateLimitBurst)

	client := &Client{
		httpClient:  httpClient,
		config:      config,
		rateLimiter: rateLimiter,
		cache:       make(map[string]cachedItem),
		genres:      make(map[int]string),
		tvGenres:    make(map[int]string),
	}

	return client
}
