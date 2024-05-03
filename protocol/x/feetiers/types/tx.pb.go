// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/feetiers/tx.proto

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

// MsgUpdatePerpetualFeeParams is the Msg/UpdatePerpetualFeeParams request type.
type MsgUpdatePerpetualFeeParams struct {
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	// Defines the parameters to update. All parameters must be supplied.
	Params PerpetualFeeParams `protobuf:"bytes,2,opt,name=params,proto3" json:"params"`
}

func (m *MsgUpdatePerpetualFeeParams) Reset()         { *m = MsgUpdatePerpetualFeeParams{} }
func (m *MsgUpdatePerpetualFeeParams) String() string { return proto.CompactTextString(m) }
func (*MsgUpdatePerpetualFeeParams) ProtoMessage()    {}
func (*MsgUpdatePerpetualFeeParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_caa74a3b986b7fd9, []int{0}
}
func (m *MsgUpdatePerpetualFeeParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdatePerpetualFeeParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdatePerpetualFeeParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdatePerpetualFeeParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdatePerpetualFeeParams.Merge(m, src)
}
func (m *MsgUpdatePerpetualFeeParams) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdatePerpetualFeeParams) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdatePerpetualFeeParams.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdatePerpetualFeeParams proto.InternalMessageInfo

func (m *MsgUpdatePerpetualFeeParams) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *MsgUpdatePerpetualFeeParams) GetParams() PerpetualFeeParams {
	if m != nil {
		return m.Params
	}
	return PerpetualFeeParams{}
}

// MsgUpdatePerpetualFeeParamsResponse is the Msg/UpdatePerpetualFeeParams
// response type.
type MsgUpdatePerpetualFeeParamsResponse struct {
}

func (m *MsgUpdatePerpetualFeeParamsResponse) Reset()         { *m = MsgUpdatePerpetualFeeParamsResponse{} }
func (m *MsgUpdatePerpetualFeeParamsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgUpdatePerpetualFeeParamsResponse) ProtoMessage()    {}
func (*MsgUpdatePerpetualFeeParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_caa74a3b986b7fd9, []int{1}
}
func (m *MsgUpdatePerpetualFeeParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdatePerpetualFeeParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdatePerpetualFeeParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdatePerpetualFeeParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdatePerpetualFeeParamsResponse.Merge(m, src)
}
func (m *MsgUpdatePerpetualFeeParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdatePerpetualFeeParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdatePerpetualFeeParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdatePerpetualFeeParamsResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgUpdatePerpetualFeeParams)(nil), "dydxprotocol.feetiers.MsgUpdatePerpetualFeeParams")
	proto.RegisterType((*MsgUpdatePerpetualFeeParamsResponse)(nil), "dydxprotocol.feetiers.MsgUpdatePerpetualFeeParamsResponse")
}

func init() { proto.RegisterFile("dydxprotocol/feetiers/tx.proto", fileDescriptor_caa74a3b986b7fd9) }

