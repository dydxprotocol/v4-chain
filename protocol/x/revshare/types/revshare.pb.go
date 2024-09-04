// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/revshare/revshare.proto

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

// MarketMapperRevShareDetails specifies any details associated with the market
// mapper revenue share
type MarketMapperRevShareDetails struct {
	// Unix timestamp recorded when the market revenue share expires
	ExpirationTs uint64 `protobuf:"varint,1,opt,name=expiration_ts,json=expirationTs,proto3" json:"expiration_ts,omitempty"`
}

func (m *MarketMapperRevShareDetails) Reset()         { *m = MarketMapperRevShareDetails{} }
func (m *MarketMapperRevShareDetails) String() string { return proto.CompactTextString(m) }
func (*MarketMapperRevShareDetails) ProtoMessage()    {}
func (*MarketMapperRevShareDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_5b9759663d195798, []int{0}
}
func (m *MarketMapperRevShareDetails) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MarketMapperRevShareDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MarketMapperRevShareDetails.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MarketMapperRevShareDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MarketMapperRevShareDetails.Merge(m, src)
}
func (m *MarketMapperRevShareDetails) XXX_Size() int {
	return m.Size()
}
func (m *MarketMapperRevShareDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_MarketMapperRevShareDetails.DiscardUnknown(m)
}

var xxx_messageInfo_MarketMapperRevShareDetails proto.InternalMessageInfo

func (m *MarketMapperRevShareDetails) GetExpirationTs() uint64 {
	if m != nil {
		return m.ExpirationTs
	}
	return 0
}

// UnconditionalRevShareConfig stores recipients that
// receive a share of net revenue unconditionally.
type UnconditionalRevShareConfig struct {
	// Configs for each recipient.
	Configs []UnconditionalRevShareConfig_RecipientConfig `protobuf:"bytes,1,rep,name=configs,proto3" json:"configs"`
}

func (m *UnconditionalRevShareConfig) Reset()         { *m = UnconditionalRevShareConfig{} }
func (m *UnconditionalRevShareConfig) String() string { return proto.CompactTextString(m) }
func (*UnconditionalRevShareConfig) ProtoMessage()    {}
func (*UnconditionalRevShareConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_5b9759663d195798, []int{1}
}
func (m *UnconditionalRevShareConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UnconditionalRevShareConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UnconditionalRevShareConfig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UnconditionalRevShareConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnconditionalRevShareConfig.Merge(m, src)
}
func (m *UnconditionalRevShareConfig) XXX_Size() int {
	return m.Size()
}
func (m *UnconditionalRevShareConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_UnconditionalRevShareConfig.DiscardUnknown(m)
}

var xxx_messageInfo_UnconditionalRevShareConfig proto.InternalMessageInfo

func (m *UnconditionalRevShareConfig) GetConfigs() []UnconditionalRevShareConfig_RecipientConfig {
	if m != nil {
		return m.Configs
	}
	return nil
}

// Describes the config of a recipient
type UnconditionalRevShareConfig_RecipientConfig struct {
	// Address of the recepient.
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// Percentage of net revenue to share with recipient, in parts-per-million.
	SharePpm uint32 `protobuf:"varint,2,opt,name=share_ppm,json=sharePpm,proto3" json:"share_ppm,omitempty"`
}

func (m *UnconditionalRevShareConfig_RecipientConfig) Reset() {
	*m = UnconditionalRevShareConfig_RecipientConfig{}
}
func (m *UnconditionalRevShareConfig_RecipientConfig) String() string {
	return proto.CompactTextString(m)
}
func (*UnconditionalRevShareConfig_RecipientConfig) ProtoMessage() {}
func (*UnconditionalRevShareConfig_RecipientConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_5b9759663d195798, []int{1, 0}
}
func (m *UnconditionalRevShareConfig_RecipientConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UnconditionalRevShareConfig_RecipientConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UnconditionalRevShareConfig_RecipientConfig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UnconditionalRevShareConfig_RecipientConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UnconditionalRevShareConfig_RecipientConfig.Merge(m, src)
}
func (m *UnconditionalRevShareConfig_RecipientConfig) XXX_Size() int {
	return m.Size()
}
func (m *UnconditionalRevShareConfig_RecipientConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_UnconditionalRevShareConfig_RecipientConfig.DiscardUnknown(m)
}

