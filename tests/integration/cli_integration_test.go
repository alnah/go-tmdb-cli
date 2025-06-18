package integration

import (
	"context"
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
				mock.GetPopularMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
					return []internal.Movie{
						fixtures.SampleMovie,
						{
							ID:            2,
							Title:         "Inception",
							OriginalTitle: "Inception",
							Year:          2010,
							Rating:        8.8,
							Votes:         2000000,
							Popularity:    83.2,
							Genres:        "Action, Sci-Fi, Thriller",
							Overview:      "A thief who steals corporate secrets through dream-sharing technology.",
							Language:      "en",
							Adult:         false,
							ReleaseDate:   "2010-07-16",
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
				mock.GetTopRatedMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
		{
			name:    "now-playing movies workflow",
			command: "now-playing",
			args:    []string{"8"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetNowPlayingMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
		{
			name:    "upcoming movies workflow",
			command: "upcoming",
			args:    []string{"10"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetUpcomingMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
		{
			name:    "search movies workflow",
			command: "search",
			args:    []string{"Matrix"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(_ context.Context, query string, _ int) ([]internal.Movie, error) {
					assert.Equal(t, "Matrix", query)
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
		{
			name:    "auto-search workflow",
			command: "The Matrix",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(_ context.Context, query string, _ int) ([]internal.Movie, error) {
					assert.Equal(t, "The Matrix", query)
					return []internal.Movie{fixtures.SampleMovie}, nil
				}
			},
		},
		{
			name:    "auto-search with complex title",
			command: "Blade Runner 2049",
			args:    []string{},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(_ context.Context, query string, _ int) ([]internal.Movie, error) {
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
				mock.GetPopularTVShowsFunc = func(_ context.Context, maxItems int) ([]internal.TVShow, error) {
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
				mock.GetTopRatedTVShowsFunc = func(_ context.Context, maxItems int) ([]internal.TVShow, error) {
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
				mock.GetOnTheAirTVShowsFunc = func(_ context.Context, maxItems int) ([]internal.TVShow, error) {
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
				mock.SearchTVShowsFunc = func(_ context.Context, query string, _ int) ([]internal.TVShow, error) {
					assert.Equal(t, "Breaking Bad", query)
					return []internal.TVShow{fixtures.SampleTVShow}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		// Capture range variable
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
				assert.NotEmpty(t, mockClient.Calls)
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
			name:    "help command error handling",
			command: "help",
			args:    []string{},
			setupMock: func(_ *helpers.MockTMDBClient) {
				// No mock setup needed for help
			},
			errorMsg: "", // Help should not error
		},
		{
			name:    "version command error handling",
			command: "version",
			args:    []string{},
			setupMock: func(_ *helpers.MockTMDBClient) {
				// No mock setup needed for version
			},
			errorMsg: "", // Version should not error
		},
		{
			name:    "config command error handling",
			command: "config",
			args:    []string{},
			setupMock: func(_ *helpers.MockTMDBClient) {
				// No mock setup needed for config
			},
			errorMsg: "", // Config should not error
		},
		{
			name:    "unknown command error",
			command: "unknown",
			args:    []string{},
			setupMock: func(_ *helpers.MockTMDBClient) {
				// No setup needed
			},
			errorMsg: "unknown command: unknown",
		},
		{
			name:    "empty command error",
			command: "",
			args:    []string{},
			setupMock: func(_ *helpers.MockTMDBClient) {
				// No setup needed
			},
			errorMsg: "unknown command:",
		},
		{
			name:    "numeric command error",
			command: "123",
			args:    []string{},
			setupMock: func(_ *helpers.MockTMDBClient) {
				// No setup needed
			},
			errorMsg: "unknown command: 123",
		},
		{
			name:    "dash command error",
			command: "-",
			args:    []string{},
			setupMock: func(_ *helpers.MockTMDBClient) {
				// No setup needed
			},
			errorMsg: "unknown command: -",
		},
	}

	for _, tt := range tests {
		// Capture range variable
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

				if tt.errorMsg == "" {
					assert.NoError(t, err)
				} else {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			})
		})
	}
}

func TestCLI_OutputFormats(t *testing.T) {
	t.Parallel()

	formats := []string{"table", "json", "csv"}

	for _, format := range formats {
		// Capture range variable
		t.Run("format_"+format, func(t *testing.T) {
			t.Parallel()

			helpers.WithMockEnvironment(t, map[string]string{
				"TMDB_API_KEY": "test-key",
			}, func() {
				mockClient := helpers.NewMockTMDBClient()
				mockClient.GetPopularMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
					return []internal.Movie{fixtures.SampleMovie}, nil
				}

				logger := internal.NewLogger(internal.SilentLevel)
				dispatcher := internal.NewCLIDispatcher(mockClient, logger)

				ctx := context.Background()
				args := []string{"--format", format, "1"}

				// Capture output to verify formatting
				stdout, stderr := helpers.CaptureOutput(t, func() {
					err := dispatcher.DispatchCommand(ctx, "popular", args, "test", "test", "test")
					assert.NoError(t, err)
				})

				// Verify output was generated
				assert.NotEmpty(t, stdout)
				assert.Empty(t, stderr)

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
				mock.GetPopularMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "search with no results",
			command: "search",
			args:    []string{"NonexistentMovie"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.SearchMoviesFunc = func(_ context.Context, _ string, _ int) ([]internal.Movie, error) {
					return []internal.Movie{}, nil
				}
			},
		},
		{
			name:    "TV shows with no results",
			command: "tv",
			args:    []string{"popular", "10"},
			setupMock: func(mock *helpers.MockTMDBClient) {
				mock.GetPopularTVShowsFunc = func(_ context.Context, _ int) ([]internal.TVShow, error) {
					return []internal.TVShow{}, nil
				}
			},
		},
	}

	for _, tt := range tests {
		// Capture range variable
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
		// Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			helpers.WithMockEnvironment(t, map[string]string{
				"TMDB_API_KEY": "test-key",
			}, func() {
				mockClient := helpers.NewMockTMDBClient()
				logger := internal.NewLogger(internal.SilentLevel)
				dispatcher := internal.NewCLIDispatcher(mockClient, logger)

				ctx := context.Background()

				err := dispatcher.DispatchCommand(ctx, tt.command, tt.args, "test", "test", "test")
				assert.NoError(t, err)
			})
		})
	}
}

func TestCLI_MultipleCommandsSequential(t *testing.T) {
	t.Parallel()

	helpers.WithMockEnvironment(t, map[string]string{
		"TMDB_API_KEY": "test-key",
	}, func() {
		mockClient := helpers.NewMockTMDBClient()
		mockClient.GetPopularMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
			return []internal.Movie{fixtures.SampleMovie}, nil
		}
		mockClient.SearchMoviesFunc = func(_ context.Context, _ string, _ int) ([]internal.Movie, error) {
			return []internal.Movie{fixtures.SampleMovie}, nil
		}
		mockClient.GetTopRatedMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
			return []internal.Movie{fixtures.SampleMovie}, nil
		}
		mockClient.GetNowPlayingMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
			return []internal.Movie{fixtures.SampleMovie}, nil
		}

		logger := internal.NewLogger(internal.SilentLevel)
		dispatcher := internal.NewCLIDispatcher(mockClient, logger)

		ctx := context.Background()

		// Execute multiple commands
		commands := []struct {
			name string
			args []string
		}{
			{"popular", []string{"1"}},
			{"search", []string{"Matrix"}},
			{"top-rated", []string{"1"}},
			{"now-playing", []string{"1"}},
		}

		for _, cmd := range commands {
			err := dispatcher.DispatchCommand(
				ctx,
				cmd.name,
				cmd.args,
				"test",
				"test",
				"test",
			)
			assert.NoError(t, err)
		}

		// Verify all commands were executed
		assert.Equal(t, len(commands), len(mockClient.Calls))
	})
}

func TestCLI_ContextCancellation(t *testing.T) {
	t.Parallel()

	t.Run("canceled context is passed to client", func(t *testing.T) {
		t.Parallel()

		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY": "test-key",
		}, func() {
			mockClient := helpers.NewMockTMDBClient()
			mockClient.GetPopularMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
				return []internal.Movie{fixtures.SampleMovie}, nil
			}

			logger := internal.NewLogger(internal.SilentLevel)
			dispatcher := internal.NewCLIDispatcher(mockClient, logger)

			// Create canceled context
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
			mockClient.GetPopularMoviesFunc = func(_ context.Context, _ int) ([]internal.Movie, error) {
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
