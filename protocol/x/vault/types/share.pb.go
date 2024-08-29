// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/vault/share.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

// NumShares represents the number of shares.
type NumShares struct {
	// Number of shares.
	NumShares github_com_dydxprotocol_v4_chain_protocol_dtypes.SerializableInt `protobuf:"bytes,2,opt,name=num_shares,json=numShares,proto3,customtype=github.com/dydxprotocol/v4-chain/protocol/dtypes.SerializableInt" json:"num_shares"`
}

func (m *NumShares) Reset()         { *m = NumShares{} }
func (m *NumShares) String() string { return proto.CompactTextString(m) }
func (*NumShares) ProtoMessage()    {}
func (*NumShares) Descriptor() ([]byte, []int) {
	return fileDescriptor_b65bab3182bbf56c, []int{0}
}
func (m *NumShares) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NumShares) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NumShares.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NumShares) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NumShares.Merge(m, src)
}
func (m *NumShares) XXX_Size() int {
	return m.Size()
}
func (m *NumShares) XXX_DiscardUnknown() {
	xxx_messageInfo_NumShares.DiscardUnknown(m)
}

var xxx_messageInfo_NumShares proto.InternalMessageInfo

// OwnerShare is a type for owner shares.
type OwnerShare struct {
	Owner  string    `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	Shares NumShares `protobuf:"bytes,2,opt,name=shares,proto3" json:"shares"`
}

func (m *OwnerShare) Reset()         { *m = OwnerShare{} }
func (m *OwnerShare) String() string { return proto.CompactTextString(m) }
func (*OwnerShare) ProtoMessage()    {}
func (*OwnerShare) Descriptor() ([]byte, []int) {
	return fileDescriptor_b65bab3182bbf56c, []int{1}
}
func (m *OwnerShare) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OwnerShare) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OwnerShare.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OwnerShare) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OwnerShare.Merge(m, src)
}
func (m *OwnerShare) XXX_Size() int {
	return m.Size()
}
func (m *OwnerShare) XXX_DiscardUnknown() {
	xxx_messageInfo_OwnerShare.DiscardUnknown(m)
}

var xxx_messageInfo_OwnerShare proto.InternalMessageInfo

func (m *OwnerShare) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *OwnerShare) GetShares() NumShares {
	if m != nil {
		return m.Shares
	}
	return NumShares{}
}

// OwnerShareUnlocks stores share unlocks for an owner.
type OwnerShareUnlocks struct {
	// Address of the owner of below shares.
	OwnerAddress string `protobuf:"bytes,1,opt,name=owner_address,json=ownerAddress,proto3" json:"owner_address,omitempty"`
	// All share unlocks.
	ShareUnlocks []ShareUnlock `protobuf:"bytes,2,rep,name=share_unlocks,json=shareUnlocks,proto3" json:"share_unlocks"`
}

func (m *OwnerShareUnlocks) Reset()         { *m = OwnerShareUnlocks{} }
func (m *OwnerShareUnlocks) String() string { return proto.CompactTextString(m) }
func (*OwnerShareUnlocks) ProtoMessage()    {}
func (*OwnerShareUnlocks) Descriptor() ([]byte, []int) {
	return fileDescriptor_b65bab3182bbf56c, []int{2}
}
func (m *OwnerShareUnlocks) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OwnerShareUnlocks) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OwnerShareUnlocks.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OwnerShareUnlocks) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OwnerShareUnlocks.Merge(m, src)
}
func (m *OwnerShareUnlocks) XXX_Size() int {
	return m.Size()
}
func (m *OwnerShareUnlocks) XXX_DiscardUnknown() {
	xxx_messageInfo_OwnerShareUnlocks.DiscardUnknown(m)
}

var xxx_messageInfo_OwnerShareUnlocks proto.InternalMessageInfo

func (m *OwnerShareUnlocks) GetOwnerAddress() string {
	if m != nil {
		return m.OwnerAddress
	}
	return ""
}

func (m *OwnerShareUnlocks) GetShareUnlocks() []ShareUnlock {
	if m != nil {
		return m.ShareUnlocks
	}
	return nil
}

// ShareUnlock stores a single instance of `shares` number of shares
// unlocking at block height `unlock_block_height`.
type ShareUnlock struct {
	// Number of shares to unlock.
	Shares NumShares `protobuf:"bytes,1,opt,name=shares,proto3" json:"shares"`
	// Block height at which above shares unlock.
	UnlockBlockHeight uint32 `protobuf:"varint,2,opt,name=unlock_block_height,json=unlockBlockHeight,proto3" json:"unlock_block_height,omitempty"`
}

func (m *ShareUnlock) Reset()         { *m = ShareUnlock{} }
func (m *ShareUnlock) String() string { return proto.CompactTextString(m) }
func (*ShareUnlock) ProtoMessage()    {}
func (*ShareUnlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_b65bab3182bbf56c, []int{3}
}
func (m *ShareUnlock) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ShareUnlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ShareUnlock.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ShareUnlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShareUnlock.Merge(m, src)
}
func (m *ShareUnlock) XXX_Size() int {
	return m.Size()
}
func (m *ShareUnlock) XXX_DiscardUnknown() {
	xxx_messageInfo_ShareUnlock.DiscardUnknown(m)
}

var xxx_messageInfo_ShareUnlock proto.InternalMessageInfo

func (m *ShareUnlock) GetShares() NumShares {
	if m != nil {
		return m.Shares
	}
	return NumShares{}
}

func (m *ShareUnlock) GetUnlockBlockHeight() uint32 {
	if m != nil {
		return m.UnlockBlockHeight
	}
	return 0
}

func init() {
	proto.RegisterType((*NumShares)(nil), "dydxprotocol.vault.NumShares")
	proto.RegisterType((*OwnerShare)(nil), "dydxprotocol.vault.OwnerShare")
	proto.RegisterType((*OwnerShareUnlocks)(nil), "dydxprotocol.vault.OwnerShareUnlocks")
	proto.RegisterType((*ShareUnlock)(nil), "dydxprotocol.vault.ShareUnlock")
}

func init() { proto.RegisterFile("dydxprotocol/vault/share.proto", fileDescriptor_b65bab3182bbf56c) }

var fileDescriptor_b65bab3182bbf56c = []byte{
	// 391 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0xbd, 0x4e, 0xeb, 0x30,
	0x18, 0x8d, 0x7b, 0xef, 0xad, 0x54, 0xb7, 0x1d, 0xea, 0xdb, 0x21, 0x54, 0x22, 0xad, 0xc2, 0xd2,
	0xa5, 0x8e, 0x54, 0x90, 0x18, 0x58, 0x20, 0x53, 0x61, 0x00, 0x91, 0x8a, 0x85, 0x25, 0xca, 0x9f,
	0x92, 0x88, 0x24, 0xae, 0x62, 0x07, 0xda, 0xce, 0x3c, 0x00, 0x0f, 0xc3, 0x43, 0x74, 0xac, 0x98,
	0x10, 0x43, 0x85, 0xda, 0x17, 0x41, 0xb1, 0x43, 0x09, 0x82, 0x01, 0x16, 0x2b, 0xdf, 0x39, 0x9f,
	0xbf, 0x73, 0x8e, 0xbf, 0x40, 0xc5, 0x9d, 0xb9, 0xd3, 0x49, 0x4a, 0x18, 0x71, 0x48, 0xa4, 0xdd,
	0x5a, 0x59, 0xc4, 0x34, 0x1a, 0x58, 0xa9, 0x87, 0x39, 0x88, 0x50, 0x99, 0xc7, 0x9c, 0xef, 0xb4,
	0x7d, 0xe2, 0x13, 0x8e, 0x69, 0xf9, 0x97, 0xe8, 0xec, 0xec, 0x38, 0x84, 0xc6, 0x84, 0x9a, 0x82,
	0x10, 0x85, 0xa0, 0x54, 0x06, 0x6b, 0xe7, 0x59, 0x3c, 0xce, 0xc7, 0x52, 0xe4, 0x43, 0x98, 0x64,
	0xb1, 0xc9, 0x45, 0xa8, 0x5c, 0xe9, 0x81, 0x7e, 0x43, 0x1f, 0x2d, 0x56, 0x5d, 0xe9, 0x65, 0xd5,
	0x3d, 0xf6, 0x43, 0x16, 0x64, 0x36, 0x76, 0x48, 0xac, 0x7d, 0x36, 0x76, 0x30, 0x70, 0x02, 0x2b,
	0x4c, 0xb4, 0x2d, 0xe2, 0xb2, 0xd9, 0xc4, 0xa3, 0x78, 0xec, 0xa5, 0xa1, 0x15, 0x85, 0x73, 0xcb,
	0x8e, 0xbc, 0xd3, 0x84, 0x19, 0xb5, 0xe4, 0x5d, 0x48, 0x9d, 0x41, 0x78, 0x71, 0x97, 0x78, 0x29,
	0x2f, 0x11, 0x86, 0xff, 0x48, 0x5e, 0xc9, 0xa0, 0x07, 0xfa, 0x35, 0x5d, 0x7e, 0x7a, 0x1c, 0xb4,
	0x0b, 0x93, 0x27, 0xae, 0x9b, 0x7a, 0x94, 0x8e, 0x59, 0x1a, 0x26, 0xbe, 0x21, 0xda, 0xd0, 0x11,
	0xac, 0x96, 0x2c, 0xd6, 0x87, 0xbb, 0xf8, 0xeb, 0x4b, 0xe0, 0x6d, 0x2a, 0xfd, 0x6f, 0x9e, 0xc0,
	0x28, 0xae, 0xa8, 0xf7, 0x00, 0xb6, 0x3e, 0xb4, 0xaf, 0x92, 0x88, 0x38, 0x37, 0x14, 0xed, 0xc1,
	0x26, 0x9f, 0x6d, 0x5a, 0x42, 0x50, 0x58, 0x31, 0x1a, 0x1c, 0x2c, 0x4c, 0xa0, 0x33, 0xd8, 0xe4,
	0x43, 0xcc, 0x4c, 0xdc, 0x92, 0x2b, 0xbd, 0x3f, 0xfd, 0xfa, 0xb0, 0xfb, 0x9d, 0x7c, 0x69, 0x7a,
	0x61, 0xa0, 0x41, 0x4b, 0x82, 0xea, 0x1c, 0xd6, 0x4b, 0x2d, 0xa5, 0x48, 0xe0, 0xd7, 0x91, 0x10,
	0x86, 0xff, 0x85, 0x23, 0xd3, 0xe6, 0x67, 0xe0, 0x85, 0x7e, 0xc0, 0xf8, 0xe3, 0x34, 0x8d, 0x96,
	0xa0, 0xf4, 0xfc, 0x18, 0x71, 0x42, 0xbf, 0x5c, 0xac, 0x15, 0xb0, 0x5c, 0x2b, 0xe0, 0x75, 0xad,
	0x80, 0x87, 0x8d, 0x22, 0x2d, 0x37, 0x8a, 0xf4, 0xbc, 0x51, 0xa4, 0xeb, 0xc3, 0x9f, 0x2f, 0x79,
	0x5a, 0xfc, 0x91, 0x7c, 0xd7, 0x76, 0x95, 0xe3, 0xfb, 0x6f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xbd,
	0x29, 0xf7, 0x40, 0xb4, 0x02, 0x00, 0x00,
}

func (m *NumShares) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NumShares) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NumShares) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.NumShares.Size()
		i -= size
		if _, err := m.NumShares.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintShare(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	return len(dAtA) - i, nil
}

func (m *OwnerShare) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OwnerShare) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OwnerShare) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Shares.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintShare(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintShare(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *OwnerShareUnlocks) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OwnerShareUnlocks) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OwnerShareUnlocks) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ShareUnlocks) > 0 {
		for iNdEx := len(m.ShareUnlocks) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ShareUnlocks[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintShare(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.OwnerAddress) > 0 {
		i -= len(m.OwnerAddress)
		copy(dAtA[i:], m.OwnerAddress)
		i = encodeVarintShare(dAtA, i, uint64(len(m.OwnerAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ShareUnlock) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ShareUnlock) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ShareUnlock) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.UnlockBlockHeight != 0 {
		i = encodeVarintShare(dAtA, i, uint64(m.UnlockBlockHeight))
		i--
		dAtA[i] = 0x10
	}
	{
		size, err := m.Shares.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintShare(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintShare(dAtA []byte, offset int, v uint64) int {
	offset -= sovShare(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *NumShares) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.NumShares.Size()
	n += 1 + l + sovShare(uint64(l))
	return n
}

func (m *OwnerShare) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovShare(uint64(l))
	}
	l = m.Shares.Size()
	n += 1 + l + sovShare(uint64(l))
	return n
}

func (m *OwnerShareUnlocks) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.OwnerAddress)
	if l > 0 {
		n += 1 + l + sovShare(uint64(l))
	}
	if len(m.ShareUnlocks) > 0 {
		for _, e := range m.ShareUnlocks {
			l = e.Size()
			n += 1 + l + sovShare(uint64(l))
		}
	}
	return n
}

func (m *ShareUnlock) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Shares.Size()
	n += 1 + l + sovShare(uint64(l))
	if m.UnlockBlockHeight != 0 {
		n += 1 + sovShare(uint64(m.UnlockBlockHeight))
	}
	return n
}

func sovShare(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozShare(x uint64) (n int) {
	return sovShare(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *NumShares) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowShare
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
			return fmt.Errorf("proto: NumShares: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NumShares: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumShares", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShare
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
				return ErrInvalidLengthShare
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthShare
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.NumShares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipShare(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthShare
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
func (m *OwnerShare) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowShare
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
			return fmt.Errorf("proto: OwnerShare: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OwnerShare: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShare
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
				return ErrInvalidLengthShare
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthShare
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Shares", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShare
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
				return ErrInvalidLengthShare
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthShare
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Shares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipShare(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthShare
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
func (m *OwnerShareUnlocks) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowShare
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
			return fmt.Errorf("proto: OwnerShareUnlocks: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OwnerShareUnlocks: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OwnerAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShare
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
				return ErrInvalidLengthShare
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthShare
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OwnerAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ShareUnlocks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShare
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
				return ErrInvalidLengthShare
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthShare
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ShareUnlocks = append(m.ShareUnlocks, ShareUnlock{})
			if err := m.ShareUnlocks[len(m.ShareUnlocks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipShare(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthShare
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
func (m *ShareUnlock) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowShare
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
			return fmt.Errorf("proto: ShareUnlock: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ShareUnlock: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Shares", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShare
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
				return ErrInvalidLengthShare
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthShare
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Shares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnlockBlockHeight", wireType)
			}
			m.UnlockBlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowShare
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.UnlockBlockHeight |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipShare(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthShare
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
func skipShare(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowShare
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
					return 0, ErrIntOverflowShare
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
					return 0, ErrIntOverflowShare
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
				return 0, ErrInvalidLengthShare
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupShare
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthShare
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthShare        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowShare          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupShare = fmt.Errorf("proto: unexpected end of group")
)
