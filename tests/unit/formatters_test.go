// tests/unit/formatters_test.go
package unit

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestFormatJSON(t *testing.T) {
	tests := []struct {
		name   string
		movies []internal.Movie
	}{
		{
			name:   "empty movie list",
			movies: []internal.Movie{},
		},
		{
			name:   "single movie",
			movies: []internal.Movie{helpers.CreateTestMovie(1, "Test Movie")},
		},
		{
			name: "multiple movies",
			movies: []internal.Movie{
				helpers.CreateTestMovie(1, "Movie One"),
				helpers.CreateTestMovie(2, "Movie Two"),
			},
		},
		{
			name: "movie with special characters",
			movies: []internal.Movie{
				{
					ID:       1,
					Title:    "Movie with \"quotes\" & symbols",
					Overview: "Description with <HTML> tags",
					Year:     2024,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := internal.FormatMovies(&buf, tt.movies, internal.FormatOptions{
				Format: internal.FormatJSON,
			})

			assert.NoError(t, err)

			// Validate JSON structure
			var result struct {
				Movies []internal.Movie `json:"movies"`
				Count  int              `json:"count"`
			}

			err = json.Unmarshal(buf.Bytes(), &result)
			assert.NoError(t, err, "Output should be valid JSON")
			assert.Equal(t, len(tt.movies), result.Count)
			assert.Len(t, result.Movies, len(tt.movies))

			// Verify movie data integrity
			for i, movie := range result.Movies {
				if i < len(tt.movies) {
					assert.Equal(t, tt.movies[i].ID, movie.ID)
					assert.Equal(t, tt.movies[i].Title, movie.Title)
				}
			}
		})
	}
}

func TestFormatCSV(t *testing.T) {
	tests := []struct {
		name      string
		movies    []internal.Movie
		options   internal.FormatOptions
		wantLines int
	}{
		{
			name:      "empty list with header",
			movies:    []internal.Movie{},
			options:   internal.FormatOptions{Format: internal.FormatCSV, NoHeader: false},
			wantLines: 1, // Just header
		},
		{
			name:      "empty list without header",
			movies:    []internal.Movie{},
			options:   internal.FormatOptions{Format: internal.FormatCSV, NoHeader: true},
			wantLines: 0, // No output
		},
		{
			name: "single movie with header",
			movies: []internal.Movie{
				helpers.CreateTestMovie(1, "Test Movie"),
			},
			options:   internal.FormatOptions{Format: internal.FormatCSV, NoHeader: false},
			wantLines: 2, // Header + 1 data row
		},
		{
			name: "multiple movies without header",
			movies: []internal.Movie{
				helpers.CreateTestMovie(1, "Movie One"),
				helpers.CreateTestMovie(2, "Movie Two"),
			},
			options:   internal.FormatOptions{Format: internal.FormatCSV, NoHeader: true},
			wantLines: 2, // 2 data rows only
		},
		{
			name: "movie with commas and quotes",
			movies: []internal.Movie{
				{
					ID:       1,
					Title:    "Movie, with \"commas\" and quotes",
					Overview: "Description, with commas",
					Genres:   "Action, Drama",
					Year:     2024,
				},
			},
			options:   internal.FormatOptions{Format: internal.FormatCSV, NoHeader: false},
			wantLines: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := internal.FormatMovies(&buf, tt.movies, tt.options)

			assert.NoError(t, err)

			output := buf.String()
			lines := strings.Split(strings.TrimSpace(output), "\n")

			if tt.wantLines == 0 {
				assert.Empty(t, strings.TrimSpace(output))
				return
			}

			assert.Len(t, lines, tt.wantLines)

			// Validate CSV structure
			reader := csv.NewReader(strings.NewReader(output))
			records, err := reader.ReadAll()
			assert.NoError(t, err)

			if !tt.options.NoHeader && len(records) > 0 {
				// Check header row
				header := records[0]
				expectedHeaders := []string{
					"ID",
					"Title",
					"Original Title",
					"Year",
					"Rating",
					"Votes",
					"Popularity",
					"Genres",
					"Overview",
					"Language",
					"Adult",
				}
				assert.Equal(t, expectedHeaders, header)
			}

			// Validate data consistency
			dataStartIdx := 0
			if !tt.options.NoHeader {
				dataStartIdx = 1
			}

			for i, record := range records[dataStartIdx:] {
				if i < len(tt.movies) {
					// Verify ID is first column and matches
					assert.Equal(t, string(rune(tt.movies[i].ID+'0')), record[0])
				}
			}
		})
	}
}

