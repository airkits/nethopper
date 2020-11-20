package redis

import (
	"context"
	"errors"

	"github.com/gomodule/redigo/redis"
)

const (
	redisMustExistField = "must-exist-1234567890" //借用预添加k,v对来处理hmget为空的处理
	redisMustExistValue = "must-exist-value-1234567890"
	redisEmptyField     = "empty-1234567890"       //空记录占位
	redisEmptyValue     = "empty-value-1234567890" //空记录值
)

// ErrEmptyCached 击穿命中的时候触发
var ErrEmptyCached = errors.New("redis: cached nil")

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

//ScanStruct 为hmget返回值拼装数据，values返回非[k,v,k,v]数组
func (c *RedisCache) ScanStruct(src []interface{}, dest interface{}, fields []string) (err error) {
	if src == nil || len(src) <= 0 {
		return errors.New("src is nil")
	}
	src = src[2:] //移除以一个must key
	if fields != nil && len(fields) > 0 {
		vs := make([]interface{}, 0, len(src)*2)
		for i := 0; i < len(fields); i++ {
			if src[i] != nil {
				vs = append(vs, []byte(fields[i]), src[i])
			}
		}
		err = redis.ScanStruct(vs, dest)
	}
	return err
}

//HMGet 批量获取字段 1.get failed ,return nil, err 2.not found return nil,nil,3.get result return result,nil
func (c *RedisCache) HMGet(ctx context.Context, key interface{}, fieldNames ...interface{}) (map[string]interface{}, error) {

	values, err := redis.Values(c.Do(ctx, "HMGET", redis.Args{}.Add(key).Add(redisEmptyField).Add(redisMustExistField).AddFlat(fieldNames)...))
	if err != nil {
		return nil, err
	}
	if s, e := redis.String(values[0], nil); e == redis.ErrNil || s != redisEmptyValue {
		return nil, ErrEmptyCached
	}

	if s, e := redis.String(values[1], nil); e == redis.ErrNil || s != redisMustExistValue {
		return nil, nil
	}
	results := make(map[string]interface{})
	for index, value := range values {
		results[fieldNames[index].(string)] = value
	}
	return results, nil
}

//HMGetObj 获取结构体
func (c *RedisCache) HMGetObj(ctx context.Context, obj interface{}, key string, fieldNames ...string) (notFound bool, err error) {
	if fieldNames != nil && len(fieldNames) > 0 {
		values, err := redis.Values(c.Do(ctx, "HMGET", redis.Args{}.Add(key).Add(redisEmptyField).Add(redisMustExistField).AddFlat(fieldNames)...))
		if err != nil {
			return false, err
		}
		if s, e := redis.String(values[0], nil); e == redis.ErrNil || s != redisEmptyValue {
			return true, ErrEmptyCached
		}
		if s, e := redis.String(values[1], nil); e == redis.ErrNil || s != redisMustExistValue {
			return true, nil
		}
		return false, c.ScanStruct(values, obj, fieldNames)
	}

	values, err := redis.Values(c.Do(ctx, "HGETALL", key))
	if err == redis.ErrNil || len(values) == 0 {
		return true, err
	}
	return false, redis.ScanStruct(values, obj)

}

//HMSetObj 批量添加字段
func (c *RedisCache) HMSetObj(ctx context.Context, key interface{}, obj interface{}) error {
	_, err := c.Do(ctx, "HMSET", redis.Args{}.Add(key).AddFlat(obj).AddFlat(map[string]interface{}{redisMustExistField: redisMustExistValue})...)
	return err
}

//HMSet 批量添加字段
func (c *RedisCache) HMSet(ctx context.Context, key interface{}, fields map[string]interface{}) error {
	_, err := c.Do(ctx, "HMSET", redis.Args{}.Add(key).AddFlat(fields).AddFlat(map[string]interface{}{redisMustExistField: redisMustExistValue})...)
	return err
}

//HSetEmptyValue 设置空值，防止缓存击穿
func (c *RedisCache) HSetEmptyValue(ctx context.Context, key interface{}, expire int64) error {
	if _, err := c.Do(ctx, "HSET", key, redisEmptyField, redisEmptyValue); err != nil {
		return err
	}
	return c.SetExpire(ctx, key.(string), expire)
}
