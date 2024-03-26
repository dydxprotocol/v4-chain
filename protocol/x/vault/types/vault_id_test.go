package types_test

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestToStateKey(t *testing.T) {
	b, _ := constants.Vault_Clob_0.Marshal()
	require.Equal(t, b, constants.Vault_Clob_0.ToStateKey())

	b, _ = constants.Vault_Clob_1.Marshal()
	require.Equal(t, b, constants.Vault_Clob_1.ToStateKey())
}

func TestToModuleAccountAddress(t *testing.T) {
	require.Equal(
		t,
		authtypes.NewModuleAddress("vault-VAULT_TYPE_CLOB-0").String(),
		constants.Vault_Clob_0.ToModuleAccountAddress(),
	)
	require.Equal(
		t,
		authtypes.NewModuleAddress("vault-VAULT_TYPE_CLOB-1").String(),
		constants.Vault_Clob_1.ToModuleAccountAddress(),
	)
}

func TestToSubaccountId(t *testing.T) {
	require.Equal(
		t,
		satypes.SubaccountId{
			Owner:  constants.Vault_Clob_0.ToModuleAccountAddress(),
			Number: 0,
		},
		*constants.Vault_Clob_0.ToSubaccountId(),
	)
	require.Equal(
		t,
		satypes.SubaccountId{
			Owner:  constants.Vault_Clob_1.ToModuleAccountAddress(),
			Number: 0,
		},
		*constants.Vault_Clob_1.ToSubaccountId(),
	)
}
