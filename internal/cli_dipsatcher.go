// Package internal provides core functionality for the TMDB CLI application.
package internal

import (
	"context"
	"fmt"
)

// CLIDispatcher handles command routing and execution.
type CLIDispatcher struct {
	movieCommands  *MovieCommands
	tvCommands     *TVCommands
	configCommands *ConfigCommands
	client         TMDBClient
	logger         *Logger
}

// NewCLIDispatcher creates a new command dispatcher.
func NewCLIDispatcher(client TMDBClient, logger *Logger) *CLIDispatcher {
	return &CLIDispatcher{
		movieCommands:  NewMovieCommands(client, logger),
		tvCommands:     NewTVCommands(client, logger),
		configCommands: NewConfigCommands(),
		client:         client,
		logger:         logger,
	}
}

// DispatchCommand routes commands to appropriate handlers.
func (d *CLIDispatcher) DispatchCommand(
	ctx context.Context,
	command string,
	args []string,
	version, buildDate, gitCommit string,
) error {
	// Handle special commands
	switch command {
	case "help", "--help", "-h":
		ShowUsage(version)
		return nil
	case "version", "--version", "-v":
		ShowVersion(version, buildDate, gitCommit)
		return nil
	case "config":
		return d.configCommands.HandleConfig()
	}

	// Handle movie commands
	switch command {
	case CommandPopular, "pop":
		return d.movieCommands.HandleList(ctx, CommandPopular, args)
	case CommandTopRated, "top", "rated":
		return d.movieCommands.HandleList(ctx, CommandTopRated, args)
	case CommandNowPlaying, "now", "playing":
		return d.movieCommands.HandleList(ctx, CommandNowPlaying, args)
	case CommandUpcoming, "soon":
		return d.movieCommands.HandleList(ctx, CommandUpcoming, args)
	case CommandSearch, "find":
		return d.movieCommands.HandleSearch(ctx, args)
	case "tv":
		return d.tvCommands.Dispatch(ctx, args)
	default:
		// Check if it looks like a search query
		if LooksLikeSearch(command) {
			// Reconstruct full args for auto-search
			fullArgs := append([]string{command}, args...)
			return d.movieCommands.HandleAutoSearch(ctx, fullArgs)
		}
		return fmt.Errorf("unknown command: %s", command)
	}
}
