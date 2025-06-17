package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

// GenreLoadTestTransport implements http.RoundTripper for genre loading tests.
type GenreLoadTestTransport struct {
	mu        sync.Mutex
	responses map[string]*http.Response
	errors    map[string]error
	requests  []string
}

func (g *GenreLoadTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.requests = append(g.requests, req.URL.Path)

	if err, ok := g.errors[req.URL.Path]; ok && err != nil {
		return nil, err
	}

	if resp, ok := g.responses[req.URL.Path]; ok {
		return resp, nil
	}

	return &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(bytes.NewBufferString("not found")),
	}, nil
}

func TestClient_GenreMapping(t *testing.T) {
	t.Run("map movie genres correctly", func(t *testing.T) {
		// Setup
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		config := internal.DefaultConfig()
		config.APIKey = testKey
		client := internal.NewTestClient(config, nil, mockClock)

		// Manually set genres for testing
		client.SetMovieGenre(28, "Action")
		client.SetMovieGenre(35, "Comedy")
		client.SetMovieGenre(18, "Drama")

		// Test mapping
		genreIDs := []int{28, 35, 99} // 99 doesn't exist
		names := client.MapGenres(genreIDs)

		assert.Equal(t, []string{"Action", "Comedy"}, names)
	})

	t.Run("map TV genres correctly", func(t *testing.T) {
		// Setup
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		config := internal.DefaultConfig()
		config.APIKey = testKey
		client := internal.NewTestClient(config, nil, mockClock)

		// Manually set genres for testing
		client.SetTVGenre(10759, "Action & Adventure")
		client.SetTVGenre(16, "Animation")
		client.SetTVGenre(35, "Comedy")

		// Test mapping
		genreIDs := []int{10759, 16, 99} // 99 doesn't exist
		names := client.MapTVGenres(genreIDs)

		assert.Equal(t, []string{"Action & Adventure", "Animation"}, names)
	})

	t.Run("handle empty genre IDs", func(t *testing.T) {
		// Setup
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		config := internal.DefaultConfig()
		config.APIKey = testKey
		client := internal.NewTestClient(config, nil, mockClock)

		// Test empty mapping
		names := client.MapGenres([]int{})
		assert.Empty(t, names)

		names = client.MapTVGenres([]int{})
		assert.Empty(t, names)
	})
}

