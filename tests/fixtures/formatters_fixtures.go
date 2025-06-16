// Package fixtures provides test data for unit tests
package fixtures

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/alnah/tmdb-cli/internal"
)

// SampleMovies contains basic test movies for formatter testing.:w.
var SampleMovies = []internal.Movie{
	{
		ID:            2, // Test expects The Matrix to have ID=2
		Title:         "The Matrix",
		OriginalTitle: "The Matrix",
		Year:          1999,
		Rating:        8.7,
		Votes:         1847304,
		Popularity:    78.45,
		Genres:        "Action, Science Fiction",
		Overview:      "A computer hacker learns from mysterious rebels about the true nature of his reality.",
		Language:      "en",
		Adult:         false,
	},
	{
		ID:            1, // Spirited Away has ID=1
		Title:         "Spirited Away",
		OriginalTitle: "Sen to Chihiro no Kamikakushi",
		Year:          2001,
		Rating:        8.6,
		Votes:         692394,
		Popularity:    73.17,
		Genres:        "Animation, Family, Fantasy",
		Overview: "During her family's move to the suburbs, a sullen 10-year-old girl " +
			"wanders into a world ruled by gods.",
		Language: "ja",
		Adult:    false,
	},
	{
		ID:            3,
		Title:         "Pulp Fiction",
		OriginalTitle: "Pulp Fiction",
		Year:          1994,
		Rating:        8.9,
		Votes:         2070296,
		Popularity:    65.32,
		Genres:        "Crime, Drama",
		Overview:      "The lives of two mob hitmen, a boxer, and a pair of diner bandits intertwine.",
		Language:      "en",
		Adult:         false,
	},
}

// SampleTVShows contains basic test TV shows for formatter testing.
var SampleTVShows = []internal.TVShow{
	{
		ID:           1,
		Name:         "Breaking Bad",
		OriginalName: "Breaking Bad",
		Year:         2008,
		Rating:       9.5,
		Votes:        1587394,
		Popularity:   92.35, // Matches test expectation exactly
		Genres:       "Crime, Drama, Thriller",
		Overview:     "A high school chemistry teacher diagnosed with inoperable lung cancer.",
		Language:     "en",
		Adult:        false,
	},
	{
		ID:           2,
		Name:         "Squid Game",
		OriginalName: "오징어 게임",
		Year:         2021,
		Rating:       8.0,
		Votes:        500000,
		Popularity:   95.0,
		Genres:       "Drama, Thriller",
		Overview:     "Hundreds of cash-strapped players accept an invitation to compete in children's games.",
		Language:     "ko",
		Adult:        false,
	},
	{
		ID:           3,
		Name:         "The Sopranos",
		OriginalName: "The Sopranos",
		Year:         1999,
		Rating:       9.2,
		Votes:        384729,
		Popularity:   85.12,
		Genres:       "Crime, Drama",
		Overview: "This is a very long overview that should be truncated when displayed " +
			"in tables or other constrained formats to ensure proper formatting.",
		Language: "en",
		Adult:    false,
	},
}

// EmptyMovies represents an empty slice for testing empty results.
var EmptyMovies = []internal.Movie{}

// EmptyTVShows represents an empty slice for testing empty results.
var EmptyTVShows = []internal.TVShow{}

// SingleMovie contains one movie for testing single-item scenarios.
var SingleMovie = []internal.Movie{
	{
		ID:            100,
		Title:         "Test Movie",
		OriginalTitle: "Test Movie Original",
		Year:          2023,
		Rating:        7.5,
		Votes:         1000,
		Popularity:    25.5,
		Genres:        "Drama",
		Overview:      "A test movie for unit testing purposes.",
		Language:      "en",
		Adult:         false,
	},
}

// SingleTVShow contains one TV show for testing single-item scenarios.
var SingleTVShow = []internal.TVShow{
	{
		ID:           200,
		Name:         "Test Show",
		OriginalName: "Test Show Original",
		Year:         2023,
		Rating:       8.0,
		Votes:        2000,
		Popularity:   30.0,
		Genres:       "Comedy",
		Overview:     "A test TV show for unit testing purposes.",
		Language:     "en",
		Adult:        false,
	},
}

// FormatOptionsVariations contains different formatting scenarios.
var FormatOptionsVariations = []internal.FormatOptions{
	{
		Format:           "table",
		MaxWidth:         80,
		UseOriginalTitle: false,
		NoHeader:         false,
	},
	{
		Format:           "json",
		MaxWidth:         0,
		UseOriginalTitle: true,
		NoHeader:         false,
	},
	{
		Format:           "csv",
		MaxWidth:         120,
		UseOriginalTitle: false,
		NoHeader:         true,
	},
}

