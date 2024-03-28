package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// Params should have default values at genesis.
	params := k.GetParams(ctx)
	require.Equal(t, types.DefaultParams(), params)

	// Set new params and get.
	newParams := types.Params{
		Layers:                 3,
		SpreadMinPpm:           4_000,
		SpreadBufferPpm:        2_000,
		SkewFactorPpm:          999_999,
		OrderSizePctPpm:        200_000,
		OrderExpirationSeconds: 10,
	}
	err := k.SetParams(ctx, newParams)
	require.NoError(t, err)
	require.Equal(t, newParams, k.GetParams(ctx))

	// Set invalid params and get.
	invalidParams := types.Params{
		Layers:                 3,
		SpreadMinPpm:           4_000,
		SpreadBufferPpm:        2_000,
		SkewFactorPpm:          1_000_000,
		OrderSizePctPpm:        200_000,
		OrderExpirationSeconds: 0, // invalid
	}
	err = k.SetParams(ctx, invalidParams)
	require.Error(t, err)
	require.Equal(t, newParams, k.GetParams(ctx))
}
