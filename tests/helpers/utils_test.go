// tests/helpers/utils_test.go
package helpers_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/alnah/tmdb-cli/internal"
)

// SetupTestConfig creates a test configuration with sensible defaults
func SetupTestConfig() internal.Config {
	return internal.Config{
		APIKey:     "test-api-key-12345",
		BaseURL:    "https://api.themoviedb.org/3",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		CacheTTL:   5 * time.Minute,
		LogLevel:   "error", // Reduce noise in tests
		Format:     "table",
	}
}

// SetupTestConfigWithAPIKey creates test config with a real API key from environment
func SetupTestConfigWithAPIKey() internal.Config {
	config := SetupTestConfig()
	if apiKey := os.Getenv("TMDB_API_KEY"); apiKey != "" {
		return internal.Config{}
	}
	return config
}

// CreateTempConfigFile creates a temporary configuration file with given content
func CreateTempConfigFile(t *testing.T, content string) string {
	t.Helper()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}

	return configPath
}

// CreateTempConfigFileWithCleanup creates a temp config file and returns cleanup function
func CreateTempConfigFileWithCleanup(content string) (string, func()) {
	tempDir, err := os.MkdirTemp("", "tmdb-test-*")
	if err != nil {
		panic(fmt.Sprintf("Failed to create temp dir: %v", err))
	}

	configPath := filepath.Join(tempDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(content), 0o644)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		panic(fmt.Sprintf("Failed to create temp config file: %v", err))
	}

	cleanup := func() {
		_ = os.RemoveAll(tempDir)
	}

	return configPath, cleanup
}

// SetupTestEnv sets up test environment variables
func SetupTestEnv(t *testing.T, envVars map[string]string) {
	t.Helper()

	// Store original values for cleanup
	originalVars := make(map[string]string)

	for key, value := range envVars {
		originalVars[key] = os.Getenv(key)
		_ = os.Setenv(key, value)
	}

	// Cleanup function
	t.Cleanup(func() {
		for key, originalValue := range originalVars {
			if originalValue == "" {
				_ = os.Unsetenv(key)
			} else {
				_ = os.Setenv(key, originalValue)
			}
		}
	})
}

// CaptureOutput captures stdout/stderr for testing output
func CaptureOutput(t *testing.T, fn func()) (stdout, stderr string) {
	t.Helper()

	// Create pipes
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	stdoutR, stdoutW, _ := os.Pipe()
	stderrR, stderrW, _ := os.Pipe()

	os.Stdout = stdoutW
	os.Stderr = stderrW

	// Capture output
	stdoutCh := make(chan string)
	stderrCh := make(chan string)

	go func() {
		defer close(stdoutCh)
		buf, _ := io.ReadAll(stdoutR)
		stdoutCh <- string(buf)
	}()

	go func() {
		defer close(stderrCh)
		buf, _ := io.ReadAll(stderrR)
		stderrCh <- string(buf)
	}()

	// Execute function
	fn()

	// Restore and close
	os.Stdout = oldStdout
	os.Stderr = oldStderr
	_ = stdoutW.Close()
	_ = stderrW.Close()

	// Get captured output
	stdout = <-stdoutCh
	stderr = <-stderrCh

	return stdout, stderr
}

// SkipIfNoAPIKey skips test if TMDB_API_KEY is not set
func SkipIfNoAPIKey(t *testing.T) {
	t.Helper()

	if os.Getenv("TMDB_API_KEY") == "" {
		t.Skip("TMDB_API_KEY not set, skipping integration test")
	}
}

// SkipIfShort skips test if running with -short flag
func SkipIfShort(t *testing.T, reason string) {
	t.Helper()

	if testing.Short() {
		if reason != "" {
			t.Skipf("skipping test in short mode: %s", reason)
		} else {
			t.Skip("skipping test in short mode")
		}
	}
}

// RequireAPIKey panics if TMDB_API_KEY is not set (for tests that must run)
func RequireAPIKey(t *testing.T) string {
	t.Helper()

	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		t.Fatal("TMDB_API_KEY environment variable is required for this test")
	}
	return apiKey
}

