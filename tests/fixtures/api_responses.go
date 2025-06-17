package fixtures

import (
	"encoding/json"
	"fmt"

	"github.com/alnah/tmdb-cli/internal"
)

// CreateMoviePageResponse creates a mock TMDB movie response with specified parameters.
func CreateMoviePageResponse(page, totalPages, totalResults, itemCount int) []byte {
	movies := make([]internal.TMDBMovie, itemCount)
	for i := range itemCount {
		movies[i] = internal.TMDBMovie{
			ID:            (page-1)*20 + i + 1,
			Title:         fmt.Sprintf("Test Movie %d", (page-1)*20+i+1),
			OriginalTitle: fmt.Sprintf("Original Test Movie %d", (page-1)*20+i+1),
			ReleaseDate:   "2023-06-15",
			VoteAverage:   7.5 + float64(i%30)/10,
			VoteCount:     1000 + i*100,
			Popularity:    50.0 + float64(i),
			GenreIDs:      []int{28, 12}, // Action, Adventure
			Overview:      fmt.Sprintf("Test overview for movie %d", (page-1)*20+i+1),
			Adult:         false,
		}
	}

	response := internal.TMDBResponse{
		Page:         page,
		TotalPages:   totalPages,
		TotalResults: totalResults,
		Results:      movies,
	}

	data, _ := json.Marshal(response)
	return data
}

// CreateTVPageResponse creates a mock TMDB TV response with specified parameters.
func CreateTVPageResponse(page, totalPages, totalResults, itemCount int) []byte {
	shows := make([]internal.TMDBTVShow, itemCount)
	for i := range itemCount {
		shows[i] = internal.TMDBTVShow{
			ID:               (page-1)*20 + i + 1,
			Name:             fmt.Sprintf("Test TV Show %d", (page-1)*20+i+1),
			OriginalName:     fmt.Sprintf("Original Test TV Show %d", (page-1)*20+i+1),
			FirstAirDate:     "2023-06-15",
			VoteAverage:      7.5 + float64(i%30)/10,
			VoteCount:        1000 + i*100,
			Popularity:       50.0 + float64(i),
			GenreIDs:         []int{10759, 16}, // Action & Adventure, Animation
			Overview:         fmt.Sprintf("Test overview for TV show %d", (page-1)*20+i+1),
			OriginalLanguage: "en",
			Adult:            false,
		}
	}

	response := internal.TMDBTVResponse{
		Page:         page,
		TotalPages:   totalPages,
		TotalResults: totalResults,
		Results:      shows,
	}

	data, _ := json.Marshal(response)
	return data
}

// CreateEmptyMovieResponse creates an empty movie response.
func CreateEmptyMovieResponse() []byte {
	response := internal.TMDBResponse{
		Page:         1,
		TotalPages:   0,
		TotalResults: 0,
		Results:      []internal.TMDBMovie{},
	}

	data, _ := json.Marshal(response)
	return data
}

// CreateEmptyTVResponse creates an empty TV response.
func CreateEmptyTVResponse() []byte {
	response := internal.TMDBTVResponse{
		Page:         1,
		TotalPages:   0,
		TotalResults: 0,
		Results:      []internal.TMDBTVShow{},
	}

	data, _ := json.Marshal(response)
	return data
}

// CreateDetailedTVResponse creates a more detailed TV response for testing.
func CreateDetailedTVResponse() []byte {
	response := internal.TMDBTVResponse{
		Page:         1,
		TotalPages:   1,
		TotalResults: 2,
		Results: []internal.TMDBTVShow{
			{
				ID:               1,
				Name:             "Breaking Bad",
				OriginalName:     "Breaking Bad",
				FirstAirDate:     "2008-01-20",
				VoteAverage:      9.5,
				VoteCount:        12000,
				Popularity:       250.5,
				GenreIDs:         []int{18, 80}, // Drama, Crime
				Overview:         "A chemistry teacher turned methamphetamine manufacturer",
				OriginalLanguage: "en",
				Adult:            false,
			},
			{
				ID:               2,
				Name:             "Squid Game",
				OriginalName:     "오징어 게임",
				FirstAirDate:     "2021-09-17",
				VoteAverage:      7.8,
				VoteCount:        11000,
				Popularity:       450.3,
				GenreIDs:         []int{10759, 9648, 18}, // Action & Adventure, Mystery, Drama
				Overview:         "Hundreds of cash-strapped players accept a strange invitation",
				OriginalLanguage: "ko",
				Adult:            true,
			},
		},
	}

	data, _ := json.Marshal(response)
	return data
}

// ValidMovieSearchResponse is a pre-defined valid movie search response.
var ValidMovieSearchResponse = CreateMoviePageResponse(1, 1, 1, 1)

// ValidTVSearchResponse is a pre-defined valid TV search response.
var ValidTVSearchResponse = CreateTVPageResponse(1, 1, 1, 1)