// EdgeCaseMovies contains special test cases for edge conditions including Unicode.
var EdgeCaseMovies = []internal.Movie{
	{
		ID:            999,
		Title:         "",
		OriginalTitle: "Original Empty Title",
		Year:          0,
		Rating:        0.0,
		Votes:         0,
		Popularity:    0.0,
		Genres:        "",
		Overview:      "",
		Language:      "",
		Adult:         false,
	},
	{
		ID:            1000,
		Title:         "Movie with Special Characters: !@#$%^&*()",
		OriginalTitle: "Фильм с специальными символами",
		Year:          2024,
		Rating:        10.0,
		Votes:         999999,
		Popularity:    100.0,
		Genres:        "Action, Comedy, Drama, Horror, Romance, Sci-Fi, Thriller",
		Overview:      "A movie with unusual characters and maximum values for testing edge cases.",
		Language:      "ru",
		Adult:         true,
	},
	{
		ID:            1001,
		Title:         "Very Long Movie Title That Exceeds Normal Length Expectations",
		OriginalTitle: "Très Long Titre de Film Qui Dépasse les Attentes de Longueur Normale",
		Year:          1900,
		Rating:        5.5,
		Votes:         1,
		Popularity:    0.1,
		Genres:        "Documentary",
		Overview:      "This movie has a very long title and serves as an edge case test.",
		Language:      "fr",
		Adult:         false,
	},
	{
		ID:            1002,
		Title:         "Unicode Title 🎬",
		OriginalTitle: "Título Unicódé ñáéíóú",
		Year:          2023,
		Rating:        7.8,
		Votes:         50000,
		Popularity:    45.67,
		Genres:        "Comedy",
		Overview:      "A movie with Unicode characters in title and original title for testing.",
		Language:      "es",
		Adult:         false,
	},
}

// EdgeCaseTVShows contains special test cases for TV shows edge conditions.
var EdgeCaseTVShows = []internal.TVShow{
	{
		ID:           2000,
		Name:         "",
		OriginalName: "Original Empty Name",
		Year:         0,
		Rating:       0.0,
		Votes:        0,
		Popularity:   0.0,
		Genres:       "",
		Overview:     "",
		Language:     "",
		Adult:        false,
	},
	{
		ID:           2001,
		Name:         "Show with Unicode: 你好世界",
		OriginalName: "Show with Unicode: 你好世界",
		Year:         2024,
		Rating:       10.0,
		Votes:        999999,
		Popularity:   100.0,
		Genres:       "Drama, Fantasy",
		Overview:     "A TV show testing unicode character handling in formatters.",
		Language:     "zh",
		Adult:        false,
	},
}

// ExpectedOutputPatterns contains expected output patterns for validation.
var ExpectedOutputPatterns = map[string][]string{
	"table_headers": {
		"ID",
		"Title",
		"Year",
		"Rating",
		"Votes",
		"Popularity",
		"Genres",
	},
	"json_keys": {
		"movies",
		"count",
		"id",
		"title",
		"year",
		"rating",
	},
	"csv_headers": {
		"ID,Title,Original Title,Year,Rating,Votes,Popularity,Genres,Overview,Language,Adult",
	},
}

// GenerateLargeMovieDataset creates a large dataset for performance testing.
func GenerateLargeMovieDataset(size int) []internal.Movie {
	movies := make([]internal.Movie, size)
	//nolint:gosec // Using weak random generator is acceptable for test data generation
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	genres := []string{"Action", "Comedy", "Drama", "Horror", "Romance", "Sci-Fi", "Thriller"}
	languages := []string{"en", "es", "fr", "de", "it", "ja", "ko", "zh"}

	for i := range size {
		movies[i] = internal.Movie{
			ID:            i + 1,
			Title:         fmt.Sprintf("Generated Movie %d", i+1),
			OriginalTitle: fmt.Sprintf("Original Generated Movie %d", i+1),
			Year:          1990 + rng.Intn(34), // 1990-2023
			Rating:        float64(rng.Intn(100)) / 10.0,
			Votes:         rng.Intn(1000000),
			Popularity:    float64(rng.Intn(10000)) / 100.0,
			Genres:        genres[rng.Intn(len(genres))],
			Overview: fmt.Sprintf(
				"Generated overview for movie %d with some descriptive text.",
				i+1,
			),
			Language: languages[rng.Intn(len(languages))],
			Adult:    rng.Float32() < 0.1, // 10% chance of adult content
		}
	}

	return movies
}

// GenerateLargeTVShowDataset creates a large TV show dataset for performance testing.
func GenerateLargeTVShowDataset(size int) []internal.TVShow {
	shows := make([]internal.TVShow, size)
	//nolint:gosec // Using weak random generator is acceptable for test data generation
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	genres := []string{"Action", "Comedy", "Drama", "Horror", "Romance", "Sci-Fi", "Thriller"}
	languages := []string{"en", "es", "fr", "de", "it", "ja", "ko", "zh"}

	for i := range size {
		shows[i] = internal.TVShow{
			ID:           i + 1,
			Name:         fmt.Sprintf("Generated Show %d", i+1),
			OriginalName: fmt.Sprintf("Original Generated Show %d", i+1),
			Year:         1990 + rng.Intn(34), // 1990-2023
			Rating:       float64(rng.Intn(100)) / 10.0,
			Votes:        rng.Intn(1000000),
			Popularity:   float64(rng.Intn(10000)) / 100.0,
			Genres:       genres[rng.Intn(len(genres))],
			Overview: fmt.Sprintf(
				"Generated overview for show %d with some descriptive text.",
				i+1,
			),
			Language: languages[rng.Intn(len(languages))],
			Adult:    rng.Float32() < 0.1, // 10% chance of adult content
		}
	}

	return shows
}
