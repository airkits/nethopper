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
// * @Date: 2019-12-11 10:13:10
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-11 10:13:10

package rpc_test

import (
	"testing"

	"github.com/gonethopper/nethopper/rpc/http"
)

func TestHTTPRequest(t *testing.T) {
	var content string
	if err := http.Request("http://baidu.com", http.GET, http.RequestTypeText, nil, nil, http.ResponseTypeText, &content, http.ConnTimeoutMS, http.ServeTimeoutMS); err != nil {
		t.Error(err)
		return
	}
	if len(content) <= 0 {
		t.Error("get content failed")
	}
}
