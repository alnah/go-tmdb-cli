// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"fmt"
	"slices"
	"strings"
)

// ValidateConfig validates the loaded configuration.
func ValidateConfig(config Config) error {
	// Check for empty API key (after trimming whitespace)
	if strings.TrimSpace(config.APIKey) == "" {
		return fmt.Errorf(
			"TMDB API key is required. Set TMDB_API_KEY environment variable or add to config file",
		)
	}

	if config.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive, got %v", config.Timeout)
	}

	if config.MaxRetries < 0 {
		return fmt.Errorf("max retries must be non-negative, got %d", config.MaxRetries)
	}

	if config.CacheTTL <= 0 {
		return fmt.Errorf("cache TTL must be positive, got %v", config.CacheTTL)
	}

	// Validate log level exactly as provided (no trimming)
	validLogLevels := []string{"debug", "info", "warn", "error", "silent"}
	if !contains(validLogLevels, config.LogLevel) {
		return fmt.Errorf("invalid log level '%s', must be one of: %s",
			config.LogLevel, strings.Join(validLogLevels, ", "))
	}

	// Validate format exactly as provided (no trimming)
	if !ValidateFormat(config.Format) {
		return fmt.Errorf("invalid default format '%s', must be one of: %s",
			config.Format, strings.Join(SupportedFormats(), ", "))
	}

	return nil
}

// Helper function to check if slice contains string.
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
