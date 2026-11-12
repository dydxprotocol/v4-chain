package gov_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	canonicalusdctypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"
)

// The EndBlocker pauses GRADUAL routing as soon as the supply equation breaks,
// and leaves a consistent ledger alone.
func TestEndBlockerPausesOnBackingBreach(t *testing.T) {
	tests := map[string]struct {
		drift    int64
		wantMode canonicalusdctypes.Mode
	}{
		"consistent ledger stays gradual": {drift: 0, wantMode: canonicalusdctypes.Mode_MODE_GRADUAL},
		"supply drift pauses":             {drift: 1, wantMode: canonicalusdctypes.Mode_MODE_PAUSED},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() types.GenesisDoc {
				genesis := canonicalUsdcGenesis(true)
				supply := math.ZeroInt()
				testapp.UpdateGenesisDocWithAppStateForModule(&genesis, func(state *banktypes.GenesisState) {
					for _, balance := range state.Balances {
						supply = supply.Add(balance.Coins.AmountOf(assettypes.UusdcDenom))
					}
				})
				require.True(t, supply.IsPositive())
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(state *canonicalusdctypes.GenesisState) {
						state.Controls = gradualControls()
						state.Ledger = canonicalusdctypes.DefaultLedger()
						state.Ledger.NobleBacking = dtypes.NewIntFromBigInt(supply.SubRaw(tc.drift).BigInt())
					},
				)
				return genesis
			}).Build()
			// InitChain also executes block 1, including its EndBlocker.
			ctx := tApp.InitChain()
			require.Equal(t, tc.wantMode, tApp.App.CanonicalUsdcKeeper.GetControls(ctx).Mode)

			ctx = tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			require.Equal(t, tc.wantMode, tApp.App.CanonicalUsdcKeeper.GetControls(ctx).Mode)
		})
	}
}
