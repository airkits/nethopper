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
// * @Date: 2019-06-14 19:56:49
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-14 19:56:49

package main

import (

	//"github.com/gonethopper/nethopper/cache/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gonethopper/nethopper/examples/simple_server/services/db"
	"github.com/gonethopper/nethopper/examples/simple_server/services/http"
	"github.com/gonethopper/nethopper/examples/simple_server/services/logic"
	"github.com/gonethopper/nethopper/examples/simple_server/services/redis"
	"github.com/gonethopper/nethopper/log"
	. "github.com/gonethopper/nethopper/server"
)

func main() {

	//runtime.GOMAXPROCS(1)
	m := map[string]interface{}{
		"filename":    "logs/server.log",
		"level":       7,
		"maxSize":     50,
		"maxLines":    1000,
		"hourEnabled": false,
		"dailyEnable": true,
		"queueSize":   1000,
		"driver":      "mysql",
		"dsn":         "root:godankye@tcp(127.0.0.1:3306)/game?charset=utf8&parseTime=True&loc=Asia%2FShanghai",
	}
	RegisterService("log", log.LogServiceCreate)
	RegisterService("mysql", db.DBServiceCreate)
	//	RegisterService("tcp", tcp.SocketServiceCreate)
	RegisterService("logic", logic.LogicServiceCreate)
	RegisterService("web_http", http.HTTPServiceCreate)
	RegisterService("redis", redis.RedisServiceCreate)
	//	RegisterService("redis", redis.RedisServiceCreate)
	NewNamedService(ServiceIDLog, "log", nil, m)
	NewNamedService(ServiceIDDB, "mysql", nil, m)
	NewNamedService(ServiceIDRedis, "redis", nil, m)
	//NewNamedService(ServiceIDTCP, "tcp", nil, m)
	NewNamedService(ServiceIDLogic, "logic", nil, m)
	NewNamedService(ServiceIDHTTP, "web_http", nil, m)

	//	NewNamedService(ServiceIDRedis, "redis", nil, m)
	InitSignal()
	//GracefulExit()
}
