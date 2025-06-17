// Package internal provides core data structures and utilities for the TMDB CLI application.
package internal

// Cache constants.
const (
	MaxCacheSize = 100
)

// getFromCache retrieves data from cache if available and not expired.
func (c *Client) getFromCache(key string) []byte {
	c.cacheMu.RLock()
	defer c.cacheMu.RUnlock()

	item, ok := c.cache[key]
	if !ok || c.clock.Now().After(item.expiry) {
		return nil
	}
	return item.data
}

// putInCache stores data in cache with expiry time.
func (c *Client) putInCache(key string, data []byte) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	c.cache[key] = cachedItem{
		data:   data,
		expiry: c.clock.Now().Add(c.config.CacheTTL),
	}

	// Simple cleanup: remove expired items periodically
	if len(c.cache) > MaxCacheSize {
		go c.cleanupCache()
	}
}

// cleanupCache removes expired items from cache.
func (c *Client) cleanupCache() {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()

	now := c.clock.Now()
	for key, item := range c.cache {
		if now.After(item.expiry) {
			delete(c.cache, key)
		}
	}
}
