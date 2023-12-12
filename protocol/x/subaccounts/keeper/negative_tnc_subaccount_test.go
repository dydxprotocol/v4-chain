package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/tracer"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetNegativeTncSubaccountSeenAtBlock(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		setupTestAndPerformAssertions func(ctx sdk.Context, s keeper.Keeper)

		// Expectations.
		expectedMultiStoreWrites []string
	}{
		"Block height defaults to zero if not set": {
			setupTestAndPerformAssertions: func(ctx sdk.Context, k keeper.Keeper) {
				require.Equal(
					t,
					uint32(0),
					k.GetNegativeTncSubaccountSeenAtBlock(ctx),
				)
			},

			expectedMultiStoreWrites: []string{},
		},
		"Block height can be updated": {
			setupTestAndPerformAssertions: func(ctx sdk.Context, k keeper.Keeper) {
				k.SetNegativeTncSubaccountSeenAtBlock(ctx, 1)
				require.Equal(
					t,
					uint32(1),
					k.GetNegativeTncSubaccountSeenAtBlock(ctx),
				)
			},

			expectedMultiStoreWrites: []string{
				types.NegativeTncSubaccountSeenAtBlockKey,
			},
		},
		"Block height can be updated more than once": {
			setupTestAndPerformAssertions: func(ctx sdk.Context, k keeper.Keeper) {
				k.SetNegativeTncSubaccountSeenAtBlock(ctx, 1)
				require.Equal(
					t,
					uint32(1),
					k.GetNegativeTncSubaccountSeenAtBlock(ctx),
				)

				k.SetNegativeTncSubaccountSeenAtBlock(ctx, 2)
				require.Equal(
					t,
					uint32(2),
					k.GetNegativeTncSubaccountSeenAtBlock(ctx),
				)

				k.SetNegativeTncSubaccountSeenAtBlock(ctx, 3)
				require.Equal(
					t,
					uint32(3),
					k.GetNegativeTncSubaccountSeenAtBlock(ctx),
				)

				k.SetNegativeTncSubaccountSeenAtBlock(ctx, 10)
				require.Equal(
					t,
					uint32(10),
					k.GetNegativeTncSubaccountSeenAtBlock(ctx),
				)
			},

			expectedMultiStoreWrites: []string{
				types.NegativeTncSubaccountSeenAtBlockKey,
				types.NegativeTncSubaccountSeenAtBlockKey,
				types.NegativeTncSubaccountSeenAtBlockKey,
				types.NegativeTncSubaccountSeenAtBlockKey,
			},
		},
		"Block height can be updated to same block height": {
			setupTestAndPerformAssertions: func(ctx sdk.Context, k keeper.Keeper) {
				k.SetNegativeTncSubaccountSeenAtBlock(ctx, 0)
				require.Equal(
					t,
					uint32(0),
					k.GetNegativeTncSubaccountSeenAtBlock(ctx),
				)

				k.SetNegativeTncSubaccountSeenAtBlock(ctx, 0)
				require.Equal(
					t,
					uint32(0),
					k.GetNegativeTncSubaccountSeenAtBlock(ctx),
				)
			},

			expectedMultiStoreWrites: []string{
				types.NegativeTncSubaccountSeenAtBlockKey,
				types.NegativeTncSubaccountSeenAtBlockKey,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state and test parameters.
			ctx, subaccountsKeeper, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, true)

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
	ctx, subaccountsKeeper, _, _, _, _, _, _ := keepertest.SubaccountsKeepers(t, true)

	subaccountsKeeper.SetNegativeTncSubaccountSeenAtBlock(ctx, 2)
	require.PanicsWithValue(
		t,
		"SetNegativeTncSubaccountSeenAtBlock: new block height (1) is less than the current block height (2)",
		func() { subaccountsKeeper.SetNegativeTncSubaccountSeenAtBlock(ctx, 1) },
	)
}
