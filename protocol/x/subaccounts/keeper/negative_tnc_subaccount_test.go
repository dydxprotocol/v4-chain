package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/tracer"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetNegativeTncSubaccountSeenAtBlock(t *testing.T) {
	testCollateralPoolAddresses := []sdk.AccAddress{
		constants.IsoCollateralPoolAddress,
		constants.Iso2CollateralPoolAddress,
		types.ModuleAddress,
	}
	tests := map[string]struct {
		// Setup.
		setupTestAndPerformAssertions func(ctx sdk.Context, s keeper.Keeper)

		// Expectations.
		expectedMultiStoreWrites []string
	}{
		"Block height defaults to zero if not set and doesn't exist": {
			setupTestAndPerformAssertions: func(ctx sdk.Context, k keeper.Keeper) {
				for _, collateralPoolAddr := range testCollateralPoolAddresses {
					block, exists := k.GetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr)
					require.False(t, exists)
					require.Equal(
						t,
						uint32(0),
						block,
					)
				}
			},

			expectedMultiStoreWrites: []string{},
		},
		"Block height can be updated": {
			setupTestAndPerformAssertions: func(ctx sdk.Context, k keeper.Keeper) {
				for _, collateralPoolAddr := range testCollateralPoolAddresses {
					k.SetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr, 1)
					block, exists := k.GetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(1),
						block,
					)
				}
			},

			expectedMultiStoreWrites: []string{
				types.NegativeTncSubaccountForCollateralPoolSeenAtBlockKeyPrefix + string(constants.IsoCollateralPoolAddress),
				types.NegativeTncSubaccountForCollateralPoolSeenAtBlockKeyPrefix + string(constants.Iso2CollateralPoolAddress),
				types.NegativeTncSubaccountForCollateralPoolSeenAtBlockKeyPrefix + string(types.ModuleAddress),
			},
		},
		"Block height can be updated more than once": {
			setupTestAndPerformAssertions: func(ctx sdk.Context, k keeper.Keeper) {
				for _, collateralPoolAddr := range testCollateralPoolAddresses {
					k.SetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr, 1)
					block, exists := k.GetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(1),
						block,
					)

					k.SetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr, 2)
					block, exists = k.GetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(2),
						block,
					)

					k.SetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr, 3)
					block, exists = k.GetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(3),
						block,
					)

					k.SetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr, 10)
					block, exists = k.GetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(10),
						block,
					)
				}
			},

			expectedMultiStoreWrites: append(
				getWriteKeys(constants.IsoCollateralPoolAddress, 4),
				append(
					getWriteKeys(constants.Iso2CollateralPoolAddress, 4),
					getWriteKeys(types.ModuleAddress, 4)...,
				)...,
			),
		},
		"Block height can be updated to same block height": {
			setupTestAndPerformAssertions: func(ctx sdk.Context, k keeper.Keeper) {
				for _, collateralPoolAddr := range testCollateralPoolAddresses {
					k.SetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr, 0)
					block, exists := k.GetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(0),
						block,
					)

					k.SetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr, 0)
					block, exists = k.GetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(0),
						block,
					)
				}
			},

			expectedMultiStoreWrites: append(
				getWriteKeys(constants.IsoCollateralPoolAddress, 2),
				append(
					getWriteKeys(constants.Iso2CollateralPoolAddress, 2),
					getWriteKeys(types.ModuleAddress, 2)...,
				)...,
			),
		},
		"Block height can be updated to different block heights for each collateral pool address": {
			setupTestAndPerformAssertions: func(ctx sdk.Context, k keeper.Keeper) {
				for i, collateralPoolAddr := range testCollateralPoolAddresses {
					k.SetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr, uint32(i))
					block, exists := k.GetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(i),
						block,
					)

					k.SetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr, uint32(2*i+1))
					block, exists = k.GetNegativeTncSubaccountSeenAtBlock(ctx, collateralPoolAddr)
					require.True(t, exists)
					require.Equal(
						t,
						uint32(2*i+1),
						block,
					)
				}
			},

			expectedMultiStoreWrites: append(
				getWriteKeys(constants.IsoCollateralPoolAddress, 2),
				append(
					getWriteKeys(constants.Iso2CollateralPoolAddress, 2),
					getWriteKeys(types.ModuleAddress, 2)...,
				)...,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state and test parameters.
			ctx, subaccountsKeeper, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, true)

			// Set the tracer on the multistore to verify the performed writes are correct.
			traceDecoder := &tracer.TraceDecoder{}
			ctx.MultiStore().SetTracer(traceDecoder)

			tc.setupTestAndPerformAssertions(
				ctx,
				*subaccountsKeeper,
			)

			// Verify the writes were done in the expected order.
			traceDecoder.RequireKeyPrefixWrittenInSequence(t, tc.expectedMultiStoreWrites)
		})
	}
}

func TestGetSetNegativeTncSubaccountSeenAtBlock_PanicsOnDecreasingBlock(t *testing.T) {
	// Setup keeper state and test parameters.
	ctx, subaccountsKeeper, _, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, true)

	subaccountsKeeper.SetNegativeTncSubaccountSeenAtBlock(ctx, types.ModuleAddress, 2)
	require.PanicsWithValue(
		t,
		"SetNegativeTncSubaccountSeenAtBlock: new block height (1) is less than the current block height (2)",
		func() { subaccountsKeeper.SetNegativeTncSubaccountSeenAtBlock(ctx, types.ModuleAddress, 1) },
	)
}

func getWriteKeys(address sdk.AccAddress, times int) []string {
	writeKeys := make([]string, times)
	for i := 0; i < times; i++ {
		writeKeys[i] = types.NegativeTncSubaccountForCollateralPoolSeenAtBlockKeyPrefix + string(address)
	}
	return writeKeys
}
