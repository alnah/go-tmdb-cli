package unit

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestTVCommands_HandleSearch(t *testing.T) {
	t.Parallel()

	t.Run("successful TV search", func(t *testing.T) {
		t.Parallel()

		mockClient := helpers.NewMockTMDBClient()
		mockClient.SearchTVShowsFunc = func(_ context.Context, query string, maxItems int) ([]internal.TVShow, error) {
			assert.Equal(t, "Breaking Bad", query)
			assert.Equal(t, 20, maxItems)
			return []internal.TVShow{
				{
					ID:           1,
					Name:         "Breaking Bad",
					OriginalName: "Breaking Bad",
					Year:         2008,
					Rating:       9.3,
					Votes:        250000,
					Genres:       "Drama, Crime",
				},
			}, nil
		}

		logger := internal.NewLogger(internal.SilentLevel)
		tvCommands := internal.NewTVCommands(mockClient, logger)

		args := []string{"Breaking Bad"}
		err := tvCommands.HandleSearch(context.Background(), args)

		assert.NoError(t, err)
		require.Len(t, mockClient.Calls, 1)
		assert.Equal(t, "SearchTVShows", mockClient.Calls[0])
	})

	t.Run("TV search with format flags", func(t *testing.T) {
		t.Parallel()

		mockClient := helpers.NewMockTMDBClient()
		mockClient.SearchTVShowsFunc = func(_ context.Context, query string, _ int) ([]internal.TVShow, error) {
			assert.Equal(t, "The Office", query)
			return []internal.TVShow{
				{
					ID:           2,
					Name:         "The Office",
					OriginalName: "The Office",
					Year:         2005,
					Rating:       8.7,
					Votes:        180000,
					Genres:       "Comedy",
				},
			}, nil
		}

		logger := internal.NewLogger(internal.SilentLevel)
		tvCommands := internal.NewTVCommands(mockClient, logger)

		args := []string{"The Office", "--format", "json", "--max", "10"}
		err := tvCommands.HandleSearch(context.Background(), args)

		assert.NoError(t, err)
		require.Len(t, mockClient.Calls, 1)
		assert.Equal(t, "SearchTVShows", mockClient.Calls[0])
	})

	t.Run("TV search with empty results", func(t *testing.T) {
		t.Parallel()

		mockClient := helpers.NewMockTMDBClient()
		mockClient.SearchTVShowsFunc = func(_ context.Context, _ string, _ int) ([]internal.TVShow, error) {
			return []internal.TVShow{}, nil
		}

		logger := internal.NewLogger(internal.SilentLevel)
		tvCommands := internal.NewTVCommands(mockClient, logger)

		args := []string{"NonexistentShow"}
		err := tvCommands.HandleSearch(context.Background(), args)

		assert.NoError(t, err)
		require.Len(t, mockClient.Calls, 1)
		assert.Equal(t, "SearchTVShows", mockClient.Calls[0])
	})

	t.Run("TV search API error", func(t *testing.T) {
		t.Parallel()

		mockClient := helpers.NewMockTMDBClient()
		mockClient.SearchTVShowsFunc = func(_ context.Context, _ string, _ int) ([]internal.TVShow, error) {
			return nil, errors.New("API timeout")
		}

		logger := internal.NewLogger(internal.SilentLevel)
		tvCommands := internal.NewTVCommands(mockClient, logger)

		args := []string{"Test Show"}
		err := tvCommands.HandleSearch(context.Background(), args)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TV search failed")
		assert.Contains(t, err.Error(), "API timeout")
	})

	t.Run("TV search with invalid arguments", func(t *testing.T) {
		t.Parallel()

		mockClient := helpers.NewMockTMDBClient()
		logger := internal.NewLogger(internal.SilentLevel)
		tvCommands := internal.NewTVCommands(mockClient, logger)

		// Empty args should cause error
		args := []string{}
		err := tvCommands.HandleSearch(context.Background(), args)

		assert.Error(t, err)
	})
}

