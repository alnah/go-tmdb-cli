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

// MovieFormatter implements MediaFormatter for movies.
type MovieFormatter struct {
	Movie Movie
}

// GetTitle returns the movie title based on preference.
func (mf MovieFormatter) GetTitle(useOriginal bool) string {
	return formatTitle(mf.Movie.Title, mf.Movie.OriginalTitle, useOriginal)
}

// GetOriginalTitle returns the original movie title.
func (mf MovieFormatter) GetOriginalTitle() string {
	return mf.Movie.OriginalTitle
}

// GetFields returns CSV fields for a movie.
func (mf MovieFormatter) GetFields() []string {
	return []string{
		strconv.Itoa(mf.Movie.ID),
		mf.Movie.Title,
		mf.Movie.OriginalTitle,
		strconv.Itoa(mf.Movie.Year),
		fmt.Sprintf("%.1f", mf.Movie.Rating),
		strconv.Itoa(mf.Movie.Votes),
		fmt.Sprintf("%.2f", mf.Movie.Popularity),
		mf.Movie.Genres,
		mf.Movie.Overview,
		mf.Movie.Language,
		strconv.FormatBool(mf.Movie.Adult),
	}
}

// GetEmptyMessage returns the message for empty movie results.
func (mf MovieFormatter) GetEmptyMessage() string {
	return "No movies found."
}

// GetHeaderName returns the header name for movie title column.
func (mf MovieFormatter) GetHeaderName() string {
	return "Title"
}

// TVShowFormatter implements MediaFormatter for TV shows.
type TVShowFormatter struct {
	TVShow TVShow
}

// GetTitle returns the TV show name based on preference.
func (tvf TVShowFormatter) GetTitle(useOriginal bool) string {
	return formatTitle(tvf.TVShow.Name, tvf.TVShow.OriginalName, useOriginal)
}

// GetOriginalTitle returns the original TV show name.
func (tvf TVShowFormatter) GetOriginalTitle() string {
	return tvf.TVShow.OriginalName
}

// GetFields returns CSV fields for a TV show.
func (tvf TVShowFormatter) GetFields() []string {
	return []string{
		strconv.Itoa(tvf.TVShow.ID),
		tvf.TVShow.Name,
		tvf.TVShow.OriginalName,
		strconv.Itoa(tvf.TVShow.Year),
		fmt.Sprintf("%.1f", tvf.TVShow.Rating),
		strconv.Itoa(tvf.TVShow.Votes),
		fmt.Sprintf("%.2f", tvf.TVShow.Popularity),
		tvf.TVShow.Genres,
		tvf.TVShow.Overview,
		tvf.TVShow.Language,
		strconv.FormatBool(tvf.TVShow.Adult),
	}
}

// GetEmptyMessage returns the message for empty TV show results.
func (tvf TVShowFormatter) GetEmptyMessage() string {
	return "No TV shows found."
}

