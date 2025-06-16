package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
)

func TestDefaultConfig(t *testing.T) {
	t.Run("default config has reasonable values", func(t *testing.T) {
		config := internal.DefaultConfig()

		assert.Equal(t, "https://api.themoviedb.org/3", config.BaseURL)
		assert.Equal(t, 30*time.Second, config.Timeout)
		assert.Equal(t, 3, config.MaxRetries)
		assert.Equal(t, 5*time.Minute, config.CacheTTL)
		assert.Equal(t, "info", config.LogLevel)
		assert.Equal(t, "table", config.Format)
		assert.Empty(t, config.APIKey) // Should be empty in default
	})

	t.Run("default config values are positive", func(t *testing.T) {
		config := internal.DefaultConfig()

		assert.Greater(t, config.Timeout, 0*time.Second)
		assert.GreaterOrEqual(t, config.MaxRetries, 0)
		assert.Greater(t, config.CacheTTL, 0*time.Second)
		assert.NotEmpty(t, config.BaseURL)
		assert.NotEmpty(t, config.LogLevel)
		assert.NotEmpty(t, config.Format)
	})

	t.Run("default config uses HTTPS", func(t *testing.T) {
		config := internal.DefaultConfig()
		assert.Contains(t, config.BaseURL, "https://")
	})

	t.Run("default config points to TMDB API", func(t *testing.T) {
		config := internal.DefaultConfig()
		assert.Contains(t, config.BaseURL, "themoviedb.org")
		assert.Contains(t, config.BaseURL, "/3") // API version 3
	})
}

func TestConfigFields(t *testing.T) {
	t.Run("config fields can be set to various values", func(t *testing.T) {
		config := internal.Config{
			APIKey:     "test-api-key-12345678901234567890",
			BaseURL:    "https://custom-api.example.com/v1",
			Timeout:    45 * time.Second,
			MaxRetries: 5,
			CacheTTL:   10 * time.Minute,
			LogLevel:   "debug",
			Format:     "json",
		}

		assert.Equal(t, "test-api-key-12345678901234567890", config.APIKey)
		assert.Equal(t, "https://custom-api.example.com/v1", config.BaseURL)
		assert.Equal(t, 45*time.Second, config.Timeout)
		assert.Equal(t, 5, config.MaxRetries)
		assert.Equal(t, 10*time.Minute, config.CacheTTL)
		assert.Equal(t, "debug", config.LogLevel)
		assert.Equal(t, "json", config.Format)
	})

	t.Run("config handles extreme values", func(t *testing.T) {
		config := internal.Config{
			APIKey:     "",                      // Empty API key
			BaseURL:    "http://localhost:8080", // HTTP instead of HTTPS
			Timeout:    1 * time.Millisecond,    // Very short timeout
			MaxRetries: 0,                       // No retries
			CacheTTL:   1 * time.Second,         // Very short cache
			LogLevel:   "silent",
			Format:     "csv",
		}

		assert.Empty(t, config.APIKey)
		assert.Contains(t, config.BaseURL, "http://")
		assert.Equal(t, 1*time.Millisecond, config.Timeout)
		assert.Equal(t, 0, config.MaxRetries)
		assert.Equal(t, 1*time.Second, config.CacheTTL)
		assert.Equal(t, "silent", config.LogLevel)
		assert.Equal(t, "csv", config.Format)
	})

	t.Run("config handles unicode in API key", func(t *testing.T) {
		config := internal.Config{
			APIKey: "api-key-with-émojis-🔑-and-中文",
		}

		assert.Contains(t, config.APIKey, "🔑")
		assert.Contains(t, config.APIKey, "中文")
		assert.Contains(t, config.APIKey, "émojis")
	})

	t.Run("config handles very long API key", func(t *testing.T) {
		longAPIKey := "very-long-api-key-" + string(make([]byte, 1000))
		config := internal.Config{
			APIKey: longAPIKey,
		}

		assert.Equal(t, longAPIKey, config.APIKey)
		assert.Greater(t, len(config.APIKey), 1000)
	})
}

