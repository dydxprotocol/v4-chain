package client_test

import (
	"context"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/client"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/grpc"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRunLiquidationDaemonTaskLoop(t *testing.T) {
	tests := map[string]struct {
		// mocks
		setupMocks func(ctx context.Context, mck *mocks.QueryClient)

		// expectations
		expectedLiquidatableSubaccountIds []satypes.SubaccountId
		expectedError                     error
	}{
		"Can get liquidatable subaccount with short position": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				// Block height.
				res := &blocktimetypes.QueryPreviousBlockInfoResponse{
					Info: &blocktimetypes.BlockInfo{
						Height:    uint32(50),
						Timestamp: constants.TimeTen,
					},
				}
				mck.On("PreviousBlockInfo", mock.Anything, mock.Anything).Return(res, nil)

				// Subaccount.
				res2 := &satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Carl_Num0_1BTC_Short_54999USD,
					},
				}
				mck.On("SubaccountAll", mock.Anything, mock.Anything).Return(res2, nil)

				// Market prices.
				res3 := &pricestypes.QueryAllMarketPricesResponse{
					MarketPrices: constants.TestMarketPrices,
				}
				mck.On("AllMarketPrices", mock.Anything, mock.Anything).Return(res3, nil)

				// Perpetuals.
				res4 := &perptypes.QueryAllPerpetualsResponse{
					Perpetual: []perptypes.Perpetual{
						constants.BtcUsd_20PercentInitial_10PercentMaintenance,
					},
				}
				mck.On("AllPerpetuals", mock.Anything, mock.Anything).Return(res4, nil)

				// Liquidity tiers.
				res5 := &perptypes.QueryAllLiquidityTiersResponse{
					LiquidityTiers: constants.LiquidityTiers,
				}
				mck.On("AllLiquidityTiers", mock.Anything, mock.Anything).Return(res5, nil)

				// Sends liquidatable subaccount ids to the server.
				req := &api.LiquidateSubaccountsRequest{
					SubaccountIds: []satypes.SubaccountId{
						constants.Carl_Num0,
					},
				}
				response3 := &api.LiquidateSubaccountsResponse{}
				mck.On("LiquidateSubaccounts", ctx, req).Return(response3, nil)
			},
		},
		"Can get liquidatable subaccount with long position": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				// Block height.
				res := &blocktimetypes.QueryPreviousBlockInfoResponse{
					Info: &blocktimetypes.BlockInfo{
						Height:    uint32(50),
						Timestamp: constants.TimeTen,
					},
				}
				mck.On("PreviousBlockInfo", mock.Anything, mock.Anything).Return(res, nil)

				// Subaccount.
				res2 := &satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Dave_Num0_1BTC_Long_45001USD_Short,
					},
				}
				mck.On("SubaccountAll", mock.Anything, mock.Anything).Return(res2, nil)

				// Market prices.
				res3 := &pricestypes.QueryAllMarketPricesResponse{
					MarketPrices: constants.TestMarketPrices,
				}
				mck.On("AllMarketPrices", mock.Anything, mock.Anything).Return(res3, nil)

				// Perpetuals.
				res4 := &perptypes.QueryAllPerpetualsResponse{
					Perpetual: []perptypes.Perpetual{
						constants.BtcUsd_20PercentInitial_10PercentMaintenance,
					},
				}
				mck.On("AllPerpetuals", mock.Anything, mock.Anything).Return(res4, nil)

				// Liquidity tiers.
				res5 := &perptypes.QueryAllLiquidityTiersResponse{
					LiquidityTiers: constants.LiquidityTiers,
				}
				mck.On("AllLiquidityTiers", mock.Anything, mock.Anything).Return(res5, nil)

				// Sends liquidatable subaccount ids to the server.
				req := &api.LiquidateSubaccountsRequest{
					SubaccountIds: []satypes.SubaccountId{
						constants.Dave_Num0,
					},
				}
				response3 := &api.LiquidateSubaccountsResponse{}
				mck.On("LiquidateSubaccounts", ctx, req).Return(response3, nil)
			},
		},
		"Skip well collateralized subaccounts": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				// Block height.
				res := &blocktimetypes.QueryPreviousBlockInfoResponse{
					Info: &blocktimetypes.BlockInfo{
						Height:    uint32(50),
						Timestamp: constants.TimeTen,
					},
				}
				mck.On("PreviousBlockInfo", mock.Anything, mock.Anything).Return(res, nil)

				// Subaccount.
				res2 := &satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Carl_Num0_1BTC_Short_55000USD,
						constants.Dave_Num0_1BTC_Long_45000USD_Short,
					},
				}
				mck.On("SubaccountAll", mock.Anything, mock.Anything).Return(res2, nil)

				// Market prices.
				res3 := &pricestypes.QueryAllMarketPricesResponse{
					MarketPrices: constants.TestMarketPrices,
				}
				mck.On("AllMarketPrices", mock.Anything, mock.Anything).Return(res3, nil)

				// Perpetuals.
				res4 := &perptypes.QueryAllPerpetualsResponse{
					Perpetual: []perptypes.Perpetual{
						constants.BtcUsd_20PercentInitial_10PercentMaintenance,
					},
				}
				mck.On("AllPerpetuals", mock.Anything, mock.Anything).Return(res4, nil)

				// Liquidity tiers.
				res5 := &perptypes.QueryAllLiquidityTiersResponse{
					LiquidityTiers: constants.LiquidityTiers,
				}
				mck.On("AllLiquidityTiers", mock.Anything, mock.Anything).Return(res5, nil)

				// Sends liquidatable subaccount ids to the server.
				req := &api.LiquidateSubaccountsRequest{
					SubaccountIds: []satypes.SubaccountId{},
				}
				response3 := &api.LiquidateSubaccountsResponse{}
				mck.On("LiquidateSubaccounts", ctx, req).Return(response3, nil)
			},
		},
		"Skip subaccounts with no open positions": {
			setupMocks: func(ctx context.Context, mck *mocks.QueryClient) {
				// Block height.
				res := &blocktimetypes.QueryPreviousBlockInfoResponse{
					Info: &blocktimetypes.BlockInfo{
						Height:    uint32(50),
						Timestamp: constants.TimeTen,
					},
				}
				mck.On("PreviousBlockInfo", mock.Anything, mock.Anything).Return(res, nil)

				// Subaccount.
				res2 := &satypes.QuerySubaccountAllResponse{
					Subaccount: []satypes.Subaccount{
						constants.Alice_Num0_100_000USD,
					},
				}
				mck.On("SubaccountAll", mock.Anything, mock.Anything).Return(res2, nil)

				// Market prices.
				res3 := &pricestypes.QueryAllMarketPricesResponse{
					MarketPrices: constants.TestMarketPrices,
				}
				mck.On("AllMarketPrices", mock.Anything, mock.Anything).Return(res3, nil)

				// Perpetuals.
				res4 := &perptypes.QueryAllPerpetualsResponse{
					Perpetual: []perptypes.Perpetual{
						constants.BtcUsd_20PercentInitial_10PercentMaintenance,
					},
				}
				mck.On("AllPerpetuals", mock.Anything, mock.Anything).Return(res4, nil)

				// Liquidity tiers.
				res5 := &perptypes.QueryAllLiquidityTiersResponse{
					LiquidityTiers: constants.LiquidityTiers,
				}
				mck.On("AllLiquidityTiers", mock.Anything, mock.Anything).Return(res5, nil)

				// Sends liquidatable subaccount ids to the server.
				req := &api.LiquidateSubaccountsRequest{
					SubaccountIds: []satypes.SubaccountId{},
				}
				response3 := &api.LiquidateSubaccountsResponse{}
				mck.On("LiquidateSubaccounts", ctx, req).Return(response3, nil)
			},
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
			c.PerpetualsQueryClient = queryClientMock
			c.PricesQueryClient = queryClientMock
			c.BlocktimeQueryClient = queryClientMock

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
