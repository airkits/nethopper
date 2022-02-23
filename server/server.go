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
// * @Date: 2019-06-12 15:53:22
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-12 15:53:22

package server

import (
	"fmt"

	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/mediator"
)

// global app context
func init() {
	fmt.Println("Nethopper Framework init")
}

// GracefulExit server exit by call root context close
func GracefulExit() {
	mediator.GracefulExit()
	log.GLoggerModule.Close()

	// wait root context done
	// for {
	// 	if _, exitFlag := log.GLoggerModule.CanExit(true); exitFlag {
	// 		return
	// 	}
	// 	time.Sleep(1 * time.Second)
	// }
}
