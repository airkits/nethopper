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
// * @Date: 2019-06-24 12:11:02
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-24 12:11:02

package json

import (
	"encoding/json"

	"github.com/gonethopper/nethopper/codec/common"
)

// JSONCodec use gob encode/decode

// Marshal encode message
func Marshal(v interface{}, template interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal decode message
func Unmarshal(buf []byte, v interface{}, template interface{}) error {
	return json.Unmarshal(buf, v)
}

// Name of codec
func Name() string {
	return "JSONCodec"
}

// Type return codec type
func Type() int {
	return common.CodecTypeJSON
}
