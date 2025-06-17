package unit

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

// MockTransport implements http.RoundTripper for testing.
type MockTransport struct {
	responses    map[string]helpers.MockHTTPResponse
	requestCount int
	lastURL      string
	calls        []string
}

func NewMockTransport() *MockTransport {
	return &MockTransport{
		responses: make(map[string]helpers.MockHTTPResponse),
		calls:     make([]string, 0),
	}
}

func (mt *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	mt.requestCount++
	mt.lastURL = req.URL.String()
	mt.calls = append(mt.calls, req.URL.Path)

	for pattern, mockResp := range mt.responses {
		if contains(req.URL.String(), pattern) {
			if mockResp.Error != nil {
				return nil, mockResp.Error
			}
			return &http.Response{
				StatusCode: mockResp.StatusCode,
				Body:       helpers.CreateResponseBody(mockResp.Body),
				Header:     make(http.Header),
			}, nil
		}
	}

	return &http.Response{
		StatusCode: 404,
		Body:       helpers.CreateResponseBody([]byte(`{"error": "Not found"}`)),
		Header:     make(http.Header),
	}, nil
}

func (mt *MockTransport) SetResponse(pattern string, response helpers.MockHTTPResponse) {
	mt.responses[pattern] = response
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && hasSubstring(s, substr)
}

func hasSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestGetPopularMovies_SinglePage(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/movie/popular", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 1, 20, 1),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	movies, err := client.GetPopularMovies(context.Background(), 20)

	assert.NoError(t, err)
	assert.Len(t, movies, 1)
	assert.Equal(t, "Test Movie 1", movies[0].Title)
	assert.Equal(t, 1, transport.requestCount)
}

func TestGetPopularMovies_Pagination(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("page=1", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 3, 60, 20),
	})
	transport.SetResponse("page=2", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(2, 3, 60, 20),
	})
	transport.SetResponse("page=3", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(3, 3, 60, 10),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	movies, err := client.GetPopularMovies(context.Background(), 50)

	assert.NoError(t, err)
	assert.Len(t, movies, 50)
	assert.Equal(t, 3, transport.requestCount)
}

func TestGetPopularMovies_MaxItemsLimit(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/movie/popular", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 1, 20, 20),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	movies, err := client.GetPopularMovies(context.Background(), 5)

	assert.NoError(t, err)
	assert.Len(t, movies, 5)
	assert.Equal(t, 1, transport.requestCount)
}

func TestGetPopularMovies_APIError(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/movie/popular", helpers.MockHTTPResponse{
		StatusCode: 401,
		Body:       []byte(`{"status_message":"Invalid API key","status_code":7}`),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	movies, err := client.GetPopularMovies(context.Background(), 20)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid TMDB API key")
	assert.Nil(t, movies)
}

func TestSearchMovies_EmptyQuery(t *testing.T) {
	t.Parallel()

	config := helpers.MockValidConfig()
	client := internal.NewClient(config)

	movies, err := client.SearchMovies(context.Background(), "", 20)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "search query cannot be empty")
	assert.Nil(t, movies)
}

func TestSearchMovies_Success(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/search/movie", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 1, 5, 5),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	movies, err := client.SearchMovies(context.Background(), "Matrix", 20)

	assert.NoError(t, err)
	assert.NotNil(t, movies)
	assert.Len(t, movies, 5)
	assert.Contains(t, transport.lastURL, "query=")
	assert.Contains(t, transport.lastURL, "Matrix")
}

func TestSearchMovies_Pagination(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("page=1", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 2, 30, 20),
	})
	transport.SetResponse("page=2", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(2, 2, 30, 10),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	movies, err := client.SearchMovies(context.Background(), "action", 25)

	assert.NoError(t, err)
	assert.Len(t, movies, 25)
	assert.Equal(t, 2, transport.requestCount)
}

