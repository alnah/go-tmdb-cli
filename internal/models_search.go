// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"fmt"
)

// SearchOptions represents search parameters.
type SearchOptions struct {
	Query         string
	Page          int
	Language      string
	Year          int
	MinRating     float64
	MaxRating     float64
	MinVotes      int // Minimum votes filter
	MaxVotes      int // Maximum votes filter
	IncludeGenres []int
	ExcludeGenres []int
	SortBy        string
	SortOrder     string
	MaxItems      int
}

// ValidateSearchOptions validates search parameters.
func ValidateSearchOptions(opts *SearchOptions) error {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.MaxItems < 1 || opts.MaxItems > OneThousand {
		return fmt.Errorf("max-items must be between 1 and 1000, got %d", opts.MaxItems)
	}
	if opts.MinRating < 0 || opts.MinRating > 10 {
		return fmt.Errorf("min-rating must be between 0 and 10, got %.1f", opts.MinRating)
	}
	if opts.MaxRating < 0 || opts.MaxRating > 10 {
		return fmt.Errorf("max-rating must be between 0 and 10, got %.1f", opts.MaxRating)
	}
	if opts.MinRating > 0 && opts.MaxRating > 0 && opts.MinRating > opts.MaxRating {
		return fmt.Errorf(
			"min-rating (%.1f) cannot be greater than max-rating (%.1f)",
			opts.MinRating,
			opts.MaxRating,
		)
	}
	return nil
}

// Genre represents a movie/TV genre.
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MediaItem interface for unified processing of movies and TV shows.
type MediaItem interface {
	GetID() int
	GetTitle() string         // Name for TV, Title for Movie
	GetOriginalTitle() string // OriginalName for TV, OriginalTitle for Movie
	GetYear() int
	GetRating() float64
	GetVotes() int
	GetPopularity() float64
	GetGenres() string
	GetOverview() string
}
