// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/clob/finalize_block.proto

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

// ClobStagedFinalizeBlockEvent defines a CLOB event staged during
// FinalizeBlock.
type ClobStagedFinalizeBlockEvent struct {
	// event is the staged event.
	//
	// Types that are valid to be assigned to Event:
	//	*ClobStagedFinalizeBlockEvent_CreateClobPair
	Event isClobStagedFinalizeBlockEvent_Event `protobuf_oneof:"event"`
}

func (m *ClobStagedFinalizeBlockEvent) Reset()         { *m = ClobStagedFinalizeBlockEvent{} }
func (m *ClobStagedFinalizeBlockEvent) String() string { return proto.CompactTextString(m) }
func (*ClobStagedFinalizeBlockEvent) ProtoMessage()    {}
func (*ClobStagedFinalizeBlockEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce1d49660993e938, []int{0}
}
func (m *ClobStagedFinalizeBlockEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ClobStagedFinalizeBlockEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ClobStagedFinalizeBlockEvent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ClobStagedFinalizeBlockEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClobStagedFinalizeBlockEvent.Merge(m, src)
}
func (m *ClobStagedFinalizeBlockEvent) XXX_Size() int {
	return m.Size()
}
func (m *ClobStagedFinalizeBlockEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_ClobStagedFinalizeBlockEvent.DiscardUnknown(m)
}

var xxx_messageInfo_ClobStagedFinalizeBlockEvent proto.InternalMessageInfo

type isClobStagedFinalizeBlockEvent_Event interface {
	isClobStagedFinalizeBlockEvent_Event()
	MarshalTo([]byte) (int, error)
	Size() int
}

type ClobStagedFinalizeBlockEvent_CreateClobPair struct {
	CreateClobPair *ClobPair `protobuf:"bytes,1,opt,name=create_clob_pair,json=createClobPair,proto3,oneof" json:"create_clob_pair,omitempty"`
}

func (*ClobStagedFinalizeBlockEvent_CreateClobPair) isClobStagedFinalizeBlockEvent_Event() {}

func (m *ClobStagedFinalizeBlockEvent) GetEvent() isClobStagedFinalizeBlockEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (m *ClobStagedFinalizeBlockEvent) GetCreateClobPair() *ClobPair {
	if x, ok := m.GetEvent().(*ClobStagedFinalizeBlockEvent_CreateClobPair); ok {
		return x.CreateClobPair
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*ClobStagedFinalizeBlockEvent) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*ClobStagedFinalizeBlockEvent_CreateClobPair)(nil),
	}
}

func init() {
	proto.RegisterType((*ClobStagedFinalizeBlockEvent)(nil), "dydxprotocol.clob.ClobStagedFinalizeBlockEvent")
}

func init() {
	proto.RegisterFile("dydxprotocol/clob/finalize_block.proto", fileDescriptor_ce1d49660993e938)
}

var fileDescriptor_ce1d49660993e938 = []byte{
	// 219 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4b, 0xa9, 0x4c, 0xa9,
	0x28, 0x28, 0xca, 0x2f, 0xc9, 0x4f, 0xce, 0xcf, 0xd1, 0x4f, 0xce, 0xc9, 0x4f, 0xd2, 0x4f, 0xcb,
	0xcc, 0x4b, 0xcc, 0xc9, 0xac, 0x4a, 0x8d, 0x4f, 0xca, 0xc9, 0x4f, 0xce, 0xd6, 0x03, 0x4b, 0x0a,
	0x09, 0x22, 0xab, 0xd3, 0x03, 0xa9, 0x93, 0x52, 0xc4, 0xd4, 0x0a, 0x22, 0xe2, 0x0b, 0x12, 0x33,
	0x8b, 0x20, 0xba, 0x94, 0x0a, 0xb8, 0x64, 0x9c, 0x73, 0xf2, 0x93, 0x82, 0x4b, 0x12, 0xd3, 0x53,
	0x53, 0xdc, 0xa0, 0xe6, 0x3a, 0x81, 0x8c, 0x75, 0x2d, 0x4b, 0xcd, 0x2b, 0x11, 0x72, 0xe7, 0x12,
	0x48, 0x2e, 0x4a, 0x4d, 0x2c, 0x49, 0x8d, 0x87, 0xeb, 0x94, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x36,
	0x92, 0xd6, 0xc3, 0xb0, 0x50, 0x0f, 0x64, 0x54, 0x40, 0x62, 0x66, 0x91, 0x07, 0x43, 0x10, 0x1f,
	0x44, 0x1b, 0x4c, 0xc4, 0x89, 0x9d, 0x8b, 0x35, 0x15, 0x64, 0xa2, 0x53, 0xc0, 0x89, 0x47, 0x72,
	0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24, 0xc7, 0x38, 0xe1, 0xb1, 0x1c, 0xc3, 0x85, 0xc7,
	0x72, 0x0c, 0x37, 0x1e, 0xcb, 0x31, 0x44, 0x99, 0xa5, 0x67, 0x96, 0x64, 0x94, 0x26, 0xe9, 0x25,
	0xe7, 0xe7, 0xea, 0xa3, 0xb8, 0xbc, 0xcc, 0x44, 0x37, 0x39, 0x23, 0x31, 0x33, 0x4f, 0x1f, 0x2e,
	0x52, 0x01, 0xf1, 0x4d, 0x49, 0x65, 0x41, 0x6a, 0x71, 0x12, 0x1b, 0x58, 0xd8, 0x18, 0x10, 0x00,
	0x00, 0xff, 0xff, 0xc0, 0x74, 0xef, 0x17, 0x2a, 0x01, 0x00, 0x00,
}

func (m *ClobStagedFinalizeBlockEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClobStagedFinalizeBlockEvent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ClobStagedFinalizeBlockEvent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Event != nil {
		{
			size := m.Event.Size()
			i -= size
			if _, err := m.Event.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	return len(dAtA) - i, nil
}

func (m *ClobStagedFinalizeBlockEvent_CreateClobPair) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ClobStagedFinalizeBlockEvent_CreateClobPair) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.CreateClobPair != nil {
		{
			size, err := m.CreateClobPair.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintFinalizeBlock(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func encodeVarintFinalizeBlock(dAtA []byte, offset int, v uint64) int {
	offset -= sovFinalizeBlock(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ClobStagedFinalizeBlockEvent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Event != nil {
		n += m.Event.Size()
	}
	return n
}

func (m *ClobStagedFinalizeBlockEvent_CreateClobPair) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.CreateClobPair != nil {
		l = m.CreateClobPair.Size()
		n += 1 + l + sovFinalizeBlock(uint64(l))
	}
	return n
}

func sovFinalizeBlock(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozFinalizeBlock(x uint64) (n int) {
	return sovFinalizeBlock(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ClobStagedFinalizeBlockEvent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFinalizeBlock
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
			return fmt.Errorf("proto: ClobStagedFinalizeBlockEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClobStagedFinalizeBlockEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreateClobPair", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFinalizeBlock
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
				return ErrInvalidLengthFinalizeBlock
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFinalizeBlock
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &ClobPair{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Event = &ClobStagedFinalizeBlockEvent_CreateClobPair{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFinalizeBlock(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthFinalizeBlock
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
func skipFinalizeBlock(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFinalizeBlock
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
					return 0, ErrIntOverflowFinalizeBlock
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
					return 0, ErrIntOverflowFinalizeBlock
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
				return 0, ErrInvalidLengthFinalizeBlock
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupFinalizeBlock
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthFinalizeBlock
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthFinalizeBlock        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFinalizeBlock          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupFinalizeBlock = fmt.Errorf("proto: unexpected end of group")
)
