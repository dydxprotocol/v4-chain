package client_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/client"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/grpc"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetPreviousBlockInfo(t *testing.T) {
	tests := map[string]struct {
		// mocks
		setupMocks func(
			ctx context.Context,
			mck *mocks.QueryClient,
		)

		// expectations
		expectedBlockHeight uint32
		expectedError       error
	}{
		"Success": {
			setupMocks: func(
				ctx context.Context,
				mck *mocks.QueryClient,
			) {
				response := &blocktimetypes.QueryPreviousBlockInfoResponse{
					Info: &blocktimetypes.BlockInfo{
						Height:    uint32(50),
						Timestamp: constants.TimeTen,
					},
				}
				mck.On("PreviousBlockInfo", ctx, mock.Anything).Return(response, nil)
			},
			expectedBlockHeight: 50,
		},
		"Errors are propagated": {
			setupMocks: func(
				ctx context.Context,
				mck *mocks.QueryClient,
			) {
				mck.On("PreviousBlockInfo", ctx, mock.Anything).Return(nil, errors.New("test error"))
			},
			expectedBlockHeight: 0,
			expectedError:       errors.New("test error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			queryClientMock := &mocks.QueryClient{}
			tc.setupMocks(grpc.Ctx, queryClientMock)

			daemon := client.NewClient(log.NewNopLogger())
			daemon.BlocktimeQueryClient = queryClientMock
			actualBlockHeight, err := daemon.GetPreviousBlockInfo(grpc.Ctx)

			if err != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.Equal(t, tc.expectedBlockHeight, actualBlockHeight)
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

			daemon := client.NewClient(log.NewNopLogger())
			daemon.SubaccountQueryClient = queryClientMock
			actual, err := daemon.GetAllSubaccounts(
				grpc.Ctx,
				df.Liquidation.SubaccountPageLimit,
			)
			if err != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.Equal(t, tc.expectedSubaccounts, actual)
			}
		})
	}
}

func TestGetAllPerpetuals(t *testing.T) {
	tests := map[string]struct {
		// mocks
		setupMocks func(ctx context.Context, mck *mocks.QueryClient)
		limit      uint64

		// expectations
		expectedPerpetuals []perptypes.Perpetual
		expectedError      error
	}{
		"Success": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &perptypes.QueryAllPerpetualsRequest{
					Pagination: &query.PageRequest{
						Limit: 1_000,
					},
				}
				response := &perptypes.QueryAllPerpetualsResponse{
					Perpetual: constants.Perpetuals_DefaultGenesisState.Perpetuals,
				}
				mck.On("AllPerpetuals", mock.Anything, req).Return(response, nil)
			},
			limit:              1_000,
			expectedPerpetuals: constants.Perpetuals_DefaultGenesisState.Perpetuals,
		},
		"Success Paginated": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &perptypes.QueryAllPerpetualsRequest{
					Pagination: &query.PageRequest{
						Limit: 1,
					},
				}
				nextKey := []byte("next key")
				response := &perptypes.QueryAllPerpetualsResponse{
					Perpetual: []perptypes.Perpetual{
						constants.Perpetuals_DefaultGenesisState.Perpetuals[0],
					},
					Pagination: &query.PageResponse{
						NextKey: nextKey,
					},
				}
				mck.On("AllPerpetuals", mock.Anything, req).Return(response, nil)
				req2 := &perptypes.QueryAllPerpetualsRequest{
					Pagination: &query.PageRequest{
						Key:   nextKey,
						Limit: 1,
					},
				}
				response2 := &perptypes.QueryAllPerpetualsResponse{
					Perpetual: []perptypes.Perpetual{
						constants.Perpetuals_DefaultGenesisState.Perpetuals[1],
					},
				}
				mck.On("AllPerpetuals", mock.Anything, req2).Return(response2, nil)
			},
			limit:              1,
			expectedPerpetuals: constants.Perpetuals_DefaultGenesisState.Perpetuals,
		},
		"Errors are propagated": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &perptypes.QueryAllPerpetualsRequest{
					Pagination: &query.PageRequest{
						Limit: 1_000,
					},
				}
				mck.On("AllPerpetuals", mock.Anything, req).Return(nil, errors.New("test error"))
			},
			limit:         1_000,
			expectedError: errors.New("test error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			queryClientMock := &mocks.QueryClient{}
			tc.setupMocks(grpc.Ctx, queryClientMock)

			daemon := client.NewClient(log.NewNopLogger())
			daemon.PerpetualsQueryClient = queryClientMock
			actual, err := daemon.GetAllPerpetuals(
				grpc.Ctx,
				uint32(50),
				tc.limit,
			)
			if err != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.Equal(t, tc.expectedPerpetuals, actual)
			}
		})
	}
}