// CreateTestSearchOptions creates test search options with defaults
func CreateTestSearchOptions() internal.SearchOptions {
	return internal.SearchOptions{
		Query:         "test",
		Page:          1,
		Language:      "en",
		Year:          0,
		MinRating:     0,
		MaxRating:     0,
		IncludeGenres: nil,
		ExcludeGenres: nil,
		SortBy:        "popularity",
		SortOrder:     "desc",
		MaxItems:      20,
	}
}

// CreateTestFormatOptions creates test format options with defaults
func CreateTestFormatOptions() internal.FormatOptions {
	return internal.FormatOptions{
		Format:           "table",
		UseOriginalTitle: false,
		NoHeader:         false,
		MaxWidth:         120,
	}
}

// CreateSampleTMDBMovie creates a sample TMDB movie for testing
func CreateSampleTMDBMovie(id int, title string) TMDBMovie {
	return TMDBMovie{
		ID:            id,
		Title:         title,
		OriginalTitle: title,
		Overview:      fmt.Sprintf("Overview for %s", title),
		ReleaseDate:   "2023-06-15",
		VoteAverage:   7.5,
		VoteCount:     1500,
		GenreIDs:      []int{28, 35}, // Action, Comedy
		Popularity:    85.5,
		Adult:         false,
		Video:         false,
	}
}

// CreateSampleMovieList creates a list of sample movies for testing
func CreateSampleMovieList(count int) []internal.Movie {
	movies := make([]internal.Movie, count)
	for i := range count {
		movies[i] = internal.Movie{
			ID:          i + 1,
			Title:       fmt.Sprintf("Test Movie %d", i+1),
			Year:        2020 + (i % 5),
			Rating:      5.0 + float64(i%6),
			Votes:       1000 + i*100,
			Popularity:  float64(50 + i*10),
			Genres:      "Action, Drama",
			Overview:    fmt.Sprintf("Overview for test movie %d", i+1),
			Language:    "en",
			Adult:       false,
			ReleaseDate: fmt.Sprintf("202%d-06-15", i%5),
		}
	}
	return movies
}

// CreateSampleTMDBMovieList creates a list of sample TMDB movies
func CreateSampleTMDBMovieList(count int) []TMDBMovie {
	movies := make([]TMDBMovie, count)
	for i := range count {
		movies[i] = CreateSampleTMDBMovie(i+1, fmt.Sprintf("Test Movie %d", i+1))
	}
	return movies
}

// TimeoutTest runs a test function with a timeout
func TimeoutTest(t *testing.T, timeout time.Duration, testFunc func()) {
	t.Helper()

	done := make(chan bool, 1)

	go func() {
		testFunc()
		done <- true
	}()

	select {
	case <-done:
		// Test completed successfully
	case <-time.After(timeout):
		t.Fatalf("Test timed out after %v", timeout)
	}
}

// RetryTest retries a test function until it passes or max attempts reached
func RetryTest(t *testing.T, maxAttempts int, testFunc func() error) {
	t.Helper()

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := testFunc()
		if err == nil {
			return // Success
		}

		lastErr = err
		if attempt < maxAttempts {
			t.Logf("Attempt %d failed: %v, retrying...", attempt, err)
			time.Sleep(time.Duration(attempt) * time.Second) // Exponential backoff
		}
	}

	t.Fatalf("Test failed after %d attempts. Last error: %v", maxAttempts, lastErr)
}

