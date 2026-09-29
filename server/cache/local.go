package cache

import (
	"sync"
	"time"
)

type entry struct {
	value     any
	expiresAt time.Time
}

type LocalCache struct {
	data    sync.Map
	maxSize int
	size    int64
	mu      sync.Mutex
}

func NewLocalCache(maxSize int, evictInternal time.Duration) *LocalCache {
	cache := &LocalCache{
		maxSize: maxSize,
	}
	go cache.evictionLoop(evictInternal)
	return cache
}

func (c *LocalCache) Get(key string) (any, bool) {
	raw, ok := c.data.Load(key)
	if !ok {
		return nil, false
	}

	e := raw.(*entry)
	if time.Now().After(e.expiresAt) {
		c.data.Delete(key)
		c.decreaseSize()
		return nil, false
	}

	return e.value, true
}

func (c *LocalCache) Set(key string, value any, ttl time.Duration) {
	_, loaded := c.data.LoadOrStore(key, &entry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	})
	if loaded {
		c.increaseSize()
	} else {
		c.data.Store(key, &entry{
			value:     value,
			expiresAt: time.Now().Add(ttl),
		})
	}
}

func (c *LocalCache) Delete(key string) {
	if _, loaded := c.data.LoadAndDelete(key); loaded {
		c.decreaseSize()
	}
}

func (c *LocalCache) evictionLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		c.data.Range(func(key, value any) bool {
			if now.After(value.(entry).expiresAt) {
				c.data.Delete(key)
				c.decreaseSize()
			}

			return true
		})
	}
}

func (c *LocalCache) increaseSize() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.size++
}

func (c *LocalCache) decreaseSize() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.size--
}
