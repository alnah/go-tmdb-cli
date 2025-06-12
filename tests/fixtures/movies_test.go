// tests/fixtures/movies_test.go
package fixtures_test

import (
	"github.com/alnah/tmdb-cli/internal"
)

// PopularMoviesResponse simulates TMDB API popular movies response
var PopularMoviesResponse = internal.TMDBResponse{
	Page:         1,
	TotalPages:   500,
	TotalResults: 10000,
	Results: []internal.TMDBMovie{
		{
			ID:            550,
			Title:         "Fight Club",
			OriginalTitle: "Fight Club",
			Overview:      "A ticking-time-bomb insomniac and a slippery soap salesman channel primal male aggression into a shocking new form of therapy.",
			ReleaseDate:   "1999-10-15",
			VoteAverage:   8.433,
			VoteCount:     26280,
			GenreIDs:      []int{18},
			Popularity:    61.416,
			Adult:         false,
			Video:         false,
		},
		{
			ID:            13,
			Title:         "Forrest Gump",
			OriginalTitle: "Forrest Gump",
			Overview:      "A man with a low IQ has accomplished great things in his life and been present during significant historic events.",
			ReleaseDate:   "1994-07-06",
			VoteAverage:   8.471,
			VoteCount:     24951,
			GenreIDs:      []int{35, 18, 10749},
			Popularity:    56.983,
			Adult:         false,
			Video:         false,
		},
		{
			ID:            278,
			Title:         "The Shawshank Redemption",
			OriginalTitle: "The Shawshank Redemption",
			Overview:      "Framed in the 1940s for the double murder of his wife and her lover, upstanding banker Andy Dufresne begins a new life at the Shawshank prison.",
			ReleaseDate:   "1994-09-23",
			VoteAverage:   8.707,
			VoteCount:     23849,
			GenreIDs:      []int{18, 80},
			Popularity:    88.065,
			Adult:         false,
			Video:         false,
		},
	},
}

// TopRatedMoviesResponse simulates TMDB API top-rated movies response
var TopRatedMoviesResponse = internal.TMDBResponse{
	Page:         1,
	TotalPages:   543,
	TotalResults: 10858,
	Results: []internal.TMDBMovie{
		{
			ID:            278,
			Title:         "The Shawshank Redemption",
			OriginalTitle: "The Shawshank Redemption",
			Overview:      "Framed in the 1940s for the double murder of his wife and her lover, upstanding banker Andy Dufresne begins a new life at the Shawshank prison.",
			ReleaseDate:   "1994-09-23",
			VoteAverage:   8.707,
			VoteCount:     23849,
			GenreIDs:      []int{18, 80},
			Popularity:    88.065,
			Adult:         false,
			Video:         false,
		},
		{
			ID:            238,
			Title:         "The Godfather",
			OriginalTitle: "The Godfather",
			Overview:      "Spanning the years 1945 to 1955, a chronicle of the fictional Italian-American Corleone crime family.",
			ReleaseDate:   "1972-03-14",
			VoteAverage:   8.690,
			VoteCount:     17770,
			GenreIDs:      []int{18, 80},
			Popularity:    104.265,
			Adult:         false,
			Video:         false,
		},
		{
			ID:            424,
			Title:         "Schindler's List",
			OriginalTitle: "Schindler's List",
			Overview:      "The true story of how businessman Oskar Schindler saved over a thousand Jewish lives from the Nazis.",
			ReleaseDate:   "1993-12-15",
			VoteAverage:   8.565,
			VoteCount:     14149,
			GenreIDs:      []int{18, 36, 10752},
			Popularity:    52.781,
			Adult:         false,
			Video:         false,
		},
	},
}

