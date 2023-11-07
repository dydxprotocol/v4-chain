// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/bridge/params.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	_ "github.com/cosmos/gogoproto/types"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// EventParams stores parameters about which events to recognize and which
// tokens to mint.
type EventParams struct {
	// The denom of the token to mint.
	Denom string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	// The numerical chain ID of the Ethereum chain to query.
	EthChainId uint64 `protobuf:"varint,2,opt,name=eth_chain_id,json=ethChainId,proto3" json:"eth_chain_id,omitempty"`
	// The address of the Ethereum contract to monitor for logs.
	EthAddress string `protobuf:"bytes,3,opt,name=eth_address,json=ethAddress,proto3" json:"eth_address,omitempty"`
}

func (m *EventParams) Reset()         { *m = EventParams{} }
func (m *EventParams) String() string { return proto.CompactTextString(m) }
func (*EventParams) ProtoMessage()    {}
func (*EventParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_29afb5e8a05168cd, []int{0}
}
func (m *EventParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventParams.Merge(m, src)
}
func (m *EventParams) XXX_Size() int {
	return m.Size()
}
func (m *EventParams) XXX_DiscardUnknown() {
	xxx_messageInfo_EventParams.DiscardUnknown(m)
}

var xxx_messageInfo_EventParams proto.InternalMessageInfo

func (m *EventParams) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *EventParams) GetEthChainId() uint64 {
	if m != nil {
		return m.EthChainId
	}
	return 0
}

func (m *EventParams) GetEthAddress() string {
	if m != nil {
		return m.EthAddress
	}
	return ""
}

// ProposeParams stores parameters for proposing to the module.
type ProposeParams struct {
	// The maximum number of bridge events to propose per block.
	// Limits the number of events to propose in a single block
	// in-order to smooth out the flow of events.
	MaxBridgesPerBlock uint32 `protobuf:"varint,1,opt,name=max_bridges_per_block,json=maxBridgesPerBlock,proto3" json:"max_bridges_per_block,omitempty"`
	// The minimum duration to wait between a finalized bridge and
	// proposing it. This allows other validators to have enough time to
	// also recognize its occurence. Therefore the bridge daemon should
	// pool for new finalized events at least as often as this parameter.
	ProposeDelayDuration time.Duration `protobuf:"bytes,2,opt,name=propose_delay_duration,json=proposeDelayDuration,proto3,stdduration" json:"propose_delay_duration"`
	// Do not propose any events if a [0, 1_000_000) random number generator
	// generates a number smaller than this number.
	// Setting this parameter to 1_000_000 means always skipping proposing events.
	SkipRatePpm uint32 `protobuf:"varint,3,opt,name=skip_rate_ppm,json=skipRatePpm,proto3" json:"skip_rate_ppm,omitempty"`
	// Do not propose any events if the timestamp of the proposal block is
	// behind the proposers' wall-clock by at least this duration.
	SkipIfBlockDelayedByDuration time.Duration `protobuf:"bytes,4,opt,name=skip_if_block_delayed_by_duration,json=skipIfBlockDelayedByDuration,proto3,stdduration" json:"skip_if_block_delayed_by_duration"`
}

func (m *ProposeParams) Reset()         { *m = ProposeParams{} }
func (m *ProposeParams) String() string { return proto.CompactTextString(m) }
func (*ProposeParams) ProtoMessage()    {}
func (*ProposeParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_29afb5e8a05168cd, []int{1}
}
func (m *ProposeParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProposeParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProposeParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProposeParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProposeParams.Merge(m, src)
}
func (m *ProposeParams) XXX_Size() int {
	return m.Size()
}
func (m *ProposeParams) XXX_DiscardUnknown() {
	xxx_messageInfo_ProposeParams.DiscardUnknown(m)
}

var xxx_messageInfo_ProposeParams proto.InternalMessageInfo

