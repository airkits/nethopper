package binary

import (
	"encoding/binary"
	"math"
)

// BinaryCodec use binary encode/decode

const (
	// BinaryCodecType type enum
	BinaryCodecType = iota
	// BEuint8 big endian uint8 type
	BEuint8
	// BEuint16 big endian uint16 type
	BEuint16
	// BEuint32 big endian uint32 type
	BEuint32
	// BEuint54 big endian uint64 type
	BEuint54
	// BEstring big endian string type
	BEstring
	// BEbytes big endian bytes type
	BEbytes

	// LEuint8 little endian uint8 type
	LEuint8
	// LEuint16 little endian uint16 type
	LEuint16
	// LEuint32 little endian uint32 type
	LEuint32
	// LEuint54 little endian uint64 type
	LEuint54
	// LEstring little endian string type
	LEstring
	// LEbytes little endian bytes type
	LEbytes
)

// Marshal encode message
func Marshal(v interface{}, template []byte) ([]byte, error) {
	return nil, nil

}

// Unmarshal decode message
func Unmarshal(buf []byte, v interface{}, template []byte) error {
	return nil
}

// Name of codec
func Name() string {
	return "BinaryCodec"
}

//////// big endian //////////////////////////////

// DecodeBEuint8 get uint8 from big endian bytes
func DecodeBEuint8(b []byte) uint8 {
	return uint8(b[0])
}

// EncodeBEuint8 put uint8 to big endian bytes
func EncodeBEuint8(b []byte, v uint8) {
	b[0] = byte(v)
}

// DecodeBEuint16 get uint16 from big endian bytes
func DecodeBEuint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

// EncodeBEuint16 put uint16 to big endian bytes
func EncodeBEuint16(b []byte, v uint16) {
	binary.BigEndian.PutUint16(b, v)
}

// DecodeBEuint32 get uint32 from big endian bytes
func DecodeBEuint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

// EncodeBEuint32 put uint32 to big endian bytes
func EncodeBEuint32(b []byte, v uint32) {
	binary.BigEndian.PutUint32(b, v)
}

// DecodeBEuint64 get uint64 from big endian bytes
func DecodeBEuint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// EncodeBEuint64 put uint64 to big endian bytes
func EncodeBEuint64(b []byte, v uint64) {
	binary.BigEndian.PutUint64(b, v)

}

// DecodeBEfloat32 get float32 from big endian bytes
func DecodeBEfloat32(b []byte) float32 {
	return math.Float32frombits(DecodeBEuint32(b))
}

// EncodeBEfloat32 put float32 to big endian bytes
func EncodeBEfloat32(b []byte, v float32) {
	EncodeBEuint32(b, math.Float32bits(v))

}

// DecodeBEfloat64 get float64 from big endian bytes
func DecodeBEfloat64(b []byte) float64 {
	return math.Float64frombits(DecodeBEuint64(b))
}

// EncodeBEfloat64 put uint64 to big endian bytes
func EncodeBEfloat64(b []byte, v float64) {
	EncodeBEuint64(b, math.Float64bits(v))
}

// DecodeBEString get string from big endian bytes
func DecodeBEString(b []byte) (string, uint16) {
	len := DecodeBEuint16(b)
	var str string
	if len > 0 {
		str = string(b[:len])
	}
	return str, len
}

// EncodeBEString put string to big endian bytes
func EncodeBEString(b []byte, s string) {
	EncodeBEuint16(b, uint16(len(s)))
	copy(b, []byte(s))
}

// DecodeBEBytes get bytes from big endian bytes
func DecodeBEBytes(b []byte) ([]byte, uint16) {
	len := DecodeBEuint16(b)
	var buf []byte
	if len > 0 {
		buf = make([]byte, len)
		copy(buf, b)
	}
	return buf, len
}

// EncodeBEBytes put bytes to big endian bytes
func EncodeBEBytes(b []byte, buf []byte) {
	EncodeBEuint16(b, uint16(len(buf)))
	copy(b, buf)
}

///////// little endian //////////////////////////////

// DecodeLEuint8 get uint8 from little endian bytes
func DecodeLEuint8(b []byte) uint8 {
	return uint8(b[0])
}

// EncodeLEuint8 put uint8 to little endian bytes
func EncodeLEuint8(b []byte, v uint8) {
	b[0] = byte(v)
}

// DecodeLEuint16 get uint16 from little endian bytes
func DecodeLEuint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

// EncodeLEuint16 put uint16 to little endian bytes
func EncodeLEuint16(b []byte, v uint16) {
	binary.LittleEndian.PutUint16(b, v)
}

// DecodeLEuint32 get uint32 from little endian bytes
func DecodeLEuint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

// EncodeLEuint32 put uint32 to little endian bytes
func EncodeLEuint32(b []byte, v uint32) {
	binary.LittleEndian.PutUint32(b, v)
}

// DecodeLEuint64 get uint64 from little endian bytes
func DecodeLEuint64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

// EncodeLEuint64 put uint64 to little endian bytes
func EncodeLEuint64(b []byte, v uint64) {
	binary.LittleEndian.PutUint64(b, v)
}

// DecodeLEfloat32 get float32 from little endian bytes
func DecodeLEfloat32(b []byte) float32 {
	return math.Float32frombits(DecodeLEuint32(b))
}

// EncodeLEfloat32 put float32 to little endian bytes
func EncodeLEfloat32(b []byte, v float32) {
	EncodeLEuint32(b, math.Float32bits(v))
}

// DecodeLEfloat64 get float64 from little endian bytes
func DecodeLEfloat64(b []byte) float64 {
	return math.Float64frombits(DecodeLEuint64(b))
}

// EncodeLEfloat64 put float64 to little endian bytes
func EncodeLEfloat64(b []byte, v float64) {
	EncodeLEuint64(b, math.Float64bits(v))
}

// DecodeLEString get string from little endian bytes
func DecodeLEString(b []byte) (string, uint16) {
	len := DecodeLEuint16(b)
	var str string
	if len > 0 {
		str = string(b[:len])
	}
	return str, len
}

// EncodeLEString put string to little endian bytes
func EncodeLEString(b []byte, s string) {
	EncodeLEuint16(b, uint16(len(s)))
	copy(b, []byte(s))
}

// DecodeLEBytes get bytes from little endian bytes
func DecodeLEBytes(b []byte) ([]byte, uint16) {
	len := DecodeLEuint16(b)
	var buf []byte
	if len > 0 {
		buf = make([]byte, len)
		copy(buf, b)
	}
	return buf, len
}

// EncodeLEBytes put bytes to little endian bytes
func EncodeLEBytes(b []byte, buf []byte) {
	EncodeLEuint16(b, uint16(len(buf)))
	copy(b, buf)
}