// NowPlayingMoviesResponse simulates TMDB API now playing movies response
var NowPlayingMoviesResponse = internal.TMDBResponse{
	Page:         1,
	TotalPages:   15,
	TotalResults: 297,
	Results: []internal.TMDBMovie{
		{
			ID:            912649,
			Title:         "Venom: The Last Dance",
			OriginalTitle: "Venom: The Last Dance",
			Overview:      "Eddie and Venom are on the run. Hunted by both of their worlds and with the net closing in, the duo are forced into a devastating decision.",
			ReleaseDate:   "2024-10-22",
			VoteAverage:   6.542,
			VoteCount:     847,
			GenreIDs:      []int{878, 28, 12},
			Popularity:    4423.925,
			Adult:         false,
			Video:         false,
		},
		{
			ID:            1184918,
			Title:         "The Wild Robot",
			OriginalTitle: "The Wild Robot",
			Overview:      "After a shipwreck, an intelligent robot called Roz is stranded on an uninhabited island.",
			ReleaseDate:   "2024-09-12",
			VoteAverage:   8.568,
			VoteCount:     2847,
			GenreIDs:      []int{16, 878, 10751},
			Popularity:    3426.789,
			Adult:         false,
			Video:         false,
		},
		{
			ID:            933260,
			Title:         "The Substance",
			OriginalTitle: "The Substance",
			Overview:      "A fading celebrity decides to use a black market drug, a cell-replicating substance that temporarily creates a younger, better version of herself.",
			ReleaseDate:   "2024-09-07",
			VoteAverage:   7.302,
			VoteCount:     1847,
			GenreIDs:      []int{27, 878, 53},
			Popularity:    2156.432,
			Adult:         false,
			Video:         false,
		},
	},
}

// UpcomingMoviesResponse simulates TMDB API upcoming movies response
var UpcomingMoviesResponse = internal.TMDBResponse{
	Page:         1,
	TotalPages:   8,
	TotalResults: 156,
	Results: []internal.TMDBMovie{
		{
			ID:            558449,
			Title:         "Gladiator II",
			OriginalTitle: "Gladiator II",
			Overview:      "Years after witnessing the death of the revered hero Maximus at the hands of his uncle, Lucius is forced to enter the Colosseum.",
			ReleaseDate:   "2024-11-13",
			VoteAverage:   0.0,
			VoteCount:     0,
			GenreIDs:      []int{28, 12, 18},
			Popularity:    1247.893,
			Adult:         false,
			Video:         false,
		},
		{
			ID:            1100782,
			Title:         "Smile 2",
			OriginalTitle: "Smile 2",
			Overview:      "About to embark on a new world tour, global pop sensation Skye Riley begins experiencing increasingly terrifying and inexplicable events.",
			ReleaseDate:   "2024-10-16",
			VoteAverage:   6.8,
			VoteCount:     145,
			GenreIDs:      []int{27, 9648},
			Popularity:    987.654,
			Adult:         false,
			Video:         false,
		},
		{
			ID:            1034541,
			Title:         "Terrifier 3",
			OriginalTitle: "Terrifier 3",
			Overview:      "Five years after surviving Art the Clown's Halloween massacre, Sienna and Jonathan are still struggling to rebuild their shattered lives.",
			ReleaseDate:   "2024-10-09",
			VoteAverage:   7.2,
			VoteCount:     432,
			GenreIDs:      []int{27, 53},
			Popularity:    756.321,
			Adult:         false,
			Video:         false,
		},
	},
}

