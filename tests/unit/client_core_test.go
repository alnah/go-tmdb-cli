package unit

import (
	"context"
	"net/http"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestNewClient(t *testing.T) {
	t.Run("creates client with valid configuration", func(t *testing.T) {
		config := helpers.MockValidConfig()
		client := internal.NewClient(config)

		require.NotNil(t, client)

		// Verify client configuration
		assert.Equal(t, config, client.GetConfig())

		// Verify HTTP client is created
		httpClient := client.GetHTTPClient()
		require.NotNil(t, httpClient)
		assert.Equal(t, config.Timeout, httpClient.Timeout)

		// Verify clock is RealClock by default
		clock := client.GetClock()
		require.NotNil(t, clock)
		_, ok := clock.(internal.RealClock)
		assert.True(t, ok, "Default clock should be RealClock")
	})

	t.Run("creates client with custom timeout", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("test-key").
			WithTimeout(60 * time.Second).
			Build()

		client := internal.NewClient(config)
		require.NotNil(t, client)

		httpClient := client.GetHTTPClient()
		assert.Equal(t, 60*time.Second, httpClient.Timeout)
	})

	t.Run("creates client with max retries configuration", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("test-key").
			WithMaxRetries(5).
			Build()

		client := internal.NewClient(config)
		require.NotNil(t, client)

		clientConfig := client.GetConfig()
		assert.Equal(t, 5, clientConfig.MaxRetries)
	})

	t.Run("initializes empty caches", func(t *testing.T) {
		config := helpers.MockValidConfig()
		client := internal.NewClient(config)

		// Test cache is empty initially
		data := client.GetFromCache("non-existent")
		assert.Nil(t, data)

		// Test genre caches are empty initially
		movieGenres := client.MapGenres([]int{28})
		assert.Empty(t, movieGenres)

		tvGenres := client.MapTVGenres([]int{10759})
		assert.Empty(t, tvGenres)
	})

	t.Run("starts background genre loading", func(t *testing.T) {
		// Create a mock HTTP transport that tracks requests
		transport := &mockTransport{
			responses: make(map[string]*http.Response),
		}

		httpClient := &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}

		config := helpers.MockValidConfig()
		client := internal.NewClientWithHTTPClient(config, httpClient)
		require.NotNil(t, client)

		// Give goroutines time to start
		time.Sleep(100 * time.Millisecond)

		// Verify that genre loading requests were attempted
		assert.True(
			t,
			transport.hasRequest("/3/genre/movie/list"),
			"Should attempt to load movie genres",
		)
		assert.True(t, transport.hasRequest("/3/genre/tv/list"), "Should attempt to load TV genres")
	})

	t.Run("http client has proper transport configuration", func(t *testing.T) {
		config := helpers.MockValidConfig()
		client := internal.NewClient(config)

		httpClient := client.GetHTTPClient()
		transport, ok := httpClient.Transport.(*http.Transport)
		require.True(t, ok, "Transport should be *http.Transport")

		assert.Equal(t, 10, transport.MaxIdleConns)
		assert.Equal(t, 90*time.Second, transport.IdleConnTimeout)
		assert.False(t, transport.DisableCompression)
	})
}

func TestNewClientWithClock(t *testing.T) {
	t.Run("creates client with mock clock", func(t *testing.T) {
		mockClock := &internal.MockClock{
			CurrentTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		}

		config := helpers.MockValidConfig()
		client := internal.NewClientWithClock(config, mockClock)

		require.NotNil(t, client)

		// Verify mock clock is used
		clock := client.GetClock()
		assert.Equal(t, mockClock, clock)
		assert.Equal(t, mockClock.CurrentTime, clock.Now())
	})

	t.Run("mock clock affects cache expiry", func(t *testing.T) {
		startTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		mockClock := &internal.MockClock{
			CurrentTime: startTime,
		}

		config := helpers.NewConfigBuilder().
			WithAPIKey("test-key").
			WithCacheTTL(5 * time.Minute).
			Build()

		client := internal.NewClientWithClock(config, mockClock)

		// Put item in cache
		client.PutInCache("test-key", []byte("test-data"))

		// Verify item is in cache
		data := client.GetFromCache("test-key")
		assert.Equal(t, []byte("test-data"), data)

		// Advance clock past TTL
		mockClock.SetTime(startTime.Add(6 * time.Minute))

		// Verify item is expired
		data = client.GetFromCache("test-key")
		assert.Nil(t, data)
	})

	t.Run("background genre loading with mock clock", func(t *testing.T) {
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Create a mock HTTP transport
		transport := &mockTransport{
			responses: make(map[string]*http.Response),
		}

		httpClient := &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}

		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, httpClient, mockClock)
		require.NotNil(t, client)

		// Manually trigger genre loading since background goroutines don't start in test client
		client.LoadGenres()
		client.LoadTVGenres()

		// Give time for requests
		time.Sleep(50 * time.Millisecond)

		// Verify requests were made
		assert.True(t, transport.hasRequest("/3/genre/movie/list"))
		assert.True(t, transport.hasRequest("/3/genre/tv/list"))
	})
}

func TestNewTestClient(t *testing.T) {
	t.Run("creates test client with nil clock uses RealClock", func(t *testing.T) {
		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, nil, nil)

		require.NotNil(t, client)

		clock := client.GetClock()
		_, ok := clock.(internal.RealClock)
		assert.True(t, ok, "Nil clock should default to RealClock")
	})

	t.Run("creates test client with nil httpClient creates default", func(t *testing.T) {
		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, nil, nil)

		require.NotNil(t, client)

		httpClient := client.GetHTTPClient()
		require.NotNil(t, httpClient)
		assert.Equal(t, config.Timeout, httpClient.Timeout)
	})

	t.Run("creates test client with custom components", func(t *testing.T) {
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		customHTTPClient := &http.Client{
			Timeout: 45 * time.Second,
		}

		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, customHTTPClient, mockClock)

		require.NotNil(t, client)

		// Verify custom components are used
		assert.Equal(t, mockClock, client.GetClock())
		assert.Equal(t, customHTTPClient, client.GetHTTPClient())
	})
}

func TestClientRateLimiting(t *testing.T) {
	t.Run("rate limiter is initialized with correct values", func(t *testing.T) {
		config := helpers.MockValidConfig()
		client := internal.NewClient(config)

		// Make rapid requests to test rate limiting
		start := time.Now()

		// Try to make requests faster than rate limit allows
		for range 5 {
			err := client.WaitForRateLimit(context.Background())
			assert.NoError(t, err)
		}

		elapsed := time.Since(start)

		// With 3.5 req/sec rate limit, 5 requests should take at least 1 second
		assert.GreaterOrEqual(t, elapsed, 1*time.Second, "Rate limiting should slow down requests")
	})

	t.Run("rate limiter respects context cancellation", func(t *testing.T) {
		config := helpers.MockValidConfig()
		client := internal.NewClient(config)

		ctx, cancel := context.WithCancel(context.Background())

		// Cancel immediately
		cancel()

		err := client.WaitForRateLimit(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")
	})
}

// mockTransport is a test helper for tracking HTTP requests.
type mockTransport struct {
	requests  []string
	responses map[string]*http.Response
}

func (mt *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	mt.requests = append(mt.requests, req.URL.Path)

	if resp, ok := mt.responses[req.URL.Path]; ok {
		return resp, nil
	}

	// Return a basic error response
	return &http.Response{
		StatusCode: 401,
		Body:       http.NoBody,
	}, nil
}

func (mt *mockTransport) hasRequest(path string) bool {
	return slices.Contains(mt.requests, path)
}
