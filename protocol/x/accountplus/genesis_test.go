package accountplus_test

import (
	"math"
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	"github.com/stretchr/testify/require"
)

func TestImportExportGenesis(t *testing.T) {
	baseTsNonce := uint64(math.Pow(2, 40))
	tests := map[string]struct {
		genesisState *types.GenesisState
		// The order of this list may not match the order in GenesisState. We want our tests to be deterministic so
		// order of expectedAccountStates was manually set based on test debug. This ordering should only be changed if
		// additional accounts added to genesisState. If a feature breaks the existing ordering, should look into why.
		expectedAccountStates []types.AccountState
	}{
		"non-empty genesis": {
			genesisState: &types.GenesisState{
				Accounts: []types.AccountState{
					{
						Address: constants.AliceAccAddress.String(),
						TimestampNonceDetails: types.TimestampNonceDetails{
							TimestampNonces: []uint64{baseTsNonce + 1, baseTsNonce + 2, baseTsNonce + 3},
							MaxEjectedNonce: baseTsNonce,
						},
					},
					{
						Address: constants.BobAccAddress.String(),
						TimestampNonceDetails: types.TimestampNonceDetails{
							TimestampNonces: []uint64{baseTsNonce + 5, baseTsNonce + 6, baseTsNonce + 7},
							MaxEjectedNonce: baseTsNonce + 1,
						},
					},
					{
						Address: constants.CarlAccAddress.String(),
						TimestampNonceDetails: types.TimestampNonceDetails{
							TimestampNonces: []uint64{baseTsNonce + 5, baseTsNonce + 6, baseTsNonce + 7},
							MaxEjectedNonce: baseTsNonce + 1,
						},
					},
					{
						Address: constants.CarlAccAddress.String(),
						TimestampNonceDetails: types.TimestampNonceDetails{
							TimestampNonces: []uint64{baseTsNonce + 5, baseTsNonce + 6, baseTsNonce + 7},
							MaxEjectedNonce: baseTsNonce + 1,
						},
					},
				},
			},
			expectedAccountStates: []types.AccountState{
				{
					Address: constants.AliceAccAddress.String(),
					TimestampNonceDetails: types.TimestampNonceDetails{
						TimestampNonces: []uint64{baseTsNonce + 1, baseTsNonce + 2, baseTsNonce + 3},
						MaxEjectedNonce: baseTsNonce,
					},
				},
				{
					Address: constants.CarlAccAddress.String(),
					TimestampNonceDetails: types.TimestampNonceDetails{
						TimestampNonces: []uint64{baseTsNonce + 5, baseTsNonce + 6, baseTsNonce + 7},
						MaxEjectedNonce: baseTsNonce + 1,
					},
				},
				{
					Address: constants.BobAccAddress.String(),
					TimestampNonceDetails: types.TimestampNonceDetails{
						TimestampNonces: []uint64{baseTsNonce + 5, baseTsNonce + 6, baseTsNonce + 7},
						MaxEjectedNonce: baseTsNonce + 1,
					},
				},
			},
		},
		"empty genesis": {
			genesisState: &types.GenesisState{
				Accounts: []types.AccountState{},
			},
			expectedAccountStates: []types.AccountState{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.AccountPlusKeeper

			accountplus.InitGenesis(ctx, k, *tc.genesisState)

			// Check that keeper accounts states are correct
			actualAccountStates, _ := k.GetAllAccountStates(ctx)
			require.Equal(
				t,
				tc.expectedAccountStates,
				actualAccountStates,
				"Keeper account states do not match Genesis account states",
			)

			exportedGenesis := accountplus.ExportGenesis(ctx, k)

			// Check that the exported state matches the expected state
			expectedGenesis := &types.GenesisState{
				Accounts: tc.expectedAccountStates,
			}
			require.Equal(t, *exportedGenesis, *expectedGenesis)
		})
	}
}