// SearchMatrixResponse simulates search results for "Matrix"
var SearchMatrixResponse = internal.TMDBResponse{
	Page:         1,
	TotalPages:   3,
	TotalResults: 42,
	Results: []internal.TMDBMovie{
		{
			ID:            603,
			Title:         "The Matrix",
			OriginalTitle: "The Matrix",
			Overview:      "Set in the 22nd century, The Matrix tells the story of a computer hacker who joins a group of underground insurgents fighting the vast and powerful computers who now rule the earth.",
			ReleaseDate:   "1999-03-30",
			VoteAverage:   8.2,
			VoteCount:     23853,
			GenreIDs:      []int{28, 878},
			Popularity:    82.965,
			Adult:         false,
			Video:         false,
		},
		{
			ID:            604,
			Title:         "The Matrix Reloaded",
			OriginalTitle: "The Matrix Reloaded",
			Overview:      "Six months after the events depicted in The Matrix, Neo has proved to be a good omen for the free humans, as more and more humans are being freed from the matrix.",
			ReleaseDate:   "2003-05-07",
			VoteAverage:   7.051,
			VoteCount:     9853,
			GenreIDs:      []int{12, 28, 53, 878},
			Popularity:    45.123,
			Adult:         false,
			Video:         false,
		},
		{
			ID:            605,
			Title:         "The Matrix Revolutions",
			OriginalTitle: "The Matrix Revolutions",
			Overview:      "The human city of Zion defends itself against the massive invasion of the machines as Neo fights to end the war at another front while also opposing the rogue Agent Smith.",
			ReleaseDate:   "2003-10-27",
			VoteAverage:   6.728,
			VoteCount:     8765,
			GenreIDs:      []int{12, 28, 53, 878},
			Popularity:    41.789,
			Adult:         false,
			Video:         false,
		},
		{
			ID:            624860,
			Title:         "The Matrix Resurrections",
			OriginalTitle: "The Matrix Resurrections",
			Overview:      "Plagued by strange memories, Neo's life takes an unexpected turn when he finds himself back inside the Matrix.",
			ReleaseDate:   "2021-12-16",
			VoteAverage:   5.534,
			VoteCount:     4567,
			GenreIDs:      []int{28, 878},
			Popularity:    28.456,
			Adult:         false,
			Video:         false,
		},
	},
}

// EmptySearchResponse simulates no results found
var EmptySearchResponse = internal.TMDBResponse{
	Page:         1,
	TotalPages:   1,
	TotalResults: 0,
	Results:      []internal.TMDBMovie{},
}

// GenresResponse simulates TMDB API genres response
var GenresResponse = internal.TMDBGenresResponse{
	Genres: []internal.Genre{
		{ID: 28, Name: "Action"},
		{ID: 12, Name: "Adventure"},
		{ID: 16, Name: "Animation"},
		{ID: 35, Name: "Comedy"},
		{ID: 80, Name: "Crime"},
		{ID: 99, Name: "Documentary"},
		{ID: 18, Name: "Drama"},
		{ID: 10751, Name: "Family"},
		{ID: 14, Name: "Fantasy"},
		{ID: 36, Name: "History"},
		{ID: 27, Name: "Horror"},
		{ID: 10402, Name: "Music"},
		{ID: 9648, Name: "Mystery"},
		{ID: 10749, Name: "Romance"},
		{ID: 878, Name: "Science Fiction"},
		{ID: 10770, Name: "TV Movie"},
		{ID: 53, Name: "Thriller"},
		{ID: 10752, Name: "War"},
		{ID: 37, Name: "Western"},
	},
}

// DiscoverActionMoviesResponse simulates discovery with action genre filter
var DiscoverActionMoviesResponse = internal.TMDBResponse{
	Page:         1,
	TotalPages:   1000,
	TotalResults: 20000,
	Results: []internal.TMDBMovie{
		{
			ID:            157336,
			Title:         "Interstellar",
			OriginalTitle: "Interstellar",
			Overview:      "The adventures of a group of explorers who make use of a newly discovered wormhole to surpass the limitations on human space travel.",
			ReleaseDate:   "2014-11-05",
			VoteAverage:   8.442,
			VoteCount:     32847,
			GenreIDs:      []int{12, 18, 878},
			Popularity:    112.345,
			Adult:         false,
			Video:         false,
		},
		{
			ID:            324857,
			Title:         "Spider-Man: Into the Spider-Verse",
			OriginalTitle: "Spider-Man: Into the Spider-Verse",
			Overview:      "Struggling to find his place in the world while juggling his dual identity, Miles Morales soon discovers a cosmic threat to the multiverse.",
			ReleaseDate:   "2018-12-06",
			VoteAverage:   8.4,
			VoteCount:     15789,
			GenreIDs:      []int{28, 12, 16, 878},
			Popularity:    87.654,
			Adult:         false,
			Video:         false,
		},
	},
}

