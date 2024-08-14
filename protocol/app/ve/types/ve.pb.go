// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/ve/ve.proto

package types

import (
	fmt "fmt"
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

// PricePair defines a pair of prices for a market.
type PricePair struct {
	SpotPrice []byte `protobuf:"bytes,1,opt,name=spot_price,json=spotPrice,proto3" json:"spot_price,omitempty"`
	PnlPrice  []byte `protobuf:"bytes,2,opt,name=pnl_price,json=pnlPrice,proto3" json:"pnl_price,omitempty"`
}

func (m *PricePair) Reset()         { *m = PricePair{} }
func (m *PricePair) String() string { return proto.CompactTextString(m) }
func (*PricePair) ProtoMessage()    {}
func (*PricePair) Descriptor() ([]byte, []int) {
	return fileDescriptor_fac2326008e9fb0f, []int{0}
}
func (m *PricePair) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PricePair) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PricePair.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PricePair) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PricePair.Merge(m, src)
}
func (m *PricePair) XXX_Size() int {
	return m.Size()
}
func (m *PricePair) XXX_DiscardUnknown() {
	xxx_messageInfo_PricePair.DiscardUnknown(m)
}

var xxx_messageInfo_PricePair proto.InternalMessageInfo

func (m *PricePair) GetSpotPrice() []byte {
	if m != nil {
		return m.SpotPrice
	}
	return nil
}

func (m *PricePair) GetPnlPrice() []byte {
	if m != nil {
		return m.PnlPrice
	}
	return nil
}

