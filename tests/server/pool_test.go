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
// * @Date: 2019-06-21 13:42:42
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-21 13:42:42

package server_test

import (
	"testing"

	. "github.com/airkits/nethopper/server"
)

func TestCalcPower(t *testing.T) {
	pw := NewBytesPool()
	if pw.CalcIndex(1) != 0 || pw.CalcIndex(0) != 0 || pw.CalcIndex(32) != 0 || pw.CalcIndex(1<<MinSizePower) != 0 {
		t.Error("calc Index 0 failed")
	}
	var power uint8 = 10
	var value int32 = 1 << power
	var value2 int32 = 1 << (power + 1)
	if pw.CalcIndex(value+1) != int32(power-MinSizePower+1) || pw.CalcIndex(value2-1) != int32(power-MinSizePower+1) {
		t.Errorf("calc Index %d failed", power-MinSizePower+1)
	}
	if pw.CalcIndex(MaxBufferSize) != MaxSizePower-MinSizePower {
		t.Error("calc Index MaxSizePower failed")
	}
	if pw.CalcIndex(MaxBufferSize-1) != MaxSizePower-MinSizePower {
		t.Error("calc Index MaxSizePower failed")
	}
	if pw.CalcIndex(MaxBufferSize*2) != OutMaxBufferPower {
		t.Error("calc Index OutMaxBufferPower failed")
	}

}
