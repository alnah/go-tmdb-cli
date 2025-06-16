package helpers

import "github.com/alnah/tmdb-cli/internal"

// CreateValidMovie creates a valid movie for testing with optional overrides.
func CreateValidMovie(overrides map[string]any) internal.Movie {
	movie := internal.Movie{
		ID:            123,
		Title:         "Test Movie",
		OriginalTitle: "Test Movie",
		Year:          2023,
		Rating:        7.5,
		Votes:         1000,
		Popularity:    50.0,
		Genres:        "Action, Drama",
		Overview:      "Test movie overview",
		Language:      "en",
		Adult:         false,
		ReleaseDate:   "2023-01-01",
	}

	return applyMovieOverrides(movie, overrides)
}

// CreateValidTVShow creates a valid TV show for testing with optional overrides.
func CreateValidTVShow(overrides map[string]any) internal.TVShow {
	show := internal.TVShow{
		ID:           456,
		Name:         "Test Show",
		OriginalName: "Test Show",
		Year:         2023,
		Rating:       8.0,
		Votes:        1500,
		Popularity:   60.0,
		Genres:       "Drama, Thriller",
		Overview:     "Test TV show overview",
		Language:     "en",
		Adult:        false,
		FirstAirDate: "2023-01-15",
	}

	return applyTVShowOverrides(show, overrides)
}

// CreateValidSearchOptions creates valid search options for testing.
func CreateValidSearchOptions() internal.SearchOptions {
	return internal.SearchOptions{
		Query:         "test",
		Page:          1,
		Language:      "en",
		Year:          2023,
		MinRating:     7.0,
		MaxRating:     9.0,
		MinVotes:      100,
		MaxVotes:      100000,
		IncludeGenres: []int{28, 35}, // Action, Comedy
		ExcludeGenres: []int{27},     // Horror
		SortBy:        "popularity",
		SortOrder:     "desc",
		MaxItems:      20,
	}
}

// applyMovieOverrides applies override values to a movie with reduced complexity.
func applyMovieOverrides(movie internal.Movie, overrides map[string]any) internal.Movie {
	if overrides == nil {
		return movie
	}

	// Apply overrides using helper functions to reduce complexity
	movie.ID = applyIntOverride(overrides, "ID", movie.ID)
	movie.Title = applyStringOverride(overrides, "Title", movie.Title)
	movie.OriginalTitle = applyStringOverride(overrides, "OriginalTitle", movie.OriginalTitle)
	movie.Year = applyIntOverride(overrides, "Year", movie.Year)
	movie.Rating = applyFloat64Override(overrides, "Rating", movie.Rating)
	movie.Votes = applyIntOverride(overrides, "Votes", movie.Votes)
	movie.Popularity = applyFloat64Override(overrides, "Popularity", movie.Popularity)
	movie.Genres = applyStringOverride(overrides, "Genres", movie.Genres)
	movie.Overview = applyStringOverride(overrides, "Overview", movie.Overview)
	movie.Language = applyStringOverride(overrides, "Language", movie.Language)
	movie.Adult = applyBoolOverride(overrides, "Adult", movie.Adult)
	movie.ReleaseDate = applyStringOverride(overrides, "ReleaseDate", movie.ReleaseDate)

	return movie
}

// applyTVShowOverrides applies override values to a TV show with reduced complexity.
func applyTVShowOverrides(show internal.TVShow, overrides map[string]any) internal.TVShow {
	if overrides == nil {
		return show
	}

	// Apply overrides using helper functions to reduce complexity
	show.ID = applyIntOverride(overrides, "ID", show.ID)
	show.Name = applyStringOverride(overrides, "Name", show.Name)
	show.OriginalName = applyStringOverride(overrides, "OriginalName", show.OriginalName)
	show.Year = applyIntOverride(overrides, "Year", show.Year)
	show.Rating = applyFloat64Override(overrides, "Rating", show.Rating)
	show.Votes = applyIntOverride(overrides, "Votes", show.Votes)
	show.Popularity = applyFloat64Override(overrides, "Popularity", show.Popularity)
	show.Genres = applyStringOverride(overrides, "Genres", show.Genres)
	show.Overview = applyStringOverride(overrides, "Overview", show.Overview)
	show.Language = applyStringOverride(overrides, "Language", show.Language)
	show.Adult = applyBoolOverride(overrides, "Adult", show.Adult)
	show.FirstAirDate = applyStringOverride(overrides, "FirstAirDate", show.FirstAirDate)

	return show
}

// Helper functions for applying overrides by type.
func applyIntOverride(overrides map[string]any, key string, defaultValue int) int {
	if value, ok := overrides[key].(int); ok {
		return value
	}
	return defaultValue
}

func applyStringOverride(overrides map[string]any, key string, defaultValue string) string {
	if value, ok := overrides[key].(string); ok {
		return value
	}
	return defaultValue
}

func applyFloat64Override(
	overrides map[string]any,
	key string,
	defaultValue float64,
) float64 {
	if value, ok := overrides[key].(float64); ok {
		return value
	}
	return defaultValue
}

func applyBoolOverride(overrides map[string]any, key string, defaultValue bool) bool {
	if value, ok := overrides[key].(bool); ok {
		return value
	}
	return defaultValue
}
