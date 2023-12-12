package client_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/types/query"
	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	d_constants "github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/client"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/appoptions"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	daemontestutils "github.com/dydxprotocol/v4-chain/protocol/testutil/daemons"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/grpc"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
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

func TestRunLiquidationDaemonTaskLoop(t *testing.T) {
	df := flags.GetDefaultDaemonFlags()
	tests := map[string]struct {
		// mocks
		setupMocks func(ctx context.Context, mck *mocks.QueryClient)

		// expectations
		expectedLiquidatableSubaccountIds []satypes.SubaccountId
		expectedError                     error
	}{
		"Success": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &satypes.QueryAllSubaccountRequest{
					Pagination: &query.PageRequest{
						Limit: df.Liquidation.SubaccountPageLimit,
					},
				}
				response := &satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Carl_Num0_1BTC_Short,
						constants.Dave_Num0_1BTC_Long_50000USD,
					},
				}
				mck.On("SubaccountAll", ctx, req).Return(response, nil)

				req2 := &clobtypes.AreSubaccountsLiquidatableRequest{
					SubaccountIds: []satypes.SubaccountId{
						constants.Carl_Num0,
						constants.Dave_Num0,
					},
				}
				response2 := &clobtypes.AreSubaccountsLiquidatableResponse{
					Results: []clobtypes.AreSubaccountsLiquidatableResponse_Result{
						{
							SubaccountId:   constants.Carl_Num0,
							IsLiquidatable: true,
						},
						{
							SubaccountId:   constants.Dave_Num0,
							IsLiquidatable: false,
						},
					},
				}
				mck.On("AreSubaccountsLiquidatable", ctx, req2).Return(response2, nil)

				req3 := &api.LiquidateSubaccountsRequest{
					SubaccountIds: []satypes.SubaccountId{
						constants.Carl_Num0,
					},
				}
				response3 := &api.LiquidateSubaccountsResponse{}
				mck.On("LiquidateSubaccounts", ctx, req3).Return(response3, nil)
			},
		},
		"Success - no open position": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &satypes.QueryAllSubaccountRequest{
					Pagination: &query.PageRequest{
						Limit: df.Liquidation.SubaccountPageLimit,
					},
				}
				response := &satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Carl_Num0_599USD, // no open positions
						constants.Dave_Num0_599USD, // no open positions
					},
				}
				mck.On("SubaccountAll", ctx, req).Return(response, nil)
				req2 := &api.LiquidateSubaccountsRequest{
					SubaccountIds: []satypes.SubaccountId{},
				}
				response2 := &api.LiquidateSubaccountsResponse{}
				mck.On("LiquidateSubaccounts", ctx, req2).Return(response2, nil)
			},
		},
		"Success - no liquidatable subaccounts": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &satypes.QueryAllSubaccountRequest{
					Pagination: &query.PageRequest{
						Limit: df.Liquidation.SubaccountPageLimit,
					},
				}
				response := &satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Carl_Num0_1BTC_Short,
						constants.Dave_Num0_1BTC_Long_50000USD,
					},
				}
				mck.On("SubaccountAll", ctx, req).Return(response, nil)

				req2 := &clobtypes.AreSubaccountsLiquidatableRequest{
					SubaccountIds: []satypes.SubaccountId{
						constants.Carl_Num0,
						constants.Dave_Num0,
					},
				}
				response2 := &clobtypes.AreSubaccountsLiquidatableResponse{
					Results: []clobtypes.AreSubaccountsLiquidatableResponse_Result{
						{
							SubaccountId:   constants.Carl_Num0,
							IsLiquidatable: false,
						},
						{
							SubaccountId:   constants.Dave_Num0,
							IsLiquidatable: false,
						},
					},
				}
				mck.On("AreSubaccountsLiquidatable", ctx, req2).Return(response2, nil)
				req3 := &api.LiquidateSubaccountsRequest{
					SubaccountIds: []satypes.SubaccountId{},
				}
				response3 := &api.LiquidateSubaccountsResponse{}
				mck.On("LiquidateSubaccounts", ctx, req3).Return(response3, nil)
			},
		},
		"Panics on error - SubaccountAll": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				mck.On("SubaccountAll", mock.Anything, mock.Anything).Return(nil, errors.New("test error"))
			},
			expectedError: errors.New("test error"),
		},
		"Panics on error - AreSubaccountsLiquidatable": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				mck.On("SubaccountAll", mock.Anything, mock.Anything).Return(&satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Carl_Num0_1BTC_Short,
					},
				}, nil)
				mck.On("AreSubaccountsLiquidatable", mock.Anything, mock.Anything).Return(nil, errors.New("test error"))
			},
			expectedError: errors.New("test error"),
		},
		"Panics on error - LiquidateSubaccounts": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				mck.On("SubaccountAll", mock.Anything, mock.Anything).Return(&satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Carl_Num0_1BTC_Short,
					},
				}, nil,
				)
				mck.On("AreSubaccountsLiquidatable", mock.Anything, mock.Anything).Return(
					&clobtypes.AreSubaccountsLiquidatableResponse{},
					nil,
				)
				mck.On("LiquidateSubaccounts", mock.Anything, mock.Anything).Return(nil, errors.New("test error"))
			},
			expectedError: errors.New("test error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			queryClientMock := &mocks.QueryClient{}
			tc.setupMocks(grpc.Ctx, queryClientMock)
			s := client.SubTaskRunnerImpl{}

			c := client.NewClient(log.NewNopLogger())
			c.SubaccountQueryClient = queryClientMock
			c.ClobQueryClient = queryClientMock
			c.LiquidationServiceClient = queryClientMock

			err := s.RunLiquidationDaemonTaskLoop(
				grpc.Ctx,
				c,
				flags.GetDefaultDaemonFlags().Liquidation,
			)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				queryClientMock.AssertExpectations(t)
			}
		})
	}
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
