// cmd/commands.go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alnah/tmdb-cli/internal"
)

func handleList(
	ctx context.Context,
	client *internal.Client,
	logger *internal.Logger,
	listType string,
	args []string,
) error {
	flags, posArgs, err := parseCommonFlags(args)
	if err != nil {
		return err
	}

	// Parse count from positional args
	count, err := parseCountArg(posArgs, flags.MaxItems)
	if err != nil {
		return err
	}

	if e := validateMaxItems(count); e != nil {
		return e
	}

	// Show progress for larger requests
	showProgress(fmt.Sprintf("Fetching %d %s movies", count, listType), count)

	// Fetch movies based on type
	var movies []internal.Movie
	switch listType {
	case CommandPopular:
		movies, err = client.GetPopularMovies(ctx, count)
		logger.Info("Fetched %d popular movies", len(movies))
	case CommandTopRated:
		movies, err = client.GetTopRatedMovies(ctx, count)
		logger.Info("Fetched %d top-rated movies", len(movies))
	case CommandNowPlaying:
		movies, err = client.GetNowPlayingMovies(ctx, count)
		logger.Info("Fetched %d now-playing movies", len(movies))
	case CommandUpcoming:
		movies, err = client.GetUpcomingMovies(ctx, count)
		logger.Info("Fetched %d upcoming movies", len(movies))
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
	options := buildFormatOptions(flags)

	if options.Format == internal.FormatTable {
		internal.FormatSummary(os.Stdout, movies, listType, flags.OriginalTitle)
	}

	return internal.FormatMovies(os.Stdout, movies, options)
}

func handleSearch(
	ctx context.Context,
	client *internal.Client,
	logger *internal.Logger,
	args []string,
) error {
	return genericSearch(ctx, client, logger, args, SearchTypeMovie)
}

func handleDiscover(
	ctx context.Context,
	client *internal.Client,
	logger *internal.Logger,
	args []string,
) error {
	// Create discover flag set
	fs, flags := createDiscoverFlagSet("discover")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := validateMaxItems(*flags.MaxItems); err != nil {
		return err
	}

	// Build search options
	opts, err := buildSearchOptions(flags)
	if err != nil {
		return err
	}

	// Validate options
	if err := internal.ValidateSearchOptions(&opts); err != nil {
		return err
	}

	showProgress("Discovering movies with filters", *flags.MaxItems)

	// Perform discovery
	movies, err := client.DiscoverMovies(ctx, opts)
	if err != nil {
		return fmt.Errorf("discovery failed: %w", err)
	}

	logger.Info("Discovery returned %d movies", len(movies))

	if len(movies) == 0 {
		fmt.Println("No movies found matching the criteria.")
		fmt.Println("Try adjusting your filters or removing some restrictions.")
		return nil
	}

	// Format and display
	options := internal.FormatOptions{
		Format:           *flags.Format,
		UseOriginalTitle: *flags.OriginalTitle,
		NoHeader:         *flags.NoHeader,
		MaxWidth:         120,
	}

	if options.Format == internal.FormatTable {
		titleType := "titles"
		if *flags.OriginalTitle {
			titleType = "original language titles"
		}
		fmt.Printf("Discovered %d movies matching your criteria (%s)\n\n", len(movies), titleType)
	}

	return internal.FormatMovies(os.Stdout, movies, options)
}
