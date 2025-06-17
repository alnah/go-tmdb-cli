// Package fixtures provides test data and fixtures for models layer tests.
package fixtures

import (
	"time"

	"github.com/alnah/tmdb-cli/internal"
)

// SampleMovie fixture.
var SampleMovie = internal.Movie{
	ID:            123,
	Title:         "The Matrix",
	OriginalTitle: "The Matrix",
	Year:          1999,
	Rating:        8.7,
	Votes:         1800000,
	Popularity:    85.5,
	Genres:        "Action, Sci-Fi",
	Overview:      "A computer hacker learns about the true nature of his reality.",
	Language:      "en",
	Adult:         false,
	ReleaseDate:   "1999-03-31",
}

// SampleMovieWithDifferentOriginal fixture.
var SampleMovieWithDifferentOriginal = internal.Movie{
	ID:            456,
	Title:         "Spirited Away",
	OriginalTitle: "千と千尋の神隠し",
	Year:          2001,
	Rating:        9.2,
	Votes:         750000,
	Popularity:    92.3,
	Genres:        "Animation, Family, Fantasy",
	Overview:      "A young girl enters a magical world.",
	Language:      "ja",
	Adult:         false,
	ReleaseDate:   "2001-07-20",
}

// InvalidMovie fixture.
var InvalidMovie = internal.Movie{
	ID:            0,  // Invalid ID
	Title:         "", // Empty title
	OriginalTitle: "Original",
	Year:          2023,
	Rating:        7.5,
	Votes:         100,
	Popularity:    50.0,
	Genres:        "Drama",
	Overview:      "Test movie",
	Language:      "en",
	Adult:         false,
}

// SampleTVShow fixture.
var SampleTVShow = internal.TVShow{
	ID:           789,
	Name:         "Breaking Bad",
	OriginalName: "Breaking Bad",
	Year:         2008,
	Rating:       9.5,
	Votes:        1600000,
	Popularity:   95.2,
	Genres:       "Crime, Drama, Thriller",
	Overview:     "A high school chemistry teacher turned methamphetamine manufacturer.",
	Language:     "en",
	Adult:        false,
	FirstAirDate: "2008-01-20",
}

// SampleTVShowWithDifferentOriginal fixture.
var SampleTVShowWithDifferentOriginal = internal.TVShow{
	ID:           101,
	Name:         "Squid Game",
	OriginalName: "오징어 게임",
	Year:         2021,
	Rating:       8.0,
	Votes:        500000,
	Popularity:   88.7,
	Genres:       "Drama, Thriller",
	Overview:     "Players compete in children's games for a cash prize.",
	Language:     "ko",
	Adult:        false,
	FirstAirDate: "2021-09-17",
}

// InvalidTVShow fixture.
var InvalidTVShow = internal.TVShow{
	ID:           -1, // Invalid ID
	Name:         "", // Empty name
	OriginalName: "Original Name",
	Year:         2023,
	Rating:       8.0,
	Votes:        1000,
	Popularity:   60.0,
	Genres:       "Drama",
	Overview:     "Test show",
	Language:     "en",
	Adult:        false,
}

// SampleTMDBMovie fixture.
var SampleTMDBMovie = internal.TMDBMovie{
	ID:            123,
	Title:         "The Matrix",
	OriginalTitle: "The Matrix",
	Overview:      "A computer hacker learns about the true nature of his reality.",
	ReleaseDate:   "1999-03-31",
	VoteAverage:   8.7,
	VoteCount:     1800000,
	GenreIDs:      []int{28, 878}, // Action, Sci-Fi
	Popularity:    85.5,
	Adult:         false,
	Video:         false,
}

// SampleTMDBTVShow fixture.
var SampleTMDBTVShow = internal.TMDBTVShow{
	ID:               789,
	Name:             "Breaking Bad",
	OriginalName:     "Breaking Bad",
	Overview:         "A high school chemistry teacher turned methamphetamine manufacturer.",
	FirstAirDate:     "2008-01-20",
	VoteAverage:      9.5,
	VoteCount:        1600000,
	GenreIDs:         []int{80, 18}, // Crime, Drama
	Popularity:       95.2,
	Adult:            false,
	OriginalLanguage: "en",
}

