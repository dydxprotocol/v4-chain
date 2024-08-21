package client_test

import (
	"context"
	"math/big"
	"testing"

	"cosmossdk.io/log"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/client"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/grpc"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
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
					BlockHeight: uint32(50),
					LiquidatableSubaccountIds: []satypes.SubaccountId{
						constants.Carl_Num0,
					},
					NegativeTncSubaccountIds: []satypes.SubaccountId{},
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId:                 0,
							SubaccountsWithLongPosition: []satypes.SubaccountId{},
							SubaccountsWithShortPosition: []satypes.SubaccountId{
								constants.Carl_Num0,
							},
						},
					},
				}
				setupMockLiquidateSubaccountRequests(mck, ctx, req)
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
					BlockHeight: uint32(50),
					LiquidatableSubaccountIds: []satypes.SubaccountId{
						constants.Dave_Num0,
					},
					NegativeTncSubaccountIds: []satypes.SubaccountId{},
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId: 0,
							SubaccountsWithLongPosition: []satypes.SubaccountId{
								constants.Dave_Num0,
							},
							SubaccountsWithShortPosition: []satypes.SubaccountId{},
						},
					},
				}
				setupMockLiquidateSubaccountRequests(mck, ctx, req)
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
					BlockHeight:               uint32(50),
					LiquidatableSubaccountIds: []satypes.SubaccountId{},
					NegativeTncSubaccountIds:  []satypes.SubaccountId{},
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId: 0,
							SubaccountsWithLongPosition: []satypes.SubaccountId{
								constants.Dave_Num0,
							},
							SubaccountsWithShortPosition: []satypes.SubaccountId{
								constants.Carl_Num0,
							},
						},
					},
				}
				setupMockLiquidateSubaccountRequests(mck, ctx, req)
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
					BlockHeight:                uint32(50),
					LiquidatableSubaccountIds:  []satypes.SubaccountId{},
					NegativeTncSubaccountIds:   []satypes.SubaccountId{},
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{},
				}
				setupMockLiquidateSubaccountRequests(mck, ctx, req)
			},
		},
		"Can get subaccount that become undercollateralized with funding payments (short)": {
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
						// Without funding, Carl has a TNC of $5,000, MMR of $5,000, and is
						// well-collateralized.
						// However, funding index for Carl's position is 10,000 and perpetual's funding index
						// is 0. Index delta is -10,000, so Carl has to make a funding payment of $1 and
						// become under-collateralized.
						{
							Id: &constants.Carl_Num0,
							AssetPositions: []*satypes.AssetPosition{
								testutil.CreateSingleAssetPosition(
									0,
									big.NewInt(55_000_000_000), // $55,000
								),
							},
							PerpetualPositions: []*satypes.PerpetualPosition{
								testutil.CreateSinglePerpetualPosition(
									0,
									big.NewInt(-100_000_000), // -1 BTC
									big.NewInt(10_000),
									big.NewInt(0),
								),
							},
						},
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
					BlockHeight: uint32(50),
					LiquidatableSubaccountIds: []satypes.SubaccountId{
						constants.Carl_Num0,
					},
					NegativeTncSubaccountIds: []satypes.SubaccountId{},
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId:                 0,
							SubaccountsWithLongPosition: []satypes.SubaccountId{},
							SubaccountsWithShortPosition: []satypes.SubaccountId{
								constants.Carl_Num0,
							},
						},
					},
				}
				setupMockLiquidateSubaccountRequests(mck, ctx, req)
			},
		},
		"Can get subaccount that become liquidatable with funding payments (long)": {
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
						// Without funding, Dave has a TNC of $5,000, MMR of $5,000, and is
						// well-collateralized.
						// However, funding index for Dave's position is -10,000 and perpetual's funding index
						// is 0. Index delta is 10,000, so Dave has to make a funding payment of $1 and
						// become under-collateralized.
						{
							Id: &constants.Dave_Num0,
							AssetPositions: []*satypes.AssetPosition{
								testutil.CreateSingleAssetPosition(
									0,
									big.NewInt(-45_000_000_000), // -$45,000
								),
							},
							PerpetualPositions: []*satypes.PerpetualPosition{
								testutil.CreateSinglePerpetualPosition(
									0,
									big.NewInt(100_000_000), // 1 BTC
									big.NewInt(-10_000),
									big.NewInt(0),
								),
							},
						},
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
					BlockHeight: uint32(50),
					LiquidatableSubaccountIds: []satypes.SubaccountId{
						constants.Dave_Num0,
					},
					NegativeTncSubaccountIds: []satypes.SubaccountId{},
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId: 0,
							SubaccountsWithLongPosition: []satypes.SubaccountId{
								constants.Dave_Num0,
							},
							SubaccountsWithShortPosition: []satypes.SubaccountId{},
						},
					},
				}
				setupMockLiquidateSubaccountRequests(mck, ctx, req)
			},
		},
		"Skips subaccount that become well-collateralized with funding payments (short)": {
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
						// Without funding, Carl has a TNC of $4,999, MMR of $5,000, and is
						// under-collateralized.
						// However, funding index for Carl's position is -10,000 and perpetual's funding index
						// is 0. Index delta is 10,000, so Carl would receive a funding payment of $1 and
						// become well-collateralized.
						{
							Id: &constants.Carl_Num0,
							AssetPositions: []*satypes.AssetPosition{
								testutil.CreateSingleAssetPosition(
									0,
									big.NewInt(54_999_000_000), // $54,999
								),
							},
							PerpetualPositions: []*satypes.PerpetualPosition{
								testutil.CreateSinglePerpetualPosition(
									0,
									big.NewInt(-100_000_000), // -1 BTC
									big.NewInt(-10_000),
									big.NewInt(0),
								),
							},
						},
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
					BlockHeight:               uint32(50),
					LiquidatableSubaccountIds: []satypes.SubaccountId{},
					NegativeTncSubaccountIds:  []satypes.SubaccountId{},
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId:                 0,
							SubaccountsWithLongPosition: []satypes.SubaccountId{},
							SubaccountsWithShortPosition: []satypes.SubaccountId{
								constants.Carl_Num0,
							},
						},
					},
				}
				setupMockLiquidateSubaccountRequests(mck, ctx, req)
			},
		},
		"Skips subaccount that become well-collateralized with funding payments (long)": {
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
						// Without funding, Dave has a TNC of $4,999, MMR of $5,000, and is
						// under-collateralized.
						// However, funding index for Dave's position is 10,000 and perpetual's funding index
						// is 0. Index delta is -10,000, so Dave would receive a funding payment of $1 and
						// become well-collateralized.
						{
							Id: &constants.Dave_Num0,
							AssetPositions: []*satypes.AssetPosition{
								testutil.CreateSingleAssetPosition(
									0,
									big.NewInt(-44_999_000_000), // -$44,999
								),
							},
							PerpetualPositions: []*satypes.PerpetualPosition{
								testutil.CreateSinglePerpetualPosition(
									0,
									big.NewInt(100_000_000), // 1 BTC
									big.NewInt(10_000),
									big.NewInt(0),
								),
							},
						},
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
					BlockHeight:               uint32(50),
					LiquidatableSubaccountIds: []satypes.SubaccountId{},
					NegativeTncSubaccountIds:  []satypes.SubaccountId{},
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId: 0,
							SubaccountsWithLongPosition: []satypes.SubaccountId{
								constants.Dave_Num0,
							},
							SubaccountsWithShortPosition: []satypes.SubaccountId{},
						},
					},
				}
				setupMockLiquidateSubaccountRequests(mck, ctx, req)
			},
		},
		"Can get negative tnc subaccount with short position": {
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
						// Carl has TNC of -$1.
						constants.Carl_Num0_1BTC_Short_49999USD,
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
					BlockHeight: uint32(50),
					LiquidatableSubaccountIds: []satypes.SubaccountId{
						constants.Carl_Num0,
					},
					NegativeTncSubaccountIds: []satypes.SubaccountId{
						constants.Carl_Num0,
					},
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId:                 0,
							SubaccountsWithLongPosition: []satypes.SubaccountId{},
							SubaccountsWithShortPosition: []satypes.SubaccountId{
								constants.Carl_Num0,
							},
						},
					},
				}
				setupMockLiquidateSubaccountRequests(mck, ctx, req)
			},
		},
		"Can get negative tnc subaccount with long position": {
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
						// Dave has TNC of -$1.
						constants.Dave_Num0_1BTC_Long_50001USD_Short,
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
					BlockHeight: uint32(50),
					LiquidatableSubaccountIds: []satypes.SubaccountId{
						constants.Dave_Num0,
					},
					NegativeTncSubaccountIds: []satypes.SubaccountId{
						constants.Dave_Num0,
					},
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId: 0,
							SubaccountsWithLongPosition: []satypes.SubaccountId{
								constants.Dave_Num0,
							},
							SubaccountsWithShortPosition: []satypes.SubaccountId{},
						},
					},
				}
				setupMockLiquidateSubaccountRequests(mck, ctx, req)
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

