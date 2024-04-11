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

func TestMsgUpdateParams(t *testing.T) {
	tests := map[string]struct {
		// Msg.
		msg *types.MsgUpdateParams
		// Expected error
		expectedErr string
	}{
		"Success": {
			msg: &types.MsgUpdateParams{
				Authority: lib.GovModuleAddress.String(),
				Params:    types.DefaultParams(),
			},
		},
		"Failure - Invalid Authority": {
			msg: &types.MsgUpdateParams{
				Authority: constants.AliceAccAddress.String(),
				Params:    types.DefaultParams(),
			},
			expectedErr: "invalid authority",
		},
		"Failure - Invalid Params": {
			msg: &types.MsgUpdateParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.Params{
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

			_, err := ms.UpdateParams(ctx, tc.msg)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				require.Equal(t, types.DefaultParams(), k.GetParams(ctx))
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.msg.Params, k.GetParams(ctx))
			}
		})
	}
}
