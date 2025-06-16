// Package helpers provides testing utilities for mock objects.
package helpers

import (
	"errors"
	"io"
	"sync"

	"github.com/alnah/tmdb-cli/internal"
)

// MockTMDBMovie creates a mock TMDB movie response for testing.
func MockTMDBMovie(id int, title string) internal.TMDBMovie {
	return internal.TMDBMovie{
		ID:            id,
		Title:         title,
		OriginalTitle: title,
		Overview:      "Mock movie overview",
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
		Overview:         "Mock TV show overview",
		FirstAirDate:     "2023-01-01",
		VoteAverage:      8.0,
		VoteCount:        1500,
		GenreIDs:         []int{18, 53}, // Drama, Thriller
		Popularity:       60.0,
		Adult:            false,
		OriginalLanguage: "en",
	}
}

// MockWriter simulates an io.Writer that can be configured to fail.
// Useful for testing error handling in formatters.
type MockWriter struct {
	mu          sync.Mutex
	buffer      []byte
	shouldError bool
	errorMsg    string
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

// Reset clears the internal buffer and resets error state.
func (mw *MockWriter) Reset() {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	mw.buffer = mw.buffer[:0]
	mw.shouldError = false
}

// Len returns the number of bytes written to the buffer.
func (mw *MockWriter) Len() int {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	return len(mw.buffer)
}

// MockReadCloser simulates an io.ReadCloser for testing.
type MockReadCloser struct {
	data      []byte
	position  int
	closed    bool
	readErr   error
	closeErr  error
	closeFunc func() error
}

// NewMockReadCloser creates a new MockReadCloser with the given data.
func NewMockReadCloser(data []byte) *MockReadCloser {
	return &MockReadCloser{
		data:     data,
		position: 0,
	}
}

// Read implements io.Reader interface.
func (mrc *MockReadCloser) Read(buffer []byte) (n int, err error) {
	if mrc.readErr != nil {
		return 0, mrc.readErr
	}

	if mrc.closed {
		return 0, errors.New("reader is closed")
	}

	if mrc.position >= len(mrc.data) {
		return 0, io.EOF
	}

	n = copy(buffer, mrc.data[mrc.position:])
	mrc.position += n
	return n, nil
}

// Close implements io.Closer interface.
func (mrc *MockReadCloser) Close() error {
	if mrc.closeFunc != nil {
		return mrc.closeFunc()
	}

	if mrc.closeErr != nil {
		return mrc.closeErr
	}

	mrc.closed = true
	return nil
}

// SetReadError configures the Read method to return an error.
func (mrc *MockReadCloser) SetReadError(err error) {
	mrc.readErr = err
}

// SetCloseError configures the Close method to return an error.
func (mrc *MockReadCloser) SetCloseError(err error) {
	mrc.closeErr = err
}

// SetCloseFunc sets a custom close function.
func (mrc *MockReadCloser) SetCloseFunc(fn func() error) {
	mrc.closeFunc = fn
}

// IsClosed returns whether the reader has been closed.
func (mrc *MockReadCloser) IsClosed() bool {
	return mrc.closed
}

// Reset resets the reader position and clears errors.
func (mrc *MockReadCloser) Reset() {
	mrc.position = 0
	mrc.closed = false
	mrc.readErr = nil
	mrc.closeErr = nil
	mrc.closeFunc = nil
}