func (m *ProposeParams) GetMaxBridgesPerBlock() uint32 {
	if m != nil {
		return m.MaxBridgesPerBlock
	}
	return 0
}

func (m *ProposeParams) GetProposeDelayDuration() time.Duration {
	if m != nil {
		return m.ProposeDelayDuration
	}
	return 0
}

func (m *ProposeParams) GetSkipRatePpm() uint32 {
	if m != nil {
		return m.SkipRatePpm
	}
	return 0
}

func (m *ProposeParams) GetSkipIfBlockDelayedByDuration() time.Duration {
	if m != nil {
		return m.SkipIfBlockDelayedByDuration
	}
	return 0
}

// SafetyParams stores safety parameters for the module.
type SafetyParams struct {
	// True if bridging is disabled.
	IsDisabled bool `protobuf:"varint,1,opt,name=is_disabled,json=isDisabled,proto3" json:"is_disabled,omitempty"`
	// The number of blocks that bridges accepted in-consensus will be pending
	// until the minted tokens are granted.
	DelayBlocks uint32 `protobuf:"varint,2,opt,name=delay_blocks,json=delayBlocks,proto3" json:"delay_blocks,omitempty"`
}

func (m *SafetyParams) Reset()         { *m = SafetyParams{} }
func (m *SafetyParams) String() string { return proto.CompactTextString(m) }
func (*SafetyParams) ProtoMessage()    {}
func (*SafetyParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_29afb5e8a05168cd, []int{2}
}
func (m *SafetyParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SafetyParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SafetyParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SafetyParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SafetyParams.Merge(m, src)
}
func (m *SafetyParams) XXX_Size() int {
	return m.Size()
}
func (m *SafetyParams) XXX_DiscardUnknown() {
	xxx_messageInfo_SafetyParams.DiscardUnknown(m)
}

var xxx_messageInfo_SafetyParams proto.InternalMessageInfo

func (m *SafetyParams) GetIsDisabled() bool {
	if m != nil {
		return m.IsDisabled
	}
	return false
}

func (m *SafetyParams) GetDelayBlocks() uint32 {
	if m != nil {
		return m.DelayBlocks
	}
	return 0
}

func init() {
	proto.RegisterType((*EventParams)(nil), "dydxprotocol.bridge.EventParams")
	proto.RegisterType((*ProposeParams)(nil), "dydxprotocol.bridge.ProposeParams")
	proto.RegisterType((*SafetyParams)(nil), "dydxprotocol.bridge.SafetyParams")
}

func init() { proto.RegisterFile("dydxprotocol/bridge/params.proto", fileDescriptor_29afb5e8a05168cd) }

