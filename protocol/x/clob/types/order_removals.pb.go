// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/clob/order_removals.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

type OrderRemoval_RemovalReason int32

const (
	// REMOVAL_REASON_UNSPECIFIED represents an unspecified removal reason. This
	// removal reason is used as a catchall and should never appear on an
	// OrderRemoval in the operations queue.
	OrderRemoval_REMOVAL_REASON_UNSPECIFIED OrderRemoval_RemovalReason = 0
	// REMOVAL_REASON_UNDERCOLLATERALIZED represents a removal of an order which
	// if filled in isolation with respect to the current state of the
	// subaccount would leave the subaccount undercollateralized.
	OrderRemoval_REMOVAL_REASON_UNDERCOLLATERALIZED OrderRemoval_RemovalReason = 1
	// REMOVAL_REASON_INVALID_REDUCE_ONLY represents a removal of a reduce-only
	// order which if filled in isolation with respect to the current state of
	// the subaccount would cause the subaccount's existing position to increase
	// or change sides.
	OrderRemoval_REMOVAL_REASON_INVALID_REDUCE_ONLY OrderRemoval_RemovalReason = 2
	// REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER represents a removal of
	// a stateful post-only order that was deemed invalid because it crossed
	// maker orders on the book of the proposer.
	OrderRemoval_REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER OrderRemoval_RemovalReason = 3
	// REMOVAL_REASON_INVALID_SELF_TRADE represents a removal of a stateful
	// order that was deemed invalid because it constituted a self trade on the
	// proposers orderbook.
	OrderRemoval_REMOVAL_REASON_INVALID_SELF_TRADE OrderRemoval_RemovalReason = 4
	// REMOVAL_REASON_CONDITIONAL_FOK_COULD_NOT_BE_FULLY_FILLED represents a
	// removal of a conditional FOK order that was deemed invalid because it
	// could not be completely filled. Conditional FOK orders should always be
	// fully-filled or removed in the block after they are triggered.
	OrderRemoval_REMOVAL_REASON_CONDITIONAL_FOK_COULD_NOT_BE_FULLY_FILLED OrderRemoval_RemovalReason = 5
	// REMOVAL_REASON_CONDITIONAL_IOC_WOULD_REST_ON_BOOK represents a removal
	// of a conditional IOC order.
	// Conditional IOC orders should always have their remaining size removed
	// in the block after they are triggered.
	OrderRemoval_REMOVAL_REASON_CONDITIONAL_IOC_WOULD_REST_ON_BOOK OrderRemoval_RemovalReason = 6
	// REMOVAL_REASON_FULLY_FILLED represents a removal of an order that
	// was fully filled and should therefore be removed from state.
	OrderRemoval_REMOVAL_REASON_FULLY_FILLED OrderRemoval_RemovalReason = 7
	// REMOVAL_REASON_FULLY_FILLED represents a removal of an order that
	//  would lead to the subaccount violating isolated subaccount constraints.
	OrderRemoval_REMOVAL_REASON_VIOLATES_ISOLATED_SUBACCOUNT_CONSTRAINTS OrderRemoval_RemovalReason = 8
	// REMOVAL_REASON_PERMISSIONED_KEY_EXPIRED represents a removal of an order
	// that was placed using an expired permissioned key.
	OrderRemoval_REMOVAL_REASON_PERMISSIONED_KEY_EXPIRED OrderRemoval_RemovalReason = 9
)

var OrderRemoval_RemovalReason_name = map[int32]string{
	0: "REMOVAL_REASON_UNSPECIFIED",
	1: "REMOVAL_REASON_UNDERCOLLATERALIZED",
	2: "REMOVAL_REASON_INVALID_REDUCE_ONLY",
	3: "REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER",
	4: "REMOVAL_REASON_INVALID_SELF_TRADE",
	5: "REMOVAL_REASON_CONDITIONAL_FOK_COULD_NOT_BE_FULLY_FILLED",
	6: "REMOVAL_REASON_CONDITIONAL_IOC_WOULD_REST_ON_BOOK",
	7: "REMOVAL_REASON_FULLY_FILLED",
	8: "REMOVAL_REASON_VIOLATES_ISOLATED_SUBACCOUNT_CONSTRAINTS",
	9: "REMOVAL_REASON_PERMISSIONED_KEY_EXPIRED",
}

