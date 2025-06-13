// cmd/tmdb/main.go
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
