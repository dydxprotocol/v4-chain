// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/vault/params.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_dydxprotocol_v4_chain_protocol_dtypes "github.com/dydxprotocol/v4-chain/protocol/dtypes"
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

// QuotingParams stores vault quoting parameters.
type QuotingParams struct {
	// The number of layers of orders a vault places. For example if
	// `layers=2`, a vault places 2 asks and 2 bids.
	Layers uint32 `protobuf:"varint,1,opt,name=layers,proto3" json:"layers,omitempty"`
	// The minimum base spread when a vault quotes around reservation price.
	SpreadMinPpm uint32 `protobuf:"varint,2,opt,name=spread_min_ppm,json=spreadMinPpm,proto3" json:"spread_min_ppm,omitempty"`
	// The buffer amount to add to min_price_change_ppm to arrive at `spread`
	// according to formula:
	// `spread = max(spread_min_ppm, min_price_change_ppm + spread_buffer_ppm)`.
	SpreadBufferPpm uint32 `protobuf:"varint,3,opt,name=spread_buffer_ppm,json=spreadBufferPpm,proto3" json:"spread_buffer_ppm,omitempty"`
	// The factor that determines how aggressive a vault skews its orders.
	SkewFactorPpm uint32 `protobuf:"varint,4,opt,name=skew_factor_ppm,json=skewFactorPpm,proto3" json:"skew_factor_ppm,omitempty"`
	// The percentage of vault equity that each order is sized at.
	OrderSizePctPpm uint32 `protobuf:"varint,5,opt,name=order_size_pct_ppm,json=orderSizePctPpm,proto3" json:"order_size_pct_ppm,omitempty"`
	// The duration that a vault's orders are valid for.
	OrderExpirationSeconds uint32 `protobuf:"varint,6,opt,name=order_expiration_seconds,json=orderExpirationSeconds,proto3" json:"order_expiration_seconds,omitempty"`
	// The number of quote quantums in quote asset that a vault with no perpetual
	// positions must have to activate, i.e. if a vault has no perpetual positions
	// and has strictly less than this amount of quote asset, it will not
	// activate.
	ActivationThresholdQuoteQuantums github_com_dydxprotocol_v4_chain_protocol_dtypes.SerializableInt `protobuf:"bytes,7,opt,name=activation_threshold_quote_quantums,json=activationThresholdQuoteQuantums,proto3,customtype=github.com/dydxprotocol/v4-chain/protocol/dtypes.SerializableInt" json:"activation_threshold_quote_quantums"`
}

func (m *QuotingParams) Reset()         { *m = QuotingParams{} }
func (m *QuotingParams) String() string { return proto.CompactTextString(m) }
func (*QuotingParams) ProtoMessage()    {}
func (*QuotingParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_6043e0b8bfdbca9f, []int{0}
}
func (m *QuotingParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QuotingParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QuotingParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QuotingParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QuotingParams.Merge(m, src)
}
func (m *QuotingParams) XXX_Size() int {
	return m.Size()
}
func (m *QuotingParams) XXX_DiscardUnknown() {
	xxx_messageInfo_QuotingParams.DiscardUnknown(m)
}

var xxx_messageInfo_QuotingParams proto.InternalMessageInfo

func (m *QuotingParams) GetLayers() uint32 {
	if m != nil {
		return m.Layers
	}
	return 0
}

func (m *QuotingParams) GetSpreadMinPpm() uint32 {
	if m != nil {
		return m.SpreadMinPpm
	}
	return 0
}

func (m *QuotingParams) GetSpreadBufferPpm() uint32 {
	if m != nil {
		return m.SpreadBufferPpm
	}
	return 0
}

func (m *QuotingParams) GetSkewFactorPpm() uint32 {
	if m != nil {
		return m.SkewFactorPpm
	}
	return 0
}

func (m *QuotingParams) GetOrderSizePctPpm() uint32 {
	if m != nil {
		return m.OrderSizePctPpm
	}
	return 0
}

func (m *QuotingParams) GetOrderExpirationSeconds() uint32 {
	if m != nil {
		return m.OrderExpirationSeconds
	}
	return 0
}

// VaultParams stores individual parameters of a vault.
type VaultParams struct {
	// The quoting parameters specific to this vault.
	QuotingParams QuotingParams `protobuf:"bytes,1,opt,name=quoting_params,json=quotingParams,proto3" json:"quoting_params"`
}

func (m *VaultParams) Reset()         { *m = VaultParams{} }
func (m *VaultParams) String() string { return proto.CompactTextString(m) }
func (*VaultParams) ProtoMessage()    {}
func (*VaultParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_6043e0b8bfdbca9f, []int{1}
}
func (m *VaultParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *VaultParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_VaultParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *VaultParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VaultParams.Merge(m, src)
}
func (m *VaultParams) XXX_Size() int {
	return m.Size()
}
func (m *VaultParams) XXX_DiscardUnknown() {
	xxx_messageInfo_VaultParams.DiscardUnknown(m)
}

