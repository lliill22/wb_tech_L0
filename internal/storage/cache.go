package storage

import (
	"context"
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	store map[string]Order
}

func NewCache(ctx context.Context, repo *OrderRepository) (*Cache, error) {
	orders, err := repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	c := &Cache{
		store: make(map[string]Order),
	}

	for _, o := range orders {
		c.store[o.OrderUID] = o
	}

	return c, nil
}

func (c *Cache) Set(key string, value Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

func (c *Cache) Get(key string) (*Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.store[key]
	return &val, ok
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}
