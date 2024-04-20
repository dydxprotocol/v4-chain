// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/feetiers/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// QueryPerpetualFeeParamsRequest is a request type for the PerpetualFeeParams
// RPC method.
type QueryPerpetualFeeParamsRequest struct {
}

func (m *QueryPerpetualFeeParamsRequest) Reset()         { *m = QueryPerpetualFeeParamsRequest{} }
func (m *QueryPerpetualFeeParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryPerpetualFeeParamsRequest) ProtoMessage()    {}
func (*QueryPerpetualFeeParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f31456045d64644f, []int{0}
}
func (m *QueryPerpetualFeeParamsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryPerpetualFeeParamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryPerpetualFeeParamsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryPerpetualFeeParamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryPerpetualFeeParamsRequest.Merge(m, src)
}
func (m *QueryPerpetualFeeParamsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryPerpetualFeeParamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryPerpetualFeeParamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryPerpetualFeeParamsRequest proto.InternalMessageInfo

// QueryPerpetualFeeParamsResponse is a response type for the PerpetualFeeParams
// RPC method.
type QueryPerpetualFeeParamsResponse struct {
	Params PerpetualFeeParams `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QueryPerpetualFeeParamsResponse) Reset()         { *m = QueryPerpetualFeeParamsResponse{} }
func (m *QueryPerpetualFeeParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryPerpetualFeeParamsResponse) ProtoMessage()    {}
func (*QueryPerpetualFeeParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f31456045d64644f, []int{1}
}
func (m *QueryPerpetualFeeParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryPerpetualFeeParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryPerpetualFeeParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryPerpetualFeeParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryPerpetualFeeParamsResponse.Merge(m, src)
}
func (m *QueryPerpetualFeeParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryPerpetualFeeParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryPerpetualFeeParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryPerpetualFeeParamsResponse proto.InternalMessageInfo

func (m *QueryPerpetualFeeParamsResponse) GetParams() PerpetualFeeParams {
	if m != nil {
		return m.Params
	}
	return PerpetualFeeParams{}
}

// QueryUserFeeTierRequest is a request type for the UserFeeTier RPC method.
type QueryUserFeeTierRequest struct {
	User string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (m *QueryUserFeeTierRequest) Reset()         { *m = QueryUserFeeTierRequest{} }
func (m *QueryUserFeeTierRequest) String() string { return proto.CompactTextString(m) }
func (*QueryUserFeeTierRequest) ProtoMessage()    {}
func (*QueryUserFeeTierRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f31456045d64644f, []int{2}
}
func (m *QueryUserFeeTierRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryUserFeeTierRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryUserFeeTierRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryUserFeeTierRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryUserFeeTierRequest.Merge(m, src)
}
func (m *QueryUserFeeTierRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryUserFeeTierRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryUserFeeTierRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryUserFeeTierRequest proto.InternalMessageInfo

func (m *QueryUserFeeTierRequest) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

// QueryUserFeeTierResponse is a request type for the UserFeeTier RPC method.
type QueryUserFeeTierResponse struct {
	// Index of the fee tier in the list queried from PerpetualFeeParams.
	Index uint32            `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	Tier  *PerpetualFeeTier `protobuf:"bytes,2,opt,name=tier,proto3" json:"tier,omitempty"`
}

func (m *QueryUserFeeTierResponse) Reset()         { *m = QueryUserFeeTierResponse{} }
func (m *QueryUserFeeTierResponse) String() string { return proto.CompactTextString(m) }
func (*QueryUserFeeTierResponse) ProtoMessage()    {}
func (*QueryUserFeeTierResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f31456045d64644f, []int{3}
}
func (m *QueryUserFeeTierResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryUserFeeTierResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryUserFeeTierResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryUserFeeTierResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryUserFeeTierResponse.Merge(m, src)
}
func (m *QueryUserFeeTierResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryUserFeeTierResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryUserFeeTierResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryUserFeeTierResponse proto.InternalMessageInfo

func (m *QueryUserFeeTierResponse) GetIndex() uint32 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *QueryUserFeeTierResponse) GetTier() *PerpetualFeeTier {
	if m != nil {
		return m.Tier
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryPerpetualFeeParamsRequest)(nil), "dydxprotocol.feetiers.QueryPerpetualFeeParamsRequest")
	proto.RegisterType((*QueryPerpetualFeeParamsResponse)(nil), "dydxprotocol.feetiers.QueryPerpetualFeeParamsResponse")
	proto.RegisterType((*QueryUserFeeTierRequest)(nil), "dydxprotocol.feetiers.QueryUserFeeTierRequest")
	proto.RegisterType((*QueryUserFeeTierResponse)(nil), "dydxprotocol.feetiers.QueryUserFeeTierResponse")
}

func init() { proto.RegisterFile("dydxprotocol/feetiers/query.proto", fileDescriptor_f31456045d64644f) }

