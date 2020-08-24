package redis

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

// Keys command
func (c *RedisCache) Keys(ctx context.Context, key string) ([]string, error) {
	return redis.Strings(c.Do(ctx, "keys", key))
}

//HMGet 批量获取字段
func (c *RedisCache) HMGet(ctx context.Context, key interface{}, fieldNames ...interface{}) (map[string]string, error) {

	params := []interface{}{key}
	params = append(params, fieldNames...)

	values, err := redis.Strings(c.Do(ctx, "HMGET", params...))
	if err != nil {
		return nil, err
	}
	results := make(map[string]string)
	for index, value := range values {
		results[fieldNames[index].(string)] = value
	}
	return results, nil
}

//HMSet 批量添加字段
func (c *RedisCache) HMSet(ctx context.Context, key interface{}, fields map[string]interface{}) error {

	params := []interface{}{key}
	for key, value := range fields {
		params = append(params, key, value)
	}
	_, err := c.Do(ctx, "HMSET", params...)
	return err
}
