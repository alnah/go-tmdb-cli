package unit

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestLoadConfig(t *testing.T) {
	t.Run("loads default config with API key from environment", func(t *testing.T) {
		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY": "test-api-key",
		}, func() {
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			assert.Equal(t, "test-api-key", config.APIKey)
			assert.Equal(t, "https://api.themoviedb.org/3", config.BaseURL)
			assert.Equal(t, 30*time.Second, config.Timeout)
			assert.Equal(t, 3, config.MaxRetries)
			assert.Equal(t, 5*time.Minute, config.CacheTTL)
			assert.Equal(t, "info", config.LogLevel)
			assert.Equal(t, "table", config.Format)
		})
	})

	t.Run("loads config from environment variables", func(t *testing.T) {
		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY":     "env-api-key",
			"TMDB_BASE_URL":    "https://custom.api.com/v1",
			"TMDB_TIMEOUT":     "60s",
			"TMDB_MAX_RETRIES": "5",
			"TMDB_CACHE_TTL":   "10m",
			"TMDB_LOG_LEVEL":   "debug",
			"TMDB_FORMAT":      "json",
		}, func() {
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			assert.Equal(t, "env-api-key", config.APIKey)
			assert.Equal(t, "https://custom.api.com/v1", config.BaseURL)
			assert.Equal(t, 60*time.Second, config.Timeout)
			assert.Equal(t, 5, config.MaxRetries)
			assert.Equal(t, 10*time.Minute, config.CacheTTL)
			assert.Equal(t, "debug", config.LogLevel)
			assert.Equal(t, "json", config.Format)
		})
	})

	t.Run("validation failure returns error", func(t *testing.T) {
		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY": "", // Invalid: empty API key
		}, func() {
			_, err := internal.LoadConfig()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "configuration validation failed")
			assert.Contains(t, err.Error(), "TMDB API key is required")
		})
	})

	t.Run("loads config without environment variables", func(t *testing.T) {
		helpers.WithCleanEnvironment(t, []string{
			"TMDB_API_KEY", "TMDB_BASE_URL", "TMDB_TIMEOUT",
			"TMDB_MAX_RETRIES", "TMDB_CACHE_TTL", "TMDB_LOG_LEVEL", "TMDB_FORMAT",
		}, func() {
			// This should fail validation due to missing API key
			_, err := internal.LoadConfig()
			require.Error(t, err)
			assert.Contains(t, err.Error(), "TMDB API key is required")
		})
	})
}

