package keeper_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	"github.com/stretchr/testify/require"
)

func TestQueryListingVaultDepositParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.ListingKeeper

	params := types.ListingVaultDepositParams{
		NewVaultDepositAmount:  dtypes.NewIntFromBigInt(big.NewInt(100_000_000)),
		MainVaultDepositAmount: dtypes.NewIntFromBigInt(big.NewInt(0)),
		NumBlocksToLockShares:  30 * 24 * 3600, // 30 days
	}

	err := k.SetListingVaultDepositParams(ctx, params)
	require.NoError(t, err)

	resp, err := k.ListingVaultDepositParams(ctx, &types.QueryListingVaultDepositParams{})
	require.NoError(t, err)
	require.Equal(t, resp.Params, params)
}
