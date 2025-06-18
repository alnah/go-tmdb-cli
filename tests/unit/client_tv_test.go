package unit

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/fixtures"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

func createInstantTestClientTV(config internal.Config, httpClient *http.Client) *internal.Client {
	client := internal.NewTestClient(config, httpClient, nil)

	// Replace the rate limiter with one that allows instant requests
	instantRateLimiter := rate.NewLimiter(rate.Limit(1000), 1000) // Very high rate
	client.SetRateLimiter(instantRateLimiter)

	return client
}

func TestGetPopularTVShows_SinglePage(t *testing.T) {
	t.Parallel()

	transport := helpers.NewMockTransport()
	transport.SetResponse("/tv/popular", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(1, 1, 20, 1),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClientTV(config, httpClient)

	tvShows, err := client.GetPopularTVShows(context.Background(), 20)

	assert.NoError(t, err)
	assert.Len(t, tvShows, 1)
	assert.Equal(t, "Test TV Show 1", tvShows[0].Name)
	assert.Equal(t, 1, transport.RequestCount)
}

func TestGetPopularTVShows_Pagination(t *testing.T) {
	t.Parallel()

	transport := helpers.NewMockTransport()
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
	client := createInstantTestClientTV(config, httpClient)

	tvShows, err := client.GetPopularTVShows(context.Background(), 50)

	assert.NoError(t, err)
	assert.Len(t, tvShows, 50)
	assert.Equal(t, 3, transport.RequestCount)
}

func TestGetPopularTVShows_APIError(t *testing.T) {
	t.Parallel()

	transport := helpers.NewMockTransport()
	transport.SetResponse("/tv/popular", helpers.MockHTTPResponse{
		StatusCode: 401,
		Body:       []byte(`{"status_message":"Invalid API key","status_code":7}`),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClientTV(config, httpClient)

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

	transport := helpers.NewMockTransport()
	transport.SetResponse("/search/tv", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(1, 1, 5, 5),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClientTV(config, httpClient)

	tvShows, err := client.SearchTVShows(context.Background(), "Office", 20)

	assert.NoError(t, err)
	assert.NotNil(t, tvShows)
	assert.Len(t, tvShows, 5)
}

func TestSearchTVShows_Pagination(t *testing.T) {
	t.Parallel()

	transport := helpers.NewMockTransport()
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
	client := createInstantTestClientTV(config, httpClient)

	tvShows, err := client.SearchTVShows(context.Background(), "comedy", 25)

	assert.NoError(t, err)
	assert.Len(t, tvShows, 25)
	assert.Equal(t, 2, transport.RequestCount)
}

func TestDiscoverTVShows_Parameters(t *testing.T) {
	t.Parallel()

	transport := helpers.NewMockTransport()
	transport.SetResponse("/discover/tv", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(1, 1, 10, 10),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClientTV(config, httpClient)

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
	url := transport.LastURL
	assert.Contains(t, url, "with_genres=")
	assert.Contains(t, url, "vote_average.gte=7")
	assert.Contains(t, url, "vote_average.lte=9")
	assert.Contains(t, url, "first_air_date_year=2023")
}

func TestGetTopRatedTVShows(t *testing.T) {
	t.Parallel()

	transport := helpers.NewMockTransport()
	transport.SetResponse("/tv/top_rated", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(1, 1, 15, 15),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClientTV(config, httpClient)

	tvShows, err := client.GetTopRatedTVShows(context.Background(), 15)

	assert.NoError(t, err)
	assert.Len(t, tvShows, 15)
}

func TestGetOnTheAirTVShows(t *testing.T) {
	t.Parallel()

	transport := helpers.NewMockTransport()
	transport.SetResponse("/tv/on_the_air", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(1, 1, 10, 10),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClientTV(config, httpClient)

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

	transport := helpers.NewMockTransport()
	transport.SetResponse("/tv/popular", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(1, 1, 20, 20),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	config.CacheTTL = 5 * time.Minute
	client := createInstantTestClientTV(config, httpClient)

	// First request
	tvShows1, err := client.GetPopularTVShows(context.Background(), 20)
	assert.NoError(t, err)
	assert.Len(t, tvShows1, 20)

	// Second request should use cache
	tvShows2, err := client.GetPopularTVShows(context.Background(), 20)
	assert.NoError(t, err)
	assert.Len(t, tvShows2, 20)

	// Only one HTTP request should have been made
	assert.Equal(t, 1, transport.RequestCount)
}

func TestGetPopularTVShows_NetworkError(t *testing.T) {
	t.Parallel()

	// Using the failingTransport from movies test or define it here
	transport := &failingTransport{
		err: errors.New("connection refused"),
	}

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	config.MaxRetries = 0 // No retries for faster test
	client := createInstantTestClientTV(config, httpClient)

	tvShows, err := client.GetPopularTVShows(context.Background(), 20)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
	assert.Nil(t, tvShows)
}

// Test concurrent access for TV operations.
func TestTVOperations_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	transport := helpers.NewMockTransport()
	// Set up responses for all TV endpoints
	endpoints := []string{
		"/tv/popular",
		"/tv/top_rated",
		"/tv/on_the_air",
		"/search/tv",
	}

	for _, endpoint := range endpoints {
		transport.SetResponse(endpoint, helpers.MockHTTPResponse{
			StatusCode: 200,
			Body:       fixtures.CreateTVPageResponse(1, 1, 5, 5),
		})
	}

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClientTV(config, httpClient)

	// Run concurrent requests
	var wg sync.WaitGroup
	errors := make(chan error, len(endpoints))

	for _, endpoint := range endpoints {
		wg.Add(1)
		go func(ep string) {
			defer wg.Done()

			var err error
			switch ep {
			case "/tv/popular":
				_, err = client.GetPopularTVShows(context.Background(), 5)
			case "/tv/top_rated":
				_, err = client.GetTopRatedTVShows(context.Background(), 5)
			case "/tv/on_the_air":
				_, err = client.GetOnTheAirTVShows(context.Background(), 5)
			case "/search/tv":
				_, err = client.SearchTVShows(context.Background(), "test", 5)
			}

			if err != nil {
				errors <- err
			}
		}(endpoint)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		assert.NoError(t, err)
	}
}

// Test TV-specific parameters.
func TestDiscoverTVShows_TVSpecificParams(t *testing.T) {
	t.Parallel()

	transport := helpers.NewMockTransport()
	transport.SetResponse("/discover/tv", helpers.MockHTTPResponse{
		StatusCode: 200,
		Body:       fixtures.CreateTVPageResponse(1, 1, 10, 10),
	})

	httpClient := &http.Client{Transport: transport}
	config := helpers.MockValidConfig()
	client := createInstantTestClientTV(config, httpClient)

	// Test with TV-specific search options
	opts := helpers.NewSearchOptionsBuilder().
		WithYear(2023). // Should map to first_air_date_year for TV
		WithMaxItems(10).
		Build()

	tvShows, err := client.DiscoverTVShows(context.Background(), opts)

	assert.NoError(t, err)
	assert.NotNil(t, tvShows)

	// Verify TV-specific parameter mapping
	url := transport.LastURL
	assert.Contains(t, url, "first_air_date_year=2023") // Not primary_release_year
}

// Test invalid date formats.
func TestConvertTVShows_InvalidDates(t *testing.T) {
	t.Parallel()

	config := helpers.MockValidConfig()
	client := internal.NewClient(config)

	tmdbTVShows := []internal.TMDBTVShow{
		{
			ID:           1,
			Name:         "Test Show",
			FirstAirDate: "invalid-date", // Invalid date format
		},
		{
			ID:           2,
			Name:         "Test Show 2",
			FirstAirDate: "", // Empty date
		},
	}

	tvShows := client.ConvertTVShows(tmdbTVShows)

	require.Len(t, tvShows, 2)
	assert.Equal(t, 0, tvShows[0].Year) // Should default to 0 for invalid date
	assert.Equal(t, 0, tvShows[1].Year) // Should default to 0 for empty date
}

// Benchmark test for TV pagination.
func BenchmarkSearchTVShows_Pagination(b *testing.B) {
	transport := helpers.NewMockTransport()
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
	client := createInstantTestClientTV(config, httpClient)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.SearchTVShows(ctx, "comedy", 25)
	}
}
