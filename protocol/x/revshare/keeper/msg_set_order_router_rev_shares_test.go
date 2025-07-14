package keeper_test

import (
	"testing"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/keeper"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

func TestSetOrderRouterRevShares(t *testing.T) {
	tests := map[string]struct {
		// Msg
		msg *types.MsgSetOrderRouterRevShare
		// Expected error
		expectedErr string
	}{
		"Success - Set revenue share": {
			msg: &types.MsgSetOrderRouterRevShare{
				Authority: lib.GovModuleAddress.String(),
				OrderRouterRevShare: types.OrderRouterRevShare{
					Address:  constants.AliceAccAddress.String(),
					SharePpm: 100_000,
				},
			},
			expectedErr: "",
		},
		"Failure - Invalid Authority": {
			msg: &types.MsgSetOrderRouterRevShare{
				Authority: constants.AliceAccAddress.String(),
				OrderRouterRevShare: types.OrderRouterRevShare{
					Address:  constants.AliceAccAddress.String(),
					SharePpm: 100_000,
				},
			},
			expectedErr: "invalid authority",
		},
		"Failure - Empty Authority": {
			msg: &types.MsgSetOrderRouterRevShare{
				OrderRouterRevShare: types.OrderRouterRevShare{
					Address:  constants.AliceAccAddress.String(),
					SharePpm: 100_000,
				},
			},
			expectedErr: "invalid authority",
		},
		"Failure - Invalid revenue share address": {
			msg: &types.MsgSetOrderRouterRevShare{
				Authority: lib.GovModuleAddress.String(),
				OrderRouterRevShare: types.OrderRouterRevShare{
					Address:  "invalid_address",
					SharePpm: 100_000,
				},
			},
			expectedErr: "invalid address",
		},
		"Failure - Invalid revenue share ppm": {
			msg: &types.MsgSetOrderRouterRevShare{
				Authority: lib.GovModuleAddress.String(),
				OrderRouterRevShare: types.OrderRouterRevShare{
					Address:  constants.AliceAccAddress.String(),
					SharePpm: 1_000_000,
				},
			},
			expectedErr: "rev share safety violation: rev shares greater than or equal to allowed amount:",
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				tApp := testapp.NewTestAppBuilder(t).Build()
				ctx := tApp.InitChain()
				k := tApp.App.RevShareKeeper
				ms := keeper.NewMsgServerImpl(k)
				_, err := ms.SetOrderRouterRevShares(ctx, tc.msg)
				if tc.expectedErr != "" {
					require.Error(t, err)
					require.Contains(t, err.Error(), tc.expectedErr)
				} else {
					require.NoError(t, err)
					for _, revShare := range tc.msg.OrderRouterRevShares {
						params, err := k.GetOrderRouterRevShares(ctx, revShare.Address)
						require.NoError(t, err)
						require.Equal(t, revShare.SharePpm, params)
					}
				}
			},
		)
	}
}
