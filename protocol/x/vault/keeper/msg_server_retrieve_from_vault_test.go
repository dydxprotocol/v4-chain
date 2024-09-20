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

func TestMsgRetrieveFromVault(t *testing.T) {
	tests := map[string]struct {
		// Operator.
		operator string
		// Number of quote quantums main vault has.
		mainVaultQuantums uint64
		// Number of quote quantums sub vault has.
		subVaultQuantums uint64
		// Number of base quantums of sub vault's position.
		subVaultPositionBaseQuantums int64
		// Existing vault params, if any.
		vaultParams *vaulttypes.VaultParams
		// Msg.
		msg *vaulttypes.MsgRetrieveFromVault
		// Expected error.
		expectedErr string
	}{
		"Success - Gov Authority, Retrieve 50 From Vault Clob 0": {
			operator:          constants.AliceAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  200,
			vaultParams: &vaulttypes.VaultParams{
				Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			},
			msg: &vaulttypes.MsgRetrieveFromVault{
				Authority:     lib.GovModuleAddress.String(),
				VaultId:       constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewInt(50),
			},
		},
		"Success - Operator Authority, Retrieve all from Vault Clob 1": {
			operator:          constants.AliceAccAddress.String(),
			mainVaultQuantums: 0,
			subVaultQuantums:  3_500_000,
			vaultParams: &vaulttypes.VaultParams{
				Status: vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
			},
			msg: &vaulttypes.MsgRetrieveFromVault{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob1,
				QuoteQuantums: dtypes.NewInt(3_500_000),
			},
		},
		"Failure - Insufficient quantums to retrieve from Vault Clob 0 with no open position": {
			operator:          constants.AliceAccAddress.String(),
			mainVaultQuantums: 0,
			subVaultQuantums:  26,
			vaultParams: &vaulttypes.VaultParams{
				Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			},
			msg: &vaulttypes.MsgRetrieveFromVault{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewInt(27),
			},
			expectedErr: "failed to apply subaccount updates",
		},
		"Success - Retrieval from vault with open position exactly meets initial margin requirement": {
			operator:                     constants.AliceAccAddress.String(),
			mainVaultQuantums:            0,
			subVaultQuantums:             3_500_000,
			subVaultPositionBaseQuantums: -1_000_000,
			// open_notional = -1_000_000 * 10^-9 * 1_500 * 10^6 = = -1_500_000
			// equity = 3_500_000 - 1_500_000 = 2_000_000
			// initial_margin_requirement = position_size * imf
			// = |-1_500_000| * 0.05 = 75_000
			vaultParams: &vaulttypes.VaultParams{
				Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			},
			msg: &vaulttypes.MsgRetrieveFromVault{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob1,
				QuoteQuantums: dtypes.NewInt(1_925_000),
			},
		},
		"Failure - Retrieval from vault with open position would result in undercollateralization": {
			operator:                     constants.AliceAccAddress.String(),
			mainVaultQuantums:            0,
			subVaultQuantums:             3_500_000,
			subVaultPositionBaseQuantums: -1_000_000,
			// open_notional = -1_000_000 * 10^-9 * 1_500 * 10^6 = = -1_500_000
			// equity = 3_500_000 - 1_500_000 = 2_000_000
			// initial_margin_requirement = position_size * imf
			// = |-1_500_000| * 0.05 = 75_000
			vaultParams: &vaulttypes.VaultParams{
				Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			},
			msg: &vaulttypes.MsgRetrieveFromVault{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob1,
				QuoteQuantums: dtypes.NewInt(1_925_001),
			},
			expectedErr: satypes.NewlyUndercollateralized.String(),
		},
		"Failure - Retrieve from non-existent vault": {
			operator:          constants.AliceAccAddress.String(),
			mainVaultQuantums: 0,
			subVaultQuantums:  15,
			msg: &vaulttypes.MsgRetrieveFromVault{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewInt(10),
			},
			expectedErr: vaulttypes.ErrVaultParamsNotFound.Error(),
		},
		"Failure - Invalid Authority": {
			operator:          constants.BobAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  15,
			msg: &vaulttypes.MsgRetrieveFromVault{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob1,
				QuoteQuantums: dtypes.NewInt(10),
			},
			expectedErr: vaulttypes.ErrInvalidAuthority.Error(),
		},
		"Failure - Empty Authority": {
			operator:          constants.BobAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  15,
			msg: &vaulttypes.MsgRetrieveFromVault{
				Authority:     "",
				VaultId:       constants.Vault_Clob1,
				QuoteQuantums: dtypes.NewInt(10),
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
								PerpetualPositions: []*satypes.PerpetualPosition{
									{
										PerpetualId: tc.msg.VaultId.Number,
										Quantums:    dtypes.NewInt(tc.subVaultPositionBaseQuantums),
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

			// Retrieve from vault.
			_, err := ms.RetrieveFromVault(ctx, tc.msg)

			// Check expectations.
			mainVault := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, vaulttypes.MegavaultMainSubaccount)
			subVault := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *tc.msg.VaultId.ToSubaccountId())

			require.Len(t, mainVault.AssetPositions, 1)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)

				// Verify that main vault and sub vault balances are unchanged.
				require.Len(t, subVault.AssetPositions, 1)
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
			} else {
				require.NoError(t, err)

				// Verify that main vault and sub vault balances are updated.
				expectedSubVaultQuantums := tc.subVaultQuantums - tc.msg.QuoteQuantums.BigInt().Uint64()
				if expectedSubVaultQuantums == 0 {
					require.Len(t, subVault.AssetPositions, 0)
				} else {
					require.Len(t, subVault.AssetPositions, 1)
					require.Equal(
						t,
						expectedSubVaultQuantums,
						subVault.AssetPositions[0].Quantums.BigInt().Uint64(),
					)
				}
				require.Equal(
					t,
					tc.mainVaultQuantums+tc.msg.QuoteQuantums.BigInt().Uint64(),
					mainVault.AssetPositions[0].Quantums.BigInt().Uint64(),
				)
			}
		})
	}
}
