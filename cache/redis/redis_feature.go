package redis

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

// Keys command
func (c *RedisCache) Keys(ctx context.Context, key string) ([]string, error) {
	return redis.Strings(c.Do(ctx, "keys", key))
}
