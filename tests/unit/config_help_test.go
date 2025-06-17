package unit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestCreateExampleConfig(t *testing.T) {
	t.Run("creates valid example config file", func(t *testing.T) {
		helpers.WithTempDir(t, "tmdb-config-test-*", func(tempDir string) {
			configPath := filepath.Join(tempDir, "example-config.yaml")

			err := internal.CreateExampleConfig(configPath)
			require.NoError(t, err)

			// Verify file exists
			helpers.AssertFileExists(t, configPath)

			// Read and verify content
			content, err := os.ReadFile(configPath)
			require.NoError(t, err)

			validateExampleConfigContent(t, string(content))
		})
	})

	t.Run("creates directory if it does not exist", func(t *testing.T) {
		helpers.WithTempDir(t, "tmdb-config-test-*", func(tempDir string) {
			nestedPath := filepath.Join(tempDir, "nested", "deep", "config.yaml")

			err := internal.CreateExampleConfig(nestedPath)
			require.NoError(t, err)

			helpers.AssertFileExists(t, nestedPath)

			// Verify parent directories were created
			assert.DirExists(t, filepath.Dir(nestedPath))
		})
	})

	t.Run("overwrites existing config file", func(t *testing.T) {
		helpers.WithTempDir(t, "tmdb-config-test-*", func(tempDir string) {
			configPath := filepath.Join(tempDir, "config.yaml")

			// Create existing file with different content
			existingContent := "existing: content"
			err := os.WriteFile(configPath, []byte(existingContent), 0o644)
			require.NoError(t, err)

			// Create example config (should overwrite)
			err = internal.CreateExampleConfig(configPath)
			require.NoError(t, err)

			// Verify content was overwritten
			content, err := os.ReadFile(configPath)
			require.NoError(t, err)

			contentStr := string(content)
			assert.NotContains(t, contentStr, "existing: content")
			assert.Contains(t, contentStr, "# TMDB CLI Configuration")
		})
	})

	t.Run("creates valid YAML that can be parsed", func(t *testing.T) {
		helpers.WithTempDir(t, "tmdb-config-test-*", func(tempDir string) {
			configPath := filepath.Join(tempDir, "config.yaml")

			err := internal.CreateExampleConfig(configPath)
			require.NoError(t, err)

			validateExampleConfigYAML(t, configPath)
		})
	})

	t.Run("example config uses default values", func(t *testing.T) {
		helpers.WithTempDir(t, "tmdb-config-test-*", func(tempDir string) {
			configPath := filepath.Join(tempDir, "config.yaml")

			err := internal.CreateExampleConfig(configPath)
			require.NoError(t, err)

			validateExampleConfigDefaults(t, configPath)
		})
	})

	t.Run("handles permission errors gracefully", func(t *testing.T) {
		helpers.WithTempDir(t, "tmdb-config-test-*", func(tempDir string) {
			// Create a read-only directory
			readOnlyDir := filepath.Join(tempDir, "readonly")
			err := os.Mkdir(readOnlyDir, 0o555) // #nosec
			require.NoError(t, err)

			configPath := filepath.Join(readOnlyDir, "config.yaml")

			err = internal.CreateExampleConfig(configPath)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "write config file")
		})
	})

	t.Run("handles invalid path gracefully", func(t *testing.T) {
		// Try to create config in a non-existent directory with permission issues
		invalidPath := "/root/non-existent/config.yaml"

		err := internal.CreateExampleConfig(invalidPath)
		assert.Error(t, err)
	})

	t.Run("preserves file permissions", func(t *testing.T) {
		helpers.WithTempDir(t, "tmdb-config-test-*", func(tempDir string) {
			configPath := filepath.Join(tempDir, "config.yaml")

			err := internal.CreateExampleConfig(configPath)
			require.NoError(t, err)

			// Check file permissions
			fileInfo, err := os.Stat(configPath)
			require.NoError(t, err)

			// Should be readable and writable by owner, readable by group/others
			assert.Equal(t, os.FileMode(0o644), fileInfo.Mode().Perm())
		})
	})

	t.Run("example config file structure", func(t *testing.T) {
		helpers.WithTempDir(t, "tmdb-config-test-*", func(tempDir string) {
			configPath := filepath.Join(tempDir, "config.yaml")

			err := internal.CreateExampleConfig(configPath)
			require.NoError(t, err)

			content, err := os.ReadFile(configPath)
			require.NoError(t, err)

			validateExampleConfigStructure(t, string(content))
		})
	})
}

