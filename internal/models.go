// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
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

// SearchResult represents paginated search results.
type SearchResult struct {
	Movies       []Movie `json:"movies"`
	Page         int     `json:"page"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
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

// Genre represents a movie genre.
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

// IsHighlyRated determines if a movie is highly rated.
func IsHighlyRated(rating float64, votes int) bool {
	return rating >= HighRatingThreshold && votes >= MinVotesForUncertain
}

// IsPopular determines if a movie is popular.
func IsPopular(popularity float64, votes int) bool {
	return popularity >= PopularityThreshold || votes >= PopularVotesThreshold
}

// IsRecent determines if a movie is recent.
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
