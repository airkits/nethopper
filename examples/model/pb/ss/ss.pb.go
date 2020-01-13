// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ss.proto

package ss

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
	return fileDescriptor_dc711c54c28c22d7, []int{0}
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
	return fileDescriptor_dc711c54c28c22d7, []int{1}
}

type SSMessage struct {
	Cmd                  string            `protobuf:"bytes,1,opt,name=cmd,proto3" json:"cmd,omitempty"`
	Uid                  string            `protobuf:"bytes,2,opt,name=uid,proto3" json:"uid,omitempty"`
	MsgType              int32             `protobuf:"varint,3,opt,name=msgType,proto3" json:"msgType,omitempty"`
	Seq                  int32             `protobuf:"varint,4,opt,name=seq,proto3" json:"seq,omitempty"`
	Userdata             int32             `protobuf:"varint,5,opt,name=userdata,proto3" json:"userdata,omitempty"`
	Header               map[string][]byte `protobuf:"bytes,6,rep,name=header,proto3" json:"header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Payload              []byte            `protobuf:"bytes,11,opt,name=payload,proto3" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *SSMessage) Reset()         { *m = SSMessage{} }
func (m *SSMessage) String() string { return proto.CompactTextString(m) }
func (*SSMessage) ProtoMessage()    {}
func (*SSMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_dc711c54c28c22d7, []int{0}
}

func (m *SSMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SSMessage.Unmarshal(m, b)
}
func (m *SSMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SSMessage.Marshal(b, m, deterministic)
}
func (m *SSMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SSMessage.Merge(m, src)
}
func (m *SSMessage) XXX_Size() int {
	return xxx_messageInfo_SSMessage.Size(m)
}
func (m *SSMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_SSMessage.DiscardUnknown(m)
}

var xxx_messageInfo_SSMessage proto.InternalMessageInfo

func (m *SSMessage) GetCmd() string {
	if m != nil {
		return m.Cmd
	}
	return ""
}

func (m *SSMessage) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *SSMessage) GetMsgType() int32 {
	if m != nil {
		return m.MsgType
	}
	return 0
}

func (m *SSMessage) GetSeq() int32 {
	if m != nil {
		return m.Seq
	}
	return 0
}

func (m *SSMessage) GetUserdata() int32 {
	if m != nil {
		return m.Userdata
	}
	return 0
}

func (m *SSMessage) GetHeader() map[string][]byte {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *SSMessage) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
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
	return fileDescriptor_dc711c54c28c22d7, []int{1}
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
	return fileDescriptor_dc711c54c28c22d7, []int{2}
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
	return fileDescriptor_dc711c54c28c22d7, []int{3}
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
	proto.RegisterEnum("ss.MessageCmd", MessageCmd_name, MessageCmd_value)
	proto.RegisterEnum("ss.ErrorCode", ErrorCode_name, ErrorCode_value)
	proto.RegisterType((*SSMessage)(nil), "ss.SSMessage")
	proto.RegisterMapType((map[string][]byte)(nil), "ss.SSMessage.HeaderEntry")
	proto.RegisterType((*LoginReq)(nil), "ss.LoginReq")
	proto.RegisterType((*Result)(nil), "ss.Result")
	proto.RegisterType((*LoginResp)(nil), "ss.LoginResp")
}

func init() { proto.RegisterFile("ss.proto", fileDescriptor_dc711c54c28c22d7) }

