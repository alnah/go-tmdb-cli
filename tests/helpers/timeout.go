// Package helpers provides testing utilities for timeout handling.
package helpers

import (
	"testing"
)

// TestTimeout provides a helper for testing functions with timeouts.
type TestTimeout struct {
	t       *testing.T
	timeout chan struct{}
}

// NewTestTimeout creates a new test timeout helper.
func NewTestTimeout(t *testing.T) *TestTimeout {
	t.Helper()

	return &TestTimeout{
		t:       t,
		timeout: make(chan struct{}),
	}
}

// Start starts the timeout mechanism.
func (tt *TestTimeout) Start() {
	go func() {
		<-tt.timeout
	}()
}

// Stop stops the timeout mechanism.
func (tt *TestTimeout) Stop() {
	close(tt.timeout)
}
