package unit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alnah/tmdb-cli/internal"
	"github.com/alnah/tmdb-cli/tests/helpers"
)

const testKey = "test-key"

func TestClient_CacheExpiry(t *testing.T) {
	t.Run("cache returns data before expiry", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Create client with 5 minute cache TTL
		config := internal.DefaultConfig()
		config.CacheTTL = 5 * time.Minute

		client := internal.NewTestClient(config, nil, mockClock)

		// Put item in cache
		testData := []byte("test-data")
		client.PutInCache(testKey, testData)

		// Should be in cache immediately
		data := client.GetFromCache(testKey)
		assert.Equal(t, testData, data)

		// Advance clock by 4 minutes (still within TTL)
		mockClock.Advance(4 * time.Minute)

		// Should still be in cache
		data = client.GetFromCache(testKey)
		assert.Equal(t, testData, data)
	})

	t.Run("cache returns nil after expiry", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Create client with 5 minute cache TTL
		config := internal.DefaultConfig()
		config.CacheTTL = 5 * time.Minute

		client := internal.NewTestClient(config, nil, mockClock)

		// Put item in cache
		testData := []byte("test-data")
		client.PutInCache(testKey, testData)

		// Advance clock past TTL
		mockClock.Advance(6 * time.Minute)

		// Should not be in cache
		data := client.GetFromCache(testKey)
		assert.Nil(t, data)
	})

	t.Run("cache handles exact expiry boundary", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Create client with 1 hour cache TTL
		config := internal.DefaultConfig()
		config.CacheTTL = 1 * time.Hour

		client := internal.NewTestClient(config, nil, mockClock)

		// Put item in cache
		testKey := "boundary-test"
		testData := []byte("boundary-data")
		client.PutInCache(testKey, testData)

		// Advance clock to exact expiry time
		mockClock.Advance(1 * time.Hour)

		// At exact expiry, item should still be valid (not After)
		data := client.GetFromCache(testKey)
		assert.Equal(t, testData, data)

		// One nanosecond past expiry, item should be expired
		mockClock.Advance(1 * time.Nanosecond)
		data = client.GetFromCache(testKey)
		assert.Nil(t, data)
	})
}

func TestClient_CacheMultipleItems(t *testing.T) {
	t.Run("cache handles multiple items with different expiry times", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Create client with 10 minute cache TTL
		config := internal.DefaultConfig()
		config.CacheTTL = 10 * time.Minute

		client := internal.NewTestClient(config, nil, mockClock)

		// Put first item
		client.PutInCache("item1", []byte("data1"))

		// Advance clock by 5 minutes
		mockClock.Advance(5 * time.Minute)

		// Put second item (will expire 5 minutes after first)
		client.PutInCache("item2", []byte("data2"))

		// Advance clock by 6 minutes (11 minutes total)
		mockClock.Advance(6 * time.Minute)

		// First item should be expired, second should still be valid
		data1 := client.GetFromCache("item1")
		data2 := client.GetFromCache("item2")

		assert.Nil(t, data1, "item1 should be expired")
		assert.Equal(t, []byte("data2"), data2, "item2 should still be valid")
	})
}

func TestClient_CacheCleanup(t *testing.T) {
	t.Run("expired items are not returned from cache", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Create client with short cache TTL
		config := internal.DefaultConfig()
		config.CacheTTL = 1 * time.Minute

		client := internal.NewTestClient(config, nil, mockClock)

		// Add multiple items
		for i := range 10 {
			key := helpers.CreateKey("item", i)
			data := []byte(helpers.CreateValue("data", i))
			client.PutInCache(key, data)
		}

		// Verify all items are in cache
		for i := range 10 {
			key := helpers.CreateKey("item", i)
			data := client.GetFromCache(key)
			expected := []byte(helpers.CreateValue("data", i))
			assert.Equal(t, expected, data, "item %d should be in cache", i)
		}

		// Advance clock past TTL
		mockClock.Advance(2 * time.Minute)

		// Verify all items are expired and not returned
		for i := range 10 {
			key := helpers.CreateKey("item", i)
			data := client.GetFromCache(key)
			assert.Nil(t, data, "expired item %d should not be returned", i)
		}
	})

	t.Run("cleanup goroutine can be triggered", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		// Create client
		config := internal.DefaultConfig()
		config.CacheTTL = 1 * time.Minute

		client := internal.NewTestClient(config, nil, mockClock)

		// Add exactly MaxCacheSize items
		for i := range 100 { // MaxCacheSize = 100
			key := helpers.CreateKey("item", i)
			data := []byte(helpers.CreateValue("data", i))
			client.PutInCache(key, data)
		}

		// Adding one more item should trigger cleanup
		// This verifies the cleanup goroutine is launched but doesn't test its execution
		client.PutInCache("trigger", []byte("cleanup"))

		// Verify the trigger item was added
		data := client.GetFromCache("trigger")
		assert.Equal(t, []byte("cleanup"), data)
	})
}

func TestClient_CacheKeyNotFound(t *testing.T) {
	t.Run("cache returns nil for non-existent key", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		config := internal.DefaultConfig()
		client := internal.NewTestClient(config, nil, mockClock)

		// Try to get non-existent key
		data := client.GetFromCache("non-existent-key")
		assert.Nil(t, data)
	})
}

