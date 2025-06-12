// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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

// FormatOptions controls output formatting.
type FormatOptions struct {
	Format           string
	UseOriginalTitle bool
	NoHeader         bool
	MaxWidth         int
}

// FormatMovies formats movies according to the specified format.
func FormatMovies(writer io.Writer, movies []Movie, options FormatOptions) error {
	switch options.Format {
	case FormatJSON:
		return formatJSON(writer, movies)
	case FormatCSV:
		return formatCSV(writer, movies, options)
	case FormatTable:
		return formatTable(writer, movies, options)
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
		return formatCSV(writer, result.Movies, options)
	case FormatTable:
		return formatSearchResultTable(writer, result, options)
	default:
		return fmt.Errorf("unsupported format: %s", options.Format)
	}
}

// JSON formatting.
func formatJSON(writer io.Writer, movies []Movie) error {
	output := struct {
		Movies []Movie `json:"movies"`
		Count  int     `json:"count"`
	}{
		Movies: movies,
		Count:  len(movies),
	}

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(output)
}

func formatSearchResultJSON(writer io.Writer, result SearchResult) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(result)
}

// CSV formatting.
func formatCSV(writer io.Writer, movies []Movie, options FormatOptions) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header if not disabled
	if !options.NoHeader {
		headers := []string{
			"ID", "Title", "Original Title", "Year", "Rating", "Votes",
			"Popularity", "Genres", "Overview", "Language", "Adult",
		}
		if err := csvWriter.Write(headers); err != nil {
			return fmt.Errorf("write CSV header: %w", err)
		}
	}

	// Write movie rows
	for _, movie := range movies {
		row := []string{
			strconv.Itoa(movie.ID),
			movie.Title,
			movie.OriginalTitle,
			strconv.Itoa(movie.Year),
			fmt.Sprintf("%.1f", movie.Rating),
			strconv.Itoa(movie.Votes),
			fmt.Sprintf("%.2f", movie.Popularity),
			movie.Genres,
			movie.Overview,
			movie.Language,
			strconv.FormatBool(movie.Adult),
		}

		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("write CSV row: %w", err)
		}
	}

	return nil
}

// Table formatting.
func formatTable(writer io.Writer, movies []Movie, options FormatOptions) error {
	if len(movies) == 0 {
		_, err := fmt.Fprintln(writer, "No movies found.")
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(writer)
	t.SetStyle(table.StyleDefault)

	// Add headers if not disabled
	if !options.NoHeader {
		headers := []string{"#", "Title", "Year", "Rating", "Votes", "Genres"}
		headerRow := make(table.Row, len(headers))
		for i, h := range headers {
			headerRow[i] = h
		}
		t.AppendHeader(headerRow)
	}

	// Add movie rows
	for i, movie := range movies {
		title := formatMovieTitle(movie, options.UseOriginalTitle, TitleColumnWidth)
		rating := FormatRating(movie.Rating, movie.Votes)
		votes := FormatVotes(movie.Votes)
		genres := FormatGenres(movie.Genres, GenresColumnWidth)

		row := table.Row{
			i + 1,
			addMovieIndicators(movie, title),
			formatYear(movie.Year),
			rating,
			votes,
			genres,
		}
		t.AppendRow(row)
	}

	// Configure column properties
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignCenter, WidthMax: NumberColWidth}, // #
		{Number: 2, Align: text.AlignLeft, WidthMax: TitleColWidth},    // Title
		{Number: 3, Align: text.AlignCenter, WidthMax: YearColWidth},   // Year
		{Number: 4, Align: text.AlignRight, WidthMax: RatingColWidth},  // Rating
		{Number: 5, Align: text.AlignRight, WidthMax: VotesColWidth},   // Votes
		{Number: 6, Align: text.AlignLeft, WidthMax: GenresColWidth},   // Genres
	})

	t.Render()
	return nil
}

func formatSearchResultTable(writer io.Writer, result SearchResult, options FormatOptions) error {
	if len(result.Movies) == 0 {
		_, err := fmt.Fprintln(writer, "No movies found.")
		return err
	}

	// Show pagination info first
	if result.TotalPages > 1 {
		if _, err := fmt.Fprintf(writer, "Page %d of %d (%d total results)\n\n",
			result.Page, result.TotalPages, result.TotalResults); err != nil {
			return err
		}
	}

	return formatTable(writer, result.Movies, options)
}

// Helper functions for table formatting.
func formatMovieTitle(movie Movie, useOriginal bool, maxWidth int) string {
	var title string
	if useOriginal && movie.OriginalTitle != "" && movie.OriginalTitle != movie.Title {
		title = movie.OriginalTitle
		// Add English title in parentheses if space allows
		fullTitle := title + " (" + movie.Title + ")"
		if len(fullTitle) <= maxWidth {
			title = fullTitle
		}
	} else {
		title = movie.Title
		// Add original title in parentheses if different and space allows
		if movie.OriginalTitle != "" && movie.OriginalTitle != movie.Title {
			fullTitle := title + " (" + movie.OriginalTitle + ")"
			if len(fullTitle) <= maxWidth {
				title = fullTitle
			}
		}
	}

	// Truncate if too long
	if len(title) > maxWidth {
		title = title[:maxWidth-3] + "..."
	}

	return title
}

func addMovieIndicators(movie Movie, title string) string {
	var indicators []string

	if IsRecent(movie.Year) {
		indicators = append(indicators, "[NEW]")
	}
	if IsHighlyRated(movie.Rating, movie.Votes) {
		indicators = append(indicators, "[TOP]")
	}
	if IsPopular(movie.Popularity, movie.Votes) {
		indicators = append(indicators, "[HOT]")
	}

	if len(indicators) > 0 {
		return strings.Join(indicators, "") + " " + title
	}
	return title
}

func formatYear(year int) string {
	if year == 0 {
		return DefaultNA
	}
	return strconv.Itoa(year)
}

// Error formatting.
func FormatError(writer io.Writer, err error) error {
	_, writeErr := fmt.Fprintf(writer, "Error: %s\n", err.Error())
	return writeErr
}

// Summary formatting for different commands.
func FormatSummary(writer io.Writer, movies []Movie, command string, useOriginal bool) {
	if len(movies) == 0 {
		return
	}

	titleType := "titles"
	if useOriginal {
		titleType = "original language titles"
	}

	switch command {
	case "popular":
		if _, err := fmt.Fprintf(writer, "Showing %d popular movies (%s)\n\n", len(movies), titleType); err != nil {
			// Error is intentionally ignored in summary formatting
		}
	case "top-rated":
		if _, err := fmt.Fprintf(writer, "Showing %d top-rated movies (%s)\n\n", len(movies), titleType); err != nil {
			// Error is intentionally ignored in summary formatting
		}
	case "now-playing":
		if _, err := fmt.Fprintf(writer, "Showing %d movies now playing (%s)\n\n", len(movies), titleType); err != nil {
			// Error is intentionally ignored in summary formatting
		}
	case "upcoming":
		if _, err := fmt.Fprintf(writer, "Showing %d upcoming movies (%s)\n\n", len(movies), titleType); err != nil {
			// Error is intentionally ignored in summary formatting
		}
	case "search":
		if _, err := fmt.Fprintf(writer, "Found %d movies (%s)\n\n", len(movies), titleType); err != nil {
			// Error is intentionally ignored in summary formatting
		}
	case "discover":
		if _, err := fmt.Fprintf(writer, "Discovered %d movies (%s)\n\n", len(movies), titleType); err != nil {
			// Error is intentionally ignored in summary formatting
		}
	}
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
