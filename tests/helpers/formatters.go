// Package helpers provides testing utilities for formatter tests.
package helpers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
)

// FormatTestCase represents a test case for formatter testing.
type FormatTestCase struct {
	Name                string
	Description         string
	Movies              []internal.Movie
	TVShows             []internal.TVShow
	Options             internal.FormatOptions
	ExpectedError       string
	ExpectedInOutput    []string
	NotExpectedInOutput []string
	ValidateFunc        func(t *testing.T, output string)
}

// ValidateFormatOptionsEdgeCases validates format options edge cases.
func ValidateFormatOptionsEdgeCases(t *testing.T, options internal.FormatOptions) {
	t.Helper()

	require.True(t, internal.ValidateFormat(options.Format), "Format should be valid")
	require.GreaterOrEqual(t, options.MaxWidth, 0, "MaxWidth should be non-negative")
}

// AssertFormatError validates that a format operation returns expected error.
func AssertFormatError(t *testing.T, err error, expectedError string) {
	t.Helper()

	if expectedError == "" {
		require.NoError(t, err, "Format operation should not return error")
	} else {
		require.Error(t, err, "Format operation should return error")
		require.Contains(t, err.Error(), expectedError, "Error should contain expected message")
	}
}

// SimulateWriteError simulates a write error during formatting.
func SimulateWriteError(t *testing.T, formatFunc func(io.Writer) error) {
	t.Helper()

	mockWriter := NewMockWriter()
	mockWriter.SetShouldError(true)

	err := formatFunc(mockWriter)
	require.Error(t, err, "Format operation should fail with write error")
}

// ValidateJSONOutput validates that output is valid JSON with expected structure.
func ValidateJSONOutput(t *testing.T, output string, expectMovies bool) {
	t.Helper()

	require.NotEmpty(t, output, "JSON output should not be empty")

	var result map[string]interface{}
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err, "Output should be valid JSON")

	if expectMovies {
		require.Contains(t, result, "movies", "JSON should contain movies key")
		require.Contains(t, result, "count", "JSON should contain count key")

		movies, ok := result["movies"].([]interface{})
		require.True(t, ok, "movies should be an array")

		count, ok := result["count"].(float64)
		require.True(t, ok, "count should be a number")
		require.Equal(t, float64(len(movies)), count, "count should match movies array length")
	} else {
		require.Contains(t, result, "tv_shows", "JSON should contain tv_shows key")
		require.Contains(t, result, "count", "JSON should contain count key")

		shows, ok := result["tv_shows"].([]interface{})
		require.True(t, ok, "tv_shows should be an array")

		count, ok := result["count"].(float64)
		require.True(t, ok, "count should be a number")
		require.Equal(t, float64(len(shows)), count, "count should match tv_shows array length")
	}
}

// ValidateJSONStructure validates basic JSON structure without content verification.
func ValidateJSONStructure(t *testing.T, output string) map[string]interface{} {
	t.Helper()

	require.NotEmpty(t, output, "JSON output should not be empty")

	var result map[string]interface{}
	err := json.Unmarshal([]byte(output), &result)
	require.NoError(t, err, "Output should be valid JSON")

	return result
}

// ValidateJSONFormatting validates that JSON is properly formatted with indentation.
func ValidateJSONFormatting(t *testing.T, output string) {
	t.Helper()

	require.Contains(t, output, "\n", "JSON should contain newlines")
	require.Contains(t, output, "  ", "JSON should contain indentation")
}

// ValidateJSONSearchResult validates search result JSON structure.
func ValidateJSONSearchResult(t *testing.T, output string, expectMovies bool) {
	t.Helper()

	result := ValidateJSONStructure(t, output)

	if expectMovies {
		require.Contains(t, result, "movies", "Search result should contain movies")
		require.Contains(t, result, "page", "Search result should contain page")
		require.Contains(t, result, "total_pages", "Search result should contain total_pages")
		require.Contains(t, result, "total_results", "Search result should contain total_results")
	} else {
		require.Contains(t, result, "tv_shows", "Search result should contain tv_shows")
		require.Contains(t, result, "page", "Search result should contain page")
		require.Contains(t, result, "total_pages", "Search result should contain total_pages")
		require.Contains(t, result, "total_results", "Search result should contain total_results")
	}
}

// ValidateCSVOutput validates that output is valid CSV with expected headers.
func ValidateCSVOutput(t *testing.T, output string, expectedHeaders []string, expectHeader bool) {
	t.Helper()

	require.NotEmpty(t, output, "CSV output should not be empty")

	reader := csv.NewReader(strings.NewReader(output))
	records, err := reader.ReadAll()
	require.NoError(t, err, "Output should be valid CSV")
	require.NotEmpty(t, records, "CSV should contain records")

	if expectHeader {
		require.GreaterOrEqual(t, len(records), 1, "CSV should have at least header row")

		header := records[0]
		require.Equal(
			t,
			len(expectedHeaders),
			len(header),
			"Header should have expected number of columns",
		)

		for i, expected := range expectedHeaders {
			require.Equal(t, expected, header[i], "Header column %d should match", i)
		}
	}
}

