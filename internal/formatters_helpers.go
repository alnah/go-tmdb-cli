// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"fmt"
	"io"
	"strconv"
)

// formatTitle formats a title with original title handling.
func formatTitle(title, originalTitle string, useOriginal bool) string {
	var displayTitle string
	if useOriginal && originalTitle != "" && originalTitle != title {
		displayTitle = originalTitle
		// Add English title/name in parentheses if space allows
		fullTitle := displayTitle + " (" + title + ")"
		if len(fullTitle) <= TitleColumnWidth {
			displayTitle = fullTitle
		}
	} else {
		displayTitle = title
		// Add original title/name in parentheses if different and space allows
		if originalTitle != "" && originalTitle != title {
			fullTitle := displayTitle + " (" + originalTitle + ")"
			if len(fullTitle) <= TitleColumnWidth {
				displayTitle = fullTitle
			}
		}
	}

	// Truncate if too long
	if len(displayTitle) > TitleColumnWidth {
		displayTitle = displayTitle[:TitleColumnWidth-3] + "..."
	}

	return displayTitle
}

// formatYear formats year with N/A handling.
func formatYear(year int) string {
	if year == 0 {
		return DefaultNA
	}
	return strconv.Itoa(year)
}

// FormatSummary formats summary information for different commands.
func FormatSummary(writer io.Writer, movies []Movie, command string, useOriginal bool) {
	if len(movies) == 0 {
		return
	}

	titleType := "titles"
	if useOriginal {
		titleType = "original language titles"
	}

	switch command {
	case "popular":
		_, _ = fmt.Fprintf(writer, "Showing %d popular movies (%s)\n\n", len(movies), titleType)
	case "top-rated":
		_, _ = fmt.Fprintf(writer, "Showing %d top-rated movies (%s)\n\n", len(movies), titleType)
	case "now-playing":
		_, _ = fmt.Fprintf(writer, "Showing %d movies now playing (%s)\n\n", len(movies), titleType)
	case "upcoming":
		_, _ = fmt.Fprintf(writer, "Showing %d upcoming movies (%s)\n\n", len(movies), titleType)
	case "search":
		_, _ = fmt.Fprintf(writer, "Found %d movies (%s)\n\n", len(movies), titleType)
	case "discover":
		_, _ = fmt.Fprintf(writer, "Discovered %d movies (%s)\n\n", len(movies), titleType)
	}
}

// FormatTVSummary formats summary information for different TV commands.
func FormatTVSummary(writer io.Writer, shows []TVShow, command string, useOriginal bool) {
	if len(shows) == 0 {
		return
	}

	titleType := "names"
	if useOriginal {
		titleType = "original language names"
	}

	switch command {
	case "popular":
		_, _ = fmt.Fprintf(writer, "Showing %d popular TV shows (%s)\n\n", len(shows), titleType)
	case "top-rated":
		_, _ = fmt.Fprintf(writer, "Showing %d top-rated TV shows (%s)\n\n", len(shows), titleType)
	case "on-the-air":
		_, _ = fmt.Fprintf(writer, "Showing %d TV shows on the air (%s)\n\n", len(shows), titleType)
	case "search":
		_, _ = fmt.Fprintf(writer, "Found %d TV shows (%s)\n\n", len(shows), titleType)
	case "discover":
		_, _ = fmt.Fprintf(writer, "Discovered %d TV shows (%s)\n\n", len(shows), titleType)
	}
}
