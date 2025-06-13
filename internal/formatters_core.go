// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"fmt"
	"io"
)

// Format constants.
const (
	FormatJSON  = "json"
	FormatCSV   = "csv"
	FormatTable = "table"
	DefaultNA   = "N/A"
)

// Table column configuration constants.
const (
	TitleColumnWidth  = 40
	GenresColumnWidth = 25
	NumberColWidth    = 4
	TitleColWidth     = 45
	YearColWidth      = 6
	RatingColWidth    = 8
	VotesColWidth     = 8
	GenresColWidth    = 25
)

// Table column numbers for configuration.
const (
	ColNumber = iota + 1
	ColTitle
	ColYear
	ColRating
	ColVotes
	ColGenres
)

// FormatOptions controls output formatting.
type FormatOptions struct {
	Format           string
	UseOriginalTitle bool
	NoHeader         bool
	MaxWidth         int
}

// MediaFormatter defines the interface for formatting different media types.
type MediaFormatter interface {
	GetTitle(useOriginal bool) string
	GetOriginalTitle() string
	GetFields() []string
	GetEmptyMessage() string
	GetHeaderName() string
}

// FormatMovies formats movies according to the specified format.
func FormatMovies(writer io.Writer, movies []Movie, options FormatOptions) error {
	switch options.Format {
	case FormatJSON:
		return formatMoviesJSON(writer, movies)
	case FormatCSV:
		formatters := make([]MediaFormatter, len(movies))
		for i, movie := range movies {
			formatters[i] = MovieFormatter{Movie: movie}
		}
		return formatCSVGeneric(writer, formatters, getMovieHeaders(), options)
	case FormatTable:
		return formatMoviesTable(writer, movies, options)
	default:
		return fmt.Errorf("unsupported format: %s", options.Format)
	}
}

// FormatTVShows formats TV shows according to the specified format.
func FormatTVShows(writer io.Writer, shows []TVShow, options FormatOptions) error {
	switch options.Format {
	case FormatJSON:
		return formatTVShowsJSON(writer, shows)
	case FormatCSV:
		formatters := make([]MediaFormatter, len(shows))
		for i, show := range shows {
			formatters[i] = TVShowFormatter{TVShow: show}
		}
		return formatCSVGeneric(writer, formatters, getTVShowHeaders(), options)
	case FormatTable:
		return formatTVShowsTable(writer, shows, options)
	default:
		return fmt.Errorf("unsupported format: %s", options.Format)
	}
}

// FormatSearchResult formats search results with pagination info.
func FormatSearchResult(writer io.Writer, result SearchResult, options FormatOptions) error {
	switch options.Format {
	case FormatJSON:
		return formatSearchResultJSON(writer, result)
	case FormatCSV:
		formatters := make([]MediaFormatter, len(result.Movies))
		for i, movie := range result.Movies {
			formatters[i] = MovieFormatter{Movie: movie}
		}
		return formatCSVGeneric(writer, formatters, getMovieHeaders(), options)
	case FormatTable:
		return formatSearchResultTable(writer, result, options)
	default:
		return fmt.Errorf("unsupported format: %s", options.Format)
	}
}

// FormatTVSearchResult formats TV search results with pagination info.
func FormatTVSearchResult(writer io.Writer, result TVSearchResult, options FormatOptions) error {
	switch options.Format {
	case FormatJSON:
		return formatTVSearchResultJSON(writer, result)
	case FormatCSV:
		formatters := make([]MediaFormatter, len(result.TVShows))
		for i, show := range result.TVShows {
			formatters[i] = TVShowFormatter{TVShow: show}
		}
		return formatCSVGeneric(writer, formatters, getTVShowHeaders(), options)
	case FormatTable:
		return formatTVSearchResultTable(writer, result, options)
	default:
		return fmt.Errorf("unsupported format: %s", options.Format)
	}
}

// FormatError formats error messages for output.
func FormatError(writer io.Writer, err error) error {
	_, writeErr := fmt.Fprintf(writer, "Error: %s\n", err.Error())
	return writeErr
}

// ValidateFormat validates format option.
func ValidateFormat(format string) bool {
	switch format {
	case FormatTable, FormatJSON, FormatCSV:
		return true
	default:
		return false
	}
}

// SupportedFormats returns supported formats.
func SupportedFormats() []string {
	return []string{FormatTable, FormatJSON, FormatCSV}
}
