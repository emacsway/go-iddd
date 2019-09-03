// Code generated by protoc-gen-go. DO NOT EDIT.
// source: customer.proto

package customer

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type RegisterRequest struct {
	EmailAddress         string   `protobuf:"bytes,1,opt,name=emailAddress,proto3" json:"emailAddress,omitempty"`
	GivenName            string   `protobuf:"bytes,2,opt,name=givenName,proto3" json:"givenName,omitempty"`
	FamilyName           string   `protobuf:"bytes,3,opt,name=familyName,proto3" json:"familyName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterRequest) Reset()         { *m = RegisterRequest{} }
func (m *RegisterRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterRequest) ProtoMessage()    {}
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9efa92dae3d6ec46, []int{0}
}

func (m *RegisterRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterRequest.Unmarshal(m, b)
}
func (m *RegisterRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterRequest.Marshal(b, m, deterministic)
}
func (m *RegisterRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterRequest.Merge(m, src)
}
func (m *RegisterRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterRequest.Size(m)
}
func (m *RegisterRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterRequest proto.InternalMessageInfo

func (m *RegisterRequest) GetEmailAddress() string {
	if m != nil {
		return m.EmailAddress
	}
	return ""
}

func (m *RegisterRequest) GetGivenName() string {
	if m != nil {
		return m.GivenName
	}
	return ""
}

func (m *RegisterRequest) GetFamilyName() string {
	if m != nil {
		return m.FamilyName
	}
	return ""
}

type RegisterResponse struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterResponse) Reset()         { *m = RegisterResponse{} }
func (m *RegisterResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterResponse) ProtoMessage()    {}
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9efa92dae3d6ec46, []int{1}
}

func (m *RegisterResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterResponse.Unmarshal(m, b)
}
func (m *RegisterResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterResponse.Marshal(b, m, deterministic)
}
func (m *RegisterResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterResponse.Merge(m, src)
}
func (m *RegisterResponse) XXX_Size() int {
	return xxx_messageInfo_RegisterResponse.Size(m)
}
func (m *RegisterResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterResponse proto.InternalMessageInfo

func (m *RegisterResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func init() {
	proto.RegisterType((*RegisterRequest)(nil), "customer.RegisterRequest")
	proto.RegisterType((*RegisterResponse)(nil), "customer.RegisterResponse")
}

func init() { proto.RegisterFile("customer.proto", fileDescriptor_9efa92dae3d6ec46) }

var fileDescriptor_9efa92dae3d6ec46 = []byte{
	// 202 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x8f, 0xc1, 0x4a, 0xc4, 0x30,
	0x10, 0x86, 0xdd, 0x0a, 0xd2, 0x1d, 0x64, 0x95, 0x9c, 0xd6, 0x52, 0x44, 0x72, 0xf2, 0xd4, 0x82,
	0x3e, 0x81, 0xf4, 0xae, 0xd0, 0x37, 0x88, 0x66, 0x0c, 0x03, 0x4d, 0xa6, 0x66, 0x52, 0xc1, 0xb7,
	0x17, 0xd2, 0xd6, 0xaa, 0xec, 0x31, 0xff, 0xff, 0x91, 0x7f, 0x3e, 0x38, 0xbc, 0x4d, 0x92, 0xd8,
	0x63, 0x6c, 0xc6, 0xc8, 0x89, 0x55, 0xb9, 0xbe, 0xab, 0xda, 0x31, 0xbb, 0x01, 0x5b, 0x33, 0x52,
	0x6b, 0x42, 0xe0, 0x64, 0x12, 0x71, 0x90, 0x99, 0xd3, 0x02, 0x57, 0x3d, 0x3a, 0x92, 0x84, 0xb1,
	0xc7, 0x8f, 0x09, 0x25, 0x29, 0x0d, 0x97, 0xe8, 0x0d, 0x0d, 0x4f, 0xd6, 0x46, 0x14, 0x39, 0xee,
	0xee, 0x76, 0xf7, 0xfb, 0xfe, 0x4f, 0xa6, 0x6a, 0xd8, 0x3b, 0xfa, 0xc4, 0xf0, 0x6c, 0x3c, 0x1e,
	0x8b, 0x0c, 0x6c, 0x81, 0xba, 0x05, 0x78, 0x37, 0x9e, 0x86, 0xaf, 0x5c, 0x9f, 0xe7, 0xfa, 0x57,
	0xa2, 0x35, 0x5c, 0x6f, 0xa3, 0x32, 0x72, 0x10, 0x54, 0x07, 0x28, 0xc8, 0x2e, 0x5b, 0x05, 0xd9,
	0x87, 0x17, 0x28, 0xbb, 0x45, 0x41, 0x75, 0x50, 0xae, 0xbc, 0xba, 0x69, 0x7e, 0x4c, 0xff, 0x1d,
	0x5e, 0x55, 0xa7, 0xaa, 0xf9, 0x7b, 0x7d, 0xf6, 0x7a, 0x91, 0x85, 0x1f, 0xbf, 0x03, 0x00, 0x00,
	0xff, 0xff, 0x84, 0xea, 0xce, 0xd8, 0x2a, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CustomerClient is the client API for Customer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CustomerClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
}

type customerClient struct {
	cc *grpc.ClientConn
}

func NewCustomerClient(cc *grpc.ClientConn) CustomerClient {
	return &customerClient{cc}
}

func (c *customerClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/customer.Customer/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CustomerServer is the server API for Customer service.
type CustomerServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
}

// UnimplementedCustomerServer can be embedded to have forward compatible implementations.
type UnimplementedCustomerServer struct {
}

func (*UnimplementedCustomerServer) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}

func RegisterCustomerServer(s *grpc.Server, srv CustomerServer) {
	s.RegisterService(&_Customer_serviceDesc, srv)
}

func _Customer_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CustomerServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/customer.Customer/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CustomerServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Customer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "customer.Customer",
	HandlerType: (*CustomerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Customer_Register_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "customer.proto",
}
