// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"encoding/csv"
	"fmt"
	"io"
)

// formatCSVGeneric formats any media type as CSV using MediaFormatter interface.
func formatCSVGeneric(
	writer io.Writer,
	formatters []MediaFormatter,
	headers []string,
	options FormatOptions,
) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header if not disabled
	if !options.NoHeader {
		if err := csvWriter.Write(headers); err != nil {
			return fmt.Errorf("write CSV header: %w", err)
		}
	}

	// Write data rows
	for _, formatter := range formatters {
		if err := csvWriter.Write(formatter.GetFields()); err != nil {
			return fmt.Errorf("write CSV row: %w", err)
		}
	}

	return nil
}
