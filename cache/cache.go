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
// * @Date: 2019-06-11 21:49:35
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-11 21:49:35

package cache

import (
	"context"

	"github.com/airkits/nethopper/config"
)

// ICache interface
type ICache interface {
	// Version cache version
	Version() string
	// Setup init cache with config
	Setup(conf config.IConfig) (ICache, error)
	// Ping to check connection is alive
	Ping() error
	// Get command to get value from cache, control with context
	Get(ctx context.Context, key string) (interface{}, error)
	// GetInt command
	GetInt(ctx context.Context, key string) (int, error)
	// GetInt64 command
	GetInt64(ctx context.Context, key string) (int64, error)
	// GetFloat64 command
	GetFloat64(ctx context.Context, key string) (float64, error)
	// GetString command
	GetString(ctx context.Context, key string) (string, error)
	// GetInts command
	GetInts(ctx context.Context, keys ...interface{}) (map[string]int, error)
	// GetInt64s command
	GetInt64s(ctx context.Context, keys ...interface{}) (map[string]int64, error)
	// GetStrings command
	GetStrings(ctx context.Context, keys ...interface{}) (map[string]string, error)
	// Set command to set value to cache,key is string, if expire(in seconds) is setted, than key will have Expire, in seconds,
	Set(ctx context.Context, key string, val interface{}, expire int64) error
	// Del key from cache
	Del(ctx context.Context, key string) error
	// Exists key in redis, exist return true
	Exists(key string) (bool, error)
	// SetExpire set expire time for key,in seconds
	SetExpire(ctx context.Context, key string, expire int64) error
	// Incr auto-Increment get key and set v++
	Incr(ctx context.Context, key string) (int64, error)
	// Decr auto-Decrement get key and set v--
	Decr(ctx context.Context, key string) (int64, error)
	// Gets command to get multi keys from cache
	Gets(ctx context.Context, keys ...string) (map[string]interface{}, error)
	// Do command to exec custom command
	Do(ctx context.Context, commandName string, args ...interface{}) (reply interface{}, err error)
}
