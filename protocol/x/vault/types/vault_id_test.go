package types_test

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestToString(t *testing.T) {
	tests := map[string]struct {
		// Vault ID.
		vaultId types.VaultId
		// Expected string.
		expectedStr string
	}{
		"Vault for Clob Pair 0": {
			vaultId:     constants.Vault_Clob0,
			expectedStr: "VAULT_TYPE_CLOB-0",
		},
		"Vault for Clob Pair 1": {
			vaultId:     constants.Vault_Clob1,
			expectedStr: "VAULT_TYPE_CLOB-1",
		},
		"Vault, missing type and number": {
			vaultId:     types.VaultId{},
			expectedStr: "VAULT_TYPE_UNSPECIFIED-0",
		},
		"Vault, missing type": {
			vaultId: types.VaultId{
				Number: 1,
			},
			expectedStr: "VAULT_TYPE_UNSPECIFIED-1",
		},
		"Vault, missing number": {
			vaultId: types.VaultId{
				Type: types.VaultType_VAULT_TYPE_CLOB,
			},
			expectedStr: "VAULT_TYPE_CLOB-0",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.vaultId.ToString(),
				tc.expectedStr,
			)
		})
	}
}

func TestToStateKey(t *testing.T) {
	require.Equal(
		t,
		[]byte("VAULT_TYPE_CLOB-0"),
		constants.Vault_Clob0.ToStateKey(),
	)

	require.Equal(
		t,
		[]byte("VAULT_TYPE_CLOB-1"),
		constants.Vault_Clob1.ToStateKey(),
	)
}

func TestGetVaultIdFromStateKey(t *testing.T) {
	tests := map[string]struct {
		// State key.
		stateKey []byte
		// Expected vault ID.
		expectedVaultId types.VaultId
		// Expected error.
		expectedErr string
	}{
		"Vault for Clob Pair 0": {
			stateKey:        []byte("VAULT_TYPE_CLOB-0"),
			expectedVaultId: constants.Vault_Clob0,
		},
		"Vault for Clob Pair 1": {
			stateKey:        []byte("VAULT_TYPE_CLOB-1"),
			expectedVaultId: constants.Vault_Clob1,
		},
		"Empty bytes": {
			stateKey:    []byte{},
			expectedErr: "stateKey in string must follow format <type>-<number>",
		},
		"Nil bytes": {
			stateKey:    nil,
			expectedErr: "stateKey in string must follow format <type>-<number>",
		},
		"Non-existent vault type": {
			stateKey:    []byte("VAULT_TYPE_SPOT-1"),
			expectedErr: "unknown vault type",
		},
		"Malformed vault number": {
			stateKey:    []byte("VAULT_TYPE_CLOB-abc"),
			expectedErr: "failed to parse number",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vaultId, err := types.GetVaultIdFromStateKey(tc.stateKey)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedVaultId, *vaultId)
			}
		})
	}
}

func TestToModuleAccountAddress(t *testing.T) {
	require.Equal(
		t,
		authtypes.NewModuleAddress("vault-VAULT_TYPE_CLOB-0").String(),
		constants.Vault_Clob0.ToModuleAccountAddress(),
	)
	require.Equal(
		t,
		authtypes.NewModuleAddress("vault-VAULT_TYPE_CLOB-1").String(),
		constants.Vault_Clob1.ToModuleAccountAddress(),
	)
}

func TestToSubaccountId(t *testing.T) {
	require.Equal(
		t,
		satypes.SubaccountId{
			Owner:  constants.Vault_Clob0.ToModuleAccountAddress(),
			Number: 0,
		},
		*constants.Vault_Clob0.ToSubaccountId(),
	)
	require.Equal(
		t,
		satypes.SubaccountId{
			Owner:  constants.Vault_Clob1.ToModuleAccountAddress(),
			Number: 0,
		},
		*constants.Vault_Clob1.ToSubaccountId(),
	)
}
