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
// * @Date: 2019-06-21 12:05:28
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-21 12:05:28

package server

import "sync"

const (
	// MinSizePower min buffer get from pool, means buf size = 1<<7
	MinSizePower = 7
	// MaxSizePower max buffer get from pool, means buf size = 1<<16
	MaxSizePower = 16
	// OutMaxBufferPower means buffer size > MaxBufferSize
	OutMaxBufferPower = MaxSizePower + 1
	// MinBufferSize buffer size 1 << 6
	MinBufferSize = 1 << MinSizePower
	// MaxBufferSize buffer size 1 << 16
	MaxBufferSize = 1 << MaxSizePower
)

// CreatePool []byte pool
func CreatePool(size uint8) *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			buf := make([]byte, size)
			return &buf
		}}
}

// NewBytesPool create buffer pool
func NewBytesPool() *BytesPool {
	bp := &BytesPool{}
	bp.Pools = make([]*sync.Pool, MaxSizePower-MinSizePower+1)
	for index := MinSizePower; index <= MaxSizePower; index++ {
		bp.Pools[index-MinSizePower] = CreatePool(1 << uint8(index))
	}
	return bp
}

// BytesPool alloc and manger buffer
type BytesPool struct {
	Pools []*sync.Pool
}

// CalcIndex get the pool index,if size > MaxBufferSize, return OutMaxBufferPower
func (p *BytesPool) CalcIndex(size int32) int32 {
	if size > MaxBufferSize {
		return OutMaxBufferPower
	}
	if size <= MinBufferSize {
		return 0
	}
	power := int32(0)
	value := size
	for {
		if value <= 1 {
			break
		}
		value >>= 1
		power++
	}
	if size&(size-1) == 0 { //is power of 2
		return power - MinSizePower
	}
	return power - MinSizePower + 1
}

// Alloc borrow buffer from pool,if size > MaxBufferSize, dynamic alloc buffer
func (p *BytesPool) Alloc(size int32) []byte {
	index := p.CalcIndex(size)
	if index == OutMaxBufferPower {
		return make([]byte, size)
	}
	return p.Pools[index].Get().([]byte)
}

// Free return buffer to pool
func (p *BytesPool) Free(buf []byte) {
	size := int32(len(buf))
	index := p.CalcIndex(size)
	if index == OutMaxBufferPower {
		// do nothing
		return
	}
	p.Pools[index].Put(buf)
}
