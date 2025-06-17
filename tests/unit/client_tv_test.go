package unit

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func TestGetPopularTVShows_SinglePage(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/tv/popular", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(1, 1, 20, 1),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	tvShows, err := client.GetPopularTVShows(context.Background(), 20)

	assert.NoError(t, err)
	assert.Len(t, tvShows, 1)
	assert.Equal(t, "Test TV Show 1", tvShows[0].Name)
	assert.Equal(t, 1, transport.requestCount)
}

func TestGetPopularTVShows_Pagination(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("page=1", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(1, 3, 60, 20),
	})
	transport.SetResponse("page=2", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(2, 3, 60, 20),
	})
	transport.SetResponse("page=3", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(3, 3, 60, 10),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	tvShows, err := client.GetPopularTVShows(context.Background(), 50)

	assert.NoError(t, err)
	assert.Len(t, tvShows, 50)
	assert.Equal(t, 3, transport.requestCount)
}

func TestGetPopularTVShows_APIError(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/tv/popular", helpers.MockHTTPResponse{
		StatusCode: 401,
		Body:       []byte(`{"status_message":"Invalid API key","status_code":7}`),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	tvShows, err := client.GetPopularTVShows(context.Background(), 20)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid TMDB API key")
	assert.Nil(t, tvShows)
}

func TestSearchTVShows_EmptyQuery(t *testing.T) {
	t.Parallel()

	config := helpers.MockValidConfig()
	client := internal.NewClient(config)

	tvShows, err := client.SearchTVShows(context.Background(), "", 20)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "search query cannot be empty")
	assert.Nil(t, tvShows)
}

func TestSearchTVShows_Success(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/search/tv", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(1, 1, 5, 5),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	tvShows, err := client.SearchTVShows(context.Background(), "Office", 20)

	assert.NoError(t, err)
	assert.NotNil(t, tvShows)
	assert.Len(t, tvShows, 5)
}

func TestSearchTVShows_Pagination(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("page=1", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(1, 2, 30, 20),
	})
	transport.SetResponse("page=2", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(2, 2, 30, 10),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	tvShows, err := client.SearchTVShows(context.Background(), "comedy", 25)

	assert.NoError(t, err)
	assert.Len(t, tvShows, 25)
	assert.Equal(t, 2, transport.requestCount)
}

func TestDiscoverTVShows_Parameters(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/discover/tv", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(1, 1, 10, 10),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	opts := helpers.NewSearchOptionsBuilder().
		WithIncludeGenres([]int{10759, 16}).
		WithMinRating(7.0).
		WithMaxRating(9.0).
		WithYear(2023).
		WithMaxItems(10).
		Build()

	tvShows, err := client.DiscoverTVShows(context.Background(), opts)

	assert.NoError(t, err)
	assert.NotNil(t, tvShows)
	assert.Len(t, tvShows, 10)

	// URL should contain TV-specific parameters
	url := transport.lastURL
	assert.Contains(t, url, "with_genres=")
	assert.Contains(t, url, "vote_average.gte=7")
	assert.Contains(t, url, "vote_average.lte=9")
	assert.Contains(t, url, "first_air_date_year=2023")
}

func TestGetTopRatedTVShows(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/tv/top_rated", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(1, 1, 15, 15),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	tvShows, err := client.GetTopRatedTVShows(context.Background(), 15)

	assert.NoError(t, err)
	assert.Len(t, tvShows, 15)
}

func TestGetOnTheAirTVShows(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/tv/on_the_air", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(1, 1, 10, 10),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	tvShows, err := client.GetOnTheAirTVShows(context.Background(), 10)

	assert.NoError(t, err)
	assert.Len(t, tvShows, 10)
}

func TestConvertTVShows(t *testing.T) {
	t.Parallel()

	config := helpers.MockValidConfig()
	client := internal.NewClient(config)

	client.SetTVGenre(10759, "Action & Adventure")
	client.SetTVGenre(16, "Animation")

	tmdbTVShows := []internal.TMDBTVShow{
		{
			ID:               1,
			Name:             "Test Show",
			OriginalName:     "Original Test Show",
			FirstAirDate:     "2023-06-15",
			VoteAverage:      8.5,
			VoteCount:        1000,
			Popularity:       75.5,
			GenreIDs:         []int{10759, 16},
			Overview:         "Test overview",
			OriginalLanguage: "en",
			Adult:            false,
		},
	}

	tvShows := client.ConvertTVShows(tmdbTVShows)

	require.Len(t, tvShows, 1)
	assert.Equal(t, 1, tvShows[0].ID)
	assert.Equal(t, "Test Show", tvShows[0].Name)
	assert.Equal(t, 2023, tvShows[0].Year)
	assert.Equal(t, "Action & Adventure, Animation", tvShows[0].Genres)
}

func TestConvertTVShows_EmptyName(t *testing.T) {
	t.Parallel()

	config := helpers.MockValidConfig()
	client := internal.NewClient(config)

	tmdbTVShows := []internal.TMDBTVShow{
		{
			ID:           1,
			Name:         "", // Empty name
			FirstAirDate: "2023-01-01",
		},
	}

	tvShows := client.ConvertTVShows(tmdbTVShows)

	require.Len(t, tvShows, 1)
	assert.Equal(t, "Unknown TV Show", tvShows[0].Name)
}

func TestGetPopularTVShows_Cache(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/tv/popular", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(1, 1, 20, 20),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	config.CacheTTL = 5 * time.Minute
	client := internal.NewTestClient(config, httpClient, nil)

	// First request
	tvShows1, err := client.GetPopularTVShows(context.Background(), 20)
	assert.NoError(t, err)
	assert.Len(t, tvShows1, 20)

	// Second request should use cache
	tvShows2, err := client.GetPopularTVShows(context.Background(), 20)
	assert.NoError(t, err)
	assert.Len(t, tvShows2, 20)

	// Only one HTTP request should have been made
	assert.Equal(t, 1, transport.requestCount)
}

func TestGetPopularTVShows_NetworkError(t *testing.T) {
	t.Parallel()

	transport := NewMockTransport()
	transport.SetResponse("/tv/popular", helpers.MockHTTPResponse{
		Error: errors.New("connection refused"),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := internal.NewTestClient(config, httpClient, nil)

	tvShows, err := client.GetPopularTVShows(context.Background(), 20)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
	assert.Nil(t, tvShows)
}