var fileDescriptor_dc711c54c28c22d7 = []byte{
	// 377 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0xdd, 0xef, 0x93, 0x30,
	0x14, 0x5d, 0x61, 0x20, 0x5c, 0x7e, 0x26, 0xa4, 0x1a, 0x53, 0xf7, 0x84, 0x3c, 0x91, 0x3d, 0xa0,
	0x4e, 0x63, 0xd4, 0xd7, 0x65, 0x89, 0x89, 0xf3, 0x23, 0xdd, 0xfc, 0x03, 0xea, 0xda, 0xe0, 0xb2,
	0xf1, 0xb1, 0x5e, 0xd0, 0xf0, 0xd7, 0x6b, 0xda, 0xc1, 0x3e, 0x12, 0xdf, 0xee, 0x39, 0xed, 0xe9,
	0x3d, 0xe7, 0x00, 0x04, 0x88, 0x79, 0xa3, 0xeb, 0xb6, 0xa6, 0x0e, 0x62, 0xfa, 0x97, 0x40, 0xb8,
	0xd9, 0x7c, 0x51, 0x88, 0xa2, 0x50, 0x34, 0x06, 0x77, 0x57, 0x4a, 0x46, 0x12, 0x92, 0x85, 0xdc,
	0x8c, 0x86, 0xe9, 0xf6, 0x92, 0x39, 0x67, 0xa6, 0xdb, 0x4b, 0xca, 0xe0, 0x51, 0x89, 0xc5, 0xb6,
	0x6f, 0x14, 0x73, 0x13, 0x92, 0x79, 0x7c, 0x84, 0xe6, 0x2e, 0xaa, 0x13, 0x9b, 0x5a, 0xd6, 0x8c,
	0x74, 0x06, 0x41, 0x87, 0x4a, 0x4b, 0xd1, 0x0a, 0xe6, 0x59, 0xfa, 0x82, 0xe9, 0x6b, 0xf0, 0x7f,
	0x29, 0x21, 0x95, 0x66, 0x7e, 0xe2, 0x66, 0xd1, 0xe2, 0x79, 0x8e, 0x98, 0x5f, 0xac, 0xe4, 0x9f,
	0xec, 0xd9, 0xaa, 0x6a, 0x75, 0xcf, 0x87, 0x8b, 0x66, 0x75, 0x23, 0xfa, 0x63, 0x2d, 0x24, 0x8b,
	0x12, 0x92, 0x3d, 0xf0, 0x11, 0xce, 0x3e, 0x40, 0x74, 0x23, 0x30, 0x4e, 0x0e, 0xaa, 0x1f, 0x73,
	0x1c, 0x54, 0x4f, 0x9f, 0x82, 0xf7, 0x5b, 0x1c, 0x3b, 0x65, 0x93, 0x3c, 0xf0, 0x33, 0xf8, 0xe8,
	0xbc, 0x27, 0xe9, 0x5b, 0x08, 0xd6, 0x75, 0xb1, 0xaf, 0xb8, 0x3a, 0x8d, 0x69, 0xc9, 0x35, 0xed,
	0x33, 0xf0, 0x1b, 0x81, 0xf8, 0x67, 0xac, 0x60, 0x40, 0x69, 0x0e, 0x3e, 0x57, 0xd8, 0x1d, 0x5b,
	0x4a, 0x61, 0xba, 0xab, 0xa5, 0xb2, 0x22, 0x8f, 0xdb, 0xd9, 0xbc, 0x53, 0x62, 0x31, 0xb6, 0x56,
	0x62, 0x91, 0x96, 0x10, 0x0e, 0x5b, 0xb0, 0xa1, 0x29, 0xf8, 0xda, 0x8a, 0xad, 0x28, 0x5a, 0x80,
	0x89, 0x7e, 0x7e, 0x8e, 0x0f, 0x27, 0xff, 0x29, 0xfe, 0x6a, 0xc5, 0xbd, 0xb5, 0x62, 0x0c, 0x54,
	0xa2, 0x54, 0xb6, 0xf7, 0x90, 0xdb, 0x79, 0xfe, 0x02, 0x60, 0x28, 0x72, 0x59, 0x4a, 0x1a, 0xc0,
	0xf4, 0xeb, 0x8f, 0xf5, 0x3a, 0x9e, 0xd0, 0x10, 0x3c, 0x6b, 0x23, 0x26, 0xf3, 0x27, 0x10, 0xae,
	0xb4, 0xae, 0xf5, 0xd2, 0x18, 0xf6, 0xc1, 0xf9, 0xf6, 0x39, 0x9e, 0x2c, 0xde, 0x81, 0xcb, 0xbf,
	0x2f, 0xe9, 0x4b, 0x08, 0xb7, 0x5a, 0x54, 0xd8, 0xd4, 0xba, 0xa5, 0x8f, 0xef, 0x3e, 0xcc, 0xec,
	0x1e, 0xa6, 0x93, 0x8c, 0xbc, 0x22, 0x3f, 0x7d, 0xfb, 0x47, 0xbd, 0xf9, 0x17, 0x00, 0x00, 0xff,
	0xff, 0x6f, 0x7d, 0x1c, 0x35, 0x5d, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RPCClient is the client API for RPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RPCClient interface {
	Transport(ctx context.Context, opts ...grpc.CallOption) (RPC_TransportClient, error)
}

type rPCClient struct {
	cc *grpc.ClientConn
}

func NewRPCClient(cc *grpc.ClientConn) RPCClient {
	return &rPCClient{cc}
}

func (c *rPCClient) Transport(ctx context.Context, opts ...grpc.CallOption) (RPC_TransportClient, error) {
	stream, err := c.cc.NewStream(ctx, &_RPC_serviceDesc.Streams[0], "/ss.RPC/Transport", opts...)
	if err != nil {
		return nil, err
	}
	x := &rPCTransportClient{stream}
	return x, nil
}

type RPC_TransportClient interface {
	Send(*SSMessage) error
	Recv() (*SSMessage, error)
	grpc.ClientStream
}

type rPCTransportClient struct {
	grpc.ClientStream
}

func (x *rPCTransportClient) Send(m *SSMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *rPCTransportClient) Recv() (*SSMessage, error) {
	m := new(SSMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RPCServer is the server API for RPC service.
type RPCServer interface {
	Transport(RPC_TransportServer) error
}

// UnimplementedRPCServer can be embedded to have forward compatible implementations.
type UnimplementedRPCServer struct {
}

func (*UnimplementedRPCServer) Transport(srv RPC_TransportServer) error {
	return status.Errorf(codes.Unimplemented, "method Transport not implemented")
}

func RegisterRPCServer(s *grpc.Server, srv RPCServer) {
	s.RegisterService(&_RPC_serviceDesc, srv)
}

func _RPC_Transport_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RPCServer).Transport(&rPCTransportServer{stream})
}

type RPC_TransportServer interface {
	Send(*SSMessage) error
	Recv() (*SSMessage, error)
	grpc.ServerStream
}

type rPCTransportServer struct {
	grpc.ServerStream
}

func (x *rPCTransportServer) Send(m *SSMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *rPCTransportServer) Recv() (*SSMessage, error) {
	m := new(SSMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _RPC_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ss.RPC",
	HandlerType: (*RPCServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Transport",
			Handler:       _RPC_Transport_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "ss.proto",
}
