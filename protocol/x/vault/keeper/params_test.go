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
	newParams := types.QuotingParams{
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
	invalidParams := types.QuotingParams{
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

func TestGetSetVaultParams(t *testing.T) {
	tests := map[string]struct {
		// Vault id.
		vaultId types.VaultId
		// Vault params to set.
		vaultParams *types.VaultParams
		// Expected error.
		expectedErr error
	}{
		"Success - Vault Clob 0": {
			vaultId:     constants.Vault_Clob0,
			vaultParams: &constants.VaultParams,
		},
		"Success - Vault Clob 1": {
			vaultId:     constants.Vault_Clob1,
			vaultParams: &constants.VaultParams,
		},
		"Success - Non-existent Vault Params": {
			vaultId:     constants.Vault_Clob1,
			vaultParams: nil,
		},
		"Failure - Unspecified Status": {
			vaultId: constants.Vault_Clob0,
			vaultParams: &types.VaultParams{
				QuotingParams: &constants.QuotingParams,
			},
			expectedErr: types.ErrUnspecifiedVaultStatus,
		},
		"Failure - Invalid Quoting Params": {
			vaultId: constants.Vault_Clob0,
			vaultParams: &types.VaultParams{
				Status:        types.VaultStatus_VAULT_STATUS_STAND_BY,
				QuotingParams: &constants.InvalidQuotingParams,
			},
			expectedErr: types.ErrInvalidOrderExpirationSeconds,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			if tc.vaultParams == nil {
				_, exists := k.GetVaultParams(ctx, tc.vaultId)
				require.False(t, exists)
				return
			}

			err := k.SetVaultParams(ctx, tc.vaultId, *tc.vaultParams)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				_, exists := k.GetVaultParams(ctx, tc.vaultId)
				require.False(t, exists)
			} else {
				require.NoError(t, err)
				p, exists := k.GetVaultParams(ctx, tc.vaultId)
				require.True(t, exists)
				require.Equal(t, *tc.vaultParams, p)
			}
		})
	}
}

func TestGetVaultQuotingParams(t *testing.T) {
	tests := map[string]struct {
		/* Setup */
		// Vault id.
		vaultId types.VaultId
		// Vault params to set.
		vaultParams *types.VaultParams
		/* Expectations */
		// Whether quoting params should be default.
		shouldBeDefault bool
	}{
		"Default Quoting Params": {
			vaultId: constants.Vault_Clob0,
			vaultParams: &types.VaultParams{
				Status: types.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
			},
			shouldBeDefault: true,
		},
		"Custom Quoting Params": {
			vaultId:         constants.Vault_Clob1,
			vaultParams:     &constants.VaultParams,
			shouldBeDefault: false,
		},
		"Non-existent Vault Params": {
			vaultId:     constants.Vault_Clob1,
			vaultParams: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			if tc.vaultParams != nil {
				err := k.SetVaultParams(ctx, tc.vaultId, *tc.vaultParams)
				require.NoError(t, err)
				p, exists := k.GetVaultQuotingParams(ctx, tc.vaultId)
				require.True(t, exists)
				if tc.shouldBeDefault {
					require.Equal(t, types.DefaultQuotingParams(), p)
				} else {
					require.Equal(t, *tc.vaultParams.QuotingParams, p)
				}
			} else {
				_, exists := k.GetVaultQuotingParams(ctx, tc.vaultId)
				require.False(t, exists)
			}
		})
	}
}
