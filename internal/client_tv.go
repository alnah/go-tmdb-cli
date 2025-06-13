// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// GetPopularTVShows fetches popular TV shows.
func (c *Client) GetPopularTVShows(ctx context.Context, maxItems int) ([]TVShow, error) {
	return c.fetchTVShowPages(ctx, TVPopularEndpoint, nil, maxItems)
}

// GetTopRatedTVShows fetches top-rated TV shows.
func (c *Client) GetTopRatedTVShows(ctx context.Context, maxItems int) ([]TVShow, error) {
	return c.fetchTVShowPages(ctx, TVTopRatedEndpoint, nil, maxItems)
}

// GetOnTheAirTVShows fetches TV shows currently on the air.
func (c *Client) GetOnTheAirTVShows(ctx context.Context, maxItems int) ([]TVShow, error) {
	return c.fetchTVShowPages(ctx, TVOnTheAirEndpoint, nil, maxItems)
}

// SearchTVShows searches for TV shows by query.
func (c *Client) SearchTVShows(ctx context.Context, query string, maxItems int) ([]TVShow, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	params := url.Values{}
	params.Set("query", query)

	return c.fetchTVShowPages(ctx, TVSearchEndpoint, params, maxItems)
}

// DiscoverTVShows discovers TV shows with filters.
func (c *Client) DiscoverTVShows(ctx context.Context, opts SearchOptions) ([]TVShow, error) {
	params := c.buildTVDiscoverParams(opts)
	return c.fetchTVShowPages(ctx, TVDiscoverEndpoint, params, opts.MaxItems)
}

// fetchTVShowPages handles pagination and returns consolidated TV show results.
//
//nolint:dupl // Similar structure to fetchMoviePages but different types
func (c *Client) fetchTVShowPages(
	ctx context.Context,
	endpoint string,
	params url.Values,
	maxItems int,
) ([]TVShow, error) {
	if maxItems <= 0 {
		maxItems = 20
	}

	var allTVShows []TVShow
	page := 1

	for len(allTVShows) < maxItems {
		pageData, err := c.fetchPageData(ctx, endpoint, params, page, "tv_")
		if err != nil {
			return nil, err
		}
		if pageData == nil {
			break
		}

		var response TMDBTVResponse
		if err := json.Unmarshal(pageData, &response); err != nil {
			return nil, fmt.Errorf("parse response: %w", err)
		}

		tvShows := c.convertTVShows(response.Results)
		allTVShows = append(allTVShows, tvShows...)

		if page >= response.TotalPages || len(allTVShows) >= maxItems {
			break
		}
		page++
	}

	if len(allTVShows) > maxItems {
		allTVShows = allTVShows[:maxItems]
	}

	return allTVShows, nil
}

// convertTVShows converts TMDB TV API response to our simplified format.
func (c *Client) convertTVShows(tmdbTVShows []TMDBTVShow) []TVShow {
	shows := make([]TVShow, len(tmdbTVShows))

	for i, tm := range tmdbTVShows {
		shows[i] = c.convertTVShowSafely(tm)
	}

	return shows
}

// convertTVShowSafely converts a single TMDB TV show with error handling.
func (c *Client) convertTVShowSafely(tmdbShow TMDBTVShow) TVShow {
	// Map genres
	genreNames := c.mapTVGenres(tmdbShow.GenreIDs)

	show := TVShow{
		ID:           tmdbShow.ID,
		Name:         tmdbShow.Name,
		OriginalName: tmdbShow.OriginalName,
		Year:         ParseTVYear(tmdbShow.FirstAirDate),
		Rating:       tmdbShow.VoteAverage,
		Votes:        tmdbShow.VoteCount,
		Popularity:   tmdbShow.Popularity,
		Genres:       strings.Join(genreNames, ", "),
		Overview:     tmdbShow.Overview,
		Language:     tmdbShow.OriginalLanguage,
		Adult:        tmdbShow.Adult,
		FirstAirDate: tmdbShow.FirstAirDate,
	}

	// Fallback for empty name
	if show.Name == "" {
		show.Name = "Unknown TV Show"
	}

	return show
}

// buildTVDiscoverParams converts search options to API parameters for TV shows.
func (c *Client) buildTVDiscoverParams(opts SearchOptions) url.Values {
	params := c.buildBaseDiscoverParams(opts)

	// TV-specific parameters
	if opts.Year > 0 {
		params.Set("first_air_date_year", fmt.Sprintf("%d", opts.Year))
	}

	return params
}
