package keeper_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func setVaultState(t *testing.T, tApp *testapp.TestApp, ctx sdk.Context, vaultState constants.VaultState) {
	totalShares := big.NewInt(0)
	for _, ownerShare := range vaultState.OwnerShares {
		totalShares.Add(totalShares, ownerShare.Shares.NumShares.BigInt())

		err := tApp.App.VaultKeeper.SetOwnerShares(
			ctx,
			vaultState.VaultId,
			ownerShare.Owner,
			*ownerShare.Shares,
		)
		require.NoError(t, err)
	}

	// Set vault's total shares.
	err := tApp.App.VaultKeeper.SetTotalShares(
		ctx,
		vaultState.VaultId,
		vaulttypes.BigIntToNumShares(totalShares),
	)
	require.NoError(t, err)
}

func TestValidateWithdrawFromVault(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		vaultState constants.VaultState

		/* --- Inputs --- */
		msg vaulttypes.MsgWithdrawFromVault

		/* --- Expectations --- */
		expectedErr error
	}{
		"Success: single owner": {
			vaultState: constants.Vault_Clob0_SingleOwner_Alice0_1000,
			msg: vaulttypes.MsgWithdrawFromVault{
				VaultId:       &constants.Vault_Clob0,
				SubaccountId:  &constants.Alice_Num0,
				QuoteQuantums: dtypes.NewInt(100),
			},
			expectedErr: nil,
		},
		"Success: multiple owners": {
			vaultState: constants.Vault_Clob0_MultiOwner_Alice0_1000_Bob0_2500,
			msg: vaulttypes.MsgWithdrawFromVault{
				VaultId:       &constants.Vault_Clob0,
				SubaccountId:  &constants.Bob_Num0,
				QuoteQuantums: dtypes.NewInt(100),
			},
			expectedErr: nil,
		},
		"Failure: no vault": {
			vaultState: constants.VaultState{}, // nil vault state.
			msg: vaulttypes.MsgWithdrawFromVault{
				VaultId:       &constants.Vault_Clob0, // vault 0 doesn't exist.
				SubaccountId:  &constants.Alice_Num0,
				QuoteQuantums: dtypes.NewInt(100),
			},
			expectedErr: vaulttypes.ErrVaultNotFound,
		},
		"Failure: no matching vault": {
			vaultState: constants.Vault_Clob0_SingleOwner_Alice0_1000,
			msg: vaulttypes.MsgWithdrawFromVault{
				VaultId:       &constants.Vault_Clob1, // vault 1 doesn't exist.
				SubaccountId:  &constants.Alice_Num0,
				QuoteQuantums: dtypes.NewInt(100),
			},
			expectedErr: vaulttypes.ErrVaultNotFound,
		},
		"Failure: no matching shares": {
			vaultState: constants.Vault_Clob0_MultiOwner_Alice0_1000_Bob0_2500,
			msg: vaulttypes.MsgWithdrawFromVault{
				VaultId:       &constants.Vault_Clob0,
				SubaccountId:  &constants.Carl_Num0, // Carl doesn't have any shares.
				QuoteQuantums: dtypes.NewInt(100),
			},
			expectedErr: vaulttypes.ErrOwnerShareNotFound,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize tApp and ctx.
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				return testapp.DefaultGenesis()
			}).Build()
			ctx := tApp.InitChain()

			// Set vault's state if not nil.
			if tc.vaultState.VaultId != (vaulttypes.VaultId{}) {
				setVaultState(t, tApp, ctx, tc.vaultState)
			}

			// Run.
			err := tApp.App.VaultKeeper.ValidateWithdrawFromVault(ctx, &tc.msg)

			// Validate
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
