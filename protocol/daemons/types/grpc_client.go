package types

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
)

// GrpcClientImpl is the struct that implements the `GrpcClient` interface.
type GrpcClientImpl struct{}

// Ensure the `GrpcClient` interface is implemented at compile time.
var _ GrpcClient = (*GrpcClientImpl)(nil)

// GrpcClient is an interface that encapsulates the `NewGrpcConnection` function and `CloseConnection`.
type GrpcClient interface {
	NewGrpcConnection(ctx context.Context, socketAddress string) (*grpc.ClientConn, error)
	NewTcpConnection(ctx context.Context, endpoint string) (*grpc.ClientConn, error)
	CloseConnection(grpcConn *grpc.ClientConn) error
}

// NewGrpcConnection calls `grpc.Dial` with custom parameters to create a secure connection
// with context that blocks until the underlying connection is up.
func (g *GrpcClientImpl) NewGrpcConnection(
	ctx context.Context,
	socketAddress string,
) (*grpc.ClientConn, error) {
	return grpc.DialContext( //nolint:staticcheck
		ctx,
		socketAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// https://github.com/grpc/grpc-go/blob/master/dialoptions.go#L264
		grpc.WithBlock(), //nolint:staticcheck
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			// Create a custom `net.Dialer` in order to specify `unix` as the desired network.
			var dialer net.Dialer
			return dialer.DialContext(ctx, constants.UnixProtocol, addr)
		}),
	)
}

// NewTcpConnection calls `grpc.Dial` to create an insecure tcp connection.
func (g *GrpcClientImpl) NewTcpConnection(
	ctx context.Context,
	endpoint string,
) (*grpc.ClientConn, error) {
	return grpc.DialContext( //nolint:staticcheck
		ctx,
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), //nolint:staticcheck
	)
}

// CloseConnection calls `grpc.ClientConn.Close()` to close a grpc connection.
func (g *GrpcClientImpl) CloseConnection(grpcConn *grpc.ClientConn) error {
	return grpcConn.Close()
}
