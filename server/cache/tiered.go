package cache

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type TieredCache struct {
	local    *LocalCache
	redis    *redis.Client
	localTTL time.Duration
	redisTTL time.Duration
	sf       singleflight.Group
}

func NewTieredCache(rc *redis.Client, localTTL, redisTTL time.Duration, jitterFactor time.Duration) *TieredCache {
	return &TieredCache{
		local:    NewLocalCache(100000, 30*time.Second),
		redis:    rc,
		localTTL: localTTL + rand.N(jitterFactor),
		redisTTL: redisTTL + rand.N(jitterFactor),
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

	// Use singleflight for cache stampede protection
	result, err, shared := c.sf.Do(key, func() (any, error) {
		// Double check
		// anothe rgoroutine might have populated the local cache
		if value, ok := c.Get(ctx, key); ok {
			return value, nil
		}
		value, err := loader(ctx)
		if err != nil {
			return nil, fmt.Errorf("loader for key %s: %w", key, err)
		}
		_ = c.Set(ctx, key, value)
		return value, nil
	})
	if err != nil {
		return nil, err
	}

	// TODO: Add this is metrics for tracking cache stampedes
	_ = shared

	return result.([]byte), nil
}

func (c *TieredCache) Delete(ctx context.Context, key string) error {
	c.local.Delete(key)
	return c.redis.Del(ctx, key).Err()
}
