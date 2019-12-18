// MIT License

// Copyright (c) 2019 gonethopper

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// * @Author: ankye
// * @Date: 2019-06-24 11:02:59
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-24 11:02:59

package redis

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gonethopper/nethopper/server"
	"github.com/pkg/errors"
)

// NewRedisCache create redis cache instance
func NewRedisCache(m map[string]interface{}) (*RedisCache, error) {
	cache := &RedisCache{}
	return cache.Setup(m)

}

// NewRedisPool create redis pool by address(ip:port) and pwd
func NewRedisPool(addr string, pwd string, db int, maxIdle int, maxActive int, idleTimeout int) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     maxIdle,   // 最大链接 default 8
		MaxActive:   maxActive, //0：表示最大空闲连接个数
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, errors.Wrap(err, "[backend] Redis Dial failed")
			}
			if pwd != "" {
				if _, err := c.Do("AUTH", pwd); err != nil {
					c.Close()
					return nil, errors.Wrap(err, "[backend] Redis AUTH failed")
				}
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, errors.Wrap(err, "[backend] Redis DB select failed")
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return pool
}

// RedisCache use redis as cache
type RedisCache struct {
	Address     string
	Password    string
	maxIdle     int
	maxActive   int
	idleTimeout int
	db          int
	pool        *redis.Pool
}

// Setup init cache with config
func (c *RedisCache) Setup(m map[string]interface{}) (*RedisCache, error) {
	if err := c.ReadConfig(m); err != nil {
		return nil, err
	}
	c.pool = NewRedisPool(c.Address, c.Password, c.db, c.maxIdle, c.maxActive, c.idleTimeout)

	return c, nil
}

// config map
// maxIdle default 8
// maxActive default 10
// idleTimeout default 300
// address default 127.0.0.1:6379
// password default ""
func (c *RedisCache) ReadConfig(m map[string]interface{}) error {
	maxIdle, err := server.ParseValue(m, "maxIdle", 8)
	if err != nil {
		return err
	}
	c.maxIdle = maxIdle.(int)

	maxActive, err := server.ParseValue(m, "maxActive", 10)
	if err != nil {
		return err
	}
	c.maxActive = maxActive.(int)

	idleTimeout, err := server.ParseValue(m, "idleTimeout", 300)
	if err != nil {
		return err
	}
	c.idleTimeout = idleTimeout.(int)

	address, err := server.ParseValue(m, "address", "127.0.0.1:6379")
	if err != nil {
		return err
	}
	c.Address = address.(string)
	password, err := server.ParseValue(m, "password", "")
	if err != nil {
		return err
	}
	c.Password = password.(string)
	db, err := server.ParseValue(m, "db", 0)
	if err != nil {
		return err
	}
	c.db = db.(int)

	return nil
}

// Version cache version
func (c *RedisCache) Version() string {
	conn := c.pool.Get()
	r, e := redis.String(conn.Do("INFO"))
	if e != nil && e != redis.ErrNil {
		return e.Error()
	}
	return r
}

// Ping to check connection is alive
func (c *RedisCache) Ping() error {
	return errors.New("redis pool internal processing")
}

// Get command to get value from cache, control with context
func (c *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	return c.Do(ctx, "GET", key)
}

// GetInt command
func (c *RedisCache) GetInt(ctx context.Context, key string) (int, error) {
	return redis.Int(c.Get(ctx, key))
}

// GetInt64 command
func (c *RedisCache) GetInt64(ctx context.Context, key string) (int64, error) {
	return redis.Int64(c.Get(ctx, key))
}

// GetFloat64 command
func (c *RedisCache) GetFloat64(ctx context.Context, key string) (float64, error) {
	return redis.Float64(c.Get(ctx, key))
}

// GetString command
func (c *RedisCache) GetString(ctx context.Context, key string) (string, error) {
	return redis.String(c.Get(ctx, key))
}

// GetInts command
func (c *RedisCache) GetInts(ctx context.Context, keys ...interface{}) (map[string]int, error) {
	return redis.IntMap(c.Do(ctx, "MGET", keys))
}

// GetInt64s command
func (c *RedisCache) GetInt64s(ctx context.Context, keys ...interface{}) (map[string]int64, error) {
	return redis.Int64Map(c.Do(ctx, "MGET", keys))
}

// GetStrings command
func (c *RedisCache) GetStrings(ctx context.Context, keys ...interface{}) (map[string]string, error) {
	return redis.StringMap(c.Do(ctx, "MGET", keys))
}

// Gets command to get multi keys from cache
func (c *RedisCache) Gets(ctx context.Context, keys ...string) (map[string]interface{}, error) {

	v, err := redis.Values(c.Do(ctx, "MGET", keys))
	// If field is not found then return map with fields as nil
	if len(v) == 0 || err == redis.ErrNil {
		v = make([]interface{}, len(keys))
	}

	// Form a map with returned results
	res := make(map[string]interface{})
	for i, k := range keys {
		res[k] = v[i]
	}
	return res, err

}

// Set command to set value to cache,key is string, if expire(second) is setted bigger than 0, than key will have Expire time, in seconds,
func (c *RedisCache) Set(ctx context.Context, key string, val interface{}, expire int64) error {

	if expire > 0 {
		_, err := c.Do(ctx, "SETEX", key, expire, val)
		if err != nil {
			return err
		}
	} else {
		_, err := c.Do(ctx, "SET", key, val)
		if err != nil {
			return err
		}
	}
	return nil
}

// Del key from cache
func (c *RedisCache) Del(ctx context.Context, key string) error {
	if _, err := c.Do(ctx, "DEL", key); err != nil {
		return err
	}
	return nil
}

// Exists key in redis, exist return true
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	ret, err := c.Do(ctx, "EXISTS", key)
	if err != nil {
		return false, err
	}
	return redis.Bool(ret, err)
}

// SetExpire set expire time for key,expire(in seconds)
func (c *RedisCache) SetExpire(ctx context.Context, key string, expire int64) error {
	_, err := c.Do(ctx, "EXPIRE", key, expire)
	if err != nil {
		return err
	}
	return nil
}

// Incr auto-Increment get key and set v++
func (c *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	ret, err := c.Do(ctx, "INCR", key)
	if err != nil {
		return -1, err
	}
	return redis.Int64(ret, err)
}

// Decr auto-Decrement get key and set v--
func (c *RedisCache) Decr(ctx context.Context, key string) (int64, error) {
	ret, err := c.Do(ctx, "DECR", key)
	if err != nil {
		return -1, err
	}
	return redis.Int64(ret, err)
}

// Do command to exec custom command
// if redis return redis.ErrNil should convert to value null and err null
func (c *RedisCache) Do(ctx context.Context, commandName string, args ...interface{}) (reply interface{}, err error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	ret, err := conn.Do(commandName, args...)
	if err != nil && err == redis.ErrNil {
		return nil, nil
	}
	return ret, err
}
