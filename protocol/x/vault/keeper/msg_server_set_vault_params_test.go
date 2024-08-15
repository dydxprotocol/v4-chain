package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	"github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestMsgSetVaultParams(t *testing.T) {
	tests := map[string]struct {
		// Msg.
		msg *types.MsgSetVaultParams
		// Expected error
		expectedErr string
	}{
		"Success - Vault Clob 0": {
			msg: &types.MsgSetVaultParams{
				Authority:   lib.GovModuleAddress.String(),
				VaultId:     constants.Vault_Clob0,
				VaultParams: constants.VaultParams,
			},
		},
		"Success - Vault Clob 1": {
			msg: &types.MsgSetVaultParams{
				Authority:   lib.GovModuleAddress.String(),
				VaultId:     constants.Vault_Clob1,
				VaultParams: constants.VaultParams,
			},
		},
		"Failure - Invalid Authority": {
			msg: &types.MsgSetVaultParams{
				Authority:   constants.AliceAccAddress.String(),
				VaultId:     constants.Vault_Clob0,
				VaultParams: constants.VaultParams,
			},
			expectedErr: "invalid authority",
		},
		"Failure - Empty Authority": {
			msg: &types.MsgSetVaultParams{
				Authority:   "",
				VaultId:     constants.Vault_Clob0,
				VaultParams: constants.VaultParams,
			},
			expectedErr: "invalid authority",
		},
		"Failure - Vault Clob 0. Invalid Quoting Params": {
			msg: &types.MsgSetVaultParams{
				Authority: lib.GovModuleAddress.String(),
				VaultId:   constants.Vault_Clob0,
				VaultParams: types.VaultParams{
					Status: types.VaultStatus_VAULT_STATUS_STAND_BY,
					QuotingParams: &types.QuotingParams{
						Layers:                           3,
						SpreadMinPpm:                     4_000,
						SpreadBufferPpm:                  2_000,
						SkewFactorPpm:                    500_000,
						OrderSizePctPpm:                  100_000,
						OrderExpirationSeconds:           5,
						ActivationThresholdQuoteQuantums: dtypes.NewInt(-1), // invalid
					},
				},
			},
			expectedErr: types.ErrInvalidActivationThresholdQuoteQuantums.Error(),
		},
		"Failure - Vault Clob 1. Unspecified status": {
			msg: &types.MsgSetVaultParams{
				Authority: lib.GovModuleAddress.String(),
				VaultId:   constants.Vault_Clob0,
				VaultParams: types.VaultParams{
					QuotingParams: &constants.QuotingParams,
				},
			},
			expectedErr: types.ErrUnspecifiedVaultStatus.Error(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper
			ms := keeper.NewMsgServerImpl(k)

			_, err := ms.SetVaultParams(ctx, tc.msg)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				_, exists := k.GetVaultParams(ctx, tc.msg.VaultId)
				require.False(t, exists)
			} else {
				require.NoError(t, err)
				p, exists := k.GetVaultParams(ctx, tc.msg.VaultId)
				require.True(t, exists)
				require.Equal(
					t,
					tc.msg.VaultParams,
					p,
				)
			}
		})
	}
}
