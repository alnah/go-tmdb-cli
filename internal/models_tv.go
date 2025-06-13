// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"errors"
)

// TV-specific constants.
const (
	TVMinValidYear      = 1928 // First television broadcast
	TVMaxReasonableYear = 2030
	TVDefaultTimeout    = 30
	TVCacheTTL          = 5 // minutes
)

// TV Show endpoints.
const (
	TVPopularEndpoint  = "/tv/popular"
	TVTopRatedEndpoint = "/tv/top_rated"
	TVOnTheAirEndpoint = "/tv/on_the_air"
	TVSearchEndpoint   = "/search/tv"
	TVDiscoverEndpoint = "/discover/tv"
	TVGenresEndpoint   = "/genre/tv/list"
)

// TVShow represents a simplified TV show structure.
type TVShow struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`          // Equivalent to Movie.Title
	OriginalName string  `json:"original_name"` // Equivalent to Movie.OriginalTitle
	Year         int     `json:"year"`          // Extracted from first_air_date
	Rating       float64 `json:"rating"`        // vote_average
	Votes        int     `json:"votes"`         // vote_count
	Popularity   float64 `json:"popularity"`
	Genres       string  `json:"genres"` // Mapped genres as string
	Overview     string  `json:"overview,omitempty"`
	Language     string  `json:"language"`
	Adult        bool    `json:"adult"`
	FirstAirDate string  `json:"first_air_date,omitempty"` // First broadcast date
}

// Validate validates the TVShow data structure.
func (tv *TVShow) Validate() error {
	if tv.Name == "" {
		return errors.New("TV show name is required")
	}
	if tv.ID <= 0 {
		return errors.New("TV show ID must be positive")
	}
	return nil
}

// Implement MediaItem interface for TVShow

// GetID returns the TV show ID.
func (tv TVShow) GetID() int { return tv.ID }

// GetTitle returns the TV show name.
func (tv TVShow) GetTitle() string { return tv.Name }

// GetOriginalTitle returns the TV show original name.
func (tv TVShow) GetOriginalTitle() string { return tv.OriginalName }

// GetYear returns the TV show year.
func (tv TVShow) GetYear() int { return tv.Year }

// GetRating returns the TV show rating.
func (tv TVShow) GetRating() float64 { return tv.Rating }

// GetVotes returns the TV show votes.
func (tv TVShow) GetVotes() int { return tv.Votes }

// GetPopularity returns the TV show popularity.
func (tv TVShow) GetPopularity() float64 { return tv.Popularity }

// GetGenres returns the TV show genres.
func (tv TVShow) GetGenres() string { return tv.Genres }

// GetOverview returns the TV show overview.
func (tv TVShow) GetOverview() string { return tv.Overview }

// TVSearchResult represents paginated TV search results.
type TVSearchResult struct {
	TVShows      []TVShow `json:"tv_shows"`
	Page         int      `json:"page"`
	TotalPages   int      `json:"total_pages"`
	TotalResults int      `json:"total_results"`
}
