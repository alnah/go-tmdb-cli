// Package internal provides core functionality for the TMDB CLI application.
package internal

import "fmt"

// Progress constants for when to show progress indicators.
const (
	ProgressThreshold = 40
)

// ShowProgress displays progress message for large requests.
func ShowProgress(message string, maxItems int) {
	if maxItems > ProgressThreshold {
		fmt.Printf("%s...\n", message)
	}
}

// ShowFetchingProgress displays progress for fetching operations.
func ShowFetchingProgress(count int, listType string) {
	ShowProgress(fmt.Sprintf("Fetching %d %s movies", count, listType), count)
}

// ShowSearchProgress displays progress for search operations.
func ShowSearchProgress(query string, maxItems int) {
	ShowProgress(fmt.Sprintf("Searching for \"%s\"", query), maxItems)
}

// ShowTVSearchProgress displays progress for TV search operations.
func ShowTVSearchProgress(query string, maxItems int) {
	ShowProgress(fmt.Sprintf("Searching TV shows for \"%s\"", query), maxItems)
}
