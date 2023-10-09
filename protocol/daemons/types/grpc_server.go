package types

import (
	"net"

	"google.golang.org/grpc"
)

// Ensure the `GrpcServer` interface is implemented at compile time.
var _ GrpcServer = (*grpc.Server)(nil)

// GrpcServer is an interface that encapsulates a `Grpc.Server` object.
type GrpcServer interface {
	Serve(lis net.Listener) error
	Stop()
	RegisterService(sd *grpc.ServiceDesc, ss interface{})
}
