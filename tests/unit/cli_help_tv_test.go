package unit

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestShowTVUsage(t *testing.T) {
	t.Parallel()

	t.Run("shows complete TV usage information", func(t *testing.T) {
		t.Parallel()

		// Capture stdout to verify TV usage output
		output := helpers.CaptureStdout(t, func() {
			internal.ShowTVUsage()
		})

		// Verify header
		assert.Contains(t, output, "TV Show Commands:")

		// Verify usage section
		assert.Contains(t, output, "USAGE:")
		assert.Contains(t, output, "tmdb-cli tv [subcommand] [arguments] [flags]")

		// Verify subcommands section
		assert.Contains(t, output, "SUBCOMMANDS:")
		assert.Contains(t, output, "popular, pop")
		assert.Contains(t, output, "Get popular TV shows")
		assert.Contains(t, output, "top-rated, top, rated")
		assert.Contains(t, output, "Get top-rated TV shows")
		assert.Contains(t, output, "on-the-air, air, airing")
		assert.Contains(t, output, "Get TV shows currently on the air")
		assert.Contains(t, output, "search [query]")
		assert.Contains(t, output, "Search for TV shows")

		// Verify examples section
		assert.Contains(t, output, "EXAMPLES:")
		assert.Contains(t, output, "tmdb-cli tv popular 10")
		assert.Contains(t, output, "tmdb-cli tv search \"Breaking Bad\"")
	})

	t.Run("TV usage output is properly formatted", func(t *testing.T) {
		t.Parallel()

		output := helpers.CaptureStdout(t, func() {
			internal.ShowTVUsage()
		})

		// Check for proper line breaks and formatting
		lines := strings.Split(output, "\n")
		assert.Greater(t, len(lines), 10, "Should have multiple lines of output")

		// Verify empty lines for formatting
		assert.Contains(t, output, "\n\n")

		// Verify proper indentation for subcommands
		assert.Contains(t, output, "  popular, pop")
		assert.Contains(t, output, "  top-rated, top, rated")
		assert.Contains(t, output, "  on-the-air, air, airing")
		assert.Contains(t, output, "  search [query]")
	})

	t.Run("TV usage includes all required subcommands", func(t *testing.T) {
		t.Parallel()

		output := helpers.CaptureStdout(t, func() {
			internal.ShowTVUsage()
		})

		// Check that all main TV subcommands are documented
		expectedSubcommands := []string{
			"popular",
			"top-rated",
			"on-the-air",
			"search",
		}

		for _, cmd := range expectedSubcommands {
			assert.Contains(t, output, cmd, "TV usage should document %s subcommand", cmd)
		}

		// Check that aliases are documented
		expectedAliases := []string{
			"pop",
			"top",
			"rated",
			"air",
			"airing",
		}

		for _, alias := range expectedAliases {
			assert.Contains(t, output, alias, "TV usage should document %s alias", alias)
		}
	})

	t.Run("TV usage includes practical examples", func(t *testing.T) {
		t.Parallel()

		output := helpers.CaptureStdout(t, func() {
			internal.ShowTVUsage()
		})

		// Verify examples are practical and accurate
		expectedExamples := []string{
			"tmdb-cli tv popular 10",
			"tmdb-cli tv search \"Breaking Bad\"",
		}

		for _, example := range expectedExamples {
			assert.Contains(t, output, example, "TV usage should include example: %s", example)
		}
	})
}

func TestShowSearchHelp(t *testing.T) {
	t.Parallel()

	t.Run("shows search help suggestions", func(t *testing.T) {
		t.Parallel()

		output := helpers.CaptureStdout(t, func() {
			internal.ShowSearchHelp()
		})

		// Verify the help suggestions
		assert.Contains(t, output, "Try:")
		assert.Contains(t, output, "Check spelling and try again")
		assert.Contains(t, output, "Use fewer, more common words")
		assert.Contains(t, output, "Try the original language title")
	})

	t.Run("search help formatting", func(t *testing.T) {
		t.Parallel()

		output := helpers.CaptureStdout(t, func() {
			internal.ShowSearchHelp()
		})

		// Check for proper bullet point formatting
		assert.Contains(t, output, "  - Check spelling")
		assert.Contains(t, output, "  - Use fewer")
		assert.Contains(t, output, "  - Try the original")
	})
}

func TestShowUsage(t *testing.T) {
	t.Parallel()

	t.Run("main usage includes TV command", func(t *testing.T) {
		t.Parallel()

		output := helpers.CaptureStdout(t, func() {
			version := "1.0.0"
			internal.ShowUsage(version)
		})

		// Verify TV command is documented in main usage
		assert.Contains(t, output, "tv [subcommand]")
		assert.Contains(t, output, "TV show commands")
	})

	t.Run("main usage shows version", func(t *testing.T) {
		t.Parallel()

		output := helpers.CaptureStdout(t, func() {
			version := "2.1.3"
			internal.ShowUsage(version)
		})

		// Verify version is shown
		assert.Contains(t, output, "TMDB CLI 2.1.3")
	})
}

func TestShowVersion(t *testing.T) {
	t.Parallel()

	t.Run("shows complete version information", func(t *testing.T) {
		t.Parallel()

		output := helpers.CaptureStdout(t, func() {
			version := "1.2.3"
			buildDate := "2024-01-15T10:30:00Z"
			gitCommit := "abc123def456"

			internal.ShowVersion(version, buildDate, gitCommit)
		})

		// Verify all version information is displayed
		assert.Contains(t, output, "TMDB CLI 1.2.3")
		assert.Contains(t, output, "Build Date: 2024-01-15T10:30:00Z")
		assert.Contains(t, output, "Git Commit: abc123def456")
	})

	t.Run("handles empty version information", func(t *testing.T) {
		t.Parallel()

		output := helpers.CaptureStdout(t, func() {
			version := ""
			buildDate := ""
			gitCommit := ""

			internal.ShowVersion(version, buildDate, gitCommit)
		})

		// Should still show structure even with empty values
		assert.Contains(t, output, "TMDB CLI")
		assert.Contains(t, output, "Build Date:")
		assert.Contains(t, output, "Git Commit:")
	})

	t.Run("handles development version", func(t *testing.T) {
		t.Parallel()

		output := helpers.CaptureStdout(t, func() {
			version := "dev"
			buildDate := "unknown"
			gitCommit := "dirty"

			internal.ShowVersion(version, buildDate, gitCommit)
		})

		assert.Contains(t, output, "TMDB CLI dev")
		assert.Contains(t, output, "Build Date: unknown")
		assert.Contains(t, output, "Git Commit: dirty")
	})
}
