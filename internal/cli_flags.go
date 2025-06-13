// Package internal provides core functionality for the TMDB CLI application.
package internal

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
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

// SearchType represents the type of search (movie or TV).
type SearchType int

const (
	// SearchTypeMovie represents movie search.
	SearchTypeMovie SearchType = iota
	// SearchTypeTV represents TV show search.
	SearchTypeTV
)

// ParseCommonFlags parses common command line flags from arguments.
func ParseCommonFlags(args []string) (CommonFlags, []string, error) {
	var flags CommonFlags

	fs := flag.NewFlagSet("common", flag.ContinueOnError)
	fs.StringVar(&flags.Format, "format", FormatTable, "Output format")
	fs.IntVar(&flags.MaxItems, "max-items", DefaultMaxItems, "Maximum items")
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

// BuildFormatOptions builds format options from common flags.
func BuildFormatOptions(flags CommonFlags) FormatOptions {
	return FormatOptions{
		Format:           flags.Format,
		UseOriginalTitle: flags.OriginalTitle,
		NoHeader:         flags.NoHeader,
		MaxWidth:         DefaultMaxWidth,
	}
}

// ParseCountArg parses count argument from positional args.
func ParseCountArg(args []string, defaultCount int) (int, error) {
	if len(args) == 0 {
		return defaultCount, nil
	}

	count, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("invalid count '%s', must be a number", args[0])
	}

	if err := ValidateMaxItems(count); err != nil {
		return 0, err
	}

	return count, nil
}

// ValidateMaxItems validates the max items parameter.
func ValidateMaxItems(maxItems int) error {
	if maxItems < 1 || maxItems > 1000 {
		return fmt.Errorf("max-items must be between 1 and 1000, got %d", maxItems)
	}
	return nil
}

// ParseSearchCommand parses search command arguments and flags.
func ParseSearchCommand(
	args []string,
	searchType SearchType,
) (string, int, string, bool, bool, error) {
	if len(args) == 0 {
		return "", 0, "", false, false, fmt.Errorf("search requires a query argument")
	}

	// Parse search arguments
	query, flagArgs := parseSearchArgs(args)
	if query == "" {
		return "", 0, "", false, false, fmt.Errorf("search query cannot be empty")
	}

	// Create search flag set
	flagSetName := CommandSearch
	if searchType == SearchTypeTV {
		flagSetName = "tv search"
	}
	fs, format, maxItems, originalTitle, noHeader := createSearchFlagSet(flagSetName)

	// Parse flags
	if err := parseSearchFlags(fs, flagArgs); err != nil {
		return "", 0, "", false, false, err
	}

	if err := ValidateMaxItems(*maxItems); err != nil {
		return "", 0, "", false, false, err
	}

	// Clean up quoted query
	query = strings.Trim(query, `"'`)

	return query, *maxItems, *format, *originalTitle, *noHeader, nil
}

// parseSearchArgs separates query from flag arguments.
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

// parseSearchFlags parses search-specific flags.
func parseSearchFlags(fs *flag.FlagSet, flagArgs []string) error {
	if len(flagArgs) > 0 {
		return fs.Parse(flagArgs)
	}
	return nil
}

// createSearchFlagSet creates a flag set for search commands.
func createSearchFlagSet(name string) (*flag.FlagSet, *string, *int, *bool, *bool) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	format := fs.String("format", FormatTable, "Output format")
	maxItems := fs.Int("max-items", DefaultMaxItems, "Maximum items")
	originalTitle := fs.Bool("original-title", false, "Show original titles")
	noHeader := fs.Bool("no-header", false, "No table headers")

	return fs, format, maxItems, originalTitle, noHeader
}
