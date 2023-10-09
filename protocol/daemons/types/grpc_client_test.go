package types

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	socketPath      = "/tmp/daemons.sock"
	defaultEndpoint = "localhost:9090"
)

const (
	timeout = time.Second * 3
)

var (
	client = (*GrpcClientImpl)(nil)
)

func TestNewGrpcConnection_withTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()
	_, err := client.NewGrpcConnection(
		ctx,
		socketPath,
	)

	require.EqualError(
		t,
		err,
		"context deadline exceeded",
	)
}

func TestNewGrpcTcpConnection_withTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()
	_, err := client.NewTcpConnection(
		ctx,
		defaultEndpoint,
	)

	require.EqualError(
		t,
		err,
		"context deadline exceeded",
	)
}
