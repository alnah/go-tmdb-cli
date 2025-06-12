// tests/fixtures/configs.go
package fixtures

import (
	"time"

	"github.com/alnah/tmdb-cli/internal"
)

// ValidConfig provides a valid configuration for testing
var ValidConfig = internal.Config{
	APIKey:     "test-api-key-123456789",
	BaseURL:    "https://api.themoviedb.org/3",
	Timeout:    30 * time.Second,
	MaxRetries: 3,
	CacheTTL:   5 * time.Minute,
	LogLevel:   "info",
	Format:     "table",
}

// MinimalConfig provides minimal required configuration
var MinimalConfig = internal.Config{
	APIKey:  "test-key",
	BaseURL: "https://api.example.com",
	Timeout: 10 * time.Second,
}

// InvalidConfigs for testing validation
var InvalidConfigs = []struct {
	Name   string
	Config internal.Config
	Error  string
}{
	{
		Name: "missing API key",
		Config: internal.Config{
			BaseURL: "https://api.themoviedb.org/3",
			Timeout: 30 * time.Second,
		},
		Error: "API key is required",
	},
	{
		Name: "negative timeout",
		Config: internal.Config{
			APIKey:  "test-key",
			BaseURL: "https://api.themoviedb.org/3",
			Timeout: -5 * time.Second,
		},
		Error: "timeout must be positive",
	},
	{
		Name: "invalid log level",
		Config: internal.Config{
			APIKey:   "test-key",
			BaseURL:  "https://api.themoviedb.org/3",
			Timeout:  30 * time.Second,
			LogLevel: "invalid-level",
		},
		Error: "invalid log level",
	},
	{
		Name: "invalid format",
		Config: internal.Config{
			APIKey:  "test-key",
			BaseURL: "https://api.themoviedb.org/3",
			Timeout: 30 * time.Second,
			Format:  "invalid-format",
		},
		Error: "invalid default format",
	},
}

// ConfigYAMLExamples for testing YAML configuration loading
var ValidConfigYAML = `
api_key: "test-api-key-from-yaml"
base_url: "https://api.themoviedb.org/3"
timeout: "45s"
max_retries: 5
cache_ttl: "10m"
log_level: "debug"
format: "json"
`

var InvalidConfigYAML = `
api_key: "test-key"
timeout: "invalid-duration"
max_retries: "not-a-number"
`

var PartialConfigYAML = `
api_key: "test-key"
log_level: "warn"
# Other fields will use defaults
`

// tests/fixtures/search_options.go

// ValidSearchOptions provides valid search options for testing
var ValidSearchOptions = internal.SearchOptions{
	Query:         "Matrix",
	Page:          1,
	Language:      "en",
	Year:          1999,
	MinRating:     7.0,
	MaxRating:     9.0,
	IncludeGenres: []int{28, 878}, // Action, Sci-Fi
	ExcludeGenres: []int{27},      // Horror
	SortBy:        "popularity",
	SortOrder:     "desc",
	MaxItems:      20,
}

// InvalidSearchOptions for testing validation
var InvalidSearchOptions = []struct {
	Name    string
	Options internal.SearchOptions
	Error   string
}{
	{
		Name: "negative max items",
		Options: internal.SearchOptions{
			MaxItems: -5,
		},
		Error: "max-items must be between 1 and 1000",
	},
	{
		Name: "too many max items",
		Options: internal.SearchOptions{
			MaxItems: 1500,
		},
		Error: "max-items must be between 1 and 1000",
	},
	{
		Name: "invalid min rating",
		Options: internal.SearchOptions{
			MinRating: -1.0,
			MaxItems:  20,
		},
		Error: "min-rating must be between 0 and 10",
	},
	{
		Name: "invalid max rating",
		Options: internal.SearchOptions{
			MaxRating: 15.0,
			MaxItems:  20,
		},
		Error: "max-rating must be between 0 and 10",
	},
	{
		Name: "min rating greater than max rating",
		Options: internal.SearchOptions{
			MinRating: 8.0,
			MaxRating: 6.0,
			MaxItems:  20,
		},
		Error: "min-rating (8.0) cannot be greater than max-rating (6.0)",
	},
}

// tests/fixtures/format_options.go

// FormatOptionsTestCases provides various format options for testing
var FormatOptionsTestCases = []struct {
	Name     string
	Options  internal.FormatOptions
	Expected string // Expected substring in output
}{
	{
		Name: "table format with headers",
		Options: internal.FormatOptions{
			Format:   "table",
			NoHeader: false,
			MaxWidth: 120,
		},
		Expected: "#",
	},
	{
		Name: "table format without headers",
		Options: internal.FormatOptions{
			Format:   "table",
			NoHeader: true,
			MaxWidth: 120,
		},
		Expected: "Fight Club",
	},
	{
		Name: "json format",
		Options: internal.FormatOptions{
			Format: "json",
		},
		Expected: `"movies"`,
	},
	{
		Name: "csv format",
		Options: internal.FormatOptions{
			Format: "csv",
		},
		Expected: "ID,Title",
	},
	{
		Name: "original titles",
		Options: internal.FormatOptions{
			Format:           "table",
			UseOriginalTitle: true,
		},
		Expected: "Fight Club",
	},
}

