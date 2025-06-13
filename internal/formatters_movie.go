// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"fmt"
	"strconv"
)

//nolint:dupl // Similar to formatters_tv.go but for different types

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

// getMovieHeaders returns CSV headers for movies.
func getMovieHeaders() []string {
	return []string{
		"ID", "Title", "Original Title", "Year", "Rating", "Votes",
		"Popularity", "Genres", "Overview", "Language", "Adult",
	}
}
