// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pb/chat.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	pb/chat.proto

It has these top-level messages:
	MessagePack
	OutgoingMessage
	IncomingMessage
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MessagePack struct {
	Messages []*IncomingMessage `protobuf:"bytes,1,rep,name=messages" json:"messages,omitempty"`
}

func (m *MessagePack) Reset()                    { *m = MessagePack{} }
func (m *MessagePack) String() string            { return proto.CompactTextString(m) }
func (*MessagePack) ProtoMessage()               {}
func (*MessagePack) Link() ([]byte, []int) { return fileLink0, []int{0} }

func (m *MessagePack) GetMessages() []*IncomingMessage {
	if m != nil {
		return m.Messages
	}
	return nil
}

type OutgoingMessage struct {
	Room      string `protobuf:"bytes,1,opt,name=room" json:"room,omitempty"`
	Subscribe bool   `protobuf:"varint,2,opt,name=subscribe" json:"subscribe,omitempty"`
	Username  string `protobuf:"bytes,3,opt,name=username" json:"username,omitempty"`
	Text      string `protobuf:"bytes,4,opt,name=text" json:"text,omitempty"`
}

func (m *OutgoingMessage) Reset()                    { *m = OutgoingMessage{} }
func (m *OutgoingMessage) String() string            { return proto.CompactTextString(m) }
func (*OutgoingMessage) ProtoMessage()               {}
func (*OutgoingMessage) Link() ([]byte, []int) { return fileLink0, []int{1} }

func (m *OutgoingMessage) GetRoom() string {
	if m != nil {
		return m.Room
	}
	return ""
}

func (m *OutgoingMessage) GetSubscribe() bool {
	if m != nil {
		return m.Subscribe
	}
	return false
}

func (m *OutgoingMessage) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *OutgoingMessage) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

type IncomingMessage struct {
	Room      string `protobuf:"bytes,1,opt,name=room" json:"room,omitempty"`
	Username  string `protobuf:"bytes,2,opt,name=username" json:"username,omitempty"`
	Timestamp int64  `protobuf:"varint,3,opt,name=timestamp" json:"timestamp,omitempty"`
	Text      string `protobuf:"bytes,4,opt,name=text" json:"text,omitempty"`
}

func (m *IncomingMessage) Reset()                    { *m = IncomingMessage{} }
func (m *IncomingMessage) String() string            { return proto.CompactTextString(m) }
func (*IncomingMessage) ProtoMessage()               {}
func (*IncomingMessage) Link() ([]byte, []int) { return fileLink0, []int{2} }

func (m *IncomingMessage) GetRoom() string {
	if m != nil {
		return m.Room
	}
	return ""
}

func (m *IncomingMessage) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *IncomingMessage) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *IncomingMessage) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func init() {
	proto.RegisterType((*MessagePack)(nil), "pb.MessagePack")
	proto.RegisterType((*OutgoingMessage)(nil), "pb.OutgoingMessage")
	proto.RegisterType((*IncomingMessage)(nil), "pb.IncomingMessage")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Chat service

type ChatClient interface {
	Subscribe(ctx context.Context, opts ...grpc.CallOption) (Chat_SubscribeClient, error)
}

type chatClient struct {
	cc *grpc.ClientConn
}

func NewChatClient(cc *grpc.ClientConn) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) Subscribe(ctx context.Context, opts ...grpc.CallOption) (Chat_SubscribeClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Chat_serviceDesc.Streams[0], c.cc, "/pb.Chat/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatSubscribeClient{stream}
	return x, nil
}

type Chat_SubscribeClient interface {
	Send(*OutgoingMessage) error
	Recv() (*MessagePack, error)
	grpc.ClientStream
}

type chatSubscribeClient struct {
	grpc.ClientStream
}

func (x *chatSubscribeClient) Send(m *OutgoingMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chatSubscribeClient) Recv() (*MessagePack, error) {
	m := new(MessagePack)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Chat service

type ChatServer interface {
	Subscribe(Chat_SubscribeServer) error
}

func RegisterChatServer(s *grpc.Server, srv ChatServer) {
	s.RegisterService(&_Chat_serviceDesc, srv)
}

func _Chat_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChatServer).Subscribe(&chatSubscribeServer{stream})
}