func TestGetConfigHelp(t *testing.T) {
	t.Run("returns comprehensive help text", func(t *testing.T) {
		help := internal.GetConfigHelp()
		assert.NotEmpty(t, help)

		// Check for main sections
		assert.Contains(t, help, "Configuration can be set via:")
		assert.Contains(t, help, "1. Environment variables:")
		assert.Contains(t, help, "2. Config file (YAML) at:")
		assert.Contains(t, help, "3. Command line flags")

		// Check for environment variables documentation
		assert.Contains(t, help, "TMDB_API_KEY")
		assert.Contains(t, help, "TMDB_BASE_URL")
		assert.Contains(t, help, "TMDB_TIMEOUT")
		assert.Contains(t, help, "TMDB_MAX_RETRIES")
		assert.Contains(t, help, "TMDB_CACHE_TTL")
		assert.Contains(t, help, "TMDB_LOG_LEVEL")
		assert.Contains(t, help, "TMDB_FORMAT")

		// Check for config file paths
		assert.Contains(t, help, "./config.yaml")
		assert.Contains(t, help, "./tmdb.yaml")
		assert.Contains(t, help, "~/.tmdb/config.yaml")
		assert.Contains(t, help, "~/.config/tmdb/config.yaml")

		// Check for example config
		assert.Contains(t, help, "Example config file:")
		assert.Contains(t, help, "api_key:")
		assert.Contains(t, help, "timeout:")
		assert.Contains(t, help, "max_retries:")
		assert.Contains(t, help, "cache_ttl:")
		assert.Contains(t, help, "log_level:")
		assert.Contains(t, help, "format:")

		// Check for API key link
		assert.Contains(t, help, "https://www.themoviedb.org/settings/api")
	})

	t.Run("contains all valid log levels", func(t *testing.T) {
		help := internal.GetConfigHelp()

		validLogLevels := []string{"debug", "info", "warn", "error", "silent"}
		for _, level := range validLogLevels {
			assert.Contains(t, help, level)
		}
	})

	t.Run("contains all valid formats", func(t *testing.T) {
		help := internal.GetConfigHelp()

		validFormats := []string{"table", "json", "csv"}
		for _, format := range validFormats {
			assert.Contains(t, help, format)
		}
	})

	t.Run("includes duration format examples", func(t *testing.T) {
		help := internal.GetConfigHelp()

		assert.Contains(t, help, "\"30s\"")
		assert.Contains(t, help, "\"5m\"")
	})

	t.Run("help text is properly formatted", func(t *testing.T) {
		help := internal.GetConfigHelp()

		// Should have reasonable length
		assert.Greater(t, len(help), 500, "Help text should be comprehensive")

		// Should have multiple lines
		lines := strings.Split(help, "\n")
		assert.Greater(t, len(lines), 10, "Help should have multiple lines")

		// Should not have trailing whitespace on lines
		for i, line := range lines {
			if line != "" {
				expected := strings.TrimRight(line, " \t")
				assert.Equal(t, expected, line,
					"Line %d should not have trailing whitespace", i+1)
			}
		}
	})
	t.Run("help text consistency", func(t *testing.T) {
		help := internal.GetConfigHelp()

		// Should mention all config file paths consistently
		paths := []string{
			"./config.yaml", "./tmdb.yaml",
			"~/.tmdb/config.yaml", "~/.config/tmdb/config.yaml",
		}
		for _, path := range paths {
			assert.Contains(t, help, path, "Help should mention config path: %s", path)
		}

		// Should use consistent naming
		assert.Contains(t, help, "TMDB API key")
		assert.Contains(t, help, "Environment variables") // Capital E
		assert.Contains(t, help, "Command line flags")
	})
}