func TestTVCommands_HandleSearchAdvanced(t *testing.T) {
	t.Parallel()

	t.Run("TV search with progress indicator", func(t *testing.T) {
		t.Parallel()

		mockClient := helpers.NewMockTMDBClient()
		mockClient.SearchTVShowsFunc = func(_ context.Context, query string, maxItems int) ([]internal.TVShow, error) {
			assert.Equal(t, "Game of Thrones", query)
			assert.Equal(t, 50, maxItems) // Should trigger progress indicator
			return []internal.TVShow{
				{
					ID:           3,
					Name:         "Game of Thrones",
					OriginalName: "Game of Thrones",
					Year:         2011,
					Rating:       9.2,
					Votes:        500000,
					Genres:       "Drama, Fantasy",
				},
			}, nil
		}

		logger := internal.NewLogger(internal.SilentLevel)
		tvCommands := internal.NewTVCommands(mockClient, logger)

		// Capture stdout to verify progress indicator
		output := helpers.CaptureStdout(t, func() {
			// Test the HandleSearch function with high max items
			args := []string{"Game of Thrones", "--max", "50"}
			err := tvCommands.HandleSearch(context.Background(), args)
			assert.NoError(t, err)
		})
		require.Len(t, mockClient.Calls, 1)
		assert.Equal(t, "SearchTVShows", mockClient.Calls[0])

		// Check that progress was shown (function would show progress for 50 items)
		assert.Contains(t, output, "TV shows for")
	})

	t.Run("TV search with original title flag", func(t *testing.T) {
		t.Parallel()

		mockClient := helpers.NewMockTMDBClient()
		mockClient.SearchTVShowsFunc = func(_ context.Context, _ string, _ int) ([]internal.TVShow, error) {
			return []internal.TVShow{
				{
					ID:           4,
					Name:         "Money Heist",
					OriginalName: "La Casa de Papel",
					Year:         2017,
					Rating:       8.3,
					Votes:        200000,
					Genres:       "Drama, Crime",
				},
			}, nil
		}

		logger := internal.NewLogger(internal.SilentLevel)
		tvCommands := internal.NewTVCommands(mockClient, logger)

		args := []string{"Money Heist", "--original-title"}
		err := tvCommands.HandleSearch(context.Background(), args)

		assert.NoError(t, err)
		require.Len(t, mockClient.Calls, 1)
		assert.Equal(t, "SearchTVShows", mockClient.Calls[0])
	})

	t.Run("TV search with CSV format", func(t *testing.T) {
		t.Parallel()

		mockClient := helpers.NewMockTMDBClient()
		mockClient.SearchTVShowsFunc = func(_ context.Context, _ string, _ int) ([]internal.TVShow, error) {
			return []internal.TVShow{
				{
					ID:           5,
					Name:         "Stranger Things",
					OriginalName: "Stranger Things",
					Year:         2016,
					Rating:       8.7,
					Votes:        300000,
					Genres:       "Drama, Fantasy, Horror",
				},
			}, nil
		}

		logger := internal.NewLogger(internal.SilentLevel)
		tvCommands := internal.NewTVCommands(mockClient, logger)

		args := []string{"Stranger Things", "--format", "csv"}
		err := tvCommands.HandleSearch(context.Background(), args)

		assert.NoError(t, err)
		require.Len(t, mockClient.Calls, 1)
		assert.Equal(t, "SearchTVShows", mockClient.Calls[0])
	})

	t.Run("TV search handles empty search results correctly", func(t *testing.T) {
		t.Parallel()

		mockClient := helpers.NewMockTMDBClient()
		mockClient.SearchTVShowsFunc = func(_ context.Context, _ string, _ int) ([]internal.TVShow, error) {
			return []internal.TVShow{}, nil
		}

		logger := internal.NewLogger(internal.SilentLevel)
		tvCommands := internal.NewTVCommands(mockClient, logger)

		// Capture stdout to verify empty results handling
		output := helpers.CaptureStdout(t, func() {
			args := []string{"Unknown TV Show"}
			err := tvCommands.HandleSearch(context.Background(), args)
			assert.NoError(t, err)
		})

		// Check that appropriate message was shown for empty results
		assert.Contains(t, output, "No TV shows found")
		assert.Contains(t, output, "Try:")
		assert.Contains(t, output, "Include the first air year")
	})
}

