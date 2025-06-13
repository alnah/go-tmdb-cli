// Package main provides the entry point for the TMDB CLI application.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
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

	// Dispatch command to appropriate handler
	err = dispatchCommand(ctx, client, logger, command, os.Args[2:])
	if err != nil {
		// Only show usage on unknown command errors
		if err.Error() == fmt.Sprintf("unknown command: %s", command) {
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			showUsage()
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}

func showVersion() {
	fmt.Printf("TMDB CLI %s\n", version)
	fmt.Printf("Build Date: %s\n", buildDate)
	fmt.Printf("Git Commit: %s\n", gitCommit)
}

func handleConfigCommand() {
	fmt.Print(internal.GetConfigHelp())
}

// CommonFlags represents common command line flags used across different commands.
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

// dispatchCommand routes commands to appropriate handlers.
func dispatchCommand(
	ctx context.Context,
	client *internal.Client,
	logger *internal.Logger,
	command string,
	args []string,
) error {
	dispatcher := internal.NewDispatcher(client, logger)
	return dispatcher.DispatchCommand(ctx, command, args)
}

func showUsage() {
	fmt.Printf("TMDB CLI %s\n\n", version)
	fmt.Println("A command-line interface for The Movie Database (TMDB)")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  tmdb-cli [command] [arguments] [flags]")
	fmt.Println("")
	fmt.Println("COMMANDS:")
	fmt.Println("  popular, pop              Get popular movies")
	fmt.Println("  top-rated, top, rated     Get top-rated movies")
	fmt.Println("  now-playing, playing      Get now playing movies")
	fmt.Println("  upcoming                  Get upcoming movies")
	fmt.Println("  search [query]            Search for movies")
	fmt.Println("  discover                  Discover movies with filters")
	fmt.Println("  tv-popular                Get popular TV shows")
	fmt.Println("  tv-top-rated              Get top-rated TV shows")
	fmt.Println("  tv-airing-today           Get TV shows airing today")
	fmt.Println("  tv-on-the-air             Get TV shows on the air")
	fmt.Println("  tv-search [query]         Search for TV shows")
	fmt.Println("  config                    Show configuration help")
	fmt.Println("  version, --version, -v    Show version information")
	fmt.Println("  help, --help, -h          Show this help message")
	fmt.Println("")
	fmt.Println("GLOBAL FLAGS:")
	fmt.Println("  --format string           Output format (table, json, csv) (default \"table\")")
	fmt.Println("  --max-items int           Maximum number of items to display (default 20)")
	fmt.Println("  --original-title          Show original titles instead of localized titles")
	fmt.Println("  --no-header               Don't show table headers")
	fmt.Println("  --verbose                 Enable verbose logging")
	fmt.Println("  --debug                   Enable debug logging")
	fmt.Println("")
	fmt.Println("EXAMPLES:")
	fmt.Println("  tmdb-cli popular 10")
	fmt.Println("  tmdb-cli search \"inception\" --format json")
	fmt.Println("  tmdb-cli discover --year 2023 --min-rating 7.0")
	fmt.Println("  tmdb-cli tv-popular --max-items 5")
	fmt.Println("")
	fmt.Println("For more information, visit: https://github.com/alnah/tmdb-cli")
}
