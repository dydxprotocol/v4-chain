// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/govplus/query.proto

package types

import (
	context "context"
	fmt "fmt"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	grpc "google.golang.org/grpc"
	math "math"
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

func init() { proto.RegisterFile("dydxprotocol/govplus/query.proto", fileDescriptor_b04547eded25b028) }

var fileDescriptor_b04547eded25b028 = []byte{
	// 152 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x48, 0xa9, 0x4c, 0xa9,
	0x28, 0x28, 0xca, 0x2f, 0xc9, 0x4f, 0xce, 0xcf, 0xd1, 0x4f, 0xcf, 0x2f, 0x2b, 0xc8, 0x29, 0x2d,
	0xd6, 0x2f, 0x2c, 0x4d, 0x2d, 0xaa, 0xd4, 0x03, 0x0b, 0x0b, 0x89, 0x20, 0xab, 0xd0, 0x83, 0xaa,
	0x30, 0x62, 0xe7, 0x62, 0x0d, 0x04, 0x29, 0x72, 0x4a, 0x3c, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23,
	0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6,
	0x63, 0x39, 0x86, 0x28, 0xf7, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xfd,
	0xe0, 0x92, 0xa2, 0xd4, 0xc4, 0x5c, 0xb7, 0xcc, 0xbc, 0xc4, 0xbc, 0xe4, 0x54, 0xdd, 0x00, 0x98,
	0x7d, 0xc5, 0x60, 0x61, 0xdd, 0xe4, 0x8c, 0xc4, 0xcc, 0x3c, 0x7d, 0xb8, 0x2b, 0x2a, 0xe0, 0xee,
	0x28, 0xa9, 0x2c, 0x48, 0x2d, 0x4e, 0x62, 0x03, 0xcb, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff,
	0xa6, 0x86, 0xe9, 0xb0, 0xac, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

// QueryServer is the server API for Query service.
type QueryServer interface {
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "dydxprotocol.govplus.Query",
	HandlerType: (*QueryServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "dydxprotocol/govplus/query.proto",
}
