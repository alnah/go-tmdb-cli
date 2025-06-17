// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"context"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"time"
)

// fetchPageData fetches a single page of data, either from cache or API.
func (c *Client) fetchPageData(
	ctx context.Context,
	endpoint string,
	params url.Values,
	page int,
	cachePrefix string,
) ([]byte, error) {
	// Build cache key
	pageParams := make(url.Values)
	if params != nil {
		maps.Copy(pageParams, params)
	}
	pageParams.Set("page", fmt.Sprintf("%d", page))
	cacheKey := fmt.Sprintf("%s%s?%s", cachePrefix, endpoint, pageParams.Encode())

	// Try cache first
	if data := c.getFromCache(cacheKey); data != nil {
		return data, nil
	}

	// Make API call with rate limiting
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	// Make HTTP request
	data, err := c.makeGenericRequest(ctx, endpoint, params, page)
	if err != nil {
		return nil, err
	}

	// Cache the response
	c.putInCache(cacheKey, data)

	return data, nil
}

// makeGenericRequest makes a generic paginated request to TMDB API.
func (c *Client) makeGenericRequest(
	ctx context.Context,
	endpoint string,
	params url.Values,
	page int,
) ([]byte, error) {
	pageParams := make(url.Values)
	maps.Copy(pageParams, params)
	pageParams.Set("page", fmt.Sprintf("%d", page))

	return c.makeHTTPRequest(ctx, endpoint, pageParams)
}

// makeHTTPRequest makes the actual HTTP request to TMDB API.
func (c *Client) makeHTTPRequest(
	ctx context.Context,
	endpoint string,
	params url.Values,
) ([]byte, error) {
	// Build URL
	apiURL, err := url.Parse(c.config.BaseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Add API key and params
	q := apiURL.Query()
	q.Set("api_key", c.config.APIKey)
	for key, values := range params {
		for _, value := range values {
			q.Add(key, value)
		}
	}
	apiURL.RawQuery = q.Encode()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "tmdb-cli/2.0")

	// Make request with retries
	var resp *http.Response
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			if attempt == c.config.MaxRetries {
				return nil, fmt.Errorf(
					"request failed after %d attempts: %w",
					c.config.MaxRetries+1,
					err,
				)
			}
			// Use clock for retry delay
			c.clock.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}
		break
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Check for API errors
	if resp.StatusCode >= BadRequestStatus {
		return nil, c.handleAPIError(resp.StatusCode, body)
	}

	return body, nil
}

// handleAPIError creates user-friendly error messages.
func (c *Client) handleAPIError(statusCode int, body []byte) error {
	switch statusCode {
	case UnauthorizedStatus:
		return fmt.Errorf(
			"invalid TMDB API key - get one from https://www.themoviedb.org/settings/api",
		)
	case NotFoundStatus:
		return fmt.Errorf("resource not found")
	case RateLimitStatus:
		return fmt.Errorf("rate limit exceeded - please wait and try again")
	default:
		return fmt.Errorf("API error (status %d): %s", statusCode, string(body))
	}
}
