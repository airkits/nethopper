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
	"flag"

	"github.com/gonethopper/nethopper/config"
	"github.com/gonethopper/nethopper/examples/simple_client/modules/logic"
	"github.com/gonethopper/nethopper/examples/simple_client/modules/wsjson"
	"github.com/gonethopper/nethopper/log"
	"github.com/gonethopper/nethopper/network/grpc"
	"github.com/gonethopper/nethopper/network/kcp"
	"github.com/gonethopper/nethopper/network/quic"
	"github.com/gonethopper/nethopper/network/tcp"
	"github.com/gonethopper/nethopper/network/ws"
	. "github.com/gonethopper/nethopper/server"
)

// Config server config
type Config struct {
	Env  string            `default:"env"`
	Log  log.Config        `yarm:"log"`
	GPRC grpc.ClientConfig `yarm:"grpc_client"`
	KCP  kcp.ClientConfig  `yarm:"kcp_client"`
	QUIC quic.ClientConfig `yarm:"quic_client"`
	TCP  tcp.ClientConfig  `yarm:"tcp_client"`
	WS   ws.ClientConfig   `yarm:"ws_client"`
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

func main() {

	m := map[string]interface{}{
		"filename":     "logs/server.log",
		"level":        DEBUG,
		"maxSize":      5000,
		"maxLines":     10000000,
		"hourEnabled":  false,
		"dailyEnable":  true,
		"queueSize":    1000,
		"wsServerID_0": 0,
		"wsAddress_0":  "ws://127.0.0.1:12080",
		"wsName_0":     "wsclient_0",
	}
	RegisterModule("log", log.LogModuleCreate)
	RegisterModule("logic", logic.ModuleCreate)
	RegisterModule("wsjson", wsjson.ModuleCreate)
	//RegisterModule("grpc", grpc.ModuleCreate)
	//RegisterModule("wspb", wspb.ModuleCreate)
	//RegisterModule("tcp", tcp.ModuleCreate)
	//RegisterModule("kcp", kcp.ModuleCreate)
	//RegisterModule("quic", quic.ModuleCreate)
	NewNamedModule(ModuleIDLog, "log", nil, m)
	NewNamedModule(ModuleIDLogic, "logic", nil, m)
	//NewNamedModule(ModuleIDWSClient, "wspb", nil, m)
	NewNamedModule(ModuleIDWSClient, "wsjson", nil, m)
	//	NewNamedModule(ModuleIDGRPCClient, "grpc", nil, m)
	//	NewNamedModule(ModuleIDTCPClient, "tcp", nil, m)
	//NewNamedModule(ModuleIDKCPClient, "kcp", nil, m)
	//NewNamedModule(ModuleIDQUICClient, "quic", nil, m)
	InitSignal()
	//GracefulExit()
}
