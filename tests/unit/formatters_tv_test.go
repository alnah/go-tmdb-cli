package unit

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
)

func TestFormatTVSearchResult(t *testing.T) {
	t.Parallel()

	sampleTVShows := []internal.TVShow{
		{
			ID:           1,
			Name:         "Breaking Bad",
			OriginalName: "Breaking Bad",
			Year:         2008,
			Rating:       9.3,
			Votes:        500000,
			Genres:       "Drama, Crime",
			Popularity:   85.5,
		},
		{
			ID:           2,
			Name:         "The Office",
			OriginalName: "The Office",
			Year:         2005,
			Rating:       8.7,
			Votes:        300000,
			Genres:       "Comedy",
			Popularity:   75.2,
		},
	}

	t.Run("formats TV search result as JSON", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		result := internal.TVSearchResult{
			TVShows:      sampleTVShows,
			Page:         1,
			TotalPages:   2,
			TotalResults: 25,
		}

		options := internal.FormatOptions{
			Format: internal.FormatJSON,
		}

		err := internal.FormatTVSearchResult(&buf, result, options)

		assert.NoError(t, err)

		// Verify JSON is valid
		var jsonResult internal.TVSearchResult
		err = json.Unmarshal(buf.Bytes(), &jsonResult)
		require.NoError(t, err)

		assert.Equal(t, 1, jsonResult.Page)
		assert.Equal(t, 2, jsonResult.TotalPages)
		assert.Equal(t, 25, jsonResult.TotalResults)
		assert.Len(t, jsonResult.TVShows, 2)
		assert.Equal(t, "Breaking Bad", jsonResult.TVShows[0].Name)
	})

	t.Run("formats TV search result as CSV", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		result := internal.TVSearchResult{
			TVShows:      sampleTVShows,
			Page:         1,
			TotalPages:   1,
			TotalResults: 2,
		}

		options := internal.FormatOptions{
			Format: internal.FormatCSV,
		}

		err := internal.FormatTVSearchResult(&buf, result, options)

		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "ID,Name,Year,Rating,Votes,Genres")
		assert.Contains(t, output, "1,Breaking Bad,2008,9.3,500000,\"Drama,Crime\"")
		assert.Contains(t, output, "2,The Office,2005,8.7,300000,Comedy")
	})

	t.Run("formats TV search result as table with pagination", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		result := internal.TVSearchResult{
			TVShows:      sampleTVShows,
			Page:         2,
			TotalPages:   5,
			TotalResults: 100,
		}

		options := internal.FormatOptions{
			Format: internal.FormatTable,
		}

		err := internal.FormatTVSearchResult(&buf, result, options)

		assert.NoError(t, err)

		output := buf.String()
		// Should show pagination info
		assert.Contains(t, output, "Page 2 of 5 (100 total results)")
		// Should show table content
		assert.Contains(t, output, "Breaking Bad")
		assert.Contains(t, output, "The Office")
	})

	t.Run("formats empty TV search result", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		result := internal.TVSearchResult{
			TVShows:      []internal.TVShow{},
			Page:         1,
			TotalPages:   0,
			TotalResults: 0,
		}

		options := internal.FormatOptions{
			Format: internal.FormatTable,
		}

		err := internal.FormatTVSearchResult(&buf, result, options)

		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "No TV shows found.")
	})

	t.Run("formats TV search result with original titles", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		tvShows := []internal.TVShow{
			{
				ID:           3,
				Name:         "Money Heist",
				OriginalName: "La Casa de Papel",
				Year:         2017,
				Rating:       8.3,
				Votes:        200000,
				Genres:       "Drama, Crime",
			},
		}

		result := internal.TVSearchResult{
			TVShows:      tvShows,
			Page:         1,
			TotalPages:   1,
			TotalResults: 1,
		}

		options := internal.FormatOptions{
			Format:           internal.FormatTable,
			UseOriginalTitle: true,
		}

		err := internal.FormatTVSearchResult(&buf, result, options)

		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "La Casa de Papel")
	})

	t.Run("handles unsupported format", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		result := internal.TVSearchResult{
			TVShows: sampleTVShows,
		}

		options := internal.FormatOptions{
			Format: "xml",
		}

		err := internal.FormatTVSearchResult(&buf, result, options)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported format: xml")
	})
}