// SinglePageResponse simulates a response with only one page
var SinglePageResponse = internal.TMDBResponse{
	Page:         1,
	TotalPages:   1,
	TotalResults: 2,
	Results: []internal.TMDBMovie{
		{
			ID:            123,
			Title:         "Test Movie 1",
			OriginalTitle: "Test Movie 1",
			Overview:      "This is a test movie for single page response",
			ReleaseDate:   "2024-01-01",
			VoteAverage:   7.5,
			VoteCount:     100,
			GenreIDs:      []int{18},
			Popularity:    50.0,
			Adult:         false,
			Video:         false,
		},
		{
			ID:            124,
			Title:         "Test Movie 2",
			OriginalTitle: "Test Movie 2",
			Overview:      "This is another test movie for single page response",
			ReleaseDate:   "2024-01-02",
			VoteAverage:   6.8,
			VoteCount:     75,
			GenreIDs:      []int{35},
			Popularity:    45.0,
			Adult:         false,
			Video:         false,
		},
	},
}

// ErrorResponses for testing various API error conditions

// UnauthorizedResponse simulates 401 error
var UnauthorizedResponse = `{
  "success": false,
  "status_code": 7,
  "status_message": "Invalid API key: You must be granted a valid key."
}`

// NotFoundResponse simulates 404 error
var NotFoundResponse = `{
  "success": false,
  "status_code": 34,
  "status_message": "The resource you requested could not be found."
}`

// RateLimitResponse simulates 429 error
var RateLimitResponse = `{
  "success": false,
  "status_code": 25,
  "status_message": "Your request count (41) is over the allowed limit of 40."
}`

// InternalServerErrorResponse simulates 500 error
var InternalServerErrorResponse = `{
  "success": false,
  "status_code": 11,
  "status_message": "Internal error: Something went wrong, contact TMDb."
}`

// ConvertedMovies provides pre-converted Movie objects for testing
var ConvertedMovies = []internal.Movie{
	{
		ID:            550,
		Title:         "Fight Club",
		OriginalTitle: "Fight Club",
		Year:          1999,
		Rating:        8.433,
		Votes:         26280,
		Popularity:    61.416,
		Genres:        "Drama",
		Overview:      "A ticking-time-bomb insomniac and a slippery soap salesman channel primal male aggression into a shocking new form of therapy.",
		Language:      "en",
		Adult:         false,
		ReleaseDate:   "1999-10-15",
	},
	{
		ID:            13,
		Title:         "Forrest Gump",
		OriginalTitle: "Forrest Gump",
		Year:          1994,
		Rating:        8.471,
		Votes:         24951,
		Popularity:    56.983,
		Genres:        "Comedy, Drama, Romance",
		Overview:      "A man with a low IQ has accomplished great things in his life and been present during significant historic events.",
		Language:      "en",
		Adult:         false,
		ReleaseDate:   "1994-07-06",
	},
	{
		ID:            278,
		Title:         "The Shawshank Redemption",
		OriginalTitle: "The Shawshank Redemption",
		Year:          1994,
		Rating:        8.707,
		Votes:         23849,
		Popularity:    88.065,
		Genres:        "Drama, Crime",
		Overview:      "Framed in the 1940s for the double murder of his wife and her lover, upstanding banker Andy Dufresne begins a new life at the Shawshank prison.",
		Language:      "en",
		Adult:         false,
		ReleaseDate:   "1994-09-23",
	},
}