// tests/fixtures/api_responses.go

// APIResponseTestCases provides various API response scenarios
var APIResponseTestCases = []struct {
	Name           string
	StatusCode     int
	ResponseBody   string
	ExpectedError  string
	ExpectedMovies int
}{
	{
		Name:           "successful response",
		StatusCode:     200,
		ResponseBody:   ValidMovieListResponseBody,
		ExpectedMovies: 1,
	},
	{
		Name:          "unauthorized",
		StatusCode:    401,
		ResponseBody:  UnauthorizedResponse,
		ExpectedError: "invalid TMDB API key",
	},
	{
		Name:          "not found",
		StatusCode:    404,
		ResponseBody:  NotFoundResponse,
		ExpectedError: "resource not found",
	},
	{
		Name:          "rate limited",
		StatusCode:    429,
		ResponseBody:  RateLimitResponse,
		ExpectedError: "rate limit exceeded",
	},
	{
		Name:          "internal server error",
		StatusCode:    500,
		ResponseBody:  InternalServerErrorResponse,
		ExpectedError: "API error (status 500)",
	},
	{
		Name:          "invalid json",
		StatusCode:    200,
		ResponseBody:  InvalidJSONResponseBody,
		ExpectedError: "parse response",
	},
	{
		Name:          "malformed data",
		StatusCode:    200,
		ResponseBody:  MalformedDataResponseBody,
		ExpectedError: "json: cannot unmarshal",
	},
}

// tests/fixtures/command_test_cases.go

// CommandTestCases provides test cases for CLI commands
var CommandTestCases = []struct {
	Name         string
	Args         []string
	ExpectedCode int
	Contains     []string
	NotContains  []string
}{
	{
		Name:         "popular movies default",
		Args:         []string{"popular"},
		ExpectedCode: 0,
		Contains:     []string{"#", "Title", "Year", "Rating"},
	},
	{
		Name:         "popular movies with count",
		Args:         []string{"popular", "5"},
		ExpectedCode: 0,
		Contains:     []string{"Showing 5 popular movies"},
	},
	{
		Name:         "search with query",
		Args:         []string{"search", "Matrix"},
		ExpectedCode: 0,
		Contains:     []string{"Found", "movies"},
	},
	{
		Name:         "search with quoted query",
		Args:         []string{"search", "\"The Matrix\""},
		ExpectedCode: 0,
		Contains:     []string{"Matrix"},
	},
	{
		Name:         "json output",
		Args:         []string{"popular", "2", "--format", "json"},
		ExpectedCode: 0,
		Contains:     []string{`"movies"`, `"count"`},
		NotContains:  []string{"#", "Title"},
	},
	{
		Name:         "csv output",
		Args:         []string{"top-rated", "3", "--format", "csv"},
		ExpectedCode: 0,
		Contains:     []string{"ID,Title,Original Title"},
	},
	{
		Name:         "original titles",
		Args:         []string{"popular", "5", "--original-title"},
		ExpectedCode: 0,
		Contains:     []string{"original language titles"},
	},
	{
		Name:         "discover with genres",
		Args:         []string{"discover", "--genres", "action,sci-fi", "--max-items", "5"},
		ExpectedCode: 0,
		Contains:     []string{"Discovered", "movies"},
	},
	{
		Name:         "invalid command",
		Args:         []string{"invalid-command"},
		ExpectedCode: 1,
		Contains:     []string{"Unknown command"},
	},
	{
		Name:         "search without query",
		Args:         []string{"search"},
		ExpectedCode: 1,
		Contains:     []string{"search requires a query"},
	},
	{
		Name:         "invalid max items",
		Args:         []string{"popular", "2000"},
		ExpectedCode: 1,
		Contains:     []string{"max-items must be between 1 and 1000"},
	},
	{
		Name:         "help command",
		Args:         []string{"help"},
		ExpectedCode: 0,
		Contains:     []string{"USAGE", "COMMANDS", "OPTIONS"},
	},
	{
		Name:         "version command",
		Args:         []string{"version"},
		ExpectedCode: 0,
		Contains:     []string{"TMDB CLI", "Build Date"},
	},
}

// tests/fixtures/performance_data.go

// PerformanceTestData provides data for performance/benchmark testing
var (
	SmallMovieList  = ConvertedMovies[:1]
	MediumMovieList = append(ConvertedMovies, ConvertedMovies...) // 6 movies
	LargeMovieList  = make([]internal.Movie, 100)                 // Will be filled by test
)

// BenchmarkTestCases for performance testing
var BenchmarkTestCases = []struct {
	Name   string
	Movies []internal.Movie
	Size   string
}{
	{
		Name:   "format small list",
		Movies: SmallMovieList,
		Size:   "small",
	},
	{
		Name:   "format medium list",
		Movies: MediumMovieList,
		Size:   "medium",
	},
	{
		Name:   "format large list",
		Movies: LargeMovieList,
		Size:   "large",
	},
}

