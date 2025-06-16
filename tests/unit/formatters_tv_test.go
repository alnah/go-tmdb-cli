package unit

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
)

func TestTVShowFormatterBasicFields(t *testing.T) {
	tvShow := fixtures.SampleTVShows[0] // Breaking Bad
	formatter := internal.TVShowFormatter{TVShow: tvShow}

	t.Run("GetTitle", func(t *testing.T) {
		name := formatter.GetTitle(false)
		assert.Equal(t, "Breaking Bad", name)

		originalName := formatter.GetTitle(true)
		assert.Equal(t, "Breaking Bad", originalName) // Same for this show
	})

	t.Run("GetOriginalTitle", func(t *testing.T) {
		originalName := formatter.GetOriginalTitle()
		assert.Equal(t, "Breaking Bad", originalName)
	})

	t.Run("GetFields", func(t *testing.T) {
		fields := formatter.GetFields()
		expectedFields := []string{
			"1",                      // ID
			"Breaking Bad",           // Name
			"Breaking Bad",           // Original Name
			"2008",                   // Year
			"9.5",                    // Rating
			"1587394",                // Votes
			"92.35",                  // Popularity
			"Crime, Drama, Thriller", // Genres
			"A high school chemistry teacher diagnosed with inoperable lung cancer.", // Overview
			"en",    // Language
			"false", // Adult
		}

		require.Len(t, fields, len(expectedFields), "Should have correct number of fields")
		for i, expected := range expectedFields {
			assert.Equal(t, expected, fields[i], "Field %d should match", i)
		}
	})

	t.Run("GetEmptyMessage", func(t *testing.T) {
		message := formatter.GetEmptyMessage()
		assert.Equal(t, "No TV shows found.", message)
	})

	t.Run("GetHeaderName", func(t *testing.T) {
		headerName := formatter.GetHeaderName()
		assert.Equal(t, "Name", headerName)
	})
}

func TestTVShowFormatterWithDifferentNames(t *testing.T) {
	// Create a TV show with different original name
	tvShow := internal.TVShow{
		ID:           100,
		Name:         "Squid Game",
		OriginalName: "오징어 게임",
		Year:         2021,
	}
	formatter := internal.TVShowFormatter{TVShow: tvShow}

	t.Run("GetTitle shows original when requested", func(t *testing.T) {
		name := formatter.GetTitle(true)
		assert.Contains(t, name, "오징어 게임")
	})

	t.Run("GetTitle shows English when not requested", func(t *testing.T) {
		name := formatter.GetTitle(false)
		assert.Contains(t, name, "Squid Game")
	})
}

func TestTVShowFormatterEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		tvShow   internal.TVShow
		testFunc func(t *testing.T, formatter internal.TVShowFormatter)
	}{
		{
			name: "empty names",
			tvShow: internal.TVShow{
				ID:           1,
				Name:         "",
				OriginalName: "",
			},
			testFunc: func(t *testing.T, formatter internal.TVShowFormatter) {
				t.Helper()
				name := formatter.GetTitle(false)
				assert.Equal(t, "", name)

				originalName := formatter.GetOriginalTitle()
				assert.Equal(t, "", originalName)
			},
		},
		{
			name: "only original name",
			tvShow: internal.TVShow{
				ID:           2,
				Name:         "",
				OriginalName: "Original Only",
			},
			testFunc: func(t *testing.T, formatter internal.TVShowFormatter) {
				t.Helper()
				name := formatter.GetTitle(true)
				assert.Equal(t, "Original Only", name)

				name = formatter.GetTitle(false)
				assert.Equal(t, "", name) // Empty name
			},
		},
		{
			name: "only regular name",
			tvShow: internal.TVShow{
				ID:           3,
				Name:         "Regular Only",
				OriginalName: "",
			},
			testFunc: func(t *testing.T, formatter internal.TVShowFormatter) {
				t.Helper()
				name := formatter.GetTitle(false)
				assert.Equal(t, "Regular Only", name)

				name = formatter.GetTitle(true)
				assert.Equal(t, "Regular Only", name) // Falls back to regular name
			},
		},
		{
			name: "very long names",
			tvShow: internal.TVShow{
				ID:           4,
				Name:         "This Is An Extremely Long TV Show Name That Should Be Truncated When Displayed",
				OriginalName: "This Is An Even Longer Original TV Show Name That Should Also Be Truncated",
			},
			testFunc: func(t *testing.T, formatter internal.TVShowFormatter) {
				t.Helper()
				name := formatter.GetTitle(false)
				// Should be truncated
				assert.Contains(t, name, "...")

				originalName := formatter.GetTitle(true)
				// Should also be truncated
				assert.Contains(t, originalName, "...")
			},
		},
		{
			name: "unicode characters",
			tvShow: internal.TVShow{
				ID:           5,
				Name:         "Show 📺 with émojis",
				OriginalName: "Programa con émojis 📺",
			},
			testFunc: func(t *testing.T, formatter internal.TVShowFormatter) {
				t.Helper()
				name := formatter.GetTitle(false)
				assert.Contains(t, name, "📺")
				assert.Contains(t, name, "émojis")

				originalName := formatter.GetTitle(true)
				assert.Contains(t, originalName, "📺")
				assert.Contains(t, originalName, "Programa")
			},
		},
		{
			name: "zero and negative values",
			tvShow: internal.TVShow{
				ID:         6,
				Name:       "Edge Case Show",
				Year:       0,
				Rating:     0.0,
				Votes:      0,
				Popularity: 0.0,
			},
			testFunc: func(t *testing.T, formatter internal.TVShowFormatter) {
				t.Helper()
				fields := formatter.GetFields()
				assert.Equal(t, "0", fields[3])    // Year
				assert.Equal(t, "0.0", fields[4])  // Rating
				assert.Equal(t, "0", fields[5])    // Votes
				assert.Equal(t, "0.00", fields[6]) // Popularity
			},
		},
		{
			name: "adult content flag",
			tvShow: internal.TVShow{
				ID:    7,
				Name:  "Adult Show",
				Adult: true,
			},
			testFunc: func(t *testing.T, formatter internal.TVShowFormatter) {
				t.Helper()
				fields := formatter.GetFields()
				assert.Equal(t, "true", fields[10]) // Adult field
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := internal.TVShowFormatter{TVShow: tt.tvShow}
			tt.testFunc(t, formatter)
		})
	}
}

