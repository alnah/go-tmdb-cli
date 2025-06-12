// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Constants for common values.
const (
	MinValidYear          = 1888
	MaxReasonableYear     = 2030
	YearDigits            = 4
	OneThousand           = 1000
	OneMillion            = 1000000
	MinVotesForRating     = 10
	MinVotesForUncertain  = 100
	HighRatingThreshold   = 8.0
	PopularityThreshold   = 50.0
	PopularVotesThreshold = 1000
	RecentYearsBack       = 2
	DefaultTimeout        = 30
	DefaultCacheTTL       = 5
)

// TV-specific constants.
const (
	TVMinValidYear      = 1928 // First television broadcast
	TVMaxReasonableYear = 2030
	TVDefaultTimeout    = 30
	TVCacheTTL          = 5 // minutes
)

// TV Show endpoints.
const (
	TVPopularEndpoint  = "/tv/popular"
	TVTopRatedEndpoint = "/tv/top_rated"
	TVOnTheAirEndpoint = "/tv/on_the_air"
	TVSearchEndpoint   = "/search/tv"
	TVDiscoverEndpoint = "/discover/tv"
	TVGenresEndpoint   = "/genre/tv/list"
)

// TMDB genre ID constants.
const (
	GenreAction      = 28
	GenreAdventure   = 12
	GenreAnimation   = 16
	GenreComedy      = 35
	GenreCrime       = 80
	GenreDocumentary = 99
	GenreDrama       = 18
	GenreFamily      = 10751
	GenreFantasy     = 14
	GenreHistory     = 36
	GenreHorror      = 27
	GenreMusic       = 10402
	GenreMystery     = 9648
	GenreRomance     = 10749
	GenreSciFi       = 878
	GenreThriller    = 53
	GenreWar         = 10752
	GenreWestern     = 37
)

// Movie represents a simplified movie structure.
type Movie struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title,omitempty"`
	Year          int     `json:"year"`
	Rating        float64 `json:"rating"`
	Votes         int     `json:"votes"`
	Popularity    float64 `json:"popularity"`
	Genres        string  `json:"genres"`
	Overview      string  `json:"overview,omitempty"`
	Language      string  `json:"language"`
	Adult         bool    `json:"adult"`
	ReleaseDate   string  `json:"release_date,omitempty"`
}

// TVShow represents a simplified TV show structure.
type TVShow struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`          // Equivalent to Movie.Title
	OriginalName string  `json:"original_name"` // Equivalent to Movie.OriginalTitle
	Year         int     `json:"year"`          // Extracted from first_air_date
	Rating       float64 `json:"rating"`        // vote_average
	Votes        int     `json:"votes"`         // vote_count
	Popularity   float64 `json:"popularity"`
	Genres       string  `json:"genres"` // Mapped genres as string
	Overview     string  `json:"overview,omitempty"`
	Language     string  `json:"language"`
	Adult        bool    `json:"adult"`
	FirstAirDate string  `json:"first_air_date,omitempty"` // First broadcast date
}

// Validate validates the TVShow data structure.
func (tv *TVShow) Validate() error {
	if tv.Name == "" {
		return errors.New("TV show name is required")
	}
	if tv.ID <= 0 {
		return errors.New("TV show ID must be positive")
	}
	return nil
}

// SearchResult represents paginated search results.
type SearchResult struct {
	Movies       []Movie `json:"movies"`
	Page         int     `json:"page"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}

// TVSearchResult represents paginated TV search results.
type TVSearchResult struct {
	TVShows      []TVShow `json:"tv_shows"`
	Page         int      `json:"page"`
	TotalPages   int      `json:"total_pages"`
	TotalResults int      `json:"total_results"`
}

// SearchOptions represents search parameters.
type SearchOptions struct {
	Query         string
	Page          int
	Language      string
	Year          int
	MinRating     float64
	MaxRating     float64
	IncludeGenres []int
	ExcludeGenres []int
	SortBy        string
	SortOrder     string
	MaxItems      int
}

// Genre represents a movie/TV genre.
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Config represents application configuration.
type Config struct {
	APIKey     string        `yaml:"api_key"     env:"TMDB_API_KEY"`
	BaseURL    string        `yaml:"base_url"`
	Timeout    time.Duration `yaml:"timeout"`
	MaxRetries int           `yaml:"max_retries"`
	CacheTTL   time.Duration `yaml:"cache_ttl"`
	LogLevel   string        `yaml:"log_level"`
	Format     string        `yaml:"format"`
}

// TMDBResponse represents raw API response structure for movies.
type TMDBResponse struct {
	Page         int         `json:"page"`
	Results      []TMDBMovie `json:"results"`
	TotalPages   int         `json:"total_pages"`
	TotalResults int         `json:"total_results"`
}