func TestGetAPIKeyHelp(t *testing.T) {
	t.Run("returns comprehensive API key help", func(t *testing.T) {
		help := internal.GetAPIKeyHelp()
		assert.NotEmpty(t, help)

		// Check for main steps
		assert.Contains(t, help, "To get a TMDB API key:")
		assert.Contains(t, help, "1. Go to https://www.themoviedb.org/")
		assert.Contains(t, help, "2. Create a free account")
		assert.Contains(t, help, "3. Go to https://www.themoviedb.org/settings/api")
		assert.Contains(t, help, "4. Request an API key")
		assert.Contains(t, help, "5. Fill out the application form")
		assert.Contains(t, help, "6. Once approved, copy your API key")

		// Check for usage instructions
		assert.Contains(t, help, "Then set it via:")
		assert.Contains(t, help, "Environment variable:")
		assert.Contains(t, help, "export TMDB_API_KEY=")
		assert.Contains(t, help, "Config file:")
		assert.Contains(t, help, "api_key:")
		assert.Contains(t, help, "~/.tmdb/config.yaml")
	})

	t.Run("contains all necessary URLs", func(t *testing.T) {
		help := internal.GetAPIKeyHelp()

		// Main TMDB website
		assert.Contains(t, help, "https://www.themoviedb.org/")

		// API settings page
		assert.Contains(t, help, "https://www.themoviedb.org/settings/api")
	})

	t.Run("mentions developer option", func(t *testing.T) {
		help := internal.GetAPIKeyHelp()

		assert.Contains(t, help, "Developer")
	})

	t.Run("provides multiple setup methods", func(t *testing.T) {
		help := internal.GetAPIKeyHelp()

		// Environment variable method
		assert.Contains(t, help, "TMDB_API_KEY")
		assert.Contains(t, help, "export")

		// Config file method
		assert.Contains(t, help, "config.yaml")
		assert.Contains(t, help, "api_key:")
	})

	t.Run("help text is properly formatted", func(t *testing.T) {
		help := internal.GetAPIKeyHelp()

		// Should have reasonable length
		assert.Greater(t, len(help), 200, "API key help should be comprehensive")

		// Should have multiple lines
		lines := strings.Split(help, "\n")
		assert.Greater(t, len(lines), 5, "Help should have multiple lines")

		// Should start with clear instruction
		assert.Contains(t, lines[0], "To get a TMDB API key:")
	})

	t.Run("includes step-by-step instructions", func(t *testing.T) {
		help := internal.GetAPIKeyHelp()

		// Should have numbered steps
		for i := 1; i <= 6; i++ {
			stepPrefix := formatInt(i) + "."
			assert.Contains(t, help, stepPrefix, "Should contain step %d", i)
		}
	})

	t.Run("mentions free account", func(t *testing.T) {
		help := internal.GetAPIKeyHelp()

		assert.Contains(t, help, "free account")
	})

	t.Run("provides complete workflow", func(t *testing.T) {
		help := internal.GetAPIKeyHelp()

		// Should cover the complete process
		workflow := []string{
			"Go to",
			"Create",
			"Request",
			"Fill out",
			"approved",
			"copy",
			"set it via",
		}

		for _, step := range workflow {
			assert.Contains(t, help, step, "Help should mention workflow step: %s", step)
		}
	})
}

