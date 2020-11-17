package redis

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

// Keys command
func (c *RedisCache) Keys(ctx context.Context, key string) ([]string, error) {
	return redis.Strings(c.Do(ctx, "keys", key))
}

//HGetAll get OBJECT from redis
func (c *RedisCache) HGetAll(ctx context.Context, key interface{}, obj interface{}) error {
	//获取缓存
	value, _ := redis.Values(c.Do(ctx, "HGETALL", key))
	return redis.ScanStruct(value, obj)
}

//HMGet 批量获取字段
func (c *RedisCache) HMGet(ctx context.Context, key interface{}, fieldNames ...interface{}) (map[string]interface{}, error) {

	params := []interface{}{key}
	params = append(params, fieldNames...)

	values, err := redis.Values(c.Do(ctx, "HMGET", params...))
	if err != nil {
		return nil, err
	}
	results := make(map[string]interface{})
	for index, value := range values {
		results[fieldNames[index].(string)] = value
	}
	return results, nil
}

//ScanStruct HGETALL：fields传空 HMGET: fields 按照顺序传入field
func (c *RedisCache) ScanStruct(src []interface{}, dest interface{}, fields []string) (err error) {
	if len(src) > 0 {
		//redigo 返回值使用[]byte, go-redis 返回值使用string
		//此处统一转为[]byte
		for i := 0; i < len(src); i++ {
			if src[i] != nil {
				if _, ok := src[i].(string); ok {
					src[i] = []byte(src[i].(string))
				}
			}
		}
	}

	if fields != nil && len(fields) > 0 {
		vs := make([]interface{}, 0, len(src)*2)
		for i := 0; i < len(fields); i++ {
			if src[i] != nil {
				vs = append(vs, []byte(fields[i]), src[i])
			}
		}

		err = redis.ScanStruct(vs, dest)
	} else {
		err = redis.ScanStruct(src, dest)
	}
	return err
}

//HMGetObj 获取结构体
func (c *RedisCache) HMGetObj(ctx context.Context, out interface{}, key string, fieldNames ...string) (notFound bool, err error) {

	values, err := redis.Values(c.Do(ctx, "HMGET", redis.Args{}.Add(key).AddFlat(fieldNames)...))
	if err != nil {
		return false, err
	}
	// if _, ok := values[0].(string); ok {
	// 	return true, nil
	// }

	return false, c.ScanStruct(values, out, fieldNames)
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
