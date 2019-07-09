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
	conn := &RedisCache{}
	return conn, nil
}

// NewRedisPool create redis pool by address(ip:port) and pwd
func NewRedisPool(addr string, pwd string, db int, maxIdle int, maxActive int, idleTimeout int) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     maxIdle,   // 最大链接 default 8
		MaxActive:   maxActive, //0：表示最大空闲连接个数
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Wait:        false,
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
	if err := c.readConfig(m); err != nil {
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
func (c *RedisCache) readConfig(m map[string]interface{}) error {
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
	return nil
}

// Get command to get value from cache, control with context
func (c *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	return nil, nil
}

// Set command to set value to cache,key is string, if timeout is setted, than key will have Expire, in seconds,
func (c *RedisCache) Set(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	return nil
}

// Del key from cache
func (c *RedisCache) Del(ctx context.Context, key string) error {
	return nil
}

// Exists key in redis, exist return true
func (c *RedisCache) Exists(key string) bool {
	return true
}

// SetExpire set expire time for key,in seconds
func (c *RedisCache) SetExpire(ctx context.Context, key string, timeout time.Duration) error {
	return nil
}

// Incr auto-Increment get key and set v++
func (c *RedisCache) Incr(ctx context.Context, key string) error {
	return nil
}

// Decr auto-Decrement get key and set v--
func (c *RedisCache) Decr(ctx context.Context, key string) error {
	return nil
}

// Gets command to get multi keys from cache
func (c *RedisCache) Gets(ctx context.Context, keys ...string) (map[string]interface{}, error) {
	return nil, nil
}

// Do command to exec custom command
func (c *RedisCache) Do(ctx context.Context, commandName string, args ...interface{}) (reply interface{}, err error) {
	return nil, nil
}
