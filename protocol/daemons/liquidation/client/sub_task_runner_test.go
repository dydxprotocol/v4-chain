package client_test

import (
	"context"
	"errors"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/client"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/grpc"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	daemonInitializingErrorString = "no successful update has occurred"
)

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

			daemonClient := client.NewClient(log.NewNopLogger())
			actual, err := client.GetAllSubaccounts(
				daemonClient,
				grpc.Ctx,
				queryClientMock,
				df.Liquidation.SubaccountPageLimit,
			)
			if err != nil {
				require.EqualError(t, err, tc.expectedError.Error())
				// The daemon initializes as unhealthy.
				// If a request fails, the daemon will not be toggled to healthy.
				require.ErrorContains(t, daemonClient.HealthCheck(), daemonInitializingErrorString)
			} else {
				require.Equal(t, tc.expectedSubaccounts, actual)
				// If the request(s) succeeded, expect a healthy daemon.
				require.NoError(t, daemonClient.HealthCheck())
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

			daemon := client.NewClient(log.NewNopLogger())
			actual, err := client.CheckCollateralizationForSubaccounts(
				daemon,
				grpc.Ctx,
				queryClientMock,
				tc.subaccountIds,
			)

			if err != nil {
				require.EqualError(t, err, tc.expectedError.Error())
				// The daemon initializes as unhealthy.
				// If a request fails, the daemon will not be toggled to healthy.
				require.ErrorContains(t, daemon.HealthCheck(), daemonInitializingErrorString)
			} else {
				require.Equal(t, tc.expectedResults, actual)
				// If the request(s) succeeded, expect a healthy daemon.
				require.NoError(t, daemon.HealthCheck())
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