func TestClient_LoadGenres(t *testing.T) {
	t.Run("successful genre loading", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Create mock response
		movieGenres := internal.TMDBGenresResponse{
			Genres: []internal.Genre{
				{ID: 28, Name: "Action"},
				{ID: 35, Name: "Comedy"},
				{ID: 18, Name: "Drama"},
			},
		}

		tvGenres := internal.TMDBGenresResponse{
			Genres: []internal.Genre{
				{ID: 10759, Name: "Action & Adventure"},
				{ID: 16, Name: "Animation"},
				{ID: 35, Name: "Comedy"},
			},
		}

		movieBody, _ := json.Marshal(movieGenres)
		tvBody, _ := json.Marshal(tvGenres)

		// Setup transport
		transport := &GenreLoadTestTransport{
			responses: map[string]*http.Response{
				"/3/genre/movie/list": {
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBuffer(movieBody)),
					Header:     make(http.Header),
				},
				"/3/genre/tv/list": {
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBuffer(tvBody)),
					Header:     make(http.Header),
				},
			},
			errors: make(map[string]error),
		}

		httpClient := &http.Client{
			Transport: transport,
		}

		// Create client
		config := internal.DefaultConfig()
		config.APIKey = testKey
		client := internal.NewTestClient(config, httpClient, mockClock)

		// Load genres
		client.LoadGenres()
		client.LoadTVGenres()

		// Allow time for loading
		time.Sleep(100 * time.Millisecond)

		// Verify movie genres loaded
		movieNames := client.MapGenres([]int{28, 35, 18})
		assert.Equal(t, []string{"Action", "Comedy", "Drama"}, movieNames)

		// Verify TV genres loaded
		tvNames := client.MapTVGenres([]int{10759, 16, 35})
		assert.Equal(t, []string{"Action & Adventure", "Animation", "Comedy"}, tvNames)

		// Verify requests were made
		assert.Contains(t, transport.requests, "/3/genre/movie/list")
		assert.Contains(t, transport.requests, "/3/genre/tv/list")
	})

	t.Run("genre loading with timeout", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Setup transport that delays response
		transport := &GenreLoadTestTransport{
			responses: make(map[string]*http.Response),
			errors: map[string]error{
				"/3/genre/movie/list": context.DeadlineExceeded,
				"/3/genre/tv/list":    context.DeadlineExceeded,
			},
		}

		httpClient := &http.Client{
			Transport: transport,
		}

		// Create client
		config := internal.DefaultConfig()
		config.APIKey = testKey
		client := internal.NewTestClient(config, httpClient, mockClock)

		// Load genres
		client.LoadGenres()
		client.LoadTVGenres()

		// Allow time for loading attempts
		time.Sleep(50 * time.Millisecond)

		// Verify genres not loaded due to timeout
		movieNames := client.MapGenres([]int{28, 35})
		assert.Empty(t, movieNames)

		tvNames := client.MapTVGenres([]int{10759, 16})
		assert.Empty(t, tvNames)
	})

	t.Run("genre loading with API error", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Setup transport with error responses
		transport := &GenreLoadTestTransport{
			responses: map[string]*http.Response{
				"/3/genre/movie/list": {
					StatusCode: 401,
					Body:       io.NopCloser(bytes.NewBufferString(`{"error": "Invalid API key"}`)),
					Header:     make(http.Header),
				},
				"/3/genre/tv/list": {
					StatusCode: 500,
					Body:       io.NopCloser(bytes.NewBufferString(`{"error": "Server error"}`)),
					Header:     make(http.Header),
				},
			},
			errors: make(map[string]error),
		}

		httpClient := &http.Client{
			Transport: transport,
		}

		// Create client
		config := internal.DefaultConfig()
		config.APIKey = "invalid-key"
		client := internal.NewTestClient(config, httpClient, mockClock)

		// Load genres
		client.LoadGenres()
		client.LoadTVGenres()

		// Allow time for loading
		time.Sleep(50 * time.Millisecond)

		// Verify genres not loaded due to errors
		movieNames := client.MapGenres([]int{28, 35})
		assert.Empty(t, movieNames)

		tvNames := client.MapTVGenres([]int{10759, 16})
		assert.Empty(t, tvNames)
	})

	t.Run("genre loading with network error", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Setup transport with network errors
		transport := &GenreLoadTestTransport{
			responses: make(map[string]*http.Response),
			errors: map[string]error{
				"/3/genre/movie/list": errors.New("network unreachable"),
				"/3/genre/tv/list":    errors.New("connection refused"),
			},
		}

		httpClient := &http.Client{
			Transport: transport,
		}

		// Create client
		config := internal.DefaultConfig()
		config.APIKey = testKey
		client := internal.NewTestClient(config, httpClient, mockClock)

		// Load genres
		client.LoadGenres()
		client.LoadTVGenres()

		// Allow time for loading
		time.Sleep(50 * time.Millisecond)

		// Verify genres not loaded due to network errors
		movieNames := client.MapGenres([]int{28, 35})
		assert.Empty(t, movieNames)

		tvNames := client.MapTVGenres([]int{10759, 16})
		assert.Empty(t, tvNames)
	})

	t.Run("genre loading with malformed JSON", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Setup transport with malformed responses
		transport := &GenreLoadTestTransport{
			responses: map[string]*http.Response{
				"/3/genre/movie/list": {
					StatusCode: 200,
					Body: io.NopCloser(
						bytes.NewBufferString(`{"genres": [{"id": 28, "name": "Action"`),
					), // Missing closing brackets
					Header: make(http.Header),
				},
				"/3/genre/tv/list": {
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBufferString(`not json at all`)),
					Header:     make(http.Header),
				},
			},
			errors: make(map[string]error),
		}

		httpClient := &http.Client{
			Transport: transport,
		}

		// Create client
		config := internal.DefaultConfig()
		config.APIKey = testKey
		client := internal.NewTestClient(config, httpClient, mockClock)

		// Load genres
		client.LoadGenres()
		client.LoadTVGenres()

		// Allow time for loading
		time.Sleep(50 * time.Millisecond)

		// Verify genres not loaded due to JSON errors
		movieNames := client.MapGenres([]int{28, 35})
		assert.Empty(t, movieNames)

		tvNames := client.MapTVGenres([]int{10759, 16})
		assert.Empty(t, tvNames)
	})
}

