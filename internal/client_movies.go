// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Movie API endpoints.
const (
	MoviePopularEndpoint    = "/movie/popular"
	MovieTopRatedEndpoint   = "/movie/top_rated"
	MovieNowPlayingEndpoint = "/movie/now_playing"
	MovieUpcomingEndpoint   = "/movie/upcoming"
	MovieSearchEndpoint     = "/search/movie"
	MovieDiscoverEndpoint   = "/discover/movie"
	MovieGenresEndpoint     = "/genre/movie/list"
)

// GetPopularMovies fetches popular movies.
func (c *Client) GetPopularMovies(ctx context.Context, maxItems int) ([]Movie, error) {
	return c.fetchMoviePages(ctx, MoviePopularEndpoint, nil, maxItems)
}

// GetTopRatedMovies fetches top-rated movies.
func (c *Client) GetTopRatedMovies(ctx context.Context, maxItems int) ([]Movie, error) {
	return c.fetchMoviePages(ctx, MovieTopRatedEndpoint, nil, maxItems)
}

// GetNowPlayingMovies fetches movies currently in theaters.
func (c *Client) GetNowPlayingMovies(ctx context.Context, maxItems int) ([]Movie, error) {
	return c.fetchMoviePages(ctx, MovieNowPlayingEndpoint, nil, maxItems)
}

// GetUpcomingMovies fetches upcoming movie releases.
func (c *Client) GetUpcomingMovies(ctx context.Context, maxItems int) ([]Movie, error) {
	return c.fetchMoviePages(ctx, MovieUpcomingEndpoint, nil, maxItems)
}

// SearchMovies searches for movies by query.
func (c *Client) SearchMovies(ctx context.Context, query string, maxItems int) ([]Movie, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	params := url.Values{}
	params.Set("query", query)

	return c.fetchMoviePages(ctx, MovieSearchEndpoint, params, maxItems)
}

// DiscoverMovies discovers movies with filters.
func (c *Client) DiscoverMovies(ctx context.Context, opts SearchOptions) ([]Movie, error) {
	params := c.buildDiscoverParams(opts)
	return c.fetchMoviePages(ctx, MovieDiscoverEndpoint, params, opts.MaxItems)
}

// fetchMoviePages handles pagination and returns consolidated movie results.
//
//nolint:dupl // Similar structure to fetchTVShowPages but different types
func (c *Client) fetchMoviePages(
	ctx context.Context,
	endpoint string,
	params url.Values,
	maxItems int,
) ([]Movie, error) {
	if maxItems <= 0 {
		maxItems = 20
	}

	var allMovies []Movie
	page := 1

	for len(allMovies) < maxItems {
		pageData, err := c.fetchPageData(ctx, endpoint, params, page, "")
		if err != nil {
			return nil, err
		}
		if pageData == nil {
			break
		}

		var response TMDBResponse
		if err := json.Unmarshal(pageData, &response); err != nil {
			return nil, fmt.Errorf("parse response: %w", err)
		}

		movies := c.convertMovies(response.Results)
		allMovies = append(allMovies, movies...)

		if page >= response.TotalPages || len(allMovies) >= maxItems {
			break
		}
		page++
	}

	if len(allMovies) > maxItems {
		allMovies = allMovies[:maxItems]
	}

	return allMovies, nil
}

// convertMovies converts TMDB API response to our simplified format.
func (c *Client) convertMovies(tmdbMovies []TMDBMovie) []Movie {
	movies := make([]Movie, len(tmdbMovies))

	for i, tm := range tmdbMovies {
		// Map genres
		genreNames := c.mapGenres(tm.GenreIDs)

		movies[i] = Movie{
			ID:            tm.ID,
			Title:         tm.Title,
			OriginalTitle: tm.OriginalTitle,
			Year:          ParseYear(tm.ReleaseDate),
			Rating:        tm.VoteAverage,
			Votes:         tm.VoteCount,
			Popularity:    tm.Popularity,
			Genres:        strings.Join(genreNames, ", "),
			Overview:      tm.Overview,
			Language:      "en", // Simplified for now
			Adult:         tm.Adult,
			ReleaseDate:   tm.ReleaseDate,
		}
	}

	return movies
}

// buildDiscoverParams converts search options to API parameters for movies.
func (c *Client) buildDiscoverParams(opts SearchOptions) url.Values {
	params := c.buildBaseDiscoverParams(opts)

	// Movie-specific parameters
	if opts.Year > 0 {
		params.Set("primary_release_year", fmt.Sprintf("%d", opts.Year))
	}

	return params
}
