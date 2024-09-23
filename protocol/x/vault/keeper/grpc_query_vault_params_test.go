package keeper_test

import (
	"math/big"
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestVaultParamsQuery(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Vault ID.
		vaultId vaulttypes.VaultId
		// Vault params.
		vaultParams vaulttypes.VaultParams
		// Query request.
		req *vaulttypes.QueryVaultParamsRequest

		/* --- Expectations --- */
		expectedEquity *big.Int
		expectedErr    string
	}{
		"Success": {
			req: &vaulttypes.QueryVaultParamsRequest{
				Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
				Number: 0,
			},
			vaultId:     constants.Vault_Clob0,
			vaultParams: constants.VaultParams,
		},
		"Success: close only vault status": {
			req: &vaulttypes.QueryVaultParamsRequest{
				Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
				Number: 0,
			},
			vaultId: constants.Vault_Clob0,
			vaultParams: vaulttypes.VaultParams{
				Status: vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
			},
		},
		"Error: query non-existent vault": {
			req: &vaulttypes.QueryVaultParamsRequest{
				Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
				Number: 1, // Non-existent vault.
			},
			vaultId:     constants.Vault_Clob0,
			vaultParams: constants.VaultParams,
			expectedErr: "vault not found",
		},
		"Error: nil request": {
			req:         nil,
			vaultId:     constants.Vault_Clob0,
			vaultParams: constants.VaultParams,
			expectedErr: "invalid request",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			// Set vault params.
			err := k.SetVaultParams(ctx, tc.vaultId, tc.vaultParams)
			require.NoError(t, err)

			// Check Vault query response is as expected.
			response, err := k.VaultParams(ctx, tc.req)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				expectedResponse := vaulttypes.QueryVaultParamsResponse{
					VaultId:     tc.vaultId,
					VaultParams: tc.vaultParams,
				}
				require.Equal(t, expectedResponse, *response)
			}
		})
	}
}
