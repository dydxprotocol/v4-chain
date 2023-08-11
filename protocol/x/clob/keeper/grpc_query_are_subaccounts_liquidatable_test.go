package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4/dtypes"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/testutil/nullify"
	"github.com/dydxprotocol/v4/x/clob/memclob"
	"github.com/dydxprotocol/v4/x/clob/types"
	perptypes "github.com/dydxprotocol/v4/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
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
			ctx,
				clobKeeper,
				pricesKeeper,
				_,
				perpetualsKeeper,
				subaccountsKeeper,
				_,
				_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			// Create the default markets.
			keepertest.CreateTestMarketsAndExchangeFeeds(t, ctx, pricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			// Create all perpetuals.
			for _, p := range tc.perpetuals {
				_, err := perpetualsKeeper.CreatePerpetual(
					ctx,
					p.Ticker,
					p.MarketId,
					p.AtomicResolution,
					p.DefaultFundingPpm,
					p.LiquidityTier,
				)
				require.NoError(t, err)
			}

			for _, subaccount := range tc.subaccounts {
				subaccountsKeeper.SetSubaccount(ctx, subaccount)
			}

			wctx := sdk.WrapSDKContext(ctx)
			response, err := clobKeeper.AreSubaccountsLiquidatable(wctx, tc.request)

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
