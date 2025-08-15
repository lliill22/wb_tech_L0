package storage

import (
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	store map[string]Order
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string]Order),
	}
}

func (c *Cache) Set(key string, value Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

func (c *Cache) Get(key string) (Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.store[key]
	return val, ok
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}