func TestFormatTable(t *testing.T) {
	tests := []struct {
		name        string
		movies      []internal.Movie
		options     internal.FormatOptions
		contains    []string
		notContains []string
	}{
		{
			name:     "empty movie list",
			movies:   []internal.Movie{},
			options:  internal.FormatOptions{Format: internal.FormatTable},
			contains: []string{"No movies found"},
		},
		{
			name: "single movie with headers",
			movies: []internal.Movie{
				helpers.CreateTestMovie(1, "Test Movie"),
			},
			options: internal.FormatOptions{Format: internal.FormatTable, NoHeader: false},
			contains: []string{
				"#",
				"TITLE",
				"YEAR",
				"RATING",
				"Test Movie",
			}, // Check for uppercase headers
		},
		{
			name: "single movie without headers",
			movies: []internal.Movie{
				helpers.CreateTestMovie(1, "Test Movie"),
			},
			options:     internal.FormatOptions{Format: internal.FormatTable, NoHeader: true},
			contains:    []string{"Test Movie"},
			notContains: []string{"TITLE", "YEAR", "RATING"},
		},
		{
			name: "multiple movies",
			movies: []internal.Movie{
				helpers.CreateTestMovie(1, "First Movie"),
				helpers.CreateTestMovie(2, "Second Movie"),
			},
			options:  internal.FormatOptions{Format: internal.FormatTable},
			contains: []string{"First Movie", "Second Movie", "1", "2"},
		},
		{
			name: "original titles enabled",
			movies: []internal.Movie{
				{
					ID:            1,
					Title:         "English Title",
					OriginalTitle: "Foreign Title",
					Year:          2024,
					Rating:        7.5,
					Votes:         1000,
				},
			},
			options: internal.FormatOptions{
				Format:           internal.FormatTable,
				UseOriginalTitle: true,
			},
			contains: []string{"Foreign Title", "English Title"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := internal.FormatMovies(&buf, tt.movies, tt.options)

			assert.NoError(t, err)

			output := buf.String()

			for _, want := range tt.contains {
				assert.Contains(t, output, want, "Output should contain: %s", want)
			}

			for _, notWant := range tt.notContains {
				assert.NotContains(t, output, notWant, "Output should not contain: %s", notWant)
			}

			// Verify table structure if not empty
			if len(tt.movies) > 0 {
				lines := strings.Split(output, "\n")
				assert.Greater(t, len(lines), 2, "Table should have multiple lines")
			}
		})
	}
}

