package cache

import (
	"log"
	"sync"
	"time"
)

type Cache struct {
	cacheEntries map[string]cacheEntry
	ticker       *time.Ticker
	mu           sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	content   interface{}
	ttl       time.Duration
}

func (c *Cache) Add(key string, val interface{}, ttl time.Duration) {
	log.Println("adding to cache")
	c.mu.Lock()
	defer c.mu.Unlock()
	if ttl == 0 {
		ttl = time.Duration(5 * time.Minute)
	}
	c.cacheEntries[key] = cacheEntry{
		createdAt: time.Now(),
		ttl:       ttl,
		content:   val,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	log.Println("retrieving from cache")
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, exists := c.cacheEntries[key]
	if !exists {
		log.Println("not found in cache")
		return nil, false
	}
	now := time.Now()
	if now.Sub(val.createdAt) >= val.ttl {
		log.Println("cache entry expired")
		return nil, false
	}
	log.Println("found in cache")
	return val.content, true
}

func NewCache(ttl time.Duration) *Cache {
	cache := &Cache{
		cacheEntries: make(map[string]cacheEntry),
		ticker:       time.NewTicker(ttl),
	}
	go cache.reapLoop()
	return cache
}

func (c *Cache) reapLoop() {
	for range c.ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, entry := range c.cacheEntries {
			if now.Sub(entry.createdAt) >= entry.ttl {
				log.Printf("deleting %s from cache\n", key)
				delete(c.cacheEntries, key)
			}
		}
		c.mu.Unlock()
	}
}
