package client_test

import (
	"context"
	"errors"
	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/appoptions"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/types/query"
	d_constants "github.com/dydxprotocol/v4-chain/protocol/daemons/constants"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/client"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
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

	require.EqualError(
		t,
		client.Start(
			grpc.Ctx,
			flags.GetDefaultDaemonFlags(),
			appflags.GetFlagValuesFromOptions(appoptions.GetDefaultTestAppOptions("", nil)),
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
			appflags.GetFlagValuesFromOptions(appoptions.GetDefaultTestAppOptions("", nil)),
			log.NewNopLogger(),
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

			err := client.RunLiquidationDaemonTaskLoop(
				grpc.Ctx,
				flags.GetDefaultDaemonFlags().Liquidation,
				queryClientMock,
				queryClientMock,
				queryClientMock,
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

func TestGetAllSubaccounts(t *testing.T) {
	df := flags.GetDefaultDaemonFlags()
	tests := map[string]struct {
		// mocks
		setupMocks func(ctx context.Context, mck *mocks.QueryClient)

		// expectations
		expectedSubaccounts []satypes.Subaccount
		expectedError       error
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
						constants.Carl_Num0_599USD,
						constants.Dave_Num0_599USD,
					},
				}
				mck.On("SubaccountAll", ctx, req).Return(response, nil)
			},
			expectedSubaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
				constants.Dave_Num0_599USD,
			},
		},
		"Success Paginated": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &satypes.QueryAllSubaccountRequest{
					Pagination: &query.PageRequest{
						Limit: df.Liquidation.SubaccountPageLimit,
					},
				}
				nextKey := []byte("next key")
				response := &satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Carl_Num0_599USD,
					},
					Pagination: &query.PageResponse{
						NextKey: nextKey,
					},
				}
				mck.On("SubaccountAll", ctx, req).Return(response, nil)
				req2 := &satypes.QueryAllSubaccountRequest{
					Pagination: &query.PageRequest{
						Key:   nextKey,
						Limit: df.Liquidation.SubaccountPageLimit,
					},
				}
				response2 := &satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Dave_Num0_599USD,
					},
				}
				mck.On("SubaccountAll", ctx, req2).Return(response2, nil)
			},
			expectedSubaccounts: []satypes.Subaccount{
				constants.Carl_Num0_599USD,
				constants.Dave_Num0_599USD,
			},
		},
		"Errors are propagated": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &satypes.QueryAllSubaccountRequest{
					Pagination: &query.PageRequest{
						Limit: df.Liquidation.SubaccountPageLimit,
					},
				}
				mck.On("SubaccountAll", ctx, req).Return(nil, errors.New("test error"))
			},
			expectedError: errors.New("test error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			queryClientMock := &mocks.QueryClient{}
			tc.setupMocks(grpc.Ctx, queryClientMock)

			actual, err := client.GetAllSubaccounts(grpc.Ctx, queryClientMock, df.Liquidation.SubaccountPageLimit)
			if err != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.Equal(t, tc.expectedSubaccounts, actual)
			}
		})
	}
}

