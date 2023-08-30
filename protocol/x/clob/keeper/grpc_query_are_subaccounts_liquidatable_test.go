package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

func TestAreSubaccountsLiquidatable(t *testing.T) {
	for _, tc := range []struct {
		desc        string
		subaccounts []satypes.Subaccount
		perpetuals  []perptypes.Perpetual
		request     *types.AreSubaccountsLiquidatableRequest
		response    *types.AreSubaccountsLiquidatableResponse
		err         error
	}{
		{
			desc: "No errors",
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_10_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(-1_000_000_000), // 1 BTC
						},
					},
				},
				{
					Id: &constants.Bob_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_10_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
						},
					},
				},
			},
			request: &types.AreSubaccountsLiquidatableRequest{
				SubaccountIds: []satypes.SubaccountId{
					constants.Alice_Num0,
					constants.Bob_Num0,
				},
			},
			response: &types.AreSubaccountsLiquidatableResponse{
				Results: []types.AreSubaccountsLiquidatableResponse_Result{
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
		},
		{
			desc: "Non-existent subaccount",
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			subaccounts: []satypes.Subaccount{},
			request: &types.AreSubaccountsLiquidatableRequest{
				SubaccountIds: []satypes.SubaccountId{
					constants.Alice_Num0,
				},
			},
			response: &types.AreSubaccountsLiquidatableResponse{
				Results: []types.AreSubaccountsLiquidatableResponse_Result{
					{
						SubaccountId:   constants.Alice_Num0,
						IsLiquidatable: false,
					},
				},
			},
		},
		{
			desc: "Errors are propagated",
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Alice_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_10_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(-1_000_000_000), // 1 BTC
						},
					},
				},
			},
			perpetuals: []perptypes.Perpetual{},
			request: &types.AreSubaccountsLiquidatableRequest{
				SubaccountIds: []satypes.SubaccountId{
					constants.Alice_Num0,
				},
			},
			err: perptypes.ErrPerpetualDoesNotExist,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ks.Ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
				)
				require.NoError(t, err)
			}

			for _, subaccount := range tc.subaccounts {
				ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, subaccount)
			}

			wctx := sdk.WrapSDKContext(ks.Ctx)
			response, err := ks.ClobKeeper.AreSubaccountsLiquidatable(wctx, tc.request)

			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response), //nolint:staticcheck
					nullify.Fill(response),    //nolint:staticcheck
				)
			}
		})
	}
}