func TestConfigHelpIntegration(t *testing.T) {
	t.Run("config help and API key help are consistent", func(t *testing.T) {
		configHelp := internal.GetConfigHelp()
		apiKeyHelp := internal.GetAPIKeyHelp()

		// Both should mention the same API URL
		apiURL := "https://www.themoviedb.org/settings/api"
		assert.Contains(t, configHelp, apiURL)
		assert.Contains(t, apiKeyHelp, apiURL)

		// Both should mention environment variables
		assert.Contains(t, configHelp, "TMDB_API_KEY")
		assert.Contains(t, apiKeyHelp, "TMDB_API_KEY")

		// Both should mention config file
		assert.Contains(t, configHelp, "config.yaml")
		assert.Contains(t, apiKeyHelp, "config.yaml")
	})

	t.Run("example config matches help documentation", func(t *testing.T) {
		helpers.WithTempDir(t, "tmdb-config-test-*", func(tempDir string) {
			configPath := filepath.Join(tempDir, "config.yaml")

			err := internal.CreateExampleConfig(configPath)
			require.NoError(t, err)

			content, err := os.ReadFile(configPath)
			require.NoError(t, err)

			configHelp := internal.GetConfigHelp()

			// Example in help should match created example
			exampleFields := []string{
				"api_key:", "timeout:", "max_retries:",
				"cache_ttl:", "log_level:", "format:",
			}
			for _, field := range exampleFields {
				assert.Contains(t, string(content), field,
					"Example config should contain field: %s", field)
				assert.Contains(t, configHelp, field,
					"Config help should mention field: %s", field)
			}
		})
	})

	t.Run("help mentions all supported config file locations", func(t *testing.T) {
		configHelp := internal.GetConfigHelp()

		// All paths that are checked in loadConfigFile
		configPaths := []string{
			"./config.yaml",
			"./tmdb.yaml",
			"~/.tmdb/config.yaml",
			"~/.config/tmdb/config.yaml",
		}

		for _, path := range configPaths {
			assert.Contains(t, configHelp, path,
				"Help should document config path: %s", path)
		}
	})

	t.Run("help environment variables match actual environment loading", func(t *testing.T) {
		configHelp := internal.GetConfigHelp()

		// All environment variables that are actually loaded
		envVars := []string{
			"TMDB_API_KEY",
			"TMDB_BASE_URL",
			"TMDB_TIMEOUT",
			"TMDB_MAX_RETRIES",
			"TMDB_CACHE_TTL",
			"TMDB_LOG_LEVEL",
			"TMDB_FORMAT",
		}

		for _, envVar := range envVars {
			assert.Contains(t, configHelp, envVar,
				"Help should document environment variable: %s", envVar)
		}
	})

	t.Run("example config creates loadable configuration", func(t *testing.T) {
		helpers.WithTempDir(t, "tmdb-config-test-*", func(tempDir string) {
			configPath := filepath.Join(tempDir, "config.yaml")

			// Create example config
			err := internal.CreateExampleConfig(configPath)
			require.NoError(t, err)

			// Verify the source file has content
			sourceContent, err := os.ReadFile(configPath)
			require.NoError(t, err)
			require.NotEmpty(t, sourceContent, "Source config file should not be empty")

			// Change to temp directory so config file is found
			originalDir, err := os.Getwd()
			require.NoError(t, err)
			err = os.Chdir(tempDir)
			require.NoError(t, err)
			defer func() {
				_ = os.Chdir(originalDir)
			}()

			// Since we're already in tempDir and the file is already named config.yaml,
			// we don't need to copy it. LoadConfig will find it in the current directory.

			// Debug: Check what's in the config file
			content, err := os.ReadFile("config.yaml")
			require.NoError(t, err)
			t.Logf("Config file content:\n%s", content)

			// Load config (should work without errors)
			config, err := internal.LoadConfig()
			require.NoError(t, err)

			// Should have placeholder API key
			assert.Equal(t, "your-api-key-here", config.APIKey)

			// Should have default values
			assert.Equal(t, "https://api.themoviedb.org/3", config.BaseURL)
			assert.Equal(t, "info", config.LogLevel)
			assert.Equal(t, "table", config.Format)
		})
	})
}

func TestConfigHelpErrorCases(t *testing.T) {
	t.Run("help functions never return empty strings", func(t *testing.T) {
		configHelp := internal.GetConfigHelp()
		apiKeyHelp := internal.GetAPIKeyHelp()

		assert.NotEmpty(t, configHelp, "Config help should never be empty")
		assert.NotEmpty(t, apiKeyHelp, "API key help should never be empty")

		assert.Greater(t, len(configHelp), 100, "Config help should be substantial")
		assert.Greater(t, len(apiKeyHelp), 100, "API key help should be substantial")
	})

	t.Run("help functions are deterministic", func(t *testing.T) {
		// Call multiple times and ensure same result
		configHelp1 := internal.GetConfigHelp()
		configHelp2 := internal.GetConfigHelp()
		apiKeyHelp1 := internal.GetAPIKeyHelp()
		apiKeyHelp2 := internal.GetAPIKeyHelp()

		assert.Equal(t, configHelp1, configHelp2, "Config help should be deterministic")
		assert.Equal(t, apiKeyHelp1, apiKeyHelp2, "API key help should be deterministic")
	})

	t.Run("help text contains no sensitive information", func(t *testing.T) {
		configHelp := internal.GetConfigHelp()
		apiKeyHelp := internal.GetAPIKeyHelp()

		// Should not contain actual API keys or passwords
		sensitivePatterns := []string{
			"password",
			"secret",
			"token",
			"bearer",
			"auth",
			// Removed "key=" as it appears in legitimate examples
		}

		allHelp := configHelp + apiKeyHelp
		for _, pattern := range sensitivePatterns {
			assert.NotContains(t, strings.ToLower(allHelp), pattern,
				"Help should not contain sensitive pattern: %s", pattern)
		}
	})

	t.Run("help text uses consistent terminology", func(t *testing.T) {
		configHelp := internal.GetConfigHelp()
		apiKeyHelp := internal.GetAPIKeyHelp()

		// Should use consistent terms
		allHelp := configHelp + apiKeyHelp

		// Should consistently use "TMDB API key" not variations
		assert.Contains(t, allHelp, "TMDB API key")
		assert.NotContains(t, allHelp, "tmdb api key")
		// Removed check for "api-key" as YAML uses "api_key:"

		// Should consistently use "Environment variable" not "env var"
		assert.Contains(t, allHelp, "Environment variable")
	})
}

