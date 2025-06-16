package unit

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestTMDBResponse(t *testing.T) {
	t.Run("valid TMDB response structure", func(t *testing.T) {
		response := fixtures.SampleTMDBResponse

		assert.Equal(t, 1, response.Page)
		assert.Len(t, response.Results, 1)
		assert.Equal(t, 10, response.TotalPages)
		assert.Equal(t, 200, response.TotalResults)

		// Verify the movie in results
		movie := response.Results[0]
		assert.Equal(t, 123, movie.ID)
		assert.Equal(t, "The Matrix", movie.Title)
	})

	t.Run("empty TMDB response", func(t *testing.T) {
		response := internal.TMDBResponse{
			Page:         1,
			Results:      []internal.TMDBMovie{},
			TotalPages:   0,
			TotalResults: 0,
		}

		assert.Equal(t, 1, response.Page)
		assert.Empty(t, response.Results)
		assert.Equal(t, 0, response.TotalPages)
		assert.Equal(t, 0, response.TotalResults)
	})

	t.Run("large TMDB response", func(t *testing.T) {
		movies := make([]internal.TMDBMovie, 20)
		for i := range movies {
			movies[i] = helpers.MockTMDBMovie(i+1, fmt.Sprintf("Movie %d", i+1))
		}

		response := internal.TMDBResponse{
			Page:         5,
			Results:      movies,
			TotalPages:   100,
			TotalResults: 2000,
		}

		assert.Equal(t, 5, response.Page)
		assert.Len(t, response.Results, 20)
		assert.Equal(t, 100, response.TotalPages)
		assert.Equal(t, 2000, response.TotalResults)

		// Verify all movies have valid IDs
		for i, movie := range response.Results {
			assert.Equal(t, i+1, movie.ID)
			assert.Contains(t, movie.Title, fmt.Sprintf("Movie %d", i+1))
		}
	})

	t.Run("response pagination boundaries", func(t *testing.T) {
		tests := []struct {
			name         string
			page         int
			totalPages   int
			totalResults int
		}{
			{"first page", 1, 10, 200},
			{"middle page", 5, 10, 200},
			{"last page", 10, 10, 200},
			{"single page", 1, 1, 20},
			{"many pages", 50, 100, 2000},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				response := internal.TMDBResponse{
					Page:         test.page,
					TotalPages:   test.totalPages,
					TotalResults: test.totalResults,
				}

				assert.Equal(t, test.page, response.Page)
				assert.Equal(t, test.totalPages, response.TotalPages)
				assert.Equal(t, test.totalResults, response.TotalResults)
				assert.Greater(t, response.Page, 0)
				assert.LessOrEqual(t, response.Page, response.TotalPages)
			})
		}
	})
}

