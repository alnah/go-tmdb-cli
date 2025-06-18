// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import "context"

// MovieLister provides methods for listing movies.
type MovieLister interface {
	GetPopularMovies(ctx context.Context, maxItems int) ([]Movie, error)
	GetTopRatedMovies(ctx context.Context, maxItems int) ([]Movie, error)
	GetNowPlayingMovies(ctx context.Context, maxItems int) ([]Movie, error)
	GetUpcomingMovies(ctx context.Context, maxItems int) ([]Movie, error)
}

// MovieSearcher provides methods for searching movies.
type MovieSearcher interface {
	SearchMovies(ctx context.Context, query string, maxItems int) ([]Movie, error)
	DiscoverMovies(ctx context.Context, opts SearchOptions) ([]Movie, error)
}

// TVShowLister provides methods for listing TV shows.
type TVShowLister interface {
	GetPopularTVShows(ctx context.Context, maxItems int) ([]TVShow, error)
	GetTopRatedTVShows(ctx context.Context, maxItems int) ([]TVShow, error)
	GetOnTheAirTVShows(ctx context.Context, maxItems int) ([]TVShow, error)
}

// TVShowSearcher provides methods for searching TV shows.
type TVShowSearcher interface {
	SearchTVShows(ctx context.Context, query string, maxItems int) ([]TVShow, error)
	DiscoverTVShows(ctx context.Context, opts SearchOptions) ([]TVShow, error)
}

// MovieClient combines all movie-related operations.
type MovieClient interface {
	MovieLister
	MovieSearcher
}

// TVClient combines all TV show-related operations.
type TVClient interface {
	TVShowLister
	TVShowSearcher
}

// Searcher combines movie and TV search capabilities.
type Searcher interface {
	MovieSearcher
	TVShowSearcher
}

// TMDBClient combines all TMDB operations.
type TMDBClient interface {
	MovieClient
	TVClient
}

// Ensure Client implements all interfaces.
var (
	_ MovieLister    = (*Client)(nil)
	_ MovieSearcher  = (*Client)(nil)
	_ TVShowLister   = (*Client)(nil)
	_ TVShowSearcher = (*Client)(nil)
	_ MovieClient    = (*Client)(nil)
	_ TVClient       = (*Client)(nil)
	_ TMDBClient     = (*Client)(nil)
)