var xxx_messageInfo_UnconditionalRevShareConfig_RecipientConfig proto.InternalMessageInfo

func (m *UnconditionalRevShareConfig_RecipientConfig) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *UnconditionalRevShareConfig_RecipientConfig) GetSharePpm() uint32 {
	if m != nil {
		return m.SharePpm
	}
	return 0
}

func init() {
	proto.RegisterType((*MarketMapperRevShareDetails)(nil), "dydxprotocol.revshare.MarketMapperRevShareDetails")
	proto.RegisterType((*UnconditionalRevShareConfig)(nil), "dydxprotocol.revshare.UnconditionalRevShareConfig")
	proto.RegisterType((*UnconditionalRevShareConfig_RecipientConfig)(nil), "dydxprotocol.revshare.UnconditionalRevShareConfig.RecipientConfig")
}

func init() {
	proto.RegisterFile("dydxprotocol/revshare/revshare.proto", fileDescriptor_5b9759663d195798)
}

var fileDescriptor_5b9759663d195798 = []byte{
	// 307 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0x31, 0x4b, 0x03, 0x31,
	0x1c, 0xc5, 0x2f, 0x5a, 0xac, 0x8d, 0x16, 0xe1, 0x50, 0x28, 0x2d, 0xc4, 0x52, 0x1d, 0xba, 0x98,
	0x03, 0x75, 0x72, 0x3c, 0x1d, 0x5c, 0x0a, 0x12, 0xeb, 0xe2, 0x52, 0xd2, 0x5c, 0xbc, 0x06, 0xdb,
	0x24, 0x24, 0xb1, 0xb4, 0xdf, 0xc2, 0x8f, 0x55, 0x70, 0xe9, 0xe8, 0x24, 0xd2, 0xfb, 0x22, 0x72,
	0x39, 0xae, 0x55, 0x11, 0xb7, 0xf7, 0xff, 0xbd, 0x97, 0x97, 0x3f, 0xfc, 0xe1, 0x69, 0x32, 0x4f,
	0x66, 0xda, 0x28, 0xa7, 0x98, 0x1a, 0x47, 0x86, 0x4f, 0xed, 0x88, 0x1a, 0xbe, 0x16, 0xd8, 0x5b,
	0xe1, 0xd1, 0xf7, 0x14, 0x2e, 0xcd, 0xe6, 0x61, 0xaa, 0x52, 0xe5, 0x71, 0x94, 0xab, 0x22, 0xdc,
	0x89, 0x61, 0xab, 0x47, 0xcd, 0x33, 0x77, 0x3d, 0xaa, 0x35, 0x37, 0x84, 0x4f, 0xef, 0xf3, 0xf4,
	0x0d, 0x77, 0x54, 0x8c, 0x6d, 0x78, 0x02, 0xeb, 0x7c, 0xa6, 0x85, 0xa1, 0x4e, 0x28, 0x39, 0x70,
	0xb6, 0x01, 0xda, 0xa0, 0x5b, 0x21, 0xfb, 0x1b, 0xd8, 0xb7, 0x9d, 0x37, 0x00, 0x5b, 0x0f, 0x92,
	0x29, 0x99, 0x88, 0x9c, 0xd0, 0x71, 0xd9, 0x72, 0xad, 0xe4, 0x93, 0x48, 0xc3, 0x21, 0xac, 0x32,
	0xaf, 0xf2, 0xe7, 0xdb, 0xdd, 0xbd, 0xf3, 0x18, 0xff, 0xb9, 0x22, 0xfe, 0xa7, 0x04, 0x13, 0xce,
	0x84, 0x16, 0x5c, 0xba, 0x62, 0x8e, 0x2b, 0x8b, 0x8f, 0xe3, 0x80, 0x94, 0xc5, 0xcd, 0x5b, 0x78,
	0xf0, 0x2b, 0x11, 0x36, 0x60, 0x95, 0x26, 0x89, 0xe1, 0xb6, 0xd8, 0xba, 0x46, 0xca, 0x31, 0x6c,
	0xc1, 0x9a, 0xff, 0x70, 0xa0, 0xf5, 0xa4, 0xb1, 0xd5, 0x06, 0xdd, 0x3a, 0xd9, 0xf5, 0xe0, 0x4e,
	0x4f, 0xe2, 0xfe, 0x62, 0x85, 0xc0, 0x72, 0x85, 0xc0, 0xe7, 0x0a, 0x81, 0xd7, 0x0c, 0x05, 0xcb,
	0x0c, 0x05, 0xef, 0x19, 0x0a, 0x1e, 0xaf, 0x52, 0xe1, 0x46, 0x2f, 0x43, 0xcc, 0xd4, 0x24, 0xfa,
	0x71, 0x89, 0xe9, 0xe5, 0x19, 0x1b, 0x51, 0x21, 0xa3, 0x35, 0x99, 0x6d, 0xae, 0xe3, 0xe6, 0x9a,
	0xdb, 0xe1, 0x8e, 0xb7, 0x2e, 0xbe, 0x02, 0x00, 0x00, 0xff, 0xff, 0x92, 0xdf, 0x87, 0xec, 0xc3,
	0x01, 0x00, 0x00,
}

