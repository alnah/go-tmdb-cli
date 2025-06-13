// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"fmt"
	"io"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// formatMoviesTable formats movies as a table.
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

// formatTVShowsTable formats TV shows as a table.
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

// formatTableGeneric creates a table for any media type.
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

// formatSearchResultTable formats search results as a table with pagination info.
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

// formatTVSearchResultTable formats TV search results as a table with pagination info.
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

// addMovieIndicators adds visual indicators for movies.
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

// addTVShowIndicators adds visual indicators for TV shows.
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
