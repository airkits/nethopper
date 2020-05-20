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
	"github.com/gonethopper/nethopper/config"
	"github.com/gonethopper/nethopper/network/common"
	"github.com/gonethopper/nethopper/network/kcp"
	"github.com/gonethopper/nethopper/network/tcp"
	"github.com/gonethopper/nethopper/network/ws"

	//"github.com/gonethopper/nethopper/cache/redis"
	"flag"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gonethopper/nethopper/examples/simple_server/docs"
	"github.com/gonethopper/nethopper/examples/simple_server/modules/db"
	"github.com/gonethopper/nethopper/examples/simple_server/modules/grpc"
	"github.com/gonethopper/nethopper/examples/simple_server/modules/http"
	"github.com/gonethopper/nethopper/examples/simple_server/modules/logic"
	"github.com/gonethopper/nethopper/examples/simple_server/modules/quic"
	"github.com/gonethopper/nethopper/examples/simple_server/modules/redis"
	"github.com/gonethopper/nethopper/examples/simple_server/modules/wsjson"
	"github.com/gonethopper/nethopper/log"
	. "github.com/gonethopper/nethopper/server"
)

// Config server config
type Config struct {
	Env   string             `default:"env"`
	Log   log.Config         `mapstructure:"log"`
	GPRC  grpc.ServerConfig  `mapstructure:"grpc"`
	KCP   kcp.ServerConfig   `mapstructure:"kcp"`
	QUIC  quic.ServerConfig  `mapstructure:"quic"`
	TCP   tcp.ServerConfig   `mapstructure:"tcp"`
	WS    ws.ServerConfig    `mapstructure:"ws"`
	Logic common.LogicConfig `mapstructure:"logic"`
}

var cfg Config

//GetViper get config
func GetViper() *Config {
	return &cfg
}

func init() {

	flag.StringVar(&cfg.Env, "env", "dev", "the environment and config that used")
	flag.Parse()
	if err := config.InitViper("simple_client", "./conf", cfg.Env, &cfg, false); err != nil {
		panic(err.Error())
	}
}

// @title Nethopper Simple Server
// @version 1.0.2
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:11080
// @BasePath
func main() {

	//runtime.GOMAXPROCS(1)

	NewNamedModule(ModuleIDLog, "log", log.LogModuleCreate, nil, &cfg.Log)
	NewNamedModule(ModuleIDDB, "mysql", db.ModuleCreate, nil, m)
	NewNamedModule(ModuleIDRedis, "redis", redis.ModuleCreate, nil, m)
	NewNamedModule(ModuleIDLogic, "logic", logic.ModuleCreate, nil, m)
	NewNamedModule(ModuleIDHTTP, "http", http.ModuleCreate, nil, m)
	NewNamedModule(ModuleIDWSServer, "wsjson", wsjson.ModuleCreate, nil, m)
	//NewNamedModule(ModuleIDWSServer, "wspb", wspb.ModuleCreate, nil, m)
	NewNamedModule(ModuleIDGRPCServer, "grpc", grpc.ModuleCreate, nil, m)
	//NewNamedModule(ModuleIDTCP, "tcp",tcp.ModuleCreate, nil, m)
	//NewNamedModule(ModuleIDKCP, "kcp",kcp.ModuleCreate, nil, m)
	//NewNamedModule(ModuleIDQUIC, "quic", quic.ModuleCreate, nil, m)
	InitSignal()
	//GracefulExit()
}