// CompareStringSlices compares two string slices for equality
func CompareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// ContainsString checks if a string slice contains a specific string
func ContainsString(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

// FilterMoviesByRating filters movies by minimum rating
func FilterMoviesByRating(movies []internal.Movie, minRating float64) []internal.Movie {
	filtered := make([]internal.Movie, 0)
	for _, movie := range movies {
		if movie.Rating >= minRating {
			filtered = append(filtered, movie)
		}
	}
	return filtered
}

// FilterMoviesByYear filters movies by year
func FilterMoviesByYear(movies []internal.Movie, year int) []internal.Movie {
	filtered := make([]internal.Movie, 0)
	for _, movie := range movies {
		if movie.Year == year {
			filtered = append(filtered, movie)
		}
	}
	return filtered
}

// FindMovieByTitle finds a movie by title (case-insensitive)
func FindMovieByTitle(movies []internal.Movie, title string) *internal.Movie {
	lowerTitle := strings.ToLower(title)
	for _, movie := range movies {
		if strings.ToLower(movie.Title) == lowerTitle {
			return &movie
		}
	}
	return nil
}

// ParseTableOutput parses table output and returns rows as string slices
func ParseTableOutput(output string) [][]string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	rows := make([][]string, 0)

	for _, line := range lines {
		// Skip empty lines and decorative lines
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "+") ||
			strings.HasPrefix(line, "|") && strings.Count(line, "|") < 3 {
			continue
		}

		// Parse table row
		if strings.Contains(line, "|") {
			fields := strings.Split(line, "|")
			row := make([]string, 0)
			for _, field := range fields {
				field = strings.TrimSpace(field)
				if field != "" {
					row = append(row, field)
				}
			}
			if len(row) > 0 {
				rows = append(rows, row)
			}
		}
	}

	return rows
}

// GetTestDataDir returns the path to test data directory
func GetTestDataDir(t *testing.T) string {
	t.Helper()

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Look for tests directory
	testDataDir := filepath.Join(cwd, "tests", "fixtures")
	if _, err := os.Stat(testDataDir); os.IsNotExist(err) {
		// Try parent directory
		testDataDir = filepath.Join(filepath.Dir(cwd), "tests", "fixtures")
		if _, err := os.Stat(testDataDir); os.IsNotExist(err) {
			t.Fatalf("Could not find test data directory")
		}
	}

	return testDataDir
}

// LoadTestFixture loads a test fixture file
func LoadTestFixture(t *testing.T, filename string) []byte {
	t.Helper()

	testDataDir := GetTestDataDir(t)
	filePath := filepath.Join(testDataDir, filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to load test fixture %s: %v", filename, err)
	}

	return data
}

// CreateTestLogFile creates a temporary log file for testing
func CreateTestLogFile(t *testing.T) *os.File {
	t.Helper()

	file, err := os.CreateTemp("", "tmdb-test-*.log")
	if err != nil {
		t.Fatalf("Failed to create test log file: %v", err)
	}

	t.Cleanup(func() {
		_ = file.Close()
		_ = os.Remove(file.Name())
	})

	return file
}

// AssertEventuallyTrue retries a condition until it becomes true or timeout
func AssertEventuallyTrue(
	t *testing.T,
	condition func() bool,
	timeout time.Duration,
	message string,
) {
	t.Helper()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeoutCh := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			if condition() {
				return // Success
			}
		case <-timeoutCh:
			t.Fatalf("Condition not met within %v: %s", timeout, message)
		}
	}
}

// GetFreePort finds a free port for testing
func GetFreePort(t *testing.T) int {
	t.Helper()

	// This is a simplified implementation
	// In a real scenario, you might use net.Listen with port 0
	return 8080 + int(t.Name()[0])%100 // Simple hash-based port
}

// Benchmark utilities

// BenchmarkFunc is a helper for creating benchmark functions
type BenchmarkFunc func(b *testing.B, size int)

// RunBenchmarkSizes runs a benchmark with different input sizes
func RunBenchmarkSizes(b *testing.B, sizes []int, benchFunc BenchmarkFunc) {
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			benchFunc(b, size)
		})
	}
}

// MeasureMemory measures memory allocation for a function
func MeasureMemory(b *testing.B, fn func()) {
	b.Helper()

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	fn()

	runtime.GC()
	runtime.ReadMemStats(&m2)

	b.ReportMetric(float64(m2.Alloc-m1.Alloc), "bytes/op")
}