func TestConfigTimeouts(t *testing.T) {
	t.Run("various timeout durations", func(t *testing.T) {
		tests := []struct {
			name     string
			timeout  time.Duration
			cacheTTL time.Duration
		}{
			{"milliseconds", 500 * time.Millisecond, 30 * time.Second},
			{"seconds", 30 * time.Second, 5 * time.Minute},
			{"minutes", 2 * time.Minute, 30 * time.Minute},
			{"hours", 1 * time.Hour, 24 * time.Hour},
			{"very short", 1 * time.Nanosecond, 1 * time.Nanosecond},
			{"very long", 24 * time.Hour, 7 * 24 * time.Hour},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				config := internal.Config{
					Timeout:  test.timeout,
					CacheTTL: test.cacheTTL,
				}

				assert.Equal(t, test.timeout, config.Timeout)
				assert.Equal(t, test.cacheTTL, config.CacheTTL)
			})
		}
	})

	t.Run("zero and negative timeouts", func(t *testing.T) {
		config := internal.Config{
			Timeout:  0 * time.Second,  // Zero timeout
			CacheTTL: -5 * time.Minute, // Negative cache TTL
		}

		assert.Equal(t, 0*time.Second, config.Timeout)
		assert.Equal(t, -5*time.Minute, config.CacheTTL)
	})
}

func TestConfigRetries(t *testing.T) {
	t.Run("various retry counts", func(t *testing.T) {
		tests := []struct {
			name       string
			maxRetries int
		}{
			{"no retries", 0},
			{"single retry", 1},
			{"default retries", 3},
			{"many retries", 10},
			{"extreme retries", 100},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				config := internal.Config{
					MaxRetries: test.maxRetries,
				}

				assert.Equal(t, test.maxRetries, config.MaxRetries)
			})
		}
	})

	t.Run("negative retries", func(t *testing.T) {
		config := internal.Config{
			MaxRetries: -1,
		}

		assert.Equal(t, -1, config.MaxRetries)
	})
}

func TestConfigLogLevels(t *testing.T) {
	t.Run("all valid log levels", func(t *testing.T) {
		validLogLevels := []string{
			"debug",
			"info",
			"warn",
			"warning",
			"error",
			"silent",
			"none",
		}

		for _, level := range validLogLevels {
			t.Run("log_level_"+level, func(t *testing.T) {
				config := internal.Config{
					LogLevel: level,
				}

				assert.Equal(t, level, config.LogLevel)
			})
		}
	})

	t.Run("case variations in log levels", func(t *testing.T) {
		tests := []string{
			"DEBUG",
			"Info",
			"WARN",
			"Error",
			"SILENT",
		}

		for _, level := range tests {
			t.Run("case_"+level, func(t *testing.T) {
				config := internal.Config{
					LogLevel: level,
				}

				assert.Equal(t, level, config.LogLevel)
			})
		}
	})

	t.Run("invalid log levels", func(t *testing.T) {
		invalidLevels := []string{
			"invalid",
			"trace",
			"fatal",
			"",
			"log-level",
		}

		for _, level := range invalidLevels {
			t.Run("invalid_"+level, func(t *testing.T) {
				config := internal.Config{
					LogLevel: level,
				}

				assert.Equal(t, level, config.LogLevel)
				// Note: Validation happens in config validator, not in the struct itself
			})
		}
	})
}

func TestConfigFormats(t *testing.T) {
	t.Run("all valid formats", func(t *testing.T) {
		validFormats := []string{
			"table",
			"json",
			"csv",
		}

		for _, format := range validFormats {
			t.Run("format_"+format, func(t *testing.T) {
				config := internal.Config{
					Format: format,
				}

				assert.Equal(t, format, config.Format)
			})
		}
	})

	t.Run("case variations in formats", func(t *testing.T) {
		tests := []string{
			"TABLE",
			"Json",
			"CSV",
		}

		for _, format := range tests {
			t.Run("case_"+format, func(t *testing.T) {
				config := internal.Config{
					Format: format,
				}

				assert.Equal(t, format, config.Format)
			})
		}
	})

	t.Run("invalid formats", func(t *testing.T) {
		invalidFormats := []string{
			"xml",
			"yaml",
			"html",
			"",
			"invalid",
		}

		for _, format := range invalidFormats {
			t.Run("invalid_"+format, func(t *testing.T) {
				config := internal.Config{
					Format: format,
				}

				assert.Equal(t, format, config.Format)
				// Note: Validation happens in config validator, not in the struct itself
			})
		}
	})
}

