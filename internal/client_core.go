// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"net/http"
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