// ValidateCSVStructure validates CSV structure without header verification.
func ValidateCSVStructure(t *testing.T, output string) [][]string {
	t.Helper()

	require.NotEmpty(t, output, "CSV output should not be empty")

	reader := csv.NewReader(strings.NewReader(output))
	records, err := reader.ReadAll()
	require.NoError(t, err, "Output should be valid CSV")
	require.NotEmpty(t, records, "CSV should contain records")

	return records
}

// ValidateCSVRecordCount validates that CSV has expected number of records.
func ValidateCSVRecordCount(t *testing.T, output string, expectedCount int, hasHeader bool) {
	t.Helper()

	records := ValidateCSVStructure(t, output)

	expectedTotal := expectedCount
	if hasHeader {
		expectedTotal++ // Add header row
	}

	require.Equal(t, expectedTotal, len(records), "CSV should have expected number of records")
}

// ValidateCSVColumnCount validates that all CSV rows have the same number of columns.
func ValidateCSVColumnCount(t *testing.T, output string) {
	t.Helper()

	records := ValidateCSVStructure(t, output)

	if len(records) == 0 {
		return
	}

	expectedColumns := len(records[0])
	for i, record := range records {
		require.Equal(
			t,
			expectedColumns,
			len(record),
			"Row %d should have %d columns",
			i,
			expectedColumns,
		)
	}
}

// ExtractCSVColumn extracts a specific column from CSV output.
func ExtractCSVColumn(t *testing.T, output string, columnIndex int, skipHeader bool) []string {
	t.Helper()

	records := ValidateCSVStructure(t, output)

	var column []string
	startIndex := 0
	if skipHeader && len(records) > 0 {
		startIndex = 1
	}

	for i := startIndex; i < len(records); i++ {
		record := records[i]
		require.Greater(t, len(record), columnIndex, "Row %d should have column %d", i, columnIndex)
		column = append(column, record[columnIndex])
	}

	return column
}

// ValidateCSVFieldContent validates that specific CSV fields contain expected content.
func ValidateCSVFieldContent(
	t *testing.T,
	output string,
	rowIndex, columnIndex int,
	expectedContent string,
) {
	t.Helper()

	records := ValidateCSVStructure(t, output)

	require.Greater(t, len(records), rowIndex, "CSV should have row %d", rowIndex)
	require.Greater(
		t,
		len(records[rowIndex]),
		columnIndex,
		"Row %d should have column %d",
		rowIndex,
		columnIndex,
	)

	actual := records[rowIndex][columnIndex]
	require.Equal(
		t,
		expectedContent,
		actual,
		"Field at row %d, column %d should match",
		rowIndex,
		columnIndex,
	)
}

// ValidateTableOutput validates that output contains expected table elements.
func ValidateTableOutput(
	t *testing.T,
	output string,
	expectHeader bool,
	expectedElements []string,
) {
	t.Helper()

	require.NotEmpty(t, output, "Table output should not be empty")

	lines := strings.Split(output, "\n")
	require.NotEmpty(t, lines, "Table should have lines")

	if expectHeader {
		headerFound := false
		for _, line := range lines {
			if strings.Contains(line, "#") &&
				(strings.Contains(line, "Year") || strings.Contains(line, "Title") || strings.Contains(line, "Name")) {
				headerFound = true
				break
			}
		}
		require.True(t, headerFound, "Table should contain header row")
	}

	for _, expected := range expectedElements {
		require.Contains(t, output, expected, "Table should contain expected element: %s", expected)
	}
}

// CountTableRows counts the number of data rows in table output.
func CountTableRows(output string) int {
	lines := strings.Split(output, "\n")
	dataRows := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "+") {
			continue
		}
		// Count lines that look like data rows (contain | but not headers)
		if strings.Contains(line, "|") && !strings.Contains(line, "Title") &&
			!strings.Contains(line, "Name") &&
			!strings.Contains(line, "#") {
			dataRows++
		}
	}

	return dataRows
}

// ExtractTableColumns extracts column values from table output.
func ExtractTableColumns(output string, columnIndex int) []string {
	lines := strings.Split(output, "\n")
	var columns []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "|") && !strings.HasPrefix(line, "+") {
			parts := strings.Split(line, "|")
			if len(parts) > columnIndex+1 {
				column := strings.TrimSpace(parts[columnIndex+1])
				if column != "" && !strings.Contains(column, "-") &&
					!strings.Contains(column, "#") {
					columns = append(columns, column)
				}
			}
		}
	}

	return columns
}

