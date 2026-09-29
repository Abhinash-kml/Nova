package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TieredCache struct {
	local    *LocalCache
	redis    *redis.Client
	localTTL time.Duration
	redisTTL time.Duration
}

func NewTieredCache(rc *redis.Client, localTTL, redisTTL time.Duration) *TieredCache {
	return &TieredCache{
		local:    NewLocalCache(100000, 30*time.Second),
		redis:    rc,
		localTTL: localTTL,
		redisTTL: redisTTL,
	}
}

func (c *TieredCache) Get(ctx context.Context, key string) ([]byte, bool) {
	// Check local first
	if val, ok := c.local.Get(key); ok {
		return val.([]byte), true
	}

	// Check redis
	value, err := c.redis.Get(ctx, key).Bytes()
	if err == nil {
		c.local.Set(key, value, c.localTTL)
		return value, true
	}
	return nil, false
}

func (c *TieredCache) Set(ctx context.Context, key string, value []byte) error {
	c.local.Set(key, value, c.localTTL)
	return c.redis.Set(ctx, key, value, c.redisTTL).Err()
}

func (c *TieredCache) GetOrLoad(ctx context.Context, key string, loader func(ctx context.Context) ([]byte, error)) ([]byte, error) {
	if value, ok := c.Get(ctx, key); ok {
		return value, nil
	}
	value, err := loader(ctx)
	if err != nil {
		return nil, fmt.Errorf("loader for key %s: %w", key, err)
	}
	_ = c.Set(ctx, key, value)
	return value, nil
}
