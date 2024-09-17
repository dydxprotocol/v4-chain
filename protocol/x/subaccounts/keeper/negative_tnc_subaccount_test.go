package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/tracer"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetNegativeTncSubaccountSeenAtBlock(t *testing.T) {
	testPerpetualIds := []uint32{
		constants.IsoUsd_IsolatedMarket.Params.Id,
		constants.Iso2Usd_IsolatedMarket.Params.Id,
		constants.BtcUsd_NoMarginRequirement.Params.Id,
	}
	testsuffixes := []string{
		lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
		lib.UintToString(constants.Iso2Usd_IsolatedMarket.Params.Id),
		types.CrossCollateralSuffix,
	}
	tests := map[string]struct {
		// Setup.
		setupTestAndPerformAssertions func(ctx sdk.Context, s keeper.Keeper) error

		// Expectations.
		expectedMultiStoreWrites []string
	}{
		"Block height defaults to zero if not set and doesn't exist": {
			setupTestAndPerformAssertions: func(ctx sdk.Context, k keeper.Keeper) error {
				for _, perpetualId := range testPerpetualIds {
					block, exists, err := k.GetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId)
					require.NoError(t, err)
					require.False(t, exists)
					require.Equal(
						t,
						uint32(0),
						block,
					)
				}
				return nil
			},

			expectedMultiStoreWrites: []string{},
		},
		"Block height can be updated": {
			setupTestAndPerformAssertions: func(ctx sdk.Context, k keeper.Keeper) error {
				for _, perpetualId := range testPerpetualIds {
					err := k.SetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId, 1)
					if err != nil {
						return err
					}
					block, exists, err := k.GetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId)
					require.NoError(t, err)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(1),
						block,
					)
				}
				return nil
			},

			expectedMultiStoreWrites: []string{
				types.NegativeTncSubaccountForCollateralPoolSeenAtBlockKeyPrefix + testsuffixes[0],
				types.NegativeTncSubaccountForCollateralPoolSeenAtBlockKeyPrefix + testsuffixes[1],
				types.NegativeTncSubaccountForCollateralPoolSeenAtBlockKeyPrefix + testsuffixes[2],
			},
		},
		"Block height can be updated more than once": {
			setupTestAndPerformAssertions: func(ctx sdk.Context, k keeper.Keeper) error {
				for _, perpetualId := range testPerpetualIds {
					err := k.SetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId, 1)
					if err != nil {
						return nil
					}
					block, exists, err := k.GetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId)
					require.NoError(t, err)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(1),
						block,
					)

					err = k.SetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId, 2)
					if err != nil {
						return nil
					}
					block, exists, err = k.GetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId)
					require.NoError(t, err)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(2),
						block,
					)

					err = k.SetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId, 3)
					if err != nil {
						return nil
					}
					block, exists, err = k.GetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId)
					require.NoError(t, err)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(3),
						block,
					)

					err = k.SetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId, 10)
					if err != nil {
						return nil
					}
					block, exists, err = k.GetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId)
					require.NoError(t, err)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(10),
						block,
					)
				}
				return nil
			},

			expectedMultiStoreWrites: append(
				getWriteKeys(testsuffixes[0], 4),
				append(
					getWriteKeys(testsuffixes[1], 4),
					getWriteKeys(testsuffixes[2], 4)...,
				)...,
			),
		},
		"Block height can be updated to same block height": {
			setupTestAndPerformAssertions: func(ctx sdk.Context, k keeper.Keeper) error {
				for _, perpetualId := range testPerpetualIds {
					err := k.SetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId, 0)
					if err != nil {
						return err
					}
					block, exists, err := k.GetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId)
					require.NoError(t, err)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(0),
						block,
					)

					err = k.SetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId, 0)
					if err != nil {
						return err
					}
					block, exists, err = k.GetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId)
					require.NoError(t, err)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(0),
						block,
					)
				}
				return nil
			},

			expectedMultiStoreWrites: append(
				getWriteKeys(testsuffixes[0], 2),
				append(
					getWriteKeys(testsuffixes[1], 2),
					getWriteKeys(testsuffixes[2], 2)...,
				)...,
			),
		},
		"Block height can be updated to different block heights for each collateral pool address": {
			setupTestAndPerformAssertions: func(ctx sdk.Context, k keeper.Keeper) error {
				for i, perpetualId := range testPerpetualIds {
					err := k.SetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId, uint32(i))
					if err != nil {
						return err
					}
					block, exists, err := k.GetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId)
					require.NoError(t, err)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(i),
						block,
					)

					err = k.SetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId, uint32(2*i+1))
					if err != nil {
						return err
					}
					block, exists, err = k.GetNegativeTncSubaccountSeenAtBlock(ctx, perpetualId)
					require.NoError(t, err)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(2*i+1),
						block,
					)
				}
				return nil
			},

			expectedMultiStoreWrites: append(
				getWriteKeys(testsuffixes[0], 2),
				append(
					getWriteKeys(testsuffixes[1], 2),
					getWriteKeys(testsuffixes[2], 2)...,
				)...,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state and test parameters.
			ctx, subaccountsKeeper, pricesKeeper, perpetualsKeeper, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(
				t,
				true,
			)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)
			keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)
			keepertest.CreateTestPerpetuals(t, ctx, perpetualsKeeper)

			// Set the tracer on the multistore to verify the performed writes are correct.
			traceDecoder := &tracer.TraceDecoder{}
			ctx.MultiStore().SetTracer(traceDecoder)

			err := tc.setupTestAndPerformAssertions(
				ctx,
				*subaccountsKeeper,
			)
			require.NoError(t, err)

			// Verify the writes were done in the expected order.
			traceDecoder.RequireKeyPrefixWrittenInSequence(t, tc.expectedMultiStoreWrites)
		})
	}
}

func TestGetSetNegativeTncSubaccountSeenAtBlock_PanicsOnDecreasingBlock(t *testing.T) {
	// Setup keeper state and test parameters.
	ctx, subaccountsKeeper, pricesKeeper, perpetualsKeeper, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, true)
	keepertest.CreateTestMarkets(t, ctx, pricesKeeper)
	keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)
	keepertest.CreateTestPerpetuals(t, ctx, perpetualsKeeper)
	err := subaccountsKeeper.SetNegativeTncSubaccountSeenAtBlock(ctx, uint32(0), 2)
	require.NoError(t, err)
	require.PanicsWithValue(
		t,
		"SetNegativeTncSubaccountSeenAtBlock: new block height (1) is less than the current block height (2)",
		func() { _ = subaccountsKeeper.SetNegativeTncSubaccountSeenAtBlock(ctx, uint32(0), 1) },
	)
}

func getWriteKeys(suffix string, times int) []string {
	writeKeys := make([]string, times)
	for i := 0; i < times; i++ {
		writeKeys[i] = types.NegativeTncSubaccountForCollateralPoolSeenAtBlockKeyPrefix + suffix
	}
	return writeKeys
}
