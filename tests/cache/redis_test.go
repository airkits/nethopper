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
// * @Date: 2019-07-04 09:08:18
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-07-04 09:08:18

package cache_test

import (
	"context"
	"testing"

	"github.com/airkits/nethopper/cache"
	"github.com/airkits/nethopper/cache/redis"
)

func TestRedis(t *testing.T) {

	node := cache.NodeInfo{
		ID:       0,
		Password: "",
		Address:  "127.0.0.1:6379",
		DB:       0,
	}
	conf := cache.Config{
		Nodes:           []cache.NodeInfo{node},
		QueueSize:       1000,
		ConnectInterval: 3,
		MaxActive:       10,
		MaxIdle:         8,
		AutoReconnect:   true,
	}

	redisCache, err := redis.NewRedisCache(&conf)
	if err != nil {
		t.Error(err)
	}
	//fmt.Println(redisCache.Version())

	ctx := context.Background()
	key := "aaaa"
	val := "bbb"
	var expr int64 = 0
	if err := redisCache.Set(ctx, key, val, expr); err != nil {
		t.Error(err)
	}
	if v, err := redisCache.GetString(ctx, key); err != nil || val != v {
		t.Error(err)
	}
	if err := redisCache.Del(ctx, key); err != nil {
		t.Error(err)
	}
	var valInt int64 = 99
	if err := redisCache.Set(ctx, key, valInt, expr); err != nil {
		t.Error(err)
	}
	if ret, err := redisCache.Incr(ctx, key); err != nil || ret != valInt+1 {
		t.Error(err)
	}
	if ret, err := redisCache.Decr(ctx, key); err != nil || ret != valInt {
		t.Error(err)
	}

}
