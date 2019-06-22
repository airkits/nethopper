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

	. "github.com/gonethopper/nethopper/server"
	"github.com/gonethopper/nethopper/service"
)

type Factory struct {
	Name string
}

func (g *Factory) CallStructName0() int {
	Debug("CallStructName0")
	return 1
}

func (g *Factory) CallStructName1(value int) int {
	Debug("CallStructName1 %d", value)
	return value
}

func (g *Factory) CallStructName2(value int, name string) int {
	Debug("CallStructName2 %d %s", value, name)
	return value
}

func (g *Factory) CallStructNameArgs(v ...interface{}) int {
	Debug("CallStructName3 %v", v)
	return v[0].(int)
}
func initServer() error {
	RegisterService("log", service.LogServiceCreate)
	m := map[string]interface{}{
		"filename":    "logs/server_call.log",
		"level":       7,
		"maxSize":     50,
		"maxLines":    1000,
		"hourEnabled": false,
		"dailyEnable": true,
		"queueSize":   1000,
	}
	NewNamedService(ServiceIDLog, "log", nil, m)

	return nil
}
func init() {
	initServer()
}

func TestGO(t *testing.T) {

	f := &Factory{Name: "Factory"}

	GO(f.CallStructName0)
	GO(f.CallStructName1, 1)
	GO(f.CallStructName2, 2, "hello2")
	GO(f.CallStructNameArgs, 3, 4, 5, 6, 7)

	WG.Wait()
	GracefulExit()

}

func CallUserFunc0() int {
	return 0
}
func CallUserFunc1(i int) int {
	return i
}
func CallUserFunc2(i int, j int) int {
	return i + j
}
func CallUserFunc3(i int, j int, k string) int {
	return i + j + 1
}

func TestCallUserFunc(t *testing.T) {
	if CallUserFunc(CallUserFunc0)[0].Int() != 0 {
		t.Errorf("CallUserFunc0 != 0")
	}
	if CallUserFunc(CallUserFunc1, 1)[0].Int() != 1 {
		t.Errorf("CallUserFunc1 != 1")
	}
	if CallUserFunc(CallUserFunc2, 1, 1)[0].Int() != 2 {
		t.Errorf("CallUserFunc2 != 2")
	}
	if CallUserFunc(CallUserFunc3, 1, 1, "hello")[0].Int() != 3 {
		t.Errorf("CallUserFunc3 != 3")
	}
}

func TestCallUserMethod(t *testing.T) {
	f := &Factory{Name: "Factory"}
	if CallUserMethod(f, "CallStructName0")[0].Int() != 1 {
		t.Error("CallStructName0 error")
	}
	if CallUserMethod(f, "CallStructName1", 1)[0].Int() != 1 {
		t.Error("CallStructName1 error")
	}
	if CallUserMethod(f, "CallStructName2", 2, "hello")[0].Int() != 2 {
		t.Error("CallStructName2 error")
	}
	if CallUserMethod(f, "CallStructNameArgs", 3, 1, "hello")[0].Int() != 3 {
		t.Error("CallStructNameArgs error")
	}

}
