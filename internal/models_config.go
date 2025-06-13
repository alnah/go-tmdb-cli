// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"time"
)

// Config represents application configuration.
type Config struct {
	APIKey     string        `yaml:"api_key"     env:"TMDB_API_KEY"`
	BaseURL    string        `yaml:"base_url"`
	Timeout    time.Duration `yaml:"timeout"`
	MaxRetries int           `yaml:"max_retries"`
	CacheTTL   time.Duration `yaml:"cache_ttl"`
	LogLevel   string        `yaml:"log_level"`
	Format     string        `yaml:"format"`
}

// DefaultConfig returns default configuration.
func DefaultConfig() Config {
	return Config{
		BaseURL:    "https://api.themoviedb.org/3",
		Timeout:    DefaultTimeout * time.Second,
		MaxRetries: 3,
		CacheTTL:   DefaultCacheTTL * time.Minute,
		LogLevel:   "info",
		Format:     "table",
	}
}