// SampleTMDBResponse fixture.
var SampleTMDBResponse = internal.TMDBResponse{
	Page:         1,
	Results:      []internal.TMDBMovie{SampleTMDBMovie},
	TotalPages:   10,
	TotalResults: 200,
}

// SampleTMDBTVResponse fixture.
var SampleTMDBTVResponse = internal.TMDBTVResponse{
	Page:         1,
	Results:      []internal.TMDBTVShow{SampleTMDBTVShow},
	TotalPages:   5,
	TotalResults: 100,
}

// SampleSearchResult fixture.
var SampleSearchResult = internal.SearchResult{
	Movies:       []internal.Movie{SampleMovie},
	Page:         1,
	TotalPages:   3,
	TotalResults: 50,
}

// SampleTVSearchResult fixture.
var SampleTVSearchResult = internal.TVSearchResult{
	TVShows:      []internal.TVShow{SampleTVShow},
	Page:         1,
	TotalPages:   2,
	TotalResults: 25,
}

// SampleConfig fixture.
var SampleConfig = internal.Config{
	APIKey:     "test-api-key-12345",
	BaseURL:    "https://api.themoviedb.org/3",
	Timeout:    30 * time.Second,
	MaxRetries: 3,
	CacheTTL:   5 * time.Minute,
	LogLevel:   "info",
	Format:     "table",
}

// InvalidConfig fixture.
var InvalidConfig = internal.Config{
	APIKey:     "", // Missing API key
	BaseURL:    "invalid-url",
	Timeout:    -1 * time.Second, // Invalid timeout
	MaxRetries: -1,               // Invalid retries
	CacheTTL:   -1 * time.Minute, // Invalid cache TTL
	LogLevel:   "invalid",        // Invalid log level
	Format:     "invalid",        // Invalid format
}

// ValidSearchOptions fixture.
var ValidSearchOptions = internal.SearchOptions{
	Query:         "matrix",
	Page:          1,
	Language:      "en",
	Year:          1999,
	MinRating:     7.0,
	MaxRating:     10.0,
	MinVotes:      1000,
	MaxVotes:      2000000,
	IncludeGenres: []int{28, 878}, // Action, Sci-Fi
	ExcludeGenres: []int{27},      // Horror
	SortBy:        "popularity",
	SortOrder:     "desc",
	MaxItems:      20,
}

// InvalidSearchOptions fixture.
var InvalidSearchOptions = internal.SearchOptions{
	Query:         "", // Empty query
	Page:          -1, // Invalid page
	Language:      "invalid-lang",
	Year:          1800,       // Invalid year
	MinRating:     -1.0,       // Invalid rating
	MaxRating:     15.0,       // Invalid rating
	MinVotes:      -100,       // Invalid votes
	MaxVotes:      -1,         // Invalid votes
	IncludeGenres: []int{999}, // Invalid genre
	ExcludeGenres: []int{-1},  // Invalid genre
	SortBy:        "invalid",
	SortOrder:     "invalid",
	MaxItems:      -1, // Invalid max items
}

// GenreTestCases fixture.
var GenreTestCases = []struct {
	Name     string
	Input    []string
	Expected []int
	HasError bool
}{
	{
		Name:     "valid single genre",
		Input:    []string{"action"},
		Expected: []int{28},
		HasError: false,
	},
	{
		Name:     "valid multiple genres",
		Input:    []string{"action", "comedy", "drama"},
		Expected: []int{28, 35, 18},
		HasError: false,
	},
	{
		Name:     "case insensitive genres",
		Input:    []string{"ACTION", "Comedy", "dRaMa"},
		Expected: []int{28, 35, 18},
		HasError: false,
	},
	{
		Name:     "sci-fi variations",
		Input:    []string{"sci-fi", "science-fiction"},
		Expected: []int{878, 878},
		HasError: false,
	},
	{
		Name:     "invalid genre",
		Input:    []string{"invalid-genre"},
		Expected: nil,
		HasError: true,
	},
	{
		Name:     "mixed valid and invalid",
		Input:    []string{"action", "invalid-genre"},
		Expected: nil,
		HasError: true,
	},
	{
		Name:     "genres with whitespace",
		Input:    []string{" action ", " comedy "},
		Expected: []int{28, 35},
		HasError: false,
	},
}

