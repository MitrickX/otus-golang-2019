// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package grpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type Event struct {
	Id                   int32                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string               `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Start                *timestamp.Timestamp `protobuf:"bytes,3,opt,name=start,proto3" json:"start,omitempty"`
	End                  *timestamp.Timestamp `protobuf:"bytes,4,opt,name=end,proto3" json:"end,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Event) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Event) GetStart() *timestamp.Timestamp {
	if m != nil {
		return m.Start
	}
	return nil
}

func (m *Event) GetEnd() *timestamp.Timestamp {
	if m != nil {
		return m.End
	}
	return nil
}

type SimpleResponse struct {
	Result               string   `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SimpleResponse) Reset()         { *m = SimpleResponse{} }
func (m *SimpleResponse) String() string { return proto.CompactTextString(m) }
func (*SimpleResponse) ProtoMessage()    {}
func (*SimpleResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

func (m *SimpleResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SimpleResponse.Unmarshal(m, b)
}
func (m *SimpleResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SimpleResponse.Marshal(b, m, deterministic)
}
func (m *SimpleResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SimpleResponse.Merge(m, src)
}
func (m *SimpleResponse) XXX_Size() int {
	return xxx_messageInfo_SimpleResponse.Size(m)
}
func (m *SimpleResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SimpleResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SimpleResponse proto.InternalMessageInfo

func (m *SimpleResponse) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

type EventListResponse struct {
	Events               []*Event `protobuf:"bytes,1,rep,name=events,proto3" json:"events,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventListResponse) Reset()         { *m = EventListResponse{} }
func (m *EventListResponse) String() string { return proto.CompactTextString(m) }
func (*EventListResponse) ProtoMessage()    {}
func (*EventListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}

func (m *EventListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventListResponse.Unmarshal(m, b)
}
func (m *EventListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventListResponse.Marshal(b, m, deterministic)
}
func (m *EventListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventListResponse.Merge(m, src)
}
func (m *EventListResponse) XXX_Size() int {
	return xxx_messageInfo_EventListResponse.Size(m)
}
func (m *EventListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EventListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EventListResponse proto.InternalMessageInfo

func (m *EventListResponse) GetEvents() []*Event {
	if m != nil {
		return m.Events
	}
	return nil
}

type CreateEventRequest struct {
	Name                 string               `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Start                *timestamp.Timestamp `protobuf:"bytes,2,opt,name=start,proto3" json:"start,omitempty"`
	End                  *timestamp.Timestamp `protobuf:"bytes,3,opt,name=end,proto3" json:"end,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *CreateEventRequest) Reset()         { *m = CreateEventRequest{} }
func (m *CreateEventRequest) String() string { return proto.CompactTextString(m) }
func (*CreateEventRequest) ProtoMessage()    {}
func (*CreateEventRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{3}
}

func (m *CreateEventRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateEventRequest.Unmarshal(m, b)
}
func (m *CreateEventRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateEventRequest.Marshal(b, m, deterministic)
}
func (m *CreateEventRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateEventRequest.Merge(m, src)
}
func (m *CreateEventRequest) XXX_Size() int {
	return xxx_messageInfo_CreateEventRequest.Size(m)
}
func (m *CreateEventRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateEventRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateEventRequest proto.InternalMessageInfo

func (m *CreateEventRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateEventRequest) GetStart() *timestamp.Timestamp {
	if m != nil {
		return m.Start
	}
	return nil
}

func (m *CreateEventRequest) GetEnd() *timestamp.Timestamp {
	if m != nil {
		return m.End
	}
	return nil
}

type UpdateEventRequest struct {
	Id                   int32                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string               `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Start                *timestamp.Timestamp `protobuf:"bytes,3,opt,name=start,proto3" json:"start,omitempty"`
	End                  *timestamp.Timestamp `protobuf:"bytes,4,opt,name=end,proto3" json:"end,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *UpdateEventRequest) Reset()         { *m = UpdateEventRequest{} }
func (m *UpdateEventRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateEventRequest) ProtoMessage()    {}
func (*UpdateEventRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{4}
}

func (m *UpdateEventRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateEventRequest.Unmarshal(m, b)
}
func (m *UpdateEventRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateEventRequest.Marshal(b, m, deterministic)
}
func (m *UpdateEventRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateEventRequest.Merge(m, src)
}
func (m *UpdateEventRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateEventRequest.Size(m)
}
func (m *UpdateEventRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateEventRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateEventRequest proto.InternalMessageInfo

func (m *UpdateEventRequest) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateEventRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *UpdateEventRequest) GetStart() *timestamp.Timestamp {
	if m != nil {
		return m.Start
	}
	return nil
}

func (m *UpdateEventRequest) GetEnd() *timestamp.Timestamp {
	if m != nil {
		return m.End
	}
	return nil
}

type DeleteEventRequest struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteEventRequest) Reset()         { *m = DeleteEventRequest{} }
func (m *DeleteEventRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteEventRequest) ProtoMessage()    {}
func (*DeleteEventRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{5}
}

func (m *DeleteEventRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteEventRequest.Unmarshal(m, b)
}
func (m *DeleteEventRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteEventRequest.Marshal(b, m, deterministic)
}
func (m *DeleteEventRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteEventRequest.Merge(m, src)
}
func (m *DeleteEventRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteEventRequest.Size(m)
}
func (m *DeleteEventRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteEventRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteEventRequest proto.InternalMessageInfo

func (m *DeleteEventRequest) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

type Nothing struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Nothing) Reset()         { *m = Nothing{} }
func (m *Nothing) String() string { return proto.CompactTextString(m) }
func (*Nothing) ProtoMessage()    {}
func (*Nothing) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{6}
}

func (m *Nothing) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Nothing.Unmarshal(m, b)
}
func (m *Nothing) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Nothing.Marshal(b, m, deterministic)
}
func (m *Nothing) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Nothing.Merge(m, src)
}
func (m *Nothing) XXX_Size() int {
	return xxx_messageInfo_Nothing.Size(m)
}
func (m *Nothing) XXX_DiscardUnknown() {
	xxx_messageInfo_Nothing.DiscardUnknown(m)
}

