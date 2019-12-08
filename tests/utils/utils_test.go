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
// * @Date: 2019-06-17 11:58:21
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-17 11:58:21

package utils_test

import (
	"path/filepath"
	"testing"

	"github.com/gonethopper/nethopper/utils"
)

func TestGetWorkDirectory(t *testing.T) {
	_, err := utils.GetWorkDirectory()
	if err != nil {
		t.Error(err)
	}

}

func TestGetDirectory(t *testing.T) {
	p := utils.GetAbsDirectory("../logs/utils_test.log")
	workDir, err := utils.GetWorkDirectory()
	if err != nil {
		t.Error(err)
	}
	if p != filepath.Join(workDir, "../logs") {
		t.Error("file path failed")
	}
	p = utils.GetAbsDirectory("/Users/admin/work/nethopper/log/server.info")
	if p != "/Users/admin/work/nethopper/log/" {
		t.Error("file path failed 2")
	}
}

func TestGetAbsFilePath(t *testing.T) {
	p := utils.GetAbsFilePath("log/utils_test.log")
	workDir, err := utils.GetWorkDirectory()
	if err != nil {
		t.Error(err)
	}
	if p != filepath.Join(workDir, "log/utils_test.log") {
		t.Error("file path failed")
	}
	p = utils.GetAbsFilePath("/Users/admin/work/nethopper/log/server.info")
	if p != "/Users/admin/work/nethopper/log/server.info" {
		t.Error("file path failed 2")
	}
}

func TestPowerCalc(t *testing.T) {
	c, p := utils.PowerCalc(8)
	if c != 8 || p != 3 {
		t.Errorf("PowerCalc error %d %d", c, p)
	}

	c, p = utils.PowerCalc(9)
	if c != 16 || p != 4 {
		t.Errorf("PowerCalc error %d %d", c, p)
	}

	c, p = utils.PowerCalc(0)
	if c != 0 || p != 0 {
		t.Errorf("PowerCalc error %d %d", c, p)
	}
	c, p = utils.PowerCalc(1)
	if c != 1 || p != 0 {
		t.Errorf("PowerCalc error %d %d", c, p)
	}
}
