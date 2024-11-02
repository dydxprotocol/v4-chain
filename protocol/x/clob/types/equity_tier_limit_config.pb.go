// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: klyraprotocol/clob/equity_tier_limit_config.proto

package types

import (
	fmt "fmt"
	github_com_StreamFinance_Protocol_stream_chain_protocol_dtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
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

// Defines the set of equity tiers to limit how many open orders
// a subaccount is allowed to have.
type EquityTierLimitConfiguration struct {
	// How many short term stateful orders are allowed per equity tier.
	// Specifying 0 values disables this limit.
	ShortTermOrderEquityTiers []EquityTierLimit `protobuf:"bytes,1,rep,name=short_term_order_equity_tiers,json=shortTermOrderEquityTiers,proto3" json:"short_term_order_equity_tiers"`
	// How many open stateful orders are allowed per equity tier.
	// Specifying 0 values disables this limit.
	StatefulOrderEquityTiers []EquityTierLimit `protobuf:"bytes,2,rep,name=stateful_order_equity_tiers,json=statefulOrderEquityTiers,proto3" json:"stateful_order_equity_tiers"`
}

func (m *EquityTierLimitConfiguration) Reset()         { *m = EquityTierLimitConfiguration{} }
func (m *EquityTierLimitConfiguration) String() string { return proto.CompactTextString(m) }
func (*EquityTierLimitConfiguration) ProtoMessage()    {}
func (*EquityTierLimitConfiguration) Descriptor() ([]byte, []int) {
	return fileDescriptor_b04fb95694171826, []int{0}
}
func (m *EquityTierLimitConfiguration) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EquityTierLimitConfiguration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EquityTierLimitConfiguration.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EquityTierLimitConfiguration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EquityTierLimitConfiguration.Merge(m, src)
}
func (m *EquityTierLimitConfiguration) XXX_Size() int {
	return m.Size()
}
func (m *EquityTierLimitConfiguration) XXX_DiscardUnknown() {
	xxx_messageInfo_EquityTierLimitConfiguration.DiscardUnknown(m)
}

var xxx_messageInfo_EquityTierLimitConfiguration proto.InternalMessageInfo

func (m *EquityTierLimitConfiguration) GetShortTermOrderEquityTiers() []EquityTierLimit {
	if m != nil {
		return m.ShortTermOrderEquityTiers
	}
	return nil
}

func (m *EquityTierLimitConfiguration) GetStatefulOrderEquityTiers() []EquityTierLimit {
	if m != nil {
		return m.StatefulOrderEquityTiers
	}
	return nil
}

// Defines an equity tier limit.
type EquityTierLimit struct {
	// The total net collateral in TDAI quote quantums of equity required.
	UsdTncRequired github_com_StreamFinance_Protocol_stream_chain_protocol_dtypes.SerializableInt `protobuf:"bytes,1,opt,name=usd_tnc_required,json=usdTncRequired,proto3,customtype=github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes.SerializableInt" json:"usd_tnc_required"`
	// What the limit is for `usd_tnc_required`.
	Limit uint32 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (m *EquityTierLimit) Reset()         { *m = EquityTierLimit{} }
func (m *EquityTierLimit) String() string { return proto.CompactTextString(m) }
func (*EquityTierLimit) ProtoMessage()    {}
func (*EquityTierLimit) Descriptor() ([]byte, []int) {
	return fileDescriptor_b04fb95694171826, []int{1}
}
func (m *EquityTierLimit) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EquityTierLimit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EquityTierLimit.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EquityTierLimit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EquityTierLimit.Merge(m, src)
}
func (m *EquityTierLimit) XXX_Size() int {
	return m.Size()
}
func (m *EquityTierLimit) XXX_DiscardUnknown() {
	xxx_messageInfo_EquityTierLimit.DiscardUnknown(m)
}

var xxx_messageInfo_EquityTierLimit proto.InternalMessageInfo

func (m *EquityTierLimit) GetLimit() uint32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func init() {
	proto.RegisterType((*EquityTierLimitConfiguration)(nil), "klyraprotocol.clob.EquityTierLimitConfiguration")
	proto.RegisterType((*EquityTierLimit)(nil), "klyraprotocol.clob.EquityTierLimit")
}

func init() {
	proto.RegisterFile("klyraprotocol/clob/equity_tier_limit_config.proto", fileDescriptor_b04fb95694171826)
}

