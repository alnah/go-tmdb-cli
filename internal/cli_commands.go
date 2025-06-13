// Package internal provides core functionality for the TMDB CLI application.
package internal

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

// Command constants for consistent string usage.
const (
	CommandPopular    = "popular"
	CommandTopRated   = "top-rated"
	CommandNowPlaying = "now-playing"
	CommandUpcoming   = "upcoming"
	CommandOnTheAir   = "on-the-air"
	CommandSearch     = "search"
)

// Default values and limits.
const (
	DefaultMaxItems    = 20
	DefaultMaxWidth    = 120
	AutoSearchMaxItems = 20
)

// MovieCommands handles movie-related CLI commands.
type MovieCommands struct {
	client *Client
	logger *Logger
}

// NewMovieCommands creates a new MovieCommands instance.
func NewMovieCommands(client *Client, logger *Logger) *MovieCommands {
	return &MovieCommands{
		client: client,
		logger: logger,
	}
}

// TVCommands handles TV show-related CLI commands.
type TVCommands struct {
	client *Client
	logger *Logger
}

// NewTVCommands creates a new TVCommands instance.
func NewTVCommands(client *Client, logger *Logger) *TVCommands {
	return &TVCommands{
		client: client,
		logger: logger,
	}
}

// ConfigCommands handles configuration-related CLI commands.
type ConfigCommands struct{}

// NewConfigCommands creates a new ConfigCommands instance.
func NewConfigCommands() *ConfigCommands {
	return &ConfigCommands{}
}

// HandleConfig displays configuration help information.
func (cc *ConfigCommands) HandleConfig() error {
	fmt.Print(GetConfigHelp())
	return nil
}