func TestTMDBMovie(t *testing.T) {
	t.Run("valid TMDB movie structure", func(t *testing.T) {
		movie := fixtures.SampleTMDBMovie

		assert.Equal(t, 123, movie.ID)
		assert.Equal(t, "The Matrix", movie.Title)
		assert.Equal(t, "The Matrix", movie.OriginalTitle)
		assert.Equal(
			t,
			"A computer hacker learns about the true nature of his reality.",
			movie.Overview,
		)
		assert.Equal(t, "1999-03-31", movie.ReleaseDate)
		assert.Equal(t, 8.7, movie.VoteAverage)
		assert.Equal(t, 1800000, movie.VoteCount)
		assert.Equal(t, []int{28, 878}, movie.GenreIDs) // Action, Sci-Fi
		assert.Equal(t, 85.5, movie.Popularity)
		assert.False(t, movie.Adult)
		assert.False(t, movie.Video)
	})

	t.Run("TMDB movie with different original title", func(t *testing.T) {
		movie := internal.TMDBMovie{
			Title:         "Spirited Away",
			OriginalTitle: "千と千尋の神隠し",
		}

		assert.Equal(t, "Spirited Away", movie.Title)
		assert.Equal(t, "千と千尋の神隠し", movie.OriginalTitle)
		assert.NotEqual(t, movie.Title, movie.OriginalTitle)
		assert.Contains(t, movie.OriginalTitle, "千と千尋の神隠し")
	})

	t.Run("TMDB movie edge cases", func(t *testing.T) {
		movie := internal.TMDBMovie{
			ID:            1,
			Title:         "",      // Empty title
			OriginalTitle: "",      // Empty original title
			Overview:      "",      // Empty overview
			ReleaseDate:   "",      // Empty release date
			VoteAverage:   0.0,     // Zero rating
			VoteCount:     0,       // Zero votes
			GenreIDs:      []int{}, // Empty genres
			Popularity:    0.0,     // Zero popularity
			Adult:         false,
			Video:         true, // Is a video
		}

		assert.Equal(t, 1, movie.ID)
		assert.Empty(t, movie.Title)
		assert.Empty(t, movie.OriginalTitle)
		assert.Empty(t, movie.Overview)
		assert.Empty(t, movie.ReleaseDate)
		assert.Equal(t, 0.0, movie.VoteAverage)
		assert.Equal(t, 0, movie.VoteCount)
		assert.Empty(t, movie.GenreIDs)
		assert.Equal(t, 0.0, movie.Popularity)
		assert.False(t, movie.Adult)
		assert.True(t, movie.Video)
	})

	t.Run("TMDB movie with extreme values", func(t *testing.T) {
		movie := internal.TMDBMovie{
			ID:          999999999,
			VoteAverage: 10.0,      // Maximum rating
			VoteCount:   999999999, // Very high vote count
			Popularity:  999.99,    // Very high popularity
			Adult:       true,      // Adult content
		}

		assert.Equal(t, 999999999, movie.ID)
		assert.Equal(t, 10.0, movie.VoteAverage)
		assert.Equal(t, 999999999, movie.VoteCount)
		assert.Equal(t, 999.99, movie.Popularity)
		assert.True(t, movie.Adult)
	})

	t.Run("TMDB movie with many genres", func(t *testing.T) {
		movie := internal.TMDBMovie{
			GenreIDs: []int{28, 12, 16, 35, 80, 99, 18, 10751, 14, 36}, // Many genres
		}

		assert.Len(t, movie.GenreIDs, 10)
		assert.Contains(t, movie.GenreIDs, 28) // Action
		assert.Contains(t, movie.GenreIDs, 35) // Comedy
		assert.Contains(t, movie.GenreIDs, 18) // Drama
	})

	t.Run("TMDB movie with unicode content", func(t *testing.T) {
		movie := internal.TMDBMovie{
			Title:         "Movie with émojis 🎬 and åccénts",
			OriginalTitle: "Фильм с русскими символами",
			Overview:      "Description with 中文 characters and العربية text",
		}

		assert.Contains(t, movie.Title, "🎬")
		assert.Contains(t, movie.Title, "åccénts")
		assert.Contains(t, movie.OriginalTitle, "Фильм")
		assert.Contains(t, movie.Overview, "中文")
		assert.Contains(t, movie.Overview, "العربية")
	})
}

func TestTMDBTVResponse(t *testing.T) {
	t.Run("valid TMDB TV response structure", func(t *testing.T) {
		response := fixtures.SampleTMDBTVResponse

		assert.Equal(t, 1, response.Page)
		assert.Len(t, response.Results, 1)
		assert.Equal(t, 5, response.TotalPages)
		assert.Equal(t, 100, response.TotalResults)

		// Verify the TV show in results
		show := response.Results[0]
		assert.Equal(t, 789, show.ID)
		assert.Equal(t, "Breaking Bad", show.Name)
	})

	t.Run("empty TMDB TV response", func(t *testing.T) {
		response := internal.TMDBTVResponse{
			Page:         1,
			Results:      []internal.TMDBTVShow{},
			TotalPages:   0,
			TotalResults: 0,
		}

		assert.Equal(t, 1, response.Page)
		assert.Empty(t, response.Results)
		assert.Equal(t, 0, response.TotalPages)
		assert.Equal(t, 0, response.TotalResults)
	})

	t.Run("large TMDB TV response", func(t *testing.T) {
		shows := make([]internal.TMDBTVShow, 15)
		for i := range shows {
			shows[i] = helpers.MockTMDBTVShow(i+1, fmt.Sprintf("TV Show %d", i+1))
		}

		response := internal.TMDBTVResponse{
			Page:         3,
			Results:      shows,
			TotalPages:   25,
			TotalResults: 500,
		}

		assert.Equal(t, 3, response.Page)
		assert.Len(t, response.Results, 15)
		assert.Equal(t, 25, response.TotalPages)
		assert.Equal(t, 500, response.TotalResults)

		// Verify all shows have valid IDs
		for i, show := range response.Results {
			assert.Equal(t, i+1, show.ID)
			assert.Contains(t, show.Name, fmt.Sprintf("TV Show %d", i+1))
		}
	})
}

