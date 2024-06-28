// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/revshare/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
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

// Message to set the market mapper revenue share
type MsgSetMarketMapperRevenueShare struct {
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	// Parameters for the revenue share
	Params MarketMapperRevenueShareParams `protobuf:"bytes,2,opt,name=params,proto3" json:"params"`
}

func (m *MsgSetMarketMapperRevenueShare) Reset()         { *m = MsgSetMarketMapperRevenueShare{} }
func (m *MsgSetMarketMapperRevenueShare) String() string { return proto.CompactTextString(m) }
func (*MsgSetMarketMapperRevenueShare) ProtoMessage()    {}
func (*MsgSetMarketMapperRevenueShare) Descriptor() ([]byte, []int) {
	return fileDescriptor_460d8062a262197e, []int{0}
}
func (m *MsgSetMarketMapperRevenueShare) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetMarketMapperRevenueShare) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetMarketMapperRevenueShare.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetMarketMapperRevenueShare) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetMarketMapperRevenueShare.Merge(m, src)
}
func (m *MsgSetMarketMapperRevenueShare) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetMarketMapperRevenueShare) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetMarketMapperRevenueShare.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetMarketMapperRevenueShare proto.InternalMessageInfo

func (m *MsgSetMarketMapperRevenueShare) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *MsgSetMarketMapperRevenueShare) GetParams() MarketMapperRevenueShareParams {
	if m != nil {
		return m.Params
	}
	return MarketMapperRevenueShareParams{}
}

// Response to a MsgSetMarketMapperRevenueShare
type MsgSetMarketMapperRevenueShareResponse struct {
}

func (m *MsgSetMarketMapperRevenueShareResponse) Reset() {
	*m = MsgSetMarketMapperRevenueShareResponse{}
}
func (m *MsgSetMarketMapperRevenueShareResponse) String() string { return proto.CompactTextString(m) }
func (*MsgSetMarketMapperRevenueShareResponse) ProtoMessage()    {}
func (*MsgSetMarketMapperRevenueShareResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_460d8062a262197e, []int{1}
}
func (m *MsgSetMarketMapperRevenueShareResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetMarketMapperRevenueShareResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetMarketMapperRevenueShareResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetMarketMapperRevenueShareResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetMarketMapperRevenueShareResponse.Merge(m, src)
}
func (m *MsgSetMarketMapperRevenueShareResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetMarketMapperRevenueShareResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetMarketMapperRevenueShareResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetMarketMapperRevenueShareResponse proto.InternalMessageInfo

// Msg to set market mapper revenue share details (e.g. expiration timestamp) for a
// specific market. To be used as an override for existing revenue share
// settings set by the MsgSetMarketMapperRevenueShare msg
type MsgSetMarketMapperRevShareDetailsForMarket struct {
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	// Parameters for the revenue share details
	Params MarketRevShareDetailsParams `protobuf:"bytes,2,opt,name=params,proto3" json:"params"`
}

func (m *MsgSetMarketMapperRevShareDetailsForMarket) Reset() {
	*m = MsgSetMarketMapperRevShareDetailsForMarket{}
}
func (m *MsgSetMarketMapperRevShareDetailsForMarket) String() string {
	return proto.CompactTextString(m)
}
func (*MsgSetMarketMapperRevShareDetailsForMarket) ProtoMessage() {}
func (*MsgSetMarketMapperRevShareDetailsForMarket) Descriptor() ([]byte, []int) {
	return fileDescriptor_460d8062a262197e, []int{2}
}
func (m *MsgSetMarketMapperRevShareDetailsForMarket) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetMarketMapperRevShareDetailsForMarket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetMarketMapperRevShareDetailsForMarket.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetMarketMapperRevShareDetailsForMarket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetMarketMapperRevShareDetailsForMarket.Merge(m, src)
}
func (m *MsgSetMarketMapperRevShareDetailsForMarket) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetMarketMapperRevShareDetailsForMarket) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetMarketMapperRevShareDetailsForMarket.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetMarketMapperRevShareDetailsForMarket proto.InternalMessageInfo

