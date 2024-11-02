// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: klyraprotocol/subaccounts/subaccount.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

// SubaccountId defines a unique identifier for a Subaccount.
type SubaccountId struct {
	// The address of the wallet that owns this subaccount.
	Owner string `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	// The unique number of this subaccount for the owner.
	// Currently limited to 128*1000 subaccounts per owner.
	Number uint32 `protobuf:"varint,2,opt,name=number,proto3" json:"number,omitempty"`
}

func (m *SubaccountId) Reset()         { *m = SubaccountId{} }
func (m *SubaccountId) String() string { return proto.CompactTextString(m) }
func (*SubaccountId) ProtoMessage()    {}
func (*SubaccountId) Descriptor() ([]byte, []int) {
	return fileDescriptor_66d58a5d9356fb3d, []int{0}
}
func (m *SubaccountId) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SubaccountId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SubaccountId.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SubaccountId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubaccountId.Merge(m, src)
}
func (m *SubaccountId) XXX_Size() int {
	return m.Size()
}
func (m *SubaccountId) XXX_DiscardUnknown() {
	xxx_messageInfo_SubaccountId.DiscardUnknown(m)
}

var xxx_messageInfo_SubaccountId proto.InternalMessageInfo

func (m *SubaccountId) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *SubaccountId) GetNumber() uint32 {
	if m != nil {
		return m.Number
	}
	return 0
}

// Subaccount defines a single sub-account for a given address.
// Subaccounts are uniquely indexed by a subaccountNumber/owner pair.
type Subaccount struct {
	// The Id of the Subaccount
	Id *SubaccountId `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// All `AssetPosition`s associated with this subaccount.
	// Always sorted ascending by `asset_id`.
	AssetPositions []*AssetPosition `protobuf:"bytes,2,rep,name=asset_positions,json=assetPositions,proto3" json:"asset_positions,omitempty"`
	// All `PerpetualPosition`s associated with this subaccount.
	// Always sorted ascending by `perpetual_id.
	PerpetualPositions []*PerpetualPosition `protobuf:"bytes,3,rep,name=perpetual_positions,json=perpetualPositions,proto3" json:"perpetual_positions,omitempty"`
	// Set by the owner. If true, then margin trades can be made in this
	// subaccount.
	MarginEnabled bool `protobuf:"varint,4,opt,name=margin_enabled,json=marginEnabled,proto3" json:"margin_enabled,omitempty"`
	// The current yield index is determined by the cumulative
	// all-time history of the yield mechanism for assets.
	// Starts at 0. This string should always be converted big.Rat.
	AssetYieldIndex string `protobuf:"bytes,5,opt,name=asset_yield_index,json=assetYieldIndex,proto3" json:"asset_yield_index,omitempty"`
}

func (m *Subaccount) Reset()         { *m = Subaccount{} }
func (m *Subaccount) String() string { return proto.CompactTextString(m) }
func (*Subaccount) ProtoMessage()    {}
func (*Subaccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_66d58a5d9356fb3d, []int{1}
}
func (m *Subaccount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Subaccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Subaccount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Subaccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Subaccount.Merge(m, src)
}
func (m *Subaccount) XXX_Size() int {
	return m.Size()
}
func (m *Subaccount) XXX_DiscardUnknown() {
	xxx_messageInfo_Subaccount.DiscardUnknown(m)
}

var xxx_messageInfo_Subaccount proto.InternalMessageInfo

func (m *Subaccount) GetId() *SubaccountId {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Subaccount) GetAssetPositions() []*AssetPosition {
	if m != nil {
		return m.AssetPositions
	}
	return nil
}

func (m *Subaccount) GetPerpetualPositions() []*PerpetualPosition {
	if m != nil {
		return m.PerpetualPositions
	}
	return nil
}

func (m *Subaccount) GetMarginEnabled() bool {
	if m != nil {
		return m.MarginEnabled
	}
	return false
}

func (m *Subaccount) GetAssetYieldIndex() string {
	if m != nil {
		return m.AssetYieldIndex
	}
	return ""
}

func init() {
	proto.RegisterType((*SubaccountId)(nil), "klyraprotocol.subaccounts.SubaccountId")
	proto.RegisterType((*Subaccount)(nil), "klyraprotocol.subaccounts.Subaccount")
}

func init() {
	proto.RegisterFile("klyraprotocol/subaccounts/subaccount.proto", fileDescriptor_66d58a5d9356fb3d)
}