// GetHeaderName returns the header name for TV show name column.
func (tvf TVShowFormatter) GetHeaderName() string {
	return "Name"
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

// JSON formatting functions

func formatMoviesJSON(writer io.Writer, movies []Movie) error {
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

func formatTVShowsJSON(writer io.Writer, shows []TVShow) error {
	output := struct {
		TVShows []TVShow `json:"tv_shows"`
		Count   int      `json:"count"`
	}{
		TVShows: shows,
		Count:   len(shows),
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

func formatTVSearchResultJSON(writer io.Writer, result TVSearchResult) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(result)
}

// CSV formatting functions

func getMovieHeaders() []string {
	return []string{
		"ID", "Title", "Original Title", "Year", "Rating", "Votes",
		"Popularity", "Genres", "Overview", "Language", "Adult",
	}
}

func getTVShowHeaders() []string {
	return []string{
		"ID", "Name", "Original Name", "Year", "Rating", "Votes",
		"Popularity", "Genres", "Overview", "Language", "Adult",
	}
}

func formatCSVGeneric(
	writer io.Writer,
	formatters []MediaFormatter,
	headers []string,
	options FormatOptions,
) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header if not disabled
	if !options.NoHeader {
		if err := csvWriter.Write(headers); err != nil {
			return fmt.Errorf("write CSV header: %w", err)
		}
	}

	// Write data rows
	for _, formatter := range formatters {
		if err := csvWriter.Write(formatter.GetFields()); err != nil {
			return fmt.Errorf("write CSV row: %w", err)
		}
	}

	return nil
}

// Table formatting functions

func formatMoviesTable(writer io.Writer, movies []Movie, options FormatOptions) error {
	if len(movies) == 0 {
		_, err := fmt.Fprintln(writer, "No movies found.")
		return err
	}

	formatters := make([]any, len(movies))
	for i, movie := range movies {
		formatters[i] = movie
	}

	return formatTableGeneric(writer, formatters, options, "Title")
}

func formatTVShowsTable(writer io.Writer, shows []TVShow, options FormatOptions) error {
	if len(shows) == 0 {
		_, err := fmt.Fprintln(writer, "No TV shows found.")
		return err
	}

	formatters := make([]any, len(shows))
	for i, show := range shows {
		formatters[i] = show
	}

	return formatTableGeneric(writer, formatters, options, "Name")
}

func formatTableGeneric(
	writer io.Writer,
	items []any,
	options FormatOptions,
	titleHeader string,
) error {
	t := table.NewWriter()
	t.SetOutputMirror(writer)
	t.SetStyle(table.StyleDefault)

	// Add headers if not disabled
	if !options.NoHeader {
		headers := []string{"#", titleHeader, "Year", "Rating", "Votes", "Genres"}
		headerRow := make(table.Row, len(headers))
		for i, h := range headers {
			headerRow[i] = h
		}
		t.AppendHeader(headerRow)
	}

	// Add rows
	for i, item := range items {
		var row table.Row

		switch v := item.(type) {
		case Movie:
			title := formatTitle(v.Title, v.OriginalTitle, options.UseOriginalTitle)
			rating := FormatRating(v.Rating, v.Votes)
			votes := FormatVotes(v.Votes)
			genres := FormatGenres(v.Genres, GenresColumnWidth)

			row = table.Row{
				i + 1,
				addMovieIndicators(v, title),
				formatYear(v.Year),
				rating,
				votes,
				genres,
			}
		case TVShow:
			title := formatTitle(v.Name, v.OriginalName, options.UseOriginalTitle)
			rating := FormatRating(v.Rating, v.Votes)
			votes := FormatVotes(v.Votes)
			genres := FormatGenres(v.Genres, GenresColumnWidth)

			row = table.Row{
				i + 1,
				addTVShowIndicators(v, title),
				formatYear(v.Year),
				rating,
				votes,
				genres,
			}
		}

		t.AppendRow(row)
	}

	// Configure column properties
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: ColNumber, Align: text.AlignCenter, WidthMax: NumberColWidth}, // #
		{Number: ColTitle, Align: text.AlignLeft, WidthMax: TitleColWidth},     // Title/Name
		{Number: ColYear, Align: text.AlignCenter, WidthMax: YearColWidth},     // Year
		{Number: ColRating, Align: text.AlignRight, WidthMax: RatingColWidth},  // Rating
		{Number: ColVotes, Align: text.AlignRight, WidthMax: VotesColWidth},    // Votes
		{Number: ColGenres, Align: text.AlignLeft, WidthMax: GenresColWidth},   // Genres
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

	return formatMoviesTable(writer, result.Movies, options)
}

func formatTVSearchResultTable(
	writer io.Writer,
	result TVSearchResult,
	options FormatOptions,
) error {
	if len(result.TVShows) == 0 {
		_, err := fmt.Fprintln(writer, "No TV shows found.")
		return err
	}

	// Show pagination info first
	if result.TotalPages > 1 {
		if _, err := fmt.Fprintf(writer, "Page %d of %d (%d total results)\n\n",
			result.Page, result.TotalPages, result.TotalResults); err != nil {
			return err
		}
	}

	return formatTVShowsTable(writer, result.TVShows, options)
}

// Helper functions

func formatTitle(title, originalTitle string, useOriginal bool) string {
	var displayTitle string
	if useOriginal && originalTitle != "" && originalTitle != title {
		displayTitle = originalTitle
		// Add English title/name in parentheses if space allows
		fullTitle := displayTitle + " (" + title + ")"
		if len(fullTitle) <= TitleColumnWidth {
			displayTitle = fullTitle
		}
	} else {
		displayTitle = title
		// Add original title/name in parentheses if different and space allows
		if originalTitle != "" && originalTitle != title {
			fullTitle := displayTitle + " (" + originalTitle + ")"
			if len(fullTitle) <= TitleColumnWidth {
				displayTitle = fullTitle
			}
		}
	}

	// Truncate if too long
	if len(displayTitle) > TitleColumnWidth {
		displayTitle = displayTitle[:TitleColumnWidth-3] + "..."
	}

	return displayTitle
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

func addTVShowIndicators(show TVShow, title string) string {
	var indicators []string

	if IsRecent(show.Year) {
		indicators = append(indicators, "[NEW]")
	}
	if IsHighlyRated(show.Rating, show.Votes) {
		indicators = append(indicators, "[TOP]")
	}
	if IsPopular(show.Popularity, show.Votes) {
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

// FormatError formats error messages for output.
func FormatError(writer io.Writer, err error) error {
	_, writeErr := fmt.Fprintf(writer, "Error: %s\n", err.Error())
	return writeErr
}

// FormatSummary formats summary information for different commands.
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
		_, _ = fmt.Fprintf(writer, "Showing %d popular movies (%s)\n\n", len(movies), titleType)
	case "top-rated":
		_, _ = fmt.Fprintf(writer, "Showing %d top-rated movies (%s)\n\n", len(movies), titleType)
	case "now-playing":
		_, _ = fmt.Fprintf(writer, "Showing %d movies now playing (%s)\n\n", len(movies), titleType)
	case "upcoming":
		_, _ = fmt.Fprintf(writer, "Showing %d upcoming movies (%s)\n\n", len(movies), titleType)
	case "search":
		_, _ = fmt.Fprintf(writer, "Found %d movies (%s)\n\n", len(movies), titleType)
	case "discover":
		_, _ = fmt.Fprintf(writer, "Discovered %d movies (%s)\n\n", len(movies), titleType)
	}
}

// FormatTVSummary formats summary information for different TV commands.
func FormatTVSummary(writer io.Writer, shows []TVShow, command string, useOriginal bool) {
	if len(shows) == 0 {
		return
	}

	titleType := "names"
	if useOriginal {
		titleType = "original language names"
	}

	switch command {
	case "popular":
		_, _ = fmt.Fprintf(writer, "Showing %d popular TV shows (%s)\n\n", len(shows), titleType)
	case "top-rated":
		_, _ = fmt.Fprintf(writer, "Showing %d top-rated TV shows (%s)\n\n", len(shows), titleType)
	case "on-the-air":
		_, _ = fmt.Fprintf(writer, "Showing %d TV shows on the air (%s)\n\n", len(shows), titleType)
	case "search":
		_, _ = fmt.Fprintf(writer, "Found %d TV shows (%s)\n\n", len(shows), titleType)
	case "discover":
		_, _ = fmt.Fprintf(writer, "Discovered %d TV shows (%s)\n\n", len(shows), titleType)
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