func (m *MsgSetMarketMapperRevShareDetailsForMarket) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *MsgSetMarketMapperRevShareDetailsForMarket) GetParams() MarketRevShareDetailsParams {
	if m != nil {
		return m.Params
	}
	return MarketRevShareDetailsParams{}
}

// Response to a MsgSetMarketMapperRevShareDetailsForMarket
type MsgSetMarketMapperRevShareDetailsForMarketResponse struct {
}

func (m *MsgSetMarketMapperRevShareDetailsForMarketResponse) Reset() {
	*m = MsgSetMarketMapperRevShareDetailsForMarketResponse{}
}
func (m *MsgSetMarketMapperRevShareDetailsForMarketResponse) String() string {
	return proto.CompactTextString(m)
}
func (*MsgSetMarketMapperRevShareDetailsForMarketResponse) ProtoMessage() {}
func (*MsgSetMarketMapperRevShareDetailsForMarketResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_460d8062a262197e, []int{3}
}
func (m *MsgSetMarketMapperRevShareDetailsForMarketResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetMarketMapperRevShareDetailsForMarketResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetMarketMapperRevShareDetailsForMarketResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetMarketMapperRevShareDetailsForMarketResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetMarketMapperRevShareDetailsForMarketResponse.Merge(m, src)
}
func (m *MsgSetMarketMapperRevShareDetailsForMarketResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetMarketMapperRevShareDetailsForMarketResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetMarketMapperRevShareDetailsForMarketResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetMarketMapperRevShareDetailsForMarketResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgSetMarketMapperRevenueShare)(nil), "dydxprotocol.revshare.MsgSetMarketMapperRevenueShare")
	proto.RegisterType((*MsgSetMarketMapperRevenueShareResponse)(nil), "dydxprotocol.revshare.MsgSetMarketMapperRevenueShareResponse")
	proto.RegisterType((*MsgSetMarketMapperRevShareDetailsForMarket)(nil), "dydxprotocol.revshare.MsgSetMarketMapperRevShareDetailsForMarket")
	proto.RegisterType((*MsgSetMarketMapperRevShareDetailsForMarketResponse)(nil), "dydxprotocol.revshare.MsgSetMarketMapperRevShareDetailsForMarketResponse")
}

func init() { proto.RegisterFile("dydxprotocol/revshare/tx.proto", fileDescriptor_460d8062a262197e) }

