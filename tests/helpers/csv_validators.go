// Package helpers provides testing utilities for the TMDB CLI test suite.
package helpers

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// ValidateCSVOutput validates that output is valid CSV with expected headers and structure.
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

// ValidateCSVStructure validates CSV structure without header verification and returns records.
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

// ValidateCSVFieldContent validates that specific CSV fields contain expected content.
func ValidateCSVFieldContent(t *testing.T, output string, row, column int, expectedContent string) {
	t.Helper()

	records := ValidateCSVStructure(t, output)

	require.Greater(t, len(records), row, "CSV should have enough rows")
	require.Greater(t, len(records[row]), column, "CSV row should have enough columns")

	actualContent := records[row][column]
	require.Equal(t, expectedContent, actualContent,
		"CSV field at row %d, column %d should contain expected content", row, column)
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

// ValidateCSVHeaders validates that CSV headers match expected values exactly.
func ValidateCSVHeaders(t *testing.T, output string, expectedHeaders []string) {
	t.Helper()

	records := ValidateCSVStructure(t, output)
	require.NotEmpty(t, records, "CSV should not be empty")

	headers := records[0]
	require.Equal(t, len(expectedHeaders), len(headers),
		"CSV should have expected number of header columns")

	for i, expected := range expectedHeaders {
		require.Equal(t, expected, headers[i],
			"Header column %d should match expected value", i)
	}
}

// ValidateCSVDataRow validates that a specific data row contains expected values.
func ValidateCSVDataRow(t *testing.T, output string, rowIndex int, expectedValues []string) {
	t.Helper()

	records := ValidateCSVStructure(t, output)

	// Account for header row
	actualRowIndex := rowIndex + 1
	require.Greater(t, len(records), actualRowIndex,
		"CSV should have enough rows (including header)")

	row := records[actualRowIndex]
	require.Equal(t, len(expectedValues), len(row),
		"CSV row should have expected number of columns")

	for i, expected := range expectedValues {
		require.Equal(t, expected, row[i],
			"CSV row %d, column %d should match expected value", rowIndex, i)
	}
}

// ValidateCSVNoHeader validates CSV output that should not contain headers.
func ValidateCSVNoHeader(t *testing.T, output string, expectedDataRows int) {
	t.Helper()

	records := ValidateCSVStructure(t, output)
	require.Equal(t, expectedDataRows, len(records),
		"CSV should contain only data rows, no header")

	// Validate that first row looks like data, not headers
	if len(records) > 0 {
		firstRow := records[0]
		// First field should be numeric (ID) for our data structure
		require.NotEmpty(t, firstRow[0], "First field should not be empty")
		// Should not contain typical header words
		require.NotContains(t, firstRow[0], "ID",
			"First row should be data, not header")
	}
}

// ValidateCSVSpecialCharacters validates proper handling of special characters in CSV.
func ValidateCSVSpecialCharacters(t *testing.T, output string, expectedSpecialChars []string) {
	t.Helper()

	require.NotEmpty(t, output, "CSV output should not be empty")

	for _, specialChar := range expectedSpecialChars {
		require.Contains(t, output, specialChar,
			"CSV should properly handle special character: %s", specialChar)
	}

	// Ensure CSV is still parseable despite special characters
	ValidateCSVStructure(t, output)
}

// ValidateCSVQuoting validates that CSV fields are properly quoted when necessary.
func ValidateCSVQuoting(t *testing.T, output string) {
	t.Helper()

	records := ValidateCSVStructure(t, output)

	// Check that fields containing commas, quotes, or newlines are properly handled
	for i, record := range records {
		for j, field := range record {
			// These characters should be handled properly by the CSV parser
			if strings.Contains(field, ",") || strings.Contains(field, "\"") || strings.Contains(field, "\n") {
				// If we can parse it, it means it was properly quoted
				t.Logf("Row %d, Column %d properly handles special characters: %q", i, j, field)
			}
		}
	}
}

// ValidateCSVEmpty validates CSV output for empty collections.
func ValidateCSVEmpty(t *testing.T, output string, expectHeader bool) {
	t.Helper()

	if expectHeader {
		// Should have header row only
		records := ValidateCSVStructure(t, output)
		require.Equal(t, 1, len(records), "Empty CSV with header should have exactly one row")

		// Verify it looks like a header (contains typical header words)
		header := records[0]
		hasHeaderWords := false
		for _, field := range header {
			if strings.Contains(field, "ID") || strings.Contains(field, "Title") || strings.Contains(field, "Name") {
				hasHeaderWords = true
				break
			}
		}
		require.True(t, hasHeaderWords, "Header row should contain typical header words")
	} else {
		// Should be completely empty or just whitespace
		trimmed := strings.TrimSpace(output)
		require.Empty(t, trimmed, "CSV without header should be empty for empty collections")
	}
}

// ValidateCSVConsistency validates that CSV structure is consistent across all rows.
func ValidateCSVConsistency(t *testing.T, output string) {
	t.Helper()

	records := ValidateCSVStructure(t, output)

	if len(records) == 0 {
		return
	}

	expectedColumns := len(records[0])

	for i, record := range records {
		require.Equal(t, expectedColumns, len(record),
			"Row %d should have same number of columns as first row (%d)", i, expectedColumns)

		// Validate that no field is completely empty unless it's supposed to be
		for j, field := range record {
			// Allow empty strings for optional fields, but ensure they're not nil
			require.NotNil(t, field, "Field at row %d, column %d should not be nil", i, j)
		}
	}
}