var fileDescriptor_f31456045d64644f = []byte{
	// 453 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0xbf, 0x8e, 0xd3, 0x30,
	0x18, 0x4f, 0xaa, 0xde, 0x49, 0xf8, 0xc4, 0x62, 0x15, 0x11, 0x22, 0x94, 0x3b, 0xb2, 0x1c, 0x27,
	0xd1, 0x18, 0x1d, 0x70, 0x0b, 0x13, 0x1d, 0x7a, 0x8c, 0x25, 0x07, 0x0b, 0x4b, 0xe5, 0x26, 0xdf,
	0xe5, 0x8c, 0x1a, 0x3b, 0x67, 0x3b, 0xa8, 0x5d, 0x79, 0x02, 0x24, 0x1e, 0x80, 0x97, 0x60, 0xe3,
	0x05, 0x3a, 0x56, 0xb0, 0x30, 0x21, 0xd4, 0x22, 0xf1, 0x1a, 0x28, 0x4e, 0x02, 0x45, 0x69, 0xab,
	0xb2, 0xf9, 0xfb, 0xfc, 0xfb, 0xf7, 0x7d, 0x36, 0xba, 0x17, 0x4f, 0xe3, 0x49, 0x26, 0x85, 0x16,
	0x91, 0x18, 0x93, 0x4b, 0x00, 0xcd, 0x40, 0x2a, 0x72, 0x9d, 0x83, 0x9c, 0x06, 0xa6, 0x8f, 0x6f,
	0xad, 0x42, 0x82, 0x1a, 0xe2, 0xde, 0x89, 0x84, 0x4a, 0x85, 0x1a, 0x9a, 0x1b, 0x52, 0x16, 0x25,
	0xc3, 0xed, 0x24, 0x22, 0x11, 0x65, 0xbf, 0x38, 0x55, 0xdd, 0xbb, 0x89, 0x10, 0xc9, 0x18, 0x08,
	0xcd, 0x18, 0xa1, 0x9c, 0x0b, 0x4d, 0x35, 0x13, 0xbc, 0xe6, 0xf8, 0xeb, 0x83, 0x64, 0x54, 0xd2,
	0xb4, 0xc2, 0xf8, 0x47, 0xc8, 0x7b, 0x51, 0x04, 0x1b, 0x80, 0xcc, 0x40, 0xe7, 0x74, 0xdc, 0x07,
	0x18, 0x18, 0x40, 0x08, 0xd7, 0x39, 0x28, 0xed, 0xbf, 0x41, 0x87, 0x1b, 0x11, 0x2a, 0x13, 0x5c,
	0x01, 0x3e, 0x47, 0xfb, 0xa5, 0xa8, 0x63, 0x1f, 0xd9, 0xf7, 0x0f, 0x4e, 0x4f, 0x82, 0xb5, 0xf3,
	0x05, 0x4d, 0x89, 0x5e, 0x7b, 0xf6, 0xfd, 0xd0, 0x0a, 0x2b, 0xba, 0x7f, 0x8e, 0x6e, 0x1b, 0xaf,
	0x57, 0x0a, 0x64, 0x1f, 0xe0, 0x25, 0x03, 0x59, 0xc5, 0xc0, 0x0f, 0x50, 0x3b, 0x57, 0x20, 0x8d,
	0xc3, 0x8d, 0x9e, 0xf3, 0xe5, 0x53, 0xb7, 0x53, 0x2d, 0xe8, 0x59, 0x1c, 0x4b, 0x50, 0xea, 0x42,
	0x4b, 0xc6, 0x93, 0xd0, 0xa0, 0xfc, 0x14, 0x39, 0x4d, 0xa1, 0x2a, 0x6d, 0x07, 0xed, 0x31, 0x1e,
	0xc3, 0xc4, 0x48, 0xdd, 0x0c, 0xcb, 0x02, 0x3f, 0x45, 0xed, 0x22, 0xa4, 0xd3, 0x32, 0x13, 0x1c,
	0xef, 0x30, 0x81, 0x11, 0x35, 0xa4, 0xd3, 0x5f, 0x2d, 0xb4, 0x67, 0xfc, 0xf0, 0x67, 0x1b, 0xe1,
	0xe6, 0x98, 0xf8, 0xc9, 0x06, 0xbd, 0xed, 0xbb, 0x77, 0xcf, 0xfe, 0x97, 0x56, 0x8e, 0xe8, 0x9f,
	0xbd, 0xfb, 0xfa, 0xf3, 0x43, 0xeb, 0x21, 0x0e, 0xc8, 0x3f, 0x5f, 0xe0, 0xed, 0xe3, 0x95, 0x5f,
	0x50, 0xb3, 0x87, 0x97, 0x00, 0xc3, 0x72, 0xff, 0xf8, 0xa3, 0x8d, 0x0e, 0x56, 0x56, 0x86, 0x83,
	0x6d, 0xfe, 0xcd, 0x47, 0x72, 0xc9, 0xce, 0xf8, 0x2a, 0x28, 0x31, 0x41, 0x4f, 0xf0, 0xf1, 0xe6,
	0xa0, 0xc5, 0x7b, 0x9a, 0x8c, 0x45, 0xd9, 0x1b, 0xcd, 0x16, 0x9e, 0x3d, 0x5f, 0x78, 0xf6, 0x8f,
	0x85, 0x67, 0xbf, 0x5f, 0x7a, 0xd6, 0x7c, 0xe9, 0x59, 0xdf, 0x96, 0x9e, 0xf5, 0xfa, 0x79, 0xc2,
	0xf4, 0x55, 0x3e, 0x0a, 0x22, 0x91, 0x92, 0x0b, 0x2d, 0x81, 0xa6, 0x7d, 0xc6, 0x29, 0x8f, 0xa0,
	0x3b, 0xa8, 0x65, 0x95, 0x69, 0x77, 0xa3, 0x2b, 0xca, 0x38, 0xf9, 0x63, 0x36, 0xf9, 0xeb, 0xa5,
	0xa7, 0x19, 0xa8, 0xd1, 0xbe, 0xb9, 0x7a, 0xf4, 0x3b, 0x00, 0x00, 0xff, 0xff, 0x38, 0x34, 0xb7,
	0x9a, 0xc9, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Queries the PerpetualFeeParams.
	PerpetualFeeParams(ctx context.Context, in *QueryPerpetualFeeParamsRequest, opts ...grpc.CallOption) (*QueryPerpetualFeeParamsResponse, error)
	// Queries a user's fee tier
	UserFeeTier(ctx context.Context, in *QueryUserFeeTierRequest, opts ...grpc.CallOption) (*QueryUserFeeTierResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) PerpetualFeeParams(ctx context.Context, in *QueryPerpetualFeeParamsRequest, opts ...grpc.CallOption) (*QueryPerpetualFeeParamsResponse, error) {
	out := new(QueryPerpetualFeeParamsResponse)
	err := c.cc.Invoke(ctx, "/dydxprotocol.feetiers.Query/PerpetualFeeParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) UserFeeTier(ctx context.Context, in *QueryUserFeeTierRequest, opts ...grpc.CallOption) (*QueryUserFeeTierResponse, error) {
	out := new(QueryUserFeeTierResponse)
	err := c.cc.Invoke(ctx, "/dydxprotocol.feetiers.Query/UserFeeTier", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Queries the PerpetualFeeParams.
	PerpetualFeeParams(context.Context, *QueryPerpetualFeeParamsRequest) (*QueryPerpetualFeeParamsResponse, error)
	// Queries a user's fee tier
	UserFeeTier(context.Context, *QueryUserFeeTierRequest) (*QueryUserFeeTierResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) PerpetualFeeParams(ctx context.Context, req *QueryPerpetualFeeParamsRequest) (*QueryPerpetualFeeParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PerpetualFeeParams not implemented")
}
func (*UnimplementedQueryServer) UserFeeTier(ctx context.Context, req *QueryUserFeeTierRequest) (*QueryUserFeeTierResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserFeeTier not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_PerpetualFeeParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPerpetualFeeParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).PerpetualFeeParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dydxprotocol.feetiers.Query/PerpetualFeeParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).PerpetualFeeParams(ctx, req.(*QueryPerpetualFeeParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_UserFeeTier_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryUserFeeTierRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).UserFeeTier(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dydxprotocol.feetiers.Query/UserFeeTier",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).UserFeeTier(ctx, req.(*QueryUserFeeTierRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dydxprotocol.feetiers.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PerpetualFeeParams",
			Handler:    _Query_PerpetualFeeParams_Handler,
		},
		{
			MethodName: "UserFeeTier",
			Handler:    _Query_UserFeeTier_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dydxprotocol/feetiers/query.proto",
}

func (m *QueryPerpetualFeeParamsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryPerpetualFeeParamsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryPerpetualFeeParamsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryPerpetualFeeParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryPerpetualFeeParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryPerpetualFeeParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryUserFeeTierRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryUserFeeTierRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryUserFeeTierRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.User) > 0 {
		i -= len(m.User)
		copy(dAtA[i:], m.User)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.User)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryUserFeeTierResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryUserFeeTierResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryUserFeeTierResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Tier != nil {
		{
			size, err := m.Tier.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.Index != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryPerpetualFeeParamsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryPerpetualFeeParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryUserFeeTierRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.User)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryUserFeeTierResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Index != 0 {
		n += 1 + sovQuery(uint64(m.Index))
	}
	if m.Tier != nil {
		l = m.Tier.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryPerpetualFeeParamsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryPerpetualFeeParamsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryPerpetualFeeParamsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryPerpetualFeeParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryPerpetualFeeParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryPerpetualFeeParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryUserFeeTierRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryUserFeeTierRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryUserFeeTierRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field User", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.User = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryUserFeeTierResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryUserFeeTierResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryUserFeeTierResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tier", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Tier == nil {
				m.Tier = &PerpetualFeeTier{}
			}
			if err := m.Tier.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
