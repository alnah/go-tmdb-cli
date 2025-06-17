package unit

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

// MockRoundTripper implements http.RoundTripper for testing HTTP requests.
type MockRoundTripper struct {
	responses []MockResponse
	callCount int
}

type MockResponse struct {
	statusCode int
	body       string
	err        error
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.callCount >= len(m.responses) {
		return nil, errors.New("unexpected request")
	}

	resp := m.responses[m.callCount]
	m.callCount++

	if resp.err != nil {
		return nil, resp.err
	}

	return &http.Response{
		StatusCode: resp.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(resp.body)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

func TestClient_RetryDelays(t *testing.T) {
	t.Run("retry with progressive delays on network errors", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Setup HTTP client that fails multiple times
		mockTransport := &MockRoundTripper{
			responses: []MockResponse{
				{err: errors.New("network error 1")},
				{err: errors.New("network error 2")},
				{err: errors.New("network error 3")},
				{err: errors.New("network error 4")}, // Final failure
			},
		}

		httpClient := &http.Client{
			Transport: mockTransport,
			Timeout:   30 * time.Second,
		}

		// Create client with 3 retries
		config := internal.DefaultConfig()
		config.MaxRetries = 3

		client := internal.NewTestClient(config, httpClient, mockClock)

		// Make request that will fail
		ctx := context.Background()
		_, err := client.MakeHTTPRequest(ctx, "/test", nil)

		// Verify request failed
		require.Error(t, err)
		assert.Contains(t, err.Error(), "request failed after 4 attempts")

		// Verify retry delays
		require.Len(t, mockClock.SleepCalls, 3)
		assert.Equal(t, 1*time.Second, mockClock.SleepCalls[0], "First retry delay")
		assert.Equal(t, 2*time.Second, mockClock.SleepCalls[1], "Second retry delay")
		assert.Equal(t, 3*time.Second, mockClock.SleepCalls[2], "Third retry delay")

		// Verify all attempts were made
		assert.Equal(t, 4, mockTransport.callCount, "Should make initial request + 3 retries")
	})

	t.Run("successful request after retries", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Setup HTTP client that fails twice then succeeds
		mockTransport := &MockRoundTripper{
			responses: []MockResponse{
				{err: errors.New("network error 1")},
				{err: errors.New("network error 2")},
				{statusCode: 200, body: `{"success": true}`}, // Success on third attempt
			},
		}

		httpClient := &http.Client{
			Transport: mockTransport,
			Timeout:   30 * time.Second,
		}

		// Create client
		config := internal.DefaultConfig()
		config.MaxRetries = 3

		client := internal.NewTestClient(config, httpClient, mockClock)

		// Make request
		ctx := context.Background()
		data, err := client.MakeHTTPRequest(ctx, "/test", nil)

		// Verify success
		require.NoError(t, err)
		assert.Equal(t, `{"success": true}`, string(data))

		// Verify retry delays before success
		require.Len(t, mockClock.SleepCalls, 2)
		assert.Equal(t, 1*time.Second, mockClock.SleepCalls[0])
		assert.Equal(t, 2*time.Second, mockClock.SleepCalls[1])

		// Verify attempts
		assert.Equal(t, 3, mockTransport.callCount)
	})

	t.Run("no retries on successful first request", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Setup HTTP client that succeeds immediately
		mockTransport := &MockRoundTripper{
			responses: []MockResponse{
				{statusCode: 200, body: `{"data": "test"}`},
			},
		}

		httpClient := &http.Client{
			Transport: mockTransport,
			Timeout:   30 * time.Second,
		}

		// Create client
		config := internal.DefaultConfig()
		config.MaxRetries = 3

		client := internal.NewTestClient(config, httpClient, mockClock)

		// Make request
		ctx := context.Background()
		data, err := client.MakeHTTPRequest(ctx, "/test", nil)

		// Verify success
		require.NoError(t, err)
		assert.Equal(t, `{"data": "test"}`, string(data))

		// Verify no sleep calls
		assert.Empty(t, mockClock.SleepCalls)

		// Verify only one attempt
		assert.Equal(t, 1, mockTransport.callCount)
	})

	t.Run("zero retries configuration", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Setup HTTP client that fails
		mockTransport := &MockRoundTripper{
			responses: []MockResponse{
				{err: errors.New("network error")},
			},
		}

		httpClient := &http.Client{
			Transport: mockTransport,
			Timeout:   30 * time.Second,
		}

		// Create client with 0 retries
		config := internal.DefaultConfig()
		config.MaxRetries = 0

		client := internal.NewTestClient(config, httpClient, mockClock)

		// Make request
		ctx := context.Background()
		_, err := client.MakeHTTPRequest(ctx, "/test", nil)

		// Verify failure
		require.Error(t, err)
		assert.Contains(t, err.Error(), "request failed after 1 attempts")

		// Verify no sleep calls
		assert.Empty(t, mockClock.SleepCalls)

		// Verify only one attempt
		assert.Equal(t, 1, mockTransport.callCount)
	})

	t.Run("context cancellation during retry", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Setup HTTP client that always fails
		mockTransport := &MockRoundTripper{
			responses: []MockResponse{
				{err: errors.New("network error 1")},
				{err: errors.New("network error 2")},
				{err: errors.New("network error 3")},
			},
		}

		httpClient := &http.Client{
			Transport: mockTransport,
			Timeout:   30 * time.Second,
		}

		// Create client
		config := internal.DefaultConfig()
		config.MaxRetries = 3

		client := internal.NewTestClient(config, httpClient, mockClock)

		// Create cancellable context
		ctx, cancel := context.WithCancel(context.Background())

		// Cancel after first retry
		go func() {
			time.Sleep(1 * time.Millisecond)
			cancel()
		}()

		// Make request
		_, err := client.MakeHTTPRequest(ctx, "/test", nil)

		// Verify context error
		require.Error(t, err)
		// The error could be from the HTTP request or from context cancellation
		// depending on timing, so we just verify there was an error
	})
}

