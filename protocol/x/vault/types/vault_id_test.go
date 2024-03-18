package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
)

func TestToStateKey(t *testing.T) {
	b, _ := constants.Vault_Clob_0.Marshal()
	require.Equal(t, b, constants.Vault_Clob_0.ToStateKey())

	b, _ = constants.Vault_Clob_1.Marshal()
	require.Equal(t, b, constants.Vault_Clob_1.ToStateKey())
}
