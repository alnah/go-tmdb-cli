package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestShowProgress(t *testing.T) {
	t.Parallel()

	t.Run("shows progress for large requests", func(t *testing.T) {
		t.Parallel()

		// Capture stdout to verify progress message
		output := helpers.CaptureStdout(t, func() {
			message := "Fetching 50 popular movies"
			maxItems := 50

			internal.ShowProgress(message, maxItems)
		})

		assert.Contains(t, output, "Fetching 50 popular movies...")
	})

	t.Run("does not show progress for small requests", func(t *testing.T) {
		t.Parallel()

		// Capture stdout to verify no output
		output := helpers.CaptureStdout(t, func() {
			message := "Fetching 10 popular movies"
			maxItems := 10

			internal.ShowProgress(message, maxItems)
		})

		assert.Empty(t, output)
	})

	t.Run("threshold boundary test", func(t *testing.T) {
		t.Parallel()

		// Test exactly at threshold
		output := helpers.CaptureStdout(t, func() {
			message := "Fetching movies"
			maxItems := 40 // Equal to ProgressThreshold

			internal.ShowProgress(message, maxItems)
		})

		assert.Empty(t, output) // Should not show at threshold

		// Test just above threshold
		output2 := helpers.CaptureStdout(t, func() {
			message := "Fetching movies"
			maxItems := 41 // Above ProgressThreshold
			internal.ShowProgress(message, maxItems)
		})

		assert.Contains(t, output2, "Fetching movies...")
	})
}

func TestShowFetchingProgress(t *testing.T) {
	t.Parallel()

	t.Run("shows progress for fetching movies", func(t *testing.T) {
		t.Parallel()

		// Capture stdout to verify progress message
		output := helpers.CaptureStdout(t, func() {
			count := 60
			listType := "popular"

			internal.ShowFetchingProgress(count, listType)
		})

		expectedMessage := "Fetching 60 popular movies..."
		assert.Contains(t, output, expectedMessage)
	})

	t.Run("does not show progress for small fetch requests", func(t *testing.T) {
		t.Parallel()

		// Capture stdout to verify no output
		output := helpers.CaptureStdout(t, func() {
			count := 20
			listType := "top-rated"

			internal.ShowFetchingProgress(count, listType)
		})

		assert.Empty(t, output)
	})

	t.Run("shows progress for various list types", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name     string
			count    int
			listType string
			expected string
		}{
			{
				name:     "popular movies",
				count:    100,
				listType: "popular",
				expected: "Fetching 100 popular movies...",
			},
			{
				name:     "top-rated movies",
				count:    75,
				listType: "top-rated",
				expected: "Fetching 75 top-rated movies...",
			},
			{
				name:     "now-playing movies",
				count:    50,
				listType: "now-playing",
				expected: "Fetching 50 now-playing movies...",
			},
			{
				name:     "upcoming movies",
				count:    45,
				listType: "upcoming",
				expected: "Fetching 45 upcoming movies...",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				output := helpers.CaptureStdout(t, func() {
					internal.ShowFetchingProgress(tc.count, tc.listType)
				})

				assert.Contains(t, output, tc.expected)
			})
		}
	})
}

func TestShowSearchProgress(t *testing.T) {
	t.Parallel()

	t.Run("shows progress for large search requests", func(t *testing.T) {
		t.Parallel()

		// Capture stdout to verify progress message
		output := helpers.CaptureStdout(t, func() {
			query := "The Matrix"
			maxItems := 80

			internal.ShowSearchProgress(query, maxItems)
		})

		expectedMessage := "Searching for \"The Matrix\"..."
		assert.Contains(t, output, expectedMessage)
	})

	t.Run("does not show progress for small search requests", func(t *testing.T) {
		t.Parallel()

		// Capture stdout to verify no output
		output := helpers.CaptureStdout(t, func() {
			query := "Inception"
			maxItems := 30

			internal.ShowSearchProgress(query, maxItems)
		})

		assert.Empty(t, output)
	})

	t.Run("handles special characters in query", func(t *testing.T) {
		t.Parallel()

		output := helpers.CaptureStdout(t, func() {
			query := "Spider-Man: Into the Spider-Verse"
			maxItems := 60

			internal.ShowSearchProgress(query, maxItems)
		})

		expectedMessage := "Searching for \"Spider-Man: Into the Spider-Verse\"..."
		assert.Contains(t, output, expectedMessage)
	})

	t.Run("handles empty query", func(t *testing.T) {
		t.Parallel()

		output := helpers.CaptureStdout(t, func() {
			query := ""
			maxItems := 50

			internal.ShowSearchProgress(query, maxItems)
		})

		expectedMessage := "Searching for \"\"..."
		assert.Contains(t, output, expectedMessage)
	})
}

func TestShowTVSearchProgress(t *testing.T) {
	t.Parallel()

	t.Run("shows TV search progress for large requests", func(t *testing.T) {
		t.Parallel()

		// Capture stdout to verify progress message
		output := helpers.CaptureStdout(t, func() {
			query := "Breaking Bad"
			maxItems := 70

			internal.ShowTVSearchProgress(query, maxItems)
		})

		expectedMessage := "Searching TV shows for \"Breaking Bad\"..."
		assert.Contains(t, output, expectedMessage)
	})

	t.Run("does not show TV search progress for small requests", func(t *testing.T) {
		t.Parallel()

		// Capture stdout to verify no output
		output := helpers.CaptureStdout(t, func() {
			query := "The Office"
			maxItems := 25

			internal.ShowTVSearchProgress(query, maxItems)
		})

		assert.Empty(t, output)
	})

	t.Run("handles various TV show queries", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name     string
			query    string
			maxItems int
			expected string
		}{
			{
				name:     "popular TV show",
				query:    "Game of Thrones",
				maxItems: 100,
				expected: "Searching TV shows for \"Game of Thrones\"...",
			},
			{
				name:     "comedy series",
				query:    "Friends",
				maxItems: 60,
				expected: "Searching TV shows for \"Friends\"...",
			},
			{
				name:     "international show",
				query:    "Money Heist",
				maxItems: 45,
				expected: "Searching TV shows for \"Money Heist\"...",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				output := helpers.CaptureStdout(t, func() {
					internal.ShowTVSearchProgress(tc.query, tc.maxItems)
				})

				assert.Contains(t, output, tc.expected)
			})
		}
	})

	t.Run("TV search progress threshold boundary", func(t *testing.T) {
		t.Parallel()

		// Test exactly at threshold
		output := helpers.CaptureStdout(t, func() {
			query := "Test Show"
			maxItems := 40 // Equal to ProgressThreshold

			internal.ShowTVSearchProgress(query, maxItems)
		})

		assert.Empty(t, output) // Should not show at threshold

		// Test just above threshold
		output2 := helpers.CaptureStdout(t, func() {
			query := "Test Show"
			maxItems := 41 // Above ProgressThreshold
			internal.ShowTVSearchProgress(query, maxItems)
		})

		expectedMessage := "Searching TV shows for \"Test Show\"..."
		assert.Contains(t, output2, expectedMessage)
	})

	t.Run("handles Unicode characters in TV query", func(t *testing.T) {
		t.Parallel()

		output := helpers.CaptureStdout(t, func() {
			query := "La Casa de Papel"
			maxItems := 50

			internal.ShowTVSearchProgress(query, maxItems)
		})

		expectedMessage := "Searching TV shows for \"La Casa de Papel\"..."
		assert.Contains(t, output, expectedMessage)
	})
}
