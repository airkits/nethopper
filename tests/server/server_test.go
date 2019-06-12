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
// * @Date: 2019-06-12 17:11:57
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-12 17:11:57

package server_test

import (
	"testing"

	"github.com/gonethopper/nethopper/log"
	"github.com/gonethopper/nethopper/server"
)

type Factory struct {
	Name string
}

func (g *Factory) CallStructName0() {
	server.Logger.Debug("CallStructName0")
}

func (g *Factory) CallStructName1(value int) {
	server.Logger.Debug("CallStructName1 %d \n", value)
}

func (g *Factory) CallStructName2(value int, name string) {
	server.Logger.Debug("CallStructName2 %d %s \n", value, name)
}

func (g *Factory) CallStructNameArgs(v ...interface{}) {
	server.Logger.Debug("CallStructNameArgs %v \n", v)
}

func TestGO(t *testing.T) {
	m := map[string]interface{}{
		"filename":    "server.log",
		"level":       7,
		"maxSize":     50,
		"maxLines":    1000,
		"hourEnabled": false,
		"dailyEnable": true,
	}
	logger, err := log.NewFileLogger(m)
	if err != nil {
		t.Error(err)
	}
	server.Logger = logger

	f := &Factory{Name: "Factory"}

	server.GO(f.CallStructName0)
	server.GO(f.CallStructName1, 1)
	//	runtime.GO(f.CallStructName2, 2, "hello")
	//server.GO(f.CallStructNameArgs, 3, 4, 5, 6, 7)

	server.WG.Wait()
}