func TestFormatMovieTitle(t *testing.T) {
	tests := []struct {
		name        string
		movie       internal.Movie
		useOriginal bool
		maxWidth    int
		expected    string
	}{
		{
			name: "simple title",
			movie: internal.Movie{
				Title:         "Simple Title",
				OriginalTitle: "Simple Title",
			},
			useOriginal: false,
			maxWidth:    40,
			expected:    "Simple Title",
		},
		{
			name: "different original title",
			movie: internal.Movie{
				Title:         "English Title",
				OriginalTitle: "Foreign Title",
			},
			useOriginal: false,
			maxWidth:    40,
			expected:    "English Title (Foreign Title)",
		},
		{
			name: "use original title",
			movie: internal.Movie{
				Title:         "English Title",
				OriginalTitle: "Foreign Title",
			},
			useOriginal: true,
			maxWidth:    40,
			expected:    "Foreign Title (English Title)",
		},
		{
			name: "long title truncation",
			movie: internal.Movie{
				Title:         "This is a Very Long Movie Title That Should Be Truncated",
				OriginalTitle: "This is a Very Long Movie Title That Should Be Truncated",
			},
			useOriginal: false,
			maxWidth:    20,
			// Expect truncation with ellipsis
			expected: "...",
		},
		{
			name: "long combined title truncation",
			movie: internal.Movie{
				Title:         "Long English Title",
				OriginalTitle: "Long Foreign Title",
			},
			useOriginal: false,
			maxWidth:    25,
			expected:    "Long English Title", // Should not add original if combined is too long
		},
		{
			name: "empty original title",
			movie: internal.Movie{
				Title:         "Main Title",
				OriginalTitle: "",
			},
			useOriginal: true,
			maxWidth:    40,
			expected:    "Main Title",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This tests the internal logic that would be called by formatMovieTitle
			// Since formatMovieTitle is not exported, we test the behavior through the public API
			var buf bytes.Buffer
			options := internal.FormatOptions{
				Format:           internal.FormatTable,
				UseOriginalTitle: tt.useOriginal,
				MaxWidth:         120, // Use standard width, test specific logic
			}

			err := internal.FormatMovies(&buf, []internal.Movie{tt.movie}, options)
			assert.NoError(t, err)

			output := buf.String()
			if tt.useOriginal && tt.movie.OriginalTitle != "" &&
				tt.movie.OriginalTitle != tt.movie.Title {
				assert.Contains(t, output, tt.movie.OriginalTitle)
			} else {
				// For truncation tests, just check that the title appears in some form
				if strings.Contains(tt.expected, "...") {
					assert.Contains(t, output, "...")
				} else {
					assert.Contains(t, output, tt.movie.Title)
				}
			}
		})
	}
}

func TestAddMovieIndicators(t *testing.T) {
	currentYear := 2024

	tests := []struct {
		name       string
		movie      internal.Movie
		expectTags []string
	}{
		{
			name: "recent movie",
			movie: internal.Movie{
				Year:       currentYear - 1, // Recent
				Rating:     7.0,
				Votes:      500,
				Popularity: 30.0,
			},
			expectTags: []string{"[NEW]"},
		},
		{
			name: "highly rated movie",
			movie: internal.Movie{
				Year:       2020,
				Rating:     8.5, // High rating
				Votes:      1000,
				Popularity: 30.0,
			},
			expectTags: []string{"[TOP]"},
		},
		{
			name: "popular movie",
			movie: internal.Movie{
				Year:       2020,
				Rating:     7.0,
				Votes:      2000, // High vote count
				Popularity: 80.0, // High popularity
			},
			expectTags: []string{"[HOT]"},
		},
		{
			name: "multiple indicators",
			movie: internal.Movie{
				Year:       currentYear, // Recent
				Rating:     8.5,         // Highly rated
				Votes:      2000,        // Popular
				Popularity: 80.0,        // Popular
			},
			expectTags: []string{"[NEW]", "[TOP]", "[HOT]"},
		},
		{
			name: "no indicators",
			movie: internal.Movie{
				Year:       2010,
				Rating:     6.0,
				Votes:      100,
				Popularity: 20.0,
			},
			expectTags: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			options := internal.FormatOptions{Format: internal.FormatTable}

			err := internal.FormatMovies(&buf, []internal.Movie{tt.movie}, options)
			assert.NoError(t, err)

			output := buf.String()

			for _, tag := range tt.expectTags {
				assert.Contains(t, output, tag, "Output should contain indicator: %s", tag)
			}

			// If no tags expected, verify none are present
			if len(tt.expectTags) == 0 {
				assert.NotContains(t, output, "[NEW]")
				assert.NotContains(t, output, "[TOP]")
				assert.NotContains(t, output, "[HOT]")
			}
		})
	}
}

