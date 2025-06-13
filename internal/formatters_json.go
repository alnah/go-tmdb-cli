// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"encoding/json"
	"io"
)

// formatMoviesJSON formats movies as JSON.
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

// formatTVShowsJSON formats TV shows as JSON.
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

// formatSearchResultJSON formats movie search results as JSON.
func formatSearchResultJSON(writer io.Writer, result SearchResult) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(result)
}

// formatTVSearchResultJSON formats TV search results as JSON.
func formatTVSearchResultJSON(writer io.Writer, result TVSearchResult) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(result)
}