var OrderRemoval_RemovalReason_value = map[string]int32{
	"REMOVAL_REASON_UNSPECIFIED":                               0,
	"REMOVAL_REASON_UNDERCOLLATERALIZED":                       1,
	"REMOVAL_REASON_INVALID_REDUCE_ONLY":                       2,
	"REMOVAL_REASON_POST_ONLY_WOULD_CROSS_MAKER_ORDER":         3,
	"REMOVAL_REASON_INVALID_SELF_TRADE":                        4,
	"REMOVAL_REASON_CONDITIONAL_FOK_COULD_NOT_BE_FULLY_FILLED": 5,
	"REMOVAL_REASON_CONDITIONAL_IOC_WOULD_REST_ON_BOOK":        6,
	"REMOVAL_REASON_FULLY_FILLED":                              7,
	"REMOVAL_REASON_VIOLATES_ISOLATED_SUBACCOUNT_CONSTRAINTS":  8,
	"REMOVAL_REASON_PERMISSIONED_KEY_EXPIRED":                  9,
}

func (x OrderRemoval_RemovalReason) String() string {
	return proto.EnumName(OrderRemoval_RemovalReason_name, int32(x))
}

func (OrderRemoval_RemovalReason) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_60fa12f781955c9f, []int{0, 0}
}

// OrderRemoval is a request type used for forced removal of stateful orders.
type OrderRemoval struct {
	OrderId       OrderId                    `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id"`
	RemovalReason OrderRemoval_RemovalReason `protobuf:"varint,2,opt,name=removal_reason,json=removalReason,proto3,enum=dydxprotocol.clob.OrderRemoval_RemovalReason" json:"removal_reason,omitempty"`
}

func (m *OrderRemoval) Reset()         { *m = OrderRemoval{} }
func (m *OrderRemoval) String() string { return proto.CompactTextString(m) }
func (*OrderRemoval) ProtoMessage()    {}
func (*OrderRemoval) Descriptor() ([]byte, []int) {
	return fileDescriptor_60fa12f781955c9f, []int{0}
}
func (m *OrderRemoval) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderRemoval) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderRemoval.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderRemoval) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderRemoval.Merge(m, src)
}
func (m *OrderRemoval) XXX_Size() int {
	return m.Size()
}
func (m *OrderRemoval) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderRemoval.DiscardUnknown(m)
}

var xxx_messageInfo_OrderRemoval proto.InternalMessageInfo

func (m *OrderRemoval) GetOrderId() OrderId {
	if m != nil {
		return m.OrderId
	}
	return OrderId{}
}

func (m *OrderRemoval) GetRemovalReason() OrderRemoval_RemovalReason {
	if m != nil {
		return m.RemovalReason
	}
	return OrderRemoval_REMOVAL_REASON_UNSPECIFIED
}

func init() {
	proto.RegisterEnum("dydxprotocol.clob.OrderRemoval_RemovalReason", OrderRemoval_RemovalReason_name, OrderRemoval_RemovalReason_value)
	proto.RegisterType((*OrderRemoval)(nil), "dydxprotocol.clob.OrderRemoval")
}

func init() {
	proto.RegisterFile("dydxprotocol/clob/order_removals.proto", fileDescriptor_60fa12f781955c9f)
}

