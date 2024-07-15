// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/accountplus/accountplus.proto

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

type TimestampNonceDetails struct {
	// unsorted list of n most recent timestamp nonces
	TimestampNonces []uint64 `protobuf:"varint,1,rep,packed,name=timestamp_nonces,json=timestampNonces,proto3" json:"timestamp_nonces,omitempty"`
	// most recent timestamp nonce that was ejected from list above
	LatestEjectedNonce uint64 `protobuf:"varint,2,opt,name=latest_ejected_nonce,json=latestEjectedNonce,proto3" json:"latest_ejected_nonce,omitempty"`
}

func (m *TimestampNonceDetails) Reset()         { *m = TimestampNonceDetails{} }
func (m *TimestampNonceDetails) String() string { return proto.CompactTextString(m) }
func (*TimestampNonceDetails) ProtoMessage()    {}
func (*TimestampNonceDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_391b06af1cfe6fb0, []int{0}
}
func (m *TimestampNonceDetails) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TimestampNonceDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TimestampNonceDetails.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TimestampNonceDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TimestampNonceDetails.Merge(m, src)
}
func (m *TimestampNonceDetails) XXX_Size() int {
	return m.Size()
}
func (m *TimestampNonceDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_TimestampNonceDetails.DiscardUnknown(m)
}

var xxx_messageInfo_TimestampNonceDetails proto.InternalMessageInfo

func (m *TimestampNonceDetails) GetTimestampNonces() []uint64 {
	if m != nil {
		return m.TimestampNonces
	}
	return nil
}

func (m *TimestampNonceDetails) GetLatestEjectedNonce() uint64 {
	if m != nil {
		return m.LatestEjectedNonce
	}
	return 0
}

func init() {
	proto.RegisterType((*TimestampNonceDetails)(nil), "dydxprotocol.accountplus.TimestampNonceDetails")
}

func init() {
	proto.RegisterFile("dydxprotocol/accountplus/accountplus.proto", fileDescriptor_391b06af1cfe6fb0)
}

var fileDescriptor_391b06af1cfe6fb0 = []byte{
	// 213 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4a, 0xa9, 0x4c, 0xa9,
	0x28, 0x28, 0xca, 0x2f, 0xc9, 0x4f, 0xce, 0xcf, 0xd1, 0x4f, 0x4c, 0x4e, 0xce, 0x2f, 0xcd, 0x2b,
	0x29, 0xc8, 0x29, 0x2d, 0x46, 0x66, 0xeb, 0x81, 0x15, 0x08, 0x49, 0x20, 0xab, 0xd5, 0x43, 0x92,
	0x57, 0x2a, 0xe1, 0x12, 0x0d, 0xc9, 0xcc, 0x4d, 0x2d, 0x2e, 0x49, 0xcc, 0x2d, 0xf0, 0xcb, 0xcf,
	0x4b, 0x4e, 0x75, 0x49, 0x2d, 0x49, 0xcc, 0xcc, 0x29, 0x16, 0xd2, 0xe4, 0x12, 0x28, 0x81, 0x49,
	0xc4, 0xe7, 0x81, 0x64, 0x8a, 0x25, 0x18, 0x15, 0x98, 0x35, 0x58, 0x82, 0xf8, 0x4b, 0x50, 0x34,
	0x14, 0x0b, 0x19, 0x70, 0x89, 0xe4, 0x24, 0x96, 0xa4, 0x16, 0x97, 0xc4, 0xa7, 0x66, 0xa5, 0x26,
	0x97, 0xa4, 0xa6, 0x40, 0xd4, 0x4b, 0x30, 0x29, 0x30, 0x6a, 0xb0, 0x04, 0x09, 0x41, 0xe4, 0x5c,
	0x21, 0x52, 0x60, 0x2d, 0x4e, 0xe1, 0x27, 0x1e, 0xc9, 0x31, 0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0,
	0x91, 0x1c, 0xe3, 0x84, 0xc7, 0x72, 0x0c, 0x17, 0x1e, 0xcb, 0x31, 0xdc, 0x78, 0x2c, 0xc7, 0x10,
	0x65, 0x9b, 0x9e, 0x59, 0x92, 0x51, 0x9a, 0xa4, 0x97, 0x9c, 0x9f, 0xab, 0x8f, 0xe2, 0xc1, 0x32,
	0x13, 0xdd, 0xe4, 0x8c, 0xc4, 0xcc, 0x3c, 0x7d, 0xb8, 0x48, 0x05, 0x8a, 0xa7, 0x4b, 0x2a, 0x0b,
	0x52, 0x8b, 0x93, 0xd8, 0xc0, 0xb2, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf9, 0xc6, 0xa9,
	0xa4, 0x1d, 0x01, 0x00, 0x00,
}

func (m *TimestampNonceDetails) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TimestampNonceDetails) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TimestampNonceDetails) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LatestEjectedNonce != 0 {
		i = encodeVarintAccountplus(dAtA, i, uint64(m.LatestEjectedNonce))
		i--
		dAtA[i] = 0x10
	}
	if len(m.TimestampNonces) > 0 {
		dAtA2 := make([]byte, len(m.TimestampNonces)*10)
		var j1 int
		for _, num := range m.TimestampNonces {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		i -= j1
		copy(dAtA[i:], dAtA2[:j1])
		i = encodeVarintAccountplus(dAtA, i, uint64(j1))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintAccountplus(dAtA []byte, offset int, v uint64) int {
	offset -= sovAccountplus(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TimestampNonceDetails) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.TimestampNonces) > 0 {
		l = 0
		for _, e := range m.TimestampNonces {
			l += sovAccountplus(uint64(e))
		}
		n += 1 + sovAccountplus(uint64(l)) + l
	}
	if m.LatestEjectedNonce != 0 {
		n += 1 + sovAccountplus(uint64(m.LatestEjectedNonce))
	}
	return n
}

func sovAccountplus(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAccountplus(x uint64) (n int) {
	return sovAccountplus(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TimestampNonceDetails) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAccountplus
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
			return fmt.Errorf("proto: TimestampNonceDetails: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TimestampNonceDetails: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowAccountplus
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.TimestampNonces = append(m.TimestampNonces, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowAccountplus
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthAccountplus
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthAccountplus
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.TimestampNonces) == 0 {
					m.TimestampNonces = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowAccountplus
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.TimestampNonces = append(m.TimestampNonces, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field TimestampNonces", wireType)
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LatestEjectedNonce", wireType)
			}
			m.LatestEjectedNonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAccountplus
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LatestEjectedNonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipAccountplus(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAccountplus
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
func skipAccountplus(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAccountplus
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
					return 0, ErrIntOverflowAccountplus
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
					return 0, ErrIntOverflowAccountplus
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
				return 0, ErrInvalidLengthAccountplus
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAccountplus
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAccountplus
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAccountplus        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAccountplus          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAccountplus = fmt.Errorf("proto: unexpected end of group")
)