var fileDescriptor_caa74a3b986b7fd9 = []byte{
	// 354 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4b, 0xa9, 0x4c, 0xa9,
	0x28, 0x28, 0xca, 0x2f, 0xc9, 0x4f, 0xce, 0xcf, 0xd1, 0x4f, 0x4b, 0x4d, 0x2d, 0xc9, 0x4c, 0x2d,
	0x2a, 0xd6, 0x2f, 0xa9, 0xd0, 0x03, 0x0b, 0x0a, 0x89, 0x22, 0xcb, 0xeb, 0xc1, 0xe4, 0xa5, 0x24,
	0x93, 0xf3, 0x8b, 0x73, 0xf3, 0x8b, 0xe3, 0xc1, 0x32, 0xfa, 0x10, 0x0e, 0x44, 0x87, 0x94, 0x38,
	0x84, 0xa7, 0x9f, 0x5b, 0x9c, 0xae, 0x5f, 0x66, 0x08, 0xa2, 0xa0, 0x12, 0x4a, 0xd8, 0xad, 0x2a,
	0x48, 0x2c, 0x4a, 0xcc, 0x85, 0x69, 0x16, 0x49, 0xcf, 0x4f, 0xcf, 0x87, 0x18, 0x0a, 0x62, 0x41,
	0x44, 0x95, 0xd6, 0x31, 0x72, 0x49, 0xfb, 0x16, 0xa7, 0x87, 0x16, 0xa4, 0x24, 0x96, 0xa4, 0x06,
	0xa4, 0x16, 0x15, 0xa4, 0x96, 0x94, 0x26, 0xe6, 0xb8, 0xa5, 0xa6, 0x06, 0x80, 0xf5, 0x0a, 0x99,
	0x71, 0x71, 0x26, 0x96, 0x96, 0x64, 0xe4, 0x17, 0x65, 0x96, 0x54, 0x4a, 0x30, 0x2a, 0x30, 0x6a,
	0x70, 0x3a, 0x49, 0x5c, 0xda, 0xa2, 0x2b, 0x02, 0x75, 0x97, 0x63, 0x4a, 0x4a, 0x51, 0x6a, 0x71,
	0x71, 0x70, 0x49, 0x51, 0x66, 0x5e, 0x7a, 0x10, 0x42, 0xa9, 0x90, 0x3b, 0x17, 0x1b, 0xc4, 0x76,
	0x09, 0x26, 0x05, 0x46, 0x0d, 0x6e, 0x23, 0x4d, 0x3d, 0xac, 0xbe, 0xd5, 0xc3, 0xb4, 0xd2, 0x89,
	0xe5, 0xc4, 0x3d, 0x79, 0x86, 0x20, 0xa8, 0x76, 0x2b, 0xbe, 0xa6, 0xe7, 0x1b, 0xb4, 0x10, 0x06,
	0x2b, 0xa9, 0x72, 0x29, 0xe3, 0x71, 0x6f, 0x50, 0x6a, 0x71, 0x41, 0x7e, 0x5e, 0x71, 0xaa, 0xd1,
	0x4c, 0x46, 0x2e, 0x66, 0xdf, 0xe2, 0x74, 0xa1, 0x2e, 0x46, 0x2e, 0x09, 0x9c, 0x9e, 0x33, 0xc2,
	0xe1, 0x28, 0x3c, 0x16, 0x48, 0x59, 0x91, 0xae, 0x07, 0xe6, 0x28, 0x29, 0xd6, 0x86, 0xe7, 0x1b,
	0xb4, 0x18, 0x9d, 0x92, 0x4e, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39,
	0xc6, 0x09, 0x8f, 0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0xca, 0x23,
	0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0x3f, 0xb8, 0xa4, 0x28, 0x35, 0x31,
	0xd7, 0x2d, 0x33, 0x2f, 0x31, 0x2f, 0x39, 0x55, 0x37, 0x00, 0x16, 0xb9, 0xc5, 0x60, 0x61, 0xdd,
	0xe4, 0x8c, 0xc4, 0xcc, 0x3c, 0x7d, 0x78, 0x94, 0x57, 0x20, 0xa5, 0xaf, 0xca, 0x82, 0xd4, 0xe2,
	0x24, 0x36, 0xb0, 0x94, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x67, 0x2e, 0x23, 0xa6, 0x85, 0x02,
	0x00, 0x00,
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
	// UpdatePerpetualFeeParams updates the PerpetualFeeParams in state.
	UpdatePerpetualFeeParams(ctx context.Context, in *MsgUpdatePerpetualFeeParams, opts ...grpc.CallOption) (*MsgUpdatePerpetualFeeParamsResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) UpdatePerpetualFeeParams(ctx context.Context, in *MsgUpdatePerpetualFeeParams, opts ...grpc.CallOption) (*MsgUpdatePerpetualFeeParamsResponse, error) {
	out := new(MsgUpdatePerpetualFeeParamsResponse)
	err := c.cc.Invoke(ctx, "/dydxprotocol.feetiers.Msg/UpdatePerpetualFeeParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// UpdatePerpetualFeeParams updates the PerpetualFeeParams in state.
	UpdatePerpetualFeeParams(context.Context, *MsgUpdatePerpetualFeeParams) (*MsgUpdatePerpetualFeeParamsResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) UpdatePerpetualFeeParams(ctx context.Context, req *MsgUpdatePerpetualFeeParams) (*MsgUpdatePerpetualFeeParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePerpetualFeeParams not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_UpdatePerpetualFeeParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdatePerpetualFeeParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdatePerpetualFeeParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dydxprotocol.feetiers.Msg/UpdatePerpetualFeeParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdatePerpetualFeeParams(ctx, req.(*MsgUpdatePerpetualFeeParams))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dydxprotocol.feetiers.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdatePerpetualFeeParams",
			Handler:    _Msg_UpdatePerpetualFeeParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dydxprotocol/feetiers/tx.proto",
}

func (m *MsgUpdatePerpetualFeeParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdatePerpetualFeeParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdatePerpetualFeeParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *MsgUpdatePerpetualFeeParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdatePerpetualFeeParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdatePerpetualFeeParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
func (m *MsgUpdatePerpetualFeeParams) Size() (n int) {
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

func (m *MsgUpdatePerpetualFeeParamsResponse) Size() (n int) {
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
func (m *MsgUpdatePerpetualFeeParams) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgUpdatePerpetualFeeParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdatePerpetualFeeParams: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgUpdatePerpetualFeeParamsResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgUpdatePerpetualFeeParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdatePerpetualFeeParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
