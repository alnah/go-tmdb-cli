package integration

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestCLI_FullWorkflow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		command   string
		args      []string
		setupMock func(*helpers.MockTMDBClient)
	}{
		{
			name:    "popular movies workflow",
			command: "popular",
			args:    []string{"5"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetPopularMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{
						fixtures.SampleMovie,
						{
							ID:          2,
							Title:       "Inception",
							ReleaseDate: "2010-07-16",
							Rating:      8.8,
							Votes:       2000000,
							Popularity:  83.2,
							Genres:      "Action, Sci-Fi, Thriller",
							Overview:    "A thief who steals corporate secrets through dream-sharing technology.",
						},
					}, nil
				}
			},
		},
		{
			name:    "top-rated movies workflow",
			command: "top-rated",
			args:    []string{"3"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetTopRatedMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					assert.Equal(t, 3, maxItems)
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
		{
			name:    "now-playing movies workflow",
			command: "now-playing",
			args:    []string{"10"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetNowPlayingMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					assert.Equal(t, 10, maxItems)
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
		{
			name:    "upcoming movies workflow",
			command: "upcoming",
			args:    []string{"15"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetUpcomingMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					assert.Equal(t, 15, maxItems)
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
		{
			name:    "search movies workflow",
			command: "search",
			args:    []string{"Matrix"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(ctx context.Context, query string, maxItems int) ([]internal.Movie, error) {
					assert.Equal(t, "Matrix", query)
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
		{
			name:    "search movies with multiple words",
			command: "search",
			args:    []string{"The", "Dark", "Knight"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(ctx context.Context, query string, maxItems int) ([]internal.Movie, error) {
					assert.Equal(t, "The Dark Knight", query)
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
		{
			name:    "auto-search workflow",
			command: "The Matrix",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(ctx context.Context, query string, maxItems int) ([]internal.Movie, error) {
					assert.Equal(t, "The Matrix", query)
					assert.Equal(t, internal.AutoSearchMaxItems, maxItems)
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
		{
			name:    "auto-search with multiple words",
			command: "Blade Runner 2049",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(ctx context.Context, query string, maxItems int) ([]internal.Movie, error) {
					assert.Equal(t, "Blade Runner 2049", query)
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
		{
			name:    "TV popular shows workflow",
			command: "tv",
			args:    []string{"popular", "3"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetPopularTVShowsFunc = func(ctx context.Context, maxItems int) ([]internal.TVShow, error) {
					assert.Equal(t, 3, maxItems)
					return []internal.TVShow{fixtures.SampleTVShow}, nil
				}
			},
		},
		{
			name:    "TV top-rated shows workflow",
			command: "tv",
			args:    []string{"top-rated", "5"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetTopRatedTVShowsFunc = func(ctx context.Context, maxItems int) ([]internal.TVShow, error) {
					assert.Equal(t, 5, maxItems)
					return []internal.TVShow{fixtures.SampleTVShow}, nil
				}
			},
		},
		{
			name:    "TV on-the-air shows workflow",
			command: "tv",
			args:    []string{"on-the-air", "8"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetOnTheAirTVShowsFunc = func(ctx context.Context, maxItems int) ([]internal.TVShow, error) {
					assert.Equal(t, 8, maxItems)
					return []internal.TVShow{fixtures.SampleTVShow}, nil
				}
			},
		},
		{
			name:    "TV search workflow",
			command: "tv",
			args:    []string{"search", "Breaking Bad"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchTVShowsFunc = func(ctx context.Context, query string, maxItems int) ([]internal.TVShow, error) {
					assert.Equal(t, "Breaking Bad", query)
					return []internal.TVShow{fixtures.SampleTVShow}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Setup environment
			helpers.WithMockEnvironment(t, map[string]string{
				"TMDB_API_KEY": "test-key",
			}, func() {
				// Setup mock client
				mockClient := helpers.NewMockTMDBClient()
				tt.setupMock(mockClient)

				logger := internal.NewLogger(internal.SilentLevel)
				dispatcher := internal.NewCLIDispatcher(mockClient, logger)

				ctx := context.Background()

				// Execute command
				err := dispatcher.DispatchCommand(
					ctx,
					tt.command,
					tt.args,
					"test-version",
					"test-date",
					"test-commit",
				)
				require.NoError(t, err)

				// Verify mock was called
				assert.NotEmpty(t, mockClient.Calls, "Mock client should have been called")
			})
		})
	}
}

func TestCLI_ErrorHandling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		command   string
		args      []string
		setupMock func(*helpers.MockTMDBClient)
		errorMsg  string
	}{
		{
			name:    "API error from popular movies",
			command: "popular",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetPopularMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return nil, errors.New("API rate limit exceeded")
				}
			},
			errorMsg: "API rate limit exceeded",
		},
		{
			name:    "API error from search",
			command: "search",
			args:    []string{"Matrix"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(ctx context.Context, query string, maxItems int) ([]internal.Movie, error) {
					return nil, errors.New("network timeout")
				}
			},
			errorMsg: "network timeout",
		},
		{
			name:    "API error from TV shows",
			command: "tv",
			args:    []string{"popular", "10"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetPopularTVShowsFunc = func(ctx context.Context, maxItems int) ([]internal.TVShow, error) {
					return nil, errors.New("service unavailable")
				}
			},
			errorMsg: "service unavailable",
		},
		{
			name:    "invalid count argument",
			command: "popular",
			args:    []string{"invalid"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				// Should not be called
			},
			errorMsg: "invalid count",
		},
		{
			name:    "zero count argument",
			command: "top-rated",
			args:    []string{"0"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				// Should not be called
			},
			errorMsg: "max-items must be between 1 and 1000",
		},
		{
			name:    "negative count argument",
			command: "now-playing",
			args:    []string{"-5"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				// Should not be called
			},
			errorMsg: "max-items must be between 1 and 1000",
		},
		{
			name:    "count too large",
			command: "upcoming",
			args:    []string{"1001"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				// Should not be called
			},
			errorMsg: "max-items must be between 1 and 1000",
		},
		{
			name:    "search without query",
			command: "search",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				// Should not be called
			},
			errorMsg: "search query required",
		},
		{
			name:    "TV command without subcommand",
			command: "tv",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				// Should not be called
			},
			errorMsg: "missing TV subcommand",
		},
		{
			name:    "TV command with invalid subcommand",
			command: "tv",
			args:    []string{"invalid-subcommand"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				// Should not be called
			},
			errorMsg: "unknown TV subcommand",
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			helpers.WithMockEnvironment(t, map[string]string{
				"TMDB_API_KEY": "test-key",
			}, func() {
				mockClient := helpers.NewMockTMDBClient()
				tt.setupMock(mockClient)

				logger := internal.NewLogger(internal.SilentLevel)
				dispatcher := internal.NewCLIDispatcher(mockClient, logger)

				ctx := context.Background()

				// Execute command
				err := dispatcher.DispatchCommand(ctx, tt.command, tt.args, "test", "test", "test")
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			})
		})
	}
}

func TestCLI_OutputFormats(t *testing.T) {
	t.Parallel()

	formats := []string{"table", "json", "csv"}

	for _, format := range formats {
		format := format // Capture range variable
		t.Run("format_"+format, func(t *testing.T) {
			t.Parallel()

			helpers.WithMockEnvironment(t, map[string]string{
				"TMDB_API_KEY": "test-key",
			}, func() {
				mockClient := helpers.NewMockTMDBClient()
				mockClient.GetPopularMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{fixtures.SampleMovie}, nil
				}

				logger := internal.NewLogger(internal.SilentLevel)
				dispatcher := internal.NewCLIDispatcher(mockClient, logger)

				ctx := context.Background()
				args := []string{"--format", format, "1"}

				// Capture output to verify formatting
				stdout := helpers.CaptureStdout(t, func() {
					err := dispatcher.DispatchCommand(ctx, "popular", args, "test", "test", "test")
					assert.NoError(t, err)
				})

				// Verify output was generated
				assert.NotEmpty(t, stdout)

				// Verify format-specific output
				switch format {
				case "json":
					assert.Contains(t, stdout, `"title"`)
				case "csv":
					assert.Contains(t, stdout, fixtures.SampleMovie.Title)
				case "table":
					assert.Contains(t, stdout, fixtures.SampleMovie.Title)
				}
			})
		})
	}
}

func TestCLI_EmptyResults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		command   string
		args      []string
		setupMock func(*helpers.MockTMDBClient)
	}{
		{
			name:    "popular movies with no results",
			command: "popular",
			args:    []string{"10"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetPopularMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "search with no results",
			command: "search",
			args:    []string{"NonexistentMovie"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(ctx context.Context, query string, maxItems int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "TV shows with no results",
			command: "tv",
			args:    []string{"popular", "10"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetPopularTVShowsFunc = func(ctx context.Context, maxItems int) ([]internal.TVShow, error) {
					return []internal.TVShow{}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			helpers.WithMockEnvironment(t, map[string]string{
				"TMDB_API_KEY": "test-key",
			}, func() {
				mockClient := helpers.NewMockTMDBClient()
				tt.setupMock(mockClient)

				logger := internal.NewLogger(internal.SilentLevel)
				dispatcher := internal.NewCLIDispatcher(mockClient, logger)

				ctx := context.Background()

				err := dispatcher.DispatchCommand(ctx, tt.command, tt.args, "test", "test", "test")
				assert.NoError(t, err)
			})
		})
	}
}

func TestCLI_SpecialCommands(t *testing.T) {
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
			name:    "version command",
			command: "version",
			args:    []string{},
		},
		{
			name:    "config command",
			command: "config",
			args:    []string{},
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
				tt.command,
				tt.args,
				"1.0.0",
				"2024-01-01",
				"abc123",
			)
			assert.NoError(t, err)

			// Special commands should not call the API
			assert.Empty(t, mockClient.Calls, "Special commands should not call API")
		})
	}
}

func TestCLI_WithFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		command   string
		args      []string
		setupMock func(*helpers.MockTMDBClient)
	}{
		{
			name:    "popular movies with verbose flag",
			command: "popular",
			args:    []string{"--verbose", "5"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetPopularMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					assert.Equal(t, 5, maxItems)
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
		{
			name:    "search with max-items flag",
			command: "search",
			args:    []string{"--max-items", "50", "Matrix"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(ctx context.Context, query string, maxItems int) ([]internal.Movie, error) {
					assert.Equal(t, "Matrix", query)
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
		{
			name:    "top-rated with original-title flag",
			command: "top-rated",
			args:    []string{"--original-title", "10"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetTopRatedMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					assert.Equal(t, 10, maxItems)
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
		{
			name:    "now-playing with no-header flag",
			command: "now-playing",
			args:    []string{"--no-header", "--format", "table", "5"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetNowPlayingMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
					assert.Equal(t, 5, maxItems)
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			helpers.WithMockEnvironment(t, map[string]string{
				"TMDB_API_KEY": "test-key",
			}, func() {
				mockClient := helpers.NewMockTMDBClient()
				tt.setupMock(mockClient)

				logger := internal.NewLogger(internal.SilentLevel)
				dispatcher := internal.NewCLIDispatcher(mockClient, logger)

				ctx := context.Background()

				err := dispatcher.DispatchCommand(ctx, tt.command, tt.args, "test", "test", "test")
				assert.NoError(t, err)
				assert.NotEmpty(t, mockClient.Calls)
			})
		})
	}
}

func TestCLI_ContextCancellation(t *testing.T) {
	t.Parallel()

	t.Run("cancelled context is passed to client", func(t *testing.T) {
		t.Parallel()

		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY": "test-key",
		}, func() {
			mockClient := helpers.NewMockTMDBClient()
			mockClient.GetPopularMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
				// Verify context is passed through
				assert.NotNil(t, ctx)
				return []internal.Movie{fixtures.SampleMovie}, nil
			}

			logger := internal.NewLogger(internal.SilentLevel)
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			// Create cancelled context
			ctx, cancel := context.WithCancel(context.Background())
			cancel() // Cancel immediately

			err := dispatcher.DispatchCommand(
				ctx,
				"popular",
				[]string{"10"},
				"test",
				"test",
				"test",
			)
			assert.NoError(t, err) // Mock doesn't check cancellation
		})
	})
}

func TestCLI_Integration_Performance(t *testing.T) {
	t.Parallel()

	t.Run("CLI handles multiple rapid commands", func(t *testing.T) {
		t.Parallel()

		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY": "test-key",
		}, func() {
			mockClient := helpers.NewMockTMDBClient()
			mockClient.GetPopularMoviesFunc = func(ctx context.Context, maxItems int) ([]internal.Movie, error) {
				return []internal.Movie{fixtures.SampleMovie}, nil
			}

			logger := internal.NewLogger(internal.SilentLevel)
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			ctx := context.Background()

			// Execute multiple commands rapidly
			commands := []string{"popular", "top-rated", "now-playing", "upcoming"}
			for _, command := range commands {
				err := dispatcher.DispatchCommand(
					ctx,
					command,
					[]string{"1"},
					"test",
					"test",
					"test",
				)
				assert.NoError(t, err)
			}

			// Verify all commands were executed
			assert.Equal(t, len(commands), len(mockClient.Calls))
		})
	})
}
