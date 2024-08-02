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

func TestMsgSetVaultQuotingParams(t *testing.T) {
	tests := map[string]struct {
		// Msg.
		msg *types.MsgSetVaultQuotingParams
		// Expected error
		expectedErr string
	}{
		"Success - Vault Clob 0": {
			msg: &types.MsgSetVaultQuotingParams{
				Authority:     lib.GovModuleAddress.String(),
				VaultId:       constants.Vault_Clob0,
				QuotingParams: constants.QuotingParams,
			},
		},
		"Success - Vault Clob 1": {
			msg: &types.MsgSetVaultQuotingParams{
				Authority:     lib.GovModuleAddress.String(),
				VaultId:       constants.Vault_Clob1,
				QuotingParams: constants.QuotingParams,
			},
		},
		"Failure - Invalid Authority": {
			msg: &types.MsgSetVaultQuotingParams{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob0,
				QuotingParams: constants.QuotingParams,
			},
			expectedErr: "invalid authority",
		},
		"Failure - Vault Clob 0. Invalid Quoting Params": {
			msg: &types.MsgSetVaultQuotingParams{
				Authority: lib.GovModuleAddress.String(),
				VaultId:   constants.Vault_Clob0,
				QuotingParams: types.QuotingParams{
					Layers:                           3,
					SpreadMinPpm:                     4_000,
					SpreadBufferPpm:                  2_000,
					SkewFactorPpm:                    500_000,
					OrderSizePctPpm:                  100_000,
					OrderExpirationSeconds:           5,
					ActivationThresholdQuoteQuantums: dtypes.NewInt(-1), // invalid
				},
			},
			expectedErr: types.ErrInvalidActivationThresholdQuoteQuantums.Error(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper
			ms := keeper.NewMsgServerImpl(k)

			_, err := ms.SetVaultQuotingParams(ctx, tc.msg)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				require.Equal(
					t,
					types.DefaultQuotingParams(),
					k.GetVaultQuotingParams(ctx, tc.msg.VaultId),
				)
			} else {
				require.NoError(t, err)
				require.Equal(
					t,
					tc.msg.QuotingParams,
					k.GetVaultQuotingParams(ctx, tc.msg.VaultId),
				)
			}
		})
	}
}