func TestLoadFromEnvironment(t *testing.T) {
	t.Run("loads all environment variables correctly", func(t *testing.T) {
		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY":     "test-key-from-env",
			"TMDB_BASE_URL":    "https://env.api.com",
			"TMDB_TIMEOUT":     "45s",
			"TMDB_MAX_RETRIES": "7",
			"TMDB_CACHE_TTL":   "15m",
			"TMDB_LOG_LEVEL":   "warn",
			"TMDB_FORMAT":      "csv",
		}, func() {
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			assert.Equal(t, "test-key-from-env", config.APIKey)
			assert.Equal(t, "https://env.api.com", config.BaseURL)
			assert.Equal(t, 45*time.Second, config.Timeout)
			assert.Equal(t, 7, config.MaxRetries)
			assert.Equal(t, 15*time.Minute, config.CacheTTL)
			assert.Equal(t, "warn", config.LogLevel)
			assert.Equal(t, "csv", config.Format)
		})
	})

	t.Run("partial environment variables override defaults", func(t *testing.T) {
		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY":   "partial-key",
			"TMDB_LOG_LEVEL": "error",
		}, func() {
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			// Overridden values
			assert.Equal(t, "partial-key", config.APIKey)
			assert.Equal(t, "error", config.LogLevel)

			// Default values
			assert.Equal(t, "https://api.themoviedb.org/3", config.BaseURL)
			assert.Equal(t, 30*time.Second, config.Timeout)
			assert.Equal(t, 3, config.MaxRetries)
			assert.Equal(t, 5*time.Minute, config.CacheTTL)
			assert.Equal(t, "table", config.Format)
		})
	})

	t.Run("empty environment variables are ignored", func(t *testing.T) {
		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY":     "valid-key",
			"TMDB_BASE_URL":    "", // Empty, should be ignored
			"TMDB_TIMEOUT":     "", // Empty, should be ignored
			"TMDB_MAX_RETRIES": "", // Empty, should be ignored
			"TMDB_CACHE_TTL":   "", // Empty, should be ignored
			"TMDB_LOG_LEVEL":   "", // Empty, should be ignored
			"TMDB_FORMAT":      "", // Empty, should be ignored
		}, func() {
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			// Only API key should be overridden
			assert.Equal(t, "valid-key", config.APIKey)

			// All others should remain defaults
			assert.Equal(t, "https://api.themoviedb.org/3", config.BaseURL)
			assert.Equal(t, 30*time.Second, config.Timeout)
			assert.Equal(t, 3, config.MaxRetries)
			assert.Equal(t, 5*time.Minute, config.CacheTTL)
			assert.Equal(t, "info", config.LogLevel)
			assert.Equal(t, "table", config.Format)
		})
	})

	t.Run("invalid duration values are ignored", func(t *testing.T) {
		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY":   "valid-key",
			"TMDB_TIMEOUT":   "invalid-duration",
			"TMDB_CACHE_TTL": "not-a-duration",
		}, func() {
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			// Should use default values for invalid durations
			assert.Equal(t, 30*time.Second, config.Timeout)
			assert.Equal(t, 5*time.Minute, config.CacheTTL)
		})
	})

	t.Run("invalid integer values are ignored", func(t *testing.T) {
		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY":     "valid-key",
			"TMDB_MAX_RETRIES": "not-a-number",
		}, func() {
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			// Should use default value for invalid integer
			assert.Equal(t, 3, config.MaxRetries)
		})
	})

	t.Run("edge case environment values", func(t *testing.T) {
		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY":     "key-with-special-chars!@#$%",
			"TMDB_BASE_URL":    "https://api.example.com:8443/path?query=value",
			"TMDB_TIMEOUT":     "1ns",
			"TMDB_MAX_RETRIES": "0",
			"TMDB_CACHE_TTL":   "24h",
			"TMDB_LOG_LEVEL":   "silent",
			"TMDB_FORMAT":      "json",
		}, func() {
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			assert.Equal(t, "key-with-special-chars!@#$%", config.APIKey)
			assert.Equal(t, "https://api.example.com:8443/path?query=value", config.BaseURL)
			assert.Equal(t, 1*time.Nanosecond, config.Timeout)
			assert.Equal(t, 0, config.MaxRetries)
			assert.Equal(t, 24*time.Hour, config.CacheTTL)
			assert.Equal(t, "silent", config.LogLevel)
			assert.Equal(t, "json", config.Format)
		})
	})

	t.Run("whitespace in environment variables", func(t *testing.T) {
		helpers.WithTempDir(t, "tmdb-test-*", func(tempDir string) {
			originalDir, err := os.Getwd()
			require.NoError(t, err)
			err = os.Chdir(tempDir)
			require.NoError(t, err)
			defer func() {
				_ = os.Chdir(originalDir)
			}()

			helpers.WithMockEnvironment(t, map[string]string{
				"TMDB_API_KEY":   "  key-with-spaces  ",
				"TMDB_LOG_LEVEL": "  info  ",
				"TMDB_FORMAT":    "  table  ",
			}, func() {
				config, err := internal.LoadConfig()

				// The implementation should trim whitespace and load successfully
				require.NoError(t, err)
				assert.Equal(t, "  key-with-spaces  ", config.APIKey) // API key preserves spaces
				assert.Equal(t, "info", config.LogLevel)              // Log level is trimmed
				assert.Equal(t, "table", config.Format)               // Format is trimmed
			})
		})
	})
}

