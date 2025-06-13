// cmd/dispatcher.go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alnah/tmdb-cli/internal"
)

// dispatchCommand routes commands to appropriate handlers.
func dispatchCommand(
	ctx context.Context,
	client *internal.Client,
	logger *internal.Logger,
	command string,
	args []string,
) error {
	// Handle special commands
	switch command {
	case "help", "--help", "-h":
		showUsage()
		return nil
	case "version", "--version", "-v":
		showVersion()
		return nil
	case "config":
		handleConfigCommand()
		return nil
	}

	// Handle movie commands
	// Handle movie commands
	switch command {
	case "popular", "pop":
		return handleList(ctx, client, logger, CommandPopular, args)
	case "top-rated", "top", "rated":
		return handleList(ctx, client, logger, CommandTopRated, args)
	case "now-playing", "now", "playing":
		return handleList(ctx, client, logger, CommandNowPlaying, args)
	case "upcoming", "soon":
		return handleList(ctx, client, logger, CommandUpcoming, args)
	case CommandSearch, "find":
		return handleSearch(ctx, client, logger, args)
	case CommandDiscover, "filter":
		return handleDiscover(ctx, client, logger, args)
	case CommandTV:
		return dispatchTVCommand(ctx, client, logger, args)
	default:
		// Check if it looks like a search query
		if looksLikeSearch(command) {
			// Reconstruct full args for auto-search
			fullArgs := append([]string{command}, args...)
			return handleAutoSearch(ctx, client, logger, fullArgs)
		}
		return fmt.Errorf("unknown command: %s", command)
	}
}

// dispatchTVCommand handles TV show subcommands.
func dispatchTVCommand(
	ctx context.Context,
	client *internal.Client,
	logger *internal.Logger,
	args []string,
) error {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "TV command requires a subcommand\n\n")
		showTVUsage()
		return fmt.Errorf("missing TV subcommand")
	}

	tvCommand := args[0]
	tvArgs := args[1:]

	switch tvCommand {
	case "popular", "pop":
		return handleTVList(ctx, client, logger, CommandPopular, tvArgs)
	case "top-rated", "top", "rated":
		return handleTVList(ctx, client, logger, CommandTopRated, tvArgs)
	case "on-the-air", "air", "airing":
		return handleTVList(ctx, client, logger, CommandOnTheAir, tvArgs)
	case "search", "find":
		return handleTVSearch(ctx, client, logger, tvArgs)
	case "discover", "filter":
		return handleTVDiscover(ctx, client, logger, tvArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown TV subcommand: %s\n\n", tvCommand)
		showTVUsage()
		return fmt.Errorf("unknown TV subcommand: %s", tvCommand)
	}
}
