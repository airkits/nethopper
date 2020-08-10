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
	"flag"

	"github.com/gonethopper/nethopper/config"
	_ "github.com/gonethopper/nethopper/examples/micro_server/gate/docs"
	"github.com/gonethopper/nethopper/examples/micro_server/gate/global"
	"github.com/gonethopper/nethopper/examples/micro_server/gate/modules/gclient"
	"github.com/gonethopper/nethopper/examples/micro_server/gate/modules/http"
	"github.com/gonethopper/nethopper/examples/micro_server/gate/modules/logic"
	"github.com/gonethopper/nethopper/examples/micro_server/gate/modules/wsjson"
	"github.com/gonethopper/nethopper/log"
	. "github.com/gonethopper/nethopper/server"
)

func init() {
	cfg := global.GetInstance().GetConfig()
	flag.StringVar(&cfg.Env, "env", "dev", "the environment and config that used")
	flag.Parse()
	if err := config.InitViper("gate", "./conf", cfg.Env, &cfg, false); err != nil {
		panic(err.Error())
	}
}

// @title Nethopper Micro Server - Gate
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
	cfg := global.GetInstance().GetConfig()
	NewNamedModule(MIDLog, "log", log.LogModuleCreate, nil, &cfg.Log)
	NewNamedModule(MIDLogic, "logic", logic.ModuleCreate, nil, &cfg.Logic)
	NewNamedModule(MIDHTTP, "http", http.ModuleCreate, nil, &cfg.HTTP)
	NewNamedModule(MIDWSServer, "wsjson", wsjson.ModuleCreate, nil, &cfg.WSJSON)
	NewNamedModule(MIDGRPCClient, "gclient", gclient.ModuleCreate, nil, &cfg.GPRCClient)
	InitSignal()
	//GracefulExit()
}
