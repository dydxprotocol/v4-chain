// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/revshare/query.proto

package types

import (
	context "context"
	fmt "fmt"
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

// Queries market mapper revenue share details for a specific market
type QueryMarketMapperRevShareDetails struct {
	MarketId uint32 `protobuf:"varint,1,opt,name=market_id,json=marketId,proto3" json:"market_id,omitempty"`
}

func (m *QueryMarketMapperRevShareDetails) Reset()         { *m = QueryMarketMapperRevShareDetails{} }
func (m *QueryMarketMapperRevShareDetails) String() string { return proto.CompactTextString(m) }
func (*QueryMarketMapperRevShareDetails) ProtoMessage()    {}
func (*QueryMarketMapperRevShareDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_13d50c6e3048e744, []int{0}
}
func (m *QueryMarketMapperRevShareDetails) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryMarketMapperRevShareDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryMarketMapperRevShareDetails.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryMarketMapperRevShareDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryMarketMapperRevShareDetails.Merge(m, src)
}
func (m *QueryMarketMapperRevShareDetails) XXX_Size() int {
	return m.Size()
}
func (m *QueryMarketMapperRevShareDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryMarketMapperRevShareDetails.DiscardUnknown(m)
}

var xxx_messageInfo_QueryMarketMapperRevShareDetails proto.InternalMessageInfo

func (m *QueryMarketMapperRevShareDetails) GetMarketId() uint32 {
	if m != nil {
		return m.MarketId
	}
	return 0
}

// Response type for QueryMarketMapperRevShareDetails
type QueryMarketMapperRevShareDetailsResponse struct {
	Details *MarketMapperRevShareDetails `protobuf:"bytes,1,opt,name=details,proto3" json:"details,omitempty"`
}

func (m *QueryMarketMapperRevShareDetailsResponse) Reset() {
	*m = QueryMarketMapperRevShareDetailsResponse{}
}
func (m *QueryMarketMapperRevShareDetailsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryMarketMapperRevShareDetailsResponse) ProtoMessage()    {}
func (*QueryMarketMapperRevShareDetailsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_13d50c6e3048e744, []int{1}
}
func (m *QueryMarketMapperRevShareDetailsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryMarketMapperRevShareDetailsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryMarketMapperRevShareDetailsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryMarketMapperRevShareDetailsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryMarketMapperRevShareDetailsResponse.Merge(m, src)
}
func (m *QueryMarketMapperRevShareDetailsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryMarketMapperRevShareDetailsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryMarketMapperRevShareDetailsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryMarketMapperRevShareDetailsResponse proto.InternalMessageInfo

func (m *QueryMarketMapperRevShareDetailsResponse) GetDetails() *MarketMapperRevShareDetails {
	if m != nil {
		return m.Details
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryMarketMapperRevShareDetails)(nil), "dydxprotocol.revshare.QueryMarketMapperRevShareDetails")
	proto.RegisterType((*QueryMarketMapperRevShareDetailsResponse)(nil), "dydxprotocol.revshare.QueryMarketMapperRevShareDetailsResponse")
}

func init() { proto.RegisterFile("dydxprotocol/revshare/query.proto", fileDescriptor_13d50c6e3048e744) }

