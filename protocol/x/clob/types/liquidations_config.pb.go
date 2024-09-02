// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/clob/liquidations_config.proto

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

// LiquidationsConfig stores all configurable fields related to liquidations.
type LiquidationsConfig struct {
	// The maximum liquidation fee (in parts-per-million). This fee goes
	// 100% to the insurance fund.
	MaxLiquidationFeePpm uint32 `protobuf:"varint,1,opt,name=max_liquidation_fee_ppm,json=maxLiquidationFeePpm,proto3" json:"max_liquidation_fee_ppm,omitempty"`
	// Limits around how many quote quantums from a single subaccount can
	// be liquidated within a single block.
	SubaccountBlockLimits SubaccountBlockLimits `protobuf:"bytes,2,opt,name=subaccount_block_limits,json=subaccountBlockLimits,proto3" json:"subaccount_block_limits"`
}

func (m *LiquidationsConfig) Reset()         { *m = LiquidationsConfig{} }
func (m *LiquidationsConfig) String() string { return proto.CompactTextString(m) }
func (*LiquidationsConfig) ProtoMessage()    {}
func (*LiquidationsConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_d11e0d49099a14b4, []int{0}
}
func (m *LiquidationsConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LiquidationsConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LiquidationsConfig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LiquidationsConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LiquidationsConfig.Merge(m, src)
}
func (m *LiquidationsConfig) XXX_Size() int {
	return m.Size()
}
func (m *LiquidationsConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_LiquidationsConfig.DiscardUnknown(m)
}

var xxx_messageInfo_LiquidationsConfig proto.InternalMessageInfo

func (m *LiquidationsConfig) GetMaxLiquidationFeePpm() uint32 {
	if m != nil {
		return m.MaxLiquidationFeePpm
	}
	return 0
}

func (m *LiquidationsConfig) GetSubaccountBlockLimits() SubaccountBlockLimits {
	if m != nil {
		return m.SubaccountBlockLimits
	}
	return SubaccountBlockLimits{}
}

// SubaccountBlockLimits stores all configurable fields related to limits
// around how many quote quantums from a single subaccount can
// be liquidated within a single block.
type SubaccountBlockLimits struct {
	// The maximum insurance-fund payout amount for a given subaccount
	// per block. I.e. how much it can cover for that subaccount.
	MaxQuantumsInsuranceLost uint64 `protobuf:"varint,1,opt,name=max_quantums_insurance_lost,json=maxQuantumsInsuranceLost,proto3" json:"max_quantums_insurance_lost,omitempty"`
}

func (m *SubaccountBlockLimits) Reset()         { *m = SubaccountBlockLimits{} }
func (m *SubaccountBlockLimits) String() string { return proto.CompactTextString(m) }
func (*SubaccountBlockLimits) ProtoMessage()    {}
func (*SubaccountBlockLimits) Descriptor() ([]byte, []int) {
	return fileDescriptor_d11e0d49099a14b4, []int{1}
}
func (m *SubaccountBlockLimits) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SubaccountBlockLimits) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SubaccountBlockLimits.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SubaccountBlockLimits) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubaccountBlockLimits.Merge(m, src)
}
func (m *SubaccountBlockLimits) XXX_Size() int {
	return m.Size()
}
func (m *SubaccountBlockLimits) XXX_DiscardUnknown() {
	xxx_messageInfo_SubaccountBlockLimits.DiscardUnknown(m)
}

var xxx_messageInfo_SubaccountBlockLimits proto.InternalMessageInfo

func (m *SubaccountBlockLimits) GetMaxQuantumsInsuranceLost() uint64 {
	if m != nil {
		return m.MaxQuantumsInsuranceLost
	}
	return 0
}

func init() {
	proto.RegisterType((*LiquidationsConfig)(nil), "dydxprotocol.clob.LiquidationsConfig")
	proto.RegisterType((*SubaccountBlockLimits)(nil), "dydxprotocol.clob.SubaccountBlockLimits")
}

func init() {
	proto.RegisterFile("dydxprotocol/clob/liquidations_config.proto", fileDescriptor_d11e0d49099a14b4)
}

