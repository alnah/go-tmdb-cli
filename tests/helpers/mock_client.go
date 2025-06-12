// tests/helpers/mock_client.go
package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// MockHTTPClient provides a controllable HTTP client for testing
type MockHTTPClient struct {
	// Response configuration
	responses       map[string]*http.Response
	defaultResponse *http.Response

	// Request tracking
	requests    []*http.Request
	requestURLs []string
	callCount   int
	mu          sync.RWMutex

	// Behavior configuration
	shouldTimeout bool
	shouldError   bool
	errorMessage  string
	delayDuration time.Duration
	// rateLimitCount int
	rateLimitAfter int
}

// NewMockHTTPClient creates a new mock HTTP client
func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{
		responses: make(map[string]*http.Response),
		defaultResponse: &http.Response{
			StatusCode: 404,
			Body:       io.NopCloser(strings.NewReader(`{"error": "Not Found"}`)),
			Header:     make(http.Header),
		},
	}
}

// Do implements the http.Client interface
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Track request
	m.requests = append(m.requests, req)
	m.requestURLs = append(m.requestURLs, req.URL.String())
	m.callCount++

	// Simulate delay if configured
	if m.delayDuration > 0 {
		time.Sleep(m.delayDuration)
	}

	// Check for timeout simulation
	if m.shouldTimeout {
		return nil, context.DeadlineExceeded
	}

	// Check for error simulation
	if m.shouldError {
		if m.errorMessage != "" {
			return nil, fmt.Errorf(m.errorMessage)
		}
		return nil, fmt.Errorf("mock client error")
	}

	// Rate limiting simulation
	if m.rateLimitAfter > 0 && m.callCount > m.rateLimitAfter {
		return &http.Response{
			StatusCode: 429,
			Body:       io.NopCloser(strings.NewReader(`{"error": "Rate limit exceeded"}`)),
			Header:     make(http.Header),
			Request:    req,
		}, nil
	}

	// Generate cache key from URL path and query
	key := m.generateCacheKey(req.URL)

	// Check for specific response
	if resp, exists := m.responses[key]; exists {
		// Clone response to avoid reuse issues
		return m.cloneResponse(resp, req), nil
	}

	// Return default response
	return m.cloneResponse(m.defaultResponse, req), nil
}

// SetResponse sets a mock response for a specific URL pattern
func (m *MockHTTPClient) SetResponse(urlPattern string, response *http.Response) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses[urlPattern] = response
}

// SetDefaultResponse sets the default response for unmatched URLs
func (m *MockHTTPClient) SetDefaultResponse(response *http.Response) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.defaultResponse = response
}

// SetTimeout configures timeout simulation
func (m *MockHTTPClient) SetTimeout(shouldTimeout bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldTimeout = shouldTimeout
}

// SetError configures error simulation
func (m *MockHTTPClient) SetError(shouldError bool, message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldError = shouldError
	m.errorMessage = message
}

// SetDelay configures response delay simulation
func (m *MockHTTPClient) SetDelay(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.delayDuration = duration
}

// SetRateLimit configures rate limiting simulation
func (m *MockHTTPClient) SetRateLimit(afterCalls int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rateLimitAfter = afterCalls
}

// SetStatusCode sets status code for default response
func (m *MockHTTPClient) SetStatusCode(code int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	body := fmt.Sprintf(`{"error": "HTTP %d"}`, code)
	m.defaultResponse = &http.Response{
		StatusCode: code,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

// GetRequests returns all captured requests
func (m *MockHTTPClient) GetRequests() []*http.Request {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return copy to prevent race conditions
	requests := make([]*http.Request, len(m.requests))
	copy(requests, m.requests)
	return requests
}

// GetRequestURLs returns all captured request URLs
func (m *MockHTTPClient) GetRequestURLs() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	urls := make([]string, len(m.requestURLs))
	copy(urls, m.requestURLs)
	return urls
}

// GetCallCount returns the number of requests made
func (m *MockHTTPClient) GetCallCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.callCount
}

// Reset clears all captured data and resets state
func (m *MockHTTPClient) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requests = nil
	m.requestURLs = nil
	m.callCount = 0
	m.shouldTimeout = false
	m.shouldError = false
	m.errorMessage = ""
	m.delayDuration = 0
	m.rateLimitAfter = 0
	m.responses = make(map[string]*http.Response)
}

// HasRequestWithURL checks if a request was made to a specific URL pattern
func (m *MockHTTPClient) HasRequestWithURL(pattern string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, url := range m.requestURLs {
		if strings.Contains(url, pattern) {
			return true
		}
	}
	return false
}

// GetLastRequest returns the most recent request
func (m *MockHTTPClient) GetLastRequest() *http.Request {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.requests) == 0 {
		return nil
	}
	return m.requests[len(m.requests)-1]
}

// Helper methods for response creation