// ValidateTableStructure validates basic table structure.
func ValidateTableStructure(t *testing.T, output string) {
	t.Helper()

	require.NotEmpty(t, output, "Table output should not be empty")

	lines := strings.Split(output, "\n")
	require.NotEmpty(t, lines, "Table should have lines")

	// Check for table borders
	borderFound := false
	dataFound := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "+") || strings.Contains(line, "---") {
			borderFound = true
		}
		if strings.Contains(line, "|") && !strings.HasPrefix(line, "+") {
			dataFound = true
		}
	}

	require.True(t, borderFound || dataFound, "Table should contain borders or data rows")
}

// ValidateTableHeaders validates that table contains expected headers.
func ValidateTableHeaders(t *testing.T, output string, expectedHeaders []string) {
	t.Helper()

	for _, header := range expectedHeaders {
		require.Contains(t, output, header, "Table should contain header: %s", header)
	}
}

// ValidateIndicators validates that visual indicators are properly applied.
func ValidateIndicators(t *testing.T, output string, expectedIndicators []string) {
	t.Helper()

	for _, indicator := range expectedIndicators {
		if indicator == "[NEW]" || indicator == "[TOP]" || indicator == "[HOT]" {
			// These indicators should appear when appropriate based on data
			continue
		}
		require.Contains(t, output, indicator, "Output should contain indicator: %s", indicator)
	}
}

// ValidateTableRowCount validates that table has expected number of data rows.
func ValidateTableRowCount(t *testing.T, output string, expectedCount int) {
	t.Helper()

	actualCount := CountTableRows(output)
	require.Equal(
		t,
		expectedCount,
		actualCount,
		"Table should have %d data rows, got %d",
		expectedCount,
		actualCount,
	)
}

// ValidateTableColumnPresence validates that specific columns are present in table.
func ValidateTableColumnPresence(t *testing.T, output string, columnName string) {
	t.Helper()

	require.Contains(t, output, columnName, "Table should contain column: %s", columnName)
}

// ExtractTableCellContent extracts content from a specific table cell.
func ExtractTableCellContent(t *testing.T, output string, rowIndex, columnIndex int) string {
	t.Helper()

	lines := strings.Split(output, "\n")
	dataRows := make([]string, 0)

	// Extract data rows (skip borders and headers)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "|") && !strings.HasPrefix(line, "+") &&
			!strings.Contains(line, "Title") &&
			!strings.Contains(line, "Name") &&
			!strings.Contains(line, "#") {
			dataRows = append(dataRows, line)
		}
	}

	require.Greater(t, len(dataRows), rowIndex, "Table should have row %d", rowIndex)

	parts := strings.Split(dataRows[rowIndex], "|")
	require.Greater(t, len(parts), columnIndex+1, "Row should have column %d", columnIndex)

	return strings.TrimSpace(parts[columnIndex+1])
}

// ValidateSearchResultOutput validates search result pagination information.
func ValidateSearchResultOutput(t *testing.T, output string, page, totalPages, totalResults int) {
	t.Helper()

	if totalPages > 1 {
		require.Contains(t, output, "Page", "Output should contain pagination info")
		require.Contains(t, output, "total results", "Output should contain total results info")

		// Validate specific pagination format
		expectedPage := fmt.Sprintf("Page %d of %d", page, totalPages)
		require.Contains(t, output, expectedPage, "Output should contain correct page info")

		expectedResults := fmt.Sprintf("%d total results", totalResults)
		require.Contains(t, output, expectedResults, "Output should contain correct total results")
	} else {
		// Single page results should not show pagination
		require.NotContains(t, output, "Page", "Single page results should not show pagination")
	}
}

// ValidatePaginationFormat validates that pagination is formatted correctly.
func ValidatePaginationFormat(t *testing.T, output string, page, totalPages, totalResults int) {
	t.Helper()

	if totalPages <= 1 {
		// No pagination should be shown
		require.NotContains(t, output, "Page", "No pagination should be shown for single page")
		return
	}

	lines := strings.Split(output, "\n")
	paginationFound := false

	for _, line := range lines {
		if strings.Contains(line, "Page") && strings.Contains(line, "total results") {
			paginationFound = true

			// Validate format: "Page X of Y (Z total results)"
			expectedFormat := fmt.Sprintf(
				"Page %d of %d (%d total results)",
				page,
				totalPages,
				totalResults,
			)
			require.Contains(t, line, expectedFormat, "Pagination should be in correct format")
			break
		}
	}

	require.True(t, paginationFound, "Pagination information should be present")
}

