// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// mapGenres converts genre IDs to names for movies.
func (c *Client) mapGenres(genreIDs []int) []string {
	c.genresMu.RLock()
	defer c.genresMu.RUnlock()

	names := make([]string, 0, len(genreIDs))
	for _, id := range genreIDs {
		if name, ok := c.genres[id]; ok {
			names = append(names, name)
		}
	}
	return names
}

// mapTVGenres converts genre IDs to names for TV shows.
func (c *Client) mapTVGenres(genreIDs []int) []string {
	c.tvGenresMu.RLock()
	defer c.tvGenresMu.RUnlock()

	names := make([]string, 0, len(genreIDs))
	for _, id := range genreIDs {
		if name, ok := c.tvGenres[id]; ok {
			names = append(names, name)
		}
	}
	return names
}

// loadGenres loads movie genre mappings from API.
func (c *Client) loadGenres() {
	c.loadGenreMapping(MovieGenresEndpoint, func(genres []Genre) {
		c.genresMu.Lock()
		defer c.genresMu.Unlock()
		for _, genre := range genres {
			c.genres[genre.ID] = genre.Name
		}
	})
}

// loadTVGenres loads TV show genre mappings from API.
func (c *Client) loadTVGenres() {
	c.loadGenreMapping(TVGenresEndpoint, func(genres []Genre) {
		c.tvGenresMu.Lock()
		defer c.tvGenresMu.Unlock()
		for _, genre := range genres {
			c.tvGenres[genre.ID] = genre.Name
		}
	})
}

// loadGenreMapping is a generic function to load genre mappings.
func (c *Client) loadGenreMapping(endpoint string, handler func([]Genre)) {
	ctx, cancel := context.WithTimeout(context.Background(), LoadGenresTimeout*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		c.config.BaseURL+endpoint+"?api_key="+c.config.APIKey,
		nil,
	)
	if err != nil {
		return
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var genresResp TMDBGenresResponse
	if err := json.NewDecoder(resp.Body).Decode(&genresResp); err != nil {
		return
	}

	handler(genresResp.Genres)
}
