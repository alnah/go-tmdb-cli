package helpers

import "github.com/alnah/tmdb-cli/internal"

// MockTMDBMovie creates a mock TMDB movie response for testing.
func MockTMDBMovie(id int, title string) internal.TMDBMovie {
	return internal.TMDBMovie{
		ID:            id,
		Title:         title,
		OriginalTitle: title,
		Overview:      "Mock movie overview",
		ReleaseDate:   "2023-01-01",
		VoteAverage:   7.5,
		VoteCount:     1000,
		GenreIDs:      []int{28, 35}, // Action, Comedy
		Popularity:    50.0,
		Adult:         false,
		Video:         false,
	}
}

// MockTMDBTVShow creates a mock TMDB TV show response for testing.
func MockTMDBTVShow(id int, name string) internal.TMDBTVShow {
	return internal.TMDBTVShow{
		ID:               id,
		Name:             name,
		OriginalName:     name,
		Overview:         "Mock TV show overview",
		FirstAirDate:     "2023-01-01",
		VoteAverage:      8.0,
		VoteCount:        1500,
		GenreIDs:         []int{18, 53}, // Drama, Thriller
		Popularity:       60.0,
		Adult:            false,
		OriginalLanguage: "en",
	}
}
