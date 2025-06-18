// Package helpers provides testing utilities for the TMDB CLI test suite.
package helpers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/alnah/tmdb-cli/internal"
)

// MockWriter simulates an io.Writer that can be configured to fail.
// Useful for testing error handling in formatters.
type MockWriter struct {
	mu          sync.Mutex
	buffer      []byte
	shouldError bool
	errorMsg    string
	writeCount  int
}

// NewMockWriter creates a new MockWriter instance.
func NewMockWriter() *MockWriter {
	return &MockWriter{
		buffer:   make([]byte, 0),
		errorMsg: "mock write error",
	}
}

// Write implements io.Writer interface with controllable error behavior.
func (mw *MockWriter) Write(data []byte) (n int, err error) {
	mw.mu.Lock()
	defer mw.mu.Unlock()

	mw.writeCount++

	if mw.shouldError {
		return 0, errors.New(mw.errorMsg)
	}

	mw.buffer = append(mw.buffer, data...)
	return len(data), nil
}

// SetShouldError configures whether Write should return an error.
func (mw *MockWriter) SetShouldError(shouldError bool) {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	mw.shouldError = shouldError
}

// SetErrorMessage sets the error message returned when shouldError is true.
func (mw *MockWriter) SetErrorMessage(msg string) {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	mw.errorMsg = msg
}

// GetWritten returns the data that was written to the mock writer.
func (mw *MockWriter) GetWritten() []byte {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	result := make([]byte, len(mw.buffer))
	copy(result, mw.buffer)
	return result
}

// GetWrittenString returns the written data as a string.
func (mw *MockWriter) GetWrittenString() string {
	return string(mw.GetWritten())
}

// GetWriteCount returns the number of times Write was called.
func (mw *MockWriter) GetWriteCount() int {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	return mw.writeCount
}

// Reset clears the internal buffer and resets error state.
func (mw *MockWriter) Reset() {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	mw.buffer = mw.buffer[:0]
	mw.shouldError = false
	mw.writeCount = 0
}

// Len returns the number of bytes written to the buffer.
func (mw *MockWriter) Len() int {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	return len(mw.buffer)
}

// MockReadCloser simulates an io.ReadCloser for testing.
type MockReadCloser struct {
	data       []byte
	pos        int
	closeCount int
	shouldErr  bool
	errorMsg   string
}

// NewMockReadCloser creates a new MockReadCloser with the given data.
func NewMockReadCloser(data []byte) *MockReadCloser {
	return &MockReadCloser{
		data:     data,
		errorMsg: "mock read error",
	}
}

// Read implements io.Reader interface.
func (mrc *MockReadCloser) Read(data []byte) (n int, err error) {
	if mrc.shouldErr {
		return 0, errors.New(mrc.errorMsg)
	}

	if mrc.pos >= len(mrc.data) {
		return 0, io.EOF
	}

	n = copy(data, mrc.data[mrc.pos:])
	mrc.pos += n
	return n, nil
}

// Close implements io.Closer interface.
func (mrc *MockReadCloser) Close() error {
	mrc.closeCount++
	if mrc.shouldErr {
		return errors.New("mock close error")
	}
	return nil
}

// SetShouldError configures whether Read/Close should return errors.
func (mrc *MockReadCloser) SetShouldError(shouldErr bool) {
	mrc.shouldErr = shouldErr
}

// SetErrorMessage sets the error message for Read/Close operations.
func (mrc *MockReadCloser) SetErrorMessage(msg string) {
	mrc.errorMsg = msg
}

// GetCloseCount returns the number of times Close was called.
func (mrc *MockReadCloser) GetCloseCount() int {
	return mrc.closeCount
}

// Reset resets the reader position and error state.
func (mrc *MockReadCloser) Reset() {
	mrc.pos = 0
	mrc.shouldErr = false
	mrc.closeCount = 0
}

// MockTMDBMovie creates a mock TMDB movie response for testing.
func MockTMDBMovie(id int, title string) internal.TMDBMovie {
	return internal.TMDBMovie{
		ID:            id,
		Title:         title,
		OriginalTitle: title,
		Overview:      "Mock movie overview for testing purposes",
		ReleaseDate:   "2023-01-01",
		VoteAverage:   7.5,
		VoteCount:     1000,
		GenreIDs:      []int{28, 35}, // Action, Comedy
		Popularity:    50.0,
		Adult:         false,
		Video:         false,
	}
}

// MockTMDBTVShow creates a mock TMDB TV show response for testing.
func MockTMDBTVShow(id int, name string) internal.TMDBTVShow {
	return internal.TMDBTVShow{
		ID:               id,
		Name:             name,
		OriginalName:     name,
		Overview:         "Mock TV show overview for testing purposes",
		FirstAirDate:     "2023-01-01",
		VoteAverage:      8.0,
		VoteCount:        1500,
		GenreIDs:         []int{18, 53}, // Drama, Thriller
		Popularity:       60.0,
		Adult:            false,
		OriginalLanguage: "en",
	}
}