type Chat_SubscribeServer interface {
	Send(*MessagePack) error
	Recv() (*OutgoingMessage, error)
	grpc.ServerStream
}

type chatSubscribeServer struct {
	grpc.ServerStream
}

func (x *chatSubscribeServer) Send(m *MessagePack) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chatSubscribeServer) Recv() (*OutgoingMessage, error) {
	m := new(OutgoingMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Chat_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Chat",
	HandlerType: (*ChatServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _Chat_Subscribe_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pb/chat.proto",
}

func init() { proto.RegisterFile("pb/chat.proto", fileLink0) }

var fileLink0 = []byte{
	// 237 bytes of a gzipped FileLinkProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x31, 0x6b, 0xc3, 0x30,
	0x10, 0x85, 0x2b, 0xdb, 0x14, 0xfb, 0x42, 0x09, 0xa8, 0x8b, 0x08, 0x19, 0x8c, 0x27, 0x4f, 0x4e,
	0x49, 0x87, 0x6e, 0xed, 0xd0, 0xa9, 0x43, 0x69, 0x71, 0x7f, 0x81, 0x64, 0x84, 0x63, 0x8a, 0x2c,
	0xa1, 0x3b, 0x43, 0x7f, 0x7e, 0x91, 0x12, 0xec, 0x60, 0xb2, 0xdd, 0x3d, 0x9d, 0xde, 0x7b, 0x7c,
	0xf0, 0xe0, 0xd4, 0xa1, 0x3b, 0x49, 0x6a, 0x9c, 0xb7, 0x64, 0x79, 0xe2, 0x54, 0xf5, 0x0a, 0x9b,
	0x4f, 0x8d, 0x28, 0x7b, 0xfd, 0x2d, 0xbb, 0x5f, 0x7e, 0x80, 0xdc, 0x9c, 0x57, 0x14, 0xac, 0x4c,
	0xeb, 0xcd, 0xf1, 0xb1, 0x71, 0xaa, 0xf9, 0x18, 0x3b, 0x6b, 0x86, 0xb1, 0xbf, 0x9c, 0xb6, 0xf3,
	0x51, 0x85, 0xb0, 0xfd, 0x9a, 0xa8, 0xb7, 0xcb, 0x23, 0xe7, 0x90, 0x79, 0x6b, 0x8d, 0x60, 0x25,
	0xab, 0x8b, 0x36, 0xce, 0x7c, 0x0f, 0x05, 0x4e, 0x0a, 0x3b, 0x3f, 0x28, 0x2d, 0x92, 0x92, 0xd5,
	0x79, 0xbb, 0x08, 0x7c, 0x07, 0xf9, 0x84, 0xda, 0x8f, 0xd2, 0x68, 0x91, 0xc6, 0x5f, 0xf3, 0x1e,
	0xdc, 0x48, 0xff, 0x91, 0xc8, 0xce, 0x6e, 0x61, 0x0e, 0xa1, 0xab, 0x46, 0x37, 0x43, 0xaf, 0x6d,
	0x93, 0x95, 0xed, 0x1e, 0x0a, 0x1a, 0x8c, 0x46, 0x92, 0xc6, 0xc5, 0xcc, 0xb4, 0x5d, 0x84, 0x5b,
	0xa1, 0xc7, 0x37, 0xc8, 0xde, 0x4f, 0x92, 0xf8, 0x0b, 0x14, 0x3f, 0x73, 0xf3, 0x48, 0x67, 0x05,
	0x60, 0xb7, 0x0d, 0xe2, 0x15, 0xd5, 0xea, 0xae, 0x66, 0x4f, 0x4c, 0xdd, 0x47, 0xea, 0xcf, 0xff,
	0x01, 0x00, 0x00, 0xff, 0xff, 0xe1, 0x20, 0x8a, 0x49, 0x86, 0x01, 0x00, 0x00,
}