// Helper functions to reduce cognitive complexity

func validateExampleConfigContent(t *testing.T, content string) {
	t.Helper()

	assert.Contains(t, content, "# TMDB CLI Configuration")
	assert.Contains(t, content, "# Get your API key from: https://www.themoviedb.org/settings/api")
	assert.Contains(t, content, "api_key: your-api-key-here") // No quotes
	assert.Contains(t, content, "base_url:")
	assert.Contains(t, content, "timeout:")
	assert.Contains(t, content, "max_retries:")
	assert.Contains(t, content, "cache_ttl:")
	assert.Contains(t, content, "log_level:")
	assert.Contains(t, content, "format:")
}

func validateExampleConfigYAML(t *testing.T, configPath string) {
	t.Helper()

	// Read file content
	content, err := os.ReadFile(configPath)
	require.NoError(t, err)

	// Extract YAML part (skip comments at the beginning)
	yamlContent := extractYAMLContent(string(content))

	// Parse YAML to verify it's valid
	var config internal.Config
	err = yaml.Unmarshal([]byte(yamlContent), &config)
	require.NoError(t, err)

	// Verify config has expected values
	assert.Equal(t, "your-api-key-here", config.APIKey)
	assert.Equal(t, "https://api.themoviedb.org/3", config.BaseURL)
	assert.Equal(t, "info", config.LogLevel)
	assert.Equal(t, "table", config.Format)
}

func validateExampleConfigDefaults(t *testing.T, configPath string) {
	t.Helper()

	content, err := os.ReadFile(configPath)
	require.NoError(t, err)

	// Extract and parse YAML
	yamlContent := extractYAMLContent(string(content))
	var exampleConfig internal.Config
	err = yaml.Unmarshal([]byte(yamlContent), &exampleConfig)
	require.NoError(t, err)

	// Compare with default config (except API key)
	defaultConfig := internal.DefaultConfig()
	assert.Equal(t, defaultConfig.BaseURL, exampleConfig.BaseURL)
	assert.Equal(t, defaultConfig.Timeout, exampleConfig.Timeout)
	assert.Equal(t, defaultConfig.MaxRetries, exampleConfig.MaxRetries)
	assert.Equal(t, defaultConfig.CacheTTL, exampleConfig.CacheTTL)
	assert.Equal(t, defaultConfig.LogLevel, exampleConfig.LogLevel)
	assert.Equal(t, defaultConfig.Format, exampleConfig.Format)
}

func validateExampleConfigStructure(t *testing.T, content string) {
	t.Helper()

	lines := strings.Split(content, "\n")

	// Verify structure
	assert.True(t, len(lines) > 5, "Config should have multiple lines")
	assert.Contains(t, lines[0], "# TMDB CLI Configuration")
	assert.Contains(t, lines[1], "# Get your API key from:")

	// Verify all required fields are present
	requiredFields := []string{
		"api_key:", "base_url:", "timeout:", "max_retries:",
		"cache_ttl:", "log_level:", "format:",
	}
	for _, field := range requiredFields {
		assert.Contains(t, content, field, "Config should contain field: %s", field)
	}
}

func extractYAMLContent(content string) string {
	lines := strings.Split(content, "\n")
	var yamlLines []string
	foundYAML := false
	for _, line := range lines {
		if strings.HasPrefix(line, "api_key:") {
			foundYAML = true
		}
		if foundYAML {
			yamlLines = append(yamlLines, line)
		}
	}
	return strings.Join(yamlLines, "\n")
}

// itoa converts integer to string without using strconv package.
func itoa(i int) string {
	if i == 0 {
		return "0"
	}

	if i < 0 {
		return "-" + itoa(-i)
	}

	var result []byte
	for i > 0 {
		result = append([]byte{byte('0' + i%10)}, result...)
		i /= 10
	}
	return string(result)
}

// formatInt is an alias for itoa to avoid naming conflicts.
func formatInt(i int) string {
	return itoa(i)
}