func setupMockLiquidateSubaccountRequests(
	mck *mocks.QueryClient,
	ctx context.Context,
	request *api.LiquidateSubaccountsRequest,
) {
	req := &api.LiquidateSubaccountsRequest{
		BlockHeight:               request.BlockHeight,
		LiquidatableSubaccountIds: request.LiquidatableSubaccountIds,
	}
	response := &api.LiquidateSubaccountsResponse{}
	mck.On("LiquidateSubaccounts", ctx, req).Return(response, nil)

	req = &api.LiquidateSubaccountsRequest{
		BlockHeight:              request.BlockHeight,
		NegativeTncSubaccountIds: request.NegativeTncSubaccountIds,
	}
	mck.On("LiquidateSubaccounts", ctx, req).Return(response, nil)

	if len(request.SubaccountOpenPositionInfo) == 0 {
		req = &api.LiquidateSubaccountsRequest{
			BlockHeight:                request.BlockHeight,
			SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{},
		}
		mck.On("LiquidateSubaccounts", ctx, req).Return(response, nil)
	} else {
		for _, info := range request.SubaccountOpenPositionInfo {
			if len(info.SubaccountsWithLongPosition) > 0 {
				req = &api.LiquidateSubaccountsRequest{
					BlockHeight: request.BlockHeight,
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId:                 info.PerpetualId,
							SubaccountsWithLongPosition: info.SubaccountsWithLongPosition,
						},
					},
				}
				mck.On("LiquidateSubaccounts", ctx, req).Return(response, nil)
			}

			if len(info.SubaccountsWithShortPosition) > 0 {
				req = &api.LiquidateSubaccountsRequest{
					BlockHeight: request.BlockHeight,
					SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
						{
							PerpetualId:                  info.PerpetualId,
							SubaccountsWithShortPosition: info.SubaccountsWithShortPosition,
						},
					},
				}
				mck.On("LiquidateSubaccounts", ctx, req).Return(response, nil)
			}
		}
	}
}
