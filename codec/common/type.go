package common

const (
	// CodecType codec type define,json/pb/gob/custom binary
	CodecType = iota
	//CodecTypeJSON type json
	CodecTypeJSON
	//CodecTypePB type protobuf
	CodecTypePB
	//CodecTypeGOB type gob
	CodecTypeGOB
	//CodecTypeBinary type binary
	CodecTypeBinary
)
