// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"fmt"
	"slices"
	"strings"
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
			slices.Sort(available)
			return nil, fmt.Errorf("unknown genre: %s. Available: %s", name, strings.Join(available, ", "))
		}
	}
	return ids, nil
}