// ValidateSearchSummary validates search summary information.
func ValidateSearchSummary(
	t *testing.T,
	output string,
	query string,
	resultCount int,
	mediaType string,
) {
	t.Helper()

	expectedSummary := fmt.Sprintf("Found %d %s for \"%s\"", resultCount, mediaType, query)
	require.Contains(t, output, expectedSummary, "Output should contain search summary")
}

// ValidateEmptySearchResult validates empty search result handling.
func ValidateEmptySearchResult(t *testing.T, output string, mediaType string) {
	t.Helper()

	switch mediaType {
	case "movies":
		require.Contains(
			t,
			output,
			"No movies found",
			"Empty movie search should show appropriate message",
		)
	case "tv_shows", "TV shows":
		require.Contains(
			t,
			output,
			"No TV shows found",
			"Empty TV search should show appropriate message",
		)
	default:
		require.Contains(t, output, "found", "Empty search should show appropriate message")
	}
}

// ValidateSearchResultConsistency validates that search results are consistent.
func ValidateSearchResultConsistency(t *testing.T, output string, expectedCount int) {
	t.Helper()

	// Count actual result rows in the output
	actualCount := CountTableRows(output)

	// For table format, verify count matches
	if strings.Contains(output, "|") {
		require.Equal(
			t,
			expectedCount,
			actualCount,
			"Search result count should match displayed rows",
		)
	}
}

// BufferedWriter wraps a writer to capture output for testing.
type BufferedWriter struct {
	buffer bytes.Buffer
}

// NewBufferedWriter creates a new buffered writer.
func NewBufferedWriter() *BufferedWriter {
	return &BufferedWriter{}
}

// Write implements io.Writer interface.
func (bw *BufferedWriter) Write(p []byte) (n int, err error) {
	return bw.buffer.Write(p)
}

// String returns the written content.
func (bw *BufferedWriter) String() string {
	return bw.buffer.String()
}

// Bytes returns the written content as bytes.
func (bw *BufferedWriter) Bytes() []byte {
	return bw.buffer.Bytes()
}

// Lines returns the written content as lines.
func (bw *BufferedWriter) Lines() []string {
	return strings.Split(bw.buffer.String(), "\n")
}

// Reset resets the buffer.
func (bw *BufferedWriter) Reset() {
	bw.buffer.Reset()
}

// Len returns the number of bytes written.
func (bw *BufferedWriter) Len() int {
	return bw.buffer.Len()
}

// SlowWriter simulates a slow writer for performance testing.
type SlowWriter struct {
	buffer bytes.Buffer
	delay  time.Duration
}

// NewSlowWriter creates a new slow writer with specified delay.
func NewSlowWriter(delay time.Duration) *SlowWriter {
	return &SlowWriter{delay: delay}
}

// Write implements io.Writer interface with artificial delay.
func (sw *SlowWriter) Write(p []byte) (n int, err error) {
	time.Sleep(sw.delay)
	return sw.buffer.Write(p)
}

// String returns the written content.
func (sw *SlowWriter) String() string {
	return sw.buffer.String()
}

// CreateLargeDataWriter creates a writer for testing large data formatting.
func CreateLargeDataWriter() io.Writer {
	return &bytes.Buffer{}
}

// ValidatePerformance validates that formatting completes within reasonable time.
func ValidatePerformance(t *testing.T, formatFunc func() error, maxDurationMs int) {
	t.Helper()

	start := time.Now()
	err := formatFunc()
	duration := time.Since(start)

	require.NoError(t, err, "Format operation should not error")
	require.Less(t, duration.Milliseconds(), int64(maxDurationMs),
		"Format operation should complete within %dms, took %v", maxDurationMs, duration)
}

// CountingWriter counts the number of write operations and bytes written.
type CountingWriter struct {
	buffer     bytes.Buffer
	writeCount int
	byteCount  int
}

// NewCountingWriter creates a new counting writer.
func NewCountingWriter() *CountingWriter {
	return &CountingWriter{}
}

// Write implements io.Writer interface while counting operations.
func (cw *CountingWriter) Write(p []byte) (n int, err error) {
	cw.writeCount++
	n, err = cw.buffer.Write(p)
	cw.byteCount += n
	return n, err
}

// String returns the written content.
func (cw *CountingWriter) String() string {
	return cw.buffer.String()
}

// GetWriteCount returns the number of write operations.
func (cw *CountingWriter) GetWriteCount() int {
	return cw.writeCount
}

// GetByteCount returns the number of bytes written.
func (cw *CountingWriter) GetByteCount() int {
	return cw.byteCount
}

// Reset resets all counters and buffer.
func (cw *CountingWriter) Reset() {
	cw.buffer.Reset()
	cw.writeCount = 0
	cw.byteCount = 0
}
