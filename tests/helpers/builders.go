// Package helpers provides testing utilities for the TMDB CLI test suite.
package helpers

import (
	"time"

	"github.com/alnah/tmdb-cli/internal"
)

// ConfigBuilder provides a fluent interface for building test configs.
type ConfigBuilder struct {
	config internal.Config
}

// NewConfigBuilder creates a new config builder with default values.
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: internal.DefaultConfig(),
	}
}

// WithAPIKey sets the API key.
func (cb *ConfigBuilder) WithAPIKey(apiKey string) *ConfigBuilder {
	cb.config.APIKey = apiKey
	return cb
}

// WithBaseURL sets the base URL.
func (cb *ConfigBuilder) WithBaseURL(baseURL string) *ConfigBuilder {
	cb.config.BaseURL = baseURL
	return cb
}

// WithTimeout sets the timeout duration.
func (cb *ConfigBuilder) WithTimeout(timeout time.Duration) *ConfigBuilder {
	cb.config.Timeout = timeout
	return cb
}

// WithMaxRetries sets the maximum number of retries.
func (cb *ConfigBuilder) WithMaxRetries(maxRetries int) *ConfigBuilder {
	cb.config.MaxRetries = maxRetries
	return cb
}

// WithCacheTTL sets the cache TTL duration.
func (cb *ConfigBuilder) WithCacheTTL(cacheTTL time.Duration) *ConfigBuilder {
	cb.config.CacheTTL = cacheTTL
	return cb
}

// WithLogLevel sets the log level.
func (cb *ConfigBuilder) WithLogLevel(logLevel string) *ConfigBuilder {
	cb.config.LogLevel = logLevel
	return cb
}

// WithFormat sets the default format.
func (cb *ConfigBuilder) WithFormat(format string) *ConfigBuilder {
	cb.config.Format = format
	return cb
}

// Build returns the constructed config.
func (cb *ConfigBuilder) Build() internal.Config {
	return cb.config
}

// MovieBuilder provides a fluent interface for building test movies.
type MovieBuilder struct {
	movie internal.Movie
}

// NewMovieBuilder creates a new movie builder with default values.
func NewMovieBuilder() *MovieBuilder {
	return &MovieBuilder{
		movie: internal.Movie{
			ID:       1,
			Title:    "Test Movie",
			Year:     2023,
			Rating:   7.5,
			Votes:    1000,
			Genres:   "Action, Adventure",
			Overview: "A test movie for unit testing",
			Language: "en",
			Adult:    false,
		},
	}
}

// WithID sets the movie ID.
func (mb *MovieBuilder) WithID(id int) *MovieBuilder {
	mb.movie.ID = id
	return mb
}

// WithTitle sets the movie title.
func (mb *MovieBuilder) WithTitle(title string) *MovieBuilder {
	mb.movie.Title = title
	return mb
}

// WithOriginalTitle sets the original title.
func (mb *MovieBuilder) WithOriginalTitle(originalTitle string) *MovieBuilder {
	mb.movie.OriginalTitle = originalTitle
	return mb
}

// WithYear sets the release year.
func (mb *MovieBuilder) WithYear(year int) *MovieBuilder {
	mb.movie.Year = year
	return mb
}

// WithRating sets the rating.
func (mb *MovieBuilder) WithRating(rating float64) *MovieBuilder {
	mb.movie.Rating = rating
	return mb
}

// WithVotes sets the vote count.
func (mb *MovieBuilder) WithVotes(votes int) *MovieBuilder {
	mb.movie.Votes = votes
	return mb
}

// WithPopularity sets the popularity score.
func (mb *MovieBuilder) WithPopularity(popularity float64) *MovieBuilder {
	mb.movie.Popularity = popularity
	return mb
}

// WithGenres sets the genres.
func (mb *MovieBuilder) WithGenres(genres string) *MovieBuilder {
	mb.movie.Genres = genres
	return mb
}

