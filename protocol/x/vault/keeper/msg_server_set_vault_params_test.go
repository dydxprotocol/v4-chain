package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/keeper"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestMsgSetVaultParams(t *testing.T) {
	tests := map[string]struct {
		// Operator.
		operator string
		// Msg.
		msg *vaulttypes.MsgSetVaultParams
		// Expected error
		expectedErr string
	}{
		"Success - Gov Authority, Vault Clob 0": {
			operator: constants.AliceAccAddress.String(),
			msg: &vaulttypes.MsgSetVaultParams{
				Authority:   lib.GovModuleAddress.String(),
				VaultId:     constants.Vault_Clob0,
				VaultParams: constants.VaultParams,
			},
		},
		"Success - Gov Authority, Vault Clob 1": {
			operator: constants.AliceAccAddress.String(),
			msg: &vaulttypes.MsgSetVaultParams{
				Authority:   lib.GovModuleAddress.String(),
				VaultId:     constants.Vault_Clob1,
				VaultParams: constants.VaultParams,
			},
		},
		"Success - Operator Authority, Vault Clob 1": {
			operator: constants.AliceAccAddress.String(),
			msg: &vaulttypes.MsgSetVaultParams{
				Authority:   constants.AliceAccAddress.String(),
				VaultId:     constants.Vault_Clob1,
				VaultParams: constants.VaultParams,
			},
		},
		"Failure - Invalid Authority": {
			operator: constants.AliceAccAddress.String(),
			msg: &vaulttypes.MsgSetVaultParams{
				Authority:   constants.BobAccAddress.String(), // not a module authority or operator.
				VaultId:     constants.Vault_Clob0,
				VaultParams: constants.VaultParams,
			},
			expectedErr: vaulttypes.ErrInvalidAuthority.Error(),
		},
		"Failure - Empty Authority": {
			operator: constants.AliceAccAddress.String(),
			msg: &vaulttypes.MsgSetVaultParams{
				Authority:   "",
				VaultId:     constants.Vault_Clob0,
				VaultParams: constants.VaultParams,
			},
			expectedErr: vaulttypes.ErrInvalidAuthority.Error(),
		},
		"Failure - Vault Clob 0. Invalid Quoting Params": {
			operator: constants.AliceAccAddress.String(),
			msg: &vaulttypes.MsgSetVaultParams{
				Authority: lib.GovModuleAddress.String(),
				VaultId:   constants.Vault_Clob0,
				VaultParams: vaulttypes.VaultParams{
					Status: vaulttypes.VaultStatus_VAULT_STATUS_STAND_BY,
					QuotingParams: &vaulttypes.QuotingParams{
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
			expectedErr: vaulttypes.ErrInvalidActivationThresholdQuoteQuantums.Error(),
		},
		"Failure - Vault Clob 1. Unspecified status": {
			operator: constants.AliceAccAddress.String(),
			msg: &vaulttypes.MsgSetVaultParams{
				Authority: lib.GovModuleAddress.String(),
				VaultId:   constants.Vault_Clob0,
				VaultParams: vaulttypes.VaultParams{
					QuotingParams: &constants.QuotingParams,
				},
			},
			expectedErr: vaulttypes.ErrUnspecifiedVaultStatus.Error(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Set megavault operator.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vaulttypes.GenesisState) {
						genesisState.OperatorParams = vaulttypes.OperatorParams{
							Operator: tc.operator,
						}
					},
				)
				return genesis
			}).Build()
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