// TestMoviesWithVariousRatings provides movies with different rating ranges for testing
var TestMoviesWithVariousRatings = []internal.Movie{
	{
		ID:     1,
		Title:  "Low Rated Movie",
		Rating: 3.2,
		Votes:  50,
		Year:   2020,
	},
	{
		ID:     2,
		Title:  "Average Movie",
		Rating: 6.5,
		Votes:  500,
		Year:   2021,
	},
	{
		ID:     3,
		Title:  "Highly Rated Movie",
		Rating: 8.7,
		Votes:  5000,
		Year:   2022,
	},
	{
		ID:     4,
		Title:  "Uncertain Rating",
		Rating: 7.8,
		Votes:  8, // Too few votes
		Year:   2023,
	},
}

// TestMoviesWithLongTitles for testing title truncation
var TestMoviesWithLongTitles = []internal.Movie{
	{
		ID:            1,
		Title:         "This is an Extremely Long Movie Title That Should Be Truncated When Displayed",
		OriginalTitle: "Dette er en Ekstremt Lang Filmtittel Som Burde Bli Forkortet Når Den Vises",
		Year:          2024,
		Rating:        7.5,
		Votes:         1000,
	},
	{
		ID:     2,
		Title:  "Short Title",
		Year:   2024,
		Rating: 8.0,
		Votes:  2000,
	},
}

// TestMoviesWithMixedData for comprehensive testing
var TestMoviesWithMixedData = []internal.Movie{
	{
		ID:          1,
		Title:       "Recent Blockbuster",
		Year:        2024,
		Rating:      8.5,
		Votes:       15000,
		Popularity:  95.5,
		Genres:      "Action, Adventure",
		Overview:    "A recent blockbuster with high ratings and popularity.",
		Adult:       false,
		ReleaseDate: "2024-06-01",
	},
	{
		ID:          2,
		Title:       "Classic Drama",
		Year:        1975,
		Rating:      8.9,
		Votes:       8000,
		Popularity:  25.3,
		Genres:      "Drama",
		Overview:    "A classic drama from the 1970s with excellent ratings.",
		Adult:       false,
		ReleaseDate: "1975-12-25",
	},
	{
		ID:          3,
		Title:       "B-Movie Horror",
		Year:        2010,
		Rating:      4.2,
		Votes:       150,
		Popularity:  12.1,
		Genres:      "Horror",
		Overview:    "A low-budget horror movie with poor ratings.",
		Adult:       false,
		ReleaseDate: "2010-10-31",
	},
	{
		ID:          4,
		Title:       "Adult Content",
		Year:        2015,
		Rating:      6.0,
		Votes:       500,
		Popularity:  30.0,
		Genres:      "Drama",
		Overview:    "An adult-rated drama.",
		Adult:       true,
		ReleaseDate: "2015-03-15",
	},
}

// HTTPResponseBodies for mocking HTTP responses

// ValidMovieListResponseBody provides a valid JSON response body
var ValidMovieListResponseBody = `{
  "page": 1,
  "results": [
    {
      "id": 550,
      "title": "Fight Club",
      "original_title": "Fight Club",
      "overview": "A ticking-time-bomb insomniac and a slippery soap salesman channel primal male aggression into a shocking new form of therapy.",
      "release_date": "1999-10-15",
      "vote_average": 8.433,
      "vote_count": 26280,
      "genre_ids": [18],
      "popularity": 61.416,
      "adult": false,
      "video": false
    }
  ],
  "total_pages": 500,
  "total_results": 10000
}`

// InvalidJSONResponseBody provides malformed JSON for testing error handling
var InvalidJSONResponseBody = `{
  "page": 1,
  "results": [
    {
      "id": 550,
      "title": "Fight Club",
      // Invalid JSON comment
      "vote_average": 8.433,
    }
  ]
}`

// MalformedDataResponseBody provides valid JSON with unexpected data structure
var MalformedDataResponseBody = `{
  "page": "not_a_number",
  "results": "not_an_array",
  "total_pages": null
}`