func (m *MarketMapperRevShareDetails) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MarketMapperRevShareDetails) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MarketMapperRevShareDetails) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ExpirationTs != 0 {
		i = encodeVarintRevshare(dAtA, i, uint64(m.ExpirationTs))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *UnconditionalRevShareConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UnconditionalRevShareConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UnconditionalRevShareConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Configs) > 0 {
		for iNdEx := len(m.Configs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Configs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintRevshare(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *UnconditionalRevShareConfig_RecipientConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UnconditionalRevShareConfig_RecipientConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UnconditionalRevShareConfig_RecipientConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SharePpm != 0 {
		i = encodeVarintRevshare(dAtA, i, uint64(m.SharePpm))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintRevshare(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintRevshare(dAtA []byte, offset int, v uint64) int {
	offset -= sovRevshare(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MarketMapperRevShareDetails) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ExpirationTs != 0 {
		n += 1 + sovRevshare(uint64(m.ExpirationTs))
	}
	return n
}

func (m *UnconditionalRevShareConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Configs) > 0 {
		for _, e := range m.Configs {
			l = e.Size()
			n += 1 + l + sovRevshare(uint64(l))
		}
	}
	return n
}

func (m *UnconditionalRevShareConfig_RecipientConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovRevshare(uint64(l))
	}
	if m.SharePpm != 0 {
		n += 1 + sovRevshare(uint64(m.SharePpm))
	}
	return n
}

func sovRevshare(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRevshare(x uint64) (n int) {
	return sovRevshare(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MarketMapperRevShareDetails) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRevshare
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
			return fmt.Errorf("proto: MarketMapperRevShareDetails: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MarketMapperRevShareDetails: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExpirationTs", wireType)
			}
			m.ExpirationTs = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRevshare
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExpirationTs |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRevshare(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRevshare
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
func (m *UnconditionalRevShareConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRevshare
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
			return fmt.Errorf("proto: UnconditionalRevShareConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UnconditionalRevShareConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Configs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRevshare
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
				return ErrInvalidLengthRevshare
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRevshare
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Configs = append(m.Configs, UnconditionalRevShareConfig_RecipientConfig{})
			if err := m.Configs[len(m.Configs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRevshare(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRevshare
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
func (m *UnconditionalRevShareConfig_RecipientConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRevshare
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
			return fmt.Errorf("proto: RecipientConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RecipientConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRevshare
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
				return ErrInvalidLengthRevshare
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRevshare
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SharePpm", wireType)
			}
			m.SharePpm = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRevshare
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SharePpm |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipRevshare(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRevshare
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
func skipRevshare(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRevshare
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
					return 0, ErrIntOverflowRevshare
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
					return 0, ErrIntOverflowRevshare
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
				return 0, ErrInvalidLengthRevshare
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRevshare
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRevshare
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRevshare        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRevshare          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRevshare = fmt.Errorf("proto: unexpected end of group")
)
