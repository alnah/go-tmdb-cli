package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestCLIDispatcher_DispatchCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		command     string
		args        []string
		setupMock   func(*helpers.MockTMDBClient)
		expectError bool
		errorMsg    string
	}{
		{
			name:        "help command",
			command:     "help",
			args:        []string{},
			setupMock:   func(mock *helpers.MockTMDBClient) {},
			expectError: false,
		},
		{
			name:        "help command with --help flag",
			command:     "--help",
			args:        []string{},
			setupMock:   func(mock *helpers.MockTMDBClient) {},
			expectError: false,
		},
		{
			name:        "help command with -h flag",
			command:     "-h",
			args:        []string{},
			setupMock:   func(mock *helpers.MockTMDBClient) {},
			expectError: false,
		},
		{
			name:        "version command",
			command:     "version",
			args:        []string{},
			setupMock:   func(mock *helpers.MockTMDBClient) {},
			expectError: false,
		},
		{
			name:        "version command with --version flag",
			command:     "--version",
			args:        []string{},
			setupMock:   func(mock *helpers.MockTMDBClient) {},
			expectError: false,
		},
		{
			name:        "version command with -v flag",
			command:     "-v",
			args:        []string{},
			setupMock:   func(mock *helpers.MockTMDBClient) {},
			expectError: false,
		},
		{
			name:        "config command",
			command:     "config",
			args:        []string{},
			setupMock:   func(mock *helpers.MockTMDBClient) {},
			expectError: false,
		},
		{
			name:    "popular movies command",
			command: "popular",
			args:    []string{"10"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetPopularMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
			expectError: false,
		},
		{
			name:    "popular movies command with alias",
			command: "pop",
			args:    []string{"5"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetPopularMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
			expectError: false,
		},
		{
			name:    "top-rated movies command",
			command: "top-rated",
			args:    []string{"10"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetTopRatedMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
			expectError: false,
		},
		{
			name:    "top-rated movies with alias top",
			command: "top",
			args:    []string{"15"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetTopRatedMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
			expectError: false,
		},
		{
			name:    "top-rated movies with alias rated",
			command: "rated",
			args:    []string{"20"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetTopRatedMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
			expectError: false,
		},
		{
			name:    "now-playing movies command",
			command: "now-playing",
			args:    []string{"10"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetNowPlayingMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
			expectError: false,
		},
		{
			name:    "now-playing movies with alias now",
			command: "now",
			args:    []string{"10"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetNowPlayingMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
			expectError: false,
		},
		{
			name:    "now-playing movies with alias playing",
			command: "playing",
			args:    []string{"10"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetNowPlayingMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
			expectError: false,
		},
		{
			name:    "upcoming movies command",
			command: "upcoming",
			args:    []string{"10"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetUpcomingMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
			expectError: false,
		},
		{
			name:    "upcoming movies with alias soon",
			command: "soon",
			args:    []string{"10"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetUpcomingMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
			expectError: false,
		},
		{
			name:    "search movies command",
			command: "search",
			args:    []string{"Matrix"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(ctx context.Context, query string, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
			expectError: false,
		},
		{
			name:    "search movies with alias find",
			command: "find",
			args:    []string{"Matrix"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(ctx context.Context, query string, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
			expectError: false,
		},
		{
			name:    "TV command popular",
			command: "tv",
			args:    []string{"popular", "10"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetPopularTVShowsFunc = func(ctx context.Context, maxItems int) ([]internal.TVShow, error) {
					return []internal.TVShow{}, nil
				}
			},
			expectError: false,
		},
		{
			name:    "auto-search with movie title",
			command: "The Matrix",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(ctx context.Context, query string, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
			expectError: false,
		},
		{
			name:    "auto-search with multi-word title",
			command: "The Dark Knight",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(ctx context.Context, query string, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
			expectError: false,
		},
		{
			name:        "unknown command",
			command:     "unknown-command",
			args:        []string{},
			setupMock:   func(mock *helpers.MockTMDBClient) {},
			expectError: true,
			errorMsg:    "unknown command: unknown-command",
		},
		{
			name:        "short command that doesn't look like search",
			command:     "xy",
			args:        []string{},
			setupMock:   func(mock *helpers.MockTMDBClient) {},
			expectError: true,
			errorMsg:    "unknown command: xy",
		},
		{
			name:        "numeric command that doesn't look like search",
			command:     "123",
			args:        []string{},
			setupMock:   func(mock *helpers.MockTMDBClient) {},
			expectError: true,
			errorMsg:    "unknown command: 123",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable for parallel execution
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Make each subtest parallel

			// Setup mock client
			mockClient := helpers.NewMockTMDBClient()
			tt.setupMock(mockClient)

			// Create logger with silent level for tests
			logger := internal.NewLogger(internal.SilentLevel)

			// Create dispatcher
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			// Execute command
			ctx := context.Background()
			err := dispatcher.DispatchCommand(
				ctx,
				tt.command,
				tt.args,
				"test-version",
				"test-date",
				"test-commit",
			)

			// Verify results
			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCLIDispatcher_CommandAliases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		command   string
		args      []string
		setupMock func(*helpers.MockTMDBClient)
	}{
		{
			name:    "popular alias pop",
			command: "pop",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetPopularMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "top-rated alias top",
			command: "top",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetTopRatedMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "top-rated alias rated",
			command: "rated",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetTopRatedMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "now-playing alias now",
			command: "now",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetNowPlayingMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "now-playing alias playing",
			command: "playing",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetNowPlayingMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "upcoming alias soon",
			command: "soon",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetUpcomingMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "search alias find",
			command: "find",
			args:    []string{"Matrix"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(ctx context.Context, query string, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := helpers.NewMockTMDBClient()
			tt.setupMock(mockClient)

			logger := internal.NewLogger(internal.SilentLevel)
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			ctx := context.Background()
			err := dispatcher.DispatchCommand(ctx, tt.command, tt.args, "test", "test", "test")
			assert.NoError(t, err)
		})
	}
}

func TestCLIDispatcher_AutoSearchDetection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		command      string
		expectedCall bool
	}{
		{
			name:         "movie title with spaces",
			command:      "The Matrix",
			expectedCall: true,
		},
		{
			name:         "single word movie title",
			command:      "Inception",
			expectedCall: true,
		},
		{
			name:         "movie title with numbers and letters",
			command:      "Blade Runner 2049",
			expectedCall: true,
		},
		{
			name:         "lowercase movie title",
			command:      "avatar",
			expectedCall: true,
		},
		{
			name:         "mixed case movie title",
			command:      "SpIdEr-MaN",
			expectedCall: true,
		},
		{
			name:         "command that doesn't look like search",
			command:      "popular",
			expectedCall: false,
		},
		{
			name:         "short non-search command",
			command:      "xy",
			expectedCall: false,
		},
		{
			name:         "numeric only command",
			command:      "123",
			expectedCall: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := helpers.NewMockTMDBClient()
			mockClient.SearchMoviesFunc = func(ctx context.Context, query string, maxItems int) ([]internal.Movie, error) {
				return []internal.Movie{}, nil
			}

			logger := internal.NewLogger(internal.SilentLevel)
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			ctx := context.Background()
			_ = dispatcher.DispatchCommand(ctx, tt.command, []string{}, "test", "test", "test")

			if tt.expectedCall {
				assert.True(
					t,
					mockClient.WasCalled("SearchMovies"),
					"SearchMovies should have been called for: %s",
					tt.command,
				)
			} else {
				assert.False(t, mockClient.WasCalled("SearchMovies"), "SearchMovies should not have been called for: %s", tt.command)
			}
		})
	}
}

func TestCLIDispatcher_ContextPropagation(t *testing.T) {
	t.Parallel()

	mockClient := helpers.NewMockTMDBClient()
	mockClient.GetPopularMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
		// Verify context is passed through
		assert.NotNil(t, ctx)
		return []internal.Movie{}, nil
	}

	logger := internal.NewLogger(internal.SilentLevel)
	dispatcher := internal.NewCLIDispatcher(mockClient, logger)

	ctx := context.Background()
	err := dispatcher.DispatchCommand(ctx, "popular", []string{"10"}, "test", "test", "test")
	assert.NoError(t, err)
}

func TestCLIDispatcher_VersionHandling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		version   string
		buildDate string
		gitCommit string
	}{
		{
			name:      "complete version info",
			version:   "1.0.0",
			buildDate: "2024-01-01",
			gitCommit: "abc123",
		},
		{
			name:      "version with unknown build info",
			version:   "1.0.0",
			buildDate: "unknown",
			gitCommit: "unknown",
		},
		{
			name:      "development version",
			version:   "dev",
			buildDate: "dev",
			gitCommit: "dev",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := helpers.NewMockTMDBClient()
			logger := internal.NewLogger(internal.SilentLevel)
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			ctx := context.Background()
			err := dispatcher.DispatchCommand(
				ctx,
				"version",
				[]string{},
				tt.version,
				tt.buildDate,
				tt.gitCommit,
			)
			assert.NoError(t, err)
		})
	}
}