var xxx_messageInfo_VaultParams proto.InternalMessageInfo

func (m *VaultParams) GetQuotingParams() QuotingParams {
	if m != nil {
		return m.QuotingParams
	}
	return QuotingParams{}
}

func init() {
	proto.RegisterType((*QuotingParams)(nil), "dydxprotocol.vault.QuotingParams")
	proto.RegisterType((*VaultParams)(nil), "dydxprotocol.vault.VaultParams")
}

func init() { proto.RegisterFile("dydxprotocol/vault/params.proto", fileDescriptor_6043e0b8bfdbca9f) }

var fileDescriptor_6043e0b8bfdbca9f = []byte{
	// 433 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x4f, 0x8b, 0xd4, 0x30,
	0x18, 0xc6, 0xa7, 0xee, 0x3a, 0x42, 0x76, 0x67, 0x17, 0x8b, 0x2c, 0xc5, 0x43, 0x67, 0x5c, 0x45,
	0x16, 0xc5, 0x16, 0x54, 0xd0, 0xa3, 0x0c, 0x28, 0x7a, 0x50, 0xe6, 0x8f, 0x78, 0x10, 0x24, 0x64,
	0xd2, 0x4c, 0x1b, 0x6c, 0x93, 0x34, 0x49, 0xd7, 0x99, 0xf9, 0x14, 0xde, 0xbc, 0xfb, 0x69, 0xf6,
	0xb8, 0x47, 0xf1, 0xb0, 0xc8, 0xcc, 0x17, 0x91, 0xbc, 0xa9, 0xab, 0x83, 0x97, 0xbd, 0xb5, 0xcf,
	0xf3, 0x4b, 0xde, 0xbe, 0xfc, 0x8a, 0xfa, 0xd9, 0x32, 0x5b, 0x28, 0x2d, 0xad, 0xa4, 0xb2, 0x4c,
	0x4f, 0x49, 0x53, 0xda, 0x54, 0x11, 0x4d, 0x2a, 0x93, 0x40, 0x1a, 0x86, 0xff, 0x02, 0x09, 0x00,
	0xb7, 0x6f, 0xe5, 0x32, 0x97, 0x90, 0xa5, 0xee, 0xc9, 0x93, 0xc7, 0xdf, 0x77, 0x50, 0x6f, 0xdc,
	0x48, 0xcb, 0x45, 0x3e, 0x82, 0x1b, 0xc2, 0x23, 0xd4, 0x2d, 0xc9, 0x92, 0x69, 0x13, 0x05, 0x83,
	0xe0, 0xa4, 0x37, 0x69, 0xdf, 0xc2, 0x7b, 0xe8, 0xc0, 0x28, 0xcd, 0x48, 0x86, 0x2b, 0x2e, 0xb0,
	0x52, 0x55, 0x74, 0x0d, 0xfa, 0x7d, 0x9f, 0xbe, 0xe5, 0x62, 0xa4, 0xaa, 0xf0, 0x01, 0xba, 0xd9,
	0x52, 0xb3, 0x66, 0x3e, 0x67, 0x1a, 0xc0, 0x1d, 0x00, 0x0f, 0x7d, 0x31, 0x84, 0xdc, 0xb1, 0xf7,
	0xd1, 0xa1, 0xf9, 0xcc, 0xbe, 0xe0, 0x39, 0xa1, 0x56, 0x7a, 0x72, 0x17, 0xc8, 0x9e, 0x8b, 0x5f,
	0x41, 0xea, 0xb8, 0x87, 0x28, 0x94, 0x3a, 0x63, 0x1a, 0x1b, 0xbe, 0x62, 0x58, 0x51, 0x0b, 0xe8,
	0x75, 0x7f, 0x29, 0x34, 0x53, 0xbe, 0x62, 0x23, 0x6a, 0x1d, 0xfc, 0x1c, 0x45, 0x1e, 0x66, 0x0b,
	0xc5, 0x35, 0xb1, 0x5c, 0x0a, 0x6c, 0x18, 0x95, 0x22, 0x33, 0x51, 0x17, 0x8e, 0x1c, 0x41, 0xff,
	0xf2, 0xb2, 0x9e, 0xfa, 0x36, 0xfc, 0x16, 0xa0, 0xbb, 0x84, 0x5a, 0x7e, 0xea, 0x0f, 0xd9, 0x42,
	0x33, 0x53, 0xc8, 0x32, 0xc3, 0x75, 0x23, 0x2d, 0xc3, 0x75, 0x43, 0x84, 0x6d, 0x2a, 0x13, 0xdd,
	0x18, 0x04, 0x27, 0xfb, 0xc3, 0xd7, 0x67, 0x17, 0xfd, 0xce, 0xcf, 0x8b, 0xfe, 0x8b, 0x9c, 0xdb,
	0xa2, 0x99, 0x25, 0x54, 0x56, 0xe9, 0xb6, 0x96, 0xa7, 0x8f, 0x68, 0x41, 0xb8, 0x48, 0x2f, 0x93,
	0xcc, 0x2e, 0x15, 0x33, 0xc9, 0x94, 0x69, 0x4e, 0x4a, 0xbe, 0x22, 0xb3, 0x92, 0xbd, 0x11, 0x76,
	0x32, 0xf8, 0x3b, 0xf4, 0xfd, 0x9f, 0x99, 0x4e, 0x09, 0x1b, 0xb7, 0x13, 0x8f, 0x3f, 0xa1, 0xbd,
	0x0f, 0xce, 0x61, 0x6b, 0xe8, 0x1d, 0x3a, 0xa8, 0xbd, 0x32, 0xec, 0xad, 0x83, 0xa9, 0xbd, 0xc7,
	0x77, 0x92, 0xff, 0xb5, 0x27, 0x5b, 0x72, 0x87, 0xbb, 0xee, 0xab, 0x27, 0xbd, 0x7a, 0x2b, 0x1c,
	0x9f, 0xad, 0xe3, 0xe0, 0x7c, 0x1d, 0x07, 0xbf, 0xd6, 0x71, 0xf0, 0x75, 0x13, 0x77, 0xce, 0x37,
	0x71, 0xe7, 0xc7, 0x26, 0xee, 0x7c, 0x7c, 0x76, 0xf5, 0xe5, 0x16, 0xed, 0x7f, 0x08, 0x3b, 0xce,
	0xba, 0x90, 0x3f, 0xf9, 0x1d, 0x00, 0x00, 0xff, 0xff, 0xe3, 0xd1, 0x51, 0x65, 0xaa, 0x02, 0x00,
	0x00,
}