func TestCheckCollateralizationForSubaccounts(t *testing.T) {
	tests := map[string]struct {
		// mocks
		setupMocks func(
			ctx context.Context,
			mck *mocks.QueryClient,
			results []clobtypes.AreSubaccountsLiquidatableResponse_Result,
		)
		subaccountIds []satypes.SubaccountId

		// expectations
		expectedResults []clobtypes.AreSubaccountsLiquidatableResponse_Result
		expectedError   error
	}{
		"Success": {
			setupMocks: func(
				ctx context.Context,
				mck *mocks.QueryClient,
				results []clobtypes.AreSubaccountsLiquidatableResponse_Result,
			) {
				query := &clobtypes.AreSubaccountsLiquidatableRequest{
					SubaccountIds: []satypes.SubaccountId{
						constants.Alice_Num0,
						constants.Bob_Num0,
					},
				}
				response := &clobtypes.AreSubaccountsLiquidatableResponse{
					Results: results,
				}
				mck.On("AreSubaccountsLiquidatable", ctx, query).Return(response, nil)
			},
			subaccountIds: []satypes.SubaccountId{
				constants.Alice_Num0,
				constants.Bob_Num0,
			},
			expectedResults: []clobtypes.AreSubaccountsLiquidatableResponse_Result{
				{
					SubaccountId:   constants.Alice_Num0,
					IsLiquidatable: true,
				},
				{
					SubaccountId:   constants.Bob_Num0,
					IsLiquidatable: false,
				},
			},
		},
		"Success - Empty": {
			setupMocks: func(
				ctx context.Context,
				mck *mocks.QueryClient,
				results []clobtypes.AreSubaccountsLiquidatableResponse_Result,
			) {
				query := &clobtypes.AreSubaccountsLiquidatableRequest{
					SubaccountIds: []satypes.SubaccountId{},
				}
				response := &clobtypes.AreSubaccountsLiquidatableResponse{
					Results: results,
				}
				mck.On("AreSubaccountsLiquidatable", ctx, query).Return(response, nil)
			},
			subaccountIds:   []satypes.SubaccountId{},
			expectedResults: []clobtypes.AreSubaccountsLiquidatableResponse_Result{},
		},
		"Errors are propagated": {
			setupMocks: func(
				ctx context.Context,
				mck *mocks.QueryClient,
				results []clobtypes.AreSubaccountsLiquidatableResponse_Result,
			) {
				query := &clobtypes.AreSubaccountsLiquidatableRequest{
					SubaccountIds: []satypes.SubaccountId{},
				}
				mck.On("AreSubaccountsLiquidatable", ctx, query).Return(nil, errors.New("test error"))
			},
			subaccountIds:   []satypes.SubaccountId{},
			expectedResults: []clobtypes.AreSubaccountsLiquidatableResponse_Result{},
			expectedError:   errors.New("test error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			queryClientMock := &mocks.QueryClient{}
			tc.setupMocks(grpc.Ctx, queryClientMock, tc.expectedResults)

			actual, err := client.CheckCollateralizationForSubaccounts(grpc.Ctx, queryClientMock, tc.subaccountIds)
			if err != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.Equal(t, tc.expectedResults, actual)
			}
		})
	}
}

func TestSendLiquidatableSubaccountIds(t *testing.T) {
	tests := map[string]struct {
		// mocks
		setupMocks    func(ctx context.Context, mck *mocks.QueryClient, ids []satypes.SubaccountId)
		subaccountIds []satypes.SubaccountId

		// expectations
		expectedError error
	}{
		"Success": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient, ids []satypes.SubaccountId) {
				req := &api.LiquidateSubaccountsRequest{
					SubaccountIds: ids,
				}
				response := &api.LiquidateSubaccountsResponse{}
				mck.On("LiquidateSubaccounts", ctx, req).Return(response, nil)
			},
			subaccountIds: []satypes.SubaccountId{
				constants.Alice_Num0,
				constants.Bob_Num0,
			},
		},
		"Success Empty": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient, ids []satypes.SubaccountId) {
				req := &api.LiquidateSubaccountsRequest{
					SubaccountIds: ids,
				}
				response := &api.LiquidateSubaccountsResponse{}
				mck.On("LiquidateSubaccounts", ctx, req).Return(response, nil)
			},
			subaccountIds: []satypes.SubaccountId{},
		},
		"Errors are propagated": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient, ids []satypes.SubaccountId) {
				req := &api.LiquidateSubaccountsRequest{
					SubaccountIds: ids,
				}
				mck.On("LiquidateSubaccounts", ctx, req).Return(nil, errors.New("test error"))
			},
			subaccountIds: []satypes.SubaccountId{},
			expectedError: errors.New("test error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			queryClientMock := &mocks.QueryClient{}
			tc.setupMocks(grpc.Ctx, queryClientMock, tc.subaccountIds)

			err := client.SendLiquidatableSubaccountIds(grpc.Ctx, queryClientMock, tc.subaccountIds)
			require.Equal(t, tc.expectedError, err)
		})
	}
}