// tests/fixtures/error_scenarios.go

// ErrorScenarios provides comprehensive error testing scenarios
var ErrorScenarios = []struct {
	Name        string
	SetupError  string
	Command     []string
	ExpectedErr string
}{
	{
		Name:        "network timeout",
		SetupError:  "timeout",
		Command:     []string{"popular"},
		ExpectedErr: "timeout",
	},
	{
		Name:        "connection refused",
		SetupError:  "connection_refused",
		Command:     []string{"search", "Matrix"},
		ExpectedErr: "connection refused",
	},
	{
		Name:        "dns resolution failure",
		SetupError:  "dns_failure",
		Command:     []string{"top-rated"},
		ExpectedErr: "no such host",
	},
	{
		Name:        "invalid API response",
		SetupError:  "invalid_response",
		Command:     []string{"now-playing"},
		ExpectedErr: "parse response",
	},
	{
		Name:        "empty response body",
		SetupError:  "empty_response",
		Command:     []string{"upcoming"},
		ExpectedErr: "unexpected end of JSON input",
	},
}

// tests/fixtures/pagination_data.go

// PaginationTestData provides data for testing pagination scenarios
var FirstPageResponse = internal.TMDBResponse{
	Page:         1,
	TotalPages:   3,
	TotalResults: 50,
	Results:      PopularMoviesResponse.Results[:2],
}

var SecondPageResponse = internal.TMDBResponse{
	Page:         2,
	TotalPages:   3,
	TotalResults: 50,
	Results:      TopRatedMoviesResponse.Results[:2],
}

var LastPageResponse = internal.TMDBResponse{
	Page:         3,
	TotalPages:   3,
	TotalResults: 50,
	Results:      NowPlayingMoviesResponse.Results[:1],
}

// Cache test data
var CacheTestScenarios = []struct {
	Name        string
	CacheTTL    time.Duration
	WaitTime    time.Duration
	ShouldHit   bool
	Description string
}{
	{
		Name:        "cache hit within TTL",
		CacheTTL:    5 * time.Second,
		WaitTime:    1 * time.Second,
		ShouldHit:   true,
		Description: "Second request should use cache",
	},
	{
		Name:        "cache miss after TTL",
		CacheTTL:    1 * time.Second,
		WaitTime:    2 * time.Second,
		ShouldHit:   false,
		Description: "Cache should expire, make new request",
	},
	{
		Name:        "cache disabled",
		CacheTTL:    0,
		WaitTime:    0,
		ShouldHit:   false,
		Description: "No caching when TTL is 0",
	},
}

// Environment variable test data
var EnvVarTestCases = []struct {
	Name     string
	EnvVars  map[string]string
	Expected internal.Config
}{
	{
		Name: "all environment variables set",
		EnvVars: map[string]string{
			"TMDB_API_KEY":     "env-api-key",
			"TMDB_BASE_URL":    "https://env-api.example.com",
			"TMDB_TIMEOUT":     "60s",
			"TMDB_MAX_RETRIES": "5",
			"TMDB_CACHE_TTL":   "15m",
			"TMDB_LOG_LEVEL":   "debug",
			"TMDB_FORMAT":      "json",
		},
		Expected: internal.Config{
			APIKey:     "env-api-key",
			BaseURL:    "https://env-api.example.com",
			Timeout:    60 * time.Second,
			MaxRetries: 5,
			CacheTTL:   15 * time.Minute,
			LogLevel:   "debug",
			Format:     "json",
		},
	},
	{
		Name: "partial environment variables",
		EnvVars: map[string]string{
			"TMDB_API_KEY":   "partial-key",
			"TMDB_LOG_LEVEL": "warn",
		},
		Expected: internal.Config{
			APIKey:   "partial-key",
			LogLevel: "warn",
			// Other fields should use defaults
		},
	},
}

// Genre mapping test data
var GenreTestCases = []struct {
	Name        string
	Input       []string
	ExpectedIDs []int
	ExpectedErr string
}{
	{
		Name:        "valid single genre",
		Input:       []string{"action"},
		ExpectedIDs: []int{28},
	},
	{
		Name:        "valid multiple genres",
		Input:       []string{"action", "comedy", "drama"},
		ExpectedIDs: []int{28, 35, 18},
	},
	{
		Name:        "sci-fi alias",
		Input:       []string{"sci-fi"},
		ExpectedIDs: []int{878},
	},
	{
		Name:        "case insensitive",
		Input:       []string{"ACTION", "Comedy", "dRaMa"},
		ExpectedIDs: []int{28, 35, 18},
	},
	{
		Name:        "invalid genre",
		Input:       []string{"invalid-genre"},
		ExpectedErr: "unknown genre: invalid-genre",
	},
	{
		Name:        "mixed valid and invalid",
		Input:       []string{"action", "invalid"},
		ExpectedErr: "unknown genre: invalid",
	},
}
