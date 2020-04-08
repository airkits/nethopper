// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ss.proto

package ss

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
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

type Message struct {
	ID                   int32             `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Cmd                  string            `protobuf:"bytes,2,opt,name=cmd,proto3" json:"cmd,omitempty"`
	MsgType              int32             `protobuf:"varint,3,opt,name=msgType,proto3" json:"msgType,omitempty"`
	Options              map[string][]byte `protobuf:"bytes,4,rep,name=options,proto3" json:"options,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body                 *any.Any          `protobuf:"bytes,11,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_dc711c54c28c22d7, []int{0}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetID() int32 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Message) GetCmd() string {
	if m != nil {
		return m.Cmd
	}
	return ""
}

func (m *Message) GetMsgType() int32 {
	if m != nil {
		return m.MsgType
	}
	return 0
}

func (m *Message) GetOptions() map[string][]byte {
	if m != nil {
		return m.Options
	}
	return nil
}

func (m *Message) GetBody() *any.Any {
	if m != nil {
		return m.Body
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "ss.Message")
	proto.RegisterMapType((map[string][]byte)(nil), "ss.Message.OptionsEntry")
}

func init() { proto.RegisterFile("ss.proto", fileDescriptor_dc711c54c28c22d7) }

var fileDescriptor_dc711c54c28c22d7 = []byte{
	// 251 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0xce, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x06, 0x60, 0x37, 0x69, 0x8d, 0x99, 0x14, 0x91, 0xa1, 0x87, 0x35, 0xa7, 0xd0, 0xd3, 0x82,
	0xb0, 0x95, 0xf5, 0x22, 0xbd, 0x89, 0xf5, 0xd0, 0x83, 0x28, 0x4b, 0x5f, 0x20, 0xb1, 0x6b, 0x10,
	0xdb, 0xdd, 0x90, 0x49, 0x85, 0x7d, 0x5a, 0x5f, 0x45, 0xb2, 0x31, 0x90, 0xdb, 0xfc, 0x33, 0x3f,
	0xc3, 0x07, 0x57, 0x44, 0xb2, 0x69, 0x5d, 0xe7, 0x30, 0x22, 0xca, 0x6f, 0x6b, 0xe7, 0xea, 0xa3,
	0x59, 0x87, 0x4d, 0x75, 0xfe, 0x5c, 0x97, 0xd6, 0x0f, 0xe7, 0xd5, 0x2f, 0x83, 0xe4, 0xd5, 0x10,
	0x95, 0xb5, 0xc1, 0x6b, 0x88, 0x76, 0x5b, 0xce, 0x0a, 0x26, 0xe6, 0x3a, 0xda, 0x6d, 0xf1, 0x06,
	0xe2, 0x8f, 0xd3, 0x81, 0x47, 0x05, 0x13, 0xa9, 0xee, 0x47, 0xe4, 0x90, 0x9c, 0xa8, 0xde, 0xfb,
	0xc6, 0xf0, 0x38, 0xd4, 0xc6, 0x88, 0x0a, 0x12, 0xd7, 0x74, 0x5f, 0xce, 0x12, 0x9f, 0x15, 0xb1,
	0xc8, 0x14, 0x97, 0x44, 0xf2, 0xff, 0xb3, 0x7c, 0x1b, 0x4e, 0x2f, 0xb6, 0x6b, 0xbd, 0x1e, 0x8b,
	0x28, 0x60, 0x56, 0xb9, 0x83, 0xe7, 0x59, 0xc1, 0x44, 0xa6, 0x96, 0x72, 0x50, 0xca, 0x51, 0x29,
	0x9f, 0xac, 0xd7, 0xa1, 0x91, 0x6f, 0x60, 0x31, 0x7d, 0xd1, 0xcb, 0xbe, 0x8d, 0x0f, 0xd4, 0x54,
	0xf7, 0x23, 0x2e, 0x61, 0xfe, 0x53, 0x1e, 0xcf, 0x26, 0x68, 0x17, 0x7a, 0x08, 0x9b, 0xe8, 0x91,
	0x29, 0x05, 0xb1, 0x7e, 0x7f, 0xc6, 0x3b, 0x48, 0xf7, 0x6d, 0x69, 0xa9, 0x71, 0x6d, 0x87, 0xd9,
	0x04, 0x97, 0x4f, 0xc3, 0xea, 0x42, 0xb0, 0x7b, 0x56, 0x5d, 0x06, 0xc3, 0xc3, 0x5f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x43, 0xa6, 0x2e, 0x8e, 0x47, 0x01, 0x00, 0x00,
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
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type rPCTransportClient struct {
	grpc.ClientStream
}

func (x *rPCTransportClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *rPCTransportClient) Recv() (*Message, error) {
	m := new(Message)
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
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type rPCTransportServer struct {
	grpc.ServerStream
}

func (x *rPCTransportServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *rPCTransportServer) Recv() (*Message, error) {
	m := new(Message)
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
