package unit

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func formatInt(i int) string {
	return strconv.Itoa(i)
}

func TestValidateConfig(t *testing.T) {
	t.Run("valid config passes validation", func(t *testing.T) {
		config := fixtures.SampleConfig
		err := internal.ValidateConfig(config)
		assert.NoError(t, err)
	})

	t.Run("default config passes validation", func(t *testing.T) {
		config := internal.DefaultConfig()
		config.APIKey = "test-api-key"
		err := internal.ValidateConfig(config)
		assert.NoError(t, err)
	})

	t.Run("valid config variations pass validation", func(t *testing.T) {
		tests := []struct {
			name   string
			config internal.Config
		}{
			{
				name: "minimal valid config",
				config: helpers.NewConfigBuilder().
					WithAPIKey("test-key").
					WithTimeout(1 * time.Second).
					WithMaxRetries(0).
					WithCacheTTL(1 * time.Second).
					WithLogLevel("debug").
					WithFormat("json").
					Build(),
			},
			{
				name: "maximum values config",
				config: helpers.NewConfigBuilder().
					WithAPIKey("very-long-api-key-with-special-chars-123456789").
					WithTimeout(60 * time.Minute).
					WithMaxRetries(10).
					WithCacheTTL(24 * time.Hour).
					WithLogLevel("silent").
					WithFormat("csv").
					Build(),
			},
			{
				name: "different log levels",
				config: helpers.NewConfigBuilder().
					WithAPIKey("test").
					WithLogLevel("warn").
					Build(),
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				err := internal.ValidateConfig(test.config)
				assert.NoError(t, err)
			})
		}
	})
}

func TestValidateConfig_APIKeyValidation(t *testing.T) {
	t.Run("empty API key fails validation", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("").
			Build()

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TMDB API key is required")
		assert.Contains(t, err.Error(), "TMDB_API_KEY")
	})

	t.Run("whitespace-only API key fails validation", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("   ").
			Build()

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "TMDB API key is required")
	})

	t.Run("various valid API keys pass validation", func(t *testing.T) {
		validAPIKeys := []string{
			"a",                                   // Single character
			"test-api-key",                        // Standard format
			"1234567890abcdef",                    // Alphanumeric
			"api-key-with-dashes-and-numbers-123", // With dashes
			"api_key_with_underscores",            // With underscores
			"UPPERCASE-API-KEY",                   // Uppercase
			"MiXeD-cAsE-aPi-KeY",                  // Mixed case
			"key-with-special-chars!@#$%^&*()",    // Special characters
			"very-long-api-key-" + strings.Repeat("x", 100), // Very long
		}

		for _, apiKey := range validAPIKeys {
			t.Run("api_key_"+apiKey[:min(10, len(apiKey))], func(t *testing.T) {
				config := helpers.NewConfigBuilder().
					WithAPIKey(apiKey).
					Build()

				err := internal.ValidateConfig(config)
				assert.NoError(t, err)
			})
		}
	})
}

func TestValidateConfig_TimeoutValidation(t *testing.T) {
	t.Run("zero timeout fails validation", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("test").
			WithTimeout(0).
			Build()

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout must be positive")
		assert.Contains(t, err.Error(), "got 0s")
	})

	t.Run("negative timeout fails validation", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("test").
			WithTimeout(-5 * time.Second).
			Build()

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout must be positive")
		assert.Contains(t, err.Error(), "got -5s")
	})

	t.Run("valid timeouts pass validation", func(t *testing.T) {
		validTimeouts := []time.Duration{
			1 * time.Nanosecond,  // Minimum positive duration
			1 * time.Millisecond, // Short timeout
			30 * time.Second,     // Default timeout
			5 * time.Minute,      // Long timeout
			24 * time.Hour,       // Very long timeout
		}

		for _, timeout := range validTimeouts {
			t.Run("timeout_"+timeout.String(), func(t *testing.T) {
				config := helpers.NewConfigBuilder().
					WithAPIKey("test").
					WithTimeout(timeout).
					Build()

				err := internal.ValidateConfig(config)
				assert.NoError(t, err)
			})
		}
	})
}

func TestValidateConfig_MaxRetriesValidation(t *testing.T) {
	t.Run("negative max retries fails validation", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("test").
			WithMaxRetries(-1).
			Build()

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "max retries must be non-negative")
		assert.Contains(t, err.Error(), "got -1")
	})

	t.Run("very negative max retries fails validation", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("test").
			WithMaxRetries(-100).
			Build()

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "max retries must be non-negative")
		assert.Contains(t, err.Error(), "got -100")
	})

	t.Run("valid max retries pass validation", func(t *testing.T) {
		validRetries := []int{0, 1, 3, 5, 10, 100}

		for _, retries := range validRetries {
			t.Run("retries_"+formatInt(retries), func(t *testing.T) {
				config := helpers.NewConfigBuilder().
					WithAPIKey("test").
					WithMaxRetries(retries).
					Build()

				err := internal.ValidateConfig(config)
				assert.NoError(t, err)
			})
		}
	})
}