// WithOverview sets the overview.
func (mb *MovieBuilder) WithOverview(overview string) *MovieBuilder {
	mb.movie.Overview = overview
	return mb
}

// WithLanguage sets the language.
func (mb *MovieBuilder) WithLanguage(language string) *MovieBuilder {
	mb.movie.Language = language
	return mb
}

// WithAdult sets the adult flag.
func (mb *MovieBuilder) WithAdult(adult bool) *MovieBuilder {
	mb.movie.Adult = adult
	return mb
}

// Build returns the constructed movie.
func (mb *MovieBuilder) Build() internal.Movie {
	return mb.movie
}

// TVShowBuilder provides a fluent interface for building test TV shows.
type TVShowBuilder struct {
	show internal.TVShow
}

// NewTVShowBuilder creates a new TV show builder with default values.
func NewTVShowBuilder() *TVShowBuilder {
	return &TVShowBuilder{
		show: internal.TVShow{
			ID:       1,
			Name:     "Test TV Show",
			Year:     2023,
			Rating:   8.0,
			Votes:    1500,
			Genres:   "Drama, Thriller",
			Overview: "A test TV show for unit testing",
			Language: "en",
			Adult:    false,
		},
	}
}

// WithID sets the TV show ID.
func (tb *TVShowBuilder) WithID(id int) *TVShowBuilder {
	tb.show.ID = id
	return tb
}

// WithName sets the TV show name.
func (tb *TVShowBuilder) WithName(name string) *TVShowBuilder {
	tb.show.Name = name
	return tb
}

// WithOriginalName sets the original name.
func (tb *TVShowBuilder) WithOriginalName(originalName string) *TVShowBuilder {
	tb.show.OriginalName = originalName
	return tb
}

// WithYear sets the first air year.
func (tb *TVShowBuilder) WithYear(year int) *TVShowBuilder {
	tb.show.Year = year
	return tb
}

// WithRating sets the rating.
func (tb *TVShowBuilder) WithRating(rating float64) *TVShowBuilder {
	tb.show.Rating = rating
	return tb
}

// WithVotes sets the vote count.
func (tb *TVShowBuilder) WithVotes(votes int) *TVShowBuilder {
	tb.show.Votes = votes
	return tb
}

// WithPopularity sets the popularity score.
func (tb *TVShowBuilder) WithPopularity(popularity float64) *TVShowBuilder {
	tb.show.Popularity = popularity
	return tb
}

// WithGenres sets the genres.
func (tb *TVShowBuilder) WithGenres(genres string) *TVShowBuilder {
	tb.show.Genres = genres
	return tb
}

// WithOverview sets the overview.
func (tb *TVShowBuilder) WithOverview(overview string) *TVShowBuilder {
	tb.show.Overview = overview
	return tb
}

// WithLanguage sets the language.
func (tb *TVShowBuilder) WithLanguage(language string) *TVShowBuilder {
	tb.show.Language = language
	return tb
}

// WithAdult sets the adult flag.
func (tb *TVShowBuilder) WithAdult(adult bool) *TVShowBuilder {
	tb.show.Adult = adult
	return tb
}

// Build returns the constructed TV show.
func (tb *TVShowBuilder) Build() internal.TVShow {
	return tb.show
}

// SearchOptionsBuilder provides a fluent interface for building search options.
type SearchOptionsBuilder struct {
	options internal.SearchOptions
}

// NewSearchOptionsBuilder creates a new search options builder with default values.
func NewSearchOptionsBuilder() *SearchOptionsBuilder {
	return &SearchOptionsBuilder{
		options: internal.SearchOptions{
			MaxItems:  20,
			MinRating: 0.0,
			MaxRating: 10.0,
			MinVotes:  0,
			MaxVotes:  0,
			Page:      1,
		},
	}
}

// WithMaxItems sets the maximum number of items.
func (sob *SearchOptionsBuilder) WithMaxItems(maxItems int) *SearchOptionsBuilder {
	sob.options.MaxItems = maxItems
	return sob
}

