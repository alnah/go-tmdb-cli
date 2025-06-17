package unit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestClient_fetchPageData(t *testing.T) {
	t.Run("fetches from API when cache miss", func(t *testing.T) {
		// Setup
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		transport := &mockHTTPTransport{
			responses: make(map[string]mockResponse),
		}

		// Prepare response
		tmdbResp := internal.TMDBResponse{
			Page:         1,
			Results:      []internal.TMDBMovie{fixtures.SampleTMDBMovie},
			TotalPages:   1,
			TotalResults: 1,
		}
		respBody, _ := json.Marshal(tmdbResp)

		transport.setResponse("/movie/popular", mockResponse{
			statusCode: 200,
			body:       respBody,
		})

		httpClient := &http.Client{
			Transport: transport,
		}

		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, httpClient, mockClock)

		// Test
		data, err := client.FetchPageData(
			context.Background(),
			"/movie/popular",
			nil,
			1,
			"movie_",
		)

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, data)

		// Verify API was called
		assert.Equal(t, 1, transport.getRequestCount())
		assert.Contains(t, transport.getLastRequestURL(), "/movie/popular")
		assert.Contains(t, transport.getLastRequestURL(), "page=1")
		assert.Contains(t, transport.getLastRequestURL(), "api_key=")
	})

	t.Run("returns cached data on cache hit", func(t *testing.T) {
		// Setup
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		transport := &mockHTTPTransport{
			responses: make(map[string]mockResponse),
		}

		httpClient := &http.Client{
			Transport: transport,
		}

		config := helpers.NewConfigBuilder().
			WithAPIKey("test-key").
			WithCacheTTL(10 * time.Minute).
			Build()

		client := internal.NewTestClient(config, httpClient, mockClock)

		// Pre-populate cache
		cachedData := []byte(`{"cached": "data"}`)
		client.PutInCache("movie_/movie/popular?page=1", cachedData)

		// Test
		data, err := client.FetchPageData(
			context.Background(),
			"/movie/popular",
			nil,
			1,
			"movie_",
		)

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, cachedData, data)

		// Verify no API call was made
		assert.Equal(t, 0, transport.getRequestCount())
	})

	t.Run("respects rate limiting", func(t *testing.T) {
		// Setup
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		transport := &mockHTTPTransport{
			responses: make(map[string]mockResponse),
		}

		// Set response for all requests
		tmdbResp := internal.TMDBResponse{
			Page:         1,
			Results:      []internal.TMDBMovie{},
			TotalPages:   1,
			TotalResults: 0,
		}
		respBody, _ := json.Marshal(tmdbResp)

		transport.setResponse("/movie/popular", mockResponse{
			statusCode: 200,
			body:       respBody,
		})

		httpClient := &http.Client{
			Transport: transport,
		}

		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, httpClient, mockClock)

		// Make multiple rapid requests
		start := time.Now()

		for i := 1; i <= 3; i++ {
			_, err := client.FetchPageData(
				context.Background(),
				"/movie/popular",
				nil,
				i, // Different pages to avoid cache
				"movie_",
			)
			assert.NoError(t, err)
		}

		elapsed := time.Since(start)

		// With rate limiting, 3 requests should take some time
		assert.GreaterOrEqual(t, elapsed, 500*time.Millisecond)
		assert.Equal(t, 3, transport.getRequestCount())
	})

	t.Run("handles context cancellation", func(t *testing.T) {
		// Setup
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		transport := &mockHTTPTransport{
			responses: make(map[string]mockResponse),
			delay:     100 * time.Millisecond,
		}

		httpClient := &http.Client{
			Transport: transport,
		}

		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, httpClient, mockClock)

		// Create canceled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Test
		data, err := client.FetchPageData(
			ctx,
			"/movie/popular",
			nil,
			1,
			"movie_",
		)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Contains(t, err.Error(), "rate limit")
	})

	t.Run("builds cache key with parameters", func(t *testing.T) {
		// Setup
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		transport := &mockHTTPTransport{
			responses: make(map[string]mockResponse),
		}

		httpClient := &http.Client{
			Transport: transport,
		}

		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, httpClient, mockClock)

		// Pre-populate cache with specific key
		params := url.Values{}
		params.Set("with_genres", "28,35")
		params.Set("sort_by", "popularity.desc")

		expectedKey := "movie_/discover/movie?page=2&sort_by=popularity.desc&with_genres=28%2C35"
		client.PutInCache(expectedKey, []byte(`{"cached": "with params"}`))

		// Test
		data, err := client.FetchPageData(
			context.Background(),
			"/discover/movie",
			params,
			2,
			"movie_",
		)

		// Verify cache hit with parameters
		assert.NoError(t, err)
		assert.Equal(t, []byte(`{"cached": "with params"}`), data)
		assert.Equal(t, 0, transport.getRequestCount())
	})

	t.Run("caches successful API responses", func(t *testing.T) {
		// Setup
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		transport := &mockHTTPTransport{
			responses: make(map[string]mockResponse),
		}

		tmdbResp := internal.TMDBResponse{
			Page:         1,
			Results:      []internal.TMDBMovie{fixtures.SampleTMDBMovie},
			TotalPages:   1,
			TotalResults: 1,
		}
		respBody, _ := json.Marshal(tmdbResp)

		transport.setResponse("/movie/popular", mockResponse{
			statusCode: 200,
			body:       respBody,
		})

		httpClient := &http.Client{
			Transport: transport,
		}

		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, httpClient, mockClock)

		// First request - should hit API
		data1, err := client.FetchPageData(
			context.Background(),
			"/movie/popular",
			nil,
			1,
			"movie_",
		)
		assert.NoError(t, err)
		assert.NotNil(t, data1)
		assert.Equal(t, 1, transport.getRequestCount())

		// Second request - should hit cache
		data2, err := client.FetchPageData(
			context.Background(),
			"/movie/popular",
			nil,
			1,
			"movie_",
		)
		assert.NoError(t, err)
		assert.Equal(t, data1, data2)
		assert.Equal(t, 1, transport.getRequestCount()) // No additional request
	})
}