func TestTMDBTVShow(t *testing.T) {
	t.Run("valid TMDB TV show structure", func(t *testing.T) {
		show := fixtures.SampleTMDBTVShow

		assert.Equal(t, 789, show.ID)
		assert.Equal(t, "Breaking Bad", show.Name)
		assert.Equal(t, "Breaking Bad", show.OriginalName)
		assert.Equal(
			t,
			"A high school chemistry teacher turned methamphetamine manufacturer.",
			show.Overview,
		)
		assert.Equal(t, "2008-01-20", show.FirstAirDate)
		assert.Equal(t, 9.5, show.VoteAverage)
		assert.Equal(t, 1600000, show.VoteCount)
		assert.Equal(t, []int{80, 18}, show.GenreIDs) // Crime, Drama
		assert.Equal(t, 95.2, show.Popularity)
		assert.False(t, show.Adult)
		assert.Equal(t, "en", show.OriginalLanguage)
	})

	t.Run("TMDB TV show with different original name", func(t *testing.T) {
		show := internal.TMDBTVShow{
			Name:             "Squid Game",
			OriginalName:     "오징어 게임",
			OriginalLanguage: "ko",
		}

		assert.Equal(t, "Squid Game", show.Name)
		assert.Equal(t, "오징어 게임", show.OriginalName)
		assert.NotEqual(t, show.Name, show.OriginalName)
		assert.Equal(t, "ko", show.OriginalLanguage)
		assert.Contains(t, show.OriginalName, "오징어")
	})

	t.Run("TMDB TV show edge cases", func(t *testing.T) {
		show := internal.TMDBTVShow{
			ID:               1,
			Name:             "",      // Empty name
			OriginalName:     "",      // Empty original name
			Overview:         "",      // Empty overview
			FirstAirDate:     "",      // Empty air date
			VoteAverage:      0.0,     // Zero rating
			VoteCount:        0,       // Zero votes
			GenreIDs:         []int{}, // Empty genres
			Popularity:       0.0,     // Zero popularity
			Adult:            false,
			OriginalLanguage: "", // Empty language
		}

		assert.Equal(t, 1, show.ID)
		assert.Empty(t, show.Name)
		assert.Empty(t, show.OriginalName)
		assert.Empty(t, show.Overview)
		assert.Empty(t, show.FirstAirDate)
		assert.Equal(t, 0.0, show.VoteAverage)
		assert.Equal(t, 0, show.VoteCount)
		assert.Empty(t, show.GenreIDs)
		assert.Equal(t, 0.0, show.Popularity)
		assert.False(t, show.Adult)
		assert.Empty(t, show.OriginalLanguage)
	})

	t.Run("TMDB TV show with extreme values", func(t *testing.T) {
		show := internal.TMDBTVShow{
			ID:               999999999,
			VoteAverage:      10.0,      // Maximum rating
			VoteCount:        999999999, // Very high vote count
			Popularity:       999.99,    // Very high popularity
			Adult:            true,      // Adult content
			OriginalLanguage: "xx",      // Unusual language code
		}

		assert.Equal(t, 999999999, show.ID)
		assert.Equal(t, 10.0, show.VoteAverage)
		assert.Equal(t, 999999999, show.VoteCount)
		assert.Equal(t, 999.99, show.Popularity)
		assert.True(t, show.Adult)
		assert.Equal(t, "xx", show.OriginalLanguage)
	})

	t.Run("TMDB TV show with many genres", func(t *testing.T) {
		show := internal.TMDBTVShow{
			GenreIDs: []int{28, 12, 16, 35, 80, 99, 18, 10751}, // Many genres
		}

		assert.Len(t, show.GenreIDs, 8)
		assert.Contains(t, show.GenreIDs, 28) // Action
		assert.Contains(t, show.GenreIDs, 35) // Comedy
		assert.Contains(t, show.GenreIDs, 18) // Drama
		assert.Contains(t, show.GenreIDs, 80) // Crime
	})

	t.Run("TMDB TV show with unicode content", func(t *testing.T) {
		show := internal.TMDBTVShow{
			Name:             "TV Show with émojis 📺 and åccénts",
			OriginalName:     "Сериал с русскими символами",
			Overview:         "Description with 中文 characters and العربية text",
			OriginalLanguage: "ru",
		}

		assert.Contains(t, show.Name, "📺")
		assert.Contains(t, show.Name, "åccénts")
		assert.Contains(t, show.OriginalName, "Сериал")
		assert.Contains(t, show.Overview, "中文")
		assert.Contains(t, show.Overview, "العربية")
		assert.Equal(t, "ru", show.OriginalLanguage)
	})

	t.Run("various language codes", func(t *testing.T) {
		languages := []string{
			"en", "es", "fr", "de", "it", "pt", "ru", "ja", "ko", "zh", "ar", "hi",
		}

		for _, lang := range languages {
			t.Run("language_"+lang, func(t *testing.T) {
				show := internal.TMDBTVShow{
					OriginalLanguage: lang,
				}

				assert.Equal(t, lang, show.OriginalLanguage)
			})
		}
	})
}

