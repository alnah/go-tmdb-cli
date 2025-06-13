// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"fmt"
	"net/url"
	"strings"
)

// buildBaseDiscoverParams builds common discover parameters for both movies and TV shows.
func (c *Client) buildBaseDiscoverParams(opts SearchOptions) url.Values {
	params := make(url.Values)

	if opts.Language != "" {
		params.Set("with_original_language", opts.Language)
	}
	if opts.MinRating > 0 {
		params.Set("vote_average.gte", fmt.Sprintf("%.1f", opts.MinRating))
	}
	if opts.MaxRating > 0 {
		params.Set("vote_average.lte", fmt.Sprintf("%.1f", opts.MaxRating))
	}

	// Add vote count filters
	if opts.MinVotes > 0 {
		params.Set("vote_count.gte", fmt.Sprintf("%d", opts.MinVotes))
	}
	if opts.MaxVotes > 0 {
		params.Set("vote_count.lte", fmt.Sprintf("%d", opts.MaxVotes))
	}

	// Genres
	if len(opts.IncludeGenres) > 0 {
		genreStr := make([]string, len(opts.IncludeGenres))
		for i, id := range opts.IncludeGenres {
			genreStr[i] = fmt.Sprintf("%d", id)
		}
		params.Set("with_genres", strings.Join(genreStr, ","))
	}

	if len(opts.ExcludeGenres) > 0 {
		genreStr := make([]string, len(opts.ExcludeGenres))
		for i, id := range opts.ExcludeGenres {
			genreStr[i] = fmt.Sprintf("%d", id)
		}
		params.Set("without_genres", strings.Join(genreStr, ","))
	}

	// Default sorting
	sortBy := "popularity"
	if opts.SortBy != "" {
		sortBy = opts.SortBy
		// Handle special case for vote sorting
		if sortBy == "votes" {
			sortBy = "vote_count"
		}
	}
	sortOrder := "desc"
	if opts.SortOrder != "" {
		sortOrder = opts.SortOrder
	}
	params.Set("sort_by", sortBy+"."+sortOrder)

	return params
}