// Note: makeGenericRequest is tested indirectly through fetchPageData tests
// since it's an unexported method. The fetchPageData tests verify that:
// - Page parameter is correctly added to requests
// - Original parameters are preserved
// - The request is properly constructed

func TestClient_makeHTTPRequest_Success(t *testing.T) {
	t.Run("successful API request", func(t *testing.T) {
		// Setup
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		transport := &mockHTTPTransport{
			responses: make(map[string]mockResponse),
		}

		expectedResp := internal.TMDBResponse{
			Page:         1,
			TotalPages:   10,
			TotalResults: 200,
		}
		respBody, _ := json.Marshal(expectedResp)

		transport.setResponse("/movie/popular", mockResponse{
			statusCode: 200,
			body:       respBody,
		})

		httpClient := &http.Client{
			Transport: transport,
		}

		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, httpClient, mockClock)

		// Test
		data, err := client.MakeHTTPRequest(
			context.Background(),
			"/movie/popular",
			nil,
		)

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, data)

		var resp internal.TMDBResponse
		_ = json.Unmarshal(data, &resp)
		assert.Equal(t, expectedResp, resp)

		// Verify request headers
		headers := transport.getLastRequestHeaders()
		assert.Equal(t, "application/json", headers.Get("Accept"))
		assert.Equal(t, "tmdb-cli/2.0", headers.Get("User-Agent"))
	})

	t.Run("constructs URL with parameters correctly", func(t *testing.T) {
		// Setup
		transport := &mockHTTPTransport{
			responses: make(map[string]mockResponse),
		}

		transport.setResponse("/discover/movie", mockResponse{
			statusCode: 200,
			body:       []byte(`{}`),
		})

		httpClient := &http.Client{
			Transport: transport,
		}

		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, httpClient, nil)

		// Test with complex parameters
		params := url.Values{}
		params.Set("language", "en-US")
		params.Set("sort_by", "popularity.desc")
		params.Add("with_genres", "28")
		params.Add("with_genres", "35")
		params.Set("primary_release_year", "2023")

		_, err := client.MakeHTTPRequest(
			context.Background(),
			"/discover/movie",
			params,
		)

		// Verify
		assert.NoError(t, err)

		lastURL := transport.getLastRequestURL()
		assert.Contains(t, lastURL, "api_key=")
		assert.Contains(t, lastURL, "language=en-US")
		assert.Contains(t, lastURL, "sort_by=popularity.desc")
		assert.Contains(t, lastURL, "with_genres=28")
		assert.Contains(t, lastURL, "with_genres=35")
		assert.Contains(t, lastURL, "primary_release_year=2023")
	})
}