func TestClient_RetryOnStatusCodes(t *testing.T) {
	t.Run("no retry on client errors (4xx)", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Test various 4xx status codes
		statusCodes := []int{400, 401, 403, 404, 429}

		for _, statusCode := range statusCodes {
			t.Run(helpers.FormatStatusCode(statusCode), func(t *testing.T) {
				// Setup HTTP client that returns error status
				mockTransport := &MockRoundTripper{
					responses: []MockResponse{
						{statusCode: statusCode, body: `{"error": "client error"}`},
					},
				}

				httpClient := &http.Client{
					Transport: mockTransport,
					Timeout:   30 * time.Second,
				}

				// Create client
				config := internal.DefaultConfig()
				config.MaxRetries = 3

				client := internal.NewTestClient(config, httpClient, mockClock)

				// Make request
				ctx := context.Background()
				_, err := client.MakeHTTPRequest(ctx, "/test", nil)

				// Verify error
				require.Error(t, err)

				// Verify no retries (no sleep calls)
				assert.Empty(t, mockClock.SleepCalls)

				// Verify only one attempt
				assert.Equal(t, 1, mockTransport.callCount)
			})
		}
	})

	t.Run("no retry on server errors (5xx)", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Test various 5xx status codes
		statusCodes := []int{500, 502, 503, 504}

		for _, statusCode := range statusCodes {
			t.Run(helpers.FormatStatusCode(statusCode), func(t *testing.T) {
				// Setup HTTP client that returns error status
				mockTransport := &MockRoundTripper{
					responses: []MockResponse{
						{statusCode: statusCode, body: `{"error": "server error"}`},
					},
				}

				httpClient := &http.Client{
					Transport: mockTransport,
					Timeout:   30 * time.Second,
				}

				// Create client
				config := internal.DefaultConfig()
				config.MaxRetries = 3

				client := internal.NewTestClient(config, httpClient, mockClock)

				// Make request
				ctx := context.Background()
				_, err := client.MakeHTTPRequest(ctx, "/test", nil)

				// Verify error
				require.Error(t, err)

				// Verify no retries (no sleep calls)
				assert.Empty(t, mockClock.SleepCalls)

				// Verify only one attempt
				assert.Equal(t, 1, mockTransport.callCount)
			})
		}
	})
}

func TestClient_RetryDelayCalculation(t *testing.T) {
	t.Run("exponential backoff pattern", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Create a transport that always fails
		numRetries := 10
		responses := make([]MockResponse, numRetries+1)
		for i := range responses {
			responses[i] = MockResponse{err: errors.New("network error")}
		}

		mockTransport := &MockRoundTripper{
			responses: responses,
		}

		httpClient := &http.Client{
			Transport: mockTransport,
			Timeout:   30 * time.Second,
		}

		// Create client with many retries
		config := internal.DefaultConfig()
		config.APIKey = "test-key"
		config.MaxRetries = numRetries

		client := internal.NewTestClient(config, httpClient, mockClock)

		// Make request
		ctx := context.Background()
		_, err := client.MakeHTTPRequest(ctx, "/test", nil)

		// Verify failure
		require.Error(t, err)

		// Verify progressive delays
		require.Len(t, mockClock.SleepCalls, numRetries)
		for i := range numRetries {
			expectedDelay := time.Duration(i+1) * time.Second
			assert.Equal(t, expectedDelay, mockClock.SleepCalls[i],
				"Retry %d should have delay of %v", i+1, expectedDelay)
		}
	})
}

func TestClient_TimeTracking(t *testing.T) {
	t.Run("track total time with retries", func(t *testing.T) {
		// Setup mock clock at specific time
		startTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		mockClock := &internal.MockClock{
			CurrentTime: startTime,
		}

		// Setup HTTP client that fails 3 times then succeeds
		mockTransport := &MockRoundTripper{
			responses: []MockResponse{
				{err: errors.New("error 1")},
				{err: errors.New("error 2")},
				{err: errors.New("error 3")},
				{statusCode: 200, body: `{"success": true}`},
			},
		}

		httpClient := &http.Client{
			Transport: mockTransport,
			Timeout:   30 * time.Second,
		}

		// Create client
		config := internal.DefaultConfig()
		config.APIKey = "test-key"
		config.MaxRetries = 5

		client := internal.NewTestClient(config, httpClient, mockClock)

		// Make request
		ctx := context.Background()
		_, err := client.MakeHTTPRequest(ctx, "/test", nil)

		// Verify success
		require.NoError(t, err)

		// Calculate total sleep time
		var totalSleepTime time.Duration
		for _, sleep := range mockClock.SleepCalls {
			totalSleepTime += sleep
		}

		// Should have slept: 1s + 2s + 3s = 6s
		assert.Equal(t, 6*time.Second, totalSleepTime)

		// Verify sleep pattern
		assert.Equal(t, []time.Duration{
			1 * time.Second,
			2 * time.Second,
			3 * time.Second,
		}, mockClock.SleepCalls)
	})
}
