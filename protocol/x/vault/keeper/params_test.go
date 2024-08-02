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
	require.Equal(t, *newParams, k.GetDefaultQuotingParams(ctx))

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
	require.Equal(t, *newParams, k.GetDefaultQuotingParams(ctx))
}

func TestGetSetVaultQuotingParams(t *testing.T) {
	tests := map[string]struct {
		// Vault id.
		vaultId types.VaultId
		// Vault quoting params to set.
		vaultQuotingParams *types.QuotingParams
	}{
		"Vault Clob 0. Default quoting params": {
			vaultId:            constants.Vault_Clob0,
			vaultQuotingParams: nil,
		},
		"Vault Clob 0. Non-default quoting params": {
			vaultId: constants.Vault_Clob0,
			vaultQuotingParams: &types.QuotingParams{
				Layers:                           3,
				SpreadMinPpm:                     12_345,
				SpreadBufferPpm:                  5_678,
				SkewFactorPpm:                    4_121_787,
				OrderSizePctPpm:                  232_121,
				OrderExpirationSeconds:           120,
				ActivationThresholdQuoteQuantums: dtypes.NewInt(2_123_456_789),
			},
		},
		"Vault Clob 1. Default quoting params": {
			vaultId:            constants.Vault_Clob0,
			vaultQuotingParams: nil,
		},
		"Vault Clob 1. Non-default quoting params": {
			vaultId: constants.Vault_Clob0,
			vaultQuotingParams: &types.QuotingParams{
				Layers:                           4,
				SpreadMinPpm:                     123_456,
				SpreadBufferPpm:                  87_654,
				SkewFactorPpm:                    5_432_111,
				OrderSizePctPpm:                  444_333,
				OrderExpirationSeconds:           90,
				ActivationThresholdQuoteQuantums: dtypes.NewInt(1_111_111_111),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			if tc.vaultQuotingParams != nil {
				// Set quoting params.
				err := k.SetVaultQuotingParams(ctx, tc.vaultId, *tc.vaultQuotingParams)
				require.NoError(t, err)
				// Verify quoting params are as set.
				require.Equal(
					t,
					*tc.vaultQuotingParams,
					k.GetVaultQuotingParams(ctx, tc.vaultId),
				)
			} else {
				// Verify quoting params are default.
				require.Equal(
					t,
					k.GetDefaultQuotingParams(ctx),
					k.GetVaultQuotingParams(ctx, tc.vaultId),
				)
			}
		})
	}
}