func TestClient_makeHTTPRequest_ErrorHandling(t *testing.T) {
	t.Run("handles 401 unauthorized error", func(t *testing.T) {
		// Setup
		transport := &mockHTTPTransport{
			responses: make(map[string]mockResponse),
		}

		transport.setResponse("/movie/popular", mockResponse{
			statusCode: 401,
			body:       []byte(`{"status_message": "Invalid API key"}`),
		})

		httpClient := &http.Client{
			Transport: transport,
		}

		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, httpClient, nil)

		// Test
		data, err := client.MakeHTTPRequest(
			context.Background(),
			"/movie/popular",
			nil,
		)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Contains(t, err.Error(), "invalid TMDB API key")
		assert.Contains(t, err.Error(), "https://www.themoviedb.org/settings/api")
	})

	t.Run("handles 404 not found error", func(t *testing.T) {
		// Setup
		transport := &mockHTTPTransport{
			responses: make(map[string]mockResponse),
		}

		transport.setResponse("/invalid/endpoint", mockResponse{
			statusCode: 404,
			body: []byte(
				`{"status_message": "The resource you requested could not be found."}`,
			),
		})

		httpClient := &http.Client{
			Transport: transport,
		}

		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, httpClient, nil)

		// Test
		data, err := client.MakeHTTPRequest(
			context.Background(),
			"/invalid/endpoint",
			nil,
		)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Contains(t, err.Error(), "resource not found")
	})

	t.Run("handles 429 rate limit error", func(t *testing.T) {
		// Setup
		transport := &mockHTTPTransport{
			responses: make(map[string]mockResponse),
		}

		transport.setResponse("/movie/popular", mockResponse{
			statusCode: 429,
			body:       []byte(`{"status_message": "Rate limit exceeded"}`),
		})

		httpClient := &http.Client{
			Transport: transport,
		}

		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, httpClient, nil)

		// Test
		data, err := client.MakeHTTPRequest(
			context.Background(),
			"/movie/popular",
			nil,
		)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Contains(t, err.Error(), "rate limit exceeded")
		assert.Contains(t, err.Error(), "please wait and try again")
	})

	t.Run("handles generic API errors", func(t *testing.T) {
		// Setup
		transport := &mockHTTPTransport{
			responses: make(map[string]mockResponse),
		}

		transport.setResponse("/movie/popular", mockResponse{
			statusCode: 500,
			body:       []byte(`{"status_message": "Internal server error"}`),
		})

		httpClient := &http.Client{
			Transport: transport,
		}

		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, httpClient, nil)

		// Test
		data, err := client.MakeHTTPRequest(
			context.Background(),
			"/movie/popular",
			nil,
		)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Contains(t, err.Error(), "API error (status 500)")
		assert.Contains(t, err.Error(), "Internal server error")
	})

	t.Run("handles invalid response body", func(t *testing.T) {
		// Setup
		transport := &mockHTTPTransport{
			responses: make(map[string]mockResponse),
			bodyError: errors.New("read error"),
		}

		httpClient := &http.Client{
			Transport: transport,
		}

		config := helpers.MockValidConfig()
		client := internal.NewTestClient(config, httpClient, nil)

		// Test
		data, err := client.MakeHTTPRequest(
			context.Background(),
			"/movie/popular",
			nil,
		)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Contains(t, err.Error(), "read response")
		assert.Contains(t, err.Error(), "read error")
	})
}