func TestGetAllLiquidityTiers(t *testing.T) {
	tests := map[string]struct {
		// mocks
		setupMocks func(ctx context.Context, mck *mocks.QueryClient)
		limit      uint64

		// expectations
		expectedLiquidityTiers []perptypes.LiquidityTier
		expectedError          error
	}{
		"Success": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &perptypes.QueryAllLiquidityTiersRequest{
					Pagination: &query.PageRequest{
						Limit: 1_000,
					},
				}
				response := &perptypes.QueryAllLiquidityTiersResponse{
					LiquidityTiers: constants.LiquidityTiers,
				}
				mck.On("AllLiquidityTiers", mock.Anything, req).Return(response, nil)
			},
			limit:                  1_000,
			expectedLiquidityTiers: constants.LiquidityTiers,
		},
		"Success Paginated": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &perptypes.QueryAllLiquidityTiersRequest{
					Pagination: &query.PageRequest{
						Limit: 5,
					},
				}
				nextKey := []byte("next key")
				response := &perptypes.QueryAllLiquidityTiersResponse{
					LiquidityTiers: constants.LiquidityTiers[0:5],
					Pagination: &query.PageResponse{
						NextKey: nextKey,
					},
				}
				mck.On("AllLiquidityTiers", mock.Anything, req).Return(response, nil)
				req2 := &perptypes.QueryAllLiquidityTiersRequest{
					Pagination: &query.PageRequest{
						Key:   nextKey,
						Limit: 5,
					},
				}
				response2 := &perptypes.QueryAllLiquidityTiersResponse{
					LiquidityTiers: constants.LiquidityTiers[5:],
				}
				mck.On("AllLiquidityTiers", mock.Anything, req2).Return(response2, nil)
			},
			limit:                  5,
			expectedLiquidityTiers: constants.LiquidityTiers,
		},
		"Errors are propagated": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &perptypes.QueryAllLiquidityTiersRequest{
					Pagination: &query.PageRequest{
						Limit: 1_000,
					},
				}
				mck.On("AllLiquidityTiers", mock.Anything, req).Return(nil, errors.New("test error"))
			},
			limit:         1_000,
			expectedError: errors.New("test error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			queryClientMock := &mocks.QueryClient{}
			tc.setupMocks(grpc.Ctx, queryClientMock)

			daemon := client.NewClient(log.NewNopLogger())
			daemon.PerpetualsQueryClient = queryClientMock
			actual, err := daemon.GetAllLiquidityTiers(
				grpc.Ctx,
				uint32(50),
				tc.limit,
			)
			if err != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.Equal(t, tc.expectedLiquidityTiers, actual)
			}
		})
	}
}

func TestGetAllMarketPrices(t *testing.T) {
	tests := map[string]struct {
		// mocks
		setupMocks func(ctx context.Context, mck *mocks.QueryClient)
		limit      uint64

		// expectations
		expectedMarketPrices []pricestypes.MarketPrice
		expectedError        error
	}{
		"Success": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &pricestypes.QueryAllMarketPricesRequest{
					Pagination: &query.PageRequest{
						Limit: 1_000,
					},
				}
				response := &pricestypes.QueryAllMarketPricesResponse{
					MarketPrices: constants.TestMarketPrices,
				}
				mck.On("AllMarketPrices", mock.Anything, req).Return(response, nil)
			},
			limit:                1_000,
			expectedMarketPrices: constants.TestMarketPrices,
		},
		"Success Paginated": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &pricestypes.QueryAllMarketPricesRequest{
					Pagination: &query.PageRequest{
						Limit: 2,
					},
				}
				nextKey := []byte("next key")
				response := &pricestypes.QueryAllMarketPricesResponse{
					MarketPrices: []pricestypes.MarketPrice{
						constants.TestMarketPrices[0],
						constants.TestMarketPrices[1],
					},
					Pagination: &query.PageResponse{
						NextKey: nextKey,
					},
				}
				mck.On("AllMarketPrices", mock.Anything, req).Return(response, nil)
				req2 := &pricestypes.QueryAllMarketPricesRequest{
					Pagination: &query.PageRequest{
						Key:   nextKey,
						Limit: 2,
					},
				}
				response2 := &pricestypes.QueryAllMarketPricesResponse{
					MarketPrices: []pricestypes.MarketPrice{
						constants.TestMarketPrices[2],
					},
				}
				mck.On("AllMarketPrices", mock.Anything, req2).Return(response2, nil)
			},
			limit:                2,
			expectedMarketPrices: constants.TestMarketPrices,
		},
		"Errors are propagated": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				req := &pricestypes.QueryAllMarketPricesRequest{
					Pagination: &query.PageRequest{
						Limit: 1_000,
					},
				}
				mck.On("AllMarketPrices", mock.Anything, req).Return(nil, errors.New("test error"))
			},
			limit:         1_000,
			expectedError: errors.New("test error"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			queryClientMock := &mocks.QueryClient{}
			tc.setupMocks(grpc.Ctx, queryClientMock)

			daemon := client.NewClient(log.NewNopLogger())
			daemon.PricesQueryClient = queryClientMock
			actual, err := daemon.GetAllMarketPrices(
				grpc.Ctx,
				uint32(50),
				tc.limit,
			)
			if err != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.Equal(t, tc.expectedMarketPrices, actual)
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
			daemon.ClobQueryClient = queryClientMock
			actual, err := daemon.CheckCollateralizationForSubaccounts(
				grpc.Ctx,
				tc.subaccountIds,
			)

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
					LiquidatableSubaccountIds: ids,
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
					LiquidatableSubaccountIds: ids,
				}
				response := &api.LiquidateSubaccountsResponse{}
				mck.On("LiquidateSubaccounts", ctx, req).Return(response, nil)
			},
			subaccountIds: []satypes.SubaccountId{},
		},
		"Errors are propagated": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient, ids []satypes.SubaccountId) {
				req := &api.LiquidateSubaccountsRequest{
					LiquidatableSubaccountIds: ids,
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

			daemon := client.NewClient(log.NewNopLogger())
			daemon.LiquidationServiceClient = queryClientMock

			err := daemon.SendLiquidatableSubaccountIds(grpc.Ctx, tc.subaccountIds)
			require.Equal(t, tc.expectedError, err)
		})
	}
}
