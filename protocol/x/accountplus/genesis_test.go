package accountplus_test

import (
	"sort"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	"github.com/stretchr/testify/require"
)

func TestImportExportGenesis(t *testing.T) {
	tests := map[string]struct {
		genesisState *types.GenesisState
	}{
		"non-empty genesis": {
			genesisState: &types.GenesisState{
				Accounts: []*types.AccountState{
					{
						Address: constants.AliceAccAddress.String(),
						Details: &types.TimestampNonceDetails{
							TimestampNonces:    []uint64{2 ^ 40 + 1, 2 ^ 40 + 2, 2 ^ 40 + 3},
							LatestEjectedNonce: 0,
						},
					},
					{
						Address: constants.BobAccAddress.String(),
						Details: &types.TimestampNonceDetails{
							TimestampNonces:    []uint64{2 ^ 40 + 5, 2 ^ 40 + 6, 2 ^ 40 + 7},
							LatestEjectedNonce: 3,
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
	accountStates := k.GetAllAccoutDetails(ctx)

	compareAccountStateLists(t, accountStates, genesisState.GetAccounts())
}

func requireGenesisStatesEqual(t *testing.T, actualGenesisState, expectedGenesisState *types.GenesisState) {
	compareAccountStateLists(t, actualGenesisState.GetAccounts(), expectedGenesisState.GetAccounts())
}

func compareAccountStateLists(t *testing.T, actualAccountStates, expectedAccountStates []*types.AccountState) {
	require.Equal(t, len(actualAccountStates), len(expectedAccountStates), "GenesisState.Accounts length mismatch")

	// Sort the Accounts by ID to ensure we compare the correct Accounts
	sort.Slice(actualAccountStates, func(i, j int) bool {
		return actualAccountStates[i].Address < actualAccountStates[j].Address
	})
	sort.Slice(expectedAccountStates, func(i, j int) bool {
		return expectedAccountStates[i].Address < expectedAccountStates[j].Address
	})

	for i := range actualAccountStates {
		require.Equal(t, actualAccountStates[i].Address, expectedAccountStates[i].Address, "Account address mismatch at index %d", i)
		compareDetailsLists(t, actualAccountStates[i].GetDetails(), expectedAccountStates[i].GetDetails())
	}
}

func compareDetailsLists(t *testing.T, actualDetail, expectedDetail *types.TimestampNonceDetails) {
	actualTsNonces := actualDetail.GetTimestampNonces()
	expectedTsNonces := expectedDetail.GetTimestampNonces()

	require.Equal(t, len(actualTsNonces), len(expectedTsNonces), "GenesisState.Accounts.TimestampNonceDetails.TimestampNonces length mismatch")

	// Sort the values to ignore order
	sort.Slice(actualTsNonces, func(i, j int) bool { return actualTsNonces[i] < actualTsNonces[j] })
	sort.Slice(expectedTsNonces, func(i, j int) bool { return expectedTsNonces[i] < expectedTsNonces[j] })

	require.Equal(t, actualDetail.GetLatestEjectedNonce(), expectedDetail.GetLatestEjectedNonce(), "LastEjectedNonce mismatch")
	require.Equal(t, actualTsNonces, expectedTsNonces, "TimestampNonces mismatch")
}
