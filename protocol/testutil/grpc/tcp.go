package grpc

import "github.com/dydxprotocol/v4-chain/protocol/mocks"

var (
	TcpEndpoint = "localhost:9090"
)

// GenerateMockGrpcClientWithOptionalTcpConnectionErrors generates a mock gRPC client that mocks Tcp connections and
// optionally returns the given connection and close errors. This mock also mocks grpc connections if the tcp
// connection is expected to be closed.
func GenerateMockGrpcClientWithOptionalTcpConnectionErrors(
	connectionErr error,
	closeErr error,
	closeConnectionIsCalled bool,
) *mocks.GrpcClient {
	mockGrpcClient := &mocks.GrpcClient{}

	// Conditionally set up Tcp connection to return the given connection and close errors.
	mockGrpcClient.On("NewTcpConnection", Ctx, TcpEndpoint).
		Return(TcpConn, connectionErr)

	if closeErr != nil || closeConnectionIsCalled {
		mockGrpcClient.On("NewGrpcConnection", Ctx, SocketPath).
			Return(GrpcConn, nil)

		mockGrpcClient.On("CloseConnection", GrpcConn).
			Return(nil)

		mockGrpcClient.On("CloseConnection", TcpConn).
			Return(closeErr)
	}

	return mockGrpcClient
}
