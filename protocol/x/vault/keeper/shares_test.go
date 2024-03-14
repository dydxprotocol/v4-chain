package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetTotalShares(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// Get total shares for a non-existing vault.
	_, exists := k.GetTotalShares(ctx, constants.Vault_Clob_0)
	require.Equal(t, false, exists)

	// Set total shares for a vault and then get.
	k.SetTotalShares(ctx, constants.Vault_Clob_0, types.NumShares{
		NumShares: dtypes.NewInt(7),
	})
	numShares, exists := k.GetTotalShares(ctx, constants.Vault_Clob_0)
	require.Equal(t, true, exists)
	require.Equal(t, dtypes.NewInt(7), numShares.NumShares)

	// Set total shares for another vault and then get.
	k.SetTotalShares(ctx, constants.Vault_Clob_1, types.NumShares{
		NumShares: dtypes.NewInt(456),
	})
	numShares, exists = k.GetTotalShares(ctx, constants.Vault_Clob_1)
	require.Equal(t, true, exists)
	require.Equal(t, dtypes.NewInt(456), numShares.NumShares)

	// Set total shares for the first vault again and then get.
	k.SetTotalShares(ctx, constants.Vault_Clob_0, types.NumShares{
		NumShares: dtypes.NewInt(123),
	})
	numShares, exists = k.GetTotalShares(ctx, constants.Vault_Clob_0)
	require.Equal(t, true, exists)
	require.Equal(t, dtypes.NewInt(123), numShares.NumShares)
}