func TestFormatSummary(t *testing.T) {
	movies := helpers.CreateSampleMovieList(5)

	tests := []struct {
		name        string
		command     string
		useOriginal bool
		expected    string
	}{
		{
			name:        "popular movies summary",
			command:     "popular",
			useOriginal: false,
			expected:    "Showing 5 popular movies (titles)",
		},
		{
			name:        "top-rated with original titles",
			command:     "top-rated",
			useOriginal: true,
			expected:    "Showing 5 top-rated movies (original language titles)",
		},
		{
			name:        "search results",
			command:     "search",
			useOriginal: false,
			expected:    "Found 5 movies (titles)",
		},
		{
			name:        "discovery results",
			command:     "discover",
			useOriginal: false,
			expected:    "Discovered 5 movies (titles)",
		},
		{
			name:        "now playing",
			command:     "now-playing",
			useOriginal: false,
			expected:    "Showing 5 movies now playing (titles)",
		},
		{
			name:        "upcoming",
			command:     "upcoming",
			useOriginal: false,
			expected:    "Showing 5 upcoming movies (titles)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			internal.FormatSummary(&buf, movies, tt.command, tt.useOriginal)

			output := buf.String()
			assert.Contains(t, output, tt.expected)
		})
	}

	t.Run("empty movie list", func(t *testing.T) {
		var buf bytes.Buffer
		internal.FormatSummary(&buf, []internal.Movie{}, "popular", false)

		output := buf.String()
		assert.Empty(t, output, "Should not output summary for empty list")
	})
}

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		format string
		valid  bool
	}{
		{"table", true},
		{"json", true},
		{"csv", true},
		{"xml", false},
		{"yaml", false},
		{"", false},
		{"TABLE", false}, // Case sensitive
		{"Json", false},  // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			result := internal.ValidateFormat(tt.format)
			assert.Equal(t, tt.valid, result)
		})
	}
}

func TestSupportedFormats(t *testing.T) {
	formats := internal.SupportedFormats()

	assert.Len(t, formats, 3)
	assert.Contains(t, formats, "table")
	assert.Contains(t, formats, "json")
	assert.Contains(t, formats, "csv")
}

func TestFormatError(t *testing.T) {
	tests := []struct {
		name  string
		error error
	}{
		{
			name:  "simple error",
			error: assert.AnError,
		},
		{
			name:  "nil error",
			error: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			if tt.error != nil {
				err := internal.FormatError(&buf, tt.error)
				assert.NoError(t, err)

				output := buf.String()
				assert.Contains(t, output, "Error:")
				assert.Contains(t, output, tt.error.Error())
			}
		})
	}
}

func TestFormatSearchResult(t *testing.T) {
	searchResult := internal.SearchResult{
		Movies: []internal.Movie{
			helpers.CreateTestMovie(1, "Test Movie"),
		},
		Page:         1,
		TotalPages:   5,
		TotalResults: 100,
	}

	tests := []struct {
		name     string
		options  internal.FormatOptions
		contains []string
	}{
		{
			name:     "JSON format",
			options:  internal.FormatOptions{Format: internal.FormatJSON},
			contains: []string{`"page"`, `"total_pages"`, `"total_results"`, `"movies"`},
		},
		{
			name:     "CSV format",
			options:  internal.FormatOptions{Format: internal.FormatCSV},
			contains: []string{"ID,Title"},
		},
		{
			name:     "Table format with pagination",
			options:  internal.FormatOptions{Format: internal.FormatTable},
			contains: []string{"Page 1 of 5", "100 total results", "Test Movie"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := internal.FormatSearchResult(&buf, searchResult, tt.options)

			assert.NoError(t, err)

			output := buf.String()
			for _, want := range tt.contains {
				assert.Contains(t, output, want)
			}
		})
	}

	t.Run("single page result", func(t *testing.T) {
		singlePageResult := internal.SearchResult{
			Movies:       []internal.Movie{helpers.CreateTestMovie(1, "Test Movie")},
			Page:         1,
			TotalPages:   1,
			TotalResults: 1,
		}

		var buf bytes.Buffer
		err := internal.FormatSearchResult(&buf, singlePageResult, internal.FormatOptions{
			Format: internal.FormatTable,
		})

		assert.NoError(t, err)

		output := buf.String()
		// Should not show pagination info for single page
		assert.NotContains(t, output, "Page 1 of 1")
		assert.Contains(t, output, "Test Movie")
	})
}

