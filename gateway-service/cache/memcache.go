package cache

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
)

type MemCache struct {
	client *cache.Cache
	ctx    context.Context
}

func NewMemCache(cache *cache.Cache) Cache {
	return &MemCache{
		client: cache,
	}
}

func (c *MemCache) Context(ctx context.Context) Cache {
	c.ctx = ctx
	return c
}

func (c *MemCache) Put(key string, data interface{}, expiration time.Duration) error {
	return c.client.Add(key, data, expiration)
}

func (c *MemCache) Get(key string) (interface{}, time.Time, error) {
	if _, ok := c.client.Get(key); !ok {
		return nil, time.Time{}, ErrKeyNotFound
	}
	cache, duration, _ := c.client.GetWithExpiration(key)
	return cache, duration, nil
}

func (c *MemCache) Delete(key string) error {
	if _, ok := c.client.Get(key); !ok {
		return ErrKeyNotFound
	}
	c.client.Delete(key)
	return nil
}

func (c *MemCache) Expired(expiration time.Time) bool {
	if expiration.IsZero() {
		return true
	}

	return time.Now().UnixNano() > expiration.UnixNano()
}