// CreateJSONResponse creates an HTTP response with JSON body
func CreateJSONResponse(statusCode int, data any) *http.Response {
	jsonData, _ := json.Marshal(data)

	resp := &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(jsonData)),
		Header:     make(http.Header),
	}
	resp.Header.Set("Content-Type", "application/json")

	return resp
}

// CreateErrorResponse creates an HTTP error response
func CreateErrorResponse(statusCode int, message string) *http.Response {
	errorData := map[string]string{"error": message}
	return CreateJSONResponse(statusCode, errorData)
}

// CreateTMDBResponse creates a mock TMDB API response
func CreateTMDBResponse(movies []TMDBMovie, page, totalPages int) *http.Response {
	response := TMDBResponse{
		Page:         page,
		TotalPages:   totalPages,
		TotalResults: len(movies) * totalPages,
		Results:      movies,
	}
	return CreateJSONResponse(200, response)
}

// CreateGenresResponse creates a mock genres API response
func CreateGenresResponse(genres []Genre) *http.Response {
	response := TMDBGenresResponse{
		Genres: genres,
	}
	return CreateJSONResponse(200, response)
}

// Helper structures for TMDB responses
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

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Utility functions for common test scenarios

// SetupPopularMoviesResponse configures mock for popular movies endpoint
func (m *MockHTTPClient) SetupPopularMoviesResponse(movies []TMDBMovie) {
	response := CreateTMDBResponse(movies, 1, 1)
	m.SetResponse("/movie/popular", response)
}

// SetupSearchResponse configures mock for search endpoint
func (m *MockHTTPClient) SetupSearchResponse(query string, movies []TMDBMovie) {
	response := CreateTMDBResponse(movies, 1, 1)
	m.SetResponse("/search/movie", response)
}

// SetupGenresResponse configures mock for genres endpoint
func (m *MockHTTPClient) SetupGenresResponse() {
	genres := []Genre{
		{ID: 28, Name: "Action"},
		{ID: 35, Name: "Comedy"},
		{ID: 18, Name: "Drama"},
		{ID: 27, Name: "Horror"},
		{ID: 878, Name: "Science Fiction"},
		{ID: 53, Name: "Thriller"},
	}
	response := CreateGenresResponse(genres)
	m.SetResponse("/genre/movie/list", response)
}

// SetupErrorScenarios configures common error scenarios
func (m *MockHTTPClient) SetupUnauthorizedError() {
	response := CreateErrorResponse(401, "Invalid API key")
	m.SetDefaultResponse(response)
}

func (m *MockHTTPClient) SetupRateLimitError() {
	response := CreateErrorResponse(429, "Rate limit exceeded")
	m.SetDefaultResponse(response)
}

func (m *MockHTTPClient) SetupNotFoundError() {
	response := CreateErrorResponse(404, "Not found")
	m.SetDefaultResponse(response)
}

func (m *MockHTTPClient) SetupServerError() {
	response := CreateErrorResponse(500, "Internal server error")
	m.SetDefaultResponse(response)
}

// Private helper methods

func (m *MockHTTPClient) generateCacheKey(u *url.URL) string {
	// Generate a simple cache key from path and relevant query parameters
	key := u.Path

	// Add important query parameters
	query := u.Query()
	if apiKey := query.Get("api_key"); apiKey != "" {
		// Don't include API key in cache key for security
		query.Del("api_key")
	}

	if len(query) > 0 {
		key += "?" + query.Encode()
	}

	return key
}

func (m *MockHTTPClient) cloneResponse(original *http.Response, req *http.Request) *http.Response {
	// Read original body
	var bodyBytes []byte
	if original.Body != nil {
		bodyBytes, _ = io.ReadAll(original.Body)
		original.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	// Create new response with cloned body
	clone := &http.Response{
		Status:           original.Status,
		StatusCode:       original.StatusCode,
		Proto:            original.Proto,
		ProtoMajor:       original.ProtoMajor,
		ProtoMinor:       original.ProtoMinor,
		Header:           original.Header,
		Body:             io.NopCloser(bytes.NewReader(bodyBytes)),
		ContentLength:    original.ContentLength,
		TransferEncoding: original.TransferEncoding,
		Close:            original.Close,
		Uncompressed:     original.Uncompressed,
		Trailer:          original.Trailer,
		Request:          req,
		TLS:              original.TLS,
	}

	return clone
}

// MockRoundTripper implements http.RoundTripper for more advanced scenarios
type MockRoundTripper struct {
	client *MockHTTPClient
}

func NewMockRoundTripper() *MockRoundTripper {
	return &MockRoundTripper{
		client: NewMockHTTPClient(),
	}
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.client.Do(req)
}

func (m *MockRoundTripper) GetClient() *MockHTTPClient {
	return m.client
}

// Helper function to create a real http.Client with mock transport
func CreateMockHTTPClientWithTransport() (*http.Client, *MockHTTPClient) {
	mockClient := NewMockHTTPClient()
	transport := &MockRoundTripper{client: mockClient}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return httpClient, mockClient
}