var fileDescriptor_b04fb95694171826 = []byte{
	// 366 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0xbb, 0x4e, 0xf3, 0x30,
	0x18, 0x86, 0xe3, 0xfe, 0x87, 0x21, 0xff, 0xcf, 0x41, 0x51, 0x87, 0x70, 0x4a, 0xab, 0xb2, 0x74,
	0x69, 0x22, 0xe0, 0x0e, 0xca, 0x41, 0x42, 0x42, 0x80, 0xd2, 0x8a, 0x81, 0x01, 0xcb, 0x71, 0xdc,
	0xc4, 0xaa, 0x63, 0x17, 0xdb, 0x91, 0x28, 0x13, 0x97, 0xc0, 0x75, 0x70, 0x25, 0x1d, 0x3b, 0x22,
	0x86, 0x0a, 0xb5, 0xd7, 0xc0, 0x8e, 0xe2, 0x52, 0x0e, 0x85, 0x05, 0xb6, 0xe4, 0xd3, 0xf7, 0x3e,
	0xcf, 0x27, 0xbd, 0xb6, 0xb7, 0xba, 0xac, 0x2f, 0x51, 0x4f, 0x0a, 0x2d, 0xb0, 0x60, 0x01, 0x66,
	0x22, 0x0a, 0xc8, 0x65, 0x4e, 0x75, 0x1f, 0x6a, 0x4a, 0x24, 0x64, 0x34, 0xa3, 0x1a, 0x62, 0xc1,
	0x3b, 0x34, 0xf1, 0xcd, 0x9a, 0xe3, 0x7c, 0x88, 0xf8, 0x45, 0x64, 0xb5, 0x9c, 0x88, 0x44, 0x98,
	0x51, 0x50, 0x7c, 0x4d, 0x37, 0x6b, 0x4f, 0xc0, 0x5e, 0xdf, 0x37, 0xb0, 0x36, 0x25, 0xf2, 0xa8,
	0x40, 0xed, 0x1a, 0x52, 0x2e, 0x91, 0xa6, 0x82, 0x3b, 0x5d, 0x7b, 0x43, 0xa5, 0x42, 0x6a, 0xa8,
	0x89, 0xcc, 0xa0, 0x90, 0x31, 0x91, 0xf0, 0x9d, 0x5d, 0xb9, 0xa0, 0xfa, 0xab, 0xfe, 0x6f, 0x7b,
	0xd3, 0xff, 0xac, 0xf4, 0xe7, 0xc0, 0xcd, 0xdf, 0x83, 0x51, 0xc5, 0x0a, 0x57, 0x0c, 0xaf, 0x4d,
	0x64, 0x76, 0x52, 0xd0, 0xde, 0x96, 0x94, 0x93, 0xda, 0x6b, 0x4a, 0x23, 0x4d, 0x3a, 0x39, 0xfb,
	0x4a, 0x55, 0xfa, 0xae, 0xca, 0x9d, 0xd1, 0xe6, 0x4d, 0xb5, 0x3b, 0x60, 0x2f, 0xcd, 0x65, 0x9c,
	0x1b, 0x60, 0x2f, 0xe7, 0x2a, 0x86, 0x9a, 0x63, 0x28, 0x0b, 0xb1, 0x24, 0xb1, 0x0b, 0xaa, 0xa0,
	0xfe, 0xbf, 0x79, 0x56, 0xe0, 0x1e, 0x46, 0x95, 0xe3, 0x84, 0xea, 0x34, 0x8f, 0x7c, 0x2c, 0xb2,
	0xa0, 0xa5, 0x25, 0x41, 0xd9, 0x01, 0xe5, 0x88, 0x63, 0xd2, 0x38, 0x9d, 0xf5, 0xa3, 0xcc, 0xb8,
	0x81, 0x53, 0x44, 0x79, 0xf0, 0xda, 0x5a, 0xac, 0xfb, 0x3d, 0xa2, 0xfc, 0x16, 0x91, 0x14, 0x31,
	0x7a, 0x8d, 0x22, 0x46, 0x0e, 0xb9, 0x0e, 0x17, 0x73, 0x15, 0xb7, 0x39, 0x0e, 0x5f, 0x6c, 0x4e,
	0xd9, 0xfe, 0x63, 0xea, 0x74, 0x4b, 0x55, 0x50, 0x5f, 0x08, 0xa7, 0x3f, 0xcd, 0x8b, 0xc1, 0xd8,
	0x03, 0xc3, 0xb1, 0x07, 0x1e, 0xc7, 0x1e, 0xb8, 0x9d, 0x78, 0xd6, 0x70, 0xe2, 0x59, 0xf7, 0x13,
	0xcf, 0x3a, 0xdf, 0xfb, 0xe9, 0x3d, 0x57, 0xd3, 0x77, 0x64, 0xae, 0x8a, 0xfe, 0x9a, 0xf1, 0xce,
	0x73, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8d, 0x59, 0x6e, 0x35, 0x6a, 0x02, 0x00, 0x00,
}

