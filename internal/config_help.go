// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

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
	const dirPerm = 0o755
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