// MockTMDBResponse creates a mock TMDB API response for testing.
func MockTMDBResponse(movies []internal.TMDBMovie) internal.TMDBResponse {
	return internal.TMDBResponse{
		Page:         1,
		Results:      movies,
		TotalPages:   1,
		TotalResults: len(movies),
	}
}

// MockHTTPClient simulates an HTTP client for testing API interactions.
type MockHTTPClient struct {
	mu           sync.Mutex
	responses    map[string]*MockHTTPResponse
	requestCount int
	lastURL      string
	lastMethod   string
	shouldError  bool
	errorMsg     string
}

// MockHTTPResponse represents a mock HTTP response.
type MockHTTPResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
	Error      error
}

// NewMockHTTPClient creates a new mock HTTP client.
func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{
		responses: make(map[string]*MockHTTPResponse),
		errorMsg:  "mock HTTP error",
	}
}

// SetResponse configures a mock response for a specific URL pattern.
func (mhc *MockHTTPClient) SetResponse(urlPattern string, response *MockHTTPResponse) {
	mhc.mu.Lock()
	defer mhc.mu.Unlock()
	mhc.responses[urlPattern] = response
}

// SetError configures the client to return an error for all requests.
func (mhc *MockHTTPClient) SetError(shouldError bool, errorMsg string) {
	mhc.mu.Lock()
	defer mhc.mu.Unlock()
	mhc.shouldError = shouldError
	if errorMsg != "" {
		mhc.errorMsg = errorMsg
	}
}

// GetRequestCount returns the number of requests made.
func (mhc *MockHTTPClient) GetRequestCount() int {
	mhc.mu.Lock()
	defer mhc.mu.Unlock()
	return mhc.requestCount
}

// GetLastRequest returns information about the last request made.
func (mhc *MockHTTPClient) GetLastRequest() (method, url string) {
	mhc.mu.Lock()
	defer mhc.mu.Unlock()
	return mhc.lastMethod, mhc.lastURL
}

// Reset clears all mock responses and resets counters.
func (mhc *MockHTTPClient) Reset() {
	mhc.mu.Lock()
	defer mhc.mu.Unlock()
	mhc.responses = make(map[string]*MockHTTPResponse)
	mhc.requestCount = 0
	mhc.lastURL = ""
	mhc.lastMethod = ""
	mhc.shouldError = false
}

// MockConfig creates a mock configuration for testing.
func MockConfig() internal.Config {
	return NewConfigBuilder().
		WithAPIKey("mock-api-key").
		WithBaseURL("https://mock.api.com").
		WithLogLevel("info").
		WithFormat("json").
		Build()
}

// MockValidConfig creates a valid mock configuration for testing.
func MockValidConfig() internal.Config {
	return NewConfigBuilder().
		WithAPIKey("valid-api-key-12345").
		WithBaseURL("https://api.themoviedb.org/3").
		WithLogLevel("info").
		WithFormat("table").
		Build()
}

// MockInvalidConfig creates an invalid mock configuration for testing validation.
func MockInvalidConfig() internal.Config {
	config := internal.Config{}
	// Leave fields empty or with invalid values
	config.APIKey = ""
	config.BaseURL = ""
	config.LogLevel = "invalid"
	config.Format = "invalid"
	return config
}

// MockSearchOptions creates mock search options for testing.
func MockSearchOptions() internal.SearchOptions {
	return NewSearchOptionsBuilder().
		WithMaxItems(10).
		WithMinRating(7.0).
		WithMaxRating(9.0).
		Build()
}

// MockFormatOptions creates mock format options for testing.
func MockFormatOptions(format string) internal.FormatOptions {
	return internal.FormatOptions{
		Format:   format,
		MaxWidth: 120,
		NoHeader: false,
	}
}

// MockLogger simulates a logger for testing.
type MockLogger struct {
	mu       sync.Mutex
	messages []LogMessage
}

// LogMessage represents a logged message.
type LogMessage struct {
	Level   string
	Message string
	Fields  map[string]any
}

// NewMockLogger creates a new mock logger.
func NewMockLogger() *MockLogger {
	return &MockLogger{
		messages: make([]LogMessage, 0),
	}
}

// Debug logs a debug message.
func (ml *MockLogger) Debug(msg string, fields ...any) {
	ml.log("DEBUG", msg, fields...)
}

// Info logs an info message.
func (ml *MockLogger) Info(msg string, fields ...any) {
	ml.log("INFO", msg, fields...)
}

// Warn logs a warning message.
func (ml *MockLogger) Warn(msg string, fields ...any) {
	ml.log("WARN", msg, fields...)
}

// Error logs an error message.
func (ml *MockLogger) Error(msg string, fields ...any) {
	ml.log("ERROR", msg, fields...)
}

// log is a helper method to log messages.
func (ml *MockLogger) log(level, msg string, fields ...any) {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	fieldMap := make(map[string]any)
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			if key, ok := fields[i].(string); ok {
				fieldMap[key] = fields[i+1]
			}
		}
	}

	ml.messages = append(ml.messages, LogMessage{
		Level:   level,
		Message: msg,
		Fields:  fieldMap,
	})
}

