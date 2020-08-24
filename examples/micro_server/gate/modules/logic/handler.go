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
// * @Date: 2020-01-09 11:01:34
// * @Last Modified by:   ankye
// * @Last Modified time: 2020-01-09 11:01:34

package logic

import (
	"strconv"

	"github.com/gonethopper/nethopper/examples/micro_server/gamedb/cmd"
	"github.com/gonethopper/nethopper/server"
)

// Login user to login
// @Summary Login
// @Tags LogicModule
// @version 1.0
// @Accept  plain
// @Produce plain
// @Param uid query string true "UserID"
// @Param pwd query string true "Password"
// @Success 200 {string} string 成功后返回值
// @Router /call/Login [put]
func Login(s *Module, obj *server.CallObject, uid string, pwd string) (string, server.Ret) {
	defer server.TraceCost(server.RunModuleFuncName(s))()
	opt, err := strconv.Atoi(uid)

	v, result := server.Call(server.MIDGRPCClient, cmd.GRPCLogin, int32(opt), uid, pwd)
	if err != nil {
		return "", result
	}

	server.Info("get from mysql")
	return v.(string), result
}