func (m *QuotingParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QuotingParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QuotingParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.ActivationThresholdQuoteQuantums.Size()
		i -= size
		if _, err := m.ActivationThresholdQuoteQuantums.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	if m.OrderExpirationSeconds != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.OrderExpirationSeconds))
		i--
		dAtA[i] = 0x30
	}
	if m.OrderSizePctPpm != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.OrderSizePctPpm))
		i--
		dAtA[i] = 0x28
	}
	if m.SkewFactorPpm != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SkewFactorPpm))
		i--
		dAtA[i] = 0x20
	}
	if m.SpreadBufferPpm != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SpreadBufferPpm))
		i--
		dAtA[i] = 0x18
	}
	if m.SpreadMinPpm != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SpreadMinPpm))
		i--
		dAtA[i] = 0x10
	}
	if m.Layers != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.Layers))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *VaultParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *VaultParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *VaultParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.QuotingParams.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QuotingParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Layers != 0 {
		n += 1 + sovParams(uint64(m.Layers))
	}
	if m.SpreadMinPpm != 0 {
		n += 1 + sovParams(uint64(m.SpreadMinPpm))
	}
	if m.SpreadBufferPpm != 0 {
		n += 1 + sovParams(uint64(m.SpreadBufferPpm))
	}
	if m.SkewFactorPpm != 0 {
		n += 1 + sovParams(uint64(m.SkewFactorPpm))
	}
	if m.OrderSizePctPpm != 0 {
		n += 1 + sovParams(uint64(m.OrderSizePctPpm))
	}
	if m.OrderExpirationSeconds != 0 {
		n += 1 + sovParams(uint64(m.OrderExpirationSeconds))
	}
	l = m.ActivationThresholdQuoteQuantums.Size()
	n += 1 + l + sovParams(uint64(l))
	return n
}

func (m *VaultParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.QuotingParams.Size()
	n += 1 + l + sovParams(uint64(l))
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QuotingParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: QuotingParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QuotingParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Layers", wireType)
			}
			m.Layers = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Layers |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SpreadMinPpm", wireType)
			}
			m.SpreadMinPpm = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SpreadMinPpm |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SpreadBufferPpm", wireType)
			}
			m.SpreadBufferPpm = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SpreadBufferPpm |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SkewFactorPpm", wireType)
			}
			m.SkewFactorPpm = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SkewFactorPpm |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderSizePctPpm", wireType)
			}
			m.OrderSizePctPpm = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrderSizePctPpm |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderExpirationSeconds", wireType)
			}
			m.OrderExpirationSeconds = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrderExpirationSeconds |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ActivationThresholdQuoteQuantums", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ActivationThresholdQuoteQuantums.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func (m *VaultParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: VaultParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: VaultParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field QuotingParams", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.QuotingParams.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
