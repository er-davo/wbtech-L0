package cache

import (
	"sync"
)

type Cache[K comparable, V any] struct {
	mu       sync.Mutex
	cache    map[K]V
	list     []K
	capacity int
}

func New[K comparable, V any](capacity int) *Cache[K, V] {
	return &Cache[K, V]{
		cache:    make(map[K]V),
		capacity: capacity,
	}
}

func (c *Cache[K, V]) Add(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.cache[key]; ok {
		return
	}

	if len(c.cache) >= c.capacity {
		delete(c.cache, c.list[0])
		c.list = c.list[1:]
	}

	c.cache[key] = value
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	value, ok := c.cache[key]
	c.list = append(c.list, key)
	c.mu.Unlock()
	return value, ok
}

func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	c.cache = make(map[K]V)
	c.list = []K{}
	c.mu.Unlock()
}
