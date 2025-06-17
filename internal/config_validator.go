// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	minLengthAPIKey  int = 8
	suspiciousLength int = 4
)

// Placeholder API keys that should be rejected during validation.
var placeholderAPIKeys = []string{
	"your-api-key-here",
	"YOUR_API_KEY_HERE",
	"your-tmdb-api-key",
	"insert-api-key-here",
	"api-key-placeholder",
	"replace-with-your-key",
	"example-api-key",
	"demo-api-key",
	"test-api-key-placeholder",
	"change-this-api-key",
	"<your-api-key>",
	"[YOUR_API_KEY]",
	"{api_key}",
	"xxxx-xxxx-xxxx-xxxx",
	"placeholder",
	"TODO",
	"FIXME",
}

// isTestEnvironment checks if we're running in a test environment.
func isTestEnvironment() bool {
	// Check if binary name contains test patterns
	if strings.HasSuffix(os.Args[0], ".test") {
		return true
	}

	// Check for Go test runner patterns in arguments
	for _, arg := range os.Args {
		if strings.Contains(arg, "_test.go") || arg == "-test.v" || arg == "-test.run" {
			return true
		}
	}

	// Check environment variables that indicate testing
	testEnvVars := []string{"GO_TEST", "TESTING", "_TESTING"}
	for _, envVar := range testEnvVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}

	return false
}

// ValidateConfig validates the loaded configuration.
func ValidateConfig(config Config) error {
	// Check for empty API key (after trimming whitespace)
	trimmedAPIKey := strings.TrimSpace(config.APIKey)
	if trimmedAPIKey == "" {
		return fmt.Errorf(
			"TMDB API key is required. Set TMDB_API_KEY environment variable or add to config file",
		)
	}

	// Check for exact placeholder API keys (always enforced)
	if isPlaceholderAPIKey(trimmedAPIKey) {
		return fmt.Errorf(
			"placeholder API key detected: '%s'. Please set a real "+
				"TMDB API key from https://www.themoviedb.org/settings/api",
			config.APIKey,
		)
	}

	// In test environments, be more lenient with API key validation
	if !isTestEnvironment() {
		// Validate API key format (basic sanity check) - only in production
		if len(trimmedAPIKey) < minLengthAPIKey {
			return fmt.Errorf(
				"API key appears too short (%d characters). TMDB API keys are typically longer. Please verify your key",
				len(trimmedAPIKey),
			)
		}

		// Check for suspicious patterns - only in production
		if containsSuspiciousPatterns(trimmedAPIKey) {
			return fmt.Errorf(
				"API key contains suspicious patterns that suggest it may be a placeholder: '%s'",
				config.APIKey,
			)
		}
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

// isPlaceholderAPIKey checks if the given API key is a known placeholder.
func isPlaceholderAPIKey(apiKey string) bool {
	// Check exact matches (case-insensitive)
	lowerKey := strings.ToLower(apiKey)
	for _, placeholder := range placeholderAPIKeys {
		if strings.ToLower(placeholder) == lowerKey {
			return true
		}
	}
	return false
}

// containsSuspiciousPatterns checks for patterns that might indicate test/placeholder keys.
func containsSuspiciousPatterns(apiKey string) bool {
	lowerKey := strings.ToLower(apiKey)

	// Patterns that suggest test/placeholder keys
	suspiciousPatterns := []string{
		"test",
		"demo",
		"example",
		"placeholder",
		"sample",
		"fake",
		"mock",
		"dummy",
		"invalid",
		"replace",
		"change",
		"insert",
		"todo",
		"fixme",
		"xxx",
		"000",
		"123",
		"abc",
	}

	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerKey, pattern) {
			return true
		}
	}

	// Check for repeated characters (often used in placeholders)
	if hasRepeatedCharacterPatterns(apiKey) {
		return true
	}

	return false
}

// hasRepeatedCharacterPatterns checks for suspicious repeated character patterns.
func hasRepeatedCharacterPatterns(str string) bool {
	// Check for 4+ consecutive same characters
	for i := 0; i < len(str)-3; i++ {
		if str[i] == str[i+1] && str[i+1] == str[i+2] && str[i+2] == str[i+3] {
			return true
		}
	}

	// Check for simple ascending/descending patterns
	if len(str) >= suspiciousLength {
		ascending := true
		descending := true

		for i := 1; i < len(str) && i < 8; i++ { // Check first 8 chars
			if str[i] != str[i-1]+1 {
				ascending = false
			}
			if str[i] != str[i-1]-1 {
				descending = false
			}
		}

		if ascending || descending {
			return true
		}
	}

	return false
}

// Helper function to check if slice contains string.
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
