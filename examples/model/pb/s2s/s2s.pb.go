// Code generated by protoc-gen-go. DO NOT EDIT.
// source: s2s.proto

package s2s

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// client  <-> server  10000 - 19999
type MessageCmd int32

const (
	MessageCmd_NULL  MessageCmd = 0
	MessageCmd_Login MessageCmd = 1
)

var MessageCmd_name = map[int32]string{
	0: "NULL",
	1: "Login",
}

var MessageCmd_value = map[string]int32{
	"NULL":  0,
	"Login": 1,
}

func (x MessageCmd) String() string {
	return proto.EnumName(MessageCmd_name, int32(x))
}

func (MessageCmd) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_953aa047daafed82, []int{0}
}

type ErrorCode int32

const (
	ErrorCode_OK ErrorCode = 0
)

var ErrorCode_name = map[int32]string{
	0: "OK",
}

var ErrorCode_value = map[string]int32{
	"OK": 0,
}

func (x ErrorCode) String() string {
	return proto.EnumName(ErrorCode_name, int32(x))
}

func (ErrorCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_953aa047daafed82, []int{1}
}

type LoginReq struct {
	Uid                  string   `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Passwd               string   `protobuf:"bytes,2,opt,name=passwd,proto3" json:"passwd,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginReq) Reset()         { *m = LoginReq{} }
func (m *LoginReq) String() string { return proto.CompactTextString(m) }
func (*LoginReq) ProtoMessage()    {}
func (*LoginReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_953aa047daafed82, []int{0}
}

func (m *LoginReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginReq.Unmarshal(m, b)
}
func (m *LoginReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginReq.Marshal(b, m, deterministic)
}
func (m *LoginReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginReq.Merge(m, src)
}
func (m *LoginReq) XXX_Size() int {
	return xxx_messageInfo_LoginReq.Size(m)
}
func (m *LoginReq) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginReq.DiscardUnknown(m)
}

var xxx_messageInfo_LoginReq proto.InternalMessageInfo

func (m *LoginReq) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *LoginReq) GetPasswd() string {
	if m != nil {
		return m.Passwd
	}
	return ""
}

type Result struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Result) Reset()         { *m = Result{} }
func (m *Result) String() string { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()    {}
func (*Result) Descriptor() ([]byte, []int) {
	return fileDescriptor_953aa047daafed82, []int{1}
}

func (m *Result) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Result.Unmarshal(m, b)
}
func (m *Result) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Result.Marshal(b, m, deterministic)
}
func (m *Result) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Result.Merge(m, src)
}
func (m *Result) XXX_Size() int {
	return xxx_messageInfo_Result.Size(m)
}
func (m *Result) XXX_DiscardUnknown() {
	xxx_messageInfo_Result.DiscardUnknown(m)
}

var xxx_messageInfo_Result proto.InternalMessageInfo

func (m *Result) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Result) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type LoginResp struct {
	Result               *Result  `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	Uid                  string   `protobuf:"bytes,2,opt,name=uid,proto3" json:"uid,omitempty"`
	Passwd               string   `protobuf:"bytes,3,opt,name=passwd,proto3" json:"passwd,omitempty"`
	Name                 string   `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginResp) Reset()         { *m = LoginResp{} }
func (m *LoginResp) String() string { return proto.CompactTextString(m) }
func (*LoginResp) ProtoMessage()    {}
func (*LoginResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_953aa047daafed82, []int{2}
}

func (m *LoginResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginResp.Unmarshal(m, b)
}
func (m *LoginResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginResp.Marshal(b, m, deterministic)
}
func (m *LoginResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginResp.Merge(m, src)
}
func (m *LoginResp) XXX_Size() int {
	return xxx_messageInfo_LoginResp.Size(m)
}
func (m *LoginResp) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginResp.DiscardUnknown(m)
}

var xxx_messageInfo_LoginResp proto.InternalMessageInfo

func (m *LoginResp) GetResult() *Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *LoginResp) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *LoginResp) GetPasswd() string {
	if m != nil {
		return m.Passwd
	}
	return ""
}

func (m *LoginResp) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterEnum("s2s.MessageCmd", MessageCmd_name, MessageCmd_value)
	proto.RegisterEnum("s2s.ErrorCode", ErrorCode_name, ErrorCode_value)
	proto.RegisterType((*LoginReq)(nil), "s2s.LoginReq")
	proto.RegisterType((*Result)(nil), "s2s.Result")
	proto.RegisterType((*LoginResp)(nil), "s2s.LoginResp")
}

func init() { proto.RegisterFile("s2s.proto", fileDescriptor_953aa047daafed82) }

var fileDescriptor_953aa047daafed82 = []byte{
	// 224 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x3d, 0x4b, 0x04, 0x31,
	0x10, 0x40, 0x6f, 0xbf, 0xe2, 0x65, 0xae, 0x09, 0x23, 0xc8, 0x96, 0xba, 0x36, 0x72, 0x45, 0x8a,
	0x68, 0x67, 0xe7, 0x61, 0xe5, 0xaa, 0x10, 0xb0, 0xb1, 0x5b, 0x4d, 0x58, 0x0e, 0xdc, 0xcd, 0x9a,
	0xb9, 0xc3, 0xbf, 0x2f, 0x19, 0x23, 0x58, 0xd8, 0xbd, 0x4c, 0x78, 0x93, 0x47, 0x40, 0x92, 0x21,
	0xbd, 0xc4, 0x70, 0x08, 0x58, 0x91, 0xa1, 0xee, 0x06, 0xd6, 0x7d, 0x18, 0xf7, 0xb3, 0xf5, 0x9f,
	0xa8, 0xa0, 0x3a, 0xee, 0x5d, 0x5b, 0x9c, 0x17, 0x57, 0xd2, 0x26, 0xc4, 0x33, 0x10, 0xcb, 0x40,
	0xf4, 0xe5, 0xda, 0x92, 0x87, 0xf9, 0xd4, 0x69, 0x10, 0xd6, 0xd3, 0xf1, 0xe3, 0x80, 0x08, 0xf5,
	0x7b, 0x70, 0x9e, 0xa5, 0xc6, 0x32, 0xa7, 0x3d, 0x13, 0x8d, 0x59, 0x49, 0xd8, 0xcd, 0x20, 0xf3,
	0x2b, 0xb4, 0xe0, 0x25, 0x88, 0xc8, 0x32, 0x4b, 0x1b, 0xb3, 0xd1, 0xa9, 0xe9, 0x67, 0x9f, 0xcd,
	0x57, 0xbf, 0x2d, 0xe5, 0x7f, 0x2d, 0xd5, 0xdf, 0x96, 0x54, 0x30, 0x0f, 0x93, 0x6f, 0x6b, 0x9e,
	0x32, 0x6f, 0x2f, 0x00, 0x1e, 0x3d, 0xd1, 0x30, 0xfa, 0xdd, 0xe4, 0x70, 0x0d, 0xf5, 0xd3, 0x4b,
	0xdf, 0xab, 0x15, 0x4a, 0x68, 0xb8, 0x43, 0x15, 0xdb, 0x53, 0x90, 0xf7, 0x31, 0x86, 0xb8, 0x4b,
	0xc5, 0x02, 0xca, 0xe7, 0x07, 0xb5, 0xba, 0x3b, 0x79, 0x6d, 0xf4, 0x2d, 0x19, 0x7a, 0x13, 0xfc,
	0x45, 0xd7, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xc2, 0x04, 0x10, 0x28, 0x2f, 0x01, 0x00, 0x00,
}
