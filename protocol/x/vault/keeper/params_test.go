package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
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
		Layers:                           3,
		SpreadMinPpm:                     4_000,
		SpreadBufferPpm:                  2_000,
		SkewFactorPpm:                    999_999,
		OrderSizePctPpm:                  200_000,
		OrderExpirationSeconds:           10,
		ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
	}
	err := k.SetParams(ctx, newParams)
	require.NoError(t, err)
	require.Equal(t, newParams, k.GetParams(ctx))

	// Set invalid params and get.
	invalidParams := types.Params{
		Layers:                           3,
		SpreadMinPpm:                     4_000,
		SpreadBufferPpm:                  2_000,
		SkewFactorPpm:                    1_000_000,
		OrderSizePctPpm:                  200_000,
		OrderExpirationSeconds:           0, // invalid
		ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
	}
	err = k.SetParams(ctx, invalidParams)
	require.Error(t, err)
	require.Equal(t, newParams, k.GetParams(ctx))
}

func TestGetSetVaultParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// Get non-existent vault params.
	_, exists := k.GetVaultParams(ctx, constants.Vault_Clob_0)
	require.False(t, exists)

	// Set vault params of vault clob 0.
	vaultClob0Params := types.VaultParams{
		LaggedPrice: &pricestypes.MarketPrice{
			Id:       uint32(0),
			Exponent: -5,
			Price:    123_456_789,
		},
	}
	err := k.SetVaultParams(ctx, constants.Vault_Clob_0, vaultClob0Params)
	require.NoError(t, err)

	// Get vault params of vault clob 0.
	params, exists := k.GetVaultParams(ctx, constants.Vault_Clob_0)
	require.True(t, exists)
	require.Equal(t, vaultClob0Params, params)

	// Get vault params of vault clob 1.
	_, exists = k.GetVaultParams(ctx, constants.Vault_Clob_1)
	require.False(t, exists)
}
