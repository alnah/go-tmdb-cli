package helpers

import (
	"context"
	"fmt"
	"sync"

	"github.com/alnah/tmdb-cli/internal"
)

// MockTMDBClient is a mock implementation of TMDBClient for testing.
type MockTMDBClient struct {
	mu sync.Mutex

	// Movie operation mocks
	GetPopularMoviesFunc    func(ctx context.Context, maxItems int) ([]internal.Movie, error)
	GetTopRatedMoviesFunc   func(ctx context.Context, maxItems int) ([]internal.Movie, error)
	GetNowPlayingMoviesFunc func(ctx context.Context, maxItems int) ([]internal.Movie, error)
	GetUpcomingMoviesFunc   func(ctx context.Context, maxItems int) ([]internal.Movie, error)
	SearchMoviesFunc        func(ctx context.Context, query string, maxItems int) ([]internal.Movie, error)
	DiscoverMoviesFunc      func(ctx context.Context, opts internal.SearchOptions) ([]internal.Movie, error)

	// TV show operation mocks
	GetPopularTVShowsFunc  func(ctx context.Context, maxItems int) ([]internal.TVShow, error)
	GetTopRatedTVShowsFunc func(ctx context.Context, maxItems int) ([]internal.TVShow, error)
	GetOnTheAirTVShowsFunc func(ctx context.Context, maxItems int) ([]internal.TVShow, error)
	SearchTVShowsFunc      func(ctx context.Context, query string, maxItems int) ([]internal.TVShow, error)
	DiscoverTVShowsFunc    func(ctx context.Context, opts internal.SearchOptions) ([]internal.TVShow, error)

	// Call tracking
	Calls []MethodCall
}

// MethodCall represents a method call made to the mock.
type MethodCall struct {
	Method string
	Args   []any
}

// NewMockTMDBClient creates a new mock TMDB client with default no-op implementations.
func NewMockTMDBClient() *MockTMDBClient {
	return &MockTMDBClient{
		// Default implementations return empty slices
		GetPopularMoviesFunc: func(_ context.Context, _ int) ([]internal.Movie, error) {
			return []internal.Movie{}, nil
		},
		GetTopRatedMoviesFunc: func(_ context.Context, _ int) ([]internal.Movie, error) {
			return []internal.Movie{}, nil
		},
		GetNowPlayingMoviesFunc: func(_ context.Context, _ int) ([]internal.Movie, error) {
			return []internal.Movie{}, nil
		},
		GetUpcomingMoviesFunc: func(_ context.Context, _ int) ([]internal.Movie, error) {
			return []internal.Movie{}, nil
		},
		SearchMoviesFunc: func(_ context.Context, _ string, _ int) ([]internal.Movie, error) {
			return []internal.Movie{}, nil
		},
		DiscoverMoviesFunc: func(_ context.Context, _ internal.SearchOptions) ([]internal.Movie, error) {
			return []internal.Movie{}, nil
		},
		GetPopularTVShowsFunc: func(_ context.Context, _ int) ([]internal.TVShow, error) {
			return []internal.TVShow{}, nil
		},
		GetTopRatedTVShowsFunc: func(_ context.Context, _ int) ([]internal.TVShow, error) {
			return []internal.TVShow{}, nil
		},
		GetOnTheAirTVShowsFunc: func(_ context.Context, _ int) ([]internal.TVShow, error) {
			return []internal.TVShow{}, nil
		},
		SearchTVShowsFunc: func(_ context.Context, _ string, _ int) ([]internal.TVShow, error) {
			return []internal.TVShow{}, nil
		},
		DiscoverTVShowsFunc: func(_ context.Context, _ internal.SearchOptions) ([]internal.TVShow, error) {
			return []internal.TVShow{}, nil
		},
	}
}

// recordCall records a method call for verification.
func (m *MockTMDBClient) recordCall(method string, args ...any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = append(m.Calls, MethodCall{Method: method, Args: args})
}

// Movie operations

// GetPopularMovies calls the mock function.
func (m *MockTMDBClient) GetPopularMovies(
	ctx context.Context,
	maxItems int,
) ([]internal.Movie, error) {
	m.recordCall("GetPopularMovies", ctx, maxItems)
	if m.GetPopularMoviesFunc == nil {
		return nil, fmt.Errorf("GetPopularMoviesFunc not implemented")
	}
	return m.GetPopularMoviesFunc(ctx, maxItems)
}

// GetTopRatedMovies calls the mock function.
func (m *MockTMDBClient) GetTopRatedMovies(
	ctx context.Context,
	maxItems int,
) ([]internal.Movie, error) {
	m.recordCall("GetTopRatedMovies", ctx, maxItems)
	if m.GetTopRatedMoviesFunc == nil {
		return nil, fmt.Errorf("GetTopRatedMoviesFunc not implemented")
	}
	return m.GetTopRatedMoviesFunc(ctx, maxItems)
}

// GetNowPlayingMovies calls the mock function.
func (m *MockTMDBClient) GetNowPlayingMovies(
	ctx context.Context,
	maxItems int,
) ([]internal.Movie, error) {
	m.recordCall("GetNowPlayingMovies", ctx, maxItems)
	if m.GetNowPlayingMoviesFunc == nil {
		return nil, fmt.Errorf("GetNowPlayingMoviesFunc not implemented")
	}
	return m.GetNowPlayingMoviesFunc(ctx, maxItems)
}

// GetUpcomingMovies calls the mock function.
func (m *MockTMDBClient) GetUpcomingMovies(
	ctx context.Context,
	maxItems int,
) ([]internal.Movie, error) {
	m.recordCall("GetUpcomingMovies", ctx, maxItems)
	if m.GetUpcomingMoviesFunc == nil {
		return nil, fmt.Errorf("GetUpcomingMoviesFunc not implemented")
	}
	return m.GetUpcomingMoviesFunc(ctx, maxItems)
}

