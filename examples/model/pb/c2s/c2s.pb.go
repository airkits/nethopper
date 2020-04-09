// Code generated by protoc-gen-go. DO NOT EDIT.
// source: c2s.proto

package c2s

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
	return fileDescriptor_d62ebb484b362a6d, []int{0}
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
	return fileDescriptor_d62ebb484b362a6d, []int{1}
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
	return fileDescriptor_d62ebb484b362a6d, []int{0}
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
	return fileDescriptor_d62ebb484b362a6d, []int{1}
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
	return fileDescriptor_d62ebb484b362a6d, []int{2}
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
	proto.RegisterEnum("c2s.MessageCmd", MessageCmd_name, MessageCmd_value)
	proto.RegisterEnum("c2s.ErrorCode", ErrorCode_name, ErrorCode_value)
	proto.RegisterType((*LoginReq)(nil), "c2s.LoginReq")
	proto.RegisterType((*Result)(nil), "c2s.Result")
	proto.RegisterType((*LoginResp)(nil), "c2s.LoginResp")
}

func init() { proto.RegisterFile("c2s.proto", fileDescriptor_d62ebb484b362a6d) }

var fileDescriptor_d62ebb484b362a6d = []byte{
	// 216 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0xcd, 0x4a, 0xc4, 0x30,
	0x10, 0x80, 0xb7, 0x3f, 0x1b, 0x36, 0xb3, 0x97, 0x30, 0x82, 0xf4, 0xa8, 0xf5, 0x22, 0x3d, 0xf4,
	0x50, 0x7d, 0x83, 0xe2, 0xc9, 0xaa, 0x10, 0xf0, 0x01, 0x6a, 0x13, 0x4a, 0xc1, 0x36, 0x35, 0xd3,
	0xe2, 0xeb, 0x4b, 0xc6, 0x08, 0x1e, 0xf6, 0xf6, 0x65, 0xc2, 0x37, 0xf9, 0x08, 0xc8, 0xa1, 0xa1,
	0x7a, 0xf5, 0x6e, 0x73, 0x98, 0x0d, 0x0d, 0x95, 0x8f, 0x70, 0xea, 0xdc, 0x38, 0x2d, 0xda, 0x7e,
	0xa1, 0x82, 0x6c, 0x9f, 0x4c, 0x91, 0xdc, 0x24, 0xf7, 0x52, 0x07, 0xc4, 0x6b, 0x10, 0x6b, 0x4f,
	0xf4, 0x6d, 0x8a, 0x94, 0x87, 0xf1, 0x54, 0xd6, 0x20, 0xb4, 0xa5, 0xfd, 0x73, 0x43, 0x84, 0x7c,
	0x70, 0xc6, 0xb2, 0x74, 0xd4, 0xcc, 0x61, 0xcf, 0x4c, 0x63, 0x54, 0x02, 0x96, 0x0b, 0xc8, 0xf8,
	0x0a, 0xad, 0x78, 0x07, 0xc2, 0xb3, 0xcc, 0xd2, 0xb9, 0x39, 0xd7, 0xa1, 0xe9, 0x77, 0x9f, 0x8e,
	0x57, 0x7f, 0x2d, 0xe9, 0xa5, 0x96, 0xec, 0x7f, 0x4b, 0x28, 0x58, 0xfa, 0xd9, 0x16, 0x39, 0x4f,
	0x99, 0xab, 0x5b, 0x80, 0x17, 0x4b, 0xd4, 0x8f, 0xb6, 0x9d, 0x0d, 0x9e, 0x20, 0x7f, 0x7d, 0xef,
	0x3a, 0x75, 0x40, 0x09, 0x47, 0xee, 0x50, 0x49, 0x75, 0x05, 0xf2, 0xc9, 0x7b, 0xe7, 0xdb, 0x50,
	0x2c, 0x20, 0x7d, 0x7b, 0x56, 0x87, 0x0f, 0xc1, 0x3f, 0xf3, 0xf0, 0x13, 0x00, 0x00, 0xff, 0xff,
	0xcc, 0xbe, 0x41, 0x2e, 0x26, 0x01, 0x00, 0x00,
}