// YearTestCases fixture.
var YearTestCases = []struct {
	Name        string
	Input       string
	Expected    int
	Description string
}{
	{
		Name:        "valid full date",
		Input:       "1999-03-31",
		Expected:    1999,
		Description: "Full release date",
	},
	{
		Name:        "valid year only",
		Input:       "2023",
		Expected:    2023,
		Description: "Year only",
	},
	{
		Name:        "empty string",
		Input:       "",
		Expected:    0,
		Description: "Empty date string",
	},
	{
		Name:        "invalid format",
		Input:       "invalid-date",
		Expected:    0,
		Description: "Invalid date format",
	},
	{
		Name:        "partial date",
		Input:       "20",
		Expected:    0,
		Description: "Too short for year",
	},
	{
		Name:        "future date",
		Input:       "2030-12-25",
		Expected:    2030,
		Description: "Future release date",
	},
	{
		Name:        "old date",
		Input:       "1920-01-01",
		Expected:    1920,
		Description: "Very old release date",
	},
}

// RatingTestCases fixture.
var RatingTestCases = []struct {
	Name        string
	Rating      float64
	Votes       int
	Expected    string
	Description string
}{
	{
		Name:        "high rating with many votes",
		Rating:      8.7,
		Votes:       1800000,
		Expected:    "8.7",
		Description: "Confident rating",
	},
	{
		Name:        "low rating with few votes",
		Rating:      5.2,
		Votes:       5, // Less than MinVotesForRating (10)
		Expected:    "N/A",
		Description: "Too few votes",
	},
	{
		Name:        "uncertain rating",
		Rating:      7.1,
		Votes:       75, // Between MinVotesForRating (10) and MinVotesForUncertain (100)
		Expected:    "7.1?",
		Description: "Uncertain rating",
	},
	{
		Name:        "perfect rating",
		Rating:      10.0,
		Votes:       500000,
		Expected:    "10.0",
		Description: "Perfect score",
	},
	{
		Name:        "zero rating",
		Rating:      0.0,
		Votes:       5,
		Expected:    "N/A",
		Description: "No rating available",
	},
}

// VoteTestCases fixtures.
var VoteTestCases = []struct {
	Name        string
	Votes       int
	Expected    string
	Description string
}{
	{
		Name:        "small number",
		Votes:       123,
		Expected:    "123",
		Description: "Under 1000",
	},
	{
		Name:        "thousands",
		Votes:       1500,
		Expected:    "1.5K",
		Description: "Thousands format",
	},
	{
		Name:        "millions",
		Votes:       2500000,
		Expected:    "2.5M",
		Description: "Millions format",
	},
	{
		Name:        "exact thousand",
		Votes:       1000,
		Expected:    "1.0K",
		Description: "Exact thousand",
	},
	{
		Name:        "exact million",
		Votes:       1000000,
		Expected:    "1.0M",
		Description: "Exact million",
	},
	{
		Name:        "zero votes",
		Votes:       0,
		Expected:    "0",
		Description: "No votes",
	},
}

// HelperTestCases fixtures.
var HelperTestCases = struct {
	HighlyRated []struct {
		Rating   float64
		Votes    int
		Expected bool
	}
	Popular []struct {
		Popularity float64
		Votes      int
		Expected   bool
	}
	Recent []struct {
		Year     int
		Expected bool
	}
}{
	HighlyRated: []struct {
		Rating   float64
		Votes    int
		Expected bool
	}{
		{8.5, 200, true},  // High rating, sufficient votes
		{8.5, 50, false},  // High rating, insufficient votes
		{7.0, 200, false}, // Low rating, sufficient votes
		{7.0, 50, false},  // Low rating, insufficient votes
	},
	Popular: []struct {
		Popularity float64
		Votes      int
		Expected   bool
	}{
		{60.0, 500, true},  // High popularity
		{30.0, 1500, true}, // High votes
		{30.0, 500, false}, // Neither high popularity nor high votes
		{60.0, 1500, true}, // Both high
	},
	Recent: []struct {
		Year     int
		Expected bool
	}{
		{time.Now().Year(), true},       // Current year
		{time.Now().Year() - 1, true},   // Last year
		{time.Now().Year() - 2, true},   // Two years ago
		{time.Now().Year() - 3, false},  // Three years ago
		{time.Now().Year() - 10, false}, // Old
	},
}