var fileDescriptor_d11e0d49099a14b4 = []byte{
	// 323 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0x41, 0x6b, 0xf2, 0x30,
	0x1c, 0xc6, 0x9b, 0x17, 0x79, 0x0f, 0x1d, 0x3b, 0xac, 0x28, 0xca, 0x06, 0x9d, 0x78, 0x12, 0x86,
	0x2d, 0x6c, 0xec, 0xb8, 0x8b, 0x1b, 0xc2, 0xc0, 0x83, 0x53, 0xd8, 0x61, 0x87, 0x85, 0x34, 0xc6,
	0x1a, 0x96, 0xe4, 0x5f, 0x4d, 0x02, 0xf5, 0x5b, 0xec, 0xbb, 0xec, 0x4b, 0x78, 0xf4, 0xb8, 0xd3,
	0x18, 0xf6, 0x8b, 0x8c, 0xc6, 0x29, 0x85, 0x79, 0x2b, 0x4f, 0x7f, 0xbf, 0x7f, 0x1e, 0x78, 0xfc,
	0xab, 0xe9, 0x6a, 0x9a, 0x67, 0x4b, 0x30, 0x40, 0x41, 0xc4, 0x54, 0x40, 0x12, 0x0b, 0xbe, 0xb0,
	0x7c, 0x4a, 0x0c, 0x07, 0xa5, 0x31, 0x05, 0x35, 0xe3, 0x69, 0xe4, 0x88, 0xe0, 0xac, 0x0a, 0x47,
	0x25, 0x7c, 0x5e, 0x4f, 0x21, 0x05, 0x17, 0xc5, 0xe5, 0xd7, 0x0e, 0xec, 0x7c, 0x20, 0x3f, 0x18,
	0x56, 0xce, 0xdc, 0xbb, 0x2b, 0xc1, 0xad, 0xdf, 0x94, 0x24, 0xc7, 0x95, 0x07, 0xf0, 0x8c, 0x31,
	0x9c, 0x65, 0xb2, 0x85, 0xda, 0xa8, 0x7b, 0x3a, 0xae, 0x4b, 0x92, 0x57, 0xbc, 0x01, 0x63, 0xa3,
	0x4c, 0x06, 0x33, 0xbf, 0xa9, 0x6d, 0x42, 0x28, 0x05, 0xab, 0x0c, 0x4e, 0x04, 0xd0, 0x37, 0x2c,
	0xb8, 0xe4, 0x46, 0xb7, 0xfe, 0xb5, 0x51, 0xf7, 0xe4, 0xba, 0x1b, 0xfd, 0x29, 0x16, 0x4d, 0x0e,
	0x46, 0xbf, 0x14, 0x86, 0x8e, 0xef, 0xd7, 0xd6, 0x5f, 0x97, 0xde, 0xb8, 0xa1, 0x8f, 0xfd, 0xec,
	0x3c, 0xfb, 0x8d, 0xa3, 0x56, 0x70, 0xe7, 0x5f, 0x94, 0xbd, 0x17, 0x96, 0x28, 0x63, 0xa5, 0xc6,
	0x5c, 0x69, 0xbb, 0x24, 0x8a, 0x32, 0x2c, 0x40, 0x1b, 0xd7, 0xbd, 0x36, 0x6e, 0x49, 0x92, 0x3f,
	0xfd, 0x12, 0x8f, 0x7b, 0x60, 0x08, 0xda, 0xf4, 0x5f, 0xd7, 0xdb, 0x10, 0x6d, 0xb6, 0x21, 0xfa,
	0xde, 0x86, 0xe8, 0xbd, 0x08, 0xbd, 0x4d, 0x11, 0x7a, 0x9f, 0x45, 0xe8, 0xbd, 0x3c, 0xa4, 0xdc,
	0xcc, 0x6d, 0x12, 0x51, 0x90, 0xf1, 0xc4, 0x2c, 0x19, 0x91, 0x03, 0xae, 0x4a, 0xaf, 0x37, 0xda,
	0x4f, 0xa2, 0x5d, 0xdc, 0xa3, 0x73, 0xc2, 0x55, 0x7c, 0x18, 0x2a, 0xdf, 0x4d, 0x65, 0x56, 0x19,
	0xd3, 0xc9, 0x7f, 0x17, 0xdf, 0xfc, 0x04, 0x00, 0x00, 0xff, 0xff, 0xbc, 0xe5, 0x2c, 0xd0, 0xcc,
	0x01, 0x00, 0x00,
}

func (m *LiquidationsConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LiquidationsConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LiquidationsConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.SubaccountBlockLimits.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintLiquidationsConfig(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.MaxLiquidationFeePpm != 0 {
		i = encodeVarintLiquidationsConfig(dAtA, i, uint64(m.MaxLiquidationFeePpm))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SubaccountBlockLimits) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SubaccountBlockLimits) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SubaccountBlockLimits) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MaxQuantumsInsuranceLost != 0 {
		i = encodeVarintLiquidationsConfig(dAtA, i, uint64(m.MaxQuantumsInsuranceLost))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintLiquidationsConfig(dAtA []byte, offset int, v uint64) int {
	offset -= sovLiquidationsConfig(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *LiquidationsConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MaxLiquidationFeePpm != 0 {
		n += 1 + sovLiquidationsConfig(uint64(m.MaxLiquidationFeePpm))
	}
	l = m.SubaccountBlockLimits.Size()
	n += 1 + l + sovLiquidationsConfig(uint64(l))
	return n
}

func (m *SubaccountBlockLimits) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MaxQuantumsInsuranceLost != 0 {
		n += 1 + sovLiquidationsConfig(uint64(m.MaxQuantumsInsuranceLost))
	}
	return n
}

func sovLiquidationsConfig(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozLiquidationsConfig(x uint64) (n int) {
	return sovLiquidationsConfig(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *LiquidationsConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLiquidationsConfig
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
			return fmt.Errorf("proto: LiquidationsConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LiquidationsConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxLiquidationFeePpm", wireType)
			}
			m.MaxLiquidationFeePpm = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidationsConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxLiquidationFeePpm |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubaccountBlockLimits", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidationsConfig
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
				return ErrInvalidLengthLiquidationsConfig
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthLiquidationsConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SubaccountBlockLimits.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLiquidationsConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLiquidationsConfig
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
func (m *SubaccountBlockLimits) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLiquidationsConfig
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
			return fmt.Errorf("proto: SubaccountBlockLimits: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SubaccountBlockLimits: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxQuantumsInsuranceLost", wireType)
			}
			m.MaxQuantumsInsuranceLost = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLiquidationsConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxQuantumsInsuranceLost |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipLiquidationsConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLiquidationsConfig
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
func skipLiquidationsConfig(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLiquidationsConfig
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
					return 0, ErrIntOverflowLiquidationsConfig
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
					return 0, ErrIntOverflowLiquidationsConfig
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
				return 0, ErrInvalidLengthLiquidationsConfig
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupLiquidationsConfig
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthLiquidationsConfig
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthLiquidationsConfig        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLiquidationsConfig          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupLiquidationsConfig = fmt.Errorf("proto: unexpected end of group")
)