func TestTVShowFormatterFieldsConsistency(t *testing.T) {
	// Test that GetFields always returns the same number of fields
	testTVShows := []internal.TVShow{
		fixtures.SampleTVShows[0],
		fixtures.EdgeCaseTVShows[0],
		{ID: 999}, // Minimal TV show
	}

	expectedFieldCount := 11 // As defined in the CSV headers

	for i, tvShow := range testTVShows {
		t.Run(fmt.Sprintf("tvshow_%d", i), func(t *testing.T) {
			formatter := internal.TVShowFormatter{TVShow: tvShow}
			fields := formatter.GetFields()
			assert.Len(t, fields, expectedFieldCount, "All TV shows should have same field count")
		})
	}
}

func TestTVShowFormatterFloatFormatting(t *testing.T) {
	tvShow := internal.TVShow{
		ID:         1,
		Name:       "Test Show",
		Rating:     8.123456,  // Many decimal places
		Popularity: 75.987654, // Many decimal places
	}

	formatter := internal.TVShowFormatter{TVShow: tvShow}
	fields := formatter.GetFields()

	// Rating should be formatted to 1 decimal place
	assert.Equal(t, "8.1", fields[4])

	// Popularity should be formatted to 2 decimal places
	assert.Equal(t, "75.99", fields[6])
}

func TestTVShowFormatterSpecialCharacters(t *testing.T) {
	tvShow := internal.TVShow{
		ID:       1,
		Name:     "Show with \"quotes\" and 'apostrophes'",
		Genres:   "Comedy",
		Overview: "Overview with\nnewlines\tand\ttabs & special chars <html>",
		Language: "ko-KR",
	}

	formatter := internal.TVShowFormatter{TVShow: tvShow}
	fields := formatter.GetFields()

	// Fields should preserve special characters
	assert.Contains(t, fields[1], "quotes")   // Name
	assert.Contains(t, fields[7], "Comedy")   // Genres
	assert.Contains(t, fields[8], "newlines") // Overview
	assert.Contains(t, fields[8], "<html>")   // HTML in overview
	assert.Equal(t, "ko-KR", fields[9])       // Language
}

func TestGetTVShowHeaders(t *testing.T) {
	// This tests the unexported function through the CSV formatter
	tvShow := fixtures.SampleTVShows[0]

	var buf bytes.Buffer
	err := internal.FormatTVShows(&buf, []internal.TVShow{tvShow}, internal.FormatOptions{
		Format:   internal.FormatCSV,
		NoHeader: false,
	})
	require.NoError(t, err)

	output := buf.String()
	lines := strings.Split(output, "\n")
	require.NotEmpty(t, lines)

	// First line should contain headers
	headerLine := lines[0]
	expectedHeaders := []string{
		"ID", "Name", "Original Name", "Year", "Rating",
		"Votes", "Popularity", "Genres", "Overview", "Language", "Adult",
	}

	for _, header := range expectedHeaders {
		assert.Contains(t, headerLine, header, "Header line should contain: %s", header)
	}
}

func TestTVShowFormatterInterface(t *testing.T) {
	tvShow := fixtures.SampleTVShows[0]
	formatter := internal.TVShowFormatter{TVShow: tvShow}

	// Test that TVShowFormatter implements MediaFormatter interface
	var mediaFormatter internal.MediaFormatter = formatter

	// Test interface methods
	name := mediaFormatter.GetTitle(false)
	assert.NotEmpty(t, name)

	originalName := mediaFormatter.GetOriginalTitle()
	assert.NotEmpty(t, originalName)

	fields := mediaFormatter.GetFields()
	assert.NotEmpty(t, fields)

	emptyMessage := mediaFormatter.GetEmptyMessage()
	assert.Equal(t, "No TV shows found.", emptyMessage)

	headerName := mediaFormatter.GetHeaderName()
	assert.Equal(t, "Name", headerName)
}