func TestDiscoverMovies_Parameters(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/discover/movie", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 1, 10, 10),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	opts := helpers.NewSearchOptionsBuilder().
		WithIncludeGenres([]int{28, 35}).
		WithMinRating(7.0).
		WithMaxRating(9.0).
		WithYear(2023).
		WithMinVotes(1000).
		WithMaxItems(10).
		Build()

	movies, err := client.DiscoverMovies(context.Background(), opts)

	assert.NoError(t, err)
	assert.NotNil(t, movies)
	assert.Len(t, movies, 10)

	// URL should contain parameters
	url := transport.lastURL
	assert.Contains(t, url, "with_genres=")
	assert.Contains(t, url, "vote_average.gte=7")
	assert.Contains(t, url, "vote_average.lte=9")
	assert.Contains(t, url, "primary_release_year=2023")
}

func TestConvertMovies(t *testing.T) {
	t.Parallel()

	config := helpers.MockValidConfig()
	client := internal.NewClient(config)

	client.SetMovieGenre(28, "Action")
	client.SetMovieGenre(12, "Adventure")

	tmdbMovies := []internal.TMDBMovie{
		{
			ID:            1,
			Title:         "Test Movie",
			OriginalTitle: "Original Test Movie",
			ReleaseDate:   "2023-06-15",
			VoteAverage:   8.5,
			VoteCount:     1000,
			Popularity:    75.5,
			GenreIDs:      []int{28, 12},
			Overview:      "Test overview",
			Adult:         false,
		},
	}

	movies := client.ConvertMovies(tmdbMovies)

	require.Len(t, movies, 1)
	assert.Equal(t, 1, movies[0].ID)
	assert.Equal(t, "Test Movie", movies[0].Title)
	assert.Equal(t, 2023, movies[0].Year)
	assert.Equal(t, "Action, Adventure", movies[0].Genres)
}

func TestBuildBaseDiscoverParams(t *testing.T) {
	t.Parallel()

	config := helpers.MockValidConfig()
	client := internal.NewClient(config)

	opts := internal.SearchOptions{
		MinRating:     7.0,
		MaxRating:     9.0,
		MinVotes:      1000,
		IncludeGenres: []int{28, 35},
		ExcludeGenres: []int{27},
		Language:      "en",
		SortBy:        "popularity",
		SortOrder:     "desc",
	}

	params := client.BuildBaseDiscoverParams(opts)

	assert.Equal(t, "7.0", params.Get("vote_average.gte"))
	assert.Equal(t, "9.0", params.Get("vote_average.lte"))
	assert.Equal(t, "1000", params.Get("vote_count.gte"))
	assert.Equal(t, "28,35", params.Get("with_genres"))
	assert.Equal(t, "27", params.Get("without_genres"))
	assert.Equal(t, "en", params.Get("with_original_language"))
	assert.Equal(t, "popularity.desc", params.Get("sort_by"))
}

func TestGetPopularMovies_Cache(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/movie/popular", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateMoviePageResponse(1, 1, 20, 20),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	config.CacheTTL = 5 * time.Minute
	client := internal.NewTestClient(config, httpClient, nil)

	// First request
	movies1, err := client.GetPopularMovies(context.Background(), 20)
	assert.NoError(t, err)
	assert.Len(t, movies1, 20)

	// Second request should use cache
	movies2, err := client.GetPopularMovies(context.Background(), 20)
	assert.NoError(t, err)
	assert.Len(t, movies2, 20)

	// Only one HTTP request should have been made
	assert.Equal(t, 1, transport.requestCount)
}

func TestSearchMovies_RateLimit(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	// First request returns rate limit error
	transport.SetResponse("/search/movie", helpers.MockHTTPResponse{
		StatusCode: 429,
		Body:       []byte(`{"status_message":"Rate limit exceeded"}`),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	config.MaxRetries = 0 // No retries
	client := internal.NewTestClient(config, httpClient, nil)

	movies, err := client.SearchMovies(context.Background(), "test", 5)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rate limit exceeded")
	assert.Nil(t, movies)
}