func TestTMDBGenresResponse(t *testing.T) {
	t.Run("valid genres response", func(t *testing.T) {
		genres := []internal.Genre{
			{ID: 28, Name: "Action"},
			{ID: 35, Name: "Comedy"},
			{ID: 18, Name: "Drama"},
		}

		response := internal.TMDBGenresResponse{
			Genres: genres,
		}

		assert.Len(t, response.Genres, 3)
		assert.Equal(t, 28, response.Genres[0].ID)
		assert.Equal(t, "Action", response.Genres[0].Name)
		assert.Equal(t, 35, response.Genres[1].ID)
		assert.Equal(t, "Comedy", response.Genres[1].Name)
		assert.Equal(t, 18, response.Genres[2].ID)
		assert.Equal(t, "Drama", response.Genres[2].Name)
	})

	t.Run("empty genres response", func(t *testing.T) {
		response := internal.TMDBGenresResponse{
			Genres: []internal.Genre{},
		}

		assert.Empty(t, response.Genres)
	})

	t.Run("genres with unicode names", func(t *testing.T) {
		genres := []internal.Genre{
			{ID: 1, Name: "Действие"}, // Russian
			{ID: 2, Name: "コメディ"},     // Japanese
			{ID: 3, Name: "戏剧"},       // Chinese
			{ID: 4, Name: "أكشن"},     // Arabic
		}

		response := internal.TMDBGenresResponse{
			Genres: genres,
		}

		assert.Len(t, response.Genres, 4)
		assert.Contains(t, response.Genres[0].Name, "Действие")
		assert.Contains(t, response.Genres[1].Name, "コメディ")
		assert.Contains(t, response.Genres[2].Name, "戏剧")
		assert.Contains(t, response.Genres[3].Name, "أكشن")
	})

	t.Run("large genres response", func(t *testing.T) {
		genres := make([]internal.Genre, 50)
		for i := range genres {
			genres[i] = internal.Genre{
				ID:   i + 1,
				Name: fmt.Sprintf("Genre %d", i+1),
			}
		}

		response := internal.TMDBGenresResponse{
			Genres: genres,
		}

		assert.Len(t, response.Genres, 50)
		assert.Equal(t, 1, response.Genres[0].ID)
		assert.Equal(t, "Genre 1", response.Genres[0].Name)
		assert.Equal(t, 50, response.Genres[49].ID)
		assert.Equal(t, "Genre 50", response.Genres[49].Name)
	})
}
