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

	"github.com/gonethopper/nethopper/examples/simple_client/global"

	"github.com/gonethopper/nethopper/config"
	"github.com/gonethopper/nethopper/examples/simple_client/modules/logic"
	"github.com/gonethopper/nethopper/examples/simple_client/modules/wsjson"

	"github.com/gonethopper/nethopper/log"
	. "github.com/gonethopper/nethopper/server"
)

func init() {
	cfg := global.GetInstance().GetConfig()
	flag.StringVar(&cfg.Env, "env", "dev", "the environment and config that used")
	flag.Parse()
	if err := config.InitViper("simple_client", "./conf", cfg.Env, &cfg, false); err != nil {
		panic(err.Error())
	}
}

func main() {
	cfg := global.GetInstance().GetConfig()

	NewNamedModule(ModuleIDLog, "log", log.LogModuleCreate, nil, &cfg.Log)
	NewNamedModule(ModuleIDLogic, "logic", logic.ModuleCreate, nil, &cfg.Logic)
	//NewNamedModule(ModuleIDWSClient, "wspb",wspb.ModuleCreate, nil, &cfg.WS)
	NewNamedModule(ModuleIDWSClient, "wsjson", wsjson.ModuleCreate, nil, &cfg.WS)
	//NewNamedModule(ModuleIDGRPCClient, "grpc", grpc_client.ModuleCreate, nil, &cfg.GPRC)
	// NewNamedModule(ModuleIDTCPClient, "tcp", tcp_client.ModuleCreate, nil, &cfg.TCP)
	// NewNamedModule(ModuleIDKCPClient, "kcp", kcp_client.ModuleCreate, nil, &cfg.KCP)
	// NewNamedModule(ModuleIDQUICClient, "quic", quic_client.ModuleCreate, nil, &cfg.QUIC)
	InitSignal()
	//GracefulExit()
}
