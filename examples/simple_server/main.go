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
	"github.com/gonethopper/nethopper/cache"
	"github.com/gonethopper/nethopper/config"
	"github.com/gonethopper/nethopper/database"
	"github.com/gonethopper/nethopper/network/common"
	"github.com/gonethopper/nethopper/network/grpc"
	"github.com/gonethopper/nethopper/network/quic"

	//	"github.com/gonethopper/nethopper/network/grpc"
	"github.com/gonethopper/nethopper/network/kcp"
	//"github.com/gonethopper/nethopper/network/quic"
	"github.com/gonethopper/nethopper/network/tcp"
	"github.com/gonethopper/nethopper/network/ws"

	//"github.com/gonethopper/nethopper/cache/redis"
	"flag"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gonethopper/nethopper/examples/simple_server/docs"
	"github.com/gonethopper/nethopper/examples/simple_server/modules/db"

	//	grpc_server "github.com/gonethopper/nethopper/examples/simple_server/modules/grpc"
	"github.com/gonethopper/nethopper/examples/simple_server/modules/http"
	"github.com/gonethopper/nethopper/examples/simple_server/modules/logic"

	//	quic_server "github.com/gonethopper/nethopper/examples/simple_server/modules/quic"
	"github.com/gonethopper/nethopper/examples/simple_server/modules/redis"
	"github.com/gonethopper/nethopper/examples/simple_server/modules/wsjson"
	"github.com/gonethopper/nethopper/log"
	http_server "github.com/gonethopper/nethopper/network/http"
	. "github.com/gonethopper/nethopper/server"
)

// Config server config
type Config struct {
	Env   string                   `default:"env"`
	Log   log.Config               `mapstructure:"log"`
	GPRC  grpc.ServerConfig        `mapstructure:"grpc"`
	KCP   kcp.ServerConfig         `mapstructure:"kcp"`
	QUIC  quic.ServerConfig        `mapstructure:"quic"`
	TCP   tcp.ServerConfig         `mapstructure:"tcp"`
	WS    ws.ServerConfig          `mapstructure:"wsjson"`
	Logic common.LogicConfig       `mapstructure:"logic"`
	Mysql database.Config          `mapstructure:"mysql"`
	Redis cache.Config             `mapstructure:"redis"`
	HTTP  http_server.ServerConfig `mapstructure:"http"`
}

var cfg Config

//GetViper get config
func GetViper() *Config {
	return &cfg
}

func init() {

	flag.StringVar(&cfg.Env, "env", "dev", "the environment and config that used")
	flag.Parse()
	if err := config.InitViper("simple_server", "./conf", cfg.Env, &cfg, false); err != nil {
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
	NewNamedModule(ModuleIDDB, "mysql", db.ModuleCreate, nil, &cfg.Mysql)
	NewNamedModule(ModuleIDRedis, "redis", redis.ModuleCreate, nil, &cfg.Redis)
	NewNamedModule(ModuleIDLogic, "logic", logic.ModuleCreate, nil, &cfg.Logic)
	NewNamedModule(ModuleIDHTTP, "http", http.ModuleCreate, nil, &cfg.HTTP)
	NewNamedModule(ModuleIDWSServer, "wsjson", wsjson.ModuleCreate, nil, &cfg.WS)
	//NewNamedModule(ModuleIDWSServer, "wspb", wspb.ModuleCreate, nil,&cfg.Log)
	//NewNamedModule(ModuleIDGRPCServer, "grpc", grpc_server.ModuleCreate, nil, &cfg.GPRC)
	//NewNamedModule(ModuleIDTCP, "tcp",tcp.ModuleCreate, nil, &cfg.Tcp)
	//NewNamedModule(ModuleIDKCP, "kcp",kcp.ModuleCreate, nil, &cfg.Kcp)
	//NewNamedModule(ModuleIDQUIC, "quic", quic_server.ModuleCreate, nil, &cfg.QUIC)
	InitSignal()
	//GracefulExit()
}