// LooksLikeSearch determines if input looks like a search query.
func LooksLikeSearch(input string) bool {
	// Contains spaces, letters, or looks like a movie/TV title
	return strings.Contains(input, " ") ||
		(len(input) > 2 && strings.ContainsAny(input, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"))
}

// HandleList handles movie list commands (popular, top-rated, now-playing, upcoming).
func (mc *MovieCommands) HandleList(
	ctx context.Context,
	listType string,
	args []string,
) error {
	flags, posArgs, err := ParseCommonFlags(args)
	if err != nil {
		return err
	}

	// Parse count from positional args
	count, err := ParseCountArg(posArgs, flags.MaxItems)
	if err != nil {
		return err
	}

	if e := ValidateMaxItems(count); e != nil {
		return e
	}

	// Show progress for larger requests
	ShowFetchingProgress(count, listType)

	// Fetch movies based on type
	var movies []Movie
	switch listType {
	case CommandPopular:
		movies, err = mc.client.GetPopularMovies(ctx, count)
		mc.logger.Info("Fetched %d popular movies", len(movies))
	case CommandTopRated:
		movies, err = mc.client.GetTopRatedMovies(ctx, count)
		mc.logger.Info("Fetched %d top-rated movies", len(movies))
	case CommandNowPlaying:
		movies, err = mc.client.GetNowPlayingMovies(ctx, count)
		mc.logger.Info("Fetched %d now-playing movies", len(movies))
	case CommandUpcoming:
		movies, err = mc.client.GetUpcomingMovies(ctx, count)
		mc.logger.Info("Fetched %d upcoming movies", len(movies))
	default:
		return fmt.Errorf("unknown list type: %s", listType)
	}

	if err != nil {
		return fmt.Errorf("failed to fetch %s movies: %w", listType, err)
	}

	if len(movies) == 0 {
		fmt.Printf("No %s movies found.\n", listType)
		return nil
	}

	// Format and display
	options := BuildFormatOptions(flags)

	if options.Format == FormatTable {
		FormatSummary(os.Stdout, movies, listType, flags.OriginalTitle)
	}

	return FormatMovies(os.Stdout, movies, options)
}

// HandleSearch handles movie search commands.
func (mc *MovieCommands) HandleSearch(
	ctx context.Context,
	args []string,
) error {
	return genericSearch(ctx, mc.client, mc.logger, args, SearchTypeMovie)
}

// HandleAutoSearch handles automatic search detection for movie titles.
func (mc *MovieCommands) HandleAutoSearch(
	ctx context.Context,
	args []string,
) error {
	query := strings.Join(args, " ")
	query = CleanQuery(query)

	mc.logger.Info("Auto-searching for: %s", query)

	movies, err := mc.client.SearchMovies(ctx, query, AutoSearchMaxItems)
	if err != nil {
		return err
	}

	if len(movies) == 0 {
		fmt.Printf("No movies found for \"%s\"\n", query)
		return nil
	}

	fmt.Printf("Found %d movies for \"%s\"\n\n", len(movies), query)

	options := FormatOptions{
		Format:   FormatTable,
		NoHeader: false,
		MaxWidth: DefaultMaxWidth,
	}

	return FormatMovies(os.Stdout, movies, options)
}

// Dispatch handles TV subcommand routing.
func (tc *TVCommands) Dispatch(
	ctx context.Context,
	args []string,
) error {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "TV command requires a subcommand\n\n")
		ShowTVUsage()
		return fmt.Errorf("missing TV subcommand")
	}

	tvCommand := args[0]
	tvArgs := args[1:]

	switch tvCommand {
	case CommandPopular, "pop":
		return tc.HandleList(ctx, CommandPopular, tvArgs)
	case CommandTopRated, "top", "rated":
		return tc.HandleList(ctx, CommandTopRated, tvArgs)
	case CommandOnTheAir, "air", "airing":
		return tc.HandleList(ctx, CommandOnTheAir, tvArgs)
	case CommandSearch, "find":
		return tc.HandleSearch(ctx, tvArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown TV subcommand: %s\n\n", tvCommand)
		ShowTVUsage()
		return fmt.Errorf("unknown TV subcommand: %s", tvCommand)
	}
}

// HandleList handles TV show list commands.
func (tc *TVCommands) HandleList(
	ctx context.Context,
	listType string,
	args []string,
) error {
	flags, posArgs, err := ParseCommonFlags(args)
	if err != nil {
		return err
	}

	count, err := ParseCountArg(posArgs, flags.MaxItems)
	if err != nil {
		return err
	}

	if e := ValidateMaxItems(count); e != nil {
		return e
	}

	ShowProgress(fmt.Sprintf("Fetching %d %s TV shows", count, listType), count)

	var tvShows []TVShow
	switch listType {
	case CommandPopular:
		tvShows, err = tc.client.GetPopularTVShows(ctx, count)
		tc.logger.Info("Fetched %d popular TV shows", len(tvShows))
	case CommandTopRated:
		tvShows, err = tc.client.GetTopRatedTVShows(ctx, count)
		tc.logger.Info("Fetched %d top-rated TV shows", len(tvShows))
	case CommandOnTheAir:
		tvShows, err = tc.client.GetOnTheAirTVShows(ctx, count)
		tc.logger.Info("Fetched %d on-the-air TV shows", len(tvShows))
	default:
		return fmt.Errorf("unknown TV show list type: %s", listType)
	}

	if err != nil {
		return fmt.Errorf("failed to fetch %s TV shows: %w", listType, err)
	}

	if len(tvShows) == 0 {
		fmt.Printf("No %s TV shows found.\n", listType)
		return nil
	}

	options := BuildFormatOptions(flags)

	if options.Format == FormatTable {
		FormatTVSummary(os.Stdout, tvShows, listType, flags.OriginalTitle)
	}

	return FormatTVShows(os.Stdout, tvShows, options)
}

// HandleSearch handles TV show search commands.
func (tc *TVCommands) HandleSearch(
	ctx context.Context,
	args []string,
) error {
	return genericSearch(ctx, tc.client, tc.logger, args, SearchTypeTV)
}

// genericSearch handles both movie and TV search with shared logic.
func genericSearch(
	ctx context.Context,
	client *Client,
	logger *Logger,
	args []string,
	searchType SearchType,
) error {
	query, maxItems, format, originalTitle, noHeader, err := ParseSearchCommand(args, searchType)
	if err != nil {
		return err
	}

	// Perform search and get results
	resultCount, displayFunc, err := performSearch(ctx, client, logger, query, maxItems, searchType)
	if err != nil {
		return err
	}

	// Handle empty results
	if resultCount == 0 {
		return handleEmptySearchResults(query, searchType)
	}

	// Format and display results
	return displaySearchResults(
		displayFunc,
		format,
		originalTitle,
		noHeader,
		query,
		resultCount,
		searchType,
	)
}

// performSearch executes the search based on type and returns results.
func performSearch(
	ctx context.Context,
	client *Client,
	logger *Logger,
	query string,
	maxItems int,
	searchType SearchType,
) (int, func(io.Writer, FormatOptions) error, error) {
	if searchType == SearchTypeMovie {
		ShowSearchProgress(query, maxItems)

		movies, err := client.SearchMovies(ctx, query, maxItems)
		if err != nil {
			return 0, nil, fmt.Errorf("search failed: %w", err)
		}

		logger.Info("Search for \"%s\" returned %d movies", query, len(movies))

		displayFunc := func(w io.Writer, opts FormatOptions) error {
			return FormatMovies(w, movies, opts)
		}

		return len(movies), displayFunc, nil
	}

	// TV search
	ShowTVSearchProgress(query, maxItems)

	tvShows, err := client.SearchTVShows(ctx, query, maxItems)
	if err != nil {
		return 0, nil, fmt.Errorf("TV search failed: %w", err)
	}

	logger.Info("TV search for \"%s\" returned %d shows", query, len(tvShows))

	displayFunc := func(w io.Writer, opts FormatOptions) error {
		return FormatTVShows(w, tvShows, opts)
	}

	return len(tvShows), displayFunc, nil
}

// handleEmptySearchResults handles the case when no search results are found.
func handleEmptySearchResults(query string, searchType SearchType) error {
	if searchType == SearchTypeMovie {
		fmt.Printf("No movies found for \"%s\"\n\n", query)
		ShowSearchHelp()
		fmt.Printf("  - Include the release year\n")
	} else {
		fmt.Printf("No TV shows found for \"%s\"\n\n", query)
		ShowSearchHelp()
		fmt.Printf("  - Include the first air year\n")
	}
	return nil
}

// displaySearchResults formats and displays search results.
func displaySearchResults(
	displayFunc func(io.Writer, FormatOptions) error,
	format string,
	originalTitle, noHeader bool,
	query string,
	resultCount int,
	searchType SearchType,
) error {
	// Format and display
	options := FormatOptions{
		Format:           format,
		UseOriginalTitle: originalTitle,
		NoHeader:         noHeader,
		MaxWidth:         DefaultMaxWidth,
	}

	if options.Format == FormatTable {
		var itemType, titleType string
		if searchType == SearchTypeMovie {
			itemType = "movies"
			titleType = "titles"
		} else {
			itemType = "TV shows"
			titleType = "names"
		}

		if originalTitle {
			titleType = "original language " + titleType
		}
		fmt.Printf("Found %d %s for \"%s\" (%s)\n\n", resultCount, itemType, query, titleType)
	}

	return displayFunc(os.Stdout, options)
}
