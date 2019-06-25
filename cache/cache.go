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
	"time"
)

// Cache interface
type Cache interface {
	// Version cache version
	Version() string
	// Setup init cache with config
	Setup(m map[string]interface{}) (Cache, error)
	// Ping to check connection is alive
	Ping() error
	// Get command to get value from cache, control with context
	Get(ctx context.Context, key string) (interface{}, error)
	// Set command to set value to cache,key is string, if timeout is setted, than key will have Expire, in seconds,
	Set(ctx context.Context, key string, val interface{}, timeout time.Duration) error
	// Del key from cache
	Del(ctx context.Context, key string) error
	// Exists key in redis, exist return true
	Exists(key string) bool
	// SetExpire set expire time for key,in seconds
	SetExpire(ctx context.Context, key string, timeout time.Duration) error
	// Incr auto-Increment get key and set v++
	Incr(ctx context.Context, key string) error
	// Decr auto-Decrement get key and set v--
	Decr(ctx context.Context, key string) error
	// Gets command to get multi keys from cache
	Gets(ctx context.Context, keys ...string) (map[string]interface{}, error)
	// Do command to exec custom command
	Do(ctx context.Context, commandName string, args ...interface{}) (reply interface{}, err error)
}
