package client_test

import (
	"errors"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/dydxprotocol/v4/daemons/bridge/client"
	d_constants "github.com/dydxprotocol/v4/daemons/constants"
	"github.com/dydxprotocol/v4/daemons/flags"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/grpc"
	"github.com/stretchr/testify/require"
)

func TestStart_TcpConnectionFails(t *testing.T) {
	errorMsg := "Failed to create connection"

	mockGrpcClient := &mocks.GrpcClient{}
	mockGrpcClient.On("NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint).Return(nil, errors.New(errorMsg))

	require.EqualError(
		t,
		client.Start(
			grpc.Ctx,
			flags.GetDefaultDaemonFlags(),
			log.NewNopLogger(),
			mockGrpcClient,
		),
		errorMsg,
	)
	mockGrpcClient.AssertCalled(t, "NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint)
	mockGrpcClient.AssertNotCalled(t, "NewGrpcConnection", grpc.Ctx, grpc.SocketPath)
	mockGrpcClient.AssertNotCalled(t, "CloseConnection", grpc.GrpcConn)
}

func TestStart_UnixSocketConnectionFails(t *testing.T) {
	errorMsg := "Failed to create connection"

	mockGrpcClient := &mocks.GrpcClient{}
	mockGrpcClient.On("NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint).Return(grpc.GrpcConn, nil)
	mockGrpcClient.On("NewGrpcConnection", grpc.Ctx, grpc.SocketPath).Return(nil, errors.New(errorMsg))
	mockGrpcClient.On("CloseConnection", grpc.GrpcConn).Return(nil)

	require.EqualError(
		t,
		client.Start(
			grpc.Ctx,
			flags.GetDefaultDaemonFlags(),
			log.NewNopLogger(),
			mockGrpcClient,
		),
		errorMsg,
	)
	mockGrpcClient.AssertCalled(t, "NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint)
	mockGrpcClient.AssertCalled(t, "NewGrpcConnection", grpc.Ctx, grpc.SocketPath)
	mockGrpcClient.AssertNumberOfCalls(t, "CloseConnection", 1)
}
