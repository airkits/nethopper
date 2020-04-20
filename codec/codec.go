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
// * @Date: 2019-06-24 11:50:26
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-24 11:50:26

package codec

import (
	"github.com/gonethopper/nethopper/codec/gob"
	"github.com/gonethopper/nethopper/codec/json"
	"github.com/gonethopper/nethopper/codec/pb"
	"github.com/gonethopper/nethopper/codec/raw"
)

// Codec encodes/decodes
type Codec struct {
	// marshal encode message
	Marshal func(v interface{}, template interface{}) ([]byte, error)
	// Unmarshal decode message
	Unmarshal func(buf []byte, v interface{}, template interface{}) error
	// Name of codec
	Name func() string
	Type func() int
}

var (
	// PBCodec protobuf encode/decode
	PBCodec = Codec{pb.Marshal, pb.Unmarshal, pb.Name, pb.Type}
	// JSONCodec json encode/decode
	JSONCodec = Codec{json.Marshal, json.Unmarshal, json.Name, json.Type}
	// GobCodec gob encode/decode
	GobCodec = Codec{gob.Marshal, gob.Unmarshal, gob.Name, gob.Type}
	// RawCodec binary encode/decode
	RawCodec = Codec{raw.Marshal, raw.Unmarshal, raw.Name, raw.Type}
)
