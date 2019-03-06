// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package message

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_struct "github.com/golang/protobuf/ptypes/struct"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type Message_Direction int32

const (
	Message_PASS   Message_Direction = 0
	Message_ACCEPT Message_Direction = 1
	Message_REJECT Message_Direction = 2
)

var Message_Direction_name = map[int32]string{
	0: "PASS",
	1: "ACCEPT",
	2: "REJECT",
}

var Message_Direction_value = map[string]int32{
	"PASS":   0,
	"ACCEPT": 1,
	"REJECT": 2,
}

func (x Message_Direction) String() string {
	return proto.EnumName(Message_Direction_name, int32(x))
}

func (Message_Direction) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0, 0}
}

type Message struct {
	Direction            Message_Direction    `protobuf:"varint,1,opt,name=direction,proto3,enum=Message_Direction" json:"direction,omitempty"`
	Received             *timestamp.Timestamp `protobuf:"bytes,2,opt,name=received,proto3" json:"received,omitempty"`
	Content              *_struct.Value       `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
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

func (m *Message) GetDirection() Message_Direction {
	if m != nil {
		return m.Direction
	}
	return Message_PASS
}

func (m *Message) GetReceived() *timestamp.Timestamp {
	if m != nil {
		return m.Received
	}
	return nil
}

func (m *Message) GetContent() *_struct.Value {
	if m != nil {
		return m.Content
	}
	return nil
}

func init() {
	proto.RegisterEnum("Message_Direction", Message_Direction_name, Message_Direction_value)
	proto.RegisterType((*Message)(nil), "Message")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 218 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x4d, 0x2d, 0x2e,
	0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x97, 0x92, 0x4f, 0xcf, 0xcf, 0x4f, 0xcf,
	0x49, 0xd5, 0x07, 0xf3, 0x92, 0x4a, 0xd3, 0xf4, 0x4b, 0x32, 0x73, 0x53, 0x8b, 0x4b, 0x12, 0x73,
	0x0b, 0xa0, 0x0a, 0x64, 0xd0, 0x15, 0x14, 0x97, 0x14, 0x95, 0x26, 0x97, 0x40, 0x64, 0x95, 0xae,
	0x30, 0x72, 0xb1, 0xfb, 0x42, 0x0c, 0x14, 0x32, 0xe0, 0xe2, 0x4c, 0xc9, 0x2c, 0x4a, 0x4d, 0x2e,
	0xc9, 0xcc, 0xcf, 0x93, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x33, 0x12, 0xd2, 0x83, 0x4a, 0xea, 0xb9,
	0xc0, 0x64, 0x82, 0x10, 0x8a, 0x84, 0xcc, 0xb8, 0x38, 0x8a, 0x52, 0x93, 0x53, 0x33, 0xcb, 0x52,
	0x53, 0x24, 0x98, 0x14, 0x18, 0x35, 0xb8, 0x8d, 0xa4, 0xf4, 0x20, 0xd6, 0xe9, 0xc1, 0xac, 0xd3,
	0x0b, 0x81, 0xb9, 0x27, 0x08, 0xae, 0x56, 0xc8, 0x80, 0x8b, 0x3d, 0x39, 0x3f, 0xaf, 0x24, 0x35,
	0xaf, 0x44, 0x82, 0x19, 0xac, 0x4d, 0x0c, 0x43, 0x5b, 0x58, 0x62, 0x4e, 0x69, 0x6a, 0x10, 0x4c,
	0x99, 0x92, 0x2e, 0x17, 0x27, 0xdc, 0x05, 0x42, 0x1c, 0x5c, 0x2c, 0x01, 0x8e, 0xc1, 0xc1, 0x02,
	0x0c, 0x42, 0x5c, 0x5c, 0x6c, 0x8e, 0xce, 0xce, 0xae, 0x01, 0x21, 0x02, 0x8c, 0x20, 0x76, 0x90,
	0xab, 0x97, 0xab, 0x73, 0x88, 0x00, 0x53, 0x12, 0x1b, 0xd8, 0x1c, 0x63, 0x40, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xb2, 0xf7, 0x90, 0x58, 0x2d, 0x01, 0x00, 0x00,
}
