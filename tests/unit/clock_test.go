package unit

import (
	"context"
	"testing"
	"time"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRealClock_Now(t *testing.T) {
	t.Parallel()

	t.Run("returns current time", func(t *testing.T) {
		t.Parallel()
		clock := internal.RealClock{}
		before := time.Now()
		now := clock.Now()
		after := time.Now()

		assert.True(t, now.After(before) || now.Equal(before))
		assert.True(t, now.Before(after) || now.Equal(after))
	})
}

func TestRealClock_Sleep(t *testing.T) {
	t.Parallel()

	t.Run("sleeps for specified duration", func(t *testing.T) {
		t.Parallel()
		clock := internal.RealClock{}
		duration := 10 * time.Millisecond

		start := time.Now()
		clock.Sleep(duration)
		elapsed := time.Since(start)

		assert.GreaterOrEqual(t, elapsed, duration)
		assert.Less(t, elapsed, duration+50*time.Millisecond) // Allow some tolerance
	})

	t.Run("sleep with zero duration returns immediately", func(t *testing.T) {
		t.Parallel()
		clock := internal.RealClock{}

		start := time.Now()
		clock.Sleep(0)
		elapsed := time.Since(start)

		assert.Less(t, elapsed, 10*time.Millisecond)
	})
}

func TestRealClock_After(t *testing.T) {
	t.Parallel()

	t.Run("sends time on channel after duration", func(t *testing.T) {
		t.Parallel()
		clock := internal.RealClock{}
		duration := 10 * time.Millisecond

		start := time.Now()
		ch := clock.After(duration)

		select {
		case receivedTime := <-ch:
			elapsed := time.Since(start)
			assert.GreaterOrEqual(t, elapsed, duration)
			assert.Less(t, elapsed, duration+50*time.Millisecond)
			assert.True(t, receivedTime.After(start))
		case <-time.After(100 * time.Millisecond):
			t.Fatal("timeout waiting for After channel")
		}
	})

	t.Run("after with zero duration sends immediately", func(t *testing.T) {
		t.Parallel()
		clock := internal.RealClock{}

		ch := clock.After(0)

		select {
		case <-ch:
			// Expected to receive immediately
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for After channel with zero duration")
		}
	})
}

func TestRealClock_NewTimer(t *testing.T) {
	t.Parallel()

	t.Run("creates timer that fires after duration", func(t *testing.T) {
		t.Parallel()
		clock := internal.RealClock{}
		duration := 10 * time.Millisecond

		start := time.Now()
		timer := clock.NewTimer(duration)
		defer timer.Stop()

		select {
		case receivedTime := <-timer.C:
			elapsed := time.Since(start)
			assert.GreaterOrEqual(t, elapsed, duration)
			assert.Less(t, elapsed, duration+50*time.Millisecond)
			assert.True(t, receivedTime.After(start))
		case <-time.After(100 * time.Millisecond):
			t.Fatal("timeout waiting for timer")
		}
	})

	t.Run("timer can be stopped", func(t *testing.T) {
		t.Parallel()
		clock := internal.RealClock{}
		timer := clock.NewTimer(50 * time.Millisecond)

		// Stop timer immediately
		stopped := timer.Stop()
		assert.True(t, stopped)

		// Timer should not fire
		select {
		case <-timer.C:
			t.Fatal("timer fired after being stopped")
		case <-time.After(100 * time.Millisecond):
			// Expected - timer was stopped
		}
	})

	t.Run("timer with zero duration fires immediately", func(t *testing.T) {
		t.Parallel()
		clock := internal.RealClock{}
		timer := clock.NewTimer(0)
		defer timer.Stop()

		select {
		case <-timer.C:
			// Expected to fire immediately
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for immediate timer")
		}
	})
}

func TestRealClock_NewTicker(t *testing.T) {
	t.Parallel()

	t.Run("creates ticker that ticks at intervals", func(t *testing.T) {
		t.Parallel()
		clock := internal.RealClock{}
		interval := 20 * time.Millisecond

		ticker := clock.NewTicker(interval)
		defer ticker.Stop()

		start := time.Now()
		tickCount := 0
		maxTicks := 3

		for tickCount < maxTicks {
			select {
			case <-ticker.C:
				tickCount++
				elapsed := time.Since(start)
				expectedMin := time.Duration(tickCount) * interval
				expectedMax := expectedMin + 50*time.Millisecond
				assert.GreaterOrEqual(t, elapsed, expectedMin)
				assert.Less(t, elapsed, expectedMax)
			case <-time.After(200 * time.Millisecond):
				t.Fatalf("timeout waiting for tick %d", tickCount+1)
			}
		}
	})

	t.Run("ticker can be stopped", func(t *testing.T) {
		t.Parallel()
		clock := internal.RealClock{}
		ticker := clock.NewTicker(10 * time.Millisecond)

		// Let it tick once
		select {
		case <-ticker.C:
			// Expected first tick
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for first tick")
		}

		// Stop ticker
		ticker.Stop()

		// Should not tick again
		select {
		case <-ticker.C:
			t.Fatal("ticker ticked after being stopped")
		case <-time.After(50 * time.Millisecond):
			// Expected - ticker was stopped
		}
	})
}

func TestRealClock_WithTimeout(t *testing.T) {
	t.Parallel()

	t.Run("creates context with timeout", func(t *testing.T) {
		t.Parallel()
		clock := internal.RealClock{}
		parent := context.Background()
		timeout := 20 * time.Millisecond

		ctx, cancel := clock.WithTimeout(parent, timeout)
		defer cancel()

		start := time.Now()

		select {
		case <-ctx.Done():
			elapsed := time.Since(start)
			assert.GreaterOrEqual(t, elapsed, timeout)
			assert.Less(t, elapsed, timeout+50*time.Millisecond)
			assert.Equal(t, context.DeadlineExceeded, ctx.Err())
		case <-time.After(100 * time.Millisecond):
			t.Fatal("timeout waiting for context timeout")
		}
	})

	t.Run("context can be canceled before timeout", func(t *testing.T) {
		t.Parallel()
		clock := internal.RealClock{}
		parent := context.Background()
		timeout := 100 * time.Millisecond

		ctx, cancel := clock.WithTimeout(parent, timeout)

		// Cancel immediately
		cancel()

		select {
		case <-ctx.Done():
			assert.Equal(t, context.Canceled, ctx.Err())
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for context cancellation")
		}
	})

	t.Run("inherits parent context cancellation", func(t *testing.T) {
		t.Parallel()
		clock := internal.RealClock{}
		parent, parentCancel := context.WithCancel(context.Background())
		timeout := 100 * time.Millisecond

		ctx, cancel := clock.WithTimeout(parent, timeout)
		defer cancel()

		// Cancel parent
		parentCancel()

		select {
		case <-ctx.Done():
			assert.Equal(t, context.Canceled, ctx.Err())
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for parent context cancellation")
		}
	})
}

func TestMockClock_Now(t *testing.T) {
	t.Parallel()

	t.Run("returns set time", func(t *testing.T) {
		t.Parallel()
		expectedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		clock := &internal.MockClock{
			CurrentTime: expectedTime,
		}

		assert.Equal(t, expectedTime, clock.Now())
	})
}

func TestMockClock_Sleep(t *testing.T) {
	t.Parallel()

	t.Run("records sleep calls", func(t *testing.T) {
		t.Parallel()
		clock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		durations := []time.Duration{
			10 * time.Millisecond,
			100 * time.Millisecond,
			1 * time.Second,
		}

		for _, duration := range durations {
			clock.Sleep(duration)
		}

		require.Len(t, clock.SleepCalls, len(durations))
		for i, expected := range durations {
			assert.Equal(t, expected, clock.SleepCalls[i])
		}
	})

	t.Run("multiple sleep calls accumulate", func(t *testing.T) {
		t.Parallel()
		clock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		assert.Empty(t, clock.SleepCalls)

		clock.Sleep(1 * time.Millisecond)
		assert.Len(t, clock.SleepCalls, 1)

		clock.Sleep(2 * time.Millisecond)
		assert.Len(t, clock.SleepCalls, 2)

		assert.Equal(t, 1*time.Millisecond, clock.SleepCalls[0])
		assert.Equal(t, 2*time.Millisecond, clock.SleepCalls[1])
	})
}

func TestMockClock_After(t *testing.T) {
	t.Parallel()

	t.Run("returns channel with future time", func(t *testing.T) {
		t.Parallel()
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		clock := &internal.MockClock{
			CurrentTime: baseTime,
		}

		duration := 5 * time.Minute
		ch := clock.After(duration)

		select {
		case receivedTime := <-ch:
			expectedTime := baseTime.Add(duration)
			assert.Equal(t, expectedTime, receivedTime)
		case <-time.After(10 * time.Millisecond):
			t.Fatal("timeout waiting for After channel")
		}
	})

	t.Run("works with zero duration", func(t *testing.T) {
		t.Parallel()
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		clock := &internal.MockClock{
			CurrentTime: baseTime,
		}

		ch := clock.After(0)

		select {
		case receivedTime := <-ch:
			assert.Equal(t, baseTime, receivedTime)
		case <-time.After(10 * time.Millisecond):
			t.Fatal("timeout waiting for zero duration After channel")
		}
	})
}

func TestMockClock_NewTimer(t *testing.T) {
	t.Parallel()

	t.Run("creates timer that fires with future time", func(t *testing.T) {
		t.Parallel()
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		clock := &internal.MockClock{
			CurrentTime: baseTime,
		}

		duration := 10 * time.Minute
		timer := clock.NewTimer(duration)

		select {
		case receivedTime := <-timer.C:
			expectedTime := baseTime.Add(duration)
			assert.Equal(t, expectedTime, receivedTime)
		case <-time.After(10 * time.Millisecond):
			t.Fatal("timeout waiting for timer")
		}

		// Verify timer was recorded
		require.Len(t, clock.Timers, 1)
		assert.Equal(t, duration, clock.Timers[0].Duration)
		assert.Equal(t, clock, clock.Timers[0].Clock)
	})

	t.Run("multiple timers are tracked", func(t *testing.T) {
		t.Parallel()
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		clock := &internal.MockClock{
			CurrentTime: baseTime,
		}

		durations := []time.Duration{
			1 * time.Minute,
			5 * time.Minute,
			10 * time.Minute,
		}

		for i, duration := range durations {
			timer := clock.NewTimer(duration)

			// Consume the timer channel
			select {
			case <-timer.C:
				// Expected
			case <-time.After(10 * time.Millisecond):
				t.Fatalf("timeout waiting for timer %d", i)
			}
		}

		require.Len(t, clock.Timers, len(durations))
		for i, expected := range durations {
			assert.Equal(t, expected, clock.Timers[i].Duration)
		}
	})
}

func TestMockClock_NewTicker(t *testing.T) {
	t.Parallel()

	t.Run("creates ticker that is immediately stopped", func(t *testing.T) {
		t.Parallel()
		clock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		ticker := clock.NewTicker(1 * time.Second)

		// Mock ticker should not actually tick
		select {
		case <-ticker.C:
			t.Fatal("mock ticker should not tick")
		case <-time.After(10 * time.Millisecond):
			// Expected - mock ticker doesn't actually tick
		}

		ticker.Stop() // Should not panic
	})
}

func TestMockClock_WithTimeout(t *testing.T) {
	t.Parallel()

	t.Run("creates context with timeout", func(t *testing.T) {
		t.Parallel()
		clock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		parent := context.Background()
		timeout := 10 * time.Millisecond

		ctx, cancel := clock.WithTimeout(parent, timeout)
		defer cancel()

		// Should behave like a real timeout context
		select {
		case <-ctx.Done():
			assert.Equal(t, context.DeadlineExceeded, ctx.Err())
		case <-time.After(50 * time.Millisecond):
			t.Fatal("timeout waiting for context timeout")
		}
	})
}

func TestMockClock_Advance(t *testing.T) {
	t.Parallel()

	t.Run("advances time by specified duration", func(t *testing.T) {
		t.Parallel()
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		clock := &internal.MockClock{
			CurrentTime: baseTime,
		}

		duration := 2 * time.Hour
		clock.Advance(duration)

		expectedTime := baseTime.Add(duration)
		assert.Equal(t, expectedTime, clock.Now())
	})

	t.Run("multiple advances accumulate", func(t *testing.T) {
		t.Parallel()
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		clock := &internal.MockClock{
			CurrentTime: baseTime,
		}

		clock.Advance(1 * time.Hour)
		clock.Advance(30 * time.Minute)
		clock.Advance(15 * time.Minute)

		expectedTime := baseTime.Add(1*time.Hour + 30*time.Minute + 15*time.Minute)
		assert.Equal(t, expectedTime, clock.Now())
	})

	t.Run("advance with zero duration does not change time", func(t *testing.T) {
		t.Parallel()
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		clock := &internal.MockClock{
			CurrentTime: baseTime,
		}

		clock.Advance(0)
		assert.Equal(t, baseTime, clock.Now())
	})

	t.Run("advance with negative duration moves backward", func(t *testing.T) {
		t.Parallel()
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		clock := &internal.MockClock{
			CurrentTime: baseTime,
		}

		duration := -1 * time.Hour
		clock.Advance(duration)

		expectedTime := baseTime.Add(duration)
		assert.Equal(t, expectedTime, clock.Now())
	})
}

func TestMockClock_SetTime(t *testing.T) {
	t.Parallel()

	t.Run("sets time to specific value", func(t *testing.T) {
		t.Parallel()
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		newTime := time.Date(2025, 6, 15, 8, 30, 45, 0, time.UTC)

		clock := &internal.MockClock{
			CurrentTime: baseTime,
		}

		clock.SetTime(newTime)
		assert.Equal(t, newTime, clock.Now())
	})

	t.Run("set time overwrites previous time", func(t *testing.T) {
		t.Parallel()
		clock := &internal.MockClock{
			CurrentTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		}

		times := []time.Time{
			time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 7, 4, 15, 30, 0, 0, time.UTC),
			time.Date(2025, 1, 1, 23, 59, 59, 0, time.UTC),
		}

		for _, expectedTime := range times {
			clock.SetTime(expectedTime)
			assert.Equal(t, expectedTime, clock.Now())
		}
	})
}

func TestMockClock_IntegrationWithClient(t *testing.T) {
	t.Parallel()

	t.Run("mock clock integrates with client cache expiry", func(t *testing.T) {
		t.Parallel()
		baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		clock := &internal.MockClock{
			CurrentTime: baseTime,
		}

		config := internal.DefaultConfig()
		config.APIKey = "test-api-key"
		config.CacheTTL = 5 * time.Minute

		client := internal.NewClientWithClock(config, clock)

		// Put item in cache
		data := []byte("test-data")
		client.PutInCache(testKey, data)

		// Verify item is in cache
		retrieved := client.GetFromCache(testKey)
		assert.Equal(t, data, retrieved)

		// Advance clock past TTL
		clock.Advance(6 * time.Minute)

		// Item should be expired
		retrieved = client.GetFromCache(testKey)
		assert.Nil(t, retrieved)
	})
}