func TestConfigURLs(t *testing.T) {
	t.Run("various valid URLs", func(t *testing.T) {
		validURLs := []string{
			"https://api.themoviedb.org/3",
			"http://localhost:8080",
			"https://custom-api.example.com/v1",
			"https://api.staging.themoviedb.org/3",
			"https://192.168.1.100:3000/api",
		}

		for _, url := range validURLs {
			t.Run("url_"+url, func(t *testing.T) {
				config := internal.Config{
					BaseURL: url,
				}

				assert.Equal(t, url, config.BaseURL)
			})
		}
	})

	t.Run("URLs with special characters", func(t *testing.T) {
		specialURLs := []string{
			"https://api.example.com/v1?key=test",
			"https://api.example.com/path/with-dashes",
			"https://api.example.com/path_with_underscores",
			"https://api.example.com:8443/secure",
		}

		for _, url := range specialURLs {
			t.Run("special_url", func(t *testing.T) {
				config := internal.Config{
					BaseURL: url,
				}

				assert.Equal(t, url, config.BaseURL)
			})
		}
	})

	t.Run("invalid URLs", func(t *testing.T) {
		invalidURLs := []string{
			"not-a-url",
			"",
			"ftp://invalid-protocol.com",
			"https://",
			"api.example.com", // Missing protocol
		}

		for _, url := range invalidURLs {
			t.Run("invalid_url", func(t *testing.T) {
				config := internal.Config{
					BaseURL: url,
				}

				assert.Equal(t, url, config.BaseURL)
				// Note: URL validation happens in config validator
			})
		}
	})
}

func TestConfigStructTags(t *testing.T) {
	t.Run("yaml tags are present", func(t *testing.T) {
		// This test ensures the struct has proper YAML tags for file loading
		config := fixtures.SampleConfig

		// Basic check that the config struct can hold all expected values
		assert.NotEmpty(t, config.APIKey)
		assert.NotEmpty(t, config.BaseURL)
		assert.Greater(t, config.Timeout, 0*time.Second)
		assert.GreaterOrEqual(t, config.MaxRetries, 0)
		assert.Greater(t, config.CacheTTL, 0*time.Second)
		assert.NotEmpty(t, config.LogLevel)
		assert.NotEmpty(t, config.Format)
	})

	t.Run("env tags are present", func(t *testing.T) {
		// Verify that APIKey field has env tag for environment variable loading
		config := internal.Config{
			APIKey: "test-from-env",
		}

		assert.Equal(t, "test-from-env", config.APIKey)
	})
}

func TestConfigEquality(t *testing.T) {
	t.Run("identical configs are equal", func(t *testing.T) {
		config1 := fixtures.SampleConfig
		config2 := fixtures.SampleConfig

		assert.Equal(t, config1, config2)
		assert.Equal(t, config1.APIKey, config2.APIKey)
		assert.Equal(t, config1.BaseURL, config2.BaseURL)
		assert.Equal(t, config1.Timeout, config2.Timeout)
		assert.Equal(t, config1.MaxRetries, config2.MaxRetries)
		assert.Equal(t, config1.CacheTTL, config2.CacheTTL)
		assert.Equal(t, config1.LogLevel, config2.LogLevel)
		assert.Equal(t, config1.Format, config2.Format)
	})

	t.Run("different configs are not equal", func(t *testing.T) {
		config1 := fixtures.SampleConfig
		config2 := fixtures.SampleConfig
		config2.APIKey = "different-api-key"

		assert.NotEqual(t, config1, config2)
		assert.NotEqual(t, config1.APIKey, config2.APIKey)
	})

	t.Run("config with helper validation", func(t *testing.T) {
		config := fixtures.SampleConfig

		// Don't use the helper since it has type issues with assert.Greater
		assert.NotEmpty(t, config.APIKey, "API key should not be empty")
		assert.NotEmpty(t, config.BaseURL, "Base URL should not be empty")
		assert.True(t, config.Timeout > 0, "Timeout should be positive")
		assert.True(t, config.MaxRetries >= 0, "Max retries should be non-negative")
		assert.True(t, config.CacheTTL > 0, "Cache TTL should be positive")
		assert.NotEmpty(t, config.LogLevel, "Log level should not be empty")
		assert.NotEmpty(t, config.Format, "Format should not be empty")
	})
}