func (m *EquityTierLimitConfiguration) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EquityTierLimitConfiguration) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EquityTierLimitConfiguration) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.StatefulOrderEquityTiers) > 0 {
		for iNdEx := len(m.StatefulOrderEquityTiers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.StatefulOrderEquityTiers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintEquityTierLimitConfig(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.ShortTermOrderEquityTiers) > 0 {
		for iNdEx := len(m.ShortTermOrderEquityTiers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ShortTermOrderEquityTiers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintEquityTierLimitConfig(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *EquityTierLimit) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EquityTierLimit) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EquityTierLimit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Limit != 0 {
		i = encodeVarintEquityTierLimitConfig(dAtA, i, uint64(m.Limit))
		i--
		dAtA[i] = 0x10
	}
	{
		size := m.UsdTncRequired.Size()
		i -= size
		if _, err := m.UsdTncRequired.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintEquityTierLimitConfig(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintEquityTierLimitConfig(dAtA []byte, offset int, v uint64) int {
	offset -= sovEquityTierLimitConfig(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *EquityTierLimitConfiguration) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ShortTermOrderEquityTiers) > 0 {
		for _, e := range m.ShortTermOrderEquityTiers {
			l = e.Size()
			n += 1 + l + sovEquityTierLimitConfig(uint64(l))
		}
	}
	if len(m.StatefulOrderEquityTiers) > 0 {
		for _, e := range m.StatefulOrderEquityTiers {
			l = e.Size()
			n += 1 + l + sovEquityTierLimitConfig(uint64(l))
		}
	}
	return n
}

func (m *EquityTierLimit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.UsdTncRequired.Size()
	n += 1 + l + sovEquityTierLimitConfig(uint64(l))
	if m.Limit != 0 {
		n += 1 + sovEquityTierLimitConfig(uint64(m.Limit))
	}
	return n
}

func sovEquityTierLimitConfig(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEquityTierLimitConfig(x uint64) (n int) {
	return sovEquityTierLimitConfig(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EquityTierLimitConfiguration) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEquityTierLimitConfig
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
			return fmt.Errorf("proto: EquityTierLimitConfiguration: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EquityTierLimitConfiguration: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ShortTermOrderEquityTiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEquityTierLimitConfig
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
				return ErrInvalidLengthEquityTierLimitConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEquityTierLimitConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ShortTermOrderEquityTiers = append(m.ShortTermOrderEquityTiers, EquityTierLimit{})
			if err := m.ShortTermOrderEquityTiers[len(m.ShortTermOrderEquityTiers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StatefulOrderEquityTiers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEquityTierLimitConfig
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
				return ErrInvalidLengthEquityTierLimitConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEquityTierLimitConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StatefulOrderEquityTiers = append(m.StatefulOrderEquityTiers, EquityTierLimit{})
			if err := m.StatefulOrderEquityTiers[len(m.StatefulOrderEquityTiers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEquityTierLimitConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEquityTierLimitConfig
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
func (m *EquityTierLimit) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEquityTierLimitConfig
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
			return fmt.Errorf("proto: EquityTierLimit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EquityTierLimit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UsdTncRequired", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEquityTierLimitConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthEquityTierLimitConfig
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthEquityTierLimitConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.UsdTncRequired.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Limit", wireType)
			}
			m.Limit = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEquityTierLimitConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Limit |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipEquityTierLimitConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEquityTierLimitConfig
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
func skipEquityTierLimitConfig(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEquityTierLimitConfig
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
					return 0, ErrIntOverflowEquityTierLimitConfig
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
					return 0, ErrIntOverflowEquityTierLimitConfig
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
				return 0, ErrInvalidLengthEquityTierLimitConfig
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEquityTierLimitConfig
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEquityTierLimitConfig
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEquityTierLimitConfig        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEquityTierLimitConfig          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEquityTierLimitConfig = fmt.Errorf("proto: unexpected end of group")
)
