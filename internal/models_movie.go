// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"errors"
)

// Movie represents a simplified movie structure.
type Movie struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title,omitempty"`
	Year          int     `json:"year"`
	Rating        float64 `json:"rating"`
	Votes         int     `json:"votes"`
	Popularity    float64 `json:"popularity"`
	Genres        string  `json:"genres"`
	Overview      string  `json:"overview,omitempty"`
	Language      string  `json:"language"`
	Adult         bool    `json:"adult"`
	ReleaseDate   string  `json:"release_date,omitempty"`
}

// Validate validates the Movie data structure.
func (m *Movie) Validate() error {
	if m.Title == "" {
		return errors.New("movie title is required")
	}
	if m.ID <= 0 {
		return errors.New("movie ID must be positive")
	}
	return nil
}

// Implement MediaItem interface for Movie

// GetID returns the movie ID.
func (m Movie) GetID() int { return m.ID }

// GetTitle returns the movie title.
func (m Movie) GetTitle() string { return m.Title }

// GetOriginalTitle returns the movie original title.
func (m Movie) GetOriginalTitle() string { return m.OriginalTitle }

// GetYear returns the movie year.
func (m Movie) GetYear() int { return m.Year }

// GetRating returns the movie rating.
func (m Movie) GetRating() float64 { return m.Rating }

// GetVotes returns the movie votes.
func (m Movie) GetVotes() int { return m.Votes }

// GetPopularity returns the movie popularity.
func (m Movie) GetPopularity() float64 { return m.Popularity }

// GetGenres returns the movie genres.
func (m Movie) GetGenres() string { return m.Genres }

// GetOverview returns the movie overview.
func (m Movie) GetOverview() string { return m.Overview }

// SearchResult represents paginated movie search results.
type SearchResult struct {
	Movies       []Movie `json:"movies"`
	Page         int     `json:"page"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}
