package middleware

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// CacheEntry represents a cached response
type CacheEntry struct {
	Response  *CachedResponse
	ExpiresAt time.Time
}

// CachedResponse represents a cached HTTP response
type CachedResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// Cache middleware for HTTP responses
type cacheMiddleware struct {
	cache map[string]*CacheEntry
	ttl   time.Duration
	mu    sync.RWMutex
}

// NewCache creates a new cache middleware
func NewCache(ttl time.Duration) Middleware {
	cm := &cacheMiddleware{
		cache: make(map[string]*CacheEntry),
		ttl:   ttl,
	}
	
	// Start cleanup goroutine
	go cm.cleanup()
	
	return cm
}

func (c *cacheMiddleware) Before(req *http.Request) error {
	// Only cache GET requests
	if req.Method != "GET" {
		return nil
	}
	
	key := c.generateKey(req)
	
	c.mu.RLock()
	entry, exists := c.cache[key]
	c.mu.RUnlock()
	
	if exists && time.Now().Before(entry.ExpiresAt) {
		// Cache hit - we'll handle this in a custom way
		// For now, just mark the request as cacheable
		req.Header.Set("X-Cache-Key", key)
	}
	
	return nil
}

func (c *cacheMiddleware) After(resp *http.Response) {
	// Only cache successful GET responses
	if resp.Request.Method != "GET" || resp.StatusCode >= 400 {
		return
	}
	
	key := resp.Request.Header.Get("X-Cache-Key")
	if key == "" {
		key = c.generateKey(resp.Request)
	}
	
	// Read and cache the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()
	
	// Create cached response
	cachedResp := &CachedResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header.Clone(),
		Body:       body,
	}
	
	// Store in cache
	c.mu.Lock()
	c.cache[key] = &CacheEntry{
		Response:  cachedResp,
		ExpiresAt: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()
	
	// Restore body for the original response
	resp.Body = io.NopCloser(bytes.NewReader(body))
}

func (c *cacheMiddleware) generateKey(req *http.Request) string {
	key := fmt.Sprintf("%s:%s", req.Method, req.URL.String())
	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("%x", hash)
}

func (c *cacheMiddleware) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		now := time.Now()
		c.mu.Lock()
		for key, entry := range c.cache {
			if now.After(entry.ExpiresAt) {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}

// GetCachedResponse retrieves a cached response if available
func (c *cacheMiddleware) GetCachedResponse(req *http.Request) (*CachedResponse, bool) {
	if req.Method != "GET" {
		return nil, false
	}
	
	key := c.generateKey(req)
	
	c.mu.RLock()
	entry, exists := c.cache[key]
	c.mu.RUnlock()
	
	if !exists || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	
	return entry.Response, true
}