func TestLoadConfigFile(t *testing.T) {
	t.Run("loads config from YAML file", func(t *testing.T) {
		yamlContent := `
api_key: "file-api-key"
base_url: "https://file.api.com"
timeout: "120s"
max_retries: 2
cache_ttl: "30m"
log_level: "debug"
format: "json"
`

		helpers.WithTempFile(t, yamlContent, func(configPath string) {
			helpers.WithMockEnvironment(t, map[string]string{
				"TMDB_API_KEY": "env-key-override", // Should override file
			}, func() {
				// Test requires config file to be in current directory
				workingDir, err := os.Getwd()
				require.NoError(t, err)

				configFile := filepath.Join(workingDir, "config.yaml")
				err = os.WriteFile(configFile, []byte(yamlContent), 0o644)
				require.NoError(t, err)
				defer os.Remove(configFile)

				config, err := internal.LoadConfig()
				require.NoError(t, err)

				// Environment should override file
				assert.Equal(t, "env-key-override", config.APIKey)

				// File values should be loaded for non-overridden fields
				assert.Equal(t, "https://file.api.com", config.BaseURL)
				assert.Equal(t, 120*time.Second, config.Timeout)
				assert.Equal(t, 2, config.MaxRetries)
				assert.Equal(t, 30*time.Minute, config.CacheTTL)
				assert.Equal(t, "debug", config.LogLevel)
				assert.Equal(t, "json", config.Format)
			})
		})
	})

	t.Run("loads config from tmdb.yaml file", func(t *testing.T) {
		yamlContent := `
api_key: "tmdb-yaml-key"
log_level: "warn"
`

		workingDir, err := os.Getwd()
		require.NoError(t, err)

		configFile := filepath.Join(workingDir, "tmdb.yaml")
		err = os.WriteFile(configFile, []byte(yamlContent), 0o644)
		require.NoError(t, err)
		defer os.Remove(configFile)

		config, err := internal.LoadConfig()
		require.NoError(t, err)

		assert.Equal(t, "tmdb-yaml-key", config.APIKey)
		assert.Equal(t, "warn", config.LogLevel)
	})

	t.Run("handles missing config file gracefully", func(t *testing.T) {
		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY": "env-only-key",
		}, func() {
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			// Should load from environment and defaults
			assert.Equal(t, "env-only-key", config.APIKey)
			assert.Equal(t, "https://api.themoviedb.org/3", config.BaseURL)
		})
	})

	t.Run("handles invalid YAML gracefully", func(t *testing.T) {
		invalidYAML := `
api_key: "valid-key"
invalid: yaml: content: [
  - missing
  closing
`

		workingDir, err := os.Getwd()
		require.NoError(t, err)

		configFile := filepath.Join(workingDir, "config.yaml")
		err = os.WriteFile(configFile, []byte(invalidYAML), 0o644)
		require.NoError(t, err)
		defer os.Remove(configFile)

		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY": "fallback-key",
		}, func() {
			// Should still succeed with defaults and environment
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			assert.Equal(t, "fallback-key", config.APIKey)
		})
	})

	t.Run("partial YAML config merges with defaults", func(t *testing.T) {
		partialYAML := `
api_key: "partial-key"
timeout: "90s"
`

		workingDir, err := os.Getwd()
		require.NoError(t, err)

		configFile := filepath.Join(workingDir, "config.yaml")
		err = os.WriteFile(configFile, []byte(partialYAML), 0o644)
		require.NoError(t, err)
		defer os.Remove(configFile)

		config, err := internal.LoadConfig()
		require.NoError(t, err)

		// From file
		assert.Equal(t, "partial-key", config.APIKey)
		assert.Equal(t, 90*time.Second, config.Timeout)

		// From defaults
		assert.Equal(t, "https://api.themoviedb.org/3", config.BaseURL)
		assert.Equal(t, 3, config.MaxRetries)
		assert.Equal(t, 5*time.Minute, config.CacheTTL)
		assert.Equal(t, "info", config.LogLevel)
		assert.Equal(t, "table", config.Format)
	})

	t.Run("environment overrides config file", func(t *testing.T) {
		yamlContent := `
api_key: "file-key"
base_url: "https://file.api.com"
timeout: "60s"
max_retries: 1
cache_ttl: "10m"
log_level: "error"
format: "csv"
`

		workingDir, err := os.Getwd()
		require.NoError(t, err)

		configFile := filepath.Join(workingDir, "config.yaml")
		err = os.WriteFile(configFile, []byte(yamlContent), 0o644)
		require.NoError(t, err)
		defer os.Remove(configFile)

		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY":   "env-override-key",
			"TMDB_TIMEOUT":   "30s",
			"TMDB_LOG_LEVEL": "debug",
		}, func() {
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			// Environment overrides
			assert.Equal(t, "env-override-key", config.APIKey)
			assert.Equal(t, 30*time.Second, config.Timeout)
			assert.Equal(t, "debug", config.LogLevel)

			// File values for non-overridden fields
			assert.Equal(t, "https://file.api.com", config.BaseURL)
			assert.Equal(t, 1, config.MaxRetries)
			assert.Equal(t, 10*time.Minute, config.CacheTTL)
			assert.Equal(t, "csv", config.Format)
		})
	})
}

