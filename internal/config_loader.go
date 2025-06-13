// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from multiple sources with precedence:
// 1. Command line flags (handled by caller)
// 2. Environment variables
// 3. Config file
// 4. Defaults.
func LoadConfig() (Config, error) {
	// Start with defaults
	config := DefaultConfig()

	// Load from config file if it exists
	if err := loadConfigFile(&config); err != nil {
		// Config file errors are not fatal, just use defaults
		fmt.Fprintf(os.Stderr, "Warning: Could not load config file: %v\n", err)
	}

	// Override with environment variables
	loadFromEnvironment(&config)

	// Validate the final configuration
	if err := ValidateConfig(config); err != nil {
		return config, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// loadConfigFile attempts to load configuration from YAML file.
func loadConfigFile(config *Config) error {
	configPaths := []string{
		"./config.yaml",
		"./tmdb.yaml",
		filepath.Join(os.Getenv("HOME"), ".tmdb", "config.yaml"),
		filepath.Join(os.Getenv("HOME"), ".config", "tmdb", "config.yaml"),
	}

	for _, path := range configPaths {
		if data, err := os.ReadFile(path); err == nil {
			if err := yaml.Unmarshal(data, config); err != nil {
				return fmt.Errorf("parse config file %s: %w", path, err)
			}
			return nil
		}
	}

	// No config file found - this is okay
	return nil
}

// loadFromEnvironment loads configuration from environment variables.
func loadFromEnvironment(config *Config) {
	// API Key
	if apiKey := os.Getenv("TMDB_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
	}

	// Base URL
	if baseURL := os.Getenv("TMDB_BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	// Timeout
	if timeoutStr := os.Getenv("TMDB_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			config.Timeout = timeout
		}
	}

	// Max Retries
	if retriesStr := os.Getenv("TMDB_MAX_RETRIES"); retriesStr != "" {
		if retries, err := strconv.Atoi(retriesStr); err == nil {
			config.MaxRetries = retries
		}
	}

	// Cache TTL
	if cacheTTLStr := os.Getenv("TMDB_CACHE_TTL"); cacheTTLStr != "" {
		if cacheTTL, err := time.ParseDuration(cacheTTLStr); err == nil {
			config.CacheTTL = cacheTTL
		}
	}

	// Log Level
	if logLevel := os.Getenv("TMDB_LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}

	// Default Format
	if format := os.Getenv("TMDB_FORMAT"); format != "" {
		config.Format = format
	}
}
