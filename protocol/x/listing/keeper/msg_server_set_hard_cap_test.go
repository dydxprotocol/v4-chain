package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	"github.com/stretchr/testify/require"
)

func TestMsgSetMarketsHardCap(t *testing.T) {
	tests := map[string]struct {
		// Msg.
		msg *types.MsgSetMarketsHardCap
		// Expected error
		expectedErr string
	}{
		"Success - Hard cap to 100": {
			msg: &types.MsgSetMarketsHardCap{
				Authority:         lib.GovModuleAddress.String(),
				HardCapForMarkets: 100,
			},
		},
		"Success - Disabled": {
			msg: &types.MsgSetMarketsHardCap{
				Authority:         lib.GovModuleAddress.String(),
				HardCapForMarkets: 0,
			},
		},
		"Failure - Invalid Authority": {
			msg: &types.MsgSetMarketsHardCap{
				Authority:         constants.AliceAccAddress.String(),
				HardCapForMarkets: 100,
			},
			expectedErr: "invalid authority",
		},
		"Failure - Empty authority": {
			msg: &types.MsgSetMarketsHardCap{
				HardCapForMarkets: 100,
			},
			expectedErr: "invalid authority",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.ListingKeeper
			ms := keeper.NewMsgServerImpl(k)
			_, err := ms.SetMarketsHardCap(ctx, tc.msg)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				enabledFlag := k.GetMarketsHardCap(ctx)
				require.Equal(t, tc.msg.HardCapForMarkets, enabledFlag)
			}
		})
	}
}