var xxx_messageInfo_Nothing proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Event)(nil), "grpc.Event")
	proto.RegisterType((*SimpleResponse)(nil), "grpc.SimpleResponse")
	proto.RegisterType((*EventListResponse)(nil), "grpc.EventListResponse")
	proto.RegisterType((*CreateEventRequest)(nil), "grpc.CreateEventRequest")
	proto.RegisterType((*UpdateEventRequest)(nil), "grpc.UpdateEventRequest")
	proto.RegisterType((*DeleteEventRequest)(nil), "grpc.DeleteEventRequest")
	proto.RegisterType((*Nothing)(nil), "grpc.Nothing")
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c) }

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 388 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x52, 0xcd, 0x4e, 0xa3, 0x50,
	0x14, 0xee, 0x85, 0xfe, 0x84, 0x43, 0xa6, 0x33, 0xbd, 0x99, 0xcc, 0x10, 0x36, 0x43, 0x98, 0x59,
	0xb0, 0x98, 0xd0, 0x49, 0x67, 0x63, 0xa2, 0xc6, 0x85, 0x55, 0x37, 0xea, 0x82, 0x6a, 0x5c, 0xd3,
	0x72, 0xa4, 0x37, 0x02, 0x17, 0xb9, 0xb7, 0x4d, 0x7c, 0x01, 0x63, 0xe2, 0x0b, 0xf8, 0xb8, 0x06,
	0x28, 0x0d, 0x5a, 0xad, 0xed, 0xce, 0x1d, 0x9c, 0xfb, 0x9d, 0x2f, 0xdf, 0xcf, 0x01, 0xcd, 0x4f,
	0x99, 0x9b, 0x66, 0x5c, 0x72, 0xda, 0x0c, 0xb3, 0x74, 0x62, 0xfe, 0x0a, 0x39, 0x0f, 0x23, 0xec,
	0x17, 0xb3, 0xf1, 0xec, 0xba, 0x2f, 0x59, 0x8c, 0x42, 0xfa, 0x71, 0x5a, 0xc2, 0xec, 0x47, 0x02,
	0xad, 0xa3, 0x39, 0x26, 0x92, 0x76, 0x41, 0x61, 0x81, 0x41, 0x2c, 0xe2, 0xb4, 0x3c, 0x85, 0x05,
	0x94, 0x42, 0x33, 0xf1, 0x63, 0x34, 0x14, 0x8b, 0x38, 0x9a, 0x57, 0x7c, 0xd3, 0x7f, 0xd0, 0x12,
	0xd2, 0xcf, 0xa4, 0xa1, 0x5a, 0xc4, 0xd1, 0x07, 0xa6, 0x5b, 0xd2, 0xbb, 0x15, 0xbd, 0x7b, 0x51,
	0xd1, 0x7b, 0x25, 0x90, 0xfe, 0x05, 0x15, 0x93, 0xc0, 0x68, 0x7e, 0x88, 0xcf, 0x61, 0xb6, 0x03,
	0xdd, 0x11, 0x8b, 0xd3, 0x08, 0x3d, 0x14, 0x29, 0x4f, 0x04, 0xd2, 0x1f, 0xd0, 0xce, 0x50, 0xcc,
	0x22, 0x59, 0x28, 0xd3, 0xbc, 0xc5, 0x9f, 0xbd, 0x03, 0xbd, 0x42, 0xf6, 0x29, 0x13, 0x72, 0x09,
	0xfe, 0x0d, 0x6d, 0xcc, 0x87, 0xc2, 0x20, 0x96, 0xea, 0xe8, 0x03, 0xdd, 0xcd, 0x43, 0x70, 0x0b,
	0xa0, 0xb7, 0x78, 0xb2, 0x1f, 0x08, 0xd0, 0xc3, 0x0c, 0x7d, 0x89, 0xe5, 0x1c, 0x6f, 0x67, 0x28,
	0xe4, 0xd2, 0x2e, 0x79, 0xcb, 0xae, 0xb2, 0xa5, 0x5d, 0x75, 0x33, 0xbb, 0x4f, 0x04, 0xe8, 0x65,
	0x1a, 0xbc, 0x96, 0xf2, 0x19, 0x9a, 0xf8, 0x03, 0x74, 0x88, 0x11, 0xae, 0x57, 0x66, 0x6b, 0xd0,
	0x39, 0xe7, 0x72, 0xca, 0x92, 0x70, 0x70, 0xaf, 0x42, 0x67, 0x84, 0xd9, 0x9c, 0x4d, 0x90, 0x1e,
	0x80, 0x5e, 0x4b, 0x98, 0x1a, 0x65, 0x0d, 0xab, 0xa1, 0x9b, 0xdf, 0xcb, 0x97, 0x97, 0x9d, 0xdb,
	0x8d, 0x9c, 0xa0, 0x96, 0x4b, 0x45, 0xb0, 0x1a, 0xd5, 0x3a, 0x82, 0x9a, 0xfc, 0x8a, 0x60, 0xd5,
	0xd1, 0xbb, 0x04, 0xbb, 0xf0, 0xf5, 0x04, 0x65, 0x01, 0x15, 0xc7, 0x3c, 0x1b, 0xfa, 0x77, 0xf4,
	0x4b, 0x09, 0x5d, 0x18, 0x36, 0x7f, 0xd6, 0x8e, 0xab, 0x7e, 0x85, 0x76, 0x83, 0xee, 0xc1, 0xb7,
	0xfa, 0xf2, 0x15, 0xe2, 0xcd, 0x16, 0xdb, 0xfb, 0xd0, 0xab, 0x6f, 0x9f, 0xf1, 0x44, 0x4e, 0x37,
	0x5f, 0x1f, 0xb7, 0x8b, 0x4a, 0xff, 0x3f, 0x07, 0x00, 0x00, 0xff, 0xff, 0xc3, 0x91, 0x08, 0x2c,
	0x0c, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServiceClient interface {
	CreateEvent(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*SimpleResponse, error)
	UpdateEvent(ctx context.Context, in *UpdateEventRequest, opts ...grpc.CallOption) (*SimpleResponse, error)
	DeleteEvent(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*SimpleResponse, error)
	GetEventsForDay(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*EventListResponse, error)
	GetEventsForWeek(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*EventListResponse, error)
	GetEventsForMonth(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*EventListResponse, error)
}