func TestTVCommands_Dispatch(t *testing.T) {
	t.Parallel()

	t.Run("TV dispatch with search subcommand", func(t *testing.T) {
		t.Parallel()

		mockClient := helpers.NewMockTMDBClient()
		mockClient.SearchTVShowsFunc = func(_ context.Context, _ string, _ int) ([]internal.TVShow, error) {
			return []internal.TVShow{}, nil
		}

		logger := internal.NewLogger(internal.SilentLevel)
		tvCommands := internal.NewTVCommands(mockClient, logger)

		args := []string{"search", "Test Show"}
		err := tvCommands.Dispatch(context.Background(), args)

		assert.NoError(t, err)
		require.Len(t, mockClient.Calls, 1)
		assert.Equal(t, "SearchTVShows", mockClient.Calls[0])
	})

	t.Run("TV dispatch with find alias", func(t *testing.T) {
		t.Parallel()

		mockClient := helpers.NewMockTMDBClient()
		mockClient.SearchTVShowsFunc = func(_ context.Context, _ string, _ int) ([]internal.TVShow, error) {
			return []internal.TVShow{}, nil
		}

		logger := internal.NewLogger(internal.SilentLevel)
		tvCommands := internal.NewTVCommands(mockClient, logger)

		args := []string{"find", "Test Show"}
		err := tvCommands.Dispatch(context.Background(), args)

		assert.NoError(t, err)
		require.Len(t, mockClient.Calls, 1)
		assert.Equal(t, "SearchTVShows", mockClient.Calls[0])
	})

	t.Run("TV dispatch with empty args returns error", func(t *testing.T) {
		t.Parallel()

		mockClient := helpers.NewMockTMDBClient()
		logger := internal.NewLogger(internal.SilentLevel)
		tvCommands := internal.NewTVCommands(mockClient, logger)

		args := []string{}
		err := tvCommands.Dispatch(context.Background(), args)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tv command requires a subcommand")
	})

	t.Run("TV dispatch with unknown subcommand returns error", func(t *testing.T) {
		t.Parallel()

		mockClient := helpers.NewMockTMDBClient()
		logger := internal.NewLogger(internal.SilentLevel)
		tvCommands := internal.NewTVCommands(mockClient, logger)

		args := []string{"unknown", "arg"}
		err := tvCommands.Dispatch(context.Background(), args)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown tv subcommand: unknown")
	})
}

func TestMovieCommands_HandleSearch(t *testing.T) {
	t.Parallel()

	t.Run("movie search handles empty results", func(t *testing.T) {
		t.Parallel()

		mockClient := helpers.NewMockTMDBClient()
		mockClient.SearchMoviesFunc = func(_ context.Context, _ string, _ int) ([]internal.Movie, error) {
			return []internal.Movie{}, nil
		}

		logger := internal.NewLogger(internal.SilentLevel)
		movieCommands := internal.NewMovieCommands(mockClient, logger)

		// Capture stdout to verify empty results handling
		output := helpers.CaptureStdout(t, func() {
			args := []string{"Unknown Movie"}
			err := movieCommands.HandleSearch(context.Background(), args)
			assert.NoError(t, err)
		})

		// Check that appropriate message was shown for empty results
		assert.Contains(t, output, "No movies found")
		assert.Contains(t, output, "Try:")
		assert.Contains(t, output, "Include the release year")
	})

	t.Run("movie search with various flags", func(t *testing.T) {
		t.Parallel()

		mockClient := helpers.NewMockTMDBClient()
		mockClient.SearchMoviesFunc = func(_ context.Context, query string, maxItems int) ([]internal.Movie, error) {
			assert.Equal(t, "The Matrix", query)
			assert.Equal(t, 15, maxItems)
			return []internal.Movie{
				{
					ID:            1,
					Title:         "The Matrix",
					OriginalTitle: "The Matrix",
					Year:          1999,
					Rating:        8.7,
					Votes:         500000,
					Genres:        "Action, Sci-Fi",
				},
			}, nil
		}

		logger := internal.NewLogger(internal.SilentLevel)
		movieCommands := internal.NewMovieCommands(mockClient, logger)

		args := []string{"The Matrix", "--format", "json", "--max", "15", "--original-title"}
		err := movieCommands.HandleSearch(context.Background(), args)

		assert.NoError(t, err)
		require.Len(t, mockClient.Calls, 1)
		assert.Equal(t, "SearchMovies", mockClient.Calls[0])
	})
}

// Note: The handleEmptySearchResults function is not exported,
// so we test it indirectly through the HandleSearch methods which call it
// when search results are empty. This provides coverage for the function.