var fileDescriptor_66d58a5d9356fb3d = []byte{
	// 399 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x52, 0x5d, 0x8b, 0xda, 0x40,
	0x14, 0x35, 0xb1, 0x4a, 0x3b, 0x56, 0x4b, 0xa7, 0xa5, 0x44, 0x1f, 0x42, 0x10, 0x4a, 0x83, 0xd4,
	0x04, 0xec, 0x43, 0x9f, 0x15, 0x5a, 0x90, 0xbe, 0xd8, 0x08, 0x85, 0x16, 0x4a, 0x98, 0x64, 0x06,
	0x1d, 0x9a, 0xcc, 0x84, 0x99, 0x09, 0xd5, 0x7f, 0xd1, 0x1f, 0xd3, 0x1f, 0xb1, 0xec, 0x93, 0xec,
	0xd3, 0x3e, 0x2e, 0xfa, 0x47, 0x96, 0x7c, 0x6c, 0x36, 0xba, 0xe8, 0xdb, 0xdc, 0x73, 0xcf, 0x39,
	0x73, 0xef, 0xe5, 0x80, 0xd1, 0x9f, 0x68, 0x2b, 0x50, 0x22, 0xb8, 0xe2, 0x21, 0x8f, 0x5c, 0x99,
	0x06, 0x28, 0x0c, 0x79, 0xca, 0x94, 0xac, 0xbd, 0x9d, 0xbc, 0x0f, 0xfb, 0x47, 0x5c, 0xa7, 0xc6,
	0x1d, 0xf4, 0x43, 0x2e, 0x63, 0x2e, 0xfd, 0xbc, 0xe9, 0x16, 0x45, 0xa1, 0x1a, 0x38, 0xe7, 0x7f,
	0x40, 0x52, 0x12, 0xe5, 0x27, 0x5c, 0x52, 0x45, 0x39, 0x2b, 0xf9, 0x93, 0xf3, 0xfc, 0x84, 0x88,
	0x84, 0xa8, 0x14, 0x45, 0x27, 0x9a, 0xe1, 0x0f, 0xf0, 0x72, 0x59, 0xf1, 0xe6, 0x18, 0x3a, 0xa0,
	0xc5, 0xff, 0x32, 0x22, 0x0c, 0xcd, 0xd2, 0xec, 0x17, 0x33, 0xe3, 0xe6, 0xff, 0xf8, 0x6d, 0x39,
	0xd4, 0x14, 0x63, 0x41, 0xa4, 0x5c, 0x2a, 0x41, 0xd9, 0xca, 0x2b, 0x68, 0xf0, 0x1d, 0x68, 0xb3,
	0x34, 0x0e, 0x88, 0x30, 0x74, 0x4b, 0xb3, 0xbb, 0x5e, 0x59, 0x0d, 0xaf, 0x75, 0x00, 0x1e, 0x8d,
	0xe1, 0x67, 0xa0, 0x53, 0x9c, 0x7b, 0x76, 0x26, 0x1f, 0x9c, 0xb3, 0xd7, 0x70, 0xea, 0xb3, 0x78,
	0x3a, 0xc5, 0xf0, 0x3b, 0x78, 0x75, 0xbc, 0xab, 0x34, 0x74, 0xab, 0x69, 0x77, 0x26, 0xf6, 0x05,
	0x97, 0x69, 0xa6, 0x58, 0x94, 0x02, 0xaf, 0x87, 0xea, 0xa5, 0x84, 0xbf, 0xc1, 0x9b, 0xa7, 0xe7,
	0x90, 0x46, 0x33, 0xb7, 0xfd, 0x78, 0xc1, 0x76, 0xf1, 0xa0, 0xaa, 0xac, 0x61, 0x72, 0x0a, 0x49,
	0xf8, 0x1e, 0xf4, 0x62, 0x24, 0x56, 0x94, 0xf9, 0x84, 0xa1, 0x20, 0x22, 0xd8, 0x78, 0x66, 0x69,
	0xf6, 0x73, 0xaf, 0x5b, 0xa0, 0x5f, 0x0a, 0x10, 0x8e, 0xc0, 0xeb, 0x62, 0xb1, 0x2d, 0x25, 0x11,
	0xf6, 0x29, 0xc3, 0x64, 0x63, 0xb4, 0xb2, 0xa3, 0x7b, 0xc5, 0xc6, 0x3f, 0x33, 0x7c, 0x9e, 0xc1,
	0x33, 0x72, 0xb5, 0x37, 0xb5, 0xdd, 0xde, 0xd4, 0xee, 0xf6, 0xa6, 0xf6, 0xef, 0x60, 0x36, 0x76,
	0x07, 0xb3, 0x71, 0x7b, 0x30, 0x1b, 0xbf, 0xbe, 0xad, 0xa8, 0x5a, 0xa7, 0x81, 0x13, 0xf2, 0xd8,
	0x5d, 0x2a, 0x41, 0x50, 0xfc, 0x95, 0x32, 0xc4, 0x42, 0x32, 0x5e, 0x54, 0x31, 0xc8, 0xe1, 0x71,
	0xb8, 0x46, 0x94, 0xb9, 0x55, 0x38, 0x36, 0x47, 0xf1, 0x50, 0xdb, 0x84, 0xc8, 0xa0, 0x9d, 0x77,
	0x3f, 0xdd, 0x07, 0x00, 0x00, 0xff, 0xff, 0xf0, 0x66, 0x6f, 0x3c, 0xda, 0x02, 0x00, 0x00,
}

