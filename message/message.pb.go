// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package message

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _struct "github.com/golang/protobuf/ptypes/struct"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Message_Priority int32

const (
	Message_NORMAL Message_Priority = 0
	Message_HIGH   Message_Priority = 1
	Message_URGENT Message_Priority = 2
)

var Message_Priority_name = map[int32]string{
	0: "NORMAL",
	1: "HIGH",
	2: "URGENT",
}
var Message_Priority_value = map[string]int32{
	"NORMAL": 0,
	"HIGH":   1,
	"URGENT": 2,
}

func (x Message_Priority) String() string {
	return proto.EnumName(Message_Priority_name, int32(x))
}
func (Message_Priority) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_message_ac3706dc675b70c7, []int{0, 0}
}

type Message struct {
	Priority             Message_Priority     `protobuf:"varint,1,opt,name=priority,proto3,enum=message.Message_Priority" json:"priority,omitempty"`
	Received             *timestamp.Timestamp `protobuf:"bytes,2,opt,name=received,proto3" json:"received,omitempty"`
	Content              *_struct.Struct      `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_message_ac3706dc675b70c7, []int{0}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (dst *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(dst, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetPriority() Message_Priority {
	if m != nil {
		return m.Priority
	}
	return Message_NORMAL
}

func (m *Message) GetReceived() *timestamp.Timestamp {
	if m != nil {
		return m.Received
	}
	return nil
}

func (m *Message) GetContent() *_struct.Struct {
	if m != nil {
		return m.Content
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "message.Message")
	proto.RegisterEnum("message.Message_Priority", Message_Priority_name, Message_Priority_value)
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_message_ac3706dc675b70c7) }

var fileDescriptor_message_ac3706dc675b70c7 = []byte{
	// 228 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x4d, 0x2d, 0x2e,
	0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x72, 0xa5, 0xe4, 0xd3,
	0xf3, 0xf3, 0xd3, 0x73, 0x52, 0xf5, 0xc1, 0xc2, 0x49, 0xa5, 0x69, 0xfa, 0x25, 0x99, 0xb9, 0xa9,
	0xc5, 0x25, 0x89, 0xb9, 0x05, 0x10, 0x95, 0x52, 0x32, 0xe8, 0x0a, 0x8a, 0x4b, 0x8a, 0x4a, 0x93,
	0x4b, 0x20, 0xb2, 0x4a, 0x37, 0x19, 0xb9, 0xd8, 0x7d, 0x21, 0x46, 0x09, 0x99, 0x72, 0x71, 0x14,
	0x14, 0x65, 0xe6, 0x17, 0x65, 0x96, 0x54, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0xf0, 0x19, 0x49, 0xea,
	0xc1, 0x6c, 0x85, 0xaa, 0xd1, 0x0b, 0x80, 0x2a, 0x08, 0x82, 0x2b, 0x15, 0x32, 0xe3, 0xe2, 0x28,
	0x4a, 0x4d, 0x4e, 0xcd, 0x2c, 0x4b, 0x4d, 0x91, 0x60, 0x52, 0x60, 0xd4, 0xe0, 0x36, 0x92, 0xd2,
	0x83, 0xd8, 0xa9, 0x07, 0xb3, 0x53, 0x2f, 0x04, 0xe6, 0xa8, 0x20, 0xb8, 0x5a, 0x21, 0x43, 0x2e,
	0xf6, 0xe4, 0xfc, 0xbc, 0x92, 0xd4, 0xbc, 0x12, 0x09, 0x66, 0xb0, 0x36, 0x71, 0x0c, 0x6d, 0xc1,
	0x60, 0xa7, 0x06, 0xc1, 0xd4, 0x29, 0xe9, 0x70, 0x71, 0xc0, 0x1c, 0x20, 0xc4, 0xc5, 0xc5, 0xe6,
	0xe7, 0x1f, 0xe4, 0xeb, 0xe8, 0x23, 0xc0, 0x20, 0xc4, 0xc1, 0xc5, 0xe2, 0xe1, 0xe9, 0xee, 0x21,
	0xc0, 0x08, 0x12, 0x0d, 0x0d, 0x72, 0x77, 0xf5, 0x0b, 0x11, 0x60, 0x72, 0xe2, 0x8c, 0x82, 0x85,
	0x52, 0x12, 0x1b, 0xd8, 0x48, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xfa, 0x11, 0xa3, 0x41,
	0x46, 0x01, 0x00, 0x00,
}