func TestClient_makeHTTPRequest_Retry(t *testing.T) {
	t.Run("retries on network errors", func(t *testing.T) {
		// Setup
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		transport := &mockHTTPTransport{
			responses:    make(map[string]mockResponse),
			failCount:    2, // Fail first 2 attempts
			networkError: errors.New("connection refused"),
		}

		transport.setResponse("/movie/popular", mockResponse{
			statusCode: 200,
			body:       []byte(`{"success": true}`),
		})

		httpClient := &http.Client{
			Transport: transport,
		}

		config := helpers.NewConfigBuilder().
			WithAPIKey("test-key").
			WithMaxRetries(3).
			Build()

		client := internal.NewTestClient(config, httpClient, mockClock)

		// Test
		data, err := client.MakeHTTPRequest(
			context.Background(),
			"/movie/popular",
			nil,
		)

		// Verify - should succeed after retries
		assert.NoError(t, err)
		assert.NotNil(t, data)
		assert.Equal(t, 3, transport.getRequestCount()) // 2 failures + 1 success

		// Verify retry delays were applied
		assert.Len(t, mockClock.SleepCalls, 2)
		assert.Equal(t, 1*time.Second, mockClock.SleepCalls[0])
		assert.Equal(t, 2*time.Second, mockClock.SleepCalls[1])
	})

	t.Run("fails after max retries exceeded", func(t *testing.T) {
		// Setup
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		transport := &mockHTTPTransport{
			responses:    make(map[string]mockResponse),
			networkError: errors.New("network unreachable"),
			failCount:    999, // Always fail
		}

		httpClient := &http.Client{
			Transport: transport,
		}

		config := helpers.NewConfigBuilder().
			WithAPIKey("test-key").
			WithMaxRetries(2).
			Build()

		client := internal.NewTestClient(config, httpClient, mockClock)

		// Test
		data, err := client.MakeHTTPRequest(
			context.Background(),
			"/movie/popular",
			nil,
		)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Contains(t, err.Error(), "request failed after 3 attempts")
		assert.Contains(t, err.Error(), "network unreachable")
		assert.Equal(t, 3, transport.getRequestCount()) // Initial + 2 retries.
	})
}

// mockHTTPTransport simulates HTTP transport for testing.
type mockHTTPTransport struct {
	responses    map[string]mockResponse
	requests     []mockRequest
	failCount    int
	currentFails int
	networkError error
	bodyError    error
	delay        time.Duration
}

type mockResponse struct {
	statusCode int
	body       []byte
	headers    http.Header
}

type mockRequest struct {
	method  string
	url     string
	headers http.Header
}

func (mt *mockHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Simulate delay
	if mt.delay > 0 {
		time.Sleep(mt.delay)
	}

	// Record request
	mt.requests = append(mt.requests, mockRequest{
		method:  req.Method,
		url:     req.URL.String(),
		headers: req.Header.Clone(),
	})

	// Simulate network errors
	if mt.networkError != nil && mt.currentFails < mt.failCount {
		mt.currentFails++
		return nil, mt.networkError
	}

	// Find matching response
	for pattern, resp := range mt.responses {
		if strings.Contains(req.URL.Path, pattern) {
			// Create response body
			var body io.ReadCloser
			if mt.bodyError != nil {
				body = &errorReader{err: mt.bodyError}
			} else {
				body = io.NopCloser(bytes.NewReader(resp.body))
			}

			return &http.Response{
				StatusCode: resp.statusCode,
				Body:       body,
				Header:     resp.headers,
			}, nil
		}
	}

	// Default error response
	return &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"error": "not found"}`))),
		Header:     make(http.Header),
	}, nil
}

func (mt *mockHTTPTransport) setResponse(pattern string, response mockResponse) {
	if response.headers == nil {
		response.headers = make(http.Header)
	}
	mt.responses[pattern] = response
}

func (mt *mockHTTPTransport) getRequestCount() int {
	return len(mt.requests)
}

func (mt *mockHTTPTransport) getLastRequestURL() string {
	if len(mt.requests) == 0 {
		return ""
	}
	return mt.requests[len(mt.requests)-1].url
}

func (mt *mockHTTPTransport) getLastRequestHeaders() http.Header {
	if len(mt.requests) == 0 {
		return nil
	}
	return mt.requests[len(mt.requests)-1].headers
}

// errorReader simulates read errors.
type errorReader struct {
	err error
}

func (er *errorReader) Read(_ []byte) (n int, err error) {
	return 0, er.err
}

func (er *errorReader) Close() error {
	return nil
}