func (m *SubaccountId) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SubaccountId) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SubaccountId) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Number != 0 {
		i = encodeVarintSubaccount(dAtA, i, uint64(m.Number))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintSubaccount(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Subaccount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Subaccount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Subaccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.AssetYieldIndex) > 0 {
		i -= len(m.AssetYieldIndex)
		copy(dAtA[i:], m.AssetYieldIndex)
		i = encodeVarintSubaccount(dAtA, i, uint64(len(m.AssetYieldIndex)))
		i--
		dAtA[i] = 0x2a
	}
	if m.MarginEnabled {
		i--
		if m.MarginEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if len(m.PerpetualPositions) > 0 {
		for iNdEx := len(m.PerpetualPositions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PerpetualPositions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSubaccount(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.AssetPositions) > 0 {
		for iNdEx := len(m.AssetPositions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AssetPositions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSubaccount(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Id != nil {
		{
			size, err := m.Id.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSubaccount(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSubaccount(dAtA []byte, offset int, v uint64) int {
	offset -= sovSubaccount(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SubaccountId) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovSubaccount(uint64(l))
	}
	if m.Number != 0 {
		n += 1 + sovSubaccount(uint64(m.Number))
	}
	return n
}

func (m *Subaccount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != nil {
		l = m.Id.Size()
		n += 1 + l + sovSubaccount(uint64(l))
	}
	if len(m.AssetPositions) > 0 {
		for _, e := range m.AssetPositions {
			l = e.Size()
			n += 1 + l + sovSubaccount(uint64(l))
		}
	}
	if len(m.PerpetualPositions) > 0 {
		for _, e := range m.PerpetualPositions {
			l = e.Size()
			n += 1 + l + sovSubaccount(uint64(l))
		}
	}
	if m.MarginEnabled {
		n += 2
	}
	l = len(m.AssetYieldIndex)
	if l > 0 {
		n += 1 + l + sovSubaccount(uint64(l))
	}
	return n
}

func sovSubaccount(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSubaccount(x uint64) (n int) {
	return sovSubaccount(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SubaccountId) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSubaccount
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
			return fmt.Errorf("proto: SubaccountId: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SubaccountId: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
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
				return ErrInvalidLengthSubaccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSubaccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Number", wireType)
			}
			m.Number = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Number |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSubaccount(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSubaccount
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
func (m *Subaccount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSubaccount
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
			return fmt.Errorf("proto: Subaccount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Subaccount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
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
				return ErrInvalidLengthSubaccount
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSubaccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Id == nil {
				m.Id = &SubaccountId{}
			}
			if err := m.Id.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetPositions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
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
				return ErrInvalidLengthSubaccount
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSubaccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AssetPositions = append(m.AssetPositions, &AssetPosition{})
			if err := m.AssetPositions[len(m.AssetPositions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PerpetualPositions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
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
				return ErrInvalidLengthSubaccount
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSubaccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PerpetualPositions = append(m.PerpetualPositions, &PerpetualPosition{})
			if err := m.PerpetualPositions[len(m.PerpetualPositions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MarginEnabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.MarginEnabled = bool(v != 0)
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetYieldIndex", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSubaccount
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
				return ErrInvalidLengthSubaccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSubaccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AssetYieldIndex = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSubaccount(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSubaccount
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
func skipSubaccount(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSubaccount
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
					return 0, ErrIntOverflowSubaccount
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
					return 0, ErrIntOverflowSubaccount
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
				return 0, ErrInvalidLengthSubaccount
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSubaccount
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSubaccount
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSubaccount        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSubaccount          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSubaccount = fmt.Errorf("proto: unexpected end of group")
)
