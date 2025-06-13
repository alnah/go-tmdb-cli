// cmd/tv_commands.go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alnah/tmdb-cli/internal"
)

func handleTVList(
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
	showProgress(fmt.Sprintf("Fetching %d %s TV shows", count, listType), count)

	// Fetch TV shows based on type
	var tvShows []internal.TVShow
	switch listType {
	case CommandPopular:
		tvShows, err = client.GetPopularTVShows(ctx, count)
		logger.Info("Fetched %d popular TV shows", len(tvShows))
	case CommandTopRated:
		tvShows, err = client.GetTopRatedTVShows(ctx, count)
		logger.Info("Fetched %d top-rated TV shows", len(tvShows))
	case CommandOnTheAir:
		tvShows, err = client.GetOnTheAirTVShows(ctx, count)
		logger.Info("Fetched %d on-the-air TV shows", len(tvShows))
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

	// Format and display
	options := buildFormatOptions(flags)

	if options.Format == internal.FormatTable {
		internal.FormatTVSummary(os.Stdout, tvShows, listType, flags.OriginalTitle)
	}

	return internal.FormatTVShows(os.Stdout, tvShows, options)
}

func handleTVSearch(
	ctx context.Context,
	client *internal.Client,
	logger *internal.Logger,
	args []string,
) error {
	return genericSearch(ctx, client, logger, args, SearchTypeTV)
}

func handleTVDiscover(
	ctx context.Context,
	client *internal.Client,
	logger *internal.Logger,
	args []string,
) error {
	// Create discover flag set with custom descriptions for TV
	fs, flags := createDiscoverFlagSet("tv discover")

	// Override year description
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Note: --year flag filters by first air date year for TV shows\n")
		fs.PrintDefaults()
	}

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

	showProgress("Discovering TV shows with filters", *flags.MaxItems)

	// Perform discovery
	tvShows, err := client.DiscoverTVShows(ctx, opts)
	if err != nil {
		return fmt.Errorf("TV discovery failed: %w", err)
	}

	logger.Info("TV discovery returned %d shows", len(tvShows))

	if len(tvShows) == 0 {
		fmt.Println("No TV shows found matching the criteria.")
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
		titleType := "names"
		if *flags.OriginalTitle {
			titleType = "original language names"
		}
		fmt.Printf(
			"Discovered %d TV shows matching your criteria (%s)\n\n",
			len(tvShows),
			titleType,
		)
	}

	return internal.FormatTVShows(os.Stdout, tvShows, options)
}
