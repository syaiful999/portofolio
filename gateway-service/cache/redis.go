package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(cache *redis.Client) Cache {
	return &RedisCache{
		client: cache,
	}
}

func (c *RedisCache) Context(ctx context.Context) Cache {
	c.ctx = ctx
	return c
}

func (c *RedisCache) Put(key string, data interface{}, expiration time.Duration) error {
	err := c.client.Set(context.Background(), key, data, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisCache) Get(key string) (interface{}, time.Time, error) {
	data, err := c.client.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, time.Time{}, err
	}
	return data, time.Time{}, nil
}

func (c *RedisCache) Delete(key string) error {
	if _, _, ok := c.Get(key); ok != nil {
		return ErrKeyNotFound
	}
	_, err := c.client.Del(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *RedisCache) Expired(expiration time.Time) bool {
	if expiration.IsZero() {
		return true
	}

	return time.Now().UnixNano() > expiration.UnixNano()
}