// GetMessages returns all logged messages.
func (ml *MockLogger) GetMessages() []LogMessage {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	messages := make([]LogMessage, len(ml.messages))
	copy(messages, ml.messages)
	return messages
}

// GetMessagesByLevel returns messages filtered by level.
func (ml *MockLogger) GetMessagesByLevel(level string) []LogMessage {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	var filtered []LogMessage
	for _, msg := range ml.messages {
		if msg.Level == level {
			filtered = append(filtered, msg)
		}
	}
	return filtered
}

// Reset clears all logged messages.
func (ml *MockLogger) Reset() {
	ml.mu.Lock()
	defer ml.mu.Unlock()
	ml.messages = ml.messages[:0]
}

// HasMessage checks if a message with the given text was logged.
func (ml *MockLogger) HasMessage(text string) bool {
	ml.mu.Lock()
	defer ml.mu.Unlock()

	for _, msg := range ml.messages {
		if msg.Message == text {
			return true
		}
	}
	return false
}

// MockCache simulates a cache for testing.
type MockCache struct {
	mu     sync.Mutex
	data   map[string]any
	hits   int
	misses int
}

// NewMockCache creates a new mock cache.
func NewMockCache() *MockCache {
	return &MockCache{
		data: make(map[string]any),
	}
}

// Get retrieves a value from the cache.
func (mc *MockCache) Get(key string) (any, bool) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	value, exists := mc.data[key]
	if exists {
		mc.hits++
	} else {
		mc.misses++
	}
	return value, exists
}

// Set stores a value in the cache.
func (mc *MockCache) Set(key string, value any) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.data[key] = value
}

// Delete removes a value from the cache.
func (mc *MockCache) Delete(key string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	delete(mc.data, key)
}

// Clear removes all values from the cache.
func (mc *MockCache) Clear() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.data = make(map[string]any)
	mc.hits = 0
	mc.misses = 0
}

// GetStats returns cache statistics.
func (mc *MockCache) GetStats() (hits, misses int) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	return mc.hits, mc.misses
}

// Size returns the number of items in the cache.
func (mc *MockCache) Size() int {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	return len(mc.data)
}

// SimulateWriteError creates a function that simulates write errors for testing.
func SimulateWriteError(writer *MockWriter) func() {
	return func() {
		writer.SetShouldError(true)
	}
}

// CreateMockMovieCollection creates a collection of mock movies for testing.
func CreateMockMovieCollection(count int) []internal.Movie {
	movies := make([]internal.Movie, count)
	for i := range count {
		movies[i] = NewMovieBuilder().
			WithID(i + 1).
			WithTitle(fmt.Sprintf("Mock Movie %d", i+1)).
			Build()
	}
	return movies
}

// CreateMockTVShowCollection creates a collection of mock TV shows for testing.
func CreateMockTVShowCollection(count int) []internal.TVShow {
	shows := make([]internal.TVShow, count)
	for i := range count {
		shows[i] = NewTVShowBuilder().
			WithID(i + 1).
			WithName(fmt.Sprintf("Mock TV Show %d", i+1)).
			Build()
	}
	return shows
}

// MockRoundTripper is for testing purpose.
type MockRoundTripper struct {
	Responses []MockResponse
	CallCount int
}

// MockResponse is for testing purpose.
type MockResponse struct {
	StatusCode int
	Body       string
	Err        error
}

// RoundTrip is for testing purpose.
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.CallCount >= len(m.Responses) {
		return nil, errors.New("unexpected request")
	}

	resp := m.Responses[m.CallCount]
	m.CallCount++

	if resp.Err != nil {
		return nil, resp.Err
	}

	return &http.Response{
		StatusCode: resp.StatusCode,
		Body:       io.NopCloser(bytes.NewBufferString(resp.Body)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

// MockTransport is a mock HTTP transport for testing.
type MockTransport struct {
	mu           sync.RWMutex
	Responses    map[string]MockHTTPResponse
	RequestCount int
	LastURL      string
}

// NewMockTransport creates a new mock transport.
func NewMockTransport() *MockTransport {
	return &MockTransport{
		Responses: make(map[string]MockHTTPResponse),
	}
}

// SetResponse sets a response for a URL pattern.
func (t *MockTransport) SetResponse(pattern string, response MockHTTPResponse) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Responses[pattern] = response
}

// RoundTrip implements http.RoundTripper.
func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.mu.Lock()
	t.RequestCount++
	t.LastURL = req.URL.String()
	t.mu.Unlock()

	t.mu.RLock()
	defer t.mu.RUnlock()

	// Check for exact match first
	if resp, ok := t.Responses[req.URL.Path]; ok {
		return &http.Response{
			StatusCode: resp.StatusCode,
			Body:       io.NopCloser(bytes.NewReader(resp.Body)),
			Header:     make(http.Header),
		}, nil
	}

	// Check for pattern matches (e.g., "page=1", "page=2")
	for pattern, resp := range t.Responses {
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