func TestClient_ConcurrentGenreAccess(t *testing.T) {
	t.Run("concurrent reads and writes are safe", func(t *testing.T) {
		// Setup
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		config := internal.DefaultConfig()
		config.APIKey = testKey
		client := internal.NewTestClient(config, nil, mockClock)

		// Run concurrent operations
		var wg sync.WaitGroup
		numGoroutines := 10
		numOperations := 100

		// Writers
		for i := range numGoroutines {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := range numOperations {
					genreID := id*1000 + j
					client.SetMovieGenre(genreID, helpers.CreateValue("Genre", id, j))
					client.SetTVGenre(genreID, helpers.CreateValue("TVGenre", id, j))
				}
			}(i)
		}

		// Readers
		for i := range numGoroutines {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for range numOperations {
					// Read random genre IDs
					genreIDs := []int{id * 1000, id*1000 + 50, id*1000 + 99}
					_ = client.MapGenres(genreIDs)
					_ = client.MapTVGenres(genreIDs)
				}
			}(i)
		}

		// Wait for all operations
		wg.Wait()

		// Verify some data is present
		for i := range numGoroutines {
			genreID := i*1000 + numOperations - 1
			expectedName := helpers.CreateValue("Genre", i, numOperations-1)

			names := client.MapGenres([]int{genreID})
			require.Len(t, names, 1)
			assert.Equal(t, expectedName, names[0])

			expectedTVName := helpers.CreateValue("TVGenre", i, numOperations-1)
			tvNames := client.MapTVGenres([]int{genreID})
			require.Len(t, tvNames, 1)
			assert.Equal(t, expectedTVName, tvNames[0])
		}
	})
}

// TimeoutCapturingClock wraps MockClock to capture timeout calls.
type TimeoutCapturingClock struct {
	*internal.MockClock
	TimeoutCalled   bool
	CapturedTimeout time.Duration
}

func (t *TimeoutCapturingClock) WithTimeout(
	parent context.Context,
	timeout time.Duration,
) (context.Context, context.CancelFunc) {
	t.TimeoutCalled = true
	t.CapturedTimeout = timeout
	return context.WithTimeout(parent, timeout)
}

func TestClient_GenreTimeoutBehavior(t *testing.T) {
	t.Run("verify timeout is enforced", func(t *testing.T) {
		// Create a custom mock clock that captures timeout calls
		mockClock := &TimeoutCapturingClock{
			MockClock: &internal.MockClock{
				CurrentTime: time.Now(),
			},
		}

		// Create mock genres response
		genres := internal.TMDBGenresResponse{
			Genres: []internal.Genre{
				{ID: 28, Name: "Action"},
			},
		}
		body, _ := json.Marshal(genres)

		// Setup transport with successful response
		transport := &GenreLoadTestTransport{
			responses: map[string]*http.Response{
				"/3/genre/movie/list": {
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewBuffer(body)),
					Header:     make(http.Header),
				},
			},
			errors: make(map[string]error),
		}

		httpClient := &http.Client{
			Transport: transport,
		}

		// Create client
		config := internal.DefaultConfig()
		config.APIKey = testKey
		client := internal.NewTestClient(config, httpClient, mockClock)

		// Load genres (this will trigger WithTimeout)
		client.LoadGenres()

		// Allow time for loading
		time.Sleep(50 * time.Millisecond)

		// Verify timeout was called with correct duration
		assert.True(t, mockClock.TimeoutCalled, "WithTimeout should have been called")
		assert.Equal(
			t,
			30*time.Second,
			mockClock.CapturedTimeout,
			"Timeout should be LoadGenresTimeout",
		)
	})
}
