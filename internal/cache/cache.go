package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	RDB *redis.Client
	TTL time.Duration
}

func NewCache(rdb *redis.Client, defaultTTL time.Duration) *Cache {
	return &Cache{
		RDB: rdb,
		TTL: defaultTTL,
	}
}

func (c *Cache) Set(
	ctx context.Context,
	key string,
	value interface{},
	ttl ...time.Duration,
) error {
	expiration := c.TTL

	if len(ttl) > 0 {
		expiration = ttl[0]
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.RDB.Set(ctx, key, data, expiration).Err()
}

// / Retrieves a value from the cache
// # Parameters
// - context
// - key(cache key)
// - dest(value)
// # Return
// - true
// - error
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	data, err := c.RDB.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return false, err
	}

	return true, nil
}

// / Removes values from the cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.RDB.Del(ctx, key).Err()
}

// / Invalidate(getting rid of) outdated data from the cache
// # Return
// - error
func (c *Cache) Invalidate(ctx context.Context, pattern string) error {
	keys, err := c.RDB.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.RDB.Del(ctx, keys...).Err()
	}
	return nil
}
