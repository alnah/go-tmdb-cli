// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

// IsRecent determines if a year is recent.
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

// CleanQuery removes surrounding quotes from search queries.
func CleanQuery(query string) string {
	return strings.Trim(query, `"'`)
}
