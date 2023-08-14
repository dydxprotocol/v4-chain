package grpc

import (
	"context"
	"os"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"google.golang.org/grpc"
)

const (
	SocketPath                   = "/tmp/daemons.sock"
	UserReadWriteOnlyPermissions = os.FileMode(0600)
	DefaultPermissions           = os.FileMode(0x180)
)

var (
	Ctx      = context.TODO()
	TcpConn  = &grpc.ClientConn{}
	GrpcConn = (*grpc.ClientConn)(nil)
)

// GenerateMockGrpcClientWithOptionalGrpcConnectionErrors generates a mock gRPC client that mocks both Tcp and Grpc
// connections and optionally returns the specified errors on Grpc connections.
func GenerateMockGrpcClientWithOptionalGrpcConnectionErrors(
	connectionErr error,
	closeErr error,
	closeConnectionIsCalled bool,
) *mocks.GrpcClient {
	mockGrpcClient := &mocks.GrpcClient{}

	// Conditionally set up Grpc connection to return the given connection and close errors.
	mockGrpcClient.On("NewGrpcConnection", Ctx, SocketPath).
		Return(GrpcConn, connectionErr)

	if closeErr != nil || closeConnectionIsCalled {
		mockGrpcClient.On("CloseConnection", GrpcConn).
			Return(closeErr)
	}
	// Setup Tcp connections to return without error.
	mockGrpcClient.On("NewTcpConnection", Ctx, TcpEndpoint).
		Return(TcpConn, nil)

	mockGrpcClient.On("CloseConnection", TcpConn).
		Return(nil)

	return mockGrpcClient
}
