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

func TestMsgUpdateDefaultQuotingParams(t *testing.T) {
	tests := map[string]struct {
		// Msg.
		msg *vaulttypes.MsgUpdateDefaultQuotingParams
		// Operator.
		operator string
		// Expected error
		expectedErr string
	}{
		"Success. Update to default. Gov Authority": {
			msg: &vaulttypes.MsgUpdateDefaultQuotingParams{
				Authority:            lib.GovModuleAddress.String(),
				DefaultQuotingParams: vaulttypes.DefaultQuotingParams(),
			},
			operator: constants.AliceAccAddress.String(),
		},
		"Success. Update to default. Operator Authority": {
			msg: &vaulttypes.MsgUpdateDefaultQuotingParams{
				Authority:            constants.CarlAccAddress.String(),
				DefaultQuotingParams: vaulttypes.DefaultQuotingParams(),
			},
			operator: constants.CarlAccAddress.String(),
		},
		"Success. Update to non-default": {
			msg: &vaulttypes.MsgUpdateDefaultQuotingParams{
				Authority: lib.GovModuleAddress.String(),
				DefaultQuotingParams: vaulttypes.QuotingParams{
					Layers:                           3,
					SpreadMinPpm:                     234_567,
					SpreadBufferPpm:                  6_789,
					SkewFactorPpm:                    321_123,
					OrderSizePctPpm:                  255_678,
					OrderExpirationSeconds:           120,
					ActivationThresholdQuoteQuantums: dtypes.NewInt(2_121_343_787),
				},
			},
			operator: constants.AliceAccAddress.String(),
		},
		"Failure - Invalid Authority": {
			msg: &vaulttypes.MsgUpdateDefaultQuotingParams{
				Authority:            constants.AliceAccAddress.String(),
				DefaultQuotingParams: vaulttypes.DefaultQuotingParams(),
			},
			operator:    constants.BobAccAddress.String(),
			expectedErr: "invalid authority",
		},
		"Failure - Invalid Params": {
			msg: &vaulttypes.MsgUpdateDefaultQuotingParams{
				Authority: lib.GovModuleAddress.String(),
				DefaultQuotingParams: vaulttypes.QuotingParams{
					Layers:                           3,
					SpreadMinPpm:                     4_000,
					SpreadBufferPpm:                  2_000,
					SkewFactorPpm:                    500_000,
					OrderSizePctPpm:                  0, // invalid
					OrderExpirationSeconds:           5,
					ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
				},
			},
			operator:    constants.AliceAccAddress.String(),
			expectedErr: vaulttypes.ErrInvalidOrderSizePctPpm.Error(),
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

			_, err := ms.UpdateDefaultQuotingParams(ctx, tc.msg)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				require.Equal(t, vaulttypes.DefaultQuotingParams(), k.GetDefaultQuotingParams(ctx))
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.msg.DefaultQuotingParams, k.GetDefaultQuotingParams(ctx))
			}
		})
	}
}
