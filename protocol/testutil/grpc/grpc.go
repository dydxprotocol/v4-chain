package grpc

import (
	"context"
	"os"

	"github.com/dydxprotocol/v4/mocks"
	"google.golang.org/grpc"
)

const (
	SocketPath                   = "path.txt"
	UserReadWriteOnlyPermissions = os.FileMode(0600)
	DefaultPermissions           = os.FileMode(0x180)
)

var (
	Ctx        = context.TODO()
	ClientConn = (*grpc.ClientConn)(nil)
)

// GenerateMockGrpcClientWithReturns generates a mock gRPC client
// that returns the given connection and close errors.
func GenerateMockGrpcClientWithReturns(
	connectionErr error,
	closeErr error,
	closeConnectionIsCalled bool,
) *mocks.GrpcClient {
	mockGrpcClient := &mocks.GrpcClient{}

	mockGrpcClient.On("NewGrpcConnection", Ctx, SocketPath).
		Return(ClientConn, connectionErr)

	if closeErr != nil || closeConnectionIsCalled {
		mockGrpcClient.On("CloseConnection", ClientConn).
			Return(closeErr)
	}

	return mockGrpcClient
}
