// cmd/tmdb/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/alnah/tmdb-cli/internal"
)

var (
	version   = "2.0.0"
	buildDate = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Create context that cancels on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Load configuration
	config, err := internal.LoadConfig()
	if err != nil {
		if strings.Contains(err.Error(), "API key") {
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			fmt.Fprint(os.Stderr, internal.GetAPIKeyHelp())
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Create logger
	logger := internal.NewLogger(internal.ParseLogLevel(config.LogLevel))

	// Create TMDB client
	client := internal.NewClient(config)

	// Parse command line arguments
	if len(os.Args) < 2 {
		showUsage()
		return
	}

	command := os.Args[1]

	// Handle special commands first
	switch command {
	case "help", "--help", "-h":
		showUsage()
		return
	case "version", "--version", "-v":
		showVersion()
		return
	case "config":
		handleConfigCommand()
		return
	}

	// Handle movie commands
	switch command {
	case "popular", "pop":
		err = handleList(ctx, client, logger, "popular", os.Args[2:])
	case "top-rated", "top", "rated":
		err = handleList(ctx, client, logger, "top-rated", os.Args[2:])
	case "now-playing", "now", "playing":
		err = handleList(ctx, client, logger, "now-playing", os.Args[2:])
	case "upcoming", "soon":
		err = handleList(ctx, client, logger, "upcoming", os.Args[2:])
	case "search", "find":
		err = handleSearch(ctx, client, logger, os.Args[2:])
	case "discover", "filter":
		err = handleDiscover(ctx, client, logger, os.Args[2:])
	default:
		// Check if it looks like a search query
		if looksLikeSearch(command) {
			err = handleAutoSearch(ctx, client, logger, os.Args[1:])
		} else {
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
			showUsage()
			os.Exit(1)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func looksLikeSearch(input string) bool {
	// Contains spaces, letters, or looks like a movie title
	return strings.Contains(input, " ") ||
		(len(input) > 2 && strings.ContainsAny(input, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"))
}

func showUsage() {
	fmt.Print(`TMDB CLI - Movie Database Command Line Interface

USAGE:
    tmdb <command> [options]

COMMANDS:
    popular [count]              Show popular movies
    top-rated [count]            Show top-rated movies
    now-playing [count]          Show movies in theaters
    upcoming [count]             Show upcoming releases
    search <query>               Search for movies
    discover [options]           Discover movies with filters

    version                      Show version information
    config                       Show configuration help
    help                         Show this help

SHORTCUTS:
    tmdb "movie title"           Auto-search for movie
    tmdb pop 50                  50 popular movies
    tmdb top                     Top-rated movies
    tmdb now                     Now playing movies

OPTIONS:
    --format table|json|csv      Output format (default: table)
    --max-items N                Maximum items to show (default: 20)
    --original-title             Show original language titles
    --no-header                  Don't show table headers
    --verbose                    Enable verbose logging
    --debug                      Enable debug logging

SEARCH OPTIONS:
    --language LANG              Language filter (e.g., en, fr, es)
    --year YEAR                  Year filter
    --min-rating N               Minimum rating (0-10)
    --max-rating N               Maximum rating (0-10)

DISCOVER OPTIONS:
    --genres action,comedy       Include genres
    --exclude horror,thriller    Exclude genres
    --sort popularity|rating     Sort by field
    --order asc|desc             Sort order

EXAMPLES:
    tmdb popular                 # 20 popular movies
    tmdb "The Matrix"            # Search for The Matrix
    tmdb search "action 2023"    # Search action movies from 2023
    tmdb top-rated --format json # Top movies as JSON
    tmdb discover --genres action --year 2023 --min-rating 7
    tmdb now-playing --original-title

CONFIGURATION:
    Set TMDB_API_KEY environment variable or create config file.
    Run 'tmdb config' for detailed configuration help.

Get your API key from: https://www.themoviedb.org/settings/api
`)
}

func showVersion() {
	fmt.Printf("TMDB CLI %s\n", version)
	fmt.Printf("Build Date: %s\n", buildDate)
	fmt.Printf("Git Commit: %s\n", gitCommit)
}

func handleConfigCommand() {
	fmt.Print(internal.GetConfigHelp())
}

// Common flag parsing for all commands.
type CommonFlags struct {
	Format        string
	MaxItems      int
	OriginalTitle bool
	NoHeader      bool
	Verbose       bool
	Debug         bool
}

func parseCommonFlags(args []string) (CommonFlags, []string, error) {
	var flags CommonFlags

	fs := flag.NewFlagSet("common", flag.ContinueOnError)
	fs.StringVar(&flags.Format, "format", internal.FormatTable, "Output format")
	fs.IntVar(&flags.MaxItems, "max-items", 20, "Maximum items")
	fs.BoolVar(&flags.OriginalTitle, "original-title", false, "Show original titles")
	fs.BoolVar(&flags.NoHeader, "no-header", false, "No table headers")
	fs.BoolVar(&flags.Verbose, "verbose", false, "Verbose logging")
	fs.BoolVar(&flags.Debug, "debug", false, "Debug logging")

	// Find where flags end and positional args begin
	posArgs := make([]string, 0, len(args))
	var flagArgs []string

	for i, arg := range args {
		if strings.HasPrefix(arg, "-") {
			flagArgs = args[i:]
			break
		}
		posArgs = append(posArgs, arg)
	}

	if len(flagArgs) > 0 {
		if err := fs.Parse(flagArgs); err != nil {
			return flags, nil, err
		}
	}

	return flags, posArgs, nil
}

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

func buildFormatOptions(flags CommonFlags) internal.FormatOptions {
	return internal.FormatOptions{
		Format:           flags.Format,
		UseOriginalTitle: flags.OriginalTitle,
		NoHeader:         flags.NoHeader,
		MaxWidth:         120,
	}
}

func validateMaxItems(maxItems int) error {
	if maxItems < 1 || maxItems > 1000 {
		return fmt.Errorf("max-items must be between 1 and 1000, got %d", maxItems)
	}
	return nil
}

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

func showProgress(message string, maxItems int) {
	const progressThreshold = 40
	if maxItems > progressThreshold {
		fmt.Printf("%s...\n", message)
	}
}
