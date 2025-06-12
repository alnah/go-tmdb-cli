// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
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

// ValidateConfig validates the loaded configuration.
func ValidateConfig(config Config) error {
	if config.APIKey == "" {
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

	validLogLevels := []string{"debug", "info", "warn", "error", "silent"}
	if !contains(validLogLevels, config.LogLevel) {
		return fmt.Errorf("invalid log level '%s', must be one of: %s",
			config.LogLevel, strings.Join(validLogLevels, ", "))
	}

	if !ValidateFormat(config.Format) {
		return fmt.Errorf("invalid default format '%s', must be one of: %s",
			config.Format, strings.Join(SupportedFormats(), ", "))
	}

	return nil
}

// CreateExampleConfig creates an example configuration file.
func CreateExampleConfig(path string) error {
	config := DefaultConfig()
	config.APIKey = "your-api-key-here"

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	// Add comments to the example
	example := `# TMDB CLI Configuration
# Get your API key from: https://www.themoviedb.org/settings/api

` + string(data)

	// Ensure directory exists
	dir := filepath.Dir(path)
	const dirPerm = 0o644
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	const filePerm = 0o644
	if err := os.WriteFile(path, []byte(example), filePerm); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	return nil
}

// GetConfigHelp returns help text for configuration.
func GetConfigHelp() string {
	return `Configuration can be set via:

1. Environment variables:
   TMDB_API_KEY      - Your TMDB API key (required)
   TMDB_BASE_URL     - TMDB API base URL
   TMDB_TIMEOUT      - Request timeout (e.g., "30s")
   TMDB_MAX_RETRIES  - Maximum request retries
   TMDB_CACHE_TTL    - Cache TTL (e.g., "5m")
   TMDB_LOG_LEVEL    - Log level (debug, info, warn, error, silent)
   TMDB_FORMAT       - Default output format (table, json, csv)

2. Config file (YAML) at:
   ./config.yaml
   ./tmdb.yaml
   ~/.tmdb/config.yaml
   ~/.config/tmdb/config.yaml

3. Command line flags (override everything else)

Example config file:
   api_key: "your-api-key-here"
   timeout: "30s"
   max_retries: 3
   cache_ttl: "5m"
   log_level: "info"
   format: "table"

Get your API key from: https://www.themoviedb.org/settings/api`
}

// GetAPIKeyHelp returns help for getting an API key.
func GetAPIKeyHelp() string {
	return `To get a TMDB API key:

1. Go to https://www.themoviedb.org/
2. Create a free account if you don't have one
3. Go to https://www.themoviedb.org/settings/api
4. Request an API key (choose "Developer" option)
5. Fill out the application form
6. Once approved, copy your API key

Then set it via:
- Environment variable: export TMDB_API_KEY="your-key-here"
- Config file: add 'api_key: "your-key-here"' to config.yaml
- Or place it in ~/.tmdb/config.yaml`
}

// Helper function to check if slice contains string.
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

// Logging helpers.
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	SilentLevel
)

func ParseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "silent", "none":
		return SilentLevel
	default:
		return InfoLevel
	}
}

// Simple logger that respects the configuration.
type Logger struct {
	level LogLevel
}

func NewLogger(level LogLevel) *Logger {
	return &Logger{level: level}
}

func (l *Logger) Debug(msg string, args ...any) {
	if l.level <= DebugLevel {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+msg+"\n", args...)
	}
}

func (l *Logger) Info(msg string, args ...any) {
	if l.level <= InfoLevel {
		fmt.Fprintf(os.Stderr, "[INFO] "+msg+"\n", args...)
	}
}

func (l *Logger) Warn(msg string, args ...any) {
	if l.level <= WarnLevel {
		fmt.Fprintf(os.Stderr, "[WARN] "+msg+"\n", args...)
	}
}

func (l *Logger) Error(msg string, args ...any) {
	if l.level <= ErrorLevel {
		fmt.Fprintf(os.Stderr, "[ERROR] "+msg+"\n", args...)
	}
}
