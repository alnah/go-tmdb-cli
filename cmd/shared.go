// cmd/shared.go
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/alnah/tmdb-cli/internal"
)

// CommonFlags represents common command line flags used across different commands.
type CommonFlags struct {
	Format        string
	MaxItems      int
	OriginalTitle bool
	NoHeader      bool
	Verbose       bool
	Debug         bool
}

// parseSearchArgs parses search command arguments into query and flags.
func parseSearchArgs(args []string) (string, []string) {
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

	return query, flagArgs
}

// parseSearchFlags parses common search flags.
func parseSearchFlags(fs *flag.FlagSet, flagArgs []string) error {
	if len(flagArgs) > 0 {
		return fs.Parse(flagArgs)
	}
	return nil
}

// createSearchFlagSet creates a flag set for search commands.
func createSearchFlagSet(name string) (*flag.FlagSet, *string, *int, *bool, *bool) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	format := fs.String("format", internal.FormatTable, "Output format")
	maxItems := fs.Int("max-items", 20, "Maximum items")
	originalTitle := fs.Bool("original-title", false, "Show original titles")
	noHeader := fs.Bool("no-header", false, "No table headers")

	return fs, format, maxItems, originalTitle, noHeader
}

// DiscoverFlags contains all flags for discover commands.
type DiscoverFlags struct {
	Format        *string
	MaxItems      *int
	OriginalTitle *bool
	NoHeader      *bool
	Language      *string
	Year          *int
	MinRating     *float64
	MaxRating     *float64
	MinVotes      *int // New field
	MaxVotes      *int // New field
	Genres        *string
	ExcludeGenres *string
	SortBy        *string
	SortOrder     *string
}

// createDiscoverFlagSet creates a flag set for discover commands.
func createDiscoverFlagSet(name string) (*flag.FlagSet, *DiscoverFlags) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)

	flags := &DiscoverFlags{
		Format:        fs.String("format", internal.FormatTable, "Output format"),
		MaxItems:      fs.Int("max-items", 20, "Maximum items"),
		OriginalTitle: fs.Bool("original-title", false, "Show original titles"),
		NoHeader:      fs.Bool("no-header", false, "No table headers"),
		Language:      fs.String("language", "", "Language filter (e.g., en, fr)"),
		Year:          fs.Int("year", 0, "Year filter"),
		MinRating:     fs.Float64("min-rating", 0, "Minimum rating (0-10)"),
		MaxRating:     fs.Float64("max-rating", 0, "Maximum rating (0-10)"),
		MinVotes:      fs.Int("min-votes", 0, "Minimum vote count"), // New flag
		MaxVotes:      fs.Int("max-votes", 0, "Maximum vote count"), // New flag
		Genres:        fs.String("genres", "", "Include genres (comma-separated)"),
		ExcludeGenres: fs.String("exclude", "", "Exclude genres (comma-separated)"),
		SortBy: fs.String(
			"sort",
			"popularity",
			"Sort by field (popularity, rating, release_date, title, votes)",
		),
		SortOrder: fs.String("order", "desc", "Sort order (asc, desc)"),
	}

	return fs, flags
}

// parseGenreList parses a comma-separated genre list.
func parseGenreList(genresStr string) ([]int, error) {
	if genresStr == "" {
		return nil, nil
	}

	genreNames := strings.Split(genresStr, ",")
	for i, name := range genreNames {
		genreNames[i] = strings.TrimSpace(name)
	}

	return internal.ParseGenres(genreNames)
}

// buildSearchOptions builds search options from discover flags.
func buildSearchOptions(flags *DiscoverFlags) (internal.SearchOptions, error) {
	opts := internal.SearchOptions{
		MaxItems:  *flags.MaxItems,
		Language:  *flags.Language,
		Year:      *flags.Year,
		MinRating: *flags.MinRating,
		MaxRating: *flags.MaxRating,
		MinVotes:  *flags.MinVotes, // New field
		MaxVotes:  *flags.MaxVotes, // New field
		SortBy:    *flags.SortBy,
		SortOrder: *flags.SortOrder,
	}

	// Parse include genres
	if *flags.Genres != "" {
		genreIDs, err := parseGenreList(*flags.Genres)
		if err != nil {
			return opts, err
		}
		opts.IncludeGenres = genreIDs
	}

	// Parse exclude genres
	if *flags.ExcludeGenres != "" {
		genreIDs, err := parseGenreList(*flags.ExcludeGenres)
		if err != nil {
			return opts, err
		}
		opts.ExcludeGenres = genreIDs
	}

	return opts, nil
}

// showSearchHelp displays search help hints.
func showSearchHelp() {
	fmt.Println("Try:")
	fmt.Printf("  - Check spelling and try again\n")
	fmt.Printf("  - Use fewer, more common words\n")
	fmt.Printf("  - Try the original language title\n")
}

// showProgress displays progress message for large requests.
func showProgress(message string, maxItems int) {
	if maxItems > progressThreshold {
		fmt.Printf("%s...\n", message)
	}
}

// SearchType represents the type of search (movie or TV).
type SearchType int

const (
	SearchTypeMovie SearchType = iota
	SearchTypeTV
)

