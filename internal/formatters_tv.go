// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"fmt"
	"strconv"
)

//nolint:dupl // Similar to formatters_movie.go but for different types

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

// getTVShowHeaders returns CSV headers for TV shows.
func getTVShowHeaders() []string {
	return []string{
		"ID", "Name", "Original Name", "Year", "Rating", "Votes",
		"Popularity", "Genres", "Overview", "Language", "Adult",
	}
}
