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
		expectedAccountStates       []types.AccountState
		expectedParams              types.Params
		expectedNextAuthenticatorId uint64
		expectedAuthenticatorData   []types.AuthenticatorData
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
				},
				Params: types.Params{
					IsSmartAccountActive: true,
				},
				NextAuthenticatorId: 100,
				AuthenticatorData: []types.AuthenticatorData{
					{
						Address: constants.AliceAccAddress.String(),
						Authenticators: []types.AccountAuthenticator{
							{
								Id:     1,
								Type:   "MessageFilter",
								Config: []byte("/cosmos.bank.v1beta1.MsgSend"),
							},
						},
					},
					{
						Address: constants.BobAccAddress.String(),
						Authenticators: []types.AccountAuthenticator{
							{
								Id:     1,
								Type:   "ClobPairIdFilter",
								Config: []byte("0,1,2"),
							},
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
			expectedParams: types.Params{
				IsSmartAccountActive: true,
			},
			expectedNextAuthenticatorId: 100,
			expectedAuthenticatorData: []types.AuthenticatorData{
				{
					Address: constants.BobAccAddress.String(),
					Authenticators: []types.AccountAuthenticator{
						{
							Id:     1,
							Type:   "ClobPairIdFilter",
							Config: []byte("0,1,2"),
						},
					},
				},
				{
					Address: constants.AliceAccAddress.String(),
					Authenticators: []types.AccountAuthenticator{
						{
							Id:     1,
							Type:   "MessageFilter",
							Config: []byte("/cosmos.bank.v1beta1.MsgSend"),
						},
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

			// Check that keeper params are correct
			actualParams := k.GetParams(ctx)
			require.Equal(t, tc.expectedParams, actualParams)

			// Check that keeper next authenticator id is correct
			actualNextAuthenticatorId := k.InitializeOrGetNextAuthenticatorId(ctx)
			require.Equal(t, tc.expectedNextAuthenticatorId, actualNextAuthenticatorId)

			// Check that keeper authenticator data is correct
			actualAuthenticatorData, _ := k.GetAllAuthenticatorData(ctx)
			require.Equal(t, tc.expectedAuthenticatorData, actualAuthenticatorData)

			exportedGenesis := accountplus.ExportGenesis(ctx, k)

			// Check that the exported state matches the expected state
			expectedGenesis := &types.GenesisState{
				Accounts:            tc.expectedAccountStates,
				Params:              tc.expectedParams,
				NextAuthenticatorId: tc.expectedNextAuthenticatorId,
				AuthenticatorData:   tc.expectedAuthenticatorData,
			}
			require.Equal(t, *exportedGenesis, *expectedGenesis)
		})
	}
}
