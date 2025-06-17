// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

import (
	"context"
	"time"
)

// Clock provides an interface for time operations to facilitate testing.
type Clock interface {
	// Now returns the current time.
	Now() time.Time

	// Sleep pauses the current goroutine for at least the duration d.
	Sleep(d time.Duration)

	// After waits for the duration to elapse and then sends the current time on the returned channel.
	After(d time.Duration) <-chan time.Time

	// NewTimer creates a new Timer that will send the current time on its channel after at least duration d.
	NewTimer(d time.Duration) *time.Timer

	// NewTicker returns a new Ticker containing a channel that will send the
	// time with a period specified by the duration d.
	NewTicker(d time.Duration) *time.Ticker

	// WithTimeout returns a context with timeout.
	WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)
}

// RealClock implements Clock using standard time functions.
type RealClock struct{}

// Now returns the current time.
func (RealClock) Now() time.Time {
	return time.Now()
}

// Sleep pauses the current goroutine for at least the duration d.
func (RealClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

// After waits for the duration to elapse and then sends the current time on the returned channel.
func (RealClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// NewTimer creates a new Timer.
func (RealClock) NewTimer(d time.Duration) *time.Timer {
	return time.NewTimer(d)
}

// NewTicker returns a new Ticker.
func (RealClock) NewTicker(d time.Duration) *time.Ticker {
	return time.NewTicker(d)
}

// WithTimeout returns a context with timeout.
func (RealClock) WithTimeout(
	parent context.Context,
	timeout time.Duration,
) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

// MockClock implements Clock for testing purposes.
type MockClock struct {
	CurrentTime time.Time
	SleepCalls  []time.Duration
	Timers      []*mockTimer
}

// mockTimer represents a mock timer for testing.
type mockTimer struct {
	C        chan time.Time
	Duration time.Duration
	Clock    *MockClock
}

// Now returns the mocked current time.
func (m *MockClock) Now() time.Time {
	return m.CurrentTime
}

// Sleep records the sleep call.
func (m *MockClock) Sleep(d time.Duration) {
	m.SleepCalls = append(m.SleepCalls, d)
}

// After returns a channel that receives after the duration.
func (m *MockClock) After(d time.Duration) <-chan time.Time {
	ch := make(chan time.Time, 1)
	ch <- m.CurrentTime.Add(d)
	return ch
}

// NewTimer creates a mock timer.
func (m *MockClock) NewTimer(dur time.Duration) *time.Timer {
	timer := &mockTimer{
		C:        make(chan time.Time, 1),
		Duration: dur,
		Clock:    m,
	}
	m.Timers = append(m.Timers, timer)
	// For simplicity, immediately send on the channel
	timer.C <- m.CurrentTime.Add(dur)
	return &time.Timer{C: timer.C}
}

// NewTicker creates a mock ticker.
func (m *MockClock) NewTicker(d time.Duration) *time.Ticker {
	// Simple implementation for testing
	ticker := time.NewTicker(d)
	ticker.Stop() // Stop immediately to prevent real ticking
	return ticker
}

// WithTimeout returns a context with timeout using the mock clock.
func (m *MockClock) WithTimeout(
	parent context.Context,
	timeout time.Duration,
) (context.Context, context.CancelFunc) {
	// For testing, we can use a real timeout but with a very short duration
	// or implement a more sophisticated mock
	return context.WithTimeout(parent, timeout)
}

// Advance moves the mock clock forward by the specified duration.
func (m *MockClock) Advance(d time.Duration) {
	m.CurrentTime = m.CurrentTime.Add(d)
}

// SetTime sets the mock clock to a specific time.
func (m *MockClock) SetTime(t time.Time) {
	m.CurrentTime = t
}