var fileDescriptor_60fa12f781955c9f = []byte{
	// 500 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0x5d, 0x6b, 0x13, 0x41,
	0x14, 0x86, 0xb3, 0xfd, 0x76, 0xb4, 0x65, 0x1d, 0xbc, 0x28, 0x11, 0xb7, 0x35, 0x60, 0x2d, 0x48,
	0x37, 0x5a, 0xeb, 0x07, 0xd4, 0x9b, 0xcd, 0xce, 0x09, 0x0c, 0x99, 0xec, 0x84, 0x99, 0xdd, 0x68,
	0x7a, 0x73, 0xc8, 0x17, 0x69, 0x20, 0xed, 0x94, 0x4d, 0x2c, 0xed, 0x9d, 0x3f, 0xc1, 0x9f, 0xd5,
	0xcb, 0x5e, 0x7a, 0x21, 0x22, 0xc9, 0x1f, 0x91, 0xec, 0x06, 0x69, 0xb6, 0xc6, 0xab, 0x99, 0x73,
	0xce, 0xf3, 0xbe, 0xef, 0x61, 0x18, 0xb2, 0xd7, 0xb9, 0xee, 0x5c, 0x5d, 0xc4, 0x66, 0x64, 0xda,
	0x66, 0x50, 0x6c, 0x0f, 0x4c, 0xab, 0x68, 0xe2, 0x4e, 0x37, 0xc6, 0xb8, 0x7b, 0x66, 0x2e, 0x9b,
	0x83, 0xa1, 0x9b, 0x0c, 0xe9, 0xe3, 0xbb, 0x9c, 0x3b, 0xe5, 0xf2, 0x4f, 0x7a, 0xa6, 0x67, 0x92,
	0x56, 0x71, 0x7a, 0x4b, 0xc1, 0xfc, 0xb3, 0x05, 0x86, 0xe9, 0xb8, 0xf0, 0x6d, 0x95, 0x3c, 0x92,
	0xd3, 0x5a, 0xa5, 0xfe, 0xf4, 0x98, 0x6c, 0xa4, 0x81, 0xfd, 0xce, 0xb6, 0xb5, 0x6b, 0xed, 0x3f,
	0x3c, 0xcc, 0xbb, 0xf7, 0xb2, 0xdc, 0x44, 0xc2, 0x3b, 0xa5, 0x95, 0x9b, 0x5f, 0x3b, 0x39, 0xb5,
	0x6e, 0xd2, 0x92, 0x86, 0x64, 0x6b, 0xb6, 0x27, 0xc6, 0xdd, 0xe6, 0xd0, 0x9c, 0x6f, 0x2f, 0xed,
	0x5a, 0xfb, 0x5b, 0x87, 0x07, 0x8b, 0x2c, 0x66, 0xa9, 0xee, 0xec, 0x54, 0x89, 0x48, 0x6d, 0xc6,
	0x77, 0xcb, 0xc2, 0xcf, 0x65, 0xb2, 0x39, 0x07, 0x50, 0x87, 0xe4, 0x15, 0x54, 0x65, 0xdd, 0x13,
	0xa8, 0xc0, 0xd3, 0x32, 0xc0, 0x28, 0xd0, 0x35, 0xf0, 0x79, 0x99, 0x03, 0xb3, 0x73, 0x74, 0x8f,
	0x14, 0xee, 0xcd, 0x19, 0x28, 0x5f, 0x0a, 0xe1, 0x85, 0xa0, 0x3c, 0xc1, 0x4f, 0x80, 0xd9, 0xd6,
	0x3f, 0x38, 0x1e, 0xd4, 0x3d, 0xc1, 0x19, 0x2a, 0x60, 0x91, 0x0f, 0x28, 0x03, 0xd1, 0xb0, 0x97,
	0xe8, 0x11, 0x79, 0x9d, 0xe1, 0x6a, 0x52, 0x87, 0xc9, 0x14, 0x3f, 0xcb, 0x48, 0x30, 0xf4, 0x95,
	0xd4, 0x1a, 0xab, 0x5e, 0x05, 0x14, 0x4a, 0xc5, 0x40, 0xd9, 0xcb, 0xf4, 0x05, 0x79, 0xbe, 0xc0,
	0x5d, 0x83, 0x28, 0x63, 0xa8, 0x3c, 0x06, 0xf6, 0x0a, 0xfd, 0x44, 0x3e, 0x66, 0x30, 0x5f, 0x06,
	0x8c, 0x87, 0x5c, 0x06, 0x9e, 0xc0, 0xb2, 0xac, 0xa0, 0x9f, 0x44, 0x04, 0x32, 0xc4, 0x12, 0x60,
	0x39, 0x12, 0xa2, 0x81, 0x65, 0x2e, 0x04, 0x30, 0x7b, 0x95, 0xbe, 0x23, 0x6f, 0xfe, 0xa3, 0xe6,
	0xd2, 0x9f, 0x2d, 0xa8, 0x20, 0x59, 0x18, 0x4b, 0x52, 0x56, 0xec, 0x35, 0xba, 0x43, 0x9e, 0x66,
	0x64, 0x73, 0xbe, 0xeb, 0xf4, 0x98, 0x7c, 0xc8, 0x00, 0x75, 0x2e, 0xa7, 0xaf, 0xa7, 0x91, 0xeb,
	0xe4, 0xc2, 0x50, 0x47, 0x25, 0xcf, 0xf7, 0x65, 0x14, 0x84, 0xd3, 0x50, 0x1d, 0x2a, 0x8f, 0x07,
	0xa1, 0xb6, 0x37, 0xe8, 0x2b, 0xf2, 0x32, 0xfb, 0x5e, 0xa0, 0xaa, 0x5c, 0x6b, 0x2e, 0x03, 0x60,
	0x58, 0x81, 0x06, 0xc2, 0x97, 0x1a, 0x57, 0xc0, 0xec, 0x07, 0xa5, 0xda, 0xcd, 0xd8, 0xb1, 0x6e,
	0xc7, 0x8e, 0xf5, 0x7b, 0xec, 0x58, 0xdf, 0x27, 0x4e, 0xee, 0x76, 0xe2, 0xe4, 0x7e, 0x4c, 0x9c,
	0xdc, 0xc9, 0xfb, 0x5e, 0x7f, 0x74, 0xfa, 0xb5, 0xe5, 0xb6, 0xcd, 0x59, 0x71, 0xee, 0x1b, 0x5f,
	0x1e, 0x1d, 0xb4, 0x4f, 0x9b, 0xfd, 0xf3, 0xe2, 0xdf, 0xce, 0x55, 0xfa, 0xb5, 0x47, 0xd7, 0x17,
	0xdd, 0x61, 0x6b, 0x2d, 0x69, 0xbf, 0xfd, 0x13, 0x00, 0x00, 0xff, 0xff, 0xba, 0xd3, 0xd2, 0xeb,
	0x4d, 0x03, 0x00, 0x00,
}

