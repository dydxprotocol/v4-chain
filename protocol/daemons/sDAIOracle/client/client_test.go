package client_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"cosmossdk.io/log"
	appflags "github.com/StreamFinance-Protocol/stream-chain/protocol/app/flags"
	d_constants "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/constants"
	daemonflags "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/flags"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/api"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/client"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/appoptions"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/daemons"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/grpc"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

var (
	TestError = fmt.Errorf("test error")
)

func TestStart_EthRpcEndpointNotSet(t *testing.T) {
	errorMsg := "flag sDAI-daemon-eth-rpc-endpoint is not set"
	require.EqualError(
		t,
		client.NewClient(log.NewNopLogger()).Start(
			grpc.Ctx,
			daemonflags.GetDefaultDaemonFlags(),
			appflags.GetFlagValuesFromOptions(appoptions.GetDefaultTestAppOptions("", nil)),
			&mocks.GrpcClient{},
		),
		errorMsg,
	)
}

func TestStart_TcpConnectionFails(t *testing.T) {
	errorMsg := "Failed to create connection"

	// Mock the gRPC client to return an error when creating a TCP connection.
	mockGrpcClient := &mocks.GrpcClient{}
	mockGrpcClient.On("NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint).Return(nil, errors.New(errorMsg))

	// Override default daemon flags with a non-empty EthRpcEndpoint.
	daemonFlagsWithEthRpcEndpoint := daemonflags.GetDefaultDaemonFlags()
	daemonFlagsWithEthRpcEndpoint.SDAI.EthRpcEndpoint = "http://localhost:8545"

	require.EqualError(
		t,
		client.NewClient(log.NewNopLogger()).Start(
			grpc.Ctx,
			daemonFlagsWithEthRpcEndpoint,
			appflags.GetFlagValuesFromOptions(appoptions.GetDefaultTestAppOptions("", nil)),
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

	// Mock the gRPC client to
	// - return a successful TCP connection.
	// - return an error when creating a gRPC connection.
	// - successfully close the TCP connection.
	mockGrpcClient := &mocks.GrpcClient{}
	mockGrpcClient.On("NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint).Return(grpc.GrpcConn, nil)
	mockGrpcClient.On("NewGrpcConnection", grpc.Ctx, grpc.SocketPath).Return(nil, errors.New(errorMsg))
	mockGrpcClient.On("CloseConnection", grpc.GrpcConn).Return(nil)

	// Override default daemon flags with a non-empty EthRpcEndpoint.
	daemonFlagsWithEthRpcEndpoint := daemonflags.GetDefaultDaemonFlags()
	daemonFlagsWithEthRpcEndpoint.SDAI.EthRpcEndpoint = "http://localhost:8545"

	require.EqualError(
		t,
		client.NewClient(log.NewNopLogger()).Start(
			grpc.Ctx,
			daemonFlagsWithEthRpcEndpoint,
			appflags.GetFlagValuesFromOptions(appoptions.GetDefaultTestAppOptions("", nil)),
			mockGrpcClient,
		),
		errorMsg,
	)
	mockGrpcClient.AssertCalled(t, "NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint)
	mockGrpcClient.AssertCalled(t, "NewGrpcConnection", grpc.Ctx, grpc.SocketPath)

	// Assert that the connection from NewTcpConnection is closed.
	mockGrpcClient.AssertNumberOfCalls(t, "CloseConnection", 1)
	mockGrpcClient.AssertCalled(t, "CloseConnection", grpc.GrpcConn)
}

// FakeSubTaskRunner is a mock implementation of SubTaskRunner that returns the specified results in order.
type FakeSubTaskRunner struct {
	results   []error
	callIndex int
}

func NewFakeSubTaskRunnerWithResults(results []error) *FakeSubTaskRunner {
	return &FakeSubTaskRunner{
		results:   results,
		callIndex: -1,
	}
}

func (f *FakeSubTaskRunner) RunsDAIDaemonTaskLoop(
	_ context.Context,
	_ *client.Client,
	_ log.Logger,
	_ *ethclient.Client,
	_ api.SDAIServiceClient,
) error {
	f.callIndex += 1
	return f.results[f.callIndex]
}

func TestHealthCheck_Mixed(t *testing.T) {
	tests := map[string]struct {
		// updateResult represents the list of responses for individual daemon task loops. Add a nil value to represent
		// a successful update.
		updateResults        []error
		expectedHealthStatus error
	}{
		"unhealthy: no updates": {
			updateResults:        []error{},
			expectedHealthStatus: fmt.Errorf("no successful update has occurred"),
		},
		"unhealthy: no successful updates": {
			updateResults: []error{
				TestError, // failed update
			},
			expectedHealthStatus: fmt.Errorf("no successful update has occurred"),
		},
		"healthy: one recent successful update": {
			updateResults: []error{
				nil, // successful update
			},
			expectedHealthStatus: nil,
		},
		"unhealthy: one recent successful update, followed by a failed update": {
			updateResults: []error{
				nil,       // successful update
				TestError, // failed update
			},
			expectedHealthStatus: fmt.Errorf("last update failed"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c := client.NewClient(log.NewNopLogger())

			fakeSubTaskRunner := NewFakeSubTaskRunnerWithResults(tc.updateResults)

			for i := 0; i < len(tc.updateResults); i++ {
				ticker, stop := daemons.SingleTickTickerAndStop()
				client.StartsDAIDaemonTaskLoop(
					grpc.Ctx,
					c,
					ticker,
					stop,
					fakeSubTaskRunner,
					nil,
					nil,
				)
			}

			if tc.expectedHealthStatus == nil {
				require.NoError(t, c.HealthCheck())
			} else {
				require.ErrorContains(t, c.HealthCheck(), tc.expectedHealthStatus.Error())
			}
		})
	}
}
