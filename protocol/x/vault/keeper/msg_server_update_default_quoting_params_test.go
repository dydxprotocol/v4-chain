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

func TestMsgUpdateDefaultQuotingParams(t *testing.T) {
	tests := map[string]struct {
		// Msg.
		msg *types.MsgUpdateDefaultQuotingParams
		// Expected error
		expectedErr string
	}{
		"Success. Update to default": {
			msg: &types.MsgUpdateDefaultQuotingParams{
				Authority:            lib.GovModuleAddress.String(),
				DefaultQuotingParams: types.DefaultQuotingParams(),
			},
		},
		"Success. Update to non-default": {
			msg: &types.MsgUpdateDefaultQuotingParams{
				Authority: lib.GovModuleAddress.String(),
				DefaultQuotingParams: types.QuotingParams{
					Layers:                           3,
					SpreadMinPpm:                     234_567,
					SpreadBufferPpm:                  6_789,
					SkewFactorPpm:                    321_123,
					OrderSizePctPpm:                  255_678,
					OrderExpirationSeconds:           120,
					ActivationThresholdQuoteQuantums: dtypes.NewInt(2_121_343_787),
				},
			},
		},
		"Failure - Invalid Authority": {
			msg: &types.MsgUpdateDefaultQuotingParams{
				Authority:            constants.AliceAccAddress.String(),
				DefaultQuotingParams: types.DefaultQuotingParams(),
			},
			expectedErr: "invalid authority",
		},
		"Failure - Invalid Params": {
			msg: &types.MsgUpdateDefaultQuotingParams{
				Authority: lib.GovModuleAddress.String(),
				DefaultQuotingParams: types.QuotingParams{
					Layers:                           3,
					SpreadMinPpm:                     4_000,
					SpreadBufferPpm:                  2_000,
					SkewFactorPpm:                    500_000,
					OrderSizePctPpm:                  0, // invalid
					OrderExpirationSeconds:           5,
					ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
				},
			},
			expectedErr: types.ErrInvalidOrderSizePctPpm.Error(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper
			ms := keeper.NewMsgServerImpl(k)

			_, err := ms.UpdateDefaultQuotingParams(ctx, tc.msg)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				require.Equal(t, types.DefaultQuotingParams(), k.GetDefaultQuotingParams(ctx))
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.msg.DefaultQuotingParams, k.GetDefaultQuotingParams(ctx))
			}
		})
	}
}
