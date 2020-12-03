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

/**
Redis 有序集合和集合一样也是string类型元素的集合,且不允许重复的成员。
不同的是每个元素都会关联一个double类型的分数。redis正是通过分数来为集合中的成员进行从小到大的排序。
有序集合的成员是唯一的,但分数(score)却可以重复。
集合是通过哈希表实现的，所以添加，删除，查找的复杂度都是O(1)。
**/

// ZAdd 将一个 member 元素及其 score 值加入到有序集 key 当中。
func (c *RedisCache) ZAdd(ctx context.Context, key string, score int64, member string) (reply interface{}, err error) {
	return c.Do(ctx, "ZADD", key, score, member)
}

// ZRem 移除有序集 key 中的一个成员，不存在的成员将被忽略。
func (c *RedisCache) ZRem(ctx context.Context, key string, member string) (reply interface{}, err error) {
	return c.Do(ctx, "ZREM", key, member)
}

// ZScore 返回有序集 key 中，成员 member 的 score 值。 如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 nil 。
func (c *RedisCache) ZScore(ctx context.Context, key string, member string) (int64, error) {
	return redis.Int64(c.Do(ctx, "ZSCORE", key, member))
}

// ZRank 返回有序集中指定成员的排名。其中有序集成员按分数值递增(从小到大)顺序排列。score 值最小的成员排名为 0
func (c *RedisCache) ZRank(ctx context.Context, key, member string) (int64, error) {
	return redis.Int64(c.Do(ctx, "ZRANK", key, member))
}

// ZRevrank 返回有序集中成员的排名。其中有序集成员按分数值递减(从大到小)排序。分数值最大的成员排名为 0 。
func (c *RedisCache) ZRevrank(ctx context.Context, key, member string) (int64, error) {
	return redis.Int64(c.Do(ctx, "ZREVRANK", key, member))
}

// ZRange 返回有序集中，指定区间内的成员。其中成员的位置按分数值递增(从小到大)来排序。具有相同分数值的成员按字典序(lexicographical order )来排列。
// 以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。或 以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (c *RedisCache) ZRange(ctx context.Context, key string, from, to int64) (map[string]int64, error) {
	return redis.Int64Map(c.Do(ctx, "ZRANGE", key, from, to, "WITHSCORES"))
}

// ZRevrange 返回有序集中，指定区间内的成员。其中成员的位置按分数值递减(从大到小)来排列。具有相同分数值的成员按字典序(lexicographical order )来排列。
// 以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。或 以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (c *RedisCache) ZRevrange(ctx context.Context, key string, from, to int64) (map[string]int64, error) {
	return redis.Int64Map(c.Do(ctx, "ZREVRANGE", key, from, to, "WITHSCORES"))
}

// ZRangeByScore 返回有序集合中指定分数区间的成员列表。有序集成员按分数值递增(从小到大)次序排列。
// 具有相同分数值的成员按字典序来排列
func (c *RedisCache) ZRangeByScore(ctx context.Context, key string, from, to, offset int64, count int) (map[string]int64, error) {
	return redis.Int64Map(c.Do(ctx, "ZRANGEBYSCORE", key, from, to, "WITHSCORES", "LIMIT", offset, count))
}

// ZRevrangeByScore 返回有序集中指定分数区间内的所有的成员。有序集成员按分数值递减(从大到小)的次序排列。
// 具有相同分数值的成员按字典序来排列
func (c *RedisCache) ZRevrangeByScore(ctx context.Context, key string, from, to, offset int64, count int) (map[string]int64, error) {
	return redis.Int64Map(c.Do(ctx, "ZREVRANGEBYSCORE", key, from, to, "WITHSCORES", "LIMIT", offset, count))
}
