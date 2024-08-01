package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestGetSetDefaultQuotingParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// Params should have default values at genesis.
	params := k.GetDefaultQuotingParams(ctx)
	require.Equal(t, types.DefaultQuotingParams(), params)

	// Set new params and get.
	newParams := &types.QuotingParams{
		Layers:                           3,
		SpreadMinPpm:                     4_000,
		SpreadBufferPpm:                  2_000,
		SkewFactorPpm:                    999_999,
		OrderSizePctPpm:                  200_000,
		OrderExpirationSeconds:           10,
		ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
	}
	err := k.SetDefaultQuotingParams(ctx, newParams)
	require.NoError(t, err)
	require.Equal(t, newParams, k.GetDefaultQuotingParams(ctx))

	// Set invalid params and get.
	invalidParams := &types.QuotingParams{
		Layers:                           3,
		SpreadMinPpm:                     4_000,
		SpreadBufferPpm:                  2_000,
		SkewFactorPpm:                    1_000_000,
		OrderSizePctPpm:                  200_000,
		OrderExpirationSeconds:           0, // invalid
		ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
	}
	err = k.SetDefaultQuotingParams(ctx, invalidParams)
	require.Error(t, err)
	require.Equal(t, newParams, k.GetDefaultQuotingParams(ctx))
}

func TestGetSetVaultQuotingParams(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// Should get default quoting params.
	quotingParams := k.GetVaultQuotingParams(ctx, constants.Vault_Clob0)
	require.Equal(t, types.DefaultQuotingParams(), quotingParams)

	// Set quoting params of vault clob 0.
	vaultClob0QuotingParams := &types.QuotingParams{
		Layers:                           3,
		SpreadMinPpm:                     12_345,
		SpreadBufferPpm:                  5_678,
		SkewFactorPpm:                    4_121_787,
		OrderSizePctPpm:                  232_121,
		OrderExpirationSeconds:           120,
		ActivationThresholdQuoteQuantums: dtypes.NewInt(2_123_456_789),
	}
	err := k.SetVaultQuotingParams(ctx, constants.Vault_Clob0, vaultClob0QuotingParams)
	require.NoError(t, err)

	// Get quoting params of vault clob 0.
	quotingParams = k.GetVaultQuotingParams(ctx, constants.Vault_Clob0)
	require.Equal(t, vaultClob0QuotingParams, quotingParams)

	// Set quoting params of vault clob 1.
	vaultClob1QuotingParams := &types.QuotingParams{
		Layers:                           4,
		SpreadMinPpm:                     123_456,
		SpreadBufferPpm:                  87_654,
		SkewFactorPpm:                    5_432_111,
		OrderSizePctPpm:                  444_333,
		OrderExpirationSeconds:           90,
		ActivationThresholdQuoteQuantums: dtypes.NewInt(1_111_111_111),
	}
	err = k.SetVaultQuotingParams(ctx, constants.Vault_Clob1, vaultClob1QuotingParams)
	require.NoError(t, err)

	// Get quoting params of vault clob 1.
	quotingParams = k.GetVaultQuotingParams(ctx, constants.Vault_Clob1)
	require.Equal(t, vaultClob1QuotingParams, quotingParams)
}