type serviceClient struct {
	cc *grpc.ClientConn
}

func NewServiceClient(cc *grpc.ClientConn) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) CreateEvent(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*SimpleResponse, error) {
	out := new(SimpleResponse)
	err := c.cc.Invoke(ctx, "/grpc.Service/CreateEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) UpdateEvent(ctx context.Context, in *UpdateEventRequest, opts ...grpc.CallOption) (*SimpleResponse, error) {
	out := new(SimpleResponse)
	err := c.cc.Invoke(ctx, "/grpc.Service/UpdateEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DeleteEvent(ctx context.Context, in *DeleteEventRequest, opts ...grpc.CallOption) (*SimpleResponse, error) {
	out := new(SimpleResponse)
	err := c.cc.Invoke(ctx, "/grpc.Service/DeleteEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetEventsForDay(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*EventListResponse, error) {
	out := new(EventListResponse)
	err := c.cc.Invoke(ctx, "/grpc.Service/GetEventsForDay", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetEventsForWeek(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*EventListResponse, error) {
	out := new(EventListResponse)
	err := c.cc.Invoke(ctx, "/grpc.Service/GetEventsForWeek", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetEventsForMonth(ctx context.Context, in *Nothing, opts ...grpc.CallOption) (*EventListResponse, error) {
	out := new(EventListResponse)
	err := c.cc.Invoke(ctx, "/grpc.Service/GetEventsForMonth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
type ServiceServer interface {
	CreateEvent(context.Context, *CreateEventRequest) (*SimpleResponse, error)
	UpdateEvent(context.Context, *UpdateEventRequest) (*SimpleResponse, error)
	DeleteEvent(context.Context, *DeleteEventRequest) (*SimpleResponse, error)
	GetEventsForDay(context.Context, *Nothing) (*EventListResponse, error)
	GetEventsForWeek(context.Context, *Nothing) (*EventListResponse, error)
	GetEventsForMonth(context.Context, *Nothing) (*EventListResponse, error)
}

// UnimplementedServiceServer can be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (*UnimplementedServiceServer) CreateEvent(ctx context.Context, req *CreateEventRequest) (*SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEvent not implemented")
}
func (*UnimplementedServiceServer) UpdateEvent(ctx context.Context, req *UpdateEventRequest) (*SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEvent not implemented")
}
func (*UnimplementedServiceServer) DeleteEvent(ctx context.Context, req *DeleteEventRequest) (*SimpleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteEvent not implemented")
}
func (*UnimplementedServiceServer) GetEventsForDay(ctx context.Context, req *Nothing) (*EventListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventsForDay not implemented")
}
func (*UnimplementedServiceServer) GetEventsForWeek(ctx context.Context, req *Nothing) (*EventListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventsForWeek not implemented")
}
func (*UnimplementedServiceServer) GetEventsForMonth(ctx context.Context, req *Nothing) (*EventListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventsForMonth not implemented")
}

func RegisterServiceServer(s *grpc.Server, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_CreateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).CreateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Service/CreateEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).CreateEvent(ctx, req.(*CreateEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_UpdateEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).UpdateEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Service/UpdateEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).UpdateEvent(ctx, req.(*UpdateEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DeleteEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DeleteEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Service/DeleteEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DeleteEvent(ctx, req.(*DeleteEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetEventsForDay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nothing)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetEventsForDay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Service/GetEventsForDay",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetEventsForDay(ctx, req.(*Nothing))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetEventsForWeek_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nothing)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetEventsForWeek(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Service/GetEventsForWeek",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetEventsForWeek(ctx, req.(*Nothing))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetEventsForMonth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Nothing)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetEventsForMonth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.Service/GetEventsForMonth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetEventsForMonth(ctx, req.(*Nothing))
	}
	return interceptor(ctx, in, info, handler)
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateEvent",
			Handler:    _Service_CreateEvent_Handler,
		},
		{
			MethodName: "UpdateEvent",
			Handler:    _Service_UpdateEvent_Handler,
		},
		{
			MethodName: "DeleteEvent",
			Handler:    _Service_DeleteEvent_Handler,
		},
		{
			MethodName: "GetEventsForDay",
			Handler:    _Service_GetEventsForDay_Handler,
		},
		{
			MethodName: "GetEventsForWeek",
			Handler:    _Service_GetEventsForWeek_Handler,
		},
		{
			MethodName: "GetEventsForMonth",
			Handler:    _Service_GetEventsForMonth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}