func TestValidateConfig_CacheTTLValidation(t *testing.T) {
	t.Run("zero cache TTL fails validation", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("test").
			WithCacheTTL(0).
			Build()

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache TTL must be positive")
		assert.Contains(t, err.Error(), "got 0s")
	})

	t.Run("negative cache TTL fails validation", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("test").
			WithCacheTTL(-10 * time.Minute).
			Build()

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache TTL must be positive")
		assert.Contains(t, err.Error(), "got -10m0s")
	})

	t.Run("valid cache TTL values pass validation", func(t *testing.T) {
		validCacheTTLs := []time.Duration{
			1 * time.Nanosecond, // Minimum positive duration
			1 * time.Second,     // Short cache
			5 * time.Minute,     // Default cache
			1 * time.Hour,       // Long cache
			7 * 24 * time.Hour,  // Very long cache
		}

		for _, cacheTTL := range validCacheTTLs {
			t.Run("cache_ttl_"+cacheTTL.String(), func(t *testing.T) {
				config := helpers.NewConfigBuilder().
					WithAPIKey("test").
					WithCacheTTL(cacheTTL).
					Build()

				err := internal.ValidateConfig(config)
				assert.NoError(t, err)
			})
		}
	})
}

func TestValidateConfig_LogLevelValidation(t *testing.T) {
	t.Run("invalid log level fails validation", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("test").
			WithLogLevel("invalid").
			Build()

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid log level 'invalid'")
		assert.Contains(t, err.Error(), "must be one of:")
		assert.Contains(t, err.Error(), "debug, info, warn, error, silent")
	})

	t.Run("empty log level fails validation", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("test").
			WithLogLevel("").
			Build()

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid log level ''")
	})

	t.Run("case sensitive log level validation", func(t *testing.T) {
		invalidCaseLogLevels := []string{"DEBUG", "Info", "WARN", "Error", "SILENT"}

		for _, logLevel := range invalidCaseLogLevels {
			t.Run("invalid_case_"+logLevel, func(t *testing.T) {
				config := helpers.NewConfigBuilder().
					WithAPIKey("test").
					WithLogLevel(logLevel).
					Build()

				err := internal.ValidateConfig(config)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid log level")
			})
		}
	})

	t.Run("valid log levels pass validation", func(t *testing.T) {
		validLogLevels := []string{"debug", "info", "warn", "error", "silent"}

		for _, logLevel := range validLogLevels {
			t.Run("valid_"+logLevel, func(t *testing.T) {
				config := helpers.NewConfigBuilder().
					WithAPIKey("test").
					WithLogLevel(logLevel).
					Build()

				err := internal.ValidateConfig(config)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("whitespace-only log level fails validation", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("test").
			WithLogLevel("   ").
			Build()

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid log level")
	})

	t.Run("log level with extra characters fails validation", func(t *testing.T) {
		invalidLogLevels := []string{
			"infox",
			"debugg",
			"warn1",
			"error ",
			" silent",
			"info-level",
		}

		for _, logLevel := range invalidLogLevels {
			t.Run("invalid_extra_"+logLevel, func(t *testing.T) {
				config := helpers.NewConfigBuilder().
					WithAPIKey("test").
					WithLogLevel(logLevel).
					Build()

				err := internal.ValidateConfig(config)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid log level")
			})
		}
	})
}

func TestValidateConfig_FormatValidation(t *testing.T) {
	t.Run("invalid format fails validation", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("test").
			WithFormat("invalid").
			Build()

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid default format 'invalid'")
		assert.Contains(t, err.Error(), "must be one of:")
		assert.Contains(t, err.Error(), "table, json, csv")
	})

	t.Run("empty format fails validation", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("test").
			WithFormat("").
			Build()

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid default format ''")
	})

	t.Run("case sensitive format validation", func(t *testing.T) {
		invalidCaseFormats := []string{"TABLE", "Json", "CSV"}

		for _, format := range invalidCaseFormats {
			t.Run("invalid_case_"+format, func(t *testing.T) {
				config := helpers.NewConfigBuilder().
					WithAPIKey("test").
					WithFormat(format).
					Build()

				err := internal.ValidateConfig(config)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid default format")
			})
		}
	})

	t.Run("valid formats pass validation", func(t *testing.T) {
		validFormats := []string{"table", "json", "csv"}

		for _, format := range validFormats {
			t.Run("valid_"+format, func(t *testing.T) {
				config := helpers.NewConfigBuilder().
					WithAPIKey("test").
					WithFormat(format).
					Build()

				err := internal.ValidateConfig(config)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("whitespace format fails validation", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("test").
			WithFormat("   ").
			Build()

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid default format")
	})

	t.Run("format with extra characters fails validation", func(t *testing.T) {
		invalidFormats := []string{
			"tablex",
			"jsonn",
			"csv1",
			"table ",
			" json",
			"csv-format",
		}

		for _, format := range invalidFormats {
			t.Run("invalid_extra_"+format, func(t *testing.T) {
				config := helpers.NewConfigBuilder().
					WithAPIKey("test").
					WithFormat(format).
					Build()

				err := internal.ValidateConfig(config)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid default format")
			})
		}
	})
}

func TestValidateConfig_MultipleErrors(t *testing.T) {
	t.Run("multiple validation errors return first error", func(t *testing.T) {
		config := internal.Config{
			APIKey:     "",        // Invalid: empty
			Timeout:    0,         // Invalid: zero
			MaxRetries: -1,        // Invalid: negative
			CacheTTL:   0,         // Invalid: zero
			LogLevel:   "invalid", // Invalid: unknown level
			Format:     "invalid", // Invalid: unknown format
		}

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		// Should return the first validation error (API key)
		assert.Contains(t, err.Error(), "TMDB API key is required")
	})

	t.Run("validation stops at first error", func(t *testing.T) {
		config := internal.Config{
			APIKey:     "valid-key",     // Valid
			Timeout:    0,               // Invalid: zero
			MaxRetries: -1,              // Invalid: negative
			CacheTTL:   1 * time.Second, // Valid
			LogLevel:   "invalid",       // Invalid: unknown level
			Format:     "table",         // Valid
		}

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		// Should return timeout error (first invalid field)
		assert.Contains(t, err.Error(), "timeout must be positive")
		assert.NotContains(t, err.Error(), "max retries")
	})
}

func TestValidateConfig_EdgeCases(t *testing.T) {
	t.Run("config with unicode characters", func(t *testing.T) {
		config := helpers.NewConfigBuilder().
			WithAPIKey("api-key-with-émojis-🔑-and-中文").
			WithBaseURL("https://api.example.com/测试").
			Build()

		err := internal.ValidateConfig(config)
		assert.NoError(t, err)
	})

	t.Run("config with very long values", func(t *testing.T) {
		longAPIKey := "test-" + strings.Repeat("x", 1000)
		config := helpers.NewConfigBuilder().
			WithAPIKey(longAPIKey).
			Build()

		err := internal.ValidateConfig(config)
		assert.NoError(t, err)
	})

	t.Run("boundary timeout values", func(t *testing.T) {
		tests := []struct {
			name    string
			timeout time.Duration
			valid   bool
		}{
			{"minimum valid", 1 * time.Nanosecond, true},
			{"just above zero", 1 * time.Microsecond, true},
			{"exactly zero", 0, false},
			{"just below zero", -1 * time.Nanosecond, false},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				config := helpers.NewConfigBuilder().
					WithAPIKey("test").
					WithTimeout(test.timeout).
					Build()

				err := internal.ValidateConfig(config)
				if test.valid {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			})
		}
	})

	t.Run("boundary retry values", func(t *testing.T) {
		tests := []struct {
			name       string
			maxRetries int
			valid      bool
		}{
			{"zero retries", 0, true},
			{"positive retries", 1, true},
			{"negative retries", -1, false},
			{"very negative", -100, false},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				config := helpers.NewConfigBuilder().
					WithAPIKey("test").
					WithMaxRetries(test.maxRetries).
					Build()

				err := internal.ValidateConfig(config)
				if test.valid {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			})
		}
	})
}

func TestContainsHelper(t *testing.T) {
	t.Run("contains function works correctly", func(t *testing.T) {
		// Note: This tests the internal contains function indirectly
		// through ValidateConfig since it's not exported
		validLogLevels := []string{"debug", "info", "warn", "error", "silent"}

		for _, level := range validLogLevels {
			config := helpers.NewConfigBuilder().
				WithAPIKey("test").
				WithLogLevel(level).
				Build()

			err := internal.ValidateConfig(config)
			assert.NoError(t, err, "Level %s should be valid", level)
		}

		invalidLevel := "nonexistent"
		config := helpers.NewConfigBuilder().
			WithAPIKey("test").
			WithLogLevel(invalidLevel).
			Build()

		err := internal.ValidateConfig(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid log level")
	})
}