// TMDBMovie represents the TMDB API movie structure.
type TMDBMovie struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title"`
	Overview      string  `json:"overview"`
	ReleaseDate   string  `json:"release_date"`
	VoteAverage   float64 `json:"vote_average"`
	VoteCount     int     `json:"vote_count"`
	GenreIDs      []int   `json:"genre_ids"`
	Popularity    float64 `json:"popularity"`
	Adult         bool    `json:"adult"`
	Video         bool    `json:"video"`
}

// TMDBTVResponse represents raw API response structure for TV shows.
type TMDBTVResponse struct {
	Page         int          `json:"page"`
	Results      []TMDBTVShow `json:"results"`
	TotalPages   int          `json:"total_pages"`
	TotalResults int          `json:"total_results"`
}

// TMDBTVShow represents the TMDB API TV show structure.
type TMDBTVShow struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	OriginalName     string  `json:"original_name"`
	Overview         string  `json:"overview"`
	FirstAirDate     string  `json:"first_air_date"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
	GenreIDs         []int   `json:"genre_ids"`
	Popularity       float64 `json:"popularity"`
	Adult            bool    `json:"adult"`
	OriginalLanguage string  `json:"original_language"`
}

// TMDBGenresResponse represents the TMDB genres API response.
type TMDBGenresResponse struct {
	Genres []Genre `json:"genres"`
}

// Helper functions (not methods to keep it simple)

// FormatRating formats rating with vote count context.
func FormatRating(rating float64, votes int) string {
	if votes < MinVotesForRating {
		return "N/A"
	}
	if votes < MinVotesForUncertain {
		return fmt.Sprintf("%.1f?", rating)
	}
	return fmt.Sprintf("%.1f", rating)
}

// FormatVotes formats vote count with human-readable suffixes.
func FormatVotes(votes int) string {
	if votes < OneThousand {
		return strconv.Itoa(votes)
	}
	if votes < OneMillion {
		return fmt.Sprintf("%.1fK", float64(votes)/OneThousand)
	}
	return fmt.Sprintf("%.1fM", float64(votes)/OneMillion)
}

// FormatGenres cleans up genre formatting.
func FormatGenres(genres string, maxWidth int) string {
	if genres == "" {
		return "N/A"
	}
	if len(genres) <= maxWidth {
		return genres
	}
	return genres[:maxWidth-3] + "..."
}

// ParseYear extracts year from release date string.
func ParseYear(releaseDate string) int {
	if releaseDate == "" {
		return 0
	}
	if len(releaseDate) >= YearDigits {
		if year, err := strconv.Atoi(releaseDate[:YearDigits]); err == nil {
			return year
		}
	}
	return 0
}

// ParseTVYear extracts year from first_air_date string.
func ParseTVYear(firstAirDate string) int {
	return ParseYear(firstAirDate)
}

// ShortenOverview truncates overview to specified length.
func ShortenOverview(overview string, maxLength int) string {
	if len(overview) <= maxLength {
		return overview
	}
	if maxLength < 10 {
		return overview[:maxLength] + "..."
	}
	truncated := overview[:maxLength-3]
	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > maxLength/2 {
		truncated = truncated[:lastSpace]
	}
	return truncated + "..."
}

// IsHighlyRated determines if a movie/TV show is highly rated.
func IsHighlyRated(rating float64, votes int) bool {
	return rating >= HighRatingThreshold && votes >= MinVotesForUncertain
}

// IsPopular determines if a movie/TV show is popular.
func IsPopular(popularity float64, votes int) bool {
	return popularity >= PopularityThreshold || votes >= PopularVotesThreshold
}

// IsRecent determines if a movie/TV show is recent.
func IsRecent(year int) bool {
	currentYear := time.Now().Year()
	return year >= currentYear-RecentYearsBack
}

// BuildDisplayTitle creates display title with original title if different.
func BuildDisplayTitle(title, originalTitle string) string {
	if originalTitle == "" || originalTitle == title {
		return title
	}
	return fmt.Sprintf("%s (%s)", title, originalTitle)
}

// BuildDisplayTVTitle creates display title for TV shows with original name if different.
func BuildDisplayTVTitle(name, originalName string) string {
	if originalName == "" || originalName == name {
		return name
	}
	return fmt.Sprintf("%s (%s)", name, originalName)
}

// ValidateSearchOptions validates search parameters.
func ValidateSearchOptions(opts *SearchOptions) error {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.MaxItems < 1 || opts.MaxItems > OneThousand {
		return fmt.Errorf("max-items must be between 1 and 1000, got %d", opts.MaxItems)
	}
	if opts.MinRating < 0 || opts.MinRating > 10 {
		return fmt.Errorf("min-rating must be between 0 and 10, got %.1f", opts.MinRating)
	}
	if opts.MaxRating < 0 || opts.MaxRating > 10 {
		return fmt.Errorf("max-rating must be between 0 and 10, got %.1f", opts.MaxRating)
	}
	if opts.MinRating > 0 && opts.MaxRating > 0 && opts.MinRating > opts.MaxRating {
		return fmt.Errorf(
			"min-rating (%.1f) cannot be greater than max-rating (%.1f)",
			opts.MinRating,
			opts.MaxRating,
		)
	}
	return nil
}

// DefaultConfig returns default configuration.
func DefaultConfig() Config {
	return Config{
		BaseURL:    "https://api.themoviedb.org/3",
		Timeout:    DefaultTimeout * time.Second,
		MaxRetries: 3,
		CacheTTL:   DefaultCacheTTL * time.Minute,
		LogLevel:   "info",
		Format:     "table",
	}
}

// GenreMap provides a mapping from genre names to TMDB genre IDs for discovery functionality.
var GenreMap = map[string]int{
	"action":          GenreAction,
	"adventure":       GenreAdventure,
	"animation":       GenreAnimation,
	"comedy":          GenreComedy,
	"crime":           GenreCrime,
	"documentary":     GenreDocumentary,
	"drama":           GenreDrama,
	"family":          GenreFamily,
	"fantasy":         GenreFantasy,
	"history":         GenreHistory,
	"horror":          GenreHorror,
	"music":           GenreMusic,
	"mystery":         GenreMystery,
	"romance":         GenreRomance,
	"science-fiction": GenreSciFi,
	"sci-fi":          GenreSciFi,
	"thriller":        GenreThriller,
	"war":             GenreWar,
	"western":         GenreWestern,
}

// ParseGenres converts genre names to IDs.
func ParseGenres(genreNames []string) ([]int, error) {
	var ids []int
	for _, name := range genreNames {
		name = strings.ToLower(strings.TrimSpace(name))
		if id, ok := GenreMap[name]; ok {
			ids = append(ids, id)
		} else {
			available := make([]string, 0, len(GenreMap))
			for k := range GenreMap {
				available = append(available, k)
			}
			return nil, fmt.Errorf("unknown genre: %s. Available: %s", name, strings.Join(available, ", "))
		}
	}
	return ids, nil
}

// MediaItem interface for unified processing of movies and TV shows.
type MediaItem interface {
	GetID() int
	GetTitle() string         // Name for TV, Title for Movie
	GetOriginalTitle() string // OriginalName for TV, OriginalTitle for Movie
	GetYear() int
	GetRating() float64
	GetVotes() int
	GetPopularity() float64
	GetGenres() string
	GetOverview() string
}

// GetID returns the movie ID.
func (m Movie) GetID() int { return m.ID }

// GetTitle returns the movie title.
func (m Movie) GetTitle() string { return m.Title }

// GetOriginalTitle returns the movie original title.
func (m Movie) GetOriginalTitle() string { return m.OriginalTitle }

// GetYear returns the movie year.
func (m Movie) GetYear() int { return m.Year }

// GetRating returns the movie rating.
func (m Movie) GetRating() float64 { return m.Rating }

// GetVotes returns the movie votes.
func (m Movie) GetVotes() int { return m.Votes }

// GetPopularity returns the movie popularity.
func (m Movie) GetPopularity() float64 { return m.Popularity }

// GetGenres returns the movie genres.
func (m Movie) GetGenres() string { return m.Genres }

// GetOverview returns the movie overview.
func (m Movie) GetOverview() string { return m.Overview }

// GetID returns the TV show ID.
func (tv TVShow) GetID() int { return tv.ID }

// GetTitle returns the TV show name.
func (tv TVShow) GetTitle() string { return tv.Name }

// GetOriginalTitle returns the TV show original name.
func (tv TVShow) GetOriginalTitle() string { return tv.OriginalName }

// GetYear returns the TV show year.
func (tv TVShow) GetYear() int { return tv.Year }

// GetRating returns the TV show rating.
func (tv TVShow) GetRating() float64 { return tv.Rating }

// GetVotes returns the TV show votes.
func (tv TVShow) GetVotes() int { return tv.Votes }

// GetPopularity returns the TV show popularity.
func (tv TVShow) GetPopularity() float64 { return tv.Popularity }

// GetGenres returns the TV show genres.
func (tv TVShow) GetGenres() string { return tv.Genres }

// GetOverview returns the TV show overview.
func (tv TVShow) GetOverview() string { return tv.Overview }