var fileDescriptor_460d8062a262197e = []byte{
	// 400 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4b, 0xa9, 0x4c, 0xa9,
	0x28, 0x28, 0xca, 0x2f, 0xc9, 0x4f, 0xce, 0xcf, 0xd1, 0x2f, 0x4a, 0x2d, 0x2b, 0xce, 0x48, 0x2c,
	0x4a, 0xd5, 0x2f, 0xa9, 0xd0, 0x03, 0x0b, 0x0a, 0x89, 0x22, 0xcb, 0xeb, 0xc1, 0xe4, 0xa5, 0x24,
	0x93, 0xf3, 0x8b, 0x73, 0xf3, 0x8b, 0xe3, 0xc1, 0x32, 0xfa, 0x10, 0x0e, 0x44, 0x87, 0x94, 0x38,
	0x84, 0xa7, 0x9f, 0x5b, 0x9c, 0xae, 0x5f, 0x66, 0x08, 0xa2, 0xa0, 0x12, 0x22, 0xe9, 0xf9, 0xe9,
	0xf9, 0x10, 0x0d, 0x20, 0x16, 0x54, 0x54, 0x09, 0xbb, 0x03, 0x0a, 0x12, 0x8b, 0x12, 0x73, 0xa1,
	0x46, 0x2a, 0xed, 0x65, 0xe4, 0x92, 0xf3, 0x2d, 0x4e, 0x0f, 0x4e, 0x2d, 0xf1, 0x4d, 0x2c, 0xca,
	0x06, 0x91, 0x05, 0x05, 0xa9, 0x45, 0x41, 0xa9, 0x65, 0xa9, 0x79, 0xa5, 0xa9, 0xc1, 0x20, 0xf5,
	0x42, 0x66, 0x5c, 0x9c, 0x89, 0xa5, 0x25, 0x19, 0xf9, 0x45, 0x99, 0x25, 0x95, 0x12, 0x8c, 0x0a,
	0x8c, 0x1a, 0x9c, 0x4e, 0x12, 0x97, 0xb6, 0xe8, 0x8a, 0x40, 0x9d, 0xe6, 0x98, 0x92, 0x52, 0x94,
	0x5a, 0x5c, 0x1c, 0x5c, 0x52, 0x94, 0x99, 0x97, 0x1e, 0x84, 0x50, 0x2a, 0x14, 0xcc, 0xc5, 0x06,
	0xb1, 0x4a, 0x82, 0x49, 0x81, 0x51, 0x83, 0xdb, 0xc8, 0x54, 0x0f, 0xab, 0x87, 0xf5, 0x70, 0x59,
	0x1c, 0x00, 0xd6, 0xec, 0xc4, 0x72, 0xe2, 0x9e, 0x3c, 0x43, 0x10, 0xd4, 0x28, 0x2b, 0xbe, 0xa6,
	0xe7, 0x1b, 0xb4, 0x10, 0x96, 0x28, 0x69, 0x70, 0xa9, 0xe1, 0x77, 0x7e, 0x50, 0x6a, 0x71, 0x41,
	0x7e, 0x5e, 0x71, 0xaa, 0xd2, 0x31, 0x46, 0x2e, 0x2d, 0xac, 0x4a, 0xc1, 0xca, 0x5c, 0x52, 0x4b,
	0x12, 0x33, 0x73, 0x8a, 0xdd, 0xf2, 0x8b, 0x20, 0xb2, 0x64, 0xfb, 0x3a, 0x00, 0xcd, 0xd7, 0x46,
	0x78, 0x7d, 0x8d, 0x66, 0x3d, 0x51, 0x5e, 0x36, 0xe1, 0x32, 0x22, 0xde, 0x1f, 0x30, 0xef, 0x1b,
	0x5d, 0x65, 0xe2, 0x62, 0xf6, 0x2d, 0x4e, 0x17, 0x9a, 0xcc, 0xc8, 0x25, 0x8d, 0x2f, 0xb6, 0x71,
	0xc6, 0x12, 0xde, 0x50, 0x96, 0xb2, 0x25, 0x4b, 0x1b, 0xcc, 0x75, 0x42, 0xdb, 0x19, 0xb9, 0xd4,
	0x89, 0x8d, 0x19, 0x47, 0x52, 0xac, 0xc2, 0x6a, 0x84, 0x94, 0x27, 0xc5, 0x46, 0xc0, 0x5c, 0xee,
	0x14, 0x72, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0x4e, 0x78,
	0x2c, 0xc7, 0x70, 0xe1, 0xb1, 0x1c, 0xc3, 0x8d, 0xc7, 0x72, 0x0c, 0x51, 0x56, 0xe9, 0x99, 0x25,
	0x19, 0xa5, 0x49, 0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0x28, 0x39, 0xb1, 0xcc, 0x44, 0x37, 0x39, 0x23,
	0x31, 0x33, 0x4f, 0x1f, 0x2e, 0x52, 0x81, 0x54, 0x3c, 0x54, 0x16, 0xa4, 0x16, 0x27, 0xb1, 0x81,
	0xa5, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x1c, 0xa8, 0xca, 0x89, 0x44, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// SetMarketMapperRevenueShare sets the revenue share for a market
	// mapper.
	SetMarketMapperRevenueShare(ctx context.Context, in *MsgSetMarketMapperRevenueShare, opts ...grpc.CallOption) (*MsgSetMarketMapperRevenueShareResponse, error)
	SetMarketMapperRevShareDetailsForMarket(ctx context.Context, in *MsgSetMarketMapperRevShareDetailsForMarket, opts ...grpc.CallOption) (*MsgSetMarketMapperRevShareDetailsForMarketResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) SetMarketMapperRevenueShare(ctx context.Context, in *MsgSetMarketMapperRevenueShare, opts ...grpc.CallOption) (*MsgSetMarketMapperRevenueShareResponse, error) {
	out := new(MsgSetMarketMapperRevenueShareResponse)
	err := c.cc.Invoke(ctx, "/dydxprotocol.revshare.Msg/SetMarketMapperRevenueShare", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SetMarketMapperRevShareDetailsForMarket(ctx context.Context, in *MsgSetMarketMapperRevShareDetailsForMarket, opts ...grpc.CallOption) (*MsgSetMarketMapperRevShareDetailsForMarketResponse, error) {
	out := new(MsgSetMarketMapperRevShareDetailsForMarketResponse)
	err := c.cc.Invoke(ctx, "/dydxprotocol.revshare.Msg/SetMarketMapperRevShareDetailsForMarket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// SetMarketMapperRevenueShare sets the revenue share for a market
	// mapper.
	SetMarketMapperRevenueShare(context.Context, *MsgSetMarketMapperRevenueShare) (*MsgSetMarketMapperRevenueShareResponse, error)
	SetMarketMapperRevShareDetailsForMarket(context.Context, *MsgSetMarketMapperRevShareDetailsForMarket) (*MsgSetMarketMapperRevShareDetailsForMarketResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) SetMarketMapperRevenueShare(ctx context.Context, req *MsgSetMarketMapperRevenueShare) (*MsgSetMarketMapperRevenueShareResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetMarketMapperRevenueShare not implemented")
}
func (*UnimplementedMsgServer) SetMarketMapperRevShareDetailsForMarket(ctx context.Context, req *MsgSetMarketMapperRevShareDetailsForMarket) (*MsgSetMarketMapperRevShareDetailsForMarketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetMarketMapperRevShareDetailsForMarket not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_SetMarketMapperRevenueShare_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetMarketMapperRevenueShare)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetMarketMapperRevenueShare(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dydxprotocol.revshare.Msg/SetMarketMapperRevenueShare",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetMarketMapperRevenueShare(ctx, req.(*MsgSetMarketMapperRevenueShare))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SetMarketMapperRevShareDetailsForMarket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetMarketMapperRevShareDetailsForMarket)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetMarketMapperRevShareDetailsForMarket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dydxprotocol.revshare.Msg/SetMarketMapperRevShareDetailsForMarket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetMarketMapperRevShareDetailsForMarket(ctx, req.(*MsgSetMarketMapperRevShareDetailsForMarket))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dydxprotocol.revshare.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetMarketMapperRevenueShare",
			Handler:    _Msg_SetMarketMapperRevenueShare_Handler,
		},
		{
			MethodName: "SetMarketMapperRevShareDetailsForMarket",
			Handler:    _Msg_SetMarketMapperRevShareDetailsForMarket_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dydxprotocol/revshare/tx.proto",
}

func (m *MsgSetMarketMapperRevenueShare) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetMarketMapperRevenueShare) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetMarketMapperRevenueShare) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgSetMarketMapperRevenueShareResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetMarketMapperRevenueShareResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetMarketMapperRevenueShareResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgSetMarketMapperRevShareDetailsForMarket) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetMarketMapperRevShareDetailsForMarket) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetMarketMapperRevShareDetailsForMarket) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgSetMarketMapperRevShareDetailsForMarketResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetMarketMapperRevShareDetailsForMarketResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetMarketMapperRevShareDetailsForMarketResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgSetMarketMapperRevenueShare) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Params.Size()
	n += 1 + l + sovTx(uint64(l))
	return n
}

func (m *MsgSetMarketMapperRevenueShareResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgSetMarketMapperRevShareDetailsForMarket) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Params.Size()
	n += 1 + l + sovTx(uint64(l))
	return n
}

func (m *MsgSetMarketMapperRevShareDetailsForMarketResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgSetMarketMapperRevenueShare) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgSetMarketMapperRevenueShare: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetMarketMapperRevenueShare: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
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
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgSetMarketMapperRevenueShareResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgSetMarketMapperRevenueShareResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetMarketMapperRevenueShareResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgSetMarketMapperRevShareDetailsForMarket) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgSetMarketMapperRevShareDetailsForMarket: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetMarketMapperRevShareDetailsForMarket: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
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
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgSetMarketMapperRevShareDetailsForMarketResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgSetMarketMapperRevShareDetailsForMarketResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetMarketMapperRevShareDetailsForMarketResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
