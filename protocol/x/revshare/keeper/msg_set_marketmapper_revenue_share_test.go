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

func TestSetMarketMapperRevenueShareParams(t *testing.T) {
	tests := map[string]struct {
		// Msg
		msg *types.MsgSetMarketMapperRevenueShare
		// Expected error
		expectedErr string
	}{
		"Success - Set revenue share": {
			msg: &types.MsgSetMarketMapperRevenueShare{
				Authority: lib.GovModuleAddress.String(),
				Params: types.MarketMapperRevenueShareParams{
					Address:         constants.AliceAccAddress.String(),
					RevenueSharePpm: 100_000,
					ValidDays:       240,
				},
			},
			expectedErr: "",
		},
		"Failure - Invalid Authority": {
			msg: &types.MsgSetMarketMapperRevenueShare{
				Authority: constants.AliceAccAddress.String(),
				Params: types.MarketMapperRevenueShareParams{
					Address:         constants.AliceAccAddress.String(),
					RevenueSharePpm: 100_000,
					ValidDays:       240,
				},
			},
			expectedErr: "invalid authority",
		},
		"Failure - Empty Authority": {
			msg: &types.MsgSetMarketMapperRevenueShare{
				Params: types.MarketMapperRevenueShareParams{
					Address:         constants.AliceAccAddress.String(),
					RevenueSharePpm: 100_000,
					ValidDays:       240,
				},
			},
			expectedErr: "invalid authority",
		},
		"Failure - Invalid revenue share address": {
			msg: &types.MsgSetMarketMapperRevenueShare{
				Authority: lib.GovModuleAddress.String(),
				Params: types.MarketMapperRevenueShareParams{
					Address:         "invalid_address",
					RevenueSharePpm: 100_000,
					ValidDays:       240,
				},
			},
			expectedErr: "invalid address",
		},
		"Failure - Invalid revenue share ppm": {
			msg: &types.MsgSetMarketMapperRevenueShare{
				Authority: lib.GovModuleAddress.String(),
				Params: types.MarketMapperRevenueShareParams{
					Address:         constants.AliceAccAddress.String(),
					RevenueSharePpm: 1_000_000,
					ValidDays:       240,
				},
			},
			expectedErr: "rev share safety violation: rev shares greater than or equal to 100%",
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				tApp := testapp.NewTestAppBuilder(t).Build()
				ctx := tApp.InitChain()
				k := tApp.App.RevShareKeeper
				ms := keeper.NewMsgServerImpl(k)
				_, err := ms.SetMarketMapperRevenueShare(ctx, tc.msg)
				if tc.expectedErr != "" {
					require.Error(t, err)
					require.Contains(t, err.Error(), tc.expectedErr)
				} else {
					require.NoError(t, err)
					params := k.GetMarketMapperRevenueShareParams(ctx)
					require.Equal(t, tc.msg.Params, params)
				}
			},
		)
	}
}
