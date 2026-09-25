package callouts

import (
	"sync"
	"time"
)

type cache struct {
	mu   sync.RWMutex
	data map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *cache {
	result := cache{
		data: map[string]cacheEntry{},
	}

	go result.reapLoop(interval)

	return &result
}

func (c *cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if value, ok := c.data[key]; ok {
		return value.val, true
	}

	return []byte{}, false
}

func (c *cache) reapLoop(interval time.Duration) {
	if interval == 0 {
		interval = 10 * time.Second
	}

	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, _ := range c.data {
			if time.Since(c.data[key].createdAt) > interval {
				//fmt.Println("Killing key in cache:", key)
				delete(c.data, key)
			}
		}
		c.mu.Unlock()
	}
}
