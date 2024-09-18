package keeper_test

import (
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	assetstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

func TestMsgAllocateToVault(t *testing.T) {
	tests := map[string]struct {
		// Operator.
		operator string
		// Number of quote quantums main vault has.
		mainVaultQuantums uint64
		// Number of quote quantums sub vault has.
		subVaultQuantums uint64
		// Existing vault params, if any.
		vaultParams *vaulttypes.VaultParams
		// Msg.
		msg *vaulttypes.MsgAllocateToVault
		// Expected error.
		expectedErr string
	}{
		"Success - Gov Authority, Allocate 50 to Vault Clob 0, Existing vault params": {
			operator:          constants.AliceAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  0,
			vaultParams: &vaulttypes.VaultParams{
				Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			},
			msg: &vaulttypes.MsgAllocateToVault{
				Authority:     lib.GovModuleAddress.String(),
				VaultId:       constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewInt(50),
			},
		},
		"Success - Gov Authority, Allocate 77 to Vault Clob 1, Non-existent Vault Params": {
			operator:          constants.AliceAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  15,
			msg: &vaulttypes.MsgAllocateToVault{
				Authority:     lib.GovModuleAddress.String(),
				VaultId:       constants.Vault_Clob1,
				QuoteQuantums: dtypes.NewInt(77),
			},
		},
		"Success - Operator Authority, Allocate all to Vault Clob 1, Existing vault params": {
			operator:          constants.AliceAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  15,
			vaultParams: &vaulttypes.VaultParams{
				Status: vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
			},
			msg: &vaulttypes.MsgAllocateToVault{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob1,
				QuoteQuantums: dtypes.NewInt(100),
			},
		},
		"Failure - Operator Authority, Insufficient quantums to allocate to Vault Clob 0, Existing vault params": {
			operator:          constants.AliceAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  15,
			vaultParams: &vaulttypes.VaultParams{
				Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			},
			msg: &vaulttypes.MsgAllocateToVault{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewInt(101),
			},
			expectedErr: "failed to apply subaccount updates",
		},
		"Failure - Operator Authority, No corresponding clob pair": {
			operator:          constants.AliceAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  0,
			msg: &vaulttypes.MsgAllocateToVault{
				Authority: constants.AliceAccAddress.String(),
				VaultId: vaulttypes.VaultId{
					Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
					Number: 727,
				},
				QuoteQuantums: dtypes.NewInt(1),
			},
			expectedErr: vaulttypes.ErrClobPairNotFound.Error(),
		},
		"Failure - Invalid Authority, Non-existent Vault Params": {
			operator:          constants.BobAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  15,
			msg: &vaulttypes.MsgAllocateToVault{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob1,
				QuoteQuantums: dtypes.NewInt(77),
			},
			expectedErr: vaulttypes.ErrInvalidAuthority.Error(),
		},
		"Failure - Empty Authority, Existing vault params": {
			operator:          constants.BobAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  15,
			msg: &vaulttypes.MsgAllocateToVault{
				Authority:     "",
				VaultId:       constants.Vault_Clob1,
				QuoteQuantums: dtypes.NewInt(77),
			},
			expectedErr: vaulttypes.ErrInvalidAuthority.Error(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Set megavault operator.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vaulttypes.GenesisState) {
						genesisState.OperatorParams = vaulttypes.OperatorParams{
							Operator: tc.operator,
						}
						if tc.vaultParams != nil {
							genesisState.Vaults = []vaulttypes.Vault{
								{
									VaultId:     tc.msg.VaultId,
									VaultParams: *tc.vaultParams,
								},
							}
						}
					},
				)
				// Set balances of main vault and sub vault.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id: &vaulttypes.MegavaultMainSubaccount,
								AssetPositions: []*satypes.AssetPosition{
									{
										AssetId:  assetstypes.AssetUsdc.Id,
										Quantums: dtypes.NewIntFromUint64(tc.mainVaultQuantums),
									},
								},
							},
							{
								Id: tc.msg.VaultId.ToSubaccountId(),
								AssetPositions: []*satypes.AssetPosition{
									{
										AssetId:  assetstypes.AssetUsdc.Id,
										Quantums: dtypes.NewIntFromUint64(tc.subVaultQuantums),
									},
								},
							},
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper
			ms := keeper.NewMsgServerImpl(k)

			// Allocate to vault.
			_, err := ms.AllocateToVault(ctx, tc.msg)

			// Check expectations.
			mainVault := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, vaulttypes.MegavaultMainSubaccount)
			subVault := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *tc.msg.VaultId.ToSubaccountId())
			require.Len(t, subVault.AssetPositions, 1)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)

				// Verify that main vault and sub vault balances are unchanged.
				require.Len(t, mainVault.AssetPositions, 1)
				require.Equal(
					t,
					tc.mainVaultQuantums,
					mainVault.AssetPositions[0].Quantums.BigInt().Uint64(),
				)
				require.Equal(
					t,
					tc.subVaultQuantums,
					subVault.AssetPositions[0].Quantums.BigInt().Uint64(),
				)

				// Verify that vault params is unchanged.
				vaultParams, exists := k.GetVaultParams(ctx, tc.msg.VaultId)
				if tc.vaultParams != nil {
					require.True(t, exists)
					require.Equal(t, *tc.vaultParams, vaultParams)
				} else {
					require.False(t, exists)
				}
			} else {
				require.NoError(t, err)

				// Verify that main vault and sub vault balances are updated.
				expectedMainVaultQuantums := tc.mainVaultQuantums - tc.msg.QuoteQuantums.BigInt().Uint64()
				if expectedMainVaultQuantums == 0 {
					require.Len(t, mainVault.AssetPositions, 0)
				} else {
					require.Len(t, mainVault.AssetPositions, 1)
					require.Equal(
						t,
						expectedMainVaultQuantums,
						mainVault.AssetPositions[0].Quantums.BigInt().Uint64(),
					)
				}
				require.Equal(
					t,
					tc.subVaultQuantums+tc.msg.QuoteQuantums.BigInt().Uint64(),
					subVault.AssetPositions[0].Quantums.BigInt().Uint64(),
				)

				// Verify that vault params is initialized if didn't exist before.
				vaultParams, exists := k.GetVaultParams(ctx, tc.msg.VaultId)
				require.True(t, exists)
				if tc.vaultParams != nil {
					require.Equal(t, *tc.vaultParams, vaultParams)
				} else {
					require.Equal(
						t,
						vaulttypes.VaultParams{
							Status: vaulttypes.VaultStatus_VAULT_STATUS_STAND_BY,
						},
						vaultParams,
					)
				}
			}
		})
	}
}
