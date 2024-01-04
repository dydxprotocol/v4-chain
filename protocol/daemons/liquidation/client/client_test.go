package client_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"cosmossdk.io/log"
	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	d_constants "github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/client"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/appoptions"
	daemontestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/daemons"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/grpc"
	"github.com/stretchr/testify/require"
)

func TestStart_TcpConnectionFails(t *testing.T) {
	errorMsg := "Failed to create connection"

	mockGrpcClient := &mocks.GrpcClient{}
	mockGrpcClient.On("NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint).Return(nil, errors.New(errorMsg))

	liquidationsClient := client.NewClient(log.NewNopLogger())
	require.EqualError(
		t,
		liquidationsClient.Start(
			grpc.Ctx,
			flags.GetDefaultDaemonFlags(),
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

	mockGrpcClient := &mocks.GrpcClient{}
	mockGrpcClient.On("NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint).Return(grpc.GrpcConn, nil)
	mockGrpcClient.On("NewGrpcConnection", grpc.Ctx, grpc.SocketPath).Return(nil, errors.New(errorMsg))
	mockGrpcClient.On("CloseConnection", grpc.GrpcConn).Return(nil)

	liquidationsClient := client.NewClient(log.NewNopLogger())
	require.EqualError(
		t,
		liquidationsClient.Start(
			grpc.Ctx,
			flags.GetDefaultDaemonFlags(),
			appflags.GetFlagValuesFromOptions(appoptions.GetDefaultTestAppOptions("", nil)),
			mockGrpcClient,
		),
		errorMsg,
	)
	mockGrpcClient.AssertCalled(t, "NewTcpConnection", grpc.Ctx, d_constants.DefaultGrpcEndpoint)
	mockGrpcClient.AssertCalled(t, "NewGrpcConnection", grpc.Ctx, grpc.SocketPath)
	mockGrpcClient.AssertNumberOfCalls(t, "CloseConnection", 1)
}

// FakeSubTaskRunner is a mock implementation of the SubTaskRunner interface for testing.
type FakeSubTaskRunner struct {
	err    error
	called bool
}

func NewFakeSubTaskRunnerWithError(err error) *FakeSubTaskRunner {
	return &FakeSubTaskRunner{
		err: err,
	}
}

// RunLiquidationDaemonTaskLoop is a mock implementation of the SubTaskRunner interface. It records the
// call as a sanity check, and returns the error set by NewFakeSubTaskRunnerWithError.
func (f *FakeSubTaskRunner) RunLiquidationDaemonTaskLoop(
	_ context.Context,
	_ *client.Client,
	_ flags.LiquidationFlags,
) error {
	f.called = true
	return f.err
}

func TestHealthCheck_Mixed(t *testing.T) {
	tests := map[string]struct {
		// taskLoopResponses is a list of errors returned by the task loop. If the error is nil, the task loop is
		// considered to have succeeded.
		taskLoopResponses    []error
		expectedHealthStatus error
	}{
		"Healthy - successful update": {
			taskLoopResponses: []error{
				nil, // 1 successful update
			},
			expectedHealthStatus: nil, // healthy status
		},
		"Unhealthy - failed update": {
			taskLoopResponses: []error{
				fmt.Errorf("failed to update"), // 1 failed update
			},
			expectedHealthStatus: fmt.Errorf("no successful update has occurred"),
		},
		"Unhealthy - failed update after successful update": {
			taskLoopResponses: []error{
				nil,                            // 1 successful update
				fmt.Errorf("failed to update"), // 1 failed update
			},
			expectedHealthStatus: fmt.Errorf("last update failed"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			c := client.NewClient(log.NewNopLogger())

			// Sanity check - the client should be unhealthy before the first successful update.
			require.ErrorContains(
				t,
				c.HealthCheck(),
				"no successful update has occurred",
			)

			// Run the sequence of task loop responses.
			for _, taskLoopError := range tc.taskLoopResponses {
				ticker, stop := daemontestutils.SingleTickTickerAndStop()

				c.SubaccountQueryClient = &mocks.QueryClient{}
				c.ClobQueryClient = &mocks.QueryClient{}
				c.LiquidationServiceClient = &mocks.QueryClient{}

				// Start the daemon task loop. Since we created a single-tick ticker, this will run for one iteration and
				// return.
				client.StartLiquidationsDaemonTaskLoop(
					c,
					grpc.Ctx,
					NewFakeSubTaskRunnerWithError(taskLoopError),
					flags.GetDefaultDaemonFlags(),
					ticker,
					stop,
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