func TestLoadConfig_EdgeCases(t *testing.T) {
	t.Run("very long environment values", func(t *testing.T) {
		longAPIKey := "key-" + strings.Repeat("x", 1000)
		longBaseURL := "https://very-long-subdomain-" + strings.Repeat(
			"x",
			100,
		) + ".example.com/api/v1"

		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY":  longAPIKey,
			"TMDB_BASE_URL": longBaseURL,
		}, func() {
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			assert.Equal(t, longAPIKey, config.APIKey)
			assert.Equal(t, longBaseURL, config.BaseURL)
		})
	})

	t.Run("unicode in environment values", func(t *testing.T) {
		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY":  "key-with-émojis-🔑-and-中文",
			"TMDB_BASE_URL": "https://api.example.com/测试/路径",
		}, func() {
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			assert.Equal(t, "key-with-émojis-🔑-and-中文", config.APIKey)
			assert.Equal(t, "https://api.example.com/测试/路径", config.BaseURL)
		})
	})

	t.Run("extreme duration values in environment", func(t *testing.T) {
		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY":   "test-key",
			"TMDB_TIMEOUT":   "1ns",
			"TMDB_CACHE_TTL": "8760h", // 1 year
		}, func() {
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			assert.Equal(t, 1*time.Nanosecond, config.Timeout)
			assert.Equal(t, 8760*time.Hour, config.CacheTTL)
		})
	})

	t.Run("extreme integer values in environment", func(t *testing.T) {
		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY":     "test-key",
			"TMDB_MAX_RETRIES": strconv.Itoa(int(^uint(0) >> 1)), // Max int
		}, func() {
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			assert.Equal(t, int(^uint(0)>>1), config.MaxRetries)
		})
	})

	t.Run("precedence order: env > file > defaults", func(t *testing.T) {
		yamlContent := `
api_key: "file-key"
timeout: "60s"
log_level: "warn"
`

		workingDir, err := os.Getwd()
		require.NoError(t, err)

		configFile := filepath.Join(workingDir, "config.yaml")
		err = os.WriteFile(configFile, []byte(yamlContent), 0o644)
		require.NoError(t, err)
		defer os.Remove(configFile)

		helpers.WithMockEnvironment(t, map[string]string{
			"TMDB_API_KEY": "env-key", // Should override file
			// TMDB_TIMEOUT not set, should use file value
			// TMDB_LOG_LEVEL not set, should use file value
			// Other values should use defaults
		}, func() {
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			// Environment override
			assert.Equal(t, "env-key", config.APIKey)

			// File values
			assert.Equal(t, 60*time.Second, config.Timeout)
			assert.Equal(t, "warn", config.LogLevel)

			// Default values
			assert.Equal(t, "https://api.themoviedb.org/3", config.BaseURL)
			assert.Equal(t, 3, config.MaxRetries)
			assert.Equal(t, 5*time.Minute, config.CacheTTL)
			assert.Equal(t, "table", config.Format)
		})
	})

	t.Run("config validation with all edge cases", func(t *testing.T) {
		tests := []struct {
			name        string
			envVars     map[string]string
			expectError bool
			errorMsg    string
		}{
			{
				name: "valid edge case config",
				envVars: map[string]string{
					"TMDB_API_KEY":     "a",
					"TMDB_TIMEOUT":     "1ns",
					"TMDB_MAX_RETRIES": "0",
					"TMDB_CACHE_TTL":   "1ns",
					"TMDB_LOG_LEVEL":   "silent",
					"TMDB_FORMAT":      "csv",
				},
				expectError: false,
			},
			{
				name: "invalid timeout from environment",
				envVars: map[string]string{
					"TMDB_API_KEY": "valid",
					"TMDB_TIMEOUT": "0s",
				},
				expectError: true,
				errorMsg:    "timeout must be positive",
			},
			{
				name: "invalid retries from environment",
				envVars: map[string]string{
					"TMDB_API_KEY":     "valid",
					"TMDB_MAX_RETRIES": "-1",
				},
				expectError: true,
				errorMsg:    "max retries must be non-negative",
			},
			{
				name: "invalid cache TTL from environment",
				envVars: map[string]string{
					"TMDB_API_KEY":   "valid",
					"TMDB_CACHE_TTL": "0s",
				},
				expectError: true,
				errorMsg:    "cache TTL must be positive",
			},
			{
				name: "invalid log level from environment",
				envVars: map[string]string{
					"TMDB_API_KEY":   "valid",
					"TMDB_LOG_LEVEL": "invalid",
				},
				expectError: true,
				errorMsg:    "invalid log level",
			},
			{
				name: "invalid format from environment",
				envVars: map[string]string{
					"TMDB_API_KEY": "valid",
					"TMDB_FORMAT":  "invalid",
				},
				expectError: true,
				errorMsg:    "invalid default format",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				helpers.WithMockEnvironment(t, test.envVars, func() {
					config, err := internal.LoadConfig()

					if test.expectError {
						require.Error(t, err)
						assert.Contains(t, err.Error(), test.errorMsg)
					} else {
						require.NoError(t, err)
						assert.NotEmpty(t, config.APIKey)
					}
				})
			})
		}
	})
}