var fileDescriptor_29afb5e8a05168cd = []byte{
	// 442 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xbd, 0x92, 0xd3, 0x3c,
	0x14, 0x8d, 0xf7, 0xdb, 0x8f, 0x59, 0xe4, 0xb8, 0x11, 0x81, 0x09, 0x3b, 0x8c, 0x93, 0x4d, 0xb5,
	0x0d, 0xf6, 0xf0, 0x53, 0xd0, 0x62, 0x42, 0xb1, 0x5d, 0x46, 0x54, 0xd0, 0x68, 0xe4, 0xe8, 0xda,
	0xd1, 0xac, 0x1d, 0x69, 0x24, 0x65, 0x27, 0x79, 0x0b, 0x4a, 0xde, 0x82, 0xd7, 0xd8, 0x72, 0x4b,
	0x2a, 0x60, 0x92, 0x17, 0x61, 0x7c, 0xe5, 0xec, 0x2c, 0x1d, 0x9d, 0x75, 0xce, 0xf1, 0xb9, 0xe7,
	0x5c, 0x89, 0x4c, 0xe5, 0x4e, 0x6e, 0x8d, 0xd5, 0x5e, 0x2f, 0x75, 0x93, 0x97, 0x56, 0xc9, 0x1a,
	0x72, 0x23, 0xac, 0x68, 0x5d, 0x86, 0x30, 0x7d, 0xf2, 0x50, 0x91, 0x05, 0xc5, 0xf9, 0xa8, 0xd6,
	0xb5, 0x46, 0x30, 0xef, 0xbe, 0x82, 0xf4, 0x3c, 0xad, 0xb5, 0xae, 0x1b, 0xc8, 0xf1, 0x54, 0x6e,
	0xaa, 0x5c, 0x6e, 0xac, 0xf0, 0x4a, 0xaf, 0x03, 0x3f, 0xab, 0x48, 0xfc, 0xf1, 0x06, 0xd6, 0x7e,
	0x81, 0xfe, 0x74, 0x44, 0xfe, 0x97, 0xb0, 0xd6, 0xed, 0x38, 0x9a, 0x46, 0x97, 0x8f, 0x59, 0x38,
	0xd0, 0x29, 0x19, 0x82, 0x5f, 0xf1, 0xe5, 0x4a, 0xa8, 0x35, 0x57, 0x72, 0x7c, 0x32, 0x8d, 0x2e,
	0x4f, 0x19, 0x01, 0xbf, 0xfa, 0xd0, 0x41, 0x57, 0x92, 0x4e, 0x48, 0xdc, 0x29, 0x84, 0x94, 0x16,
	0x9c, 0x1b, 0xff, 0x87, 0x7f, 0x77, 0x82, 0xf7, 0x01, 0x99, 0x7d, 0x3f, 0x21, 0xc9, 0xc2, 0x6a,
	0xa3, 0x1d, 0xf4, 0xa3, 0x5e, 0x91, 0xa7, 0xad, 0xd8, 0xf2, 0x90, 0xde, 0x71, 0x03, 0x96, 0x97,
	0x8d, 0x5e, 0x5e, 0xe3, 0xe8, 0x84, 0xd1, 0x56, 0x6c, 0x8b, 0xc0, 0x2d, 0xc0, 0x16, 0x1d, 0x43,
	0x3f, 0x93, 0x67, 0x26, 0x78, 0x70, 0x09, 0x8d, 0xd8, 0xf1, 0x63, 0x19, 0x4c, 0x14, 0xbf, 0x7e,
	0x9e, 0x85, 0xb6, 0xd9, 0xb1, 0x6d, 0x36, 0xef, 0x05, 0xc5, 0xd9, 0xed, 0xcf, 0xc9, 0xe0, 0xdb,
	0xaf, 0x49, 0xc4, 0x46, 0xbd, 0xc5, 0xbc, 0x73, 0x38, 0xf2, 0x74, 0x46, 0x12, 0x77, 0xad, 0x0c,
	0xb7, 0xc2, 0x03, 0x37, 0xa6, 0xc5, 0x0a, 0x09, 0x8b, 0x3b, 0x90, 0x09, 0x0f, 0x0b, 0xd3, 0xd2,
	0x86, 0x5c, 0xa0, 0x46, 0x55, 0x21, 0x69, 0x08, 0x01, 0x92, 0x97, 0x0f, 0x92, 0x9c, 0xfe, 0x7b,
	0x92, 0x17, 0x9d, 0xdb, 0x55, 0x85, 0xdd, 0xe6, 0xc1, 0xaa, 0xb8, 0x4f, 0x34, 0x63, 0x64, 0xf8,
	0x49, 0x54, 0xe0, 0x77, 0xfd, 0xbe, 0x26, 0x24, 0x56, 0x8e, 0x4b, 0xe5, 0x44, 0xd9, 0x80, 0xc4,
	0x2d, 0x9d, 0x31, 0xa2, 0xdc, 0xbc, 0x47, 0xe8, 0x05, 0x19, 0x86, 0xad, 0x60, 0x38, 0x87, 0x3b,
	0x49, 0x58, 0x8c, 0x18, 0xce, 0x70, 0x05, 0xbb, 0xdd, 0xa7, 0xd1, 0xdd, 0x3e, 0x8d, 0x7e, 0xef,
	0xd3, 0xe8, 0xeb, 0x21, 0x1d, 0xdc, 0x1d, 0xd2, 0xc1, 0x8f, 0x43, 0x3a, 0xf8, 0xf2, 0xae, 0x56,
	0x7e, 0xb5, 0x29, 0xb3, 0xa5, 0x6e, 0xf3, 0xbf, 0xde, 0xdf, 0xcd, 0xdb, 0x97, 0x78, 0xef, 0xf9,
	0x3d, 0xb2, 0x3d, 0xbe, 0x49, 0xbf, 0x33, 0xe0, 0xca, 0x47, 0x48, 0xbc, 0xf9, 0x13, 0x00, 0x00,
	0xff, 0xff, 0x2b, 0x28, 0x4e, 0x34, 0xb7, 0x02, 0x00, 0x00,
}