var fileDescriptor_13d50c6e3048e744 = []byte{
	// 314 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4c, 0xa9, 0x4c, 0xa9,
	0x28, 0x28, 0xca, 0x2f, 0xc9, 0x4f, 0xce, 0xcf, 0xd1, 0x2f, 0x4a, 0x2d, 0x2b, 0xce, 0x48, 0x2c,
	0x4a, 0xd5, 0x2f, 0x2c, 0x4d, 0x2d, 0xaa, 0xd4, 0x03, 0x8b, 0x0b, 0x89, 0x22, 0x2b, 0xd1, 0x83,
	0x29, 0x91, 0x92, 0x49, 0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0xd5, 0x4f, 0x2c, 0xc8, 0xd4, 0x4f, 0xcc,
	0xcb, 0xcb, 0x2f, 0x49, 0x2c, 0xc9, 0xcc, 0xcf, 0x2b, 0x86, 0x68, 0x92, 0x52, 0xc1, 0x6e, 0x2e,
	0x8c, 0x01, 0x51, 0xa5, 0x64, 0xcf, 0xa5, 0x10, 0x08, 0xb2, 0xc9, 0x37, 0xb1, 0x28, 0x3b, 0xb5,
	0xc4, 0x37, 0xb1, 0xa0, 0x20, 0xb5, 0x28, 0x28, 0xb5, 0x2c, 0x18, 0xa4, 0xc4, 0x25, 0xb5, 0x24,
	0x31, 0x33, 0xa7, 0x58, 0x48, 0x9a, 0x8b, 0x33, 0x17, 0x2c, 0x1d, 0x9f, 0x99, 0x22, 0xc1, 0xa8,
	0xc0, 0xa8, 0xc1, 0x1b, 0xc4, 0x01, 0x11, 0xf0, 0x4c, 0x51, 0xaa, 0xe0, 0xd2, 0x20, 0x64, 0x40,
	0x50, 0x6a, 0x71, 0x41, 0x7e, 0x5e, 0x71, 0xaa, 0x90, 0x0f, 0x17, 0x7b, 0x0a, 0x44, 0x08, 0x6c,
	0x0c, 0xb7, 0x91, 0x91, 0x1e, 0x56, 0x9f, 0xe9, 0xe1, 0x33, 0x0c, 0x66, 0x84, 0xd1, 0x5b, 0x46,
	0x2e, 0x56, 0xb0, 0xd5, 0x42, 0x8f, 0x19, 0xb9, 0xa4, 0xf1, 0x79, 0xc0, 0x1c, 0x87, 0x35, 0x84,
	0x1c, 0x2e, 0x65, 0x4f, 0xa6, 0x46, 0x98, 0x8f, 0x95, 0xbc, 0x9a, 0x2e, 0x3f, 0x99, 0xcc, 0xe4,
	0x22, 0xe4, 0xa4, 0x8f, 0x3d, 0x36, 0xa0, 0xe1, 0x9a, 0x0b, 0x36, 0x24, 0xbe, 0x28, 0xb5, 0x2c,
	0x1e, 0x2c, 0x1e, 0x0f, 0xf5, 0xa3, 0x7e, 0x35, 0x3c, 0xe0, 0x6b, 0x9d, 0x42, 0x4e, 0x3c, 0x92,
	0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18, 0x2e, 0x3c,
	0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0xca, 0x2a, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f,
	0x39, 0x3f, 0x17, 0xd5, 0x9e, 0x32, 0x13, 0xdd, 0xe4, 0x8c, 0xc4, 0xcc, 0x3c, 0x7d, 0xb8, 0x48,
	0x05, 0xc2, 0xee, 0x92, 0xca, 0x82, 0xd4, 0xe2, 0x24, 0x36, 0xb0, 0x94, 0x31, 0x20, 0x00, 0x00,
	0xff, 0xff, 0x18, 0x56, 0xfe, 0x50, 0x87, 0x02, 0x00, 0x00,
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
	// Queries market mapper revenue share details for a specific market
	MarketMapperRevShareDetails(ctx context.Context, in *QueryMarketMapperRevShareDetails, opts ...grpc.CallOption) (*QueryMarketMapperRevShareDetailsResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) MarketMapperRevShareDetails(ctx context.Context, in *QueryMarketMapperRevShareDetails, opts ...grpc.CallOption) (*QueryMarketMapperRevShareDetailsResponse, error) {
	out := new(QueryMarketMapperRevShareDetailsResponse)
	err := c.cc.Invoke(ctx, "/dydxprotocol.revshare.Query/MarketMapperRevShareDetails", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Queries market mapper revenue share details for a specific market
	MarketMapperRevShareDetails(context.Context, *QueryMarketMapperRevShareDetails) (*QueryMarketMapperRevShareDetailsResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) MarketMapperRevShareDetails(ctx context.Context, req *QueryMarketMapperRevShareDetails) (*QueryMarketMapperRevShareDetailsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarketMapperRevShareDetails not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_MarketMapperRevShareDetails_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryMarketMapperRevShareDetails)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).MarketMapperRevShareDetails(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dydxprotocol.revshare.Query/MarketMapperRevShareDetails",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).MarketMapperRevShareDetails(ctx, req.(*QueryMarketMapperRevShareDetails))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dydxprotocol.revshare.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MarketMapperRevShareDetails",
			Handler:    _Query_MarketMapperRevShareDetails_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dydxprotocol/revshare/query.proto",
}

func (m *QueryMarketMapperRevShareDetails) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryMarketMapperRevShareDetails) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryMarketMapperRevShareDetails) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MarketId != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.MarketId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *QueryMarketMapperRevShareDetailsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryMarketMapperRevShareDetailsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryMarketMapperRevShareDetailsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Details != nil {
		{
			size, err := m.Details.MarshalToSizedBuffer(dAtA[:i])
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
func (m *QueryMarketMapperRevShareDetails) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MarketId != 0 {
		n += 1 + sovQuery(uint64(m.MarketId))
	}
	return n
}

func (m *QueryMarketMapperRevShareDetailsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Details != nil {
		l = m.Details.Size()
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
func (m *QueryMarketMapperRevShareDetails) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryMarketMapperRevShareDetails: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryMarketMapperRevShareDetails: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MarketId", wireType)
			}
			m.MarketId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MarketId |= uint32(b&0x7F) << shift
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
func (m *QueryMarketMapperRevShareDetailsResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryMarketMapperRevShareDetailsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryMarketMapperRevShareDetailsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Details", wireType)
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
			if m.Details == nil {
				m.Details = &MarketMapperRevShareDetails{}
			}
			if err := m.Details.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
