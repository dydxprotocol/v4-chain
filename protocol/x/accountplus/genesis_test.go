package accountplus_test

import (
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/testutils"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	"github.com/stretchr/testify/require"
)

func TestImportExportGenesis(t *testing.T) {
	baseTsNonce := uint64(math.Pow(2, 40))
	tests := map[string]struct {
		genesisState *types.GenesisState
		// The order of this list may not match the order in GenesisState. We want our tests to be deterministic so
		// order of expectedAccountStates was manually set based on test debug. This ordering should only be changed if
		// additional accounts added to genesisState. If a feature breaks the existing ordering, should look into ßwhy.
		expectedAccountStates []*types.AccountState
	}{
		"non-empty genesis": {
			genesisState: &types.GenesisState{
				Accounts: []*types.AccountState{
					{
						Address: constants.AliceAccAddress.String(),
						TimestampNonceDetails: &types.TimestampNonceDetails{
							TimestampNonces: []uint64{baseTsNonce + 1, baseTsNonce + 2, baseTsNonce + 3},
							MaxEjectedNonce: baseTsNonce,
						},
					},
					{
						Address: constants.BobAccAddress.String(),
						TimestampNonceDetails: &types.TimestampNonceDetails{
							TimestampNonces: []uint64{baseTsNonce + 5, baseTsNonce + 6, baseTsNonce + 7},
							MaxEjectedNonce: baseTsNonce + 1,
						},
					},
					{
						Address: constants.CarlAccAddress.String(),
						TimestampNonceDetails: &types.TimestampNonceDetails{
							TimestampNonces: []uint64{baseTsNonce + 5, baseTsNonce + 6, baseTsNonce + 7},
							MaxEjectedNonce: baseTsNonce + 1,
						},
					},
				},
			},
			expectedAccountStates: []*types.AccountState{
				{
					Address: constants.AliceAccAddress.String(),
					TimestampNonceDetails: &types.TimestampNonceDetails{
						TimestampNonces: []uint64{baseTsNonce + 1, baseTsNonce + 2, baseTsNonce + 3},
						MaxEjectedNonce: baseTsNonce,
					},
				},
				{
					Address: constants.CarlAccAddress.String(),
					TimestampNonceDetails: &types.TimestampNonceDetails{
						TimestampNonces: []uint64{baseTsNonce + 5, baseTsNonce + 6, baseTsNonce + 7},
						MaxEjectedNonce: baseTsNonce + 1,
					},
				},
				{
					Address: constants.BobAccAddress.String(),
					TimestampNonceDetails: &types.TimestampNonceDetails{
						TimestampNonces: []uint64{baseTsNonce + 5, baseTsNonce + 6, baseTsNonce + 7},
						MaxEjectedNonce: baseTsNonce + 1,
					},
				},
			},
		},
		"empty genesis": {
			genesisState: &types.GenesisState{
				Accounts: []*types.AccountState{},
			},
			expectedAccountStates: []*types.AccountState{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.AccountPlusKeeper

			// Initialize genesis state
			accountplus.InitGenesis(ctx, k, *tc.genesisState)

			// Check that keeper state is correct
			compareKeeperWithGenesisState(t, ctx, &k, tc.expectedAccountStates)

			// Export the genesis state
			exportedGenesis := accountplus.ExportGenesis(ctx, k)

			// Ensure the exported state matches the expected state
			expectedGenesis := &types.GenesisState{
				Accounts: tc.expectedAccountStates,
			}
			requireGenesisStatesEqual(t, exportedGenesis, expectedGenesis)
		})
	}
}

func compareKeeperWithGenesisState(
	t *testing.T,
	ctx sdk.Context,
	k *keeper.Keeper,
	expectedAccountStates []*types.AccountState,
) {
	// Compare states. Order matters.
	isEqual := testutils.CompareAccountStateLists(k.GetAllAccountStates(ctx), expectedAccountStates)

	require.True(t, isEqual, "Keeper account states does not match Genesis account states")
}

func requireGenesisStatesEqual(t *testing.T, actualGenesisState, expectedGenesisState *types.GenesisState) {
	isEqual := testutils.CompareAccountStateLists(actualGenesisState.GetAccounts(), expectedGenesisState.GetAccounts())

	require.True(t, isEqual, "GenesisState mismatch")
}