func (m *EventParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.EthAddress) > 0 {
		i -= len(m.EthAddress)
		copy(dAtA[i:], m.EthAddress)
		i = encodeVarintParams(dAtA, i, uint64(len(m.EthAddress)))
		i--
		dAtA[i] = 0x1a
	}
	if m.EthChainId != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.EthChainId))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintParams(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ProposeParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProposeParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProposeParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.SkipIfBlockDelayedByDuration, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.SkipIfBlockDelayedByDuration):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintParams(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x22
	if m.SkipRatePpm != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SkipRatePpm))
		i--
		dAtA[i] = 0x18
	}
	n2, err2 := github_com_cosmos_gogoproto_types.StdDurationMarshalTo(m.ProposeDelayDuration, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.ProposeDelayDuration):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintParams(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x12
	if m.MaxBridgesPerBlock != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxBridgesPerBlock))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SafetyParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SafetyParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SafetyParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.DelayBlocks != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.DelayBlocks))
		i--
		dAtA[i] = 0x10
	}
	if m.IsDisabled {
		i--
		if m.IsDisabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
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
func (m *EventParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.EthChainId != 0 {
		n += 1 + sovParams(uint64(m.EthChainId))
	}
	l = len(m.EthAddress)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	return n
}

func (m *ProposeParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MaxBridgesPerBlock != 0 {
		n += 1 + sovParams(uint64(m.MaxBridgesPerBlock))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.ProposeDelayDuration)
	n += 1 + l + sovParams(uint64(l))
	if m.SkipRatePpm != 0 {
		n += 1 + sovParams(uint64(m.SkipRatePpm))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdDuration(m.SkipIfBlockDelayedByDuration)
	n += 1 + l + sovParams(uint64(l))
	return n
}

func (m *SafetyParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.IsDisabled {
		n += 2
	}
	if m.DelayBlocks != 0 {
		n += 1 + sovParams(uint64(m.DelayBlocks))
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *EventParams) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: EventParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EthChainId", wireType)
			}
			m.EthChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EthChainId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EthAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EthAddress = string(dAtA[iNdEx:postIndex])
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
func (m *ProposeParams) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: ProposeParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProposeParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxBridgesPerBlock", wireType)
			}
			m.MaxBridgesPerBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxBridgesPerBlock |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProposeDelayDuration", wireType)
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
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.ProposeDelayDuration, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SkipRatePpm", wireType)
			}
			m.SkipRatePpm = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SkipRatePpm |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SkipIfBlockDelayedByDuration", wireType)
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
			if err := github_com_cosmos_gogoproto_types.StdDurationUnmarshal(&m.SkipIfBlockDelayedByDuration, dAtA[iNdEx:postIndex]); err != nil {
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
func (m *SafetyParams) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: SafetyParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SafetyParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsDisabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
			m.IsDisabled = bool(v != 0)
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DelayBlocks", wireType)
			}
			m.DelayBlocks = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DelayBlocks |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
