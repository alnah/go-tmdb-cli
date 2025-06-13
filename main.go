// Package main provides the entry point for the TMDB CLI application.
package main

import (
	"context"
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
		internal.ShowUsage(version)
		return
	}

	command := os.Args[1]

	// Create dispatcher and handle command
	dispatcher := internal.NewCLIDispatcher(client, logger)
	err = dispatcher.DispatchCommand(ctx, command, os.Args[2:], version, buildDate, gitCommit)
	if err != nil {
		// Only show usage on unknown command errors
		if err.Error() == fmt.Sprintf("unknown command: %s", command) {
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			internal.ShowUsage(version)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}
