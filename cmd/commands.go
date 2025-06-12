// cmd/tmdb/commands.go
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

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
	case "popular":
		movies, err = client.GetPopularMovies(ctx, count)
		logger.Info("Fetched %d popular movies", len(movies))
	case "top-rated":
		movies, err = client.GetTopRatedMovies(ctx, count)
		logger.Info("Fetched %d top-rated movies", len(movies))
	case "now-playing":
		movies, err = client.GetNowPlayingMovies(ctx, count)
		logger.Info("Fetched %d now-playing movies", len(movies))
	case "upcoming":
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
	if len(args) == 0 {
		return fmt.Errorf("search requires a query argument")
	}

	// Parse search-specific flags
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	format := fs.String("format", internal.FormatTable, "Output format")
	maxItems := fs.Int("max-items", 20, "Maximum items")
	originalTitle := fs.Bool("original-title", false, "Show original titles")
	noHeader := fs.Bool("no-header", false, "No table headers")

	// Find the query and flags
	var query string
	var flagArgs []string

	for i, arg := range args {
		if strings.HasPrefix(arg, "-") {
			flagArgs = args[i:]
			break
		}
		if query == "" {
			query = arg
		} else {
			query += " " + arg
		}
	}

	if query == "" {
		return fmt.Errorf("search query cannot be empty")
	}

	// Parse flags
	if len(flagArgs) > 0 {
		if err := fs.Parse(flagArgs); err != nil {
			return err
		}
	}

	if err := validateMaxItems(*maxItems); err != nil {
		return err
	}

	// Clean up quoted query
	query = strings.Trim(query, `"'`)

	showProgress(fmt.Sprintf("Searching for \"%s\"", query), *maxItems)

	// Perform search
	movies, err := client.SearchMovies(ctx, query, *maxItems)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	logger.Info("Search for \"%s\" returned %d movies", query, len(movies))

	if len(movies) == 0 {
		fmt.Printf("No movies found for \"%s\"\n\n", query)
		fmt.Println("Try:")
		fmt.Printf("  - Check spelling and try again\n")
		fmt.Printf("  - Use fewer, more common words\n")
		fmt.Printf("  - Try the original language title\n")
		fmt.Printf("  - Include the release year\n")
		return nil
	}

	// Format and display
	options := internal.FormatOptions{
		Format:           *format,
		UseOriginalTitle: *originalTitle,
		NoHeader:         *noHeader,
		MaxWidth:         120,
	}

	if options.Format == internal.FormatTable {
		titleType := "titles"
		if *originalTitle {
			titleType = "original language titles"
		}
		fmt.Printf("Found %d movies for \"%s\" (%s)\n\n", len(movies), query, titleType)
	}

	return internal.FormatMovies(os.Stdout, movies, options)
}

func handleDiscover(
	ctx context.Context,
	client *internal.Client,
	logger *internal.Logger,
	args []string,
) error {
	// Parse discover-specific flags
	fs := flag.NewFlagSet("discover", flag.ContinueOnError)
	format := fs.String("format", internal.FormatTable, "Output format")
	maxItems := fs.Int("max-items", 20, "Maximum items")
	originalTitle := fs.Bool("original-title", false, "Show original titles")
	noHeader := fs.Bool("no-header", false, "No table headers")

	language := fs.String("language", "", "Language filter (e.g., en, fr)")
	year := fs.Int("year", 0, "Year filter")
	minRating := fs.Float64("min-rating", 0, "Minimum rating (0-10)")
	maxRating := fs.Float64("max-rating", 0, "Maximum rating (0-10)")
	genres := fs.String("genres", "", "Include genres (comma-separated)")
	excludeGenres := fs.String("exclude", "", "Exclude genres (comma-separated)")
	sortBy := fs.String("sort", "popularity", "Sort by (popularity, rating, release_date)")
	sortOrder := fs.String("order", "desc", "Sort order (asc, desc)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := validateMaxItems(*maxItems); err != nil {
		return err
	}

	// Build search options
	opts := internal.SearchOptions{
		MaxItems:  *maxItems,
		Language:  *language,
		Year:      *year,
		MinRating: *minRating,
		MaxRating: *maxRating,
		SortBy:    *sortBy,
		SortOrder: *sortOrder,
	}

	// Parse genres
	if *genres != "" {
		genreNames := strings.Split(*genres, ",")
		for i, name := range genreNames {
			genreNames[i] = strings.TrimSpace(name)
		}
		genreIDs, err := internal.ParseGenres(genreNames)
		if err != nil {
			return err
		}
		opts.IncludeGenres = genreIDs
	}

	// Parse exclude genres
	if *excludeGenres != "" {
		genreNames := strings.Split(*excludeGenres, ",")
		for i, name := range genreNames {
			genreNames[i] = strings.TrimSpace(name)
		}
		genreIDs, err := internal.ParseGenres(genreNames)
		if err != nil {
			return err
		}
		opts.ExcludeGenres = genreIDs
	}

	// Validate options
	if err := internal.ValidateSearchOptions(&opts); err != nil {
		return err
	}

	showProgress("Discovering movies with filters", *maxItems)

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
		Format:           *format,
		UseOriginalTitle: *originalTitle,
		NoHeader:         *noHeader,
		MaxWidth:         120,
	}

	if options.Format == internal.FormatTable {
		titleType := "titles"
		if *originalTitle {
			titleType = "original language titles"
		}
		fmt.Printf("Discovered %d movies matching your criteria (%s)\n\n", len(movies), titleType)
	}

	return internal.FormatMovies(os.Stdout, movies, options)
}
