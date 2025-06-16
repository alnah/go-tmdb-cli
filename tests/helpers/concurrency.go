// Package helpers provides testing utilities for concurrent execution.
package helpers

import (
	"sync"
	"testing"
)

// RunConcurrently executes a function concurrently in multiple goroutines.
// This is useful for testing thread safety and concurrent behavior.
func RunConcurrently(t *testing.T, numGoroutines int, fn func(goroutineID int)) {
	t.Helper()

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			fn(id)
		}(i)
	}

	wg.Wait()
}
