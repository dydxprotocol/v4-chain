// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/assets/query.proto

package types

import (
	context "context"
	fmt "fmt"
	query "github.com/cosmos/cosmos-sdk/types/query"
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

// Queries an Asset by id.
type QueryAssetRequest struct {
	Id uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (m *QueryAssetRequest) Reset()         { *m = QueryAssetRequest{} }
func (m *QueryAssetRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAssetRequest) ProtoMessage()    {}
func (*QueryAssetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8e6c21d5bfb3fef3, []int{0}
}
func (m *QueryAssetRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAssetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAssetRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAssetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAssetRequest.Merge(m, src)
}
func (m *QueryAssetRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAssetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAssetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAssetRequest proto.InternalMessageInfo

func (m *QueryAssetRequest) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

// QueryAssetResponse is response type for the Asset RPC method.
type QueryAssetResponse struct {
	Asset Asset `protobuf:"bytes,1,opt,name=asset,proto3" json:"asset"`
}

func (m *QueryAssetResponse) Reset()         { *m = QueryAssetResponse{} }
func (m *QueryAssetResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAssetResponse) ProtoMessage()    {}
func (*QueryAssetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8e6c21d5bfb3fef3, []int{1}
}
func (m *QueryAssetResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAssetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAssetResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAssetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAssetResponse.Merge(m, src)
}
func (m *QueryAssetResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAssetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAssetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAssetResponse proto.InternalMessageInfo

func (m *QueryAssetResponse) GetAsset() Asset {
	if m != nil {
		return m.Asset
	}
	return Asset{}
}

// Queries a list of Asset items.
type QueryAllAssetsRequest struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllAssetsRequest) Reset()         { *m = QueryAllAssetsRequest{} }
func (m *QueryAllAssetsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAllAssetsRequest) ProtoMessage()    {}
func (*QueryAllAssetsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8e6c21d5bfb3fef3, []int{2}
}
func (m *QueryAllAssetsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllAssetsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllAssetsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllAssetsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllAssetsRequest.Merge(m, src)
}
func (m *QueryAllAssetsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllAssetsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllAssetsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllAssetsRequest proto.InternalMessageInfo

func (m *QueryAllAssetsRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryAllAssetsResponse is response type for the AllAssets RPC method.
type QueryAllAssetsResponse struct {
	Asset      []Asset             `protobuf:"bytes,1,rep,name=asset,proto3" json:"asset"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllAssetsResponse) Reset()         { *m = QueryAllAssetsResponse{} }
func (m *QueryAllAssetsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAllAssetsResponse) ProtoMessage()    {}
func (*QueryAllAssetsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8e6c21d5bfb3fef3, []int{3}
}
func (m *QueryAllAssetsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllAssetsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllAssetsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllAssetsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllAssetsResponse.Merge(m, src)
}
func (m *QueryAllAssetsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllAssetsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllAssetsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllAssetsResponse proto.InternalMessageInfo

func (m *QueryAllAssetsResponse) GetAsset() []Asset {
	if m != nil {
		return m.Asset
	}
	return nil
}

func (m *QueryAllAssetsResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryAssetRequest)(nil), "dydxprotocol.assets.QueryAssetRequest")
	proto.RegisterType((*QueryAssetResponse)(nil), "dydxprotocol.assets.QueryAssetResponse")
	proto.RegisterType((*QueryAllAssetsRequest)(nil), "dydxprotocol.assets.QueryAllAssetsRequest")
	proto.RegisterType((*QueryAllAssetsResponse)(nil), "dydxprotocol.assets.QueryAllAssetsResponse")
}

func init() { proto.RegisterFile("dydxprotocol/assets/query.proto", fileDescriptor_8e6c21d5bfb3fef3) }

var fileDescriptor_8e6c21d5bfb3fef3 = []byte{
	// 431 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xcf, 0xeb, 0xd3, 0x30,
	0x18, 0xc6, 0x9b, 0xea, 0x04, 0x23, 0x0a, 0xc6, 0x1f, 0x48, 0xf9, 0xd2, 0x69, 0x85, 0x4d, 0x26,
	0x4b, 0xd8, 0x04, 0xef, 0xee, 0x30, 0x2f, 0x1e, 0x66, 0xbd, 0x79, 0xd1, 0xb4, 0x0d, 0x5d, 0xa0,
	0x4b, 0xba, 0x25, 0x93, 0x0d, 0xf1, 0xa0, 0x27, 0x8f, 0x82, 0x20, 0xf8, 0x1f, 0xed, 0x38, 0xf0,
	0xe2, 0x49, 0x64, 0xf3, 0x0f, 0x91, 0x26, 0xdd, 0xac, 0xfb, 0xe1, 0xfc, 0x9e, 0x5a, 0xde, 0xf7,
	0x79, 0x9f, 0xe7, 0xf3, 0xe6, 0x85, 0xf5, 0x64, 0x9e, 0xcc, 0xf2, 0x89, 0xd4, 0x32, 0x96, 0x19,
	0xa1, 0x4a, 0x31, 0xad, 0xc8, 0x78, 0xca, 0x26, 0x73, 0x6c, 0xaa, 0xe8, 0x46, 0x55, 0x80, 0xad,
	0xc0, 0xbb, 0x99, 0xca, 0x54, 0x9a, 0x22, 0x29, 0xfe, 0xac, 0xd4, 0x3b, 0x4b, 0xa5, 0x4c, 0x33,
	0x46, 0x68, 0xce, 0x09, 0x15, 0x42, 0x6a, 0xaa, 0xb9, 0x14, 0xaa, 0xec, 0xb6, 0x62, 0xa9, 0x46,
	0x52, 0x91, 0x88, 0x2a, 0x66, 0x13, 0xc8, 0x9b, 0x4e, 0xc4, 0x34, 0xed, 0x90, 0x9c, 0xa6, 0x5c,
	0x18, 0x71, 0xa9, 0x3d, 0x48, 0x65, 0x3e, 0x56, 0x10, 0xdc, 0x87, 0xd7, 0x9f, 0x17, 0x16, 0x4f,
	0x8a, 0x5a, 0xc8, 0xc6, 0x53, 0xa6, 0x34, 0xba, 0x06, 0x5d, 0x9e, 0xdc, 0x01, 0x77, 0xc1, 0x83,
	0xab, 0xa1, 0xcb, 0x93, 0xe0, 0x19, 0x44, 0x55, 0x91, 0xca, 0xa5, 0x50, 0x0c, 0x3d, 0x86, 0x35,
	0xe3, 0x64, 0x84, 0x57, 0xba, 0x1e, 0x3e, 0xb0, 0x20, 0x36, 0x23, 0xbd, 0x8b, 0x8b, 0x1f, 0x75,
	0x27, 0xb4, 0xf2, 0xe0, 0x15, 0xbc, 0x65, 0xdd, 0xb2, 0xcc, 0x74, 0xd5, 0x26, 0xb6, 0x0f, 0xe1,
	0x9f, 0x05, 0x4a, 0xd7, 0x06, 0xb6, 0xdb, 0xe2, 0x62, 0x5b, 0x6c, 0xdf, 0xb3, 0xdc, 0x16, 0x0f,
	0x68, 0xca, 0xca, 0xd9, 0xb0, 0x32, 0x19, 0x7c, 0x05, 0xf0, 0xf6, 0x6e, 0xc2, 0x3e, 0xf3, 0x85,
	0x73, 0x30, 0xa3, 0xa7, 0x7f, 0xa1, 0xb9, 0x06, 0xad, 0x79, 0x12, 0xcd, 0x86, 0x56, 0xd9, 0xba,
	0x5f, 0x5c, 0x58, 0x33, 0x6c, 0xe8, 0x3d, 0x80, 0x35, 0x93, 0x84, 0x1a, 0x07, 0x29, 0xf6, 0xce,
	0xe2, 0x35, 0x4f, 0xea, 0x6c, 0x60, 0xd0, 0xfc, 0xf0, 0xed, 0xd7, 0x67, 0xf7, 0x1e, 0xaa, 0x93,
	0xa3, 0xe7, 0x27, 0x6f, 0x79, 0xf2, 0x0e, 0x7d, 0x04, 0xf0, 0xf2, 0xf6, 0x91, 0x50, 0xeb, 0x1f,
	0xfe, 0x3b, 0xb7, 0xf2, 0x1e, 0xfe, 0x97, 0xb6, 0xe4, 0x09, 0x0c, 0xcf, 0x19, 0xf2, 0x8e, 0xf3,
	0xf4, 0x5e, 0x2f, 0x56, 0x3e, 0x58, 0xae, 0x7c, 0xf0, 0x73, 0xe5, 0x83, 0x4f, 0x6b, 0xdf, 0x59,
	0xae, 0x7d, 0xe7, 0xfb, 0xda, 0x77, 0x5e, 0xf6, 0x53, 0xae, 0x87, 0xd3, 0x08, 0xc7, 0x72, 0x44,
	0x5e, 0xe8, 0x09, 0xa3, 0xa3, 0x3e, 0x17, 0x54, 0xc4, 0xac, 0x3d, 0xd8, 0x38, 0x29, 0x53, 0x6e,
	0xc7, 0x43, 0xca, 0x05, 0xd9, 0xfa, 0xcf, 0x36, 0x09, 0x7a, 0x9e, 0x33, 0x15, 0x5d, 0x32, 0x8d,
	0x47, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0x8e, 0x29, 0xea, 0xa1, 0xaa, 0x03, 0x00, 0x00,
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
	// Queries a Asset by id.
	Asset(ctx context.Context, in *QueryAssetRequest, opts ...grpc.CallOption) (*QueryAssetResponse, error)
	// Queries a list of Asset items.
	AllAssets(ctx context.Context, in *QueryAllAssetsRequest, opts ...grpc.CallOption) (*QueryAllAssetsResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Asset(ctx context.Context, in *QueryAssetRequest, opts ...grpc.CallOption) (*QueryAssetResponse, error) {
	out := new(QueryAssetResponse)
	err := c.cc.Invoke(ctx, "/dydxprotocol.assets.Query/Asset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) AllAssets(ctx context.Context, in *QueryAllAssetsRequest, opts ...grpc.CallOption) (*QueryAllAssetsResponse, error) {
	out := new(QueryAllAssetsResponse)
	err := c.cc.Invoke(ctx, "/dydxprotocol.assets.Query/AllAssets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Queries a Asset by id.
	Asset(context.Context, *QueryAssetRequest) (*QueryAssetResponse, error)
	// Queries a list of Asset items.
	AllAssets(context.Context, *QueryAllAssetsRequest) (*QueryAllAssetsResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Asset(ctx context.Context, req *QueryAssetRequest) (*QueryAssetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Asset not implemented")
}
func (*UnimplementedQueryServer) AllAssets(ctx context.Context, req *QueryAllAssetsRequest) (*QueryAllAssetsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllAssets not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Asset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAssetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Asset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dydxprotocol.assets.Query/Asset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Asset(ctx, req.(*QueryAssetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_AllAssets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllAssetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).AllAssets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dydxprotocol.assets.Query/AllAssets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).AllAssets(ctx, req.(*QueryAllAssetsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dydxprotocol.assets.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Asset",
			Handler:    _Query_Asset_Handler,
		},
		{
			MethodName: "AllAssets",
			Handler:    _Query_AllAssets_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dydxprotocol/assets/query.proto",
}

func (m *QueryAssetRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAssetRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAssetRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryAssetResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAssetResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAssetResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Asset.MarshalToSizedBuffer(dAtA[:i])
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

func (m *QueryAllAssetsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllAssetsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllAssetsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryAllAssetsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllAssetsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllAssetsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Asset) > 0 {
		for iNdEx := len(m.Asset) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Asset[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
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
func (m *QueryAssetRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovQuery(uint64(m.Id))
	}
	return n
}

func (m *QueryAssetResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Asset.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryAllAssetsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryAllAssetsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Asset) > 0 {
		for _, e := range m.Asset {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
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
func (m *QueryAssetRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryAssetRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAssetRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
func (m *QueryAssetResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryAssetResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAssetResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asset", wireType)
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
			if err := m.Asset.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *QueryAllAssetsRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryAllAssetsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllAssetsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
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
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *QueryAllAssetsResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryAllAssetsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllAssetsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Asset", wireType)
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
			m.Asset = append(m.Asset, Asset{})
			if err := m.Asset[len(m.Asset)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
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
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
