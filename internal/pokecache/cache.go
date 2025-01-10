package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	sync.Mutex
	cache    map[string]cacheEntry
	duration time.Duration
}

func NewCache(time time.Duration) *Cache {
	cache := Cache{
		cache:    map[string]cacheEntry{},
		duration: time,
	}
	cache.reapLoop()
	return &cache
}

func (c *Cache) Add(key string, value []byte) {
	c.Lock()
	defer c.Unlock()

	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}

}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Lock()
	defer c.Unlock()

	v, ok := c.cache[key]
	if !ok {
		return nil, ok
	}

	return v.val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.duration)
	// defer ticker.Stop()

	// spawns a new routine which reads from ticker pipe and delete expired keys
	go func() {
		for {
			select {
			case <-ticker.C:
				now := time.Now()
				c.Lock()
				for k := range c.cache {
					if now.Sub(c.cache[k].createdAt) >= c.duration {
						delete(c.cache, k)
					}
				}
				c.Unlock()
			}
		}

	}()

}
