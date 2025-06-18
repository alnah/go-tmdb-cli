package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestCLIDispatcher_HelpCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{
			name:    "help command",
			command: "help",
			args:    []string{},
		},
		{
			name:    "help command with --help flag",
			command: "--help",
			args:    []string{},
		},
		{
			name:    "help command with -h flag",
			command: "-h",
			args:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := helpers.NewMockTMDBClient()
			logger := internal.NewLogger(internal.SilentLevel)
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			ctx := context.Background()
			err := dispatcher.DispatchCommand(
				ctx,
				tt.command,
				tt.args,
				"test-version",
				"test-date",
				"test-commit",
			)

			assert.NoError(t, err)
		})
	}
}

func TestCLIDispatcher_VersionCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{
			name:    "version command",
			command: "version",
			args:    []string{},
		},
		{
			name:    "version command with --version flag",
			command: "--version",
			args:    []string{},
		},
		{
			name:    "version command with -v flag",
			command: "-v",
			args:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := helpers.NewMockTMDBClient()
			logger := internal.NewLogger(internal.SilentLevel)
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			ctx := context.Background()
			err := dispatcher.DispatchCommand(
				ctx,
				tt.command,
				tt.args,
				"test-version",
				"test-date",
				"test-commit",
			)

			assert.NoError(t, err)
		})
	}
}

func TestCLIDispatcher_ConfigCommand(t *testing.T) {
	t.Parallel()

	mockClient := helpers.NewMockTMDBClient()
	logger := internal.NewLogger(internal.SilentLevel)
	dispatcher := internal.NewCLIDispatcher(mockClient, logger)

	ctx := context.Background()
	err := dispatcher.DispatchCommand(
		ctx,
		"config",
		[]string{},
		"test-version",
		"test-date",
		"test-commit",
	)

	assert.NoError(t, err)
}

func TestCLIDispatcher_PopularMovieCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{
			name:    "popular movies command",
			command: "popular",
			args:    []string{"10"},
		},
		{
			name:    "popular movies command with alias",
			command: "pop",
			args:    []string{"5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := helpers.NewMockTMDBClient()
			mockClient.GetPopularMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
				return []internal.Movie{}, nil
			}

			logger := internal.NewLogger(internal.SilentLevel)
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			ctx := context.Background()
			err := dispatcher.DispatchCommand(
				ctx,
				tt.command,
				tt.args,
				"test-version",
				"test-date",
				"test-commit",
			)

			assert.NoError(t, err)
		})
	}
}

func TestCLIDispatcher_TopRatedMovieCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{
			name:    "top-rated movies command",
			command: "top-rated",
			args:    []string{"3"},
		},
		{
			name:    "top-rated movies with alias",
			command: "top",
			args:    []string{"8"},
		},
		{
			name:    "top-rated movies with rated alias",
			command: "rated",
			args:    []string{"12"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := helpers.NewMockTMDBClient()
			mockClient.GetTopRatedMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
				return []internal.Movie{}, nil
			}

			logger := internal.NewLogger(internal.SilentLevel)
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			ctx := context.Background()
			err := dispatcher.DispatchCommand(
				ctx,
				tt.command,
				tt.args,
				"test-version",
				"test-date",
				"test-commit",
			)

			assert.NoError(t, err)
		})
	}
}

func TestCLIDispatcher_NowPlayingMovieCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{
			name:    "now-playing movies command",
			command: "now-playing",
			args:    []string{"15"},
		},
		{
			name:    "now-playing movies with now alias",
			command: "now",
			args:    []string{"7"},
		},
		{
			name:    "now-playing movies with playing alias",
			command: "playing",
			args:    []string{"4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := helpers.NewMockTMDBClient()
			mockClient.GetNowPlayingMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
				return []internal.Movie{}, nil
			}

			logger := internal.NewLogger(internal.SilentLevel)
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			ctx := context.Background()
			err := dispatcher.DispatchCommand(
				ctx,
				tt.command,
				tt.args,
				"test-version",
				"test-date",
				"test-commit",
			)

			assert.NoError(t, err)
		})
	}
}

func TestCLIDispatcher_UpcomingMovieCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{
			name:    "upcoming movies command",
			command: "upcoming",
			args:    []string{"20"},
		},
		{
			name:    "upcoming movies with soon alias",
			command: "soon",
			args:    []string{"6"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := helpers.NewMockTMDBClient()
			mockClient.GetUpcomingMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
				return []internal.Movie{}, nil
			}

			logger := internal.NewLogger(internal.SilentLevel)
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			ctx := context.Background()
			err := dispatcher.DispatchCommand(
				ctx,
				tt.command,
				tt.args,
				"test-version",
				"test-date",
				"test-commit",
			)

			assert.NoError(t, err)
		})
	}
}

