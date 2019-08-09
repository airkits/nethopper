package binary

import (
	"encoding/binary"
	"math"
)

const initialBufferSize = 1024

// BinaryCodec use binary Write/Read
const (
	// BinaryCodecType type enum
	BinaryCodecType = iota
	// BCuint8 big endian uint8 type
	BCuint8 = 0x01
	// BCuint16 big endian uint16 type
	BCuint16 = 0x02
	// BCuint32 big endian uint32 type
	BCuint32 = 0x03
	// BCuint64 big endian uint64 type
	BCuint64 = 0x04
	// BCstring big endian string type
	BCstring = 0x05
	// BCbytes big endian bytes type
	BCbytes = 0x06
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

// NewCoder create binary coder to Write/Read
func NewCoder(buffer []byte, bigEndian bool) *Coder {
	c := &Coder{
		Offset:    0,
		BigEndian: bigEndian,
	}
	if buffer == nil {
		c.Buffer = make([]byte, initialBufferSize)
	} else {
		c.Buffer = buffer
	}

	return c
}

// Coder Write/Read codec
type Coder struct {
	Buffer    []byte //raw buffer
	Offset    uint32 //Write/Read offset
	BigEndian bool
}

func (c *Coder) ensureCapacity(n uint32) {
	length := uint32(len(c.Buffer))
	if c.Offset+n <= length {
		return
	}
	oldBuffer := c.Buffer
	for c.Offset+n > length {
		length = length * 2
	}
	c.Buffer = make([]byte, length)
	copy(c.Buffer, oldBuffer)
}

// RawData get current pos Write raw data
func (c *Coder) RawData() []byte {
	return c.Buffer[0:c.Offset]
}

// ByteSlice get slice from buffer
func (c *Coder) ByteSlice(start, end uint32) []byte {
	return c.Buffer[start:end]
}

// Reset the codec to init status
func (c *Coder) Reset() {
	c.Offset = 0
}

// Pos get the buffer offset position
func (c *Coder) Pos() uint32 {
	return c.Offset
}

// SkipUint8 skip uint8
func (c *Coder) SkipUint8() {
	c.Offset++
}

// SkipUint16 skip uint16
func (c *Coder) SkipUint16() {
	c.Offset += 2
}

// SkipUint32 skip uint32
func (c *Coder) SkipUint32() {
	c.Offset += 4
}

// SkipUint64 skip uint64
func (c *Coder) SkipUint64() {
	c.Offset += 8
}

// SkipFloat32 skip float32
func (c *Coder) SkipFloat32() {
	c.Offset += 4
}

// SkipFloat64 skip float64
func (c *Coder) SkipFloat64() {
	c.Offset += 8
}

// SkipString skip string length
func (c *Coder) SkipString() {
	size := uint32(c.ReadUint16())
	if size == 0 {
		c.Offset++
	} else {
		c.Offset += size
	}
}

// SkipRaw skip raw buffer length
func (c *Coder) SkipRaw() {
	size := uint32(c.ReadUint16())
	c.Offset += size
}

// ReadUint8 get uint8 value
func (c *Coder) ReadUint8() uint8 {
	v := uint8(c.Buffer[c.Offset])
	c.Offset++
	return v
}

// WriteUint8 set uint8 to byte
func (c *Coder) WriteUint8(v uint8) {
	c.ensureCapacity(1)
	c.Buffer[c.Offset] = byte(v)
	c.Offset++
}

// ReadUint16 get uint16 value
func (c *Coder) ReadUint16() uint16 {
	var v uint16
	if c.BigEndian {
		v = binary.BigEndian.Uint16(c.Buffer[c.Offset:])
	} else {
		v = binary.LittleEndian.Uint16(c.Buffer[c.Offset:])
	}
	c.Offset += 2
	return v
}

// WriteUint16 set uint16 to byte
func (c *Coder) WriteUint16(v uint16) {
	c.ensureCapacity(2)
	if c.BigEndian {
		binary.BigEndian.PutUint16(c.Buffer[c.Offset:], v)
	} else {
		binary.LittleEndian.PutUint16(c.Buffer[c.Offset:], v)
	}
	c.Offset += 2
}

// ReadUint32 get uint32 value
func (c *Coder) ReadUint32() uint32 {
	var v uint32
	if c.BigEndian {
		v = binary.BigEndian.Uint32(c.Buffer[c.Offset:])
	} else {
		v = binary.LittleEndian.Uint32(c.Buffer[c.Offset:])
	}

	c.Offset += 4
	return v
}

// WriteUint32 set uint32 to byte
func (c *Coder) WriteUint32(v uint32) {
	c.ensureCapacity(4)
	if c.BigEndian {
		binary.BigEndian.PutUint32(c.Buffer[c.Offset:], v)
	} else {
		binary.LittleEndian.PutUint32(c.Buffer[c.Offset:], v)
	}
	c.Offset += 4
}

// ReadUint64 get uint64 value
func (c *Coder) ReadUint64() uint64 {
	var v uint64
	if c.BigEndian {
		v = binary.BigEndian.Uint64(c.Buffer[c.Offset:])
	} else {
		v = binary.LittleEndian.Uint64(c.Buffer[c.Offset:])
	}
	c.Offset += 8
	return v
}

// WriteUint64 set uint64 to byte
func (c *Coder) WriteUint64(v uint64) {
	c.ensureCapacity(8)
	if c.BigEndian {
		binary.BigEndian.PutUint64(c.Buffer[c.Offset:], v)
	} else {
		binary.LittleEndian.PutUint64(c.Buffer[c.Offset:], v)
	}
	c.Offset += 8
}

// ReadFloat32 get float32 value
func (c *Coder) ReadFloat32() float32 {
	var v float32
	if c.BigEndian {
		v = math.Float32frombits(binary.BigEndian.Uint32(c.Buffer[c.Offset:]))
	} else {
		v = math.Float32frombits(binary.LittleEndian.Uint32(c.Buffer[c.Offset:]))
	}
	c.Offset += 4
	return v
}

// WriteFloat32 set float32 to byte
func (c *Coder) WriteFloat32(v float32) {
	c.ensureCapacity(4)
	if c.BigEndian {
		binary.BigEndian.PutUint32(c.Buffer[c.Offset:], math.Float32bits(v))
	} else {
		binary.LittleEndian.PutUint32(c.Buffer[c.Offset:], math.Float32bits(v))
	}
	c.Offset += 4
}

// ReadFloat64 get float64 value
func (c *Coder) ReadFloat64() float64 {
	var v float64
	if c.BigEndian {
		v = math.Float64frombits(binary.BigEndian.Uint64(c.Buffer[c.Offset:]))
	} else {
		v = math.Float64frombits(binary.LittleEndian.Uint64(c.Buffer[c.Offset:]))
	}
	c.Offset += 8
	return v
}

// WriteFloat64 set float64 to byte
func (c *Coder) WriteFloat64(v float64) {
	c.ensureCapacity(8)
	if c.BigEndian {
		binary.BigEndian.PutUint64(c.Buffer[c.Offset:], math.Float64bits(v))
	} else {
		binary.LittleEndian.PutUint64(c.Buffer[c.Offset:], math.Float64bits(v))
	}
	c.Offset += 8
}

// ReadString get string value
// Use a zero byte to represent an empty string.
func (c *Coder) ReadString() string {
	size := uint32(c.ReadUint16())
	if size == 0 {
		c.Offset++
		return ""
	}
	v := string(c.Buffer[c.Offset : c.Offset+size])
	c.Offset += size
	return v
}

// WriteString set string to bytes
func (c *Coder) WriteString(s string) {
	size := uint32(len(s))
	c.WriteUint16(uint16(size))
	if size == 0 {
		c.WriteUint8(0x00)
	} else {
		c.ensureCapacity(size)
		copy(c.Buffer[c.Offset:], []byte(s))
		c.Offset += size
	}
}

// ReadRaw use uint16 length and raw bytes
func (c *Coder) ReadRaw() []byte {
	size := uint32(c.ReadUint16())
	if size == 0 {
		return []byte{}
	}
	b := make([]byte, size)
	copy(b, c.Buffer[c.Offset:c.Offset+size])
	c.Offset += size
	return b
}

// WriteRaw use uint16 length and copy raw bytes to buffer
func (c *Coder) WriteRaw(b []byte) {
	size := uint32(len(b))
	c.ensureCapacity(uint32(size))
	copy(c.Buffer[c.Offset:], b)
	c.Offset += size
}

///// seek function ////////////////////////

// SeekReadUint8 get uint8 value
func (c *Coder) SeekReadUint8(offset uint32) uint8 {
	return uint8(c.Buffer[offset])
}

// SeekWriteUint8 set uint8 to byte
func (c *Coder) SeekWriteUint8(offset uint32, v uint8) {
	c.Buffer[offset] = byte(v)
}

// SeekReadUint16 get uint16 value
func (c *Coder) SeekReadUint16(offset uint32) uint16 {
	var v uint16
	if c.BigEndian {
		v = binary.BigEndian.Uint16(c.Buffer[offset:])
	} else {
		v = binary.LittleEndian.Uint16(c.Buffer[offset:])
	}
	return v
}

// SeekWriteUint16 set uint16 to byte
func (c *Coder) SeekWriteUint16(offset uint32, v uint16) {
	if c.BigEndian {
		binary.BigEndian.PutUint16(c.Buffer[offset:], v)
	} else {
		binary.LittleEndian.PutUint16(c.Buffer[offset:], v)
	}
}

// SeekReadUint32 get uint32 value
func (c *Coder) SeekReadUint32(offset uint32) uint32 {
	var v uint32
	if c.BigEndian {
		v = binary.BigEndian.Uint32(c.Buffer[offset:])
	} else {
		v = binary.LittleEndian.Uint32(c.Buffer[offset:])
	}
	return v
}

// SeekWriteUint32 set uint32 to byte
func (c *Coder) SeekWriteUint32(offset uint32, v uint32) {
	if c.BigEndian {
		binary.BigEndian.PutUint32(c.Buffer[offset:], v)
	} else {
		binary.LittleEndian.PutUint32(c.Buffer[offset:], v)
	}
}

// SeekReadUint64 get uint64 value
func (c *Coder) SeekReadUint64(offset uint32) uint64 {
	var v uint64
	if c.BigEndian {
		v = binary.BigEndian.Uint64(c.Buffer[offset:])
	} else {
		v = binary.LittleEndian.Uint64(c.Buffer[offset:])
	}
	return v
}

// SeekWriteUint64 set uint64 to byte
func (c *Coder) SeekWriteUint64(offset uint32, v uint64) {
	if c.BigEndian {
		binary.BigEndian.PutUint64(c.Buffer[offset:], v)
	} else {
		binary.LittleEndian.PutUint64(c.Buffer[offset:], v)
	}
}

// SeekReadFloat32 get float32 value
func (c *Coder) SeekReadFloat32(offset uint32) float32 {
	var v float32
	if c.BigEndian {
		v = math.Float32frombits(binary.BigEndian.Uint32(c.Buffer[offset:]))
	} else {
		v = math.Float32frombits(binary.LittleEndian.Uint32(c.Buffer[offset:]))
	}
	return v
}

// SeekWriteFloat32 set float32 to byte
func (c *Coder) SeekWriteFloat32(offset uint32, v float32) {
	if c.BigEndian {
		binary.BigEndian.PutUint32(c.Buffer[offset:], math.Float32bits(v))
	} else {
		binary.LittleEndian.PutUint32(c.Buffer[offset:], math.Float32bits(v))
	}

}

// SeekReadFloat64 get float64 value
func (c *Coder) SeekReadFloat64(offset uint32) float64 {
	var v float64
	if c.BigEndian {
		v = math.Float64frombits(binary.BigEndian.Uint64(c.Buffer[offset:]))
	} else {
		v = math.Float64frombits(binary.LittleEndian.Uint64(c.Buffer[offset:]))
	}
	return v
}

// SeekWriteFloat64 set float64 to byte
func (c *Coder) SeekWriteFloat64(offset uint32, v float64) {
	if c.BigEndian {
		binary.BigEndian.PutUint64(c.Buffer[offset:], math.Float64bits(v))
	} else {
		binary.LittleEndian.PutUint64(c.Buffer[offset:], math.Float64bits(v))
	}
}

// SeekReadString get string value
// Use a zero byte to represent an empty string.
func (c *Coder) SeekReadString(offset uint32) string {
	size := uint32(c.SeekReadUint16(offset))
	if size == 0 {
		return ""
	}
	offset += 2
	v := string(c.Buffer[offset : offset+size])
	return v
}

// SeekReadRaw use uint16 length and raw bytes
func (c *Coder) SeekReadRaw(offset uint32) []byte {
	size := uint32(c.SeekReadUint16(offset))
	if size == 0 {
		return []byte{}
	}
	b := make([]byte, size)
	offset += 2
	copy(b, c.Buffer[offset:offset+size])
	return b
}