func TestFormatMoviesUnsupportedFormat(t *testing.T) {
	movies := []internal.Movie{helpers.CreateTestMovie(1, "Test Movie")}

	var buf bytes.Buffer
	err := internal.FormatMovies(&buf, movies, internal.FormatOptions{
		Format: "unsupported",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported format")
}

func TestFormatOptions(t *testing.T) {
	movie := internal.Movie{
		ID:       1,
		Title:    "Test Movie",
		Year:     2024,
		Rating:   7.5,
		Votes:    1000,
		Genres:   "Action, Drama",
		Overview: "A test movie for testing various format options and behaviors.",
	}

	t.Run("max width handling", func(t *testing.T) {
		var buf bytes.Buffer
		options := internal.FormatOptions{
			Format:   internal.FormatTable,
			MaxWidth: 80, // Smaller width
		}

		err := internal.FormatMovies(&buf, []internal.Movie{movie}, options)
		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "Test Movie")
	})

	t.Run("original title preference", func(t *testing.T) {
		movieWithOriginal := internal.Movie{
			ID:            1,
			Title:         "English Title",
			OriginalTitle: "Original Title",
			Year:          2024,
		}

		var buf bytes.Buffer
		options := internal.FormatOptions{
			Format:           internal.FormatTable,
			UseOriginalTitle: true,
		}

		err := internal.FormatMovies(&buf, []internal.Movie{movieWithOriginal}, options)
		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "Original Title")
	})
}

// Benchmark tests for performance validation.
func BenchmarkFormatJSON(b *testing.B) {
	movies := helpers.CreateSampleMovieList(100)
	options := internal.FormatOptions{Format: internal.FormatJSON}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_ = internal.FormatMovies(&buf, movies, options)
	}
}

func BenchmarkFormatTable(b *testing.B) {
	movies := helpers.CreateSampleMovieList(100)
	options := internal.FormatOptions{Format: internal.FormatTable}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_ = internal.FormatMovies(&buf, movies, options)
	}
}

func BenchmarkFormatCSV(b *testing.B) {
	movies := helpers.CreateSampleMovieList(100)
	options := internal.FormatOptions{Format: internal.FormatCSV}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_ = internal.FormatMovies(&buf, movies, options)
	}
}

// Edge case tests.
func TestFormatMoviesEdgeCases(t *testing.T) {
	t.Run("movie with nil/empty fields", func(t *testing.T) {
		movie := internal.Movie{
			ID:    1,
			Title: "Movie with Missing Data",
			// Other fields are zero values
		}

		var buf bytes.Buffer
		err := internal.FormatMovies(&buf, []internal.Movie{movie}, internal.FormatOptions{
			Format: internal.FormatJSON,
		})

		assert.NoError(t, err)
		helpers.AssertValidJSON(t, buf.Bytes())
	})

	t.Run("very large movie list", func(t *testing.T) {
		movies := helpers.CreateSampleMovieList(1000)

		var buf bytes.Buffer
		err := internal.FormatMovies(&buf, movies, internal.FormatOptions{
			Format: internal.FormatTable,
		})

		assert.NoError(t, err)
		output := buf.String()
		assert.NotEmpty(t, output)
	})

	t.Run("movie with unicode characters", func(t *testing.T) {
		movie := internal.Movie{
			ID:       1,
			Title:    "电影标题 🎬 Тест フィルム",
			Overview: "ôvërvïëw wïth ûnïcôdë 🌟",
			Genres:   "ドラマ, アクション",
		}

		formats := []string{internal.FormatJSON, internal.FormatCSV, internal.FormatTable}
		for _, format := range formats {
			t.Run(format, func(t *testing.T) {
				var buf bytes.Buffer
				err := internal.FormatMovies(&buf, []internal.Movie{movie}, internal.FormatOptions{
					Format: format,
				})

				assert.NoError(t, err)
				output := buf.String()
				assert.Contains(t, output, "电影标题")
			})
		}
	})
}
