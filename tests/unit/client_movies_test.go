package unit

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

// MockTransport is a mock HTTP transport for testing.
type MockTransport struct {
	mu           sync.RWMutex
	responses    map[string]helpers.MockHTTPResponse
	requestCount int
	lastURL      string
}

// NewMockTransport creates a new mock transport.
func NewMockTransport() *MockTransport {
	return &MockTransport{
		responses: make(map[string]helpers.MockHTTPResponse),
	}
}

// SetResponse sets a response for a URL pattern.
func (t *MockTransport) SetResponse(pattern string, response helpers.MockHTTPResponse) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.responses[pattern] = response
}

// RoundTrip implements http.RoundTripper.
func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mu.Lock()
	t.requestCount++
	t.lastURL = req.URL.String()
	t.mu.Unlock()

	t.mu.RLock()
	defer t.mu.RUnlock()

	// Check for exact match first
	if resp, ok := t.responses[req.URL.Path]; ok {
		return &http.Response{
			StatusCode: resp.StatusCode,
			Body:       io.NopCloser(bytes.NewReader(resp.Body)),
			Header:     make(http.Header),
		}, nil
	}

	// Check for pattern matches (e.g., "page=1", "page=2")
	for pattern, resp := range t.responses {
		if bytes.Contains([]byte(req.URL.String()), []byte(pattern)) {
			return &http.Response{
				StatusCode: resp.StatusCode,
				Body:       io.NopCloser(bytes.NewReader(resp.Body)),
				Header:     make(http.Header),
			}, nil
		}
	}

	// Default error response
	return &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"error":"not found"}`))),
		Header:     make(http.Header),
	}, nil
}

// createInstantTestClient creates a test client with instant rate limiting for fast tests.
func createInstantTestClient(config internal.Config, httpClient *http.Client) *internal.Client {
	// Create a client with an instant rate limiter for tests
	client := internal.NewTestClient(config, httpClient, nil)

	// Replace the rate limiter with one that allows instant requests
	// This requires adding a SetRateLimiter method to the Client
	// For now, we'll use a high rate limit to minimize delays
	instantRateLimiter := rate.NewLimiter(rate.Limit(1000), 1000) // Very high rate
	client.SetRateLimiter(instantRateLimiter)

	return client
}

func TestGetPopularMovies_SinglePage(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/movie/popular", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 1, 20, 1),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClient(config, httpClient)

	movies, err := client.GetPopularMovies(context.Background(), 20)

	assert.NoError(t, err)
	assert.Len(t, movies, 1)
	assert.Equal(t, "Test Movie 1", movies[0].Title)
	assert.Equal(t, 1, transport.requestCount)
}

func TestGetPopularMovies_Pagination(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("page=1", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 3, 60, 20),
	})
	transport.SetResponse("page=2", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(2, 3, 60, 20),
	})
	transport.SetResponse("page=3", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(3, 3, 60, 10),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClient(config, httpClient)

	movies, err := client.GetPopularMovies(context.Background(), 50)

	assert.NoError(t, err)
	assert.Len(t, movies, 50)
	assert.Equal(t, 3, transport.requestCount)
}

func TestGetPopularMovies_MaxItemsLimit(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/movie/popular", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 1, 20, 20),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClient(config, httpClient)

	movies, err := client.GetPopularMovies(context.Background(), 5)

	assert.NoError(t, err)
	assert.Len(t, movies, 5)
	assert.Equal(t, 1, transport.requestCount)
}

func TestGetPopularMovies_APIError(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/movie/popular", helpers.MockHTTPResponse{
		StatusCode: 401,
		Body:       []byte(`{"status_message":"Invalid API key","status_code":7}`),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClient(config, httpClient)

	movies, err := client.GetPopularMovies(context.Background(), 20)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid TMDB API key")
	assert.Nil(t, movies)
}

func TestSearchMovies_EmptyQuery(t *testing.T) {
	t.Parallel()

	config := helpers.MockValidConfig()
	client := internal.NewClient(config)

	movies, err := client.SearchMovies(context.Background(), "", 20)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "search query cannot be empty")
	assert.Nil(t, movies)
}

func TestSearchMovies_Success(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/search/movie", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 1, 5, 5),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClient(config, httpClient)

	movies, err := client.SearchMovies(context.Background(), "Matrix", 20)

	assert.NoError(t, err)
	assert.NotNil(t, movies)
	assert.Len(t, movies, 5)
	assert.Contains(t, transport.lastURL, "query=")
	assert.Contains(t, transport.lastURL, "Matrix")
}

func TestSearchMovies_Pagination(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("page=1", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 2, 30, 20),
	})
	transport.SetResponse("page=2", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(2, 2, 30, 10),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClient(config, httpClient)

	movies, err := client.SearchMovies(context.Background(), "action", 25)

	assert.NoError(t, err)
	assert.Len(t, movies, 25)
	assert.Equal(t, 2, transport.requestCount)
}

func TestDiscoverMovies_Parameters(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/discover/movie", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 1, 10, 10),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClient(config, httpClient)

	opts := helpers.NewSearchOptionsBuilder().
		WithIncludeGenres([]int{28, 35}).
		WithMinRating(7.0).
		WithMaxRating(9.0).
		WithYear(2023).
		WithMinVotes(1000).
		WithMaxItems(10).
		Build()

	movies, err := client.DiscoverMovies(context.Background(), opts)

	assert.NoError(t, err)
	assert.NotNil(t, movies)
	assert.Len(t, movies, 10)

	// URL should contain parameters
	url := transport.lastURL
	assert.Contains(t, url, "with_genres=")
	assert.Contains(t, url, "vote_average.gte=7")
	assert.Contains(t, url, "vote_average.lte=9")
	assert.Contains(t, url, "primary_release_year=2023")
}

func TestConvertMovies(t *testing.T) {
	t.Parallel()

	config := helpers.MockValidConfig()
	client := internal.NewClient(config)

	client.SetMovieGenre(28, "Action")
	client.SetMovieGenre(12, "Adventure")

	tmdbMovies := []internal.TMDBMovie{
		{
			ID:            1,
			Title:         "Test Movie",
			OriginalTitle: "Original Test Movie",
			ReleaseDate:   "2023-06-15",
			VoteAverage:   8.5,
			VoteCount:     1000,
			Popularity:    75.5,
			GenreIDs:      []int{28, 12},
			Overview:      "Test overview",
			Adult:         false,
		},
	}

	movies := client.ConvertMovies(tmdbMovies)

	require.Len(t, movies, 1)
	assert.Equal(t, 1, movies[0].ID)
	assert.Equal(t, "Test Movie", movies[0].Title)
	assert.Equal(t, 2023, movies[0].Year)
	assert.Equal(t, "Action, Adventure", movies[0].Genres)
}

func TestBuildBaseDiscoverParams(t *testing.T) {
	t.Parallel()

	config := helpers.MockValidConfig()
	client := internal.NewClient(config)

	opts := internal.SearchOptions{
		MinRating:     7.0,
		MaxRating:     9.0,
		MinVotes:      1000,
		IncludeGenres: []int{28, 35},
		ExcludeGenres: []int{27},
		Language:      "en",
		SortBy:        "popularity",
		SortOrder:     "desc",
	}

	params := client.BuildBaseDiscoverParams(opts)

	assert.Equal(t, "7.0", params.Get("vote_average.gte"))
	assert.Equal(t, "9.0", params.Get("vote_average.lte"))
	assert.Equal(t, "1000", params.Get("vote_count.gte"))
	assert.Equal(t, "28,35", params.Get("with_genres"))
	assert.Equal(t, "27", params.Get("without_genres"))
	assert.Equal(t, "en", params.Get("with_original_language"))
	assert.Equal(t, "popularity.desc", params.Get("sort_by"))
}

func TestGetPopularMovies_Cache(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/movie/popular", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 1, 20, 20),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	config.CacheTTL = 5 * time.Minute
	client := createInstantTestClient(config, httpClient)

	// First request
	movies1, err := client.GetPopularMovies(context.Background(), 20)
	assert.NoError(t, err)
	assert.Len(t, movies1, 20)

	// Second request should use cache
	movies2, err := client.GetPopularMovies(context.Background(), 20)
	assert.NoError(t, err)
	assert.Len(t, movies2, 20)

	// Only one HTTP request should have been made
	assert.Equal(t, 1, transport.requestCount)
}

func TestSearchMovies_RateLimit(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	// First request returns rate limit error
	transport.SetResponse("/search/movie", helpers.MockHTTPResponse{
		StatusCode: 429,
		Body:       []byte(`{"status_message":"Rate limit exceeded"}`),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	config.MaxRetries = 0 // No retries
	client := createInstantTestClient(config, httpClient)

	movies, err := client.SearchMovies(context.Background(), "test", 5)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rate limit exceeded")
	assert.Nil(t, movies)
}

// Additional test cases for other movie endpoints.
func TestGetTopRatedMovies(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/movie/top_rated", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 1, 15, 15),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClient(config, httpClient)

	movies, err := client.GetTopRatedMovies(context.Background(), 15)

	assert.NoError(t, err)
	assert.Len(t, movies, 15)
}

func TestGetNowPlayingMovies(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/movie/now_playing", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 1, 10, 10),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClient(config, httpClient)

	movies, err := client.GetNowPlayingMovies(context.Background(), 10)

	assert.NoError(t, err)
	assert.Len(t, movies, 10)
}

func TestGetUpcomingMovies(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/movie/upcoming", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 1, 8, 8),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClient(config, httpClient)

	movies, err := client.GetUpcomingMovies(context.Background(), 8)

	assert.NoError(t, err)
	assert.Len(t, movies, 8)
}

// Test error handling for various scenarios.
func TestMovieOperations_NetworkError(t *testing.T) {
	t.Parallel()

	// Create a transport that always fails
	transport := &failingTransport{
		err: errors.New("network error"),
	}

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	config.MaxRetries = 0 // No retries for faster test
	client := createInstantTestClient(config, httpClient)

	movies, err := client.GetPopularMovies(context.Background(), 20)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "network error")
	assert.Nil(t, movies)
}

// Helper transport that always fails.
type failingTransport struct {
	err error
}

func (t *failingTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	return nil, t.err
}

// Test concurrent access.
func TestMovieOperations_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	// Set up responses for all endpoints
	endpoints := []string{
		"/movie/popular",
		"/movie/top_rated",
		"/movie/now_playing",
		"/movie/upcoming",
	}

	for _, endpoint := range endpoints {
		transport.SetResponse(endpoint, helpers.MockHTTPResponse{
			StatusCode: 200,
			Body:       fixtures.CreateMoviePageResponse(1, 1, 5, 5),
		})
	}

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClient(config, httpClient)

	// Run concurrent requests
	var wg sync.WaitGroup
	errors := make(chan error, len(endpoints))

	for _, endpoint := range endpoints {
		wg.Add(1)
		go func(ep string) {
			defer wg.Done()

			var err error
			switch ep {
			case "/movie/popular":
				_, err = client.GetPopularMovies(context.Background(), 5)
			case "/movie/top_rated":
				_, err = client.GetTopRatedMovies(context.Background(), 5)
			case "/movie/now_playing":
				_, err = client.GetNowPlayingMovies(context.Background(), 5)
			case "/movie/upcoming":
				_, err = client.GetUpcomingMovies(context.Background(), 5)
			}

			if err != nil {
				errors <- err
			}
		}(endpoint)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		assert.NoError(t, err)
	}
}

// Benchmark test for pagination.
func BenchmarkSearchMovies_Pagination(b *testing.B) {
	transport := NewMockTransport()
	transport.SetResponse("page=1", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 2, 30, 20),
	})
	transport.SetResponse("page=2", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(2, 2, 30, 10),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClient(config, httpClient)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.SearchMovies(ctx, "action", 25)
	}
}