func TestClient_CacheOverwrite(t *testing.T) {
	t.Run("putting same key overwrites previous value", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		config := internal.DefaultConfig()
		config.CacheTTL = 10 * time.Minute

		client := internal.NewTestClient(config, nil, mockClock)

		key := "overwrite-test"

		// Put initial value
		client.PutInCache(key, []byte("initial-value"))

		// Verify initial value
		data := client.GetFromCache(key)
		assert.Equal(t, []byte("initial-value"), data)

		// Advance clock by 5 minutes
		mockClock.Advance(5 * time.Minute)

		// Overwrite with new value
		client.PutInCache(key, []byte("new-value"))

		// Verify new value
		data = client.GetFromCache(key)
		assert.Equal(t, []byte("new-value"), data)

		// Advance clock by 6 more minutes (11 total from start)
		mockClock.Advance(6 * time.Minute)

		// Original would be expired, but new value should still be valid
		data = client.GetFromCache(key)
		assert.Equal(t, []byte("new-value"), data)
	})
}

func TestClient_CacheConcurrency(t *testing.T) {
	t.Run("cache handles concurrent reads and writes", func(t *testing.T) {
		// Setup mock clock
		mockClock := &internal.MockClock{
			CurrentTime: time.Now(),
		}

		config := internal.DefaultConfig()
		config.CacheTTL = 1 * time.Hour

		client := internal.NewTestClient(config, nil, mockClock)

		// Run concurrent operations
		done := make(chan bool)
		numGoroutines := 10
		numOperations := 100

		// Writers
		for i := range numGoroutines {
			go func(id int) {
				for j := range numOperations {
					key := helpers.CreateKey("concurrent", id, j)
					data := []byte(helpers.CreateValue("data", id, j))
					client.PutInCache(key, data)
				}
				done <- true
			}(i)
		}

		// Readers
		for i := range numGoroutines {
			go func(id int) {
				for j := range numOperations {
					key := helpers.CreateKey("concurrent", id, j)
					_ = client.GetFromCache(key)
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < numGoroutines*2; i++ {
			<-done
		}

		// Verify some data is present
		for i := range numGoroutines {
			key := helpers.CreateKey("concurrent", i, numOperations-1)
			data := client.GetFromCache(key)
			expected := []byte(helpers.CreateValue("data", i, numOperations-1))
			assert.Equal(t, expected, data)
		}
	})
}

func TestClient_CacheWithDifferentTTLs(t *testing.T) {
	tests := []struct {
		name     string
		ttl      time.Duration
		waitTime time.Duration
		expired  bool
	}{
		{
			name:     "1 second TTL - not expired",
			ttl:      1 * time.Second,
			waitTime: 500 * time.Millisecond,
			expired:  false,
		},
		{
			name:     "1 second TTL - expired",
			ttl:      1 * time.Second,
			waitTime: 2 * time.Second,
			expired:  true,
		},
		{
			name:     "1 minute TTL - not expired",
			ttl:      1 * time.Minute,
			waitTime: 30 * time.Second,
			expired:  false,
		},
		{
			name:     "1 hour TTL - not expired",
			ttl:      1 * time.Hour,
			waitTime: 30 * time.Minute,
			expired:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock clock
			mockClock := &internal.MockClock{
				CurrentTime: time.Now(),
			}

			config := internal.DefaultConfig()
			config.CacheTTL = tt.ttl

			client := internal.NewTestClient(config, nil, mockClock)

			// Put item in cache
			key := "ttl-test"
			data := []byte("ttl-data")
			client.PutInCache(key, data)

			// Advance clock
			mockClock.Advance(tt.waitTime)

			// Check if item is expired
			result := client.GetFromCache(key)
			if tt.expired {
				assert.Nil(t, result, "item should be expired")
			} else {
				assert.Equal(t, data, result, "item should not be expired")
			}
		})
	}
}

// TestClient_CacheTimeProgression tests cache behavior with realistic time progression.
func TestClient_CacheTimeProgression(t *testing.T) {
	t.Run("cache items expire in correct order", func(t *testing.T) {
		// Setup mock clock at specific time
		startTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		mockClock := &internal.MockClock{
			CurrentTime: startTime,
		}

		config := internal.DefaultConfig()
		config.CacheTTL = 5 * time.Minute

		client := internal.NewTestClient(config, nil, mockClock)

		// Add items at different times
		items := []struct {
			key      string
			data     string
			addDelay time.Duration
		}{
			{"item1", "data1", 0},
			{"item2", "data2", 1 * time.Minute},
			{"item3", "data3", 2 * time.Minute},
			{"item4", "data4", 3 * time.Minute},
		}

		for _, item := range items {
			mockClock.SetTime(startTime.Add(item.addDelay))
			client.PutInCache(item.key, []byte(item.data))
		}

		// Check at various time points
		checkPoints := []struct {
			time            time.Duration
			expectedPresent []string
			expectedAbsent  []string
		}{
			{
				time:            4 * time.Minute,
				expectedPresent: []string{"item1", "item2", "item3", "item4"},
				expectedAbsent:  []string{},
			},
			{
				time:            6 * time.Minute,
				expectedPresent: []string{"item2", "item3", "item4"},
				expectedAbsent:  []string{"item1"},
			},
			{
				time:            7 * time.Minute,
				expectedPresent: []string{"item3", "item4"},
				expectedAbsent:  []string{"item1", "item2"},
			},
			{
				time:            8 * time.Minute,
				expectedPresent: []string{"item4"},
				expectedAbsent:  []string{"item1", "item2", "item3"},
			},
			{
				time:            9 * time.Minute,
				expectedPresent: []string{},
				expectedAbsent:  []string{"item1", "item2", "item3", "item4"},
			},
		}

		for _, cp := range checkPoints {
			mockClock.SetTime(startTime.Add(cp.time))

			for _, key := range cp.expectedPresent {
				data := client.GetFromCache(key)
				require.NotNil(t, data, "at %v, %s should be present", cp.time, key)
			}

			for _, key := range cp.expectedAbsent {
				data := client.GetFromCache(key)
				require.Nil(t, data, "at %v, %s should be absent", cp.time, key)
			}
		}
	})
}