func TestFormatTVSearchResultJSON_Internal(t *testing.T) {
	t.Parallel()

	t.Run("formats complete TV search result as JSON", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		result := internal.TVSearchResult{
			TVShows: []internal.TVShow{
				{
					ID:           1,
					Name:         "Stranger Things",
					OriginalName: "Stranger Things",
					Year:         2016,
					Rating:       8.7,
					Votes:        400000,
					Genres:       "Drama, Fantasy, Horror",
					Popularity:   95.8,
				},
			},
			Page:         1,
			TotalPages:   3,
			TotalResults: 45,
		}

		options := internal.FormatOptions{
			Format: internal.FormatJSON,
		}

		err := internal.FormatTVSearchResult(&buf, result, options)

		assert.NoError(t, err)

		// Verify JSON structure
		var jsonResult internal.TVSearchResult
		err = json.Unmarshal(buf.Bytes(), &jsonResult)
		require.NoError(t, err)

		assert.Equal(t, 1, jsonResult.Page)
		assert.Equal(t, 3, jsonResult.TotalPages)
		assert.Equal(t, 45, jsonResult.TotalResults)
		require.Len(t, jsonResult.TVShows, 1)

		show := jsonResult.TVShows[0]
		assert.Equal(t, 1, show.ID)
		assert.Equal(t, "Stranger Things", show.Name)
		assert.Equal(t, 2016, show.Year)
		assert.Equal(t, 8.7, show.Rating)
		assert.Equal(t, 400000, show.Votes)
		assert.Contains(t, show.Genres, "Drama")
		assert.Contains(t, show.Genres, "Fantasy")
		assert.Contains(t, show.Genres, "Horror")
	})

	t.Run("formats empty TV search result as JSON", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		result := internal.TVSearchResult{
			TVShows:      []internal.TVShow{},
			Page:         1,
			TotalPages:   0,
			TotalResults: 0,
		}

		options := internal.FormatOptions{
			Format: internal.FormatJSON,
		}

		err := internal.FormatTVSearchResult(&buf, result, options)

		assert.NoError(t, err)

		var jsonResult internal.TVSearchResult
		err = json.Unmarshal(buf.Bytes(), &jsonResult)
		require.NoError(t, err)

		assert.Equal(t, 0, jsonResult.TotalResults)
		assert.Empty(t, jsonResult.TVShows)
	})
}

func TestFormatTVSearchResultTable_Internal(t *testing.T) {
	t.Parallel()

	t.Run("formats TV search result as table", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		result := internal.TVSearchResult{
			TVShows: []internal.TVShow{
				{
					ID:     1,
					Name:   "Game of Thrones",
					Year:   2011,
					Rating: 9.2,
					Votes:  1000000,
					Genres: "Drama, Fantasy",
				},
			},
			Page:         1,
			TotalPages:   1,
			TotalResults: 1,
		}

		options := internal.FormatOptions{
			Format: internal.FormatTable,
		}

		err := internal.FormatTVSearchResult(&buf, result, options)

		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "Game of Thrones")
		assert.Contains(t, output, "2011")
		assert.Contains(t, output, "9.2")
	})

	t.Run("shows pagination info for multi-page results", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		result := internal.TVSearchResult{
			TVShows: []internal.TVShow{
				{
					ID:     1,
					Name:   "Test Show",
					Year:   2020,
					Rating: 8.0,
					Votes:  10000,
					Genres: "Comedy",
				},
			},
			Page:         3,
			TotalPages:   7,
			TotalResults: 135,
		}

		options := internal.FormatOptions{
			Format: internal.FormatTable,
		}

		err := internal.FormatTVSearchResult(&buf, result, options)

		assert.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "Page 3 of 7 (135 total results)")
	})

	t.Run("does not show pagination for single page", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		result := internal.TVSearchResult{
			TVShows: []internal.TVShow{
				{
					ID:     1,
					Name:   "Single Show",
					Year:   2023,
					Rating: 7.5,
					Votes:  5000,
					Genres: "Drama",
				},
			},
			Page:         1,
			TotalPages:   1,
			TotalResults: 1,
		}

		options := internal.FormatOptions{
			Format: internal.FormatTable,
		}

		err := internal.FormatTVSearchResult(&buf, result, options)

		assert.NoError(t, err)

		output := buf.String()
		assert.NotContains(t, output, "Page 1 of 1")
		assert.Contains(t, output, "Single Show")
	})
}

func TestFormatTVError(t *testing.T) {
	t.Parallel()

	t.Run("formats basic error message", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		err := errors.New("API connection failed")

		writeErr := internal.FormatError(&buf, err)

		assert.NoError(t, writeErr)

		output := buf.String()
		assert.Equal(t, "Error: API connection failed\n", output)
	})

	t.Run("formats complex error message", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		err := errors.New("search failed: invalid API key (status: 401)")

		writeErr := internal.FormatError(&buf, err)

		assert.NoError(t, writeErr)

		output := buf.String()
		assert.Equal(t, "Error: search failed: invalid API key (status: 401)\n", output)
	})

	t.Run("formats nil error", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		var err error // nil error

		writeErr := internal.FormatError(&buf, err)

		assert.NoError(t, writeErr)

		output := buf.String()
		assert.Equal(t, "Error: <nil>\n", output)
	})

	t.Run("handles write errors", func(t *testing.T) {
		t.Parallel()

		// Use a writer that always fails
		failingWriter := &failingWriter{err: errors.New("write failed")}
		err := errors.New("test error")

		writeErr := internal.FormatError(failingWriter, err)

		assert.Error(t, writeErr)
		assert.Contains(t, writeErr.Error(), "write failed")
	})
}

// Helper type for testing write errors.
type failingWriter struct {
	err error
}

func (fw *failingWriter) Write(_ []byte) (n int, err error) {
	return 0, fw.err
}