// WithMinRating sets the minimum rating.
func (sob *SearchOptionsBuilder) WithMinRating(minRating float64) *SearchOptionsBuilder {
	sob.options.MinRating = minRating
	return sob
}

// WithMaxRating sets the maximum rating.
func (sob *SearchOptionsBuilder) WithMaxRating(maxRating float64) *SearchOptionsBuilder {
	sob.options.MaxRating = maxRating
	return sob
}

// WithMinVotes sets the minimum vote count.
func (sob *SearchOptionsBuilder) WithMinVotes(minVotes int) *SearchOptionsBuilder {
	sob.options.MinVotes = minVotes
	return sob
}

// WithMaxVotes sets the maximum vote count.
func (sob *SearchOptionsBuilder) WithMaxVotes(maxVotes int) *SearchOptionsBuilder {
	sob.options.MaxVotes = maxVotes
	return sob
}

// WithYear sets the year (SearchOptions only has Year, not MinYear/MaxYear).
func (sob *SearchOptionsBuilder) WithYear(year int) *SearchOptionsBuilder {
	sob.options.Year = year
	return sob
}

// WithPage sets the page number.
func (sob *SearchOptionsBuilder) WithPage(page int) *SearchOptionsBuilder {
	sob.options.Page = page
	return sob
}

// WithQuery sets the search query.
func (sob *SearchOptionsBuilder) WithQuery(query string) *SearchOptionsBuilder {
	sob.options.Query = query
	return sob
}

// WithLanguage sets the language.
func (sob *SearchOptionsBuilder) WithLanguage(language string) *SearchOptionsBuilder {
	sob.options.Language = language
	return sob
}

// WithIncludeGenres sets the genres to include.
func (sob *SearchOptionsBuilder) WithIncludeGenres(genres []int) *SearchOptionsBuilder {
	sob.options.IncludeGenres = genres
	return sob
}

// WithExcludeGenres sets the genres to exclude.
func (sob *SearchOptionsBuilder) WithExcludeGenres(genres []int) *SearchOptionsBuilder {
	sob.options.ExcludeGenres = genres
	return sob
}

// WithSortBy sets the sort field.
func (sob *SearchOptionsBuilder) WithSortBy(sortBy string) *SearchOptionsBuilder {
	sob.options.SortBy = sortBy
	return sob
}

// WithSortOrder sets the sort order.
func (sob *SearchOptionsBuilder) WithSortOrder(sortOrder string) *SearchOptionsBuilder {
	sob.options.SortOrder = sortOrder
	return sob
}

// Build returns the constructed search options.
func (sob *SearchOptionsBuilder) Build() internal.SearchOptions {
	return sob.options
}

// Factory functions for common test scenarios

// CreateValidSearchOptions creates valid search options for testing.
func CreateValidSearchOptions() internal.SearchOptions {
	return NewSearchOptionsBuilder().
		WithMaxItems(20).
		WithMinRating(7.0).
		WithMaxRating(9.0).
		WithMinVotes(100).
		WithQuery("test movie").
		Build()
}

// CreateInvalidSearchOptions creates invalid search options for testing validation.
func CreateInvalidSearchOptions() internal.SearchOptions {
	return NewSearchOptionsBuilder().
		WithMaxItems(0).     // Invalid: should be >= 1
		WithMinRating(-1.0). // Invalid: should be >= 0
		WithMaxRating(11.0). // Invalid: should be <= 10
		Build()
}

// CreateHighRatedMovie creates a movie with high rating and sufficient votes.
func CreateHighRatedMovie() internal.Movie {
	return NewMovieBuilder().
		WithRating(8.5).
		WithVotes(1000).
		WithTitle("High Rated Movie").
		Build()
}

// CreatePopularMovie creates a movie with high popularity.
func CreatePopularMovie() internal.Movie {
	return NewMovieBuilder().
		WithPopularity(75.0).
		WithVotes(2000).
		WithTitle("Popular Movie").
		Build()
}