func (m *OrderRemoval) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderRemoval) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderRemoval) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.RemovalReason != 0 {
		i = encodeVarintOrderRemovals(dAtA, i, uint64(m.RemovalReason))
		i--
		dAtA[i] = 0x10
	}
	{
		size, err := m.OrderId.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintOrderRemovals(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintOrderRemovals(dAtA []byte, offset int, v uint64) int {
	offset -= sovOrderRemovals(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *OrderRemoval) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.OrderId.Size()
	n += 1 + l + sovOrderRemovals(uint64(l))
	if m.RemovalReason != 0 {
		n += 1 + sovOrderRemovals(uint64(m.RemovalReason))
	}
	return n
}

func sovOrderRemovals(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozOrderRemovals(x uint64) (n int) {
	return sovOrderRemovals(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *OrderRemoval) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOrderRemovals
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
			return fmt.Errorf("proto: OrderRemoval: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderRemoval: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderId", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrderRemovals
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
				return ErrInvalidLengthOrderRemovals
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOrderRemovals
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.OrderId.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RemovalReason", wireType)
			}
			m.RemovalReason = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrderRemovals
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RemovalReason |= OrderRemoval_RemovalReason(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipOrderRemovals(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthOrderRemovals
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
func skipOrderRemovals(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowOrderRemovals
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
					return 0, ErrIntOverflowOrderRemovals
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
					return 0, ErrIntOverflowOrderRemovals
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
				return 0, ErrInvalidLengthOrderRemovals
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupOrderRemovals
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthOrderRemovals
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthOrderRemovals        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowOrderRemovals          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupOrderRemovals = fmt.Errorf("proto: unexpected end of group")
)