// SearchMovies calls the mock function.
func (m *MockTMDBClient) SearchMovies(
	ctx context.Context,
	query string,
	maxItems int,
) ([]internal.Movie, error) {
	m.recordCall("SearchMovies", ctx, query, maxItems)
	if m.SearchMoviesFunc == nil {
		return nil, fmt.Errorf("SearchMoviesFunc not implemented")
	}
	return m.SearchMoviesFunc(ctx, query, maxItems)
}

// DiscoverMovies calls the mock function.
func (m *MockTMDBClient) DiscoverMovies(
	ctx context.Context,
	opts internal.SearchOptions,
) ([]internal.Movie, error) {
	m.recordCall("DiscoverMovies", ctx, opts)
	if m.DiscoverMoviesFunc == nil {
		return nil, fmt.Errorf("DiscoverMoviesFunc not implemented")
	}
	return m.DiscoverMoviesFunc(ctx, opts)
}

// TV show operations

// GetPopularTVShows calls the mock function.
func (m *MockTMDBClient) GetPopularTVShows(
	ctx context.Context,
	maxItems int,
) ([]internal.TVShow, error) {
	m.recordCall("GetPopularTVShows", ctx, maxItems)
	if m.GetPopularTVShowsFunc == nil {
		return nil, fmt.Errorf("GetPopularTVShowsFunc not implemented")
	}
	return m.GetPopularTVShowsFunc(ctx, maxItems)
}

// GetTopRatedTVShows calls the mock function.
func (m *MockTMDBClient) GetTopRatedTVShows(
	ctx context.Context,
	maxItems int,
) ([]internal.TVShow, error) {
	m.recordCall("GetTopRatedTVShows", ctx, maxItems)
	if m.GetTopRatedTVShowsFunc == nil {
		return nil, fmt.Errorf("GetTopRatedTVShowsFunc not implemented")
	}
	return m.GetTopRatedTVShowsFunc(ctx, maxItems)
}

// GetOnTheAirTVShows calls the mock function.
func (m *MockTMDBClient) GetOnTheAirTVShows(
	ctx context.Context,
	maxItems int,
) ([]internal.TVShow, error) {
	m.recordCall("GetOnTheAirTVShows", ctx, maxItems)
	if m.GetOnTheAirTVShowsFunc == nil {
		return nil, fmt.Errorf("GetOnTheAirTVShowsFunc not implemented")
	}
	return m.GetOnTheAirTVShowsFunc(ctx, maxItems)
}

// SearchTVShows calls the mock function.
func (m *MockTMDBClient) SearchTVShows(
	ctx context.Context,
	query string,
	maxItems int,
) ([]internal.TVShow, error) {
	m.recordCall("SearchTVShows", ctx, query, maxItems)
	if m.SearchTVShowsFunc == nil {
		return nil, fmt.Errorf("SearchTVShowsFunc not implemented")
	}
	return m.SearchTVShowsFunc(ctx, query, maxItems)
}

// DiscoverTVShows calls the mock function.
func (m *MockTMDBClient) DiscoverTVShows(
	ctx context.Context,
	opts internal.SearchOptions,
) ([]internal.TVShow, error) {
	m.recordCall("DiscoverTVShows", ctx, opts)
	if m.DiscoverTVShowsFunc == nil {
		return nil, fmt.Errorf("DiscoverTVShowsFunc not implemented")
	}
	return m.DiscoverTVShowsFunc(ctx, opts)
}

// Helper methods for verification

// CallCount returns the number of times a method was called.
func (m *MockTMDBClient) CallCount(method string) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	count := 0
	for _, call := range m.Calls {
		if call.Method == method {
			count++
		}
	}
	return count
}

// WasCalled checks if a method was called at least once.
func (m *MockTMDBClient) WasCalled(method string) bool {
	return m.CallCount(method) > 0
}

// GetCallArgs returns the arguments for a specific call.
func (m *MockTMDBClient) GetCallArgs(method string, callIndex int) ([]any, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	calls := []MethodCall{}
	for _, call := range m.Calls {
		if call.Method == method {
			calls = append(calls, call)
		}
	}

	if callIndex >= len(calls) {
		return nil, fmt.Errorf("call index %d out of range for method %s (total calls: %d)",
			callIndex, method, len(calls))
	}

	return calls[callIndex].Args, nil
}

// Reset clears all recorded calls.
func (m *MockTMDBClient) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Calls = []MethodCall{}
}

// findMatchingCall searches for a call with the given method name.
func (m *MockTMDBClient) findMatchingCall(method string) *MethodCall {
	for i := range m.Calls {
		if m.Calls[i].Method == method {
			return &m.Calls[i]
		}
	}
	return nil
}

// argsMatch compares expected arguments with actual arguments.
func argsMatch(actual []any, expected []any) bool {
	if len(actual) < len(expected) {
		return false
	}

	for i, exp := range expected {
		if fmt.Sprintf("%v", actual[i]) != fmt.Sprintf("%v", exp) {
			return false
		}
	}
	return true
}

// AssertCalled verifies a method was called with specific arguments.
func (m *MockTMDBClient) AssertCalled(method string, expectedArgs ...any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	call := m.findMatchingCall(method)
	if call == nil {
		return fmt.Errorf("method %s was not called", method)
	}

	if len(expectedArgs) == 0 {
		return nil // Method was called, no args to check
	}

	if argsMatch(call.Args, expectedArgs) {
		return nil
	}

	return fmt.Errorf("method %s was not called with expected arguments", method)
}

// Ensure MockTMDBClient implements TMDBClient interface.
var _ internal.TMDBClient = (*MockTMDBClient)(nil)
