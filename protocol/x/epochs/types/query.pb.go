// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: klyraprotocol/epochs/query.proto

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

// QueryGetEpochInfoRequest is request type for the GetEpochInfo RPC method.
type QueryGetEpochInfoRequest struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (m *QueryGetEpochInfoRequest) Reset()         { *m = QueryGetEpochInfoRequest{} }
func (m *QueryGetEpochInfoRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetEpochInfoRequest) ProtoMessage()    {}
func (*QueryGetEpochInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dc7b48a99645dfa, []int{0}
}
func (m *QueryGetEpochInfoRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetEpochInfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetEpochInfoRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetEpochInfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetEpochInfoRequest.Merge(m, src)
}
func (m *QueryGetEpochInfoRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetEpochInfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetEpochInfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetEpochInfoRequest proto.InternalMessageInfo

func (m *QueryGetEpochInfoRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// QueryEpochInfoResponse is response type for the GetEpochInfo RPC method.
type QueryEpochInfoResponse struct {
	EpochInfo EpochInfo `protobuf:"bytes,1,opt,name=epoch_info,json=epochInfo,proto3" json:"epoch_info"`
}

func (m *QueryEpochInfoResponse) Reset()         { *m = QueryEpochInfoResponse{} }
func (m *QueryEpochInfoResponse) String() string { return proto.CompactTextString(m) }
func (*QueryEpochInfoResponse) ProtoMessage()    {}
func (*QueryEpochInfoResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dc7b48a99645dfa, []int{1}
}
func (m *QueryEpochInfoResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryEpochInfoResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryEpochInfoResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryEpochInfoResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryEpochInfoResponse.Merge(m, src)
}
func (m *QueryEpochInfoResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryEpochInfoResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryEpochInfoResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryEpochInfoResponse proto.InternalMessageInfo

func (m *QueryEpochInfoResponse) GetEpochInfo() EpochInfo {
	if m != nil {
		return m.EpochInfo
	}
	return EpochInfo{}
}

// QueryAllEpochInfoRequest is request type for the AllEpochInfo RPC method.
type QueryAllEpochInfoRequest struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllEpochInfoRequest) Reset()         { *m = QueryAllEpochInfoRequest{} }
func (m *QueryAllEpochInfoRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAllEpochInfoRequest) ProtoMessage()    {}
func (*QueryAllEpochInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dc7b48a99645dfa, []int{2}
}
func (m *QueryAllEpochInfoRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllEpochInfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllEpochInfoRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllEpochInfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllEpochInfoRequest.Merge(m, src)
}
func (m *QueryAllEpochInfoRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllEpochInfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllEpochInfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllEpochInfoRequest proto.InternalMessageInfo

func (m *QueryAllEpochInfoRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryEpochInfoAllResponse is response type for the AllEpochInfo RPC method.
type QueryEpochInfoAllResponse struct {
	EpochInfo  []EpochInfo         `protobuf:"bytes,1,rep,name=epoch_info,json=epochInfo,proto3" json:"epoch_info"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryEpochInfoAllResponse) Reset()         { *m = QueryEpochInfoAllResponse{} }
func (m *QueryEpochInfoAllResponse) String() string { return proto.CompactTextString(m) }
func (*QueryEpochInfoAllResponse) ProtoMessage()    {}
func (*QueryEpochInfoAllResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2dc7b48a99645dfa, []int{3}
}
func (m *QueryEpochInfoAllResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryEpochInfoAllResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryEpochInfoAllResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryEpochInfoAllResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryEpochInfoAllResponse.Merge(m, src)
}
func (m *QueryEpochInfoAllResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryEpochInfoAllResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryEpochInfoAllResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryEpochInfoAllResponse proto.InternalMessageInfo

func (m *QueryEpochInfoAllResponse) GetEpochInfo() []EpochInfo {
	if m != nil {
		return m.EpochInfo
	}
	return nil
}

func (m *QueryEpochInfoAllResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryGetEpochInfoRequest)(nil), "klyraprotocol.epochs.QueryGetEpochInfoRequest")
	proto.RegisterType((*QueryEpochInfoResponse)(nil), "klyraprotocol.epochs.QueryEpochInfoResponse")
	proto.RegisterType((*QueryAllEpochInfoRequest)(nil), "klyraprotocol.epochs.QueryAllEpochInfoRequest")
	proto.RegisterType((*QueryEpochInfoAllResponse)(nil), "klyraprotocol.epochs.QueryEpochInfoAllResponse")
}

func init() { proto.RegisterFile("klyraprotocol/epochs/query.proto", fileDescriptor_2dc7b48a99645dfa) }

var fileDescriptor_2dc7b48a99645dfa = []byte{
	// 457 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0x4f, 0x8b, 0xd3, 0x40,
	0x18, 0xc6, 0x33, 0x75, 0x15, 0x3a, 0x7a, 0x1a, 0x16, 0x59, 0x83, 0x64, 0x97, 0xc8, 0xaa, 0xac,
	0xdb, 0x19, 0x5a, 0xfd, 0x02, 0x2d, 0xda, 0xe2, 0xad, 0xc6, 0x9b, 0x07, 0x75, 0x12, 0xa6, 0x69,
	0x30, 0x99, 0x49, 0x33, 0xd3, 0x62, 0x11, 0x2f, 0x7e, 0x02, 0xc1, 0xab, 0x7e, 0x01, 0x3f, 0x88,
	0xf4, 0x58, 0xf0, 0xe2, 0x49, 0xa4, 0xf5, 0x83, 0x48, 0x26, 0x69, 0x9b, 0xd0, 0xb4, 0x96, 0x3d,
	0x25, 0xbc, 0xf3, 0xbc, 0xef, 0xf3, 0x7b, 0xff, 0xc0, 0xb3, 0x77, 0xe1, 0x34, 0xa1, 0x71, 0x22,
	0x94, 0xf0, 0x44, 0x48, 0x58, 0x2c, 0xbc, 0xa1, 0x24, 0xa3, 0x31, 0x4b, 0xa6, 0x58, 0x47, 0xd1,
	0x71, 0x49, 0x81, 0x33, 0x85, 0x79, 0xec, 0x0b, 0x5f, 0xe8, 0x20, 0x49, 0xff, 0x32, 0xad, 0x79,
	0xd7, 0x17, 0xc2, 0x0f, 0x19, 0xa1, 0x71, 0x40, 0x28, 0xe7, 0x42, 0x51, 0x15, 0x08, 0x2e, 0xf3,
	0xd7, 0x0b, 0x4f, 0xc8, 0x48, 0x48, 0xe2, 0x52, 0xc9, 0x32, 0x0b, 0x32, 0x69, 0xba, 0x4c, 0xd1,
	0x26, 0x89, 0xa9, 0x1f, 0x70, 0x2d, 0xce, 0xb5, 0xe7, 0x95, 0x5c, 0xfa, 0xf3, 0x26, 0xe0, 0x83,
	0xdc, 0xd0, 0xc6, 0xf0, 0xe4, 0x45, 0x5a, 0xa8, 0xc7, 0xd4, 0xb3, 0xf4, 0xed, 0x39, 0x1f, 0x08,
	0x87, 0x8d, 0xc6, 0x4c, 0x2a, 0x84, 0xe0, 0x11, 0xa7, 0x11, 0x3b, 0x01, 0x67, 0xe0, 0x61, 0xdd,
	0xd1, 0xff, 0xf6, 0x6b, 0x78, 0x5b, 0xeb, 0x0b, 0x62, 0x19, 0x0b, 0x2e, 0x19, 0x7a, 0x0a, 0xe1,
	0xa6, 0xba, 0xce, 0xb9, 0xd9, 0x3a, 0xc5, 0x55, 0xbd, 0xe3, 0x75, 0x72, 0xe7, 0x68, 0xf6, 0xfb,
	0xd4, 0x70, 0xea, 0x6c, 0x15, 0xb0, 0xdd, 0x9c, 0xa7, 0x1d, 0x86, 0x5b, 0x3c, 0x5d, 0x08, 0x37,
	0x6d, 0xe6, 0x0e, 0xf7, 0x71, 0x36, 0x13, 0x9c, 0xce, 0x04, 0x67, 0x63, 0xcf, 0x67, 0x82, 0xfb,
	0xd4, 0x67, 0x79, 0xae, 0x53, 0xc8, 0xb4, 0xbf, 0x03, 0x78, 0xa7, 0xdc, 0x44, 0x3b, 0x0c, 0x77,
	0xf6, 0x71, 0xed, 0x2a, 0x7d, 0xa0, 0x5e, 0x89, 0xb5, 0xa6, 0x59, 0x1f, 0xfc, 0x97, 0x35, 0x43,
	0x28, 0xc2, 0xb6, 0x7e, 0xd4, 0xe0, 0x75, 0x0d, 0x8b, 0xbe, 0x01, 0x58, 0x5f, 0x3b, 0x22, 0x5c,
	0x8d, 0xb4, 0x6b, 0x99, 0xe6, 0xe5, 0x1e, 0xfd, 0xd6, 0x32, 0xed, 0xd6, 0xa7, 0x9f, 0x7f, 0xbf,
	0xd4, 0x2e, 0xd1, 0x05, 0x29, 0x9f, 0xd1, 0xe4, 0xc9, 0xf6, 0x25, 0x91, 0x0f, 0xe9, 0x65, 0x7c,
	0x44, 0x5f, 0x01, 0xbc, 0x55, 0x9c, 0xe8, 0x5e, 0xc4, 0x8a, 0xfd, 0x9a, 0xe4, 0x10, 0xc4, 0xc2,
	0xaa, 0xec, 0x47, 0x9a, 0xf2, 0x1c, 0xdd, 0x3b, 0x80, 0xb2, 0xf3, 0x76, 0xb6, 0xb0, 0xc0, 0x7c,
	0x61, 0x81, 0x3f, 0x0b, 0x0b, 0x7c, 0x5e, 0x5a, 0xc6, 0x7c, 0x69, 0x19, 0xbf, 0x96, 0x96, 0xf1,
	0xaa, 0xeb, 0x07, 0x6a, 0x38, 0x76, 0xb1, 0x27, 0x22, 0xf2, 0x52, 0x25, 0x8c, 0x46, 0xdd, 0x80,
	0x53, 0xee, 0xb1, 0x46, 0x7f, 0x55, 0x51, 0xea, 0x70, 0xc3, 0x1b, 0xd2, 0x80, 0x93, 0xb5, 0xcf,
	0xfb, 0x95, 0x8d, 0x9a, 0xc6, 0x4c, 0xba, 0x37, 0xf4, 0xc3, 0xe3, 0x7f, 0x01, 0x00, 0x00, 0xff,
	0xff, 0x47, 0x5f, 0xfb, 0x20, 0x13, 0x04, 0x00, 0x00,
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
	// Queries a EpochInfo by name.
	EpochInfo(ctx context.Context, in *QueryGetEpochInfoRequest, opts ...grpc.CallOption) (*QueryEpochInfoResponse, error)
	// Queries a list of EpochInfo items.
	EpochInfoAll(ctx context.Context, in *QueryAllEpochInfoRequest, opts ...grpc.CallOption) (*QueryEpochInfoAllResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) EpochInfo(ctx context.Context, in *QueryGetEpochInfoRequest, opts ...grpc.CallOption) (*QueryEpochInfoResponse, error) {
	out := new(QueryEpochInfoResponse)
	err := c.cc.Invoke(ctx, "/klyraprotocol.epochs.Query/EpochInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) EpochInfoAll(ctx context.Context, in *QueryAllEpochInfoRequest, opts ...grpc.CallOption) (*QueryEpochInfoAllResponse, error) {
	out := new(QueryEpochInfoAllResponse)
	err := c.cc.Invoke(ctx, "/klyraprotocol.epochs.Query/EpochInfoAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Queries a EpochInfo by name.
	EpochInfo(context.Context, *QueryGetEpochInfoRequest) (*QueryEpochInfoResponse, error)
	// Queries a list of EpochInfo items.
	EpochInfoAll(context.Context, *QueryAllEpochInfoRequest) (*QueryEpochInfoAllResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) EpochInfo(ctx context.Context, req *QueryGetEpochInfoRequest) (*QueryEpochInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EpochInfo not implemented")
}
func (*UnimplementedQueryServer) EpochInfoAll(ctx context.Context, req *QueryAllEpochInfoRequest) (*QueryEpochInfoAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EpochInfoAll not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_EpochInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetEpochInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).EpochInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/klyraprotocol.epochs.Query/EpochInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).EpochInfo(ctx, req.(*QueryGetEpochInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_EpochInfoAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllEpochInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).EpochInfoAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/klyraprotocol.epochs.Query/EpochInfoAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).EpochInfoAll(ctx, req.(*QueryAllEpochInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "klyraprotocol.epochs.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EpochInfo",
			Handler:    _Query_EpochInfo_Handler,
		},
		{
			MethodName: "EpochInfoAll",
			Handler:    _Query_EpochInfoAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "klyraprotocol/epochs/query.proto",
}

func (m *QueryGetEpochInfoRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetEpochInfoRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetEpochInfoRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryEpochInfoResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryEpochInfoResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryEpochInfoResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.EpochInfo.MarshalToSizedBuffer(dAtA[:i])
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

func (m *QueryAllEpochInfoRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllEpochInfoRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllEpochInfoRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *QueryEpochInfoAllResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryEpochInfoAllResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryEpochInfoAllResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if len(m.EpochInfo) > 0 {
		for iNdEx := len(m.EpochInfo) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.EpochInfo[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
func (m *QueryGetEpochInfoRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryEpochInfoResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.EpochInfo.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryAllEpochInfoRequest) Size() (n int) {
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

func (m *QueryEpochInfoAllResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.EpochInfo) > 0 {
		for _, e := range m.EpochInfo {
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
func (m *QueryGetEpochInfoRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryGetEpochInfoRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetEpochInfoRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
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
			m.Name = string(dAtA[iNdEx:postIndex])
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
func (m *QueryEpochInfoResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryEpochInfoResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryEpochInfoResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EpochInfo", wireType)
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
			if err := m.EpochInfo.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *QueryAllEpochInfoRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryAllEpochInfoRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllEpochInfoRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *QueryEpochInfoAllResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryEpochInfoAllResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryEpochInfoAllResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EpochInfo", wireType)
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
			m.EpochInfo = append(m.EpochInfo, EpochInfo{})
			if err := m.EpochInfo[len(m.EpochInfo)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