// CreateRecentMovie creates a movie from recent years.
func CreateRecentMovie() internal.Movie {
	currentYear := time.Now().Year()
	return NewMovieBuilder().
		WithYear(currentYear - 1).
		WithTitle("Recent Movie").
		Build()
}

// CreateEdgeCaseMovies creates movies with special characters and edge cases.
func CreateEdgeCaseMovies() []internal.Movie {
	return []internal.Movie{
		NewMovieBuilder().
			WithID(1).
			WithTitle("Movie with \"quotes\" and, commas").
			WithOverview("Overview with\nnewlines and special chars: àáâãäå").
			Build(),
		NewMovieBuilder().
			WithID(2).
			WithTitle("Unicode Title 🎬").
			WithOverview("Description with émojis and ñáéíóú characters").
			Build(),
		NewMovieBuilder().
			WithID(3).
			WithTitle("Título Unicódé ñáéíóú").
			WithOverview("Resumo com caracteres especiais").
			Build(),
	}
}

// CreateValidMovie creates a valid movie with optional field overrides.
func CreateValidMovie(overrides map[string]any) internal.Movie {
	builder := NewMovieBuilder()

	if overrides == nil {
		return builder.Build()
	}

	// Apply overrides using helper functions to reduce complexity
	if id, ok := overrides["ID"].(int); ok {
		builder = builder.WithID(id)
	}
	if title, ok := overrides["Title"].(string); ok {
		builder = builder.WithTitle(title)
	}
	if originalTitle, ok := overrides["OriginalTitle"].(string); ok {
		builder = builder.WithOriginalTitle(originalTitle)
	}
	if year, ok := overrides["Year"].(int); ok {
		builder = builder.WithYear(year)
	}
	if rating, ok := overrides["Rating"].(float64); ok {
		builder = builder.WithRating(rating)
	}
	if votes, ok := overrides["Votes"].(int); ok {
		builder = builder.WithVotes(votes)
	}
	if popularity, ok := overrides["Popularity"].(float64); ok {
		builder = builder.WithPopularity(popularity)
	}
	if genres, ok := overrides["Genres"].(string); ok {
		builder = builder.WithGenres(genres)
	}
	if overview, ok := overrides["Overview"].(string); ok {
		builder = builder.WithOverview(overview)
	}
	if language, ok := overrides["Language"].(string); ok {
		builder = builder.WithLanguage(language)
	}
	if adult, ok := overrides["Adult"].(bool); ok {
		builder = builder.WithAdult(adult)
	}

	return builder.Build()
}

// CreateValidTVShow creates a valid TV show with optional field overrides.
func CreateValidTVShow(overrides map[string]any) internal.TVShow {
	builder := NewTVShowBuilder()

	if overrides == nil {
		return builder.Build()
	}

	// Apply overrides using helper functions to reduce complexity
	if id, ok := overrides["ID"].(int); ok {
		builder = builder.WithID(id)
	}
	if name, ok := overrides["Name"].(string); ok {
		builder = builder.WithName(name)
	}
	if originalName, ok := overrides["OriginalName"].(string); ok {
		builder = builder.WithOriginalName(originalName)
	}
	if year, ok := overrides["Year"].(int); ok {
		builder = builder.WithYear(year)
	}
	if rating, ok := overrides["Rating"].(float64); ok {
		builder = builder.WithRating(rating)
	}
	if votes, ok := overrides["Votes"].(int); ok {
		builder = builder.WithVotes(votes)
	}
	if popularity, ok := overrides["Popularity"].(float64); ok {
		builder = builder.WithPopularity(popularity)
	}
	if genres, ok := overrides["Genres"].(string); ok {
		builder = builder.WithGenres(genres)
	}
	if overview, ok := overrides["Overview"].(string); ok {
		builder = builder.WithOverview(overview)
	}
	if language, ok := overrides["Language"].(string); ok {
		builder = builder.WithLanguage(language)
	}
	if adult, ok := overrides["Adult"].(bool); ok {
		builder = builder.WithAdult(adult)
	}

	return builder.Build()
}