// genericSearch handles both movie and TV search with shared logic.
func genericSearch(
	ctx context.Context,
	client *internal.Client,
	logger *internal.Logger,
	args []string,
	searchType SearchType,
) error {
	if len(args) == 0 {
		return fmt.Errorf("search requires a query argument")
	}

	// Create search flag set
	flagSetName := "search"
	if searchType == SearchTypeTV {
		flagSetName = "tv search"
	}
	fs, format, maxItems, originalTitle, noHeader := createSearchFlagSet(flagSetName)

	// Parse query and flags
	query, flagArgs := parseSearchArgs(args)

	if query == "" {
		return fmt.Errorf("search query cannot be empty")
	}

	// Parse flags
	if err := parseSearchFlags(fs, flagArgs); err != nil {
		return err
	}

	if err := validateMaxItems(*maxItems); err != nil {
		return err
	}

	// Clean up quoted query
	query = strings.Trim(query, `"'`)

	// Perform search based on type
	var err error
	var resultCount int
	var displayFunc func(io.Writer, internal.FormatOptions) error

	if searchType == SearchTypeMovie {
		showProgress(fmt.Sprintf("Searching for \"%s\"", query), *maxItems)

		movies, searchErr := client.SearchMovies(ctx, query, *maxItems)
		err = searchErr
		resultCount = len(movies)

		logger.Info("Search for \"%s\" returned %d movies", query, resultCount)

		displayFunc = func(w io.Writer, opts internal.FormatOptions) error {
			return internal.FormatMovies(w, movies, opts)
		}

		if resultCount == 0 {
			fmt.Printf("No movies found for \"%s\"\n\n", query)
			showSearchHelp()
			fmt.Printf("  - Include the release year\n")
			return nil
		}
	} else {
		showProgress(fmt.Sprintf("Searching TV shows for \"%s\"", query), *maxItems)

		tvShows, searchErr := client.SearchTVShows(ctx, query, *maxItems)
		err = searchErr
		resultCount = len(tvShows)

		logger.Info("TV search for \"%s\" returned %d shows", query, resultCount)

		displayFunc = func(w io.Writer, opts internal.FormatOptions) error {
			return internal.FormatTVShows(w, tvShows, opts)
		}

		if resultCount == 0 {
			fmt.Printf("No TV shows found for \"%s\"\n\n", query)
			showSearchHelp()
			fmt.Printf("  - Include the first air year\n")
			return nil
		}
	}

	if err != nil {
		if searchType == SearchTypeMovie {
			return fmt.Errorf("search failed: %w", err)
		}
		return fmt.Errorf("TV search failed: %w", err)
	}

	// Format and display
	options := internal.FormatOptions{
		Format:           *format,
		UseOriginalTitle: *originalTitle,
		NoHeader:         *noHeader,
		MaxWidth:         120,
	}

	if options.Format == internal.FormatTable {
		var itemType, titleType string
		if searchType == SearchTypeMovie {
			itemType = "movies"
			titleType = "titles"
		} else {
			itemType = "TV shows"
			titleType = "names"
		}

		if *originalTitle {
			titleType = "original language " + titleType
		}
		fmt.Printf("Found %d %s for \"%s\" (%s)\n\n", resultCount, itemType, query, titleType)
	}

	return displayFunc(os.Stdout, options)
}

// validateMaxItems validates the max items parameter.
func validateMaxItems(maxItems int) error {
	if maxItems < 1 || maxItems > 1000 {
		return fmt.Errorf("max-items must be between 1 and 1000, got %d", maxItems)
	}
	return nil
}

// handleAutoSearch performs auto-search for movie titles.
func handleAutoSearch(
	ctx context.Context,
	client *internal.Client,
	logger *internal.Logger,
	args []string,
) error {
	query := strings.Join(args, " ")
	query = strings.Trim(query, `"'`) // Remove quotes if present

	logger.Info("Auto-searching for: %s", query)

	movies, err := client.SearchMovies(ctx, query, 20)
	if err != nil {
		return err
	}

	if len(movies) == 0 {
		fmt.Printf("No movies found for \"%s\"\n", query)
		return nil
	}

	fmt.Printf("Found %d movies for \"%s\"\n\n", len(movies), query)

	options := internal.FormatOptions{
		Format:   internal.FormatTable,
		NoHeader: false,
		MaxWidth: 120,
	}

	return internal.FormatMovies(os.Stdout, movies, options)
}

// buildFormatOptions builds format options from common flags.
func buildFormatOptions(flags CommonFlags) internal.FormatOptions {
	return internal.FormatOptions{
		Format:           flags.Format,
		UseOriginalTitle: flags.OriginalTitle,
		NoHeader:         flags.NoHeader,
		MaxWidth:         120,
	}
}

// parseCountArg parses count argument from positional args.
func parseCountArg(args []string, defaultCount int) (int, error) {
	if len(args) == 0 {
		return defaultCount, nil
	}

	count, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("invalid count '%s', must be a number", args[0])
	}

	if err := validateMaxItems(count); err != nil {
		return 0, err
	}

	return count, nil
}

// looksLikeSearch determines if input looks like a movie search query.
func looksLikeSearch(input string) bool {
	// Contains spaces, letters, or looks like a movie title
	return strings.Contains(input, " ") ||
		(len(input) > 2 && strings.ContainsAny(input, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"))
}