// Daemon VoteExtension defines the vote extension structure for daemon prices.
type DaemonVoteExtension struct {
	// Prices defines a map of marketId -> PricePair.
	Prices map[uint32]*PricePair `protobuf:"bytes,1,rep,name=prices,proto3" json:"prices,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *DaemonVoteExtension) Reset()         { *m = DaemonVoteExtension{} }
func (m *DaemonVoteExtension) String() string { return proto.CompactTextString(m) }
func (*DaemonVoteExtension) ProtoMessage()    {}
func (*DaemonVoteExtension) Descriptor() ([]byte, []int) {
	return fileDescriptor_fac2326008e9fb0f, []int{1}
}
func (m *DaemonVoteExtension) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DaemonVoteExtension) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DaemonVoteExtension.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DaemonVoteExtension) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DaemonVoteExtension.Merge(m, src)
}
func (m *DaemonVoteExtension) XXX_Size() int {
	return m.Size()
}
func (m *DaemonVoteExtension) XXX_DiscardUnknown() {
	xxx_messageInfo_DaemonVoteExtension.DiscardUnknown(m)
}

var xxx_messageInfo_DaemonVoteExtension proto.InternalMessageInfo

func (m *DaemonVoteExtension) GetPrices() map[uint32]*PricePair {
	if m != nil {
		return m.Prices
	}
	return nil
}

func init() {
	proto.RegisterType((*PricePair)(nil), "dydxprotocol.ve.PricePair")
	proto.RegisterType((*DaemonVoteExtension)(nil), "dydxprotocol.ve.DaemonVoteExtension")
	proto.RegisterMapType((map[uint32]*PricePair)(nil), "dydxprotocol.ve.DaemonVoteExtension.PricesEntry")
}

func init() { proto.RegisterFile("dydxprotocol/ve/ve.proto", fileDescriptor_fac2326008e9fb0f) }

var fileDescriptor_fac2326008e9fb0f = []byte{
	// 295 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x48, 0xa9, 0x4c, 0xa9,
	0x28, 0x28, 0xca, 0x2f, 0xc9, 0x4f, 0xce, 0xcf, 0xd1, 0x2f, 0x4b, 0xd5, 0x2f, 0x4b, 0xd5, 0x03,
	0x73, 0x85, 0xf8, 0x91, 0x65, 0xf4, 0xca, 0x52, 0x95, 0xdc, 0xb9, 0x38, 0x03, 0x8a, 0x32, 0x93,
	0x53, 0x03, 0x12, 0x33, 0x8b, 0x84, 0x64, 0xb9, 0xb8, 0x8a, 0x0b, 0xf2, 0x4b, 0xe2, 0x0b, 0x40,
	0x22, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x3c, 0x41, 0x9c, 0x20, 0x11, 0xb0, 0x12, 0x21, 0x69, 0x2e,
	0xce, 0x82, 0xbc, 0x1c, 0xa8, 0x2c, 0x13, 0x58, 0x96, 0xa3, 0x20, 0x2f, 0x07, 0x2c, 0xa9, 0xb4,
	0x8d, 0x91, 0x4b, 0xd8, 0x25, 0x31, 0x35, 0x37, 0x3f, 0x2f, 0x2c, 0xbf, 0x24, 0xd5, 0xb5, 0xa2,
	0x24, 0x35, 0xaf, 0x38, 0x33, 0x3f, 0x4f, 0xc8, 0x83, 0x8b, 0x0d, 0xac, 0xa1, 0x58, 0x82, 0x51,
	0x81, 0x59, 0x83, 0xdb, 0xc8, 0x40, 0x0f, 0xcd, 0x09, 0x7a, 0x58, 0x74, 0xe9, 0x81, 0xcd, 0x2c,
	0x76, 0xcd, 0x2b, 0x29, 0xaa, 0x0c, 0x82, 0xea, 0x97, 0x0a, 0xe5, 0xe2, 0x46, 0x12, 0x16, 0x12,
	0xe0, 0x62, 0xce, 0x4e, 0xad, 0x04, 0xbb, 0x92, 0x37, 0x08, 0xc4, 0x14, 0x32, 0xe0, 0x62, 0x2d,
	0x4b, 0xcc, 0x29, 0x85, 0xb8, 0x8d, 0xdb, 0x48, 0x0a, 0xc3, 0x26, 0xb8, 0x4f, 0x83, 0x20, 0x0a,
	0xad, 0x98, 0x2c, 0x18, 0x9d, 0xe2, 0x4e, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1,
	0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21,
	0xca, 0x25, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57, 0x3f, 0xb8, 0xa4, 0x28,
	0x35, 0x31, 0xd7, 0x2d, 0x33, 0x2f, 0x31, 0x2f, 0x39, 0x55, 0x37, 0x00, 0x16, 0xb6, 0xc5, 0x60,
	0x61, 0xdd, 0xe4, 0x8c, 0xc4, 0xcc, 0x3c, 0x7d, 0x78, 0x88, 0x27, 0x16, 0x14, 0x80, 0x42, 0xbd,
	0xa4, 0xb2, 0x20, 0xb5, 0x38, 0x89, 0x0d, 0x2c, 0x6c, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x2b,
	0x66, 0xce, 0x42, 0x95, 0x01, 0x00, 0x00,
}

func (m *PricePair) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PricePair) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PricePair) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PnlPrice) > 0 {
		i -= len(m.PnlPrice)
		copy(dAtA[i:], m.PnlPrice)
		i = encodeVarintVe(dAtA, i, uint64(len(m.PnlPrice)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.SpotPrice) > 0 {
		i -= len(m.SpotPrice)
		copy(dAtA[i:], m.SpotPrice)
		i = encodeVarintVe(dAtA, i, uint64(len(m.SpotPrice)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DaemonVoteExtension) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DaemonVoteExtension) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DaemonVoteExtension) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Prices) > 0 {
		for k := range m.Prices {
			v := m.Prices[k]
			baseI := i
			if v != nil {
				{
					size, err := v.MarshalToSizedBuffer(dAtA[:i])
					if err != nil {
						return 0, err
					}
					i -= size
					i = encodeVarintVe(dAtA, i, uint64(size))
				}
				i--
				dAtA[i] = 0x12
			}
			i = encodeVarintVe(dAtA, i, uint64(k))
			i--
			dAtA[i] = 0x8
			i = encodeVarintVe(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintVe(dAtA []byte, offset int, v uint64) int {
	offset -= sovVe(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PricePair) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.SpotPrice)
	if l > 0 {
		n += 1 + l + sovVe(uint64(l))
	}
	l = len(m.PnlPrice)
	if l > 0 {
		n += 1 + l + sovVe(uint64(l))
	}
	return n
}

func (m *DaemonVoteExtension) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Prices) > 0 {
		for k, v := range m.Prices {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovVe(uint64(l))
			}
			mapEntrySize := 1 + sovVe(uint64(k)) + l
			n += mapEntrySize + 1 + sovVe(uint64(mapEntrySize))
		}
	}
	return n
}

func sovVe(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozVe(x uint64) (n int) {
	return sovVe(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PricePair) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVe
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
			return fmt.Errorf("proto: PricePair: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PricePair: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SpotPrice", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVe
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
				return ErrInvalidLengthVe
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthVe
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SpotPrice = append(m.SpotPrice[:0], dAtA[iNdEx:postIndex]...)
			if m.SpotPrice == nil {
				m.SpotPrice = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PnlPrice", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVe
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
				return ErrInvalidLengthVe
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthVe
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PnlPrice = append(m.PnlPrice[:0], dAtA[iNdEx:postIndex]...)
			if m.PnlPrice == nil {
				m.PnlPrice = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVe(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthVe
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
func (m *DaemonVoteExtension) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVe
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
			return fmt.Errorf("proto: DaemonVoteExtension: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DaemonVoteExtension: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prices", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVe
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
				return ErrInvalidLengthVe
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthVe
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Prices == nil {
				m.Prices = make(map[uint32]*PricePair)
			}
			var mapkey uint32
			var mapvalue *PricePair
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowVe
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
				if fieldNum == 1 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowVe
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapkey |= uint32(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowVe
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthVe
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthVe
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &PricePair{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipVe(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthVe
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Prices[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVe(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthVe
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
func skipVe(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowVe
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
					return 0, ErrIntOverflowVe
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
					return 0, ErrIntOverflowVe
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
				return 0, ErrInvalidLengthVe
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupVe
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthVe
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthVe        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowVe          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupVe = fmt.Errorf("proto: unexpected end of group")
)
