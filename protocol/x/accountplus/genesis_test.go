package accountplus_test

import (
	"math"
	"sort"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestImportExportGenesis(t *testing.T) {
	baseTsNonce := uint64(math.Pow(2, 40))
	tests := map[string]struct {
		genesisState *types.GenesisState
	}{
		"non-empty genesis": {
			genesisState: &types.GenesisState{
				Accounts: []*types.AccountState{
					{
						Address: constants.AliceAccAddress.String(),
						TimestampNonceDetails: &types.TimestampNonceDetails{
							TimestampNonces:    []uint64{baseTsNonce + 1, baseTsNonce + 2, baseTsNonce + 3},
							LatestEjectedNonce: baseTsNonce,
						},
					},
					{
						Address: constants.BobAccAddress.String(),
						TimestampNonceDetails: &types.TimestampNonceDetails{
							TimestampNonces:    []uint64{baseTsNonce + 5, baseTsNonce + 6, baseTsNonce + 7},
							LatestEjectedNonce: baseTsNonce + 1,
						},
					},
				},
			},
		},
		"empty genesis": {
			genesisState: &types.GenesisState{
				Accounts: []*types.AccountState{},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, k, _, _ := keepertest.TimestampNonceKeepers(t)

			// Initialize genesis state
			accountplus.InitGenesis(ctx, *k, *tc.genesisState)

			// Check that keeper state is correct
			compareKeeperWithGenesisState(t, ctx, k, tc.genesisState)

			// Export the genesis state
			exportedGenesis := accountplus.ExportGenesis(ctx, *k)

			// Ensure the exported state matches the expected state
			requireGenesisStatesEqual(t, tc.genesisState, &exportedGenesis)
		})
	}
}

func compareKeeperWithGenesisState(t *testing.T, ctx sdk.Context, k *keeper.Keeper, genesisState *types.GenesisState) {
	accountStates := k.GetAllAccountStates(ctx)

	compareAccountStates(t, accountStates, genesisState.GetAccounts())
}

func requireGenesisStatesEqual(t *testing.T, actualGenesisState, expectedGenesisState *types.GenesisState) {
	compareAccountStates(t, actualGenesisState.GetAccounts(), expectedGenesisState.GetAccounts())
}

func compareAccountStates(t *testing.T, actualAccountStates, expectedAccountStates []*types.AccountState) {
	require.Equal(t, len(actualAccountStates), len(expectedAccountStates), "GenesisState.Accounts length mismatch")

	// Sort the Accounts by ID to ensure we compare the correct Accounts
	sort.Slice(actualAccountStates, func(i, j int) bool {
		return actualAccountStates[i].Address < actualAccountStates[j].Address
	})
	sort.Slice(expectedAccountStates, func(i, j int) bool {
		return expectedAccountStates[i].Address < expectedAccountStates[j].Address
	})

	// Iterate through the account states and test equality on each field
	for i := range actualAccountStates {
		require.Equal(
			t,
			actualAccountStates[i].Address,
			expectedAccountStates[i].Address,
			"Account address mismatch at index %d", i,
		)
		compareTimestampNonceDetails(
			t,
			actualAccountStates[i].GetTimestampNonceDetails(),
			expectedAccountStates[i].GetTimestampNonceDetails(),
		)
	}
}

func compareTimestampNonceDetails(t *testing.T, actualDetails, expectedDetails *types.TimestampNonceDetails) {
	equal := cmp.Equal(
		actualDetails.GetTimestampNonces(),
		expectedDetails.GetTimestampNonces(),
		cmpopts.SortSlices(func(a, b uint64) bool { return a < b }),
	)

	require.True(t, equal, "TimestampNonces mismatch for account")

	require.Equal(
		t,
		actualDetails.GetLatestEjectedNonce(),
		expectedDetails.GetLatestEjectedNonce(),
		"LastEjectedNonce mismatch",
	)
}