func TestCLIDispatcher_SearchCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{
			name:    "search movies command",
			command: "search",
			args:    []string{"Matrix"},
		},
		{
			name:    "search movies with find alias",
			command: "find",
			args:    []string{"Inception"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := helpers.NewMockTMDBClient()
			mockClient.SearchMoviesFunc = func(_ context.Context, _ string, _ int) ([]internal.Movie, error) {
				return []internal.Movie{}, nil
			}

			logger := internal.NewLogger(internal.SilentLevel)
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			ctx := context.Background()
			err := dispatcher.DispatchCommand(
				ctx,
				tt.command,
				tt.args,
				"test-version",
				"test-date",
				"test-commit",
			)

			assert.NoError(t, err)
		})
	}
}

func TestCLIDispatcher_TVCommands(t *testing.T) {
	t.Parallel()

	mockClient := helpers.NewMockTMDBClient()
	mockClient.GetPopularTVShowsFunc = func(_ context.Context, _ int) ([]internal.TVShow, error) {
		return []internal.TVShow{}, nil
	}

	logger := internal.NewLogger(internal.SilentLevel)
	dispatcher := internal.NewCLIDispatcher(mockClient, logger)

	ctx := context.Background()
	err := dispatcher.DispatchCommand(
		ctx,
		"tv",
		[]string{"popular", "5"},
		"test-version",
		"test-date",
		"test-commit",
	)

	assert.NoError(t, err)
}

func TestCLIDispatcher_AutoSearchCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		command string
		args    []string
	}{
		{
			name:    "auto-search with movie title",
			command: "The Matrix",
			args:    []string{},
		},
		{
			name:    "auto-search with single word",
			command: "Inception",
			args:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := helpers.NewMockTMDBClient()
			mockClient.SearchMoviesFunc = func(_ context.Context, _ string, _ int) ([]internal.Movie, error) {
				return []internal.Movie{}, nil
			}

			logger := internal.NewLogger(internal.SilentLevel)
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			ctx := context.Background()
			err := dispatcher.DispatchCommand(
				ctx,
				tt.command,
				tt.args,
				"test-version",
				"test-date",
				"test-commit",
			)

			assert.NoError(t, err)
		})
	}
}

func TestCLIDispatcher_ErrorCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		command     string
		args        []string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "unknown command",
			command:     "unknown-command",
			args:        []string{},
			expectError: true,
			errorMsg:    "unknown command: unknown-command",
		},
		{
			name:        "empty command",
			command:     "",
			args:        []string{},
			expectError: true,
			errorMsg:    "unknown command:",
		},
		{
			name:        "numeric command",
			command:     "123",
			args:        []string{},
			expectError: true,
			errorMsg:    "unknown command: 123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := helpers.NewMockTMDBClient()
			logger := internal.NewLogger(internal.SilentLevel)
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			ctx := context.Background()
			err := dispatcher.DispatchCommand(
				ctx,
				tt.command,
				tt.args,
				"test-version",
				"test-date",
				"test-commit",
			)

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
				mock.GetPopularMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "top-rated alias top",
			command: "top",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetTopRatedMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "top-rated alias rated",
			command: "rated",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetTopRatedMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "now-playing alias now",
			command: "now",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetNowPlayingMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "now-playing alias playing",
			command: "playing",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetNowPlayingMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "upcoming alias soon",
			command: "soon",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetUpcomingMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "search alias find",
			command: "find",
			args:    []string{"Matrix"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(_ context.Context, _ string, _ int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
	}

	for _, tt := range tests {
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := helpers.NewMockTMDBClient()
			mockClient.SearchMoviesFunc = func(_ context.Context, _ string, _ int) ([]internal.Movie, error) {
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
				// Split long line to fix lll linter issue
				assert.False(t, mockClient.WasCalled("SearchMovies"),
					"SearchMovies should not have been called for: %s", tt.command)
			}
		})
	}
}

func TestCLIDispatcher_ContextPropagation(t *testing.T) {
	t.Parallel()

	mockClient := helpers.NewMockTMDBClient()
	mockClient.GetPopularMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
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
